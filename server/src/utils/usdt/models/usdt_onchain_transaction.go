package models

//auto_models_start
type UsdtOnchainTransaction struct {
	TxId             string `orm:"column(tx_id);pk;size(100)" json:"tx_id,omitempty"`
	Uaid             uint64 `orm:"column(uaid)" json:"uaid,omitempty"`
	Type             uint8  `orm:"column(type)" json:"type,omitempty"`
	PropertyId       uint32 `orm:"column(property_id)" json:"property_id,omitempty"`
	PropertyName     string `orm:"column(property_name);size(100)" json:"property_name,omitempty"`
	TxType           string `orm:"column(tx_type);size(100)" json:"tx_type,omitempty"`
	TxTypeInt        int32  `orm:"column(tx_type_int)" json:"tx_type_int,omitempty"`
	AmountInteger    int64  `orm:"column(amount_integer)" json:"amount_integer,omitempty"`
	Block            uint32 `orm:"column(block)" json:"block,omitempty"`
	BlockHash        string `orm:"column(block_hash);size(100)" json:"block_hash,omitempty"`
	BlockTime        int64  `orm:"column(block_time)" json:"block_time,omitempty"`
	Confirmations    uint32 `orm:"column(confirmations)" json:"confirmations,omitempty"`
	Divisible        bool   `orm:"column(divisible)" json:"divisible,omitempty"`
	FeeAmountInteger int64  `orm:"column(fee_amount_integer)" json:"fee_amount_integer,omitempty"`
	IsMine           bool   `orm:"column(is_mine)" json:"is_mine,omitempty"`
	PositionInBlock  uint32 `orm:"column(position_in_block)" json:"position_in_block,omitempty"`
	ReferenceAddress string `orm:"column(referenceaddress);size(100)" json:"referenceaddress,omitempty"`
	SendingAddress   string `orm:"column(sending_address);size(100)" json:"sending_address,omitempty"`
	Version          int32  `orm:"column(version)" json:"version,omitempty"`
	Mtime            int64  `orm:"column(mtime)" json:"mtime,omitempty"`
}

func (this *UsdtOnchainTransaction) TableName() string {
	return "usdt_onchain_transaction"
}

//table usdt_onchain_transaction name and attributes defination.
const TABLE_UsdtOnchainTransaction = "usdt_onchain_transaction"
const COLUMN_UsdtOnchainTransaction_TxId = "tx_id"
const COLUMN_UsdtOnchainTransaction_Uaid = "uaid"
const COLUMN_UsdtOnchainTransaction_Type = "type"
const COLUMN_UsdtOnchainTransaction_PropertyId = "property_id"
const COLUMN_UsdtOnchainTransaction_PropertyName = "property_name"
const COLUMN_UsdtOnchainTransaction_TxType = "tx_type"
const COLUMN_UsdtOnchainTransaction_TxTypeInt = "tx_type_int"
const COLUMN_UsdtOnchainTransaction_AmountInteger = "amount_integer"
const COLUMN_UsdtOnchainTransaction_Block = "block"
const COLUMN_UsdtOnchainTransaction_BlockHash = "block_hash"
const COLUMN_UsdtOnchainTransaction_BlockTime = "block_time"
const COLUMN_UsdtOnchainTransaction_Confirmations = "confirmations"
const COLUMN_UsdtOnchainTransaction_Divisible = "divisible"
const COLUMN_UsdtOnchainTransaction_FeeAmountInteger = "fee_amount_integer"
const COLUMN_UsdtOnchainTransaction_IsMine = "is_mine"
const COLUMN_UsdtOnchainTransaction_PositionInBlock = "position_in_block"
const COLUMN_UsdtOnchainTransaction_ReferenceAddress = "referenceaddress"
const COLUMN_UsdtOnchainTransaction_SendingAddress = "sending_address"
const COLUMN_UsdtOnchainTransaction_Version = "version"
const COLUMN_UsdtOnchainTransaction_Mtime = "mtime"
const ATTRIBUTE_UsdtOnchainTransaction_TxId = "TxId"
const ATTRIBUTE_UsdtOnchainTransaction_Uaid = "Uaid"
const ATTRIBUTE_UsdtOnchainTransaction_Type = "Type"
const ATTRIBUTE_UsdtOnchainTransaction_PropertyId = "PropertyId"
const ATTRIBUTE_UsdtOnchainTransaction_PropertyName = "PropertyName"
const ATTRIBUTE_UsdtOnchainTransaction_TxType = "TxType"
const ATTRIBUTE_UsdtOnchainTransaction_TxTypeInt = "TxTypeInt"
const ATTRIBUTE_UsdtOnchainTransaction_AmountInteger = "AmountInteger"
const ATTRIBUTE_UsdtOnchainTransaction_Block = "Block"
const ATTRIBUTE_UsdtOnchainTransaction_BlockHash = "BlockHash"
const ATTRIBUTE_UsdtOnchainTransaction_BlockTime = "BlockTime"
const ATTRIBUTE_UsdtOnchainTransaction_Confirmations = "Confirmations"
const ATTRIBUTE_UsdtOnchainTransaction_Divisible = "Divisible"
const ATTRIBUTE_UsdtOnchainTransaction_FeeAmountInteger = "FeeAmountInteger"
const ATTRIBUTE_UsdtOnchainTransaction_IsMine = "IsMine"
const ATTRIBUTE_UsdtOnchainTransaction_PositionInBlock = "PositionInBlock"
const ATTRIBUTE_UsdtOnchainTransaction_ReferenceAddress = "ReferenceAddress"
const ATTRIBUTE_UsdtOnchainTransaction_SendingAddress = "SendingAddress"
const ATTRIBUTE_UsdtOnchainTransaction_Version = "Version"
const ATTRIBUTE_UsdtOnchainTransaction_Mtime = "Mtime"

//auto_models_end
