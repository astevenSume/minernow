package dao

import (
	"common"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	. "otc_error"
	otcmodels "utils/otc/models"
	"utils/usdt/models"

	"github.com/astaxie/beego/orm"
)

// log type
const (
	WealthLogTypeUnkown      = iota
	WealthLogTypeTransferIn  //转入
	WealthLogTypeTransferOut //转出 提现
	WealthLogTypeMortgage    //抵押
	WealthLogTypeRelease     //赎回
	WealthLogTypeDeposit     //充值
)

// log status of mortgage & release
const (
	WealthLogStatus               uint32 = iota
	WealthLogStatusMortgaging            // 抵押中
	WealthLogStatusMortgaged             // 已抵押
	WealthLogStatusMortgageFailed        // 抵押失败
	WealthLogStatusReleasing             // 赎回中
	WealthLogStatusReleased              // 已赎回
	WealthLogStatusReleaseFailed         // 赎回失败
)

// log status of deposit
const (
	WealthLogStatusDepositUnknown uint32 = iota + 100
	WealthLogStatusDepositing            //充值中
	WealthLogStatusDeposited             //已充值
	WealthLogStatusDepositFailed         //充值失败
)

// log status of transfer out
const (
	WealthLogStatusOutUnknown uint32 = iota + 200
	WealthLogStatusOutSubmitted
	WealthLogStatusOutCanceled
	WealthLogStatusOutApproved
	WealthLogStatusOutRejected
	WealthLogStatusOutTransferred
	WealthLogStatusOutConfirmed
	WealthLogStatusOutFailure
)

// step of mortgage
const (
	MortgageStepUnkown         = "unkown"
	MortgageStepEncodeCurrency = "EncodeCurrency"
	MortgageStepMortgage       = "Mortgage"
	MortgageStepDelegateEusd   = "DelegateEusd"
	MortgageStepRelease        = "Release"
)

// step of release
const (
	ReleaseStepUnkown         = "unkown"
	ReleaseStepEncodeCurrency = "EncodeCurrency"
	ReleaseStepUndelegateEusd = "UndelegateEusd"
	ReleaseStepRelease        = "Release"
)

// step of transfer
const (
	OutStepTransferFrozen              = "TransferFrozen"
	OutStepTransferCheckEnough         = "CheckEnough"
	OutStepTransferGetPk               = "GetPriKey"
	OutStepTransferSignTx              = "SignTx"
	OutStepTransferTxPushed            = "TxPushed"
	OutStepTransferClearFrozen         = "TransferClearFrozen"
	OutStepTransferUnfrozenWhileCancel = "TransferUnfrozenWhileCancel"
	OutStepTransferUnfrozenWhileReject = "TransferUnfrozenWhileReject"
	OutStepTransferRejected            = "Rejected"
	OutStepTransferTransfering         = "Transfering"
)

const (
	DepositStepTxSynced  = "TxSynced"
	DepositStepTxUpdated = "TxUpdated"
	DepositStepDeposit   = "Deposit"
)

type WealthLogDao struct {
	common.BaseDao
}

