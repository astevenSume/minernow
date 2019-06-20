package models

//auto_models_start
 type ReportGameRecordAg struct{
	Id uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Account string `orm:"column(account);size(50)" json:"account,omitempty"`
	GameType string `orm:"column(game_type);size(50)" json:"game_type,omitempty"`
	GameName string `orm:"column(game_name);size(50)" json:"game_name,omitempty"`
	OrderId string `orm:"column(order_id);size(50)" json:"order_id,omitempty"`
	TableId string `orm:"column(table_id);size(50)" json:"table_id,omitempty"`
	Bet int64 `orm:"column(bet)" json:"bet,omitempty"`
	ValidBet int64 `orm:"column(valid_bet)" json:"valid_bet,omitempty"`
	Profit int64 `orm:"column(profit)" json:"profit,omitempty"`
	BetTime string `orm:"column(bet_time);size(32)" json:"bet_time,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Uid uint64 `orm:"column(uid)" json:"uid,omitempty"`
}

func (this *ReportGameRecordAg) TableName() string {
    return "report_game_record_ag"
}

//table report_game_record_ag name and attributes defination.
const TABLE_ReportGameRecordAg = "report_game_record_ag"
const COLUMN_ReportGameRecordAg_Id = "id"
const COLUMN_ReportGameRecordAg_Account = "account"
const COLUMN_ReportGameRecordAg_GameType = "game_type"
const COLUMN_ReportGameRecordAg_GameName = "game_name"
const COLUMN_ReportGameRecordAg_OrderId = "order_id"
const COLUMN_ReportGameRecordAg_TableId = "table_id"
const COLUMN_ReportGameRecordAg_Bet = "bet"
const COLUMN_ReportGameRecordAg_ValidBet = "valid_bet"
const COLUMN_ReportGameRecordAg_Profit = "profit"
const COLUMN_ReportGameRecordAg_BetTime = "bet_time"
const COLUMN_ReportGameRecordAg_Ctime = "ctime"
const COLUMN_ReportGameRecordAg_Uid = "uid"
const ATTRIBUTE_ReportGameRecordAg_Id = "Id"
const ATTRIBUTE_ReportGameRecordAg_Account = "Account"
const ATTRIBUTE_ReportGameRecordAg_GameType = "GameType"
const ATTRIBUTE_ReportGameRecordAg_GameName = "GameName"
const ATTRIBUTE_ReportGameRecordAg_OrderId = "OrderId"
const ATTRIBUTE_ReportGameRecordAg_TableId = "TableId"
const ATTRIBUTE_ReportGameRecordAg_Bet = "Bet"
const ATTRIBUTE_ReportGameRecordAg_ValidBet = "ValidBet"
const ATTRIBUTE_ReportGameRecordAg_Profit = "Profit"
const ATTRIBUTE_ReportGameRecordAg_BetTime = "BetTime"
const ATTRIBUTE_ReportGameRecordAg_Ctime = "Ctime"
const ATTRIBUTE_ReportGameRecordAg_Uid = "Uid"

//auto_models_end
