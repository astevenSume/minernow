package eosplus

import (
	"common"
	. "otc_error"
	"utils/eusd/dao"
	"utils/eusd/models"
)

type Transaction struct {
}

func (a *Transaction) TransferByUids(from, to uint64, quantity int64, memo string) (errCode ERROR_CODE) {
	users, err := dao.WealthDaoEntity.FetchByIds(from, to)

	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}

	fromUser, toUser := &models.EosWealth{}, &models.EosWealth{}
	for _, u := range users {
		if u.Uid == from {
			fromUser = u
		}
		if u.Uid == to {
			toUser = u
		}
	}
	errCode = a.transferring(fromUser, toUser, quantity, memo)
	if errCode != ERROR_CODE_SUCCESS {
		a.transferUnLock(fromUser.Uid, quantity)
		return
	}

	return
}

//转账资源检查
func (a *Transaction) transferring(fromUser, toUser *models.EosWealth, quantity int64, memo string) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	//判断账号是否被冻结
	if fromUser.Status == dao.WealthStatusLock {
		errCode = ERROR_CODE_WEALTH_LOCK
		return
	}
	//资金是否充足
	if fromUser.Available < quantity {
		errCode = ERROR_CODE_TRANSFER_LACK_AVAILABLE
		return
	}

	// redis-lock
	l, err := wealthRedisJamLock(fromUser.Uid)
	if err != nil {
		errCode = ERROR_CODE_REDIS_LOCK_ERR
		common.LogFuncError("%v", err)
		return
	}
	defer func() {
		_ = wealthRedisUnJamLock(l)
	}()

	//资金直接转出
	ok, _ := dao.WealthDaoEntity.TransferOutDirect(fromUser.Uid, quantity)
	if !ok {
		errCode = ERROR_CODE_TRANSFER_DB_OUT
	}

	wealthRedisUnJamLock(l)

	l, err = wealthRedisJamLock(toUser.Uid)
	if err != nil {
		errCode = ERROR_CODE_REDIS_LOCK_ERR
		common.LogFuncError("%v", err)
		return
	}
	defer func() {
		_ = wealthRedisUnJamLock(l)
	}()

	// to 资金转入
	ok, _ = dao.WealthDaoEntity.TransferInto(toUser.Uid, quantity)
	if !ok {
		errCode = ERROR_CODE_TRANSFER_DB_INTO
	}

	logIds, _ := dao.WealthLogDaoEntity.AddBoth(fromUser.Uid, toUser.Uid, dao.WealthLogTypeTransferOut, dao.WealthLogTypeTransferInto, quantity, 0)
	if len(logIds) == 2 {
		//链上转账
		go common.SafeRun(func() {
			TransferLogicMemo(logIds, fromUser, toUser, quantity, memo)
		})()
	}

	errCode = ERROR_CODE_SUCCESS
	return
}

//转账资金锁定
func (a *Transaction) transferLock(from uint64, quantity int64) bool {
	ok, err := dao.WealthDaoEntity.TransferLock(from, quantity)
	if err != nil {
		common.LogFuncError("Transfer Lock DB:%v", err)
		return false
	}
	if !ok {
		return false
	}
	return true
}

//转账资金解锁，交易失败解锁
func (a *Transaction) transferUnLock(from uint64, quantity int64) bool {
	ok, err := dao.WealthDaoEntity.TransferUnLock(from, quantity)
	if err != nil {
		common.LogFuncError("Transfer Unlock DB ERR || Uid:%d, Qunant:%d,err:%v", from, quantity, err)
		return false
	}
	if !ok {
		common.LogFuncError("Transfer Unlock Trade || Uid:%d, Qunant:%d", from, quantity)
		return false
	}
	return true
}
