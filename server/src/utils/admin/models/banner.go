package models

//auto_models_start
type Banner struct {
	Id      uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Subject string `orm:"column(subject);size(256)" json:"subject,omitempty"`
	Image   string `orm:"column(image);size(256)" json:"image,omitempty"`
	Url     string `orm:"column(url);size(256)" json:"url,omitempty"`
	Status  int8   `orm:"column(status)" json:"status,omitempty"`
	Ctime   int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime   int64  `orm:"column(utime)" json:"utime,omitempty"`
	Stime   int64  `orm:"column(stime)" json:"stime,omitempty"`
	Etime   int64  `orm:"column(etime)" json:"etime,omitempty"`
}

func (this *Banner) TableName() string {
	return "banner"
}

//table banner name and attributes defination.
const TABLE_Banner = "banner"
const COLUMN_Banner_Id = "id"
const COLUMN_Banner_Subject = "subject"
const COLUMN_Banner_Image = "image"
const COLUMN_Banner_Url = "url"
const COLUMN_Banner_Status = "status"
const COLUMN_Banner_Ctime = "ctime"
const COLUMN_Banner_Utime = "utime"
const COLUMN_Banner_Stime = "stime"
const COLUMN_Banner_Etime = "etime"
const ATTRIBUTE_Banner_Id = "Id"
const ATTRIBUTE_Banner_Subject = "Subject"
const ATTRIBUTE_Banner_Image = "Image"
const ATTRIBUTE_Banner_Url = "Url"
const ATTRIBUTE_Banner_Status = "Status"
const ATTRIBUTE_Banner_Ctime = "Ctime"
const ATTRIBUTE_Banner_Utime = "Utime"
const ATTRIBUTE_Banner_Stime = "Stime"
const ATTRIBUTE_Banner_Etime = "Etime"

//auto_models_end
