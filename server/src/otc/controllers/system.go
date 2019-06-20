package controllers

import (
	"common"
	"fmt"
	"math"
	common2 "otc/common"
	controllers "otc_error"
	"usdt/prices"
	dao1 "utils/admin/dao"
	common3 "utils/common"
)

type SystemController struct {
	BaseController
}

func (c *SystemController) Get() {
	list, err := common3.AppConfigMgr.FetchString(
		common3.BuyFeeRate, common3.SellFeeRate,
		common3.ServiceWechat, common3.InviteWechat, common3.PlatformName,
		common3.ContactTelegram, common3.ContactWeChat,
		common3.OtcTradeConfirmExpire, common3.OtcTradePayExpire,
		common3.OtcTradeLowerLimitRmb, common3.OtcTradeUpperLimitRmb,
	)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	c.SuccessResponse(list)
}

func (c *SystemController) SystemMSecTime() {
	type msg struct {
		Timestamp int64 `json:"timestamp"`
	}
	res := &msg{
		Timestamp: common.NowInt64MS(),
	}
	c.SuccessResponse(res)
}

func (c *SystemController) OtcExUsdtPrice() {
	type Msg struct {
		Currency    string `json:"currency"`
		Price       string `json:"price"`
		Precision   int    `json:"precision"`
		BuyFeeRate  string `json:"buy_fee_rate"`
		SellFeeRate string `json:"sell_fee_rate"`
	}

	price, err := prices.GetPrice(prices.PRICE_CURRENCY_TYPE_USDT)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	list, err := common3.AppConfigMgr.FetchString(common3.BuyFeeRate, common3.SellFeeRate)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	precision := common2.Cursvr.UsdtPrecision
	var msg = Msg{
		Currency:    price.Currency,
		Price:       fmt.Sprintf("%."+fmt.Sprint(precision)+"f", float64(price.PowPrice)/float64(math.Pow10(int(price.Pow)))),
		Precision:   precision,
		BuyFeeRate:  list["buy_fee_rate"],
		SellFeeRate: list["sell_fee_rate"],
	}
	c.SuccessResponse(msg)
}

func (c *SystemController) LastAppVersion() {
	system, _ := c.GetInt8(KeySystemInput)
	appVersion, err := dao1.AppVersionDaoEntity.QueryLastAppVersion(system)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	c.SuccessResponse(appVersion)
}

func (c *SystemController) GetEndpoint() {
	data, err := dao1.EndPointDaoEntity.GetAll()
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	res := map[string]interface{}{}
	res["list"] = data
	c.SuccessResponse(res)
}
