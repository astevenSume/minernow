package eosplus

import (
	"common"
	. "otc_error"
	"utils/eusd/dao"
	"utils/eusd/models"
	otcDao "utils/otc/dao"
)

type Otc struct {
}

func OtcRpc() *Otc {
	return &Otc{}
}

//成为承兑商
func (a *Wealth) BecomeExchanger(uid uint64) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	w, err := dao.WealthDaoEntity.Info(uid)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}

	if w.Ctime == 0 {
		errCode = ERROR_CODE_WEALTH_NOT_FOUND
		return
	}

	if w.Status == dao.WealthStatusLock {
		errCode = ERROR_CODE_WEALTH_LOCK
		return
	}

	acc := w.Account
	if w.Account == "" { //没有EOS地址，绑定EOS地址
		acc, errCode = WealthRpc().BindAccount(w)
	}
	_, err = dao.EosOtcDaoEntity.Add(uid, acc)
	if err != nil {
		common.LogFuncError("error:%v", err)
		errCode = ERROR_CODE_DB
		return
	}
	return
}

//OTC 账号信息
func (a *Otc) Info(uid uint64) (user *models.EosOtc, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	user, err := dao.EosOtcDaoEntity.Info(uid)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}
	return
}

func (a *Otc) SellSetting(uid uint64, able bool, dayLimit, lowerLimit int64, payType uint8) ERROR_CODE {
	if !dao.EosOtcDaoEntity.SetSell(uid, able, dayLimit, lowerLimit, payType) {
		return ERROR_CODE_OTC_SET_SELL
	}

	return ERROR_CODE_SUCCESS
}

func (a *Otc) StopSell(uid uint64, state string) ERROR_CODE {
	if !dao.EosOtcDaoEntity.StopSell(uid, state) {
		return ERROR_CODE_OTC_SET_BUY
	}

	return ERROR_CODE_SUCCESS
}

func (a *Otc) SellSettingResetState(uid uint64, today int64) ERROR_CODE {
	if !dao.EosOtcDaoEntity.ResetSellState(uid, today) {
		return ERROR_CODE_OTC_SET_SELL
	}

	return ERROR_CODE_SUCCESS
}

func (a *Otc) UpdateSellState(uid uint64, today int64) ERROR_CODE {
	if !dao.EosOtcDaoEntity.UpdateSellState(uid, today) {
		return ERROR_CODE_OTC_SET_SELL
	}

	return ERROR_CODE_SUCCESS
}

func (a *Otc) BuySetting(uid uint64, able bool, dayLimit, lowerLimit int64) ERROR_CODE {
	if !dao.EosOtcDaoEntity.SetBuy(uid, able, dayLimit, lowerLimit) {
		return ERROR_CODE_OTC_SET_BUY
	}

	return ERROR_CODE_SUCCESS
}

func (a *Otc) StopBuy(uid uint64, state string) ERROR_CODE {
	if !dao.EosOtcDaoEntity.StopBuy(uid, state) {
		return ERROR_CODE_OTC_SET_BUY
	}

	return ERROR_CODE_SUCCESS
}

func (a *Otc) BuySettingResetState(uid uint64, today int64) ERROR_CODE {
	if !dao.EosOtcDaoEntity.ResetBuyState(uid, today) {
		return ERROR_CODE_OTC_SET_BUY
	}

	return ERROR_CODE_SUCCESS
}

func (a *Otc) UpdateBuyState(uid uint64, today int64) ERROR_CODE {
	if !dao.EosOtcDaoEntity.UpdateBuyState(uid, today) {
		return ERROR_CODE_OTC_SET_BUY
	}

	return ERROR_CODE_SUCCESS
}

//承兑商账号锁定
func (a *Otc) Lock(uid uint64) (errCode ERROR_CODE) {
	ok, _ := dao.EosOtcDaoEntity.UpdateStatus(uid, dao.WealthStatusLock)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}

	err := otcDao.OtcBuyDaoEntity.DeleteByUid(uid)
	if err != nil {
		common.LogFuncError("error:%v", err)
	}
	err = otcDao.OtcSellDaoEntity.Delete(uid)
	if err != nil {
		common.LogFuncError("error:%v", err)
	}

	errCode = ERROR_CODE_SUCCESS
	return
}

