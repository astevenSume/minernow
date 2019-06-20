package usdt

import (
	"common"
	"eusd/eosplus"
	"fmt"
	. "otc_error"
	otcerror "otc_error"
	"strconv"
	"usdt/explorer"
	"utils/usdt/dao"
	"utils/usdt/models"

	"github.com/astaxie/beego/orm"
	lock "github.com/bsm/redis-lock"
	"github.com/mitchellh/mapstructure"
)

// create usdt account for user identified by uid
func CreateAccount(uid uint64) (account *models.UsdtAccount, err error) {
	//generate account data
	//var account *models.UsdtAccount
	account, err = dao.AccountDaoEntity.Create(uid)
	if err != nil {
		return
	}

	return
}

type Balance struct {
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Mortgaged string `json:"mortgaged"`
}

type UsdtMsg struct {
	Address   string  `json:"address"`
	Status    uint8   `json:"status"`
	Symbal    string  `json:"symbal"`
	Precision int     `json:"precision"`
	Balance   Balance `json:"balance"`
}

// query usdt account by uid
func QueryByUid(uid uint64) (msg UsdtMsg, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	var account *models.UsdtAccount
	var err error
	account, err = dao.AccountDaoEntity.QueryByUid(uid)
	if err != nil {
		if err == orm.ErrNoRows {
			// try create usdt account
			account, err = CreateAccount(uid)
			if err != nil {
				errCode = ERROR_CODE_USDT_CREATE_ACCOUNT_FAILED
				return
			}
		}
		errCode = ERROR_CODE_DB
		return
	}

	available, frozen, mortgaged := encodeBalance(account)
	msg = UsdtMsg{
		Address:   account.Address,
		Status:    account.Status,
		Symbal:    "USDT",
		Precision: UsdtConfig.Precision,
		Balance: Balance{
			Available: available,
			Frozen:    frozen,
			Mortgaged: mortgaged,
		},
	}

	return
}

func encodeBalance(account *models.UsdtAccount) (available, frozen, mortgaged string) {
	return common.CurrencyInt64ToStr(account.AvailableInteger, UsdtConfig.Precision),
		common.CurrencyInt64ToStr(account.FrozenInteger, UsdtConfig.Precision),
		common.CurrencyInt64ToStr(account.MortgagedInteger, UsdtConfig.Precision)
}

// query balance by uid
func BalanceByUid(uid uint64) (balance Balance, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	account, err := dao.AccountDaoEntity.QueryByUid(uid)
	if err != nil {
		if err == orm.ErrNoRows {
			errCode = ERROR_CODE_USDT_ACCOUNT_NO_FOUND
			return
		}
		errCode = ERROR_CODE_DB
		return
	}

	available, frozen, mortgaged := encodeBalance(account)
	balance = Balance{
		Frozen:    frozen,
		Mortgaged: mortgaged,
		Available: available,
	}

	return
}

func TransferBySystem(uid uint64, logId uint64, toAddr string, amount int64, memo string) (txHash string, errCode ERROR_CODE) {
	return transfer(true, uid, logId, toAddr, amount, memo)
}

func TransferByUser(uid uint64, logId uint64, toAddr string, amount int64, memo string) (txHash string, errCode ERROR_CODE) {
	return transfer(false, uid, logId, toAddr, amount, memo)
}

