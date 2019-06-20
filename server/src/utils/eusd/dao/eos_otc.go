package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/eusd/models"
)

type EosOtcDao struct {
	common.BaseDao
}

func NewEosOtcDao(db string) *EosOtcDao {
	return &EosOtcDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *EosOtcDao) Add(uid uint64, acc string) (id int64, err error) {
	data := &models.EosOtc{
		Uid:     uid,
		Account: acc,
		Status:  WealthStatusWorking,
		Ctime:   common.NowInt64MS(),
	}
	id, err = d.Orm.Insert(data)
	if err != nil {
		return
	}

	return
}

func (d *EosOtcDao) Info(uid uint64) (user *models.EosOtc, err error) {
	user = &models.EosOtc{
		Uid: uid,
	}
	err = d.Orm.Read(user)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			user.Uid = 0
			return
		}
		common.LogFuncError("DB:%v", err)
		return
	}

	return
}

// 账号OTC锁定/解锁
func (d *EosOtcDao) UpdateStatus(uid uint64, status int8) (ok bool, err error) {
	eosOtc := &models.EosOtc{
		Uid:      uid,
		Status:   status,
		BuyAble:  false,
		SellAble: false,
	}
	var num int64
	if status == WealthStatusLock {
		num, err = d.Orm.Update(eosOtc, models.COLUMN_EosOtc_Status, models.COLUMN_EosOtc_BuyAble, models.COLUMN_EosOtc_SellAble)
	} else {
		num, err = d.Orm.Update(eosOtc, models.COLUMN_EosOtc_Status)
	}

	if err != nil {
		common.LogFuncError("Wealth (Un)Lock DBERR: %v, uid:%d, s:%d", err, uid, status)
		return
	}
	if num == 0 {
		return
	}

	ok = true
	return
}

func (d *EosOtcDao) ResetAble(uid uint64) (err error) {
	eosOtc := &models.EosOtc{
		Uid:      uid,
		BuyAble:  false,
		SellAble: false,
	}
	_, err = d.Orm.Update(eosOtc, models.COLUMN_EosOtc_BuyAble, models.COLUMN_EosOtc_SellAble)
	if err != nil {
		common.LogFuncError("Wealth (Un)Lock DBERR: %v, uid:%d", err, uid)
		return
	}

	return
}

