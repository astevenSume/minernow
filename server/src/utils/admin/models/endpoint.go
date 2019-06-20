package models

//auto_models_start
type Endpoint struct {
	Id       uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Endpoint string `orm:"column(endpoint);size(100)" json:"endpoint,omitempty"`
	Ctime    int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime    int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *Endpoint) TableName() string {
	return "endpoint"
}

//table endpoint name and attributes defination.
const TABLE_Endpoint = "endpoint"
const COLUMN_Endpoint_Id = "id"
const COLUMN_Endpoint_Endpoint = "endpoint"
const COLUMN_Endpoint_Ctime = "ctime"
const COLUMN_Endpoint_Utime = "utime"
const ATTRIBUTE_Endpoint_Id = "Id"
const ATTRIBUTE_Endpoint_Endpoint = "Endpoint"
const ATTRIBUTE_Endpoint_Ctime = "Ctime"
const ATTRIBUTE_Endpoint_Utime = "Utime"

//auto_models_end
