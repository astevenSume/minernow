package eosplus

import (
	"common"
	. "otc_error"
	"strconv"
	"utils/eusd/dao"
	"utils/eusd/models"
)

type Wealth struct {
}

func WealthRpc() *Wealth {
	return &Wealth{}
}

func (a *Wealth) Create(uid uint64, bindAccount bool) (w *models.EosWealth, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	if uid < 1 {
		errCode = ERROR_CODE_UID_PARAMS_ERROR
		return
	}
	acc := &models.EosAccount{}
	if bindAccount {
		acc, errCode = AccountRpc().Bind(uid)
		if errCode != ERROR_CODE_SUCCESS {
			return
		}
	}
	_, err := dao.WealthDaoEntity.Add(uid, acc.Account)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}
	w, err = dao.WealthDaoEntity.Info(uid)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}
	return
}

func (a *Wealth) InfoMap(uid uint64) (user map[string]interface{}, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	data, _ := dao.WealthDaoEntity.Info(uid)

	user = map[string]interface{}{
		"uid":       strconv.FormatUint(data.Uid, 10),
		"address":   data.Account,
		"state":     dao.WealthStatusToString[int(data.Status)],
		"symbol":    EosConfig.Symbol,
		"precision": EosPrecision,
		"balance": map[string]interface{}{
			"available": data.Available,
			"frozen":    data.Game + data.Trade + data.Transfer,
			"game":      data.Game,
			"trade":     data.Trade,
			"transfer":  data.Transfer,
		},
	}

	return
}

func (a *Wealth) Info(uid uint64) (user *models.EosWealth, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	user, _ = dao.WealthDaoEntity.Info(uid)

	return
}

func (a *Wealth) BindAccount(w *models.EosWealth) (account string, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	if w.Uid < 1 {
		errCode = ERROR_CODE_UID_PARAMS_ERROR
		return
	}
	acc, errCode := AccountRpc().Bind(w.Uid)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}

	w.Account = acc.Account

	err := dao.WealthDaoEntity.Update(w, models.COLUMN_EosWealth_Account)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}
	account = acc.Account
	return
}

// 账号锁定
func (a *Wealth) Lock(uid uint64) (errCode ERROR_CODE) {
	ok, _ := dao.WealthDaoEntity.Lock(uid)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}
	errCode = ERROR_CODE_SUCCESS
	return
}

func (a *Wealth) Unlock(uid uint64) (errCode ERROR_CODE) {
	ok, _ := dao.WealthDaoEntity.UnLock(uid)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}
	errCode = ERROR_CODE_SUCCESS
	return
}

//资产变更记录
func (a *Wealth) Records(uid uint64, types []interface{}, page, limit int64) (list []*models.EosWealthLog, meta *common.Meta) {
	count := dao.WealthLogDaoEntity.Count(uid, types)
	meta = common.MakeMeta(count, page, limit)
	list = []*models.EosWealthLog{}
	if meta.Offset < meta.Total {
		list, _ = dao.WealthLogDaoEntity.Fetch(uid, types, meta.Limit, meta.Offset)
	}

	return
}

// 资产变更详情
func (a *Wealth) RecordInfo(uid, id uint64) (transfer *models.EosTransaction) {
	record := dao.WealthLogDaoEntity.Info(id)
	transfer = &models.EosTransaction{}
	if record.Id == 0 || record.Uid != uid {
		return
	}
	if record.Txid <= 0 {
		return
	}
	transfer = dao.TransactionDaoEntity.Info(record.Txid)

	return
}

