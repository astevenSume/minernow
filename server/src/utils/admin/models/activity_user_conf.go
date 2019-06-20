package models

//auto_models_start
type ActivityUserConf struct {
	Id                 uint64 `orm:"column(id);pk" json:"id,omitempty"`
	PlayGameDay        int32  `orm:"column(play_game_day)" json:"play_game_day,omitempty"`
	BetAmount          int64  `orm:"column(bet_amount)" json:"bet_amount,omitempty"`
	EffectiveBetAmount int64  `orm:"column(effective_bet_amount)" json:"effective_bet_amount,omitempty"`
}

func (this *ActivityUserConf) TableName() string {
	return "activity_user_conf"
}

//table activity_user_conf name and attributes defination.
const TABLE_ActivityUserConf = "activity_user_conf"
const COLUMN_ActivityUserConf_Id = "id"
const COLUMN_ActivityUserConf_PlayGameDay = "play_game_day"
const COLUMN_ActivityUserConf_BetAmount = "bet_amount"
const COLUMN_ActivityUserConf_EffectiveBetAmount = "effective_bet_amount"
const ATTRIBUTE_ActivityUserConf_Id = "Id"
const ATTRIBUTE_ActivityUserConf_PlayGameDay = "PlayGameDay"
const ATTRIBUTE_ActivityUserConf_BetAmount = "BetAmount"
const ATTRIBUTE_ActivityUserConf_EffectiveBetAmount = "EffectiveBetAmount"

//auto_models_end
