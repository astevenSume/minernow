package eosplus

import (
	"common"
)

type EusdSvr struct {
	common.SvrBase
}

func NewEusdSvr() *EusdSvr {
	return &EusdSvr{
		SvrBase: common.SvrBase{},
	}
}

func (o *EusdSvr) Init() (err error) {

	if err = o.SvrBase.Init(); err != nil {
		return
	}

	return
}

//current server
var Cursvr *EusdSvr = NewEusdSvr()
