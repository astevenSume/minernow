package models

//auto_models_start
 type GameUserMonthReport struct{
	Id uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Profit int64 `orm:"column(profit)" json:"profit,omitempty"`
	AgentsProfit int64 `orm:"column(agents_profit)" json:"agents_profit,omitempty"`
	ResultProfit int64 `orm:"column(result_profit)" json:"result_profit,omitempty"`
	BetAmount int64 `orm:"column(bet_amount)" json:"bet_amount,omitempty"`
	EffectiveBetAmount int64 `orm:"column(effective_bet_amount)" json:"effective_bet_amount,omitempty"`
	PlayGameDay int32 `orm:"column(play_game_day)" json:"play_game_day,omitempty"`
	IsActivityUser bool `orm:"column(is_activity_user)" json:"is_activity_user,omitempty"`
	AgentLevel uint32 `orm:"column(agent_level)" json:"agent_level,omitempty"`
	UpAgentUid uint64 `orm:"column(up_agent_uid)" json:"up_agent_uid,omitempty"`
	ActivityAgentNum int32 `orm:"column(activity_agent_num)" json:"activity_agent_num,omitempty"`
}

func (this *GameUserMonthReport) TableName() string {
    return "game_user_month_report"
}

//table game_user_month_report name and attributes defination.
const TABLE_GameUserMonthReport = "game_user_month_report"
const COLUMN_GameUserMonthReport_Id = "id"
const COLUMN_GameUserMonthReport_Uid = "uid"
const COLUMN_GameUserMonthReport_Ctime = "ctime"
const COLUMN_GameUserMonthReport_Profit = "profit"
const COLUMN_GameUserMonthReport_AgentsProfit = "agents_profit"
const COLUMN_GameUserMonthReport_ResultProfit = "result_profit"
const COLUMN_GameUserMonthReport_BetAmount = "bet_amount"
const COLUMN_GameUserMonthReport_EffectiveBetAmount = "effective_bet_amount"
const COLUMN_GameUserMonthReport_PlayGameDay = "play_game_day"
const COLUMN_GameUserMonthReport_IsActivityUser = "is_activity_user"
const COLUMN_GameUserMonthReport_AgentLevel = "agent_level"
const COLUMN_GameUserMonthReport_UpAgentUid = "up_agent_uid"
const COLUMN_GameUserMonthReport_ActivityAgentNum = "activity_agent_num"
const ATTRIBUTE_GameUserMonthReport_Id = "Id"
const ATTRIBUTE_GameUserMonthReport_Uid = "Uid"
const ATTRIBUTE_GameUserMonthReport_Ctime = "Ctime"
const ATTRIBUTE_GameUserMonthReport_Profit = "Profit"
const ATTRIBUTE_GameUserMonthReport_AgentsProfit = "AgentsProfit"
const ATTRIBUTE_GameUserMonthReport_ResultProfit = "ResultProfit"
const ATTRIBUTE_GameUserMonthReport_BetAmount = "BetAmount"
const ATTRIBUTE_GameUserMonthReport_EffectiveBetAmount = "EffectiveBetAmount"
const ATTRIBUTE_GameUserMonthReport_PlayGameDay = "PlayGameDay"
const ATTRIBUTE_GameUserMonthReport_IsActivityUser = "IsActivityUser"
const ATTRIBUTE_GameUserMonthReport_AgentLevel = "AgentLevel"
const ATTRIBUTE_GameUserMonthReport_UpAgentUid = "UpAgentUid"
const ATTRIBUTE_GameUserMonthReport_ActivityAgentNum = "ActivityAgentNum"

//auto_models_end
