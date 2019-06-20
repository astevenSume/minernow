package models

//auto_models_start
type EosWealth struct {
	Uid          uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Status       int8   `orm:"column(status)" json:"status,omitempty"`
	Account      string `orm:"column(account);size(100)" json:"account,omitempty"`
	Balance      int64  `orm:"column(balance)" json:"balance,omitempty"`
	Available    int64  `orm:"column(available)" json:"available,omitempty"`
	Game         int64  `orm:"column(game)" json:"game,omitempty"`
	Trade        int64  `orm:"column(trade)" json:"trade,omitempty"`
	Transfer     int64  `orm:"column(transfer)" json:"transfer,omitempty"`
	TransferGame int64  `orm:"column(transfer_game)" json:"transfer_game,omitempty"`
	IsExchanger  int8   `orm:"column(is_exchanger)" json:"is_exchanger,omitempty"`
	Ctime        int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime        int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *EosWealth) TableName() string {
	return "eos_wealth"
}

//table eos_wealth name and attributes defination.
const TABLE_EosWealth = "eos_wealth"
const COLUMN_EosWealth_Uid = "uid"
const COLUMN_EosWealth_Status = "status"
const COLUMN_EosWealth_Account = "account"
const COLUMN_EosWealth_Balance = "balance"
const COLUMN_EosWealth_Available = "available"
const COLUMN_EosWealth_Game = "game"
const COLUMN_EosWealth_Trade = "trade"
const COLUMN_EosWealth_Transfer = "transfer"
const COLUMN_EosWealth_TransferGame = "transfer_game"
const COLUMN_EosWealth_IsExchanger = "is_exchanger"
const COLUMN_EosWealth_Ctime = "ctime"
const COLUMN_EosWealth_Utime = "utime"
const ATTRIBUTE_EosWealth_Uid = "Uid"
const ATTRIBUTE_EosWealth_Status = "Status"
const ATTRIBUTE_EosWealth_Account = "Account"
const ATTRIBUTE_EosWealth_Balance = "Balance"
const ATTRIBUTE_EosWealth_Available = "Available"
const ATTRIBUTE_EosWealth_Game = "Game"
const ATTRIBUTE_EosWealth_Trade = "Trade"
const ATTRIBUTE_EosWealth_Transfer = "Transfer"
const ATTRIBUTE_EosWealth_TransferGame = "TransferGame"
const ATTRIBUTE_EosWealth_IsExchanger = "IsExchanger"
const ATTRIBUTE_EosWealth_Ctime = "Ctime"
const ATTRIBUTE_EosWealth_Utime = "Utime"

//auto_models_end
