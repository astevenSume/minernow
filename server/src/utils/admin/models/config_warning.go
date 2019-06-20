package models

//auto_models_start
type ConfigWarning struct {
	Id           uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Type         int8   `orm:"column(type)" json:"type,omitempty"`
	NationalCode string `orm:"column(national_code);size(16)" json:"national_code,omitempty"`
	Mobile       string `orm:"column(mobile);size(32)" json:"mobile,omitempty"`
	SmsType      int8   `orm:"column(sms_type)" json:"sms_type,omitempty"`
}

func (this *ConfigWarning) TableName() string {
	return "config_warning"
}

//table config_warning name and attributes defination.
const TABLE_ConfigWarning = "config_warning"
const COLUMN_ConfigWarning_Id = "id"
const COLUMN_ConfigWarning_Type = "type"
const COLUMN_ConfigWarning_NationalCode = "national_code"
const COLUMN_ConfigWarning_Mobile = "mobile"
const COLUMN_ConfigWarning_SmsType = "sms_type"
const ATTRIBUTE_ConfigWarning_Id = "Id"
const ATTRIBUTE_ConfigWarning_Type = "Type"
const ATTRIBUTE_ConfigWarning_NationalCode = "NationalCode"
const ATTRIBUTE_ConfigWarning_Mobile = "Mobile"
const ATTRIBUTE_ConfigWarning_SmsType = "SmsType"

//auto_models_end
