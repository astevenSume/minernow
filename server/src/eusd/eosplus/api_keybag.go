package eosplus

import (
	"common"
	"encoding/json"
	"github.com/eoscanada/eos-go"
	. "otc_error"
	dao2 "utils/eos/dao"
	"utils/eusd/dao"
	"utils/eusd/models"
)

type RpcKeyBag struct {
	Api              *eos.API
	KeyBag           *eos.KeyBag
	TokenAccount     eos.AccountName
	ResourcesAccount eos.AccountName
	Symbol           string
	transferAble     bool
}

func (a *RpcKeyBag) UpdateSetting() (err error) {
	a.KeyBag = &eos.KeyBag{}
	a.TokenAccount = eos.AccountName(EosConfig.TokenAccount)
	a.ResourcesAccount = eos.AccountName(EosConfig.ResourcesAccount)
	a.Symbol = EosConfig.Symbol

	if a.transferAble {
		err = a.ImportKeys(EosConfig.TokenAccount)
		if err != nil {
			return
		}
		err = a.ImportKeys(EosConfig.ResourcesAccount)
		if err != nil {
			return
		}
	}

	return
}

func (a *RpcKeyBag) InitApi(url, wallet string) {
	a.Api = eos.New(url, wallet)
	_ = a.UpdateSetting()
	a.Api.SetSigner(a.KeyBag)
	//默认不允许交易， EOS-server才开启交易功能
	a.transferAble = false
}

//开启交易功能
func (a *RpcKeyBag) SetTransferAble() (err error) {
	a.transferAble = true
	err = a.ImportKeys(EosConfig.TokenAccount)
	if err != nil {
		return
	}
	err = a.ImportKeys(EosConfig.ResourcesAccount)
	if err != nil {
		return
	}
	return
}

func (a *RpcKeyBag) ImportKeys(account string) (err error) {
	info, err := dao2.EosAccountKeysEntity.Info(account)
	if err != nil || info.Account == "" || info.PrivateKeyActive == "" {
		common.LogFuncError("Not Found Private Keys:%s", account)
		return
	}
	err = a.KeyBag.ImportPrivateKey(info.PrivateKeyActive)
	if err != nil {
		common.LogFuncError("Import Private Key ERR:%s", info.PrivateKeyActive)
		return
	}
	a.Api.SetSigner(a.KeyBag)
	return
}

func (a *RpcKeyBag) SetDebug() {
	a.Api.Debug = true
}

//获取区块链信息
func (a *RpcKeyBag) GetInfo() (res *eos.InfoResp, err error) {
	res, err = a.Api.GetInfo()
	return
}

//获取账号信息
func (a *RpcKeyBag) GetAccount(accountName string) (res *eos.AccountResp, err error) {
	account := eos.AccountName(accountName)
	res, err = a.Api.GetAccount(account)
	return
}

//获取账户余额
func (a *RpcKeyBag) GetBalance(accountName string) []eos.Asset {
	account := eos.AccountName(accountName)
	res, _ := a.Api.GetCurrencyBalance(account, "EUSD", a.TokenAccount)
	return res
}

//获取账户EOS余额
func (a *RpcKeyBag) GetBalanceEos(accountName string) []eos.Asset {
	account := eos.AccountName(accountName)
	res, _ := a.Api.GetCurrencyBalance(account, "EOS", "eosio.token")
	return res
}

// 获取交易信息
func (a *RpcKeyBag) GetTransfer(id string, blockNum int) (res *eos.TransactionRespV16, err error) {
	res, err = a.Api.GetTransaction(id, blockNum)
	//res, err = a.Api.GetTransactionRaw(id)
	return
}

// 转账信息确认 只确认状态和区块不可逆(还需要确认转账人员 & 金额 & token币）
func (a *RpcKeyBag) IsTransferOK(id string, blockNum int) (ok bool, res *eos.TransactionRespV16, err error) {
	ok = false
	res, err = a.GetTransfer(id, blockNum)
	if err != nil {
		return
	}
	//交易未完成
	if res.Trx.Receipt.Status != "executed" {
		return
	}

	//获取最新区块信息
	chain, err := a.GetInfo()
	if err != nil {
		return
	}
	//未到可逆交易块
	if chain.LastIrreversibleBlockNum < res.LastIrreversibleBlock {
		return
	}

	ok = true
	return
}

