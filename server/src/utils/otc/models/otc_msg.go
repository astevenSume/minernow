package models

//auto_models_start
type OtcMsg struct {
	Id          uint64 `orm:"column(id);pk" json:"id,omitempty"`
	OrderId     int64  `orm:"column(order_id)" json:"order_id,omitempty"`
	Uid         uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Content     string `orm:"column(content)" json:"content,omitempty"`
	IsRead      uint8  `orm:"column(is_read)" json:"is_read,omitempty"`
	MessageType string `orm:"column(msg_type);size(200)" json:"msg_type,omitempty"`
	Ctime       int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *OtcMsg) TableName() string {
	return "otc_msg"
}

//table otc_msg name and attributes defination.
const TABLE_OtcMsg = "otc_msg"
const COLUMN_OtcMsg_Id = "id"
const COLUMN_OtcMsg_OrderId = "order_id"
const COLUMN_OtcMsg_Uid = "uid"
const COLUMN_OtcMsg_Content = "content"
const COLUMN_OtcMsg_IsRead = "is_read"
const COLUMN_OtcMsg_MessageType = "msg_type"
const COLUMN_OtcMsg_Ctime = "ctime"
const ATTRIBUTE_OtcMsg_Id = "Id"
const ATTRIBUTE_OtcMsg_OrderId = "OrderId"
const ATTRIBUTE_OtcMsg_Uid = "Uid"
const ATTRIBUTE_OtcMsg_Content = "Content"
const ATTRIBUTE_OtcMsg_IsRead = "IsRead"
const ATTRIBUTE_OtcMsg_MessageType = "MessageType"
const ATTRIBUTE_OtcMsg_Ctime = "Ctime"

//auto_models_end
