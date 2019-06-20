package models

//auto_models_start
type RolePermission struct {
	Id          uint64 `orm:"column(id);pk" json:"id,omitempty"`
	RoleId      uint64 `orm:"column(roleid)" json:"roleid,omitempty"`
	Pemissionid uint64 `orm:"column(permissionid)" json:"permissionid,omitempty"`
}

func (this *RolePermission) TableName() string {
	return "role_permission"
}

//table role_permission name and attributes defination.
const TABLE_RolePermission = "role_permission"
const COLUMN_RolePermission_Id = "id"
const COLUMN_RolePermission_RoleId = "roleid"
const COLUMN_RolePermission_Pemissionid = "permissionid"
const ATTRIBUTE_RolePermission_Id = "Id"
const ATTRIBUTE_RolePermission_RoleId = "RoleId"
const ATTRIBUTE_RolePermission_Pemissionid = "Pemissionid"

//auto_models_end
