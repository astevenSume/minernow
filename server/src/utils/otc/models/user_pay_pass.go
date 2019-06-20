package models

//auto_models_start
type UserPayPassword struct {
	Uid        uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Password   string `orm:"column(pass);size(256)" json:"pass,omitempty"`
	Salt       string `orm:"column(salt);size(256)" json:"salt,omitempty"`
	SignSalt   string `orm:"column(sign_salt);size(256)" json:"sign_salt,omitempty"`
	Status     int8   `orm:"column(status)" json:"status,omitempty"`
	Method     int8   `orm:"column(method)" json:"method,omitempty"`
	VerifyStep int8   `orm:"column(verify_step)" json:"verify_step,omitempty"`
	Timestamp  int64  `orm:"column(timestamp)" json:"timestamp,omitempty"`
}

func (this *UserPayPassword) TableName() string {
	return "user_pay_pass"
}

//table user_pay_pass name and attributes defination.
const TABLE_UserPayPassword = "user_pay_pass"
const COLUMN_UserPayPassword_Uid = "uid"
const COLUMN_UserPayPassword_Password = "pass"
const COLUMN_UserPayPassword_Salt = "salt"
const COLUMN_UserPayPassword_SignSalt = "sign_salt"
const COLUMN_UserPayPassword_Status = "status"
const COLUMN_UserPayPassword_Method = "method"
const COLUMN_UserPayPassword_VerifyStep = "verify_step"
const COLUMN_UserPayPassword_Timestamp = "timestamp"
const ATTRIBUTE_UserPayPassword_Uid = "Uid"
const ATTRIBUTE_UserPayPassword_Password = "Password"
const ATTRIBUTE_UserPayPassword_Salt = "Salt"
const ATTRIBUTE_UserPayPassword_SignSalt = "SignSalt"
const ATTRIBUTE_UserPayPassword_Status = "Status"
const ATTRIBUTE_UserPayPassword_Method = "Method"
const ATTRIBUTE_UserPayPassword_VerifyStep = "VerifyStep"
const ATTRIBUTE_UserPayPassword_Timestamp = "Timestamp"

//auto_models_end
