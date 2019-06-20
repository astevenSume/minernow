package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/admin/models"
)

type SystemMessageMethodDao struct {
	common.BaseDao
}

func NewSystemMessageMethodDao(db string) *SystemMessageMethodDao {
	return &SystemMessageMethodDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var SystemMessageMethodDaoEntity *SystemMessageMethodDao

// add SystemMessage
func (d *SystemMessageMethodDao) AddSysMsg(key string, buyer string, seller string, admin string) (sm *models.SystemMessage, err error) {
	sm = &models.SystemMessage{
		Key:    key,
		Buyer:  buyer,
		Seller: seller,
		Admin:  admin,
		Ctime:  common.NowInt64MS(),
	}
	_, err = d.Orm.Insert(sm)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

// del SystemMessage
func (d *SystemMessageMethodDao) DelSysMsg(id uint64) (err error) {
	_, err = d.Orm.Delete(&models.SystemMessage{Id: id})
	if err != nil {
		common.LogFuncError("delete id= %d SystemMessage fail error is %v", id, err)
		return
	}
	return
}

// update SystemMessage
func (d *SystemMessageMethodDao) UpdateSysMsg(id uint64, key string, buyer string, seller string, admin string) (sm *models.SystemMessage, err error) {
	sm = &models.SystemMessage{
		Id: id,
	}
	if err = d.Orm.Read(sm); err != nil {
		if err == orm.ErrNoRows {
			err = nil
		}
		return
	}
	sm.Buyer = buyer
	sm.Seller = seller
	sm.Admin = admin
	sm.Utime = common.NowInt64MS()
	if _, err = d.Orm.Update(sm); err != nil {
		common.LogFuncError("update SystemMessage id = %d fail error is %v", id, err)
		return
	}
	return

}

// query SystemMessage by id
func (d *SystemMessageMethodDao) QuerySysMsgById(id uint64) (sm *models.SystemMessage, err error) {
	if id <= 0 {
		common.LogFuncCritical("query SystemMessage id <= 0")
		return
	}
	sm = &models.SystemMessage{
		Id: id,
	}

	err = d.Orm.Read(sm, models.ATTRIBUTE_AdminUser_Id)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("query SystemMessage by id error is %v", err)
		return
	}
	return
}

// query SystemMessage all
func (d *SystemMessageMethodDao) QuerySysMsgAll() (list []*models.SystemMessage, err error) {

	querySeter := d.Orm.QueryTable(models.TABLE_SystemMessage)
	if _, err = querySeter.All(&list); err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("query SystemMessage all error is %v", err)
		return
	}
	return
}

// 判断某Id的SysMsg 存在不存在
func (d *SystemMessageMethodDao) ReadSysMsgById(id uint64) (err error) {
	if id <= 0 {
		common.LogFuncCritical("read SystemMessage id <= 0")
		return
	}
	sm := &models.SystemMessage{
		Id: id,
	}
	err = d.Orm.Read(sm, models.ATTRIBUTE_SystemMessage_Id)

	if err != nil {
		common.LogFuncError("read SystemMessage by id error is %v", err)
		return
	}
	return
}
