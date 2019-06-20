package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"strconv"
	"utils/admin/dao"
	octDao "utils/otc/dao"
)

type AppealController struct {
	BaseController
}

//获取申述管理
func (c *AppealController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadAppeal, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint64(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAppeal, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := octDao.AppealDaoEntity.QueryById(id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppeal, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadAppeal, input, octDao.AppealDaoEntity.ClientAppeal(&data))
	} else {
		sType, err := c.GetInt(KEY_TYPE, 0)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppeal, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		status, err := c.GetInt(KEY_STATUS, 0)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppeal, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		page, err := c.GetInt(KEY_PAGE, 1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppeal, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppeal, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		mobile := c.GetString(KEY_MOBILE)
		otcUid, errCode := GetOtcUidByMobile(mobile)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			c.ErrorResponseAndLog(OPActionReadAppeal, errCode, string(c.Ctx.Input.RequestBody))
			return
		}

		//分页查询
		res := map[string]interface{}{}
		input = fmt.Sprintf("{\"mobile\":\"%v\",\"status\":%v,\"page\":%v,\"limit\":%v}", mobile, status, page, perPage)
		count, data, err := octDao.AppealDaoEntity.QueryPageCondition(otcUid, int8(sType), int8(status), page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppeal, controllers.ERROR_CODE_DB, input)
			return
		}

		meta := PageInfo{
			Limit: perPage,
			Total: int(count),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadAppeal, input, res)
	}
}

//更新
func (c *AppealController) Update() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditAppeal, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id      string `json:"id"`
		AdminId uint32 `json:"admin_id"`
		Status  int8   `json:"status"`
		Type    int8   `json:"type"`
		Context string `json:"context"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil || len(req.Id) <= 0 {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditAppeal, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditAppeal, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	_, err = octDao.AppealDaoEntity.Update(id, req.AdminId, req.Type, req.Status, req.Context)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionEditAppeal, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionEditAppeal, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionEditAppeal, string(c.Ctx.Input.RequestBody))
}
