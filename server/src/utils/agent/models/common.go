package models

import "github.com/astaxie/beego/orm"

func ModelsInit() (err error) {
	orm.RegisterModel(
		new(Agent),
		new(AgentChannelCommission),
		new(AgentPath),
		new(AgentWithdraw),
		new(InviteCode),
	)

	return
}
