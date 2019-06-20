package dao

import (
	"common"
	"utils/usdt/models"
)

type OnchainSyncPosDao struct {
	common.BaseDao
}

func NewOnchainSyncPosDao(db string) *OnchainSyncPosDao {
	return &OnchainSyncPosDao{
		BaseDao: common.NewBaseDao(db),
	}
}

type OnchainSyncPosDaoInterface interface {
	SetLastestTransaction(addr, tx string, page uint32) (err error)
	GetLastestTransaction(addr string) (page uint32, tx string, err error)
}

var OnchainSyncPosDaoEntity OnchainSyncPosDaoInterface

func (d *OnchainSyncPosDao) SetLastestTransaction(addr, tx string, page uint32) (err error) {
	data := models.UsdtOnChainSyncPos{
		Address: addr,
		Page:    page,
		TxId:    tx,
	}
	_, err = d.Orm.InsertOrUpdate(&data)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *OnchainSyncPosDao) GetLastestTransaction(addr string) (page uint32, tx string, err error) {
	data := models.UsdtOnChainSyncPos{
		Address: addr,
	}

	err = d.Orm.Read(&data)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	page, tx = data.Page, data.TxId

	return
}
