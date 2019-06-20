package controllers

import (
	"admin/controllers/errcode"
	"fmt"
	"utils/admin/dao"
	otcDao "utils/otc/dao"
)

type OtcUserLoginLogController struct {
	BaseController
}

//获取otc用户登录日志
func (c *OtcUserLoginLogController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadUserLoginLog, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUserLoginLog, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUserLoginLog, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	mobile := c.GetString(KEY_MOBILE)
	otcUid, errCode := GetOtcUidByMobile(mobile)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadOrder, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	//分页查询
	res := map[string]interface{}{}
	input := fmt.Sprintf("{\"mobile\":\"%v\",\"page\":%v,\"limit\":%v}", mobile, page, perPage)
	count, data, err := otcDao.UserLoginLogDaoEntity.QueryByPage(otcUid, page, perPage)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUserLoginLog, controllers.ERROR_CODE_DB, input)
		return
	}

	meta := dao.PageInfo{
		Limit: perPage,
		Total: int(count),
		Page:  page,
	}
	res["meta"] = meta
	res["list"] = data
	c.SuccessResponseAndLog(OPActionReadUserLoginLog, input, res)
}
