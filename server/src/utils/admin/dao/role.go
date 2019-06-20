package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"strconv"
	"utils/admin/models"
)

type RoleDao struct {
	common.BaseDao
}

func NewRoleDao(db string) *RoleDao {
	return &RoleDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var RoleDaoEntity *RoleDao

//创建角色
func (d *RoleDao) CreateRole(role *models.Role, slugs []string) ([]string, error) {
	var permission []string
	if role == nil {
		return permission, ErrParam
	}
	id, err := d.Orm.Insert(role)
	if err != nil {
		return permission, err
	}
	role.Id = uint64(id)

	err = RolePermissionDaoEntity.InsertMulRolePermissions(role.Id, slugs)
	if err != nil {
		return permission, err
	}

	permission, err = RolePermissionDaoEntity.GetRolePermissions(role.Id)
	if err != nil {
		return permission, err
	}

	return permission, nil
}

//更新角色 pids:角色权限
func (d *RoleDao) UpdateRole(newRole *models.Role, slugs []string, cols ...string) ([]string, error) {
	var permission []string
	if newRole == nil {
		return permission, ErrParam
	}

	role := &models.Role{Id: newRole.Id}
	err := d.Orm.ReadForUpdate(role, models.COLUMN_Role_Id)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return permission, err
	}

	role = newRole
	_, err = d.Orm.Update(role, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return permission, err
	}

	if len(slugs) > 0 {
		p, err := RolePermissionDaoEntity.UpdateRolePermissions(role.Id, slugs)
		if err != nil {
			return permission, err
		}
		permission = p
	}

	return permission, nil
}

//删除角色
func (d *RoleDao) DelRole(roleID uint64) error {
	role := models.Role{Id: roleID}
	_, err := d.Orm.Delete(&role)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	//角色权限表
	rolePermission := models.RolePermission{RoleId: roleID}
	_, err = d.Orm.Delete(&rolePermission, models.COLUMN_RolePermission_RoleId)
	if err != nil {
		common.LogFuncDebug("mysql error:%v", err)
	}

	//用户角色表
	roleAdmin := models.RoleAdmin{RoleId: roleID}
	_, err = d.Orm.Delete(&roleAdmin, models.COLUMN_RoleAdmin_RoleId)
	if err != nil {
		common.LogFuncDebug("mysql error:%v", err)
	}

	return nil
}

//查询角色信息
func (d *RoleDao) QueryRole(role *models.Role, cols ...string) (RoleInfo, error) {
	var roleInfo RoleInfo
	if role == nil {
		return roleInfo, ErrParam
	}

	//角色基本信息
	err := d.Orm.Read(role, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return roleInfo, err
	}

	//角色权限
	per, err := RolePermissionDaoEntity.GetRolePermissions(role.Id)
	if err != nil {
		return roleInfo, err
	}

	roleInfo = GetBaseRoleInfo(*role)
	roleInfo.Permissions = per

	return roleInfo, nil
}

//分页查询角色
func (d *RoleDao) QueryPageRole(page int, perPage int) ([]RoleInfo, PageInfo, error) {
	var roleInfo []RoleInfo
	var info PageInfo

	qs := d.Orm.QueryTable(models.TABLE_Role)
	if qs == nil {
		common.LogFuncError("mysql_err:Permission fail")
		return roleInfo, info, ErrParam
	}

	count, err := qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			return roleInfo, info, nil
		}
		common.LogFuncError("mysql error:%v", err)
		return roleInfo, info, err
	}
	info.Total = int(count)
	info.Page = page
	info.Limit = perPage
	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > count {
		return roleInfo, info, ErrParam
	}

	var data []models.Role
	_, err = qs.Limit(perPage, start).All(&data)
	if err != nil {
		if err == orm.ErrNoRows {
			return roleInfo, info, nil
		}
		common.LogFuncError("mysql error:%v", err)
		return roleInfo, info, err
	}

	//查询所有角色拥有的权限
	var ids []string
	for _, v := range data {
		ids = append(ids, strconv.Itoa(int(v.Id)))
	}
	p, err := RolePermissionDaoEntity.GetUserRolesInfo(ids)
	if err != nil {
		return roleInfo, info, err
	}

	for _, r := range data {
		roleInfo = append(roleInfo, getRoleInfo(r, p))
	}

	return roleInfo, info, nil
}
