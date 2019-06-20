package models

import "github.com/astaxie/beego/orm"

func ModelsInit() (err error) {
	orm.RegisterModel(
		new(ChannelDaily),
		new(GameLog),
		new(GameOrderRisk),
		new(GameRiskAlert),
		new(GameTransfer),
		new(GameUser),
		new(GameUserDaily),
	)

	return
}
