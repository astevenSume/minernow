package controllers

import (
	"common"
	"errors"
	"github.com/bsm/redis-lock"
	"math"
	common3 "otc/common"
	"otc_error"
	"sync"
	"time"
	admindao "utils/admin/dao"
	adminmodels "utils/admin/models"
	agentdao "utils/agent/dao"
	agentmodels "utils/agent/models"
	common2 "utils/common"
	"utils/game/dao/gameapi"
	reportdao "utils/report/dao"
	reportmodels "utils/report/models"
)

const (
	DailySalaryCalcPerPage = 100
)

var wg sync.WaitGroup

//日工资定时任务，计算的是昨天的日工资，包含明细拉取，工资计算两部分
func TaskDailySalary() {
	wg.Add(2) //几个协程
	timestamp := time.Now().Unix() - common.DaySeconds
	//彩票日工资
	go common.SafeRun(func() {
		defer wg.Done()
		TaskRgDailySalary(gameapi.GAME_CHANNEL_RG, timestamp)
	})()

	//开元棋牌日工资
	go common.SafeRun(func() {
		defer wg.Done()
		TaskKyDailySalary(gameapi.GAME_CHANNEL_KY, timestamp)
	})()

	wg.Wait()
}

//Ag日工资定时任务(ag数据有延迟,放到凌晨4点)，计算的是昨天的日工资，包含明细拉取，工资计算两部分
func TaskAgDailySalary() {
	timestamp := time.Now().Unix() - common.DaySeconds
	TaskRgDailySalary(gameapi.GAME_CHANNEL_AG, timestamp)
}

