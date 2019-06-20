package controllers

import (
	"common"
	"math"
	. "otc_error"
	dao3 "utils/otc/dao"
	"utils/report/dao"
	"utils/report/models"
)

type ProfitReportController struct {
	BaseController
}
type SearchCondition struct {
	Uid       uint64 `json:"uid"`
	ChannelId uint32 `json:"channel_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}
type ProfitReport struct {
	ProfitReports map[uint32]*models.ProfitReportDaily
	//todo：佣金提现
	SumWithdraw int64
	EusdBuyNum  uint32
	EusdSoldNum uint32
}

func (c *ProfitReportController) GetProfitReports() {
	//uid, errCode := c.getUidFromToken()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponse(errCode)
	//	return
	//}
	profitReport := new(ProfitReport)

	req := &SearchCondition{}
	req.Uid = 1
	channelId, err := c.GetUint32(KeyChannelId, 1)
	if err != nil {
		common.LogFuncDebug("get %s failed", KeyChannelId)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}
	req.ChannelId = channelId
	startTime, err := c.GetInt64("start_time")
	if err != nil {
		common.LogFuncDebug("get start_time failed")
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}
	req.StartTime = startTime
	endTime, err := c.GetInt64("end_time")
	if err != nil {
		common.LogFuncDebug("get end_time failed")
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}
	req.EndTime = endTime
	//分页获取这个用户这个时间段这个渠道号的分润报表数据
	totalNum, err := dao.ProfitReportDailyDaoEntity.QueryTotalByDateUidChannel(req.Uid, req.ChannelId, startTime, endTime)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	profitReport.ProfitReports = make(map[uint32]*models.ProfitReportDaily)
	if totalNum > 0 {
		reportLimit := 200
		totalPage := math.Ceil(float64(totalNum) / float64(reportLimit))
		totalPageInt := int(totalPage)
		for page := 1; page <= totalPageInt; page++ {
			reports, err := dao.ProfitReportDailyDaoEntity.FindByDataRangeAndUidChannel(req.Uid, req.ChannelId, startTime, endTime, page, reportLimit)
			if err != nil {
				c.ErrorResponse(ERROR_CODE_DB)
				return
			}
			for _, report := range reports {
				if _, ok := profitReport.ProfitReports[report.ChannelId]; !ok {
					profitReport.ProfitReports[report.ChannelId] = new(models.ProfitReportDaily)
					profitReport.ProfitReports[report.ChannelId].Uid = report.Uid
					profitReport.ProfitReports[report.ChannelId].ChannelId = report.ChannelId
				}
				profitReport.ProfitReports[report.ChannelId].Bet += report.Bet
				profitReport.ProfitReports[report.ChannelId].TotalValidBet += report.TotalValidBet
				profitReport.ProfitReports[report.ChannelId].Profit += report.Profit
				profitReport.ProfitReports[report.ChannelId].Salary += report.Salary
				profitReport.ProfitReports[report.ChannelId].SelfDividend += report.SelfDividend
				profitReport.ProfitReports[report.ChannelId].AgentDividend += report.AgentDividend
				profitReport.ProfitReports[report.ChannelId].ResultDividend += report.ResultDividend
				profitReport.ProfitReports[report.ChannelId].GameWithdrawAmount += report.GameWithdrawAmount
				profitReport.ProfitReports[report.ChannelId].GameRechargeAmount += report.GameRechargeAmount
			}
		}

		//统计这个用户，某个时间段所有渠道号的分润报表
		statReports, err := dao.ProfitReportDailyDaoEntity.StatProfitReport(req.Uid, startTime, endTime)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_DB)
			return
		}
		var totalChannelId uint32 = 0
		profitReport.ProfitReports[totalChannelId] = new(models.ProfitReportDaily)
		profitReport.ProfitReports[totalChannelId].Uid = statReports.Uid
		profitReport.ProfitReports[totalChannelId].ChannelId = totalChannelId
		profitReport.ProfitReports[totalChannelId].Bet = statReports.Bet
		profitReport.ProfitReports[totalChannelId].TotalValidBet = statReports.TotalValidBet
		profitReport.ProfitReports[totalChannelId].Profit = statReports.Profit
		profitReport.ProfitReports[totalChannelId].Salary = statReports.Salary
		profitReport.ProfitReports[totalChannelId].SelfDividend = statReports.SelfDividend
		profitReport.ProfitReports[totalChannelId].AgentDividend = statReports.AgentDividend
		profitReport.ProfitReports[totalChannelId].ResultDividend = statReports.ResultDividend
		profitReport.ProfitReports[totalChannelId].GameWithdrawAmount = statReports.WithdrawAmount
		profitReport.ProfitReports[totalChannelId].GameRechargeAmount = statReports.GameRechargeAmount
	}

	otcOrders := dao3.OrdersDaoEntity.StatPeopleByDateRange(req.Uid, startTime, endTime)
	profitReport.EusdBuyNum = otcOrders["eusd_buy"]
	profitReport.EusdSoldNum = otcOrders["eusd_sell"]
	c.SuccessResponse(map[string]interface{}{
		"profit_report": profitReport,
	})
}
