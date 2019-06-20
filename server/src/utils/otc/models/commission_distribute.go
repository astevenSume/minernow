package models

//auto_models_start
type CommissionDistribute struct {
	Id              uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Start           string `orm:"column(start);size(16)" json:"start,omitempty"`
	End             string `orm:"column(end);size(16)" json:"end,omitempty"`
	DistributeStart string `orm:"column(distribute_start);size(128)" json:"distribute_start,omitempty"`
	DistributeEnd   string `orm:"column(distribute_end);size(128)" json:"distribute_end,omitempty"`
	Status          uint8  `orm:"column(status)" json:"status,omitempty"`
	Desc            string `orm:"column(desc);size(256)" json:"desc,omitempty"`
}

func (this *CommissionDistribute) TableName() string {
	return "commission_distribute"
}

//table commission_distribute name and attributes defination.
const TABLE_CommissionDistribute = "commission_distribute"
const COLUMN_CommissionDistribute_Id = "id"
const COLUMN_CommissionDistribute_Start = "start"
const COLUMN_CommissionDistribute_End = "end"
const COLUMN_CommissionDistribute_DistributeStart = "distribute_start"
const COLUMN_CommissionDistribute_DistributeEnd = "distribute_end"
const COLUMN_CommissionDistribute_Status = "status"
const COLUMN_CommissionDistribute_Desc = "desc"
const ATTRIBUTE_CommissionDistribute_Id = "Id"
const ATTRIBUTE_CommissionDistribute_Start = "Start"
const ATTRIBUTE_CommissionDistribute_End = "End"
const ATTRIBUTE_CommissionDistribute_DistributeStart = "DistributeStart"
const ATTRIBUTE_CommissionDistribute_DistributeEnd = "DistributeEnd"
const ATTRIBUTE_CommissionDistribute_Status = "Status"
const ATTRIBUTE_CommissionDistribute_Desc = "Desc"

//auto_models_end
