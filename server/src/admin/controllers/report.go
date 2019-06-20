package controllers

import (
	"admin/controllers/errcode"
	"common"
	"errors"
	"github.com/astaxie/beego/orm"
	dao2 "utils/agent/dao"
	"utils/otc/dao"
	models2 "utils/otc/models"
	dao3 "utils/report/dao"
	"utils/report/models"
)

type ReportController struct {
	BaseController
}
type SearchCondition struct {
	UserName  string `json:"user_name"`
	ChannelId uint32 `json:"channel_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Limit     int    `json:"limit"`
	Page      int    `json:"page"`
}

var totalChannelId uint32 = 0

type PersonalReport struct {
	Uid            uint64 `json:"uid"`
	Nick           string `json:"nick"`
	RegisterTime   int64  `json:"register_time"`
	LastLoginTime  int64  `json:"last_login_time"`
	ReferrerName   string `json:"referrer_name"`
	AgentLv        uint32 `json:"agent_lv"`
	SumCanWithdraw int64  `json:"sum_can_withdraw"` //佣金余额
	SumWithdraw    int64  `json:"sum_withdraw"`     //佣金提现
	*models.ProfitReportDaily
	EusdBuyNum  uint32 `json:"eusd_buy_num"`
	EusdSoldNum uint32 `json:"eusd_sold_num"`
}

type totalReport struct {
	Bet                int64  `json:"bet"`
	TotalValidBet      int64  `json:"total_valid_bet"`
	Profit             int64  `json:"profit"`
	Salary             int64  ` json:"salary"`
	ResultDividend     int64  `json:"result_dividend"`
	GameWithdrawAmount uint32 `json:"game_withdraw_amount"`
	GameRechargeAmount uint32 ` json:"game_recharge_amount"`
	EusdBuyNum         uint32 `json:"eusd_buy_num"`
	EusdSoldNum        uint32 `json:"eusd_sold_num"`
	SumCanWithdraw     int64  `json:"sum_can_withdraw"` //佣金余额
	SumWithdraw        int64  `json:"sum_withdraw"`     //佣金提现
}

// 每个玩家只有一条数据
func (c *ReportController) GetPersonalReport() {
	//	按搜索条件获取指定数据
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionReadMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}
	req := &SearchCondition{}
	req.UserName = c.GetString("user_name", "")
	//-1代表查询所有的渠道的数据
	channelId, err := c.GetUint32("channel_id", totalChannelId)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	startTime, err := c.GetInt64("start_time")
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	req.StartTime = startTime
	endTime, err := c.GetInt64("end_time")
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	req.EndTime = endTime
	req.ChannelId = channelId
	req.Page, _ = c.GetInt(KEY_PAGE, 1)
	req.Limit, _ = c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	total := new(totalReport)
	//如果用户名为空则查询所有人的纪录
	//如果用户名为空则先查记录表
	//如果用户名不为空先查用户表
	res := map[string]interface{}{}
	totalNum := 1
	personalReports := make([]*PersonalReport, 0)
	if req.UserName != "" {
		//有输入名字 查询的是这些用户这个时间段的指定渠道号所有分润报表记录的总和情况，根据用户的创建时间排序

		personalReport := c.queryByName(req.UserName, req.ChannelId, req.StartTime, req.EndTime, req.Page, req.Limit)

		personalReports = make([]*PersonalReport, 0)
		if personalReport != nil {
			personalReports = append(personalReports, personalReport)
		}
		res["list"] = personalReports
		meta := map[string]interface{}{
			"total": totalNum,
			"page":  req.Page,
			"limit": req.Limit,
		}
		res["meta"] = meta

	} else {
		//没有输入记录查询的是这个时间段里指定渠道号的每天的分润报表情况，由前往后倒叙排列
		//不能按照上面那样子把一个用户所有数据统计起来，因为如果统计起来在分页的时候会有错误
		totalNum, personalReports = c.queryMulti(req.ChannelId, req.StartTime, req.EndTime, req.Page, req.Limit)
		meta := map[string]interface{}{
			"total": totalNum,
			"page":  req.Page,
			"limit": req.Limit,
		}
		res["meta"] = meta
		res["list"] = personalReports
	}
	totals := make([]*totalReport, 0)
	getTotal(total, personalReports)
	totals = append(totals, total)
	res["total"] = totals

	c.SuccessResponseAndLog(OPActionEditAgentDividendPosition, string(c.Ctx.Input.RequestBody), res)
}
func getTotal(total *totalReport, reports []*PersonalReport) {
	for _, report := range reports {
		total.Bet += report.Bet
		total.TotalValidBet += report.TotalValidBet
		total.Profit += report.Profit
		total.Salary += report.Salary
		total.ResultDividend += report.ResultDividend
		total.GameWithdrawAmount += report.GameWithdrawAmount
		total.GameRechargeAmount += report.GameRechargeAmount
		total.EusdBuyNum += report.EusdBuyNum
		total.EusdSoldNum += report.EusdSoldNum
		total.SumCanWithdraw += report.SumCanWithdraw
		total.SumWithdraw += report.SumWithdraw
	}
}
func (c *ReportController) queryByName(nick string, channelId uint32, startTime, endTime int64, page, limit int) (personalReport *PersonalReport) {
	user, err := dao.UserDaoEntity.GetUserByNick(nick)
	if err != nil {
		if err == orm.ErrNoRows {
			return
		}
		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	personalReport = new(PersonalReport)
	//根据用户构建personalReportsMap
	personalReport.Nick = user.Nick
	personalReport.RegisterTime = user.Ctime
	personalReport.LastLoginTime = user.LastLoginTime
	personalReport.Uid = user.Uid

	//	根据agentPath查找推介人和自己的代理等级
	agentPath, err := dao2.AgentPathDaoEntity.GetByUid(user.Uid)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	personalReport.AgentLv = agentPath.Level

	if personalReport.AgentLv > 1 {
		parentUserInfo, err := dao.UserDaoEntity.GetByUid(agentPath.ParentUid)
		if err != nil {
			c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
			return
		}
		personalReport.ReferrerName = parentUserInfo.Nick
	}

	//构建个人报表中的分润报表数据
	personalReport.ProfitReportDaily, err = dao3.ProfitReportDailyDaoEntity.StatProfitReportByUidChannelId(user.Uid, channelId, startTime, endTime)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	//获取用户eusd消费情况
	otcOders := dao.OrdersDaoEntity.StatPeopleByDateRange(personalReport.Uid, startTime, endTime)
	personalReport.EusdBuyNum = otcOders["eusd_buy"]
	personalReport.EusdSoldNum = otcOders["eusd_sell"]

	//获取佣金提现和佣金余额
	agent, err := dao2.AgentDaoEntity.Info(personalReport.Uid)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	personalReport.SumCanWithdraw = agent.SumCanWithdraw
	personalReport.SumWithdraw = agent.SumWithdraw
	return
}

//没有的话就分页查询
//如果没有输入名字进行模糊查询的话 就在PersonalReport表中查询指定渠道指定天数的前几条记录，并且查询用户表,代理表,eusdt消费表的纪录整合起来形成个人报表记录返回给前端
//传入时间和渠道,渠道可以不用传因为要计算总的值也是要把所有的渠道订单取出来
//要传渠道，不是这个渠道的数据不返回给前端，并且还要考虑到分页多个渠道号会导致查到的数据不足
//如果分到这一页的时候这个用户有多个数据，则会重复了，因为都是同一个uid，应该再加上ctime作为值
func (c *ReportController) queryMulti(channelId uint32, startTime, endTime int64, page, limit int) (total int, personReports []*PersonalReport) {
	//找出这个时间段总共有多少条数据
	personReports = make([]*PersonalReport, 0)
	reports := make([]*models.ProfitReportDaily, 0)
	err := errors.New("no error")

	total, err = dao3.ProfitReportDailyDaoEntity.QueryTotalByDateAndChannelId(channelId, startTime, endTime)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	//分页查询
	reports, err = dao3.ProfitReportDailyDaoEntity.StatProfitReportByChannelId(channelId, startTime, endTime, page, limit)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	if len(reports) == 0 {
		return
	}

	uids := make([]uint64, 0)
	uidsMap := make(map[uint64]uint64, 0)

	for _, report := range reports {
		if _, ok := uidsMap[report.Uid]; ok {
			continue
		}
		reportWithoutName := new(PersonalReport)
		reportWithoutName.Uid = report.Uid
		reportWithoutName.ProfitReportDaily = report
		personReports = append(personReports, reportWithoutName)
		uidsMap[report.Uid] = report.Uid
		uids = append(uids, report.Uid)
	}
	//获取用户信息
	usersInfo, err := dao.UserDaoEntity.FindByUids(uids)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	usersInfoMap := make(map[uint64]*models2.User, 0)
	for _, user := range usersInfo {
		usersInfoMap[user.Uid] = user
	}
	for _, report := range personReports {
		report.Nick = usersInfoMap[report.Uid].Nick
		report.RegisterTime = usersInfoMap[report.Uid].Ctime
		report.LastLoginTime = usersInfoMap[report.Uid].LastLoginTime
	}
	//获取代理情况
	agentPathsMap, err := dao2.AgentPathDaoEntity.FindByUids(uids)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	parentUids := make([]uint64, 0, len(agentPathsMap))

	for _, report := range personReports {
		report.AgentLv = agentPathsMap[report.Uid].Level
		parentUids = append(parentUids, agentPathsMap[report.Uid].ParentUid)
	}
	for _, agentPath := range agentPathsMap {
		parentUids = append(parentUids, agentPath.ParentUid)
	}
	parentUserInfos, err := dao.UserDaoEntity.FindMapByUids(parentUids)

	for _, report := range personReports {
		if agentPathsMap[report.Uid].Level == 1 {
			continue
		}
		if _, ok := parentUserInfos[report.Uid]; !ok {
			continue
		}
		report.ReferrerName = parentUserInfos[report.Uid].Nick
	}
	//eusd交易情况
	for _, personalReport := range personReports {
		otcOders := dao.OrdersDaoEntity.StatPeopleByDateRange(personalReport.Uid, startTime, endTime)
		personalReport.EusdBuyNum = otcOders["eusd_buy"]
		personalReport.EusdSoldNum = otcOders["eusd_sell"]
	}

	for _, personalReport := range personReports {
		agent, err := dao2.AgentDaoEntity.Info(personalReport.Uid)
		if err != nil {
			c.ErrorResponseAndLog(OPActionEditAgentDividendPosition, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
			return
		}
		personalReport.SumCanWithdraw = agent.SumCanWithdraw
		personalReport.SumWithdraw = agent.SumWithdraw
	}
	return
}
