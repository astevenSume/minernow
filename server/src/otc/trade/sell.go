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

type SellLogic struct {
}

var SellLogicEntity = &SellLogic{}

//承兑商otc信息
func (b *SellLogic) Info(uid uint64) (wealth *models2.EosOtc, errCode controllers.ERROR_CODE) {
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

// 承兑商设置购买token配置
func (b *SellLogic) Setting(wealth *models2.EosOtc, able bool, pay uint8, dayLimit, lowLimit int64) (res *models2.EosOtc, errCode controllers.ERROR_CODE) {
	TradeLowerLimitRMB := GetTradeLowerLimit()
	if lowLimit > 0 && lowLimit < TradeLowerLimitRMB {
		errCode = controllers.ERROR_CODE_OTC_TRADE_TOO_SMALL
		return
	}

	// 重置卖币统计
	if wealth.SellRmbToday > 0 {
		today, _ := common2.TodayTimeRange()
		if wealth.SellUTime < today {
			eosplus.EosPlusAPI.Otc.SellSettingResetState(wealth.Uid, today)
			wealth.SellRmbToday = 0
		}
	}
	if dayLimit < TradeLowerLimitRMB {
		errCode = controllers.ERROR_CODE_OTC_TRADE_TOO_SMALL
	}

	//写入配置
	errCode2 := eosplus.EosPlusAPI.Otc.SellSetting(wealth.Uid, able, dayLimit, lowLimit, pay)
	errCode = controllers.ERROR_CODE(errCode2)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		errCode = controllers.ERROR_CODE_OTC_SETTING_ERR
	}

	wealth.SellRmbDay = dayLimit
	wealth.SellRmbLowerLimit = lowLimit
	wealth.SellAble = able
	wealth.SellPayType = pay

	res = wealth
	return
}

// 开启购买
func (b *SellLogic) Start(wealth *models2.EosOtc) (errCode controllers.ERROR_CODE) {
	errCode = controllers.ERROR_CODE_SUCCESS
	uid := wealth.Uid
	defer func() {
		if errCode != controllers.ERROR_CODE_SUCCESS {
			_ = dao.OtcSellDaoEntity.Delete(uid)
			eosplus.EosPlusAPI.Otc.StopSell(wealth.Uid, controllers.ErrorMsg(errCode, "zh"))
		}
	}()
	//用户资格检查
	if wealth.Status != dao2.WealthStatusWorking {
		if wealth.Account != "" {
			errCode = controllers.ERROR_CODE_OTC_ACCOUNT_LOCK
			return
		}
		return controllers.ERROR_CODE_OTC_NOT_A_EXCHANGER
	}

	// 关闭交易
	if wealth.SellAble == false {
		_ = dao.OtcSellDaoEntity.Delete(uid)
		return
	}
	TradeLowerLimitRMB := GetTradeLowerLimit()
	TradeUpperLimitRMB := GetTradeUpperLimit()
	// 可用额度
	av := int64(0) //无限额
	if wealth.SellRmbDay > 0 {
		av = wealth.SellRmbDay - wealth.BuyRmbToday
		if av < TradeLowerLimitRMB {
			return controllers.ERROR_CODE_OTC_DAY_LIMIT_SELL
		}
	}

	//最低限额
	ll := MaxInt64(TradeLowerLimitRMB, wealth.SellRmbLowerLimit)

	// 达到限额
	if av > 0 && av < ll {
		err := dao.OtcSellDaoEntity.Delete(uid)
		if err != nil {
			errCode = controllers.ERROR_CODE_DB
		}
		// 关闭OTC交易
		b.Setting(wealth, false, wealth.SellPayType, wealth.SellRmbDay, wealth.SellRmbLowerLimit)
		return
	}

	//写入待交易列表
	data := &models.OtcSell{
		Uid:        uid,
		Available:  av,
		LowerLimit: ll,
		UpperLimit: MinInt64(av, TradeUpperLimitRMB),
		PayType:    wealth.SellPayType,
		Ctime:      common2.NowInt64MS(),
	}
	err := dao.OtcSellDaoEntity.Edit(data)
	if err != nil {
		return controllers.ERROR_CODE_DB
	}
	return controllers.ERROR_CODE_SUCCESS
}

