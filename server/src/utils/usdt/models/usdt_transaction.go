package models

//auto_models_start
type UsdtTransaction struct {
	TxId          string `orm:"column(tx_id);pk;size(100)" json:"tx_id,omitempty"`
	Uaid          uint64 `orm:"column(uaid)" json:"uaid,omitempty"`
	Type          uint8  `orm:"column(type)" json:"type,omitempty"`
	BlockNum      uint32 `orm:"column(block_num)" json:"block_num,omitempty"`
	Status        int32  `orm:"column(status)" json:"status,omitempty"`
	Payer         string `orm:"column(payer);size(100)" json:"payer,omitempty"`
	Receiver      string `orm:"column(receiver);size(100)" json:"receiver,omitempty"`
	AmountInteger int64  `orm:"column(amount_integer)" json:"amount_integer,omitempty"`
	Fee           string `orm:"column(fee);size(100)" json:"fee,omitempty"`
	Memo          string `orm:"column(memo);size(100)" json:"memo,omitempty"`
	Ctime         int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime         int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *UsdtTransaction) TableName() string {
	return "usdt_transaction"
}

//table usdt_transaction name and attributes defination.
const TABLE_UsdtTransaction = "usdt_transaction"
const COLUMN_UsdtTransaction_TxId = "tx_id"
const COLUMN_UsdtTransaction_Uaid = "uaid"
const COLUMN_UsdtTransaction_Type = "type"
const COLUMN_UsdtTransaction_BlockNum = "block_num"
const COLUMN_UsdtTransaction_Status = "status"
const COLUMN_UsdtTransaction_Payer = "payer"
const COLUMN_UsdtTransaction_Receiver = "receiver"
const COLUMN_UsdtTransaction_AmountInteger = "amount_integer"
const COLUMN_UsdtTransaction_Fee = "fee"
const COLUMN_UsdtTransaction_Memo = "memo"
const COLUMN_UsdtTransaction_Ctime = "ctime"
const COLUMN_UsdtTransaction_Utime = "utime"
const ATTRIBUTE_UsdtTransaction_TxId = "TxId"
const ATTRIBUTE_UsdtTransaction_Uaid = "Uaid"
const ATTRIBUTE_UsdtTransaction_Type = "Type"
const ATTRIBUTE_UsdtTransaction_BlockNum = "BlockNum"
const ATTRIBUTE_UsdtTransaction_Status = "Status"
const ATTRIBUTE_UsdtTransaction_Payer = "Payer"
const ATTRIBUTE_UsdtTransaction_Receiver = "Receiver"
const ATTRIBUTE_UsdtTransaction_AmountInteger = "AmountInteger"
const ATTRIBUTE_UsdtTransaction_Fee = "Fee"
const ATTRIBUTE_UsdtTransaction_Memo = "Memo"
const ATTRIBUTE_UsdtTransaction_Ctime = "Ctime"
const ATTRIBUTE_UsdtTransaction_Utime = "Utime"

//auto_models_end
