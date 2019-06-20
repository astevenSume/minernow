package models

//auto_models_start
type UserConfig struct {
	Uid          uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	WealthNotice bool   `orm:"column(wealth_notice)" json:"wealth_notice,omitempty"`
	OrderNotice  bool   `orm:"column(order_notice)" json:"order_notice,omitempty"`
}

func (this *UserConfig) TableName() string {
	return "user_config"
}

//table user_config name and attributes defination.
const TABLE_UserConfig = "user_config"
const COLUMN_UserConfig_Uid = "uid"
const COLUMN_UserConfig_WealthNotice = "wealth_notice"
const COLUMN_UserConfig_OrderNotice = "order_notice"
const ATTRIBUTE_UserConfig_Uid = "Uid"
const ATTRIBUTE_UserConfig_WealthNotice = "WealthNotice"
const ATTRIBUTE_UserConfig_OrderNotice = "OrderNotice"

//auto_models_end
