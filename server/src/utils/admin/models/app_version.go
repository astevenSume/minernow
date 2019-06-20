package models

//auto_models_start
type AppVersion struct {
	Id         int64  `orm:"column(id);pk" json:"id,omitempty"`
	Version    string `orm:"column(version);size(100)" json:"version,omitempty"`
	VersionNum int32  `orm:"column(version_num)" json:"version_num,omitempty"`
	ChangeLog  string `orm:"column(changelog);size(300)" json:"changelog,omitempty"`
	Download   string `orm:"column(download);size(300)" json:"download,omitempty"`
	System     int8   `orm:"column(system)" json:"system,omitempty"`
	Status     int8   `orm:"column(status)" json:"status,omitempty"`
	Ctime      int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime      int64  `orm:"column(utime)" json:"utime,omitempty"`
	Dtime      int64  `orm:"column(dtime)" json:"dtime,omitempty"`
}

func (this *AppVersion) TableName() string {
	return "app_version"
}

//table app_version name and attributes defination.
const TABLE_AppVersion = "app_version"
const COLUMN_AppVersion_Id = "id"
const COLUMN_AppVersion_Version = "version"
const COLUMN_AppVersion_VersionNum = "version_num"
const COLUMN_AppVersion_ChangeLog = "changelog"
const COLUMN_AppVersion_Download = "download"
const COLUMN_AppVersion_System = "system"
const COLUMN_AppVersion_Status = "status"
const COLUMN_AppVersion_Ctime = "ctime"
const COLUMN_AppVersion_Utime = "utime"
const COLUMN_AppVersion_Dtime = "dtime"
const ATTRIBUTE_AppVersion_Id = "Id"
const ATTRIBUTE_AppVersion_Version = "Version"
const ATTRIBUTE_AppVersion_VersionNum = "VersionNum"
const ATTRIBUTE_AppVersion_ChangeLog = "ChangeLog"
const ATTRIBUTE_AppVersion_Download = "Download"
const ATTRIBUTE_AppVersion_System = "System"
const ATTRIBUTE_AppVersion_Status = "Status"
const ATTRIBUTE_AppVersion_Ctime = "Ctime"
const ATTRIBUTE_AppVersion_Utime = "Utime"
const ATTRIBUTE_AppVersion_Dtime = "Dtime"

//auto_models_end
