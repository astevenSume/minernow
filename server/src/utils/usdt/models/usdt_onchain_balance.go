package models

//auto_models_start
type UsdtOnchainBalance struct {
	Address       string `orm:"column(address);pk;size(100)" json:"address,omitempty"`
	PropertyId    uint32 `orm:"column(property_id)" json:"property_id,omitempty"`
	PendingPos    string `orm:"column(pending_pos);size(100)" json:"pending_pos,omitempty"`
	Reserved      string `orm:"column(reserved);size(100)" json:"reserved,omitempty"`
	Divisible     bool   `orm:"column(divisible)" json:"divisible,omitempty"`
	AmountInteger int64  `orm:"column(amount_integer)" json:"amount_integer,omitempty"`
	Frozen        string `orm:"column(frozen);size(100)" json:"frozen,omitempty"`
	PendingNeg    string `orm:"column(pending_neg);size(100)" json:"pending_neg,omitempty"`
	Mtime         int64  `orm:"column(mtime)" json:"mtime,omitempty"`
}

func (this *UsdtOnchainBalance) TableName() string {
	return "usdt_onchain_balance"
}

//table usdt_onchain_balance name and attributes defination.
const TABLE_UsdtOnchainBalance = "usdt_onchain_balance"
const COLUMN_UsdtOnchainBalance_Address = "address"
const COLUMN_UsdtOnchainBalance_PropertyId = "property_id"
const COLUMN_UsdtOnchainBalance_PendingPos = "pending_pos"
const COLUMN_UsdtOnchainBalance_Reserved = "reserved"
const COLUMN_UsdtOnchainBalance_Divisible = "divisible"
const COLUMN_UsdtOnchainBalance_AmountInteger = "amount_integer"
const COLUMN_UsdtOnchainBalance_Frozen = "frozen"
const COLUMN_UsdtOnchainBalance_PendingNeg = "pending_neg"
const COLUMN_UsdtOnchainBalance_Mtime = "mtime"
const ATTRIBUTE_UsdtOnchainBalance_Address = "Address"
const ATTRIBUTE_UsdtOnchainBalance_PropertyId = "PropertyId"
const ATTRIBUTE_UsdtOnchainBalance_PendingPos = "PendingPos"
const ATTRIBUTE_UsdtOnchainBalance_Reserved = "Reserved"
const ATTRIBUTE_UsdtOnchainBalance_Divisible = "Divisible"
const ATTRIBUTE_UsdtOnchainBalance_AmountInteger = "AmountInteger"
const ATTRIBUTE_UsdtOnchainBalance_Frozen = "Frozen"
const ATTRIBUTE_UsdtOnchainBalance_PendingNeg = "PendingNeg"
const ATTRIBUTE_UsdtOnchainBalance_Mtime = "Mtime"

//auto_models_end
