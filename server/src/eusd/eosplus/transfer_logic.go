package eosplus

import (
	"common"
	"fmt"
	"github.com/eoscanada/eos-go"
	"otc_error"
	"strconv"
	"strings"
	"time"
	"utils/eusd/dao"
	"utils/eusd/models"
	dao2 "utils/otc/dao"
)

func TransferLogic(logIds []int64, from, to *models.EosWealth, quantity int64) {
	transferLogicBase(logIds, from, to, quantity, 0, false, "")
}

func TransferLogicMemo(logIds []int64, from, to *models.EosWealth, quantity int64, memo string) {
	transferLogicBase(logIds, from, to, quantity, 0, false, memo)
}

// otc订单转账
func TransferLogicOrderByUids(orderId uint64, logIds []int64, from, to uint64, quantity int64) {
	fromUser, toUser, err := dao.WealthDaoEntity.GetFromToUser(from, to)
	if err != nil {
		return
	}
	TransferLogicOrder(orderId, logIds, fromUser, toUser, quantity)
}

// otc订单转账
func TransferLogicOrder(orderId uint64, logIds []int64, from, to *models.EosWealth, quantity int64) {
	transferLogicBase(logIds, from, to, quantity, orderId, false, "")
}

// 转账信息入库基础逻辑
func transferLogicBase(logIds []int64, from, to *models.EosWealth, quantity int64, orderId uint64, delayDeal bool, memo string) {
	str := ""
	for _, v := range logIds {
		str += fmt.Sprint(v) + ","
	}
	if to.Account == "" {
		acc, errCode := AccountRpc().Bind(to.Uid)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			common.LogFuncError("TransferLogic Bind Account Err: %v", logIds)
			return
		}
		to.Account = acc.Account
		_ = dao.WealthDaoEntity.Update(to, models.COLUMN_EosWealth_Account)
	}
	data := &models.EosTxLog{
		OrderId:   orderId,
		From:      from.Account,
		FromUid:   from.Uid,
		To:        to.Account,
		ToUid:     to.Uid,
		Ctime:     common.NowInt64MS(),
		Quantity:  quantity,
		LogIds:    str,
		Status:    0,
		DelayDeal: delayDeal,
		Memo:      memo,
	}
	sign, err := transferSign(data)
	if err != nil {
		common.LogFuncError("TxLogErr Md5: %v", ToJson(data))
		return
	}
	data.Sign = sign

	txLogId, err := dao.EosTxLogDaoEntity.Create(data)
	if err != nil {
		common.LogFuncError("TxLogErr DB: %v", ToJson(data))
		return
	}
	data.Id = uint64(txLogId)

	err = common.RabbitMQPublish(RabbitMQEusdTransfer, RabbitMQEusdTransfer, []byte(fmt.Sprintf("%d", txLogId)))
	if err != nil {
		common.LogFuncError("TxLogErr MQ: %v", ToJson(data))
		return
	}
}

const EUDTSalt = "cFWBKOnoifn34$cef5b"

func transferSign(TxLog *models.EosTxLog) (md5Str string, err error) {
	md5Str, err = common.GenerateDoubleMD5(fmt.Sprintf("%s%s%d%d", TxLog.From, TxLog.To, TxLog.Ctime, TxLog.Quantity), EUDTSalt)
	return
}

func transferSignCheck(TxLog *models.EosTxLog) bool {
	return common.CompareMd5(TxLog.Sign, fmt.Sprintf("%s%s%d%d", TxLog.From, TxLog.To, TxLog.Ctime, TxLog.Quantity), EUDTSalt)
}

// cron 循环发起交易
func CronRunTxLog() (err error) {
	for {
		common.LogFuncDebug("CronTxLog")
		list, err := dao.EosTxLogDaoEntity.FetchCheck(100)
		if err != nil {
			common.LogFuncDebug("CronTxLog Err:%v", err)
			continue
		}
		if len(list) <= 0 {
			break
		}

		for _, v := range list {
			RealTransfer(v)
		}
		time.Sleep(1 * time.Second)
	}
	return
}

//使用MQ队列运行交易
func MQRunTxLog(id uint64) {
	txLog := dao.EosTxLogDaoEntity.Info(id)
	if txLog.Id > 0 {
		RealTransfer(txLog)
	}
}

