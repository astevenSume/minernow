package controllers

import (
	"common"
	"time"
	"usdt/prices"
	usdtdao "utils/usdt/dao"

	"github.com/astaxie/beego"
	lock "github.com/bsm/redis-lock"
)

func Init() (err error) {
	// init dao
	err = usdtdao.Init(nil)
	if err != nil {
		return
	}

	//process pk relation
	go common.SafeRun(func() {
		var interval int64
		interval, err = beego.AppConfig.Int64("usdt::PkRelationInterval")
		if err != nil {
			common.LogFuncError("%v", err)
			interval = 5 //default 5 seconds.
		}

		for {
			ProcessPkUnrelated()
			time.Sleep(time.Second * time.Duration(interval))
		}
	})()

	// go common.SafeRun(func() {
	// 	var interval int64
	// 	interval, err = beego.AppConfig.Int64("usdt::ApprovedTransferInterval")
	// 	if err != nil {
	// 		common.LogFuncError("%v", err)
	// 		interval = 30 //default 30 seconds.
	// 	}

	// 	for {
	// 		ProcessApprovedTransfer()
	// 		time.Sleep(time.Second * time.Duration(interval))
	// 	}
	// })()

	go common.SafeRun(func() {
		var interval int64
		interval, err = beego.AppConfig.Int64("usdt::SyncRecommendedFeesInterval")
		if err != nil {
			common.LogFuncError("%v", err)
			interval = 30 //default 30 seconds.
		}

		for {
			locker("usdt_sync_recommended_fees", SyncRecommendedFees, time.Second*time.Duration(interval))
		}
	})()

	go common.SafeRun(func() {
		var interval int64
		interval, err = beego.AppConfig.Int64("usdt::ProcessSweepInterval")
		if err != nil {
			common.LogFuncError("%v", err)
			interval = 30 //default 30 seconds.
		}

		var lockSecs int64
		lockSecs, err = beego.AppConfig.Int64("usdt::ProcessSweepLockSecs")
		if err != nil {
			lockSecs = 10800 //default 3 hours
		}

		for {

			l, err := common.RedisLock2("usdt_process_cash_sweep", lock.Options{
				LockTimeout: time.Second * time.Duration(lockSecs),
				RetryCount:  common.RetryCount,
				RetryDelay:  common.RetryDelay,
			})
			if err != nil { //doesn't get redis lock, skip
				if err != lock.ErrLockNotObtained {
					common.LogFuncError("%v", err)
				}
			} else { //got redis lock, do the job
				ProcessSweep()
				common.RedisUnlock(l)
			}

			time.Sleep(time.Second * time.Duration(interval))
		}
	})()

	// 初始化热钱包组
	go common.SafeRun(func() {
		var interval int64
		interval, err = beego.AppConfig.Int64("wallet::ProcessGeneratorInterval")
		if err != nil {
			common.LogFuncError("%v", err)
			interval = 30 //default 30 seconds.
		}

		var lockSecs int64
		lockSecs, err = beego.AppConfig.Int64("wallet::ProcessGeneratorInterval")
		if err != nil {
			lockSecs = 10800 //default 3 hours
		}

		for {
			l, err := common.RedisLock2("usdt_hot_wallet_group_generate", lock.Options{
				LockTimeout: time.Second * time.Duration(lockSecs),
				RetryCount:  common.RetryCount,
				RetryDelay:  common.RetryDelay,
			})
			if err != nil { //doesn't get redis lock, skip
				if err != lock.ErrLockNotObtained {
					common.LogFuncError("%v", err)
				}
			} else { //got redis lock, do the job
				HotWalletGroupGenerate()
				common.RedisUnlock(l)
			}

			time.Sleep(time.Second * time.Duration(interval))
		}
	})()

	// 定时获取冷热钱包组余额
	go common.SafeRun(func() {
		var interval int64
		interval, err = beego.AppConfig.Int64("wallet::ProcessSyncBalanceInterval")
		if err != nil {
			common.LogFuncError("%v", err)
			interval = 30 //default 30 seconds.
		}

		var lockSecs int64
		lockSecs, err = beego.AppConfig.Int64("wallet::ProcessSyncBalanceLockSecs")
		if err != nil {
			lockSecs = 10800 //default 3 hours
		}

		for {
			l, err := common.RedisLock2("usdt_wallet_group_sync_balance", lock.Options{
				LockTimeout: time.Second * time.Duration(lockSecs),
				RetryCount:  common.RetryCount,
				RetryDelay:  common.RetryDelay,
			})
			if err != nil { //doesn't get redis lock, skip
				if err != lock.ErrLockNotObtained {
					common.LogFuncError("%v", err)
				}
			} else { //got redis lock, do the job
				WalletGroupSyncBalance()
				common.RedisUnlock(l)
			}

			time.Sleep(time.Second * time.Duration(interval))
		}
	})()

	// TODO: 暂时使用死循环,后面改成定时调用
	go common.SafeRun(func() {
		for {

			l, err := common.RedisLock2("usdt_sync_prices", lock.Options{
				LockTimeout: time.Second * time.Duration(60),
				RetryCount:  common.RetryCount,
				RetryDelay:  common.RetryDelay,
			})

			if _, err = prices.SyncUsdtPrice(); err != nil {
				common.LogFuncError("%v", err)
			}
			if _, err := prices.SyncBtcPrice(); err != nil {
				common.LogFuncError("%v", err)
			}

			common.RedisUnlock(l)

			time.Sleep(time.Second * time.Duration(5*60))
		}
	})

	// init rabbitmq
	err = common.RabbitMQInit(&AmqpFuncContainer{})
	if err != nil {
		return
	}

	return
}

func locker(lockKey string, execFunc func(), interval time.Duration) {
	var (
		l   *lock.Locker
		err error
	)

	defer func() {
		time.Sleep(interval)
		if l != nil {
			common.RedisUnlock(l)
		}
	}()
	if l, err = common.RedisLock2(lockKey, lock.Options{LockTimeout: time.Second * time.Duration(60)}); err != nil {
		return
	}

	execFunc()

}
