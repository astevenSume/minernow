package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"eusd/eosplus"
	"fmt"
	"github.com/astaxie/beego/orm"
	lock "github.com/bsm/redis-lock"
	"time"

	common3 "otc/common"
	"umeng_push/uemng_plus"
	"utils/admin/dao"
	dao2 "utils/agent/dao"
	"utils/agent/models"
	dao3 "utils/otc/dao"
	models2 "utils/otc/models"
	dao4 "utils/report/dao"
)

type ProfitThresholdController struct {
	BaseController
}

func (c ProfitThresholdController) Add() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}
	type msg struct {
		Threshold float64 `json:"threshold"`
		AdminId   uint64  `json:"admin_id"`
	}
	reqMsg := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqMsg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddProfitThreshold, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//删除之前的所有配置
	err = dao.ProfitThresholdDaoEntity.DeleteAll()
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddProfitThreshold, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	threshold := eosplus.QuantityFloat64ToInt64(reqMsg.Threshold)
	_, err = dao.ProfitThresholdDaoEntity.Add(threshold, reqMsg.AdminId)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddProfitThreshold, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	c.SuccessResponseWithoutData()
}

func (c ProfitThresholdController) Del() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}
	id, err := c.GetUint32(KEY_ID)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionDelProfitThreshold, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//删除之前的所有配置
	err = dao.ProfitThresholdDaoEntity.Delete(id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelProfitThreshold, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	c.SuccessResponseWithoutData()
}

func (c ProfitThresholdController) Update() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}
	type msg struct {
		Id        uint32  `json:"id"`
		Threshold float64 `json:"threshold"`
		AdminId   uint64  `json:"admin_id"`
	}
	reqMsg := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqMsg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionUpdateProfitThreshold, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	threshold := eosplus.QuantityFloat64ToInt64(reqMsg.Threshold)
	_, err = dao.ProfitThresholdDaoEntity.Update(reqMsg.Id, threshold, reqMsg.AdminId)
	if err != nil {
		c.ErrorResponseAndLog(OPActionUpdateProfitThreshold, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	c.SuccessResponseWithoutData()
}

type ProfitThresholdToClient struct {
	Id        uint32  `json:"id,omitempty"`
	Threshold float64 `json:"threshold,omitempty"`
	AdminId   uint64  `json:"admin_id,omitempty"`
	Ctime     int64   `json:"ctime,omitempty"`
	Utime     int64   `json:"utime,omitempty"`
}

func (c ProfitThresholdController) Find() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	datas, err := dao.ProfitThresholdDaoEntity.Find()
	if err != nil {
		c.ErrorResponseAndLog(OPActionFindProfitThreshold, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	thresholdDatas := make([]*ProfitThresholdToClient, 0, len(datas))
	for _, data := range datas {
		r := new(ProfitThresholdToClient)
		r.Id = data.Id
		r.AdminId = data.AdminId
		r.Threshold = eosplus.QuantityInt64ToFloat64(data.Threshold)
		r.Utime = data.Utime
		r.Ctime = data.Ctime
		thresholdDatas = append(thresholdDatas, r)
	}

	req := map[string]interface{}{}
	req["list"] = thresholdDatas
	req[KEY_META] = map[string]interface{}{
		KEY_LIMIT: 100,
		KEY_PAGE:  1,
		"total":   len(thresholdDatas),
	}
	c.SuccessResponseAndLog(OPActionFindProfitThreshold, string(c.Ctx.Input.RequestBody), req)
}

//获取待审核记录
func (c *ProfitThresholdController) FindCheckingRecords() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}
	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelMenu, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	limit, err := c.GetInt(KEY_LIMIT, 10)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelMenu, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	total, err := dao2.AgentWithdrawDaoEntity.QueryTotalWithStatus(dao2.AgentWithdrawStatusChecking)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddMenu, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	withdrawList, err := dao2.AgentWithdrawDaoEntity.FindByStatus(dao2.AgentWithdrawStatusChecking, page, limit)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddMenu, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	uids := make([]uint64, 0, len(withdrawList))
	for _, withdraw := range withdrawList {
		uids = append(uids, withdraw.Uid)
	}
	users, err := dao3.UserDaoEntity.FindByUids(uids)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddMenu, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	userMap := make(map[uint64]*models2.User, 0)
	for _, user := range users {
		userMap[user.Uid] = user
	}
	type agentWithdraw struct {
		Nick   string `json:"nick"`
		Mobile string `json:"mobile"`
		*models.AgentWithdraw
	}
	agentWithdrawList := make([]*agentWithdraw, 0, len(withdrawList))
	for _, withdraw := range withdrawList {
		agentWith := new(agentWithdraw)
		agentWith.AgentWithdraw = withdraw
		if user, ok := userMap[agentWith.Uid]; ok {
			agentWith.Nick = user.Nick
			agentWith.Mobile = user.Mobile
		}
		agentWithdrawList = append(agentWithdrawList, agentWith)
	}
	res := map[string]interface{}{
		KEY_LIST: agentWithdrawList,
		KEY_META: map[string]interface{}{
			KEY_PAGE:  page,
			KEY_LIMIT: limit,
			KEY_TOTAL: total,
		},
	}
	c.SuccessResponseAndLog(OPActionAddMenu, string(c.Ctx.Input.RequestBody), res)
}

