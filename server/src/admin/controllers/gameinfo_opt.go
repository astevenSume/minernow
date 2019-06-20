package controllers

import (
	controllers "admin/controllers/errcode"
	"common"
	"eusd/eosplus"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"net/http"
	"strconv"
	"time"
	gamedao "utils/game/dao"
	"utils/game/dao/gameapi"
	"utils/otc/dao"
	reportdao "utils/report/dao"
	reportmodels "utils/report/models"
	tmpdao "utils/tmp/dao"
	"utils/tmp/models"
)

const TRANSFORM_MUN float64 = 100000000

type GameInfoOptController struct {
	BaseController
}

type ResultGameTransfer struct {
	Id           uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid          uint64 `orm:"column(uid)" json:"uid,omitempty"`
	ChannelId    uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	Account      string `orm:"column(account);size(50)" json:"account,omitempty"`
	TransferType uint32 `orm:"column(transfer_type)" json:"transfer_type,omitempty"`
	Order        string `orm:"column(order);size(50)" json:"order,omitempty"`
	GameOrder    string `orm:"column(game_order);size(50)" json:"game_order,omitempty"`
	CoinInteger  string `orm:"column(coin_integer)" json:"coin_integer,omitempty"`
	EusdInteger  string `orm:"column(eusd_integer)" json:"eusd_integer,omitempty"`
	Status       uint32 `orm:"column(status)" json:"status,omitempty"`
	Ctime        int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Desc         string `orm:"column(desc);size(512)" json:"desc,omitempty"`
	Step         string `orm:"column(step);size(256)" json:"step,omitempty"`
}

/**transfertype 查询 1 充值 2 提现 0充值提现记录 */
func (c *GameInfoOptController) QueryChargeRecord() {
	c.setOPAction(OPActionQueryGameCharge)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	uid := c.GetString(KEY_UID)
	mobile := c.GetString(KEY_MOBILE)
	//判断用户是否存在
	var rowNum int
	var err error

	if uid != "" {
		rowNum, err = dao.UserDaoEntity.UserIsPresence(uid)
	}

	if mobile != "" && rowNum == 0 {
		rowNum, err = dao.UserDaoEntity.NewMobileIsPresence(mobile)

	}

	if err != nil || rowNum == 0 {
		c.ErrorResponseAndLog(OPActionQueryGameCharge, controllers.ERROR_CODE_NO_USER, string(c.Ctx.Input.RequestBody))
	}
	if uid == "" {
		theUid, err := dao.UserDaoEntity.GetUidByMobile(mobile)
		if err != nil {
			if err == orm.ErrNoRows {
				c.ErrorResponseAndLog(OPActionQueryGameCharge, controllers.ERROR_CODE_NO_USER, string(c.Ctx.Input.RequestBody))
			}
			c.ErrorResponseAndLog(OPActionQueryGameCharge, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		uid = strconv.FormatInt(int64(theUid), 10)
	}

	//查询出记录
	transferType, err := c.GetInt(KEY_GAME_TRANSFER_TYPE, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionQueryGameCharge, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
	}
	page, _ := c.GetInt(KEY_PAGE, 1)
	limit, _ := c.GetInt(KEY_LIMIT, 10)

	total, data, err := gamedao.GameTransferDaoEntity.QueryPageTransferCharge(uid, transferType, page, limit)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
	}
	var resultData []ResultGameTransfer

	for _, v := range data {
		result := ResultGameTransfer{
			Id:           v.Id,
			Uid:          v.Uid,
			ChannelId:    v.ChannelId,
			Account:      v.Account,
			TransferType: v.TransferType,
			Order:        v.Order,
			GameOrder:    v.GameOrder,
			CoinInteger:  common.Init64DivisorToStr(v.CoinInteger, TRANSFORM_MUN),
			EusdInteger:  common.Init64DivisorToStr(v.EusdInteger, TRANSFORM_MUN),
			Status:       v.Status,
			Ctime:        v.Ctime,
			Desc:         v.Desc,
			Step:         v.Step,
		}
		resultData = append(resultData, result)
	}
	c.SuccessResponse(map[string]interface{}{
		KEY_LIST: resultData,
		KEY_META: PageInfo{
			Page:  page,
			Total: int(total),
			Limit: limit,
		},
	})

}

/*
 * 查询游戏数据统计报表
 */
func (c *GameInfoOptController) GetStatisticInfo() {
	c.setOPAction(OPActionGetGameReport)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	action := c.Ctx.Input.Param(KEY_KEY_ACTION)
	channelId := getChannelId(action)
	if channelId == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
	}

	start, _ := strconv.ParseInt(c.GetString(KEY_START_TIME), 10, 64)
	over, _ := strconv.ParseInt(c.GetString(KEY_END_TIME), 10, 64)

	page, _ := c.GetInt(KEY_PAGE, 1)
	limit, _ := c.GetInt(KEY_LIMIT, 10)
	total, data, err := reportdao.ReportStatisticGameAllDaoEntity.QueryGameReport(channelId, start, over, page, limit)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
	}

	var resultTotal []interface{}
	totalData, err := reportdao.ReportStatisticSumDaoEntity.QuerySumReport(channelId, start, over)
	logs.Debug("totalData", totalData, err)
	if totalData.ChannelId > 0 {
		if channelId == gameapi.GAME_CHANNEL_RG {
			monthDividend, err := reportdao.MonthDividendRecordDaoEntity.FindByRange(start, over)
			if err == nil {
				totalData.ChannelRgDividend = monthDividend.ResultDividend
			}
		}
		resultTotal = append(resultTotal, totalData)
	}

	c.SuccessResponse(map[string]interface{}{
		KEY_LIST:  data,
		KEY_TOTAL: resultTotal,
		KEY_META: PageInfo{
			Page:  page,
			Total: int(total),
			Limit: limit,
		},
	})

}

