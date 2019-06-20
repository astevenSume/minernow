package gameapi

import (
	"common"
	"encoding/xml"
	"eusd/eosplus"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"

	admindao "utils/admin/dao"
	gamedao "utils/game/dao"
	reportdao "utils/report/dao"
	reportmodels "utils/report/models"

	"github.com/astaxie/beego/httplib"
	json "github.com/mailru/easyjson"
)

type AsiaGamingAPI struct {
	baseURL, prefix, desKey string
}

func NewAsiaGamingAPI(baseURL, prefix, desKey string) GameAPI {
	return &AsiaGamingAPI{
		baseURL: baseURL,
		prefix:  prefix,
		desKey:  desKey,
	}
}

const (
	methodCheckOrCreateGameAccount = "/Ag/CheckOrCreateGameAccount"
	methodForwardGame              = "/Ag/ForwardGame"
	methodGetBalance               = "/Ag/GetBalance"
	methodPrepareTransferCredit    = "/Ag/PrepareTransferCredit"
	methodTransferCreditConfirm    = "/Ag/TransferCreditConfirm"
	methodQueryOrderStatus         = "/Ag/QueryOrderStatus"
	methodGetBetRecords            = "/Ag/GetBetRecords"

	KeyHeadAgentCode = "X-Pos-Code"
	KeyHeadAgentSign = "X-Signature-Code"
)

type agTransferType string

const (
	agTransferIn  = agTransferType("IN")
	agTransferOut = agTransferType("OUT")
)

//easyjson:json
type AGLoginReq struct {
	CagentType int    `json:"cagentType"`
	UserName   string `json:"userName"`
	Password   string `json:"password"`
	GameType   string `json:"gameType"`
	OddType    string `json:"oddtype"`
	Timespan   string `json:"timespan"`
	Sign       string `json:"sign"`
}

func (api *AsiaGamingAPI) Login(account, password, ip string, KindID string) (*LoginReply, error) {

	param := &AGLoginReq{
		CagentType: 1,
		UserName:   account,
		Password:   password,
		GameType:   KindID,
		OddType:    "A",
		Timespan:   time.Now().Format("2006-01-02 15:04:05"),
		Sign:       "",
	}
	data, err := api.request(methodForwardGame, param)
	if err != nil {
		return nil, err
	}
	return &LoginReply{URL: string(data)}, nil
}

//easyjson:json
type AGRegisterReq struct {
	CagentType int    `json:"cagentType"`
	UserName   string `json:"userName"`
	Password   string `json:"password"`
	Timespan   string `json:"timespan"`
	Sign       string `json:"sign"`
}

type AGRegisterRes struct {
	XMLName xml.Name `xml:"result"`
	Info    string   `xml:"info,attr"`
	Msg     string   `xml:"msg,attr"`
}

func (api *AsiaGamingAPI) Register(account, password, ip string) error {
	param := &AGRegisterReq{
		CagentType: 1,
		UserName:   account,
		Password:   password,
		Timespan:   time.Now().Format("2006-01-02 15:04:05"),
		Sign:       "",
	}
	res := &AGRegisterRes{}
	data, err := api.request(methodCheckOrCreateGameAccount, param)
	if err != nil {
		return err
	}
	if err = xml.Unmarshal(data, res); err != nil {
		return err
	}
	if res.Info != "0" {
		return fmt.Errorf(res.Msg)
	}
	return nil
}
func (api *AsiaGamingAPI) Logout(account, password string) {
}

type AGBalanceReq = AGRegisterRes

func (api *AsiaGamingAPI) GetBalance(account, password string) (float64, error) {
	param := &AGRegisterReq{
		CagentType: 1,
		UserName:   account,
		Password:   password,
		Timespan:   time.Now().Format("2006-01-02 15:04:05"),
		Sign:       "",
	}
	res := &AGBalanceReq{}
	data, err := api.request(methodGetBalance, param)
	if err != nil {
		return 0, err
	}
	if err = xml.Unmarshal(data, res); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(res.Info, 64)
}

//easyjson:json
type AGTransferReq struct {
	CagentType  int    `json:"cagentType"`
	UserName    string `json:"userName"`
	Password    string `json:"password"`
	OrderNumber string `json:"orderNumber"`
	Type        string `json:"type"`
	Credit      string `json:"credit"`
	Timespan    string `json:"timespan"`
	Sign        string `json:"sign"`
}
type AGPrepareTransferRes = AGRegisterRes
type AGConfirmTransferRes = AGRegisterRes

//easyjson:json
type AGQueryOrderReq struct {
	CagentType  int    `json:"cagentType"`
	OrderNumber string `json:"orderNumber"`
	Timespan    string `json:"timespan"`
	Sign        string `json:"sign"`
}

