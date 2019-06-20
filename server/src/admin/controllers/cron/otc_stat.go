package cron

import (
	"common"
	"strconv"
	dao2 "utils/admin/dao"
	"utils/admin/models"
	dao3 "utils/game/dao"
	"utils/otc/dao"
	usdtdao "utils/usdt/dao"
)

func OtcStatCron() {
	data := OtcStat(-1)
	dao2.OtcStatDaoEntity.Create(data)
}

// OTC服务数据统计  day=0 今天  day=-1 昨天
func OtcStat(day int64) (data *models.OtcStat) {
	dayStart, dayEnd, dayDate, _ := common.TheOtherDayTimeRange(day)

	dayDate = dayDate[:8] //Ymd

	data = &models.OtcStat{}
	tmp, _ := strconv.Atoi(dayDate)
	data.Date = uint32(tmp)

	// 用户登录数
	data.NumLogin = dao.UserLoginLogDaoEntity.CountByTime(dayStart*1000, dayEnd*1000)
	// 新用户数
	data.NumUserNew = dao.UserDaoEntity.CountByTime(dayStart*1000, dayEnd*1000)
	// 订单数 & 金额 & 手续费
	resMap, err := dao.OrdersDaoEntity.OtcStat(dayDate)
	if err != nil {
		common.LogFuncError("Cron OtcStat ERR:%v", err)
		return
	}

	data.NumOrder = resMap["num_order"]
	data.NumOrderDeal = resMap["num_order_sell"] + resMap["num_order_buy"]
	data.NumOrderBuy = resMap["num_order_buy"]
	data.NumOrderSell = resMap["num_order_sell"]
	data.NumFunds = resMap["num_funds_sell"] + resMap["num_funds_buy"]
	data.NumFeeBuy = resMap["num_fee_buy"]
	data.NumFeeSell = resMap["num_fee_sell"]
	data.NumAmount = resMap["num_amount_sell"] + resMap["num_amount_buy"]
	data.NumAmountBuy = resMap["num_amount_buy"]
	data.NumAmountSell = resMap["num_amount_sell"]

	// 游戏平台 充值&提现
	data.GameRecharge, data.GameWithdrawal = dao3.GameTransferDaoEntity.CountByTime(dayStart*1000, dayEnd*1000)

	// USDT 数据统计
	resMap2 := usdtdao.WealthLogDaoEntity.Stat(dayStart*1000, dayEnd*1000)
	data.UsdtRecharge = resMap2["usdt_recharge"]
	data.UsdtWithdrawal = resMap2["usdt_withdrawal"]
	data.UsdtFee = resMap2["usdt_fee"]

	dao2.OtcStatDaoEntity.Create(data)

	return data
}

func OtcStatPeople(uid uint64) (res map[string]uint32) {
	res = map[string]uint32{
		"usdt_recharge":   0,
		"usdt_withdrawal": 0,
	}
	//usdt
	tmp1 := usdtdao.WealthLogDaoEntity.StatPeople(uid)
	res["usdt_recharge"] = tmp1["usdt_recharge"]
	res["usdt_withdrawal"] = tmp1["usdt_withdrawal"]

	//eos
	tmp2 := dao.OrdersDaoEntity.OtcStatPeople(uid)
	res["num_order_buy"] = tmp2["num_order_buy"]
	res["num_order_sell"] = tmp2["num_order_sell"]
	res["num_amount_buy"] = tmp2["num_amount_buy"]
	res["num_amount_sell"] = tmp2["num_amount_sell"]

	return
}
