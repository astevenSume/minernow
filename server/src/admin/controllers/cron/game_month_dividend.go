package cron

import (
	common3 "admin/common"
	"admin/controllers/errcode"
	"common"
	"eusd/eosplus"
	"fmt"
	"github.com/bsm/redis-lock"
	"math"
	"strconv"
	"time"
	"umeng_push/uemng_plus"
	"utils/admin/dao"
	"utils/admin/models"
	dao2 "utils/agent/dao"
	"utils/game/dao/gameapi"
	otcDao "utils/otc/dao"
	dao3 "utils/report/dao"
	models2 "utils/report/models"
	usdtDao "utils/usdt/dao"
)

const DATE_FORMAT = "20060102"

// 亏损分红配置表
//		1，2，3分别代表分1级代理分红配置表，2级代理分红配置表，3级代理分红配置表
//      LowerLossLimit和UpperLossLimit单位是万元
//todo 上级代理的亏损额计算逻辑改变了
//todo 日工资有3个渠道 这里只要获取一个渠道的月分红就好了
// todo 日工资报表存的是昨天凌晨0点的时间戳看下这里的查询时间格式对不对
func GameUserMonthDividend() {
	dividendConf, err := dao.MonthDividendPositionConfDaoEntity.QueryPageDividendConfs()
	if err != nil {
		return
	}
	activityUserConf, err := dao.ActivityUserConfDaoEntity.GetConf()
	if err != nil {
		return
	}
	lowDividendAgentLevel := uint32(len(dividendConf))
	if lowDividendAgentLevel == 0 {
		common.LogFuncError("没有配置月分红配置表.")
	}
	//获取上个月的日报表
	lastMonthBeginDay, lastMonthEndDay := getLastMonthBeginAndEndDay()
	//以UID为Key,后台设置的分红等级为值
	agentWhite := make(map[uint64]uint32, 0)
	//以UID为key
	monthReportsMap := make(map[uint64]*models2.GameUserMonthReport, 0)
	monthDividendRecordsMap := make(map[uint64]*models2.MonthDividendRecord, 0)
	//直接代理关系  map[代理等级]map[代理uid][该代理的所有直接下级代理uid],目前第一层的代理等级只到2
	agentRelationship := make(map[uint32]map[uint64][]uint64, 0)
	//每个等级之下的uid
	agentLvUid := make(map[uint32][]uint64, 0)
	//可以进行分红的月报表(月分红报表是所有代理的月分红报表，dividendMonthReports是只有可以进行分红的用户的月分红报表)
	dividendMonthReports := make([]*models2.GameUserMonthReport, 0)
	//用来最终插入月流水表和月分红表的数据
	monthReports := make([]*models2.GameUserMonthReport, 0)
	monthDividendRecords := make([]*models2.MonthDividendRecord, 0)
	//获取月分红白名单配置数据
	monthDividendWhiteListMap, err := dao.MonthDividendWhiteListDaoEntity.GetAllMap()
	if err != nil {
		return
	}
	//最低的代理等级
	//通过刚刚获取的所有日流水报表记录来生成 monthReportsMap
	//获取总页数,向上取整
	limit := 200
	totalNum, err := dao3.ReportGameUserDailyDaoEntity.QueryTotalByData(lastMonthBeginDay, lastMonthEndDay)
	totalPage := math.Ceil(float64(totalNum) / float64(limit))
	totalPageInt := int(totalPage)
	ctime := time.Now().Unix()
	//dailyRecords分页获取，遍历完分页获得的数据在遍历下一个分页数据
	for page := 1; page <= totalPageInt; page++ {
		dailyRecords, _ := dao3.ReportGameUserDailyDaoEntity.FindByDayRangeAndChannelId(lastMonthBeginDay, lastMonthEndDay, gameapi.GAME_CHANNEL_RG, page, limit)
		for _, record := range dailyRecords {
			uid := record.Uid
			if _, ok := monthReportsMap[uid]; !ok {
				monthReportsMap[uid] = new(models2.GameUserMonthReport)
				monthReportsMap[uid].Uid = uid
				monthReportsMap[uid].Ctime = ctime
			}
			//计算自己的盈利值
			monthReportsMap[uid].Profit += record.Profit
			monthReportsMap[uid].BetAmount += record.Bet
			monthReportsMap[uid].EffectiveBetAmount += record.TotalValidBet
			//有一条日流水记录就代表玩家今天玩了游戏
			monthReportsMap[uid].PlayGameDay += 1
		}
	}

	//判断该用户是否是活跃用户,并且构建最终应该插入月报表的数据
	for _, report := range monthReportsMap {
		if report.PlayGameDay > activityUserConf.PlayGameDay && report.EffectiveBetAmount > activityUserConf.EffectiveBetAmount || report.BetAmount > activityUserConf.BetAmount {
			report.IsActivityUser = true
		}
		monthReports = append(monthReports, report)
	}

	//构建月流水中和代理相关的数据
	//主要是代理等级和上级代理
	for _, monthReport := range monthReports {
		selfUid := monthReport.Uid
		//获取自己的代理等级
		agentPath, err := dao2.AgentPathDaoEntity.Info(selfUid)
		if err != nil {
			return
		}
		agentWhite[monthReport.Uid] = agentPath.DividendPosition
		monthReport.AgentLevel = agentPath.Level
		upAgentUid := agentPath.ParentUid
		//构建agentRelationship,只有在指定可以分红的代理等级里的代理才有构建agentRelationship的必要，并且不是一级代理
		//todo:无限下级代理所累计的金额可以从日报表里拿
		//if monthReport.AgentLevel <= lowDividendAgentLevel && monthReport.AgentLevel != 1 {
		if monthReport.AgentLevel > 1 && monthReport.AgentLevel <= lowDividendAgentLevel {
			//获得上一级的UID
			monthReport.UpAgentUid = upAgentUid
			//构建代理关系等级，构建上级代理的下级活跃代理
			upAgentLevel := monthReport.AgentLevel - 1
			if agentRelationship[upAgentLevel] == nil {
				agentRelationship[upAgentLevel] = make(map[uint64][]uint64, 0)
			}
			if _, ok := agentRelationship[upAgentLevel][upAgentUid]; !ok {
				agentRelationship[upAgentLevel][upAgentUid] = make([]uint64, 0)
			}
			agentRelationship[upAgentLevel][upAgentUid] = append(agentRelationship[upAgentLevel][upAgentUid], selfUid)
			//}
		}
		if _, ok := agentLvUid[agentPath.Level]; !ok {
			agentLvUid[agentPath.Level] = make([]uint64, 0)
		}
		agentLvUid[agentPath.Level] = append(agentLvUid[agentPath.Level], agentPath.Uid)
		//将在指定下发月分红的代理等级里的用户加入进来
		if agentPath.Level <= lowDividendAgentLevel {
			lowAgentsUid, err := dao2.AgentPathDaoEntity.GetAllLowAgentUidByPath(agentPath.Path)
			if err != nil {
				return
			}
			//遍历自己的所有下级代理，计算自己下级代理的活跃人数
			for _, uid := range lowAgentsUid {
				//只有在日报表里面出现的用户才有可能是活跃用户
				if _, ok := monthReportsMap[uid]; !ok {
					continue

				}
				if monthReport.IsActivityUser {
					monthReport.ActivityAgentNum += 1
				}
			}
			//将在指定下发月分红的代理等级里的用户加入进来
			dividendMonthReports = append(dividendMonthReports, monthReport)
		}
	}

	//从下级往上级遍历,计算所有代理的盈利值
	for agentLv := len(agentRelationship) + 1; agentLv > 1; agentLv-- {
		lv := uint32(agentLv)
		//从最后一级的uid开始遍历
		for _, uid := range agentLvUid[lv] {
			//最后一级代理没有下级代理
			if agentLv != len(agentRelationship)+1 {
				monthReportsMap[uid].AgentsProfit = 0
				monthReportsMap[uid].ResultProfit = monthReportsMap[uid].Profit - monthReportsMap[uid].AgentsProfit
				continue
			}
			//	拿到自己所有下级代理的uid
			for _, lowerAgentUid := range agentRelationship[lv][uid] {
				monthReportsMap[uid].AgentsProfit += monthReportsMap[lowerAgentUid].ResultProfit
			}

			monthReportsMap[uid].ResultProfit += monthReportsMap[uid].Profit - monthReportsMap[uid].AgentsProfit
		}
	}
	//计算分红
	//代理还要考虑白名单的情况 如果平台给一级代理a设置在白名单里并且设置的亏损分红比例是 40% 则分红比例低于40%的时候是按40%算，如果高于40%则按高的算
	for _, monthReport := range dividendMonthReports {
		if _, ok := monthDividendRecordsMap[monthReport.Uid]; !ok {
			monthDividendRecordsMap[monthReport.Uid] = new(models2.MonthDividendRecord)
			monthDividendRecordsMap[monthReport.Uid].Uid = monthReport.Uid
			monthDividendRecordsMap[monthReport.Uid].Ctime = ctime
		}
		//profit := -monthReport.Profit
		profit := monthReport.ResultProfit
		//先获取自己分红百分比
		ratio := getDividendRation(profit, agentWhite, monthReport, dividendConf, monthDividendWhiteListMap)
		selfDividend := float64(profit) * ratio
		monthDividendRecordsMap[monthReport.Uid].SelfDividend = eosplus.QuantityFloat64ToInt64(selfDividend)

	}
	//计算自己要分给下一级代理的月分红
	for uid, dividendRecord := range monthDividendRecordsMap {
		var agentDividend int64 = 0
		//计算自己要分给代理的分红
		level := monthReportsMap[uid].AgentLevel

		//分红配置表里的最后一级代理是不需要给她下一级代理分红的
		if level < lowDividendAgentLevel {
			lowAgentIds := agentRelationship[level][uid]
			for _, lowAgentId := range lowAgentIds {
				agentDividend = agentDividend + monthDividendRecordsMap[lowAgentId].SelfDividend
			}
		}

		dividendRecord.AgentDividend = agentDividend
		dividendRecord.ResultDividend = dividendRecord.SelfDividend - agentDividend
		dividendRecord.ReceiveStatus = dao3.Received
		dividendRecord.PayStatus = dao3.Pay
		dividendRecord.Level = level

		monthDividendRecords = append(monthDividendRecords, dividendRecord)

	}
	err = dao3.GameUserMonthReportDaoEntity.Insert(monthReports)
	if err != nil {
		return
	}
	//遍历由最上级代理到要分红的最下级代理,设置月分红的领取状态
	for agentLv := 1; uint32(agentLv) < lowDividendAgentLevel; agentLv++ {
		for selfUid, lowAgentUids := range agentRelationship[uint32(agentLv)] {
			selfRecord := monthDividendRecordsMap[selfUid]
			//不是1级代理的属于自己的月分红金额是在上级代理的可提现金额为正数的时候加到 CanWithDrawAmount 上的,所以如果不是一级代理在计算自己的addCanWithDrawAmount只要发放给下级代理的月分红
			//当上级代理充钱了以后还要让自己的下级代理加上自己应得的的月分红
			addCanWithDrawAmount := selfRecord.AgentDividend
			if agentLv == 1 {
				addCanWithDrawAmount = selfRecord.ResultDividend
			}
			canWithDraw, err := dao2.AgentDaoEntity.UpdateCanWithdraw(selfRecord.Uid, addCanWithDrawAmount, lock.Options{
				LockTimeout: time.Second * time.Duration(common3.DefaultAgentLockTimeout),
				RetryCount:  common3.DefaultAgentRetryCount,
				RetryDelay:  time.Duration(common3.DefaultAgentRetryDelay) * time.Millisecond,
			})
			//canWithDraw, err := dao2.AgentDaoEntity.UpdateCanWithdraw(selfRecord.Uid, addCanWithDrawAmount)
			if err != nil {
				return
			}
			//如果上级代理的可提现余额小于0则将下级代理的接受状态设置为 NoReceived (上级未发放)状态,1级代理不可能出现未发放状态 并且还要冻结账号,并且将自己的支付状态设置成未支付
			if canWithDraw < 0 {
				selfRecord.PayStatus = dao3.NoPay
				frozenUser(selfRecord.Uid)
				for _, lowAgentUid := range lowAgentUids {
					lowAgentRecord := monthDividendRecordsMap[lowAgentUid]
					lowAgentRecord.ReceiveStatus = dao3.NoReceived
				}
			} else {
				//	如果上级代理可提现余额大等于0，将直属下级代理的canWithDraw字段加上他自己应得的月分红
				for _, lowAgentUid := range lowAgentUids {
					lowAgentRecord := monthDividendRecordsMap[lowAgentUid]
					_, err := dao2.AgentDaoEntity.UpdateCanWithdraw(selfRecord.Uid, lowAgentRecord.SelfDividend, lock.Options{
						LockTimeout: time.Second * time.Duration(common3.DefaultAgentLockTimeout),
						RetryCount:  common3.DefaultAgentRetryCount,
						RetryDelay:  time.Duration(common3.DefaultAgentRetryDelay) * time.Millisecond,
					})
					//_, err := dao2.AgentDaoEntity.UpdateCanWithdraw(lowAgentRecord.Uid, lowAgentRecord.SelfDividend)

					if err != nil {
						return
					}
				}
			}
		}
	}

	//for _, record := range monthDividendRecords {
	//	}
	//	//因为最后一级代理的最终分红不可能小于0,所以这里就把最后一级的代理给剔除了
	//	if record.ResultDividend < 0 {
	//		//如果不够下发自己的领取状态设置成NeedRecharge
	//		record.ReceiveStatus = dao3.NeedRecharge
	//		frozenUser(record.Uid)
	//		//获取自己直属下级代理的uid
	//		agentLevel := monthReportsMap[record.Uid].AgentLevel
	//		lowAgentUids := agentRelationship[agentLevel][record.Uid]
	//		for _, lowAgentUid := range lowAgentUids {
	//			//只有在自己代理的状态不为NeedRecharge才将代理状态设置成CanNotReceive
	//			if _, ok := monthDividendRecordsMap[lowAgentUid]; ok && monthDividendRecordsMap[lowAgentUid].ReceiveStatus != dao3.NeedRecharge {
	//				monthDividendRecordsMap[lowAgentUid].ReceiveStatus = dao3.CanNotReceive
	//			}
	//		}
	//		continue
	//	}
	//	//	如果他这个月的月分红够发给下级代理则去看一下之前的月分红记录是否有赊账的状态，有的话下级代理也要设置成等待上级分红的状态
	//	beforeRecords, err := dao3.MonthDividendRecordDaoEntity.FindByUid(record.Uid)
	//	if err != nil {
	//		return
	//	}
	//	//查找之前的月分红情况如果有一条记录的ReceiveStatus是NeedRecharge状态则让其所有下级代理都设置为不可领取状态
	//	for _, beforeRecord := range beforeRecords {
	//		if beforeRecord.ReceiveStatus == dao3.NeedRecharge {
	//			//获取自己直属下级代理的uid
	//			agentLevel := monthReportsMap[record.Uid].AgentLevel
	//			lowAgentUids := agentRelationship[agentLevel][record.Uid]
	//			for _, lowAgentUid := range lowAgentUids {
	//				//只有在自己代理的状态不为NeedRecharge才将代理状态设置成CanNotReceive
	//				if _, ok := monthDividendRecordsMap[lowAgentUid]; ok && monthDividendRecordsMap[lowAgentUid].ReceiveStatus != dao3.NeedRecharge {
	//					monthDividendRecordsMap[lowAgentUid].ReceiveStatus = dao3.CanNotReceive
	//				}
	//			}
	//			break
	//		}
	//	}
	//}

	err = dao3.MonthDividendRecordDaoEntity.Insert(monthDividendRecords)
	if err != nil {
		return
	}
}

