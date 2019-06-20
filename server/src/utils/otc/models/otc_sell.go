package models

//auto_models_start
type OtcSell struct {
	Uid        uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Available  int64  `orm:"column(available)" json:"available,omitempty"`
	Frozen     int64  `orm:"column(frozen)" json:"frozen,omitempty"`
	Sold       int64  `orm:"column(sold)" json:"sold,omitempty"`
	LowerLimit int64  `orm:"column(lower_limit)" json:"lower_limit,omitempty"`
	UpperLimit int64  `orm:"column(upper_limit)" json:"upper_limit,omitempty"`
	PayType    uint8  `orm:"column(pay_type)" json:"pay_type,omitempty"`
	Ctime      int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *OtcSell) TableName() string {
	return "otc_sell"
}

//table otc_sell name and attributes defination.
const TABLE_OtcSell = "otc_sell"
const COLUMN_OtcSell_Uid = "uid"
const COLUMN_OtcSell_Available = "available"
const COLUMN_OtcSell_Frozen = "frozen"
const COLUMN_OtcSell_Sold = "sold"
const COLUMN_OtcSell_LowerLimit = "lower_limit"
const COLUMN_OtcSell_UpperLimit = "upper_limit"
const COLUMN_OtcSell_PayType = "pay_type"
const COLUMN_OtcSell_Ctime = "ctime"
const ATTRIBUTE_OtcSell_Uid = "Uid"
const ATTRIBUTE_OtcSell_Available = "Available"
const ATTRIBUTE_OtcSell_Frozen = "Frozen"
const ATTRIBUTE_OtcSell_Sold = "Sold"
const ATTRIBUTE_OtcSell_LowerLimit = "LowerLimit"
const ATTRIBUTE_OtcSell_UpperLimit = "UpperLimit"
const ATTRIBUTE_OtcSell_PayType = "PayType"
const ATTRIBUTE_OtcSell_Ctime = "Ctime"

//auto_models_end
