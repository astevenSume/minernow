package models

//auto_models_start
type OtcStatAllPeople struct {
	Uid            uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	BuyOrder       uint32 `orm:"column(buy_order)" json:"buy_order,omitempty"`
	SellOrder      uint32 `orm:"column(sell_order)" json:"sell_order,omitempty"`
	BuyEusd        uint32 `orm:"column(buy_eusd)" json:"buy_eusd,omitempty"`
	SellEusd       uint32 `orm:"column(sell_eusd)" json:"sell_eusd,omitempty"`
	UsdtRecharge   uint32 `orm:"column(usdt_recharge)" json:"usdt_recharge,omitempty"`
	UsdtWithdrawal uint32 `orm:"column(usdt_withdrawal)" json:"usdt_withdrawal,omitempty"`
}

func (this *OtcStatAllPeople) TableName() string {
	return "otc_stat_all_people"
}

//table otc_stat_all_people name and attributes defination.
const TABLE_OtcStatAllPeople = "otc_stat_all_people"
const COLUMN_OtcStatAllPeople_Uid = "uid"
const COLUMN_OtcStatAllPeople_BuyOrder = "buy_order"
const COLUMN_OtcStatAllPeople_SellOrder = "sell_order"
const COLUMN_OtcStatAllPeople_BuyEusd = "buy_eusd"
const COLUMN_OtcStatAllPeople_SellEusd = "sell_eusd"
const COLUMN_OtcStatAllPeople_UsdtRecharge = "usdt_recharge"
const COLUMN_OtcStatAllPeople_UsdtWithdrawal = "usdt_withdrawal"
const ATTRIBUTE_OtcStatAllPeople_Uid = "Uid"
const ATTRIBUTE_OtcStatAllPeople_BuyOrder = "BuyOrder"
const ATTRIBUTE_OtcStatAllPeople_SellOrder = "SellOrder"
const ATTRIBUTE_OtcStatAllPeople_BuyEusd = "BuyEusd"
const ATTRIBUTE_OtcStatAllPeople_SellEusd = "SellEusd"
const ATTRIBUTE_OtcStatAllPeople_UsdtRecharge = "UsdtRecharge"
const ATTRIBUTE_OtcStatAllPeople_UsdtWithdrawal = "UsdtWithdrawal"

//auto_models_end
