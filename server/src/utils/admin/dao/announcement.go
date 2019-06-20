package dao

import (
	"common"
	"fmt"
	"utils/admin/models"

	"github.com/astaxie/beego/orm"
)

//公告类型
const (
	AnnouncementTypeNil = iota
	AnnouncementTypeSys //系统公告
	AnnouncementTypeGame
	AnnouncementTypeMax
)

type AnnouncementDao struct {
	common.BaseDao
}

func NewAnnouncementDao(db string) *AnnouncementDao {
	return &AnnouncementDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AnnouncementDaoEntity *AnnouncementDao

//获取公告信息
func (d *AnnouncementDao) QueryById(id uint32) (announcement *models.Announcement, err error) {
	announcement = &models.Announcement{
		Id: id,
	}

	err = d.Orm.Read(announcement, models.COLUMN_Announcement_Id)
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
func (d *AnnouncementDao) QueryByPage(aType int8, stime, etime int64, page int, perPage int) (total int, announcements []models.Announcement, err error) {
	qs := d.Orm.QueryTable(models.TABLE_Announcement)
	if qs == nil {
		common.LogFuncError("QueryTable fail")
		return
	}

	if aType > AnnouncementTypeNil && aType < AnnouncementTypeMax {
		qs = qs.Filter(models.COLUMN_Announcement_Type, aType)
	}
	if stime > 0 {
		qs = qs.Filter(models.COLUMN_Announcement_Stime+"__lte", stime)
	}
	if etime > 0 {
		qs = qs.Filter(models.COLUMN_Announcement_Etime+"__gte", etime)
	}

	var count int64
	count, err = qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}
	total = int(count)

	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > count {
		err = nil
		return
	}
	_, err = qs.Limit(perPage, start).All(&announcements)
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

func (d *AnnouncementDao) Create(aType int8, title, content string, stime, etime int64) (announcement models.Announcement, err error) {
	if title == "" || content == "" || aType <= AnnouncementTypeNil || aType >= AnnouncementTypeMax {
		err = ErrParam
		return
	}
	announcement = models.Announcement{
		Type:    aType,
		Title:   title,
		Content: content,
		Ctime:   common.NowInt64MS(),
		Utime:   common.NowInt64MS(),
		Stime:   stime,
		Etime:   etime,
	}

	var id int64
	id, err = d.Orm.Insert(&announcement)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}
	announcement.Id = uint32(id)

	return
}

func (d *AnnouncementDao) Update(aType int8, id uint32, title, content string, stime, etime int64) (announcement models.Announcement, err error) {
	if title == "" || content == "" || aType <= AnnouncementTypeNil || aType >= AnnouncementTypeMax {
		err = ErrParam
		return
	}

	cols := []string{
		models.COLUMN_Announcement_Type,
		models.COLUMN_Announcement_Utime,
		models.COLUMN_Announcement_Title,
		models.COLUMN_Announcement_Content,
		models.COLUMN_Announcement_Stime,
		models.COLUMN_Announcement_Etime,
	}
	announcement.Id = id
	announcement.Type = aType
	announcement.Title = title
	announcement.Content = content
	announcement.Utime = common.NowInt64MS()
	announcement.Stime = stime
	announcement.Etime = etime

	_, err = d.Orm.Update(&announcement, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

func (d *AnnouncementDao) DelById(id uint32) error {
	announcement := &models.Announcement{
		Id: id,
	}

	_, err := d.Orm.Delete(announcement, models.COLUMN_Announcement_Id)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

// EffectiveAnnouncement 当前有效的 announcement
func (d *AnnouncementDao) EffectiveAnnouncement(aType int8) (announcement []models.Announcement, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_Announcement).
		Filter(models.COLUMN_Announcement_Type, aType).
		Filter(fmt.Sprintf("%s__lte", models.COLUMN_Announcement_Stime), common.NowInt64MS()).
		Filter(fmt.Sprintf("%s__gte", models.COLUMN_Announcement_Etime), common.NowInt64MS()).
		All(&announcement)
	if err != nil {
		common.LogFuncError("%v", err)
	}

	return
}
