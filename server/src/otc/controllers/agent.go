package controllers

import (
	"common"
	"eusd/eosplus"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bsm/redis-lock"
	json "github.com/mailru/easyjson"
	common3 "otc/common"
	. "otc_error"
	"time"
	agentdao "utils/agent/dao"
	common2 "utils/common"
	eusddao "utils/eusd/dao"
	gamedao "utils/game/dao"
	reportdao "utils/report/dao"
)

type AgentController struct {
	BaseController
}

const (
	AgentSalarySearchAll = iota
	AgentSalarySearchDay
	AgentSalarySearchMonth
)

func (c *AgentController) Query() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	agentSummary, err := agentdao.AgentDaoEntity.QuerySummary(uid)
	if err != nil {
		if err != orm.ErrNoRows {
			c.ErrorResponse(ERROR_CODE_QUERY_AGENT_FAILED)
			return
		} else { //no rows error, continue
			err = nil
		}
	}

	totalBet, err := reportdao.ReportGameUserDailyDaoEntity.GetYesterdayTotalBet(uid)
	if err != nil {
		common.LogFuncError("%v", err)
	}

	var inviteUrl string
	inviteUrl, err = common2.AppConfigMgr.String(common2.AgentInviteUrl)
	if err != nil { //just record error, no return to client
		common.LogFuncError("%v", err)
	}

	c.SuccessResponse(map[string]interface{}{
		KeyCommission:     agentSummary.SumSalary,
		KeyBalance:        agentSummary.SumCanWithdraw,
		KeyYesterdayChips: totalBet,
		KeyInviteCode:     agentSummary.InviteCode,
		KeyInviteNum:      agentSummary.InviteNum,
		KeyInviteUrl:      fmt.Sprintf(inviteUrl, agentSummary.InviteCode),
	})
}

//easyjson:json
type AgentWithdrawMsg struct {
	Amount int64 `json:"amount"`
}

//提现
func (c *AgentController) Withdraw() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := AgentWithdrawMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil || msg.Amount <= 0 {
		c.ErrorResponse(ERROR_CODE_DECODE_FAILED)
		return
	}

	//如果有存在待审核的订单则不能进行提现
	exitCheckingWithdraw := agentdao.AgentWithdrawDaoEntity.ExitByUidAndStatus(uid, agentdao.AgentWithdrawStatusChecking)
	if exitCheckingWithdraw {
		c.ErrorResponse(ERROR_CODE_WITHDRAW_FROZEN)
		return
	}
	agent, err := agentdao.AgentDaoEntity.Info(uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_AGENT_WITHDRAW_ADD_RECORD_FAILED)
		return
	}
	ok, err := agentdao.AgentDaoEntity.CheckPwdMD5(agent.Uid, agent.SumCanWithdraw, agent.Pwd)
	if err != nil || !ok {
		common.LogFuncError("md5 error:[%v]", agent.Uid)
		c.ErrorResponse(ERROR_CODE_NO_AUTH)
		return
	}

	record, err := agentdao.AgentWithdrawDaoEntity.Add(uid, msg.Amount, int64(common.NowUint32()))
	if err != nil {
		c.ErrorResponse(ERROR_CODE_AGENT_WITHDRAW_ADD_RECORD_FAILED)
		return
	}
	if agent.SumCanWithdraw < msg.Amount {
		common.LogFuncWarning("amount[%v] bigger than balance [%v]", msg.Amount, agent.SumCanWithdraw)
		c.ErrorResponse(ERROR_CODE_AGENT_WITHDRAW_NO_ENOUGH)
		return
	}
	profitThresholdStr, _ := common2.AppConfigMgr.String(common2.ProfitThreshold)
	if len(profitThresholdStr) > 0 {
		profitThreshold, err := common.CurrencyStrToInt64(profitThresholdStr)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_PROFIT_THRESHOLD_PARAM_ERROR)
			return
		}
		if msg.Amount > profitThreshold {
			_ = agentdao.AgentWithdrawDaoEntity.UpdateStatus(record.Id, agentdao.AgentWithdrawStatusChecking, fmt.Sprintf("withdraw num %v bigger than threshold %v", msg.Amount, profitThreshold))
			c.ErrorResponse(ERROR_CODE_WITHDRAW_MORE_THAN_THRESHOLD)
			return
		}
	}

	// decrease commission
	err = agentdao.AgentDaoEntity.DecreaseCommission(uid, msg.Amount, lock.Options{
		LockTimeout: time.Second * time.Duration(common3.Cursvr.LockTimeout),
		RetryCount:  common3.Cursvr.RetryCount,
		RetryDelay:  time.Duration(common3.Cursvr.RetryDelay) * time.Millisecond,
	})
	if err != nil {
		_ = agentdao.AgentWithdrawDaoEntity.UpdateStatus(record.Id, agentdao.AgentWithdrawStatusFailed, fmt.Sprintf("%v while DecreaseCommission", err))
		if err == agentdao.ErrCurrencyNoEnough {
			c.ErrorResponse(ERROR_CODE_AGENT_WITHDRAW_NO_ENOUGH)
			return
		}

		c.ErrorResponse(ERROR_CODE_AGENT_DECREASE_COMMISSION_FAILED)
		return
	}

	// update withdraw status
	err = agentdao.AgentWithdrawDaoEntity.UpdateStatus(record.Id, agentdao.AgentWithdrawStatusDecreased, "")
	if err != nil {
		c.ErrorResponse(ERROR_CODE_AGENT_WITHDRAW_UPDATE_RECORD_STATUS_FAILED)
		return
	}

	// send eusd
	errCode = eosplus.EosPlusAPI.Wealth.Commission(uid, msg.Amount)
	if errCode != ERROR_CODE_SUCCESS {
		_ = agentdao.AgentWithdrawDaoEntity.UpdateStatus(record.Id, agentdao.AgentWithdrawStatusFailed, fmt.Sprintf("error code %d while send eusd", errCode))
		c.ErrorResponse(errCode)
		return
	}

	// update withdraw status
	err = agentdao.AgentWithdrawDaoEntity.UpdateStatus(record.Id, agentdao.AgentWithdrawStatusDone, "")
	if err != nil {
		c.ErrorResponse(ERROR_CODE_AGENT_WITHDRAW_UPDATE_RECORD_STATUS_FAILED)
		return
	}
	//日报
	_ = reportdao.ReportAgentDailyDaoEntity.AddWithDraw(uid, msg.Amount)

	c.SuccessResponse(map[string]interface{}{
		KeyId:     record.Id,
		KeyAmount: msg.Amount,
		KeyCtime:  record.Ctime,
		KeyUtime:  record.Mtime,
		KeyStatus: record.Status,
	})
}

