package models

import "github.com/astaxie/beego/orm"

func ModelsInit() (err error) {
	orm.RegisterModel(
		new(ActivityUserConf),
		new(AdminUser),
		new(AgentWhiteList),
		new(Announcement),
		new(AppChannel),
		new(AppType),
		new(AppVersion),
		new(AppWhitelist),
		new(AppealService),
		new(Apps),
		new(Banner),
		new(Commissionrates),
		new(Config),
		new(ConfigWarning),
		new(Endpoint),
		new(IpWhiteList),
		new(MenuAccess),
		new(MenuConf),
		new(MonthDividendPositionConf),
		new(MonthDividendWhiteList),
		new(OperationLog),
		new(OtcStat),
		new(OtcStatAllPeople),
		new(Permission),
		new(ProfitThreshold),
		new(Role),
		new(RoleAdmin),
		new(RolePermission),
		new(ServerNode),
		new(Smscodes),
		new(Smstemplates),
		new(SystemMessage),
		new(SysNotification),
		new(Task),
		new(TaskResult),
		new(TopAgent),
	)

	return
}
