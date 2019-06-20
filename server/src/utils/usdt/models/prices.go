package models

//auto_models_start
type Prices struct {
	Id       uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Currency string `orm:"column(currency);size(100)" json:"currency,omitempty"`
	PowPrice uint64 `orm:"column(pow_price)" json:"pow_price,omitempty"`
	Pow      int32  `orm:"column(pow)" json:"pow,omitempty"`
	Ctime    int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *Prices) TableName() string {
	return "prices"
}

//table prices name and attributes defination.
const TABLE_Prices = "prices"
const COLUMN_Prices_Id = "id"
const COLUMN_Prices_Currency = "currency"
const COLUMN_Prices_PowPrice = "pow_price"
const COLUMN_Prices_Pow = "pow"
const COLUMN_Prices_Ctime = "ctime"
const ATTRIBUTE_Prices_Id = "Id"
const ATTRIBUTE_Prices_Currency = "Currency"
const ATTRIBUTE_Prices_PowPrice = "PowPrice"
const ATTRIBUTE_Prices_Pow = "Pow"
const ATTRIBUTE_Prices_Ctime = "Ctime"

//auto_models_end
