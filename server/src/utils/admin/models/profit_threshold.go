package models

//auto_models_start
type ProfitThreshold struct {
	Id        uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Threshold int64  `orm:"column(threshold)" json:"threshold,omitempty"`
	AdminId   uint64 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	Ctime     int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime     int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *ProfitThreshold) TableName() string {
	return "profit_threshold"
}

//table profit_threshold name and attributes defination.
const TABLE_ProfitThreshold = "profit_threshold"
const COLUMN_ProfitThreshold_Id = "id"
const COLUMN_ProfitThreshold_Threshold = "threshold"
const COLUMN_ProfitThreshold_AdminId = "admin_id"
const COLUMN_ProfitThreshold_Ctime = "ctime"
const COLUMN_ProfitThreshold_Utime = "utime"
const ATTRIBUTE_ProfitThreshold_Id = "Id"
const ATTRIBUTE_ProfitThreshold_Threshold = "Threshold"
const ATTRIBUTE_ProfitThreshold_AdminId = "AdminId"
const ATTRIBUTE_ProfitThreshold_Ctime = "Ctime"
const ATTRIBUTE_ProfitThreshold_Utime = "Utime"

//auto_models_end
