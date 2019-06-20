package models

//auto_models_start
 type GameOrderRisk struct{
	Id uint64 `orm:"column(id);pk" json:"id,omitempty"`
	AlertId int64 `orm:"column(alert_id)" json:"alert_id,omitempty"`
	Uid uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Amount int64 `orm:"column(amount)" json:"amount,omitempty"`
	Fund int64 `orm:"column(funds)" json:"funds,omitempty"`
	PayType int8 `orm:"column(pay_type)" json:"pay_type,omitempty"`
	PayAccount string `orm:"column(pay_account);size(300)" json:"pay_account,omitempty"`
	OrderTime int64 `orm:"column(order_time)" json:"order_time,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *GameOrderRisk) TableName() string {
    return "game_order_risk"
}

//table game_order_risk name and attributes defination.
const TABLE_GameOrderRisk = "game_order_risk"
const COLUMN_GameOrderRisk_Id = "id"
const COLUMN_GameOrderRisk_AlertId = "alert_id"
const COLUMN_GameOrderRisk_Uid = "uid"
const COLUMN_GameOrderRisk_Amount = "amount"
const COLUMN_GameOrderRisk_Fund = "funds"
const COLUMN_GameOrderRisk_PayType = "pay_type"
const COLUMN_GameOrderRisk_PayAccount = "pay_account"
const COLUMN_GameOrderRisk_OrderTime = "order_time"
const COLUMN_GameOrderRisk_Ctime = "ctime"
const ATTRIBUTE_GameOrderRisk_Id = "Id"
const ATTRIBUTE_GameOrderRisk_AlertId = "AlertId"
const ATTRIBUTE_GameOrderRisk_Uid = "Uid"
const ATTRIBUTE_GameOrderRisk_Amount = "Amount"
const ATTRIBUTE_GameOrderRisk_Fund = "Fund"
const ATTRIBUTE_GameOrderRisk_PayType = "PayType"
const ATTRIBUTE_GameOrderRisk_PayAccount = "PayAccount"
const ATTRIBUTE_GameOrderRisk_OrderTime = "OrderTime"
const ATTRIBUTE_GameOrderRisk_Ctime = "Ctime"

//auto_models_end
