package dao

import (
	"common"
	"fmt"
	"utils/admin/models"
)

type MenuConfDao struct {
	common.BaseDao
}

func NewMenuConfDao(db string) *MenuConfDao {
	return &MenuConfDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var MenuConfDaoEntity *MenuConfDao

func (d *MenuConfDao) Add(conf *models.MenuConf) (err error) {
	_, err = d.Orm.Insert(conf)
	if err != nil {
		common.LogFuncError("add menuConf db err %v ", err)
		return
	}
	return
}

func (d *MenuConfDao) FindAll() (confs []*models.MenuConf, err error) {
	querySql := fmt.Sprintf("select * from %s order by %s", models.TABLE_MenuConf, models.COLUMN_MenuConf_OrderId)
	_, err = d.Orm.Raw(querySql).QueryRows(&confs)
	if err != nil {
		common.LogFuncError("FindAll db err %v ", err)
		return
	}
	return
}
func (d *MenuConfDao) FindChildrenMenu(id uint64) (menus []*models.MenuConf, err error) {
	querySql := fmt.Sprintf("select * from %s where %s=?", models.TABLE_MenuConf, models.COLUMN_MenuConf_PId)
	_, err = d.Orm.Raw(querySql, id).QueryRows(&menus)
	if err != nil {
		common.LogFuncError("FindChildrenMenu db err %v ", err)
		return
	}
	return
}

func (d *MenuConfDao) FindByMenuIds(ids []uint64) (menus []*models.MenuConf, err error) {
	qs := d.Orm.QueryTable(models.TABLE_MenuConf)
	inSql := fmt.Sprintf("%s__in", models.COLUMN_MenuConf_Id)
	_, err = qs.Filter(inSql, ids).OrderBy(models.COLUMN_MenuConf_OrderId).All(&menus)
	if err != nil {
		common.LogFuncError("FindByMenuIds db err %v ", err)
		return
	}
	return
}

func (d *MenuConfDao) Delete(id uint64) (err error) {
	querySql := fmt.Sprintf("delete from %s where %s=? or %s=?", models.TABLE_MenuConf, models.COLUMN_MenuConf_Id, models.COLUMN_MenuConf_PId)
	_, err = d.Orm.Raw(querySql, id, id).Exec()
	if err != nil {
		common.LogFuncError("Delete menu db err %v ", err)
		return
	}
	return
}

func (d *MenuConfDao) Update(conf models.MenuConf) (err error) {
	_, err = d.Orm.Update(&conf)
	if err != nil {
		common.LogFuncError("update menu db err %v ", err)
		return
	}
	return
}
