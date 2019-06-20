package models

//auto_models_start
type OtcStat struct {
	Id             uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Date           uint32 `orm:"column(date)" json:"date,omitempty"`
	NumLogin       uint32 `orm:"column(num_login)" json:"num_login,omitempty"`
	NumUserNew     uint32 `orm:"column(num_user_new)" json:"num_user_new,omitempty"`
	NumOrder       uint32 `orm:"column(num_order)" json:"num_order,omitempty"`
	NumOrderDeal   uint32 `orm:"column(num_order_deal)" json:"num_order_deal,omitempty"`
	NumOrderBuy    uint32 `orm:"column(num_order_buy)" json:"num_order_buy,omitempty"`
	NumOrderSell   uint32 `orm:"column(num_order_sell)" json:"num_order_sell,omitempty"`
	NumFunds       uint32 `orm:"column(num_funds)" json:"num_funds,omitempty"`
	NumAmount      uint32 `orm:"column(num_amount)" json:"num_amount,omitempty"`
	NumAmountBuy   uint32 `orm:"column(num_amount_buy)" json:"num_amount_buy,omitempty"`
	NumAmountSell  uint32 `orm:"column(num_amount_sell)" json:"num_amount_sell,omitempty"`
	NumFeeBuy      uint32 `orm:"column(num_fee_buy)" json:"num_fee_buy,omitempty"`
	NumFeeSell     uint32 `orm:"column(num_fee_sell)" json:"num_fee_sell,omitempty"`
	GameRecharge   uint32 `orm:"column(game_recharge)" json:"game_recharge,omitempty"`
	GameWithdrawal uint32 `orm:"column(game_withdrawal)" json:"game_withdrawal,omitempty"`
	UsdtRecharge   uint32 `orm:"column(usdt_recharge)" json:"usdt_recharge,omitempty"`
	UsdtWithdrawal uint32 `orm:"column(usdt_withdrawal)" json:"usdt_withdrawal,omitempty"`
	UsdtFee        uint32 `orm:"column(usdt_fee)" json:"usdt_fee,omitempty"`
}

func (this *OtcStat) TableName() string {
	return "otc_stat"
}

//table otc_stat name and attributes defination.
const TABLE_OtcStat = "otc_stat"
const COLUMN_OtcStat_Id = "id"
const COLUMN_OtcStat_Date = "date"
const COLUMN_OtcStat_NumLogin = "num_login"
const COLUMN_OtcStat_NumUserNew = "num_user_new"
const COLUMN_OtcStat_NumOrder = "num_order"
const COLUMN_OtcStat_NumOrderDeal = "num_order_deal"
const COLUMN_OtcStat_NumOrderBuy = "num_order_buy"
const COLUMN_OtcStat_NumOrderSell = "num_order_sell"
const COLUMN_OtcStat_NumFunds = "num_funds"
const COLUMN_OtcStat_NumAmount = "num_amount"
const COLUMN_OtcStat_NumAmountBuy = "num_amount_buy"
const COLUMN_OtcStat_NumAmountSell = "num_amount_sell"
const COLUMN_OtcStat_NumFeeBuy = "num_fee_buy"
const COLUMN_OtcStat_NumFeeSell = "num_fee_sell"
const COLUMN_OtcStat_GameRecharge = "game_recharge"
const COLUMN_OtcStat_GameWithdrawal = "game_withdrawal"
const COLUMN_OtcStat_UsdtRecharge = "usdt_recharge"
const COLUMN_OtcStat_UsdtWithdrawal = "usdt_withdrawal"
const COLUMN_OtcStat_UsdtFee = "usdt_fee"
const ATTRIBUTE_OtcStat_Id = "Id"
const ATTRIBUTE_OtcStat_Date = "Date"
const ATTRIBUTE_OtcStat_NumLogin = "NumLogin"
const ATTRIBUTE_OtcStat_NumUserNew = "NumUserNew"
const ATTRIBUTE_OtcStat_NumOrder = "NumOrder"
const ATTRIBUTE_OtcStat_NumOrderDeal = "NumOrderDeal"
const ATTRIBUTE_OtcStat_NumOrderBuy = "NumOrderBuy"
const ATTRIBUTE_OtcStat_NumOrderSell = "NumOrderSell"
const ATTRIBUTE_OtcStat_NumFunds = "NumFunds"
const ATTRIBUTE_OtcStat_NumAmount = "NumAmount"
const ATTRIBUTE_OtcStat_NumAmountBuy = "NumAmountBuy"
const ATTRIBUTE_OtcStat_NumAmountSell = "NumAmountSell"
const ATTRIBUTE_OtcStat_NumFeeBuy = "NumFeeBuy"
const ATTRIBUTE_OtcStat_NumFeeSell = "NumFeeSell"
const ATTRIBUTE_OtcStat_GameRecharge = "GameRecharge"
const ATTRIBUTE_OtcStat_GameWithdrawal = "GameWithdrawal"
const ATTRIBUTE_OtcStat_UsdtRecharge = "UsdtRecharge"
const ATTRIBUTE_OtcStat_UsdtWithdrawal = "UsdtWithdrawal"
const ATTRIBUTE_OtcStat_UsdtFee = "UsdtFee"

//auto_models_end