func NewWealthLogDao(db string) *WealthLogDao {
	return &WealthLogDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var WealthLogDaoEntity *WealthLogDao

func (d *WealthLogDao) Add(uid uint64, ttype, status uint32, amountInteger, feeInteger, feeUsdtInteger int64, txid, from, to, memo string) (id uint64, err error) {
	id, err = common.IdManagerGen(IdTypeWealthLog)
	if err != nil {
		common.LogFuncError("%v")
		return
	}

	now := common.NowInt64MS()
	data := &models.UsdtWealthLog{
		Id:             id,
		Uid:            uid,
		TType:          ttype,
		Status:         status,
		Txid:           txid,
		AmountInteger:  amountInteger,
		Ctime:          now,
		Utime:          now,
		To:             to,
		Memo:           memo,
		FeeInteger:     feeInteger,
		FeeUsdtInteger: feeUsdtInteger,
	}

	if err = d.signWealthLog(data); err != nil {
		common.LogFuncError("sign wealth log failed : %v", err)
		return
	}

	_, err = d.Orm.Insert(data)
	if err != nil {
		common.LogFuncError("WealthLog DBERR : %v", err)
		return
	}

	return
}

// generate deposit log from onchain transaction
func (d *WealthLogDao) Gen(uid uint64, ttype, status uint32, amountInteger, ctime, mtime int64,
	txid, from, to string) (log models.UsdtWealthLog, err error) {
	var id uint64
	id, err = common.IdManagerGen(IdTypeWealthLog)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	log = models.UsdtWealthLog{
		Id:            id,
		Uid:           uid,
		TType:         ttype,
		Status:        status,
		From:          from,
		To:            to,
		Txid:          txid,
		AmountInteger: amountInteger,
		Ctime:         ctime,
		Utime:         mtime,
	}

	return
}

func (d *WealthLogDao) UpdateStatus(id uint64, status uint32, step, desc string) (err error) {
	data := &models.UsdtWealthLog{
		Id:     id,
		Status: status,
		Step:   step,
		Desc:   desc,
		Utime:  common.NowInt64MS(),
	}
	var num int64
	num, err = d.Orm.Update(data, models.COLUMN_UsdtWealthLog_Status, models.COLUMN_UsdtWealthLog_Utime,
		models.COLUMN_UsdtWealthLog_Step, models.COLUMN_UsdtWealthLog_Desc)
	if err != nil {
		common.LogFuncError("WealthLog DBERR: err id %d , status %d , %v", id, status, err)
		return
	}
	if num <= 0 {
		err = orm.ErrNoRows
		common.LogFuncError("WealthLog DBERR: update status err id %d , status %d", id, status)
		return
	}

	return
}

func (d *WealthLogDao) UpdateStatusWithCheck(id uint64, new, old uint32) (err error) {
	var result sql.Result
	result, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=? AND %s=?",
		models.TABLE_UsdtWealthLog,
		models.COLUMN_UsdtWealthLog_Status,
		models.COLUMN_UsdtWealthLog_Id,
		models.COLUMN_UsdtWealthLog_Status),
		new, id, old).Exec()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var n int64
	if n, err = result.RowsAffected(); err != nil {
		common.LogFuncError("%v", err)
		return
	} else if n != 1 {
		err = common.ErrNoRowEffected
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *WealthLogDao) UpdateAmount(id uint64, integer int64) (err error) {
	data := models.UsdtWealthLog{
		Id:            id,
		AmountInteger: integer,
	}
	_, err = d.Orm.Update(&data, models.COLUMN_UsdtWealthLog_AmountInteger)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// log transferred out order
func (d *WealthLogDao) LogOutTransferredOrder(id uint64, from, txid string) (err error) {
	data := &models.UsdtWealthLog{
		Id:     id,
		Status: WealthLogStatusOutTransferred,
		Txid:   txid,
		From:   from,
		Utime:  common.NowInt64MS(),
	}
	var num int64
	num, err = d.Orm.Update(data,
		models.COLUMN_UsdtWealthLog_Status,
		models.COLUMN_UsdtWealthLog_From,
		models.COLUMN_UsdtWealthLog_Txid,
		models.COLUMN_UsdtWealthLog_Utime)
	if err != nil {
		common.LogFuncError("WealthLog DBERR: err id %d , status %d , txid %s, %v", id, WealthLogStatusOutTransferred, txid, err)
		return
	}
	if num <= 0 {
		err = orm.ErrNoRows
		common.LogFuncError("WealthLog DBERR: err id %d , status %d , txid %s, %v", id, WealthLogStatusOutTransferred, txid, err)
		return
	}

	return
}

func (d *WealthLogDao) UpdateStatusAndTxidAndInner(id uint64, status uint32, from, txid string) (err error) {
	data := &models.UsdtWealthLog{
		Id:     id,
		Status: status,
		Txid:   txid,
		From:   from,
		Utime:  common.NowInt64MS(),
	}
	var num int64
	num, err = d.Orm.Update(data, models.COLUMN_UsdtWealthLog_Status, models.COLUMN_UsdtWealthLog_From, models.COLUMN_UsdtWealthLog_Txid, models.COLUMN_UsdtWealthLog_Utime)
	if err != nil {
		common.LogFuncError("WealthLog DBERR: err id %d , status %d , txid %s, %v", id, status, txid, err)
		return
	}
	if num <= 0 {
		err = orm.ErrNoRows
		common.LogFuncError("WealthLog DBERR: err id %d , status %d , txid %s, %v", id, status, txid, err)
		return
	}

	return
}

//
func (d *WealthLogDao) UpdateStatusAndTxid(id uint64, status uint32, txid string) (err error) {
	data := &models.UsdtWealthLog{
		Id:     id,
		Status: status,
		Txid:   txid,
		Utime:  common.NowInt64MS(),
	}
	var num int64
	num, err = d.Orm.Update(data, models.COLUMN_UsdtWealthLog_Status, models.COLUMN_UsdtWealthLog_Txid, models.COLUMN_UsdtWealthLog_Utime)
	if err != nil {
		common.LogFuncError("WealthLog DBERR: err id %d , status %d , txid %s, %v", id, status, txid, err)
		return
	}
	if num <= 0 {
		err = orm.ErrNoRows
		common.LogFuncError("WealthLog DBERR: err id %d , status %d , txid %s, %v", id, status, txid, err)
		return
	}

	return
}

func (d *WealthLogDao) QueryById(id uint64) (log *models.UsdtWealthLog, err error) {
	log = &models.UsdtWealthLog{
		Id: id,
	}

	err = d.Orm.Read(log)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
		return
	}
	return
}

