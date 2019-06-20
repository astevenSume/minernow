package models

//auto_models_start
type Smstemplates struct {
	Id       int64  `orm:"column(id);pk" json:"id,omitempty"`
	Name     string `orm:"column(name);size(100)" json:"name,omitempty"`
	Type     int8   `orm:"column(type)" json:"type,omitempty"`
	Template string `orm:"column(template);size(256)" json:"template,omitempty"`
	Ctime    int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime    int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *Smstemplates) TableName() string {
	return "smstemplates"
}

//table smstemplates name and attributes defination.
const TABLE_Smstemplates = "smstemplates"
const COLUMN_Smstemplates_Id = "id"
const COLUMN_Smstemplates_Name = "name"
const COLUMN_Smstemplates_Type = "type"
const COLUMN_Smstemplates_Template = "template"
const COLUMN_Smstemplates_Ctime = "ctime"
const COLUMN_Smstemplates_Utime = "utime"
const ATTRIBUTE_Smstemplates_Id = "Id"
const ATTRIBUTE_Smstemplates_Name = "Name"
const ATTRIBUTE_Smstemplates_Type = "Type"
const ATTRIBUTE_Smstemplates_Template = "Template"
const ATTRIBUTE_Smstemplates_Ctime = "Ctime"
const ATTRIBUTE_Smstemplates_Utime = "Utime"

//auto_models_end