// 创建新账号
func (a *RpcKeyBag) AccountCreate(accountCreator, accountNew, ownerKey, activeKey string) (response *eos.PushTransactionResponse, errCode ERROR_CODE, err error) {
	if !a.transferAble {
		errCode = ERROR_CODE_TRANSFER_UNABLE
		return
	}

	quantRam := EosConfig.RamEos
	quantNet := EosConfig.NetEos
	quantCpu := EosConfig.CpuEos

	// 生成创建账户的bin字符串
	outAccount, err := a.newAccountAbi(accountCreator, accountNew, ownerKey, activeKey)
	if err != nil {
		errCode = ERROR_CODE_DELEGATE_TO_ABI_ERROR
		//common.LogFuncDebug("AbiJsonToBin: %s", err)
		return
	}

	// 生成购买内存的bin字符串
	outRam, err := a.buyRamAbi(accountCreator, accountNew, quantRam)
	if err != nil {
		errCode = ERROR_CODE_DELEGATE_TO_ABI_ERROR
		common.LogFuncDebug("Ram ERR: %s", err)
		return
	}

	// 生成购买抵押资源的bin字符串
	actionBw, errCode, err := a.delegateAction(accountCreator, accountNew, quantNet, quantCpu)
	if err != nil {
		return
	}

	//交易签名
	//拼接签名交易数据
	authorization, _ := eos.NewPermissionLevel(accountCreator + "@active")

	actionAccount := eos.Action{
		Account:       "eosio",
		Name:          "newaccount",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: outAccount.String(),
		},
	}
	actionRam := eos.Action{
		Account:       "eosio",
		Name:          "buyram",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: outRam.String(),
		},
	}

	actionList := []*eos.Action{
		&actionAccount,
		&actionRam,
		&actionBw,
	}

	response, errCode, err = a.signPushTransfer(actionList)
	if errCode == ERROR_CODE_SUCCESS {
		transaction := &models.EosTransaction{
			Type:          TransactionTypeNewAccount,
			TransactionId: response.TransactionId,
			BlockNum:      response.Processed.BlockNum,
			Payer:         accountCreator,
			Receiver:      accountNew,
			Quantity:      quantRam + " " + quantNet + " " + quantCpu,
			Ctime:         common.NowInt64MS(),
		}
		tid, err := dao.TransactionDaoEntity.Create(transaction)
		if err != nil {
			common.LogFuncError("NewAccount Into DB Err:%v", ToJson(transaction))
		}
		quantRamInt, _ := QuantityToUint64(quantRam)
		log := &models.EosUseLog{
			Type:        TransactionTypeNewAccount,
			Tid:         uint64(tid),
			Status:      EosUseStatusIng,
			QuantityNum: quantRamInt,
		}

		qNet, _ := QuantityToUint64(quantNet)
		qCpu, _ := QuantityToUint64(quantCpu)
		qunt := qNet + qCpu
		log2 := &models.EosUseLog{
			Type:        TransactionTypeDelegate,
			Tid:         uint64(tid),
			Status:      EosUseStatusIng,
			QuantityNum: qunt,
		}
		_, err = dao.EosUseLogDaoEntity.InsertMulti([]*models.EosUseLog{log, log2})
		if err != nil {
			common.LogFuncError("NewAccount UseLog Into DB Err:%v", ToJson([]*models.EosUseLog{log, log2}))
		}
	}

	return
}