//彩票、AG日工资计算模型，包含明细拉取，工资计算两部分
func TaskRgDailySalary(channelId uint32, timestamp int64) {
	//拉取游戏记录
	err := TaskGameRecords(channelId, timestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	//日工资计算
	err = TaskRgCalcSalary(channelId, timestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

//开元棋牌日工资计算模型，包含明细拉取，工资计算两部分
func TaskKyDailySalary(channelId uint32, timestamp int64) {
	//拉取游戏记录
	err := TaskGameRecords(channelId, timestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	//日工资计算
	err = TaskKyCalcSalary(channelId, timestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

//游戏明细数据拉取
func TaskGameRecords(channelId uint32, timestamp int64) (err error) {
	api := gameapi.GetApi(channelId)
	if api == nil {
		common.LogFuncError("get channelId[%v] game api fail", channelId)
		err = errors.New("game api fail")
		return
	}
	zeroTimestamp := common.GetZeroTime(timestamp)

	//应用白名单
	appsWhite, err := admindao.AppWhiteListDaoEntity.AllByChannel(channelId)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	//timestamp当天游戏明细数据拉取及投注统计
	err = api.DayBetRecords(zeroTimestamp, appsWhite)
	if err != nil {
		common.LogFuncError("error:%v", err)
		api.DelBetRecord(zeroTimestamp)
		return
	}
	return
}

//彩票日工资计算 包括三部分：回退数据、游戏投注量统计、日工资计算
func TaskRgCalcSalary(channelId uint32, timestamp int64) (err error) {
	api := gameapi.GetApi(channelId)
	if api == nil {
		common.LogFuncError("get channelId[%v] game api fail", channelId)
		return
	}
	zeroTimestamp := common.GetZeroTime(timestamp)

	//回退数据(重算才需要回退)
	err = rollBackSalary(channelId, zeroTimestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	var maxLevel uint32
	var mapUser map[uint64]*GameUserDaily
	var mapLevel map[uint32][]*GameUserDaily
	//个人投注信息统计，有效投注量统计=本人有效投注量+无限下级有效投注量
	maxLevel, mapUser, mapLevel, err = calcAllUserBetInfo(api, zeroTimestamp)
	if err != nil {
		_ = reportdao.ReportGameUserDailyDaoEntity.DelByTimestamp(channelId, zeroTimestamp)
		return
	}

	//日工资计算
	err = calcRgUserSalary(channelId, maxLevel, zeroTimestamp, mapUser, mapLevel)
	if err != nil {
		_ = rollBackSalary(channelId, zeroTimestamp)
		return
	}

	return
}

//开元棋牌日工资计算 包括三部分：回退数据、游戏投注量统计、日工资计算
func TaskKyCalcSalary(channelId uint32, timestamp int64) (err error) {
	api := gameapi.GetApi(channelId)
	if api == nil {
		common.LogFuncError("get channelId[%v] game api fail", channelId)
		return
	}
	zeroTimestamp := common.GetZeroTime(timestamp)

	//回退数据(重算才需要回退)
	err = rollBackSalary(channelId, zeroTimestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	var maxLevel uint32
	var mapUser map[uint64]*GameUserDaily
	var mapLevel map[uint32][]*GameUserDaily
	//个人投注信息统计，有效投注量统计=本人有效投注量+无限下级有效投注量
	maxLevel, mapUser, mapLevel, err = calcAllUserBetInfo(api, zeroTimestamp)
	if err != nil {
		_ = reportdao.ReportGameUserDailyDaoEntity.DelByTimestamp(channelId, zeroTimestamp)
		return
	}

	//日工资计算
	err = calcKyUserSalary(channelId, maxLevel, zeroTimestamp, mapUser, mapLevel)
	if err != nil {
		_ = rollBackSalary(channelId, zeroTimestamp)
		return
	}

	return
}

//回滚累计日工资
func rollBackSalary(channelId uint32, timestamp int64) (err error) {
	//旧数据
	var total int
	total, err = reportdao.ReportGameUserDailyDaoEntity.QueryTotalByChannelId(channelId, timestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	if total == 0 {
		return
	}
	totalPage := total / DailySalaryCalcPerPage
	if total%DailySalaryCalcPerPage > 0 {
		totalPage += 1
	}

	for page := 1; page <= totalPage; page++ {
		var list []reportmodels.ReportGameUserDaily
		list, err = reportdao.ReportGameUserDailyDaoEntity.QueryByChannelId(channelId, timestamp, page, DailySalaryCalcPerPage)
		if err != nil {
			common.LogFuncError("error:%v", err)
			return
		}
		//减少累计工资
		for _, item := range list {
			err = agentdao.AgentDaoEntity.DecreaseBalance(item.Uid, item.Salary, lock.Options{
				LockTimeout: time.Second * time.Duration(common3.Cursvr.LockTimeout),
				RetryCount:  common3.Cursvr.RetryCount,
				RetryDelay:  time.Duration(common3.Cursvr.RetryDelay) * time.Millisecond,
			})
			if err != nil {
				common.LogFuncError("AddBalance fail,uid:%v,salary:%v", item.Uid, item.Salary)
			}
		}
	}
	_ = reportdao.ReportGameUserDailyDaoEntity.DelByTimestamp(channelId, timestamp)

	return
}

type GameUserDaily struct {
	Uid           uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Bet           int64  `orm:"column(bet)" json:"bet,omitempty"`
	ValidBet      int64  `orm:"column(valid_bet)" json:"valid_bet,omitempty"`
	TotalValidBet int64  `orm:"column(total_valid_bet)" json:"total_valid_bet,omitempty"`
	TotalBetNum   int32  `orm:"column(total_bet_num)" json:"total_bet_num,omitempty"`
	Profit        int64  `orm:"column(profit)" json:"profit,omitempty"`
	TotalProfit   int64  `orm:"column(total_profit)" json:"total_profit,omitempty"`
	Salary        int64  `orm:"column(salary)" json:"salary,omitempty"`
	TeamSalary    int64  `orm:"column(team_salary)" json:"team_salary,omitempty"`
	Level         uint32 `orm:"column(level)" json:"level,omitempty"`
	ParentUid     uint64 `orm:"column(parent_uid)" json:"parent_uid,omitempty"`
}

//有效投注量计算,有效投注量=本人投注量加无限下级投注量
func calcAllUserBetInfo(api gameapi.GameAPI, timestamp int64) (maxLevel uint32,
	mapUser map[uint64]*GameUserDaily, mapLevel map[uint32][]*GameUserDaily, err error) {
	total, err := api.GetTotalByTimestamp(timestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	if total == 0 {
		common.LogFuncInfo("calcAllUserBetInfo total=0")
		return
	}
	totalPage := total / DailySalaryCalcPerPage
	if total%DailySalaryCalcPerPage > 0 {
		totalPage += 1
	}

	mapUser = make(map[uint64]*GameUserDaily)
	mapLevel = make(map[uint32][]*GameUserDaily)
	for page := 1; page <= totalPage; page++ {
		var uIds []uint64
		var list []reportdao.BetInfo
		list, err = api.GetRecordByTimestamp(timestamp, page, DailySalaryCalcPerPage)
		if err != nil {
			common.LogFuncError("calcAllUserBetInfo error:%v", err)
			return
		}

		for _, item := range list {
			//本人投注量及盈亏统计
			if userDaily, ok := mapUser[item.Uid]; ok {
				userDaily.Profit = userDaily.Profit + item.Profit
				userDaily.Bet = userDaily.Bet + item.Bet
				userDaily.ValidBet = userDaily.ValidBet + item.ValidBet
				userDaily.TotalValidBet = userDaily.ValidBet
				userDaily.TotalProfit = userDaily.Profit
			} else {
				userDaily := new(GameUserDaily)
				userDaily.Uid = item.Uid
				userDaily.Profit = item.Profit
				userDaily.Bet = item.Bet
				userDaily.ValidBet = item.ValidBet
				userDaily.TotalValidBet = item.ValidBet
				userDaily.TotalProfit = item.Profit
				mapUser[item.Uid] = userDaily
				uIds = append(uIds, item.Uid)
			}
		}

		var paths []agentmodels.AgentPath
		paths, err = agentdao.AgentPathDaoEntity.QueryByUIds(uIds)
		if err != nil {
			common.LogFuncError("calcAllUserBetInfo error:%v", err)
			return
		}
		for _, item := range paths {
			if userDaily, ok := mapUser[item.Uid]; ok {
				userDaily.Level = item.Level
				userDaily.ParentUid = item.ParentUid
				mapLevel[userDaily.Level] = append(mapLevel[userDaily.Level], userDaily)
				if userDaily.Level > maxLevel {
					maxLevel = userDaily.Level
				}
			} else {
				common.LogFuncError("calcAllUserBetInfo uid[%v] not found", item.Uid)
			}
		}
	}

	var pUserDaily *GameUserDaily
	var ok bool
	var totalBetNum int32
	//累计有效投注量=本人投注量加无限下级投注量
	for level := maxLevel; level > 0; level-- {
		var uIds []uint64
		length := len(mapLevel[level])
		for i := 0; i < length; i++ {
			sonUserDaily := mapLevel[level][i]
			if sonUserDaily.ValidBet > 0 {
				sonUserDaily.TotalBetNum = sonUserDaily.TotalBetNum + 1
			}
			if sonUserDaily.ParentUid > 0 {
				//累加到上级
				if pUserDaily, ok = mapUser[sonUserDaily.ParentUid]; ok {
					pUserDaily.TotalProfit = pUserDaily.TotalProfit + sonUserDaily.Profit
					pUserDaily.TotalValidBet = pUserDaily.TotalValidBet + sonUserDaily.TotalValidBet
				} else {
					uIds = append(uIds, sonUserDaily.ParentUid)
					pUserDaily = new(GameUserDaily)
					pUserDaily.Uid = sonUserDaily.ParentUid
					pUserDaily.Level = level - 1
					pUserDaily.TotalProfit = sonUserDaily.Profit
					pUserDaily.TotalValidBet = sonUserDaily.TotalValidBet
					mapUser[pUserDaily.Uid] = pUserDaily
					mapLevel[pUserDaily.Level] = append(mapLevel[pUserDaily.Level], pUserDaily)
				}
				totalBetNum = sonUserDaily.TotalBetNum
				if pUserDaily.ValidBet > 0 {
					totalBetNum = totalBetNum + 1
				}
				pUserDaily.TotalBetNum = pUserDaily.TotalBetNum + totalBetNum
			}
		}

		var paths []agentmodels.AgentPath
		paths, err = agentdao.AgentPathDaoEntity.QueryByUIds(uIds)
		if err != nil {
			common.LogFuncError("error:%v", err)
			return
		}
		for _, item := range paths {
			if pUserDaily, ok := mapUser[item.Uid]; ok {
				pUserDaily.ParentUid = item.ParentUid
			} else {
				common.LogFuncError("uid[%v] not save", item.Uid)
			}
		}
	}

	return
}

//彩票日工资计算 用户日工资=累计有效投注量*0.01*0.5+无限下级分红
func calcRgUserSalary(channelId, maxLevel uint32, timestamp int64, mapUser map[uint64]*GameUserDaily, mapLevel map[uint32][]*GameUserDaily) (err error) {
	gameRgRate, err := common2.AppConfigMgr.Float(common2.GameRgRate)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	gameRgRateToUser, err := common2.AppConfigMgr.Float(common2.GameRgRateToUser)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	gameRgRateToUpParam, err := common2.AppConfigMgr.Float(common2.GameRgRateToUpParam)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	gameRgRateUpParam, err := common2.AppConfigMgr.Float(common2.GameRgRateUpParam)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	common.LogFuncInfo("gameRgRate:%v,gameRgRateToUser:%v,gameRgRateToUpParam:%v,gameRgRateUpParam:%v",
		gameRgRate, gameRgRateToUser, gameRgRateToUpParam, gameRgRateUpParam)

	var (
		userDailys   []reportmodels.ReportGameUserDaily
		pUid         uint64
		curTotalBets int64
		subLevel     uint32
		other, a1    float64
		//pwd          string
	)
	for level := maxLevel; level > 0; level-- {
		length := len(mapLevel[level])
		for i := 0; i < length; i++ {
			sonUserDaily := mapLevel[level][i]
			if level == 1 {
				sonUserDaily.Salary = sonUserDaily.Salary + int64(float64(sonUserDaily.TotalValidBet)*gameRgRate)
			} else {
				sonUserDaily.Salary = sonUserDaily.Salary + int64(float64(sonUserDaily.TotalValidBet)*gameRgRate*gameRgRateToUser)
			}
			sonUserDaily.TeamSalary = sonUserDaily.TeamSalary + sonUserDaily.Salary

			pUid = sonUserDaily.ParentUid
			pUserDaily, ok := mapUser[pUid]
			if ok {
				pUserDaily.TeamSalary = pUserDaily.TeamSalary + sonUserDaily.TeamSalary
			}

			curTotalBets = sonUserDaily.TotalValidBet
			for {
				if pUid == 0 {
					break
				}
				parentUserDaily, ok := mapUser[pUid]
				if !ok {
					common.LogFuncError("cannot find pUid[%v]", pUid)
					break
				}

				subLevel = sonUserDaily.Level - parentUserDaily.Level
				other = float64(curTotalBets) * gameRgRate * (1 - gameRgRateToUser) //归属上级工资
				a1 = other * gameRgRateToUpParam
				if parentUserDaily.Level == 1 { //一级代理
					if subLevel == 1 {
						parentUserDaily.Salary = parentUserDaily.Salary + int64(other)
						//common.LogFuncDebug("userSalary2:%v", userSalary)
					} else {
						//中间代理
						//归属于中间上级的工资(等比数列求和)
						sn := a1 * (1 - math.Pow(gameRgRateUpParam, float64(subLevel-1))) / (1 - gameRgRateUpParam)
						//归属一级代理
						parentUserDaily.Salary = parentUserDaily.Salary + int64(other-sn)
					}
				} else { //中间代理
					parentUserDaily.Salary = parentUserDaily.Salary + int64(a1*math.Pow(gameRgRateUpParam, float64(subLevel-1)))
					//common.LogFuncDebug("userSalary4:%v", userSalary)
				}
				pUid = parentUserDaily.ParentUid
			}
		}
	}

	var level uint32
	for level = 1; level <= maxLevel; level++ {
		length := len(mapLevel[level])
		for i := 0; i < length; i++ {
			item := mapLevel[level][i]
			userDailys = append(userDailys, reportmodels.ReportGameUserDaily{
				Uid:           item.Uid,
				ChannelId:     channelId,
				Ctime:         timestamp,
				Status:        reportdao.ReportGameUserDailyGet,
				Profit:        item.Profit,
				Bet:           item.Bet,
				ValidBet:      item.ValidBet,
				TotalValidBet: item.TotalValidBet,
				TotalProfit:   item.TotalProfit,
				Salary:        item.Salary,
				TeamSalary:    item.TeamSalary,
				TotalBetNum:   item.TotalBetNum,
				Level:         item.Level,
				PUid:          item.ParentUid,
				//Pwd:           pwd,
			})

			if len(userDailys) >= DailySalaryCalcPerPage {
				//个人日报入库
				err = reportdao.ReportGameUserDailyDaoEntity.InsertMul(channelId, timestamp, userDailys)
				if err != nil {
					common.LogFuncError("error:%v", err)
					return
				}
				for _, item := range userDailys {
					result := AddBalance(item.Uid, item.Salary)
					if result != nil {
						common.LogFuncError("error:%v, uid:%v, salary:%v", result, item.Uid, item.Salary)
					}
				}
			}
		}
	}

	//个人日报入库
	err = reportdao.ReportGameUserDailyDaoEntity.InsertMul(channelId, timestamp, userDailys)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	for _, item := range userDailys {
		result := AddBalance(item.Uid, item.Salary)
		if result != nil {
			common.LogFuncError("error:%v, uid:%v, salary:%v", result, item.Uid, item.Salary)
		}
	}

	return
}

//开元棋牌日工资计算 用户日工资=本人有效投注量* 利润率 * 等级乘数+无限下级分红((下级有效投注量)*利润率*(本人等级乘数-下级等级乘数))
//如果上线(D)的等级乘数大于下线(E)的等级乘数，那么E分配给D的日工资为(E有效投注量)*利润率*(D等级乘数-E等级乘数),如果D的上线C等级乘数
// 小于等于D,那么E分配给C的日工资为0，直到在C之上的上线A等级乘数大于D的等级乘数，那么E分配给A的日工资为(E有效投注量)*利润率*(A等级乘数-D等级乘数)
func calcKyUserSalary(channelId, maxLevel uint32, timestamp int64, mapUser map[uint64]*GameUserDaily, mapLevel map[uint32][]*GameUserDaily) (err error) {
	appChannel, err := admindao.AppChannelDaoEntity.QueryById(channelId)
	if err != nil || appChannel == nil {
		common.LogFuncError("get channelId[%v] cfg fail", channelId)
		return
	}

	//投注量等级配置
	commissionConfig := NewCommissionConfig()
	err = commissionConfig.Load()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	//白名单配置
	whiteListCommissionConfig := NewWhiteListCommissionConfig()
	err = whiteListCommissionConfig.Load()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	var mapWhiteList agentdao.MapWhiteList
	mapWhiteList, err = agentdao.AgentPathDaoEntity.GetAllWhiteList()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	var (
		userDailys     []reportmodels.ReportGameUserDaily
		curCommission  int32
		curPrecision   int32
		nextCommission int32
		diffCommission int32
		pUid           uint64
		//pwd            string
	)
	for level := maxLevel; level > 0; level-- {
		length := len(mapLevel[level])
		for i := 0; i < length; i++ {
			sonUserDaily := mapLevel[level][i]
			sonValidBet := sonUserDaily.ValidBet
			curCommission, curPrecision, err = GetCommission(mapWhiteList[sonUserDaily.Uid], sonUserDaily.TotalValidBet, commissionConfig, whiteListCommissionConfig)
			if err != nil {
				return
			}
			sonUserDaily.Salary = sonUserDaily.Salary + int64(appChannel.ProfitRate*curCommission)*sonValidBet/int64(admindao.AppChannelPrecision*curPrecision)
			sonUserDaily.TeamSalary = sonUserDaily.TeamSalary + sonUserDaily.Salary

			pUid = sonUserDaily.ParentUid
			pUserDaily, ok := mapUser[pUid]
			if ok {
				pUserDaily.TeamSalary = pUserDaily.TeamSalary + sonUserDaily.TeamSalary
			}

			//计算分配给无限上线的日工资
			for {
				if pUid == 0 {
					break
				}
				parentUserDaily, ok := mapUser[pUid]
				if !ok {
					common.LogFuncError("cannot find pUid[%v]", pUid)
					break
				}

				nextCommission, _, err = GetCommission(mapWhiteList[parentUserDaily.Uid], parentUserDaily.TotalValidBet, commissionConfig, whiteListCommissionConfig)
				if err != nil {
					return
				}

				if nextCommission > curCommission {
					diffCommission = nextCommission - curCommission
					parentUserDaily.Salary = parentUserDaily.Salary + int64(appChannel.ProfitRate*diffCommission)*sonValidBet/int64(admindao.AppChannelPrecision*curPrecision)
					curCommission = nextCommission
				}
				pUid = parentUserDaily.ParentUid
			}
		}
	}

	var level uint32
	for level = 1; level <= maxLevel; level++ {
		length := len(mapLevel[level])
		for i := 0; i < length; i++ {
			item := mapLevel[level][i]
			userDailys = append(userDailys, reportmodels.ReportGameUserDaily{
				Uid:           item.Uid,
				ChannelId:     channelId,
				Ctime:         timestamp,
				Status:        reportdao.ReportGameUserDailyGet,
				Profit:        item.Profit,
				Bet:           item.Bet,
				ValidBet:      item.ValidBet,
				TotalValidBet: item.TotalValidBet,
				TotalProfit:   item.TotalProfit,
				Salary:        item.Salary,
				TeamSalary:    item.TeamSalary,
				TotalBetNum:   item.TotalBetNum,
				Level:         item.Level,
				PUid:          item.ParentUid,
				//Pwd:           pwd,
			})

			if len(userDailys) >= DailySalaryCalcPerPage {
				//个人日报入库
				err = reportdao.ReportGameUserDailyDaoEntity.InsertMul(channelId, timestamp, userDailys)
				if err != nil {
					common.LogFuncError("error:%v", err)
					return
				}
				for _, item := range userDailys {
					result := AddBalance(item.Uid, item.Salary)
					if result != nil {
						common.LogFuncError("error:%v, uid:%v, salary:%v", result, item.Uid, item.Salary)
					}
				}
			}
		}
	}

	//个人日报入库
	err = reportdao.ReportGameUserDailyDaoEntity.InsertMul(channelId, timestamp, userDailys)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	for _, item := range userDailys {
		result := AddBalance(item.Uid, item.Salary)
		if result != nil {
			common.LogFuncError("error:%v, uid:%v, salary:%v", result, item.Uid, item.Salary)
		}
	}

	return
}

func GetCommission(whiteListId uint32, totalValidBet int64, commissionConfig *CommissionConfig, whiteConfig *WhiteListCommissionConfig) (commission int32, precision int32, err error) {
	var (
		whiteCommission int32
		whitePrecision  int32
		isFind          bool
	)
	if whiteListId > 0 {
		whiteCommission, whitePrecision, isFind = whiteConfig.Get(whiteListId)
		if !isFind {
			common.LogFuncError("AgentWhiteList[%d] not exist", whiteListId)
		}
	}

	commission, precision = commissionConfig.GetCommission(uint64(totalValidBet))
	if whiteCommission > commission {
		commission = whiteCommission
		precision = whitePrecision
	}
	if precision == 0 {
		common.LogFuncError("err:curPrecision = 0 ")
		err = errors.New("precision can not equal zero")
		return
	}

	return
}

func AddBalance(uid uint64, amount int64) (err error) {
	var before int64
	var after int64
	before, after, err = agentdao.AgentDaoEntity.AddBalance(uid, amount, lock.Options{
		LockTimeout: time.Second * time.Duration(common3.Cursvr.LockTimeout),
		RetryCount:  common3.Cursvr.RetryCount,
		RetryDelay:  time.Duration(common3.Cursvr.RetryDelay) * time.Millisecond,
	})
	if err != nil {
		common.LogFuncError("error:%v, uid:%v, salary:%v", err, uid, amount)
		return
	}
	if before < 0 && after >= 0 {
		//解冻发放下级月分红
		errCode := UnFrozenUser(uid)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			common.LogFuncError("UnFrozenUser errCode:%v, uid:%v, before:%v,after:%v", errCode, uid, before, after)
		}

		err = DividendToLowAgents(uid)
		if err != nil {
			common.LogFuncError("error:%v, uid:%v, salary:%v", err, uid, amount)
			return
		}
	}

	return
}

//日工资投注量等级配置
type CommissionConfig struct {
	list []adminmodels.Commissionrates
}

func NewCommissionConfig() *CommissionConfig {
	return &CommissionConfig{}
}
func (c *CommissionConfig) Load() (err error) {
	c.list, err = admindao.CommissionRateDaoEntity.All()
	if err != nil {
		return
	}

	/*for _, item := range c.list {
		item.Min = eosplus.QuantityFloat64ToUint64(float64(item.Min))
		item.Max = eosplus.QuantityFloat64ToUint64(float64(item.Max))
	}*/

	return
}
func (c *CommissionConfig) GetCommission(bet uint64) (commission int32, precision int32) {
	for _, item := range c.list {
		if bet > item.Min && (bet <= item.Max || item.Max == 0) {
			commission = item.Commission
			precision = item.Precision
			return
		}
	}

	return
}

//日工资白名单
type WhiteListCommission struct {
	Commission int32
	Precision  int32
}
type WhiteListCommissionConfig struct {
	items map[uint32]WhiteListCommission
}

func NewWhiteListCommissionConfig() *WhiteListCommissionConfig {
	return &WhiteListCommissionConfig{}
}
func (c *WhiteListCommissionConfig) Load() (err error) {
	var list []adminmodels.AgentWhiteList
	list, err = admindao.AgentWhiteListDaoEntity.All()
	if err != nil {
		return
	}

	c.items = make(map[uint32]WhiteListCommission)
	for _, l := range list {
		c.items[l.Id] = WhiteListCommission{
			Commission: l.Commission,
			Precision:  l.Precision,
		}
	}

	return
}

func (c *WhiteListCommissionConfig) Get(id uint32) (commission int32, precision int32, ok bool) {
	var v WhiteListCommission
	if v, ok = c.items[id]; ok {
		commission = v.Commission
		precision = v.Precision
		return
	}

	return
}
