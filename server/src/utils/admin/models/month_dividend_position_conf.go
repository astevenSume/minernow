package models

//auto_models_start
type MonthDividendPositionConf struct {
	Id            int64 `orm:"column(id);pk" json:"id,omitempty"`
	AgentLv       int32 `orm:"column(agent_lv)" json:"agent_lv,omitempty"`
	Position      int32 `orm:"column(position)" json:"position,omitempty"`
	Min           int64 `orm:"column(min)" json:"min,omitempty"`
	Max           int64 `orm:"column(max)" json:"max,omitempty"`
	ActivityNum   int32 `orm:"column(activity_num)" json:"activity_num,omitempty"`
	DividendRatio int32 `orm:"column(dividend_ratio)" json:"dividend_ratio,omitempty"`
	Ctime         int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime         int64 `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *MonthDividendPositionConf) TableName() string {
	return "month_dividend_position_conf"
}

//table month_dividend_position_conf name and attributes defination.
const TABLE_MonthDividendPositionConf = "month_dividend_position_conf"
const COLUMN_MonthDividendPositionConf_Id = "id"
const COLUMN_MonthDividendPositionConf_AgentLv = "agent_lv"
const COLUMN_MonthDividendPositionConf_Position = "position"
const COLUMN_MonthDividendPositionConf_Min = "min"
const COLUMN_MonthDividendPositionConf_Max = "max"
const COLUMN_MonthDividendPositionConf_ActivityNum = "activity_num"
const COLUMN_MonthDividendPositionConf_DividendRatio = "dividend_ratio"
const COLUMN_MonthDividendPositionConf_Ctime = "ctime"
const COLUMN_MonthDividendPositionConf_Utime = "utime"
const ATTRIBUTE_MonthDividendPositionConf_Id = "Id"
const ATTRIBUTE_MonthDividendPositionConf_AgentLv = "AgentLv"
const ATTRIBUTE_MonthDividendPositionConf_Position = "Position"
const ATTRIBUTE_MonthDividendPositionConf_Min = "Min"
const ATTRIBUTE_MonthDividendPositionConf_Max = "Max"
const ATTRIBUTE_MonthDividendPositionConf_ActivityNum = "ActivityNum"
const ATTRIBUTE_MonthDividendPositionConf_DividendRatio = "DividendRatio"
const ATTRIBUTE_MonthDividendPositionConf_Ctime = "Ctime"
const ATTRIBUTE_MonthDividendPositionConf_Utime = "Utime"

//auto_models_end
