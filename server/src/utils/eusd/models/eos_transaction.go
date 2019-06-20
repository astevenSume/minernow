package models

//auto_models_start
type EosTransaction struct {
	Id            uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Type          uint8  `orm:"column(type)" json:"type,omitempty"`
	TransactionId string `orm:"column(transaction_id);size(100)" json:"transaction_id,omitempty"`
	BlockNum      uint32 `orm:"column(block_num)" json:"block_num,omitempty"`
	Status        int8   `orm:"column(status)" json:"status,omitempty"`
	Payer         string `orm:"column(payer);size(100)" json:"payer,omitempty"`
	Receiver      string `orm:"column(receiver);size(100)" json:"receiver,omitempty"`
	Quantity      string `orm:"column(quantity);size(100)" json:"quantity,omitempty"`
	Memo          string `orm:"column(memo);size(100)" json:"memo,omitempty"`
	Ctime         int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime         int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *EosTransaction) TableName() string {
	return "eos_transaction"
}

//table eos_transaction name and attributes defination.
const TABLE_EosTransaction = "eos_transaction"
const COLUMN_EosTransaction_Id = "id"
const COLUMN_EosTransaction_Type = "type"
const COLUMN_EosTransaction_TransactionId = "transaction_id"
const COLUMN_EosTransaction_BlockNum = "block_num"
const COLUMN_EosTransaction_Status = "status"
const COLUMN_EosTransaction_Payer = "payer"
const COLUMN_EosTransaction_Receiver = "receiver"
const COLUMN_EosTransaction_Quantity = "quantity"
const COLUMN_EosTransaction_Memo = "memo"
const COLUMN_EosTransaction_Ctime = "ctime"
const COLUMN_EosTransaction_Utime = "utime"
const ATTRIBUTE_EosTransaction_Id = "Id"
const ATTRIBUTE_EosTransaction_Type = "Type"
const ATTRIBUTE_EosTransaction_TransactionId = "TransactionId"
const ATTRIBUTE_EosTransaction_BlockNum = "BlockNum"
const ATTRIBUTE_EosTransaction_Status = "Status"
const ATTRIBUTE_EosTransaction_Payer = "Payer"
const ATTRIBUTE_EosTransaction_Receiver = "Receiver"
const ATTRIBUTE_EosTransaction_Quantity = "Quantity"
const ATTRIBUTE_EosTransaction_Memo = "Memo"
const ATTRIBUTE_EosTransaction_Ctime = "Ctime"
const ATTRIBUTE_EosTransaction_Utime = "Utime"

//auto_models_end