func (d *WealthLogDao) QueryByIdAndUid(id, uid uint64) (log *models.UsdtWealthLog, err error) {
	log = &models.UsdtWealthLog{
		Id:  id,
		Uid: uid,
	}

	err = d.Orm.Read(log, models.COLUMN_UsdtWealthLog_Id, models.COLUMN_UsdtWealthLog_Uid)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
		return
	}
	return
}

//check if tx exists and transfered
func (d *WealthLogDao) IsExistTransfered(txId string) (isExist bool, err error) {
	l := models.UsdtWealthLog{
		Txid:   txId,
		Status: TransactionStatusTransferred,
	}

	err = d.Orm.Read(&l, models.COLUMN_UsdtWealthLog_Txid, models.COLUMN_UsdtWealthLog_Status)
	if err == nil {
		isExist = true
		return
	} else if err == orm.ErrNoRows {
		err = nil
		return
	} else {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *WealthLogDao) Count(uid uint64, status uint32, types []interface{}) (num int64) {
	qs := d.Orm.QueryTable(models.TABLE_UsdtWealthLog)
	if uid > 0 {
		qs = qs.Filter(models.COLUMN_UsdtWealthLog_Uid, uid)
	}
	if status > 0 {
		qs = qs.Filter(models.COLUMN_UsdtWealthLog_Status, status)
	}
	if len(types) > 0 {
		qs = qs.Filter(models.COLUMN_UsdtWealthLog_TType+"__in", types...)
	}
	num, err := qs.Count()

	if err != nil {
		common.LogFuncError("DBERR：%v", err)
		return
	}
	return
}

func (d *WealthLogDao) Fetch(uid uint64, types []interface{}, limit int64, offset int64) (list []*models.UsdtWealthLog, err error) {
	qs := d.Orm.QueryTable(models.TABLE_UsdtWealthLog)
	if uid > 0 {
		qs = qs.Filter(models.COLUMN_UsdtWealthLog_Uid, uid)
	}
	if len(types) > 0 {
		qs = qs.Filter(models.COLUMN_UsdtWealthLog_TType+"__in", types...)
	}

	qs = qs.Limit(limit).Offset(offset).OrderBy("-" + models.COLUMN_UsdtWealthLog_Id)

	list = []*models.UsdtWealthLog{}

	n, err := qs.All(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR：%v", err)

		return
	}
	common.LogFuncError("%v", n)

	return
}

// dump by status and type
func (d *WealthLogDao) DumpByTypeAndStatus(ttype uint32, status uint32) (list []*models.UsdtWealthLog, err error) {
	qs := d.Orm.QueryTable(models.TABLE_UsdtWealthLog).
		Filter(models.COLUMN_UsdtWealthLog_Status, status).
		Filter(models.COLUMN_UsdtWealthLog_TType, ttype)

	list = []*models.UsdtWealthLog{}
	n, err := qs.All(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR：%v", err)
		return
	}
	common.LogFuncDebug("%v", n)

	return
}

func (d *WealthLogDao) CountByStatus(uid uint64, status uint32) (num int64) {
	qs := d.Orm.QueryTable(models.TABLE_UsdtWealthLog).Filter(models.COLUMN_UsdtWealthLog_Status, status)
	if uid > 0 {
		qs = qs.Filter(models.COLUMN_UsdtWealthLog_Uid, uid)
	}
	num, err := qs.Count()
	if err != nil {
		common.LogFuncError("DBERR：%v", err)
		return
	}
	return
}

//auto_models_start
type DetailUsdtWealthLog struct {
	Id                string `orm:"column(id);pk" json:"id,omitempty"`
	Uid               string `orm:"column(uid)" json:"uid,omitempty"`
	TType             uint32 `orm:"column(ttype)" json:"ttype,omitempty"`
	Status            uint32 `orm:"column(status)" json:"status,omitempty"`
	From              string `orm:"column(from);size(256)" json:"from,omitempty"`
	To                string `orm:"column(to);size(256)" json:"to,omitempty"`
	Txid              string `orm:"column(txid);size(256)" json:"txid,omitempty"`
	AmountInteger     int64  `orm:"column(amount_integer)" json:"amount_integer,omitempty"`
	FeeInteger        int64  `orm:"column(fee_integer)" json:"fee_integer,omitempty"`
	FeeOnchainInteger int64  `orm:"column(fee_onchain_integer)" json:"fee_onchain_integer,omitempty"`
	Ctime             int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime             int64  `orm:"column(utime)" json:"utime,omitempty"`
	Step              string `orm:"column(step);size(64)" json:"step,omitempty"`
	Desc              string `orm:"column(desc);size(256)" json:"desc,omitempty"`
	Address           string `orm:"column(desc);size(256)" json:"desc,omitempty"`
	Mobile            string `json:"mobile"`
}

func (d *WealthLogDao) FetchByStatus(uid uint64, status uint32, limit int64, offset int64) (list []*DetailUsdtWealthLog, err error) {
	var sql string
	if uid > 0 {
		sql = fmt.Sprintf("SELECT T1.*,T2.%s,T3.%s FROM ((SELECT * FROM %s WHERE %s=? and %s=? ORDER BY %s"+
			" ASC LIMIT ? OFFSET ?) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s LEFT JOIN %s AS T3 ON T1.%s=T3.%s)",
			otcmodels.COLUMN_User_Mobile, models.COLUMN_UsdtAccount_Address, models.TABLE_UsdtWealthLog,
			models.COLUMN_UsdtWealthLog_Status, models.COLUMN_UsdtWealthLog_Uid, models.COLUMN_UsdtWealthLog_Id,
			otcmodels.TABLE_User, models.COLUMN_UsdtWealthLog_Uid, otcmodels.COLUMN_User_Uid, models.TABLE_UsdtAccount,
			models.COLUMN_UsdtWealthLog_Uid, models.COLUMN_UsdtAccount_Uid)
		_, err = d.Orm.Raw(sql, status, uid, limit, offset).QueryRows(&list)
	} else {
		sql = fmt.Sprintf("SELECT T1.*,T2.%s,T3.%s FROM ((SELECT * FROM %s WHERE %s=? ORDER BY %s ASC LIMIT ? "+
			"OFFSET ?) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s LEFT JOIN %s AS T3 ON T1.%s=T3.%s)",
			otcmodels.COLUMN_User_Mobile, models.COLUMN_UsdtAccount_Address, models.TABLE_UsdtWealthLog,
			models.COLUMN_UsdtWealthLog_Status, models.COLUMN_UsdtWealthLog_Id, otcmodels.TABLE_User,
			models.COLUMN_UsdtWealthLog_Uid, otcmodels.COLUMN_User_Uid, models.TABLE_UsdtAccount,
			models.COLUMN_UsdtWealthLog_Uid, models.COLUMN_UsdtAccount_Uid)
		_, err = d.Orm.Raw(sql, status, limit, offset).QueryRows(&list)
	}
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *WealthLogDao) DetailFetch(uid uint64, status uint32, types []string, limit int, offset int) (list []*DetailUsdtWealthLog, err error) {
	var qbQuery orm.QueryBuilder
	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	// 构建查询对象
	qbQuery.Select("T1.*",
		"T2."+otcmodels.COLUMN_User_Mobile,
		"T3."+models.COLUMN_UsdtAccount_Address).
		From("((").Select("*").From(models.TABLE_UsdtWealthLog).Where("1=1")
	var param []interface{}
	if uid > 0 {
		qbQuery.And(models.COLUMN_UsdtWealthLog_Uid + " = ?")
		param = append(param, uid)
	}
	if status > 0 {
		qbQuery.And(models.COLUMN_UsdtWealthLog_Status + " = ?")
		param = append(param, status)
	}
	if len(types) > 0 {
		qbQuery.And(models.COLUMN_UsdtWealthLog_TType).In(types...)
		//param = append(param, types)
	}

	qbQuery.OrderBy("-" + models.COLUMN_UsdtWealthLog_Id).Limit(limit).Offset(offset)
	sqlQuery := qbQuery.String()
	sqlQuery = fmt.Sprintf("%s) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s LEFT JOIN %s AS T3 ON T1.%s=T3.%s)",
		sqlQuery, otcmodels.TABLE_User, models.COLUMN_UsdtWealthLog_Uid, otcmodels.COLUMN_User_Uid,
		models.TABLE_UsdtAccount, models.COLUMN_UsdtWealthLog_Uid, models.COLUMN_UsdtAccount_Uid)
	_, err = d.Orm.Raw(sqlQuery, param...).QueryRows(&list)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *WealthLogDao) QueryApprovedTotal() (total int, err error) {
	var tmp int64
	tmp, err = d.Orm.QueryTable(models.TABLE_UsdtWealthLog).Filter(models.COLUMN_UsdtWealthLog_Status, WealthLogStatusOutApproved).Count()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	total = int(tmp)

	return
}

func (d *WealthLogDao) QueryApproved(page, perPage int) (list []models.UsdtWealthLog, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_UsdtWealthLog).Filter(models.COLUMN_UsdtWealthLog_Status, WealthLogStatusOutApproved).Offset((page - 1) * perPage).Limit(perPage).All(&list)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
	}
	return
}

