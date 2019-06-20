package common

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/bsm/redis-lock"
	"github.com/go-redis/redis"
	"time"
)

var RedisManger *redis.Client

var (
	ExpiredTime time.Duration
	RetryCount  int
	RetryDelay  time.Duration
	TokenPrefix string
)

func RedisInit() (err error) {
	// init redis manager
	RedisManger = redis.NewClient(&redis.Options{
		Addr:     beego.AppConfig.String("redis::host") + ":" + beego.AppConfig.String("redis::port"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// init distributed lock configuration
	var (
		expiredSecs            int
		retryCount, retryDelay int
		regionId, serverId     int64
		appName                string
	)

	if expiredSecs, err = beego.AppConfig.Int("redis::LockExpired"); err != nil {
		ExpiredTime = time.Second * 5
	} else {
		ExpiredTime = time.Duration(expiredSecs) * time.Second
	}

	if retryCount, err = beego.AppConfig.Int("redis::LockRetryCount"); err != nil {
		RetryCount = 5
	} else {
		RetryCount = retryCount
	}

	if retryDelay, err = beego.AppConfig.Int("redis::LockRetryDelay"); err != nil {
		RetryDelay = 100 * time.Millisecond
	} else {
		RetryDelay = time.Duration(retryDelay) * time.Millisecond
	}

	appName = beego.AppConfig.String("appname")
	if regionId, err = beego.AppConfig.Int64("RegionId"); err != nil {
		panic("no specific RegionId")
	}

	if serverId, err = beego.AppConfig.Int64("ServerId"); err != nil {
		panic("no specific ServerId")
	}
	TokenPrefix = fmt.Sprintf("%s.%d.%d.", appName, regionId, serverId)

	return
}

//
func RedisLock(key string) (l *lock.Locker, err error) {
	opts := lock.Options{
		LockTimeout: ExpiredTime,
		RetryCount:  RetryCount,
		RetryDelay:  RetryDelay,
		TokenPrefix: TokenPrefix,
	}
	return RedisLock2(key, opts)
}

func RedisLock2(key string, opts lock.Options) (l *lock.Locker, err error) {
	l, err = lock.Obtain(RedisManger, key, &opts)
	if err != nil {
		LogFuncError("obtain lock %s err : %v", key, err)
		return
	} else if l == nil {
		LogFuncError("could not obtain lock %s", key)
		return
	}

	return
}

func RedisUnlock(l *lock.Locker) (err error) {
	err = l.Unlock()
	if err != nil {
		LogFuncError("unlock %v failed : %v", l, err)
		return
	}

	return
}

//
func RedisSetNX(key string, expireSecond time.Duration) bool {
	resp := RedisManger.SetNX(key, 1, expireSecond*time.Second)
	lockSuccess, err := resp.Result()
	if err != nil || !lockSuccess {
		return false
	}

	return true
}
