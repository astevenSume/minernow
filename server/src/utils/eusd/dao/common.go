package dao

import (
	"common"
	"encoding/json"
)

var AccountDaoEntity *AccountDao
var EosUseLogDaoEntity *EosUseLogDao
var EosOtcDaoEntity *EosOtcDao
var TransactionDaoEntity *TransactionDao
var TransactionInfoDaoEntity *TransactionInfoDao
var WealthDaoEntity *WealthDao
var WealthLogDaoEntity *WealthLogDao
var EosTxLogDaoEntity *EosTxLogDao
var EusdRetireDaoEntity *EusdRetireDao
var PlatformUserDaoEntity *PlatformUserDao
var PlatformUserCateDaoEntity *PlatformUserCateDao

func Init(entityInitFunc common.EntityInitFunc) (err error) {

	const dbOtc = "otc"
	AccountDaoEntity = NewAccountDao(dbOtc)
	EosUseLogDaoEntity = NewEosUseLogDao(dbOtc)
	EosOtcDaoEntity = NewEosOtcDao(dbOtc)
	TransactionDaoEntity = NewTransactionDao(dbOtc)
	TransactionInfoDaoEntity = NewTransactionInfoDao(dbOtc)
	WealthDaoEntity = NewWealthDao(dbOtc)
	WealthLogDaoEntity = NewWealthLogDao(dbOtc)
	EosTxLogDaoEntity = NewEosTxLogDao(dbOtc)
	EusdRetireDaoEntity = NewEusdRetireDao(dbOtc)
	PlatformUserDaoEntity = NewPlatformUserDao(dbOtc)
	PlatformUserCateDaoEntity = NewPlatformUserCateDao(dbOtc)

	if entityInitFunc != nil {
		err = entityInitFunc()
		if err != nil {
			return
		}
	}

	return
}

func ToJson(arg interface{}) string {
	data, _ := json.Marshal(arg)

	return string(data)
}

func ToJsonIndent(arg interface{}) string {
	data, _ := json.MarshalIndent(arg, "", "  ")

	return string(data)
}
