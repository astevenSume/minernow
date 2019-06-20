package trade

import (
	common2 "common"
	"eusd/eosplus"
	"fmt"
	"math/rand"
	otc_common "otc/common"
	"otc_error"
	"umeng_push/uemng_plus"
	dao2 "utils/eusd/dao"
	models2 "utils/eusd/models"
	"utils/otc/dao"
	"utils/otc/models"
)

type BuyLogic struct {
}

var BuyLogicEntity = &BuyLogic{}

//承兑商otc信息
func (b *BuyLogic) Info(uid uint64) (wealth *models2.EosOtc, errCode controllers.ERROR_CODE) {
	errCode = controllers.ERROR_CODE_SUCCESS
	wealth, errCode2 := eosplus.EosPlusAPI.Otc.Info(uid)
	errCode = controllers.ERROR_CODE(errCode2)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		errCode = controllers.ERROR_CODE_OTC_NOT_A_EXCHANGER
		return
	}
	if wealth.Account == "" {
		errCode = controllers.ERROR_CODE_OTC_NOT_A_EXCHANGER
		return
	}
	return
}

// 设置
func (b *BuyLogic) Setting(wealth *models2.EosOtc, able bool, dayLimit, lowerLimit int64) (res *models2.EosOtc, errCode controllers.ERROR_CODE) {
	errCode = controllers.ERROR_CODE_SUCCESS
	TradeLowerLimitRMB := GetTradeLowerLimit()
	TradeUpperLimitRMB := GetTradeUpperLimit()
	if dayLimit > 0 && dayLimit < TradeLowerLimitRMB {
		errCode = controllers.ERROR_CODE_OTC_TRADE_TOO_SMALL
		return
	}

	if lowerLimit > 0 && lowerLimit < TradeLowerLimitRMB {
		errCode = controllers.ERROR_CODE_OTC_TRADE_TOO_SMALL
		return
	}

	if lowerLimit > TradeUpperLimitRMB {
		errCode = controllers.ERROR_CODE_OTC_TRADE_TOO_BIG
		return
	}

	// 重置卖币统计
	if wealth.BuyRmbToday > 0 {
		today, _ := common2.TodayTimeRange()
		if wealth.BuyUTime < today {
			eosplus.EosPlusAPI.Otc.BuySettingResetState(wealth.Uid, today)
			wealth.SellRmbToday = 0
		}
	}

	//更新用户买币（承兑商出售）设置
	errCode2 := eosplus.EosPlusAPI.Otc.BuySetting(wealth.Uid, able, dayLimit, lowerLimit)
	errCode = controllers.ERROR_CODE(errCode2)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		return
	}
	wealth.BuyAble = able
	wealth.BuyRmbLowerLimit = lowerLimit
	wealth.BuyRmbDay = dayLimit

	res = wealth
	return
}

