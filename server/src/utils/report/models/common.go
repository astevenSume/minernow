package models

import "github.com/astaxie/beego/orm"

func ModelsInit() (err error) {
	orm.RegisterModel(
		new(GameUserMonthReport),
		new(MonthDividendRecord),
		new(ProfitReportDaily),
		new(ReportAgentDaily),
		new(ReportCommission),
		new(ReportEusdDaily),
		new(ReportGameRecordAg),
		new(ReportGameRecordKy),
		new(ReportGameRecordRg),
		new(ReportGameTransferDaily),
		new(ReportGameUserDaily),
		new(ReportStatisticGameAll),
		new(ReportStatisticSum),
		new(ReportTeamDaily),
		new(ReportTeamGameTransferDaily),
	)

	return
}