//承兑商账号解锁
func (a *Otc) Unlock(uid uint64) (errCode ERROR_CODE) {
	ok, _ := dao.EosOtcDaoEntity.UpdateStatus(uid, dao.WealthStatusWorking)
	if !ok {
		errCode = ERROR_CODE_DB
		return
	}
	errCode = ERROR_CODE_SUCCESS
	return
}

//转入otc
func (a *Otc) WealthTransferInto(uid uint64, quantity int64) (otcWealth *models.EosOtc, errCode ERROR_CODE) {
	l, err := wealthRedisJamLock(uid)
	if err != nil {
		errCode = ERROR_CODE_REDIS_LOCK_ERR
		common.LogFuncError("%v", err)
		return
	}
	defer func() {
		_ = wealthRedisUnJamLock(l)
	}()
	if quantity < 1 {
		errCode = ERROR_CODE_QUANTITY_ERR
		return
	}
	w, err := dao.WealthDaoEntity.Info(uid)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}
	if w.Status != dao.WealthStatusWorking {
		errCode = ERROR_CODE_WEALTH_LOCK
		return
	}
	if w.Available < quantity {
		errCode = ERROR_CODE_EUSD_LACK
		return
	}
	otc, _ := dao.EosOtcDaoEntity.Info(uid)
	if otc.Account == "" {
		errCode = ERROR_CODE_OTC_NOT_EXCHANGER
		return
	}
	if otc.Status != dao.WealthStatusWorking {
		errCode = ERROR_CODE_WEALTH_LOCK
		return
	}

	ok := dao.WealthDaoEntity.TransferOtc(uid, quantity)
	if !ok {
		errCode = ERROR_CODE_WEALTH_OTC_TRANSFER
		return
	}
	_, _ = dao.WealthLogDaoEntity.Add(uid, dao.WealthLogTypeToOtc, quantity, 0)

	ok = dao.EosOtcDaoEntity.WealthTransferInto(uid, quantity)
	if !ok {
		errCode = ERROR_CODE_WEALTH_OTC_TRANSFER
		return
	}
	otcWealth = otc
	errCode = ERROR_CODE_SUCCESS
	return
}

//otc转出
func (a *Otc) TransferToWealth(uid uint64, quantity int64) (otcWealth *models.EosOtc, errCode ERROR_CODE) {
	if quantity < 1 {
		errCode = ERROR_CODE_QUANTITY_ERR
		return
	}
	w, err := dao.EosOtcDaoEntity.Info(uid)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}
	if w.Status != dao.WealthStatusWorking {
		errCode = ERROR_CODE_WEALTH_LOCK
		return
	}
	if w.Available < int64(quantity) {
		errCode = ERROR_CODE_EUSD_LACK
		return
	}

	ok := dao.EosOtcDaoEntity.TransferToWealth(uid, quantity)
	if !ok {
		errCode = ERROR_CODE_WEALTH_OTC_TRANSFER
		return
	}
	_, _ = dao.WealthLogDaoEntity.Add(uid, dao.WealthLogTypeFromOtc, quantity, 0)

	ok = dao.WealthDaoEntity.OtcTransferInto(uid, quantity)
	if !ok {
		errCode = ERROR_CODE_WEALTH_OTC_TRANSFER
		return
	}
	otcWealth = w
	errCode = ERROR_CODE_SUCCESS
	return
}

//冻结用于otc交易的Token
func (a *Otc) TransferLock(uid uint64, quantity int64) (errCode ERROR_CODE) {
	if quantity < 1 {
		errCode = ERROR_CODE_QUANTITY_ERR
		return
	}
	w, err := dao.EosOtcDaoEntity.Info(uid)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}
	if w.Status != dao.WealthStatusWorking {
		errCode = ERROR_CODE_WEALTH_LOCK
		return
	}
	if w.Available < int64(quantity) {
		errCode = ERROR_CODE_EUSD_LACK
		return
	}

	ok := dao.EosOtcDaoEntity.TransferLock(uid, quantity)
	if !ok {
		errCode = ERROR_CODE_OTC_TRANSFER_LOCK
		return
	}
	errCode = ERROR_CODE_SUCCESS

	return
}

