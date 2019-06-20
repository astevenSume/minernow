package models

import "github.com/astaxie/beego/orm"

func ModelsInit() (err error) {
	orm.RegisterModel(
		new(EosAccount),
		new(EosOtc),
		new(EosTransaction),
		new(EosTransactionInfo),
		new(EosTxLog),
		new(EosUseLog),
		new(EosWealth),
		new(EosWealthLog),
		new(EusdRetire),
		new(PlatformUser),
		new(PlatformUserCate),
	)

	return
}
