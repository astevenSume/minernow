package dao

import (
	"common"
	"fmt"
	"utils/admin/models"

	"github.com/astaxie/beego/orm"
)

type AppWhiteListDao struct {
	common.BaseDao
}

func NewAppWhiteListDao(db string) *AppWhiteListDao {
	return &AppWhiteListDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AppWhiteListDaoEntity *AppWhiteListDao

//获取app白名单信息
func (d *AppWhiteListDao) QueryById(id uint32) (appWhitelist *models.AppWhitelist, err error) {
	appWhitelist = &models.AppWhitelist{
		Id: id,
	}

	err = d.Orm.Read(appWhitelist, models.COLUMN_AppWhitelist_Id)
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
func (d *AppWhiteListDao) QueryPageCondition(channelID int8, page int, perPage int) (total int64, appWhitelists []models.AppWhitelist, err error) {
	qs := d.Orm.QueryTable(models.TABLE_AppWhitelist)
	if qs == nil {
		common.LogFuncError("QueryTable fail")
		return
	}

	if channelID > 0 {
		qs = qs.Filter(models.COLUMN_AppWhitelist_ChannelId, channelID)
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
	qs = qs.OrderBy(models.COLUMN_AppWhitelist_ChannelId)
	_, err = qs.Limit(perPage, start).All(&appWhitelists)
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

func (d *AppWhiteListDao) Create(channelId uint32, appId string) (appWhitelist models.AppWhitelist, err error) {
	appWhitelist = models.AppWhitelist{
		ChannelId: channelId,
		AppId:     appId,
		Ctime:     common.NowInt64MS(),
	}

	var id int64
	id, err = d.Orm.Insert(&appWhitelist)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}
	appWhitelist.Id = uint32(id)

	return
}

func (d *AppWhiteListDao) DelById(id uint32) error {
	appWhitelist := &models.AppWhitelist{
		Id: id,
	}

	_, err := d.Orm.Delete(appWhitelist, models.COLUMN_AppWhitelist_Id)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

//应用是否是白名单
func (d *AppWhiteListDao) IsWhite(channelId uint32, appId string) (isWhite bool) {
	appWhitelist := &models.AppWhitelist{
		ChannelId: channelId,
		AppId:     appId,
	}

	isWhite = false
	err := d.Orm.Read(appWhitelist, models.COLUMN_AppWhitelist_ChannelId, models.COLUMN_AppWhitelist_AppId)
	if err != nil {
		return
	}
	isWhite = true
	return
}

//应用白名单
func (d *AppWhiteListDao) AllByChannel(channelId uint32) (mapWhite orm.Params, err error) {
	//d.Orm.QueryTable(models.TABLE_AppWhitelist).Filter(models.COLUMN_AppWhitelist_ChannelId, channelId).All()
	sql := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s=?", models.COLUMN_AppWhitelist_AppId,
		models.COLUMN_AppWhitelist_ChannelId, models.TABLE_AppWhitelist, models.COLUMN_AppWhitelist_ChannelId)
	mapWhite = make(orm.Params)
	_, err = d.Orm.Raw(sql, channelId).RowsToMap(&mapWhite, models.COLUMN_AppWhitelist_AppId, models.COLUMN_AppWhitelist_ChannelId)
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
