package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"strings"
	"utils/admin/models"
)

var (
	MapPermission map[string]uint64
)

type PermissonDao struct {
	common.BaseDao
}

var PermissonDaoEntity *PermissonDao

func NewPermissonDao(db string) *PermissonDao {
	return &PermissonDao{
		BaseDao: common.NewBaseDao(db),
	}
}

//获取角色slug->id映射
func (d *PermissonDao) LoadAllPermission() error {
	MapPermission = make(map[string]uint64)
	var permission []models.Permission
	_, err := d.Orm.QueryTable(models.TABLE_Permission).All(&permission)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	for _, p := range permission {
		MapPermission[p.Slug] = p.Id
	}

	return nil
}

//查询权限
func (d *PermissonDao) QueryPermission(permission *models.Permission, cols ...string) error {
	if permission == nil {
		return ErrParam
	}

	err := d.Orm.Read(permission, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	return nil
}

//创建权限
func (d *PermissonDao) CreatePermissionBySlug(permission *models.Permission) (uint64, error) {
	if permission == nil {
		return 0, ErrParam
	}
	permission.Slug = strings.ToUpper(permission.Slug)
	_, id, err := d.Orm.ReadOrCreate(permission, models.COLUMN_Permission_Slug)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return 0, err
	}
	MapPermission[permission.Slug] = uint64(id)

	return uint64(id), nil
}

//更新权限
func (d *PermissonDao) UpdatePermissionBySlug(newPermission *models.Permission, cols ...string) error {
	if newPermission == nil {
		return ErrParam
	}
	newPermission.Slug = strings.ToUpper(newPermission.Slug)

	permission := &models.Permission{Id: newPermission.Id}
	err := d.QueryPermission(permission, models.COLUMN_Permission_Id)
	if err != nil {
		return err
	}
	if _, ok := MapPermission[permission.Slug]; ok {
		delete(MapPermission, permission.Slug)
	}

	permission = newPermission
	_, err = d.Orm.Update(permission, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	MapPermission[permission.Slug] = permission.Id

	return nil
}

//删除权限
func (d *PermissonDao) DelPermission(id uint64) error {
	permission := models.Permission{Id: id}
	_, err := d.Orm.Delete(&permission)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	for k, v := range MapPermission {
		if v == id {
			delete(MapPermission, k)
		}
	}

	rolePermission := &models.RolePermission{
		Pemissionid: id,
	}
	err = RolePermissionDaoEntity.DelRolePermission(rolePermission, models.COLUMN_RolePermission_Pemissionid)
	if err != nil {
		common.LogFuncDebug("mysql error:%v", err)
	}

	return nil
}

//分页查询
func (d *PermissonDao) QueryPagePermission(page int, perPage int) ([]models.Permission, PageInfo, error) {
	var data []models.Permission
	var info PageInfo
	qs := d.Orm.QueryTable(models.TABLE_Permission)
	if qs == nil {
		common.LogFuncError("mysql_err:Permission fail")
		return data, info, ErrParam
	}

	count, err := qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			return data, info, nil
		}
		common.LogFuncError("mysql error:%v", err)
		return data, info, err
	}
	info.Total = int(count)
	info.Page = page
	info.Limit = perPage
	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > count {
		return data, info, ErrParam
	}
	_, err = qs.Limit(perPage, start).All(&data)
	if err != nil {
		if err == orm.ErrNoRows {
			return data, info, nil
		}
		common.LogFuncError("mysql error:%v", err)
		return data, info, err
	}

	return data, info, nil
}
