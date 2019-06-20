package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/otc/models"
)

type UserConfigDao struct {
	common.BaseDao
}

func NewUserConfigDao(db string) *UserConfigDao {
	return &UserConfigDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var UserConfigDaoEntity *UserConfigDao

func (d *UserConfigDao) Info(uid uint64) (data *models.UserConfig) {
	data = &models.UserConfig{
		Uid: uid,
	}

	err := d.Orm.Read(data)
	if err == orm.ErrNoRows {
		data.WealthNotice = true
		data.OrderNotice = true
		return
	}
	if err != nil {

		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

func (d *UserConfigDao) Update(uid uint64, wealthNotice, orderNotice bool) (ok bool) {
	data := &models.UserConfig{
		Uid:          uid,
		WealthNotice: wealthNotice,
		OrderNotice:  orderNotice,
	}

	_, err := d.Orm.InsertOrUpdate(data)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}

	ok = true
	return
}
