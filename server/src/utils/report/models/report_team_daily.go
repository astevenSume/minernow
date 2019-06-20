package models

//auto_models_start
type ReportTeamDaily struct {
	Uid      uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	EusdBuy  int64  `orm:"column(eusd_buy)" json:"eusd_buy,omitempty"`
	EusdSell int64  `orm:"column(eusd_sell)" json:"eusd_sell,omitempty"`
	Level    uint32 `orm:"column(level)" json:"level,omitempty"`
	Ctime    int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *ReportTeamDaily) TableName() string {
	return "report_team_daily"
}

//table report_team_daily name and attributes defination.
const TABLE_ReportTeamDaily = "report_team_daily"
const COLUMN_ReportTeamDaily_Uid = "uid"
const COLUMN_ReportTeamDaily_EusdBuy = "eusd_buy"
const COLUMN_ReportTeamDaily_EusdSell = "eusd_sell"
const COLUMN_ReportTeamDaily_Level = "level"
const COLUMN_ReportTeamDaily_Ctime = "ctime"
const ATTRIBUTE_ReportTeamDaily_Uid = "Uid"
const ATTRIBUTE_ReportTeamDaily_EusdBuy = "EusdBuy"
const ATTRIBUTE_ReportTeamDaily_EusdSell = "EusdSell"
const ATTRIBUTE_ReportTeamDaily_Level = "Level"
const ATTRIBUTE_ReportTeamDaily_Ctime = "Ctime"

//auto_models_end
