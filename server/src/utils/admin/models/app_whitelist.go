package models

//auto_models_start
type AppWhitelist struct {
	Id        uint32 `orm:"column(id);pk" json:"id,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	AppId     string `orm:"column(app_id);size(16)" json:"app_id,omitempty"`
	Ctime     int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *AppWhitelist) TableName() string {
	return "app_whitelist"
}

//table app_whitelist name and attributes defination.
const TABLE_AppWhitelist = "app_whitelist"
const COLUMN_AppWhitelist_Id = "id"
const COLUMN_AppWhitelist_ChannelId = "channel_id"
const COLUMN_AppWhitelist_AppId = "app_id"
const COLUMN_AppWhitelist_Ctime = "ctime"
const ATTRIBUTE_AppWhitelist_Id = "Id"
const ATTRIBUTE_AppWhitelist_ChannelId = "ChannelId"
const ATTRIBUTE_AppWhitelist_AppId = "AppId"
const ATTRIBUTE_AppWhitelist_Ctime = "Ctime"

//auto_models_end