//承兑商 - OTC出售
func (b *BuyLogic) Start(wealth *models2.EosOtc) (errCode controllers.ERROR_CODE) {
	errCode = controllers.ERROR_CODE_SUCCESS
	defer func() {
		if errCode != controllers.ERROR_CODE_SUCCESS {
			_ = dao.OtcBuyDaoEntity.DeleteByUid(wealth.Uid)
			eosplus.EosPlusAPI.Otc.StopBuy(wealth.Uid, controllers.ErrorMsg(errCode, "zh"))
		}
	}()
	if wealth.Status != dao2.WealthStatusWorking {
		errCode = controllers.ERROR_CODE_OTC_ACCOUNT_LOCK
	}
	TradeLowerLimitRMB := GetTradeLowerLimit()
	common2.LogFuncDebug("wealth:%v", eosplus.ToJsonIndent(wealth))
	//设置关闭直接关闭返回，如有已经匹配交易继续进行
	if !wealth.BuyAble {
		_ = dao.OtcBuyDaoEntity.DeleteByUid(wealth.Uid)
		return
	}

	//拥有EUSD转换出来的RMB
	avRmb := BuyEUSD2RMBWithPrecision(int64(wealth.Available))
	common2.LogFuncDebug("avRmb:%v , %v ", avRmb, wealth.Available)
	//达到限额
	if wealth.BuyRmbDay > 0 {
		if wealth.BuyRmbDay-wealth.BuyRmbToday < TradeLowerLimitRMB {
			_ = dao.OtcBuyDaoEntity.DeleteByUid(wealth.Uid)
			errCode = controllers.ERROR_CODE_OTC_DAY_LIMIT_SELL
			return
		}
		avRmb = MinInt64(avRmb, wealth.BuyRmbDay-wealth.BuyRmbToday)
	}

	//可用金额小于系统最小允许交易额度
	if avRmb < TradeLowerLimitRMB {
		_ = dao.OtcBuyDaoEntity.DeleteByUid(wealth.Uid)
		errCode = controllers.ERROR_CODE_OTC_LACK_AVAILABLE_TO_SELL
		return
	}

	// 获取支持的支付方式（包含剩余额度配置）
	TradeUpperLimitRMB := GetTradeUpperLimit()
	pay, llWx, ulWx, llAli, ulAli, llBank, ulBank := payAnalyse(wealth.Uid, wealth.BuyRmbLowerLimit, TradeLowerLimitRMB, TradeUpperLimitRMB)
	if pay == 0 {
		errCode = controllers.ERROR_CODE_OTC_LACK_PAYMENT
		return
	}
	// 写入待交易表 单位RMB
	sell := &models.OtcBuy{
		Uid:              wealth.Uid,
		Available:        avRmb,
		LowerLimitWechat: llWx,
		UpperLimitWechat: ulWx,
		LowerLimitAli:    llAli,
		UpperLimitAli:    ulAli,
		LowerLimitBank:   llBank,
		UpperLimitBank:   ulBank,
		PayType:          pay,
		Ctime:            common2.NowInt64MS(),
	}
	err := dao.OtcBuyDaoEntity.Edit(sell)
	if err != nil {
		errCode = controllers.ERROR_CODE_DB
		return
	}

	return
}

//承兑商 支付方式&限制整理
func payAnalyse(uid uint64, lowerLimit, TradeLowerLimitRMB, TradeUpperLimitRMB int64) (pay uint8,
	lowerLimitWechat, upperLimitWechat, lowerLimitAli, upperLimitAli, lowerLimitBank, upperLimitBank int64) {
	list := dao.PaymentMethodDaoEntity.FetchByUid(uid)

	t, _ := common2.TodayTimeRange()
	today := uint32(t)

	lowerLimit = MaxInt64(lowerLimit, TradeLowerLimitRMB)

	lowerLimitWechat = lowerLimit
	lowerLimitAli = lowerLimit
	lowerLimitBank = lowerLimit
	for _, v := range list {
		if v.Mtime < today {
			v.MoneyToday = 0
			v.TimesToday = 0
		}

		//额度不足
		residue := v.MoneyPerDayLimit - v.MoneyToday
		if residue < lowerLimit {
			continue
		}
		//剩余次数不足
		if v.TimesPerDayLimit-v.TimesToday <= 0 {
			continue
		}

		if v.MType == dao.PayModeWechat {
			pay = pay | dao.PayModeWechat
			if pay&dao.PayModeWechat == 0 {
				pay += dao.PayModeWechat
			}
			lowerLimitWechat = MinInt64(lowerLimitWechat, v.LowMoneyPerTxLimit)
			upperLimitWechat = MaxInt64(upperLimitWechat, MinInt64(v.HighMoneyPerTxLimit, residue))

		} else if v.MType == dao.PayModeAli {
			pay = pay | dao.PayModeAli
			lowerLimitAli = MinInt64(lowerLimitAli, v.LowMoneyPerTxLimit)
			upperLimitAli = MaxInt64(upperLimitAli, MinInt64(v.HighMoneyPerTxLimit, residue))

		} else if v.MType == dao.PayModeBank {
			if pay&dao.PayModeBank == 0 {
				pay += dao.PayModeBank
			}
			lowerLimitBank = MinInt64(lowerLimitBank, v.LowMoneyPerTxLimit)
			upperLimitBank = MaxInt64(upperLimitBank, MinInt64(v.HighMoneyPerTxLimit, residue))

		}
	}
	upperLimitWechat = MinInt64(upperLimitWechat, TradeUpperLimitRMB)
	upperLimitAli = MinInt64(upperLimitAli, TradeUpperLimitRMB)
	upperLimitBank = MinInt64(upperLimitBank, TradeUpperLimitRMB)

	return
}