//查询每个report表
func (c *GameInfoOptController) GetReportInfo() {
	c.setOPAction(OPActionGetGameReport)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	action := c.Ctx.Input.Param(KEY_KEY_ACTION)
	start, _ := strconv.ParseInt(c.GetString(KEY_START_TIME), 10, 64)
	over, _ := strconv.ParseInt(c.GetString(KEY_END_TIME), 10, 64)
	page, _ := c.GetInt(KEY_PAGE, 1)
	limit, _ := c.GetInt(KEY_LIMIT, 10)

	var total int64
	var data interface{}
	var err error
	switch action {
	case "ag":
		total, data, err = reportdao.ReportGameRecordAgDaoEntity.QueryAllByTime(start, over, page, limit)
		if err != nil {
			c.ErrorResponse(controllers.ERROR_CODE_DB)
		}
		break
	case "ky":
		total, data, err = reportdao.ReportGameRecordKyDaoEntity.QueryAllByTime(start, over, page, limit)
		if err != nil {
			c.ErrorResponse(controllers.ERROR_CODE_DB)
		}
		break
	case "rg":
		total, data, err = reportdao.ReportGameRecordRgDaoEntity.QueryAllByTime(start, over, page, limit)
		if err != nil {
			c.ErrorResponse(controllers.ERROR_CODE_DB)
		}
		break
	default:
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		break
	}

	c.SuccessResponse(map[string]interface{}{
		KEY_LIST: data,
		KEY_META: PageInfo{
			Page:  page,
			Total: int(total),
			Limit: limit,
		},
	})

}

/**
*手动查询游戏用户充值和提现记录
 */
func (c *GameInfoOptController) ManaulWriteRisk() {

	userRiskList := gamedao.GameTransferDaoEntity.QueryTodayRiskInfo()
	resultList, err := gamedao.GameRiskAlertDaoEntiry.InsertRiskInfo(userRiskList)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_RISK_ALERT)
	}
	c.SuccessResponse(resultList)
}

