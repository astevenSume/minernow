package models

//auto_models_start
type Announcement struct {
	Id      uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Type    int8   `orm:"column(type)" json:"type,omitempty"`
	Title   string `orm:"column(title);size(100)" json:"title,omitempty"`
	Content string `orm:"column(content);size(500)" json:"content,omitempty"`
	Stime   int64  `orm:"column(stime)" json:"stime,omitempty"`
	Etime   int64  `orm:"column(etime)" json:"etime,omitempty"`
	Ctime   int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime   int64  `orm:"column(utime)" json:"utime,omitempty"`
}

func (this *Announcement) TableName() string {
	return "announcement"
}

//table announcement name and attributes defination.
const TABLE_Announcement = "announcement"
const COLUMN_Announcement_Id = "id"
const COLUMN_Announcement_Type = "type"
const COLUMN_Announcement_Title = "title"
const COLUMN_Announcement_Content = "content"
const COLUMN_Announcement_Stime = "stime"
const COLUMN_Announcement_Etime = "etime"
const COLUMN_Announcement_Ctime = "ctime"
const COLUMN_Announcement_Utime = "utime"
const ATTRIBUTE_Announcement_Id = "Id"
const ATTRIBUTE_Announcement_Type = "Type"
const ATTRIBUTE_Announcement_Title = "Title"
const ATTRIBUTE_Announcement_Content = "Content"
const ATTRIBUTE_Announcement_Stime = "Stime"
const ATTRIBUTE_Announcement_Etime = "Etime"
const ATTRIBUTE_Announcement_Ctime = "Ctime"
const ATTRIBUTE_Announcement_Utime = "Utime"

//auto_models_end
