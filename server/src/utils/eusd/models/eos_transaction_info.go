package models

//auto_models_start
type EosTransactionInfo struct {
	Id            uint64 `orm:"column(id);pk" json:"id,omitempty"`
	TransactionId string `orm:"column(transaction_id);size(100)" json:"transaction_id,omitempty"`
	BlockNum      uint32 `orm:"column(block_num)" json:"block_num,omitempty"`
	Ctime         int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Processed     string `orm:"column(processed)" json:"processed,omitempty"`
}

func (this *EosTransactionInfo) TableName() string {
	return "eos_transaction_info"
}

//table eos_transaction_info name and attributes defination.
const TABLE_EosTransactionInfo = "eos_transaction_info"
const COLUMN_EosTransactionInfo_Id = "id"
const COLUMN_EosTransactionInfo_TransactionId = "transaction_id"
const COLUMN_EosTransactionInfo_BlockNum = "block_num"
const COLUMN_EosTransactionInfo_Ctime = "ctime"
const COLUMN_EosTransactionInfo_Processed = "processed"
const ATTRIBUTE_EosTransactionInfo_Id = "Id"
const ATTRIBUTE_EosTransactionInfo_TransactionId = "TransactionId"
const ATTRIBUTE_EosTransactionInfo_BlockNum = "BlockNum"
const ATTRIBUTE_EosTransactionInfo_Ctime = "Ctime"
const ATTRIBUTE_EosTransactionInfo_Processed = "Processed"

//auto_models_end
