package models

//auto_models_start
 type ReportEusdDaily struct{
	Uid uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Buy int64 `orm:"column(buy)" json:"buy,omitempty"`
	Sell int64 `orm:"column(sell)" json:"sell,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *ReportEusdDaily) TableName() string {
    return "report_eusd_daily"
}

//table report_eusd_daily name and attributes defination.
const TABLE_ReportEusdDaily = "report_eusd_daily"
const COLUMN_ReportEusdDaily_Uid = "uid"
const COLUMN_ReportEusdDaily_Buy = "buy"
const COLUMN_ReportEusdDaily_Sell = "sell"
const COLUMN_ReportEusdDaily_Ctime = "ctime"
const ATTRIBUTE_ReportEusdDaily_Uid = "Uid"
const ATTRIBUTE_ReportEusdDaily_Buy = "Buy"
const ATTRIBUTE_ReportEusdDaily_Sell = "Sell"
const ATTRIBUTE_ReportEusdDaily_Ctime = "Ctime"

//auto_models_end
