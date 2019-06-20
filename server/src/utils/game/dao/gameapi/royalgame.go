package gameapi

import (
	"common"
	"eusd/eosplus"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"

	admindao "utils/admin/dao"
	gamedao "utils/game/dao"
	reportdao "utils/report/dao"
	reportmodels "utils/report/models"

	"github.com/astaxie/beego/httplib"
	json "github.com/mailru/easyjson"
)

type RoyalGameAPI struct {
	baseURL, prefix, desKey, keyMD5 string
}

func NewRoyalGameAPI(baseURL, prefix, desKey string) (*RoyalGameAPI, error) {
	keyMD5, err := md5hex([]byte(desKey))
	if err != nil {
		return nil, err
	}
	return &RoyalGameAPI{
		baseURL: baseURL,
		prefix:  prefix,
		desKey:  desKey,
		keyMD5:  keyMD5,
	}, nil
}

//easyjson:json
type RGBaseRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//easyjson:json
type RGLoginRes struct {
	RGBaseRes
	Data string `json:"data"`
}

//easyjson:json
type RGLoginReq struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Type     int    `json:"type"`
	Sign     string `json:"sign"`
	Prefix   string `json:"prefix"`
	GameId   int    `json:"gameId"`
	Datetime string `json:"datetime"`
}

const (
	methodLogin       = "/open/login"
	methodBalance     = "/open/queryBalance"
	methodTransfer    = "/open/transfer"
	methodLotteryList = "/open/lotteryList"
)

type rgTransferType int

const (
	rgTransferIn  = 1
	rgTransferOut = 2
)

//下注记录状态
const (
	RgLotteryStatusNil        = iota //全部
	RgLotteryStatusUnOPen            //未开奖
	RgLotteryStatusLost              //未中奖
	RgLotteryStatusWin               //中奖
	RgLotteryStatusUserCancel        //已撤单
	RgLotteryStatusSysCancel         //系统撤单
	RgLotteryStatusCancelOpen        //撤销开奖
)

const (
	RgSearchTimeTypeOpen = iota //按开奖时间搜索
	RgSearchTimeTypeBet         //按下注时间搜索
)