// 用户下单购买 (同)
func (b *BuyLogic) MatchUp(uid uint64, ip string, userPayment *models.PaymentMethod, quantity int64) (order *models.OtcOrder, errCode controllers.ERROR_CODE) {
	// 金额检查
	funds := BuyEUSD2RMBWithPrecision(quantity)
	TradeLowerLimitRMB := GetTradeLowerLimit()
	TradeUpperLimitRMB := GetTradeUpperLimit()
	common2.LogFuncDebug("funds:%v ,lower:%v, upper:%v", funds, TradeLowerLimitRMB, TradeUpperLimitRMB)
	if funds < TradeLowerLimitRMB {
		errCode = controllers.ERROR_CODE_OTC_TRADE_TOO_SMALL
		return
	}
	if funds > TradeUpperLimitRMB {
		errCode = controllers.ERROR_CODE_OTC_TRADE_TOO_BIG
		return
	}

	w, errCode := eosplus.EosPlusAPI.Wealth.Info(uid)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		return
	}

	if w.Status != dao2.WealthStatusWorking {
		errCode = controllers.ERROR_CODE_WEALTH_LOCK
		return
	}
	// todo: 指定承兑商检查

	// 按支付方式 & 金额搜索承兑商
	list := []*models.OtcBuy{}
	var err error
	switch userPayment.MType {
	case dao.PayModeWechat:
		list, err = dao.OtcBuyDaoEntity.MatchWechat(funds)
	case dao.PayModeAli:
		list, err = dao.OtcBuyDaoEntity.MatchAli(funds)
	case dao.PayModeBank:
		list, err = dao.OtcBuyDaoEntity.MatchBank(funds)

	}
	if err != nil {
		errCode = controllers.ERROR_CODE_DB
		return
	}
	if len(list) == 0 {
		common2.LogFuncDebug("NO EXCHANGER")
		errCode = controllers.ERROR_CODE_OTC_NOT_MATCH_EXCHANGER
		return
	}
	list = ShuffleBuy(list)
	order = &models.OtcOrder{}

	// 承兑商金额 & 支付方式检查 & 锁定
	for _, v := range list {
		if v.Uid == uid {
			common2.LogFuncDebug("%v UID Same", v.Uid)
			continue
		}
		if v.Available < funds {
			common2.LogFuncDebug("%v Available Lock", v.Uid)
			continue
		}

		// 锁定承兑商
		key := fmt.Sprintf("otc_sell_%d", v.Uid)
		n := common2.RedisManger.Incr(key)

		if n.Val() != 1 {
			common2.LogFuncDebug("%v Lock", v.Uid)
			continue
		}
		//锁定时间15分钟
		common2.RedisManger.Expire(key, TradeMatchLock)
		//获取支付方式
		payment := b.selectExchangerPayment(v.Uid, userPayment.MType, funds)
		if payment.Pmid == 0 {
			common2.RedisManger.Del(key)
			common2.LogFuncDebug("%v NO Payment", v.Uid)
			continue
		}
		//锁定支付方式
		if !dao.PaymentMethodDaoEntity.Use(payment.Pmid, funds, payment.TimesToday) {
			common2.RedisManger.Del(key)
			common2.LogFuncDebug("%v Payment Lock", v.Uid)
			continue
		}
		//锁定token
		errCode2 := eosplus.EosPlusAPI.Otc.TransferLock(v.Uid, quantity)
		if controllers.ERROR_CODE(errCode2) != controllers.ERROR_CODE_SUCCESS {
			common2.RedisManger.Del(key)
			if !dao.PaymentMethodDaoEntity.RollUse(payment.Pmid, funds) {
				common2.LogFuncError("Payment Roll %v", payment.Pmid)
			}
			common2.LogFuncDebug("%v Token Lock %v, funds %v", v.Uid, errCode2, quantity)
			continue
		}
		// 修改撮合数据表
		err := dao.OtcBuyDaoEntity.EditAvailable(v.Uid, funds)
		if err != nil {
			common2.LogFuncDebug("%v Edit Otc", v.Uid)
			common2.RedisManger.Del(key)
			if !dao.PaymentMethodDaoEntity.RollUse(payment.Pmid, funds) {
				common2.LogFuncError("Payment Roll %v", payment.Pmid)
			}
			errCode2 := eosplus.EosPlusAPI.Otc.TransferUnLock(v.Uid, quantity)
			if errCode2 != controllers.ERROR_CODE_SUCCESS {
				common2.LogFuncError("TransferUnLock Fail:%v, %v", v.Uid, quantity)
			}
			continue
		}

		// 生成订单
		order = &models.OtcOrder{
			Uid:            uid,
			Uip:            ip,
			EUid:           v.Uid,
			Side:           dao.SideBuy,
			Amount:         int64(quantity),
			Price:          fmt.Sprintf("%.4f", BuyUSDT2RMBRate()),
			Funds:          funds,
			Fee:            0,
			PayId:          userPayment.Pmid,
			PayType:        int8(userPayment.MType),
			PayName:        userPayment.Name,
			PayAccount:     userPayment.Account,
			PayBank:        userPayment.Bank,
			PayBankBranch:  userPayment.BankBranch,
			Ctime:          common2.NowInt64MS(),
			Status:         dao.OrderStatusCreated,
			EPayId:         payment.Pmid,
			EPayType:       int8(payment.MType),
			EPayName:       payment.Name,
			EPayAccount:    payment.Account,
			EPayBank:       payment.Bank,
			EPayBankBranch: payment.BankBranch,
		}
		order, err = dao.OrdersDaoEntity.Create(order)
		if err != nil {
			common2.LogFuncDebug("%v Create Order", v.Uid)

			common2.RedisManger.Del(key)

			if !dao.PaymentMethodDaoEntity.RollUse(payment.Pmid, funds) {
				common2.LogFuncError("Payment Roll %v", payment.Pmid)
			}
			errCode2 := eosplus.EosPlusAPI.Otc.TransferUnLock(v.Uid, quantity)
			if errCode2 != controllers.ERROR_CODE_SUCCESS {
				common2.LogFuncError("TransferUnLock Fail:%v, %v", v.Uid, quantity)
			}
			_ = dao.OtcBuyDaoEntity.EditAvailable(v.Uid, -funds)
			continue
		}

		//解锁
		common2.RedisManger.Del(key)
		break
	}

	if err != nil || order.Id == 0 {
		common2.LogFuncDebug("No Orders")

		errCode = controllers.ERROR_CODE_OTC_NOT_MATCH_EXCHANGER
		return
	}
	//通知承兑商
	title := otc_common.ExchangerNewOrder1Title
	body := fmt.Sprintf(otc_common.ExchangerNewOrder1Body, order.Id, float64(order.Amount)/10000, float64(order.Funds)/100)
	go common2.SafeRun(func() {
		p := new(uemng_plus.UPushPlus)
		p.PushOtcOrder(order.EUid, order.Id, body, title, "1", order.Status)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.EUid)
	})()

	//设定定时检查订单 改成MQ
	//time.AfterFunc(PayExpire(), func() {
	//	b.TimeOutOrder(order.Id)
	//})
	errCode = controllers.ERROR_CODE_SUCCESS
	return
}

