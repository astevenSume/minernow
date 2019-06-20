package controllers

import (
	"common"
	json2 "encoding/json"
	"eusd/eosplus"
	"fmt"
	otccommon "otc/common"
	. "otc_error"
	"strconv"
	"time"
	admindao "utils/admin/dao"
	"utils/admin/models"
	gamedao "utils/game/dao"
	"utils/game/dao/gameapi"
	gamemodels "utils/game/models"
	otcdao "utils/otc/dao"
	reportdao "utils/report/dao"

	"github.com/astaxie/beego/orm"
	json "github.com/mailru/easyjson"
)

type GameController struct {
	BaseController
}

// game aes key
var GameAesKey = [32]byte{'A', 'k', 'i', 'J', 'D', 'g', 'p', 'Q',
	't', 'R', 'W', 'M', 'i', '1', '4', 'c',
	'Y', 'O', '3', 'c', 'a', 'K', 'c', 'h',
	'Z', 'X', 'V', 'F', 'G', 'L', 'a', 'b'}

func syncGameUserDaily(timestamp, now int64, start, end string) (err error) {
	err = gamedao.GameUserDailyDaoEntity.Sync(otccommon.Cursvr.ChannelUrl, otccommon.Cursvr.PlatId, timestamp, now, start, end)
	return
}

func syncGameDaily(timestamp, now int64, start, end string) (integer, decimals int32, err error) {
	integer, decimals, err = gamedao.GameDailyDaoEntity.Sync(otccommon.Cursvr.ChannelUrl, otccommon.Cursvr.PlatId, timestamp, now, start, end)
	return
}

//easyjson:json
type GameLoginMsg struct {
	ChannelId uint32 `json:"channel_id"`
	MaskId    string `json:"app_id"`
}

// Login 游戏登录
func (c *GameController) Login() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := GameLoginMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	if msg.ChannelId <= 0 || len(msg.MaskId) <= 0 {
		common.LogFuncDebug("parameters error : channel_id %d, mash_id %s", msg.ChannelId, msg.MaskId)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	// try get game channel user
	var (
		gameUser gamemodels.GameUser
	)
	gameUser, err = gamedao.GameUserDaoEntity.QueryByUid(uid, msg.ChannelId)
	if err != nil {
		if err == orm.ErrNoRows { //try register
			gameUser, errCode = c.register(uid, msg.ChannelId)
			if errCode != ERROR_CODE_SUCCESS {
				c.ErrorResponse(errCode)
				return
			}
		} else {
			c.ErrorResponse(ERROR_CODE_DB)
			return
		}
	}

	// try login
	var (
		reply *gameapi.LoginReply
	)

	api := gameapi.GetApi(msg.ChannelId)
	if api == nil {
		c.ErrorResponse(ERROR_CODE_GAME_UNSUPPORTED_CHANNEL)
		return
	}

	reply, err = api.Login(gameUser.Account, gameUser.Password, common.ClientIP(c.Ctx), msg.MaskId)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_GAME_LOGIN_FAILED)
		return
	}

	common.RedisManger.Del(c.generateBalanceRedisKey(uid))

	c.SuccessResponse(map[string]interface{}{
		KeyGameUrl: reply.URL,
	})
}

type GameBalance struct {
	Channel   string `json:"channel_name"`
	ChannelID uint32 `json:"channel_id"`
	Value     string `json:"value"`
}

func (c *GameController) generateBalanceRedisKey(uid uint64) string {
	return fmt.Sprintf("game_balancekey:%v", uid)
}

// GetBalance 根据 channel 获取账号余额
func (c *GameController) GetBalance() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	var (
		result []GameBalance
		err    error
		data   []byte
		force  bool
	)

	force, _ = c.GetBool("force", false)

	if !force {
		data, err = common.RedisManger.Get(c.generateBalanceRedisKey(uid)).Bytes()
		if err == nil {
			if err = json2.Unmarshal(data, &result); err == nil {
				c.SuccessResponse(result)
				return
			}
		}
	}

	var (
		accountList []gamemodels.GameUser
		balance     float64
		appChannel  *models.AppChannel
	)

	accountList, err = gamedao.GameUserDaoEntity.QueryAllByUid(uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	result = make([]GameBalance, 0, len(accountList))

	for _, gameUser := range accountList {

		api := gameapi.GetApi(gameUser.ChannelId)
		if api == nil {
			continue
		}
		balance, err = api.GetBalance(gameUser.Account, gameUser.Password)
		if err != nil {
			continue
		}

		appChannel, err = admindao.AppChannelDaoEntity.QueryById(gameUser.ChannelId)
		if err != nil || appChannel == nil {
			continue
		}

		if appChannel.ExchangeRate <= 0 || appChannel.Precision <= 0 {
			continue
		}

		eusd := common.GameGameCoin2Eusd(balance, appChannel.ExchangeRate, appChannel.Precision)

		result = append(result, GameBalance{
			Channel:   appChannel.Name,
			ChannelID: appChannel.Id,
			Value:     strconv.FormatFloat(eusd, 'f', -1, 64),
		})

	}

	if data, err = json2.Marshal(&result); err == nil {
		common.RedisManger.Set(c.generateBalanceRedisKey(uid), data, time.Second*300).Err()
	}

	c.SuccessResponse(result)
}

