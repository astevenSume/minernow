package gameapi

import (
	"common"
	"encoding/base64"
	"encoding/json"
	"eusd/eosplus"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"

	admindao "utils/admin/dao"
	gamedao "utils/game/dao"
	reportdao "utils/report/dao"
	reportmodels "utils/report/models"

	"github.com/astaxie/beego/httplib"
)

const ErrNotFind = "Not Found"

type KaiYuanAPI struct {
	baseURL, recordURL, agent, desKey, md5Key, lineCode string
	agentLen                                            int
}

//easyjson:json
type KYBaseResp struct {
	S    int    `json:"s"`
	M    string `json:"m"`
	Code int    `json:"code"`
}

//easyjson:json
type KYLoginResp struct {
	KYBaseResp
	D KYLoginResult `json:"d"`
}

//easyjson:json
type KYLoginResult struct {
	Code int    `json:"code"`
	URL  string `json:"url"`
}

//easyjson:json
type KYAccountResp struct {
	KYBaseResp
	D KYAccountResult `json:"d"`
}

//easyjson:json
type KYAccountResult struct {
	Code    int     `json:"code"`
	Money   float64 `json:"money"`
	Account string  `json:"account"`
}

//easyjson:json
type KYTransferInResp struct {
	KYBaseResp
	D KYTransferInResult `json:"d"`
}

//easyjson:json
type KYTransferInResult struct {
	Code    int    `json:"code"`
	Money   string `json:"money"`
	Account string `json:"account"`
}

//easyjson:json
type KYTransferOutResp struct {
	KYBaseResp
	D KYTransferOutResult `json:"d"`
}

//easyjson:json
type KYTransferOutResult struct {
	Code    int    `json:"code"`
	Money   string `json:"money"`
	Account string `json:"account"`
}

/*
"测试环境：
接口地址：https://kyapi.ky206.com:189/channelHandle
明细接口地址:https://kyrecord.ky206.com:190/getRecordHandle
正式环境：
接口地址：https://api.ky195.com:189/channelHandle"
拉单接口地址:https://record.ky013.com:190/getRecordHandle
*/
// NewKaiYuanAPI 开元棋牌 api
func NewKaiYuanAPI(baseURL, recordURL, agent, desKey, md5Key, lineCode string) GameAPI {
	return &KaiYuanAPI{
		baseURL:   baseURL,
		recordURL: recordURL,
		agent:     agent,
		desKey:    desKey,
		md5Key:    md5Key,
		lineCode:  lineCode,
		agentLen:  len(agent) + 1,
	}
}

const (
	subTypeLogin         = 0
	subTypeQueryBalance  = 1
	subTypeTransferIn    = 2
	subTypeTransferOut   = 3
	subTypeCheckTransfer = 4
	subTypeGetRecord     = 6
)

// Login 登陆
func (api *KaiYuanAPI) Login(account, password, ip string, KindID string) (*LoginReply, error) {

	ymd, timestamp := now()

	orderid := fmt.Sprintf("%s%s%s", api.agent, ymd, account)

	if len(orderid) > 100 {
		return nil, fmt.Errorf("orderid length out of range")
	}

	param := strings.Join([]string{
		fmt.Sprintf("s=%v", subTypeLogin),
		fmt.Sprintf("account=%v", account),
		fmt.Sprintf("money=%v", 0),
		fmt.Sprintf("orderid=%v", orderid), //订单id agent+timestamp+account
		fmt.Sprintf("ip=%v", ip),
		fmt.Sprintf("lineCode=%v", api.lineCode),
		fmt.Sprintf("KindID=%v", KindID), // 0 大厅
	}, "&")

	data, err := common.AesECBEncrypt(param, api.desKey)
	if err != nil {
		return nil, err
	}

	key, err := md5hex([]byte(api.agent), []byte(timestamp), []byte(api.md5Key))
	if err != nil {
		return nil, err
	}

	res := &KYLoginResp{}

	api.request(timestamp, base64.StdEncoding.EncodeToString(data), key, res)
	if res.D.Code != 0 {
		return nil, kyError(res.D.Code)
	}
	return &LoginReply{URL: res.D.URL}, nil
}

func (api *KaiYuanAPI) Register(account, password, ip string) error { return nil }
func (api *KaiYuanAPI) Logout(account, password string)             {}

