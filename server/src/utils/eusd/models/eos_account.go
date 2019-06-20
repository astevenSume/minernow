package models

//auto_models_start
type EosAccount struct {
	Id      uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid     uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Account string `orm:"column(account);size(100)" json:"account,omitempty"`
	Balance string `orm:"column(balance);size(100)" json:"balance,omitempty"`
	Status  int8   `orm:"column(status)" json:"status,omitempty"`
	Ctime   int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime   int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *EosAccount) TableName() string {
	return "eos_account"
}

//table eos_account name and attributes defination.
const TABLE_EosAccount = "eos_account"
const COLUMN_EosAccount_Id = "id"
const COLUMN_EosAccount_Uid = "uid"
const COLUMN_EosAccount_Account = "account"
const COLUMN_EosAccount_Balance = "balance"
const COLUMN_EosAccount_Status = "status"
const COLUMN_EosAccount_Ctime = "ctime"
const COLUMN_EosAccount_Utime = "utime"
const ATTRIBUTE_EosAccount_Id = "Id"
const ATTRIBUTE_EosAccount_Uid = "Uid"
const ATTRIBUTE_EosAccount_Account = "Account"
const ATTRIBUTE_EosAccount_Balance = "Balance"
const ATTRIBUTE_EosAccount_Status = "Status"
const ATTRIBUTE_EosAccount_Ctime = "Ctime"
const ATTRIBUTE_EosAccount_Utime = "Utime"

//auto_models_end