//用户卖币
func (b *SellLogic) MatchUp(uid uint64, ip string, payId uint64, quantity int64) (order *models.OtcOrder, errCode controllers.ERROR_CODE) {
	funds := SellEUSD2RMBWithPrecision(int64(quantity))
	TradeLowerLimitRMB := GetTradeLowerLimit()
	TradeUpperLimitRMB := GetTradeUpperLimit()
	// 金额检查
	if funds < TradeLowerLimitRMB {
		errCode = controllers.ERROR_CODE_OTC_TRADE_TOO_SMALL
		return
	}
	if funds > TradeUpperLimitRMB {
		errCode = controllers.ERROR_CODE_OTC_TRADE_TOO_BIG
		return
	}

	// 用户资金确认
	w, errCode2 := eosplus.EosPlusAPI.Wealth.Info(uid)
	errCode = controllers.ERROR_CODE(errCode2)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		errCode = controllers.ERROR_CODE_OTC_LACK_AVAILABLE_TO_SELL
		return
	}

	if w.Status != dao2.WealthStatusWorking {
		errCode = controllers.ERROR_CODE_WEALTH_LOCK
		return
	}

	if w.Available < quantity {
		errCode = controllers.ERROR_CODE_OTC_LACK_AVAILABLE_TO_SELL
		return
	}

	// 资金锁定
	eosplus.EosPlusAPI.Otc.UserTransferLock(uid, quantity)

	defer func() {
		//没有订单生成 释放资源
		if errCode != controllers.ERROR_CODE_SUCCESS {
			eosplus.EosPlusAPI.Otc.UserTransferUnlock(uid, quantity)
		}
	}()

	// 用户收款方式检查
	payment := dao.PaymentMethodDaoEntity.Info(payId)
	if payment.Uid != uid {
		errCode = controllers.ERROR_CODE_OTC_PAYMENT_NOT_FOUND
		return
	}

	// todo: 指定承兑商检查

	// 匹配承兑商
	list := dao.OtcSellDaoEntity.Match(int8(payment.MType), funds)
	if len(list) == 0 {
		common2.LogFuncDebug("NO EXCHANGER")
		errCode = controllers.ERROR_CODE_OTC_NOT_MATCH_EXCHANGER
		return
	}
	order = &models.OtcOrder{}
	list = ShuffleSell(list)

	for _, v := range list {
		if v.Uid == uid {
			common2.LogFuncDebug("EXCHANGER is self")
			continue
		}
		key := fmt.Sprintf("otc_buy_%d", v.Uid)
		n := common2.RedisManger.Incr(key)
		if n.Val() != 1 {
			common2.LogFuncDebug("EXCHANGER LOCK")
			continue
		}
		//锁定时间15分钟
		common2.RedisManger.Expire(key, TradeMatchLock)

		order = &models.OtcOrder{
			Uid:           uid,
			Uip:           ip,
			EUid:          v.Uid,
			Eip:           "",
			Side:          dao.SideSell,
			Amount:        int64(quantity),
			Price:         fmt.Sprintf("%.4f", SellUSDT2RMBRate()),
			Funds:         funds,
			Fee:           0,
			PayId:         payment.Pmid,
			PayType:       int8(payment.MType),
			PayName:       payment.Name,
			PayAccount:    payment.Account,
			PayBank:       payment.Bank,
			PayBankBranch: payment.BankBranch,
			Ctime:         common2.NowInt64MS(),
			Status:        dao.OrderStatusCreated,
		}

		order, _ = dao.OrdersDaoEntity.Create(order)
		if order.Id < 1 {
			common2.RedisManger.Del(key)
			continue
		}
		// 修改 匹配队列里的数据
		dao.OtcSellDaoEntity.EditAvailable(v.Uid, funds)

		common2.RedisManger.Del(key)
		break
	}
	if order.Id < 1 {
		errCode = controllers.ERROR_CODE_OTC_NOT_MATCH_EXCHANGER
		return
	}
	// 通知承兑商
	go common2.SafeRun(func() {
		p := new(uemng_plus.UPushPlus)
		title := otc_common.ExchangerNewOrder2Title
		body := fmt.Sprintf(otc_common.ExchangerNewOrder2Body, order.Id, float64(order.Amount)/10000, float64(order.Funds)/10000)
		p.PushOtcOrder(order.EUid, order.Id, body, title, "1", order.Status)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.EUid)
	})()

	//设定定时检查订单 - 超时取消  改成MQ
	//time.AfterFunc(PayExpire(), func() {
	//	b.TimeOutOrder(order.Id)
	//})
	errCode = controllers.ERROR_CODE_SUCCESS
	return
}

func ShuffleSell(a []*models.OtcSell) []*models.OtcSell {
	for i := len(a) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

//（承兑商） 确认支付
func (b *SellLogic) PayOrder(oid uint64, uid uint64, payId uint64) (errCode controllers.ERROR_CODE) {
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		errCode = controllers.ERROR_CODE_DB
		return
	}
	if order.Id == 0 {
		errCode = controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
		return
	}

	if order.EUid != uid {
		errCode = controllers.ERROR_CODE_NO_AUTH
		return
	}

	// 支付方式
	payment := dao.PaymentMethodDaoEntity.Info(payId)
	if payment.Status != dao.PaymentMethodStatusActivated {
		errCode = controllers.ERROR_CODE_PAYMENT_METHOD_UNABLE
		return
	}
	if payment.Uid != uid {
		errCode = controllers.ERROR_CODE_PAYMENT_METHOD_OWNER_ERROR
		return
	}
	if payment.MType != uint8(order.PayType) {
		errCode = controllers.ERROR_CODE_PAYMENT_METHOD_NO_MATCH
		return
	}

	// 更新承兑商支付方式 & 确认支付
	ok := dao.OrdersDaoEntity.ExchangerPay(oid, payment.Pmid, payment.MType, payment.Account, payment.Name, payment.Bank, payment.BankBranch)
	if !ok {
		errCode = controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR
		return
	}
	// 通知用户
	title := otc_common.ExchangerPayedTitle
	body := fmt.Sprintf(otc_common.ExchangerPayedBody, order.Id, float64(order.Amount)/10000, float64(order.Funds)/100)
	go common2.SafeRun(func() {
		p := new(uemng_plus.UPushPlus)
		p.PushOtcOrder(order.Uid, order.Id, body, title, "0", dao.OrderStatusPayed)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.Uid)
	})()
	//设定定时检查订单
	//time.AfterFunc(ConfirmExpire(), func() {
	//	b.TimeOutOrder(order.Id)
	//})
	return controllers.ERROR_CODE_SUCCESS
}

