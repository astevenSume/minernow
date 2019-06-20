package models

//auto_models_start
 type ReportStatisticGameAll struct{
	Id uint64 `orm:"column(id);pk" json:"id,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	NewerNums int64 `orm:"column(newer_nums)" json:"newer_nums,omitempty"`
	Bet int64 `orm:"column(bet)" json:"bet,omitempty"`
	ValidBet int64 `orm:"column(valid_bet)" json:"valid_bet,omitempty"`
	Profit int64 `orm:"column(profit)" json:"profit,omitempty"`
	Revenue int64 `orm:"column(revenue)" json:"revenue,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Note string `orm:"column(note);size(255)" json:"note,omitempty"`
	GameId uint32 `orm:"column(game_id)" json:"game_id,omitempty"`
}

func (this *ReportStatisticGameAll) TableName() string {
    return "report_statistic_game_all"
}

//table report_statistic_game_all name and attributes defination.
const TABLE_ReportStatisticGameAll = "report_statistic_game_all"
const COLUMN_ReportStatisticGameAll_Id = "id"
const COLUMN_ReportStatisticGameAll_ChannelId = "channel_id"
const COLUMN_ReportStatisticGameAll_NewerNums = "newer_nums"
const COLUMN_ReportStatisticGameAll_Bet = "bet"
const COLUMN_ReportStatisticGameAll_ValidBet = "valid_bet"
const COLUMN_ReportStatisticGameAll_Profit = "profit"
const COLUMN_ReportStatisticGameAll_Revenue = "revenue"
const COLUMN_ReportStatisticGameAll_Ctime = "ctime"
const COLUMN_ReportStatisticGameAll_Note = "note"
const COLUMN_ReportStatisticGameAll_GameId = "game_id"
const ATTRIBUTE_ReportStatisticGameAll_Id = "Id"
const ATTRIBUTE_ReportStatisticGameAll_ChannelId = "ChannelId"
const ATTRIBUTE_ReportStatisticGameAll_NewerNums = "NewerNums"
const ATTRIBUTE_ReportStatisticGameAll_Bet = "Bet"
const ATTRIBUTE_ReportStatisticGameAll_ValidBet = "ValidBet"
const ATTRIBUTE_ReportStatisticGameAll_Profit = "Profit"
const ATTRIBUTE_ReportStatisticGameAll_Revenue = "Revenue"
const ATTRIBUTE_ReportStatisticGameAll_Ctime = "Ctime"
const ATTRIBUTE_ReportStatisticGameAll_Note = "Note"
const ATTRIBUTE_ReportStatisticGameAll_GameId = "GameId"

//auto_models_end
