package cron

import "common"

const (
	DailyQueryPerPage = 50
)

// @Description cron function container
type FunctionContainer struct {
}

func (fc FunctionContainer) DoSomething() (err error) {
	common.LogFuncDebug("DoSomething")
	return
}

//在这里执行一个定时器任务，调用我刚刚写的函数，然后去conf下配置定时器任务
func (fc FunctionContainer) DoEosOtcReport() (err error) {
	common.LogFuncDebug("DoEosOtcReport")
	EosOtcReport()
	return
}

//在这里执行一个定时器任务，调用我刚刚写的函数，然后去conf下配置定时器任务
func (fc FunctionContainer) DoDividendTest() (err error) {
	common.LogFuncDebug("GameUserMonthDividend")
	TestMonthDividend()
	return
}

func (fc FunctionContainer) DoGameWithdrawRiskAlert() {
	common.LogFuncDebug("DoGameWithdrawRiskAlert")
	GameWithdrawRiskAlert()
	return
}

func (fc FunctionContainer) DoStatisticGameReport() {
	common.LogFuncDebug("DoStatisticGameReport")
	GameStatisticReport()
	return
}

func (fc FunctionContainer) OtcStat() (err error) {
	OtcStatCron()
	return
}

func (fc FunctionContainer) ProfitReportDaily() {
	ProfitReportDailyCron()
	return
}

func (fc FunctionContainer) TaskGameTransferDaily() (err error) {
	TaskGameTransferDaily()
	return
}

func (fc FunctionContainer) TaskTeamDaily() (err error) {
	TaskTeamDaily()
	return
}

func (fc FunctionContainer) TaskReportCommission() (err error) {
	TaskReportCommission()
	return
}