func transfer(isTransferTySystem bool, uid uint64, logId uint64, toAddr string, amount int64, memo string) (txHash string, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	account, err := dao.AccountDaoEntity.QueryByUid(uid)
	if err != nil {
		errCode = ERROR_CODE_USDT_ACCOUNT_NO_FOUND
		return
	}

	if account.Pkid <= 0 {
		errCode = ERROR_CODE_USDT_PRI_KEY_NO_FOUND
		return
	}

	// 交易数据防篡改检查
	if errCode = dao.WealthLogDaoEntity.VerifyWealthLogSign(logId); errCode != ERROR_CODE_SUCCESS {
		common.LogFuncWarning("usdt transfer tamper protection verify failed : %v", errCode)
		return
	}

	var log *models.UsdtWealthLog

	if log, err = dao.WealthLogDaoEntity.QueryById(logId); err != nil {
		common.LogFuncWarning("%v", err)
		errCode = ERROR_CODE_DB
		if err != orm.ErrNoRows {
			errCode = ERROR_CODE_USDT_REQUEST_FORM_NO_FOUND
		}
		return
	}

	// clear frozen amount
	o := orm.NewOrm()
	o.Begin()
	defer o.Rollback()

	// 系统热钱包， 提现
	var fromUaid uint64 = 1

	// 不是系统转账，从用户钱包转账
	if !isTransferTySystem {
		fromUaid = account.Uaid
	}

	var (
		pkStr, fromAddr string
		l               *lock.Locker
	)

	defer func() {
		if l != nil {
			l.Unlock()
		}
	}()

	if (account.AvailableInteger - account.CashSweepInteger - account.WaitingCashSweepInteger) < amount {
		// 用户账号内不足,需要从平台账号转出,并更新 OwnedByPlatformInteger 和 AvailableInteger
		pkStr, fromAddr, l, err = WalletMgr.GetTransferWallet()
		if err != nil {
			// todo unfrozen usdt
			errCode = ERROR_CODE_USDT_PRI_KEY_NO_FOUND
			dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusOutFailure, dao.OutStepTransferGetPk, fmt.Sprint(err))
			return
		}

		dao.AccountDaoEntity.TransferFromPlatformTx(account.Uid, amount, o)

	} else {
		pkStr, fromAddr, err = priKeyMgr.Get(fromUaid)
		if err != nil {
			if err == orm.ErrNoRows {
				common.LogFuncWarning("priKey of uaid %d no found.", account.Uaid)
			}

			// todo unfrozen usdt
			errCode = ERROR_CODE_USDT_PRI_KEY_NO_FOUND
			dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusOutFailure, dao.OutStepTransferGetPk, fmt.Sprint(err))
			return
		}
	}

	// 手续费更新到待归集和归属平台的金额上
	dao.AccountDaoEntity.TransferServiceChargeTx(account.Uid, log.FeeUsdtInteger, o)

	//generate transaction and broadcast it to mainnet
	amountTx := fmt.Sprintf("%.16x", amount-log.FeeUsdtInteger)

	var signedTx string
	signedTx, err = getOnceSignedTx(pkStr, amountTx, fromAddr, toAddr, log.FeeInteger)
	if err != nil {
		errCode = ERROR_CODE_USDT_GET_SIGNED_TX_FAILED
		dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusOutFailure, dao.OutStepTransferSignTx, fmt.Sprint(err))
		return
	}

	common.LogFuncDebug("%s", signedTx)

	var resp map[string]interface{}
	resp, err = explorer.NewExplorer().TransactionPush(signedTx)
	if err != nil {
		common.LogFuncError("%v", err)
		dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusOutFailure, dao.OutStepTransferTxPushed, fmt.Sprint(err))
		return
	}

	type Msg struct {
		Status string `json:"status"`
		Pushed string `json:"pushed"`
		Tx     string `json:"tx"`
	}

	msg := Msg{}

	err = mapstructure.Decode(resp, &msg)
	if err != nil {
		//信息验证不通过后，发出交易前：恢复提交状态。 (redo again ?)
		dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusOutFailure, dao.OutStepTransferTxPushed, fmt.Sprint(err))
		errCode = ERROR_CODE_USDT_TRANSFER_ON_CHAIN_FAILED
		common.LogFuncError("decode %v to struct failed : %v", resp, err)
		return
	}

	if msg.Status != "OK" {
		//信息验证不通过后，发出交易前：恢复提交状态。 (redo again ?)
		dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusOutFailure, dao.OutStepTransferTxPushed, fmt.Sprintf("push signedTx to %s status : %v", explorer.OMNI_EXPLORER_URL, msg.Status))
		errCode = ERROR_CODE_USDT_TRANSFER_ON_CHAIN_FAILED
		common.LogFuncError("push signedTx to %s status : %v", explorer.OMNI_EXPLORER_URL, msg.Status)
		return
	}

	if msg.Pushed != "success" {
		//信息验证不通过后，发出交易前：恢复提交状态。 (redo again ?)
		dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusOutFailure, dao.OutStepTransferTxPushed, fmt.Sprintf("push signedTx to %s result : %v", explorer.OMNI_EXPLORER_URL, msg.Pushed))
		errCode = ERROR_CODE_USDT_TRANSFER_ON_CHAIN_FAILED
		common.LogFuncError("push signedTx to %s result : %v", explorer.OMNI_EXPLORER_URL, msg.Pushed)
		return
	}

	// commit transaction
	o.Commit()

	txHash = msg.Tx

	// 更新状态日志
	dao.WealthLogDaoEntity.LogOutTransferredOrder(logId, fromAddr, txHash)

	return
}

