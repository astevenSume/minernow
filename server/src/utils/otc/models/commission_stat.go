package models

//auto_models_start
type CommissionStat struct {
	Ctime              int64 `orm:"column(ctime);pk" json:"ctime,omitempty"`
	TaxInteger         int32 `orm:"column(tax_integer)" json:"tax_integer,omitempty"`
	TaxDecimals        int32 `orm:"column(tax_decimals)" json:"tax_decimals,omitempty"`
	ChannelInteger     int32 `orm:"column(channel_integer)" json:"channel_integer,omitempty"`
	ChannelDecimals    int32 `orm:"column(channel_decimals)" json:"channel_decimals,omitempty"`
	CommissionInteger  int32 `orm:"column(commission_integer)" json:"commission_integer,omitempty"`
	CommissionDecimals int32 `orm:"column(commission_decimals)" json:"commission_decimals,omitempty"`
	ProfitInteger      int32 `orm:"column(profit_integer)" json:"profit_integer,omitempty"`
	ProfitDecimals     int32 `orm:"column(profit_decimals)" json:"profit_decimals,omitempty"`
	Mtime              int64 `orm:"column(mtime)" json:"mtime,omitempty"`
	Status             uint8 `orm:"column(status)" json:"status,omitempty"`
}

func (this *CommissionStat) TableName() string {
	return "commission_stat"
}

//table commission_stat name and attributes defination.
const TABLE_CommissionStat = "commission_stat"
const COLUMN_CommissionStat_Ctime = "ctime"
const COLUMN_CommissionStat_TaxInteger = "tax_integer"
const COLUMN_CommissionStat_TaxDecimals = "tax_decimals"
const COLUMN_CommissionStat_ChannelInteger = "channel_integer"
const COLUMN_CommissionStat_ChannelDecimals = "channel_decimals"
const COLUMN_CommissionStat_CommissionInteger = "commission_integer"
const COLUMN_CommissionStat_CommissionDecimals = "commission_decimals"
const COLUMN_CommissionStat_ProfitInteger = "profit_integer"
const COLUMN_CommissionStat_ProfitDecimals = "profit_decimals"
const COLUMN_CommissionStat_Mtime = "mtime"
const COLUMN_CommissionStat_Status = "status"
const ATTRIBUTE_CommissionStat_Ctime = "Ctime"
const ATTRIBUTE_CommissionStat_TaxInteger = "TaxInteger"
const ATTRIBUTE_CommissionStat_TaxDecimals = "TaxDecimals"
const ATTRIBUTE_CommissionStat_ChannelInteger = "ChannelInteger"
const ATTRIBUTE_CommissionStat_ChannelDecimals = "ChannelDecimals"
const ATTRIBUTE_CommissionStat_CommissionInteger = "CommissionInteger"
const ATTRIBUTE_CommissionStat_CommissionDecimals = "CommissionDecimals"
const ATTRIBUTE_CommissionStat_ProfitInteger = "ProfitInteger"
const ATTRIBUTE_CommissionStat_ProfitDecimals = "ProfitDecimals"
const ATTRIBUTE_CommissionStat_Mtime = "Mtime"
const ATTRIBUTE_CommissionStat_Status = "Status"

//auto_models_end