// 选择承兑商支付方式
func (b *BuyLogic) selectExchangerPayment(uid uint64, payType uint8, funds int64) (pay *models.PaymentMethod) {
	pay = &models.PaymentMethod{}
	list := dao.PaymentMethodDaoEntity.FetchByUidPayType(uid, payType)
	if len(list) == 0 {
		return
	}

	today, _ := common2.TodayTimeRange()
	todayUint32 := uint32(today)
	for _, v := range list {
		if v.UseTime < todayUint32 && (v.MoneyToday != 0 || v.TimesToday != 0) {
			v.MoneyToday = 0
			v.TimesToday = 0
		}
		if v.LowMoneyPerTxLimit > funds { //每笔交易最低
			continue
		}
		if v.HighMoneyPerTxLimit < funds { //每笔交易最高
			continue
		}
		if v.MoneyPerDayLimit-v.MoneyToday < funds { //每笔交易最高
			continue
		}
		if v.TimesPerDayLimit-v.TimesToday <= 0 { //每日次数
			continue
		}
		if v.MoneySumLimit-int64(v.MoneySum) < funds { //使用总额
			continue
		}

		pay = v

		if v.UseTime < todayUint32 { //支付方式重新计数
			_ = dao.PaymentMethodDaoEntity.ResetToday(v.Pmid)
		}
		break
	}

	//今天未使用 & 数据未重置 的支付方式进行重置
	if pay.UseTime < todayUint32 && pay.TimesToday == -10 {
		dao.PaymentMethodDaoEntity.FlushUse(pay.Pmid, todayUint32)
	}

	return
}

