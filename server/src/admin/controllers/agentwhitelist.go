package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
	agentdao "utils/agent/dao"
)

type AgentWhiteListController struct {
	BaseController
}

func (c *AgentWhiteListController) GetAgentWhiteList() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteList, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	id, err := c.GetInt(KEY_ID, -1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteList, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := dao.AgentWhiteListDaoEntity.QueryById(uint32(id))
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAgentWhiteList, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadAgentWhiteList, input, *data)
	} else {
		page, err := c.GetInt("page", 1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAgentWhiteList, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		perPage, err := c.GetInt("per_page", DEFAULT_PER_PAGE)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAgentWhiteList, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		//分页查询
		res := map[string]interface{}{}
		input = fmt.Sprintf("{\"page\":%v,\"per_page\":%v}", page, perPage)
		count, data, err := dao.AgentWhiteListDaoEntity.QueryByPage(page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAgentWhiteList, controllers.ERROR_CODE_DB, input)
			return
		}
		meta := dao.PageInfo{
			Limit: perPage,
			Total: int(count),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadAgentWhiteList, input, res)
	}
}

func (c *AgentWhiteListController) CreateAgentWhiteList() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAgentWhiteList, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Name       string `json:"name"`
		Commission int32  `json:"commission"`
		Precision  int32  `json:"precision"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddAgentWhiteList, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	app, err := dao.AgentWhiteListDaoEntity.Create(req.Commission, req.Precision, req.Name)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddAgentWhiteList, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionAddAgentWhiteList, string(c.Ctx.Input.RequestBody), app)
}

func (c *AgentWhiteListController) UpdateAgentWhiteList() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditAgentWhiteList, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id         uint32 `json:"id"`
		Name       string `json:"name"`
		Commission int32  `json:"commission"`
		Precision  int32  `json:"precision"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditAgentWhiteList, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	app, err := dao.AgentWhiteListDaoEntity.UpdateById(req.Id, req.Commission, req.Precision, req.Name)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditAgentWhiteList, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionEditAgentWhiteList, string(c.Ctx.Input.RequestBody), app)
}

func (c *AgentWhiteListController) DelAgentWhiteList() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelAgentWhiteList, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAgentWhiteList, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.AgentWhiteListDaoEntity.DelById(id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAgentWhiteList, controllers.ERROR_CODE_DB, input)
		return
	}
	err = agentdao.AgentPathDaoEntity.DelWhiteList(id)
	if err != nil {
		common.LogFuncDebug("AgentPathDaoEntity.DelWhiteList fail error:%v", err)
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelAgentWhiteList, input, ack)
}

func (c *AgentWhiteListController) GetAgent() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteListAgent, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	pUid, err := c.GetUint64("puid", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteListAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	wId, err := c.GetInt("whitelist_id", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteListAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	inviteCode := c.GetString("invite_code")
	mobile := c.GetString(KEY_MOBILE)

	if pUid == 0 && len(mobile) > 0 {
		pUid, errCode = GetOtcUidByMobile(mobile)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			c.ErrorResponseAndLog(OPActionReadAgentWhiteListAgent, errCode, string(c.Ctx.Input.RequestBody))
			return
		}
	}

	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteListAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteListAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"mobile\":\"%v\",\"whitelist_id\":%v,\"inviteCode\":\"%v\",\"page\":%v,\"per_page\":%v}", mobile, wId, inviteCode, page, perPage)

	count, data, err := agentdao.AgentPathDaoEntity.GetAgentPathByPage(page, perPage, wId, pUid, inviteCode)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteListAgent, controllers.ERROR_CODE_DB, input)
		return
	}

	meta := dao.PageInfo{
		Limit: perPage,
		Total: int(count),
		Page:  page,
	}

	res := map[string]interface{}{}
	res["meta"] = meta
	res["list"] = ClientAgentInfos(data)
	c.SuccessResponseAndLog(OPActionReadAgentWhiteListAgent, input, res)
}

//easyjson:json
type AddAgentMsg struct {
	Mobile string `json:"mobile"`
	WId    uint32 `json:"whitelist_id"`
}

func (c *AgentWhiteListController) AddAgent() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditAgentWhiteListAgent, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &AddAgentMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditAgentWhiteListAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUid, errCode := GetOtcUidByMobile(req.Mobile)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditAgentWhiteListAgent, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	err = agentdao.AgentPathDaoEntity.SetAgentWhiteList(otcUid, req.WId)
	if err != nil {
		if dao.ErrParam == err {
			c.ErrorResponseAndLog(OPActionEditAgentWhiteListAgent, controllers.ERROR_CODE_EXIST_WHITE_LIST_AGENT, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionEditAgentWhiteListAgent, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	data := map[string]interface{}{
		"mobile":       req.Mobile,
		"whitelist_id": req.WId,
	}
	c.SuccessResponseAndLog(OPActionEditAgentWhiteListAgent, string(c.Ctx.Input.RequestBody), data)
}

func (c *AgentWhiteListController) DelAgent() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelAgentWhiteListAgent, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUid, err := c.GetUint64(KEY_UID, 0)
	if err != nil || otcUid <= 0 {
		c.ErrorResponseAndLog(OPActionDelAgentWhiteListAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"uid\":%v}", otcUid)

	err = agentdao.AgentPathDaoEntity.SetAgentWhiteList(otcUid, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAgentWhiteListAgent, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	data := map[string]interface{}{
		"uid": fmt.Sprintf("%v", otcUid),
	}
	c.SuccessResponseAndLog(OPActionDelAgentWhiteListAgent, input, data)
}
