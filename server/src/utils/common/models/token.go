package models

//auto_models_start
type Token struct {
	Uid         uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	ClientType  uint32 `orm:"column(client_type)" json:"client_type,omitempty"`
	MTime       int64  `orm:"column(mtime)" json:"mtime,omitempty"`
	AccessToken string `orm:"column(access_token);size(256)" json:"access_token,omitempty"`
	Mac         string `orm:"column(mac);size(100)" json:"mac,omitempty"`
}

func (this *Token) TableName() string {
	return "token"
}

//table token name and attributes defination.
const TABLE_Token = "token"
const COLUMN_Token_Uid = "uid"
const COLUMN_Token_ClientType = "client_type"
const COLUMN_Token_MTime = "mtime"
const COLUMN_Token_AccessToken = "access_token"
const COLUMN_Token_Mac = "mac"
const ATTRIBUTE_Token_Uid = "Uid"
const ATTRIBUTE_Token_ClientType = "ClientType"
const ATTRIBUTE_Token_MTime = "MTime"
const ATTRIBUTE_Token_AccessToken = "AccessToken"
const ATTRIBUTE_Token_Mac = "Mac"

//auto_models_end
