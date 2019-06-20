package usdt

import (
	"common"
	"errors"
	. "otc_error"
	"time"
	otcdao "utils/otc/dao"
	otcmodels "utils/otc/models"
	usdtdao "utils/usdt/dao"
	usdtmodels "utils/usdt/models"

	"github.com/astaxie/beego/orm"
)

//只记录成功的交易
const (
	SyncRechargeTransactionInterval = time.Second * 10
)

// 新增转账（提现）交易记录
func AddTransferOutTransaction(txid string, status int32, blockNum uint32, payer, receiver string, amountInteger int64, fee, memo string) (err error) {
	transaction := &usdtmodels.UsdtTransaction{
		TxId:          txid,
		Type:          uint8(usdtdao.TransactionTypeTransferout),
		BlockNum:      blockNum,
		Status:        status,
		Payer:         payer,
		Receiver:      receiver,
		AmountInteger: amountInteger,
		Fee:           fee,
		Memo:          memo,
	}
	_, err = usdtdao.TransactionDaoEntity.Create(transaction)
	if err != nil {
		common.LogFuncError("transaction Create err, %v json: %s", err, ToJson(transaction))
	}
	return
}

// 新增转账（归集）交易记录
func AddCollectionTransaction(txid string, status int32, blockNum uint32, payer, receiver string, amountInteger int64, memo string) (err error) {
	transaction := &usdtmodels.UsdtTransaction{
		TxId:          txid,
		Type:          uint8(usdtdao.TransactionTypeCollection),
		BlockNum:      blockNum,
		Status:        status,
		Payer:         payer,
		Receiver:      receiver,
		AmountInteger: amountInteger,
		Memo:          memo,
	}
	_, err = usdtdao.TransactionDaoEntity.Create(transaction)
	if err != nil {
		common.LogFuncError("transaction Create err, %v json: %s", err, ToJson(transaction))
	}
	return
}

// 新增转账（充值）交易记录
func AddTransferInTransaction(txid string, status int32, blockNum uint32, payer, receiver string, amountInteger int64, memo string) (err error) {
	transaction := &usdtmodels.UsdtTransaction{
		TxId:          txid,
		Type:          usdtdao.TransactionTypeTransferin,
		BlockNum:      blockNum,
		Status:        status,
		Payer:         payer,
		Receiver:      receiver,
		AmountInteger: amountInteger,
		Memo:          memo,
	}
	_, err = usdtdao.TransactionDaoEntity.Create(transaction)
	if err != nil {
		common.LogFuncError("transaction Create err, %v json: %s", err, ToJson(transaction))
	}
	return
}

// 同步充值记录
func SyncRechargeTransaction() (err error) {
	var total int
	total, err = usdtdao.AccountDaoEntity.QueryTotalRelated()
	if err != nil || total <= 0 {
		return
	}

	perPage := 100
	pages := total / perPage
	if total%perPage > 0 {
		pages += 1
	}

	var list []usdtmodels.UsdtAccount
	for i := 1; i <= pages; i++ { //per page
		// get usdt account
		list, err = usdtdao.AccountDaoEntity.QueryRelated(i, perPage)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}

		for _, v := range list { // loop usdt account
			// 同步充值记录
			SyncRechargeTransactionByAccount(v)

			// 间隔10秒同步下一个帐号
			// todo use multiple omni server
			time.Sleep(SyncRechargeTransactionInterval)
		}
	}

	return
}

// 同步钱包记录
func SyncWalletTransaction() (err error) {

	list, err := usdtdao.AccountDaoEntity.QueryWalletAccount(usdtdao.HOT_WALLET_MIN_PKID, usdtdao.COLD_WALLET_MAX_PKID)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	for _, v := range list { // loop usdt account
		// 同步充值记录
		SyncRechargeTransactionByAccount(v)

		// 间隔10秒同步下一个帐号
		// todo use multiple omni server
		time.Sleep(SyncRechargeTransactionInterval)
	}

	return
}

func SyncRechargeTransactionByMobile(nationalCode, mobile string) (err error) {
	user := &otcmodels.User{}
	user, err = otcdao.UserDaoEntity.InfoByMobile(nationalCode, mobile)
	if err != nil {
		return
	}
	if user == nil || user.Uid <= 0 {
		return errors.New("no found account by nationalCode, mobile")
	}

	err = SyncRechargeTransactionByUid(user.Uid)
	return
}

func SyncRechargeTransactionByUid(uid uint64) (err error) {
	account := &usdtmodels.UsdtAccount{}
	account, err = usdtdao.AccountDaoEntity.QueryByUid(uid)
	if err != nil {
		return
	}
	if account == nil || account.Uid <= 0 || account.Uaid <= 0 {
		err = errors.New("no usdt account")
		return
	}
	SyncRechargeTransactionByAccount(*account)
	return
}

