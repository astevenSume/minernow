package models

//auto_models_start
type AgentPath struct {
	Uid              uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Level            uint32 `orm:"column(level)" json:"level,omitempty"`
	Sn               uint32 `orm:"column(sn)" json:"sn,omitempty"`
	Path             string `orm:"column(path)" json:"path,omitempty"`
	Ctime            int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Mtime            int64  `orm:"column(mtime)" json:"mtime,omitempty"`
	InviteCode       string `orm:"column(invite_code);size(100)" json:"invite_code,omitempty"`
	WhiteListId      uint32 `orm:"column(whitelist_id)" json:"whitelist_id,omitempty"`
	InviteNum        uint32 `orm:"column(invite_num)" json:"invite_num,omitempty"`
	ParentUid        uint64 `orm:"column(parent_uid)" json:"parent_uid,omitempty"`
	DividendPosition uint32 `orm:"column(dividend_position)" json:"dividend_position,omitempty"`
}

func (this *AgentPath) TableName() string {
	return "agent_path"
}

//table agent_path name and attributes defination.
const TABLE_AgentPath = "agent_path"
const COLUMN_AgentPath_Uid = "uid"
const COLUMN_AgentPath_Level = "level"
const COLUMN_AgentPath_Sn = "sn"
const COLUMN_AgentPath_Path = "path"
const COLUMN_AgentPath_Ctime = "ctime"
const COLUMN_AgentPath_Mtime = "mtime"
const COLUMN_AgentPath_InviteCode = "invite_code"
const COLUMN_AgentPath_WhiteListId = "whitelist_id"
const COLUMN_AgentPath_InviteNum = "invite_num"
const COLUMN_AgentPath_ParentUid = "parent_uid"
const COLUMN_AgentPath_DividendPosition = "dividend_position"
const ATTRIBUTE_AgentPath_Uid = "Uid"
const ATTRIBUTE_AgentPath_Level = "Level"
const ATTRIBUTE_AgentPath_Sn = "Sn"
const ATTRIBUTE_AgentPath_Path = "Path"
const ATTRIBUTE_AgentPath_Ctime = "Ctime"
const ATTRIBUTE_AgentPath_Mtime = "Mtime"
const ATTRIBUTE_AgentPath_InviteCode = "InviteCode"
const ATTRIBUTE_AgentPath_WhiteListId = "WhiteListId"
const ATTRIBUTE_AgentPath_InviteNum = "InviteNum"
const ATTRIBUTE_AgentPath_ParentUid = "ParentUid"
const ATTRIBUTE_AgentPath_DividendPosition = "DividendPosition"

//auto_models_end