func (api *RoyalGameAPI) Login(account, password, ip string, KindID string) (lr *LoginReply, err error) {
	var (
		gameID int
		data   []byte
		sign   string
	)

	gameID, err = strconv.Atoi(KindID)
	if err != nil {
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	password, err = md5hex([]byte(password))
	if err != nil {
		return
	}

	sign, err = md5hex([]byte(fmt.Sprintf("%s%s%s%s%s", now, api.prefix, account, password, api.keyMD5)))
	if err != nil {
		return
	}
	param := &RGLoginReq{
		UserName: account,
		Password: password,
		Type:     2,
		Prefix:   api.prefix,
		GameId:   gameID,
		Datetime: now,
		Sign:     sign,
	}
	data, err = api.request(methodLogin, param)
	fmt.Println(string(data))
	if err != nil {
		return
	}

	reply := &RGLoginRes{}
	err = json.Unmarshal(data, reply)
	if err != nil {
		return
	}
	if reply.Code != 0 {
		err = fmt.Errorf(reply.Message)
		common.LogFuncError("res code:%v", reply.Code)
		return
	}
	return &LoginReply{URL: reply.Data}, nil
}

func (api *RoyalGameAPI) Register(account, password, ip string) error { return nil }
func (api *RoyalGameAPI) Logout(account, password string)             {}

//easyjson:json
type RGBalanceReq struct {
	UserName string `json:"username"`
	Sign     string `json:"sign"`
	Prefix   string `json:"prefix"`
	Datetime string `json:"datetime"`
}

//easyjson:json
type RGBalanceReply struct {
	Balance float64 `json:"balance"`
}

//easyjson:json
type RGBalanceRes struct {
	RGBaseRes
	Data RGBalanceReply `json:"data"`
}

// GetAccount 查询账号信息
func (api *RoyalGameAPI) GetBalance(account, password string) (balance float64, err error) {
	var (
		data []byte
		sign string
	)

	now := time.Now().Format("2006-01-02 15:04:05")

	password, err = md5hex([]byte(password))
	if err != nil {
		return
	}

	sign, err = md5hex([]byte(fmt.Sprintf("%s%s%s%s", now, api.prefix, account, api.keyMD5)))
	if err != nil {
		return
	}
	param := &RGBalanceReq{
		UserName: account,
		Prefix:   api.prefix,
		Datetime: now,
		Sign:     sign,
	}
	data, err = api.request(methodBalance, param)
	fmt.Println(string(data))
	if err != nil {
		return
	}
	reply := &RGBalanceRes{}
	err = json.Unmarshal(data, reply)
	if err != nil {
		return
	}
	if reply.Code != 0 {
		err = fmt.Errorf(reply.Message)
		return
	}
	balance = reply.Data.Balance
	return
}

//easyjson:json
type RGTransferReq struct {
	UserName            string  `json:"username"`
	Sign                string  `json:"sign"`
	Prefix              string  `json:"prefix"`
	MerchantOrderNumber string  `json:"merchantordernumber"`
	Type                int     `json:"type"`
	Amount              float64 `json:"amount"`
	Datetime            string  `json:"datetime"`
}

//easyjson:json
type RGTransferReply struct {
	Balance     float64 `json:"balance"`
	OrderNumber string  `json:"orderNumber"`
}

//easyjson:json
type RGTransferRes struct {
	RGBaseRes
	Data RGTransferReply `json:"data"`
}

func (api *RoyalGameAPI) TransferIn(account, password, orderNum string, money float64) (tr *TransferInReply, err error) {
	var (
		data []byte
	)

	password, err = md5hex([]byte(password))
	if err != nil {
		return
	}

	data, err = api.transfer(account, password, orderNum, money, rgTransferIn)
	reply := &RGTransferRes{}
	err = json.Unmarshal(data, reply)
	if err != nil {
		return
	}
	if reply.Code != 0 {
		err = fmt.Errorf(reply.Message)
		return
	}
	return &TransferInReply{Order: reply.Data.OrderNumber, Success: true}, nil

}
func (api *RoyalGameAPI) TransferOut(account, password, orderNum string, money float64) (tr *TransferOutReply, err error) {
	var (
		data []byte
	)

	password, err = md5hex([]byte(password))
	if err != nil {
		return
	}

	data, err = api.transfer(account, password, orderNum, money, rgTransferOut)
	reply := &RGTransferRes{}
	err = json.Unmarshal(data, reply)
	if err != nil {
		return
	}
	if reply.Code != 0 {
		err = fmt.Errorf(reply.Message)
		return
	}
	return &TransferOutReply{Order: reply.Data.OrderNumber, Success: true}, nil

}
func (api *RoyalGameAPI) transfer(account, password, orderNum string, money float64, transferType rgTransferType) (data []byte, err error) {
	var (
		sign string
	)

	now := time.Now().Format("2006-01-02 15:04:05")

	sign, err = md5hex([]byte(fmt.Sprintf("%s%s%s%s", now, api.prefix, account, api.keyMD5)))
	if err != nil {
		return
	}
	param := &RGTransferReq{
		UserName:            account,
		Sign:                sign,
		Prefix:              api.prefix,
		Type:                int(transferType), //1-转入；2-转出
		Amount:              money,
		Datetime:            now,
		MerchantOrderNumber: orderNum,
	}
	data, err = api.request(methodTransfer, param)
	fmt.Println(string(data))
	if err != nil {
		return
	}

	return
}

func (api *RoyalGameAPI) request(method string, param json.Marshaler) (data []byte, err error) {
	url := fmt.Sprintf("%s%s", api.baseURL, method)
	req := httplib.Post(url)
	req.Header("Content-Type", "application/json")

	var buf []byte
	buf, err = json.Marshal(param)
	if err != nil {
		return
	}
	body := string(buf)
	req.Body(body)
	data, err = req.Bytes()
	if err != nil {
		return nil, err
	}
	common.LogFuncDebug("req [%+v] \t param [%+v] \t resp [%+v]", url, body, string(data))
	return
}

//easyjson:json
type RGLotteryListReq struct {
	Prefix         string `json:"prefix"`
	Sign           string `json:"sign"`
	Datetime       string `json:"datetime"`
	StartTime      string `json:"starttime"`
	EndTime        string `json:"endtime"`
	PageIndex      int    `json:"pageindex"`
	PageSize       int    `json:"pagesize"`
	Status         int    `json:"status"`
	SearchTimeType int    `json:"searchtimetype"`
}

//easyjson:json
type RGLotteryList struct {
	Username      string  `json:"username"`
	OrderNumber   string  `json:"orderNumber"`
	PeriodName    string  `json:"periodName"`
	GameNameID    int     `json:"gameNameID"`
	GameName      string  `json:"gameName"`
	BetCount      int     `json:"betCount"`
	BettingAmount float64 `json:"bettingAmount"`
	State         int     `json:"state"`
	GameKindID    int     `json:"gameKindID"`
	AddTime       string  `json:"addTime"`
	OpenNumber    string  `json:"openNumber"`
	OpenDate      string  `json:"openDate"`
	WinLoseAmount float64 `json:"winLoseAmount"`
	BetContent    string  `json:"betContent"`
	GameKindName  string  `json:"gameKindName"`
}

//easyjson:json
type RGLotteryListReply struct {
	Total            int             `json:"total"`
	CurrentPageIndex int             `json:"currentPageIndex"`
	TotalPage        int             `json:"totalPage"`
	DataList         []RGLotteryList `json:"dataList"`
}

//easyjson:json
type RGLotteryListRes struct {
	RGBaseRes
	Data RGLotteryListReply `json:"data"`
}

func (api *RoyalGameAPI) getBetRecords(pageIndex, pageSize, status, searceType int, startTime, endTime string) (replyData RGLotteryListReply, err error) {
	var (
		data []byte
		sign string
	)

	now := time.Now().Format("2006-01-02 15:04:05")
	sign, err = md5hex([]byte(fmt.Sprintf("%s%s%s", now, api.prefix, api.keyMD5)))
	if err != nil {
		return
	}
	param := &RGLotteryListReq{
		Prefix:         api.prefix,
		Sign:           sign,
		Datetime:       now,
		PageIndex:      pageIndex,
		PageSize:       pageSize,
		StartTime:      startTime,
		EndTime:        endTime,
		SearchTimeType: searceType,
		Status:         status,
	}
	data, err = api.request(methodLotteryList, param)
	fmt.Println(string(data))
	if err != nil {
		return
	}

	reply := &RGLotteryListRes{}
	err = json.Unmarshal(data, reply)
	if err != nil {
		return
	}
	if reply.Code != 0 {
		err = fmt.Errorf(reply.Message)
		return
	}
	replyData = reply.Data
	/*list = &LotteryListReply{}
	list.Total = reply.Data.Total
	list.CurrentPageIndex = reply.Data.CurrentPageIndex
	list.TotalPage = reply.Data.TotalPage
	for _, item := range reply.Data.DataList {
		if item.State == RgLotteryStatusLost || item.State == RgLotteryStatusWin {
			list.DataList = append(list.DataList, LotteryList{
				Username:      item.Username,
				GameNameID:    fmt.Sprintf("%v", item.GameNameID),
				GameKindID:    fmt.Sprintf("%v", item.GameKindID),
				BetCount:      item.BetCount,
				BettingAmount: item.BettingAmount,
				WinLoseAmount: item.WinLoseAmount,
			})
		}
	}*/

	return
}

func (api *RoyalGameAPI) dayLotteryRecords(exchangeRate, precision int32, zeroTimestamp int64, startTime, endTime string, mapWhite orm.Params) (err error) {
	common.LogFuncDebug("dayLotteryRecords startTime:%v,endTime:%v", startTime, endTime)
	curPageIndex := 1
	//获取开奖数据
	for {
		var reportGameRecordRgs []reportmodels.ReportGameRecordRg
		var res RGLotteryListReply
		res, err = api.getBetRecords(curPageIndex, pageLimit, RgLotteryStatusNil, RgSearchTimeTypeOpen, startTime, endTime)
		if err != nil {
			common.LogFuncError("error:%v", err)
			return
		}

		var accounts []string
		//个人投注量统计
		for _, item := range res.DataList {
			bet := float64(eosplus.QuantityFloat64ToInt64(item.BettingAmount))
			Profit := float64(eosplus.QuantityFloat64ToInt64(item.WinLoseAmount))
			reportGameRecordRg := reportmodels.ReportGameRecordRg{
				Account:      item.Username,
				GameNameID:   fmt.Sprintf("%v", item.GameNameID),
				GameName:     item.GameName,
				GameKindName: item.GameKindName,
				OrderId:      item.OrderNumber,
				OpenDate:     item.OpenDate,
				PeriodName:   item.PeriodName,
				OpenNumber:   item.OpenNumber,
				Status:       uint8(item.State),
				Bet:          int64(common.GameGameCoin2Eusd(bet, exchangeRate, precision)),    //转换eusd,
				Profit:       int64(common.GameGameCoin2Eusd(Profit, exchangeRate, precision)), //转换eusd,
				BetTime:      item.AddTime,
				BetContent:   item.BetContent,
				Ctime:        zeroTimestamp,
			}

			if item.State == RgLotteryStatusLost || item.State == RgLotteryStatusWin {
				if _, ok := mapWhite[reportGameRecordRg.GameNameID]; !ok {
					//非应用白名单
					reportGameRecordRg.ValidBet = reportGameRecordRg.Bet
				}
			}
			accounts = append(accounts, "\""+reportGameRecordRg.Account+"\"")
			reportGameRecordRgs = append(reportGameRecordRgs, reportGameRecordRg)
		}

		var mapAccounts orm.Params
		mapAccounts, err = gamedao.GameUserDaoEntity.QueryUidByAccounts(GAME_CHANNEL_RG, accounts)
		if err != nil {
			common.LogFuncError("dayLotteryRecords error:%v", err)
			return
		}
		for i := 0; i < len(reportGameRecordRgs); i++ {
			if _, ok := mapAccounts[reportGameRecordRgs[i].Account]; ok {
				strUid := mapAccounts[reportGameRecordRgs[i].Account].(string)
				uid, _ := strconv.ParseUint(strUid, 10, 64)
				reportGameRecordRgs[i].Uid = uid
			}
		}
		//入库
		err = reportdao.ReportGameRecordRgDaoEntity.InsertMul(zeroTimestamp, reportGameRecordRgs)
		if err != nil {
			common.LogFuncError("dayLotteryRecords error:%v", err)
			return
		}

		if curPageIndex >= res.TotalPage {
			break
		}
		curPageIndex = curPageIndex + 1
	}
	return
}

func (api *RoyalGameAPI) dayUnOpenRecords(exchangeRate, precision int32, zeroTimestamp int64, startTime, endTime string, mapWhite orm.Params) (err error) {
	common.LogFuncDebug("dayUnOpenRecords startTime:%v,endTime:%v", startTime, endTime)
	curPageIndex := 1
	//获取未开奖数据
	for {
		var reportGameRecordRgs []reportmodels.ReportGameRecordRg
		var res RGLotteryListReply
		res, err = api.getBetRecords(curPageIndex, pageLimit, RgLotteryStatusUnOPen, RgSearchTimeTypeBet, startTime, endTime)
		if err != nil {
			common.LogFuncError("error:%v", err)
			return
		}

		var accounts []string
		//个人投注量统计
		for _, item := range res.DataList {
			bet := float64(eosplus.QuantityFloat64ToInt64(item.BettingAmount))
			Profit := float64(eosplus.QuantityFloat64ToInt64(item.WinLoseAmount))
			reportGameRecordRg := reportmodels.ReportGameRecordRg{
				Account:      item.Username,
				GameNameID:   fmt.Sprintf("%v", item.GameNameID),
				GameName:     item.GameName,
				GameKindName: item.GameKindName,
				OrderId:      item.OrderNumber,
				OpenDate:     item.OpenDate,
				PeriodName:   item.PeriodName,
				OpenNumber:   item.OpenNumber,
				Status:       uint8(item.State),
				Bet:          int64(common.GameGameCoin2Eusd(bet, exchangeRate, precision)), //转换eusd,
				ValidBet:     0,
				Profit:       int64(common.GameGameCoin2Eusd(Profit, exchangeRate, precision)), //转换eusd,
				BetTime:      item.AddTime,
				BetContent:   item.BetContent,
				Ctime:        zeroTimestamp,
			}
			accounts = append(accounts, "\""+reportGameRecordRg.Account+"\"")
			reportGameRecordRgs = append(reportGameRecordRgs, reportGameRecordRg)
		}

		var mapAccounts orm.Params
		mapAccounts, err = gamedao.GameUserDaoEntity.QueryUidByAccounts(GAME_CHANNEL_RG, accounts)
		if err != nil {
			common.LogFuncError("dayLotteryRecords error:%v", err)
			return
		}
		for i := 0; i < len(reportGameRecordRgs); i++ {
			if _, ok := mapAccounts[reportGameRecordRgs[i].Account]; ok {
				strUid := mapAccounts[reportGameRecordRgs[i].Account].(string)
				uid, _ := strconv.ParseUint(strUid, 10, 64)
				reportGameRecordRgs[i].Uid = uid
			}
		}
		//入库
		err = reportdao.ReportGameRecordRgDaoEntity.InsertMul(zeroTimestamp, reportGameRecordRgs)
		if err != nil {
			common.LogFuncError("dayUnOpenRecords error:%v", err)
			return
		}

		if curPageIndex >= res.TotalPage {
			break
		}
		curPageIndex = curPageIndex + 1
	}

	return
}

//查询时间段左闭右开[startTime, endTime)
func (api *RoyalGameAPI) DayBetRecords(timestamp int64, mapWhite orm.Params) (err error) {
	appChannel, err := admindao.AppChannelDaoEntity.QueryById(GAME_CHANNEL_KY)
	if err != nil || appChannel == nil {
		common.LogFuncError("get channelId[%v] cfg fail", GAME_CHANNEL_KY)
		return
	}

	err = reportdao.ReportGameRecordRgDaoEntity.DelByTimestamp(timestamp)
	if err != nil {
		common.LogFuncError("DayBetRecords error:%v", err)
	}

	//startTime, endTime := common.GetBeginAndEndString(timestamp)
	startTime := common.GetTimeFormat(timestamp)
	endTime := common.GetTimeFormat(timestamp + common.DaySeconds)
	common.LogFuncDebug("startTime:%v,endTime:%v", startTime, endTime)

	//一天开奖数据
	err = api.dayLotteryRecords(appChannel.ExchangeRate, appChannel.Precision, timestamp, startTime, endTime, mapWhite)
	if err != nil {
		reportdao.ReportGameRecordRgDaoEntity.DelByTimestamp(timestamp)
		common.LogFuncError("DayBetRecords error:%v", err)
		return
	}

	//一天未开奖数据
	err = api.dayUnOpenRecords(appChannel.ExchangeRate, appChannel.Precision, timestamp, startTime, endTime, mapWhite)
	if err != nil {
		reportdao.ReportGameRecordRgDaoEntity.DelByTimestamp(timestamp)
		common.LogFuncError("DayBetRecords error:%v", err)
		return
	}

	return
}

func (api *RoyalGameAPI) GetTotalByTimestamp(timestamp int64) (total int, err error) {
	total, err = reportdao.ReportGameRecordRgDaoEntity.QueryTotalByTimestamp(timestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (api *RoyalGameAPI) GetRecordByTimestamp(timestamp int64, page, limit int) (betInfos []reportdao.BetInfo, err error) {
	betInfos, err = reportdao.ReportGameRecordRgDaoEntity.QueryByTimestamp(timestamp, page, limit)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (api *RoyalGameAPI) DelBetRecord(timestamp int64) (err error) {
	err = reportdao.ReportGameRecordRgDaoEntity.DelByTimestamp(timestamp)
	return
}
