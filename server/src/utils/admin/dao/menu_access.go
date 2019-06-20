package dao

import (
	"common"
	"fmt"
	"utils/admin/models"
)

type MenuAccessDao struct {
	common.BaseDao
}

func NewMenuAccessDao(db string) *MenuAccessDao {
	return &MenuAccessDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var MenuAccessDaoEntity *MenuAccessDao

func (d *MenuAccessDao) FindAll() (datas []*models.MenuAccess, err error) {
	querySql := fmt.Sprintf("select * from %s ", models.TABLE_MenuAccess)
	_, err = d.Orm.Raw(querySql).QueryRows(&datas)
	if err != nil {
		common.LogFuncError("FindAll delete db err %v", err)
	}
	return
}
func (d *MenuAccessDao) FindByRoleId(roleId []uint64) (datas []*models.MenuAccess, err error) {
	qs := d.Orm.QueryTable(models.TABLE_MenuAccess)
	inStr := fmt.Sprintf("%s__in", models.COLUMN_MenuAccess_RoleId)
	_, err = qs.Filter(inStr, roleId).All(&datas)
	if err != nil {
		common.LogFuncError("FindByRoleId delete db err %v", err)
	}
	return
}
func (d *MenuAccessDao) Delete(roleId uint64) (err error) {
	querySql := fmt.Sprintf("delete from %s where %s=?", models.TABLE_MenuAccess, models.COLUMN_MenuAccess_RoleId)
	_, err = d.Orm.Raw(querySql, roleId).Exec()
	if err != nil {
		common.LogFuncError("MenuAccessDaoEntity delete db err %v", err)
	}
	return
}

func (d *MenuAccessDao) Insert(rows []*models.MenuAccess) (err error) {
	_, err = d.Orm.InsertMulti(100, &rows)
	if err != nil {
		common.LogFuncError("MenuAccessDaoEntity Insert db err %v", err)
	}
	return
}
