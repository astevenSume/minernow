package controllers

import (
	"common"
	usdtsvrcommon "usdtsvr/common"
)

// @Description cron function container
type FunctionContainer struct {
}

func NewFunctionContainer() FunctionContainer {
	return FunctionContainer{}
}

// send ping message
func (f FunctionContainer) Ping() error {
	return common.Ping()
}

//同步server状态
func (fc FunctionContainer) SyncServerState() (err error) {
	usdtsvrcommon.ServerOtcState = common.IsServerStop(common.ServerOtc)
	usdtsvrcommon.ServerGameRunning = common.IsServerStop(common.ServerGame) == common.ServerStateRunning
	usdtsvrcommon.ServerOtcTradeRunning = common.IsServerStop(common.ServerOtcTrade) == common.ServerStateRunning
	return
}
