package models

//auto_models_start
type Smscodes struct {
	Id           int64  `orm:"column(id);pk" json:"id,omitempty"`
	NationalCode string `orm:"column(national_code);size(32)" json:"national_code,omitempty"`
	Mobile       string `orm:"column(mobile);size(32)" json:"mobile,omitempty"`
	Action       string `orm:"column(action);size(100)" json:"action,omitempty"`
	Code         string `orm:"column(code);size(16)" json:"code,omitempty"`
	Status       int8   `orm:"column(status)" json:"status,omitempty"`
	Ctime        int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Etime        int64  `orm:"column(etime)" json:"etime,omitempty"`
}

func (this *Smscodes) TableName() string {
	return "smscodes"
}

//table smscodes name and attributes defination.
const TABLE_Smscodes = "smscodes"
const COLUMN_Smscodes_Id = "id"
const COLUMN_Smscodes_NationalCode = "national_code"
const COLUMN_Smscodes_Mobile = "mobile"
const COLUMN_Smscodes_Action = "action"
const COLUMN_Smscodes_Code = "code"
const COLUMN_Smscodes_Status = "status"
const COLUMN_Smscodes_Ctime = "ctime"
const COLUMN_Smscodes_Etime = "etime"
const ATTRIBUTE_Smscodes_Id = "Id"
const ATTRIBUTE_Smscodes_NationalCode = "NationalCode"
const ATTRIBUTE_Smscodes_Mobile = "Mobile"
const ATTRIBUTE_Smscodes_Action = "Action"
const ATTRIBUTE_Smscodes_Code = "Code"
const ATTRIBUTE_Smscodes_Status = "Status"
const ATTRIBUTE_Smscodes_Ctime = "Ctime"
const ATTRIBUTE_Smscodes_Etime = "Etime"

//auto_models_end
