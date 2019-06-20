package models

//auto_models_start
type MarketPrices struct {
	Id          uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Market      int8   `orm:"column(market)" json:"market,omitempty"`
	Currency    string `orm:"column(currency);size(100)" json:"currency,omitempty"`
	TradeMethod int8   `orm:"column(trade_method)" json:"trade_method,omitempty"`
	PowPrice    uint64 `orm:"column(pow_price)" json:"pow_price,omitempty"`
	Pow         int32  `orm:"column(pow)" json:"pow,omitempty"`
	Ctime       int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *MarketPrices) TableName() string {
	return "market_prices"
}

//table market_prices name and attributes defination.
const TABLE_MarketPrices = "market_prices"
const COLUMN_MarketPrices_Id = "id"
const COLUMN_MarketPrices_Market = "market"
const COLUMN_MarketPrices_Currency = "currency"
const COLUMN_MarketPrices_TradeMethod = "trade_method"
const COLUMN_MarketPrices_PowPrice = "pow_price"
const COLUMN_MarketPrices_Pow = "pow"
const COLUMN_MarketPrices_Ctime = "ctime"
const ATTRIBUTE_MarketPrices_Id = "Id"
const ATTRIBUTE_MarketPrices_Market = "Market"
const ATTRIBUTE_MarketPrices_Currency = "Currency"
const ATTRIBUTE_MarketPrices_TradeMethod = "TradeMethod"
const ATTRIBUTE_MarketPrices_PowPrice = "PowPrice"
const ATTRIBUTE_MarketPrices_Pow = "Pow"
const ATTRIBUTE_MarketPrices_Ctime = "Ctime"

//auto_models_end
