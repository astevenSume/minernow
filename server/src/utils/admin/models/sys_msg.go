package models

//auto_models_start
type SystemMessage struct {
	Id     uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Key    string `orm:"column(key);size(200)" json:"key,omitempty"`
	Buyer  string `orm:"column(buyer);size(400)" json:"buyer,omitempty"`
	Seller string `orm:"column(seller);size(400)" json:"seller,omitempty"`
	Admin  string `orm:"column(admin);size(400)" json:"admin,omitempty"`
	Ctime  int64  `orm:"column(ctime);null" json:"ctime,omitempty"`
	Utime  int64  `orm:"column(utime);null" json:"utime,omitempty"`
	Dtime  int64  `orm:"column(dtime);null" json:"dtime,omitempty"`
}

func (this *SystemMessage) TableName() string {
	return "sys_msg"
}

//table sys_msg name and attributes defination.
const TABLE_SystemMessage = "sys_msg"
const COLUMN_SystemMessage_Id = "id"
const COLUMN_SystemMessage_Key = "key"
const COLUMN_SystemMessage_Buyer = "buyer"
const COLUMN_SystemMessage_Seller = "seller"
const COLUMN_SystemMessage_Admin = "admin"
const COLUMN_SystemMessage_Ctime = "ctime"
const COLUMN_SystemMessage_Utime = "utime"
const COLUMN_SystemMessage_Dtime = "dtime"
const ATTRIBUTE_SystemMessage_Id = "Id"
const ATTRIBUTE_SystemMessage_Key = "Key"
const ATTRIBUTE_SystemMessage_Buyer = "Buyer"
const ATTRIBUTE_SystemMessage_Seller = "Seller"
const ATTRIBUTE_SystemMessage_Admin = "Admin"
const ATTRIBUTE_SystemMessage_Ctime = "Ctime"
const ATTRIBUTE_SystemMessage_Utime = "Utime"
const ATTRIBUTE_SystemMessage_Dtime = "Dtime"

//auto_models_end
