package models

//auto_models_start
type EosOtcReport struct {
	Id              uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid             uint64 `orm:"column(uid)" json:"uid,omitempty"`
	TotalOrderNum   int64  `orm:"column(total_order_num)" json:"total_order_num,omitempty"`
	SuccessOrderNum int64  `orm:"column(success_order_num)" json:"success_order_num,omitempty"`
	FailOrderNum    int64  `orm:"column(fail_order_num)" json:"fail_order_num,omitempty"`
	BuyEusdNum      int64  `orm:"column(buy_eusd_num)" json:"buy_eusd_num,omitempty"`
	SellEusdNum     int64  `orm:"column(sell_eusd_num)" json:"sell_eusd_num,omitempty"`
	Date            int32  `orm:"column(date)" json:"date,omitempty"`
}

func (this *EosOtcReport) TableName() string {
	return "eos_otc_report"
}

//table eos_otc_report name and attributes defination.
const TABLE_EosOtcReport = "eos_otc_report"
const COLUMN_EosOtcReport_Id = "id"
const COLUMN_EosOtcReport_Uid = "uid"
const COLUMN_EosOtcReport_TotalOrderNum = "total_order_num"
const COLUMN_EosOtcReport_SuccessOrderNum = "success_order_num"
const COLUMN_EosOtcReport_FailOrderNum = "fail_order_num"
const COLUMN_EosOtcReport_BuyEusdNum = "buy_eusd_num"
const COLUMN_EosOtcReport_SellEusdNum = "sell_eusd_num"
const COLUMN_EosOtcReport_Date = "date"
const ATTRIBUTE_EosOtcReport_Id = "Id"
const ATTRIBUTE_EosOtcReport_Uid = "Uid"
const ATTRIBUTE_EosOtcReport_TotalOrderNum = "TotalOrderNum"
const ATTRIBUTE_EosOtcReport_SuccessOrderNum = "SuccessOrderNum"
const ATTRIBUTE_EosOtcReport_FailOrderNum = "FailOrderNum"
const ATTRIBUTE_EosOtcReport_BuyEusdNum = "BuyEusdNum"
const ATTRIBUTE_EosOtcReport_SellEusdNum = "SellEusdNum"
const ATTRIBUTE_EosOtcReport_Date = "Date"

//auto_models_end
