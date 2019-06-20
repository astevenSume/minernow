package models

//auto_models_start
type AgentWhiteList struct {
	Id         uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Name       string `orm:"column(name);size(64)" json:"name,omitempty"`
	Commission int32  `orm:"column(commission)" json:"commission,omitempty"`
	Precision  int32  `orm:"column(precision)" json:"precision,omitempty"`
	Ctime      int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime      int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *AgentWhiteList) TableName() string {
	return "agent_white_list"
}

//table agent_white_list name and attributes defination.
const TABLE_AgentWhiteList = "agent_white_list"
const COLUMN_AgentWhiteList_Id = "id"
const COLUMN_AgentWhiteList_Name = "name"
const COLUMN_AgentWhiteList_Commission = "commission"
const COLUMN_AgentWhiteList_Precision = "precision"
const COLUMN_AgentWhiteList_Ctime = "ctime"
const COLUMN_AgentWhiteList_Utime = "utime"
const ATTRIBUTE_AgentWhiteList_Id = "Id"
const ATTRIBUTE_AgentWhiteList_Name = "Name"
const ATTRIBUTE_AgentWhiteList_Commission = "Commission"
const ATTRIBUTE_AgentWhiteList_Precision = "Precision"
const ATTRIBUTE_AgentWhiteList_Ctime = "Ctime"
const ATTRIBUTE_AgentWhiteList_Utime = "Utime"

//auto_models_end
