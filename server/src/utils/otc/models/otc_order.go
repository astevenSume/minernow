package models

//auto_models_start
type OtcOrder struct {
	Id             uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid            uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Uip            string `orm:"column(uip);size(40)" json:"uip,omitempty"`
	EUid           uint64 `orm:"column(euid)" json:"euid,omitempty"`
	Eip            string `orm:"column(eip);size(40)" json:"eip,omitempty"`
	Side           int8   `orm:"column(side)" json:"side,omitempty"`
	Amount         int64  `orm:"column(amount)" json:"amount,omitempty"`
	Price          string `orm:"column(price);size(100)" json:"price,omitempty"`
	Funds          int64  `orm:"column(funds)" json:"funds,omitempty"`
	Fee            int64  `orm:"column(fee)" json:"fee,omitempty"`
	PayId          uint64 `orm:"column(pay_id)" json:"pay_id,omitempty"`
	PayType        int8   `orm:"column(pay_type)" json:"pay_type,omitempty"`
	PayName        string `orm:"column(pay_name);size(128)" json:"pay_name,omitempty"`
	PayAccount     string `orm:"column(pay_account);size(300)" json:"pay_account,omitempty"`
	PayBank        string `orm:"column(pay_bank);size(128)" json:"pay_bank,omitempty"`
	PayBankBranch  string `orm:"column(pay_bank_branch);size(300)" json:"pay_bank_branch,omitempty"`
	TransferId     uint64 `orm:"column(transfer_id)" json:"transfer_id,omitempty"`
	Ctime          int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	PayTime        int64  `orm:"column(pay_time)" json:"pay_time,omitempty"`
	FinishTime     int64  `orm:"column(finish_time)" json:"finish_time,omitempty"`
	Utime          int64  `orm:"column(utime)" json:"utime,omitempty"`
	Status         int8   `orm:"column(status)" json:"status,omitempty"`
	EPayId         uint64 `orm:"column(epay_id)" json:"epay_id,omitempty"`
	EPayType       int8   `orm:"column(epay_type)" json:"epay_type,omitempty"`
	EPayName       string `orm:"column(epay_name);size(128)" json:"epay_name,omitempty"`
	EPayAccount    string `orm:"column(epay_account);size(300)" json:"epay_account,omitempty"`
	EPayBank       string `orm:"column(epay_bank);size(128)" json:"epay_bank,omitempty"`
	EPayBankBranch string `orm:"column(epay_bank_branch);size(300)" json:"epay_bank_branch,omitempty"`
	AppealStatus   int8   `orm:"column(appeal_status)" json:"appeal_status,omitempty"`
	AdminId        uint32 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	QrCode         string `orm:"column(qr_code);size(300)" json:"qr_code,omitempty"`
	Date           int32  `orm:"column(date)" json:"date,omitempty"`
}

func (this *OtcOrder) TableName() string {
	return "otc_order"
}

//table otc_order name and attributes defination.
const TABLE_OtcOrder = "otc_order"
const COLUMN_OtcOrder_Id = "id"
const COLUMN_OtcOrder_Uid = "uid"
const COLUMN_OtcOrder_Uip = "uip"
const COLUMN_OtcOrder_EUid = "euid"
const COLUMN_OtcOrder_Eip = "eip"
const COLUMN_OtcOrder_Side = "side"
const COLUMN_OtcOrder_Amount = "amount"
const COLUMN_OtcOrder_Price = "price"
const COLUMN_OtcOrder_Funds = "funds"
const COLUMN_OtcOrder_Fee = "fee"
const COLUMN_OtcOrder_PayId = "pay_id"
const COLUMN_OtcOrder_PayType = "pay_type"
const COLUMN_OtcOrder_PayName = "pay_name"
const COLUMN_OtcOrder_PayAccount = "pay_account"
const COLUMN_OtcOrder_PayBank = "pay_bank"
const COLUMN_OtcOrder_PayBankBranch = "pay_bank_branch"
const COLUMN_OtcOrder_TransferId = "transfer_id"
const COLUMN_OtcOrder_Ctime = "ctime"
const COLUMN_OtcOrder_PayTime = "pay_time"
const COLUMN_OtcOrder_FinishTime = "finish_time"
const COLUMN_OtcOrder_Utime = "utime"
const COLUMN_OtcOrder_Status = "status"
const COLUMN_OtcOrder_EPayId = "epay_id"
const COLUMN_OtcOrder_EPayType = "epay_type"
const COLUMN_OtcOrder_EPayName = "epay_name"
const COLUMN_OtcOrder_EPayAccount = "epay_account"
const COLUMN_OtcOrder_EPayBank = "epay_bank"
const COLUMN_OtcOrder_EPayBankBranch = "epay_bank_branch"
const COLUMN_OtcOrder_AppealStatus = "appeal_status"
const COLUMN_OtcOrder_AdminId = "admin_id"
const COLUMN_OtcOrder_QrCode = "qr_code"
const COLUMN_OtcOrder_Date = "date"
const ATTRIBUTE_OtcOrder_Id = "Id"
const ATTRIBUTE_OtcOrder_Uid = "Uid"
const ATTRIBUTE_OtcOrder_Uip = "Uip"
const ATTRIBUTE_OtcOrder_EUid = "EUid"
const ATTRIBUTE_OtcOrder_Eip = "Eip"
const ATTRIBUTE_OtcOrder_Side = "Side"
const ATTRIBUTE_OtcOrder_Amount = "Amount"
const ATTRIBUTE_OtcOrder_Price = "Price"
const ATTRIBUTE_OtcOrder_Funds = "Funds"
const ATTRIBUTE_OtcOrder_Fee = "Fee"
const ATTRIBUTE_OtcOrder_PayId = "PayId"
const ATTRIBUTE_OtcOrder_PayType = "PayType"
const ATTRIBUTE_OtcOrder_PayName = "PayName"
const ATTRIBUTE_OtcOrder_PayAccount = "PayAccount"
const ATTRIBUTE_OtcOrder_PayBank = "PayBank"
const ATTRIBUTE_OtcOrder_PayBankBranch = "PayBankBranch"
const ATTRIBUTE_OtcOrder_TransferId = "TransferId"
const ATTRIBUTE_OtcOrder_Ctime = "Ctime"
const ATTRIBUTE_OtcOrder_PayTime = "PayTime"
const ATTRIBUTE_OtcOrder_FinishTime = "FinishTime"
const ATTRIBUTE_OtcOrder_Utime = "Utime"
const ATTRIBUTE_OtcOrder_Status = "Status"
const ATTRIBUTE_OtcOrder_EPayId = "EPayId"
const ATTRIBUTE_OtcOrder_EPayType = "EPayType"
const ATTRIBUTE_OtcOrder_EPayName = "EPayName"
const ATTRIBUTE_OtcOrder_EPayAccount = "EPayAccount"
const ATTRIBUTE_OtcOrder_EPayBank = "EPayBank"
const ATTRIBUTE_OtcOrder_EPayBankBranch = "EPayBankBranch"
const ATTRIBUTE_OtcOrder_AppealStatus = "AppealStatus"
const ATTRIBUTE_OtcOrder_AdminId = "AdminId"
const ATTRIBUTE_OtcOrder_QrCode = "QrCode"
const ATTRIBUTE_OtcOrder_Date = "Date"

//auto_models_end
