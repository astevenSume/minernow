package dao

import (
	"common"
)

const (
	IdTypeIndexGameTransfer = iota
	IdTypeIndexMax
)
const (
	InsertMulCount = 100
	RECHARGE_TYPE  = iota
	WITHDRAW_TYPE
)

func Init(entityInitFunc common.EntityInitFunc) (err error) {
	//init entity
	const dbOtc = "otc"
	GameUserDailyDaoEntity = NewGameUserDailyDao(dbOtc)
	GameDaoEntity = NewGameDao(dbOtc)
	GameUserDaoEntity = NewGameUserDao(dbOtc)
	GameTransferDaoEntity = NewGameTransferDao(dbOtc)
	GameDailyDaoEntity = NewGameDailyDao(dbOtc)
	GameRiskAlertDaoEntiry = NewGameRiskAlertDao(dbOtc)
	GameOrderRiskDaoEntity = NewGameOrderRiskDao(dbOtc)

	if entityInitFunc != nil {
		err = entityInitFunc()
		if err != nil {
			return
		}
	}

	return
}
