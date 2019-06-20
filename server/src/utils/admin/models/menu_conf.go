package models

//auto_models_start
type MenuConf struct {
	Id         uint64 `orm:"column(id);pk" json:"id,omitempty"`
	PId        uint64 `orm:"column(pid)" json:"pid,omitempty"`
	Level      int32  `orm:"column(level)" json:"level,omitempty"`
	Name       string `orm:"column(name);size(100)" json:"name,omitempty"`
	Path       string `orm:"column(path);size(100)" json:"path,omitempty"`
	Icon       string `orm:"column(icon);size(100)" json:"icon,omitempty"`
	HideInMenu bool   `orm:"column(hide_in_menu)" json:"hide_in_menu,omitempty"`
	Component  string `orm:"column(component);size(100)" json:"component,omitempty"`
	OrderId    uint32 `orm:"column(order_id)" json:"order_id,omitempty"`
	CTime      int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	UTime      int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *MenuConf) TableName() string {
	return "menu_conf"
}

//table menu_conf name and attributes defination.
const TABLE_MenuConf = "menu_conf"
const COLUMN_MenuConf_Id = "id"
const COLUMN_MenuConf_PId = "pid"
const COLUMN_MenuConf_Level = "level"
const COLUMN_MenuConf_Name = "name"
const COLUMN_MenuConf_Path = "path"
const COLUMN_MenuConf_Icon = "icon"
const COLUMN_MenuConf_HideInMenu = "hide_in_menu"
const COLUMN_MenuConf_Component = "component"
const COLUMN_MenuConf_OrderId = "order_id"
const COLUMN_MenuConf_CTime = "ctime"
const COLUMN_MenuConf_UTime = "utime"
const ATTRIBUTE_MenuConf_Id = "Id"
const ATTRIBUTE_MenuConf_PId = "PId"
const ATTRIBUTE_MenuConf_Level = "Level"
const ATTRIBUTE_MenuConf_Name = "Name"
const ATTRIBUTE_MenuConf_Path = "Path"
const ATTRIBUTE_MenuConf_Icon = "Icon"
const ATTRIBUTE_MenuConf_HideInMenu = "HideInMenu"
const ATTRIBUTE_MenuConf_Component = "Component"
const ATTRIBUTE_MenuConf_OrderId = "OrderId"
const ATTRIBUTE_MenuConf_CTime = "CTime"
const ATTRIBUTE_MenuConf_UTime = "UTime"

//auto_models_end