//抵押USDT    1.2345 USDT = 1.2345 EUSDT
func (a *Wealth) DelegateUsdt(uid uint64, quant float64) (errCode ERROR_CODE) {
	if quant <= 0 {
		errCode = ERROR_CODE_PARAMS_ERROR
		return
	}
	user, errCode := a.Info(uid)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}
	//token账号发币
	quantity := QuantityFloat64ToInt64(quant)
	ok, _ := dao.WealthDaoEntity.TransferInto(uid, quantity)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}

	//token账号转出
	sysUser, err := TokenAccount()
	if err != nil {
		common.LogFuncError("SYS_TOKEN_USER ERR")
		errCode = ERROR_CODE_UNKNOWN
		return
	}

	logId, _ := dao.WealthLogDaoEntity.Add(uid, dao.WealthLogTypeUsdtDelegate, quantity, 0)
	if logId > 0 {
		go common.SafeRun(func() {
			TransferLogic([]int64{logId}, sysUser, user, quantity)
		})()
	}

	errCode = ERROR_CODE_SUCCESS
	return
}

//赎回USDT
func (a *Wealth) UnDelegateUsdt(uid uint64, quant float64) (errCode ERROR_CODE) {
	quantity := QuantityFloat64ToInt64(quant)
	if quantity <= 0 {
		errCode = ERROR_CODE_PARAMS_ERROR
		return
	}
	user, errCode := a.Info(uid)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}
	if user.Status != dao.WealthStatusWorking {
		errCode = ERROR_CODE_WEALTH_LOCK
		return
	}
	if user.Available < quantity {
		errCode = ERROR_CODE_OTC_TRANSFER_AVAILABLE_LACK
		return
	}
	//token账号回收
	sysUser, err := TokenAccount()
	if err != nil {
		common.LogFuncError("SYS_TOKEN_USER ERR")
		errCode = ERROR_CODE_UNKNOWN
		return
	}

	// redis-lock
	l, err := wealthRedisJamLock(uid)
	if err != nil {
		errCode = ERROR_CODE_REDIS_LOCK_ERR
		common.LogFuncError("%v", err)
		return
	}
	defer func() {
		_ = wealthRedisUnJamLock(l)
	}()

	ok, _ := dao.WealthDaoEntity.TransferOutDirect(uid, quantity)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}

	//解锁
	_ = wealthRedisUnJamLock(l)
	l2, err := wealthRedisJamLock(sysUser.Uid)
	if err != nil {
		errCode = ERROR_CODE_REDIS_LOCK_ERR
		common.LogFuncError("%v", err)
		return
	}
	defer func() {
		_ = wealthRedisUnJamLock(l2)
	}()

	//系统账号入账
	ok, _ = dao.WealthDaoEntity.TransferInto(sysUser.Uid, quantity)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}

	logIds, _ := dao.WealthLogDaoEntity.AddBoth(uid, sysUser.Uid, dao.WealthLogTypeUsdtUnDelegate, dao.WealthLogTypeUsdtUnDelegate, quantity, 0)
	if len(logIds) > 1 {
		go common.SafeRun(func() {
			TransferLogic(logIds, sysUser, user, quantity)
		})()
	}

	errCode = ERROR_CODE_SUCCESS
	return
}

//分润
func (a *Wealth) Commission(uid uint64, quantity int64) (errCode ERROR_CODE) {
	if quantity <= 0 {
		errCode = ERROR_CODE_PARAMS_ERROR
		return
	}
	user, errCode := a.Info(uid)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}

	//quantity := QuantityFloat64ToInt64(quant)
	ok, _ := dao.WealthDaoEntity.TransferInto(uid, quantity)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}

	// 分红账号转出
	sysUser, err := CommissionAccount()
	if err != nil {
		common.LogFuncError("SYS_TOKEN_USER ERR")
		errCode = ERROR_CODE_UNKNOWN
		return
	}
	ok, _ = dao.WealthDaoEntity.TransferOutDirect(sysUser.Uid, quantity)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}

	logIds, _ := dao.WealthLogDaoEntity.AddBoth(uid, sysUser.Uid, dao.WealthLogTypeCommission, dao.WealthLogTypeCommission, quantity, 0)
	if len(logIds) > 1 {
		go common.SafeRun(func() {
			TransferLogic(logIds, sysUser, user, quantity)
		})()
	}

	errCode = ERROR_CODE_SUCCESS
	return
}
