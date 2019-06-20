package trade

import (
	common2 "common"
	"eusd/eosplus"
	"fmt"
	"math"
	controllers "otc_error"
	"strconv"
	"time"
	"usdt/prices"
	"utils/common"
	otcdao "utils/otc/dao"
	reportdao "utils/report/dao"
)

const TradeMatchLock = 30 * time.Second //锁定时间30s

func USDT2RMBRate() float64 {
	price := prices.GetPriceFloat64(prices.PRICE_CURRENCY_TYPE_USDT)
	return price
}

func BuyUSDT2RMBRate() float64 {
	feeStr, err := common.AppConfigMgr.String(common.BuyFeeRate)
	if err != nil {
		return 0
	}
	fee, err := strconv.ParseFloat(feeStr, 64)
	if err != nil {
		return 0
	}
	fee += 1

	return fee * USDT2RMBRate()
}

func SellUSDT2RMBRate() float64 {
	feeStr, err := common.AppConfigMgr.String(common.SellFeeRate)
	if err != nil {
		return 0
	}
	fee, err := strconv.ParseFloat(feeStr, 64)
	if err != nil {
		return 0
	}
	fee = 1 - fee

	return fee * USDT2RMBRate()
}

// rmb单位分 返回 EUSD * 精度
func SellEUSD2RMBWithPrecision(eusd int64) int64 {
	return int64(float64(eusd)*SellUSDT2RMBRate()/math.Pow10(eosplus.EosPrecision)*100 + 0.5)
}

// rmb单位分 返回 EUSD * 精度
func BuyEUSD2RMBWithPrecision(eusd int64) int64 {
	return int64(float64(eusd)*BuyUSDT2RMBRate()/math.Pow10(eosplus.EosPrecision)*100 + 0.5)
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func PayExpire() time.Duration {
	t, err := common.AppConfigMgr.Int(common.OtcTradePayExpire)
	if err != nil {
		//订单超时时间15分钟
		return 900 * time.Second
	}
	return time.Duration(t) * time.Second
}

func ConfirmExpire() time.Duration {
	t, err := common.AppConfigMgr.Int(common.OtcTradeConfirmExpire)
	if err != nil {
		//订单超时时间15分钟
		return 900 * time.Second
	}
	return time.Duration(t) * time.Second
}

func RedisLock(uid uint64) bool {
	key := fmt.Sprintf("trade_lock_%d", uid)
	resp := common2.RedisManger.SetNX(key, 1, 3*time.Second)
	lockSuccess, err := resp.Result()
	if err != nil || !lockSuccess {
		return false
	}
	return true
}

func GetTradeUpperLimit() int64 {
	num, err := common.AppConfigMgr.Int64(common.OtcTradeUpperLimitRmb)
	if err != nil {
		common2.LogFuncError("GetTradeUpperLimit err:%v", err)
	}
	return num
}

func GetTradeLowerLimit() int64 {
	num, err := common.AppConfigMgr.Int64(common.OtcTradeLowerLimitRmb)
	if err != nil {
		common2.LogFuncError("GetTradeLowerLimit err:%v", err)
	}
	return num
}

func CheckAppealOrder(uid, orderId uint64) controllers.ERROR_CODE {
	appeal, err := otcdao.AppealDaoEntity.QueryByOrderId(orderId)
	if err != nil {
		return controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
	}

	if appeal.UserId > 0 && appeal.UserId != uid {
		return controllers.ERROR_CODE_OTHER_SIDE_APPEAL
	}

	return controllers.ERROR_CODE_SUCCESS
}

func DayAmountReport(side int8, uid, eUid uint64, amount int64) {
	if side == otcdao.SideBuy {
		_ = reportdao.ReportEusdDailyDaoEntity.Buy(uid, amount)
		_ = reportdao.ReportEusdDailyDaoEntity.Sell(eUid, amount)
	} else {
		_ = reportdao.ReportEusdDailyDaoEntity.Buy(eUid, amount)
		_ = reportdao.ReportEusdDailyDaoEntity.Sell(uid, amount)
	}
}
