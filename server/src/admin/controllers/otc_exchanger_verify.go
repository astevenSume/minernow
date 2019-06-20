package controllers

import (
	"admin/controllers/errcode"
	"common"
	"fmt"
	"umeng_push/uemng_plus"
	"utils/admin/dao"
	otcDao "utils/otc/dao"
)

type OtcExchangerVerifyController struct {
	BaseController
}

//获取承兑商审核列表
func (c *OtcExchangerVerifyController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadExchangersVerify, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadExchangersVerify, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := otcDao.OtcExchangerVerifyDaoEntity.QueryById(id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadExchangersVerify, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadExchangersVerify, input, otcDao.OtcExchangerVerifyDaoEntity.ClientExchangerVerify(data))
	} else {
		page, err := c.GetInt("page", 1)
		if err != nil {
			common.LogFuncError("err:%v", err)
			c.ErrorResponseAndLog(OPActionReadExchangersVerify, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		}
		perPage, err := c.GetInt("per_page", DEFAULT_PER_PAGE)
		if err != nil {
			common.LogFuncError("err:%v", err)
			c.ErrorResponseAndLog(OPActionReadExchangersVerify, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		}
		status, err := c.GetInt("status", -1)
		if err != nil {
			common.LogFuncError("err:%v", err)
			c.ErrorResponseAndLog(OPActionReadExchangersVerify, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		}
		mobile := c.GetString("mobile")
		wechat := c.GetString("wechat")

		input = fmt.Sprintf("{\"mobile\":\"%s\",\"wechat\":\"%s\",\"status\":%v,\"page\":%v,\"per_page\":%v}", mobile, wechat, status, page, perPage)
		total, data, err := otcDao.OtcExchangerVerifyDaoEntity.QueryCondition(mobile, wechat, status, page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadExchangersVerify, controllers.ERROR_CODE_DB, input)
			return
		}

		res := map[string]interface{}{}
		meta := dao.PageInfo{
			Limit: perPage,
			Total: total,
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = otcDao.OtcExchangerVerifyDaoEntity.ClientExchangerVerifys(data)
		c.SuccessResponseAndLog(OPActionReadExchangersVerify, input, res)
	}
}

//审核拒绝
func (c *OtcExchangerVerifyController) Reject() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadExchangersVerifyReject, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadExchangersVerifyReject, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)

	_, err = otcDao.OtcExchangerVerifyDaoEntity.UpdateStatus(id, otcDao.ExchangerVerifyStatusReject)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadExchangersVerifyReject, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	data := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionReadExchangersVerifyReject, input, data)
}

//审核通过承兑商
func (c *OtcExchangerVerifyController) Approve() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadExchangersVerifyApprove, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadExchangersVerifyApprove, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)

	otcVerify, err := otcDao.OtcExchangerVerifyDaoEntity.UpdateStatus(id, otcDao.ExchangerVerifyStatusActive)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadExchangersVerifyApprove, controllers.ERROR_CODE_DB, input)
		return
	}

	_, err = BecomeExchanger(otcVerify.Uid, otcVerify.From, otcVerify.Mobile, otcVerify.Wechat, otcVerify.Telegram)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadExchangersVerifyApprove, controllers.ERROR_CODE_DB, input)
		return
	}
	go common.SafeRun(func() {
		title := "恭喜您成为承兑商!"
		content := "经系统审核，您提交的成为承兑商用户，已通过。"
		p := new(uemng_plus.UPushPlus)
		p.PushSysNotice(otcVerify.Uid, content, title)
		_, _ = otcDao.SystemNotificationdDaoEntity.InsertSystemNotification("system", content, otcVerify.Uid)
	})()

	//返回查询结果
	data := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionReadExchangersVerifyApprove, input, data)
}
