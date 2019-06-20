package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"time"
	"utils/admin/dao"
	otcdao "utils/otc/dao"
)

type CommissionStatController struct {
	BaseController
}

//获取佣金统计
func (c *CommissionStatController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadCommissionStat, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	status, err := c.GetInt8(KEY_STATUS, -1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadCommissionStat, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadCommissionStat, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	limit, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadCommissionStat, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"status\":%v,\"page\":%v,\"limit\":%v}", status, page, limit)
	total, data, err := otcdao.CommissionStatDaoEntity.QueryByPage(status, page, limit)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadCommissionStat, controllers.ERROR_CODE_DB, input)
		return
	}

	res := map[string]interface{}{}
	meta := dao.PageInfo{
		Limit: limit,
		Total: int(total),
		Page:  page,
	}
	res["meta"] = meta
	res["list"] = ClientCommissionStat(data)
	c.SuccessResponseAndLog(OPActionReadCommissionStat, input, res)
}

func (c *CommissionStatController) Distribute() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddCommissionDistribute, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Ctime int64 `json:"ctime"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddCommissionDistribute, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//判断是否已发放
	isDistributed, err := otcdao.CommissionStatDaoEntity.IsDistributed(req.Ctime)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddCommissionDistribute, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	if isDistributed {
		c.ErrorResponseAndLog(OPActionAddCommissionDistribute, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//otc请求包
	type otcReq struct {
		Time      int64  `json:"time"`
		Timestamp uint32 `json:"timestamp"`
	}
	reqOtc := &otcReq{
		Time:      req.Ctime,
		Timestamp: uint32(time.Now().Unix()),
	}

	//otc响应数据
	type otcAck struct {
		Code controllers.ERROR_CODE `json:"code"`
	}
	ack := otcAck{}

	//otc post请求
	errCode = PostOtc(RouterDistributeCommission, map[string]string{
		"time": fmt.Sprint(reqOtc.Time),
	}, reqOtc, &ack, reqOtc.Timestamp)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddCommissionDistribute, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if ack.Code != controllers.ERROR_CODE_SUCCESS {
		common.LogFuncError("errCode:%v", ack.Code)
		//c.ErrorResponseAndLog(OPActionAddCommissionDistribute, controllers.ERROR_CODE_COMMISSION_DISTRIBUTE_FAIL, uid, string(c.Ctx.Input.RequestBody))
		c.ResponseAndLog(OPActionAddCommissionDistribute, string(c.Ctx.Input.RequestBody), ack)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddCommissionDistribute, string(c.Ctx.Input.RequestBody))
}

func (c *CommissionStatController) CalcCommission() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddCommissionCalc, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Ctime int64 `json:"ctime"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddCommissionCalc, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//判断是否已发放
	isDistributed, err := otcdao.CommissionStatDaoEntity.IsDistributed(req.Ctime)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddCommissionCalc, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	if isDistributed {
		c.ErrorResponseAndLog(OPActionAddCommissionCalc, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//otc请求包
	type otcReq struct {
		Time      int64  `json:"time"`
		Timestamp uint32 `json:"timestamp"`
	}
	reqOtc := &otcReq{
		Time:      req.Ctime,
		Timestamp: uint32(time.Now().Unix()),
	}

	//otc响应数据
	type otcAck struct {
		Code controllers.ERROR_CODE `json:"code"`
	}
	ack := otcAck{}

	//otc post请求
	errCode = PostOtc(RouterCalcCommission, map[string]string{
		"time": fmt.Sprint(reqOtc.Time),
	}, reqOtc, &ack, reqOtc.Timestamp)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddCommissionCalc, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if ack.Code != controllers.ERROR_CODE_SUCCESS {
		common.LogFuncError("errCode:%v", ack.Code)
		c.ResponseAndLog(OPActionAddCommissionCalc, string(c.Ctx.Input.RequestBody), ack)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddCommissionCalc, string(c.Ctx.Input.RequestBody))
}
