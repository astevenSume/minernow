package dao

import (
	"common"
	adminmodels "utils/admin/models"
)

type IpWhiteListDao struct {
	common.BaseDao
}

func NewIpWhiteListDao(db string) *IpWhiteListDao {
	return &IpWhiteListDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var IpWhiteListDaoEntity *IpWhiteListDao

func (d *IpWhiteListDao) All() (list []adminmodels.IpWhiteList, err error) {
	_, err = d.Orm.QueryTable(adminmodels.TABLE_IpWhiteList).All(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