func (api *AsiaGamingAPI) TransferIn(account, password, orderNum string, money float64) (tr *TransferInReply, err error) {
	var (
		orderId string
		data    []byte
	)
	// TODO orderid生成规则
	orderId = strings.Replace(time.Now().Format("060102150405.000"), ".", "", -1)

	data, err = api.prepareTransfer(account, password, money, orderId, agTransferIn)
	tr = &TransferInReply{Order: orderId}

	res := &AGPrepareTransferRes{}
	if err = xml.Unmarshal(data, res); err != nil {
		return nil, err
	}
	if res.Info != "0" {
		return nil, fmt.Errorf(res.Msg)
	}

	data, err = api.methodTransferCreditConfirm(account, password, money, orderId, agTransferIn)

	res2 := &AGConfirmTransferRes{}
	if err = xml.Unmarshal(data, res2); err != nil {
		return nil, err
	}
	if res2.Info != "0" {
		return nil, fmt.Errorf(res2.Msg)
	}

	tr.Success = true
	return
}
func (api *AsiaGamingAPI) TransferOut(account, password, orderNum string, money float64) (tr *TransferOutReply, err error) {
	var (
		orderId string
		data    []byte
	)
	// TODO orderid生成规则
	orderId = strings.Replace(time.Now().Format("060102150405.000"), ".", "", -1)

	data, err = api.prepareTransfer(account, password, money, orderId, agTransferOut)
	tr = &TransferOutReply{Order: orderId}

	res := &AGPrepareTransferRes{}
	if err = xml.Unmarshal(data, res); err != nil {
		return nil, err
	}
	if res.Info != "0" {
		return nil, fmt.Errorf(res.Msg)
	}

	data, err = api.methodTransferCreditConfirm(account, password, money, orderId, agTransferOut)
	res2 := &AGConfirmTransferRes{}
	if err = xml.Unmarshal(data, res2); err != nil {
		return nil, err
	}
	if res2.Info != "0" {
		return nil, fmt.Errorf(res2.Msg)
	}

	tr.Success = true
	return
}

func (api *AsiaGamingAPI) prepareTransfer(account, password string, money float64, orderId string, transferType agTransferType) (data []byte, err error) {

	param := &AGTransferReq{
		CagentType:  1,
		UserName:    account,
		Password:    password,
		OrderNumber: orderId,
		Type:        string(transferType),
		Credit:      fmt.Sprintf("%.2f", money),
		Timespan:    time.Now().Format("2006-01-02 15:04:05"),
		Sign:        "",
	}
	data, err = api.request(methodPrepareTransferCredit, param)
	return
}

func (api *AsiaGamingAPI) methodTransferCreditConfirm(account, password string, money float64, orderId string, transferType agTransferType) (data []byte, err error) {

	param := &AGTransferReq{
		CagentType:  1,
		UserName:    account,
		Password:    password,
		OrderNumber: orderId,
		Type:        string(transferType),
		Credit:      fmt.Sprintf("%.2f", money),
		Timespan:    time.Now().Format("2006-01-02 15:04:05"),
		Sign:        "",
	}
	data, err = api.request(methodTransferCreditConfirm, param)
	return
}

func (api *AsiaGamingAPI) request(method string, param json.Marshaler) (data []byte, err error) {
	url := fmt.Sprintf("%s%s", api.baseURL, method)
	req := httplib.Post(url)
	req.Header("Content-Type", "application/json")
	req.Header(KeyHeadAgentCode, api.prefix)

	var buf []byte
	buf, err = json.Marshal(param)
	if err != nil {
		return
	}
	body := string(buf)

	curTime := time.Now().Unix()
	var sign string
	desStr := fmt.Sprintf("%v%v%v%v", curTime, string(buf), api.desKey, api.prefix)
	sign, err = md5hex([]byte(desStr))
	if err != nil {
		return
	}
	agentSign := fmt.Sprintf("%v_%v_%v_%v", api.desKey, sign, curTime, api.prefix)
	req.Header(KeyHeadAgentSign, agentSign)
	common.LogFuncDebug("agentCode:%v,agentSign:%v", api.prefix, agentSign)

	req.Body(body)
	data, err = req.Bytes()
	if err != nil {
		return nil, err
	}
	common.LogFuncDebug("req [%+v] \t param [%+v] \t resp [%+v]", url, body, string(data))

	return
}

//easyjson:json
type AGGetBetRecordsReq struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
	Timespan  string `json:"timespan"`
}

//easyjson:json
type AGLotteryList struct {
	PlayerName     string  `json:"playerName"`
	GameType       string  `json:"gameType"`
	GameName       string  `json:"gameName"`
	OrderId        string  `json:"billNo"`
	TableCode      string  `json:"tableCode"`
	BetAmount      float64 `json:"betAmount"`
	ValidBetAmount float64 `json:"validBetAmount"`
	NetAmount      float64 `json:"netAmount"`
	BetTime        string  `json:"betTime"`
}

//easyjson:json
type AGGetBetRecordsRes struct {
	Info       string          `json:"info"`
	Msg        string          `json:"msg"`
	PageIndex  int             `json:"pageIndex"`
	PageSize   int             `json:"pageSize"`
	TotalCount int             `json:"totalCount"`
	DataList   []AGLotteryList `json:"data"`
}

