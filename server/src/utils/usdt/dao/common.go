package dao

import (
	"common"
	"fmt"

	"github.com/astaxie/beego"
)

const (
	IdTypeUsdtAccount = iota //usdt account id
	IdTypePriKey             //usdt private key id
	IdTypeOnchainLog         //usdt on chain log id
	IdTypeWealthLog          //usdt mortgage, release record id
	IdTypeSweepLog           //usdt cash sweep record id
	IdTypeMax
)

func Init(entityInitFunc common.EntityInitFunc) (err error) {
	//init id manager
	common.IdMgrInit(IdTypeMax)

	// init on chain data region num
	var tmp int64
	tmp, err = beego.AppConfig.Int64("usdt::onchainDataRegionNum")
	if err != nil {
		common.LogFuncWarning("may be lack of usdt::onchainDataRegionNum.")
		return
	}
	DaoConfig.OnchainDataRegionNum = uint64(tmp)

	// init dao entity
	const (
		dbOtc = "otc"
	)
	AccountDaoEntity = NewAccountDao(dbOtc)
	OnchainBalanceDaoEntity = NewOnchainBalanceDao(dbOtc)
	OnchainDataDaoEntity = NewOnchainDataDao(dbOtc)
	OnchainLogDaoEntity = NewOnchainLogDao(dbOtc)
	OnchainSyncPosDaoEntity = NewOnchainSyncPosDao(dbOtc)
	WealthLogDaoEntity = NewWealthLogDao(dbOtc)
	TransactionDaoEntity = NewTransactionDao(dbOtc)
	OnchainTransactionDaoEntity = NewOnchainTransactionDao(dbOtc)
	SweepLogDaoEntity = NewSweepLogDao(dbOtc)
	PricesDaoEntity = NewPricesDao(dbOtc)
	if entityInitFunc != nil {
		err = entityInitFunc()
		if err != nil {
			return
		}
	}

	return
}

// init for usdt svr.
func Init2(entityInitFunc common.EntityInitFunc) (err error) {
	//init id manager
	common.IdMgrInit(IdTypeMax)

	const (
		dbOtc  = "otc"
		dbUsdt = "usdt"
	)
	// dao of otc database
	AccountDaoEntity = NewAccountDao(dbOtc)
	PricesDaoEntity = NewPricesDao(dbOtc)
	// dao of usdt database
	PriKeyDaoEntity = NewPriKeyDao(dbUsdt)

	if entityInitFunc != nil {
		err = entityInitFunc()
		if err != nil {
			return
		}
	}

	return
}

type OnchainTableMgrInterface interface {
	GetOnChainDataTableName(id uint64, table string) string
}

type OnchainTableManager struct {
}

func NewOnchainTableManager() *OnchainTableManager {
	return &OnchainTableManager{}
}

var OnchainTableMgr OnchainTableMgrInterface = NewOnchainTableManager()

// get the region table name.
func (m *OnchainTableManager) GetOnChainDataTableName(id uint64, table string) string {
	index := common.GetRegionId(id, DaoConfig.OnchainDataRegionNum)
	return table + fmt.Sprintf("_%d", index)
}
