package dao

import (
	"common"
	"fmt"
	"utils/otc/models"
)

type OtcBuyDao struct {
	common.BaseDao
}

func NewOtcBuyDao(db string) *OtcBuyDao {
	return &OtcBuyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var OtcBuyDaoEntity *OtcBuyDao

func (d *OtcBuyDao) Edit(data *models.OtcBuy) (err error) {
	_, err = d.Orm.InsertOrUpdate(data, models.COLUMN_OtcBuy_Available, models.COLUMN_OtcBuy_LowerLimitWechat,
		models.COLUMN_OtcBuy_LowerLimitWechat, models.COLUMN_OtcBuy_UpperLimitWechat,
		models.COLUMN_OtcBuy_LowerLimitBank, models.COLUMN_OtcBuy_UpperLimitBank,
		models.COLUMN_OtcBuy_LowerLimitAli, models.COLUMN_OtcBuy_UpperLimitAli,
		models.COLUMN_OtcBuy_PayType,
	)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
		return
	}
	return
}

func (d *OtcBuyDao) EditAvailable(uid uint64, av int64) (err error) {
	sql := fmt.Sprintf("update %s set %s=%s-? where %s=?", models.TABLE_OtcBuy,
		models.COLUMN_OtcBuy_Available, models.COLUMN_OtcBuy_Available,
		models.COLUMN_OtcBuy_Uid)
	_, err = d.Orm.Raw(sql, av, uid).Exec()
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
		return
	}
	return
}

func (d *OtcBuyDao) DeleteByUid(uid uint64) (err error) {
	data := &models.OtcBuy{
		Uid: uid,
	}
	_, err = d.Orm.Delete(data)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
		return
	}
	return
}

func (d *OtcBuyDao) MatchWechat(quantity int64) (list []*models.OtcBuy, err error) {
	sql := fmt.Sprintf("SELECT * from %s WHERE %s<=? And %s>=? AND %s>=?", models.TABLE_OtcBuy,
		models.COLUMN_OtcBuy_LowerLimitWechat, models.COLUMN_OtcBuy_UpperLimitWechat, models.COLUMN_OtcBuy_Available)
	list = []*models.OtcBuy{}
	_, err = d.Orm.Raw(sql, quantity, quantity, quantity).QueryRows(&list)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
	}
	return
}

func (d *OtcBuyDao) MatchAli(quantity int64) (list []*models.OtcBuy, err error) {
	sql := fmt.Sprintf("SELECT * from %s WHERE %s<=? And %s>=? AND %s>=?", models.TABLE_OtcBuy,
		models.COLUMN_OtcBuy_LowerLimitAli, models.COLUMN_OtcBuy_UpperLimitAli, models.COLUMN_OtcBuy_Available)
	list = []*models.OtcBuy{}
	_, err = d.Orm.Raw(sql, quantity, quantity, quantity).QueryRows(&list)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
	}
	return
}

func (d *OtcBuyDao) MatchBank(quantity int64) (list []*models.OtcBuy, err error) {
	sql := fmt.Sprintf("SELECT * from %s WHERE %s<=? And %s>=? AND %s>=? ", models.TABLE_OtcBuy,
		models.COLUMN_OtcBuy_LowerLimitBank, models.COLUMN_OtcBuy_UpperLimitBank, models.COLUMN_OtcBuy_Available)
	list = []*models.OtcBuy{}
	_, err = d.Orm.Raw(sql, quantity, quantity, quantity).QueryRows(&list)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
	}
	return
}

func (d *OtcBuyDao) All() (list []*models.OtcBuy, err error) {
	sql := fmt.Sprintf("SELECT * from %s ", models.TABLE_OtcBuy)
	list = []*models.OtcBuy{}
	_, err = d.Orm.Raw(sql).QueryRows(&list)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
	}
	return
}