//
func execWriteTmpUser(reportGameRecord []interface{}, channelId uint32, uid uint64, start int64) {
	//投了几个游戏就有几个新增用户，就写入tmpbetters
	if len(reportGameRecord) > 0 {

		for _, v := range reportGameRecord {
			var gameId = 0
			betStatus := tmpdao.NOTBET
			if k, ok := v.(reportmodels.ReportGameRecordAg); ok {
				gameId, _ = strconv.Atoi(k.GameType)
				if k.Bet > 0 {
					betStatus = tmpdao.HASBET
				}
			} else if k, ok := v.(reportmodels.ReportGameRecordKy); ok {
				gameId, _ = strconv.Atoi(k.KindId)
				if k.Bet > 0 {
					betStatus = tmpdao.HASBET
				}
			} else if k, ok := v.(reportmodels.ReportGameRecordRg); ok {
				gameId, _ = strconv.Atoi(k.GameNameID)
				if k.Bet > 0 {
					betStatus = tmpdao.HASBET
				}
			}

			tmpdao.TmpGameBetersDaoEntity.InsertUser(uid, channelId, uint32(gameId), uint32(betStatus), start)
		}
	}
}

//
func execUpTmpUser(tmpUsers models.TmpGamebeters, reportGameRecord []interface{}, start int64) {
	if len(reportGameRecord) > 0 {
		//tmpdao.TmpGameBetersDaoEntity.DeleteByObj(tmpUsers)
		for _, v := range reportGameRecord {
			var gameId = 0
			if k, ok := v.(reportmodels.ReportGameRecordAg); ok {
				gameId, _ = strconv.Atoi(k.GameType)
			} else if k, ok := v.(reportmodels.ReportGameRecordKy); ok {
				gameId, _ = strconv.Atoi(k.KindId)
			} else if k, ok := v.(reportmodels.ReportGameRecordRg); ok {
				gameId, _ = strconv.Atoi(k.GameNameID)
			}
			tmpdao.TmpGameBetersDaoEntity.UpdateUser(tmpUsers.Uid, uint32(gameId), tmpUsers.ChannelId, tmpdao.HASBET, start)
		}
	}
}

//game_user表中新注册的用户来是否投注 channelId
func writeBetToday(channelId uint32, uid uint64, start, over int64) {
	switch channelId {
	case gameapi.GAME_CHANNEL_AG:
		//看用户都投注了哪些游戏
		num, reportGameRecordAg, _ := reportdao.ReportGameRecordAgDaoEntity.GetBetUsers(start, over, uint32(uid), true, 0)
		if num <= 0 {
			tmpdao.TmpGameBetersDaoEntity.InsertUser(uid, channelId, 0, tmpdao.NOTBET, start)
		} else {
			execWriteTmpUser(reportGameRecordAg, channelId, uid, start)
		}

		break
	case gameapi.GAME_CHANNEL_KY:
		num, reportGameRecordKy, _ := reportdao.ReportGameRecordKyDaoEntity.GetBetUsers(start, over, uint32(uid), true, 0)
		if num <= 0 {
			tmpdao.TmpGameBetersDaoEntity.InsertUser(uid, channelId, 0, tmpdao.NOTBET, start)
		} else {
			execWriteTmpUser(reportGameRecordKy, channelId, uid, start)
		}

		break
	case gameapi.GAME_CHANNEL_RG:
		//看用户都投注了哪些游戏
		num, reportGameRecordRg, _ := reportdao.ReportGameRecordRgDaoEntity.GetBetUsers(start, over, uint32(uid), true, 0)
		if num <= 0 {
			tmpdao.TmpGameBetersDaoEntity.InsertUser(uid, channelId, 0, tmpdao.NOTBET, start)
		} else {
			execWriteTmpUser(reportGameRecordRg, channelId, uid, start)
		}

		break
	}
	return
}

//判断从tmp_game_beters表中过来的未投注用户今天是否有投注
func upBetToday(tmpUsers models.TmpGamebeters, start, over int64) (isBet bool) {
	switch tmpUsers.ChannelId {
	case gameapi.GAME_CHANNEL_AG:
		//看这个之前注册的用户都投注了哪些游戏
		num, userGameRecordAg, _ := reportdao.ReportGameRecordAgDaoEntity.GetBetUsers(start, over, uint32(tmpUsers.Uid), true, tmpUsers.GameId)
		if num > 0 {
			execUpTmpUser(tmpUsers, userGameRecordAg, start)
		}
		break
	case gameapi.GAME_CHANNEL_KY:
		num, userGameRecordKy, _ := reportdao.ReportGameRecordKyDaoEntity.GetBetUsers(start, over, uint32(tmpUsers.Uid), true, tmpUsers.GameId)
		if num > 0 {
			execUpTmpUser(tmpUsers, userGameRecordKy, start)
		}

		break
	case gameapi.GAME_CHANNEL_RG:
		num, userGameRecordRg, _ := reportdao.ReportGameRecordRgDaoEntity.GetBetUsers(start, over, uint32(tmpUsers.Uid), true, tmpUsers.GameId)
		if num > 0 {
			execUpTmpUser(tmpUsers, userGameRecordRg, start)
		}

		break
	}
	return
}

