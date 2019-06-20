package models

//auto_models_start
type AppChannel struct {
	Id           uint32 `orm:"column(id);pk" json:"id,omitempty"`
	IsThirdHall  int8   `orm:"column(is_third_hall)" json:"is_third_hall,omitempty"`
	Name         string `orm:"column(name);size(256)" json:"name,omitempty"`
	Desc         string `orm:"column(desc);size(256)" json:"desc,omitempty"`
	ExchangeRate int32  `orm:"column(exchangeRate)" json:"exchangeRate,omitempty"`
	Precision    int32  `orm:"column(precision)" json:"precision,omitempty"`
	ProfitRate   int32  `orm:"column(profit_rate)" json:"profit_rate,omitempty"`
	IconUrl      string `orm:"column(icon_url);size(128)" json:"icon_url,omitempty"`
	Ctime        int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime        int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *AppChannel) TableName() string {
	return "app_channel"
}

//table app_channel name and attributes defination.
const TABLE_AppChannel = "app_channel"
const COLUMN_AppChannel_Id = "id"
const COLUMN_AppChannel_IsThirdHall = "is_third_hall"
const COLUMN_AppChannel_Name = "name"
const COLUMN_AppChannel_Desc = "desc"
const COLUMN_AppChannel_ExchangeRate = "exchangeRate"
const COLUMN_AppChannel_Precision = "precision"
const COLUMN_AppChannel_ProfitRate = "profit_rate"
const COLUMN_AppChannel_IconUrl = "icon_url"
const COLUMN_AppChannel_Ctime = "ctime"
const COLUMN_AppChannel_Utime = "utime"
const ATTRIBUTE_AppChannel_Id = "Id"
const ATTRIBUTE_AppChannel_IsThirdHall = "IsThirdHall"
const ATTRIBUTE_AppChannel_Name = "Name"
const ATTRIBUTE_AppChannel_Desc = "Desc"
const ATTRIBUTE_AppChannel_ExchangeRate = "ExchangeRate"
const ATTRIBUTE_AppChannel_Precision = "Precision"
const ATTRIBUTE_AppChannel_ProfitRate = "ProfitRate"
const ATTRIBUTE_AppChannel_IconUrl = "IconUrl"
const ATTRIBUTE_AppChannel_Ctime = "Ctime"
const ATTRIBUTE_AppChannel_Utime = "Utime"

//auto_models_end
