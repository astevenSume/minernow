package models

//auto_models_start
type TopAgent struct {
	Id           uint32 `orm:"column(id);pk" json:"id,omitempty"`
	NationalCode string `orm:"column(national_code);size(16)" json:"national_code,omitempty"`
	Mobile       string `orm:"column(mobile);size(32)" json:"mobile,omitempty"`
	Status       int8   `orm:"column(status)" json:"status,omitempty"`
	Ctime        int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime        int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *TopAgent) TableName() string {
	return "top_agent"
}

//table top_agent name and attributes defination.
const TABLE_TopAgent = "top_agent"
const COLUMN_TopAgent_Id = "id"
const COLUMN_TopAgent_NationalCode = "national_code"
const COLUMN_TopAgent_Mobile = "mobile"
const COLUMN_TopAgent_Status = "status"
const COLUMN_TopAgent_Ctime = "ctime"
const COLUMN_TopAgent_Utime = "utime"
const ATTRIBUTE_TopAgent_Id = "Id"
const ATTRIBUTE_TopAgent_NationalCode = "NationalCode"
const ATTRIBUTE_TopAgent_Mobile = "Mobile"
const ATTRIBUTE_TopAgent_Status = "Status"
const ATTRIBUTE_TopAgent_Ctime = "Ctime"
const ATTRIBUTE_TopAgent_Utime = "Utime"

//auto_models_end
