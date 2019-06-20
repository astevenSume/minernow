package models

//auto_models_start
 type ChannelDaily struct{
	ChannelId uint32 `orm:"column(channel_id);pk" json:"channel_id,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	WinLoseMoneyInteger int32 `orm:"column(win_lose_money_integer)" json:"win_lose_money_integer,omitempty"`
	WinLoseMoneyDecimals int32 `orm:"column(win_lose_money_decimals)" json:"win_lose_money_decimals,omitempty"`
	ChipsInteger int32 `orm:"column(chips_integer)" json:"chips_integer,omitempty"`
	ChipsDecimals int32 `orm:"column(chips_decimals)" json:"chips_decimals,omitempty"`
	Mtime int64 `orm:"column(mtime)" json:"mtime,omitempty"`
}

func (this *ChannelDaily) TableName() string {
    return "channel_daily"
}

//table channel_daily name and attributes defination.
const TABLE_ChannelDaily = "channel_daily"
const COLUMN_ChannelDaily_ChannelId = "channel_id"
const COLUMN_ChannelDaily_Ctime = "ctime"
const COLUMN_ChannelDaily_WinLoseMoneyInteger = "win_lose_money_integer"
const COLUMN_ChannelDaily_WinLoseMoneyDecimals = "win_lose_money_decimals"
const COLUMN_ChannelDaily_ChipsInteger = "chips_integer"
const COLUMN_ChannelDaily_ChipsDecimals = "chips_decimals"
const COLUMN_ChannelDaily_Mtime = "mtime"
const ATTRIBUTE_ChannelDaily_ChannelId = "ChannelId"
const ATTRIBUTE_ChannelDaily_Ctime = "Ctime"
const ATTRIBUTE_ChannelDaily_WinLoseMoneyInteger = "WinLoseMoneyInteger"
const ATTRIBUTE_ChannelDaily_WinLoseMoneyDecimals = "WinLoseMoneyDecimals"
const ATTRIBUTE_ChannelDaily_ChipsInteger = "ChipsInteger"
const ATTRIBUTE_ChannelDaily_ChipsDecimals = "ChipsDecimals"
const ATTRIBUTE_ChannelDaily_Mtime = "Mtime"

//auto_models_end
