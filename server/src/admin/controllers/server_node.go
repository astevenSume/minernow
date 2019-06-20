package controllers

import (
	. "admin/controllers/errcode"
	"fmt"
	admindao "utils/admin/dao"
)

type ServerNodeController struct {
	BaseController
}

func (c *ServerNodeController) Query() {
	_, errCode := c.CheckPermission()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionServerNodeGet, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	appName := c.GetString(KeyAppName)

	regionId, _ := c.GetInt64(KeyRegionId, -1)
	serverId, _ := c.GetInt64(KeyServerId, -1)

	page, _ := c.GetInt(KEY_PAGE)
	perPage, _ := c.GetInt(KEY_LIMIT)
	if page == 0 {
		page = 1
	}
	if perPage == 0 {
		perPage = DEFAULT_PER_PAGE
	}

	res := map[string]interface{}{}
	input := fmt.Sprintf("{\"page\":%v,\"per_page\":%v}", page, perPage)

	total, list, err := admindao.ServerNodeDaoEntity.Query(appName, regionId, serverId, page, perPage)
	if err != nil {
		c.ErrorResponseAndLog(OPActionServerNodeGet, ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	meta := admindao.PageInfo{
		Limit: perPage,
		Total: int(total),
		Page:  page,
	}
	res["meta"] = meta
	res["list"] = list
	c.SuccessResponseAndLog(OPActionReadPrice, input, res)
}
