package models

//auto_models_start
type CommissionCalc struct {
	Id        uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Start     string `orm:"column(start);size(16)" json:"start,omitempty"`
	End       string `orm:"column(end);size(16)" json:"end,omitempty"`
	CalcStart string `orm:"column(calc_start);size(128)" json:"calc_start,omitempty"`
	CalcEnd   string `orm:"column(calc_end);size(128)" json:"calc_end,omitempty"`
	Status    uint8  `orm:"column(status)" json:"status,omitempty"`
	Desc      string `orm:"column(desc);size(256)" json:"desc,omitempty"`
}

func (this *CommissionCalc) TableName() string {
	return "commission_calc"
}

//table commission_calc name and attributes defination.
const TABLE_CommissionCalc = "commission_calc"
const COLUMN_CommissionCalc_Id = "id"
const COLUMN_CommissionCalc_Start = "start"
const COLUMN_CommissionCalc_End = "end"
const COLUMN_CommissionCalc_CalcStart = "calc_start"
const COLUMN_CommissionCalc_CalcEnd = "calc_end"
const COLUMN_CommissionCalc_Status = "status"
const COLUMN_CommissionCalc_Desc = "desc"
const ATTRIBUTE_CommissionCalc_Id = "Id"
const ATTRIBUTE_CommissionCalc_Start = "Start"
const ATTRIBUTE_CommissionCalc_End = "End"
const ATTRIBUTE_CommissionCalc_CalcStart = "CalcStart"
const ATTRIBUTE_CommissionCalc_CalcEnd = "CalcEnd"
const ATTRIBUTE_CommissionCalc_Status = "Status"
const ATTRIBUTE_CommissionCalc_Desc = "Desc"

//auto_models_end