// GetAccount 查询账号信息
func (api *KaiYuanAPI) GetBalance(account, password string) (float64, error) {

	_, timestamp := now()

	param := strings.Join([]string{
		fmt.Sprintf("s=%v", subTypeQueryBalance),
		fmt.Sprintf("account=%v", account),
	}, "&")

	data, err := common.AesECBEncrypt(param, api.desKey)
	if err != nil {
		return 0, err
	}

	key, err := md5hex([]byte(api.agent), []byte(timestamp), []byte(api.md5Key))
	if err != nil {
		return 0, err
	}
	res := &KYAccountResp{}

	if err = api.request(timestamp, base64.StdEncoding.EncodeToString(data), key, res); err != nil {
		return 0, err
	}
	if res.D.Code != 0 {
		return 0, kyError(res.D.Code)
	}
	return res.D.Money, nil

}
func (api *KaiYuanAPI) TransferIn(account, password, orderNum string, money float64) (tr *TransferInReply, err error) {

	ymd, timestamp := now()

	orderid := fmt.Sprintf("%s%s%s", api.agent, ymd, account)
	tr = &TransferInReply{Order: orderid}
	if len(orderid) > 100 {
		err = fmt.Errorf("orderid length out of range")
		return
	}

	param := strings.Join([]string{
		fmt.Sprintf("s=%v", subTypeTransferIn),
		fmt.Sprintf("account=%v", account),
		fmt.Sprintf("money=%v", money),
		fmt.Sprintf("orderid=%v", orderid), //订单id agent+timestamp+account
	}, "&")
	var (
		data []byte
		key  string
	)
	data, err = common.AesECBEncrypt(param, api.desKey)
	if err != nil {
		return
	}

	key, err = md5hex([]byte(api.agent), []byte(timestamp), []byte(api.md5Key))
	if err != nil {
		return
	}

	res := &KYTransferInResp{}
	err = api.request(timestamp, base64.StdEncoding.EncodeToString(data), key, res)
	if err == nil {
		if res.D.Code != 0 {
			err = kyError(res.D.Code)
		} else {
			tr.Success = true
			return
		}
	}
	// 上下分接口返回如果有异常,则定时检查订单状态
	if err != nil {
		tr.Success = api.checkOrderStatus(orderid)
	}
	return
}

func (api *KaiYuanAPI) TransferOut(account, password, orderNum string, money float64) (tr *TransferOutReply, err error) {

	ymd, timestamp := now()

	orderid := fmt.Sprintf("%s%s%s", api.agent, ymd, account)
	tr = &TransferOutReply{Order: orderid}
	if len(orderid) > 100 {
		return nil, fmt.Errorf("orderid length out of range")
	}

	param := strings.Join([]string{
		fmt.Sprintf("s=%v", subTypeTransferOut),
		fmt.Sprintf("account=%v", account),
		fmt.Sprintf("money=%v", money),
		fmt.Sprintf("orderid=%v", orderid), //订单id agent+timestamp+account
	}, "&")
	var (
		data []byte
		key  string
	)
	data, err = common.AesECBEncrypt(param, api.desKey)
	if err != nil {
		return
	}

	key, err = md5hex([]byte(api.agent), []byte(timestamp), []byte(api.md5Key))
	if err != nil {
		return
	}

	res := &KYTransferOutResp{}

	err = api.request(timestamp, base64.StdEncoding.EncodeToString(data), key, res)
	if err == nil {
		if res.D.Code != 0 {
			err = kyError(res.D.Code)
		} else {
			tr.Success = true
			return
		}
	}
	// 上下分接口返回如果有异常,则定时检查订单状态
	if err != nil {
		tr.Success = api.checkOrderStatus(orderid)
	}
	return
}

func (api *KaiYuanAPI) checkOrderStatus(orderid string) bool {
	var (
		status int
		retry  int
		err    error
	)
	for {
		// 查询订单结果
		status, err = api.transferInfo(orderid)
		retry++

		if err == nil {
			switch status {
			// （-1:不存在、0:成功、2:失败、3:处理中）
			// 成功或者失败就可以直接返回
			case 0:
				return true
			case 2:
				return false
			}
		}

		// 重试次数 3 次则退出
		if retry > 2 {
			return false
		}
		// 间隔5秒查询一次
		time.Sleep(time.Duration(5) * time.Second)
	}
}

//easyjson:json
type KYTransferInfoResp struct {
	KYBaseResp
	D KYTransferInfoResult `json:"d"`
}

//easyjson:json
type KYTransferInfoResult struct {
	Code   int     `json:"code"`
	Status int     `json:"status"`
	Money  float64 `json:"money"`
}

func (api *KaiYuanAPI) transferInfo(orderid string) (status int, err error) {

	param := strings.Join([]string{
		fmt.Sprintf("s=%v", subTypeCheckTransfer),
		fmt.Sprintf("orderid=%v", orderid), //订单id agent+timestamp+account
	}, "&")
	var (
		data []byte
		key  string
	)
	data, err = common.AesECBEncrypt(param, api.desKey)
	if err != nil {
		return
	}

	_, timestamp := now()

	key, err = md5hex([]byte(api.agent), []byte(timestamp), []byte(api.md5Key))
	if err != nil {
		return
	}

	res := &KYTransferInfoResp{}
	err = api.request(timestamp, base64.StdEncoding.EncodeToString(data), key, res)
	if err != nil {
		return
	}
	if res.D.Code != 0 {
		err = kyError(res.D.Code)
	} else {
		status = res.D.Status
	}

	return
}

