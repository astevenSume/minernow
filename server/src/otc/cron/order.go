package cron

import (
	"common"
	"fmt"
	"otc/trade"
	"time"
	"utils/otc/dao"
)

const (
	Limit = 1000
)

func orderCheckLock(oid uint64, duration time.Duration) bool {
	key := fmt.Sprintf("order_check_lock_%d", oid)
	resp := common.RedisManger.SetNX(key, 1, duration)
	lockSuccess, err := resp.Result()
	if err != nil || !lockSuccess {
		return false
	}
	return true
}

//检查未付款超时订单 plus
func orderCreatedCheck() {
	// 协程个数计数器加1
	common.RedisManger.Incr("OrderCheckTimeOut_GoFuncNum")
	start := 0
	orderInfoList, err := dao.OrdersDaoEntity.GetOrderTimeOut(start, Limit)
	if err != nil {
		common.LogFuncError("OrderCheckTimeOut GetOrderTimeOut ERR:%v", err)
		common.RedisManger.Decr("OrderCheckTimeOut_GoFuncNum")
		return
	}
	count := len(orderInfoList)
	if count <= 0 {
		common.RedisManger.Decr("OrderCheckTimeOut_GoFuncNum")
		return
	}
	min15 := int64(trade.PayExpire()) / 1000000 //转成对应毫秒数
	now := common.NowInt64MS()

	deal := 0
	for _, v := range orderInfoList {
		if orderCheckLock(v.Id, time.Second*15) {
			leave := v.Ctime + min15 - now
			if leave <= 0 {
				//已经超时
				if v.Side == dao.SideBuy {
					trade.BuyLogicEntity.TimeOutOrder(v.Id)
				} else if v.Side == dao.SideSell {
					trade.SellLogicEntity.TimeOutOrder(v.Id)
				}
				deal++
				continue
			}
		}
	}
	// 协程个数计数器减1
	common.RedisManger.Decr("OrderCheckTimeOut_GoFuncNum")
	return
}

func OrderCheckTimeOut() {
	goNum, err := common.RedisManger.Get("OrderCheckTimeOut_GoFuncNum").Int()
	if err != nil && err.Error() != "redis: nil" {
		common.LogFuncError("OrderCheckTimeOut RedisManger ERR:%v", err)
	}
	if goNum <= 2 {
		go orderCreatedCheck()
	}
	return
}
