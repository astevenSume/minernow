package models

//auto_models_start
type IpWhiteList struct {
	Id uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Ip string `orm:"column(ip)" json:"ip,omitempty"`
}

func (this *IpWhiteList) TableName() string {
	return "ip_white_list"
}

//table ip_white_list name and attributes defination.
const TABLE_IpWhiteList = "ip_white_list"
const COLUMN_IpWhiteList_Id = "id"
const COLUMN_IpWhiteList_Ip = "ip"
const ATTRIBUTE_IpWhiteList_Id = "Id"
const ATTRIBUTE_IpWhiteList_Ip = "Ip"

//auto_models_end
