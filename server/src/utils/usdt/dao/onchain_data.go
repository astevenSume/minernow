package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/usdt/models"
)

type OnchainDataDaoInterface interface {
	SetLastestTransaction(addr, tx string) (err error)
	GetLastestTransaction(addr string) (tx string, err error)
	SetLastSyncTimestamp(addr string, timestamp int64) (err error)
	GetLastSyncTimestamp(addr string) (timestamp int64, err error)
}

type OnchainDataDao struct {
	common.BaseDao
}

func NewOnchainDataDao(db string) *OnchainDataDao {
	return &OnchainDataDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var OnchainDataDaoEntity OnchainDataDaoInterface

const (
	AttrTypeUnkown uint32 = iota
	AttrTypeLastestTransation
	AttrTypeLastSyncTimestamp
)

func (d *OnchainDataDao) SetLastestTransaction(addr, tx string) (err error) {
	data := models.UsdtOnChainData{
		Address:  addr,
		AttrType: AttrTypeLastestTransation,
		DataStr:  tx,
	}
	_, err = d.Orm.InsertOrUpdate(&data)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *OnchainDataDao) GetLastestTransaction(addr string) (tx string, err error) {
	data := models.UsdtOnChainData{
		Address:  addr,
		AttrType: AttrTypeLastestTransation,
	}

	err = d.Orm.Read(&data)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
		}
		return
	}

	tx = data.DataStr

	return
}

func (d *OnchainDataDao) SetLastSyncTimestamp(addr string, timestamp int64) (err error) {
	data := models.UsdtOnChainData{
		Address:   addr,
		AttrType:  AttrTypeLastSyncTimestamp,
		DataInt64: timestamp,
	}
	_, err = d.Orm.InsertOrUpdate(&data)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *OnchainDataDao) GetLastSyncTimestamp(addr string) (timestamp int64, err error) {
	data := models.UsdtOnChainData{
		Address:  addr,
		AttrType: AttrTypeLastSyncTimestamp,
	}

	err = d.Orm.Read(&data)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
		}
		return
	}

	timestamp = data.DataInt64

	return
}
