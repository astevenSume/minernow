package dao

import (
	"common"
	"utils/eusd/models"
)

type PlatformUserCateDao struct {
	common.BaseDao
}

func NewPlatformUserCateDao(db string) *PlatformUserCateDao {
	return &PlatformUserCateDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *PlatformUserCateDao) All() (list []*models.PlatformUserCate, err error) {
	list = []*models.PlatformUserCate{}

	_, err = d.Orm.Raw("SELECT * from platform_user_cate").QueryRows(&list)
	if err != nil {
		if err != nil {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

func (d *PlatformUserCateDao) Add(name string, dividend int32) (data *models.PlatformUserCate, err error) {
	data = &models.PlatformUserCate{
		Name:     name,
		Dividend: dividend,
		Ctime:    common.NowUint32(),
	}

	_, err = d.Orm.Insert(data)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}