//（用户）确认收款
func (b *SellLogic) ConfirmOrder(uid, oid uint64, ip string) (order *models.OtcOrder, errCode controllers.ERROR_CODE) {
	errCode = CheckAppealOrder(uid, oid)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		return
	}

	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		errCode = controllers.ERROR_CODE_DB
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

	ok := dao.OrdersDaoEntity.Confirmed(order.Id, ip)
	if !ok {
		errCode = controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR
		return
	}
	DayAmountReport(order.Side, order.Uid, order.EUid, order.Amount)

	go common2.SafeRun(func() {
		//区块链转账
		eosplus.EosPlusAPI.Otc.UserTransferOut(order.Id, order.Uid, order.EUid, order.Amount)
	})()

	// 通知承兑商
	title := otc_common.BuerReceiptedTitle
	body := fmt.Sprintf(otc_common.BuerReceiptedBody, order.Uid, order.Id, float64(order.Amount)/10000)
	go common2.SafeRun(func() {
		p := new(uemng_plus.UPushPlus)
		p.PushOtcOrder(order.EUid, order.Id, body, title, "1", dao.OrderStatusConfirmed)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.EUid)
	})()
	errCode = controllers.ERROR_CODE_SUCCESS
	return
}

//（用户）确认收款 超时自动确认
func (b *SellLogic) timeConfirmOrder(oid uint64) {
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
		common2.LogFuncDebug("timeConfirmOrder DB:%v,%v", oid, order.Status)
		return
	}

	go common2.SafeRun(func() {
		//区块链转账
		eosplus.EosPlusAPI.Otc.UserTransferOut(order.Id, order.Uid, order.EUid, order.Amount)
	})()
	//没有自动确认步骤
	return
}

//（承兑商）取消订单
func (b *SellLogic) CancelOrder(oid uint64, uid uint64) (order *models.OtcOrder, errCode controllers.ERROR_CODE) {
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		errCode = controllers.ERROR_CODE_DB
		return
	}
	if order.Id == 0 {
		errCode = controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
		return
	}
	if order.EUid != uid {
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
	// 解锁用户token
	errCode = eosplus.EosPlusAPI.Otc.UserTransferUnlock(order.Uid, order.Amount)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		common2.LogFuncError("CancelOrder User Token,oid:%v", order.Id, order.Amount)
	}
	//还原 otc队列数据
	if !dao.OtcSellDaoEntity.EditAvailable(order.Uid, -order.Funds) {
		common2.LogFuncError("CancelOrder Otc Sell,oid:%v", order.Id, order.Amount)
	}
	//承兑商无法取消订单

	errCode = controllers.ERROR_CODE_SUCCESS
	return
}

//超时取消订单
func (b *SellLogic) TimeOutOrder(oid uint64) {
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		common2.LogFuncError("%v", err)
		return
	}
	if order.Id == 0 {
		common2.LogFuncError("Order not found %v", oid)
		return
	}

	if order.Status != dao.OrderStatusCreated {
		common2.LogFuncError("Order status err %v", oid)
		return
	}
	ok := dao.OrdersDaoEntity.Timeout(oid)
	if !ok {
		return
	}
	errCode := eosplus.EosPlusAPI.Otc.UserTransferUnlock(order.Uid, order.Amount)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		common2.LogFuncError("CancelOrder User Token,oid:%v", order.Id, order.Amount)
		return
	}
	//还原 otc队列数据
	if !dao.OtcSellDaoEntity.EditAvailable(order.Uid, -order.Funds) {
		common2.LogFuncError("CancelOrder Otc Sell,oid:%v", order.Id, order.Amount)
	}
	//通知承兑商订单被取消、 通知用户
	title := otc_common.SysterOrderCancelTitle
	body := fmt.Sprintf(otc_common.SysterOrderCancelBody, order.Id)
	go common2.SafeRun(func() {
		p := new(uemng_plus.UPushPlus)
		p.PushOtcOrder(order.Uid, order.Id, body, title, "0", dao.OrderStatusExpired)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.Uid)
		p.PushOtcOrder(order.EUid, order.Id, body, title, "1", dao.OrderStatusExpired)
		_, _ = dao.SystemNotificationdDaoEntity.InsertSystemNotification("system", body, order.EUid)
	})()
}
