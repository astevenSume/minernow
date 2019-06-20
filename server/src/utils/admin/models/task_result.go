package models

//auto_models_start
type TaskResult struct {
	Id        uint32 `orm:"column(id);pk" json:"id,omitempty"`
	AppName   string `orm:"column(app_name);size(256)" json:"app_name,omitempty"`
	RegionId  int64  `orm:"column(region_id)" json:"region_id,omitempty"`
	ServerId  int64  `orm:"column(server_id)" json:"server_id,omitempty"`
	Name      string `orm:"column(name);size(256)" json:"name,omitempty"`
	Code      int32  `orm:"column(code)" json:"code,omitempty"`
	Detail    string `orm:"column(detail);size(256)" json:"detail,omitempty"`
	EndTime   uint32 `orm:"column(end_time)" json:"end_time,omitempty"`
	BeginTime uint32 `orm:"column(begin_time)" json:"begin_time,omitempty"`
	Ctime     uint32 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *TaskResult) TableName() string {
	return "task_result"
}

//table task_result name and attributes defination.
const TABLE_TaskResult = "task_result"
const COLUMN_TaskResult_Id = "id"
const COLUMN_TaskResult_AppName = "app_name"
const COLUMN_TaskResult_RegionId = "region_id"
const COLUMN_TaskResult_ServerId = "server_id"
const COLUMN_TaskResult_Name = "name"
const COLUMN_TaskResult_Code = "code"
const COLUMN_TaskResult_Detail = "detail"
const COLUMN_TaskResult_EndTime = "end_time"
const COLUMN_TaskResult_BeginTime = "begin_time"
const COLUMN_TaskResult_Ctime = "ctime"
const ATTRIBUTE_TaskResult_Id = "Id"
const ATTRIBUTE_TaskResult_AppName = "AppName"
const ATTRIBUTE_TaskResult_RegionId = "RegionId"
const ATTRIBUTE_TaskResult_ServerId = "ServerId"
const ATTRIBUTE_TaskResult_Name = "Name"
const ATTRIBUTE_TaskResult_Code = "Code"
const ATTRIBUTE_TaskResult_Detail = "Detail"
const ATTRIBUTE_TaskResult_EndTime = "EndTime"
const ATTRIBUTE_TaskResult_BeginTime = "BeginTime"
const ATTRIBUTE_TaskResult_Ctime = "Ctime"

//auto_models_end