/**
 * 返回新用户数,两种情况下会产生新用户：
 * 1 ， 用户当天注册且投注
 * 2 ， 用户以往注册，但是在当天投注
 * 注： 投注指只要用户有投注就算，无需有效投注
 */
func getNewPlayer(start, over int64, channelId uint32) (newersList []tmpdao.NewersList, err error) {
	regUsers, err := gamedao.GameUserDaoEntity.QueryUserByTimestamp(channelId, start, over)
	if err != nil {
		return
	}
	//查询是否有投注
	for _, v := range regUsers {
		writeBetToday(v.ChannelId, v.Uid, start, over)
	}

	//临时表所有用户未bet的用户
	todayUnbetUser, err := tmpdao.TmpGameBetersDaoEntity.QueryUnbetUser(channelId)
	if err != nil {
		return
	}
	//查询是否有之前注册的用户今天投注了，如果有代表是新用户
	for _, v := range todayUnbetUser {
		//只要有投注就是新用户，更新投注状态
		upBetToday(v, start, over)
	}

	newersList, err = tmpdao.TmpGameBetersDaoEntity.GetTodayNewerList(channelId, start, over)

	return
}

func (c *GameInfoOptController) GenStatisticInfo() {
	c.setOPAction(OPActionGenReportDaily)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	action := c.Ctx.Input.Param(KEY_KEY_ACTION)

	channelId := getChannelId(action)
	if channelId == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
	}
	//var reports []reportmodels.ReportStatisticGameAll
	reports := StatisticNow(channelId, c)
	c.SuccessResponse(reports)
}

