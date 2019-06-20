package models

//auto_models_start
type UsdtOnchainLog struct {
	Oclid         uint64 `orm:"column(oclid);pk" json:"oclid,omitempty"`
	From          string `orm:"column(from);size(100)" json:"from,omitempty"`
	To            string `orm:"column(to);size(100)" json:"to,omitempty"`
	Tx            string `orm:"column(tx);size(100)" json:"tx,omitempty"`
	Status        string `orm:"column(status);size(100)" json:"status,omitempty"`
	Pushed        string `orm:"column(pushed);size(100)" json:"pushed,omitempty"`
	SignedTx      string `orm:"column(signedTx)" json:"signedTx,omitempty"`
	AmountInteger int64  `orm:"column(amount_integer)" json:"amount_integer,omitempty"`
	Ctime         int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *UsdtOnchainLog) TableName() string {
	return "usdt_onchain_log"
}

//table usdt_onchain_log name and attributes defination.
const TABLE_UsdtOnchainLog = "usdt_onchain_log"
const COLUMN_UsdtOnchainLog_Oclid = "oclid"
const COLUMN_UsdtOnchainLog_From = "from"
const COLUMN_UsdtOnchainLog_To = "to"
const COLUMN_UsdtOnchainLog_Tx = "tx"
const COLUMN_UsdtOnchainLog_Status = "status"
const COLUMN_UsdtOnchainLog_Pushed = "pushed"
const COLUMN_UsdtOnchainLog_SignedTx = "signedTx"
const COLUMN_UsdtOnchainLog_AmountInteger = "amount_integer"
const COLUMN_UsdtOnchainLog_Ctime = "ctime"
const ATTRIBUTE_UsdtOnchainLog_Oclid = "Oclid"
const ATTRIBUTE_UsdtOnchainLog_From = "From"
const ATTRIBUTE_UsdtOnchainLog_To = "To"
const ATTRIBUTE_UsdtOnchainLog_Tx = "Tx"
const ATTRIBUTE_UsdtOnchainLog_Status = "Status"
const ATTRIBUTE_UsdtOnchainLog_Pushed = "Pushed"
const ATTRIBUTE_UsdtOnchainLog_SignedTx = "SignedTx"
const ATTRIBUTE_UsdtOnchainLog_AmountInteger = "AmountInteger"
const ATTRIBUTE_UsdtOnchainLog_Ctime = "Ctime"

//auto_models_end
