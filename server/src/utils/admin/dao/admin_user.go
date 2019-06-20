package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"strconv"
	"utils/admin/models"
)

//用户状态
const (
	AdminUserStatusNil = iota
	AdminUserStatusActive
	AdminUserStatusSuspended
	AdminUserStatusDeleted
	AdminUserStatusMax
)

const AdminPwdSalt = "ZyGYFWIWO1BWYl9lpBKaNtKmXxFRrHwu5PgJD9V332AEWweZY1QdrRyTbjcAdmin"
const defaultPassword = "9d81dbda89066a79d1ac2c1de8f554b8" //"123456"

type AdminUserDao struct {
	common.BaseDao
}

func NewAdminUserDao(db string) *AdminUserDao {
	return &AdminUserDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AdminUserDaoEntity *AdminUserDao

//查询管理员
func (d *AdminUserDao) QueryAdminUser(adminUser *models.AdminUser, cols ...string) error {
	if adminUser == nil {
		return ErrParam
	}

	err := d.Orm.Read(adminUser, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	return nil
}

//查询管理员
func (d *AdminUserDao) QueryAdminUserOrCreate(adminUser *models.AdminUser, col string, cols ...string) error {
	if adminUser == nil {
		return ErrParam
	}

	inPwd := adminUser.Pwd
	adminUser.Ctime = common.NowInt64MS()
	adminUser.Utime = common.NowInt64MS()
	isNew, id, err := d.Orm.ReadOrCreate(adminUser, col, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	adminUser.Id = uint64(id)

	if isNew {
		if inPwd != "" {
			adminUser.Pwd, err = common.GenerateDoubleMD5(inPwd, AdminPwdSalt)
			if err != nil {
				common.LogFuncError("GenerateDoubleMD5 error:%v", err)
				return err
			}
			adminUser.Status = AdminUserStatusActive
			_, err := d.Orm.Update(adminUser, models.COLUMN_AdminUser_Pwd, models.COLUMN_AdminUser_Status)
			if err != nil {
				common.LogFuncError("mysql_err:%v", err)
				return err
			}
		}
	}

	return nil
}

//分页查询管理员角色
func (d *AdminUserDao) QueryPageAdminRole(id uint64, page int, perPage int) ([]RoleInfo, PageInfo, error) {
	var roleInfo []RoleInfo
	var info PageInfo

	ids, err := RoleAdminDaoEntity.QueryRoleIdByAdminId(id)
	if err != nil {
		return roleInfo, info, err
	}
	qs := d.Orm.QueryTable(models.TABLE_Role).Filter("id__in", ids)
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
	var roles []models.Role
	start := (page - 1) * perPage
	if start > info.Total {
		start = 0
		info.Page = 1
	}
	_, err = qs.Limit(perPage, start).All(&roles)
	if err != nil {
		if err == orm.ErrNoRows {
			return roleInfo, info, nil
		}
		common.LogFuncError("mysql error:%v", err)
		return roleInfo, info, err
	}

	var strIds []string
	for _, v := range ids {
		strIds = append(strIds, strconv.Itoa(int(v)))
	}
	p, err := RolePermissionDaoEntity.GetUserRolesInfo(strIds)
	if err != nil {
		return roleInfo, info, err
	}
	for _, role := range roles {
		roleInfo = append(roleInfo, getRoleInfo(role, p))
	}

	return roleInfo, info, nil
}

//分配角色给管理员
func (d *AdminUserDao) AddAdminUserRoles(adminUser *models.AdminUser, roleID uint64, grantedBy string, cols ...string) error {
	if adminUser == nil {
		return ErrParam
	}

	err := d.QueryAdminUser(adminUser, cols...)
	if err != nil {
		return err
	}

	ids, err := RoleAdminDaoEntity.QueryRoleIdByAdminId(adminUser.Id)
	if err != nil {
		return err
	}
	for _, v := range ids {
		if v == roleID {
			return ErrParam
		}
	}

	//插入角色管理员信息表
	roleAdmin := new(models.RoleAdmin)
	roleAdmin.RoleId = roleID
	roleAdmin.AdminId = adminUser.Id
	roleAdmin.GrantedBy = grantedBy
	roleAdmin.GrantedAt = common.NowInt64MS()
	_, err = d.Orm.Insert(roleAdmin)
	if err != nil {
		return err
	}

	return nil
}

//分页获取角色成员
func (d *AdminUserDao) QueryPageRoleAdminUser(roleID uint64, page int, perPage int) ([]RoleMember, PageInfo, error) {
	var members []RoleMember
	var info PageInfo
	qs := d.Orm.QueryTable(models.TABLE_RoleAdmin).Filter("roleid__in", roleID)
	count, err := qs.Count()
	if err != nil {
		common.LogFuncError("mysql_err fail%v\n", err)
		return members, info, err
	}

	info.Total = int(count)
	info.Page = page
	info.Limit = perPage
	//获取当前页数据的起始位置
	var roleAdmin []models.RoleAdmin
	start := (page - 1) * perPage
	if start > info.Total {
		return members, info, ErrParam
	}
	_, err = qs.Limit(perPage, start).All(&roleAdmin)
	if err != nil {
		if err == orm.ErrNoRows {
			return members, info, nil
		}
		common.LogFuncError("mysql error:%v", err)
		return members, info, err
	}

	var ids []uint64
	for _, r := range roleAdmin {
		ids = append(ids, r.AdminId)
	}
	if len(ids) == 0 {
		return members, info, nil
	}

	var adminUser []models.AdminUser
	_, err = d.Orm.QueryTable(models.TABLE_AdminUser).Filter("id__in", ids).All(&adminUser)
	if nil != err {
		return members, info, err
	}
	for _, ra := range roleAdmin {
		for _, a := range adminUser {
			if ra.AdminId == a.Id {
				member := RoleMember{
					Id:        int(a.Id),
					Name:      a.Name,
					Email:     a.Email,
					Status:    a.Status,
					GrantedBy: ra.GrantedBy,
					GrantedAt: ra.GrantedAt,
				}
				members = append(members, member)
			}
		}
	}

	return members, info, nil
}

//分页获取管理员信息
func (d *AdminUserDao) QueryPageAdminUser(roleID uint64, page int, perPage int, status int8) ([]AdminMember, PageInfo, error) {
	var err error
	var ids []uint64
	var roleAdmin []models.RoleAdmin
	var members []AdminMember
	var info PageInfo
	if roleID != 0 {
		//获取拥有该角色的所有admin_id
		_, err = d.Orm.QueryTable(models.TABLE_RoleAdmin).Filter(models.COLUMN_RoleAdmin_RoleId, roleID).All(&roleAdmin)
		if err != nil {
			common.LogFuncError("mysql_err:%v", err)
			return members, info, err
		}
		for _, r := range roleAdmin {
			ids = append(ids, r.AdminId)
		}
	}

	//分页查询admin_user
	qs := d.Orm.QueryTable(models.TABLE_AdminUser)
	if len(ids) > 0 {
		qs = qs.Filter("id__in", ids)
	}
	if status > AdminUserStatusNil && status < AdminUserStatusMax {
		qs = qs.Filter(models.COLUMN_AdminUser_Status, status)
	}
	count, err := qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			return members, info, nil
		}
		common.LogFuncError("mysql error:%v", err)
		return members, info, err
	}
	info.Total = int(count)
	info.Page = page
	info.Limit = perPage
	//获取当前页数据的起始位置
	var adminUser []models.AdminUser
	start := (page - 1) * perPage
	if start > info.Total {
		return members, info, ErrParam
	}
	_, err = qs.Limit(perPage, start).All(&adminUser)
	if err != nil {
		if err == orm.ErrNoRows {
			return members, info, nil
		}
		common.LogFuncError("mysql error:%v", err)
		return members, info, err
	}

	//查询管理员角色
	var strIds []string
	for _, a := range adminUser {
		strIds = append(strIds, strconv.Itoa(int(a.Id)))
	}
	roleInfo, err := RoleAdminDaoEntity.QueryRoleAdminByAdmins(strIds)
	if err != nil {
		return members, info, err
	}
	for _, a := range adminUser {
		member := AdminMember{
			Id:           a.Id,
			Name:         a.Name,
			Email:        a.Email,
			Status:       a.Status,
			WhitelistIps: a.WhitelistIps,
			CTime:        a.Ctime,
			UTime:        a.Utime,
			DTime:        a.Dtime,
			TimeLogin:    a.LoginTime,
			IsBind:       a.IsBind,
		}
		for _, r := range roleInfo {
			if r.Adminid == member.Id {
				b := RoleBase{
					Id:   r.Roleid,
					Name: r.Name,
				}
				member.RBase = append(member.RBase, b)
			}
		}
		members = append(members, member)
	}

	return members, info, nil
}

//创建管理员
func (d *AdminUserDao) CreateAdminUser(newAdmin *models.AdminUser, roleIds []uint64, by string) (AdminMember, error) {
	var member AdminMember
	if newAdmin == nil {
		return member, ErrParam
	}

	qAdmin := &models.AdminUser{Email: newAdmin.Email}
	err := d.QueryAdminUser(qAdmin, models.COLUMN_AdminUser_Email)
	if err != nil {
		common.LogFuncDebug("mysql error:%v", err)
	}
	if qAdmin.Id != 0 {
		return member, ErrParam
	}

	//插入数据
	adminUser := new(models.AdminUser)
	adminUser.Name = newAdmin.Name
	adminUser.Email = newAdmin.Email
	adminUser.Status = AdminUserStatusActive
	adminUser.WhitelistIps = newAdmin.WhitelistIps
	adminUser.Ctime = common.NowInt64MS()
	adminUser.Utime = common.NowInt64MS()
	adminUser.Pwd = defaultPassword
	id, err := d.Orm.Insert(adminUser)
	if err != nil {
		return member, err
	}
	adminUser.Id = uint64(id)
	err = RoleAdminDaoEntity.InsertRolesAdmin(uint64(id), by, roleIds)
	if err != nil {
		return member, err
	}

	member, err = GetAdminMember(*adminUser)
	if err != nil {
		return member, err
	}

	return member, nil
}

//更新管理员
func (d *AdminUserDao) UpdateAdminUser(newAdmin *models.AdminUser, roleIds []uint64, by string) (AdminMember, error) {
	var member AdminMember
	adminUser := &models.AdminUser{Id: newAdmin.Id}
	err := d.QueryAdminUser(adminUser, models.COLUMN_AdminUser_Id)
	if err != nil {
		if err == orm.ErrNoRows {
			return member, nil
		}
		common.LogFuncError("mysql error:%v", err)
		return member, err
	}

	var cols []string
	if newAdmin.Name != "" {
		adminUser.Name = newAdmin.Name
		cols = append(cols, models.COLUMN_AdminUser_Name)
	}
	if newAdmin.Email != "" && newAdmin.Email != adminUser.Email {
		qAdmin := &models.AdminUser{
			Email: newAdmin.Email,
		}
		err = d.QueryAdminUser(qAdmin, models.COLUMN_AdminUser_Email)
		if err == nil {
			//用户已被注册
			return member, ErrParam
		}

		adminUser.Email = newAdmin.Email
		cols = append(cols, models.COLUMN_AdminUser_Email)
	}
	if newAdmin.WhitelistIps != "" {
		adminUser.WhitelistIps = newAdmin.WhitelistIps
		cols = append(cols, models.COLUMN_AdminUser_WhitelistIps)
	}
	if newAdmin.Status > AdminUserStatusNil && newAdmin.Status < AdminUserStatusMax {
		adminUser.Status = newAdmin.Status
		cols = append(cols, models.COLUMN_AdminUser_Status)
		if newAdmin.Status == AdminUserStatusDeleted {
			adminUser.Dtime = common.NowInt64MS()
			cols = append(cols, models.COLUMN_AdminUser_Dtime)
		}
	}
	if len(roleIds) > 0 {
		//删除旧数据
		_, err := d.Orm.QueryTable(models.TABLE_RoleAdmin).Filter(models.COLUMN_RoleAdmin_AdminId, adminUser.Id).Delete()
		if err != nil {
			common.LogFuncError("mysql error:%v", err)
		}
		err = RoleAdminDaoEntity.InsertRolesAdmin(adminUser.Id, by, roleIds)
		if err != nil {
			return member, err
		}
	}
	if len(cols) == 0 {
		return member, nil
	}
	adminUser.Utime = common.NowInt64MS()
	cols = append(cols, models.COLUMN_AdminUser_Utime)

	_, err = d.Orm.Update(adminUser, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return member, err
	}

	q, err := GetAdminInfo(newAdmin.Id)
	if err != nil {
		return member, err
	}
	member = TransAdminMember(q)

	return member, nil
}

//更新管理员状态
func (d *AdminUserDao) UpdateAdminUserStatus(id uint64, newAdmin *models.AdminUser, cols ...string) error {
	adminUser := &models.AdminUser{Id: id}
	err := d.QueryAdminUser(adminUser, models.COLUMN_AdminUser_Id)
	if err != nil {
		return err
	}

	adminUser = newAdmin
	_, err = d.Orm.Update(adminUser, cols...)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

//更新管理员状态
func (d *AdminUserDao) UpdateTableAdminUser(adminUser *models.AdminUser, cols ...string) error {
	if adminUser == nil {
		return ErrParam
	}
	_, err := d.Orm.Update(adminUser, cols...)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

//修改密码
func (d *AdminUserDao) ChangePassword(adminUser *models.AdminUser, oldPwd string, newPwd string) error {
	if newPwd == "" {
		return ErrParam
	}

	md5NewPwd, err := common.GenerateDoubleMD5(newPwd, AdminPwdSalt)
	if err != nil {
		common.LogFuncError("GenerateDoubleMD5 error:%v", err)
		return err
	}
	if adminUser.Pwd != "" {
		md5OldPwd, err := common.GenerateDoubleMD5(oldPwd, AdminPwdSalt)
		if err != nil {
			common.LogFuncError("GenerateDoubleMD5 error:%v", err)
			return err
		}

		if md5OldPwd != adminUser.Pwd {
			return ErrParam
		}
	}

	adminUser.Pwd = md5NewPwd
	_, err = d.Orm.Update(adminUser, models.COLUMN_AdminUser_Pwd)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

func (d *AdminUserDao) UpdateGoogleAuth(id uint64, secretId, qrCode string) error {
	adminUser := &models.AdminUser{
		Id:       id,
		SecretId: secretId,
		QrCode:   qrCode,
	}
	_, err := d.Orm.Update(adminUser, models.COLUMN_AdminUser_SecretId, models.COLUMN_AdminUser_QrCode)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

func (d *AdminUserDao) UpdateGoogleAuthBind(id uint64, isBind bool) (err error) {
	adminUser := &models.AdminUser{
		Id:     id,
		IsBind: isBind,
		Utime:  common.NowInt64MS(),
	}
	if isBind {
		_, err = d.Orm.Update(adminUser, models.COLUMN_AdminUser_IsBind, models.COLUMN_AdminUser_Utime)
	} else {
		adminUser.SecretId = ""
		adminUser.QrCode = ""
		_, err = d.Orm.Update(adminUser, models.COLUMN_AdminUser_IsBind, models.COLUMN_AdminUser_Utime,
			models.COLUMN_AdminUser_QrCode, models.COLUMN_AdminUser_SecretId)
	}

	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}
