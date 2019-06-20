package models

//auto_models_start
 type TmpGamebeters struct{
	Uid uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	GameId uint32 `orm:"column(game_id)" json:"game_id,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Status string `orm:"column(status);size(20)" json:"status,omitempty"`
	BetType uint32 `orm:"column(bet_type)" json:"bet_type,omitempty"`
}

func (this *TmpGamebeters) TableName() string {
    return "tmp_game_beters"
}

//table tmp_game_beters name and attributes defination.
const TABLE_TmpGamebeters = "tmp_game_beters"
const COLUMN_TmpGamebeters_Uid = "uid"
const COLUMN_TmpGamebeters_ChannelId = "channel_id"
const COLUMN_TmpGamebeters_GameId = "game_id"
const COLUMN_TmpGamebeters_Ctime = "ctime"
const COLUMN_TmpGamebeters_Status = "status"
const COLUMN_TmpGamebeters_BetType = "bet_type"
const ATTRIBUTE_TmpGamebeters_Uid = "Uid"
const ATTRIBUTE_TmpGamebeters_ChannelId = "ChannelId"
const ATTRIBUTE_TmpGamebeters_GameId = "GameId"
const ATTRIBUTE_TmpGamebeters_Ctime = "Ctime"
const ATTRIBUTE_TmpGamebeters_Status = "Status"
const ATTRIBUTE_TmpGamebeters_BetType = "BetType"

//auto_models_end
