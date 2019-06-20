package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
)

type TopAgentController struct {
	BaseController
}

//获取一级代理
func (c *TopAgentController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadTopAgent, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadTopAgent, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	var input string
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := dao.TopAgentDaoEntity.QueryById(id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadTopAgent, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadTopAgent, input, data)
	} else {
		status, err := c.GetInt8(KEY_STATUS, -1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadTopAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		mobile := c.GetString(KEY_MOBILE)

		page, err := c.GetInt(KEY_PAGE, 1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadTopAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		limit, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadTopAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		input = fmt.Sprintf("{\"status\":%v,\"mobile\":\"%v\",\"page\":%v,\"limit\":%v}", status, mobile, page, limit)

		total, data, err := dao.TopAgentDaoEntity.QueryPageCondition(mobile, status, page, limit)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadTopAgent, controllers.ERROR_CODE_DB, input)
			return
		}

		res := map[string]interface{}{}
		meta := dao.PageInfo{
			Limit: limit,
			Total: int(total),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadTopAgent, input, res)
	}
}

func (c *TopAgentController) Create() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddTopAgent, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		NationalCode string `json:"national_code"`
		Mobile       string `json:"mobile"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddTopAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	isNew, data, err := dao.TopAgentDaoEntity.Create(req.NationalCode, req.Mobile)
	if !isNew {
		c.ErrorResponseAndLog(OPActionAddTopAgent, controllers.ERROR_CODE_MOBILE_EXSIT, string(c.Ctx.Input.RequestBody))
		return
	}
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddTopAgent, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionAddTopAgent, string(c.Ctx.Input.RequestBody), data)
}

//更新
func (c *TopAgentController) Update() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditTopAgent, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id           uint32 `json:"id"`
		NationalCode string `json:"national_code"`
		Mobile       string `json:"mobile"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditTopAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	ok, data, err := dao.TopAgentDaoEntity.Update(req.Id, req.NationalCode, req.Mobile)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditTopAgent, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	if !ok {
		c.ErrorResponseAndLog(OPActionEditTopAgent, controllers.ERROR_CODE_MOBILE_REGISTERED, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionEditTopAgent, string(c.Ctx.Input.RequestBody), data)
}

//删除
func (c *TopAgentController) Del() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelTopAgent, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionDelTopAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	ok, err := dao.TopAgentDaoEntity.DelById(id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelTopAgent, controllers.ERROR_CODE_DB, input)
		return
	}
	if !ok {
		c.ErrorResponseAndLog(OPActionDelTopAgent, controllers.ERROR_CODE_MOBILE_REGISTERED, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelTopAgent, input, ack)
}
