package dao

import (
	"common"
	"fmt"
	"utils/otc/models"
)

type OtcSellDao struct {
	common.BaseDao
}

func NewOtcSellDao(db string) *OtcSellDao {
	return &OtcSellDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var OtcSellDaoEntity *OtcSellDao

func (d *OtcSellDao) Edit(data *models.OtcSell) (err error) {
	_, err = d.Orm.InsertOrUpdate(data)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
		return
	}
	return
}

func (d *OtcSellDao) Delete(uid uint64) (err error) {
	data := &models.OtcSell{
		Uid: uid,
	}
	_, err = d.Orm.Delete(data)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
		return
	}
	return
}

func (d *OtcSellDao) Match(pay int8, funds int64) (list []*models.OtcSell) {
	sql := fmt.Sprintf("Select * from %s Where (%s>=? or %s=0) And %s&?>0",
		models.TABLE_OtcSell, models.COLUMN_OtcSell_Available, models.COLUMN_OtcSell_Available, models.COLUMN_OtcSell_PayType)

	list = []*models.OtcSell{}
	_, err := d.Orm.Raw(sql, funds, pay).QueryRows(&list)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
		return
	}
	return
}

func (d *OtcSellDao) EditAvailable(uid uint64, funds int64) (ok bool) {
	//Available=0 表示无限制 无需还原or扣除
	sql := fmt.Sprintf("Update %s set %s=%s-? Where %s=? And %s!=0",
		models.TABLE_OtcSell, models.COLUMN_OtcSell_Available, models.COLUMN_OtcSell_Available, models.COLUMN_OtcSell_PayType, models.COLUMN_OtcSell_Available)

	res, err := d.Orm.Raw(sql, funds, uid).Exec()
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
		return
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return
	}
	return true
}

func (d *OtcSellDao) All() (list []*models.OtcSell, err error) {
	sql := fmt.Sprintf("SELECT * from %s ", models.TABLE_OtcSell)
	list = []*models.OtcSell{}
	_, err = d.Orm.Raw(sql).QueryRows(&list)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
	}
	return
}
