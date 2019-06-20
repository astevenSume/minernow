package models

//auto_models_start
type ServerNode struct {
	AppName  string `orm:"column(app_name);pk;size(256)" json:"app_name,omitempty"`
	RegionId int64  `orm:"column(region_id)" json:"region_id,omitempty"`
	ServerId int64  `orm:"column(server_id)" json:"server_id,omitempty"`
	LastPing uint32 `orm:"column(last_ping)" json:"last_ping,omitempty"`
}

func (this *ServerNode) TableName() string {
	return "server_node"
}

//table server_node name and attributes defination.
const TABLE_ServerNode = "server_node"
const COLUMN_ServerNode_AppName = "app_name"
const COLUMN_ServerNode_RegionId = "region_id"
const COLUMN_ServerNode_ServerId = "server_id"
const COLUMN_ServerNode_LastPing = "last_ping"
const ATTRIBUTE_ServerNode_AppName = "AppName"
const ATTRIBUTE_ServerNode_RegionId = "RegionId"
const ATTRIBUTE_ServerNode_ServerId = "ServerId"
const ATTRIBUTE_ServerNode_LastPing = "LastPing"

//auto_models_end
