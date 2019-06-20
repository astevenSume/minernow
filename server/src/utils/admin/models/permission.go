package models

//auto_models_start
type Permission struct {
	Id    uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Slug  string `orm:"column(slug);size(100)" json:"slug,omitempty"`
	Desc  string `orm:"column(desc);size(100)" json:"desc,omitempty"`
	Ctime int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime int64  `orm:"column(utime)" json:"utime,omitempty"`
	Dtime int64  `orm:"column(dtime)" json:"dtime,omitempty"`
}

func (this *Permission) TableName() string {
	return "permission"
}

//table permission name and attributes defination.
const TABLE_Permission = "permission"
const COLUMN_Permission_Id = "id"
const COLUMN_Permission_Slug = "slug"
const COLUMN_Permission_Desc = "desc"
const COLUMN_Permission_Ctime = "ctime"
const COLUMN_Permission_Utime = "utime"
const COLUMN_Permission_Dtime = "dtime"
const ATTRIBUTE_Permission_Id = "Id"
const ATTRIBUTE_Permission_Slug = "Slug"
const ATTRIBUTE_Permission_Desc = "Desc"
const ATTRIBUTE_Permission_Ctime = "Ctime"
const ATTRIBUTE_Permission_Utime = "Utime"
const ATTRIBUTE_Permission_Dtime = "Dtime"

//auto_models_end
