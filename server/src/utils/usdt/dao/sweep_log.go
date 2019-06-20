package dao

import (
	"common"
	"utils/usdt/models"

	"github.com/astaxie/beego/orm"
)

const (
	SweepLogType           = iota
	SweepLogTypeToPlatform //转入平台
)

const (
	SweepLogStatus            = iota //转账中
	SweepLogStatusTransferred        //链上转账成功
	SweepLogStatusFailure            //失败

)

type SweepStep string

// step of sweep
const (
	SweepStepGetPriKey        = SweepStep("GetPriKey")
	SweepStepPreSignTx        = SweepStep("PreSignTx")
	SweepStepSignedTx         = SweepStep("SignedTx")
	SweepStepTransferTxPushed = SweepStep("TxPushed")
	// SweepStep = "CheckEnough"
	// SweepStep = "GetPriKey"
	// SweepStep = "SignTx"
	// SweepStep = "TxPushed"
	// SweepStep = "TransferClearFrozen"
	// SweepStep = "TransferUnfrozenWhileCancel"
	// SweepStep = "TransferUnfrozenWhileReject"
	// SweepStep = "Rejected"
	// SweepStep = "Transfering"
)

type SweepLogDao struct {
	common.BaseDao
}

func NewSweepLogDao(db string) *SweepLogDao {
	return &SweepLogDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var SweepLogDaoEntity *SweepLogDao

func (d *SweepLogDao) Add(uid uint64, ttype, status uint32, amountInteger int64, txid, from, to string) (id uint64, err error) {
	id, err = common.IdManagerGen(IdTypeSweepLog)
	if err != nil {
		common.LogFuncError("%v")
		return
	}

	now := common.NowInt64MS()
	data := &models.UsdtSweepLog{
		Id:            id,
		Uid:           uid,
		TType:         ttype,
		Status:        status,
		Txid:          txid,
		AmountInteger: amountInteger,
		Ctime:         now,
		Utime:         now,
		From:          from,
		To:            to,
	}
	_, err = d.Orm.Insert(data)
	if err != nil {
		common.LogFuncError("SweepLog DBERR : %v", err)
		return
	}

	return
}

func (d *SweepLogDao) UpdateStatus(id uint64, status uint32, step SweepStep, desc string) (err error) {
	data := &models.UsdtSweepLog{
		Id:     id,
		Status: status,
		Step:   string(step),
		Desc:   desc,
		Utime:  common.NowInt64MS(),
	}
	var num int64
	num, err = d.Orm.Update(data, models.COLUMN_UsdtSweepLog_Status, models.COLUMN_UsdtSweepLog_Utime,
		models.COLUMN_UsdtSweepLog_Step, models.COLUMN_UsdtSweepLog_Desc)
	if err != nil {
		common.LogFuncError("SweepLog DBERR: err id %d , status %d , %v", id, status, err)
		return
	}
	if num <= 0 {
		err = orm.ErrNoRows
		common.LogFuncError("SweepLog DBERR: update status err id %d , status %d", id, status)
		return
	}

	return
}

func (d *SweepLogDao) UpdateStatusAndTxid(id uint64, status uint32, txid string) (err error) {
	data := &models.UsdtSweepLog{
		Id:     id,
		Status: status,
		Txid:   txid,
		Utime:  common.NowInt64MS(),
	}
	var num int64
	num, err = d.Orm.Update(data, models.COLUMN_UsdtSweepLog_Status, models.COLUMN_UsdtSweepLog_Txid, models.COLUMN_UsdtSweepLog_Utime)
	if err != nil {
		common.LogFuncError("SweepLog DBERR: err id %d , status %d , txid %s, %v", id, status, txid, err)
		return
	}
	if num <= 0 {
		err = orm.ErrNoRows
		common.LogFuncError("SweepLog DBERR: err id %d , status %d , txid %s, %v", id, status, txid, err)
		return
	}

	return
}

func (d *SweepLogDao) IsExistTransferred(txId string) (isExist bool, err error) {
	l := models.UsdtSweepLog{
		Txid:   txId,
		Status: SweepLogStatusTransferred,
	}

	err = d.Orm.Read(&l, models.COLUMN_UsdtSweepLog_Txid, models.COLUMN_UsdtSweepLog_Status)
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

func (d *SweepLogDao) UpdateByTxID(TxID, from, to string, amount, fee, ctime int64) (err error) {
	_, err = d.Orm.
		QueryTable(models.TABLE_UsdtSweepLog).
		Filter(models.COLUMN_UsdtSweepLog_Txid, TxID).
		Update(orm.Params{
			models.COLUMN_UsdtSweepLog_From:              from,
			models.COLUMN_UsdtSweepLog_To:                to,
			models.COLUMN_UsdtSweepLog_AmountInteger:     amount,
			models.COLUMN_UsdtSweepLog_FeeOnchainInteger: fee,
			models.COLUMN_UsdtSweepLog_Ctime:             ctime,
			models.COLUMN_UsdtSweepLog_Utime:             common.NowUint32(),
		})

	return
}
