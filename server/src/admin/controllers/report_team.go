package controllers

import (
	"admin/controllers/errcode"
	"common"
	"fmt"
	agentdao "utils/agent/dao"
	otcdao "utils/otc/dao"
	reportdao "utils/report/dao"
)

type ReportTeamController struct {
	BaseController
}

type TeamReport struct {
	Uid                   string `orm:"column(uid);pk" json:"uid,omitempty"`
	Nick                  string `orm:"column(nick);size(100)" json:"nick,omitempty"`
	PeopleNum             uint32 `orm:"column(PeopleNum);size(100)" json:"people_num,omitempty"`
	ActiveNum             uint32 `orm:"column(ActiveNum);size(100)" json:"active_num,omitempty"`
	Level                 uint32 `orm:"column(level)" json:"level,omitempty"`
	EusdBuy               int64  `orm:"column(eusd_buy)" json:"eusd_buy,omitempty"`
	EusdSell              int64  `orm:"column(eusd_sell)" json:"eusd_sell,omitempty"`
	ValidBet              int64  `orm:"column(valid_bet)" json:"valid_bet,omitempty"`
	BetNum                int32  `orm:"column(bet_num)" json:"bet_num,omitempty"`
	Profit                int64  `orm:"column(profit)" json:"profit,omitempty"`
	Salary                int64  `orm:"column(salary)" json:"salary,omitempty"`
	MonthBonus            int64  `orm:"column(result_dividend)" json:"month_bonus,omitempty"`
	GameRecharge          int64  `orm:"column(recharge)" json:"game_recharge,omitempty"`
	GameWithdraw          int64  `orm:"column(withdraw)" json:"game_withdraw,omitempty"`
	CommissionWithdraw    int64  `orm:"column(team_withdraw)" json:"c_withdraw,omitempty"`
	CommissionCanWithdraw int64  `orm:"column(team_can_withdraw)" json:"c_can_withdraw,omitempty"`
}

type ReportTeamReq struct {
	StartTime int64
	EndTime   int64
	ChannelId uint32
}

func (c *ReportTeamController) getReq() (reportTeamReq ReportTeamReq, errCode controllers.ERROR_CODE) {
	var err error
	reportTeamReq.StartTime, err = c.GetInt64("starttime", 0)
	if err != nil || reportTeamReq.StartTime == 0 {
		errCode = controllers.ERROR_CODE_PARAMS_ERROR
		return
	}
	reportTeamReq.EndTime, err = c.GetInt64("endtime", 0)
	if err != nil || reportTeamReq.EndTime == 0 {
		errCode = controllers.ERROR_CODE_PARAMS_ERROR
		return
	}
	reportTeamReq.ChannelId, err = c.GetUint32("channel_id", 0)
	if err != nil {
		errCode = controllers.ERROR_CODE_PARAMS_ERROR
		return
	}
	errCode = controllers.ERROR_CODE_SUCCESS

	return
}

func (c *ReportTeamController) List() {
	c.setOPAction(OPActionReadReportTeamList)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	req, errCode := c.getReq()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	pUid, err := c.GetUint64("puid", 0)
	if err != nil {
		errCode = controllers.ERROR_CODE_PARAMS_ERROR
		return
	}

	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}
	limit, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}
	if limit > DEFAULT_PER_PAGE {
		limit = DEFAULT_PER_PAGE
	}
	nick := c.GetString("nick")
	c.setRequestData(fmt.Sprintf("{\"starttime:\":%d,\"endTime:\":%d,\"channel_id:\":%d,\"puid:\":%d,"+
		"\"page:\":%d, \"limit:\":%d, \"nick:\":%s}", req.StartTime, req.EndTime, req.ChannelId, pUid, page, limit, nick))

	var uid uint64
	if nick != "" {
		otcUser, err := otcdao.UserDaoEntity.GetUserByNick(nick)
		if err != nil {
			c.ErrorResponse(controllers.ERROR_CODE_DB)
			return
		}
		if otcUser.Uid > 0 {
			uid = otcUser.Uid
		}
	}

	//明细
	total, teamReport, err := c.teamReport(page, limit, req.ChannelId, req.StartTime, req.EndTime, uid, pUid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	common.LogFuncDebug("total:%v", total)

	meta := PageInfo{
		Limit: limit,
		Total: total,
		Page:  page,
	}
	c.SuccessResponse(map[string]interface{}{
		"list": teamReport,
		"meta": meta,
	})
}

