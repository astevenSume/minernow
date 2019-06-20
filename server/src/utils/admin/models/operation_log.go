package models

//auto_models_start
type OperationLog struct {
	Id           uint64 `orm:"column(id);pk" json:"id,omitempty"`
	AdminId      uint64 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	Method       string `orm:"column(method);size(100)" json:"method,omitempty"`
	Route        string `orm:"column(route);size(100)" json:"route,omitempty"`
	Action       int32  `orm:"column(action)" json:"action,omitempty"`
	Input        string `orm:"column(input);size(65535)" json:"input,omitempty"`
	UserAgent    string `orm:"column(user_agent);size(512)" json:"user_agent,omitempty"`
	Ips          string `orm:"column(ips);size(100)" json:"ips,omitempty"`
	ResponseCode int32  `orm:"column(response_code)" json:"response_code,omitempty"`
	Ctime        int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *OperationLog) TableName() string {
	return "operation_log"
}

//table operation_log name and attributes defination.
const TABLE_OperationLog = "operation_log"
const COLUMN_OperationLog_Id = "id"
const COLUMN_OperationLog_AdminId = "admin_id"
const COLUMN_OperationLog_Method = "method"
const COLUMN_OperationLog_Route = "route"
const COLUMN_OperationLog_Action = "action"
const COLUMN_OperationLog_Input = "input"
const COLUMN_OperationLog_UserAgent = "user_agent"
const COLUMN_OperationLog_Ips = "ips"
const COLUMN_OperationLog_ResponseCode = "response_code"
const COLUMN_OperationLog_Ctime = "ctime"
const ATTRIBUTE_OperationLog_Id = "Id"
const ATTRIBUTE_OperationLog_AdminId = "AdminId"
const ATTRIBUTE_OperationLog_Method = "Method"
const ATTRIBUTE_OperationLog_Route = "Route"
const ATTRIBUTE_OperationLog_Action = "Action"
const ATTRIBUTE_OperationLog_Input = "Input"
const ATTRIBUTE_OperationLog_UserAgent = "UserAgent"
const ATTRIBUTE_OperationLog_Ips = "Ips"
const ATTRIBUTE_OperationLog_ResponseCode = "ResponseCode"
const ATTRIBUTE_OperationLog_Ctime = "Ctime"

//auto_models_end
