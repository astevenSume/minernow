package dao

import (
	"common"
	"utils/eusd/models"
)

type EosUseLogDao struct {
	common.BaseDao
}

func NewEosUseLogDao(db string) *EosUseLogDao {
	return &EosUseLogDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *EosUseLogDao) Create(data *models.EosUseLog) (id int64, err error) {
	id, err = d.Orm.Insert(data)
	if err != nil {
		return
	}

	return
}

func (d *EosUseLogDao) InsertMulti(data []*models.EosUseLog) (id int64, err error) {
	id, err = d.Orm.InsertMulti(50, data)
	if err != nil {
		return
	}

	return
}
