package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"strconv"
	"utils/admin/dao"
	eusdDao "utils/eusd/dao"
	otcDao "utils/otc/dao"
	"utils/otc/models"
)

type OtcExchangerController struct {
	BaseController
}

//获取承兑商
func (c *OtcExchangerController) GetExchanger() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadExchanger, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUid, err := c.GetUint64(KEY_UID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadExchanger, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	if otcUid > 0 {
		input = fmt.Sprintf("{\"id\":%v}", otcUid)
		data, err := otcDao.OtcExchangerDaoEntity.QueryById(otcUid)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadExchanger, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadExchanger, input, otcDao.OtcExchangerDaoEntity.ClientExchanger(&data))
	} else {
		page, err := c.GetInt(KEY_PAGE, 1)
		if err != nil {
			common.LogFuncError("err:%v", err)
			c.ErrorResponseAndLog(OPActionReadExchanger, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		}
		perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
		if err != nil {
			common.LogFuncError("err:%v", err)
			c.ErrorResponseAndLog(OPActionReadExchanger, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		}
		status, err := c.GetInt("status", -1)
		if err != nil {
			common.LogFuncError("err:%v", err)
			c.ErrorResponseAndLog(OPActionReadExchanger, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		}
		mobile := c.GetString("mobile")
		wechat := c.GetString("wechat")

		input = fmt.Sprintf("{\"mobile\":\"%v\",\"wechat\":\"%v\",\"status\":%v,\"page\":%v,"+
			"\"limit\":%v}", mobile, wechat, status, page, perPage)
		total, data, err := otcDao.OtcExchangerDaoEntity.QueryCondition(mobile, wechat, int8(status), page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadExchanger, controllers.ERROR_CODE_DB, input)
			return
		}

		var list []otcDao.OtcExchangerAck
		for _, item := range data {
			list = append(list, otcDao.OtcExchangerDaoEntity.ClientExchanger(&item))
		}
		res := map[string]interface{}{}
		meta := dao.PageInfo{
			Limit: perPage,
			Total: int(total),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = list
		c.SuccessResponseAndLog(OPActionReadExchanger, input, res)
	}
}

//创建承兑商
func (c *OtcExchangerController) CreateExchanger() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddExchanger, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id       string `json:"user_id"`
		Wechat   string `json:"wechat"`
		Mobile   string `json:"mobile"`
		Telegram string `json:"telegram"`
	}
	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddExchanger, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var otcUid uint64
	if len(req.Id) > 0 {
		var err error
		otcUid, err = strconv.ParseUint(req.Id, 10, 64)
		if err != nil {
			c.ErrorResponseAndLog(OPActionAddExchanger, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
	}

	_, err = BecomeExchanger(otcUid, otcDao.ExchangerAssign, req.Mobile, req.Wechat, req.Telegram)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddExchanger, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionAddExchanger, string(c.Ctx.Input.RequestBody))
}

//更新承兑商
func (c *OtcExchangerController) UpdateExchanger() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditExchanger, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Uid      string `json:"user_id"`
		Wechat   string `json:"wechat"`
		Mobile   string `json:"mobile"`
		Telegram string `json:"telegram"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditExchanger, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	var otcId uint64
	if len(req.Uid) > 0 {
		otcId, err = strconv.ParseUint(req.Uid, 10, 64)
		if err != nil {
			c.ErrorResponseAndLog(OPActionEditExchanger, controllers.ERROR_CODE_PARAMS_ERROR, req.Uid)
			return
		}
	}

	_, err = otcDao.OtcExchangerDaoEntity.UpdateExchanger(otcId, req.Mobile, req.Wechat, req.Telegram)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditExchanger, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionEditExchanger, string(c.Ctx.Input.RequestBody))
}

//通过euid 获取Otc Order
func (c *OtcExchangerController) GetOtcOrder() {
	c.setOPAction(OPActionAddExchangerRestore)
	c.setRequestData(string(c.Ctx.Input.RequestBody))
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddExchangerRestore, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	euid, err := c.GetUint64(Key_EUID, 0)
	page, err := c.GetInt(KEY_PAGE, 1)
	limit, err := c.GetInt(KEY_LIMIT, 10)

	if err != nil && euid > 0 {
		c.ErrorResponseAndLog(OPActionGetOtcOrder, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	orderInfoList, err := otcDao.OrdersDaoEntity.GetOtcOrder(euid, page, limit)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
	}
	res := map[string]interface{}{}
	if len(orderInfoList) > 0 {
		res["list"] = orderInfoList
	} else {
		res["list"] = [...][]*models.OtcOrder{}

	}
	c.SuccessResponse(res)

}

//easyjson:json
type GetExchangerWealthAck struct {
	Available string `json:"available"`
	Trade     string `json:"trade"`
	Transfer  string `json:"transfer"`
}

// 通过euid 获取承兑商资产
func (c *OtcExchangerController) GetExchangerWealth() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionGetOtcWealth, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	euid, err := c.GetUint64(Key_EUID, 0)
	if err != nil || euid == 0 {
		c.ErrorResponseAndLog(OPActionGetOtcWealth, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"euid\":%v}", euid)
	exchangerWealth, err := eusdDao.EosOtcDaoEntity.GetOtcExchangerWealth(euid)
	if err != nil {
		c.ErrorResponseAndLog(OPActionGetOtcWealth, controllers.ERROR_CODE_DB, input)
		return
	}

	c.SuccessResponseAndLog(OPActionGetOtcWealth, input, exchangerWealth)
}
