package models

//auto_models_start
 type ReportGameUserDaily struct{
	Uid uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	PUid uint64 `orm:"column(p_uid)" json:"p_uid,omitempty"`
	Level uint32 `orm:"column(level)" json:"level,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	Bet int64 `orm:"column(bet)" json:"bet,omitempty"`
	ValidBet int64 `orm:"column(valid_bet)" json:"valid_bet,omitempty"`
	TotalValidBet int64 `orm:"column(total_valid_bet)" json:"total_valid_bet,omitempty"`
	TotalBetNum int32 `orm:"column(total_bet_num)" json:"total_bet_num,omitempty"`
	Profit int64 `orm:"column(profit)" json:"profit,omitempty"`
	TotalProfit int64 `orm:"column(total_profit)" json:"total_profit,omitempty"`
	Salary int64 `orm:"column(salary)" json:"salary,omitempty"`
	TeamSalary int64 `orm:"column(team_salary)" json:"team_salary,omitempty"`
	Status uint8 `orm:"column(status)" json:"status,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *ReportGameUserDaily) TableName() string {
    return "report_game_user_daily"
}

//table report_game_user_daily name and attributes defination.
const TABLE_ReportGameUserDaily = "report_game_user_daily"
const COLUMN_ReportGameUserDaily_Uid = "uid"
const COLUMN_ReportGameUserDaily_PUid = "p_uid"
const COLUMN_ReportGameUserDaily_Level = "level"
const COLUMN_ReportGameUserDaily_ChannelId = "channel_id"
const COLUMN_ReportGameUserDaily_Bet = "bet"
const COLUMN_ReportGameUserDaily_ValidBet = "valid_bet"
const COLUMN_ReportGameUserDaily_TotalValidBet = "total_valid_bet"
const COLUMN_ReportGameUserDaily_TotalBetNum = "total_bet_num"
const COLUMN_ReportGameUserDaily_Profit = "profit"
const COLUMN_ReportGameUserDaily_TotalProfit = "total_profit"
const COLUMN_ReportGameUserDaily_Salary = "salary"
const COLUMN_ReportGameUserDaily_TeamSalary = "team_salary"
const COLUMN_ReportGameUserDaily_Status = "status"
const COLUMN_ReportGameUserDaily_Ctime = "ctime"
const ATTRIBUTE_ReportGameUserDaily_Uid = "Uid"
const ATTRIBUTE_ReportGameUserDaily_PUid = "PUid"
const ATTRIBUTE_ReportGameUserDaily_Level = "Level"
const ATTRIBUTE_ReportGameUserDaily_ChannelId = "ChannelId"
const ATTRIBUTE_ReportGameUserDaily_Bet = "Bet"
const ATTRIBUTE_ReportGameUserDaily_ValidBet = "ValidBet"
const ATTRIBUTE_ReportGameUserDaily_TotalValidBet = "TotalValidBet"
const ATTRIBUTE_ReportGameUserDaily_TotalBetNum = "TotalBetNum"
const ATTRIBUTE_ReportGameUserDaily_Profit = "Profit"
const ATTRIBUTE_ReportGameUserDaily_TotalProfit = "TotalProfit"
const ATTRIBUTE_ReportGameUserDaily_Salary = "Salary"
const ATTRIBUTE_ReportGameUserDaily_TeamSalary = "TeamSalary"
const ATTRIBUTE_ReportGameUserDaily_Status = "Status"
const ATTRIBUTE_ReportGameUserDaily_Ctime = "Ctime"

//auto_models_end
