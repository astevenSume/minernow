package models

//auto_models_start
 type GameRiskAlert struct{
	Id uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Funds int64 `orm:"column(funds)" json:"funds,omitempty"`
	EusdNum int64 `orm:"column(eusd_num)" json:"eusd_num,omitempty"`
	OrderTime uint64 `orm:"column(order_time)" json:"order_time,omitempty"`
	AlertTime uint64 `orm:"column(alert_time)" json:"alert_time,omitempty"`
	DoGet uint8 `orm:"column(do_get)" json:"do_get,omitempty"`
	WarnGrade uint8 `orm:"column(warn_grade)" json:"warn_grade,omitempty"`
	RiskType uint8 `orm:"column(risk_type)" json:"risk_type,omitempty"`
	OrderRiskId uint64 `orm:"column(order_risk_id)" json:"order_risk_id,omitempty"`
}

func (this *GameRiskAlert) TableName() string {
    return "game_risk_alert"
}

//table game_risk_alert name and attributes defination.
const TABLE_GameRiskAlert = "game_risk_alert"
const COLUMN_GameRiskAlert_Id = "id"
const COLUMN_GameRiskAlert_Uid = "uid"
const COLUMN_GameRiskAlert_Funds = "funds"
const COLUMN_GameRiskAlert_EusdNum = "eusd_num"
const COLUMN_GameRiskAlert_OrderTime = "order_time"
const COLUMN_GameRiskAlert_AlertTime = "alert_time"
const COLUMN_GameRiskAlert_DoGet = "do_get"
const COLUMN_GameRiskAlert_WarnGrade = "warn_grade"
const COLUMN_GameRiskAlert_RiskType = "risk_type"
const COLUMN_GameRiskAlert_OrderRiskId = "order_risk_id"
const ATTRIBUTE_GameRiskAlert_Id = "Id"
const ATTRIBUTE_GameRiskAlert_Uid = "Uid"
const ATTRIBUTE_GameRiskAlert_Funds = "Funds"
const ATTRIBUTE_GameRiskAlert_EusdNum = "EusdNum"
const ATTRIBUTE_GameRiskAlert_OrderTime = "OrderTime"
const ATTRIBUTE_GameRiskAlert_AlertTime = "AlertTime"
const ATTRIBUTE_GameRiskAlert_DoGet = "DoGet"
const ATTRIBUTE_GameRiskAlert_WarnGrade = "WarnGrade"
const ATTRIBUTE_GameRiskAlert_RiskType = "RiskType"
const ATTRIBUTE_GameRiskAlert_OrderRiskId = "OrderRiskId"

//auto_models_end
