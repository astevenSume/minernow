package models

//auto_models_start
 type GameTransfer struct{
	Id uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid uint64 `orm:"column(uid)" json:"uid,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	Account string `orm:"column(account);size(50)" json:"account,omitempty"`
	TransferType uint32 `orm:"column(transfer_type)" json:"transfer_type,omitempty"`
	Order string `orm:"column(order);size(50)" json:"order,omitempty"`
	GameOrder string `orm:"column(game_order);size(50)" json:"game_order,omitempty"`
	CoinInteger int64 `orm:"column(coin_integer)" json:"coin_integer,omitempty"`
	EusdInteger int64 `orm:"column(eusd_integer)" json:"eusd_integer,omitempty"`
	Status uint32 `orm:"column(status)" json:"status,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Desc string `orm:"column(desc);size(512)" json:"desc,omitempty"`
	Step string `orm:"column(step);size(256)" json:"step,omitempty"`
}

func (this *GameTransfer) TableName() string {
    return "game_transfer"
}

//table game_transfer name and attributes defination.
const TABLE_GameTransfer = "game_transfer"
const COLUMN_GameTransfer_Id = "id"
const COLUMN_GameTransfer_Uid = "uid"
const COLUMN_GameTransfer_ChannelId = "channel_id"
const COLUMN_GameTransfer_Account = "account"
const COLUMN_GameTransfer_TransferType = "transfer_type"
const COLUMN_GameTransfer_Order = "order"
const COLUMN_GameTransfer_GameOrder = "game_order"
const COLUMN_GameTransfer_CoinInteger = "coin_integer"
const COLUMN_GameTransfer_EusdInteger = "eusd_integer"
const COLUMN_GameTransfer_Status = "status"
const COLUMN_GameTransfer_Ctime = "ctime"
const COLUMN_GameTransfer_Desc = "desc"
const COLUMN_GameTransfer_Step = "step"
const ATTRIBUTE_GameTransfer_Id = "Id"
const ATTRIBUTE_GameTransfer_Uid = "Uid"
const ATTRIBUTE_GameTransfer_ChannelId = "ChannelId"
const ATTRIBUTE_GameTransfer_Account = "Account"
const ATTRIBUTE_GameTransfer_TransferType = "TransferType"
const ATTRIBUTE_GameTransfer_Order = "Order"
const ATTRIBUTE_GameTransfer_GameOrder = "GameOrder"
const ATTRIBUTE_GameTransfer_CoinInteger = "CoinInteger"
const ATTRIBUTE_GameTransfer_EusdInteger = "EusdInteger"
const ATTRIBUTE_GameTransfer_Status = "Status"
const ATTRIBUTE_GameTransfer_Ctime = "Ctime"
const ATTRIBUTE_GameTransfer_Desc = "Desc"
const ATTRIBUTE_GameTransfer_Step = "Step"

//auto_models_end
