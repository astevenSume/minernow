package models

//auto_models_start
type User struct {
	Uid           uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	NationalCode  string `orm:"column(national_code);size(100)" json:"national_code,omitempty"`
	Mobile        string `orm:"column(mobile);size(100)" json:"mobile,omitempty"`
	Status        int8   `orm:"column(status)" json:"status,omitempty"`
	Nick          string `orm:"column(nick);size(100)" json:"nick,omitempty"`
	Password      string `orm:"column(pass);size(100)" json:"pass,omitempty"`
	Salt          string `orm:"column(salt);size(16)" json:"salt,omitempty"`
	Ctime         int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime         int64  `orm:"column(utime)" json:"utime,omitempty"`
	Ip            string `orm:"column(ip);size(100)" json:"ip,omitempty"`
	LastLoginTime int64  `orm:"column(last_login_time)" json:"last_login_time,omitempty"`
	LastLoginIp   string `orm:"column(last_login_ip);size(100)" json:"last_login_ip,omitempty"`
	IsExchanger   int8   `orm:"column(is_exchanger)" json:"is_exchanger,omitempty"`
	SignSalt      string `orm:"column(sign_salt);size(256)" json:"sign_salt,omitempty"`
}

func (this *User) TableName() string {
	return "user"
}

//table user name and attributes defination.
const TABLE_User = "user"
const COLUMN_User_Uid = "uid"
const COLUMN_User_NationalCode = "national_code"
const COLUMN_User_Mobile = "mobile"
const COLUMN_User_Status = "status"
const COLUMN_User_Nick = "nick"
const COLUMN_User_Password = "pass"
const COLUMN_User_Salt = "salt"
const COLUMN_User_Ctime = "ctime"
const COLUMN_User_Utime = "utime"
const COLUMN_User_Ip = "ip"
const COLUMN_User_LastLoginTime = "last_login_time"
const COLUMN_User_LastLoginIp = "last_login_ip"
const COLUMN_User_IsExchanger = "is_exchanger"
const COLUMN_User_SignSalt = "sign_salt"
const ATTRIBUTE_User_Uid = "Uid"
const ATTRIBUTE_User_NationalCode = "NationalCode"
const ATTRIBUTE_User_Mobile = "Mobile"
const ATTRIBUTE_User_Status = "Status"
const ATTRIBUTE_User_Nick = "Nick"
const ATTRIBUTE_User_Password = "Password"
const ATTRIBUTE_User_Salt = "Salt"
const ATTRIBUTE_User_Ctime = "Ctime"
const ATTRIBUTE_User_Utime = "Utime"
const ATTRIBUTE_User_Ip = "Ip"
const ATTRIBUTE_User_LastLoginTime = "LastLoginTime"
const ATTRIBUTE_User_LastLoginIp = "LastLoginIp"
const ATTRIBUTE_User_IsExchanger = "IsExchanger"
const ATTRIBUTE_User_SignSalt = "SignSalt"

//auto_models_end
