package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/admin/models"
)

//通知状态
const (
	SysNotificatioStatusNil = iota
	SysNotificatioStatusUnPublish
	SysNotificatioStatusPublish
	SysNotificatioStatusMax
)

type SysNotificationDao struct {
	common.BaseDao
}

func NewSysNotificationDao(db string) *SysNotificationDao {
	return &SysNotificationDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var SysNotificationDaoEntity *SysNotificationDao

//获取系统通知
func (d *SysNotificationDao) QueryById(id uint32) (*models.SysNotification, error) {
	sysNotification := &models.SysNotification{Id: id}
	err := d.Orm.Read(sysNotification, models.COLUMN_SysNotification_Id)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return nil, err
	}

	return sysNotification, nil
}

//分页条件查询
func (d *SysNotificationDao) QueryByPage(status int8, page int, perPage int) (total int64, sysNotifications []models.SysNotification, err error) {
	qs := d.Orm.QueryTable(models.TABLE_SysNotification)
	if status > SysNotificatioStatusNil && status < SysNotificatioStatusMax {
		qs = qs.Filter(models.COLUMN_SysNotification_Status, status)
	}
	qs = qs.OrderBy("-" + models.COLUMN_SysNotification_Utime)

	total, err = qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > total {
		err = nil
		return
	}
	_, err = qs.Limit(perPage, start).All(&sysNotifications)
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

func (d *SysNotificationDao) Create(adminId uint32, content string) (*models.SysNotification, error) {
	sysNotification := &models.SysNotification{
		AdminId: adminId,
		Content: content,
		Status:  SysNotificatioStatusUnPublish,
		Ctime:   common.NowInt64MS(),
		Utime:   common.NowInt64MS(),
	}

	id, err := d.Orm.Insert(sysNotification)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return nil, err
	}
	sysNotification.Id = uint32(id)

	return sysNotification, nil
}

func (d *SysNotificationDao) Update(id uint32, adminId uint32, content string) (*models.SysNotification, error) {
	sysNotification := &models.SysNotification{
		Id:      id,
		AdminId: adminId,
		Content: content,
		Utime:   common.NowInt64MS(),
	}
	_, err := d.Orm.Update(sysNotification, models.COLUMN_SysNotification_AdminId, models.COLUMN_SysNotification_Content, models.COLUMN_SysNotification_Utime)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return nil, err
	}

	return sysNotification, nil
}

func (d *SysNotificationDao) DelById(id uint32) (err error) {
	sysNotification := &models.SysNotification{Id: id}
	_, err = d.Orm.Delete(sysNotification, models.COLUMN_SysNotification_Id)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

func (d *SysNotificationDao) SetStatus(id uint32, status int8) (err error) {
	sysNotification := &models.SysNotification{
		Id:     id,
		Status: status,
		Utime:  common.NowInt64MS(),
	}

	_, err = d.Orm.Update(sysNotification, models.COLUMN_SysNotification_Status, models.COLUMN_SysNotification_Utime)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}
