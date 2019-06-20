package models

//auto_models_start
 type ProfitReportDaily struct{
	Id uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid uint64 `orm:"column(uid)" json:"uid,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	Bet int64 `orm:"column(bet)" json:"bet,omitempty"`
	TotalValidBet int64 `orm:"column(total_valid_bet)" json:"total_valid_bet,omitempty"`
	Profit int64 `orm:"column(profit)" json:"profit,omitempty"`
	Salary int64 `orm:"column(salary)" json:"salary,omitempty"`
	SelfDividend int64 `orm:"column(self_dividend)" json:"self_dividend,omitempty"`
	AgentDividend int64 `orm:"column(agent_dividend)" json:"agent_dividend,omitempty"`
	ResultDividend int64 `orm:"column(result_dividend)" json:"result_dividend,omitempty"`
	GameWithdrawAmount uint32 `orm:"column(game_withdraw_amount)" json:"game_withdraw_amount,omitempty"`
	GameRechargeAmount uint32 `orm:"column(game_recharge_amount)" json:"game_recharge_amount,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *ProfitReportDaily) TableName() string {
    return "profit_report_daily"
}

//table profit_report_daily name and attributes defination.
const TABLE_ProfitReportDaily = "profit_report_daily"
const COLUMN_ProfitReportDaily_Id = "id"
const COLUMN_ProfitReportDaily_Uid = "uid"
const COLUMN_ProfitReportDaily_ChannelId = "channel_id"
const COLUMN_ProfitReportDaily_Bet = "bet"
const COLUMN_ProfitReportDaily_TotalValidBet = "total_valid_bet"
const COLUMN_ProfitReportDaily_Profit = "profit"
const COLUMN_ProfitReportDaily_Salary = "salary"
const COLUMN_ProfitReportDaily_SelfDividend = "self_dividend"
const COLUMN_ProfitReportDaily_AgentDividend = "agent_dividend"
const COLUMN_ProfitReportDaily_ResultDividend = "result_dividend"
const COLUMN_ProfitReportDaily_GameWithdrawAmount = "game_withdraw_amount"
const COLUMN_ProfitReportDaily_GameRechargeAmount = "game_recharge_amount"
const COLUMN_ProfitReportDaily_Ctime = "ctime"
const ATTRIBUTE_ProfitReportDaily_Id = "Id"
const ATTRIBUTE_ProfitReportDaily_Uid = "Uid"
const ATTRIBUTE_ProfitReportDaily_ChannelId = "ChannelId"
const ATTRIBUTE_ProfitReportDaily_Bet = "Bet"
const ATTRIBUTE_ProfitReportDaily_TotalValidBet = "TotalValidBet"
const ATTRIBUTE_ProfitReportDaily_Profit = "Profit"
const ATTRIBUTE_ProfitReportDaily_Salary = "Salary"
const ATTRIBUTE_ProfitReportDaily_SelfDividend = "SelfDividend"
const ATTRIBUTE_ProfitReportDaily_AgentDividend = "AgentDividend"
const ATTRIBUTE_ProfitReportDaily_ResultDividend = "ResultDividend"
const ATTRIBUTE_ProfitReportDaily_GameWithdrawAmount = "GameWithdrawAmount"
const ATTRIBUTE_ProfitReportDaily_GameRechargeAmount = "GameRechargeAmount"
const ATTRIBUTE_ProfitReportDaily_Ctime = "Ctime"

//auto_models_end
