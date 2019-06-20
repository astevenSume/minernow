package models

//auto_models_start
type SysNotification struct {
	Id      uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Content string `orm:"column(content);size(128)" json:"content,omitempty"`
	AdminId uint32 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	Status  int8   `orm:"column(status)" json:"status,omitempty"`
	Ctime   int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime   int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *SysNotification) TableName() string {
	return "sys_notification"
}

//table sys_notification name and attributes defination.
const TABLE_SysNotification = "sys_notification"
const COLUMN_SysNotification_Id = "id"
const COLUMN_SysNotification_Content = "content"
const COLUMN_SysNotification_AdminId = "admin_id"
const COLUMN_SysNotification_Status = "status"
const COLUMN_SysNotification_Ctime = "ctime"
const COLUMN_SysNotification_Utime = "utime"
const ATTRIBUTE_SysNotification_Id = "Id"
const ATTRIBUTE_SysNotification_Content = "Content"
const ATTRIBUTE_SysNotification_AdminId = "AdminId"
const ATTRIBUTE_SysNotification_Status = "Status"
const ATTRIBUTE_SysNotification_Ctime = "Ctime"
const ATTRIBUTE_SysNotification_Utime = "Utime"

//auto_models_end
