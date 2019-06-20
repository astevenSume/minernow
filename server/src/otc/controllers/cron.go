package controllers

import (
	"common"
	common2 "otc/common"
	"otc/cron"
	"usdt"
)

// @Description cron function container
type FunctionContainer struct {
}

func NewFunctionContainer() FunctionContainer {
	return FunctionContainer{}
}

// sync game user daily
func (fc FunctionContainer) SyncGameUserDaily() (err error) {
	//_, _, err = syncGameUserDaily()
	return
}

// sync game daily
func (fc FunctionContainer) SyncGameDaily() (err error) {
	//_, _, _, err = syncGameDaily()
	return
}

//
func (fc FunctionContainer) LoadChip2CommissionConfig() (err error) {
	/*err = Chip2CommissionConfigMgr.Load()
	if err != nil {
		return
	}
	common.LogFuncDebug("%v")*/
	return
}

// detect status of usdt transfer out wealth_log
func (fc FunctionContainer) DetectUsdtWealthLogStatus() error {
	// 同步usdt 交易数据
	usdt.DetectUsdtTransferOutOrder()
	return nil
}

//订单超时检查
func (fc FunctionContainer) OrderCheckTimeOut() error {
	cron.OrderCheckTimeOut()
	return nil
}

//日工资
func (fc FunctionContainer) TaskDailySalary() error {
	TaskDailySalary()
	return nil
}

//AG日工资
func (fc FunctionContainer) TaskAgDailySalary() error {
	TaskAgDailySalary()
	return nil
}

//同步server状态
func (fc FunctionContainer) SyncServerState() (err error) {
	common2.ServerOtcState = common.IsServerStop(common.ServerOtc)
	common2.ServerGameRunning = common.IsServerStop(common.ServerGame) == common.ServerStateRunning
	common2.ServerOtcTradeRunning = common.IsServerStop(common.ServerOtcTrade) == common.ServerStateRunning
	return
}

// send ping message
func (f FunctionContainer) Ping() error {
	return common.Ping()
}