func ShuffleBuy(a []*models.OtcBuy) []*models.OtcBuy {
	for i := len(a) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

// 确认支付
func (b *BuyLogic) PayOrder(oid uint64, uid uint64) (errCode controllers.ERROR_CODE) {
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		errCode = controllers.ERROR_CODE_DB
		return
	}
	if order.Side != dao.SideBuy {
		errCode = controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
		return
	}
	if order.Id == 0 {
		errCode = controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
		return
	}

	if order.Uid != uid {
		errCode = controllers.ERROR_CODE_NO_AUTH
		return
	}
	ok := dao.OrdersDaoEntity.Pay(oid)
	if !ok {
		errCode = controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR
		return
	}

	// 通知
	title := otc_common.BuerPayedTitle
	body := fmt.Sprintf(otc_common.BuerPayedBody, order.Id, float64(order.Amount)/10000, float64(order.Funds)/100)
	go common2.SafeRun(func() {
		p := new(uemng_plus.UPushPlus)
		p.PushOtcOrder(order.EUid, order.Id, body, title, "1", dao.OrderStatusPayed)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.EUid)
	})()

	//设定定时检查订单
	//time.AfterFunc(ConfirmExpire(), func() {
	//	b.timeConfirmOrder(oid)
	//})
	return controllers.ERROR_CODE_SUCCESS
}

//(承兑商)确认收款
func (b *BuyLogic) ConfirmOrder(uid, oid uint64, ip string) (order *models.OtcOrder, errCode controllers.ERROR_CODE) {
	errCode = CheckAppealOrder(uid, oid)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		return
	}

	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		errCode = controllers.ERROR_CODE_DB
		return
	}
	if order.Side != dao.SideBuy {
		errCode = controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
		return
	}
	if order.Id == 0 {
		errCode = controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
	}
	if order.EUid != uid {
		errCode = controllers.ERROR_CODE_NO_AUTH
		return
	}
	if order.Status != dao.OrderStatusPayed {
		errCode = controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR
		return
	}

	ok := dao.OrdersDaoEntity.Confirmed(order.Id, ip)
	if !ok {
		errCode = controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR
		return
	}
	DayAmountReport(order.Side, order.Uid, order.EUid, order.Amount)

	go common2.SafeRun(func() {
		//区块链转账
		eosplus.EosPlusAPI.Otc.TransferOut(order.Id, order.EUid, order.Uid, order.Amount)
	})()
	//通知承兑商订单被取消、 通知用户
	title := otc_common.ExchangerSendedTitle
	body := fmt.Sprintf(otc_common.ExchangerSendedBody, order.EUid, order.Id, float64(order.Funds)/100)
	go common2.SafeRun(func() {
		p := new(uemng_plus.UPushPlus)
		p.PushOtcOrder(order.Uid, order.Id, body, title, "0", dao.OrderStatusConfirmed)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.Uid)
	})()

	errCode = controllers.ERROR_CODE_SUCCESS
	return
}

