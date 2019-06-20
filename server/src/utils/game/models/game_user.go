package models

//auto_models_start
 type GameUser struct{
	Uid uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	Account string `orm:"column(account);size(50)" json:"account,omitempty"`
	NickName string `orm:"column(nick_name);size(50)" json:"nick_name,omitempty"`
	Sex uint8 `orm:"column(sex)" json:"sex,omitempty"`
	Password string `orm:"column(password);size(100)" json:"password,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Mtime int64 `orm:"column(mtime)" json:"mtime,omitempty"`
	Status uint8 `orm:"column(status)" json:"status,omitempty"`
}

func (this *GameUser) TableName() string {
    return "game_user"
}

//table game_user name and attributes defination.
const TABLE_GameUser = "game_user"
const COLUMN_GameUser_Uid = "uid"
const COLUMN_GameUser_ChannelId = "channel_id"
const COLUMN_GameUser_Account = "account"
const COLUMN_GameUser_NickName = "nick_name"
const COLUMN_GameUser_Sex = "sex"
const COLUMN_GameUser_Password = "password"
const COLUMN_GameUser_Ctime = "ctime"
const COLUMN_GameUser_Mtime = "mtime"
const COLUMN_GameUser_Status = "status"
const ATTRIBUTE_GameUser_Uid = "Uid"
const ATTRIBUTE_GameUser_ChannelId = "ChannelId"
const ATTRIBUTE_GameUser_Account = "Account"
const ATTRIBUTE_GameUser_NickName = "NickName"
const ATTRIBUTE_GameUser_Sex = "Sex"
const ATTRIBUTE_GameUser_Password = "Password"
const ATTRIBUTE_GameUser_Ctime = "Ctime"
const ATTRIBUTE_GameUser_Mtime = "Mtime"
const ATTRIBUTE_GameUser_Status = "Status"

//auto_models_end