type BalanceOnChain struct {
	Id         string `json:"id"`
	Pendingpos string `json:"pendingpos"`
	Reserved   string `json:"reserved"`
	Divisible  bool   `json:"divisible"`
	Value      string `json:"value"`
	Frozen     string `json:"frozen"`
	Pendingneg string `json:"pendingneg"`
}

// try to mortgage
func Mortgage(uid uint64, amount string) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
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

	var (
		mortgagedInteger int64
		logId            uint64
	)

	mortgagedInteger, err = common.CurrencyStrToInt64(amount)
	if err != nil {
		errCode = ERROR_CODE_USDT_CURRENCY_PARAM_ERROR
		return
	}

	if mortgagedInteger <= 0 {
		errCode = ERROR_CODE_USDT_CURRENCY_PARAM_ERROR
		return
	}

	logId, err = dao.WealthLogDaoEntity.Add(account.Uid, dao.WealthLogTypeMortgage, dao.WealthLogStatusMortgaging, mortgagedInteger, 0, 0, "", "", "", "")
	if err != nil {
		errCode = ERROR_CODE_USDT_WEALTH_LOG_FAILED
		return
	}

	// try to mortgage
	o := orm.NewOrm()
	o.Begin()
	defer o.Rollback()

	err = dao.AccountDaoEntity.MortgageTx(account.Uid, mortgagedInteger, o)
	if err != nil {
		errCode = ERROR_CODE_USDT_MORTGAGE_FAILED
		dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusMortgageFailed, dao.MortgageStepMortgage, fmt.Sprint(err))
		return
	}

	// 发放eusd
	quant := common.CurrencyInt64ToFloat64(mortgagedInteger)
	errCode = eosplus.EosPlusAPI.Wealth.DelegateUsdt(uid, quant)
	//抵押失败
	if errCode != ERROR_CODE_SUCCESS {
		dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusMortgageFailed, dao.MortgageStepDelegateEusd, otcerror.ErrorMsg(errCode, "zh"))
		return
	}

	o.Commit()

	//抵押成功 添加抵押记录
	dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusMortgaged, "", "")

	return
}

// release usdt amount by PreRelease data
func Release(uid uint64, amount string) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

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

	var (
		releaseInteger int64
		logId          uint64
	)

	releaseInteger, err = common.CurrencyStrToInt64(amount)
	if err != nil {
		errCode = ERROR_CODE_USDT_CURRENCY_PARAM_ERROR
		return
	}

	if releaseInteger <= 0 {
		errCode = ERROR_CODE_USDT_CURRENCY_PARAM_ERROR
		return
	}

	logId, err = dao.WealthLogDaoEntity.Add(uid, dao.WealthLogTypeRelease, dao.WealthLogStatusReleasing, releaseInteger, 0, 0, "", "", "", "")
	if err != nil {
		errCode = ERROR_CODE_USDT_WEALTH_LOG_FAILED
		return
	}

	// try sub mortgage value
	o := orm.NewOrm()
	o.Begin()
	defer o.Rollback()
	err = dao.AccountDaoEntity.ReleaseTx(uid, releaseInteger, o)
	if err != nil {
		errCode = ERROR_CODE_USDT_RELEASE_FAILED
		dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusReleaseFailed, dao.ReleaseStepRelease, fmt.Sprint(err))
		return
	}

	// release eusd
	quant := common.CurrencyInt64ToFloat64(releaseInteger)
	errCode = eosplus.EosPlusAPI.Wealth.UnDelegateUsdt(account.Uid, quant)
	if errCode != ERROR_CODE_SUCCESS {
		dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusReleaseFailed, dao.ReleaseStepUndelegateEusd, otcerror.ErrorMsg(errCode, "zh"))
		return
	}

	// commit transaction
	o.Commit()

	dao.WealthLogDaoEntity.UpdateStatus(logId, dao.WealthLogStatusReleased, "", "")

	return
}

type Balances struct {
	Balance []BalanceOnChain `json:"balance"`
}

