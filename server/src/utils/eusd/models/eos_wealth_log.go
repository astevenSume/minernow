package models

//auto_models_start
type EosWealthLog struct {
	Id       uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid      uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Uid2     uint64 `orm:"column(uid2)" json:"uid2,omitempty"`
	TType    uint8  `orm:"column(ttype)" json:"ttype,omitempty"`
	Status   int8   `orm:"column(status)" json:"status,omitempty"`
	Txid     uint64 `orm:"column(txid)" json:"txid,omitempty"`
	Quantity int64  `orm:"column(quantity)" json:"quantity,omitempty"`
	Ctime    int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *EosWealthLog) TableName() string {
	return "eos_wealth_log"
}

//table eos_wealth_log name and attributes defination.
const TABLE_EosWealthLog = "eos_wealth_log"
const COLUMN_EosWealthLog_Id = "id"
const COLUMN_EosWealthLog_Uid = "uid"
const COLUMN_EosWealthLog_Uid2 = "uid2"
const COLUMN_EosWealthLog_TType = "ttype"
const COLUMN_EosWealthLog_Status = "status"
const COLUMN_EosWealthLog_Txid = "txid"
const COLUMN_EosWealthLog_Quantity = "quantity"
const COLUMN_EosWealthLog_Ctime = "ctime"
const ATTRIBUTE_EosWealthLog_Id = "Id"
const ATTRIBUTE_EosWealthLog_Uid = "Uid"
const ATTRIBUTE_EosWealthLog_Uid2 = "Uid2"
const ATTRIBUTE_EosWealthLog_TType = "TType"
const ATTRIBUTE_EosWealthLog_Status = "Status"
const ATTRIBUTE_EosWealthLog_Txid = "Txid"
const ATTRIBUTE_EosWealthLog_Quantity = "Quantity"
const ATTRIBUTE_EosWealthLog_Ctime = "Ctime"

//auto_models_end
