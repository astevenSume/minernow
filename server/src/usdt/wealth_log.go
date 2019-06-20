package usdt

import (
	"common"
	"fmt"
	. "otc_error"
	"time"
	"usdt/explorer"
	"usdt/prices"
	"utils/usdt/dao"
	"utils/usdt/models"

	"github.com/astaxie/beego/orm"
	"github.com/mitchellh/mapstructure"
)

type WealthLogRecord struct {
	Id       uint64 `json:"id"`
	Uid      uint64 `json:"uid"`
	Type     uint32 `json:"type"`
	Status   uint32 `json:"status"`
	Txid     string `json:"txid"`
	Amount   string `json:"amount"`
	Utime    int64  `json:"utime"`
	Ctime    int64  `json:"ctime"`
	BlockNum string `json:"block_num"`
	Fee      string `json:"fee"`
	Memo     string `json:"memo"`
	From     string `json:"from"`
	To       string `json:"to"`
}

type WealthLogRecordClient struct {
	Id         string `json:"id"`
	Uid        string `json:"uid"`
	Type       uint32 `json:"type"`
	Status     uint32 `json:"status"`
	Step       string `json:"step"`
	Desc       string `json:"desc"`
	Txid       string `json:"txid"`
	Amount     string `json:"amount"`
	Utime      int64  `json:"utime"`
	Ctime      int64  `json:"ctime"`
	Fee        string `json:"fee"`
	FeeOnChain string `json:"fee_onchain"`
	Memo       string `json:"memo"`
	From       string `json:"from"`
	To         string `json:"to"`
	Address    string `json:"address"`
	Mobile     string `json:"mobile"`
}

// 生成提现现订单
func Transfer(uid uint64, to string, amount, fee string, memo string) (logId uint64, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	// check transfer out info
	if to == "" {
		errCode = ERROR_CODE_PARAMS_ERROR
		return
	}

	// check available amount
	account, err := dao.AccountDaoEntity.QueryByUid(uid)
	if err != nil {
		errCode = ERROR_CODE_USDT_ACCOUNT_NO_FOUND
		return
	}

	// 账号锁定
	if account.Status == dao.STATUS_LOCKED {
		errCode = ERROR_CODE_USDT_ACCOUNT_LOCK
		return
	}

	// don't try to transfer to self
	if account.Address == to {
		errCode = ERROR_CODE_USDT_TRANSFER_TO_SELF_ERR
		return
	}

	var (
		transferInteger, feeInteger, feeUsdtInteger int64
		btc2usdt                                    float64
	)
	transferInteger, err = common.CurrencyStrToInt64(amount)
	if err != nil {
		errCode = ERROR_CODE_USDT_CURRENCY_PARAM_ERROR
		return
	}

	feeUsdtInteger, err = common.CurrencyStrToInt64(fee)
	if err != nil {
		errCode = ERROR_CODE_USDT_FEE_CURRENCY_PARAM_ERROR
		return
	}

	if account.AvailableInteger < transferInteger || account.AvailableInteger < feeUsdtInteger {
		errCode = ERROR_CODE_USDT_NO_ENOUGH
		return
	}

	if btc2usdt = prices.GetBtc2Usdt(); btc2usdt == 0 {
		errCode = ERROR_CODE_GET_BTC_TO_USDT_RATE_FAILED
		return
	}

	feeInteger = int64(float64(feeUsdtInteger) / btc2usdt)

	// generate transfer order
	logId, err = dao.WealthLogDaoEntity.Add(uid, dao.WealthLogTypeTransferOut, dao.WealthLogStatusOutUnknown, transferInteger, feeInteger, feeUsdtInteger, "", "", to, memo)
	if err != nil {
		errCode = ERROR_CODE_USDT_WEALTH_LOG_FAILED
		return
	}

	// transfer frozen
	err = dao.AccountDaoEntity.TransferFrozen(uid, transferInteger)
	if err != nil {
		errCode = ERROR_CODE_USDT_TRANSFER_FROZEN_FAILED
		dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusOutFailure, dao.OutStepTransferFrozen, fmt.Sprint(err))
		return
	}

	// update log status to submitted
	dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusOutSubmitted, dao.OutStepTransferFrozen, "")

	return
}

// 取消转账申请单
func CancelTransferOutOrder(uid uint64, id uint64) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	log, err := dao.WealthLogDaoEntity.QueryByIdAndUid(id, uid)
	if err != nil {
		errCode = ERROR_CODE_USDT_WEALTH_LOG_FAILED
		return
	}

	// check status
	if log.Status != dao.WealthLogStatusOutSubmitted {
		errCode = ERROR_CODE_USDT_WEALTH_LOG_STATUS_ERR
		return
	}

	// update wealth log to canceled status
	err = dao.WealthLogDaoEntity.UpdateStatusWithCheck(id, dao.WealthLogStatusOutCanceled, dao.WealthLogStatusOutSubmitted)
	if err != nil {
		errCode = ERROR_CODE_USDT_WEALTH_LOG_STATUS_ERR
		return
	}

	// unfrozen
	err = dao.AccountDaoEntity.TransferUnfrozen(log.Uid, log.AmountInteger)
	if err != nil {
		common.LogFuncError("%v", err)
		errCode = ERROR_CODE_USDT_UNFROZEN_FAILED
		dao.WealthLogDaoEntity.UpdateStatus(id, dao.WealthLogStatusOutFailure, dao.OutStepTransferUnfrozenWhileCancel, fmt.Sprint(err))
		return
	}

	return
}

