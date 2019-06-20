package dao

import (
	"common"
	"utils/usdt/models"
)

type OnchainLogDao struct {
	common.BaseDao
}

func NewOnchainLogDao(db string) *OnchainLogDao {
	return &OnchainLogDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var OnchainLogDaoEntity *OnchainLogDao

func (d *OnchainLogDao) Add(from, to, tx, status, pushed, signedTx string, amountInteger int64) (err error) {
	var id uint64
	id, err = common.IdManagerGen(IdTypeOnchainLog)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	l := models.UsdtOnchainLog{
		Oclid:         id,
		From:          from,
		To:            to,
		Tx:            tx,
		Status:        status,
		Pushed:        pushed,
		SignedTx:      signedTx,
		AmountInteger: amountInteger,
		Ctime:         common.NowInt64MS(),
	}

	_, err = d.Orm.Insert(&l)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
