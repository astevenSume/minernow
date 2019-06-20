package models

//auto_models_start
type EosUseLog struct {
	Id          uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Type        uint8  `orm:"column(type)" json:"type,omitempty"`
	Tid         uint64 `orm:"column(tid)" json:"tid,omitempty"`
	Status      int8   `orm:"column(status)" json:"status,omitempty"`
	TidRecover  uint64 `orm:"column(tid_recover)" json:"tid_recover,omitempty"`
	QuantityNum uint64 `orm:"column(quantity_num)" json:"quantity_num,omitempty"`
}

func (this *EosUseLog) TableName() string {
	return "eos_use_log"
}

//table eos_use_log name and attributes defination.
const TABLE_EosUseLog = "eos_use_log"
const COLUMN_EosUseLog_Id = "id"
const COLUMN_EosUseLog_Type = "type"
const COLUMN_EosUseLog_Tid = "tid"
const COLUMN_EosUseLog_Status = "status"
const COLUMN_EosUseLog_TidRecover = "tid_recover"
const COLUMN_EosUseLog_QuantityNum = "quantity_num"
const ATTRIBUTE_EosUseLog_Id = "Id"
const ATTRIBUTE_EosUseLog_Type = "Type"
const ATTRIBUTE_EosUseLog_Tid = "Tid"
const ATTRIBUTE_EosUseLog_Status = "Status"
const ATTRIBUTE_EosUseLog_TidRecover = "TidRecover"
const ATTRIBUTE_EosUseLog_QuantityNum = "QuantityNum"

//auto_models_end