//查询时间段左闭右开[startTime, endTime]
func (api *AsiaGamingAPI) DayBetRecords(timestamp int64, mapWhite orm.Params) (err error) {
	appChannel, err := admindao.AppChannelDaoEntity.QueryById(GAME_CHANNEL_AG)
	if err != nil || appChannel == nil {
		common.LogFuncError("get channelId[%v] cfg fail", GAME_CHANNEL_AG)
		return
	}
	err = reportdao.ReportGameRecordAgDaoEntity.DelByTimestamp(timestamp)
	if err != nil {
		common.LogFuncError("DayBetRecords error:%v", err)
	}
	startTime := common.GetTimeFormat(timestamp)
	endTime := common.GetTimeFormat(timestamp + common.DaySeconds - 1)
	common.LogFuncDebug("startTime:%v,endTime:%v", startTime, endTime)

	curPageIndex := 1
	for {
		curNow := time.Now()
		param := &AGGetBetRecordsReq{
			StartTime: startTime,
			EndTime:   endTime,
			PageIndex: curPageIndex,
			PageSize:  pageLimit,
			Timespan:  curNow.Format("2006-01-02 15:04:05"),
		}

		var data []byte
		data, err = api.request(methodGetBetRecords, param)
		if err != nil {
			return
		}

		res := &AGGetBetRecordsRes{}
		if err = json.Unmarshal(data, res); err != nil {
			common.LogFuncError("DayBetRecords error:%v", err)
			return err
		}
		if res.Info != "0" {
			err = fmt.Errorf(res.Msg)
			return
		}

		var reportGameRecordAgs []reportmodels.ReportGameRecordAg
		var accounts []string
		for _, item := range res.DataList {
			bet := float64(eosplus.QuantityFloat64ToInt64(item.BetAmount))
			ValidBet := float64(eosplus.QuantityFloat64ToInt64(item.ValidBetAmount))
			Profit := float64(eosplus.QuantityFloat64ToInt64(item.NetAmount))
			reportGameRecordAg := reportmodels.ReportGameRecordAg{
				Account:  item.PlayerName,
				GameType: item.GameType,
				GameName: item.GameName,
				OrderId:  item.OrderId,
				TableId:  item.TableCode,
				Bet:      int64(common.GameGameCoin2Eusd(bet, appChannel.ExchangeRate, appChannel.Precision)),      //转换eusd
				ValidBet: int64(common.GameGameCoin2Eusd(ValidBet, appChannel.ExchangeRate, appChannel.Precision)), //转换eusd
				Profit:   int64(common.GameGameCoin2Eusd(Profit, appChannel.ExchangeRate, appChannel.Precision)),   //转换eusd
				BetTime:  item.BetTime,
				Ctime:    timestamp,
			}

			if _, ok := mapWhite[reportGameRecordAg.GameType]; ok {
				//应用白名单
				reportGameRecordAg.ValidBet = 0
			}
			accounts = append(accounts, "\""+reportGameRecordAg.Account+"\"")
			reportGameRecordAgs = append(reportGameRecordAgs, reportGameRecordAg)
		}

		var mapAccounts orm.Params
		mapAccounts, err = gamedao.GameUserDaoEntity.QueryUidByAccounts(GAME_CHANNEL_AG, accounts)
		if err != nil {
			common.LogFuncError("dayLotteryRecords error:%v", err)
			return
		}
		for i := 0; i < len(reportGameRecordAgs); i++ {
			if _, ok := mapAccounts[reportGameRecordAgs[i].Account]; ok {
				strUid := mapAccounts[reportGameRecordAgs[i].Account].(string)
				uid, _ := strconv.ParseUint(strUid, 10, 64)
				reportGameRecordAgs[i].Uid = uid
			}
		}
		//入库
		err = reportdao.ReportGameRecordAgDaoEntity.InsertMul(timestamp, reportGameRecordAgs)
		if err != nil {
			common.LogFuncError("dayLotteryRecords error:%v", err)
			return
		}

		totalPage := res.TotalCount / pageLimit
		if res.TotalCount%pageLimit > 0 {
			totalPage += 1
		}

		if curPageIndex >= totalPage {
			break
		}
		curPageIndex = curPageIndex + 1
	}

	return
}

func (api *AsiaGamingAPI) GetTotalByTimestamp(timestamp int64) (total int, err error) {
	total, err = reportdao.ReportGameRecordAgDaoEntity.QueryTotalByTimestamp(timestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (api *AsiaGamingAPI) GetRecordByTimestamp(timestamp int64, page, limit int) (betInfos []reportdao.BetInfo, err error) {
	betInfos, err = reportdao.ReportGameRecordAgDaoEntity.QueryByTimestamp(timestamp, page, limit)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (api *AsiaGamingAPI) DelBetRecord(timestamp int64) (err error) {
	err = reportdao.ReportGameRecordAgDaoEntity.DelByTimestamp(timestamp)
	return
}