// 打包新账号交易
func (a *RpcKeyBag) newAccountAbi(accountCreator, accountNew, ownerKey, activeKey string) (out eos.HexBytes, err error) {
	// 生成创建账户的bin字符串
	args := map[string]interface{}{
		"creator": accountCreator,
		"name":    accountNew,
		"owner": map[string]interface{}{
			"threshold": 1,
			"keys": []interface{}{
				map[string]interface{}{
					"key":    ownerKey,
					"weight": 1,
				},
			},
			"accounts": []interface{}{},
			"waits":    []interface{}{},
		},
		"active": map[string]interface{}{
			"threshold": 1,
			"keys": []interface{}{
				map[string]interface{}{
					"key":    activeKey,
					"weight": 1,
				},
			},
			"accounts": []interface{}{},
			"waits":    []interface{}{},
		},
	}
	out, err = a.Api.ABIJSONToBin("eosio", "newaccount", args)
	return
}

// 购买存储ram(byte)
func (a *RpcKeyBag) buyRamBytesAbi(payer, receiver string, bytes int) (out eos.HexBytes, err error) {
	// 生成购买内存的bin字符串
	argsRam := map[string]interface{}{
		"payer":    payer,
		"receiver": receiver,
		"bytes":    bytes,
	}
	out, err = a.Api.ABIJSONToBin("eosio", "buyrambytes", argsRam)
	return
}

// 购买存储ram
func (a *RpcKeyBag) buyRamAbi(payer, receiver string, quant string) (out eos.HexBytes, err error) {
	// 生成购买内存的bin字符串
	argsRam := map[string]interface{}{
		"payer":    payer,
		"receiver": receiver,
		"quant":    quant,
	}
	out, err = a.Api.ABIJSONToBin("eosio", "buyram", argsRam)
	return
}

//购买存储 0.5%手续费
func (a *RpcKeyBag) BuyRam(payer, receiver, quantity string) (response *eos.PushTransactionResponse, errCode ERROR_CODE, err error) {
	out, err := a.buyRamAbi(payer, receiver, quantity)
	if err != nil {
		errCode = ERROR_CODE_DELEGATE_TO_ABI_ERROR
		common.LogFuncDebug("AbiJsonToBin: %s", err)
		return
	}
	authorization, _ := eos.NewPermissionLevel(payer + "@active")
	action := eos.Action{
		Account:       "eosio",
		Name:          "buyram",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: out.String(),
		},
	}

	actionList := []*eos.Action{
		&action,
	}
	response, errCode, err = a.signPushTransfer(actionList)
	if errCode == ERROR_CODE_SUCCESS {
		transaction := &models.EosTransaction{
			Type:          TransactionTypeBuyRam,
			TransactionId: response.TransactionId,
			BlockNum:      response.Processed.BlockNum,
			Payer:         payer,
			Receiver:      receiver,
			Quantity:      quantity,
			Ctime:         common.NowInt64MS(),
		}
		_, err := dao.TransactionDaoEntity.Create(transaction)
		if err != nil {
			common.LogFuncError("BuyRam Into DB Err:%v", ToJson(transaction))
		}
	}

	return
}

// 抵押CPU & NET
func (a *RpcKeyBag) delegateBwAbi(payer, receiver, quantityNet, quantityCpu string) (out eos.HexBytes, err error) {
	// 生成购买内存的bin字符串
	argsRam := map[string]interface{}{
		"from":               payer,
		"receiver":           receiver,
		"stake_net_quantity": quantityNet,
		"stake_cpu_quantity": quantityCpu,
		"transfer":           false,
	}
	out, err = a.Api.ABIJSONToBin("eosio", "delegatebw", argsRam)
	return
}

//抵押CPU & NET Action
func (a *RpcKeyBag) delegateAction(payer, receiver, quantityNet, quantityCpu string) (action eos.Action, errCode ERROR_CODE, err error) {
	errCode = ERROR_CODE_SUCCESS
	outBw, err := a.delegateBwAbi(payer, receiver, quantityNet, quantityCpu)
	if err != nil {
		errCode = ERROR_CODE_DELEGATE_TO_ABI_ERROR
		common.LogFuncDebug("AbiJsonToBin: %s", err)
		return
	}
	authorization, _ := eos.NewPermissionLevel(payer + "@active")
	action = eos.Action{
		Account:       "eosio",
		Name:          "delegatebw",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: outBw.String(),
		},
	}
	return
}

