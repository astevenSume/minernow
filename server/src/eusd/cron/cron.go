package cron

import (
	common2 "common"
	"eusd/eosplus"
	"utils/common"
)

type FunctionContainer struct {
}

var ServerRunning = true

// EOS账号数量检查
func (f *FunctionContainer) AccountNumCheck() (err error) {
	num, _ := common.AppConfigMgr.Int64(common.EosNoUseAccountLimit)
	if num == 0 {
		num = 10
	}
	AccountNumCheck(num)
	return
}

// 交易执行
func (f *FunctionContainer) TransferRun() (err error) {
	if !ServerRunning {
		common2.ServerStopWriteLog("TransferRun")
		return
	}

	err = eosplus.CronRunTxLog()
	return
}

// 交易检查
func (f *FunctionContainer) TransferCheck() (err error) {
	if !ServerRunning {
		common2.ServerStopWriteLog("TransferCheck")
		return
	}
	err = eosplus.CronCheckTransfer()

	return
}

//EOS资产表，同步余额
func (f FunctionContainer) EosWealthFetchAndUpdate() error {
	return EosWealthFetchAndUpdate()
}

//同步EOS资源使用设置
func (f FunctionContainer) SyncRamCpuNetEos() (err error) {
	data, err := common.AppConfigMgr.FetchString(
		common.EosConfigKeyCpuEos,
		common.EosConfigKeyNetEos,
		common.EosConfigKeyRamEos,
	)
	if err != nil {
		common2.LogFuncError("EosRpc Setting ERR:%v", err)
		return
	}
	eosplus.EosConfig.CpuEos = data[common.EosConfigKeyCpuEos]
	eosplus.EosConfig.NetEos = data[common.EosConfigKeyNetEos]
	eosplus.EosConfig.RamEos = data[common.EosConfigKeyRamEos]
	return
}

//同步服务是否进行停止状态
func (f FunctionContainer) SyncServerState() (err error) {
	common2.SyncServerState(common2.ServerEUSD)
	ServerRunning = common2.ServerRunning
	return
}

// send ping message
func (f FunctionContainer) Ping() error {
	return common2.Ping()
}
