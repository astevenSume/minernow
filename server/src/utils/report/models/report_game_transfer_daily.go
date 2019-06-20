package models

//auto_models_start
 type ReportGameTransferDaily struct{
	Uid uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	Recharge int64 `orm:"column(recharge)" json:"recharge,omitempty"`
	Withdraw int64 `orm:"column(withdraw)" json:"withdraw,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *ReportGameTransferDaily) TableName() string {
    return "report_game_transfer_daily"
}

//table report_game_transfer_daily name and attributes defination.
const TABLE_ReportGameTransferDaily = "report_game_transfer_daily"
const COLUMN_ReportGameTransferDaily_Uid = "uid"
const COLUMN_ReportGameTransferDaily_ChannelId = "channel_id"
const COLUMN_ReportGameTransferDaily_Recharge = "recharge"
const COLUMN_ReportGameTransferDaily_Withdraw = "withdraw"
const COLUMN_ReportGameTransferDaily_Ctime = "ctime"
const ATTRIBUTE_ReportGameTransferDaily_Uid = "Uid"
const ATTRIBUTE_ReportGameTransferDaily_ChannelId = "ChannelId"
const ATTRIBUTE_ReportGameTransferDaily_Recharge = "Recharge"
const ATTRIBUTE_ReportGameTransferDaily_Withdraw = "Withdraw"
const ATTRIBUTE_ReportGameTransferDaily_Ctime = "Ctime"

//auto_models_end
