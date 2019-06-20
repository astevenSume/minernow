package models

//auto_models_start
type UserLoginLog struct {
	Id        uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Userid    uint64 `orm:"column(user_id)" json:"user_id,omitempty"`
	UserAgent string `orm:"column(user_agent);null;size(256)" json:"user_agent,omitempty"`
	Ips       string `orm:"column(ips);null;size(256)" json:"ips,omitempty"`
	Ctime     int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *UserLoginLog) TableName() string {
	return "user_login_log"
}

//table user_login_log name and attributes defination.
const TABLE_UserLoginLog = "user_login_log"
const COLUMN_UserLoginLog_Id = "id"
const COLUMN_UserLoginLog_Userid = "user_id"
const COLUMN_UserLoginLog_UserAgent = "user_agent"
const COLUMN_UserLoginLog_Ips = "ips"
const COLUMN_UserLoginLog_Ctime = "ctime"
const ATTRIBUTE_UserLoginLog_Id = "Id"
const ATTRIBUTE_UserLoginLog_Userid = "Userid"
const ATTRIBUTE_UserLoginLog_UserAgent = "UserAgent"
const ATTRIBUTE_UserLoginLog_Ips = "Ips"
const ATTRIBUTE_UserLoginLog_Ctime = "Ctime"

//auto_models_end
