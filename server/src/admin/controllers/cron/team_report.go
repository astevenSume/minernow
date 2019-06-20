package cron

import (
	"common"
	"math"
	"pb"
	"time"
	agentdao "utils/agent/dao"
	"utils/game/dao/gameapi"
	reportdao "utils/report/dao"
	reportmodels "utils/report/models"
)

type ReportTeamGameTransferDaily struct {
	Uid          uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	TeamRecharge int64  `orm:"column(team_recharge)" json:"team_recharge,omitempty"`
	TeamWithdraw int64  `orm:"column(team_withdraw)" json:"team_withdraw,omitempty"`
}

type UserAgent struct {
	ParentUid uint64 `orm:"column(parent_uid)" json:"parent_uid,omitempty"`
	Level     uint32 `orm:"column(level)" json:"level,omitempty"`
}

func TaskGameTransferDaily() {
	timestamp := common.GetZeroTime(time.Now().Unix() - common.DaySeconds)
	err := reportdao.ReportTeamGameTransferDailyDaoEntity.DelByTimestamp(timestamp)
	if err != nil {
		return
	}

	mapChannelUser := make(map[uint32]map[uint64]*ReportTeamGameTransferDaily)
	mapChannelLevel := make(map[uint32]map[uint32][]*ReportTeamGameTransferDaily)
	for channelId := gameapi.GAME_CHANNEL_KY; channelId < gameapi.GAME_CHANNEL_MAX; channelId++ {
		mapChannelUser[channelId] = make(map[uint64]*ReportTeamGameTransferDaily)
		mapChannelLevel[channelId] = make(map[uint32][]*ReportTeamGameTransferDaily)
	}
	mapUserAgent := make(map[uint64]*UserAgent)

	total, err := reportdao.ReportGameTransferDailyDaoEntity.QueryTotalByTimestamp(timestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	if total == 0 {
		return
	}
	totalPage := int(math.Ceil(float64(total) / float64(DailyQueryPerPage)))

	//个人日报
	var (
		ok                    bool
		maxLevel              uint32
		path                  pb.AgentPath
		userAgent, pUserAgent *UserAgent
		userDaily, pUserDaily *ReportTeamGameTransferDaily
		mapUser               map[uint64]*ReportTeamGameTransferDaily
		userDailys            []reportmodels.ReportTeamGameTransferDaily
	)

	for page := 1; page <= totalPage; page++ {
		list, err := reportdao.ReportGameTransferDailyDaoEntity.QueryByTimestamp(timestamp, page, DailyQueryPerPage)
		if err != nil {
			return
		}

		for _, item := range list {
			//common.LogFuncDebug("item:%+v", item)
			if mapUser, ok = mapChannelUser[item.ChannelId]; ok {
				userDaily = new(ReportTeamGameTransferDaily)
				userDaily.Uid = item.Uid
				userDaily.TeamRecharge = item.Recharge
				userDaily.TeamWithdraw = item.Withdraw
				mapUser[item.Uid] = userDaily

				if _, ok := mapUserAgent[item.Uid]; !ok {
					path, err = agentdao.AgentPathDaoEntity.GetRedisAgentPathByUid(item.Uid)
					if err != nil {
						common.LogFuncError("cannot find uid[%v] userAgent err:%v", item.Uid, err)
					} else {
						userAgent = new(UserAgent)
						userAgent.Level = path.Level
						userAgent.ParentUid = path.PUId
						mapUserAgent[item.Uid] = userAgent

						mapChannelLevel[item.ChannelId][userAgent.Level] = append(mapChannelLevel[item.ChannelId][userAgent.Level], userDaily)
						if userAgent.Level > maxLevel {
							maxLevel = userAgent.Level
						}
					}
				}
			} else {
				common.LogFuncError("error: unknown channel[%d]", item.ChannelId)
			}
		}
	}
	common.LogFuncDebug("maxLevel:%v", maxLevel)

	//累加到上级
	for channelId := gameapi.GAME_CHANNEL_KY; channelId < gameapi.GAME_CHANNEL_MAX; channelId++ {
		for level := maxLevel; level > 0; level-- {
			length := len(mapChannelLevel[channelId][level])
			for i := 0; i < length; i++ {
				userDaily = mapChannelLevel[channelId][level][i]
				userAgent, ok = mapUserAgent[userDaily.Uid]
				if !ok {
					common.LogFuncError("cannot find uid[%v] userAgent", userDaily.Uid)
					continue
				}

				if userAgent.ParentUid > 0 {
					//累加到上级
					if pUserDaily, ok = mapChannelUser[channelId][userAgent.ParentUid]; ok {
						pUserDaily.TeamWithdraw = pUserDaily.TeamWithdraw + userDaily.TeamWithdraw
						pUserDaily.TeamRecharge = pUserDaily.TeamRecharge + userDaily.TeamRecharge
					} else {
						pUserAgent, ok = mapUserAgent[userAgent.ParentUid]
						if !ok {
							path, err = agentdao.AgentPathDaoEntity.GetRedisAgentPathByUid(userAgent.ParentUid)
							if err != nil {
								common.LogFuncError("cannot find uid[%v] userAgent err:%v", userAgent.ParentUid, err)
								continue
							} else {
								pUserAgent = new(UserAgent)
								pUserAgent.Level = path.Level
								pUserAgent.ParentUid = path.PUId
								mapUserAgent[userAgent.ParentUid] = pUserAgent
							}
						}
						pUserDaily = new(ReportTeamGameTransferDaily)
						pUserDaily.Uid = userAgent.ParentUid
						pUserDaily.TeamRecharge = userDaily.TeamRecharge
						pUserDaily.TeamWithdraw = userDaily.TeamWithdraw
						mapChannelUser[channelId][pUserDaily.Uid] = pUserDaily
						mapChannelLevel[channelId][pUserAgent.Level] = append(mapChannelLevel[channelId][pUserAgent.Level], pUserDaily)
					}
				}
			}
		}
	}

	//入库
	for channelId := gameapi.GAME_CHANNEL_KY; channelId < gameapi.GAME_CHANNEL_MAX; channelId++ {
		for level := maxLevel; level > 0; level-- {
			length := len(mapChannelLevel[channelId][level])
			for i := 0; i < length; i++ {
				userDaily = mapChannelLevel[channelId][level][i]
				userDailys = append(userDailys, reportmodels.ReportTeamGameTransferDaily{
					Uid:          userDaily.Uid,
					TeamRecharge: userDaily.TeamRecharge,
					TeamWithdraw: userDaily.TeamWithdraw,
					ChannelId:    channelId,
					Ctime:        timestamp,
					Level:        level,
				})

				if len(userDailys) >= DailyQueryPerPage {
					err = reportdao.ReportTeamGameTransferDailyDaoEntity.InsertMul(timestamp, userDailys)
					if err != nil {
						common.LogFuncError("error:%v", err)
						return
					}
				}
			}
		}
	}
	err = reportdao.ReportTeamGameTransferDailyDaoEntity.InsertMul(timestamp, userDailys)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
}

type ReportEusdDaily struct {
	Uid       uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Buy       int64  `orm:"column(buy)" json:"buy,omitempty"`
	Sell      int64  `orm:"column(sell)" json:"sell,omitempty"`
	ParentUid uint64 `orm:"column(parent_uid)" json:"parent_uid,omitempty"`
	Level     uint32 `orm:"column(level)" json:"level,omitempty"`
}

func TaskTeamDaily() {
	timestamp := common.GetZeroTime(time.Now().Unix() - common.DaySeconds)
	err := reportdao.ReportTeamDailyDaoEntity.DelByTimestamp(timestamp)
	if err != nil {
		return
	}

	total, err := reportdao.ReportEusdDailyDaoEntity.QueryTotalByTimestamp(timestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	if total == 0 {
		return
	}
	totalPage := int(math.Ceil(float64(total) / float64(DailyQueryPerPage)))

	var (
		ok                                bool
		maxLevel                          uint32
		agentPath                         pb.AgentPath
		reportEusdDaily, pReportEusdDaily *ReportEusdDaily
		reportTeamDailys                  []reportmodels.ReportTeamDaily
	)
	mapUser := make(map[uint64]*ReportEusdDaily)
	mapLevel := make(map[uint32][]*ReportEusdDaily)
	for page := 1; page <= totalPage; page++ {
		list, err := reportdao.ReportEusdDailyDaoEntity.QueryByTimestamp(timestamp, page, DailyQueryPerPage)
		if err != nil {
			return
		}

		for _, item := range list {
			//common.LogFuncDebug("item:%+v", item)
			reportEusdDaily = new(ReportEusdDaily)
			reportEusdDaily.Uid = item.Uid
			reportEusdDaily.Buy = item.Buy
			reportEusdDaily.Sell = item.Sell

			agentPath, err = agentdao.AgentPathDaoEntity.GetRedisAgentPathByUid(item.Uid)
			if err != nil {
				common.LogFuncError("cannot find uid[%v] userAgent err:%v", item.Uid, err)
			} else {
				reportEusdDaily.ParentUid = agentPath.PUId
				reportEusdDaily.Level = agentPath.Level
				if agentPath.Level > maxLevel {
					maxLevel = agentPath.Level
				}
			}
			mapUser[item.Uid] = reportEusdDaily
			mapLevel[reportEusdDaily.Level] = append(mapLevel[reportEusdDaily.Level], reportEusdDaily)
		}
	}

	//累加到上级
	for level := maxLevel; level > 0; level-- {
		length := len(mapLevel[level])
		for i := 0; i < length; i++ {
			reportEusdDaily = mapLevel[level][i]

			if reportEusdDaily.ParentUid > 0 {
				//累加到上级
				if pReportEusdDaily, ok = mapUser[reportEusdDaily.ParentUid]; ok {
					pReportEusdDaily.Buy = pReportEusdDaily.Buy + reportEusdDaily.Buy
					pReportEusdDaily.Sell = pReportEusdDaily.Sell + reportEusdDaily.Sell
				} else {
					pReportEusdDaily = new(ReportEusdDaily)
					pReportEusdDaily.Uid = reportEusdDaily.ParentUid
					pReportEusdDaily.Buy = reportEusdDaily.Buy
					pReportEusdDaily.Sell = reportEusdDaily.Sell

					agentPath, err = agentdao.AgentPathDaoEntity.GetRedisAgentPathByUid(pReportEusdDaily.Uid)
					if err != nil {
						common.LogFuncError("cannot find uid[%v] userAgent err:%v", pReportEusdDaily.Uid, err)
					} else {
						pReportEusdDaily.ParentUid = agentPath.PUId
						pReportEusdDaily.Level = agentPath.Level
					}
					mapUser[pReportEusdDaily.Uid] = pReportEusdDaily
					mapLevel[pReportEusdDaily.Level] = append(mapLevel[pReportEusdDaily.Level], pReportEusdDaily)
				}
			}
		}
	}

	//入库
	for _, reportEusdDaily = range mapUser {
		reportTeamDailys = append(reportTeamDailys, reportmodels.ReportTeamDaily{
			Uid:      reportEusdDaily.Uid,
			EusdBuy:  reportEusdDaily.Buy,
			EusdSell: reportEusdDaily.Sell,
			Ctime:    timestamp,
			Level:    reportEusdDaily.Level,
		})

		if len(reportTeamDailys) >= DailyQueryPerPage {
			err = reportdao.ReportTeamDailyDaoEntity.InsertMul(timestamp, reportTeamDailys)
			if err != nil {
				common.LogFuncError("error:%v", err)
				return
			}
		}
	}
	err = reportdao.ReportTeamDailyDaoEntity.InsertMul(timestamp, reportTeamDailys)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
}

type ReportCommission struct {
	Uid             uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	TeamWithdraw    int64  `orm:"column(team_withdraw)" json:"team_withdraw,omitempty"`
	TeamCanWithdraw int64  `orm:"column(team_can_withdraw)" json:"team_can_withdraw,omitempty"`
	ParentUid       uint64 `orm:"column(parent_uid)" json:"parent_uid,omitempty"`
	Level           uint32 `orm:"column(level)" json:"level,omitempty"`
}

func TaskReportCommission() {
	timestamp := common.GetZeroTime(time.Now().Unix() - common.DaySeconds)
	err := reportdao.ReportCommissionDaoEntity.Del()
	if err != nil {
		return
	}

	total, err := agentdao.AgentDaoEntity.Total()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	if total == 0 {
		return
	}
	totalPage := int(math.Ceil(float64(total) / float64(DailyQueryPerPage)))

	var (
		ok                        bool
		maxLevel                  uint32
		agentPath                 pb.AgentPath
		reportDaily, pReportDaily *ReportCommission
		reportTeamDailys          []reportmodels.ReportCommission
	)
	mapUser := make(map[uint64]*ReportCommission)
	mapLevel := make(map[uint32][]*ReportCommission)
	for page := 1; page <= totalPage; page++ {
		list, err := agentdao.AgentDaoEntity.QueryByPage(page, DailyQueryPerPage)
		if err != nil {
			return
		}

		for _, item := range list {
			//common.LogFuncDebug("item:%+v", item)
			reportDaily = new(ReportCommission)
			reportDaily.Uid = item.Uid
			reportDaily.TeamCanWithdraw = item.SumCanWithdraw
			reportDaily.TeamWithdraw = item.SumWithdraw

			agentPath, err = agentdao.AgentPathDaoEntity.GetRedisAgentPathByUid(item.Uid)
			if err != nil {
				common.LogFuncError("cannot find uid[%v] userAgent err:%v", item.Uid, err)
			} else {
				reportDaily.ParentUid = agentPath.PUId
				reportDaily.Level = agentPath.Level
				if agentPath.Level > maxLevel {
					maxLevel = agentPath.Level
				}
			}
			mapUser[item.Uid] = reportDaily
			mapLevel[reportDaily.Level] = append(mapLevel[reportDaily.Level], reportDaily)
		}
	}

	//累加到上级
	for level := maxLevel; level > 0; level-- {
		length := len(mapLevel[level])
		for i := 0; i < length; i++ {
			reportDaily = mapLevel[level][i]

			if reportDaily.ParentUid > 0 {
				//累加到上级
				if pReportDaily, ok = mapUser[reportDaily.ParentUid]; ok {
					pReportDaily.TeamWithdraw = pReportDaily.TeamWithdraw + reportDaily.TeamWithdraw
					pReportDaily.TeamCanWithdraw = pReportDaily.TeamCanWithdraw + reportDaily.TeamCanWithdraw
				} else {
					pReportDaily = new(ReportCommission)
					pReportDaily.Uid = reportDaily.ParentUid
					pReportDaily.TeamWithdraw = reportDaily.TeamWithdraw
					pReportDaily.TeamCanWithdraw = reportDaily.TeamCanWithdraw

					agentPath, err = agentdao.AgentPathDaoEntity.GetRedisAgentPathByUid(pReportDaily.Uid)
					if err != nil {
						common.LogFuncError("cannot find uid[%v] userAgent err:%v", pReportDaily.Uid, err)
					} else {
						pReportDaily.ParentUid = agentPath.PUId
						pReportDaily.Level = agentPath.Level
					}
					mapUser[pReportDaily.Uid] = pReportDaily
					mapLevel[pReportDaily.Level] = append(mapLevel[pReportDaily.Level], pReportDaily)
				}
			}
		}
	}

	//入库
	for _, reportDaily = range mapUser {
		if reportDaily.TeamCanWithdraw > 0 || reportDaily.TeamWithdraw > 0 {
			reportTeamDailys = append(reportTeamDailys, reportmodels.ReportCommission{
				Uid:             reportDaily.Uid,
				TeamCanWithdraw: reportDaily.TeamCanWithdraw,
				TeamWithdraw:    reportDaily.TeamWithdraw,
				Ctime:           timestamp,
				Level:           reportDaily.Level,
			})

			if len(reportTeamDailys) >= DailyQueryPerPage {
				err = reportdao.ReportCommissionDaoEntity.InsertMul(reportTeamDailys)
				if err != nil {
					common.LogFuncError("error:%v", err)
					return
				}
			}
		}
	}
	err = reportdao.ReportCommissionDaoEntity.InsertMul(reportTeamDailys)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
}
