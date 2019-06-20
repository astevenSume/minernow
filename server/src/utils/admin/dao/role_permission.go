package dao

import (
	"admin/utils"
	"common"
	"fmt"
	"strconv"
	"strings"
	"utils/admin/models"
)

type RolePermissionDao struct {
	common.BaseDao
}

func NewRolePermissionDao(db string) *RolePermissionDao {
	return &RolePermissionDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var RolePermissionDaoEntity *RolePermissionDao

//获取角色所拥有权限
func (d *RolePermissionDao) GetRolePermissions(id uint64) ([]string, error) {
	var permission []string
	_, err := d.Orm.Raw("select slug from ((select * from role_permission where roleid = ?) as T1 LEFT JOIN permission as T2 on T1.permissionid=T2.id)", id).QueryRows(&permission)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return permission, err
	}

	return permission, nil
}

//获取角色所拥有权限
func (d *RolePermissionDao) GetRolePermissionInfo(id uint64) (mapPermission map[string]uint64, err error) {
	var name []string
	var ids []uint64
	_, err = d.Orm.Raw("select permissionid,slug from ((select * from role_permission where roleid = ?) as T1 LEFT JOIN permission as T2 on T1.permissionid=T2.id)", id).QueryRows(&ids, &name)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}
	c := len(name)
	if c != len(ids) {
		common.LogFuncError("query len = 0")
		return
	}
	mapPermission = make(map[string]uint64)
	for i := 0; i < c; i++ {
		mapPermission[name[i]] = ids[i]
	}

	return
}

//获取用户角色信息
func (d *RolePermissionDao) GetUserRolesInfo(ids []string) (rolePermission []RolePermission, err error) {
	strId := strings.Join(ids, CommaSpace)
	sql := fmt.Sprintf("select * from ((select * from role_permission where roleid in(%s)) as T1 LEFT JOIN permission as T2 on T1.permissionid=T2.id)", strId)
	_, err = d.Orm.Raw(sql).QueryRows(&rolePermission)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

//更新角色权限
func (d *RolePermissionDao) UpdateRolePermissions(roleId uint64, ids []string) ([]string, error) {
	//角色权限表
	rolePermission := models.RolePermission{
		RoleId: roleId,
	}
	_, err := d.Orm.Delete(&rolePermission, models.COLUMN_RolePermission_RoleId)
	if err != nil {
		common.LogFuncDebug("mysql error:%v", err)
	}

	var permission []string
	err = d.InsertMulRolePermissions(roleId, ids)
	if err != nil {
		return permission, err
	}

	permission, err = d.GetRolePermissions(roleId)
	if err != nil {
		return permission, err
	}

	return permission, nil
}

//添加权限
func (d *RolePermissionDao) AddRolePermissions(roleID uint64, addPermissions []string) ([]string, error) {
	permissions, err := d.GetRolePermissions(roleID)
	if err != nil {
		return permissions, err
	}

	isNew := false
	var rolePermission []models.RolePermission
	for _, p := range addPermissions {
		if pid, ok := MapPermission[p]; ok && utils.DetermineIndex(permissions, p) < 0 {
			permissions = append(permissions, p)
			r := models.RolePermission{
				RoleId:      roleID,
				Pemissionid: pid,
			}
			rolePermission = append(rolePermission, r)
			isNew = true
		}
	}
	if isNew {
		_, err := d.Orm.InsertMulti(InsertMulCount, rolePermission)
		if err != nil {
			return permissions, err
		}
	}

	return permissions, nil
}

//删除权限
func (d *RolePermissionDao) DelRolePermission(rolePermission *models.RolePermission, cols ...string) error {
	if rolePermission == nil {
		return ErrParam
	}
	_, err := d.Orm.Delete(rolePermission, cols...)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

//删除权限
func (d *RolePermissionDao) DelRolePermissions(roleID uint64, permissions []string) ([]string, error) {
	var other []string
	mapPer, err := d.GetRolePermissionInfo(roleID)
	if err != nil {
		return other, err
	}

	var ids []string
	for _, right := range permissions {
		if id, ok := mapPer[right]; ok {
			ids = append(ids, strconv.Itoa(int(id)))
			delete(mapPer, right)
		}
	}

	strId := strings.Join(ids, CommaSpace)
	sql := fmt.Sprintf("delete from %s where %s=? and %s in (%s)", models.TABLE_RolePermission, models.COLUMN_RolePermission_RoleId, models.COLUMN_RolePermission_Pemissionid, strId)
	_, err = d.Orm.Raw(sql, roleID).Exec()
	if err != nil {
		return other, err
	}

	for k := range mapPer {
		other = append(other, k)
	}

	return other, nil
}

//多条插入
func (d *RolePermissionDao) InsertMulRolePermissions(roleId uint64, slugs []string) error {
	var ids []uint64
	for _, s := range slugs {
		if id, ok := MapPermission[s]; ok {
			ids = append(ids, id)
		}
	}

	var rolePermission []models.RolePermission
	for _, pid := range ids {
		r := models.RolePermission{
			RoleId:      roleId,
			Pemissionid: pid,
		}
		rolePermission = append(rolePermission, r)
	}
	_, err := d.Orm.InsertMulti(InsertMulCount, rolePermission)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return err
	}
	return nil
}