// 抵押获取资源
func (a *RpcKeyBag) DelegateBw(payer, receiver, quantityNet, quantityCpu string) (response *eos.PushTransactionResponse, errCode ERROR_CODE, err error) {
	action, errCode, err := a.delegateAction(payer, receiver, quantityNet, quantityCpu)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}

	actionList := []*eos.Action{
		&action,
	}
	response, errCode, err = a.signPushTransfer(actionList)
	if errCode == ERROR_CODE_SUCCESS {
		transaction := &models.EosTransaction{
			Type:          TransactionTypeDelegate,
			TransactionId: response.TransactionId,
			BlockNum:      response.Processed.BlockNum,
			Payer:         payer,
			Receiver:      receiver,
			Quantity:      quantityNet + " " + quantityCpu,
			Ctime:         common.NowInt64MS(),
		}
		tid, err := dao.TransactionDaoEntity.Create(transaction)
		if err != nil {
			common.LogFuncError("Delegate Into DB Err:%v", ToJson(transaction))
		}
		qNet, _ := QuantityToUint64(quantityNet)
		qCpu, _ := QuantityToUint64(quantityCpu)
		qunt := qNet + qCpu
		//eos 使用记录
		log := &models.EosUseLog{
			Type:        TransactionTypeDelegate,
			Tid:         uint64(tid),
			Status:      EosUseStatusIng,
			QuantityNum: qunt,
		}
		_, err = dao.EosUseLogDaoEntity.Create(log)
		if err != nil {
			common.LogFuncError("EosUseLog Into DB Err:%v", ToJson(transaction))
		}
	}

	return
}

// 赎回获取资源 （EOS只能返回抵押的账号） 多次赎回资源到账时间以最后一次重新计时72H！！！
// payer 抵押资源所在账号
// receiver 接收账号
func (a *RpcKeyBag) UnDelegateBw(payer, receiver, quantityNet, quantityCpu string) (response *eos.PushTransactionResponse, errCode ERROR_CODE, err error) {
	// 生成购买内存的bin字符串
	args := map[string]interface{}{
		"from":                 payer,
		"receiver":             receiver,
		"unstake_net_quantity": quantityNet,
		"unstake_cpu_quantity": quantityCpu,
	}
	out, err := a.Api.ABIJSONToBin("eosio", "undelegatebw", args)
	if err != nil {
		errCode = ERROR_CODE_UNDELEGATE_TO_ABI_ERROR
		common.LogFuncDebug("AbiJsonToBin: %s", err)
		return
	}
	authorization, _ := eos.NewPermissionLevel(payer + "@active")
	action := eos.Action{
		Account:       "eosio",
		Name:          "undelegatebw",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: out.String(),
		},
	}

	actionList := []*eos.Action{
		&action,
	}
	response, errCode, err = a.signPushTransfer(actionList)
	if errCode == ERROR_CODE_SUCCESS {
		transaction := &models.EosTransaction{
			Type:          TransactionTypeUnDelegate,
			TransactionId: response.TransactionId,
			BlockNum:      response.Processed.BlockNum,
			Payer:         payer,
			Receiver:      receiver,
			Quantity:      quantityNet + " " + quantityCpu,
			Ctime:         common.NowInt64MS(),
		}
		_, err := dao.TransactionDaoEntity.Create(transaction)
		if err != nil {
			common.LogFuncError("UnDelegate Into DB Err:%v", ToJson(transaction))
		}
	}

	return
}

// 转账Abi
func (a *RpcKeyBag) transferAbi(token, from, to, quantity, memo string) (out eos.HexBytes, err error) {
	actionData := map[string]interface{}{
		"from":     from,
		"to":       to,
		"quantity": quantity,
		"memo":     memo,
	}

	out, err = a.Api.ABIJSONToBin(eos.AccountName(token), "transfer", actionData)
	return
}

