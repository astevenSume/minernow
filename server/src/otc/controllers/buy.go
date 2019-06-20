package controllers

import (
	"common"
	"eusd/eosplus"
	"fmt"
	"otc/trade"
	"otc_error"
	"utils/otc/dao"
)

//用户购买 - 承兑商卖出

type BuyController struct {
	BaseController
}

//easyjson:json
type BuyStartMsg struct {
	Able       bool  `json:"able"`
	LowerLimit int64 `json:"lower_limit"`
	DayLimit   int64 `json:"day_limit"`
}

//承兑商 - OTC出售
func (c *BuyController) Start() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	p := BuyStartMsg{}
	err := c.GetPost(&p)
	if err != nil {
		return
	}

	wealth, errCode := trade.BuyLogicEntity.Info(uid)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	wealth, errCode = trade.BuyLogicEntity.Setting(wealth, p.Able, p.DayLimit, p.LowerLimit)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	errCode = trade.BuyLogicEntity.Start(wealth)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	c.SuccessResponseWithoutData()
}

//easyjson:json
type BuyPostMsg struct {
	Funds    int64  `json:"funds"`
	Quantity int64  `json:"quantity"`
	PayId    uint64 `json:"pay_id"`
}

//下单
func (c *BuyController) Post() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	if !trade.RedisLock(uid) {
		c.ErrorResponse(controllers.ERROR_CODE_OP_TOO_FAST)
		return
	}

	p := BuyPostMsg{}
	err := c.GetPost(&p)
	if err != nil {
		return
	}

	payMent := dao.PaymentMethodDaoEntity.Info(p.PayId)
	if payMent.Uid != uid {
		c.ErrorResponse(controllers.ERROR_CODE_OTC_PAYMENT_NOT_FOUND)
		return
	}

	order, errCode := trade.BuyLogicEntity.MatchUp(uid, c.GetIP(), payMent, p.Quantity)

	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	err = common.RabbitMQPublishDelay(RabbitMQOrderCheck, RabbitMQOrderCheck,
		[]byte(fmt.Sprintf("%d", order.Id)), fmt.Sprintf("%d", uint64(trade.PayExpire().Seconds())*1000))
	if err != nil {
		common.LogFuncError("Order2MQ:%v", err)
	}

	c.SuccessResponse(order)
}

//取消订单
func (c *BuyController) Cancel() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}
	if !trade.RedisLock(uid) {
		c.ErrorResponse(controllers.ERROR_CODE_OP_TOO_FAST)
		return
	}
	oid, _ := c.GetParamUint64(":order_id")

	errCode := trade.BuyLogicEntity.CancelOrder(oid, uid)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	c.SuccessResponseWithoutData()
}

// 用户确认支付
func (c *BuyController) Pay() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}
	if !trade.RedisLock(uid) {
		c.ErrorResponse(controllers.ERROR_CODE_OP_TOO_FAST)
		return
	}
	oid, _ := c.GetParamUint64(":order_id")
	if oid < 1 {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	errCode := trade.BuyLogicEntity.PayOrder(oid, uid)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	c.SuccessResponseWithoutData()
}

// 承兑商确认收款
func (c *BuyController) Confirm() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}
	if !trade.RedisLock(uid) {
		c.ErrorResponse(controllers.ERROR_CODE_OP_TOO_FAST)
		return
	}
	oid, _ := c.GetParamUint64(":order_id")
	// 支付密码
	// 短信验证
	if errCode := c.check2step(uid, false); errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	order, errCode := trade.BuyLogicEntity.ConfirmOrder(uid, oid, c.GetIP())
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	c.SuccessResponseWithoutData()
	//更新承兑商的限额 & otc数据
	eosplus.EosPlusAPI.Otc.UpdateBuyState(order.EUid, order.Funds)

}

// 承兑列表
func (c *BuyController) List() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	list, _ := dao.OtcBuyDaoEntity.All()
	c.SuccessResponse(list)
}
