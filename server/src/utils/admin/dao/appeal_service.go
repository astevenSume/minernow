package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"utils/admin/models"
)

//代理状态
const (
	AppealServiceStatusNil = iota
	AppealServiceStatusRestore
	AppealServiceStatusSuspend
	AppealServiceStatusMax
)

const AppealServiceMin = 1

type AppealServiceDao struct {
	common.BaseDao
}

func NewAppealServiceDao(db string) *AppealServiceDao {
	return &AppealServiceDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AppealServiceDaoEntity *AppealServiceDao

type AppealService struct {
	Id        uint32 `orm:"column(id);pk" json:"id,omitempty"`
	AdminId   uint32 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	AdminNick string `orm:"column(name)" json:"admin_nick,omitempty"`
	Wechat    string `orm:"column(wechat);size(32)" json:"wechat,omitempty"`
	QrCode    string `orm:"column(qr_code);size(300)" json:"qr_code,omitempty"`
	Status    int8   `orm:"column(status)" json:"status,omitempty"`
}

func (d *AppealServiceDao) QueryById(id uint32) (appealService AppealService, err error) {
	var qbQuery orm.QueryBuilder
	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	qbQuery.Select("T1."+models.COLUMN_AppealService_Wechat,
		"T1."+models.COLUMN_AppealService_Id,
		"T1."+models.COLUMN_AppealService_Status,
		"T1."+models.COLUMN_AppealService_AdminId,
		"T1."+models.COLUMN_AppealService_QrCode,
		"T2."+models.COLUMN_AdminUser_Name).
		From("((").Select("*").From(models.TABLE_AppealService).Where(models.COLUMN_AppealService_Id + "=?")
	sqlQuery := qbQuery.String()
	sqlQuery = fmt.Sprintf("%s) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s)", sqlQuery, models.TABLE_AdminUser, models.COLUMN_AppealService_AdminId, models.COLUMN_AdminUser_Id)
	err = d.Orm.Raw(sqlQuery, id).QueryRow(&appealService)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//获取申述列表
func (d *AppealServiceDao) QueryByPage(status int8, page int, limit int) (total int64, appealServices []AppealService, err error) {
	var qbTotal orm.QueryBuilder
	var qbQuery orm.QueryBuilder
	qbTotal, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}

	// 构建查询对象
	qbTotal.Select("Count(*) AS total").
		From(models.TABLE_AppealService).
		Where("1 = 1")
	qbQuery.Select("T1."+models.COLUMN_AppealService_Wechat,
		"T1."+models.COLUMN_AppealService_Id,
		"T1."+models.COLUMN_AppealService_Status,
		"T1."+models.COLUMN_AppealService_AdminId,
		"T1."+models.COLUMN_AppealService_QrCode,
		"T2."+models.COLUMN_AdminUser_Name).
		From("((").Select("*").From(models.TABLE_AppealService).Where("1=1")
	var param []interface{}
	if status > AppealServiceStatusNil && status < AppealServiceStatusMax {
		qbTotal.And(models.COLUMN_AppealService_Status + "=?")
		qbQuery.And(models.COLUMN_AppealService_Status + "=?")
		param = append(param, status)
	}
	qbQuery.OrderBy("-" + models.COLUMN_AppealService_Id)

	sqlTotal := qbTotal.String()
	err = d.Orm.Raw(sqlTotal, param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error: %v", err)
		return
	}

	qbQuery.Limit(limit).Offset((page - 1) * limit)
	sqlQuery := qbQuery.String()
	sqlQuery = fmt.Sprintf("%s) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s)", sqlQuery, models.TABLE_AdminUser, models.COLUMN_AppealService_AdminId, models.COLUMN_AdminUser_Id)
	_, err = d.Orm.Raw(sqlQuery, param...).QueryRows(&appealServices)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *AppealServiceDao) Create(adminId uint32, wechat string, qrCode string) (err error) {
	if wechat == "" || qrCode == "" {
		err = ErrParam
		return
	}
	adminUser := &models.AdminUser{Id: uint64(adminId)}
	err = AdminUserDaoEntity.QueryAdminUser(adminUser, models.COLUMN_AdminUser_Id)
	if err != nil {
		err = ErrParam
		return
	}

	appealService := &models.AppealService{
		AdminId: adminId,
		Wechat:  wechat,
		QrCode:  qrCode,
		Status:  AppealServiceStatusRestore,
	}
	_, err = d.Orm.Insert(appealService)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *AppealServiceDao) Update(id uint32, wechat string, qrCode string) (err error) {
	appealService := &models.AppealService{
		Id:     id,
		Wechat: wechat,
		QrCode: qrCode,
	}
	_, err = d.Orm.Update(appealService, models.COLUMN_AppealService_Wechat, models.COLUMN_AppealService_QrCode)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

func (d *AppealServiceDao) DelById(id uint32) (err error) {
	appealService := &models.AppealService{Id: id}
	_, err = d.Orm.Delete(appealService, models.COLUMN_AppealService_Id)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

func (d *AppealServiceDao) SetStatus(id uint32, status int8) (err error) {
	appealService := &models.AppealService{
		Id:     id,
		Status: status,
	}
	_, err = d.Orm.Update(appealService, models.COLUMN_AppealService_Status)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

//是否是最后一个启用客服
func (d *AppealServiceDao) IsLast(id uint32) bool {
	var appealServices []models.AppealService
	_, err := d.Orm.QueryTable(models.TABLE_AppealService).
		Filter(models.COLUMN_AppealService_Status, AppealServiceStatusRestore).
		All(&appealServices, models.COLUMN_AppealService_Id)
	if err != nil {
		common.LogFuncDebug("mysql error:%v", err)
		return false
	}

	if len(appealServices) == AppealServiceMin && appealServices[0].Id == id {
		return true
	}

	return false
}

//获取申述客服管理员ID，微信
func (d *AppealServiceDao) GetAppealAdminInfo() (uint32, string, string) {
	var appealServices []models.AppealService
	_, err := d.Orm.QueryTable(models.TABLE_AppealService).
		Filter(models.COLUMN_AppealService_Status, AppealServiceStatusRestore).
		All(&appealServices, models.COLUMN_AppealService_AdminId, models.COLUMN_AppealService_Wechat, models.COLUMN_AppealService_QrCode)
	if err != nil {
		common.LogFuncDebug("mysql error:%v", err)
		return 0, "", ""
	}

	//从启用状态中随机一个
	length := len(appealServices)
	if length > 0 {
		random := rand.Intn(length)
		return appealServices[random].AdminId, appealServices[random].Wechat, appealServices[random].QrCode
	}

	return 0, "", ""
}

func (d *AppealServiceDao) GetContactByAdminID(adminId uint32) (string, string) {
	appealService := &models.AppealService{
		AdminId: adminId,
	}
	err := d.Orm.Read(appealService, models.COLUMN_AppealService_AdminId)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return "", ""
	}

	return appealService.Wechat, appealService.QrCode
}
