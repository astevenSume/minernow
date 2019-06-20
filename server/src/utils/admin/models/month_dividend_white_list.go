package models

//auto_models_start
type MonthDividendWhiteList struct {
	Id            uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Name          string `orm:"column(name);size(64)" json:"name,omitempty"`
	DividendRatio int32  `orm:"column(dividend_ratio)" json:"dividend_ratio,omitempty"`
	Ctime         int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime         int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *MonthDividendWhiteList) TableName() string {
	return "month_dividend_white_list"
}

//table month_dividend_white_list name and attributes defination.
const TABLE_MonthDividendWhiteList = "month_dividend_white_list"
const COLUMN_MonthDividendWhiteList_Id = "id"
const COLUMN_MonthDividendWhiteList_Name = "name"
const COLUMN_MonthDividendWhiteList_DividendRatio = "dividend_ratio"
const COLUMN_MonthDividendWhiteList_Ctime = "ctime"
const COLUMN_MonthDividendWhiteList_Utime = "utime"
const ATTRIBUTE_MonthDividendWhiteList_Id = "Id"
const ATTRIBUTE_MonthDividendWhiteList_Name = "Name"
const ATTRIBUTE_MonthDividendWhiteList_DividendRatio = "DividendRatio"
const ATTRIBUTE_MonthDividendWhiteList_Ctime = "Ctime"
const ATTRIBUTE_MonthDividendWhiteList_Utime = "Utime"

//auto_models_end
