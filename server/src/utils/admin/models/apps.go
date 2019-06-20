package models

//auto_models_start
type Apps struct {
	Id          uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Position    uint32 `orm:"column(position)" json:"position,omitempty"`
	Name        string `orm:"column(name);size(64)" json:"name,omitempty"`
	Desc        string `orm:"column(desc);size(256)" json:"desc,omitempty"`
	Url         string `orm:"column(url);size(128)" json:"url,omitempty"`
	IconUrl     string `orm:"column(icon_url);size(128)" json:"icon_url,omitempty"`
	TypeId      int8   `orm:"column(type_id)" json:"type_id,omitempty"`
	ChannelId   uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	AppId       string `orm:"column(app_id);size(50)" json:"app_id,omitempty"`
	Featured    int8   `orm:"column(featured)" json:"featured,omitempty"`
	Status      int8   `orm:"column(status)" json:"status,omitempty"`
	Orientation int8   `orm:"column(orientation)" json:"orientation,omitempty"`
	Ctime       int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime       int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *Apps) TableName() string {
	return "apps"
}

//table apps name and attributes defination.
const TABLE_Apps = "apps"
const COLUMN_Apps_Id = "id"
const COLUMN_Apps_Position = "position"
const COLUMN_Apps_Name = "name"
const COLUMN_Apps_Desc = "desc"
const COLUMN_Apps_Url = "url"
const COLUMN_Apps_IconUrl = "icon_url"
const COLUMN_Apps_TypeId = "type_id"
const COLUMN_Apps_ChannelId = "channel_id"
const COLUMN_Apps_AppId = "app_id"
const COLUMN_Apps_Featured = "featured"
const COLUMN_Apps_Status = "status"
const COLUMN_Apps_Orientation = "orientation"
const COLUMN_Apps_Ctime = "ctime"
const COLUMN_Apps_Utime = "utime"
const ATTRIBUTE_Apps_Id = "Id"
const ATTRIBUTE_Apps_Position = "Position"
const ATTRIBUTE_Apps_Name = "Name"
const ATTRIBUTE_Apps_Desc = "Desc"
const ATTRIBUTE_Apps_Url = "Url"
const ATTRIBUTE_Apps_IconUrl = "IconUrl"
const ATTRIBUTE_Apps_TypeId = "TypeId"
const ATTRIBUTE_Apps_ChannelId = "ChannelId"
const ATTRIBUTE_Apps_AppId = "AppId"
const ATTRIBUTE_Apps_Featured = "Featured"
const ATTRIBUTE_Apps_Status = "Status"
const ATTRIBUTE_Apps_Orientation = "Orientation"
const ATTRIBUTE_Apps_Ctime = "Ctime"
const ATTRIBUTE_Apps_Utime = "Utime"

//auto_models_end