//每天生成ag平台游戏数据,返1003表示无数据或生成错误
func StatisticNow(gameChannelId uint32, c *GameInfoOptController) (reportList []reportmodels.ReportStatisticGameAll) {
	//生成都是前一天的数据
	var start = common.GetPreDateStartTimestamp(1)
	var over = common.GetPreDateOverTimestamp(1)
	//查看是否统计过
	num, err := reportdao.ReportStatisticGameAllDaoEntity.FindByChannalId(gameChannelId, start)
	if num > 0 {
		c.ErrorResponseAndLog(OPActionGenReportDaily, controllers.ERROR_CODE_GAME_DATA_HASED, "")
	}
	if err != nil {
		c.ErrorResponseAndLog(OPActionGenReportDaily, controllers.ERROR_CODE_GAME_REPORT_CAL, "")
	}
	//新增玩家
	newersList, err := getNewPlayer(start, over, gameChannelId)
	if err != nil {
		c.ErrorResponseAndLog(OPActionGenReportDaily, controllers.ERROR_CODE_GAME_REPORT_CAL, "")
	}
	//开始统计
	if gameChannelId == gameapi.GAME_CHANNEL_AG {
		reportAgList, _ := reportdao.ReportGameRecordAgDaoEntity.GetInfo(start, over)
		reportList = make([]reportmodels.ReportStatisticGameAll, len(reportAgList))
		for k, o := range reportAgList {
			setCommonValue(&(reportList[k]), gameChannelId, newersList, o.GameType, o.Bet, o.ValidBet, o.Profit, 0, start)
		}
	} else if gameChannelId == gameapi.GAME_CHANNEL_KY {
		reportkyList, _ := reportdao.ReportGameRecordKyDaoEntity.GetInfo(start, over)
		reportList = make([]reportmodels.ReportStatisticGameAll, len(reportkyList))
		for k, o := range reportkyList {
			setCommonValue(&(reportList[k]), gameChannelId, newersList, o.KindId, o.Bet, o.ValidBet, o.Profit, o.Revenue, start)
		}
	} else {
		reportRgList, _ := reportdao.ReportGameRecordRgDaoEntity.GetInfo(start, over)
		reportList = make([]reportmodels.ReportStatisticGameAll, len(reportRgList))
		for k, o := range reportRgList {
			setCommonValue(&(reportList[k]), gameChannelId, newersList, o.GameNameID, o.Bet, o.ValidBet, o.Profit, 0, start)
		}
	}
	err = reportdao.ReportStatisticGameAllDaoEntity.InertReport(reportList)
	if err != nil {
		c.ErrorResponseAndLog(OPActionGenReportDaily, controllers.ERROR_CODE_GAME_REPORT_CAL, "")
	}

	//生成合计数据
	//按日合计
	var dayStatistic = reportmodels.ReportStatisticSum{}
	dayStatistic.Ctime = start
	dayStatistic.ChannelId = gameChannelId
	//活跃玩家统计
	dayStatistic.ChannelPositiveNums = uint64(reportdao.ReportGameUserDailyDaoEntity.GetPositiveUserNums(gameChannelId))

	//红包月分红
	//if gameChannelId == gameapi.GAME_CHANNEL_RG {
	//	var monthDivid reportmodels.MonthDividendRecord
	//	monthDivid, err = reportdao.MonthDividendRecordDaoEntity.FindByRange(start, over)
	//	if err != nil {
	//		c.ErrorResponse(controllers.ERROR_CODE_GAME_REPORT_CAL)
	//	} else {
	//		dayStatistic.ChannelRgDividend = monthDivid.ResultDividend
	//	}
	//}
	dayStatistic.ChannelSalaryDaily, err = reportdao.ReportGameUserDailyDaoEntity.GetSalaryByChannelId(gameChannelId, start)
	if err != nil {
		c.ErrorResponseAndLog(OPActionGenReportDaily, controllers.ERROR_CODE_GAME_REPORT_CAL, "")
	}
	dayStatistic.ChannelRechargeEusd, err = gamedao.GameTransferDaoEntity.QueryFundRwByChannelID(gameChannelId, gamedao.RECHARGE_TYPE, start, over)
	if err != nil {
		c.ErrorResponseAndLog(OPActionGenReportDaily, controllers.ERROR_CODE_GAME_REPORT_CAL, "")
	}
	dayStatistic.ChannelWithdrawEusd, err = gamedao.GameTransferDaoEntity.QueryFundRwByChannelID(gameChannelId, gamedao.WITHDRAW_TYPE, start, over)
	if err != nil {
		c.ErrorResponseAndLog(OPActionGenReportDaily, controllers.ERROR_CODE_GAME_REPORT_CAL, "")
	}
	err = reportdao.ReportStatisticSumDaoEntity.InsertData(dayStatistic)
	if err != nil {
		c.ErrorResponseAndLog(OPActionGenReportDaily, controllers.ERROR_CODE_GAME_REPORT_CAL, "")
	}
	return
	/*按月合计

	//生成合计数据 //按月
	var reportMonth = reportmodels.ReportStatisticMonth{}
	reportMonth.ChannelId = gameChannelId
	reportMonth.Ctime = common.GetMonthDateStartTimestamp(time.Unix(start, 0))
	reportMonth.ChannelPositiveNums = uint64(reportdao.ReportGameUserDailyDaoEntity.GetPositiveUserNums(gameChannelId))

	//红包月分红
	if gameChannelId == gameapi.GAME_CHANNEL_RG {
		var monthDivid reportmodels.MonthDividendRecord
		monthDivid, err = reportdao.MonthDividendRecordDaoEntity.FindByRange(start, over)
		if err != nil {
			c.ErrorResponse(controllers.ERROR_CODE_GAME_REPORT_CAL)
		} else {
			reportMonth.ChannelRgDividend = monthDivid.ResultDividend
		}
	}

	reportMonth.ChannelSalaryDaily, err = reportdao.ReportGameUserDailyDaoEntity.GetMonthSalaryByChannelId(gameChannelId, start)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_GAME_REPORT_CAL)
	}
	reportMonth.ChannelRechargeEusd, err = gamedao.GameTransferDaoEntity.QueryMonthFundRwByChannelID(gameChannelId, gamedao.RECHARGE_TYPE, start, over)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_GAME_REPORT_CAL)
	}
	reportMonth.ChannelWithdrawEusd, err = gamedao.GameTransferDaoEntity.QueryMonthFundRwByChannelID(gameChannelId, gamedao.WITHDRAW_TYPE, start, over)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_GAME_REPORT_CAL)
	}

	err = reportdao.ReportStatisticMonthDaoEntity.UpOrInsertData(reportMonth)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_GAME_REPORT_CAL)
	}

	return
	*/
}

