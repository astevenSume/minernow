package models

//auto_models_start
type SystemNotification struct {
	Nid              uint64 `orm:"column(nid);pk" json:"nid,omitempty"`
	NotificationType string `orm:"column(notification_type);size(100)" json:"notification_type,omitempty"`
	Content          string `orm:"column(content)" json:"content,omitempty"`
	Uid              uint64 `orm:"column(uid)" json:"uid,omitempty"`
	IsRead           int32  `orm:"column(is_read)" json:"is_read,omitempty"`
	Ctime            int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *SystemNotification) TableName() string {
	return "system_notification"
}

//table system_notification name and attributes defination.
const TABLE_SystemNotification = "system_notification"
const COLUMN_SystemNotification_Nid = "nid"
const COLUMN_SystemNotification_NotificationType = "notification_type"
const COLUMN_SystemNotification_Content = "content"
const COLUMN_SystemNotification_Uid = "uid"
const COLUMN_SystemNotification_IsRead = "is_read"
const COLUMN_SystemNotification_Ctime = "ctime"
const ATTRIBUTE_SystemNotification_Nid = "Nid"
const ATTRIBUTE_SystemNotification_NotificationType = "NotificationType"
const ATTRIBUTE_SystemNotification_Content = "Content"
const ATTRIBUTE_SystemNotification_Uid = "Uid"
const ATTRIBUTE_SystemNotification_IsRead = "IsRead"
const ATTRIBUTE_SystemNotification_Ctime = "Ctime"

//auto_models_end
