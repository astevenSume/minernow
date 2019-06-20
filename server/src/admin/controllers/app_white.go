package controllers

import (
	controllers "admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
)

type AppWhiteController struct {
	BaseController
}

//获取应用白名单
func (c *AppWhiteController) GetWhiteApp() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadAppWhite, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAppWhite, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := dao.AppWhiteListDaoEntity.QueryById(id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppWhite, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadAppWhite, input, *data)
	} else {
		channelId, err := c.GetInt8("channel_id", -1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppWhite, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		page, err := c.GetInt(KEY_PAGE, 1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppWhite, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppWhite, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		//分页查询
		res := map[string]interface{}{}
		input = fmt.Sprintf("{\"channelId\":%v,\"page\":%v,\"limit\":%v}",
			channelId, page, perPage)
		count, data, err := dao.AppWhiteListDaoEntity.QueryPageCondition(channelId, page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppWhite, controllers.ERROR_CODE_DB, input)
			return
		}
		meta := PageInfo{
			Limit: perPage,
			Total: int(count),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadAppWhite, input, res)
	}
}

//easyjson:json
type CreateWhiteAppReq struct {
	ChannelId uint32 `json:"channel_id"`
	AppId     string `json:"app_id"`
}

//创建
func (c *AppWhiteController) CreateWhiteApp() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAppWhite, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &CreateWhiteAppReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddAppWhite, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	app, err := dao.AppWhiteListDaoEntity.Create(req.ChannelId, req.AppId)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionAddAppWhite, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionAddAppWhite, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionAddAppWhite, string(c.Ctx.Input.RequestBody), app)
}

//删除
func (c *AppWhiteController) DelWhiteApp() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelAppWhite, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAppWhite, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.AppWhiteListDaoEntity.DelById(id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAppWhite, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelAppWhite, input, ack)
}
