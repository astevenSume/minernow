package usdt

import (
	"common"
	"fmt"
	"github.com/bsm/redis-lock"
)

const UsdtAccountLockPre = "usdt_account_"

func UsdtAccountRedisLock(uid uint64) (l *lock.Locker, err error) {
	l, err = common.RedisLock(fmt.Sprintf("%s%d", UsdtAccountLockPre, uid))
	return
}

func UsdtAccountRedisUnLock(l *lock.Locker) (err error) {
	err = common.RedisUnlock(l)
	return
}
