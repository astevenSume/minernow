package models

//auto_models_start
type AdminUser struct {
	Id           uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Name         string `orm:"column(name);size(100)" json:"name,omitempty"`
	Email        string `orm:"column(email);size(100)" json:"email,omitempty"`
	Status       int8   `orm:"column(status)" json:"status,omitempty"`
	Ctime        int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime        int64  `orm:"column(utime)" json:"utime,omitempty"`
	Dtime        int64  `orm:"column(dtime)" json:"dtime,omitempty"`
	LoginTime    int64  `orm:"column(login_time)" json:"login_time,omitempty"`
	Pwd          string `orm:"column(pwd);size(100)" json:"pwd,omitempty"`
	WhitelistIps string `orm:"column(whitelist_ips);size(256)" json:"whitelist_ips,omitempty"`
	IsBind       bool   `orm:"column(is_bind)" json:"is_bind,omitempty"`
	SecretId     string `orm:"column(secret_id);size(128)" json:"secret_id,omitempty"`
	QrCode       string `orm:"column(qr_code)" json:"qr_code,omitempty"`
}

func (this *AdminUser) TableName() string {
	return "admin_user"
}

//table admin_user name and attributes defination.
const TABLE_AdminUser = "admin_user"
const COLUMN_AdminUser_Id = "id"
const COLUMN_AdminUser_Name = "name"
const COLUMN_AdminUser_Email = "email"
const COLUMN_AdminUser_Status = "status"
const COLUMN_AdminUser_Ctime = "ctime"
const COLUMN_AdminUser_Utime = "utime"
const COLUMN_AdminUser_Dtime = "dtime"
const COLUMN_AdminUser_LoginTime = "login_time"
const COLUMN_AdminUser_Pwd = "pwd"
const COLUMN_AdminUser_WhitelistIps = "whitelist_ips"
const COLUMN_AdminUser_IsBind = "is_bind"
const COLUMN_AdminUser_SecretId = "secret_id"
const COLUMN_AdminUser_QrCode = "qr_code"
const ATTRIBUTE_AdminUser_Id = "Id"
const ATTRIBUTE_AdminUser_Name = "Name"
const ATTRIBUTE_AdminUser_Email = "Email"
const ATTRIBUTE_AdminUser_Status = "Status"
const ATTRIBUTE_AdminUser_Ctime = "Ctime"
const ATTRIBUTE_AdminUser_Utime = "Utime"
const ATTRIBUTE_AdminUser_Dtime = "Dtime"
const ATTRIBUTE_AdminUser_LoginTime = "LoginTime"
const ATTRIBUTE_AdminUser_Pwd = "Pwd"
const ATTRIBUTE_AdminUser_WhitelistIps = "WhitelistIps"
const ATTRIBUTE_AdminUser_IsBind = "IsBind"
const ATTRIBUTE_AdminUser_SecretId = "SecretId"
const ATTRIBUTE_AdminUser_QrCode = "QrCode"

//auto_models_end
