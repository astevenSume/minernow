package models

//auto_models_start
type EosOtc struct {
	Uid               uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Account           string `orm:"column(account);size(100)" json:"account,omitempty"`
	Status            int8   `orm:"column(status)" json:"status,omitempty"`
	Available         int64  `orm:"column(available)" json:"available,omitempty"`
	Trade             int64  `orm:"column(trade)" json:"trade,omitempty"`
	Transfer          int64  `orm:"column(transfer)" json:"transfer,omitempty"`
	SellState         string `orm:"column(sell_state);size(200)" json:"sell_state,omitempty"`
	SellPayType       uint8  `orm:"column(sell_pay_type)" json:"sell_pay_type,omitempty"`
	SellAble          bool   `orm:"column(sell_able)" json:"sell_able,omitempty"`
	SellRmbDay        int64  `orm:"column(sell_rmb_day)" json:"sell_rmb_day,omitempty"`
	SellRmbToday      int64  `orm:"column(sell_rmb_today)" json:"sell_rmb_today,omitempty"`
	SellRmbLowerLimit int64  `orm:"column(sell_rmb_lower_limit)" json:"sell_rmb_lower_limit,omitempty"`
	SellUTime         int64  `orm:"column(sell_utime)" json:"sell_utime,omitempty"`
	BuyAble           bool   `orm:"column(buy_able)" json:"buy_able,omitempty"`
	BuyRmbDay         int64  `orm:"column(buy_rmb_day)" json:"buy_rmb_day,omitempty"`
	BuyRmbToday       int64  `orm:"column(buy_rmb_today)" json:"buy_rmb_today,omitempty"`
	BuyRmbLowerLimit  int64  `orm:"column(buy_rmb_lower_limit)" json:"buy_rmb_lower_limit,omitempty"`
	BuyUTime          int64  `orm:"column(buy_utime)" json:"buy_utime,omitempty"`
	BuyState          string `orm:"column(buy_state);size(200)" json:"buy_state,omitempty"`
	Ctime             int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime             int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *EosOtc) TableName() string {
	return "eos_otc"
}

//table eos_otc name and attributes defination.
const TABLE_EosOtc = "eos_otc"
const COLUMN_EosOtc_Uid = "uid"
const COLUMN_EosOtc_Account = "account"
const COLUMN_EosOtc_Status = "status"
const COLUMN_EosOtc_Available = "available"
const COLUMN_EosOtc_Trade = "trade"
const COLUMN_EosOtc_Transfer = "transfer"
const COLUMN_EosOtc_SellState = "sell_state"
const COLUMN_EosOtc_SellPayType = "sell_pay_type"
const COLUMN_EosOtc_SellAble = "sell_able"
const COLUMN_EosOtc_SellRmbDay = "sell_rmb_day"
const COLUMN_EosOtc_SellRmbToday = "sell_rmb_today"
const COLUMN_EosOtc_SellRmbLowerLimit = "sell_rmb_lower_limit"
const COLUMN_EosOtc_SellUTime = "sell_utime"
const COLUMN_EosOtc_BuyAble = "buy_able"
const COLUMN_EosOtc_BuyRmbDay = "buy_rmb_day"
const COLUMN_EosOtc_BuyRmbToday = "buy_rmb_today"
const COLUMN_EosOtc_BuyRmbLowerLimit = "buy_rmb_lower_limit"
const COLUMN_EosOtc_BuyUTime = "buy_utime"
const COLUMN_EosOtc_BuyState = "buy_state"
const COLUMN_EosOtc_Ctime = "ctime"
const COLUMN_EosOtc_Utime = "utime"
const ATTRIBUTE_EosOtc_Uid = "Uid"
const ATTRIBUTE_EosOtc_Account = "Account"
const ATTRIBUTE_EosOtc_Status = "Status"
const ATTRIBUTE_EosOtc_Available = "Available"
const ATTRIBUTE_EosOtc_Trade = "Trade"
const ATTRIBUTE_EosOtc_Transfer = "Transfer"
const ATTRIBUTE_EosOtc_SellState = "SellState"
const ATTRIBUTE_EosOtc_SellPayType = "SellPayType"
const ATTRIBUTE_EosOtc_SellAble = "SellAble"
const ATTRIBUTE_EosOtc_SellRmbDay = "SellRmbDay"
const ATTRIBUTE_EosOtc_SellRmbToday = "SellRmbToday"
const ATTRIBUTE_EosOtc_SellRmbLowerLimit = "SellRmbLowerLimit"
const ATTRIBUTE_EosOtc_SellUTime = "SellUTime"
const ATTRIBUTE_EosOtc_BuyAble = "BuyAble"
const ATTRIBUTE_EosOtc_BuyRmbDay = "BuyRmbDay"
const ATTRIBUTE_EosOtc_BuyRmbToday = "BuyRmbToday"
const ATTRIBUTE_EosOtc_BuyRmbLowerLimit = "BuyRmbLowerLimit"
const ATTRIBUTE_EosOtc_BuyUTime = "BuyUTime"
const ATTRIBUTE_EosOtc_BuyState = "BuyState"
const ATTRIBUTE_EosOtc_Ctime = "Ctime"
const ATTRIBUTE_EosOtc_Utime = "Utime"

//auto_models_end
