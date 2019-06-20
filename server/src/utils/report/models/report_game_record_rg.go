package models

//auto_models_start
 type ReportGameRecordRg struct{
	Id uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Account string `orm:"column(account);size(50)" json:"account,omitempty"`
	GameNameID string `orm:"column(game_name_id);size(50)" json:"game_name_id,omitempty"`
	GameName string `orm:"column(game_name);size(50)" json:"game_name,omitempty"`
	GameKindName string `orm:"column(game_kind_name);size(50)" json:"game_kind_name,omitempty"`
	OrderId string `orm:"column(order_id);size(50)" json:"order_id,omitempty"`
	OpenDate string `orm:"column(open_date);size(32)" json:"open_date,omitempty"`
	PeriodName string `orm:"column(period_name);size(50)" json:"period_name,omitempty"`
	OpenNumber string `orm:"column(open_number);size(50)" json:"open_number,omitempty"`
	Status uint8 `orm:"column(status)" json:"status,omitempty"`
	Bet int64 `orm:"column(bet)" json:"bet,omitempty"`
	ValidBet int64 `orm:"column(valid_bet)" json:"valid_bet,omitempty"`
	Profit int64 `orm:"column(profit)" json:"profit,omitempty"`
	BetTime string `orm:"column(bet_time);size(32)" json:"bet_time,omitempty"`
	BetContent string `orm:"column(bet_content);size(50)" json:"bet_content,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Uid uint64 `orm:"column(uid)" json:"uid,omitempty"`
}

func (this *ReportGameRecordRg) TableName() string {
    return "report_game_record_rg"
}

//table report_game_record_rg name and attributes defination.
const TABLE_ReportGameRecordRg = "report_game_record_rg"
const COLUMN_ReportGameRecordRg_Id = "id"
const COLUMN_ReportGameRecordRg_Account = "account"
const COLUMN_ReportGameRecordRg_GameNameID = "game_name_id"
const COLUMN_ReportGameRecordRg_GameName = "game_name"
const COLUMN_ReportGameRecordRg_GameKindName = "game_kind_name"
const COLUMN_ReportGameRecordRg_OrderId = "order_id"
const COLUMN_ReportGameRecordRg_OpenDate = "open_date"
const COLUMN_ReportGameRecordRg_PeriodName = "period_name"
const COLUMN_ReportGameRecordRg_OpenNumber = "open_number"
const COLUMN_ReportGameRecordRg_Status = "status"
const COLUMN_ReportGameRecordRg_Bet = "bet"
const COLUMN_ReportGameRecordRg_ValidBet = "valid_bet"
const COLUMN_ReportGameRecordRg_Profit = "profit"
const COLUMN_ReportGameRecordRg_BetTime = "bet_time"
const COLUMN_ReportGameRecordRg_BetContent = "bet_content"
const COLUMN_ReportGameRecordRg_Ctime = "ctime"
const COLUMN_ReportGameRecordRg_Uid = "uid"
const ATTRIBUTE_ReportGameRecordRg_Id = "Id"
const ATTRIBUTE_ReportGameRecordRg_Account = "Account"
const ATTRIBUTE_ReportGameRecordRg_GameNameID = "GameNameID"
const ATTRIBUTE_ReportGameRecordRg_GameName = "GameName"
const ATTRIBUTE_ReportGameRecordRg_GameKindName = "GameKindName"
const ATTRIBUTE_ReportGameRecordRg_OrderId = "OrderId"
const ATTRIBUTE_ReportGameRecordRg_OpenDate = "OpenDate"
const ATTRIBUTE_ReportGameRecordRg_PeriodName = "PeriodName"
const ATTRIBUTE_ReportGameRecordRg_OpenNumber = "OpenNumber"
const ATTRIBUTE_ReportGameRecordRg_Status = "Status"
const ATTRIBUTE_ReportGameRecordRg_Bet = "Bet"
const ATTRIBUTE_ReportGameRecordRg_ValidBet = "ValidBet"
const ATTRIBUTE_ReportGameRecordRg_Profit = "Profit"
const ATTRIBUTE_ReportGameRecordRg_BetTime = "BetTime"
const ATTRIBUTE_ReportGameRecordRg_BetContent = "BetContent"
const ATTRIBUTE_ReportGameRecordRg_Ctime = "Ctime"
const ATTRIBUTE_ReportGameRecordRg_Uid = "Uid"

//auto_models_end
