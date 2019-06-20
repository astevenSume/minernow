package models

import "github.com/astaxie/beego/orm"

func ModelsInit() (err error) {
	orm.RegisterModel(
		new(Appeal),
		new(AppealDealLog),
		new(CommissionCalc),
		new(CommissionDistribute),
		new(CommissionStat),
		new(EosOtcReport),
		new(OtcBuy),
		new(OtcExchanger),
		new(OtcExchangerVerify),
		new(OtcMsg),
		new(OtcOrder),
		new(OtcSell),
		new(PaymentMethod),
		new(SystemNotification),
		new(User),
		new(UserConfig),
		new(UserLoginLog),
		new(UserPayPassword),
	)

	return
}
