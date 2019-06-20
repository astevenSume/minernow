package models

//auto_models_start
type OtcExchangerVerify struct {
	Id       int32  `orm:"column(id);pk" json:"id,omitempty"`
	Uid      uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Mobile   string `orm:"column(mobile);size(100)" json:"mobile,omitempty"`
	Wechat   string `orm:"column(wechat);size(100)" json:"wechat,omitempty"`
	Telegram string `orm:"column(telegram);size(100)" json:"telegram,omitempty"`
	Status   int8   `orm:"column(status)" json:"status,omitempty"`
	From     int8   `orm:"column(from)" json:"from,omitempty"`
	Ctime    int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime    int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *OtcExchangerVerify) TableName() string {
	return "otc_exchanger_verify"
}

//table otc_exchanger_verify name and attributes defination.
const TABLE_OtcExchangerVerify = "otc_exchanger_verify"
const COLUMN_OtcExchangerVerify_Id = "id"
const COLUMN_OtcExchangerVerify_Uid = "uid"
const COLUMN_OtcExchangerVerify_Mobile = "mobile"
const COLUMN_OtcExchangerVerify_Wechat = "wechat"
const COLUMN_OtcExchangerVerify_Telegram = "telegram"
const COLUMN_OtcExchangerVerify_Status = "status"
const COLUMN_OtcExchangerVerify_From = "from"
const COLUMN_OtcExchangerVerify_Ctime = "ctime"
const COLUMN_OtcExchangerVerify_Utime = "utime"
const ATTRIBUTE_OtcExchangerVerify_Id = "Id"
const ATTRIBUTE_OtcExchangerVerify_Uid = "Uid"
const ATTRIBUTE_OtcExchangerVerify_Mobile = "Mobile"
const ATTRIBUTE_OtcExchangerVerify_Wechat = "Wechat"
const ATTRIBUTE_OtcExchangerVerify_Telegram = "Telegram"
const ATTRIBUTE_OtcExchangerVerify_Status = "Status"
const ATTRIBUTE_OtcExchangerVerify_From = "From"
const ATTRIBUTE_OtcExchangerVerify_Ctime = "Ctime"
const ATTRIBUTE_OtcExchangerVerify_Utime = "Utime"

//auto_models_end
