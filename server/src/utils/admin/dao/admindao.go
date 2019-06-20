package dao

import (
	"common"
	"strconv"
	"utils/admin/models"
)

const CommaSpace = ", "

type RoleAllInfo struct {
	Adminid     uint64
	Roleid      uint64
	Name        string
	Desc        string
	Ctime       int64
	Utime       int64
	GrantedBy   string
	GrantedAt   int64
	Permissions []string
}
type RolePermission struct {
	Roleid uint64
	Slug   string
}

type MapRole map[uint64]RoleAllInfo
type AdminDao struct {
	AdminUser      models.AdminUser //管理员信息
	MapRoleInfo    MapRole
	SliPermissions []string
}

func Init(entityInitFunc common.EntityInitFunc) (err error) {
	//init entity (todo : maybe this should be moved to business svr project?)
	const dbAdmin = "otc_admin"
	PermissonDaoEntity = NewPermissonDao(dbAdmin)
	AdminUserDaoEntity = NewAdminUserDao(dbAdmin)
	CommissionRateDaoEntity = NewCommissionRateDao(dbAdmin)
	OperationLogDaoEntity = NewOperationLogDao(dbAdmin)
	RoleDaoEntity = NewRoleDao(dbAdmin)
	RoleAdminDaoEntity = NewRoleAdminDao(dbAdmin)
	RolePermissionDaoEntity = NewRolePermissionDao(dbAdmin)
	SmsCodeDaoEntity = NewSmsCodeDao(dbAdmin)
	SmsTemplateDaoEntity = NewSmsTemplateDao(dbAdmin)
	SystemMessageMethodDaoEntity = NewSystemMessageMethodDao(dbAdmin)
	AppVersionDaoEntity = NewAppVersionDao(dbAdmin)
	AppDaoEntity = NewAppDao(dbAdmin)
	AgentWhiteListDaoEntity = NewAgentWhiteListDao(dbAdmin)
	BannerDaoEntity = NewBannerDao(dbAdmin)
	ConfigDaoEntity = NewConfigDao(dbAdmin)
	AppTypeDaoEntity = NewAppTypeDao(dbAdmin)
	AppChannelDaoEntity = NewAppChannelDao(dbAdmin)
	IpWhiteListDaoEntity = NewIpWhiteListDao(dbAdmin)
	SysNotificationDaoEntity = NewSysNotificationDao(dbAdmin)
	AppealServiceDaoEntity = NewAppealServiceDao(dbAdmin)
	TopAgentDaoEntity = NewTopAgentDao(dbAdmin)
	EndPointDaoEntity = NewEndPointDao(dbAdmin)
	AnnouncementDaoEntity = NewAnnouncementDao(dbAdmin)
	ConfigWarningDaoEntity = NewConfigWarningDao(dbAdmin)
	AppWhiteListDaoEntity = NewAppWhiteListDao(dbAdmin)
	OtcStatDaoEntity = NewOtcStatDao(dbAdmin)
	ServerNodeDaoEntity = NewServerNodeDao(dbAdmin)
	TaskDaoEntity = NewTaskDao(dbAdmin)
	TaskResultDaoEntity = NewTaskResultDao(dbAdmin)
	MonthDividendPositionConfDaoEntity = NewMonthDividendPositionConfDao(dbAdmin)
	ActivityUserConfDaoEntity = NewActivityUserConfDao(dbAdmin)
	MenuConfDaoEntity = NewMenuConfDao(dbAdmin)
	MenuAccessDaoEntity = NewMenuAccessDao(dbAdmin)
	MonthDividendWhiteListDaoEntity = NewMonthDividendWhiteListDao(dbAdmin)
	ProfitThresholdDaoEntity = NewProfitThresholdDao(dbAdmin)
	//try reinit some entity
	if entityInitFunc != nil {
		err = entityInitFunc()
		if err != nil {
			return
		}
	}

	PermissonDaoEntity.LoadAllPermission()

	return
}

func getRoleInfo(role models.Role, rolePermission []RolePermission) RoleInfo {
	roleInfo := RoleInfo{
		Id:    role.Id,
		Name:  role.Name,
		Desc:  role.Desc,
		Ctime: role.Ctime,
		Utime: role.Utime,
	}
	for _, v := range rolePermission {
		if v.Roleid == role.Id {
			roleInfo.Permissions = append(roleInfo.Permissions, v.Slug)
		}
	}

	return roleInfo
}

func GetBaseRoleInfo(role models.Role) RoleInfo {
	roleInfo := RoleInfo{
		Id:    role.Id,
		Name:  role.Name,
		Desc:  role.Desc,
		Ctime: role.Ctime,
		Utime: role.Utime,
	}

	return roleInfo
}

