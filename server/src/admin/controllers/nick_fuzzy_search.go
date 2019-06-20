package controllers

import (
	"admin/controllers/errcode"
	"common"
	"utils/otc/dao"
)

type NickFuzzySearchController struct {
	BaseController
}

func (c *NickFuzzySearchController) GetUsers() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	nick := c.GetString("nick", "")
	if nick == "" {
		common.LogFuncDebug("no nick : %v", string(c.Ctx.Input.RequestBody))
		c.ErrorResponseAndLog(OPActionGetUsersByNick, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	page, _ := c.GetInt(KEY_PAGE, 1)
	limit, _ := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	//根据用户名模糊查询有多少用户表有多少数据，获得total给前端做分页查询使用
	total, err := dao.UserDaoEntity.CountUserByNick(nick)
	if err != nil {
		c.ErrorResponseAndLog(OPActionGetUsersByNick, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	//根据用户名模糊分页查询用户表
	users, err := dao.UserDaoEntity.FindUserIdByNick(nick, page, limit)
	if err != nil {
		c.ErrorResponseAndLog(OPActionGetUsersByNick, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	res := map[string]interface{}{
		"total": total,
		"list":  users,
	}
	c.SuccessResponseAndLog(OPActionGetUsersByNick, string(c.Ctx.Input.RequestBody), res)
}
