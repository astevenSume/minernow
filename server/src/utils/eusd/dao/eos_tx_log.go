package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/eusd/models"
)

type EosTxLogDao struct {
	common.BaseDao
}

func NewEosTxLogDao(db string) *EosTxLogDao {
	return &EosTxLogDao{
		BaseDao: common.NewBaseDao(db),
	}
}
func (d *EosTxLogDao) Create(transaction *models.EosTxLog) (id int64, err error) {
	id, err = d.Orm.Insert(transaction)
	if err != nil {
		return
	}

	return
}

func (d *EosTxLogDao) Info(id uint64) (data *models.EosTxLog) {
	data = &models.EosTxLog{
		Id: id,
	}
	err := d.Orm.Read(data)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		return
	}

	return
}

const (
	EosTxLogStatusCreated = iota
	EosTxLogStatusTransferring
	EosTxLogStatusTransferred
	EosTxLogStatusFinish
	EosTxLogStatusErr     //转账失败 err
	EosTxLogStatusSignErr //sign err
)

func (d *EosTxLogDao) ResetToCreated(id uint64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=?,%s=%s+1 Where %s=?", models.TABLE_EosTxLog,
		models.COLUMN_EosTxLog_Status, models.COLUMN_EosTxLog_Retry, models.COLUMN_EosTxLog_Retry,
		models.COLUMN_EosTxLog_Id)

	res, err := d.Orm.Raw(sql, EosTxLogStatusCreated, id).Exec()
	if err != nil {
		common.LogFuncError("ResetToCreated DBERR: %v", err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncDebug("ResetToCreated Status Err: %v", id)
		return
	}
	ok = true
	return
}

func (d *EosTxLogDao) UpdateTransferring(id uint64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=?,%s=? Where %s=? And %s=?", models.TABLE_EosTxLog,
		models.COLUMN_EosTxLog_Status, models.COLUMN_EosTxLog_Utime, models.COLUMN_EosTxLog_Status, models.COLUMN_EosTxLog_Id)

	res, err := d.Orm.Raw(sql, EosTxLogStatusTransferring, common.NowInt64MS(), EosTxLogStatusCreated, id).Exec()
	if err != nil {
		common.LogFuncError("EosTxLog transferring DBERR: %v", err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncDebug("EosTxLog transferring Status Err: %v", id)
		return
	}
	ok = true
	return
}

func (d *EosTxLogDao) UpdateTransferred(id uint64, txid uint64) (err error) {
	data := &models.EosTxLog{
		Id:     id,
		Status: EosTxLogStatusTransferred,
		Utime:  common.NowInt64MS(),
		Txid:   txid,
	}
	_, err = d.Orm.Update(data, models.COLUMN_EosTxLog_Status, models.COLUMN_EosTxLog_Txid)
	if err != nil {
		return
	}

	return
}

func (d *EosTxLogDao) UpdateStatusFinish(id uint64) (err error) {
	data := &models.EosTxLog{
		Id:     id,
		Status: EosTxLogStatusFinish,
	}
	_, err = d.Orm.Update(data, models.COLUMN_EosTxLog_Status)
	if err != nil {
		return
	}

	return
}

func (d *EosTxLogDao) UpdateStatusErr(id uint64) (err error) {
	data := &models.EosTxLog{
		Id:     id,
		Status: EosTxLogStatusErr,
	}
	_, err = d.Orm.Update(data, models.COLUMN_EosTxLog_Status, models.COLUMN_EosTxLog_Txid)
	if err != nil {
		return
	}

	return
}

func (d *EosTxLogDao) UpdateStatusSignErr(id uint64) (err error) {
	data := &models.EosTxLog{
		Id:     id,
		Status: EosTxLogStatusSignErr,
	}
	_, err = d.Orm.Update(data, models.COLUMN_EosTxLog_Status, models.COLUMN_EosTxLog_Txid)
	if err != nil {
		return
	}

	return
}

func (d *EosTxLogDao) FetchCheck(limit int64) (list []*models.EosTxLog, err error) {
	qs := d.Orm.QueryTable(models.TABLE_EosTxLog).Filter(models.COLUMN_EosWealthLog_Status, EosTxLogStatusCreated)

	qs = qs.Limit(limit)
	list = []*models.EosTxLog{}

	_, err = qs.All(&list)
	if err != nil {
		common.LogFuncError("DBERR：%v", err)
		return
	}
	return
}

func (d *EosTxLogDao) FetchCheckTransfer(limit int64) (list []*models.EosTxLog, err error) {
	qs := d.Orm.QueryTable(models.TABLE_EosTxLog).Filter(models.COLUMN_EosWealthLog_Status, EosTxLogStatusTransferred)

	qs = qs.Limit(limit)
	list = []*models.EosTxLog{}

	_, err = qs.All(&list)
	if err != nil {
		common.LogFuncError("DBERR：%v", err)
		return
	}
	return
}
