package cron

import (
	"common"
	"eusd/eosplus"
	"fmt"
	"strconv"
	"time"
	dao3 "utils/admin/dao"
	models2 "utils/admin/models"
	dao2 "utils/report/dao"
	"utils/report/models"

	"utils/agent/dao"
)

func TestMonthDividend() {
	err := dao2.GameUserMonthReportDaoEntity.DeleteAllData()
	if err != nil {
		fmt.Println("dao.GameUserMonthReportDaoEntity.DeleteAllData err ", err)
	}
	err = dao2.MonthDividendRecordDaoEntity.DeleteAllData()
	if err != nil {
		fmt.Println("MonthDividendRecordDaoEntity.DeleteAllData() err ", err)
	}
	err = dao3.ActivityUserConfDaoEntity.Insert(7, 3000, 200)
	if err != nil {
		fmt.Println("dao3.ActivityUserConfDaoEntity.Insert ", err)
	}
	err = dao3.MonthDividendPositionConfDaoEntity.DeleteAllData()
	if err != nil {
		fmt.Println("dao3.MonthDividendPositionConfDaoEntity.DeleteAllData ", err)
	}
	//自己增加月分红配置信息
	createMonthDividendConf()

	firstLevelPre := "firstLevel%v"
	secondLevelPre := ":secondLevel%v"
	thirdLevelPre := ":thirdLevel%v"
	fouthLevelPre := ":fouthLevelPre%v"
	uids := make([]uint64, 0, 1200)
	dividendUserNum := 0

	//var bet, validBet, totalValidBet, profit int64 = 100000, 100000, 100000, 100000
	var bet, validBet, totalValidBet, profit int64 = 100000, 100000, 100000, -130000
	days := 8
	err = dao.AgentPathDaoEntity.DeleteForTest(uids)
	if err != nil {
		fmt.Println("dao.AgentPathDaoEntity.DeleteForTest(uids) err ", err)
	}
	//创建一级代理的时候不用传parentInviteCode，二级代理要传它属于的一级代理的InviteCode
	//一级代理总共有
	pres := make(map[uint64]string, 0)
	for i := 1; i <= 3; i++ {
		inviteCode := firstLevelPre + strconv.Itoa(i)
		pre := fmt.Sprintf(firstLevelPre, i)
		pres[uint64(i)] = pre
		creat(uint64(i), pre, 0, 1, inviteCode)
		uids = append(uids, uint64(i))
	}
	secondUids := make([]uint64, 0)
	for _, pUid := range uids {
		for i := 0; i < 5; i++ {
			suid := pUid*10 + uint64(i)
			inviteCode := secondLevelPre + strconv.Itoa(int(suid))
			pre := fmt.Sprintf(secondLevelPre, suid)
			pres[suid] = pres[pUid] + pre
			creat(suid, pres[suid], pUid, 2, inviteCode)
			uids = append(uids, suid)
			secondUids = append(secondUids, suid)
		}
	}

	thirdUids := make([]uint64, 0)
	for _, pUid := range secondUids {
		for i := 0; i < 5; i++ {
			tuid := pUid*100 + uint64(i)
			inviteCode := thirdLevelPre + strconv.Itoa(int(tuid))
			pre := fmt.Sprintf(thirdLevelPre, tuid)
			pres[tuid] = pres[pUid] + pre
			creat(tuid, pres[tuid], pUid, 3, inviteCode)
			uids = append(uids, tuid)
			thirdUids = append(thirdUids, tuid)
		}
	}
	dividendUserNum = len(uids)
	fouthUids := make([]uint64, 0)
	for _, pUid := range thirdUids {
		for i := 0; i < 8; i++ {
			fuid := pUid*1000 + uint64(i)
			inviteCode := fouthLevelPre + strconv.Itoa(int(fuid))
			pre := fmt.Sprintf(fouthLevelPre, fuid)
			pres[fuid] = pres[pUid] + pre
			creat(fuid, pres[fuid], pUid, 4, inviteCode)
			uids = append(uids, fuid)
			fouthUids = append(fouthUids, fuid)
		}
	}
	createProfitOrLossDailyReport(uids, days, bet, validBet, totalValidBet, profit)
	fmt.Println("create data end ")

	GameUserMonthDividend()
	fmt.Println("计算逻辑结束")

	//checkProfit(uids, dividendUserNum, days, bet, validBet, totalValidBet, profit)
	checkLoss(uids, dividendUserNum, days, bet, validBet, totalValidBet, profit)
}

