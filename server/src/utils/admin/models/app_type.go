package models

//auto_models_start
type AppType struct {
	Id    uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Name  string `orm:"column(name);size(256)" json:"name,omitempty"`
	Desc  string `orm:"column(desc);size(256)" json:"desc,omitempty"`
	Ctime int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *AppType) TableName() string {
	return "app_type"
}

//table app_type name and attributes defination.
const TABLE_AppType = "app_type"
const COLUMN_AppType_Id = "id"
const COLUMN_AppType_Name = "name"
const COLUMN_AppType_Desc = "desc"
const COLUMN_AppType_Ctime = "ctime"
const COLUMN_AppType_Utime = "utime"
const ATTRIBUTE_AppType_Id = "Id"
const ATTRIBUTE_AppType_Name = "Name"
const ATTRIBUTE_AppType_Desc = "Desc"
const ATTRIBUTE_AppType_Ctime = "Ctime"
const ATTRIBUTE_AppType_Utime = "Utime"

//auto_models_end