//easyjson:json
type GameTransferInMsg struct {
	ChannelId uint32 `json:"channel_id"`
	Eusd      string `json:"eusd"`
}

func (c *GameController) TransferIn() {

	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := GameTransferInMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	appChannel, err := admindao.AppChannelDaoEntity.QueryById(msg.ChannelId)
	if err != nil || appChannel == nil {
		common.LogFuncDebug("game channel: %v , error: %v", msg.ChannelId, err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	if appChannel.ExchangeRate <= 0 || appChannel.Precision <= 0 {
		common.LogFuncDebug("ExchangeRate: %v , Precision: %v", appChannel.ExchangeRate, appChannel.Precision)
		c.ErrorResponse(ERROR_CODE_GAME_CHANNEL_UNSET_RATE)
		return
	}

	var (
		gameUser gamemodels.GameUser
	)
	gameUser, err = gamedao.GameUserDaoEntity.QueryByUid(uid, msg.ChannelId)
	if err != nil {
		if err == orm.ErrNoRows {
			c.ErrorResponse(ERROR_CODE_GAME_USER_NO_EXIST)
			return
		} else {
			c.ErrorResponse(ERROR_CODE_DB)
			return
		}
	}

	api := gameapi.GetApi(msg.ChannelId)
	if api == nil {
		c.ErrorResponse(ERROR_CODE_GAME_UNSUPPORTED_CHANNEL)
		return
	}

	eusdInt64, err := common.CurrencyStrToInt64(msg.Eusd)
	if err != nil {
		common.LogFuncDebug("str2int64 %s failed %v", msg.Eusd, err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	// freeze and get user eusd balance
	var eusdFloat64 float64
	eusdFloat64 = common.CurrencyInt64ToFloat64(eusdInt64)
	// 锁定 eusd ,可能会抹掉部分小数位
	errCode = eosplus.EosPlusAPI.Game.Recharge(gameUser.Uid, eusdFloat64)
	common.LogFuncDebug("%s eusd balance : %v", gameUser.Account, eusdFloat64)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	// transfer from eusd to coin
	coin := common.GameEusd2GameCoin(eusdFloat64, appChannel.ExchangeRate, appChannel.Precision)
	if coin <= 0 {
		// 转账金额错误,解锁 eusd
		errCode = eosplus.EosPlusAPI.Game.Withdrawal(gameUser.Uid, eusdFloat64)
		if errCode != ERROR_CODE_SUCCESS {
			common.LogFuncDebug("Withdrawal failed : %v", errCode)
		}
		c.ErrorResponse(ERROR_CODE_USDT_CURRENCY_PARAM_ERROR)
		return
	}

	transfer, err := gamedao.GameTransferDaoEntity.Add(gameUser.Uid, gameUser.ChannelId, gameUser.Account, gamedao.TRANSFER_TYPE_IN, common.CurrencyFloat64ToInt64(coin), eusdInt64, "")
	if err != nil {
		// 写入订单失败,解锁 eusd
		errCode = eosplus.EosPlusAPI.Game.Withdrawal(gameUser.Uid, eusdFloat64)
		if errCode != ERROR_CODE_SUCCESS {
			common.LogFuncDebug("Withdrawal failed : %v", errCode)
		}
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	var reply *gameapi.TransferInReply
	reply, err = api.TransferIn(gameUser.Account, gameUser.Password, transfer.Order, coin)
	if reply != nil {
		// 游戏平台订单号
		transfer.GameOrder = reply.Order
	}
	if err == nil && reply.Success {
		// 订单成功设置订单状态为成功
		transfer.Status = uint32(gamedao.TRANSFER_STATUS_DONE)
	} else {
		transfer.Status = uint32(gamedao.TRANSFER_STATUS_FAILED)
		errCode = eosplus.EosPlusAPI.Game.Withdrawal(gameUser.Uid, eusdFloat64)
		if errCode != ERROR_CODE_SUCCESS {
			common.LogFuncDebug("Withdrawal failed : %v", errCode)
		}
	}

	// 更新订单信息: 订单号,订单状态
	dberr := gamedao.GameTransferDaoEntity.Update(transfer)

	if err != nil {
		common.LogFuncDebug("transfer in failed : %v", err)
		c.ErrorResponse(ERROR_CODE_GAME_ADD_REDEEM_LOG_FAILED)
		return
	}
	if dberr != nil {
		common.LogFuncDebug("transfer update game order failed : %v", transfer.GameOrder)
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	//充值个人日报
	_ = reportdao.ReportGameTransferDailyDaoEntity.Recharge(msg.ChannelId, uid, eosplus.QuantityStringToint64(msg.Eusd))

	common.RedisManger.Del(c.generateBalanceRedisKey(uid))

	c.SuccessResponse(true)
	return

}

//easyjson:json
type GameTransferOutMsg struct {
	ChannelId uint32 `json:"channel_id"`
	Eusd      string `json:"eusd"`
}

func (c *GameController) TransferOut() {

	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := GameTransferOutMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	appChannel, err := admindao.AppChannelDaoEntity.QueryById(msg.ChannelId)
	if err != nil || appChannel == nil {
		common.LogFuncDebug("game channel: %v , error: %v", msg.ChannelId, err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	if appChannel.ExchangeRate <= 0 || appChannel.Precision <= 0 {
		common.LogFuncDebug("ExchangeRate: %v , Precision: %v", appChannel.ExchangeRate, appChannel.Precision)
		c.ErrorResponse(ERROR_CODE_GAME_CHANNEL_UNSET_RATE)
		return
	}

	var (
		gameUser gamemodels.GameUser
	)
	gameUser, err = gamedao.GameUserDaoEntity.QueryByUid(uid, msg.ChannelId)
	if err != nil {
		if err == orm.ErrNoRows {
			c.ErrorResponse(ERROR_CODE_GAME_USER_NO_EXIST)
			return
		} else {
			c.ErrorResponse(ERROR_CODE_DB)
			return
		}
	}

	api := gameapi.GetApi(msg.ChannelId)
	if api == nil {
		c.ErrorResponse(ERROR_CODE_GAME_UNSUPPORTED_CHANNEL)
		return
	}

	eusdInt64, err := common.CurrencyStrToInt64(msg.Eusd)
	if err != nil {
		common.LogFuncDebug("str2int64 %s failed %v", msg.Eusd, err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	var eusdFloat64 float64
	eusdFloat64 = common.CurrencyInt64ToFloat64(eusdInt64)

	coin := common.GameEusd2GameCoin(eusdFloat64, appChannel.ExchangeRate, appChannel.Precision)
	if coin <= 0 {
		c.ErrorResponse(ERROR_CODE_USDT_CURRENCY_PARAM_ERROR)
		return
	}

	transfer, err := gamedao.GameTransferDaoEntity.Add(gameUser.Uid, gameUser.ChannelId, gameUser.Account, gamedao.TRANSFER_TYPE_OUT, common.CurrencyFloat64ToInt64(coin), eusdInt64, "")
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	var reply *gameapi.TransferOutReply
	reply, err = api.TransferOut(gameUser.Account, gameUser.Password, transfer.Order, coin)
	if reply != nil {
		// 游戏平台订单号
		transfer.GameOrder = reply.Order
	}
	if err == nil && reply.Success {
		// 订单成功设置订单状态为成功
		transfer.Status = uint32(gamedao.TRANSFER_STATUS_DONE)
	} else {
		transfer.Status = uint32(gamedao.TRANSFER_STATUS_FAILED)
	}
	// 更新订单信息: 订单号,订单状态
	dberr := gamedao.GameTransferDaoEntity.Update(transfer)
	if err != nil {
		common.LogFuncDebug("transfer in failed : %v", err)
		c.ErrorResponse(ERROR_CODE_GAME_ADD_REDEEM_LOG_FAILED)
		return
	}
	if dberr != nil {
		common.LogFuncDebug("transfer update game order failed : %v", transfer.GameOrder)
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	// 解冻
	errCode = eosplus.EosPlusAPI.Game.Withdrawal(gameUser.Uid, eusdFloat64)
	if errCode != ERROR_CODE_SUCCESS {
		common.LogFuncDebug("Withdrawal:uid %d , order %s , errCode %v", gameUser.Uid, transfer.GameOrder, errCode)
		c.ErrorResponse(errCode)
		return
	}

	//充值个人日报
	_ = reportdao.ReportGameTransferDailyDaoEntity.Withdraw(msg.ChannelId, uid, eosplus.QuantityStringToint64(msg.Eusd))

	common.RedisManger.Del(c.generateBalanceRedisKey(uid))

	c.SuccessResponse(true)
	return

}

func (c *GameController) register(uid uint64, channelId uint32) (gameUser gamemodels.GameUser, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	// generate game user
	user, err := otcdao.UserDaoEntity.InfoByUId(uid)
	if err != nil {
		common.LogFuncError("%v", err)
		errCode = ERROR_CODE_DB
		return
	}
	if user == nil {
		common.LogFuncError("user %d no exists.", uid)
		errCode = ERROR_CODE_NO_USER
		return
	}

	gameUser.Uid = uid
	gameUser.ChannelId = channelId
	gameUser.Account = fmt.Sprintf("U%d", uid)
	gameUser.Password = common.RandomStr(16)
	gameUser.NickName = gameUser.Account
	gameUser.Sex = 1

	api := gameapi.GetApi(channelId)
	if api == nil {
		errCode = ERROR_CODE_GAME_UNSUPPORTED_CHANNEL
		return
	}

	err = api.Register(gameUser.Account, gameUser.Password, common.ClientIP(c.Ctx))
	if err != nil {
		errCode = ERROR_CODE_GAME_LOGIN_FAILED
		return
	}
	if err != nil {
		common.LogFuncError("%v", err)
		c.ErrorResponse(ERROR_CODE_ENCRYPT_FAILED)
		return
	}

	// save
	err = gamedao.GameUserDaoEntity.Add(gameUser)
	if err != nil {
		errCode = ERROR_CODE_GAME_ADD_USER_FAILED
		return
	}

	return
}

func (c *GameController) Logout() {
	// uid, errCode := c.getUidFromToken()
	// if errCode != ERROR_CODE_SUCCESS {
	// 	c.ErrorResponse(errCode)
	// 	return
	// }

	// msg := GameLogoutMsg{}
	// err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	// if err != nil {
	// 	common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
	// 	c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
	// 	return
	// }

	// var gameUser gamemodels.GameUser
	// gameUser, err = gamedao.GameUserDaoEntity.QueryByUid(uid, msg.ChannelId)
	// if err != nil {
	// 	if err == orm.ErrNoRows {
	// 		c.ErrorResponse(ERROR_CODE_GAME_USER_NO_EXIST)
	// 	} else {
	// 		c.ErrorResponse(ERROR_CODE_DB)
	// 	}
	// 	return
	// }

	// req := GameLogoutApiMsg{
	// 	Account: gameUser.Account,
	// }

	// resp := GameLoginApiRespMsg{}
	// err = gamedao.GameDaoEntity.SendToChannel(otccommon.Cursvr.ChannelUrl, gamedao.PROTO_LOGOUT, otccommon.Cursvr.PlatId, &req, &resp)
	// if err != nil {
	// 	c.ErrorResponse(ERROR_CODE_GAME_LOGIN_FAILED)
	// 	return
	// }

	// if resp.Status != int(ERROR_CODE_SUCCESS) {
	// 	common.LogFuncError("resp status %d : %s", resp.Status, resp.Desc)
	// 	c.ErrorResponse(ERROR_CODE_GAME_LOGIN_FAILED)
	// 	return
	// }

	// errCode = c.tryRedeemOut(gameUser, RedeemOutLogout)
	// if errCode != ERROR_CODE_SUCCESS {
	// 	c.ErrorResponse(errCode)
	// 	return
	// }

	// c.SuccessResponseWithoutData()
}

type TransferListMsg struct {
	ChannelId string `form:"channel_id"`
	Status    string `form:"status"`
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
}
type GameTransfer struct {
	Id           uint64 `json:"id"`
	ChannelId    uint32 `json:"channel_id"`
	TransferType uint32 `json:"transfer_type"`
	Eusd         string `json:"eusd"`
	Status       uint32 `json:"status"`
	Ctime        int64  `json:"ctime"`
}

func (c *GameController) TransferList() {

	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	var (
		msg     TransferListMsg
		err     error
		total   int64
		data    []gamemodels.GameTransfer
		records []GameTransfer
	)

	if err = c.ParseForm(&msg); err != nil {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	if total, data, err = gamedao.GameTransferDaoEntity.QueryPageTransferList(uid, msg.ChannelId, msg.Status, msg.Page, msg.Limit); err != nil {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	records = make([]GameTransfer, 0, len(data))
	for _, v := range data {
		records = append(records, GameTransfer{
			Id:           v.Id,
			ChannelId:    v.ChannelId,
			TransferType: v.TransferType,
			Eusd:         fmt.Sprintf("%.4f", common.CurrencyInt64ToFloat64(v.EusdInteger)),
			Status:       v.Status,
			Ctime:        v.Ctime,
		})
	}
	c.SuccessResponse(map[string]interface{}{
		KeyList: records,
		KeyMeta: MetaMsg{
			Total: int(total),
			Page:  msg.Page,
			Limit: msg.Limit,
		},
	})

}

func (c *GameController) TransferDetail() {

}