func createMonthDividendConf() {
	confs := make(map[int][]*models2.MonthDividendPositionConf, 0)
	ctime := common.NowInt64MS()
	agentLv := 1
	oneLevelPosition := []int32{1, 2, 3, 4, 5, 6}
	oneLevelMin := map[int32]int64{1: 2000, 2: 15000, 3: 50000, 4: 150000, 5: 600000, 6: 1500000}
	oneLevelMax := map[int32]int64{1: 15000, 2: 50000, 3: 150000, 4: 600000, 5: 1500000, 6: 0}
	oneLevelDividendRatio := map[int32]int32{1: 25, 2: 30, 3: 35, 4: 40, 5: 45, 6: 50}
	oneLevelActivityNum := map[int32]int32{1: 0, 2: 5, 3: 10, 4: 15, 5: 20, 6: 30}
	for _, position := range oneLevelPosition {
		if _, ok := confs[agentLv]; !ok {
			confs[agentLv] = make([]*models2.MonthDividendPositionConf, 0)
		}
		cfg := &models2.MonthDividendPositionConf{
			AgentLv:       int32(agentLv),
			Position:      position,
			Min:           oneLevelMin[position],
			Max:           oneLevelMax[position],
			ActivityNum:   oneLevelActivityNum[position],
			DividendRatio: oneLevelDividendRatio[position],
			Ctime:         ctime,
			Utime:         ctime,
		}
		confs[agentLv] = append(confs[agentLv], cfg)
	}
	agentLv = 2
	twoLevelPosition := []int32{1, 2, 3, 4, 5, 6}
	twoLevelMin := map[int32]int64{1: 2000, 2: 12000, 3: 40000, 4: 120000, 5: 450000, 6: 1000000}
	twoLevelMax := map[int32]int64{1: 12000, 2: 40000, 3: 120000, 4: 450000, 5: 1000000, 6: 0}
	twoLevelDividendRatio := map[int32]int32{1: 15, 2: 20, 3: 25, 4: 30, 5: 35, 6: 40}
	twoLevelActivityNum := map[int32]int32{1: 0, 2: 5, 3: 10, 4: 15, 5: 20, 6: 30}
	for _, position := range twoLevelPosition {
		if _, ok := confs[agentLv]; !ok {
			confs[agentLv] = make([]*models2.MonthDividendPositionConf, 0)
		}
		cfg := &models2.MonthDividendPositionConf{
			AgentLv:       int32(agentLv),
			Position:      position,
			Min:           twoLevelMin[position],
			Max:           twoLevelMax[position],
			ActivityNum:   twoLevelActivityNum[position],
			DividendRatio: twoLevelDividendRatio[position],
			Ctime:         ctime,
			Utime:         ctime,
		}
		confs[agentLv] = append(confs[agentLv], cfg)
	}
	agentLv = 3
	threeLevelPosition := []int32{1, 2, 3, 4, 5, 6}
	threeLevelMin := map[int32]int64{1: 2000, 2: 10000, 3: 30000, 4: 80000, 5: 300000, 6: 600000}
	threeLevelMax := map[int32]int64{1: 10000, 2: 30000, 3: 80000, 4: 300000, 5: 600000, 6: 0}
	threeLevelDividendRatio := map[int32]int32{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30}
	threeLevelActivityNum := map[int32]int32{1: 0, 2: 5, 3: 10, 4: 15, 5: 20, 6: 30}
	for _, position := range threeLevelPosition {
		if _, ok := confs[agentLv]; !ok {
			confs[agentLv] = make([]*models2.MonthDividendPositionConf, 0)
		}
		cfg := &models2.MonthDividendPositionConf{
			AgentLv:       int32(agentLv),
			Position:      position,
			Min:           threeLevelMin[position],
			Max:           threeLevelMax[position],
			ActivityNum:   threeLevelActivityNum[position],
			DividendRatio: threeLevelDividendRatio[position],
			Ctime:         ctime,
			Utime:         ctime,
		}
		confs[agentLv] = append(confs[agentLv], cfg)

	}

	for _, conf := range confs {
		err := dao3.MonthDividendPositionConfDaoEntity.InsertForTest(conf)
		if err != nil {
			common.LogFuncError("mysql_err:%v", err)
		}
	}
}

//创建所有都是盈利或者都是亏损并且所有用户都是活跃用户的日流水数据
func createProfitOrLossDailyReport(uids []uint64, days int, bet, validBet, totalValidBet, profit int64) {
	err := dao2.ReportGameUserDailyDaoEntity.DeleteForTest([]uint64{})
	if err != nil {
		fmt.Println(err)
	}
	dailyReports := make([]*models.ReportGameUserDaily, 0)
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start := thisMonth.AddDate(0, -1, 0)
	for _, uid := range uids {
		for day := 0; day < days; day++ {
			ctime, _ := strconv.ParseInt(start.Format(DATE_FORMAT), 10, 64)
			report := &models.ReportGameUserDaily{
				Uid:           uid,
				ChannelId:     1,
				Bet:           bet,
				ValidBet:      validBet,
				TotalValidBet: totalValidBet,
				Profit:        profit,
				Ctime:         int64(ctime),
			}
			dailyReports = append(dailyReports, report)
			start = start.AddDate(0, 0, 1)
		}
		start = start.AddDate(0, 0, -days)
	}
	err = dao2.ReportGameUserDailyDaoEntity.InsertForTest(dailyReports)
	if err != nil {
		fmt.Println(err)
	}
}