// 转账
func (a *RpcKeyBag) Transfer(token, from, to, quantity, memo string) (response *eos.PushTransactionResponse, transaction *models.EosTransaction, errCode ERROR_CODE) {
	if !a.transferAble {
		errCode = ERROR_CODE_TRANSFER_UNABLE
		return
	}
	//errCode = a.transferCheckResource(from, to)

	//导入私钥
	err := a.ImportKeys(from)
	if err != nil {
		return
	}

	// 转账Abi
	transfer, err := a.transferAbi(token, from, to, quantity, memo)
	if err != nil {
		errCode = ERROR_CODE_ACTION_TO_ABI_ERROR
		common.LogFuncDebug("AbiJsonToBin: %s", err)
		return
	}

	//交易签名
	//拼接签名交易数据
	authorization, _ := eos.NewPermissionLevel(from + "@active")

	actionTransfer := eos.Action{
		Account:       eos.AccountName(token),
		Name:          "transfer",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: transfer.String(),
		},
	}
	actionList := []*eos.Action{
		&actionTransfer,
	}

	response, errCode, err = a.signPushTransfer(actionList)
	if errCode == ERROR_CODE_SUCCESS {
		transaction = &models.EosTransaction{
			Type:          TransactionTypeTransfer,
			TransactionId: response.TransactionId,
			BlockNum:      response.Processed.BlockNum,
			Payer:         from,
			Receiver:      to,
			Quantity:      quantity,
			Ctime:         common.NowInt64MS(),
			Memo:          memo,
		}
		id, err := dao.TransactionDaoEntity.Create(transaction)
		if err != nil {
			common.LogFuncError("Transfer Into DB Err:%v", ToJson(transaction))
		}
		transaction.Id = uint64(id)
	}

	return
}

// 交易签名 & 推送交易 PS:调用这个func的方法需要往数据库写入交易数据
func (a *RpcKeyBag) signPushTransfer(actionList []*eos.Action) (response *eos.PushTransactionResponse, errCode ERROR_CODE, err error) {
	if !a.transferAble {
		errCode = ERROR_CODE_TRANSFER_UNABLE
		return
	}

	errCode = ERROR_CODE_SUCCESS
	txOpts := &eos.TxOptions{}
	err = txOpts.FillFromChain(a.Api)
	if err != nil {
		errCode = ERROR_CODE_GET_BLOCK_INFO_ERROR
		common.LogFuncDebug("Filling_Tx_Opts: %s", err)
		return
	}

	tx := eos.NewTransaction(actionList, txOpts)
	//发送签名
	signedTx, packedTx, err := a.Api.SignTransaction(tx, txOpts.ChainID, eos.CompressionNone)

	if err != nil {
		errCode = ERROR_CODE_SIGN_TRANSACTION_FAIL
		common.LogFuncDebug("Sign_Transaction: %s", err)
		return
	}

	//推送交易
	err = a.Api.Call("chain", "push_transaction", packedTx, &response)
	if err != nil {
		errCode = ERROR_CODE_PUSH_TRANSACTION_ERROR
		content, _ := json.MarshalIndent(signedTx, "", "  ")

		common.LogFuncDebug("Push_Transaction:%v ||Content:%v", err, string(content))
		return
	}

	//交易数据入库
	data := &models.EosTransactionInfo{
		TransactionId: response.TransactionId,
		BlockNum:      response.Processed.BlockNum,
		Processed:     ToJson(response.Processed),
		Ctime:         common.NowInt64MS(),
	}

	if _, err = dao.TransactionInfoDaoEntity.Create(data); err != nil {
		common.LogFuncError("TransInfo DB:%v", ToJson(data))
	}

	return
}

//转账资源检查
func (a *RpcKeyBag) transferCheckResource(fromAccount, toAccount string) (errCode ERROR_CODE) {
	accInfo, err := a.GetAccount(fromAccount)
	if err != nil {
		common.LogFuncError("getAccount ERR:%v", err)
		return
	}
	// 如果 内存小于250byte 转账可能失败
	if accInfo.RAMQuota-accInfo.RAMUsage < 250 {
		errCode = a.transferCheckResourceBuyRam(toAccount)
		if errCode != ERROR_CODE_SUCCESS {
			return
		}
	}

	needCpu, needNet := false, false
	//转账需要几十的CPU时间
	if accInfo.CPULimit.Available < 100 {
		needCpu = true
	}
	if accInfo.NetLimit.Available < 100 {
		needNet = true
	}

	// 补充net&cpu抵押资源
	if needCpu || needNet {
		a.transferCheckResourceDelegate(toAccount)
		if errCode != ERROR_CODE_SUCCESS {
			return
		}
	}
	return ERROR_CODE_SUCCESS
}