func (c *ReportTeamController) teamReport(page, limit int, channelId uint32, startTime, endTime int64, uid, pUid uint64) (total int, teamReport []*TeamReport, err error) {
	total, err = reportdao.ReportGameUserDailyDaoEntity.QueryTotalTeamByTime(channelId, startTime, endTime, uid, pUid)
	if err != nil {
		return
	}

	mapTeamReport := make(map[uint64]*TeamReport)
	if total > 0 {
		var teamChannel []reportdao.ChannelTeamReport
		teamChannel, err = reportdao.ReportGameUserDailyDaoEntity.QueryTeamDataByTime(page, limit, channelId, startTime, endTime, uid, pUid)
		if err != nil {
			return
		}

		//日工资、投注相关
		var uIds []string
		for _, item := range teamChannel {
			mapTeamReport[item.Uid] = &TeamReport{
				Uid:      fmt.Sprintf("%d", item.Uid),
				Level:    item.Level,
				ValidBet: item.ValidBet,
				BetNum:   item.BetNum,
				Profit:   item.Profit,
				Salary:   item.Salary,
			}
			uIds = append(uIds, mapTeamReport[item.Uid].Uid)
		}

		//游戏充值
		var teamTransfer []reportdao.ReportTeamGameTransfer
		teamTransfer, err = reportdao.ReportTeamGameTransferDailyDaoEntity.InfoByUIds(channelId, uIds)
		if err != nil {
			return
		}
		for _, item := range teamTransfer {
			mapTeamReport[item.Uid].GameRecharge = item.TeamRecharge
			mapTeamReport[item.Uid].GameWithdraw = item.TeamWithdraw
		}

		//月分红
		var monthBonus []reportdao.ReportTeamMonthBonus
		monthBonus, err = reportdao.MonthDividendRecordDaoEntity.InfoByUIds(uIds)
		if err != nil {
			return
		}
		for _, item := range monthBonus {
			mapTeamReport[item.Uid].MonthBonus = item.ResultDividend
		}

		//eusd部分
		var teamEusd []reportdao.ReportTeam
		teamEusd, err = reportdao.ReportTeamDailyDaoEntity.InfoByUIds(uIds)
		if err != nil {
			return
		}
		for _, item := range teamEusd {
			mapTeamReport[item.Uid].EusdSell = item.EusdSell
			mapTeamReport[item.Uid].EusdBuy = item.EusdBuy
		}

		//用户名
		var userNick []otcdao.UserNick
		userNick, err = otcdao.UserDaoEntity.GetNickByUIds(uIds)
		if err != nil {
			return
		}
		for _, item := range userNick {
			mapTeamReport[item.Uid].Nick = item.Nick
		}

		//佣金可提现，佣金已提现
		var teamCommission []reportdao.ReportCommission
		teamCommission, err = reportdao.ReportCommissionDaoEntity.InfoByUIds(uIds)
		if err != nil {
			return
		}
		for _, item := range teamCommission {
			mapTeamReport[item.Uid].CommissionWithdraw = item.TeamWithdraw
			mapTeamReport[item.Uid].CommissionCanWithdraw = item.TeamCanWithdraw
		}
	}

	//团队人数
	for keyUid, item := range mapTeamReport {
		item.PeopleNum, _ = agentdao.AgentPathDaoEntity.GetTeamPeopleNum(keyUid)
		teamReport = append(teamReport, item)
	}

	return
}

type TeamSumReport struct {
	EusdBuy               int64 `orm:"column(eusd_buy)" json:"eusd_buy,omitempty"`
	EusdSell              int64 `orm:"column(eusd_sell)" json:"eusd_sell,omitempty"`
	ValidBet              int64 `orm:"column(valid_bet)" json:"valid_bet,omitempty"`
	BetNum                int32 `orm:"column(bet_num)" json:"bet_num,omitempty"`
	Profit                int64 `orm:"column(profit)" json:"profit,omitempty"`
	Salary                int64 `orm:"column(salary)" json:"salary,omitempty"`
	MonthBonus            int64 `orm:"column(result_dividend)" json:"month_bonus,omitempty"`
	GameRecharge          int64 `orm:"column(recharge)" json:"game_recharge,omitempty"`
	GameWithdraw          int64 `orm:"column(withdraw)" json:"game_withdraw,omitempty"`
	CommissionWithdraw    int64 `orm:"column(team_withdraw)" json:"c_withdraw,omitempty"`
	CommissionCanWithdraw int64 `orm:"column(team_can_withdraw)" json:"c_can_withdraw,omitempty"`
}

func (c *ReportTeamController) Sum() {
	c.setOPAction(OPActionReadReportTeamSum)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	req, errCode := c.getReq()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	c.setRequestData(fmt.Sprintf("{\"starttime:\":%d,\"endTime:\":%d,\"channel_id:\":%d}",
		req.StartTime, req.EndTime, req.ChannelId))

	teamSumReport, err := c.sumTeamReport(req.ChannelId, req.StartTime, req.EndTime)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	common.LogFuncDebug("%+v", teamSumReport)

	c.SuccessResponse(teamSumReport)
}