//发起转账
func RealTransfer(TxLog *models.EosTxLog) {
	//更新记录状态
	if !dao.EosTxLogDaoEntity.UpdateTransferring(TxLog.Id) {
		// 状态处于发起请求中
		return
	}
	if !transferSignCheck(TxLog) {
		// 交易sign检查错误
		common.LogFuncError("TransferLogic Sign Err:%v", TxLog.Id)
		_ = dao.EosTxLogDaoEntity.UpdateStatusSignErr(TxLog.Id)
		return
	}

	errCode := controllers.ERROR_CODE_SUCCESS
	transaction := &models.EosTransaction{}
	rpc := EosPlusAPI.Rpc

	quant := QuantityInt64ToString(TxLog.Quantity) + " " + EosConfig.Symbol

	if TxLog.From == EosConfig.TokenAccount && TxLog.To != "retire" && TxLog.ToUid != 0 {
		//发币
		_, transaction, errCode = rpc.IssueToken(TxLog.To, quant)

	} else if TxLog.To == "retire" && TxLog.ToUid == 0 {
		//销毁
		_, transaction, errCode = rpc.RetireToken(quant, "")
	} else {
		//转账
		_, transaction, errCode = rpc.Transfer(EosConfig.TokenAccount, TxLog.From, TxLog.To, quant, TxLog.Memo)
	}

	//资源检查
	_ = common.RabbitMQPublish(RabbitMQEusdCheckResource, RabbitMQEusdCheckResource, []byte(TxLog.From))

	if errCode != controllers.ERROR_CODE_SUCCESS {
		if TxLog.Retry < 3 {
			//转账失败，重新尝试
			dao.EosTxLogDaoEntity.ResetToCreated(TxLog.Id)
		} else {
			_ = dao.EosTxLogDaoEntity.UpdateStatusErr(TxLog.Id)
		}

		common.LogFuncError("TransferLogic Err: f:%v t:%v q:%v", TxLog.From, TxLog.To, quant)
		return
	}

	logIds := []int64{}
	tmp := strings.Split(TxLog.LogIds, ",")
	for _, v := range tmp {
		i, _ := strconv.ParseInt(v, 10, 64)
		logIds = append(logIds, i)
	}
	//更新 log
	dao.WealthLogDaoEntity.UpdateTxid(transaction.Id, logIds...)
	_ = dao.EosTxLogDaoEntity.UpdateTransferred(TxLog.Id, transaction.Id)
	//写MQ进行交易检查
	_ = common.RabbitMQPublishDelay(RabbitMQEusdTransferCheck, RabbitMQEusdTransferCheck, []byte(fmt.Sprintf("%d", TxLog.Id)), "180000")
}

func checkTransferRedisLock() bool {
	resp := common.RedisManger.SetNX("CheckTransferLock", 1, 59*time.Second)
	lockSuccess, err := resp.Result()
	if err != nil || !lockSuccess {
		return false
	}
	return true
}

//循环检查订单是否交易完成
func CronCheckTransfer() (err error) {
	if !checkTransferRedisLock() {
		return
	}
	for {
		common.LogFuncDebug("CheckTransfer")
		list, err := dao.EosTxLogDaoEntity.FetchCheckTransfer(100)
		if err != nil {
			common.LogFuncDebug("CheckTransfer Err: %v", err)
			continue
		}
		if len(list) <= 0 {
			break
		}

		for _, v := range list {
			if (common.NowInt64MS()-v.Ctime)/1000 < 180 {
				time.Sleep(10 * time.Second)
				break
			}
			TransferCheckLogic(EosPlusAPI.Rpc, v)
		}
		time.Sleep(1 * time.Second)
	}
	return
}

// MQ队列
func MQCheckTransfer(txLogId uint64) {
	data := dao.EosTxLogDaoEntity.Info(txLogId)
	if data.Id > 0 {
		TransferCheckLogic(EosPlusAPI.Rpc, data)
	}
}

