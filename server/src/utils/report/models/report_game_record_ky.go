package models

//auto_models_start
 type ReportGameRecordKy struct{
	Id uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Account string `orm:"column(account);size(50)" json:"account,omitempty"`
	GameId string `orm:"column(game_id);size(50)" json:"game_id,omitempty"`
	GameName string `orm:"column(game_name);size(50)" json:"game_name,omitempty"`
	ServerId int32 `orm:"column(server_id)" json:"server_id,omitempty"`
	KindId string `orm:"column(kind_id);size(50)" json:"kind_id,omitempty"`
	TableId int32 `orm:"column(table_id)" json:"table_id,omitempty"`
	ChairId int32 `orm:"column(chair_id)" json:"chair_id,omitempty"`
	Bet int64 `orm:"column(bet)" json:"bet,omitempty"`
	ValidBet int64 `orm:"column(valid_bet)" json:"valid_bet,omitempty"`
	Profit int64 `orm:"column(profit)" json:"profit,omitempty"`
	Revenue int64 `orm:"column(revenue)" json:"revenue,omitempty"`
	StartTime string `orm:"column(start_time);size(32)" json:"start_time,omitempty"`
	EndTime string `orm:"column(end_time);size(32)" json:"end_time,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Uid uint64 `orm:"column(uid)" json:"uid,omitempty"`
}

func (this *ReportGameRecordKy) TableName() string {
    return "report_game_record_ky"
}

//table report_game_record_ky name and attributes defination.
const TABLE_ReportGameRecordKy = "report_game_record_ky"
const COLUMN_ReportGameRecordKy_Id = "id"
const COLUMN_ReportGameRecordKy_Account = "account"
const COLUMN_ReportGameRecordKy_GameId = "game_id"
const COLUMN_ReportGameRecordKy_GameName = "game_name"
const COLUMN_ReportGameRecordKy_ServerId = "server_id"
const COLUMN_ReportGameRecordKy_KindId = "kind_id"
const COLUMN_ReportGameRecordKy_TableId = "table_id"
const COLUMN_ReportGameRecordKy_ChairId = "chair_id"
const COLUMN_ReportGameRecordKy_Bet = "bet"
const COLUMN_ReportGameRecordKy_ValidBet = "valid_bet"
const COLUMN_ReportGameRecordKy_Profit = "profit"
const COLUMN_ReportGameRecordKy_Revenue = "revenue"
const COLUMN_ReportGameRecordKy_StartTime = "start_time"
const COLUMN_ReportGameRecordKy_EndTime = "end_time"
const COLUMN_ReportGameRecordKy_Ctime = "ctime"
const COLUMN_ReportGameRecordKy_Uid = "uid"
const ATTRIBUTE_ReportGameRecordKy_Id = "Id"
const ATTRIBUTE_ReportGameRecordKy_Account = "Account"
const ATTRIBUTE_ReportGameRecordKy_GameId = "GameId"
const ATTRIBUTE_ReportGameRecordKy_GameName = "GameName"
const ATTRIBUTE_ReportGameRecordKy_ServerId = "ServerId"
const ATTRIBUTE_ReportGameRecordKy_KindId = "KindId"
const ATTRIBUTE_ReportGameRecordKy_TableId = "TableId"
const ATTRIBUTE_ReportGameRecordKy_ChairId = "ChairId"
const ATTRIBUTE_ReportGameRecordKy_Bet = "Bet"
const ATTRIBUTE_ReportGameRecordKy_ValidBet = "ValidBet"
const ATTRIBUTE_ReportGameRecordKy_Profit = "Profit"
const ATTRIBUTE_ReportGameRecordKy_Revenue = "Revenue"
const ATTRIBUTE_ReportGameRecordKy_StartTime = "StartTime"
const ATTRIBUTE_ReportGameRecordKy_EndTime = "EndTime"
const ATTRIBUTE_ReportGameRecordKy_Ctime = "Ctime"
const ATTRIBUTE_ReportGameRecordKy_Uid = "Uid"

//auto_models_end
