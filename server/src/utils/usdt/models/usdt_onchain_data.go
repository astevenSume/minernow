package models

//auto_models_start
type UsdtOnChainData struct {
	Address   string `orm:"column(address);pk;size(100)" json:"address,omitempty"`
	DataInt64 int64  `orm:"column(data_int64)" json:"data_int64,omitempty"`
	AttrType  uint32 `orm:"column(attr_type)" json:"attr_type,omitempty"`
	DataStr   string `orm:"column(data_str);size(256)" json:"data_str,omitempty"`
}

func (this *UsdtOnChainData) TableName() string {
	return "usdt_onchain_data"
}

//table usdt_onchain_data name and attributes defination.
const TABLE_UsdtOnChainData = "usdt_onchain_data"
const COLUMN_UsdtOnChainData_Address = "address"
const COLUMN_UsdtOnChainData_DataInt64 = "data_int64"
const COLUMN_UsdtOnChainData_AttrType = "attr_type"
const COLUMN_UsdtOnChainData_DataStr = "data_str"
const ATTRIBUTE_UsdtOnChainData_Address = "Address"
const ATTRIBUTE_UsdtOnChainData_DataInt64 = "DataInt64"
const ATTRIBUTE_UsdtOnChainData_AttrType = "AttrType"
const ATTRIBUTE_UsdtOnChainData_DataStr = "DataStr"

//auto_models_end
