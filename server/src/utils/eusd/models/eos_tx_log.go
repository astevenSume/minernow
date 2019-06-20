package models

//auto_models_start
type EosTxLog struct {
	Id        uint64 `orm:"column(id);pk" json:"id,omitempty"`
	From      string `orm:"column(from);size(100)" json:"from,omitempty"`
	FromUid   uint64 `orm:"column(from_uid)" json:"from_uid,omitempty"`
	To        string `orm:"column(to);size(100)" json:"to,omitempty"`
	ToUid     uint64 `orm:"column(to_uid)" json:"to_uid,omitempty"`
	Quantity  int64  `orm:"column(quantity)" json:"quantity,omitempty"`
	Status    int8   `orm:"column(status)" json:"status,omitempty"`
	LogIds    string `orm:"column(log_ids);size(100)" json:"log_ids,omitempty"`
	Ctime     int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Txid      uint64 `orm:"column(txid)" json:"txid,omitempty"`
	OrderId   uint64 `orm:"column(order_id);size(100)" json:"order_id,omitempty"`
	Utime     int64  `orm:"column(utime)" json:"utime,omitempty"`
	Sign      string `orm:"column(sign);size(100)" json:"sign,omitempty"`
	DelayDeal bool   `orm:"column(delay_deal)" json:"delay_deal,omitempty"`
	Retry     int32  `orm:"column(retry)" json:"retry,omitempty"`
	Memo      string `orm:"column(memo);size(100)" json:"memo,omitempty"`
}

func (this *EosTxLog) TableName() string {
	return "eos_tx_log"
}

//table eos_tx_log name and attributes defination.
const TABLE_EosTxLog = "eos_tx_log"
const COLUMN_EosTxLog_Id = "id"
const COLUMN_EosTxLog_From = "from"
const COLUMN_EosTxLog_FromUid = "from_uid"
const COLUMN_EosTxLog_To = "to"
const COLUMN_EosTxLog_ToUid = "to_uid"
const COLUMN_EosTxLog_Quantity = "quantity"
const COLUMN_EosTxLog_Status = "status"
const COLUMN_EosTxLog_LogIds = "log_ids"
const COLUMN_EosTxLog_Ctime = "ctime"
const COLUMN_EosTxLog_Txid = "txid"
const COLUMN_EosTxLog_OrderId = "order_id"
const COLUMN_EosTxLog_Utime = "utime"
const COLUMN_EosTxLog_Sign = "sign"
const COLUMN_EosTxLog_DelayDeal = "delay_deal"
const COLUMN_EosTxLog_Retry = "retry"
const COLUMN_EosTxLog_Memo = "memo"
const ATTRIBUTE_EosTxLog_Id = "Id"
const ATTRIBUTE_EosTxLog_From = "From"
const ATTRIBUTE_EosTxLog_FromUid = "FromUid"
const ATTRIBUTE_EosTxLog_To = "To"
const ATTRIBUTE_EosTxLog_ToUid = "ToUid"
const ATTRIBUTE_EosTxLog_Quantity = "Quantity"
const ATTRIBUTE_EosTxLog_Status = "Status"
const ATTRIBUTE_EosTxLog_LogIds = "LogIds"
const ATTRIBUTE_EosTxLog_Ctime = "Ctime"
const ATTRIBUTE_EosTxLog_Txid = "Txid"
const ATTRIBUTE_EosTxLog_OrderId = "OrderId"
const ATTRIBUTE_EosTxLog_Utime = "Utime"
const ATTRIBUTE_EosTxLog_Sign = "Sign"
const ATTRIBUTE_EosTxLog_DelayDeal = "DelayDeal"
const ATTRIBUTE_EosTxLog_Retry = "Retry"
const ATTRIBUTE_EosTxLog_Memo = "Memo"

//auto_models_end
