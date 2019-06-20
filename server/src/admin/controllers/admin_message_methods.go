package controllers

import (
	. "admin/controllers/errcode"
	"fmt"
	"strconv"
	otcdao "utils/otc/dao"
)

type AdminMessageMethodController struct {
	BaseController
}

// admin 通过order_id 获取聊天内容
func (c *AdminMessageMethodController) GetMessage() {
	_, errCode := c.CheckPermission()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadOrderMsg, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOrderMsg, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	perPage, err := c.GetInt(KEY_LIMIT, 10)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOrderMsg, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var orderId int64
	orderIdStr := c.Ctx.Input.Param(":order_id")
	input := fmt.Sprintf("{\"page\":%v,\"limit\":%v,\"orderId\":\"%v\"}}", page, perPage, orderIdStr)
	if len(orderIdStr) > 0 {
		var err error
		orderId, err = strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadOrderMsg, ERROR_CODE_PARAM_FAILED, input)
			return
		}
	}

	msgList, total, err := otcdao.MessageMethodDaoEntity.QueryByOrderId(orderId, page, perPage)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOrderMsg, ERROR_CODE_DB, input)
		return
	}

	meta := &PageInfo{}
	meta.Total = int(total)
	meta.Limit = perPage
	meta.Page = page

	res := map[string]interface{}{}
	res["list"] = msgList
	res["meta"] = meta

	c.SuccessResponseAndLog(OPActionReadOrderMsg, input, res)
}
