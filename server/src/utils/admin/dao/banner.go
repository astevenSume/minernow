package dao

import (
	"common"
	"fmt"
	"utils/admin/models"

	"github.com/astaxie/beego/orm"
)

//上线状态
const (
	BannerStatusNil       = iota
	BannerStatusPending   //等待发布
	BannerStatusPublished //已发布
	BannerStatusExpired   //已过期
	BannerStatusMax
)

type BannerDao struct {
	common.BaseDao
}

func NewBannerDao(db string) *BannerDao {
	return &BannerDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var BannerDaoEntity *BannerDao

//获取广告信息
func (d *BannerDao) QueryById(id uint32) (banner *models.Banner, err error) {
	banner = &models.Banner{
		Id: id,
	}

	err = d.Orm.Read(banner, models.COLUMN_Banner_Id)
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
func (d *BannerDao) QueryByPage(status int8, page int, perPage int) (total int, banners []models.Banner, err error) {
	qs := d.Orm.QueryTable(models.TABLE_Banner)
	if qs == nil {
		common.LogFuncError("QueryTable fail")
		return
	}

	if status > BannerStatusNil && status < BannerStatusMax {
		qs = qs.Filter(models.COLUMN_Banner_Status, status)
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
	_, err = qs.Limit(perPage, start).All(&banners)
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

func (d *BannerDao) Create(subject, image, url string, stime, etime int64) (banner models.Banner, err error) {
	if subject == "" || image == "" {
		err = ErrParam
		return
	}
	banner = models.Banner{
		Subject: subject,
		Image:   image,
		Url:     url,
		Status:  BannerStatusPending,
		Ctime:   common.NowInt64MS(),
		Utime:   common.NowInt64MS(),
		Stime:   stime,
		Etime:   etime,
	}

	var id int64
	id, err = d.Orm.Insert(&banner)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}
	banner.Id = uint32(id)

	return
}

func (d *BannerDao) Update(id uint32, subject, image, url string, stime, etime int64) (banner models.Banner, err error) {
	if subject == "" || image == "" {
		err = ErrParam
		return
	}

	cols := []string{
		models.COLUMN_Banner_Utime,
		models.COLUMN_Banner_Subject,
		models.COLUMN_Banner_Image,
		models.COLUMN_Banner_Stime,
		models.COLUMN_Banner_Etime,
	}
	banner.Id = id
	banner.Subject = subject
	banner.Image = image
	banner.Utime = common.NowInt64MS()
	banner.Stime = stime
	banner.Etime = etime
	if url != "" {
		banner.Url = url
		cols = append(cols, models.COLUMN_Banner_Url)
	}

	_, err = d.Orm.Update(&banner, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

func (d *BannerDao) DelById(id uint32) error {
	banner := &models.Banner{
		Id: id,
	}

	_, err := d.Orm.Delete(banner, models.COLUMN_Banner_Id)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

//更新上线状态
func (d *BannerDao) UpdateStatus(id uint32, status int8) error {
	banner := &models.Banner{
		Id:     id,
		Status: status,
		Utime:  common.NowInt64MS(),
	}

	_, err := d.Orm.Update(banner, models.COLUMN_Banner_Status, models.COLUMN_Banner_Utime)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

// get all
func (d *BannerDao) All() (banner []models.Banner, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_Banner).All(&banner)
	if err != nil {
		common.LogFuncError("%v", err)
	}

	return
}

// CurrentEffective 当前生效的 banner
func (d *BannerDao) EffectiveBanners() (banner []models.Banner, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_Banner).
		Filter(models.COLUMN_Banner_Status, BannerStatusPublished).
		Filter(fmt.Sprintf("%s__lte", models.COLUMN_Banner_Stime), common.NowInt64MS()).
		Filter(fmt.Sprintf("%s__gte", models.COLUMN_Banner_Etime), common.NowInt64MS()).
		All(&banner)
	if err != nil {
		common.LogFuncError("%v", err)
	}

	return
}
