package models

//auto_models_start
type EusdRetire struct {
	Id       uint64 `orm:"column(id);pk" json:"id,omitempty"`
	From     string `orm:"column(from);size(100)" json:"from,omitempty"`
	FromUid  uint64 `orm:"column(from_uid)" json:"from_uid,omitempty"`
	Quantity int64  `orm:"column(quantity)" json:"quantity,omitempty"`
	Status   int8   `orm:"column(status)" json:"status,omitempty"`
	Ctime    int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *EusdRetire) TableName() string {
	return "eusd_retire"
}

//table eusd_retire name and attributes defination.
const TABLE_EusdRetire = "eusd_retire"
const COLUMN_EusdRetire_Id = "id"
const COLUMN_EusdRetire_From = "from"
const COLUMN_EusdRetire_FromUid = "from_uid"
const COLUMN_EusdRetire_Quantity = "quantity"
const COLUMN_EusdRetire_Status = "status"
const COLUMN_EusdRetire_Ctime = "ctime"
const ATTRIBUTE_EusdRetire_Id = "Id"
const ATTRIBUTE_EusdRetire_From = "From"
const ATTRIBUTE_EusdRetire_FromUid = "FromUid"
const ATTRIBUTE_EusdRetire_Quantity = "Quantity"
const ATTRIBUTE_EusdRetire_Status = "Status"
const ATTRIBUTE_EusdRetire_Ctime = "Ctime"

//auto_models_end
