package eosplus

import (
	"common"
	"fmt"
	"github.com/bsm/redis-lock"
	"time"
)

const EosWealthLockPre = "eos_wealth_"
const EosOtcLockPre = "eos_otc_"

//堵塞锁
func wealthRedisJamLock(uid uint64) (l *lock.Locker, err error) {
	l, err = common.RedisLock(fmt.Sprintf("%s%d", EosWealthLockPre, uid))
	return
}

func wealthRedisUnJamLock(l *lock.Locker) (err error) {
	if l == nil {
		return
	}
	err = common.RedisUnlock(l)
	l = nil
	return
}

//非堵塞锁
func wealthRedisLock(uid uint64) bool {
	resp := common.RedisManger.SetNX(fmt.Sprintf("%s%d", EosWealthLockPre, uid), 1, 30*time.Second)
	lockSuccess, err := resp.Result()
	if err != nil || !lockSuccess {
		return false
	}
	return true
}

func wealthRedisUnLock(uid uint64) bool {
	common.RedisManger.Del(fmt.Sprintf("%s%d", EosWealthLockPre, uid))
	return true
}
