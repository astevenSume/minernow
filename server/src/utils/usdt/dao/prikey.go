package dao

import (
	"common"
	"fmt"
	"utils/usdt/models"

	"github.com/astaxie/beego/orm"
)

type PriKeyDao struct {
	common.BaseDao
}

const (
	HOT_WALLET_MIN_PKID = uint64(10000)
	HOT_WALLET_MAX_PKID = uint64(15000)

	COLD_WALLET_MIN_PKID = uint64(15001)
	COLD_WALLET_MAX_PKID = uint64(20000)
)

func NewPriKeyDao(db string) *PriKeyDao {
	return &PriKeyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var PriKeyDaoEntity *PriKeyDao

// add private key data
func (d *PriKeyDao) Add(id uint64, pri, addr string) (err error) {
	_, err = d.Orm.Insert(&models.PriKey{
		Pkid:    id,
		Pri:     pri,
		Address: addr,
	})

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// get priKey by pkid
func (d *PriKeyDao) QueryByPkid(pkid uint64) (pk models.PriKey, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s=?", models.TABLE_PriKey, models.COLUMN_PriKey_Pkid), pkid).QueryRow(&pk)
	if err != nil {
		if err == orm.ErrNoRows {
			return
		}
		common.LogFuncError("%v", err)
		return
	}

	return
}

// get priKey by addr
func (d *PriKeyDao) QueryByAddr(addr string) (pk models.PriKey, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s=?", models.TABLE_PriKey, models.COLUMN_PriKey_Address), addr).QueryRow(&pk)
	if err != nil {
		if err == orm.ErrNoRows {
			return
		}
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *PriKeyDao) FetchPriKey(min, max uint64) (keys []models.PriKey, err error) {

	_, err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s between ? and ? order by %s desc", models.TABLE_PriKey, models.COLUMN_PriKey_Pkid, models.COLUMN_PriKey_Pkid), min, max).QueryRows(&keys)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
		}
		return
	}

	return
}

func (d *PriKeyDao) AddHotWalletPriKey(pri, addr string) (priKey *models.PriKey, err error) {

	type Result struct {
		MaxPkid uint64 `orm:"column(max_pkid);" json:"max_pkid"`
	}
	var r Result
	err = d.Orm.Raw(fmt.Sprintf("SELECT IFNULL(MAX(%s),%d) as max_pkid FROM %s WHERE %s BETWEEN ? and ? ", models.COLUMN_PriKey_Pkid, HOT_WALLET_MIN_PKID, models.TABLE_PriKey, models.COLUMN_PriKey_Pkid), HOT_WALLET_MIN_PKID, HOT_WALLET_MAX_PKID).QueryRow(&r)
	if err != nil {
		return
	}
	if r.MaxPkid+1 > HOT_WALLET_MAX_PKID {
		err = fmt.Errorf("hot wallet pkid insufficient quantity available")
		return
	}
	priKey = &models.PriKey{
		Pkid:    r.MaxPkid + 1,
		Pri:     pri,
		Address: addr,
	}
	_, err = d.Orm.Insert(priKey)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
