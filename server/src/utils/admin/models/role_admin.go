package models

//auto_models_start
type RoleAdmin struct {
	Id        uint64 `orm:"column(id);pk" json:"id,omitempty"`
	RoleId    uint64 `orm:"column(roleid)" json:"roleid,omitempty"`
	AdminId   uint64 `orm:"column(adminid)" json:"adminid,omitempty"`
	GrantedBy string `orm:"column(granted_by);size(100)" json:"granted_by,omitempty"`
	GrantedAt int64  `orm:"column(granted_at)" json:"granted_at,omitempty"`
}

func (this *RoleAdmin) TableName() string {
	return "role_admin"
}

//table role_admin name and attributes defination.
const TABLE_RoleAdmin = "role_admin"
const COLUMN_RoleAdmin_Id = "id"
const COLUMN_RoleAdmin_RoleId = "roleid"
const COLUMN_RoleAdmin_AdminId = "adminid"
const COLUMN_RoleAdmin_GrantedBy = "granted_by"
const COLUMN_RoleAdmin_GrantedAt = "granted_at"
const ATTRIBUTE_RoleAdmin_Id = "Id"
const ATTRIBUTE_RoleAdmin_RoleId = "RoleId"
const ATTRIBUTE_RoleAdmin_AdminId = "AdminId"
const ATTRIBUTE_RoleAdmin_GrantedBy = "GrantedBy"
const ATTRIBUTE_RoleAdmin_GrantedAt = "GrantedAt"

//auto_models_end