func now() (ymd string, timestamp string) {
	now := time.Now()
	return strings.Replace(now.Format("20060102150405.000"), ".", "", -1), fmt.Sprint(now.UnixNano() / 1000000)
}

func (api *KaiYuanAPI) request(timestamp, param, key string, reply json.Unmarshaler) (err error) {
	reqUrl := fmt.Sprintf("%s/channelHandle?agent=%s&timestamp=%s&param=%s&key=%s", api.baseURL, api.agent, timestamp, url.QueryEscape(param), key)
	req := httplib.Get(reqUrl)

	var data []byte
	data, err = req.Bytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, reply)
	if err != nil {
		return
	}

	common.LogFuncDebug("req [%+v] \t resp [%+v] \t jsonData [%+v]", reqUrl, string(data), reply)
	return
}
func kyError(errCode int) error {
	return fmt.Errorf("kyapi return errCode %v", errCode)
}

//easyjson:json
type KYRecordBaseResp struct {
	S int    `json:"s"`
	M string `json:"m"`
}

//easyjson:json
type KYSyncResp struct {
	KYRecordBaseResp
	Result KYSyncResult `json:"d"`
}

//easyjson:json
type KYSyncResult struct {
	Code     int      `json:"code"`
	Count    int      `json:"count"`
	Start    int64    `json:"start"`
	End      int64    `json:"end"`
	DataList KYRecord `json:"list"`
}

//easyjson:json
type KYRecord struct {
	GameID        []string `json:"GameID"`
	Accounts      []string `json:"Accounts"`
	ServerID      []int    `json:"ServerID"`
	KindID        []int    `json:"KindID"`
	TableID       []int    `json:"TableID"`
	ChairID       []int    `json:"ChairID"`
	UserCount     []int    `json:"UserCount"`
	CardValue     []string `json:"CardValue"`
	CellScore     []string `json:"CellScore"`
	AllBet        []string `json:"AllBet"`
	Profit        []string `json:"Profit"`
	Revenue       []string `json:"Revenue"`
	GameStartTime []string `json:"GameStartTime"`
	GameEndTime   []string `json:"GameEndTime"`
	ChannelID     []int    `json:"ChannelID"`
	LineCode      []string `json:"LineCode"`
}

func (api *KaiYuanAPI) requestRecord(timestamp, param, key string, reply json.Unmarshaler) (err error) {
	reqUrl := fmt.Sprintf("%s/getRecordHandle?agent=%s&timestamp=%s&param=%s&key=%s", api.recordURL, api.agent, timestamp, url.QueryEscape(param), key)
	req := httplib.Get(reqUrl)

	var data []byte
	data, err = req.Bytes()
	if err != nil {
		return err
	}

	if string(data) == ErrNotFind {
		common.LogFuncDebug("req [%+v] \t resp [%+v]", reqUrl, string(data))
		return
	}

	err = json.Unmarshal(data, &reply)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	common.LogFuncDebug("req [%+v] \t resp [%+v] \t jsonData [%+v]", reqUrl, string(data), reply)
	return
}

