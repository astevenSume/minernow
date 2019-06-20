package models

//auto_models_start
type Commissionrates struct {
	Id         int64  `orm:"column(id);pk" json:"id,omitempty"`
	Min        uint64 `orm:"column(min)" json:"min,omitempty"`
	Max        uint64 `orm:"column(max)" json:"max,omitempty"`
	Commission int32  `orm:"column(commission)" json:"commission,omitempty"`
	Precision  int32  `orm:"column(precision)" json:"precision,omitempty"`
	Ctime      int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime      int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *Commissionrates) TableName() string {
	return "commissionrates"
}

//table commissionrates name and attributes defination.
const TABLE_Commissionrates = "commissionrates"
const COLUMN_Commissionrates_Id = "id"
const COLUMN_Commissionrates_Min = "min"
const COLUMN_Commissionrates_Max = "max"
const COLUMN_Commissionrates_Commission = "commission"
const COLUMN_Commissionrates_Precision = "precision"
const COLUMN_Commissionrates_Ctime = "ctime"
const COLUMN_Commissionrates_Utime = "utime"
const ATTRIBUTE_Commissionrates_Id = "Id"
const ATTRIBUTE_Commissionrates_Min = "Min"
const ATTRIBUTE_Commissionrates_Max = "Max"
const ATTRIBUTE_Commissionrates_Commission = "Commission"
const ATTRIBUTE_Commissionrates_Precision = "Precision"
const ATTRIBUTE_Commissionrates_Ctime = "Ctime"
const ATTRIBUTE_Commissionrates_Utime = "Utime"

//auto_models_end