func setCommonValue(report *reportmodels.ReportStatisticGameAll, channelId uint32, newersList []tmpdao.NewersList, gameId string, bet, validBet, profit, revenue, start int64) {
	report.Ctime = start
	if revenue > 0 {
		report.Revenue = revenue
	}
	reportGameId, _ := strconv.Atoi(gameId)
	for _, n := range newersList {
		if uint32(reportGameId) == n.GameId {
			report.NewerNums = int64(n.Nums)
		}
	}
	report.ChannelId = channelId
	report.GameId = uint32(reportGameId)
	report.Bet = bet
	report.ValidBet = validBet
	report.Profit = profit
}

func getChannelId(channelName string) (channelId uint32) {
	switch channelName {
	case "ag":
		channelId = gameapi.GAME_CHANNEL_AG
		break
	case "ky":
		channelId = gameapi.GAME_CHANNEL_KY
		break
	case "rg":
		channelId = gameapi.GAME_CHANNEL_RG
		break
	default:
		channelId = 0
		break
	}
	return
}

func (c *GameInfoOptController) DebugStatisticGen() {

	res, err := http.Get(fmt.Sprintf("%s%s/v1/admin/game/statistic/ag", "http://127.0.0.1:", beego.AppConfig.String("httpport")))
	if err != nil {
		common.LogFuncError("err is : %s", err)
	}

	cont := make([]byte, 1024)
	res.Body.Read(cont)
	c.SuccessResponse(string(cont))
}

/*
 		1 计算购买和出售各多少钱
		1-1 ??最后一笔出售时间小于购买且其金额如果大于购买的50%时继续，即先卖后买
		2 当出售总eusd大于300个，继续
		3 出售和购买不是同一个帐户
		4 出售总金额大于购买总金额的80%时
		5 风险写入
*/
func (c *GameInfoOptController) ManualOrderRisk() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	now := time.Now().Unix()
	t := now - dao.RISK_TIME_RANGE
	//order中uid查找用户购买总额和出售总额
	orders, err := dao.OrdersDaoEntity.QueryDataByTime(t)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_GAME_REPORT_CAL)
		return
	}
	//logs.Debug("init info : ", orders, err)

	var risk_num = 0
	//last风险指数
	//var lastSoldNum int64
	//查看是否有风险
	for k, v := range orders {
		var buySum, soldSum int64
		var sameName = ""
		var isSameAccount = true
		//避免重复写入
		var lastSoldId uint64
		for n, m := range v {
			if n == 0 && m.PayType == dao.SideSell {
				sameName = m.PayAccount
				lastSoldId = m.Id
				//	lastSoldNum = m.Amount
			}

			if m.PayAccount != sameName {
				isSameAccount = false
			}

			if m.PayType == dao.SideSell {
				soldSum += m.Amount
			} else if m.PayType == dao.SideBuy {
				buySum += m.Amount
			}
		}

		if isSameAccount {
			continue
		}
		if soldSum > eosplus.QuantityFloat64ToInt64(dao.RISK_EUSD_VALUE) {
			if soldSum >= int64(float64(buySum)*dao.RISK_SOLD_RATE) {
				//生成风险信息
				insertId, _ := gamedao.GameRiskAlertDaoEntiry.InsertOrderRiskInfo(k, soldSum, now, now, gamedao.WARN_GRADE_ONE, gamedao.RISK_TYPE_ORDER, lastSoldId)

				if insertId > 0 {
					_ = gamedao.GameOrderRiskDaoEntity.InsertRiskItem(v, insertId, now)
					risk_num++
				}

			}
		}

	}

	c.SuccessResponse(fmt.Sprintf("risk_num : %d", risk_num))
}
