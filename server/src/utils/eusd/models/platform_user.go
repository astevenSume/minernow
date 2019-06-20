package models

//auto_models_start
type PlatformUser struct {
	Uid    uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Pid    int32  `orm:"column(pid)" json:"pid,omitempty"`
	Status int8   `orm:"column(status)" json:"status,omitempty"`
	Ctime  uint32 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *PlatformUser) TableName() string {
	return "platform_user"
}

//table platform_user name and attributes defination.
const TABLE_PlatformUser = "platform_user"
const COLUMN_PlatformUser_Uid = "uid"
const COLUMN_PlatformUser_Pid = "pid"
const COLUMN_PlatformUser_Status = "status"
const COLUMN_PlatformUser_Ctime = "ctime"
const ATTRIBUTE_PlatformUser_Uid = "Uid"
const ATTRIBUTE_PlatformUser_Pid = "Pid"
const ATTRIBUTE_PlatformUser_Status = "Status"
const ATTRIBUTE_PlatformUser_Ctime = "Ctime"

//auto_models_end
