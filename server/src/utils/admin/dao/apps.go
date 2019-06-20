package dao

import (
	"common"
	"fmt"
	"utils/admin/models"

	"github.com/astaxie/beego/orm"
)

//推荐状态
const (
	AppFeaturedNil = iota
	AppFeaturedNo
	AppFeaturedYes
	AppFeaturedMax
)

//上线状态
const (
	AppStatusNil       = iota
	AppStatusUnPublish //未上线
	AppStatusPublish   //已上线
	AppStatusPlan      //计划上线
	AppStatusMax
)

//横屏竖屏
const (
	AppOrientationNil = iota
	AppOrientationHorizontal
	AppOrientationVertical
	AppOrientationMax
)

type AppDao struct {
	common.BaseDao
}

func NewAppDao(db string) *AppDao {
	return &AppDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AppDaoEntity *AppDao

//获取APP信息
func (d *AppDao) QueryAppById(id uint32) (apps *models.Apps, err error) {
	apps = &models.Apps{
		Id: id,
	}

	err = d.Orm.Read(apps, models.COLUMN_Apps_Id)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

//分页条件查询
func (d *AppDao) QueryPageCondition(channelID int8, typeId int8, feature int8, status int8, page int, perPage int) (total int64, apps []models.Apps, err error) {
	qs := d.Orm.QueryTable(models.TABLE_Apps)
	if qs == nil {
		common.LogFuncError("QueryTable fail")
		return
	}

	if typeId > 0 {
		qs = qs.Filter(models.COLUMN_Apps_TypeId, typeId)
	}
	if status > AppStatusNil && feature < AppStatusMax {
		qs = qs.Filter(models.COLUMN_Apps_Status, status)
	}
	if feature > AppFeaturedNil && feature < AppFeaturedMax {
		qs = qs.Filter(models.COLUMN_Apps_Featured, feature)
	}
	if channelID > 0 {
		qs = qs.Filter(models.COLUMN_Apps_ChannelId, channelID)
	}

	total, err = qs.Count()
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > total {
		err = nil
		return
	}
	qs = qs.OrderBy(models.COLUMN_Apps_ChannelId, models.COLUMN_Apps_Position)
	_, err = qs.Limit(perPage, start).All(&apps)
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

func (d *AppDao) CreateApp(typeId int8, orientation int8, featured int8, channelId uint32, appId string, position uint32, name, desc, url, iconUrl string) (apps models.Apps, err error) {
	if name == "" || orientation <= AppOrientationNil || orientation >= AppOrientationMax {
		err = ErrParam
		return
	}

	err = d.ChangePos(true, 0, channelId, position)
	if err != nil {
		return
	}

	apps = models.Apps{
		TypeId:      typeId,
		Name:        name,
		Desc:        desc,
		Url:         url,
		IconUrl:     iconUrl,
		Featured:    AppFeaturedNo,
		Status:      AppStatusUnPublish,
		Orientation: orientation,
		ChannelId:   channelId,
		AppId:       appId,
		Position:    position,
		Ctime:       common.NowInt64MS(),
		Utime:       common.NowInt64MS(),
	}
	if featured > AppFeaturedNil && featured < AppFeaturedMax {
		apps.Featured = featured
	}

	var id int64
	id, err = d.Orm.Insert(&apps)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}
	apps.Id = uint32(id)

	return
}

func (d *AppDao) UpdateApp(id uint32, orientation int8, status int8, typeId int8, featured int8, channelId uint32, appId string,
	position uint32, name, desc, url, iconUrl string) (models.Apps, error) {
	apps := models.Apps{Id: id}
	err := d.ChangePos(true, id, channelId, position)
	if err != nil {
		return apps, err
	}

	apps.TypeId = typeId
	apps.ChannelId = channelId
	apps.AppId = appId

	cols := []string{models.COLUMN_Apps_TypeId, models.COLUMN_Apps_ChannelId, models.COLUMN_Apps_AppId}
	if name != "" {
		apps.Name = name
		cols = append(cols, models.COLUMN_Apps_Name)
	}
	if desc != "" {
		apps.Desc = desc
		cols = append(cols, models.COLUMN_Apps_Desc)
	}
	if url != "" {
		apps.Url = url
		cols = append(cols, models.COLUMN_Apps_Url)
	}
	if iconUrl != "" {
		apps.IconUrl = iconUrl
		cols = append(cols, models.COLUMN_Apps_IconUrl)
	}
	if featured > AppFeaturedNil && featured < AppFeaturedMax {
		apps.Featured = featured
		cols = append(cols, models.COLUMN_Apps_Featured)
	}
	if status > AppStatusNil && status < AppStatusMax {
		apps.Status = status
		cols = append(cols, models.COLUMN_Apps_Status)
	}
	if orientation > AppOrientationNil && orientation < AppOrientationMax {
		apps.Orientation = orientation
		cols = append(cols, models.COLUMN_Apps_Orientation)
	}
	if position > 0 {
		apps.Position = position
		cols = append(cols, models.COLUMN_Apps_Position)
	}
	apps.Utime = common.NowInt64MS()
	cols = append(cols, models.COLUMN_Apps_Utime)

	_, err = d.Orm.Update(&apps, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return apps, err
	}

	return apps, nil
}

func (d *AppDao) DelAppById(id uint32) error {
	apps := &models.Apps{
		Id: id,
	}

	_, err := d.Orm.Delete(apps, models.COLUMN_Apps_Id)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

//更新是否推荐
func (d *AppDao) UpdateFeatured(id uint32, feature int8) error {
	apps := &models.Apps{
		Id:       id,
		Featured: feature,
		Utime:    common.NowInt64MS(),
	}

	_, err := d.Orm.Update(apps, models.COLUMN_Apps_Featured, models.COLUMN_Apps_Utime)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

//更新上线状态
func (d *AppDao) UpdateStatus(id uint32, status int8) error {
	apps := &models.Apps{
		Id:     id,
		Status: status,
		Utime:  common.NowInt64MS(),
	}

	_, err := d.Orm.Update(apps, models.COLUMN_Apps_Status, models.COLUMN_Apps_Utime)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

// get all apps
func (d *AppDao) AllPublished() (apps []models.Apps, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s=? OR %s=?", models.TABLE_Apps, models.COLUMN_Apps_Status, models.COLUMN_Apps_Status), AppStatusPublish, AppStatusPlan).QueryRows(&apps)
	if err != nil {
		common.LogFuncError("%v", err)
	}

	return
}

type AppFeatured struct {
	Name      string `orm:"column(name);size(64)" json:"app_name,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	AppId     string `orm:"column(app_id);size(50)" json:"app_id,omitempty"`
}

func (d *AppDao) AllFeatured() (apps []AppFeatured, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT %s, %s, %s FROM %s WHERE %s = ?", models.COLUMN_Apps_Name,
		models.COLUMN_Apps_ChannelId, models.COLUMN_Apps_AppId, models.TABLE_Apps, models.COLUMN_Apps_Featured),
		AppFeaturedYes).QueryRows(&apps)
	if err != nil {
		common.LogFuncError("%v", err)
	}

	return
}

// get all app types
func (d *AppDao) AllTypes() (types []models.AppType, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_AppType).All(&types)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// get all app types
func (d *AppDao) AllChannels() (channels []models.AppChannel, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_AppChannel).All(&channels)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *AppDao) ChangePos(isCreate bool, id, channelID, newPos uint32) error {
	var sql string
	if isCreate {
		sql = fmt.Sprintf("UPDATE %s SET %s=%s+1 WHERE %s=? and %s>=? ", models.TABLE_Apps, models.COLUMN_Apps_Position,
			models.COLUMN_Apps_Position, models.COLUMN_Apps_ChannelId, models.COLUMN_Apps_Position)
		_, err := d.Orm.Raw(sql, channelID, newPos).Exec()
		if err != nil {
			common.LogFuncError("error: %v", err)
			return err
		}
	} else {
		app, err := d.QueryAppById(id)
		if err != nil {
			common.LogFuncError("error: %v", err)
			return err
		}
		if app.ChannelId == channelID {
			if app.Position < newPos {
				sql = fmt.Sprintf("UPDATE %s SET %s=%s-1 WHERE %s=? and (%s>? and %s<=?)", models.TABLE_Apps, models.COLUMN_Apps_Position,
					models.COLUMN_Apps_Position, models.COLUMN_Apps_ChannelId, models.COLUMN_Apps_Position, models.COLUMN_Apps_Position)
				_, err = d.Orm.Raw(sql, channelID, app.Position, newPos).Exec()
			} else if app.Position > newPos {
				sql = fmt.Sprintf("UPDATE %s SET %s=%s+1 WHERE %s=? and (%s>=? and %s<?)", models.TABLE_Apps, models.COLUMN_Apps_Position,
					models.COLUMN_Apps_Position, models.COLUMN_Apps_ChannelId, models.COLUMN_Apps_Position, models.COLUMN_Apps_Position)
				_, err = d.Orm.Raw(sql, channelID, newPos, app.Position).Exec()
			} else {
				return nil
			}
			if err != nil {
				common.LogFuncError("error: %v", err)
				return err
			}
		} else {
			sql = fmt.Sprintf("UPDATE %s SET %s=%s+1 WHERE %s=? and %s>=? ", models.TABLE_Apps, models.COLUMN_Apps_Position,
				models.COLUMN_Apps_Position, models.COLUMN_Apps_ChannelId, models.COLUMN_Apps_Position)
			_, err := d.Orm.Raw(sql, channelID, newPos).Exec()
			if err != nil {
				common.LogFuncError("error: %v", err)
				return err
			}
		}
	}
	return nil
}

func (d *AppDao) ChannelApp(channelId uint32) (apps []models.Apps, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_Apps).Filter(models.COLUMN_Apps_ChannelId, channelId).OrderBy(models.COLUMN_Apps_Position).All(&apps)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (d *AppDao) GetAppNameByChannelId(channelId uint32) (mapAppName orm.Params, err error) {
	mapAppName = make(orm.Params)
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s=?", models.COLUMN_Apps_AppId,
		models.COLUMN_Apps_Name, models.TABLE_Apps, models.COLUMN_Apps_ChannelId), channelId).
		RowsToMap(&mapAppName, models.COLUMN_Apps_AppId, models.COLUMN_Apps_Name)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("error:%v", err)
		return
	}

	return
}
