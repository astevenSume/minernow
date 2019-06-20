package models

//auto_models_start
type EosAccountKeys struct {
	Id               uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Account          string `orm:"column(account);size(100)" json:"account,omitempty"`
	PublicKeyOwner   string `orm:"column(public_key_owner);size(100)" json:"public_key_owner,omitempty"`
	PrivateKeyOwner  string `orm:"column(private_key_owner);size(100)" json:"private_key_owner,omitempty"`
	PublicKeyActive  string `orm:"column(public_key_active);size(100)" json:"public_key_active,omitempty"`
	PrivateKeyActive string `orm:"column(private_key_active);size(100)" json:"private_key_active,omitempty"`
	Ctime            int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *EosAccountKeys) TableName() string {
	return "eos_account_keys"
}

//table eos_account_keys name and attributes defination.
const TABLE_EosAccountKeys = "eos_account_keys"
const COLUMN_EosAccountKeys_Id = "id"
const COLUMN_EosAccountKeys_Account = "account"
const COLUMN_EosAccountKeys_PublicKeyOwner = "public_key_owner"
const COLUMN_EosAccountKeys_PrivateKeyOwner = "private_key_owner"
const COLUMN_EosAccountKeys_PublicKeyActive = "public_key_active"
const COLUMN_EosAccountKeys_PrivateKeyActive = "private_key_active"
const COLUMN_EosAccountKeys_Ctime = "ctime"
const ATTRIBUTE_EosAccountKeys_Id = "Id"
const ATTRIBUTE_EosAccountKeys_Account = "Account"
const ATTRIBUTE_EosAccountKeys_PublicKeyOwner = "PublicKeyOwner"
const ATTRIBUTE_EosAccountKeys_PrivateKeyOwner = "PrivateKeyOwner"
const ATTRIBUTE_EosAccountKeys_PublicKeyActive = "PublicKeyActive"
const ATTRIBUTE_EosAccountKeys_PrivateKeyActive = "PrivateKeyActive"
const ATTRIBUTE_EosAccountKeys_Ctime = "Ctime"

//auto_models_end
