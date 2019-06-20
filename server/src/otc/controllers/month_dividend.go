package controllers

import (
	"common"
	"encoding/json"
	"eusd/eosplus"
	"github.com/bsm/redis-lock"
	common3 "otc/common"
	. "otc_error"
	"time"
	dao2 "utils/agent/dao"
	dao3 "utils/otc/dao"
	models2 "utils/otc/models"
	"utils/report/dao"
	"utils/report/models"
	usdtDao "utils/usdt/dao"
)

type MonthDividendController struct {
	BaseController
}

func (c MonthDividendController) NeedRechargeAmount() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	agent, err := dao2.AgentDaoEntity.Info(uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	amount := eosplus.QuantityInt64ToFloat64(agent.SumCanWithdraw)
	if amount >= 0 {
		amount = 0
	} else {
		amount = -amount
	}
	c.SuccessResponse(map[string]interface{}{
		KeyAmount: amount,
	})
}

//充值接口
//充值时需要重置自己的ReceiveStatus和所有下级代理的 ReceiveStatus
//自己的ReceiveStatus：
//1.查看自己上级代理（如果不是一级代理的话）的 ReceiveStatus 是不是 NeedRecharge 状态，如果不是则把ReceiveStatus设置成 NoReceived(可以领取)
//如果是 NeedRecharge 状态则把自己的状态设置成 CanNotReceive(等待上级分红)
//2.设置自己所有下级代理的状态,
//2.1 查看自己所有下级代理的所有月分红数据查看她的  ReceiveStatus 是不是 NeedRecharge 状态，如果是则不改动
//2.2 如果是 CanNotReceive(等待上级分红) 状态 则设置成 NoReceived(可以领取) 状态
type RechargeMsg struct {
	Amount float64 `json:"amount"`
}

func (c *MonthDividendController) Recharge() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := RechargeMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	agent, err := dao2.AgentDaoEntity.Info(uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	needRechargeAmount := eosplus.QuantityInt64ToFloat64(agent.SumCanWithdraw)
	needRechargeAmountTmp := needRechargeAmount
	if needRechargeAmount < 0 {
		needRechargeAmountTmp = -needRechargeAmount
	}
	if msg.Amount < needRechargeAmountTmp {
		c.ErrorResponse(ERROR_CODE_RECHARGE_AMOUNT_NO_RIGHT)
		return
	}
	//todo:从eusd转到分红表
	//如果之前就处于不欠费的状态则不需要解冻和分发给下级
	if needRechargeAmount >= 0 {
		c.SuccessResponseWithoutData()
		return
	}
	errCode = UnFrozenUser(uid)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	err = DividendToLowAgents(uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DIVIDEND_TO_AGENT_ERR)
		return
	}
	c.SuccessResponseWithoutData()
}

//1.将下级代理的selfDividend加到 SumCanWithdraw ;2.将下级代理的 ReceivedStatus 设置成Received
func DividendToLowAgents(parentUid uint64) (err error) {
	//解冻自己的账号
	//将自己未支付的订单全都设置成已支付
	err = dao.MonthDividendRecordDaoEntity.UpdatePayStatus(parentUid, dao.Pay, dao.NoPay)
	if err != nil {
		return
	}
	//	agentPath找到他的直属下级
	lowAgents, err := dao2.AgentPathDaoEntity.GetSubLevelByUid(parentUid)
	uids := make([]uint64, 0, len(lowAgents))
	for _, lowAgent := range lowAgents {
		uids = append(uids, lowAgent.Uid)
	}
	// 找到直属下级的月分红记录
	records, err := dao.MonthDividendRecordDaoEntity.FindByUidsAndStatus(uids, dao.NoReceived)
	if err != nil {
		return
	}
	recordMap := make(map[uint64][]*models.MonthDividendRecord, 0)
	for _, record := range records {
		if _, ok := recordMap[record.Uid]; !ok {
			recordMap[record.Uid] = make([]*models.MonthDividendRecord, 0)
		}
		recordMap[record.Uid] = append(recordMap[record.Uid], record)
	}
	totalDividendAmount := make(map[uint64]int64, 0)
	//更新月分红记录的状态成已发放和统计所有需要发放的月分红
	for _, uid := range uids {
		var dividendAmount int64 = 0
		for _, record := range recordMap[uid] {
			dividendAmount += record.SelfDividend
			//record.ReceiveStatus = dao.Received
		}
		totalDividendAmount[uid] = dividendAmount
	}

	//更新月分红记录的状态成已发放
	for _, uid := range uids {
		err = dao.MonthDividendRecordDaoEntity.UpdateReceiveStatus(uid, dao.Received, dao.NoReceived)
	}

	//增加agent表中SumCanWithdraw的值
	for uid, amount := range totalDividendAmount {
		_, err = dao2.AgentDaoEntity.UpdateCanWithdraw(uid, amount, lock.Options{
			LockTimeout: time.Second * time.Duration(common3.Cursvr.LockTimeout),
			RetryCount:  common3.Cursvr.RetryCount,
			RetryDelay:  time.Duration(common3.Cursvr.RetryDelay) * time.Millisecond,
		})
		if err != nil {
			return
		}
	}
	return
}

