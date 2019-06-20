package dao

import (
	"common"
)

const (
	InsertMulCount = 50

	//pri redis key
	ReportGameTransfer = "report_game_transfer"
	ReportAgent        = "report_agent"
	ReportEusd         = "report_eusd"

	POSITIVE_BET_LOWEST      = 200
	POSITIVE_VALIDBET_LOWEST = 3000
	POSITIVE_DAY_NUMS        = 7 //判断积极用户要在自然月内登录几天
)

func Init(entityInitFunc common.EntityInitFunc) (err error) {
	const dbOtc = "otc"
	ReportGameUserDailyDaoEntity = NewReportGameUserDailyDao(dbOtc)
	ReportGameRecordRgDaoEntity = NewReportGameRecordRgDao(dbOtc)
	ReportGameRecordAgDaoEntity = NewReportGameRecordAgDao(dbOtc)
	ReportGameRecordKyDaoEntity = NewReportGameRecordKyDao(dbOtc)
	ReportAgentDailyDaoEntity = NewReportAgentDailyDao(dbOtc)
	ReportEusdDailyDaoEntity = NewReportEusdDailyDao(dbOtc)
	ReportGameTransferDailyDaoEntity = NewReportGameTransferDailyDao(dbOtc)
	ReportTeamDailyDaoEntity = NewReportTeamDailyDao(dbOtc)
	ReportTeamGameTransferDailyDaoEntity = NewReportTeamGameTransferDailyDao(dbOtc)
	GameUserMonthReportDaoEntity = NewGameUserMonthReportDao(dbOtc)
	MonthDividendRecordDaoEntity = NewMonthDividendRecordDao(dbOtc)
	ProfitReportDailyDaoEntity = NewProfitReportDailyDao(dbOtc)
	ReportStatisticGameAllDaoEntity = NewReportStatisticGameAllDao(dbOtc)
	ReportStatisticSumDaoEntity = NewReportStatisticSumDao(dbOtc)
	ReportCommissionDaoEntity = NewReportCommissionDao(dbOtc)

	if entityInitFunc != nil {
		err = entityInitFunc()
		if err != nil {
			return
		}
	}
	return
}
