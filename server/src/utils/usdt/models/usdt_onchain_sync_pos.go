package models

//auto_models_start
type UsdtOnChainSyncPos struct {
	Address string `orm:"column(address);pk;size(100)" json:"address,omitempty"`
	Page    uint32 `orm:"column(page)" json:"page,omitempty"`
	TxId    string `orm:"column(tx_id);size(128)" json:"tx_id,omitempty"`
}

func (this *UsdtOnChainSyncPos) TableName() string {
	return "usdt_onchain_sync_pos"
}

//table usdt_onchain_sync_pos name and attributes defination.
const TABLE_UsdtOnChainSyncPos = "usdt_onchain_sync_pos"
const COLUMN_UsdtOnChainSyncPos_Address = "address"
const COLUMN_UsdtOnChainSyncPos_Page = "page"
const COLUMN_UsdtOnChainSyncPos_TxId = "tx_id"
const ATTRIBUTE_UsdtOnChainSyncPos_Address = "Address"
const ATTRIBUTE_UsdtOnChainSyncPos_Page = "Page"
const ATTRIBUTE_UsdtOnChainSyncPos_TxId = "TxId"

//auto_models_end