//闭区间 [startTime,endTime]
func (api *KaiYuanAPI) GetBetRecord(startTime, endTime int64) (res KYSyncResp, err error) {
	common.LogFuncInfo("starttime:%v", startTime)
	_, now := now()
	param := strings.Join([]string{
		fmt.Sprintf("s=%v", subTypeGetRecord),
		fmt.Sprintf("startTime=%v", startTime),
		fmt.Sprintf("endTime=%v", endTime),
	}, "&")
	common.LogFuncDebug("param:%v", param)

	var data []byte
	data, err = common.AesECBEncrypt(param, api.desKey)
	if err != nil {
		return
	}

	var key string
	key, err = md5hex([]byte(api.agent), []byte(now), []byte(api.md5Key))
	if err != nil {
		return
	}
	err = api.requestRecord(now, base64.StdEncoding.EncodeToString(data), key, &res)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (api *KaiYuanAPI) GetAccount(channelAccount string) (account string) {
	return
}

//查询时间段双闭[startTime, endTime]
func (api *KaiYuanAPI) DayBetRecords(timestamp int64, mapWhite orm.Params) (err error) {
	appChannel, err := admindao.AppChannelDaoEntity.QueryById(GAME_CHANNEL_KY)
	if err != nil || appChannel == nil {
		common.LogFuncError("get channelId[%v] cfg fail", GAME_CHANNEL_KY)
		return
	}

	err = reportdao.ReportGameRecordKyDaoEntity.DelByTimestamp(timestamp)
	if err != nil {
		common.LogFuncError("DayBetRecords error:%v", err)
	}
	mapAppName, _ := admindao.AppDaoEntity.GetAppNameByChannelId(GAME_CHANNEL_KY)

	startTime := timestamp * 1000 //毫秒
	endTime := startTime + ReqTimeSlot
	initTime := startTime //毫秒
	finishTime := (timestamp + common.DaySeconds) * 1000
	page := 1
	for {
		var res KYSyncResp
		res, err = api.GetBetRecord(startTime, endTime)
		if err != nil {
			return
		}

		var reportGameRecordKys []reportmodels.ReportGameRecordKy
		var accounts []string
		for i := 0; i < res.Result.Count; i++ {
			bet := float64(eosplus.QuantityStringToint64(res.Result.DataList.AllBet[i]))
			ValidBet := float64(eosplus.QuantityStringToint64(res.Result.DataList.CellScore[i]))
			Profit := float64(eosplus.QuantityStringToint64(res.Result.DataList.Profit[i]))
			Revenue := float64(eosplus.QuantityStringToint64(res.Result.DataList.Revenue[i]))
			reportGameRecordKy := reportmodels.ReportGameRecordKy{
				Account:   res.Result.DataList.Accounts[i][api.agentLen:],
				GameId:    res.Result.DataList.GameID[i],
				ServerId:  int32(res.Result.DataList.ServerID[i]),
				KindId:    fmt.Sprintf("%v", res.Result.DataList.KindID[i]),
				TableId:   int32(res.Result.DataList.TableID[i]),
				ChairId:   int32(res.Result.DataList.ChairID[i]),
				Bet:       int64(common.GameGameCoin2Eusd(bet, appChannel.ExchangeRate, appChannel.Precision)),      //转换eusd,
				ValidBet:  int64(common.GameGameCoin2Eusd(ValidBet, appChannel.ExchangeRate, appChannel.Precision)), //转换eusd,
				Profit:    int64(common.GameGameCoin2Eusd(Profit, appChannel.ExchangeRate, appChannel.Precision)),   //转换eusd,
				Revenue:   int64(common.GameGameCoin2Eusd(Revenue, appChannel.ExchangeRate, appChannel.Precision)),  //转换eusd,
				StartTime: res.Result.DataList.GameStartTime[i],
				EndTime:   res.Result.DataList.GameEndTime[i],
				Ctime:     timestamp,
			}

			if v, ok := mapAppName[reportGameRecordKy.KindId]; ok {
				reportGameRecordKy.GameName = v.(string)
			}

			if _, ok := mapWhite[reportGameRecordKy.KindId]; ok {
				//应用白名单
				reportGameRecordKy.ValidBet = 0
			}
			accounts = append(accounts, "\""+reportGameRecordKy.Account+"\"")
			reportGameRecordKys = append(reportGameRecordKys, reportGameRecordKy)
		}

		var mapAccounts orm.Params
		mapAccounts, err = gamedao.GameUserDaoEntity.QueryUidByAccounts(GAME_CHANNEL_KY, accounts)
		if err != nil {
			common.LogFuncError("dayLotteryRecords error:%v", err)
			return
		}
		for i := 0; i < len(reportGameRecordKys); i++ {
			if _, ok := mapAccounts[reportGameRecordKys[i].Account]; ok {
				strUid := mapAccounts[reportGameRecordKys[i].Account].(string)
				uid, _ := strconv.ParseUint(strUid, 10, 64)
				reportGameRecordKys[i].Uid = uid
			}
		}
		//入库
		err = reportdao.ReportGameRecordKyDaoEntity.InsertMul(timestamp, reportGameRecordKys)
		if err != nil {
			common.LogFuncError("dayLotteryRecords error:%v", err)
			return
		}

		if startTime >= finishTime {
			break
		}
		startTime = initTime + int64(ReqTimeSlot*page+1)
		endTime = initTime + int64(ReqTimeSlot*(page+1))
		page++
		time.Sleep(ReqInterval)
	}

	return
}

func (api *KaiYuanAPI) GetTotalByTimestamp(timestamp int64) (total int, err error) {
	total, err = reportdao.ReportGameRecordKyDaoEntity.QueryTotalByTimestamp(timestamp)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (api *KaiYuanAPI) GetRecordByTimestamp(timestamp int64, page, limit int) (betInfos []reportdao.BetInfo, err error) {
	betInfos, err = reportdao.ReportGameRecordKyDaoEntity.QueryByTimestamp(timestamp, page, limit)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (api *KaiYuanAPI) DelBetRecord(timestamp int64) (err error) {

	return
}