const (
	DefaultPage     = 1
	DefaultPerPage  = 10
	SalaryTypeDay   = 1
	SalaryTypeMonth = 2
)

//提现记录
func (c *AgentController) Withdraws() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	var page, perPage int
	var err error
	page, err = c.GetInt(KeyPage, DefaultPage)
	if err != nil {
		common.LogFuncDebug("get %s failed", KeyPage)
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	perPage, err = c.GetInt(KeyPerPage, DefaultPerPage)
	if err != nil {
		common.LogFuncDebug("get %s failed", KeyPerPage)
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	var (
		total   int
		records []agentdao.AgentWithdrawMsg
	)
	total, records, err = agentdao.AgentWithdrawDaoEntity.Query(uid, page, perPage)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_AGENT_WITHDRAW_QUERY_RECORD_FAILED)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		KeyList: records,
		KeyMeta: MetaMsg{
			Total: total,
			Page:  page,
			Limit: perPage,
		},
	})

	return
}

func (c *AgentController) SumDetail() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	wealth, err := eusddao.WealthDaoEntity.Info(uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	agent, err := agentdao.AgentDaoEntity.Info(uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_QUERY_AGENT_FAILED)
		return
	}
	timestamp := common.GetZeroTime(time.Now().Unix()) - common.DaySeconds
	salary, err := reportdao.ReportGameUserDailyDaoEntity.GetDayTotalSalary(uid, timestamp)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_QUERY_AGENT_FAILED)
		return
	}

	sumMonthDivided, err := reportdao.MonthDividendRecordDaoEntity.GetSumMonthResultDivided(uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_QUERY_AGENT_FAILED)
		return
	}

	t := time.Now()
	begin := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Unix()
	end := time.Date(t.Year(), t.Month()+1, 1, 23, 59, 59, 0, t.Location()).Unix()
	monthDivided, err := reportdao.MonthDividendRecordDaoEntity.GetPreMonthResultDivided(uid, begin, end)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_QUERY_AGENT_FAILED)
		return
	}
	/*
		transferInfo, err := gamedao.GameTransferDaoEntity.GrandTotal(uid, false, 0, common.NowInt64MS())
		if err != nil {
			c.ErrorResponse(ERROR_CODE_QUERY_AGENT_FAILED)
			return
		}
		var transferIn int64
		var transferOut int64
		for _, item := range transferInfo {
			if item.TransferType == uint32(gamedao.TRANSFER_TYPE_IN) {
				transferIn = transferIn + item.EusdInteger
			} else {
				transferOut = transferOut + item.EusdInteger
			}
		}*/
	c.SuccessResponse(map[string]interface{}{
		"status":           wealth.Status,
		"salary":           salary,
		"sum_salary":       agent.SumSalary,
		"month_bonus":      monthDivided,    //上个月分红
		"sum_month_bonus":  sumMonthDivided, //累计月分红
		"sum_can_withdraw": agent.SumCanWithdraw,
		"sum_withdraw":     agent.SumWithdraw,
		//"sum_game_recharge": fmt.Sprintf("%.4f", common.CurrencyInt64ToFloat64(transferIn)),
		//"sum_game_withdraw": fmt.Sprintf("%.4f", common.CurrencyInt64ToFloat64(transferOut)),
	})
}

