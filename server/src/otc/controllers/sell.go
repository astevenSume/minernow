package controllers

import (
	"common"
	"eusd/eosplus"
	"fmt"
	"otc/trade"
	"otc_error"
	"utils/otc/dao"
)

type SellController struct {
	BaseController
}

//easyjson:json
type SellStartMsg struct {
	Able       bool   `json:"able"`
	LowerLimit int64  `json:"lower_limit"`
	DayLimit   int64  `json:"day_limit"`
	PayType    []int8 `json:"pay_type"`
}

// 承兑商 开启OTC购买
func (c *SellController) Start() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}
	//单位：RMB
	p := SellStartMsg{}
	err := c.GetPost(&p)
	if err != nil {
		return
	}

	if len(p.PayType) < 1 {
		c.ErrorResponse(controllers.ERROR_CODE_OTC_PAY_TYPE_EMPTY)
		return
	}
	pay := int8(0)
	for _, v := range p.PayType {
		pay += v
	}
	otc, errCode := trade.SellLogicEntity.Info(uid)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	otc, errCode = trade.SellLogicEntity.Setting(otc, p.Able, uint8(pay), p.DayLimit, p.LowerLimit)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	errCode = trade.SellLogicEntity.Start(otc)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	c.SuccessResponseWithoutData()
}

//easyjson:json
type SellPostMsg struct {
	Funds    int64  `json:"funds"`
	Quantity int64  `json:"quantity"`
	PayId    uint64 `json:"pay_id"`
}

//下单
func (c *SellController) Post() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	p := SellPostMsg{}
	err := c.GetPost(&p)
	if err != nil {
		return
	}

	//支付密码 && 二次验证
	if errCode := c.check2step(uid, false); errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	order, errCode := trade.SellLogicEntity.MatchUp(uid, c.GetIP(), p.PayId, p.Quantity)

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

//easyjson:json
type SellPayMsg struct {
	PayId uint64 `json:"pay_id"`
}

// 确认订单
func (c *SellController) Pay() {
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

	p := SellPayMsg{}
	err := c.GetPost(&p)
	if err != nil {
		return
	}

	errCode := trade.SellLogicEntity.PayOrder(oid, uid, p.PayId)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	c.SuccessResponseWithoutData()
}

// 确认收款
func (c *SellController) Confirm() {
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

	// 支付密码验证
	// 二次验证 - 短信验证
	if errCode := c.check2step(uid, false); errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	order, errCode := trade.SellLogicEntity.ConfirmOrder(uid, oid, c.GetIP())
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	c.SuccessResponseWithoutData()
	//更新承兑商的今日买入
	eosplus.EosPlusAPI.Otc.UpdateSellState(order.EUid, order.Funds)
	//更新otc设置
	otc, _ := trade.SellLogicEntity.Info(uid)
	trade.SellLogicEntity.Start(otc)
}

// 取消订单
func (c *SellController) Cancel() {
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

	_, errCode := trade.SellLogicEntity.CancelOrder(oid, uid)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	c.SuccessResponseWithoutData()
}

// 承兑列表
func (c *SellController) List() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	list, _ := dao.OtcSellDaoEntity.All()
	c.SuccessResponse(list)
}
