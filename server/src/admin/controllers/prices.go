package controllers

import (
	controllers "admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"math"
	"usdt/prices"
	"utils/admin/dao"
	usdtDao "utils/usdt/dao"
	usdtModels "utils/usdt/models"
)

const PRICE_USDT_KEY = "price_otc_ex_usdt_key"

type PriceController struct {
	BaseController
}

//获取usdt指导价
func (c *PriceController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadPrice, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint64(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadPrice, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := usdtDao.PricesDaoEntity.QueryById(id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadPrice, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadPrice, input, usdtDao.PricesDaoEntity.ClientPrice(&data))
	} else {
		market, err := c.GetInt("market", -1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadPrice, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		page, err := c.GetInt(KEY_PAGE, 1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadPrice, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadPrice, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		//分页查询
		res := map[string]interface{}{}
		input = fmt.Sprintf("{\"market\":%v,\"page\":%v,\"per_page\":%v}", market, page, perPage)
		count, data, err := usdtDao.PricesDaoEntity.QueryPageCondition(int8(market), page, perPage, prices.PRICE_CURRENCY_TYPE_USDT.CurrencyName())
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadPrice, controllers.ERROR_CODE_DB, input)
			return
		}

		meta := dao.PageInfo{
			Limit: perPage,
			Total: int(count),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = usdtDao.PricesDaoEntity.ClientPrices(data)
		c.SuccessResponseAndLog(OPActionReadPrice, input, res)
	}
}

//easyjson:json
type PriceCurResMsg struct {
	Price string `json:"price"`
}

//获取当前指导价
func (c *PriceController) GetCurPrice() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadCurPrice, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	price := usdtModels.Prices{}
	if b, err := common.RedisManger.Get(PRICE_USDT_KEY).Bytes(); err != nil {
		price, err = usdtDao.PricesDaoEntity.QueryLastPrice(prices.PRICE_CURRENCY_TYPE_USDT.CurrencyName())
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadCurPrice, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
			return
		}
	} else {
		if err := json.Unmarshal(b, &price); err != nil {
			if err != nil {
				c.ErrorResponseAndLog(OPActionReadCurPrice, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
				return
			}
		}
	}
	res := PriceCurResMsg{}
	res.Price = fmt.Sprintf("%.2f", float64(price.PowPrice)/float64(math.Pow10(int(price.Pow))))
	c.SuccessResponseAndLog(OPActionReadCurPrice, "", res)
}