func frozenUser(uid uint64) {
	// 冻结EUSD账号
	errCode2 := eosplus.EosPlusAPI.Wealth.Lock(uid)
	errCode := controllers.ERROR_CODE(errCode2)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		common.LogFuncError(fmt.Sprintf("冻结EUSD账号时出错 errCode=%v", errCode))
		return
	}

	// 冻结承兑OTC账号
	errCode2 = eosplus.EosPlusAPI.Otc.Lock(uid)
	errCode = controllers.ERROR_CODE(errCode2)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		common.LogFuncError(fmt.Sprintf("冻结承兑OTC账号时出错 errCode=%v", errCode))
		return
	}

	// 冻结USDT账号
	err := usdtDao.AccountDaoEntity.UpdateStatus(uid, usdtDao.STATUS_LOCKED)
	if err != nil {
		common.LogFuncError(fmt.Sprintf("冻结USDT账号时出错 errCode=%v", errCode))
		return
	}

	//推送通知
	go func() {
		title := "您的账号被冻结了!"
		content := "您的账号因月分红资金不足以分给下级代理被冻结了，欲解冻账号请尽快充值。如有疑问，请联系客服。"
		p := new(uemng_plus.UPushPlus)
		p.PushSysNotice(uid, content, title)
		_, _ = otcDao.SystemNotificationdDaoEntity.InsertSystemNotification("system", content, uid)
	}()
}
func getLastMonthBeginAndEndDay() (beginDay, endDay int64) {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start := thisMonth.AddDate(0, -1, 0).Format(DATE_FORMAT)
	end := thisMonth.AddDate(0, 0, -1).Format(DATE_FORMAT)
	beginDayInt, _ := strconv.ParseInt(start, 10, 64)
	endDayInt, _ := strconv.ParseInt(end, 10, 64)
	beginDay = int64(beginDayInt)
	endDay = int64(endDayInt)

	return
}

