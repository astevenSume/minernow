package models

//auto_models_start
type PriKey struct {
	Pkid    uint64 `orm:"column(pkid);pk" json:"pkid,omitempty"`
	Pri     string `orm:"column(pri);size(256)" json:"pri,omitempty"`
	Address string `orm:"column(address);size(100)" json:"address,omitempty"`
}

func (this *PriKey) TableName() string {
	return "usdt_prikey"
}

//table usdt_prikey name and attributes defination.
const TABLE_PriKey = "usdt_prikey"
const COLUMN_PriKey_Pkid = "pkid"
const COLUMN_PriKey_Pri = "pri"
const COLUMN_PriKey_Address = "address"
const ATTRIBUTE_PriKey_Pkid = "Pkid"
const ATTRIBUTE_PriKey_Pri = "Pri"
const ATTRIBUTE_PriKey_Address = "Address"

//auto_models_end
