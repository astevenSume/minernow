package dao

import (
	"common"

	"github.com/astaxie/beego"
)

const SIGNATURE_SALT = "ZyGYFWIWO1BWYl9lpBKaNtKmXxFRrHwu5PgJD9V332AEWweZY1QdrRyTbjcAjmaq"

const (
	IdTypeUser                 = iota //user id
	IdTypeUsdtReq                     //usdt request id
	IdTypePaymentMethod               //payment method
	IdTypeMessageMethod               //message method
	IdTypeOrders                      //orders
	IdTypeAgentWithdraw               //agent withdraw record
	IdTypePayPassValidator            //validate payment password action
	IdTypeCommissionCalc              // commission calc
	IdTypeCommissionDistribute        // commission distribute
	IdTypeSystemNotification          //system notification
	IdTypeMax
)

type DaoConfig struct {
	PaymentMethod_LowMoneyPerTx  int64
	PaymentMethod_HighMoneyPerTx int64
	PaymentMethod_TimesPerDay    int64
	PaymentMethod_MoneyPerDay    int64
	PaymentMethod_MoneySum       int64
}

func NewDaoConfig() *DaoConfig {
	return &DaoConfig{}
}

var Config = NewDaoConfig()

func initConfig() (err error) {
	Config.PaymentMethod_LowMoneyPerTx, err = beego.AppConfig.Int64("payment::LowMoneyPerTx")
	if err != nil {
		return
	}
	Config.PaymentMethod_HighMoneyPerTx, err = beego.AppConfig.Int64("payment::HighMoneyPerTx")
	if err != nil {
		return
	}
	Config.PaymentMethod_TimesPerDay, err = beego.AppConfig.Int64("payment::TimesPerDay")
	if err != nil {
		return
	}
	Config.PaymentMethod_MoneyPerDay, err = beego.AppConfig.Int64("payment::MoneyPerDay")
	if err != nil {
		return
	}
	Config.PaymentMethod_MoneySum, err = beego.AppConfig.Int64("payment::MoneySum")
	if err != nil {
		return
	}

	return
}

func Init(entityInitFunc common.EntityInitFunc) (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	const dbOtc = "otc"
	MessageMethodDaoEntity = NewMessageMethodDao(dbOtc)
	OtcSellDaoEntity = NewOtcSellDao(dbOtc)
	PaymentMethodDaoEntity = NewPaymentMethodDao(dbOtc)
	OrdersDaoEntity = NewOrdersDao(dbOtc)
	UserDaoEntity = NewUserDao(dbOtc)
	OtcExchangerVerifyDaoEntity = NewOtcExchangerVerifyDao(dbOtc)
	OtcBuyDaoEntity = NewOtcBuyDao(dbOtc)
	UserLoginLogDaoEntity = NewUserLoginLogDao(dbOtc)
	AppealDaoEntity = NewAppealDao(dbOtc)
	UserPayPassDaoEntity = NewUserPayPassDao(dbOtc)
	CommissionStatDaoEntity = NewCommissionStatDao(dbOtc)
	OtcExchangerDaoEntity = NewOtcExchangerDao(dbOtc)
	UserConfigDaoEntity = NewUserConfigDao(dbOtc)
	CommissionCalcDaoEntity = NewCommissionCalcDao(dbOtc)
	SystemNotificationdDaoEntity = NewSystemNotificationdDao(dbOtc)
	CommissionDistributeDaoEntity = NewCommissionDistributeDao(dbOtc)
	AppealDealLogDaoEntity = NewAppealDealLogDao(dbOtc)
	EosOtcReportDaoEntity = NewEosOtcReportDao(dbOtc)

	if entityInitFunc != nil {
		err = entityInitFunc()
		if err != nil {
			return
		}
	}
	return
}