//拒绝待审核记录
func (c *ProfitThresholdController) RejectCheckingRecord() {
	type reqMsg struct {
		RecordId uint64 `json:"record_id"`
		AdminId  uint64 `json:"admin_id"`
	}

	req := &reqMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddProfitThreshold, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	fmt.Println("req.RecordId ", req.RecordId)
	record, err := dao2.AgentWithdrawDaoEntity.Get(req.RecordId)
	if err != nil {
		if err == orm.ErrNoRows {
			c.ErrorResponseAndLog(OPActionAddMenu, controllers.ERROR_CODE_DB_NO_DATA, string(c.Ctx.Input.RequestBody))
			return
		}
		c.ErrorResponseAndLog(OPActionAddMenu, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	if record.Status != dao2.AgentWithdrawStatusChecking {
		c.ErrorResponseAndLog(OPActionAddMenu, controllers.ERROR_CODE_AGENT_WITHDRAW_CHECKED, string(c.Ctx.Input.RequestBody))
		return
	}
	desc := fmt.Sprintf("reject by adminId %v", req.AdminId)
	err = dao2.AgentWithdrawDaoEntity.UpdateStatus(req.RecordId, dao2.AgentWithdrawStatusFailed, desc)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddMenu, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	//推送通知
	go func() {
		title := "分润提现审核失败"
		content := "分润提现审核失败，如有疑问，请联系客服。"
		p := new(uemng_plus.UPushPlus)
		p.PushSysNotice(record.Uid, content, title)
		_, _ = dao3.SystemNotificationdDaoEntity.InsertSystemNotification("system", content, record.Uid)
	}()
	c.SuccessResponseWithoutData()
}

//同意待审核记录
func (c *ProfitThresholdController) AgreeCheckingRecord() {
	type reqMsg struct {
		RecordId uint64 `json:"record_id"`
		AdminId  uint64 `json:"admin_id"`
	}
	req := &reqMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddProfitThreshold, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	record, err := dao2.AgentWithdrawDaoEntity.Get(req.RecordId)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddMenu, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	if record.Status != dao2.AgentWithdrawStatusChecking {
		c.ErrorResponseAndLog(OPActionAddMenu, controllers.ERROR_CODE_AGENT_WITHDRAW_CHECKED, string(c.Ctx.Input.RequestBody))
		return
	}

	// decrease commission
	err = dao2.AgentDaoEntity.DecreaseCommission(record.Uid, record.Amount, lock.Options{
		LockTimeout: time.Second * time.Duration(common3.Cursvr.LockTimeout),
		RetryCount:  common3.Cursvr.RetryCount,
		RetryDelay:  time.Duration(common3.Cursvr.RetryDelay) * time.Millisecond,
	})
	if err != nil {
		_ = dao2.AgentWithdrawDaoEntity.UpdateStatus(record.Id, dao2.AgentWithdrawStatusFailed, fmt.Sprintf("%v while DecreaseCommission", err))
		if err == dao2.ErrCurrencyNoEnough {
			c.ErrorResponseAndLog(OPActionAddProfitThreshold, controllers.ERROR_CODE_AGENT_WITHDRAW_NO_ENOUGH, string(c.Ctx.Input.RequestBody))
			return
		}

		c.ErrorResponseAndLog(OPActionAddProfitThreshold, controllers.ERROR_CODE_AGENT_DECREASE_COMMISSION_FAILED, string(c.Ctx.Input.RequestBody))
		return
	}

	// update withdraw status
	err = dao2.AgentWithdrawDaoEntity.UpdateStatus(record.Id, dao2.AgentWithdrawStatusDecreased, "")
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddProfitThreshold, controllers.ERROR_CODE_AGENT_WITHDRAW_UPDATE_RECORD_STATUS_FAILED, string(c.Ctx.Input.RequestBody))
		return
	}

	// send eusd
	errCode := eosplus.EosPlusAPI.Wealth.Commission(record.Uid, record.Amount)
	if controllers.ERROR_CODE(errCode) != controllers.ERROR_CODE_SUCCESS {
		_ = dao2.AgentWithdrawDaoEntity.UpdateStatus(record.Id, dao2.AgentWithdrawStatusFailed, fmt.Sprintf("error code %d while send eusd", errCode))
		c.ErrorResponseAndLog(OPActionAddProfitThreshold, controllers.ERROR_CODE(errCode), string(c.Ctx.Input.RequestBody))
		return
	}

	// update withdraw status
	err = dao2.AgentWithdrawDaoEntity.UpdateStatus(record.Id, dao2.AgentWithdrawStatusDone, "")
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddProfitThreshold, controllers.ERROR_CODE_AGENT_WITHDRAW_UPDATE_RECORD_STATUS_FAILED, string(c.Ctx.Input.RequestBody))
		return
	}
	//日报
	_ = dao4.ReportAgentDailyDaoEntity.AddWithDraw(record.Uid, record.Amount)
	//推送通知
	go func() {
		title := "分润提现审核成功"
		content := "分润提现审核成功，如有疑问，请联系客服。"
		p := new(uemng_plus.UPushPlus)
		p.PushSysNotice(record.Uid, content, title)
		_, _ = dao3.SystemNotificationdDaoEntity.InsertSystemNotification("system", content, record.Uid)
	}()
	c.SuccessResponseWithoutData()
}
