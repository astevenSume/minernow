package models

//auto_models_start
type PaymentMethod struct {
	Pmid                uint64 `orm:"column(pmid);pk" json:"pmid,omitempty"`
	Uid                 uint64 `orm:"column(uid)" json:"uid,omitempty"`
	MType               uint8  `orm:"column(mtype)" json:"mtype,omitempty"`
	Ord                 uint32 `orm:"column(ord)" json:"ord,omitempty"`
	Name                string `orm:"column(name);size(128)" json:"name,omitempty"`
	Account             string `orm:"column(account);size(128)" json:"account,omitempty"`
	Status              uint8  `orm:"column(status)" json:"status,omitempty"`
	Ctime               uint32 `orm:"column(ctime)" json:"ctime,omitempty"`
	Bank                string `orm:"column(bank);size(128)" json:"bank,omitempty"`
	BankBranch          string `orm:"column(bank_branch);size(128)" json:"bank_branch,omitempty"`
	QRCode              string `orm:"column(qr_code);size(256)" json:"qr_code,omitempty"`
	QRCodeContent       string `orm:"column(qr_code_content)" json:"qr_code_content,omitempty"`
	LowMoneyPerTxLimit  int64  `orm:"column(low_money_per_tx_limit)" json:"low_money_per_tx_limit,omitempty"`
	HighMoneyPerTxLimit int64  `orm:"column(high_money_per_tx_limit)" json:"high_money_per_tx_limit,omitempty"`
	TimesPerDayLimit    int64  `orm:"column(times_per_day_limit)" json:"times_per_day_limit,omitempty"`
	MoneyPerDayLimit    int64  `orm:"column(money_per_day_limit)" json:"money_per_day_limit,omitempty"`
	MoneySumLimit       int64  `orm:"column(money_sum_limit)" json:"money_sum_limit,omitempty"`
	TimesToday          int64  `orm:"column(times_today)" json:"times_today,omitempty"`
	MoneyToday          int64  `orm:"column(money_today)" json:"money_today,omitempty"`
	MoneySum            uint64 `orm:"column(money_sum)" json:"money_sum,omitempty"`
	Mtime               uint32 `orm:"column(mtime)" json:"mtime,omitempty"`
	UseTime             uint32 `orm:"column(use_time)" json:"use_time,omitempty"`
}

func (this *PaymentMethod) TableName() string {
	return "payment_method"
}

//table payment_method name and attributes defination.
const TABLE_PaymentMethod = "payment_method"
const COLUMN_PaymentMethod_Pmid = "pmid"
const COLUMN_PaymentMethod_Uid = "uid"
const COLUMN_PaymentMethod_MType = "mtype"
const COLUMN_PaymentMethod_Ord = "ord"
const COLUMN_PaymentMethod_Name = "name"
const COLUMN_PaymentMethod_Account = "account"
const COLUMN_PaymentMethod_Status = "status"
const COLUMN_PaymentMethod_Ctime = "ctime"
const COLUMN_PaymentMethod_Bank = "bank"
const COLUMN_PaymentMethod_BankBranch = "bank_branch"
const COLUMN_PaymentMethod_QRCode = "qr_code"
const COLUMN_PaymentMethod_QRCodeContent = "qr_code_content"
const COLUMN_PaymentMethod_LowMoneyPerTxLimit = "low_money_per_tx_limit"
const COLUMN_PaymentMethod_HighMoneyPerTxLimit = "high_money_per_tx_limit"
const COLUMN_PaymentMethod_TimesPerDayLimit = "times_per_day_limit"
const COLUMN_PaymentMethod_MoneyPerDayLimit = "money_per_day_limit"
const COLUMN_PaymentMethod_MoneySumLimit = "money_sum_limit"
const COLUMN_PaymentMethod_TimesToday = "times_today"
const COLUMN_PaymentMethod_MoneyToday = "money_today"
const COLUMN_PaymentMethod_MoneySum = "money_sum"
const COLUMN_PaymentMethod_Mtime = "mtime"
const COLUMN_PaymentMethod_UseTime = "use_time"
const ATTRIBUTE_PaymentMethod_Pmid = "Pmid"
const ATTRIBUTE_PaymentMethod_Uid = "Uid"
const ATTRIBUTE_PaymentMethod_MType = "MType"
const ATTRIBUTE_PaymentMethod_Ord = "Ord"
const ATTRIBUTE_PaymentMethod_Name = "Name"
const ATTRIBUTE_PaymentMethod_Account = "Account"
const ATTRIBUTE_PaymentMethod_Status = "Status"
const ATTRIBUTE_PaymentMethod_Ctime = "Ctime"
const ATTRIBUTE_PaymentMethod_Bank = "Bank"
const ATTRIBUTE_PaymentMethod_BankBranch = "BankBranch"
const ATTRIBUTE_PaymentMethod_QRCode = "QRCode"
const ATTRIBUTE_PaymentMethod_QRCodeContent = "QRCodeContent"
const ATTRIBUTE_PaymentMethod_LowMoneyPerTxLimit = "LowMoneyPerTxLimit"
const ATTRIBUTE_PaymentMethod_HighMoneyPerTxLimit = "HighMoneyPerTxLimit"
const ATTRIBUTE_PaymentMethod_TimesPerDayLimit = "TimesPerDayLimit"
const ATTRIBUTE_PaymentMethod_MoneyPerDayLimit = "MoneyPerDayLimit"
const ATTRIBUTE_PaymentMethod_MoneySumLimit = "MoneySumLimit"
const ATTRIBUTE_PaymentMethod_TimesToday = "TimesToday"
const ATTRIBUTE_PaymentMethod_MoneyToday = "MoneyToday"
const ATTRIBUTE_PaymentMethod_MoneySum = "MoneySum"
const ATTRIBUTE_PaymentMethod_Mtime = "Mtime"
const ATTRIBUTE_PaymentMethod_UseTime = "UseTime"

//auto_models_end