func UnFrozenUser(uid uint64) (errCode ERROR_CODE) {
	// 解冻EUSD账号
	errCode = eosplus.EosPlusAPI.Wealth.Unlock(uid)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}

	// 解冻承兑OTC账号
	errCode = eosplus.EosPlusAPI.Otc.Unlock(uid)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}

	// 解冻USDT账号
	err := usdtDao.AccountDaoEntity.UpdateStatus(uid, usdtDao.STATUS_WORKING)
	if err != nil {
		errCode = ERROR_CODE_USDT_CLEAR_FROZEN_FAILED
		return
	}
	errCode = ERROR_CODE_SUCCESS
	return
}

const complete = 1       //上级已经分红并且自己也分给下级分红了(上级的 SumCanWithdraw 大等于0 并且自己的 SumCanWithdraw 大等于0)
const waitBeDividend = 2 // 等待上级分红,未收到上级的分红，但是自己已经发给下级分红了(上级的 SumCanWithdraw 小于0 并且自己的 SumCanWithdraw 大等于0)
const waitRecharge = 3   //收到上级的分红，但是自己未发放分红给下级(上级的 SumCanWithdraw 大等于0 并且自己的 SumCanWithdraw 小于0)
const allWait = 4        //未收到上级分红，并且自己也未支付给下级(上级的 SumCanWithdraw 小于0 并且自己的 SumCanWithdraw 小于0)
type MonthDividendList struct {
	RecordId       uint64  `json:"record_id"`
	ResultDividend float64 `json:"result_dividend"`
	CTime          string  `json:"c_time"`
	Status         int     `json:"status"`
}

func (c *MonthDividendController) List() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	//找到自己的月分红记录,得到 resultDividend和ctime
	records, err := dao.MonthDividendRecordDaoEntity.FindByUid(uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	recordList := make([]*MonthDividendList, 0, len(records))
	for _, record := range records {
		r := new(MonthDividendList)
		r.RecordId = record.Id
		r.ResultDividend = eosplus.QuantityInt64ToFloat64(record.ResultDividend)

		if record.ReceiveStatus == dao.Received && record.PayStatus == dao.Pay {
			r.Status = complete
		} else if record.ReceiveStatus == dao.NoReceived && record.PayStatus == dao.Pay {
			r.Status = waitBeDividend
		} else if record.ReceiveStatus == dao.Received && record.PayStatus == dao.NoPay {
			r.Status = waitRecharge
		} else if record.ReceiveStatus == dao.NoReceived && record.PayStatus == dao.NoPay {
			r.Status = allWait
		}

		tm := time.Unix(record.Ctime, 0)
		r.CTime = tm.Format("2006-01-02 15:04:05")
		recordList = append(recordList, r)
	}
	res := map[string]interface{}{
		"list": recordList,
	}
	c.SuccessResponse(res)
}

type detail struct {
	Mobile   string  `json:"mobile"`
	Dividend float64 `json:"dividend"`
}

func (c *MonthDividendController) Details() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	//自己月分红的收入和应该分给下级的月分红
	timestamp, err := c.GetInt64(KeyTimestamp)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}
	selfRecord, err := dao.MonthDividendRecordDaoEntity.GetByUidAndCtime(uid, timestamp)
	if err != nil {
		if selfRecord == nil {
			c.ErrorResponse(ERROR_CODE_NO_FOUND_MONTH_DIVIDEND_RECORD)
		}
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	//	agentPath找到他的直属下级
	lowAgents, err := dao2.AgentPathDaoEntity.GetSubLevelByUid(uid)
	lowUids := make([]uint64, 0, len(lowAgents))
	for _, lowAgent := range lowAgents {
		lowUids = append(lowUids, lowAgent.Uid)
	}
	//找到 代理id和ctime是selfRecord.Ctime的数据
	lowReports, err := dao.MonthDividendRecordDaoEntity.FindByUidsAndCtime(lowUids, selfRecord.Ctime)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	//查找下级的电话
	usersInfo, err := dao3.UserDaoEntity.FindByUids(lowUids)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	usersInfoMap := make(map[uint64]*models2.User, 0)

	for _, userInfo := range usersInfo {
		usersInfoMap[userInfo.Uid] = userInfo
	}
	detailList := make([]*detail, 0, len(lowAgents))

	//所有分给下级代理的月分红,下级代理电话+分到的月分红
	for _, lowReport := range lowReports {
		detailTmp := new(detail)
		detailTmp.Mobile = usersInfoMap[lowReport.Uid].Mobile
		detailTmp.Dividend = eosplus.QuantityInt64ToFloat64(lowReport.SelfDividend)
		detailList = append(detailList, detailTmp)
	}

	res := map[string]interface{}{
		"dividend":   eosplus.QuantityInt64ToFloat64(selfRecord.SelfDividend),
		"pay":        eosplus.QuantityInt64ToFloat64(selfRecord.AgentDividend),
		"pay_detail": detailList,
	}
	c.SuccessResponse(res)
}