//// sync usdt account balance from blockchain
func SyncBalance(uid uint64) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	// check available amount
	account, err := dao.AccountDaoEntity.QueryByUid(uid)
	if err != nil {
		errCode = ERROR_CODE_USDT_ACCOUNT_NO_FOUND
		return
	}

	var resp map[string]interface{}
	resp, err = explorer.NewExplorer().Balances([]string{account.Address})
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	balances := Balances{}
	if v, ok := resp[account.Address]; ok {
		err = mapstructure.Decode(v, &balances)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
	}

	common.LogFuncDebug("resp balance info : %v", balances)

	dbBalances := []models.UsdtOnchainBalance{}
	now := common.NowInt64MS()
	for _, b := range balances.Balance {
		var tmp int64
		tmp, err = strconv.ParseInt(b.Id, 10, 32)
		if err != nil {
			errCode = ERROR_CODE_USDT_CHAIN_BALANCE_ERROR
			return
		}

		dbBalance := models.UsdtOnchainBalance{
			Address:    account.Address,
			PropertyId: uint32(tmp),
			PendingPos: b.Pendingpos,
			PendingNeg: b.Pendingneg,
			Reserved:   b.Reserved,
			Divisible:  b.Divisible,
			Frozen:     b.Frozen,
			Mtime:      now,
		}

		tmp, err = strconv.ParseInt(b.Value, 10, 64)
		if err != nil {
			errCode = ERROR_CODE_USDT_CHAIN_BALANCE_ERROR
			return
		}
		dbBalance.AmountInteger = tmp
		dbBalances = append(dbBalances, dbBalance)
	}

	err = dao.OnchainBalanceDaoEntity.UpdateMulti(uid, dbBalances)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}

	return

}

