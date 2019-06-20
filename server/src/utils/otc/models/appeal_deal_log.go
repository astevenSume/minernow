package models

//auto_models_start
type AppealDealLog struct {
	Id       uint64 `orm:"column(id);pk" json:"id,omitempty"`
	AppealId uint64 `orm:"column(appeal_id)" json:"appeal_id,omitempty"`
	AdminId  uint32 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	OrderId  uint64 `orm:"column(order_id)" json:"order_id,omitempty"`
	Action   int8   `orm:"column(action)" json:"action,omitempty"`
	Ctime    int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *AppealDealLog) TableName() string {
	return "appeal_deal_log"
}

//table appeal_deal_log name and attributes defination.
const TABLE_AppealDealLog = "appeal_deal_log"
const COLUMN_AppealDealLog_Id = "id"
const COLUMN_AppealDealLog_AppealId = "appeal_id"
const COLUMN_AppealDealLog_AdminId = "admin_id"
const COLUMN_AppealDealLog_OrderId = "order_id"
const COLUMN_AppealDealLog_Action = "action"
const COLUMN_AppealDealLog_Ctime = "ctime"
const ATTRIBUTE_AppealDealLog_Id = "Id"
const ATTRIBUTE_AppealDealLog_AppealId = "AppealId"
const ATTRIBUTE_AppealDealLog_AdminId = "AdminId"
const ATTRIBUTE_AppealDealLog_OrderId = "OrderId"
const ATTRIBUTE_AppealDealLog_Action = "Action"
const ATTRIBUTE_AppealDealLog_Ctime = "Ctime"

//auto_models_end
