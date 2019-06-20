package models

//auto_models_start
type UsdtAccount struct {
	Uaid                    uint64 `orm:"column(uaid);pk" json:"uaid,omitempty"`
	Uid                     uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Status                  uint8  `orm:"column(status)" json:"status,omitempty"`
	AvailableInteger        int64  `orm:"column(available_integer)" json:"available_integer,omitempty"`
	FrozenInteger           int64  `orm:"column(frozen_integer)" json:"frozen_integer,omitempty"`
	TransferFrozenInteger   int64  `orm:"column(transfer_frozen_integer)" json:"transfer_frozen_integer,omitempty"`
	MortgagedInteger        int64  `orm:"column(mortgaged_integer)" json:"mortgaged_integer,omitempty"`
	BtcAvailableInteger     int64  `orm:"column(btc_available_integer)" json:"btc_available_integer,omitempty"`
	BtcFrozenInteger        int64  `orm:"column(btc_frozen_integer)" json:"btc_frozen_integer,omitempty"`
	BtcMortgagedInteger     int64  `orm:"column(btc_mortgaged_integer)" json:"btc_mortgaged_integer,omitempty"`
	WaitingCashSweepInteger int64  `orm:"column(waiting_cash_sweep_integer)" json:"waiting_cash_sweep_integer,omitempty"`
	CashSweepInteger        int64  `orm:"column(cash_sweep_integer)" json:"cash_sweep_integer,omitempty"`
	OwnedByPlatformInteger  int64  `orm:"column(owned_by_platform_integer)" json:"owned_by_platform_integer,omitempty"`
	SweepStatus             uint8  `orm:"column(sweep_status)" json:"sweep_status,omitempty"`
	Pkid                    uint64 `orm:"column(pkid)" json:"pkid,omitempty"`
	Address                 string `orm:"column(address);size(100)" json:"address,omitempty"`
	Ctime                   int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Mtime                   int64  `orm:"column(mtime)" json:"mtime,omitempty"`
	Sign                    string `orm:"column(sign);size(256)" json:"sign,omitempty"`
}

func (this *UsdtAccount) TableName() string {
	return "usdt_account"
}

//table usdt_account name and attributes defination.
const TABLE_UsdtAccount = "usdt_account"
const COLUMN_UsdtAccount_Uaid = "uaid"
const COLUMN_UsdtAccount_Uid = "uid"
const COLUMN_UsdtAccount_Status = "status"
const COLUMN_UsdtAccount_AvailableInteger = "available_integer"
const COLUMN_UsdtAccount_FrozenInteger = "frozen_integer"
const COLUMN_UsdtAccount_TransferFrozenInteger = "transfer_frozen_integer"
const COLUMN_UsdtAccount_MortgagedInteger = "mortgaged_integer"
const COLUMN_UsdtAccount_BtcAvailableInteger = "btc_available_integer"
const COLUMN_UsdtAccount_BtcFrozenInteger = "btc_frozen_integer"
const COLUMN_UsdtAccount_BtcMortgagedInteger = "btc_mortgaged_integer"
const COLUMN_UsdtAccount_WaitingCashSweepInteger = "waiting_cash_sweep_integer"
const COLUMN_UsdtAccount_CashSweepInteger = "cash_sweep_integer"
const COLUMN_UsdtAccount_OwnedByPlatformInteger = "owned_by_platform_integer"
const COLUMN_UsdtAccount_SweepStatus = "sweep_status"
const COLUMN_UsdtAccount_Pkid = "pkid"
const COLUMN_UsdtAccount_Address = "address"
const COLUMN_UsdtAccount_Ctime = "ctime"
const COLUMN_UsdtAccount_Mtime = "mtime"
const COLUMN_UsdtAccount_Sign = "sign"
const ATTRIBUTE_UsdtAccount_Uaid = "Uaid"
const ATTRIBUTE_UsdtAccount_Uid = "Uid"
const ATTRIBUTE_UsdtAccount_Status = "Status"
const ATTRIBUTE_UsdtAccount_AvailableInteger = "AvailableInteger"
const ATTRIBUTE_UsdtAccount_FrozenInteger = "FrozenInteger"
const ATTRIBUTE_UsdtAccount_TransferFrozenInteger = "TransferFrozenInteger"
const ATTRIBUTE_UsdtAccount_MortgagedInteger = "MortgagedInteger"
const ATTRIBUTE_UsdtAccount_BtcAvailableInteger = "BtcAvailableInteger"
const ATTRIBUTE_UsdtAccount_BtcFrozenInteger = "BtcFrozenInteger"
const ATTRIBUTE_UsdtAccount_BtcMortgagedInteger = "BtcMortgagedInteger"
const ATTRIBUTE_UsdtAccount_WaitingCashSweepInteger = "WaitingCashSweepInteger"
const ATTRIBUTE_UsdtAccount_CashSweepInteger = "CashSweepInteger"
const ATTRIBUTE_UsdtAccount_OwnedByPlatformInteger = "OwnedByPlatformInteger"
const ATTRIBUTE_UsdtAccount_SweepStatus = "SweepStatus"
const ATTRIBUTE_UsdtAccount_Pkid = "Pkid"
const ATTRIBUTE_UsdtAccount_Address = "Address"
const ATTRIBUTE_UsdtAccount_Ctime = "Ctime"
const ATTRIBUTE_UsdtAccount_Mtime = "Mtime"
const ATTRIBUTE_UsdtAccount_Sign = "Sign"

//auto_models_end
