package models

//auto_models_start
type UsdtSweepLog struct {
	Id                uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid               uint64 `orm:"column(uid)" json:"uid,omitempty"`
	TType             uint32 `orm:"column(ttype)" json:"ttype,omitempty"`
	Status            uint32 `orm:"column(status)" json:"status,omitempty"`
	From              string `orm:"column(from);size(256)" json:"from,omitempty"`
	To                string `orm:"column(to);size(256)" json:"to,omitempty"`
	Txid              string `orm:"column(txid);size(256)" json:"txid,omitempty"`
	AmountInteger     int64  `orm:"column(amount_integer)" json:"amount_integer,omitempty"`
	FeeInteger        int64  `orm:"column(fee_integer)" json:"fee_integer,omitempty"`
	FeeOnchainInteger int64  `orm:"column(fee_onchain_integer)" json:"fee_onchain_integer,omitempty"`
	Ctime             int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime             int64  `orm:"column(utime)" json:"utime,omitempty"`
	Step              string `orm:"column(step);size(64)" json:"step,omitempty"`
	Desc              string `orm:"column(desc);size(256)" json:"desc,omitempty"`
}

func (this *UsdtSweepLog) TableName() string {
	return "usdt_sweep_log"
}

//table usdt_sweep_log name and attributes defination.
const TABLE_UsdtSweepLog = "usdt_sweep_log"
const COLUMN_UsdtSweepLog_Id = "id"
const COLUMN_UsdtSweepLog_Uid = "uid"
const COLUMN_UsdtSweepLog_TType = "ttype"
const COLUMN_UsdtSweepLog_Status = "status"
const COLUMN_UsdtSweepLog_From = "from"
const COLUMN_UsdtSweepLog_To = "to"
const COLUMN_UsdtSweepLog_Txid = "txid"
const COLUMN_UsdtSweepLog_AmountInteger = "amount_integer"
const COLUMN_UsdtSweepLog_FeeInteger = "fee_integer"
const COLUMN_UsdtSweepLog_FeeOnchainInteger = "fee_onchain_integer"
const COLUMN_UsdtSweepLog_Ctime = "ctime"
const COLUMN_UsdtSweepLog_Utime = "utime"
const COLUMN_UsdtSweepLog_Step = "step"
const COLUMN_UsdtSweepLog_Desc = "desc"
const ATTRIBUTE_UsdtSweepLog_Id = "Id"
const ATTRIBUTE_UsdtSweepLog_Uid = "Uid"
const ATTRIBUTE_UsdtSweepLog_TType = "TType"
const ATTRIBUTE_UsdtSweepLog_Status = "Status"
const ATTRIBUTE_UsdtSweepLog_From = "From"
const ATTRIBUTE_UsdtSweepLog_To = "To"
const ATTRIBUTE_UsdtSweepLog_Txid = "Txid"
const ATTRIBUTE_UsdtSweepLog_AmountInteger = "AmountInteger"
const ATTRIBUTE_UsdtSweepLog_FeeInteger = "FeeInteger"
const ATTRIBUTE_UsdtSweepLog_FeeOnchainInteger = "FeeOnchainInteger"
const ATTRIBUTE_UsdtSweepLog_Ctime = "Ctime"
const ATTRIBUTE_UsdtSweepLog_Utime = "Utime"
const ATTRIBUTE_UsdtSweepLog_Step = "Step"
const ATTRIBUTE_UsdtSweepLog_Desc = "Desc"

//auto_models_end
