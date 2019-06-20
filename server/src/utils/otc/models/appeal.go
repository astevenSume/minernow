package models

//auto_models_start
type Appeal struct {
	Id      uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Type    int8   `orm:"column(type)" json:"type,omitempty"`
	UserId  uint64 `orm:"column(user_id)" json:"user_id,omitempty"`
	AdminId uint32 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	OrderId uint64 `orm:"column(order_id)" json:"order_id,omitempty"`
	Context string `orm:"column(context);size(256)" json:"context,omitempty"`
	Status  int8   `orm:"column(status)" json:"status,omitempty"`
	Ctime   int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime   int64  `orm:"column(utime)" json:"utime,omitempty"`
	Wechat  string `orm:"column(wechat);size(100)" json:"wechat,omitempty"`
}

func (this *Appeal) TableName() string {
	return "appeal"
}

//table appeal name and attributes defination.
const TABLE_Appeal = "appeal"
const COLUMN_Appeal_Id = "id"
const COLUMN_Appeal_Type = "type"
const COLUMN_Appeal_UserId = "user_id"
const COLUMN_Appeal_AdminId = "admin_id"
const COLUMN_Appeal_OrderId = "order_id"
const COLUMN_Appeal_Context = "context"
const COLUMN_Appeal_Status = "status"
const COLUMN_Appeal_Ctime = "ctime"
const COLUMN_Appeal_Utime = "utime"
const COLUMN_Appeal_Wechat = "wechat"
const ATTRIBUTE_Appeal_Id = "Id"
const ATTRIBUTE_Appeal_Type = "Type"
const ATTRIBUTE_Appeal_UserId = "UserId"
const ATTRIBUTE_Appeal_AdminId = "AdminId"
const ATTRIBUTE_Appeal_OrderId = "OrderId"
const ATTRIBUTE_Appeal_Context = "Context"
const ATTRIBUTE_Appeal_Status = "Status"
const ATTRIBUTE_Appeal_Ctime = "Ctime"
const ATTRIBUTE_Appeal_Utime = "Utime"
const ATTRIBUTE_Appeal_Wechat = "Wechat"

//auto_models_end
