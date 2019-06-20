package models

//auto_models_start
type UsdtWealthLog struct {
	Id                uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid               uint64 `orm:"column(uid)" json:"uid,omitempty"`
	TType             uint32 `orm:"column(ttype)" json:"ttype,omitempty"`
	Status            uint32 `orm:"column(status)" json:"status,omitempty"`
	From              string `orm:"column(from);size(256)" json:"from,omitempty"`
	To                string `orm:"column(to);size(256)" json:"to,omitempty"`
	Txid              string `orm:"column(txid);size(256)" json:"txid,omitempty"`
	AmountInteger     int64  `orm:"column(amount_integer)" json:"amount_integer,omitempty"`
	FeeInteger        int64  `orm:"column(fee_integer)" json:"fee_integer,omitempty"`
	FeeUsdtInteger    int64  `orm:"column(fee_usdt_integer)" json:"fee_usdt_integer,omitempty"`
	FeeOnchainInteger int64  `orm:"column(fee_onchain_integer)" json:"fee_onchain_integer,omitempty"`
	Ctime             int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime             int64  `orm:"column(utime)" json:"utime,omitempty"`
	Step              string `orm:"column(step);size(64)" json:"step,omitempty"`
	Desc              string `orm:"column(desc);size(256)" json:"desc,omitempty"`
	Sign              string `orm:"column(sign);size(256)" json:"sign,omitempty"`
	Memo              string `orm:"column(memo);size(256)" json:"memo,omitempty"`
}

func (this *UsdtWealthLog) TableName() string {
	return "usdt_wealth_log"
}

//table usdt_wealth_log name and attributes defination.
const TABLE_UsdtWealthLog = "usdt_wealth_log"
const COLUMN_UsdtWealthLog_Id = "id"
const COLUMN_UsdtWealthLog_Uid = "uid"
const COLUMN_UsdtWealthLog_TType = "ttype"
const COLUMN_UsdtWealthLog_Status = "status"
const COLUMN_UsdtWealthLog_From = "from"
const COLUMN_UsdtWealthLog_To = "to"
const COLUMN_UsdtWealthLog_Txid = "txid"
const COLUMN_UsdtWealthLog_AmountInteger = "amount_integer"
const COLUMN_UsdtWealthLog_FeeInteger = "fee_integer"
const COLUMN_UsdtWealthLog_FeeUsdtInteger = "fee_usdt_integer"
const COLUMN_UsdtWealthLog_FeeOnchainInteger = "fee_onchain_integer"
const COLUMN_UsdtWealthLog_Ctime = "ctime"
const COLUMN_UsdtWealthLog_Utime = "utime"
const COLUMN_UsdtWealthLog_Step = "step"
const COLUMN_UsdtWealthLog_Desc = "desc"
const COLUMN_UsdtWealthLog_Sign = "sign"
const COLUMN_UsdtWealthLog_Memo = "memo"
const ATTRIBUTE_UsdtWealthLog_Id = "Id"
const ATTRIBUTE_UsdtWealthLog_Uid = "Uid"
const ATTRIBUTE_UsdtWealthLog_TType = "TType"
const ATTRIBUTE_UsdtWealthLog_Status = "Status"
const ATTRIBUTE_UsdtWealthLog_From = "From"
const ATTRIBUTE_UsdtWealthLog_To = "To"
const ATTRIBUTE_UsdtWealthLog_Txid = "Txid"
const ATTRIBUTE_UsdtWealthLog_AmountInteger = "AmountInteger"
const ATTRIBUTE_UsdtWealthLog_FeeInteger = "FeeInteger"
const ATTRIBUTE_UsdtWealthLog_FeeUsdtInteger = "FeeUsdtInteger"
const ATTRIBUTE_UsdtWealthLog_FeeOnchainInteger = "FeeOnchainInteger"
const ATTRIBUTE_UsdtWealthLog_Ctime = "Ctime"
const ATTRIBUTE_UsdtWealthLog_Utime = "Utime"
const ATTRIBUTE_UsdtWealthLog_Step = "Step"
const ATTRIBUTE_UsdtWealthLog_Desc = "Desc"
const ATTRIBUTE_UsdtWealthLog_Sign = "Sign"
const ATTRIBUTE_UsdtWealthLog_Memo = "Memo"

//auto_models_end
