package models

//auto_models_start
type Task struct {
	Id       uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Name     string `orm:"column(name);size(256)" json:"name,omitempty"`
	Alia     string `orm:"column(alia);size(256)" json:"alia,omitempty"`
	AppName  string `orm:"column(app_name);size(256)" json:"app_name,omitempty"`
	FuncName string `orm:"column(func_name);size(256)" json:"func_name,omitempty"`
	Spec     string `orm:"column(spec);size(256)" json:"spec,omitempty"`
	Status   uint8  `orm:"column(status)" json:"status,omitempty"`
	Ctime    uint32 `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime    uint32 `orm:"column(utime)" json:"utime,omitempty"`
	Desc     string `orm:"column(desc);size(256)" json:"desc,omitempty"`
}

func (this *Task) TableName() string {
	return "task"
}

//table task name and attributes defination.
const TABLE_Task = "task"
const COLUMN_Task_Id = "id"
const COLUMN_Task_Name = "name"
const COLUMN_Task_Alia = "alia"
const COLUMN_Task_AppName = "app_name"
const COLUMN_Task_FuncName = "func_name"
const COLUMN_Task_Spec = "spec"
const COLUMN_Task_Status = "status"
const COLUMN_Task_Ctime = "ctime"
const COLUMN_Task_Utime = "utime"
const COLUMN_Task_Desc = "desc"
const ATTRIBUTE_Task_Id = "Id"
const ATTRIBUTE_Task_Name = "Name"
const ATTRIBUTE_Task_Alia = "Alia"
const ATTRIBUTE_Task_AppName = "AppName"
const ATTRIBUTE_Task_FuncName = "FuncName"
const ATTRIBUTE_Task_Spec = "Spec"
const ATTRIBUTE_Task_Status = "Status"
const ATTRIBUTE_Task_Ctime = "Ctime"
const ATTRIBUTE_Task_Utime = "Utime"
const ATTRIBUTE_Task_Desc = "Desc"

//auto_models_end
