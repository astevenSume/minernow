package eosplus

import (
	"common"
	. "otc_error"
	"utils/eusd/dao"
)

type Game struct {
}

// 游戏冻结
func (a *Game) GameLock(uid uint64, limit float64) (num float64, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	if limit <= 0 {
		errCode = ERROR_CODE_PARAMS_ERROR
		return
	}
	user, err := dao.WealthDaoEntity.Info(uid)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}

	if user.Status != dao.WealthStatusWorking {
		errCode = ERROR_CODE_WEALTH_LOCK
		return
	}

	if user.Account == "" {
		return
	}

	// redisLock
	l, err := wealthRedisJamLock(uid)
	if err != nil {
		errCode = ERROR_CODE_REDIS_LOCK_ERR
		common.LogFuncError("%v", err)
		return
	}
	defer func() {
		_ = wealthRedisUnJamLock(l)
	}()

	if user.Available <= 0 && user.Transfer <= 0 {
		return
	}
	//可用token && 待转入token都可以 带入游戏; 先使用可用token
	limitQuant := QuantityFloat64ToInt64(limit)
	use := int64(0)
	useTransfer := int64(0)
	// 先使用可用Token，
	if user.Available >= limitQuant {
		use = limitQuant

	} else if user.Available > 10000 {
		// >1EUSD 去掉小数位，不带入
		use = (user.Available / 10000) * 10000
	}

	//可用token不足，使用转入token
	if limitQuant-use > 0 && user.Transfer > 0 {
		useTransfer = limitQuant - use
		if user.Transfer < useTransfer {
			useTransfer = user.Transfer
		}
		useTransfer = (user.Transfer / 10000) * 10000
	}

	ok := dao.WealthDaoEntity.GameLock(uid, use, useTransfer)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}
	num = QuantityInt64ToFloat64(use + useTransfer)

	return
}

// 结算
func (a *Game) Settlement(uid uint64, remain float64) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	user, err := dao.WealthDaoEntity.Info(uid)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}
	if user.Uid == 0 {
		errCode = ERROR_CODE_WEALTH_NOT_FOUND
		return
	}

	quant := QuantityFloat64ToInt64(remain)

	if user.Game == 0 && user.TransferGame == 0 { //无抵押资金
		if quant > 0 {
			common.LogFuncError("Game Settlement Err:Lock 0 and remain %v", quant)
			errCode = ERROR_CODE_GAME_SETTLEMENT_ERR
			return
		}
		return
	}

	// redisLock
	l, err := wealthRedisJamLock(uid)
	if err != nil {
		errCode = ERROR_CODE_REDIS_LOCK_ERR
		common.LogFuncError("%v", err)
		return
	}
	defer func() {
		_ = wealthRedisUnJamLock(l)
	}()

	if quant == user.Game+user.TransferGame { // 无变化
		ok := dao.WealthDaoEntity.GameUnlock(uid, user.Game, user.Game, user.TransferGame, user.TransferGame)
		if !ok {
			errCode = ERROR_CODE_DB
			return
		}
		return
	}

	//余额大于抵押资金 赢钱
	if quant > user.Game+user.TransferGame {
		// 还给用户
		ok := dao.WealthDaoEntity.GameUnlock(uid, user.Game, user.Game, user.TransferGame, user.TransferGame)
		if !ok {
			errCode = ERROR_CODE_DB
			return
		}

		win := quant - user.Game - user.TransferGame
		// 用户赢钱
		ok = dao.WealthDaoEntity.GameUserInto(uid, win)
		if !ok {
			errCode = ERROR_CODE_DB
			return
		}
		//游戏平台账号
		sysUser, err := GamePlatformAccount()
		if err != nil {
			common.LogFuncError("SYS_TOKEN_USER ERR")
			errCode = ERROR_CODE_UNKNOWN
			return
		}
		ok, _ = dao.WealthDaoEntity.TransferOutDirect(sysUser.Uid, win)
		if !ok {
			errCode = ERROR_CODE_DB
			return
		}

		dao.WealthDaoEntity.GamePlatformOut(sysUser.Uid, quant)
		logIds, _ := dao.WealthLogDaoEntity.AddBoth(uid, sysUser.Uid, dao.WealthLogTypeGameWin, dao.WealthLogTypeGameLost, win, 0)
		if len(logIds) == 2 {
			//链上转账
			go common.SafeRun(func() {
				TransferLogic(logIds, sysUser, user, win)
			})()
		}

		return
	}

	// 输钱
	lost := user.Game + user.TransferGame - quant
	if quant > 0 {
		transferLeave := user.TransferGame
		if user.TransferGame > quant {
			transferLeave = quant
		}

		leave := quant - transferLeave
		if leave < 0 {
			leave = 0
		}

		// 还给用户
		ok := dao.WealthDaoEntity.GameUnlock(uid, user.Game, leave, user.TransferGame, transferLeave)
		if !ok {
			errCode = ERROR_CODE_DB
			return
		}
	}

	//  平台账号入账
	sysUser, err := GamePlatformAccount()
	if err != nil {
		common.LogFuncError("SYS_TOKEN_USER ERR")
		errCode = ERROR_CODE_UNKNOWN
		return
	}
	ok, _ := dao.WealthDaoEntity.TransferInto(sysUser.Uid, lost)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}
	if ok := dao.WealthDaoEntity.GamePlatformInto(sysUser.Uid, lost); !ok {
		errCode = ERROR_CODE_DB
		return
	}
	logIds, _ := dao.WealthLogDaoEntity.AddBoth(uid, sysUser.Uid, dao.WealthLogTypeGameLost, dao.WealthLogTypeGameWin, lost, 0)
	if len(logIds) == 2 {
		//链上转账
		go common.SafeRun(func() {
			TransferLogic(logIds, user, sysUser, lost)
		})()
	}
	return
}

