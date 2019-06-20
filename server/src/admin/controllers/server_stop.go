package controllers

import (
	"admin/controllers/errcode"
	"common"
	"fmt"
	"utils/eusd/dao"
)

type ServerStopController struct {
	BaseController
}

var serverList = []int{
	common.ServerOtc,
	common.ServerEUSD,
	common.ServerGame,
	common.ServerUSDT,
	common.ServerOtcTrade,
}

func get() map[string]interface{} {

	keys := []string{}
	for _, v := range serverList {
		keys = append(keys, fmt.Sprintf("%s_%d", common.RedisStopKey, v))
	}
	res, err := common.RedisManger.MGet(keys...).Result()
	if err != nil {
		return map[string]interface{}{}
	}
	common.LogFuncDebug("%v", dao.ToJsonIndent(keys))
	common.LogFuncDebug("%v", dao.ToJsonIndent(res))
	list := map[int]interface{}{}
	for k, v := range res {
		if v == nil {
			list[serverList[k]] = 0
		} else {
			list[serverList[k]] = v
		}
	}

	opt := map[int]string{
		common.ServerOtc:      "OTC服务",
		common.ServerEUSD:     "EUSD服务",
		common.ServerGame:     "游戏服务",
		common.ServerUSDT:     "USDT服务",
		common.ServerOtcTrade: "OTC交易",
	}

	return map[string]interface{}{"list": opt, "status": list}
}

func (c *ServerStopController) Get() {
	c.setOPAction(OPActionServerStopGet)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	res := get()

	c.SuccessResponseAndLog(0, "", res)
}

func (c *ServerStopController) Set() {
	c.setOPAction(OPActionServerStopSet)
	c.setRequestData(string(c.Ctx.Input.RequestBody))
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	type post struct {
		Server int  `json:"server"`
		Able   bool `json:"able"`
	}
	msg := &post{}
	err := c.GetPost(msg)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_UNKNOWN)
	}
	if msg.Able == false {
		if msg.Server == common.ServerOtc {
			err = common.StopServer(common.ServerOtc)
			err = common.StopServer(common.ServerEUSD)
			err = common.StopServer(common.ServerGame)
			err = common.StopServer(common.ServerUSDT)

		} else {
			err = common.StopServer(msg.Server)
		}
	} else {
		err = common.StartServer(msg.Server)
	}

	res := get()
	c.SuccessResponse(res)
}

func (c *ServerStopController) Log() {
	c.setOPAction(OPActionServerStopLog)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	res, err := common.RedisManger.HGetAll(common.RedisStopLogKey).Result()
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_UNKNOWN)
	}
	str := []string{}
	for k, v := range res {
		str = append(str, k+":"+v)
	}

	c.SuccessResponse(str)
}