//解冻用于otc交易的Token
func (a *Otc) TransferUnLock(uid uint64, quantity int64) (errCode ERROR_CODE) {
	if quantity < 1 {
		errCode = ERROR_CODE_QUANTITY_ERR
		return
	}
	w, err := dao.EosOtcDaoEntity.Info(uid)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}
	if w.Status != dao.WealthStatusWorking {
		errCode = ERROR_CODE_WEALTH_LOCK
		return
	}
	if w.Trade < quantity {
		errCode = ERROR_CODE_EUSD_LACK
		return
	}

	ok := dao.EosOtcDaoEntity.TransferUnLock(uid, quantity)
	if !ok {
		errCode = ERROR_CODE_OTC_TRANSFER_UNLOCK
		return
	}
	errCode = ERROR_CODE_SUCCESS

	return
}

//确认转出
func (a *Otc) TransferOut(orderId, from, to uint64, quantity int64) (errCode ERROR_CODE) {
	if quantity < 1 {
		errCode = ERROR_CODE_QUANTITY_ERR
		return
	}
	w, err := dao.EosOtcDaoEntity.Info(from)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}

	if w.Status != dao.WealthStatusWorking {
		errCode = ERROR_CODE_WEALTH_LOCK
		return
	}
	if w.Trade < quantity {
		errCode = ERROR_CODE_EUSD_LACK
		return
	}

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

	ok := dao.EosOtcDaoEntity.TransferOut(from, quantity)
	if !ok {
		errCode = ERROR_CODE_OTC_TRANSFER_OUT
		return
	}
	//转入用户transferInto 账户
	ok, err = dao.WealthDaoEntity.TransferInto(to, quantity)
	if !ok {
		errCode = ERROR_CODE_OTC_TRANSFER_OUT
		return
	}

	logId, _ := dao.WealthLogDaoEntity.Add(toUser.Uid, dao.WealthLogTypeOtcInto, quantity, 0)
	if logId > 0 {
		//链上转账
		go common.SafeRun(func() {
			TransferLogicOrder(orderId, []int64{logId}, fromUser, toUser, quantity)
		})()
	}

	errCode = ERROR_CODE_SUCCESS
	return
}

//用户OTC卖币锁定
func (a *Otc) UserTransferLock(uid uint64, quantity int64) (errCode ERROR_CODE) {
	if quantity < 1 {
		errCode = ERROR_CODE_QUANTITY_ERR
		return
	}
	ok, _ := dao.WealthDaoEntity.TransferLock(uid, quantity)
	if !ok {
		errCode = ERROR_CODE_WEALTH_OTC_TRANSFER
		return
	}

	return ERROR_CODE_SUCCESS
}

// 用户卖币转出
func (a *Otc) UserTransferOut(orderId, from, to uint64, quantity int64) (errCode ERROR_CODE) {
	if quantity < 1 {
		errCode = ERROR_CODE_QUANTITY_ERR
		return
	}
	//用户token取出
	ok, _ := dao.WealthDaoEntity.UserOtcOut(from, quantity)
	if !ok {
		errCode = ERROR_CODE_OTC_TRANSFER_OUT
		return
	}
	_, _ = dao.WealthLogDaoEntity.Add(from, dao.WealthLogTypeOtcOut, quantity, 0)

	//承兑商买币转入
	ok = dao.EosOtcDaoEntity.OtcTransferInto(to, quantity)
	if !ok {
		errCode = ERROR_CODE_OTC_TRANSFER_OUT
		return
	}

	logId, _ := dao.WealthLogDaoEntity.Add(from, dao.WealthLogTypeOtcOut, quantity, 0)
	if logId > 0 {
		//链上转账
		go common.SafeRun(func() {
			TransferLogicOrderByUids(orderId, []int64{logId}, from, to, quantity)
		})()
	}

	return ERROR_CODE_SUCCESS
}

//用户OTC取消卖币锁定
func (a *Otc) UserTransferUnlock(uid uint64, quantity int64) (errCode ERROR_CODE) {
	if quantity < 1 {
		errCode = ERROR_CODE_QUANTITY_ERR
		return
	}
	ok, _ := dao.WealthDaoEntity.TransferUnLock(uid, quantity)
	if !ok {
		errCode = ERROR_CODE_WEALTH_OTC_TRANSFER
		return
	}

	return ERROR_CODE_SUCCESS
}