//亏损金额达到的等级和活跃用户达到的等级取小以后 再来和白名单取大 得出最终比例
func getDividendRation(profit int64, agentWhite map[uint64]uint32, monthReport *models2.GameUserMonthReport, dividendConf map[int][]models.MonthDividendPositionConf, monthDividendWhiteListMap map[uint32]*models.MonthDividendWhiteList) (ratio float64) {
	var lossPosition int32
	level := int(monthReport.AgentLevel)
	for i, dConf := range dividendConf[level] {
		if i < len(dividendConf[level])-1 {
			//从低的往上找，如果找到满足的lossPosition就break
			if dConf.Min <= profit && profit < dConf.Max {
				lossPosition = dConf.Position
				break
			}
		} else {
			//因为配置的最高级的Max是无限大，所以这里只要比较Min就好
			//如果这个都不满足则表示亏损额没有达到分红的要求
			if dConf.Min <= profit {
				lossPosition = dConf.Position
			}
		}
	}
	var activityNumPosition int32
	for _, dConf := range dividendConf[level] {
		//从低的往上找，如果一直小于则会到达最高等级，如果遇到中间不满足的情况则会break
		if dConf.ActivityNum < monthReport.ActivityAgentNum {
			activityNumPosition = dConf.Position
		} else {
			break
		}
	}
	//亏损额满足的等级和活跃人数满足的等级里取小
	position := lossPosition
	if activityNumPosition < position {
		position = activityNumPosition
	}
	var lossRatio int32 = 0
	for _, dConf := range dividendConf[level] {
		if dConf.Position == position {
			lossRatio = dConf.DividendRatio
			break
		}
	}

	//上面取出来的position和白名单的等级取大
	whiteListPosition := agentWhite[monthReport.Uid]
	whiteListRatio := monthDividendWhiteListMap[whiteListPosition].DividendRatio

	ratio = float64(whiteListRatio) / 100
	if lossRatio > whiteListRatio {
		ratio = float64(lossRatio) / 100
	}

	return
}