//定时检查订单 - 自动确认
func (b *BuyLogic) timeConfirmOrder(oid uint64) {
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil || order.Id == 0 {
		common2.LogFuncError("timeConfirmOrder not found:%v", oid)
		return
	}
	if order.Status != dao.OrderStatusPayed {
		common2.LogFuncDebug("timeConfirmOrder Status:%v,%v", oid, order.Status)
		return
	}

	ok := dao.OrdersDaoEntity.Confirmed(order.Id, "sys")
	if !ok {
		common2.LogFuncDebug("timeConfirmOrder Db:%v,%v", oid, order.Status)
		return
	}

	go common2.SafeRun(func() {
		//区块链转账
		eosplus.EosPlusAPI.Otc.TransferOut(order.Id, order.EUid, order.Uid, order.Amount)
	})()

	//通知承兑商
	p := new(uemng_plus.UPushPlus)
	title := otc_common.BuerPayedTitle
	if err != nil {
		common2.LogFuncError("%v", err)
	}
	body := fmt.Sprintf(otc_common.BuerPayedBody, order.Id, float64(order.Amount)/10000, float64(order.Funds)/100)
	go common2.SafeRun(func() {
		p.PushOtcOrder(order.EUid, order.Id, body, title, "1", dao.OrderStatusConfirmed)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.EUid)
	})()
}

//取消订单
func (b *BuyLogic) CancelOrder(oid uint64, uid uint64) (errCode controllers.ERROR_CODE) {
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		errCode = controllers.ERROR_CODE_DB
		return
	}
	if order.Side != dao.SideBuy {
		errCode = controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
		return
	}
	if order.Id == 0 {
		errCode = controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
		return
	}
	if order.Uid != uid {
		errCode = controllers.ERROR_CODE_NO_AUTH
		return
	}
	if order.Status != dao.OrderStatusCreated {
		errCode = controllers.ERROR_CODE_OTC_ORDER_CANNOT_CANCEL
		return
	}
	ok := dao.OrdersDaoEntity.Cancel(oid)
	if !ok {
		errCode = controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR
		return
	}
	// 回滚支付方式
	if !dao.PaymentMethodDaoEntity.RollUse(order.EPayId, order.Funds) {
		common2.LogFuncError("payment roll %v", order.Id)
	}

	// 解锁token
	errCode = eosplus.EosPlusAPI.Otc.TransferUnLock(order.EUid, order.Amount)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		common2.LogFuncError("%v", errCode)
	}
	// 交易列表记录
	_ = dao.OtcBuyDaoEntity.EditAvailable(order.EUid, -order.Funds)

	//通知承兑商订单被取消
	p := new(uemng_plus.UPushPlus)
	title := otc_common.BuerOrderCancelTitle
	body := fmt.Sprintf(otc_common.BuerOrderCancelBody, order.Id)
	go common2.SafeRun(func() {
		p.PushOtcOrder(order.EUid, order.Id, body, title, "1", dao.OrderStatusCanceled)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.EUid)

	})()
	return controllers.ERROR_CODE_SUCCESS
}

//超时取消订单
func (b *BuyLogic) TimeOutOrder(oid uint64) {
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		common2.LogFuncError("ERR:%v", err)
		return
	}
	if order.Id == 0 {
		common2.LogFuncError("NO FOUND ORDER")
		return
	}

	if order.Status != dao.OrderStatusCreated {
		common2.LogFuncError("ORDER STATUS ERR")
		return
	}
	ok := dao.OrdersDaoEntity.Timeout(oid)
	if !ok {
		return
	}
	if !dao.PaymentMethodDaoEntity.RollUse(order.EPayId, order.Funds) {
		common2.LogFuncError("payment roll %v", order.Id)
	}
	// 解锁token
	if eosplus.EosPlusAPI.Otc.TransferUnLock(order.EUid, order.Amount) != controllers.ERROR_CODE_SUCCESS {
		common2.LogFuncError("Token unlock fail %v", order.Id)
		return
	}
	// 交易列表记录
	_ = dao.OtcBuyDaoEntity.EditAvailable(order.EUid, -order.Funds)
	//通知承兑商订单被取消、 通知用户
	p := new(uemng_plus.UPushPlus)
	title := otc_common.SysterOrderCancelTitle
	body := fmt.Sprintf(otc_common.SysterOrderCancelBody, order.Id)
	go common2.SafeRun(func() {
		p.PushOtcOrder(order.EUid, order.Id, body, title, "1", dao.OrderStatusExpired)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.EUid)
		p.PushOtcOrder(order.Uid, order.Id, body, title, "0", dao.OrderStatusExpired)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.Uid)
	})()
}

func (b *BuyLogic) UpdateAvailable(uid uint64, quantity int64) {

}
