package models

//auto_models_start
type ReportTeamGameTransferDaily struct {
	Uid          uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	ChannelId    uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	TeamRecharge int64  `orm:"column(team_recharge)" json:"team_recharge,omitempty"`
	TeamWithdraw int64  `orm:"column(team_withdraw)" json:"team_withdraw,omitempty"`
	Level        uint32 `orm:"column(level)" json:"level,omitempty"`
	Ctime        int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *ReportTeamGameTransferDaily) TableName() string {
	return "report_team_game_transfer_daily"
}

//table report_team_game_transfer_daily name and attributes defination.
const TABLE_ReportTeamGameTransferDaily = "report_team_game_transfer_daily"
const COLUMN_ReportTeamGameTransferDaily_Uid = "uid"
const COLUMN_ReportTeamGameTransferDaily_ChannelId = "channel_id"
const COLUMN_ReportTeamGameTransferDaily_TeamRecharge = "team_recharge"
const COLUMN_ReportTeamGameTransferDaily_TeamWithdraw = "team_withdraw"
const COLUMN_ReportTeamGameTransferDaily_Level = "level"
const COLUMN_ReportTeamGameTransferDaily_Ctime = "ctime"
const ATTRIBUTE_ReportTeamGameTransferDaily_Uid = "Uid"
const ATTRIBUTE_ReportTeamGameTransferDaily_ChannelId = "ChannelId"
const ATTRIBUTE_ReportTeamGameTransferDaily_TeamRecharge = "TeamRecharge"
const ATTRIBUTE_ReportTeamGameTransferDaily_TeamWithdraw = "TeamWithdraw"
const ATTRIBUTE_ReportTeamGameTransferDaily_Level = "Level"
const ATTRIBUTE_ReportTeamGameTransferDaily_Ctime = "Ctime"

//auto_models_end
