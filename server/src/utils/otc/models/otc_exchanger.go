package models

//auto_models_start
type OtcExchanger struct {
	Uid      uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Mobile   string `orm:"column(mobile);size(100)" json:"mobile,omitempty"`
	Wechat   string `orm:"column(wechat);size(100)" json:"wechat,omitempty"`
	Telegram string `orm:"column(telegram);size(100)" json:"telegram,omitempty"`
	From     int8   `orm:"column(from)" json:"from,omitempty"`
	Ctime    int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime    int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *OtcExchanger) TableName() string {
	return "otc_exchanger"
}

//table otc_exchanger name and attributes defination.
const TABLE_OtcExchanger = "otc_exchanger"
const COLUMN_OtcExchanger_Uid = "uid"
const COLUMN_OtcExchanger_Mobile = "mobile"
const COLUMN_OtcExchanger_Wechat = "wechat"
const COLUMN_OtcExchanger_Telegram = "telegram"
const COLUMN_OtcExchanger_From = "from"
const COLUMN_OtcExchanger_Ctime = "ctime"
const COLUMN_OtcExchanger_Utime = "utime"
const ATTRIBUTE_OtcExchanger_Uid = "Uid"
const ATTRIBUTE_OtcExchanger_Mobile = "Mobile"
const ATTRIBUTE_OtcExchanger_Wechat = "Wechat"
const ATTRIBUTE_OtcExchanger_Telegram = "Telegram"
const ATTRIBUTE_OtcExchanger_From = "From"
const ATTRIBUTE_OtcExchanger_Ctime = "Ctime"
const ATTRIBUTE_OtcExchanger_Utime = "Utime"

//auto_models_end