type SalaryDaySummary struct {
	Salary int64 `orm:"column(salary)" json:"salary,omitempty"`
	Status uint8 `orm:"column(status)" json:"is_get,omitempty"`
	Ctime  int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Type   int64 `orm:"column(type)" json:"type,omitempty"`
}

//日工资月分红列表
func (c *AgentController) SalaryList() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	bTime, err := c.GetInt64("btime", 0)
	if err != nil || bTime == 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	eTime, err := c.GetInt64("etime", 0)
	if err != nil || eTime == 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	sType, err := c.GetInt("type", 0)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	var records []SalaryDaySummary

	//日工资
	if sType == AgentSalarySearchDay || sType == AgentSalarySearchAll {
		day, err := reportdao.ReportGameUserDailyDaoEntity.GetListByTime(uid, bTime, eTime)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_AGENT_WITHDRAW_QUERY_RECORD_FAILED)
			return
		}
		for _, item := range day {
			records = append(records, SalaryDaySummary{
				Salary: item.Salary,
				Ctime:  item.Ctime,
				Status: item.Status,
				Type:   SalaryTypeDay,
			})
		}
	}

	//月分红
	if sType == AgentSalarySearchMonth || sType == AgentSalarySearchAll {
		month, err := reportdao.MonthDividendRecordDaoEntity.GetListByTime(uid, bTime, eTime)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_AGENT_WITHDRAW_QUERY_RECORD_FAILED)
			return
		}
		for _, item := range month {
			records = append(records, SalaryDaySummary{
				Salary: item.SelfDividend,
				Ctime:  item.Ctime,
				Status: item.Status,
				Type:   SalaryTypeMonth,
			})
		}
	}

	c.SuccessResponse(map[string]interface{}{
		KeyList: records,
	})

	return
}

//日工资详细
func (c *AgentController) DaySalary() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	timestamp, err := c.GetInt64(KeyTimestampInput, 0)
	if err != nil || timestamp == 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	records, err := reportdao.ReportGameUserDailyDaoEntity.GetDetailDaySalaryByUid(uid, timestamp)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		KeyList: records,
	})

	return
}

//easyjson:json
/*type AgentGetSalaryMsg struct {
	Timestamp int64 `json:"timestamp"`
}

func (c *AgentController) GetSalary() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := AgentGetSalaryMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DECODE_FAILED)
		return
	}

	records, err := reportdao.ReportGameUserDailyDaoEntity.GetDetailDaySalaryByUid(uid, msg.Timestamp)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	//var pwd string
	var salary int64
	for _, item := range records {
		if item.Status == reportdao.ReportGameUserDailyGet || item.Status == reportdao.ReportGameUserDailyUnGrant {
			c.ErrorResponse(ERROR_CODE_SALARY_HAD_GET)
			return
		}

		pwd, err = common.GenerateDoubleMD5(fmt.Sprintf("%v%v%v", item.Uid, item.Salary, item.Ctime), reportdao.SalarySalt)
		if err != nil || pwd != item.Pwd {
			common.LogFuncError("md5 error:[%v]", item.Uid)
			c.ErrorResponse(ERROR_CODE_NO_AUTH)
			return
		}
		salary = salary + item.Salary
	}

	err = reportdao.ReportGameUserDailyDaoEntity.UpdateSalaryGet(uid, msg.Timestamp)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		KeyTimestamp:       msg.Timestamp,
		"sum_can_withdraw": 0,
	})
	return
}*/

func (c *AgentController) TransferInfo() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	bTime, err := c.GetInt64("btime", 0)
	if err != nil {
		common.LogFuncDebug("get %s failed", KeyPage)
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	eTime, err := c.GetInt64("etime", 0)
	if err != nil {
		common.LogFuncDebug("get %s failed", KeyPerPage)
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	transferInfo, err := gamedao.GameTransferDaoEntity.GrandTotal(uid, true, bTime, eTime)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_QUERY_AGENT_FAILED)
		return
	}

	type GameTransferInfo struct {
		ChannelId    uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
		TransferType uint32 `orm:"column(transfer_type)" json:"transfer_type,omitempty"`
		Amount       string `orm:"column(amount)" json:"amount,omitempty"`
	}

	var list []GameTransferInfo
	for _, item := range transferInfo {
		list = append(list, GameTransferInfo{
			ChannelId:    item.ChannelId,
			TransferType: item.TransferType,
			Amount:       fmt.Sprintf("%.4f", common.CurrencyInt64ToFloat64(item.EusdInteger)),
		})
	}
	c.SuccessResponse(map[string]interface{}{
		KeyList: list,
	})
}