func SyncRechargeTransactionByAccount(account usdtmodels.UsdtAccount) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	// 计算链上数据和本地数据差异， 计算差值
	// 同步链上数据
	errCode = SyncTransaction(account, true)
	if errCode != ERROR_CODE_SUCCESS {
		return errCode
	}

	// 获取本地充值记录
	listDb := []usdtmodels.UsdtTransaction{}
	var (
		rechargeIntergerDb int64
		err                error
	)
	listDb, rechargeIntergerDb, err = usdtdao.TransactionDaoEntity.QueryAllRecharge(account.Uaid)
	if err != nil && err != orm.ErrNoRows {
		errCode = ERROR_CODE_DB
		return
	}

	// 获取链上充值记录
	var (
		listOnChain, listSweep, listWealth []usdtmodels.UsdtOnchainTransaction
	)
	var rechargeIntegerOnChain, rechargeIntegerTransferFrozen int64
	listOnChain, listSweep, listWealth, rechargeIntegerOnChain, rechargeIntegerTransferFrozen, err =
		usdtdao.OnchainTransactionDaoEntity.QueryOnchainAllRecharge(account.Uaid, uint32(UsdtConfig.ConfirmationLimit))
	if err != nil && err != orm.ErrNoRows {
		errCode = ERROR_CODE_DB
		return
	}

	// 计算amount差额
	rechargeInteger := rechargeIntegerOnChain - rechargeIntergerDb

	// 计算新增充值记录
	if len(listOnChain) <= 0 {
		return
	}
	list := []usdtmodels.UsdtTransaction{}
	now := common.NowInt64MS()
	listLog := []usdtmodels.UsdtWealthLog{}
	for _, v := range listOnChain {
		isEq := false

		for _, v1 := range listDb {
			if v.TxId == v1.TxId {
				isEq = true
				break
			}
		}

		if !isEq {
			list = append(list, usdtmodels.UsdtTransaction{
				Uaid:          v.Uaid,
				Type:          v.Type,
				Status:        int32(usdtdao.TransactionStatusConfirmed),
				TxId:          v.TxId,
				BlockNum:      v.Block,
				Payer:         v.SendingAddress,
				Receiver:      v.ReferenceAddress,
				AmountInteger: v.AmountInteger,
				Ctime:         v.BlockTime,
				Utime:         now,
			})

			var logType, logStatus uint32
			switch v.Type {
			case uint8(usdtdao.TransactionTypeTransferin):
				logType = usdtdao.WealthLogTypeTransferIn
				logStatus = usdtdao.WealthLogStatusDeposited
			case uint8(usdtdao.TransactionTypeTransferout):
				logType = usdtdao.WealthLogTypeTransferOut
				logStatus = usdtdao.WealthLogStatusOutConfirmed
			default:
				logType = usdtdao.WealthLogTypeUnkown
			}

			var hasLog bool
			// 判断是否存在 sweep log，存在则跳过写入 wealth log 的步骤
			for _, sweep := range listSweep {
				if v.TxId == sweep.TxId {
					// 更新
					if err = usdtdao.SweepLogDaoEntity.UpdateByTxID(v.TxId, v.SendingAddress, v.ReferenceAddress, v.AmountInteger, v.FeeAmountInteger, v.BlockTime); err != nil {
						errCode = ERROR_CODE_DB
						return
					}
					hasLog = true
					break
				}
			}
			// 判断是否存在 wealth log ，存在则做更新，不存在则写入
			for _, wealth := range listWealth {
				if v.TxId == wealth.TxId {
					// 更新
					if err = usdtdao.WealthLogDaoEntity.UpdateByTxID(v.TxId, v.SendingAddress, v.ReferenceAddress, logType, logStatus, v.AmountInteger, v.FeeAmountInteger, v.BlockTime); err != nil {
						errCode = ERROR_CODE_DB
						return
					}
					hasLog = true
					break
				}
			}
			if !hasLog {
				// FeeAmountInteger
				var logTmp usdtmodels.UsdtWealthLog
				logTmp, err = usdtdao.WealthLogDaoEntity.Gen(account.Uid,
					logType, logStatus, v.AmountInteger, v.BlockTime, v.Mtime,
					v.TxId, v.SendingAddress, v.ReferenceAddress)

				logTmp.FeeOnchainInteger = v.FeeAmountInteger

				listLog = append(listLog, logTmp)
			}
		}
	}

	// add unsaved transactions
	if len(list) > 0 {
		if err = usdtdao.TransactionDaoEntity.InsertMulti(list); err != nil {
			return
		}

		// update available and transfer frozen
		err = usdtdao.AccountDaoEntity.DepositAndClearTransferFrozen(account.Uid, rechargeInteger, rechargeIntegerTransferFrozen)
		if err != nil {
			errCode = ERROR_CODE_USDT_DEPOSIT_FAILED
			return
		}
	}

	if len(listLog) > 0 {
		err = usdtdao.WealthLogDaoEntity.InsertMulti(listLog)
		if err != nil {
			errCode = ERROR_CODE_USDT_WEALTH_LOG_FAILED
			return
		}
	}

	return
}
