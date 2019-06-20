package models

//auto_models_start
type AgentChannelCommission struct {
	Uid       uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	Ctime     int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Mtime     int64  `orm:"column(mtime)" json:"mtime,omitempty"`
	Integer   int32  `orm:"column(integer)" json:"integer,omitempty"`
	Decimals  int32  `orm:"column(decimals)" json:"decimals,omitempty"`
	Status    uint8  `orm:"column(status)" json:"status,omitempty"`
}

func (this *AgentChannelCommission) TableName() string {
	return "agent_channel_commission"
}

//table agent_channel_commission name and attributes defination.
const TABLE_AgentChannelCommission = "agent_channel_commission"
const COLUMN_AgentChannelCommission_Uid = "uid"
const COLUMN_AgentChannelCommission_ChannelId = "channel_id"
const COLUMN_AgentChannelCommission_Ctime = "ctime"
const COLUMN_AgentChannelCommission_Mtime = "mtime"
const COLUMN_AgentChannelCommission_Integer = "integer"
const COLUMN_AgentChannelCommission_Decimals = "decimals"
const COLUMN_AgentChannelCommission_Status = "status"
const ATTRIBUTE_AgentChannelCommission_Uid = "Uid"
const ATTRIBUTE_AgentChannelCommission_ChannelId = "ChannelId"
const ATTRIBUTE_AgentChannelCommission_Ctime = "Ctime"
const ATTRIBUTE_AgentChannelCommission_Mtime = "Mtime"
const ATTRIBUTE_AgentChannelCommission_Integer = "Integer"
const ATTRIBUTE_AgentChannelCommission_Decimals = "Decimals"
const ATTRIBUTE_AgentChannelCommission_Status = "Status"

//auto_models_end
