package models

//auto_models_start
type InviteCode struct {
	Id     uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Code   string `orm:"column(code);size(16)" json:"code,omitempty"`
	Status uint8  `orm:"column(status)" json:"status,omitempty"`
}

func (this *InviteCode) TableName() string {
	return "invite_code"
}

//table invite_code name and attributes defination.
const TABLE_InviteCode = "invite_code"
const COLUMN_InviteCode_Id = "id"
const COLUMN_InviteCode_Code = "code"
const COLUMN_InviteCode_Status = "status"
const ATTRIBUTE_InviteCode_Id = "Id"
const ATTRIBUTE_InviteCode_Code = "Code"
const ATTRIBUTE_InviteCode_Status = "Status"

//auto_models_end