func TransAdminMember(a *AdminDao) AdminMember {
	return AdminMember{
		Id:     a.AdminUser.Id,
		Name:   a.AdminUser.Name,
		Email:  a.AdminUser.Email,
		Status: a.AdminUser.Status,
		RBase:  getRoleBaseInfo(a.MapRoleInfo),
		//Permissions: a.SliPermissions,
		WhitelistIps: a.AdminUser.WhitelistIps,
		CTime:        a.AdminUser.Ctime,
		UTime:        a.AdminUser.Utime,
		DTime:        a.AdminUser.Dtime,
		TimeLogin:    a.AdminUser.LoginTime,
		IsBind:       a.AdminUser.IsBind,
	}
}

func getRoleBaseInfo(mapRole MapRole) (rolebases []RoleBase) {
	for _, r := range mapRole {
		base := RoleBase{
			Id:   r.Roleid,
			Name: r.Name,
		}
		rolebases = append(rolebases, base)
	}

	return
}

func GetAdminMember(adminUser models.AdminUser) (AdminMember, error) {
	var member AdminMember
	member.Id = adminUser.Id
	member.Name = adminUser.Name
	member.Email = adminUser.Email
	member.Status = adminUser.Status
	member.WhitelistIps = adminUser.WhitelistIps
	member.DTime = adminUser.Dtime
	member.UTime = adminUser.Utime
	member.CTime = adminUser.Ctime
	member.TimeLogin = adminUser.LoginTime

	//查询管理员角色
	roleInfo, err := RoleAdminDaoEntity.QueryRoleAdminByAdminId(adminUser.Id)
	if err != nil {
		return member, err
	}
	for _, r := range roleInfo {
		b := RoleBase{
			Id:   r.Roleid,
			Name: r.Name,
		}
		member.RBase = append(member.RBase, b)
	}

	return member, nil
}

func GetLoginInfo(adminUser models.AdminUser) (LoginInfo, error) {
	var member LoginInfo
	member.Id = adminUser.Id
	member.Name = adminUser.Name
	member.Email = adminUser.Email
	member.Status = adminUser.Status
	member.WhitelistIps = adminUser.WhitelistIps
	member.UTime = adminUser.Utime
	member.CTime = adminUser.Ctime
	member.TimeLogin = adminUser.LoginTime
	member.IsBind = adminUser.IsBind
	if !member.IsBind {
		member.SecretId = adminUser.SecretId
		member.QrCode = adminUser.QrCode
	} else {
		//查询管理员角色
		roleInfo, err := RoleAdminDaoEntity.QueryRoleAdminByAdminId(adminUser.Id)
		if err != nil {
			return member, err
		}
		var ids []string
		for _, r := range roleInfo {
			b := RoleBase{
				Id:   r.Roleid,
				Name: r.Name,
			}
			member.RBase = append(member.RBase, b)
			ids = append(ids, strconv.Itoa(int(r.Roleid)))
		}

		if len(ids) == 0 {
			return member, nil
		}

		p, err := RolePermissionDaoEntity.GetUserRolesInfo(ids)
		if err != nil {
			return member, err
		}
		for _, r := range p {
			member.Permissions = append(member.Permissions, r.Slug)
		}
	}

	return member, nil
}

//获取管理员信息
func GetAdminInfo(id uint64) (*AdminDao, error) {
	//查询管理员基本信息
	adminUser := &models.AdminUser{
		Id: id,
	}

	err := AdminUserDaoEntity.QueryAdminUser(adminUser)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return nil, err
	}
	d := &AdminDao{}
	d.AdminUser = *adminUser

	//查询管理员角色
	roleInfo, err := RoleAdminDaoEntity.QueryRoleAdminByAdminId(adminUser.Id)
	if err != nil {
		return nil, err
	}

	if len(roleInfo) > 0 {
		var ids []string
		d.MapRoleInfo = make(map[uint64]RoleAllInfo)
		for _, r := range roleInfo {
			d.MapRoleInfo[r.Roleid] = r
			ids = append(ids, strconv.Itoa(int(r.Roleid)))
		}

		//角色权限
		rolePermission, err := RolePermissionDaoEntity.GetUserRolesInfo(ids)
		if err != nil {
			return nil, err
		}
		for _, r := range rolePermission {
			if item, ok := d.MapRoleInfo[r.Roleid]; ok {
				item.Permissions = append(item.Permissions, r.Slug)
			}
			d.SliPermissions = append(d.SliPermissions, r.Slug)
		}
	}

	return d, nil
}

//是否有权限
func (d *AdminDao) Has(right string) bool {
	for _, p := range d.SliPermissions {
		if p == right {
			return true
		}
	}
	return false
	//return utils.DetermineIndex(d.SliPermissions, right) >= 0
}

//是否有权限之一
func (d *AdminDao) HasOneOf(rights ...string) bool {
	for _, right := range rights {
		for _, p := range d.SliPermissions {
			if p == right {
				return true
			}
		}
		//if utils.DetermineIndex(d.SliPermissions, right) >= 0 {
		//	return true
		//}
	}

	return false
}
