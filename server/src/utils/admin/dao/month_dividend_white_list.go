package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
	"utils/admin/models"
)

type MonthDividendWhiteListDao struct {
	common.BaseDao
}

func NewMonthDividendWhiteListDao(db string) *MonthDividendWhiteListDao {
	return &MonthDividendWhiteListDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var MonthDividendWhiteListDaoEntity *MonthDividendWhiteListDao

//获取代理白名单档位
func (d *MonthDividendWhiteListDao) QueryById(id uint32) (MonthDividendWhiteList *models.MonthDividendWhiteList, err error) {
	MonthDividendWhiteList = &models.MonthDividendWhiteList{
		Id: id,
	}

	err = d.Orm.Read(MonthDividendWhiteList, models.COLUMN_MonthDividendWhiteList_Id)
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

func (d *MonthDividendWhiteListDao) GetAllMap() (monthDividendWhiteListMap map[uint32]*models.MonthDividendWhiteList, err error) {
	query_sql := fmt.Sprintf("select * from %s", models.TABLE_MonthDividendWhiteList)
	monthDividendWhiteList := make([]*models.MonthDividendWhiteList, 0)
	_, err = d.Orm.Raw(query_sql).QueryRows(&monthDividendWhiteList)
	if err != nil {
		common.LogFuncError("GetAllMap fail %v", err)
		return
	}
	monthDividendWhiteListMap = make(map[uint32]*models.MonthDividendWhiteList, 0)
	for _, m := range monthDividendWhiteList {
		monthDividendWhiteListMap[m.Id] = m
	}
	return
}

//分页查询
func (d *MonthDividendWhiteListDao) QueryByPage(page int, perPage int) (total int, MonthDividendWhiteList []models.MonthDividendWhiteList, err error) {
	qs := d.Orm.QueryTable(models.TABLE_MonthDividendWhiteList)
	if qs == nil {
		common.LogFuncError("QueryTable fail")
		return
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
	_, err = qs.Limit(perPage, start).All(&MonthDividendWhiteList)
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

func (d *MonthDividendWhiteListDao) Create(dividendRatio int32, name string) (MonthDividendWhiteList models.MonthDividendWhiteList, err error) {
	now := time.Now().Unix()
	MonthDividendWhiteList = models.MonthDividendWhiteList{
		Name:          name,
		DividendRatio: dividendRatio,
		Ctime:         now,
		Utime:         now,
	}

	var id int64
	id, err = d.Orm.Insert(&MonthDividendWhiteList)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}
	MonthDividendWhiteList.Id = uint32(id)

	return
}

func (d *MonthDividendWhiteListDao) UpdateById(id uint32, dividendRatio int32, name string) (models.MonthDividendWhiteList, error) {
	MonthDividendWhiteList := models.MonthDividendWhiteList{Id: id}
	err := d.Orm.Read(&MonthDividendWhiteList, models.COLUMN_MonthDividendWhiteList_Id)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return MonthDividendWhiteList, err
	}

	MonthDividendWhiteList.Name = name
	MonthDividendWhiteList.Utime = time.Now().Unix()
	MonthDividendWhiteList.DividendRatio = dividendRatio

	_, err = d.Orm.Update(&MonthDividendWhiteList)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return MonthDividendWhiteList, err
	}

	return MonthDividendWhiteList, nil
}

func (d *MonthDividendWhiteListDao) DelById(id uint32) error {
	MonthDividendWhiteList := &models.MonthDividendWhiteList{
		Id: id,
	}

	_, err := d.Orm.Delete(MonthDividendWhiteList, models.COLUMN_MonthDividendWhiteList_Id)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

func (d *MonthDividendWhiteListDao) All() (list []models.MonthDividendWhiteList, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_MonthDividendWhiteList).All(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}
