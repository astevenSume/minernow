package trade

import (
	common2 "common"
	"eusd/eosplus"
	"otc_error"
	"utils/otc/dao"
	"utils/otc/models"
)

type AppealLogic struct {
}

//申述逻辑
var AppealLogicEntity = &AppealLogic{}

//回滚订单
func rollBackOrder(order *models.OtcOrder) controllers.ERROR_CODE {
	if order == nil {
		return controllers.ERROR_CODE_PARAMS_ERROR
	}
	if order.Side == dao.SideBuy {
		// 回滚支付方式
		if !dao.PaymentMethodDaoEntity.RollUse(order.EPayId, order.Funds) {
			common2.LogFuncError("payment roll %v", order.Id)
		}

		// 解锁token
		errCode := eosplus.EosPlusAPI.Otc.TransferUnLock(order.EUid, order.Amount)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			common2.LogFuncError("%v", errCode)
		}
		// 交易列表记录
		_ = dao.OtcBuyDaoEntity.EditAvailable(order.EUid, -order.Funds)
		//todo 通知承兑商订单被取消
	} else {
		// 解锁用户token
		errCode := eosplus.EosPlusAPI.Otc.UserTransferUnlock(order.Uid, order.Amount)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			common2.LogFuncError("CancelOrder User Token,oid:%v", order.Id, order.Amount)
		}
		//还原 otc队列数据
		if !dao.OtcSellDaoEntity.EditAvailable(order.Uid, -order.Funds) {
			common2.LogFuncError("CancelOrder Otc Sell,oid:%v", order.Id, order.Amount)
		}
		//todo 通知用户订单被取消
	}

	return controllers.ERROR_CODE_SUCCESS
}

func lockUser(isExchanger bool, uid uint64) {
	if isExchanger {
		errCode := eosplus.EosPlusAPI.Otc.Lock(uid)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			common2.LogFuncError("lockSell errCode:%v", errCode)
		}
	}

	errCode := eosplus.EosPlusAPI.Wealth.Lock(uid)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		common2.LogFuncError("lockSell errCode:%v", errCode)
	}
}

func lockSell(order *models.OtcOrder) (controllers.ERROR_CODE, bool) {
	if order == nil {
		return controllers.ERROR_CODE_PARAMS_ERROR, false
	}

	if order.Side == dao.SideBuy {
		//承兑商卖家
		err := dao.OtcBuyDaoEntity.EditAvailable(order.EUid, order.Funds)
		if err != nil {
			//冻结卖家
			common2.LogFuncError("lockSell error:%v", err)
			lockUser(true, order.EUid)
			return controllers.ERROR_CODE_SUCCESS, true
		}
		// 资金锁定
		errCode := eosplus.EosPlusAPI.Otc.TransferLock(order.EUid, order.Amount)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			//资金锁定失败冻结卖家
			common2.LogFuncError("lockSell errCode:%v", errCode)
			lockUser(true, order.EUid)
			return controllers.ERROR_CODE_SUCCESS, true
		}
	} else {
		//非承兑商卖家
		err := dao.OtcBuyDaoEntity.EditAvailable(order.Uid, order.Funds)
		if err != nil {
			//冻结卖家
			lockUser(false, order.Uid)
			return controllers.ERROR_CODE_SUCCESS, true
		}
		// 资金锁定
		errCode := eosplus.EosPlusAPI.Otc.UserTransferLock(order.Uid, order.Amount)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			//资金锁定失败冻结卖家
			lockUser(false, order.Uid)
			return controllers.ERROR_CODE_SUCCESS, true
		}
	}

	return controllers.ERROR_CODE_SUCCESS, false
}

//客服申述取消订单(买家点已付款但实际未付款)
func (b *AppealLogic) CancelOrder(oid, uid uint64, ip string) controllers.ERROR_CODE {
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		return controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
	}
	if uid != order.EUid && uid != order.Uid {
		return controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
	}

	//买家未付款误点确认 或者买家付款被退回 或者买家故意不付款点确认
	if order.Status != dao.OrderStatusPayed {
		return controllers.ERROR_CODE_OTC_ORDER_CANNOT_CANCEL
	}

	ok := dao.OrdersDaoEntity.AppealChangeStatus(oid, ip, dao.OrderStatusPayed, dao.OrderStatusCanceled)
	if !ok {
		return controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR
	}

	errCode := rollBackOrder(order)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		//回滚订单状态
		dao.OrdersDaoEntity.AppealChangeStatus(oid, ip, dao.OrderStatusCanceled, dao.OrderStatusPayed)
		return errCode
	}

	//更新申述状态
	_, err = dao.AppealDaoEntity.Resolve(oid, uid)
	if err != nil {
		common2.LogFuncError("error:%v", err)
		return controllers.ERROR_CODE_DB
	}

	return controllers.ERROR_CODE_SUCCESS
}

//客服申述确认已付款(买家点取消或系统超时自动取消)
func (b *AppealLogic) ConfirmOrderPay(oid, uid uint64, ip string) controllers.ERROR_CODE {
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		return controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
	}
	if uid != order.EUid && uid != order.Uid {
		return controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
	}

	//买家已付款但超时取消 或者买家已付款误点取消
	if order.Status != dao.OrderStatusCanceled && order.Status != dao.OrderStatusExpired {
		return controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR
	}

	ok := dao.OrdersDaoEntity.AppealChangeStatus(order.Id, ip, order.Status, dao.OrderStatusPayed)
	if !ok {
		return controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR
	}

	//锁定卖家资金
	errCode, isLock := lockSell(order)
	if errCode != controllers.ERROR_CODE_SUCCESS || isLock {
		//回滚订单状态
		dao.OrdersDaoEntity.AppealChangeStatus(order.Id, ip, dao.OrderStatusPayed, order.Status)
		return errCode
	}

	return controllers.ERROR_CODE_SUCCESS
}

//客服申述确认放币(买家已付款卖家未确认)
func (b *AppealLogic) ConfirmOrder(oid, uid uint64, ip string) controllers.ERROR_CODE {
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		return controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
	}
	if uid != order.EUid && uid != order.Uid {
		return controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND
	}

	//买家已付款未确认
	if order.Status != dao.OrderStatusPayed {
		return controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR
	}

	ok := dao.OrdersDaoEntity.Confirmed(order.Id, ip)
	if !ok {
		return controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR
	}
	DayAmountReport(order.Side, order.Uid, order.EUid, order.Amount)

	if order.Side == dao.SideBuy {
		go common2.SafeRun(func() {
			//区块链转账
			eosplus.EosPlusAPI.Otc.TransferOut(order.Id, order.EUid, order.Uid, order.Amount)
		})()
		// todo 通知用户
	} else {
		go common2.SafeRun(func() {
			//区块链转账
			eosplus.EosPlusAPI.Otc.UserTransferOut(order.Id, order.Uid, order.EUid, order.Amount)
		})()
		// todo 通知承兑商
	}

	_, err = dao.AppealDaoEntity.Resolve(oid, uid)
	if err != nil {
		return controllers.ERROR_CODE_DB
	}

	return controllers.ERROR_CODE_SUCCESS
}
