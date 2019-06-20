package models

//auto_models_start
 type GameUserDaily struct{
	ChannelId uint32 `orm:"column(channel_id);pk" json:"channel_id,omitempty"`
	Uid uint64 `orm:"column(uid)" json:"uid,omitempty"`
	TaxInteger int32 `orm:"column(tax_integer)" json:"tax_integer,omitempty"`
	TaxDecimals int32 `orm:"column(tax_decimals)" json:"tax_decimals,omitempty"`
	ChipsInteger int32 `orm:"column(chips_integer)" json:"chips_integer,omitempty"`
	ChipsDecimals int32 `orm:"column(chips_decimals)" json:"chips_decimals,omitempty"`
	WinloseInteger int32 `orm:"column(winlose_integer)" json:"winlose_integer,omitempty"`
	WinloseDecimals int32 `orm:"column(winlose_decimals)" json:"winlose_decimals,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Mtime int64 `orm:"column(mtime)" json:"mtime,omitempty"`
}

func (this *GameUserDaily) TableName() string {
    return "game_user_daily"
}

//table game_user_daily name and attributes defination.
const TABLE_GameUserDaily = "game_user_daily"
const COLUMN_GameUserDaily_ChannelId = "channel_id"
const COLUMN_GameUserDaily_Uid = "uid"
const COLUMN_GameUserDaily_TaxInteger = "tax_integer"
const COLUMN_GameUserDaily_TaxDecimals = "tax_decimals"
const COLUMN_GameUserDaily_ChipsInteger = "chips_integer"
const COLUMN_GameUserDaily_ChipsDecimals = "chips_decimals"
const COLUMN_GameUserDaily_WinloseInteger = "winlose_integer"
const COLUMN_GameUserDaily_WinloseDecimals = "winlose_decimals"
const COLUMN_GameUserDaily_Ctime = "ctime"
const COLUMN_GameUserDaily_Mtime = "mtime"
const ATTRIBUTE_GameUserDaily_ChannelId = "ChannelId"
const ATTRIBUTE_GameUserDaily_Uid = "Uid"
const ATTRIBUTE_GameUserDaily_TaxInteger = "TaxInteger"
const ATTRIBUTE_GameUserDaily_TaxDecimals = "TaxDecimals"
const ATTRIBUTE_GameUserDaily_ChipsInteger = "ChipsInteger"
const ATTRIBUTE_GameUserDaily_ChipsDecimals = "ChipsDecimals"
const ATTRIBUTE_GameUserDaily_WinloseInteger = "WinloseInteger"
const ATTRIBUTE_GameUserDaily_WinloseDecimals = "WinloseDecimals"
const ATTRIBUTE_GameUserDaily_Ctime = "Ctime"
const ATTRIBUTE_GameUserDaily_Mtime = "Mtime"

//auto_models_end
