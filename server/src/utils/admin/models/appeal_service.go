package models

//auto_models_start
type AppealService struct {
	Id      uint32 `orm:"column(id);pk" json:"id,omitempty"`
	AdminId uint32 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	Wechat  string `orm:"column(wechat);size(32)" json:"wechat,omitempty"`
	QrCode  string `orm:"column(qr_code);size(300)" json:"qr_code,omitempty"`
	Status  int8   `orm:"column(status)" json:"status,omitempty"`
}

func (this *AppealService) TableName() string {
	return "appeal_service"
}

//table appeal_service name and attributes defination.
const TABLE_AppealService = "appeal_service"
const COLUMN_AppealService_Id = "id"
const COLUMN_AppealService_AdminId = "admin_id"
const COLUMN_AppealService_Wechat = "wechat"
const COLUMN_AppealService_QrCode = "qr_code"
const COLUMN_AppealService_Status = "status"
const ATTRIBUTE_AppealService_Id = "Id"
const ATTRIBUTE_AppealService_AdminId = "AdminId"
const ATTRIBUTE_AppealService_Wechat = "Wechat"
const ATTRIBUTE_AppealService_QrCode = "QrCode"
const ATTRIBUTE_AppealService_Status = "Status"

//auto_models_end
