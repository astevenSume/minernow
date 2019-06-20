package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/admin/models"
)

const (
	APP_VERSION_SYSTEM_UNKNOWN = iota
	APP_VERSION_SYSTEM_ANDROID
	APP_VERSION_SYSTEM_IOS
	APP_VERSION_SYSTEM_MAX
)
const (
	APP_VERSION_STATUS_UNKNOWN = iota
	APP_VERSION_STATUS_DELETED
	APP_VERSION_STATUS_PUBLISHED
	APP_VERSION_STATUS_PENDING
	APP_VERSION_STATUS_MAX
)

type AppVersionDao struct {
	common.BaseDao
}

func NewAppVersionDao(db string) *AppVersionDao {
	return &AppVersionDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AppVersionDaoEntity *AppVersionDao

//查询最新app版本信息
func (d *AppVersionDao) QueryLastAppVersion(system int8) (appVersion *models.AppVersion, err error) {
	appVersion = &models.AppVersion{}
	if err = d.CheckSystem(system); err != nil {
		return
	}

	qs := d.Orm.QueryTable(models.TABLE_AppVersion).
		Filter(models.COLUMN_AppVersion_System, system).
		Exclude(models.COLUMN_AppVersion_Status, APP_VERSION_STATUS_DELETED)
	if qs == nil {
		return
	}

	err = qs.OrderBy("-" + models.COLUMN_AppVersion_VersionNum).One(appVersion)
	if err != nil {
		if err == orm.ErrNoRows {
			common.LogFuncError("ErrNoRows:%v", err)
			err = nil
			return
		}
		common.LogFuncError("DB_ERR:%v", err)
		return
	}

	return
}

func (d *AppVersionDao) CheckSystem(system int8) error {
	switch system {
	case APP_VERSION_SYSTEM_ANDROID, APP_VERSION_SYSTEM_IOS:
		return nil
	}
	return ErrParam
}

func (d *AppVersionDao) CheckStatus(status int8) error {
	switch status {
	case APP_VERSION_STATUS_PUBLISHED, APP_VERSION_STATUS_PENDING:
		return nil
	}
	return ErrParam
}

func (d *AppVersionDao) CheckCUParams(isAdd bool, appVersion *models.AppVersion) error {
	if appVersion == nil {
		common.LogFuncDebug("check err:%v", ErrParam)
		return ErrParam
	}

	if !isAdd && appVersion.Id <= 0 {
		common.LogFuncDebug("check err:%v", ErrParam)
		return ErrParam
	}

	if err := d.CheckSystem(appVersion.System); err != nil {
		common.LogFuncDebug("error code:%v", err)
		return err
	}

	if appVersion.Version == "" || appVersion.Download == "" || appVersion.ChangeLog == "" || appVersion.VersionNum <= 0 {
		common.LogFuncDebug("check err:%v", ErrParam)
		return ErrParam
	}

	return nil
}

//分页查询app版本信息
func (d *AppVersionDao) QueryPageAppVersions(system int8, status int8, page int, perPage int) ([]models.AppVersion, PageInfo, error) {

	appVersions := []models.AppVersion{}

	meta := PageInfo{
		Total: 0,
		Page:  page,
		Limit: perPage,
	}

	qs := d.Orm.QueryTable(models.TABLE_AppVersion)
	if qs == nil {
		common.LogFuncError("mysql_err:Permission fail")
		return appVersions, meta, ErrSql
	}

	if status > APP_VERSION_STATUS_UNKNOWN && status < APP_VERSION_STATUS_MAX {
		qs = qs.Filter(models.COLUMN_AppVersion_Status, status)
	}

	if system > APP_VERSION_SYSTEM_UNKNOWN && system < APP_VERSION_SYSTEM_MAX {
		qs = qs.Filter(models.COLUMN_AppVersion_System, system)
	}

	count, err := qs.Count()
	if err != nil {
		common.LogFuncError("mysql_err:Count fail: %v", err)
		return appVersions, meta, err
	}

	meta.Total = int(count)

	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > count {
		return appVersions, meta, nil
	}
	if _, err := qs.Limit(perPage, start).All(&appVersions); err != nil {
		return appVersions, meta, err
	}

	return appVersions, meta, nil
}

//创建app版本信息
func (d *AppVersionDao) CreateAppVersion(appVersion *models.AppVersion) error {
	if err := d.CheckCUParams(true, appVersion); err != nil {
		return err
	}

	t := common.NowInt64MS()
	appVersion.Utime = t
	appVersion.Ctime = t
	appVersion.Status = APP_VERSION_STATUS_PUBLISHED

	id, err := d.Orm.Insert(appVersion)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	appVersion.Id = id
	return nil
}

//修改app版本信息
func (d *AppVersionDao) UpdateAppVersion(appVersion *models.AppVersion) error {
	if err := d.CheckCUParams(false, appVersion); err != nil {
		return err
	}

	if _, err := d.QueryAppVersionById(appVersion.Id); err != nil {
		return err
	}

	t := common.NowInt64MS()
	appVersion.Utime = t

	_, err := d.Orm.Update(appVersion,
		models.COLUMN_AppVersion_Version,
		models.COLUMN_AppVersion_System,
		models.ATTRIBUTE_AppVersion_Download,
		models.ATTRIBUTE_AppVersion_ChangeLog,
		models.COLUMN_AppVersion_VersionNum,
		models.COLUMN_AppVersion_Utime,
	)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return nil
}

//删除app版本信息
func (d *AppVersionDao) DelAppVersion(id int64) error {
	appVersion := &models.AppVersion{
		Id:     id,
		Status: APP_VERSION_STATUS_DELETED,
		Dtime:  common.NowInt64MS(),
	}

	if err := d.UpdateAppVersionStatus(id, appVersion, models.COLUMN_AppVersion_Status, models.COLUMN_AppVersion_Dtime); err != nil {
		return err
	}

	return nil
}

//上架app版本
func (d *AppVersionDao) PublishAppVersion(id int64) error {
	appVersion := &models.AppVersion{
		Id:     id,
		Status: APP_VERSION_STATUS_PUBLISHED,
		Utime:  common.NowInt64MS(),
	}

	if err := d.UpdateAppVersionStatus(id, appVersion, models.COLUMN_AppVersion_Status, models.COLUMN_AppVersion_Utime); err != nil {
		return err
	}

	return nil
}

//下架app版本
func (d *AppVersionDao) PendAppVersion(id int64) error {
	appVersion := &models.AppVersion{
		Id:     id,
		Status: APP_VERSION_STATUS_PENDING,
		Utime:  common.NowInt64MS(),
	}

	if err := d.UpdateAppVersionStatus(id, appVersion, models.COLUMN_AppVersion_Status, models.COLUMN_AppVersion_Utime); err != nil {
		return err
	}

	return nil
}

//更新appVersjon
func (d *AppVersionDao) UpdateAppVersionStatus(id int64, appVersion *models.AppVersion, cols ...string) error {
	if appVersion == nil {
		return ErrParam
	}

	if _, err := d.QueryAppVersionById(appVersion.Id); err != nil {
		return err
	}

	_, err := d.Orm.Update(appVersion, cols...)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

// 查询app版本 by id
func (d *AppVersionDao) QueryAppVersionById(id int64) (models.AppVersion, error) {
	appVersion := models.AppVersion{
		Id: id,
	}
	err := d.Orm.Read(&appVersion, models.COLUMN_AppVersion_Id)
	return appVersion, err
}