// 审批接口 isApproved: 是否通过
func ApproveTransferOutOrder(id uint64, isApproved bool) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	// check 订单是否存在
	log, err := dao.WealthLogDaoEntity.QueryById(id)
	if err != nil {
		if err == orm.ErrNoRows {
			errCode = ERROR_CODE_USDT_REQUEST_FORM_NO_FOUND
			return
		}
		errCode = ERROR_CODE_DB
		return
	}

	// 必须是提现申请单方可审批
	if log.TType != dao.WealthLogTypeTransferOut {
		errCode = ERROR_CODE_USDT_WEALTH_LOG_STATUS_ERR
		return
	}

	// 为提交的申请单方可审批
	if log.Status != dao.WealthLogStatusOutSubmitted {
		errCode = ERROR_CODE_USDT_WEALTH_LOG_STATUS_ERR
		return
	}

	// 审批拒绝，解冻用户资金，流程结束
	if !isApproved {
		// 更新申请单为拒绝
		err = dao.WealthLogDaoEntity.UpdateStatusWithCheck(id, dao.WealthLogStatusOutRejected, dao.WealthLogStatusOutSubmitted)
		if err != nil {
			errCode = ERROR_CODE_USDT_WEALTH_LOG_STATUS_ERR
			return
		}

		// 解冻资金
		err = dao.AccountDaoEntity.TransferUnfrozen(log.Uid, log.AmountInteger)
		if err != nil {
			common.LogFuncError("%v", err)
			errCode = ERROR_CODE_USDT_FROZEN_WEALTH_ERR
			dao.WealthLogDaoEntity.UpdateStatus(id, dao.WealthLogStatusOutFailure, dao.OutStepTransferUnfrozenWhileReject, fmt.Sprint(err))
			return
		}

		return
	}

	// set to approved status
	err = dao.WealthLogDaoEntity.UpdateStatus(id, dao.WealthLogStatusOutApproved, "", "")
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

const (
	SyncDetectTransactionOutInterval = time.Second * 10
)

//
func DetectUsdtTransferOutOrder() {
	// 获取已发出的转账记录
	submittedTransOrders, err := dao.WealthLogDaoEntity.DumpByTypeAndStatus(dao.WealthLogTypeTransferOut, dao.WealthLogStatusOutTransferred)
	if err != nil {
		return
	}
	for _, v := range submittedTransOrders {
		time.Sleep(SyncDetectTransactionOutInterval)
		resp, err := explorer.NewExplorer().TransactionTX(v.Txid)
		if err != nil {
			common.LogFuncDebug("%v, %v", resp, err)
			continue
		}

		msg := Transaction{}
		err = mapstructure.Decode(resp, &msg)
		if err != nil {
			continue
		}

		// 判断是否确认，确认成功，更新状态, 添加记录
		if msg.Valid == true && msg.Confirmations >= uint32(UsdtConfig.ConfirmationLimit) {
			err = dao.WealthLogDaoEntity.UpdateStatusWithCheck(v.Id, dao.WealthLogStatusOutConfirmed, dao.WealthLogStatusOutTransferred)
			if err != nil {
				common.LogFuncError("decode transfer out order on chain err:%s , %s, %v", ToJson(v), ToJson(msg), err)
				continue
			}
			var amountInteger int64
			amountInteger, err = common.CurrencyStrToInt64(msg.Amount)
			if err != nil {
				common.LogFuncError("%v", err)
				continue
			}

			// 记录提现交易
			AddTransferOutTransaction(msg.Txid, int32(dao.TransactionStatusConfirmed), msg.Block, msg.Sendingaddress, msg.ReferenceAddress, amountInteger, msg.Fee, "block confirm")
		}
	}
}

func GetRecords(uid uint64, types []interface{}, page, limit int64) (list []*WealthLogRecord, meta *common.Meta, err error) {
	count := dao.WealthLogDaoEntity.Count(uid, 0, types)
	meta = common.MakeMeta(count, page, limit)

	tmpList := []*models.UsdtWealthLog{}
	if meta.Offset < meta.Total {
		tmpList, _ = dao.WealthLogDaoEntity.Fetch(uid, types, meta.Limit, meta.Offset)
	}
	list = []*WealthLogRecord{}
	if err != nil {
		return
	}
	for _, v := range tmpList {
		record := &WealthLogRecord{
			Id:     v.Id,
			Uid:    v.Uid,
			Type:   v.TType,
			Status: v.Status,
			Txid:   v.Txid,
			Utime:  v.Utime,
			Ctime:  v.Ctime,
			From:   v.From,
			To:     v.To,
			Memo:   v.Memo,
		}
		record.Amount = common.CurrencyInt64ToStr(v.AmountInteger, UsdtConfig.Precision)
		list = append(list, record)
	}
	return
}

