package dao

import (
	"common"
	"utils/admin/models"

	"github.com/astaxie/beego/orm"
)

//预警类型
const (
	ConfigWarningTypeNil                 = iota
	ConfigWarningTypeUsdt                //usdt转账
	ConfigWarningTypeUsdtAccountTampered //usdt账号遭篡改
	ConfigWarningTypeUsdtPriKeyLeakage   //usdt私钥泄露
	ConfigWarningTypeMax
)

type ConfigWarningDao struct {
	common.BaseDao
}

func NewConfigWarningDao(db string) *ConfigWarningDao {
	return &ConfigWarningDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ConfigWarningDaoEntity *ConfigWarningDao

//分页查询
func (d *ConfigWarningDao) QueryPageConfigWarning(page int, perPage int, cType int8, mobile string) (total int64, list []models.ConfigWarning, err error) {
	qs := d.Orm.QueryTable(models.TABLE_ConfigWarning)
	if qs == nil {
		common.LogFuncError("mysql_err:TABLE_Smscodes fail")
		return
	}
	if cType > ConfigWarningTypeNil && cType < ConfigWarningTypeMax {
		qs = qs.Filter(models.COLUMN_ConfigWarning_Type, cType)
	}
	if len(mobile) > 0 {
		qs = qs.Filter(models.COLUMN_ConfigWarning_Mobile, mobile)
	}

	total, err = qs.Count()
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > total {
		return
	}
	qs = qs.OrderBy(models.COLUMN_ConfigWarning_Type)
	_, err = qs.Limit(perPage, start).All(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *ConfigWarningDao) Create(cType int8, smsType int8, nationalCode, mobile string) (err error) {
	if cType <= ConfigWarningTypeNil || cType >= ConfigWarningTypeMax {
		err = ErrParam
		return
	}

	configWarning := &models.ConfigWarning{
		Type:         cType,
		NationalCode: nationalCode,
		Mobile:       mobile,
		SmsType:      smsType,
	}
	_, err = d.Orm.Insert(configWarning)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *ConfigWarningDao) Update(cType int8, smsType int8, id uint32, nationalCode, mobile string) (err error) {
	configWarning := &models.ConfigWarning{
		Id:           id,
		Type:         cType,
		NationalCode: nationalCode,
		Mobile:       mobile,
		SmsType:      smsType,
	}
	_, err = d.Orm.Update(configWarning, models.COLUMN_ConfigWarning_Type, models.COLUMN_ConfigWarning_NationalCode,
		models.COLUMN_ConfigWarning_Mobile, models.COLUMN_ConfigWarning_SmsType)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

func (d *ConfigWarningDao) DelById(id uint32) (err error) {
	configWarning := &models.ConfigWarning{Id: id}
	_, err = d.Orm.Delete(configWarning, models.COLUMN_ConfigWarning_Id)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

func (d *ConfigWarningDao) GetConfigWarning(cType int8) (configWarning []models.ConfigWarning, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_ConfigWarning).Filter(models.COLUMN_ConfigWarning_Type, cType).All(&configWarning)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}