//充值到游戏中 - 先转账到游戏平台中
func (a *Game) Recharge(uid uint64, num float64) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	if num <= 0 {
		errCode = ERROR_CODE_PARAMS_ERROR
		return
	}

	// redisLock
	l, err := wealthRedisJamLock(uid)
	if err != nil {
		errCode = ERROR_CODE_REDIS_LOCK_ERR
		common.LogFuncError("%v", err)
		return
	}
	defer func() {
		_ = wealthRedisUnJamLock(l)
	}()

	user, err := dao.WealthDaoEntity.Info(uid)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}

	if user.Status != dao.WealthStatusWorking {
		errCode = ERROR_CODE_WEALTH_LOCK
		return
	}

	quantity := QuantityFloat64ToInt64(num)

	if user.Available+user.Transfer < quantity {
		errCode = ERROR_CODE_EUSD_LACK
		return
	}

	availableUse := quantity
	transferUse := int64(0)
	if user.Available < quantity {
		availableUse = user.Available
		transferUse = quantity - user.Available
	}
	ok := dao.WealthDaoEntity.GameRecharge(uid, availableUse, transferUse)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}
	//  平台账号入账
	sysUser, err := GamePlatformAccount()
	if err != nil {
		common.LogFuncError("PLATFORM_USER ERR")
		errCode = ERROR_CODE_UNKNOWN
		return
	}
	if ok := dao.WealthDaoEntity.GamePlatformInto(sysUser.Uid, quantity); !ok {
		errCode = ERROR_CODE_DB
		return
	}
	logIds, _ := dao.WealthLogDaoEntity.AddBoth(uid, sysUser.Uid, dao.WealthLogTypeGameLost, dao.WealthLogTypeGameLost, quantity, 0)
	if len(logIds) == 2 {
		//链上转账
		go common.SafeRun(func() {
			//todo 延时转账处理
			TransferLogic(logIds, user, sysUser, quantity)
		})()
	}

	return
}

//从游戏提现 - 从平台账号中提出EUSD
func (a *Game) Withdrawal(uid uint64, num float64) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	if num <= 0 {
		errCode = ERROR_CODE_PARAMS_ERROR
		return
	}
	quantity := QuantityFloat64ToInt64(num)

	// redisLock
	l, err := wealthRedisJamLock(uid)
	if err != nil {
		errCode = ERROR_CODE_REDIS_LOCK_ERR
		common.LogFuncError("%v", err)
		return
	}
	defer func() {
		_ = wealthRedisUnJamLock(l)
	}()

	user, err := dao.WealthDaoEntity.Info(uid)
	if err != nil {
		errCode = ERROR_CODE_NO_USER
		return
	}

	sysUser, err := GamePlatformAccount()
	if err != nil {
		common.LogFuncError("PLATFORM_USER ERR")
		errCode = ERROR_CODE_UNKNOWN
		return
	}
	//平台转出
	if ok := dao.WealthDaoEntity.GamePlatformOut(sysUser.Uid, quantity); !ok {
		errCode = ERROR_CODE_DB
		return
	}
	//用户收款
	if ok := dao.WealthDaoEntity.GameUserInto(uid, quantity); !ok {
		errCode = ERROR_CODE_DB
		return
	}
	logIds, _ := dao.WealthLogDaoEntity.AddBoth(uid, sysUser.Uid, dao.WealthLogTypeGameWin, dao.WealthLogTypeGameWin, quantity, 0)
	if len(logIds) == 2 {
		//链上转账
		go common.SafeRun(func() {
			TransferLogic(logIds, sysUser, user, quantity)
		})()
	}

	return
}
