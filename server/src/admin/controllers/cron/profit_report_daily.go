package cron

import (
	"common"
	"fmt"
	"math"
	"time"
	dao3 "utils/game/dao"
	"utils/game/dao/gameapi"
	"utils/report/dao"
	"utils/report/models"
)

//搜索各个数据库表取数据  然后存入profit_report_daily表中

func ProfitReportDailyCron() {
	insertReports := ProfitReportDaily(-1)
	err := dao.ProfitReportDailyDaoEntity.InsertMulti(insertReports)
	if err != nil {
		return
	}
}

var totalChannelId uint32 = 0

func ProfitReportDaily(day int64) (insertReports []*models.ProfitReportDaily) {
	ctime := time.Now().Unix()
	//从日工资报表里获取 昨天的 日工资，盈亏，有效投注金额，要分页获取,日工资报表存的是昨天凌晨0点的时间戳
	profitReports := make(map[uint64]map[uint32]*models.ProfitReportDaily, 0)
	insertReports = make([]*models.ProfitReportDaily, 0)
	dayStart, dayEnd, dayDate, endDate := common.TheOtherDayTimeRange(day)
	fmt.Println(dayStart, dayEnd, dayDate, endDate)
	//分页查询
	//获取总页数,向上取整
	limit := 200
	totalNum, err := dao.ReportGameUserDailyDaoEntity.QueryTotalByData(dayStart, dayEnd)
	totalPage := math.Ceil(float64(totalNum) / float64(limit))
	totalPageInt := int(totalPage)
	uids := make([]uint64, 0)
	for page := 1; page <= totalPageInt; page++ {
		reports, err := dao.ReportGameUserDailyDaoEntity.FindByDayRange(dayStart, dayEnd, page, limit)
		if err != nil {
			return
		}
		for _, report := range reports {
			uid := report.Uid
			if _, ok := profitReports[uid]; !ok {
				profitReports[uid] = make(map[uint32]*models.ProfitReportDaily, 0)
				uids = append(uids, uid)
			}
			channelId := report.ChannelId
			profitReports[uid][channelId].Uid = report.Uid
			profitReports[uid][channelId].Salary = report.Salary
			profitReports[uid][channelId].TotalValidBet = report.TotalValidBet
			profitReports[uid][channelId].Bet = report.Bet
			profitReports[uid][channelId].Ctime = ctime
			insertReports = append(insertReports, profitReports[uid][channelId])
		}
	}
	timestamp := time.Now().Unix() - common.DaySeconds
	fmt.Println(timestamp)
	// 从月分红报表里获取上个月的月分红
	// todo 这个时间格式还要改一下
	startDay, endDay := getThisMonthBeginAndEndDay()
	//再检查对应的值
	monthReports, _ := dao.MonthDividendRecordDaoEntity.FindByDataRange(startDay, endDay)
	for _, dividendRecord := range monthReports {
		uid := dividendRecord.Uid
		profitReports[uid][gameapi.GAME_CHANNEL_RG].SelfDividend = dividendRecord.SelfDividend
		profitReports[uid][gameapi.GAME_CHANNEL_RG].AgentDividend = dividendRecord.AgentDividend
		profitReports[uid][gameapi.GAME_CHANNEL_RG].ResultDividend = dividendRecord.ResultDividend
	}
	//从GameTransfer下获取昨天的游戏充值和游戏提现
	transferRecords, err := dao3.GameTransferDaoEntity.FindByTime(dayStart, dayEnd)
	if err != nil {
		return
	}
	for _, record := range transferRecords {
		uid := record.Uid
		channelId := record.ChannelId
		if record.Type == int(dao3.TRANSFER_TYPE_IN) {
			profitReports[uid][channelId].GameRechargeAmount += record.Eusd
		} else if record.Type == int(dao3.TRANSFER_TYPE_OUT) {
			profitReports[uid][channelId].GameWithdrawAmount += record.Eusd
		}
	}

	//统计所有渠道的总和
	for _, uid := range uids {
		totalReport := new(models.ProfitReportDaily)
		totalReport.ChannelId = totalChannelId
		for _, report := range profitReports[uid] {
			totalReport.AgentDividend += report.AgentDividend
			totalReport.SelfDividend += report.SelfDividend
			totalReport.ResultDividend += report.ResultDividend
			totalReport.GameRechargeAmount += report.GameRechargeAmount
			totalReport.Salary += report.Salary
			totalReport.Profit += report.Profit
			totalReport.TotalValidBet += report.TotalValidBet
			totalReport.Bet += report.Bet
			totalReport.GameWithdrawAmount += report.GameWithdrawAmount
		}
		insertReports = append(insertReports, totalReport)
	}
	return
}
