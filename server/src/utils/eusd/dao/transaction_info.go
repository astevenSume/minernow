package dao

import (
	"common"
	"utils/eusd/models"
)

type TransactionInfoDao struct {
	common.BaseDao
}

func NewTransactionInfoDao(db string) *TransactionInfoDao {
	return &TransactionInfoDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *TransactionInfoDao) Create(transaction *models.EosTransactionInfo) (id int64, err error) {
	id, err = d.Orm.Insert(transaction)
	if err != nil {
		return
	}

	return
}