func (d *WealthLogDao) InsertMulti(list []models.UsdtWealthLog) (err error) {
	_, err = d.Orm.InsertMulti(common.BulkCount, list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

func (d *WealthLogDao) UpdateByTxID(TxID, from, to string, ttype, status uint32, amount, fee, ctime int64) (err error) {
	_, err = d.Orm.
		QueryTable(models.TABLE_UsdtWealthLog).
		Filter(models.COLUMN_UsdtWealthLog_Txid, TxID).
		Update(orm.Params{
			models.COLUMN_UsdtWealthLog_TType:             ttype,
			models.COLUMN_UsdtWealthLog_Status:            status,
			models.COLUMN_UsdtWealthLog_From:              from,
			models.COLUMN_UsdtWealthLog_To:                to,
			models.COLUMN_UsdtWealthLog_AmountInteger:     amount,
			models.COLUMN_UsdtWealthLog_FeeOnchainInteger: fee,
			models.COLUMN_UsdtWealthLog_Ctime:             ctime,
			models.COLUMN_UsdtWealthLog_Utime:             common.NowUint32(),
		})

	return
}

const (
	WEALTH_LOG_SIGN_SALT = "WEALTH_LOG_SIGN_SALT"
)

func (d *WealthLogDao) getSign(log *models.UsdtWealthLog) (sign string, err error) {

	hash := md5.New()
	if _, err = hash.Write([]byte(fmt.Sprintf("%v%v%v%v%v%v%v%v", log.Uid, log.TType, log.To, log.AmountInteger, log.FeeInteger, log.FeeUsdtInteger, log.Ctime, WEALTH_LOG_SIGN_SALT))); err != nil {
		common.LogFuncError("%v", err)
		return
	}
	sign = hex.EncodeToString(hash.Sum(nil))
	return
}
func (d *WealthLogDao) signWealthLog(log *models.UsdtWealthLog) (err error) {
	log.Sign, err = d.getSign(log)
	return
}

func (d *WealthLogDao) VerifyWealthLogSign(logID uint64) (errCode ERROR_CODE) {
	var (
		err  error
		log  *models.UsdtWealthLog
		sign string
	)

	errCode = ERROR_CODE_SUCCESS

	if log, err = d.QueryById(logID); err != nil {
		if err == orm.ErrNoRows {
			errCode = ERROR_CODE_USDT_REQUEST_FORM_NO_FOUND
			return
		}
		errCode = ERROR_CODE_DB
	}

	if sign, err = d.getSign(log); err != nil {
		errCode = ERROR_CODE_USDT_TRANSFER_TAMPER_PROTECTION_VERIFY_FAILED
		return
	}

	if sign != log.Sign {
		errCode = ERROR_CODE_USDT_TRANSFER_TAMPERED
		return
	}

	return
}

//统计充值提现
func (d *WealthLogDao) Stat(start, end int64) (res map[string]uint32) {
	sql := "select `ttype`,sum(amount_integer) as amount, sum(fee_onchain_integer) as fee from usdt_wealth_log where ctime>? and ctime<? group by `ttype`"
	type stat struct {
		Type   uint8  `json:"ttype"`
		Amount uint32 `json:"amount"`
		Fee    uint32 `json:"fee"`
	}

	list := []*stat{}
	_, err := d.Orm.Raw(sql, start, end).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	res = map[string]uint32{
		"usdt_recharge":   0,
		"usdt_withdrawal": 0,
		"usdt_fee":        0,
	}

	for _, v := range list {
		if v.Type == WealthLogTypeTransferIn {
			res["usdt_recharge"] = v.Amount
		} else if v.Type == WealthLogTypeTransferOut {
			res["usdt_withdrawal"] = v.Amount
			res["usdt_fee"] = v.Fee
		}
	}

	return
}

//统计用户充值提现
func (d *WealthLogDao) StatPeople(uid uint64) (res map[string]uint32) {
	sql := "select `ttype`,sum(amount_integer) as amount from usdt_wealth_log where uid=? group by `ttype`"
	type stat struct {
		Type   uint8  `json:"ttype"`
		Amount uint32 `json:"amount"`
		Fee    uint32 `json:"fee"`
	}

	list := []*stat{}
	_, err := d.Orm.Raw(sql, uid).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	res = map[string]uint32{
		"usdt_recharge":   0,
		"usdt_withdrawal": 0,
	}

	for _, v := range list {
		if v.Type == WealthLogTypeTransferIn {
			res["usdt_recharge"] = v.Amount
		} else if v.Type == WealthLogTypeTransferOut {
			res["usdt_withdrawal"] = v.Amount
		}
	}

	return
}

//统计用户充值提现
func (d *WealthLogDao) StatPeopleByDateRange(uid uint64, startTime, endTime int64) (res map[string]uint32) {
	querySql := fmt.Sprintf("select `%s`,sum(%s) as amount from %s where %s=? and %s>=? and %s<=> group by `%s`",
		models.COLUMN_UsdtWealthLog_TType, models.COLUMN_UsdtWealthLog_AmountInteger, models.TABLE_UsdtWealthLog,
		models.COLUMN_UsdtWealthLog_Uid, models.COLUMN_UsdtWealthLog_Ctime, models.COLUMN_UsdtWealthLog_Ctime, models.COLUMN_UsdtWealthLog_TType)
	type stat struct {
		Type   uint8  `json:"ttype"`
		Amount uint32 `json:"amount"`
	}

	list := []*stat{}
	_, err := d.Orm.Raw(querySql, uid, startTime, endTime).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	res = map[string]uint32{
		"usdt_recharge":   0,
		"usdt_withdrawal": 0,
	}

	for _, v := range list {
		if v.Type == WealthLogTypeTransferIn {
			res["usdt_recharge"] = v.Amount
		} else if v.Type == WealthLogTypeTransferOut {
			res["usdt_withdrawal"] = v.Amount
		}
	}

	return
}