//从普通账户转入OTC
func (d *EosOtcDao) WealthTransferInto(uid uint64, quantity int64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=%s+? Where %s=? And %s=?",
		models.TABLE_EosOtc,
		models.COLUMN_EosOtc_Available, models.COLUMN_EosOtc_Available,
		models.COLUMN_EosOtc_Uid,
		models.COLUMN_EosOtc_Status)

	res, err := d.Orm.Raw(sql, quantity, uid, WealthStatusWorking).Exec()
	if err != nil {
		common.LogFuncError("Wealth To Otc DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("Wealth To Otc DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	ok = true
	return
}

//OTC 转到 普通账户
func (d *EosOtcDao) TransferToWealth(uid uint64, quantity int64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=%s-? Where %s=? And %s>=? And %s=?",
		models.TABLE_EosOtc,
		models.COLUMN_EosOtc_Available, models.COLUMN_EosOtc_Available,
		models.COLUMN_EosOtc_Uid,
		models.COLUMN_EosOtc_Available,
		models.COLUMN_EosOtc_Status)

	res, err := d.Orm.Raw(sql, quantity, uid, quantity, WealthStatusWorking).Exec()
	if err != nil {
		common.LogFuncError("Otc To Wealth DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("Otc To Wealth DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	ok = true
	return
}

// 交易开始时冻结
func (d *EosOtcDao) TransferLock(uid uint64, quantity int64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=%s-?,%s=%s+? Where %s=? And %s>=? And %s=?",
		models.TABLE_EosOtc,
		models.COLUMN_EosOtc_Available, models.COLUMN_EosOtc_Available,
		models.COLUMN_EosOtc_Trade, models.COLUMN_EosOtc_Trade,
		models.COLUMN_EosOtc_Uid,
		models.COLUMN_EosOtc_Available,
		models.COLUMN_EosOtc_Status)

	res, err := d.Orm.Raw(sql, quantity, quantity, uid, quantity, WealthStatusWorking).Exec()
	if err != nil {
		common.LogFuncError("Otc Lock DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("Otc Lock DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	ok = true
	return
}

// 交易取消时解冻
func (d *EosOtcDao) TransferUnLock(uid uint64, quantity int64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=%s-?,%s=%s+? Where %s=? And %s>=? And %s=?",
		models.TABLE_EosOtc,
		models.COLUMN_EosOtc_Trade, models.COLUMN_EosOtc_Trade,
		models.COLUMN_EosOtc_Available, models.COLUMN_EosOtc_Available,
		models.COLUMN_EosOtc_Uid,
		models.COLUMN_EosOtc_Trade,
		models.COLUMN_EosOtc_Status)

	res, err := d.Orm.Raw(sql, quantity, quantity, uid, quantity, WealthStatusWorking).Exec()
	if err != nil {
		common.LogFuncError("Otc Unlock DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return
	}
	ok = true
	return
}

// 交易成功时转出
func (d *EosOtcDao) TransferOut(uid uint64, quantity int64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=%s-? Where %s=? And %s>=? And %s=?",
		models.TABLE_EosOtc,
		models.COLUMN_EosOtc_Trade, models.COLUMN_EosOtc_Trade,
		models.COLUMN_EosOtc_Uid,
		models.COLUMN_EosOtc_Trade,
		models.COLUMN_EosOtc_Status)

	res, err := d.Orm.Raw(sql, quantity, uid, quantity, WealthStatusWorking).Exec()
	if err != nil {
		common.LogFuncError("Otc TransferOut DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return
	}
	ok = true
	return
}

//设置用户卖币（承兑商买币）
func (d *EosOtcDao) SetSell(uid uint64, able bool, dayLimit int64, lowerLimit int64, payType uint8) (ok bool) {
	data := models.EosOtc{
		Uid:               uid,
		SellAble:          able,
		SellRmbDay:        dayLimit,
		SellRmbLowerLimit: lowerLimit,
		SellPayType:       payType,
		SellState:         "",
	}
	_, err := d.Orm.Update(&data, models.COLUMN_EosOtc_SellAble, models.COLUMN_EosOtc_SellRmbDay,
		models.COLUMN_EosOtc_SellRmbLowerLimit, models.COLUMN_EosOtc_SellPayType, models.COLUMN_EosOtc_SellState)
	if err != nil {
		common.LogFuncError("Set Sell DBERR: u:%d, s:%d, day:%d, low:%d, e:%v", uid, able, dayLimit, lowerLimit, err)
		return
	}

	return true
}

//关闭买币（承兑商买币）
func (d *EosOtcDao) StopSell(uid uint64, state string) (ok bool) {
	data := models.EosOtc{
		Uid:       uid,
		SellAble:  false,
		SellUTime: common.NowInt64MS(),
		SellState: state,
	}
	_, err := d.Orm.Update(&data, models.COLUMN_EosOtc_SellAble, models.COLUMN_EosOtc_SellUTime,
		models.COLUMN_EosOtc_SellState)
	if err != nil {
		common.LogFuncError("Stop Sell DBERR: u:%d, e:%v", uid, err)
		return
	}

	return true
}

// sell-setting-统计 重置
func (d *EosOtcDao) ResetSellState(uid uint64, today int64) (ok bool) {
	data := models.EosOtc{
		Uid:          uid,
		SellUTime:    today,
		SellRmbToday: 0,
	}
	_, err := d.Orm.Update(&data, models.COLUMN_EosOtc_SellUTime, models.COLUMN_EosOtc_SellRmbToday)
	if err != nil {
		return
	}

	return true
}

func (d *EosOtcDao) UpdateSellState(uid uint64, rmb int64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=%s+? Where %s=?",
		models.TABLE_EosOtc,
		models.COLUMN_EosOtc_SellRmbToday, models.COLUMN_EosOtc_SellRmbToday,
		models.COLUMN_EosOtc_Uid)

	res, err := d.Orm.Raw(sql, rmb, uid).Exec()
	if err != nil {
		common.LogFuncError("User To Otc DBERR: u:%d, q:%d, e:%v", uid, rmb, err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("User To Otc DBERR: u:%d, q:%d, e:%v", uid, rmb, err)
		return
	}
	ok = true
	return true
}

//设置买币（承兑商卖币）
func (d *EosOtcDao) SetBuy(uid uint64, able bool, dayLimit int64, lowerLimit int64) (ok bool) {
	data := models.EosOtc{
		Uid:              uid,
		BuyAble:          able,
		BuyRmbDay:        dayLimit,
		BuyRmbLowerLimit: lowerLimit,
		BuyState:         "",
	}
	_, err := d.Orm.Update(&data, models.COLUMN_EosOtc_BuyAble, models.COLUMN_EosOtc_BuyRmbDay,
		models.COLUMN_EosOtc_BuyRmbLowerLimit, models.COLUMN_EosOtc_BuyState)
	if err != nil {
		common.LogFuncError("Set Sell DBERR: u:%d, s:%d, day:%d, low:%d, e:%v", uid, able, dayLimit, lowerLimit, err)
		return
	}

	return true
}

//关闭买币（承兑商卖币）
func (d *EosOtcDao) StopBuy(uid uint64, state string) (ok bool) {
	data := models.EosOtc{
		Uid:      uid,
		BuyAble:  false,
		BuyUTime: common.NowInt64MS(),
		BuyState: state,
	}
	_, err := d.Orm.Update(&data, models.COLUMN_EosOtc_BuyAble, models.COLUMN_EosOtc_BuyUTime,
		models.COLUMN_EosOtc_BuyState)
	if err != nil {
		common.LogFuncError("Stop Sell DBERR: u:%d, e:%v", uid, err)
		return
	}

	return true
}

// buy-setting-统计 重置
func (d *EosOtcDao) ResetBuyState(uid uint64, today int64) (ok bool) {
	data := models.EosOtc{
		Uid:         uid,
		BuyUTime:    today,
		BuyRmbToday: 0,
	}
	_, err := d.Orm.Update(&data, models.COLUMN_EosOtc_BuyUTime, models.COLUMN_EosOtc_BuyRmbToday)
	if err != nil {
		return
	}

	return true
}

func (d *EosOtcDao) UpdateBuyState(uid uint64, rmb int64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=%s+? Where %s=?",
		models.TABLE_EosOtc,
		models.COLUMN_EosOtc_BuyRmbToday, models.COLUMN_EosOtc_BuyRmbToday,
		models.COLUMN_EosOtc_Uid)

	res, err := d.Orm.Raw(sql, rmb, uid).Exec()
	if err != nil {
		common.LogFuncError("User To Otc DBERR: u:%d, q:%d, e:%v", uid, rmb, err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("User To Otc DBERR: u:%d, q:%d, e:%v", uid, rmb, err)
		return
	}
	ok = true
	return true
}

//从OTC买币转入
func (d *EosOtcDao) OtcTransferInto(uid uint64, quantity int64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=%s+? Where %s=? And %s=?",
		models.TABLE_EosOtc,
		models.COLUMN_EosOtc_Transfer, models.COLUMN_EosOtc_Transfer,
		models.COLUMN_EosOtc_Uid,
		models.COLUMN_EosOtc_Status)

	res, err := d.Orm.Raw(sql, quantity, uid, WealthStatusWorking).Exec()
	if err != nil {
		common.LogFuncError("User To Otc DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("User To Otc DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	ok = true
	return
}

//解冻 OTC买币转入
func (d *EosOtcDao) OtcTransferUnlock(uid uint64, quantity int64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=%s+?,%s=%s-? Where %s=? And %s>=?",
		models.TABLE_EosOtc,
		models.COLUMN_EosOtc_Available, models.COLUMN_EosOtc_Available,
		models.COLUMN_EosOtc_Transfer, models.COLUMN_EosOtc_Transfer,
		models.COLUMN_EosOtc_Uid,
		models.COLUMN_EosOtc_Transfer)

	res, err := d.Orm.Raw(sql, quantity, quantity, uid, quantity).Exec()
	if err != nil {
		common.LogFuncError("OtcTransferUnlock DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("OtcTransferUnlock DBERR: u:%d, q:%d, e:%v", uid, quantity, err)
		return
	}
	ok = true
	return
}

type ExchangerWealth struct {
	Available string `json:"available"`
	Trade     string `json:"trade"`
	Transfer  string `json:"transfer"`
}

// 获取承兑商资产 通过Euid
func (d *EosOtcDao) GetOtcExchangerWealth(euid uint64) (exchangerWealth ExchangerWealth, err error) {
	sql := fmt.Sprintf("select available,trade,transfer from %s where %s=?", models.TABLE_EosOtc, models.COLUMN_EosOtc_Uid)
	_ = d.Orm.Raw(sql, euid).QueryRow(&exchangerWealth)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}