func checkProfit(uids []uint64, dividendNum int, days int, bet, validBet, totalValidBet, profit int64) {
	//先检查月流水记录
	//先检查月流水条数，条数应该为所有用户的数量
	startDay, endDay := getThisMonthBeginAndEndDay()
	allRight := true

	//再检查对应的值
	monthReports, _ := dao2.GameUserMonthReportDaoEntity.FindByDataRange(startDay, endDay)
	activityNums := map[uint32]int32{
		1: 231,
		2: 46,
		3: 9,
	}
	if len(monthReports) != len(uids) {
		common.LogFuncError("月流水插入条数不对 应该插入的条数%v 实际插入的条数%v", len(uids), len(monthReports))
		allRight = false
	}
	for _, report := range monthReports {
		shouldActivityNum := activityNums[report.AgentLevel]
		if report.ActivityAgentNum != shouldActivityNum {
			common.LogFuncError("report.ActivityAgentNum[%v] != shouldActivityNum[%v] agentLevel[%v]", report.ActivityAgentNum, shouldActivityNum, report.AgentLevel)
			allRight = false
			break
		}
		if report.PlayGameDay != int32(days) {
			common.LogFuncError("report.PlayGameDay[%v] != int32(days)[%v] ", report.PlayGameDay, days)
			allRight = false
			break
		}
		if report.BetAmount != int64(days)*totalValidBet {
			common.LogFuncError("report.BetAmount[%v] != int64(days)*totalValidBet[%v] ", report.BetAmount, int64(days)*totalValidBet)
			allRight = false
			break
		}
		if report.EffectiveBetAmount != int64(days)*validBet {
			common.LogFuncError("report.EffectiveBetAmount[%v] != int64(days)*validBet[%v] ", report.EffectiveBetAmount, int64(days)*validBet)
			allRight = false
			break
		}
		if report.Profit != int64(days)*profit {
			common.LogFuncError("report.profit[%v] != int64(days)*profit[%v] ", report.Profit, int64(days)*profit)
			allRight = false
			break
		}

	}
	if allRight {
		fmt.Println("check GameUserMonthReport data all right")
	}
	//再检查分红表记录
	//先检查月流水条数，条数应该为所有可以分红的用户的数量
	allRight = true

	dividendRecords, _ := dao2.MonthDividendRecordDaoEntity.FindByDataRange(startDay, endDay)
	if len(dividendRecords) != dividendNum {
		common.LogFuncError("月分红插入条数不对 应该插入的条数%v 实际插入的条数%v", dividendNum, len(dividendRecords))
		allRight = false
	}
	for _, record := range dividendRecords {
		if record.SelfDividend != 0 {
			common.LogFuncError("record.SelfDividend[%v] != 0 ", record.SelfDividend)
			allRight = false
			break
		}
		if record.AgentDividend != 0 {
			common.LogFuncError("record.AgentDividend[%v] != 0 ", record.AgentDividend)
			allRight = false
			break
		}
		if record.ResultDividend != 0 {
			common.LogFuncError("record.ResultDividend[%v] != 0 ", record.ResultDividend)
			allRight = false
			break
		}
		//if record.ReceiveStatus == dao2.Received {
		//	common.LogFuncError("record.IsReceived[%v] != false ", record.ReceiveStatus)
		//	allRight = false
		//	break
		//}
	}
	if allRight {
		fmt.Println("check MonthDividendRecord data all right")
	}
}

