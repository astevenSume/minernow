package models

//auto_models_start
type Role struct {
	Id    uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Name  string `orm:"column(name);size(100)" json:"name,omitempty"`
	Desc  string `orm:"column(desc);size(100)" json:"desc,omitempty"`
	Ctime int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *Role) TableName() string {
	return "role"
}

//table role name and attributes defination.
const TABLE_Role = "role"
const COLUMN_Role_Id = "id"
const COLUMN_Role_Name = "name"
const COLUMN_Role_Desc = "desc"
const COLUMN_Role_Ctime = "ctime"
const COLUMN_Role_Utime = "utime"
const ATTRIBUTE_Role_Id = "Id"
const ATTRIBUTE_Role_Name = "Name"
const ATTRIBUTE_Role_Desc = "Desc"
const ATTRIBUTE_Role_Ctime = "Ctime"
const ATTRIBUTE_Role_Utime = "Utime"

//auto_models_end