type Err struct {
	Code  int `json:"Code"`
	Error struct {
		Code int    `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	} `json:"error"`
}

// 检查转账
// 大部分节点开启的交易查询，都要求指定区块block_num_hint。
// block_num_hint 不一定是交易返回的那个block_num，这会导致交易查询不到。
// 解决方案遍历 block_num 之后20个区块进行验证， 发起交易后3分钟进行验证，验证3次（第2次6分钟后,第3次9分钟后），3次检查都失败重新进行交易
func TransferCheckLogic(rpc *RpcKeyBag, TxLog *models.EosTxLog) {
	//取交易详情
	transaction := dao.TransactionDaoEntity.Info(TxLog.Txid)

	if transaction.Id == 0 { // 没找到交易详情
		common.LogFuncError("Transfer Check notfound transferInfo:%v", TxLog.Id)
		return
	}

	//获取区块
	blockInfo := &eos.TransactionRespV16{}
	for i := 0; i < 20; i++ {
		ok, block, _ := rpc.IsTransferOK(transaction.TransactionId, int(transaction.BlockNum)+i)
		if block != nil && block.Id == transaction.TransactionId {
			if !ok { // 未到不可逆区块
				break
			}
			transaction.BlockNum += uint32(i)
			blockInfo = block
			break
		}
	}

	ok := checkTransferBlockInfo(blockInfo, transaction, TxLog)
	if !ok {
		if common.NowInt64MS()-TxLog.Ctime > 54000 {
			// 超9分钟查询交易还是不存在 重新发起交易
			dao.EosTxLogDaoEntity.ResetToCreated(TxLog.Id)
		} else {
			// 重新进行交易检查
			err := common.RabbitMQPublishDelay(RabbitMQEusdTransferCheck, RabbitMQEusdTransferCheck, []byte(fmt.Sprintf("%d", TxLog.Id)), "180000")
			if err != nil {
				common.LogFuncError("Transfer Check to MQ:%v,%v", TxLog.Id, err)
			}
		}
		return
	}
	//交易成功
	transferFinish(transaction, TxLog)
}

//检查交易区块信息
func checkTransferBlockInfo(block *eos.TransactionRespV16, transaction *models.EosTransaction, TxLog *models.EosTxLog) (res bool) {
	res = false
	//未取到区块
	if block.Id == "" {
		return
	}

	if len(block.Trx.Trx.Actions) < 1 {
		return
	}
	action := block.Trx.Trx.Actions[0]
	if transaction.Type == TransactionTypeIssueToken && action.Name == "issue" {
		if action.Data.From != "" || action.Data.To != transaction.Receiver || transaction.Quantity != action.Data.Quantity {
			common.LogFuncDebug("transfer Check:%v, action err1", TxLog.Id)
			return
		}
	} else if transaction.Type == TransactionTypeTransfer && action.Name == "transfer" {
		if action.Data.From != transaction.Payer || action.Data.To != transaction.Receiver || transaction.Quantity != action.Data.Quantity {
			common.LogFuncDebug("transfer Check:%v, action err2", TxLog.Id)
			return
		}

	} else {
		common.LogFuncDebug("transfer Check:%v, action unknow", TxLog.Id)
		return
	}

	return true
}

//交易成功 - 后续逻辑
func transferFinish(transaction *models.EosTransaction, TxLog *models.EosTxLog) {
	//记录标记
	dao.TransactionDaoEntity.Finish(transaction.Id, transaction.BlockNum)
	_ = dao.EosTxLogDaoEntity.UpdateStatusFinish(TxLog.Id)

	//更新资产变更记录
	logIds := []int64{}
	tmp := strings.Split(TxLog.LogIds, ",")
	for _, v := range tmp {
		i, _ := strconv.ParseInt(v, 10, 64)
		logIds = append(logIds, i)
	}
	dao.WealthLogDaoEntity.UpdateStatus(logIds...)

	//解冻资金
	if TxLog.ToUid <= 0 {
		// 没有接受人
		return
	}
	if TxLog.OrderId <= 0 {
		//非交易
		WealthTransferToAvailable(TxLog.ToUid, TxLog.Quantity)
		return
	}
	//otc转账
	order, err := dao2.OrdersDaoEntity.Info(TxLog.OrderId)
	if err != nil {
		common.LogFuncError("transfer Check: Order Not Found;txid:%v", TxLog.Id)
		return
	}

	if order.Amount != TxLog.Quantity {
		common.LogFuncError("transfer Check: Quantity No Match;txid:%v", TxLog.Id)
		return
	}

	if order.Side == dao2.SideBuy {
		//用户购买token
		if order.Uid != TxLog.ToUid {
			common.LogFuncError("transfer Check: Uid No Match;txid:%v", TxLog.Id)
			return
		}
		WealthTransferToAvailable(TxLog.ToUid, TxLog.Quantity)
		return
	}
	//承兑商收购token
	if order.EUid != TxLog.ToUid {
		common.LogFuncError("transfer Check: EUid No Match;txid:%v", TxLog.Id)
		return
	}

	ok := dao.EosOtcDaoEntity.OtcTransferUnlock(order.EUid, order.Amount)
	if !ok {
		common.LogFuncError("transfer Check: unlock err;txid:%v", TxLog.Id)
		return
	}
}

// 转账中状态还原
func WealthTransferToAvailable(uid uint64, quantity int64) {
	l, err := wealthRedisJamLock(uid)
	if err != nil {
		common.LogFuncError("WealthTransferToAvailable: %v, u:%v, q:%v", err, uid, quantity)
		return
	}
	defer func() {
		_ = wealthRedisUnJamLock(l)
	}()
	//重新获取数据
	user, _ := dao.WealthDaoEntity.Info(uid)

	transfer := quantity
	transferGame := int64(0)

	if user.Transfer < quantity {
		transfer = user.Transfer
		transferGame = quantity - user.Transfer
		if user.TransferGame < transferGame {
			// 解锁时，被锁定币 少于 锁定币，可能游戏输了
			transferGame = user.Game
			common.LogFuncWarning("WealthTransferToAvailable: q-%v t-%v tg-%v", quantity, user.Transfer, user.TransferGame)
		}
	}
	dao.WealthDaoEntity.TransferToAvailable(user.Uid, transfer, transferGame)
}
