package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/admin/models"
)

type EndPointDao struct {
	common.BaseDao
}

func NewEndPointDao(db string) *EndPointDao {
	return &EndPointDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var EndPointDaoEntity *EndPointDao

//获取申述列表
func (d *EndPointDao) GetAll() (data []models.Endpoint, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_Endpoint).All(&data)
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

func (d *EndPointDao) Create(domain string) (err error) {
	endpoint := &models.Endpoint{
		Endpoint: domain,
		Ctime:    common.NowInt64MS(),
		Utime:    common.NowInt64MS(),
	}

	_, err = d.Orm.Insert(endpoint)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *EndPointDao) Update(id uint32, domain string) (err error) {
	endpoint := &models.Endpoint{
		Id:       id,
		Endpoint: domain,
		Utime:    common.NowInt64MS(),
	}
	_, err = d.Orm.Update(endpoint, models.COLUMN_Endpoint_Endpoint, models.COLUMN_Endpoint_Utime)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

func (d *EndPointDao) DelById(id uint32) (err error) {
	endpoint := &models.Endpoint{Id: id}
	_, err = d.Orm.Delete(endpoint, models.COLUMN_Endpoint_Id)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}
