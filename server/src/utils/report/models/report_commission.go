package models

//auto_models_start
type ReportCommission struct {
	Uid             uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Level           uint32 `orm:"column(level)" json:"level,omitempty"`
	TeamWithdraw    int64  `orm:"column(team_withdraw)" json:"team_withdraw,omitempty"`
	TeamCanWithdraw int64  `orm:"column(team_can_withdraw)" json:"team_can_withdraw,omitempty"`
	Ctime           int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *ReportCommission) TableName() string {
	return "report_commission"
}

//table report_commission name and attributes defination.
const TABLE_ReportCommission = "report_commission"
const COLUMN_ReportCommission_Uid = "uid"
const COLUMN_ReportCommission_Level = "level"
const COLUMN_ReportCommission_TeamWithdraw = "team_withdraw"
const COLUMN_ReportCommission_TeamCanWithdraw = "team_can_withdraw"
const COLUMN_ReportCommission_Ctime = "ctime"
const ATTRIBUTE_ReportCommission_Uid = "Uid"
const ATTRIBUTE_ReportCommission_Level = "Level"
const ATTRIBUTE_ReportCommission_TeamWithdraw = "TeamWithdraw"
const ATTRIBUTE_ReportCommission_TeamCanWithdraw = "TeamCanWithdraw"
const ATTRIBUTE_ReportCommission_Ctime = "Ctime"

//auto_models_end