//购买内存
func (a *RpcKeyBag) transferCheckResourceBuyRam(to string) (errCode ERROR_CODE) {
	admin := string(a.ResourcesAccount)
	if admin == "" {
		errCode = ERROR_CODE_CONFIG_LACK
		return
	}
	quantity := EosConfig.RamEos

	_, errCode, _ = a.BuyRam(admin, to, quantity)

	return
}

//抵押cpu
func (a *RpcKeyBag) transferCheckResourceDelegate(to string) (errCode ERROR_CODE) {
	admin := string(a.ResourcesAccount)
	if admin == "" {
		errCode = ERROR_CODE_CONFIG_LACK
		return
	}
	cpu := EosConfig.CpuEos
	net := EosConfig.NetEos
	if cpu == "" || net == "" {
		errCode = ERROR_CODE_CONFIG_LACK
		return
	}

	_, errCode, _ = a.DelegateBw(admin, to, net, cpu)
	return
}

//从最大量 issue token
func (a *RpcKeyBag) IssueToken(receiver, quantity string) (response *eos.PushTransactionResponse, transaction *models.EosTransaction, errCode ERROR_CODE) {
	argsRam := map[string]interface{}{
		"to":       receiver,
		"quantity": quantity,
		"memo":     "memo",
	}
	singer := string(a.TokenAccount)
	out, err := a.Api.ABIJSONToBin(a.TokenAccount, "issue", argsRam)
	if err != nil {
		common.LogFuncError("%v", err)
	}
	authorization, _ := eos.NewPermissionLevel(singer + "@active")
	action := eos.Action{
		Account:       a.TokenAccount,
		Name:          "issue",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: out.String(),
		},
	}

	actionList := []*eos.Action{
		&action,
	}
	response, errCode, err = a.signPushTransfer(actionList)
	if errCode == ERROR_CODE_SUCCESS {
		transaction = &models.EosTransaction{
			Type:          TransactionTypeIssueToken,
			TransactionId: response.TransactionId,
			BlockNum:      response.Processed.BlockNum,
			Payer:         singer,
			Receiver:      receiver,
			Quantity:      quantity,
			Ctime:         common.NowInt64MS(),
		}
		id, err := dao.TransactionDaoEntity.Create(transaction)
		if err != nil {
			common.LogFuncError("Transfer Into DB Err:%v", ToJson(transaction))
		}
		transaction.Id = uint64(id)
	}

	return
}

//销毁token retire token
func (a *RpcKeyBag) RetireToken(quantity, memo string) (response *eos.PushTransactionResponse, transaction *models.EosTransaction, errCode ERROR_CODE) {
	argsRam := map[string]interface{}{
		"quantity": quantity,
		"memo":     memo,
	}
	out, err := a.Api.ABIJSONToBin(a.TokenAccount, "retire", argsRam)
	if err != nil {
		common.LogFuncError("%v", err)
	}
	singer := string(a.TokenAccount)
	authorization, _ := eos.NewPermissionLevel(singer + "@active")
	action := eos.Action{
		Account:       a.TokenAccount,
		Name:          "retire",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: out.String(),
		},
	}

	actionList := []*eos.Action{
		&action,
	}
	response, errCode, err = a.signPushTransfer(actionList)
	if errCode == ERROR_CODE_SUCCESS {
		transaction = &models.EosTransaction{
			Type:          TransactionTypeRetireToken,
			TransactionId: response.TransactionId,
			BlockNum:      response.Processed.BlockNum,
			Payer:         singer,
			//Receiver:      receiver,
			Quantity: quantity,
			Ctime:    common.NowInt64MS(),
		}
		id, err := dao.TransactionDaoEntity.Create(transaction)
		if err != nil {
			common.LogFuncError("Transfer Into DB Err:%v", ToJson(transaction))
		}
		transaction.Id = uint64(id)
	}
	return
}