type Transactions struct {
	Address      string        `json:"address"`
	Pages        int           `json:"pages"`
	CurrentPage  int           `mapstructure:"current_page"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Amount           string `json:"amount"`
	Block            uint32 `json:"block"`
	Blockhash        string `json:"blockhash"`
	Blocktime        int64  `json:"blocktime"`
	Confirmations    uint32 `json:"confirmations"`
	Divisible        bool   `json:"divisible"`
	Fee              string `json:"fee"`
	Ismine           bool   `json:"ismine"`
	Positioninblock  uint32 `json:"positioninblock"`
	Propertyid       uint32 `json:"propertyid"`
	Propertyname     string `json:"propertyname"`
	ReferenceAddress string `json:"referenceaddress"`
	Sendingaddress   string `json:"sendingaddress"`
	Txid             string `json:"txid"`
	Type             string `json:"type"`
	TypeInt          int32  `json:"type_int"`
	Valid            bool   `json:"valid"`
	Version          int32  `json:"version"`
}

const (
	NO_PAGE = 0
)

const (
	FORCE_TO_SYNC    = true
	NO_FORCE_TO_SYNC = !FORCE_TO_SYNC
)

// sync usdt transactions from blockchain
// Notice : the transaction saved to db only once.
func SyncTransaction(account models.UsdtAccount, isForce bool) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	var (
		now = int64(common.NowUint32())
	)

	if !isForce { //if no force . should check sync frequency
		errCode = checkSyncFrequency(account.Address, now)
		if errCode != ERROR_CODE_SUCCESS {
			return
		}
	}

	var newLastestTx string
	// sync transactions from block chain by page.
	var nextPage = 1 //first page
	nextPage, errCode = processPageTransactions(account.Uaid, account.Address, nextPage, &newLastestTx, isForce)
	for {
		if nextPage <= NO_PAGE {
			break
		}
		nextPage, errCode = processPageTransactions(account.Uaid, account.Address, nextPage, &newLastestTx, isForce)
	}

	// if newLastestTx set, just save it.
	if newLastestTx != "" {
		err := dao.OnchainDataDaoEntity.SetLastestTransaction(account.Address, newLastestTx)
		if err != nil {
			errCode = ERROR_CODE_DB
			return
		}

		// set sync timestamp
		err = dao.OnchainDataDaoEntity.SetLastSyncTimestamp(account.Address, now)
		if err != nil {
			errCode = ERROR_CODE_DB
			return
		}
	}

	return
}

// check sync frequency
func checkSyncFrequency(addr string, now int64) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	timestamp, err := dao.OnchainDataDaoEntity.GetLastSyncTimestamp(addr)
	if err != nil {
		if err != orm.ErrNoRows {
			errCode = ERROR_CODE_DB
			return
		}
	}

	// too frequently
	if timestamp+UsdtConfig.SyncFrequency > now {
		errCode = ERROR_CODE_USDT_CHAIN_SYNC_TOO_FREQUENTLY
		return
	}

	return
}

func processPageTransactions(uaid uint64, addr string, page int, newLastestTx *string, isForce bool) (nextPage int, errCode ERROR_CODE) {
	common.LogFuncDebug("start to process page %d", page)
	errCode = ERROR_CODE_SUCCESS
	resp, err := explorer.NewExplorer().TransactionAddress(addr, page)
	if err != nil {
		common.LogFuncError("%v", err)
		errCode = ERROR_CODE_USDT_CHAIN_TRANSACTION_ERROR
		return
	}

	nextPage, errCode = processTransactions(uaid, resp, newLastestTx, isForce)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}
	return
}

func processTransactions(uaid uint64, resp map[string]interface{}, newLastestTx *string, isForce bool) (nextPage int, errCode ERROR_CODE) {
	var (
		dbTxs          []models.UsdtOnchainTransaction
		addr           string
		curPage, pages int
	)

	addr, curPage, pages, dbTxs, errCode = generateDbTransations(uaid, resp, newLastestTx, isForce)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}

	if len(dbTxs) > 0 {
		errCode = saveTransactions2DB(uaid, addr, dbTxs)
		if errCode != ERROR_CODE_SUCCESS {
			return
		}

		//set next page.
		if pages < curPage+1 { //no any more
			nextPage = NO_PAGE
		} else {
			nextPage = curPage + 1 //set to next page
		}
	}

	return
}

func generateDbTransations(uaid uint64, resp map[string]interface{}, newLastestTx *string, isForce bool) (addr string, curPage, pages int, dbTxs []models.UsdtOnchainTransaction, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	txs := Transactions{}
	err := mapstructure.Decode(resp, &txs)
	if err != nil {
		common.LogFuncError("%v", err)
		errCode = ERROR_CODE_USDT_CHAIN_TRANSACTION_ERROR
		return
	}

	addr, curPage, pages = txs.Address, txs.CurrentPage, txs.Pages

	common.LogFuncDebug("%s curr_page %d, pages %d, transactions : %v", txs.Address, curPage, pages, txs)

	now := int64(common.NowUint32())
	var (
		oldLastestTx string
	)

	if !isForce { // if no force, should care of oldLastestTx
		oldLastestTx, err = dao.OnchainDataDaoEntity.GetLastestTransaction(txs.Address)
		if err != nil {
			if err != orm.ErrNoRows {
				errCode = ERROR_CODE_DB
				return
			}
		}
	}

	for _, t := range txs.Transactions {
		//common.LogFuncDebug("%s %v %s -> %s : %s", t.Txid, t.Valid, t.Sendingaddress, t.ReferenceAddress, t.Amount)
		// skip invalid ones
		if !t.Valid {
			continue
		}

		// this means the rests are saved already.
		if oldLastestTx != "" && t.Txid == oldLastestTx {
			break
		}

		// skip the ones whose confirmation no enough.
		//common.LogFuncDebug("confirmations %d, limit %d", t.Confirmations, UsdtConfig.ConfirmationLimit)
		if t.Confirmations < uint32(UsdtConfig.ConfirmationLimit) {
			continue
		}

		dbTx := models.UsdtOnchainTransaction{
			Uaid:             uaid,
			PropertyId:       t.Propertyid,
			PropertyName:     t.Propertyname,
			TxId:             t.Txid,
			TxType:           t.Type,
			TxTypeInt:        t.TypeInt,
			Block:            t.Block,
			BlockHash:        t.Blockhash,
			BlockTime:        t.Blocktime,
			Confirmations:    t.Confirmations,
			Divisible:        t.Divisible,
			IsMine:           t.Ismine,
			PositionInBlock:  t.Positioninblock,
			ReferenceAddress: t.ReferenceAddress,
			SendingAddress:   t.Sendingaddress,
			Version:          t.Version,
			Mtime:            now,
		}
		if t.ReferenceAddress == addr {
			dbTx.Type = dao.TransactionTypeTransferin
		} else {
			dbTx.Type = dao.TransactionTypeTransferout
		}

		dbTx.AmountInteger, err = common.CurrencyStrToInt64(t.Amount)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}

		dbTx.FeeAmountInteger, err = common.CurrencyStrToInt64(t.Fee)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}

		// save the lastest, until all of the transactions save to db
		if *newLastestTx == "" {
			*newLastestTx = t.Txid
		}

		dbTxs = append(dbTxs, dbTx)
	}

	return
}

func saveTransactions2DB(uid uint64, addr string, dbTxs []models.UsdtOnchainTransaction) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	err := dao.OnchainTransactionDaoEntity.UpdateMulti(uid, dbTxs, addr)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}

	return
}
