package models

//auto_models_start
type Config struct {
	Id     uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Action int8   `orm:"column(action)" json:"action,omitempty"`
	Key    string `orm:"column(key);size(256)" json:"key,omitempty"`
	Value  string `orm:"column(value)" json:"value,omitempty"`
	Desc   string `orm:"column(desc);size(256)" json:"desc,omitempty"`
	Ctime  int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *Config) TableName() string {
	return "config"
}

//table config name and attributes defination.
const TABLE_Config = "config"
const COLUMN_Config_Id = "id"
const COLUMN_Config_Action = "action"
const COLUMN_Config_Key = "key"
const COLUMN_Config_Value = "value"
const COLUMN_Config_Desc = "desc"
const COLUMN_Config_Ctime = "ctime"
const ATTRIBUTE_Config_Id = "Id"
const ATTRIBUTE_Config_Action = "Action"
const ATTRIBUTE_Config_Key = "Key"
const ATTRIBUTE_Config_Value = "Value"
const ATTRIBUTE_Config_Desc = "Desc"
const ATTRIBUTE_Config_Ctime = "Ctime"

//auto_models_end
