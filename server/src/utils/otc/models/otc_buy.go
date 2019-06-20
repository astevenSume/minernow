package models

//auto_models_start
type OtcBuy struct {
	Uid              uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Available        int64  `orm:"column(available)" json:"available,omitempty"`
	Frozen           uint64 `orm:"column(frozen)" json:"frozen,omitempty"`
	Bought           uint64 `orm:"column(bought)" json:"bought,omitempty"`
	LowerLimitWechat int64  `orm:"column(lower_limit_wechat)" json:"lower_limit_wechat,omitempty"`
	UpperLimitWechat int64  `orm:"column(upper_limit_wechat)" json:"upper_limit_wechat,omitempty"`
	LowerLimitBank   int64  `orm:"column(lower_limit_bank)" json:"lower_limit_bank,omitempty"`
	UpperLimitBank   int64  `orm:"column(upper_limit_bank)" json:"upper_limit_bank,omitempty"`
	LowerLimitAli    int64  `orm:"column(lower_limit_ali)" json:"lower_limit_ali,omitempty"`
	UpperLimitAli    int64  `orm:"column(upper_limit_ali)" json:"upper_limit_ali,omitempty"`
	PayType          uint8  `orm:"column(pay_type)" json:"pay_type,omitempty"`
	Ctime            int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *OtcBuy) TableName() string {
	return "otc_buy"
}

//table otc_buy name and attributes defination.
const TABLE_OtcBuy = "otc_buy"
const COLUMN_OtcBuy_Uid = "uid"
const COLUMN_OtcBuy_Available = "available"
const COLUMN_OtcBuy_Frozen = "frozen"
const COLUMN_OtcBuy_Bought = "bought"
const COLUMN_OtcBuy_LowerLimitWechat = "lower_limit_wechat"
const COLUMN_OtcBuy_UpperLimitWechat = "upper_limit_wechat"
const COLUMN_OtcBuy_LowerLimitBank = "lower_limit_bank"
const COLUMN_OtcBuy_UpperLimitBank = "upper_limit_bank"
const COLUMN_OtcBuy_LowerLimitAli = "lower_limit_ali"
const COLUMN_OtcBuy_UpperLimitAli = "upper_limit_ali"
const COLUMN_OtcBuy_PayType = "pay_type"
const COLUMN_OtcBuy_Ctime = "ctime"
const ATTRIBUTE_OtcBuy_Uid = "Uid"
const ATTRIBUTE_OtcBuy_Available = "Available"
const ATTRIBUTE_OtcBuy_Frozen = "Frozen"
const ATTRIBUTE_OtcBuy_Bought = "Bought"
const ATTRIBUTE_OtcBuy_LowerLimitWechat = "LowerLimitWechat"
const ATTRIBUTE_OtcBuy_UpperLimitWechat = "UpperLimitWechat"
const ATTRIBUTE_OtcBuy_LowerLimitBank = "LowerLimitBank"
const ATTRIBUTE_OtcBuy_UpperLimitBank = "UpperLimitBank"
const ATTRIBUTE_OtcBuy_LowerLimitAli = "LowerLimitAli"
const ATTRIBUTE_OtcBuy_UpperLimitAli = "UpperLimitAli"
const ATTRIBUTE_OtcBuy_PayType = "PayType"
const ATTRIBUTE_OtcBuy_Ctime = "Ctime"

//auto_models_end
