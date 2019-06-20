package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/admin/models"
)

type AppTypeDao struct {
	common.BaseDao
}

func NewAppTypeDao(db string) *AppTypeDao {
	return &AppTypeDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AppTypeDaoEntity *AppTypeDao

//获取app类型信息
func (d *AppTypeDao) QueryById(id uint32) (appType *models.AppType, err error) {
	appType = &models.AppType{
		Id: id,
	}

	err = d.Orm.Read(appType, models.COLUMN_AppType_Id)
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

//所有类型
func (d *AppTypeDao) QueryAll() (appTypes []models.AppType, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_AppType).All(&appTypes)
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

func (d *AppTypeDao) Create(id uint32, name, desc string) (isNew bool, appType models.AppType, err error) {
	appType = models.AppType{
		Id:    id,
		Name:  name,
		Desc:  desc,
		Ctime: common.NowInt64MS(),
		Utime: common.NowInt64MS(),
	}

	isNew, _, err = d.Orm.ReadOrCreate(&appType, models.COLUMN_AppType_Id)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}

	return
}

func (d *AppTypeDao) Update(id uint32, name, desc string) (appType models.AppType, err error) {
	appType = models.AppType{Id: id}
	err = d.Orm.Read(&appType, models.COLUMN_AppType_Id)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}

	appType.Name = name
	appType.Desc = desc
	appType.Utime = common.NowInt64MS()
	_, err = d.Orm.Update(&appType, models.COLUMN_AppType_Name, models.COLUMN_AppType_Desc, models.COLUMN_AppType_Utime)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

func (d *AppTypeDao) DelById(id uint32) error {
	appType := &models.AppType{
		Id: id,
	}

	_, err := d.Orm.Delete(appType, models.COLUMN_AppType_Id)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}