// 充值EOS储备金
func (a *RpcKeyBag) RexDeposit(amount string) (response *eos.PushTransactionResponse, transaction *models.EosTransaction, errCode ERROR_CODE) {
	args := map[string]interface{}{
		"owner":  a.ResourcesAccount,
		"amount": amount,
	}
	singer := string(a.ResourcesAccount)
	out, err := a.Api.ABIJSONToBin("eosio", "deposit", args)
	if err != nil {
		common.LogFuncError("%v", err)
	}
	authorization, _ := eos.NewPermissionLevel(singer + "@active")
	action := eos.Action{
		Account:       "eosio",
		Name:          "deposit",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: out.String(),
		},
	}

	actionList := []*eos.Action{
		&action,
	}
	response, errCode, err = a.signPushTransfer(actionList)

	if errCode == ERROR_CODE_SUCCESS {
		transaction = &models.EosTransaction{
			Type:          TransactionTypeRexDeposit,
			TransactionId: response.TransactionId,
			BlockNum:      response.Processed.BlockNum,
			Payer:         singer,
			Receiver:      "",
			Quantity:      amount,
			Ctime:         common.NowInt64MS(),
		}
		id, err := dao.TransactionDaoEntity.Create(transaction)
		if err != nil {
			common.LogFuncError("RexDeposit DB Err:%v", ToJson(transaction))
		}
		transaction.Id = uint64(id)
	}

	return
}

// 提现EOS储备金
func (a *RpcKeyBag) RexWithdraw(amount string) (response *eos.PushTransactionResponse, transaction *models.EosTransaction, errCode ERROR_CODE) {
	args := map[string]interface{}{
		"owner":  a.ResourcesAccount,
		"amount": amount,
	}
	singer := string(a.ResourcesAccount)
	out, err := a.Api.ABIJSONToBin("eosio", "withdraw", args)
	if err != nil {
		common.LogFuncError("%v", err)
	}
	authorization, _ := eos.NewPermissionLevel(singer + "@active")
	action := eos.Action{
		Account:       "eosio",
		Name:          "withdraw",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: out.String(),
		},
	}

	actionList := []*eos.Action{
		&action,
	}
	response, errCode, err = a.signPushTransfer(actionList)

	if errCode == ERROR_CODE_SUCCESS {
		transaction = &models.EosTransaction{
			Type:          TransactionTypeRexWithdraw,
			TransactionId: response.TransactionId,
			BlockNum:      response.Processed.BlockNum,
			Payer:         "",
			Receiver:      singer,
			Quantity:      amount,
			Ctime:         common.NowInt64MS(),
		}
		id, err := dao.TransactionDaoEntity.Create(transaction)
		if err != nil {
			common.LogFuncError("RexWithdraw DB Err:%v", ToJson(transaction))
		}
		transaction.Id = uint64(id)
	}

	return
}

// 购买cpu
func (a *RpcKeyBag) RexRentCpu(receiver, loanPayment string) (response *eos.PushTransactionResponse, transaction *models.EosTransaction, errCode ERROR_CODE) {
	args := map[string]interface{}{
		"from":         a.ResourcesAccount,
		"receiver":     receiver,
		"loan_payment": loanPayment,
		"loan_fund":    "0.0000 EOS",
	}
	singer := string(a.ResourcesAccount)
	out, err := a.Api.ABIJSONToBin("eosio", "rentcpu", args)
	if err != nil {
		common.LogFuncError("%v", err)
	}
	authorization, _ := eos.NewPermissionLevel(singer + "@active")
	action := eos.Action{
		Account:       "eosio",
		Name:          "rentcpu",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: out.String(),
		},
	}

	actionList := []*eos.Action{
		&action,
	}
	response, errCode, err = a.signPushTransfer(actionList)

	if errCode == ERROR_CODE_SUCCESS {
		transaction = &models.EosTransaction{
			Type:          TransactionTypeRexRentCpu,
			TransactionId: response.TransactionId,
			BlockNum:      response.Processed.BlockNum,
			Payer:         "",
			Receiver:      singer,
			Quantity:      loanPayment,
			Ctime:         common.NowInt64MS(),
		}
		id, err := dao.TransactionDaoEntity.Create(transaction)
		if err != nil {
			common.LogFuncError("RexRentCpu DB Err:%v", ToJson(transaction))
		}
		transaction.Id = uint64(id)
	}

	return
}