func (c *ReportTeamController) sumTeamReport(channelId uint32, startTime, endTime int64) (teamReport TeamSumReport, err error) {
	var teamChannel reportdao.ChannelTeamReport
	teamChannel, err = reportdao.ReportGameUserDailyDaoEntity.Sum(channelId, startTime, endTime)
	if err != nil {
		return
	}

	//日工资、投注相关
	teamReport = TeamSumReport{
		ValidBet: teamChannel.ValidBet,
		BetNum:   teamChannel.BetNum,
		Profit:   teamChannel.Profit,
		Salary:   teamChannel.Salary,
	}

	//游戏充值
	var teamTransfer reportdao.ReportTeamGameTransfer
	teamTransfer, err = reportdao.ReportTeamGameTransferDailyDaoEntity.Sum(channelId, startTime, endTime)
	if err != nil {
		return
	}
	teamReport.GameRecharge = teamTransfer.TeamRecharge
	teamReport.GameWithdraw = teamTransfer.TeamWithdraw

	//月分红
	var monthBonus reportdao.ReportTeamMonthBonus
	monthBonus, err = reportdao.MonthDividendRecordDaoEntity.Sum(startTime, endTime)
	if err != nil {
		return
	}
	teamReport.MonthBonus = monthBonus.ResultDividend

	//eusd部分
	var teamEusd reportdao.ReportTeam
	teamEusd, err = reportdao.ReportTeamDailyDaoEntity.Sum(startTime, endTime)
	if err != nil {
		return
	}
	teamReport.EusdSell = teamEusd.EusdSell
	teamReport.EusdBuy = teamEusd.EusdBuy

	//佣金可提现，佣金已提现
	var teamCommission reportdao.ReportCommission
	teamCommission, err = reportdao.ReportCommissionDaoEntity.Sum()
	if err != nil {
		return
	}
	teamReport.CommissionWithdraw = teamCommission.TeamWithdraw
	teamReport.CommissionCanWithdraw = teamCommission.TeamCanWithdraw

	return
}

func (c *ReportTeamController) Personal() {
	c.setOPAction(OPActionReadReportTeamPersonal)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	req, errCode := c.getReq()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	uid, err := c.GetUint64("uid", 0)
	if err != nil || uid == 0 {
		errCode = controllers.ERROR_CODE_PARAMS_ERROR
		return
	}
	c.setRequestData(fmt.Sprintf("{\"starttime:\":%d,\"endTime:\":%d,\"channel_id:\":%d,\"uid:\":%d}",
		req.StartTime, req.EndTime, req.ChannelId, uid))

	teamPersonalReport, err := c.personal(req.ChannelId, req.StartTime, req.EndTime, uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(teamPersonalReport)
}

func (c *ReportTeamController) personal(channelId uint32, startTime, endTime int64, uid uint64) (teamPersonalReport TeamReport, err error) {
	var teamChannel reportdao.ChannelTeamReport
	teamChannel, err = reportdao.ReportGameUserDailyDaoEntity.Personal(channelId, startTime, endTime, uid)
	if err != nil {
		return
	}

	//日工资、投注相关
	teamPersonalReport = TeamReport{
		ValidBet: teamChannel.ValidBet,
		BetNum:   teamChannel.BetNum,
		Profit:   teamChannel.Profit,
		Salary:   teamChannel.Salary,
	}

	//游戏充值
	var teamTransfer reportdao.ReportTeamGameTransfer
	teamTransfer, err = reportdao.ReportTeamGameTransferDailyDaoEntity.Personal(channelId, startTime, endTime, uid)
	if err != nil {
		return
	}
	teamPersonalReport.GameRecharge = teamTransfer.TeamRecharge
	teamPersonalReport.GameWithdraw = teamTransfer.TeamWithdraw

	//月分红
	var monthBonus reportdao.ReportTeamMonthBonus
	monthBonus, err = reportdao.MonthDividendRecordDaoEntity.Personal(startTime, endTime, uid)
	if err != nil {
		return
	}
	teamPersonalReport.MonthBonus = monthBonus.ResultDividend

	//eusd部分
	var teamEusd reportdao.ReportTeam
	teamEusd, err = reportdao.ReportTeamDailyDaoEntity.Personal(startTime, endTime, uid)
	if err != nil {
		return
	}
	teamPersonalReport.EusdSell = teamEusd.EusdSell
	teamPersonalReport.EusdBuy = teamEusd.EusdBuy

	//佣金可提现，佣金已提现
	var teamCommission reportdao.ReportCommission
	teamCommission, err = reportdao.ReportCommissionDaoEntity.Personal(uid)
	if err != nil {
		return
	}
	teamPersonalReport.CommissionWithdraw = teamCommission.TeamWithdraw
	teamPersonalReport.CommissionCanWithdraw = teamCommission.TeamCanWithdraw

	return
}
