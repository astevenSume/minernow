package common

import (
	"common"
)

type UsdtSvr struct {
	common.SvrBase
}

func NewUsdtSvr() *UsdtSvr {
	return &UsdtSvr{
		SvrBase: common.SvrBase{},
	}
}

func (o *UsdtSvr) Init() (err error) {

	if err = o.SvrBase.Init(); err != nil {
		return
	}

	return
}

//current server
var Cursvr *UsdtSvr = NewUsdtSvr()

var ServerOtcState = common.ServerStateRunning
var ServerGameRunning = true
var ServerOtcTradeRunning = true
