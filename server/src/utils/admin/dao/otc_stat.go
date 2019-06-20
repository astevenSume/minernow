package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/admin/models"
)

type OtcStatDao struct {
	common.BaseDao
}

var OtcStatDaoEntity *OtcStatDao

func NewOtcStatDao(db string) *OtcStatDao {
	return &OtcStatDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *OtcStatDao) Create(data *models.OtcStat) {
	if _, err := d.Orm.Insert(data); err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
}

func (d *OtcStatDao) Count() (num int64) {
	qs := d.Orm.QueryTable(models.TABLE_OtcStat)
	num, err := qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}
	return
}

func (d *OtcStatDao) Fetch(offset int64, limit int64) (list []*models.OtcStat) {
	qs := d.Orm.QueryTable(models.TABLE_OtcStat)

	list = []*models.OtcStat{}
	qs = qs.Limit(limit, offset)
	_, err := qs.All(&list)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	return
}
