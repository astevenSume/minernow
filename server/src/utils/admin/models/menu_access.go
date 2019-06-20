package models

//auto_models_start
type MenuAccess struct {
	Id     uint64 `orm:"column(id);pk" json:"id,omitempty"`
	RoleId uint64 `orm:"column(role_id)" json:"role_id,omitempty"`
	MenuId uint64 `orm:"column(menu_id)" json:"menu_id,omitempty"`
}

func (this *MenuAccess) TableName() string {
	return "menu_access"
}

//table menu_access name and attributes defination.
const TABLE_MenuAccess = "menu_access"
const COLUMN_MenuAccess_Id = "id"
const COLUMN_MenuAccess_RoleId = "role_id"
const COLUMN_MenuAccess_MenuId = "menu_id"
const ATTRIBUTE_MenuAccess_Id = "Id"
const ATTRIBUTE_MenuAccess_RoleId = "RoleId"
const ATTRIBUTE_MenuAccess_MenuId = "MenuId"

//auto_models_end