func GetRecord(recordId uint64) (record WealthLogRecord, err error) {
	var log *models.UsdtWealthLog
	log, err = dao.WealthLogDaoEntity.QueryById(recordId)
	if err != nil {
		return
	}

	record.Id = log.Id
	record.Uid = log.Uid
	record.Type = log.TType
	record.Status = log.Status
	record.Txid = log.Txid
	record.Utime = log.Utime
	record.Ctime = log.Ctime
	record.From = log.From
	record.To = log.To
	record.Memo = log.Memo
	record.Amount = common.CurrencyInt64ToStr(log.AmountInteger, UsdtConfig.Precision)
	if log.Txid != "" {
		transaction := &models.UsdtOnchainTransaction{}
		transaction, err = dao.OnchainTransactionDaoEntity.QueryByTxid(log.Txid)
		if err != nil {
			if err == orm.ErrNoRows {
				err = nil
			}
			return
		}
		record.BlockNum = transaction.BlockHash
		record.Fee = common.CurrencyInt64ToStr(transaction.FeeAmountInteger, UsdtConfig.Precision)
	}
	return
}

func GetRecordsByStatus(uid uint64, status uint32, page, limit int64) (list []*WealthLogRecordClient, meta *common.Meta, err error) {
	count := dao.WealthLogDaoEntity.CountByStatus(uid, status)

	meta = common.MakeMeta(count, page, limit)
	var tmpList []*dao.DetailUsdtWealthLog
	if meta.Offset < meta.Total {
		tmpList, _ = dao.WealthLogDaoEntity.FetchByStatus(uid, status, meta.Limit, meta.Offset)
	}
	list = []*WealthLogRecordClient{}
	if err != nil {
		return
	}
	for _, v := range tmpList {
		record := &WealthLogRecordClient{
			Id:         v.Id,
			Uid:        v.Uid,
			Type:       v.TType,
			Status:     v.Status,
			Step:       v.Step,
			Desc:       v.Desc,
			Txid:       v.Txid,
			Amount:     common.CurrencyInt64ToStr(v.AmountInteger, UsdtConfig.Precision),
			Utime:      v.Utime,
			Ctime:      v.Ctime,
			Fee:        common.CurrencyInt64ToStr(v.FeeInteger, UsdtConfig.Precision),
			FeeOnChain: common.CurrencyInt64ToStr(v.FeeOnchainInteger, UsdtConfig.Precision),
			From:       v.From,
			To:         v.To,
			Address:    v.Address,
			Mobile:     v.Mobile,
		}
		//record.Amount, _ = common.EncodeCurrency(v.AmountInteger, v.AmountDecimals, UsdtConfig.Precision)
		list = append(list, record)
	}
	return
}

func GetDetailRecords(uid uint64, status uint32, types []interface{}, strTypes []string, page, limit int64) (list []*WealthLogRecordClient, meta *common.Meta, err error) {
	count := dao.WealthLogDaoEntity.Count(uid, status, types)

	meta = common.MakeMeta(count, page, limit)
	var tmpList []*dao.DetailUsdtWealthLog
	if meta.Offset < meta.Total {
		tmpList, _ = dao.WealthLogDaoEntity.DetailFetch(uid, status, strTypes, int(meta.Limit), int(meta.Offset))
	}
	list = []*WealthLogRecordClient{}
	if err != nil {
		return
	}
	for _, v := range tmpList {
		record := &WealthLogRecordClient{
			Id:         v.Id,
			Uid:        v.Uid,
			Type:       v.TType,
			Status:     v.Status,
			Step:       v.Step,
			Desc:       v.Desc,
			Txid:       v.Txid,
			Amount:     common.CurrencyInt64ToStr(v.AmountInteger, UsdtConfig.Precision),
			Utime:      v.Utime,
			Ctime:      v.Ctime,
			Fee:        common.CurrencyInt64ToStr(v.FeeInteger, UsdtConfig.Precision),
			FeeOnChain: common.CurrencyInt64ToStr(v.FeeOnchainInteger, UsdtConfig.Precision),
			From:       v.From,
			To:         v.To,
			Address:    v.Address,
			Mobile:     v.Mobile,
		}
		//record.Amount, _ = common.EncodeCurrency(v.AmountInteger, v.AmountDecimals, UsdtConfig.Precision)
		list = append(list, record)
	}
	return
}
