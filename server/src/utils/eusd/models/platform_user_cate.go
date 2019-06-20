package models

//auto_models_start
type PlatformUserCate struct {
	Id       int32  `orm:"column(id);pk" json:"id,omitempty"`
	Name     string `orm:"column(name);size(100)" json:"name,omitempty"`
	Dividend int32  `orm:"column(dividend)" json:"dividend,omitempty"`
	Ctime    uint32 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *PlatformUserCate) TableName() string {
	return "platform_user_cate"
}

//table platform_user_cate name and attributes defination.
const TABLE_PlatformUserCate = "platform_user_cate"
const COLUMN_PlatformUserCate_Id = "id"
const COLUMN_PlatformUserCate_Name = "name"
const COLUMN_PlatformUserCate_Dividend = "dividend"
const COLUMN_PlatformUserCate_Ctime = "ctime"
const ATTRIBUTE_PlatformUserCate_Id = "Id"
const ATTRIBUTE_PlatformUserCate_Name = "Name"
const ATTRIBUTE_PlatformUserCate_Dividend = "Dividend"
const ATTRIBUTE_PlatformUserCate_Ctime = "Ctime"

//auto_models_end