// 购买net  loan_payment租金  loan_fund备用金
func (a *RpcKeyBag) RexRentNet(receiver, loanPayment string) (response *eos.PushTransactionResponse, transaction *models.EosTransaction, errCode ERROR_CODE) {
	args := map[string]interface{}{
		"from":         a.ResourcesAccount,
		"receiver":     receiver,
		"loan_payment": loanPayment,
		"loan_fund":    "0.0000 EOS",
	}
	singer := string(a.ResourcesAccount)
	out, err := a.Api.ABIJSONToBin("eosio", "rentnet", args)
	if err != nil {
		common.LogFuncError("%v", err)
	}
	authorization, _ := eos.NewPermissionLevel(singer + "@active")
	action := eos.Action{
		Account:       "eosio",
		Name:          "rentnet",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: out.String(),
		},
	}

	actionList := []*eos.Action{
		&action,
	}
	response, errCode, err = a.signPushTransfer(actionList)
	if errCode == ERROR_CODE_SUCCESS {
		transaction = &models.EosTransaction{
			Type:          TransactionTypeRexRentNet,
			TransactionId: response.TransactionId,
			BlockNum:      response.Processed.BlockNum,
			Payer:         singer,
			Receiver:      receiver,
			Quantity:      loanPayment,
			Ctime:         common.NowInt64MS(),
		}
		id, err := dao.TransactionDaoEntity.Create(transaction)
		if err != nil {
			common.LogFuncError("RexRentNet DB Err:%v", ToJson(transaction))
		}
		transaction.Id = uint64(id)
	}

	return
}

//转账资源检查 提供给外部调用
func (a *RpcKeyBag) CheckResource(checkAccount string) (errCode ERROR_CODE) {
	accInfo, err := a.GetAccount(checkAccount)
	if err != nil {
		common.LogFuncError("getAccount ERR:%v", err)
		return
	}
	// 如果 内存小于250byte 转账可能失败
	if accInfo.RAMQuota-accInfo.RAMUsage < 250 {
		errCode = a.transferCheckResourceBuyRam(checkAccount)
		if errCode != ERROR_CODE_SUCCESS {
			return
		}
	}

	needCpu, needNet := false, false
	//转账需要几十的CPU时间
	if accInfo.CPULimit.Available < 1000 {
		needCpu = true
	}
	if accInfo.NetLimit.Available < 300 {
		needNet = true
	}

	// 补充net&cpu抵押资源
	if needCpu || needNet {
		a.transferCheckResourceDelegate(checkAccount)
		if errCode != ERROR_CODE_SUCCESS {
			return
		}
	}
	return ERROR_CODE_SUCCESS
}

func (a *RpcKeyBag) Vote() (errCode ERROR_CODE) {
	//./cleos push action eosio voteproducer '{"voter":"voterperson1","proxy":"","producers":["accountpro1"]}' -p 1234512345ac
	args := map[string]interface{}{
		"voter":     a.ResourcesAccount,
		"proxy":     "",
		"producers": []string{"eoshuobipool"},
	}
	singer := string(a.TokenAccount)
	out, err := a.Api.ABIJSONToBin("eosio", "voteproducer", args)
	if err != nil {
		common.LogFuncError("%v", err)
	}
	authorization, _ := eos.NewPermissionLevel(singer + "@active")
	action := eos.Action{
		Account:       "eosio",
		Name:          "voteproducer",
		Authorization: []eos.PermissionLevel{authorization},
		ActionData: eos.ActionData{
			Data: out.String(),
		},
	}

	actionList := []*eos.Action{
		&action,
	}
	response, errCode, err := a.signPushTransfer(actionList)
	common.LogFuncDebug("%v", ToJsonIndent(response))
	return
}
