package dao

import (
	"common"
	"fmt"
	"strings"
	"utils/admin/models"
)

type RoleAdminDao struct {
	common.BaseDao
}

func NewRoleAdminDao(db string) *RoleAdminDao {
	return &RoleAdminDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var RoleAdminDaoEntity *RoleAdminDao

//查询用户角色ID
func (d *RoleAdminDao) QueryRoleIdByAdminId(id uint64) (ids []uint64, err error) {
	mysql := fmt.Sprintf("select %s from %s where %s = ?", models.COLUMN_RoleAdmin_RoleId, models.TABLE_RoleAdmin, models.COLUMN_RoleAdmin_AdminId)
	_, err = d.Orm.Raw(mysql, id).QueryRows(&ids)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	return
}

//查询用户角色信息
func (d *RoleAdminDao) QueryRoleAdminByAdminId(id uint64) (roleInfo []RoleAllInfo, err error) {
	//mysql := fmt.Sprintf("select * from ((SELECT * from %s where %s = ?) as T1 LEFT JOIN role as T2 on T1.%s=T2.id)",models.TABLE_RoleAdmin, models.COLUMN_RoleAdmin_AdminId, models.COLUMN_RoleAdmin_RoleId,)
	_, err = d.Orm.Raw("select * from ((SELECT * from role_admin where adminid = ?) as T1 LEFT JOIN role as T2 on T1.roleid=T2.id)", id).QueryRows(&roleInfo)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	return
}

//批量查询查询用户角色信息
func (d *RoleAdminDao) QueryRoleAdminByAdmins(strIds []string) (roleInfo []RoleAllInfo, err error) {
	strId := strings.Join(strIds, CommaSpace)
	sql := fmt.Sprintf("select * from ((SELECT * from role_admin where adminid in(%s)) as T1 LEFT JOIN role as T2 on T1.roleid=T2.id)", strId)
	_, err = d.Orm.Raw(sql).QueryRows(&roleInfo)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//插入用户角色信息
func (d *RoleAdminDao) InsertRolesAdmin(id uint64, by string, ids []uint64) error {
	var roleAdmin []models.RoleAdmin
	for _, rid := range ids {
		r := models.RoleAdmin{
			RoleId:    rid,
			AdminId:   id,
			GrantedBy: by,
			GrantedAt: common.NowInt64MS(),
		}
		roleAdmin = append(roleAdmin, r)
	}
	_, err := d.Orm.InsertMulti(InsertMulCount, roleAdmin)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return nil
}

//删除角色管理员信息表
func (d *RoleAdminDao) DelRolesAdmin(roleId uint64, adminId uint64) error {
	roleAdmin := new(models.RoleAdmin)
	roleAdmin.RoleId = roleId
	roleAdmin.AdminId = adminId
	_, err := d.Orm.Delete(roleAdmin, models.COLUMN_RoleAdmin_RoleId, models.COLUMN_RoleAdmin_AdminId)
	if err != nil {
		return err
	}

	return nil
}