//创建所有用户亏欠并且都是活跃用户的日流水数据
func checkLoss(uids []uint64, dividendNum int, days int, bet, validBet, totalValidBet, profit int64) {
	//先检查月流水记录
	//先检查月流水条数，条数应该为所有用户的数量
	startDay, endDay := getThisMonthBeginAndEndDay()
	allRight := true
	activityNums := map[uint32]int32{
		1: 231,
		2: 46,
		3: 9,
	}
	//再检查对应的值
	monthReports, _ := dao2.GameUserMonthReportDaoEntity.FindByDataRange(startDay, endDay)
	if len(monthReports) != len(uids) {
		common.LogFuncError("月流水插入条数不对 应该插入的条数%v 实际插入的条数%v", len(uids), len(monthReports))
		allRight = false
	}
	monthReportsMap := make(map[uint64]*models.GameUserMonthReport, len(monthReports))
	for _, record := range monthReports {
		monthReportsMap[record.Uid] = record
	}
	for _, report := range monthReports {
		shouldActivityNum := activityNums[report.AgentLevel]
		if report.ActivityAgentNum != shouldActivityNum {
			common.LogFuncError("report.ActivityAgentNum[%v] != shouldActivityNum[%v] agentLevel[%v]", report.ActivityAgentNum, shouldActivityNum, report.AgentLevel)
			allRight = false
			break
		}
		if report.PlayGameDay != int32(days) {
			common.LogFuncError("report.PlayGameDay[%v] != int32(days)[%v] ", report.PlayGameDay, days)
			allRight = false
			break
		}
		if report.BetAmount != int64(days)*totalValidBet {
			common.LogFuncError("report.BetAmount[%v] != int64(days)*totalValidBet[%v] ", report.BetAmount, int64(days)*totalValidBet)
			allRight = false
			break
		}
		if report.EffectiveBetAmount != int64(days)*validBet {
			common.LogFuncError("report.EffectiveBetAmount[%v] != int64(days)*validBet[%v] ", report.EffectiveBetAmount, int64(days)*validBet)
			allRight = false
			break
		}
		if report.Profit != int64(days)*profit {
			common.LogFuncError("report.profit[%v] != int64(days)*profit[%v] ", report.Profit, int64(days)*profit)
			allRight = false
			break
		}
	}
	if allRight {
		fmt.Println("check GameUserMonthReport data all right")
	}
	//再检查分红表记录
	//先检查月流水条数，条数应该为所有可以分红的用户的数量
	allRight = true

	dividendRecords, _ := dao2.MonthDividendRecordDaoEntity.FindByDataRange(startDay, endDay)
	if len(dividendRecords) != dividendNum {
		common.LogFuncError("月分红插入条数不对 应该插入的条数%v 实际插入的条数%v", dividendNum, len(dividendRecords))
		allRight = false
	}

	selfDividends := map[uint32]int64{
		1: eosplus.QuantityFloat64ToInt64(104 * 0.45 * 10000),
		2: eosplus.QuantityFloat64ToInt64(104 * 0.40 * 10000),
		3: eosplus.QuantityFloat64ToInt64(104 * 0.10 * 10000),
	}
	agentDividends := map[uint32]int64{
		1: eosplus.QuantityFloat64ToInt64(5 * 104 * 0.40 * 10000),
		2: eosplus.QuantityFloat64ToInt64(5 * 104 * 0.10 * 10000),
		3: 0,
	}
	for _, record := range dividendRecords {
		agentLv := monthReportsMap[record.Uid].AgentLevel
		if record.SelfDividend != selfDividends[agentLv] {
			common.LogFuncError("record.SelfDividend[%v] != %v agentLv=%v", record.SelfDividend, selfDividends[agentLv], agentLv)
			allRight = false
			break
		}
		if record.AgentDividend != agentDividends[agentLv] {
			common.LogFuncError("record.AgentDividend[%v] != %v agentLv=%v", record.AgentDividend, agentDividends[agentLv], agentLv)
			allRight = false
			break
		}
		resultDividend := selfDividends[agentLv] - agentDividends[agentLv]
		if record.ResultDividend != resultDividend {
			common.LogFuncError("record.ResultDividend[%v] != %v agentLv=%v", record.ResultDividend, resultDividend, agentLv)
			allRight = false
			break
		}
		//if record.ReceiveStatus == dao2.Received {
		//	common.LogFuncError("record.IsReceived[%v] != false ", record.ReceiveStatus)
		//	allRight = false
		//	break
		//}
	}
	if allRight {
		fmt.Println("check MonthDividendRecord data all right")
	}
}

//创建有亏钱所有用户的是活跃用户的日流水数据
//创建有不是活跃用户的日流水数据
func creat(uid uint64, path string, pUid uint64, level uint32, invite_code string) {
	err := dao.AgentPathDaoEntity.InsertForTest(uid, path, pUid, level, invite_code)
	if err != nil {
		fmt.Println(err)
	}
}
func getThisMonthBeginAndEndDay() (beginDay, endDay int64) {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start := thisMonth.AddDate(0, 0, 0).Format(DATE_FORMAT)
	end := thisMonth.AddDate(0, 1, -1).Format(DATE_FORMAT)
	beginDayInt, _ := strconv.ParseInt(start, 10, 64)
	endDayInt, _ := strconv.ParseInt(end, 10, 64)
	beginDay = int64(beginDayInt)
	endDay = int64(endDayInt)

	return
}
