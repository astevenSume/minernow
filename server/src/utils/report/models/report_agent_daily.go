package models

//auto_models_start
 type ReportAgentDaily struct{
	Uid uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	SumWithdraw int64 `orm:"column(sum_withdraw)" json:"sum_withdraw,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *ReportAgentDaily) TableName() string {
    return "report_agent_daily"
}

//table report_agent_daily name and attributes defination.
const TABLE_ReportAgentDaily = "report_agent_daily"
const COLUMN_ReportAgentDaily_Uid = "uid"
const COLUMN_ReportAgentDaily_SumWithdraw = "sum_withdraw"
const COLUMN_ReportAgentDaily_Ctime = "ctime"
const ATTRIBUTE_ReportAgentDaily_Uid = "Uid"
const ATTRIBUTE_ReportAgentDaily_SumWithdraw = "SumWithdraw"
const ATTRIBUTE_ReportAgentDaily_Ctime = "Ctime"

//auto_models_end
