package models

//auto_models_start
 type GameLog struct{
	Id uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid uint64 `orm:"column(uid)" json:"uid,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	Account string `orm:"column(account);size(50)" json:"account,omitempty"`
	LogType uint8 `orm:"column(log_type)" json:"log_type,omitempty"`
	Desc string `orm:"column(desc);size(512)" json:"desc,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *GameLog) TableName() string {
    return "game_log"
}

//table game_log name and attributes defination.
const TABLE_GameLog = "game_log"
const COLUMN_GameLog_Id = "id"
const COLUMN_GameLog_Uid = "uid"
const COLUMN_GameLog_ChannelId = "channel_id"
const COLUMN_GameLog_Account = "account"
const COLUMN_GameLog_LogType = "log_type"
const COLUMN_GameLog_Desc = "desc"
const COLUMN_GameLog_Ctime = "ctime"
const ATTRIBUTE_GameLog_Id = "Id"
const ATTRIBUTE_GameLog_Uid = "Uid"
const ATTRIBUTE_GameLog_ChannelId = "ChannelId"
const ATTRIBUTE_GameLog_Account = "Account"
const ATTRIBUTE_GameLog_LogType = "LogType"
const ATTRIBUTE_GameLog_Desc = "Desc"
const ATTRIBUTE_GameLog_Ctime = "Ctime"

//auto_models_end
