package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"utils/admin/dao"
	agentdao "utils/agent/dao"
)

type MonthDividendConfController struct {
	BaseController
}

//获取所有月分红等级配置
func (c *MonthDividendConfController) GetMonthDividendConf() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	res := map[string]interface{}{}
	var input string
	//数据获取
	data, err := dao.MonthDividendPositionConfDaoEntity.QueryPageDividendConfs()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadMonthDividendCfg, controllers.ERROR_CODE_DB, input)
		return
	}
	res["list"] = data

	c.SuccessResponseAndLog(OPActionReadMonthDividendCfg, input, res)
}

//删除指定等级的月分红配置，被删除的等级只能是最后一级
func (c *MonthDividendConfController) DelMonthDividendConf() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	level, err := c.GetInt("level", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	//获得更低等级的条数
	row, err := dao.MonthDividendPositionConfDaoEntity.CountLowLevelConfNum(int32(level))
	if err != nil {
		c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	if row > 0 {
		c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, controllers.ERROR_CODE_HAS_LOWER_LEVEL, string(c.Ctx.Input.RequestBody))
		return
	}

	err = dao.MonthDividendPositionConfDaoEntity.DeleteByLevel(int32(level))
	if err != nil {
		c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	c.SuccessResponseWithoutData()
}

func (c *MonthDividendConfController) GetMonthDividendPosition() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	res := map[string]interface{}{}
	var input string
	//数据获取
	data, err := dao.MonthDividendPositionConfDaoEntity.QueryDividendPosition()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadMonthDividendPosition, controllers.ERROR_CODE_DB, input)
		return
	}
	res["list"] = data

	c.SuccessResponseAndLog(OPActionReadMonthDividendPosition, input, res)
}

//新增月分红等级配置
func (c *MonthDividendConfController) EditMonthDividendConf() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Data dao.MonthDividendCfgs `json:"data"`
	}
	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditMonthDividendCfg, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	//插入数据
	err = dao.MonthDividendPositionConfDaoEntity.EditMonthDividendConf(req.Data)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionEditMonthDividendCfg, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionEditMonthDividendCfg, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回c查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionEditMonthDividendCfg, string(c.Ctx.Input.RequestBody))
}

//easyjson:json
type AddMonthAgentMsg struct {
	Mobile string `json:"mobile"`
	DId    int32  `json:"position_id"`
}

//func (c *MonthDividendConfController) AddAgent() {
//	_, errCode := c.CheckPermission()
//	if errCode != controllers.ERROR_CODE_SUCCESS {
//		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, errCode, string(c.Ctx.Input.RequestBody))
//		return
//	}
//
//	req := &AddMonthAgentMsg{}
//	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
//	if err != nil {
//		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
//		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
//		return
//	}
//
//	otcUid, errCode := GetOtcUidByMobile(req.Mobile)
//	if errCode != controllers.ERROR_CODE_SUCCESS {
//		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, errCode, string(c.Ctx.Input.RequestBody))
//		return
//	}
//	err = agentdao.AgentPathDaoEntity.SetAgentDividendPosition(otcUid, req.DId)
//	if err != nil {
//		if dao.ErrParam == err {
//			c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_EXIST_WHITE_LIST_AGENT, string(c.Ctx.Input.RequestBody))
//		} else {
//			c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
//		}
//		return
//	}
//
//	//返回结果
//	data := map[string]interface{}{
//		"mobile":      req.Mobile,
//		"position_id": req.DId,
//	}
//	c.SuccessResponseAndLog(OPActionEditAgentDividendPosition, string(c.Ctx.Input.RequestBody), data)
//}

//func (c *MonthDividendConfController) DelAgent() {
//	_, errCode := c.CheckPermission()
//	if errCode != controllers.ERROR_CODE_SUCCESS {
//		c.ErrorResponseAndLog(OPActionDelAgentDividendPosition, errCode, string(c.Ctx.Input.RequestBody))
//		return
//	}
//
//	otcUid, err := c.GetUint64(KEY_UID, 0)
//	if err != nil || otcUid <= 0 {
//		c.ErrorResponseAndLog(OPActionDelAgentDividendPosition, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
//		return
//	}
//	input := fmt.Sprintf("{\"uid\":%v}", otcUid)
//
//	err = agentdao.AgentPathDaoEntity.SetAgentDividendPosition(otcUid, 0)
//	if err != nil {
//		c.ErrorResponseAndLog(OPActionDelAgentDividendPosition, controllers.ERROR_CODE_DB, input)
//		return
//	}
//
//	//返回查询结果
//	data := map[string]interface{}{
//		"uid": fmt.Sprintf("%v", otcUid),
//	}
//	c.SuccessResponseAndLog(OPActionDelAgentDividendPosition, input, data)
//}

func (c *MonthDividendConfController) GetAgents() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadAgentDividendPosition, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	level, err := c.GetInt("agent_level", 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteListAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	position, err := c.GetInt("position", 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteListAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteListAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	limit, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAgentWhiteListAgent, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	res := map[string]interface{}{}
	agentsPosition, err := agentdao.AgentPathDaoEntity.GetAgentByPosition(uint32(level), int32(position), page, limit)
	res["list"] = agentsPosition

	c.SuccessResponseAndLog(OPActionReadMonthDividendPosition, string(c.Ctx.Input.RequestBody), res)
}
