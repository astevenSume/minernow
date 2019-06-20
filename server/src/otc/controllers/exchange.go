package controllers

import (
	"common"
	"eusd/eosplus"
	json "github.com/mailru/easyjson"
	"otc/trade"
	"otc_error"
	"utils/eusd/models"
	"utils/otc/dao"
)

type ExchangeController struct {
	BaseController
}

//easyjson:json
type ExchangeApplyMsg struct {
	Wechat   string `json:"wechat"`
	Mobile   string `json:"mobile"`
	Telegram string `json:"telegram"`
}

//申请
func (c *ExchangeController) Apply() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	msg := ExchangeApplyMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	if msg.Mobile == "" {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}
	row, err := dao.OtcExchangerVerifyDaoEntity.LastRow(uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	if row.Id > 0 && row.Status == dao.ExchangerVerifyStatusPending {
		c.ErrorResponse(controllers.ERROR_CODE_OTC_APPLY_EXCHANGER_PENDING)
		return
	}

	result, err := dao.OtcExchangerVerifyDaoEntity.Create(uid, dao.ExchangerApply, msg.Mobile, msg.Wechat, msg.Telegram)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
	}
	_, _ = dao.UserDaoEntity.UpdateUser(uid, dao.ExchangerPending, dao.UserStatusNil, "", "")
	data := map[string]interface{}{
		"id":     result.Id,
		"uid":    uid,
		"status": "pending",
		"ctime":  common.NowInt64MS(),
	}
	c.SuccessResponse(data)
}

//查看otc账户信息
func (c *ExchangeController) Info() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	d, errCode := eosplus.EosPlusAPI.Otc.Info(uid)
	if controllers.ERROR_CODE(errCode) != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(controllers.ERROR_CODE(errCode))
	}
	//APP以trade做冻结显示
	d.Trade += d.Transfer
	//获取订单数
	sell, buy, _ := dao.OrdersDaoEntity.GetExchangerUnDealOrders(uid)

	type Res struct {
		*models.EosOtc
		Sell int `json:"num_sell_order"`
		Buy  int `json:"num_buy_order"`
	}

	c.SuccessResponse(Res{EosOtc: d, Sell: sell, Buy: buy})
}

//easyjson:json
type ExchangeTransferIntoMsg struct {
	Quantity int64 `json:"quantity"`
}

//划转 - 转入OTC
func (c *ExchangeController) TransferInto() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	msg := ExchangeTransferIntoMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	otcWealth, errCode := eosplus.EosPlusAPI.Otc.WealthTransferInto(uid, msg.Quantity)
	if controllers.ERROR_CODE(errCode) != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(controllers.ERROR_CODE(errCode))
		return
	}

	if otcWealth.BuyAble == true {
		funds := trade.BuyEUSD2RMBWithPrecision(msg.Quantity)
		_ = dao.OtcBuyDaoEntity.EditAvailable(otcWealth.Uid, -funds)
	}

	c.SuccessResponseWithoutData()
}

//easyjson:json
type ExchangeTransferOutMsg struct {
	Quantity int64 `json:"quantity"`
}

//划转 - 转出OTC
func (c *ExchangeController) TransferOut() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	msg := ExchangeTransferOutMsg{}
	err := c.GetPost(&msg)
	if err != nil {
		return
	}

	otcWealth, errCode := eosplus.EosPlusAPI.Otc.TransferToWealth(uid, msg.Quantity)
	if controllers.ERROR_CODE(errCode) != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(controllers.ERROR_CODE(errCode))
		return
	}

	if otcWealth.BuyAble == true {
		funds := trade.BuyEUSD2RMBWithPrecision(msg.Quantity)
		if otcWealth.Trade == 0 && trade.BuyEUSD2RMBWithPrecision(otcWealth.Available-msg.Quantity) < trade.GetTradeLowerLimit() {
			//余额不足 && 没有交易中的订单
			_ = dao.OtcBuyDaoEntity.DeleteByUid(otcWealth.Uid)

		} else {
			_ = dao.OtcBuyDaoEntity.EditAvailable(otcWealth.Uid, funds)
		}

	}

	c.SuccessResponseWithoutData()
}
