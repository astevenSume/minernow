package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/mitchellh/mapstructure"
	"utils/eusd/models"
)

const (
	WealthStatusLock    = iota // 被锁定
	WealthStatusWorking        //账号状态正常
)

const (
	IsNotExchanger = 0 //不是承兑商
	IsExchanger    = 1 //是承兑商
)

var WealthStatusToString = map[int]string{
	WealthStatusWorking: "working",
	WealthStatusLock:    "lock",
}

type BalanceSync struct {
	Uid     uint64 `json:"uid"`
	Account string `json:"account"`
	Balance int64  `json:"balance"`
}

type WealthDao struct {
	common.BaseDao
}

func NewWealthDao(db string) *WealthDao {
	return &WealthDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *WealthDao) Add(uid uint64, account string) (id int64, err error) {
	data := &models.EosWealth{
		Uid:     uid,
		Account: account,
		Status:  WealthStatusWorking,
		Ctime:   common.NowInt64MS(),
	}

	id, err = d.Orm.Insert(data)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

func (d *WealthDao) Info(uid uint64) (user *models.EosWealth, err error) {
	user = &models.EosWealth{
		Uid: uid,
	}
	err = d.Orm.Read(user)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			user.Uid = 0
			return
		}
		common.LogFuncError("DBERR: %v", err)
		return
	}

	return
}

func (d *WealthDao) FetchByIds(uids ...uint64) (users []*models.EosWealth, err error) {

	var maps []orm.Params
	qs := d.Orm.QueryTable(models.TABLE_EosWealth).Filter(models.COLUMN_EosAccount_Uid+"__in", uids)

	_, err = qs.Values(&maps)

	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
		}
		return
	}
	for _, params := range maps {
		item := &models.EosWealth{}
		err = mapstructure.Decode(params, item)
		if err != nil {
			common.LogFuncError("mapstructure.Decode %v failed : %v", params, err)
			return
		}
		users = append(users, item)
	}
	return
}

func (d *WealthDao) GetFromToUser(from, to uint64) (fromUser *models.EosWealth, toUser *models.EosWealth, err error) {
	users, err := d.FetchByIds(from, to)
	if err != nil {
		return
	}

	fromUser, toUser = &models.EosWealth{}, &models.EosWealth{}
	for _, u := range users {
		if u.Uid == from {
			fromUser = u
		}
		if u.Uid == to {
			toUser = u
		}
	}

	return
}

func (d *WealthDao) Update(data *models.EosWealth, fields ...string) (err error) {
	_, err = d.Orm.Update(data, fields...)
	return
}

// 转账冻结
func (d *WealthDao) TransferLock(uid uint64, quantity int64) (ok bool, err error) {
	ok = false
	sql := fmt.Sprintf("Update %s Set %s=%s-?,%s=%s+? Where %s=? And %s>=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Available, models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Trade, models.COLUMN_EosWealth_Trade,
		models.COLUMN_EosWealth_Uid,
		models.COLUMN_EosWealth_Available)

	res, err := d.Orm.Raw(sql, quantity, quantity, uid, quantity).Exec()
	if err != nil {
		common.LogFuncError("TransferLock DBErr: %v", err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return
	}
	ok = true
	return
}

//转账解冻
func (d *WealthDao) TransferUnLock(uid uint64, quantity int64) (ok bool, err error) {
	ok = false
	sql := fmt.Sprintf("Update %s Set %s=%s-?,%s=%s+? Where %s=? And %s>=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Trade, models.COLUMN_EosWealth_Trade,
		models.COLUMN_EosWealth_Available, models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Uid,
		models.COLUMN_EosWealth_Trade)

	res, err := d.Orm.Raw(sql, quantity, quantity, uid, quantity).Exec()
	if err != nil {
		common.LogFuncError("TransferUnLock DBErr: %v", err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return
	}
	ok = true
	return
}

//转入(待转入账户)
func (d *WealthDao) TransferInto(uid uint64, quantity int64) (ok bool, err error) {
	sql := fmt.Sprintf("Update %s Set %s=%s+? Where %s=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Transfer, models.COLUMN_EosWealth_Transfer,
		models.COLUMN_EosWealth_Uid)

	res, err := d.Orm.Raw(sql, quantity, uid).Exec()
	if err != nil {
		common.LogFuncError("TransferInto DBErr: %v,uid:%d", err, uid)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("TransferInto DBErr: %v,uid:%d", err, uid)
		return
	}
	ok = true
	return
}

//转账中-状态解除
func (d *WealthDao) TransferToAvailable(uid uint64, transfer, transferGame int64) (ok bool) {
	ok = false
	sql := fmt.Sprintf("Update %s Set %s=%s-?,%s=%s+?,%s=%s-?,%s=%s+? Where %s=? And %s>=? And %s>=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Transfer, models.COLUMN_EosWealth_Transfer,
		models.COLUMN_EosWealth_Available, models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_TransferGame, models.COLUMN_EosWealth_TransferGame,
		models.COLUMN_EosWealth_Game, models.COLUMN_EosWealth_Game,
		models.COLUMN_EosWealth_Uid,
		models.COLUMN_EosWealth_Transfer, models.COLUMN_EosWealth_TransferGame)

	res, err := d.Orm.Raw(sql, transfer, transfer, transferGame, transferGame, uid, transfer, transferGame).Exec()
	if err != nil {
		common.LogFuncError("TransferLock DBErr: %v", err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return
	}
	ok = true
	return
}

//确认转出，从转账冻结中扣除
func (d *WealthDao) UserOtcOut(uid uint64, quantity int64) (ok bool, err error) {
	sql := fmt.Sprintf("Update %s Set %s=%s-? Where %s=? And %s>=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Trade, models.COLUMN_EosWealth_Trade,
		models.COLUMN_EosWealth_Uid,
		models.COLUMN_EosWealth_Trade)

	res, err := d.Orm.Raw(sql, quantity, uid, quantity).Exec()
	if err != nil {
		common.LogFuncError("UserOtcOut DBErr: %v,uid:%d", err, uid)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return
	}
	ok = true
	return
}

//直接转出
func (d *WealthDao) TransferOutDirect(uid uint64, quantity int64) (ok bool, err error) {
	sql := fmt.Sprintf("Update %s Set %s=%s-? Where %s=? And %s>=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Available, models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Uid,
		models.COLUMN_EosWealth_Available)

	res, err := d.Orm.Raw(sql, quantity, uid, quantity).Exec()
	if err != nil {
		common.LogFuncError("TransferOutDirect DBErr: %v,uid:%d", err, uid)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("TransferOutDirect DBErr: %v,uid:%d NotAffected", err, uid)
		return
	}
	ok = true
	return
}

// 账号锁定
func (d *WealthDao) Lock(uid uint64) (ok bool, err error) {
	sql := fmt.Sprintf("Update %s Set %s=? Where %s=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Status,
		models.COLUMN_EosWealth_Uid)

	res, err := d.Orm.Raw(sql, WealthStatusLock, uid).Exec()
	if err != nil {
		common.LogFuncError("Wealth Lock DBERR: %v,uid:%d", err, uid)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return
	}
	ok = true
	return
}

// 账号解锁
func (d *WealthDao) UnLock(uid uint64) (ok bool, err error) {
	sql := fmt.Sprintf("Update %s Set %s=? Where %s=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Status,
		models.COLUMN_EosWealth_Uid)

	res, err := d.Orm.Raw(sql, WealthStatusWorking, uid).Exec()
	if err != nil {
		common.LogFuncError("Wealth Unlock DBERR: %v,uid:%d", err, uid)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncInfo("Wealth UnLock affected row num is 0: %v,uid:%d", num, uid)
		return
	}
	ok = true
	return
}

// 转出到OTC
func (d *WealthDao) TransferOtc(uid uint64, quantity int64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=%s-? Where %s=? And %s>=? And %s=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Available, models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Uid,
		models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Status)

	res, err := d.Orm.Raw(sql, quantity, uid, quantity, WealthStatusWorking).Exec()
	if err != nil {
		common.LogFuncError("Wealth to Otc DBERR: %v,uid:%d", err, uid)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return
	}
	ok = true
	return
}

// OTC转入
func (d *WealthDao) OtcTransferInto(uid uint64, quantity int64) (ok bool) {
	sql := fmt.Sprintf("Update %s Set %s=%s+? Where %s=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Available, models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Uid)

	res, err := d.Orm.Raw(sql, quantity, uid).Exec()
	if err != nil {
		common.LogFuncError("Otc to Wealth DBERR: %v,uid:%d, q:%d", err, uid, quantity)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("Otc to Wealth DBERR: %v,uid:%d, q:%d", err, uid, quantity)
		return
	}
	ok = true
	return
}

// 转入游戏
func (d *WealthDao) GameLock(uid uint64, quantity, transfer int64) (ok bool) {
	ok = false
	sql := fmt.Sprintf("Update %s Set %s=%s-?,%s=%s+?,%s=%s-?,%s=%s+? Where %s=? And %s>=? And %s>=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Available, models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Game, models.COLUMN_EosWealth_Game,
		models.COLUMN_EosWealth_Transfer, models.COLUMN_EosWealth_Transfer,
		models.COLUMN_EosWealth_TransferGame, models.COLUMN_EosWealth_TransferGame,
		models.COLUMN_EosWealth_Uid,
		models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Transfer)

	res, err := d.Orm.Raw(sql, quantity, quantity, transfer, transfer, uid, quantity, transfer).Exec()
	if err != nil {
		common.LogFuncError("GameLock DBERR: %v", err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return
	}
	ok = true
	return
}

// 充值到游戏
func (d *WealthDao) GameRecharge(uid uint64, quantity, transfer int64) (ok bool) {
	ok = false
	sql := fmt.Sprintf("Update %s Set %s=%s-?,%s=%s-? Where %s=? And %s>=? And %s>=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Available, models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Transfer, models.COLUMN_EosWealth_Transfer,
		models.COLUMN_EosWealth_Uid,
		models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Transfer)

	res, err := d.Orm.Raw(sql, quantity, transfer, uid, quantity, transfer).Exec()
	if err != nil {
		common.LogFuncError("GameLock DBERR: %v", err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return
	}
	ok = true
	return
}

// 解冻用户锁定的token，如有剩余还给用户
// quantity 可用token  leave剩余可用token  transfer转入中token  transferLeave 剩余转入中token
func (d *WealthDao) GameUnlock(uid uint64, available, leave, transfer, transferLeave int64) (ok bool) {
	ok = false
	sql := fmt.Sprintf("Update %s Set %s=0,%s=%s+?,%s=0,%s=%s+? Where %s=? And %s=? And %s=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Game,
		models.COLUMN_EosWealth_Available, models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_TransferGame,
		models.COLUMN_EosWealth_Transfer, models.COLUMN_EosWealth_Transfer,
		models.COLUMN_EosWealth_Uid,
		models.COLUMN_EosWealth_Game, models.COLUMN_EosWealth_TransferGame,
	)

	res, err := d.Orm.Raw(sql, leave, transferLeave, uid, available, transfer).Exec()
	if err != nil {
		common.LogFuncError("GameLock DBERR: %v, leave:%d, quantity:%d", err, available, leave)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return
	}
	ok = true
	return
}

// 用户游戏冻结转出 (用户输钱)
func (d *WealthDao) GameLockOut(uid uint64, quantity int64) (ok bool) {
	ok = false
	sql := fmt.Sprintf("Update %s Set %s=%s-? Where %s=? And %s>=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Game, models.COLUMN_EosWealth_Game,
		models.COLUMN_EosWealth_Uid,
		models.COLUMN_EosWealth_Game)

	res, err := d.Orm.Raw(sql, quantity, uid, quantity).Exec()
	if err != nil {
		common.LogFuncError("GameLockToPlatform DBERR: %v", err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("GameLockToPlatform DBERR: %v", err)
		return
	}
	ok = true
	return
}

// 平台收钱 (用户输钱 or 充值)
func (d *WealthDao) GamePlatformInto(uid uint64, quantity int64) (ok bool) {
	ok = false
	sql := fmt.Sprintf("Update %s Set %s=%s+? Where %s=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Available, models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Uid)

	res, err := d.Orm.Raw(sql, quantity, uid).Exec()
	if err != nil {
		common.LogFuncError("GamePlatformToUser DBERR: %v", err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("GamePlatformToUser DBERR: %v", err)
		return
	}
	ok = true
	return
}

// 平台转出
func (d *WealthDao) GamePlatformOut(uid uint64, quantity int64) (ok bool) {
	ok = false
	sql := fmt.Sprintf("Update %s Set %s=%s-? Where %s=? And %s>=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Available, models.COLUMN_EosWealth_Available,
		models.COLUMN_EosWealth_Uid, models.COLUMN_EosWealth_Available)

	res, err := d.Orm.Raw(sql, quantity, uid, quantity).Exec()
	if err != nil {
		common.LogFuncError("GamePlatformToUser DBERR: %v", err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("GamePlatformToUser DBERR: %v", err)
		return
	}
	ok = true
	return
}

// 用户转入 （提现）
func (d *WealthDao) GameUserInto(uid uint64, quantity int64) (ok bool) {
	ok = false
	sql := fmt.Sprintf("Update %s Set %s=%s+? Where %s=?",
		models.TABLE_EosWealth,
		models.COLUMN_EosWealth_Transfer, models.COLUMN_EosWealth_Transfer,
		models.COLUMN_EosWealth_Uid)

	res, err := d.Orm.Raw(sql, quantity, uid).Exec()
	if err != nil {
		common.LogFuncError("GamePlatformToUser DBERR: %v", err)
		return
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		common.LogFuncError("GamePlatformToUser DBERR: %v", err)
		return
	}
	ok = true
	return
}

//get all Eos Wealth
func (d *WealthDao) GetEosWealth(page int, limit int, uid uint64, status int8) (listWl []*models.EosWealth, total int64, err error) {
	qs := d.Orm.QueryTable(models.TABLE_EosWealth)
	if uid > 0 {
		qs = qs.Filter(models.COLUMN_EosAccount_Uid, uid)
	}
	if status >= 0 {
		qs = qs.Filter(models.COLUMN_EosAccount_Status, status)
	}
	total, err = qs.Count()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	_, err = qs.OrderBy("-"+models.COLUMN_EosAccount_Ctime).Limit(limit, (page-1)*limit).All(&listWl)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

//GetCount 获取总记录数
func (d *WealthDao) GetCount() (total int, err error) {

	sqlTotal := fmt.Sprintf("select count(*) from %s", models.TABLE_EosWealth)
	err = d.Orm.Raw(sqlTotal).QueryRow(&total)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//获取account
func (d *WealthDao) GetAccountForLimit(current, limit int) (list []*BalanceSync, err error) {

	sql := fmt.Sprintf("SELECT uid,account from %s where uid > 0 LIMIT ?,?", models.TABLE_EosWealth)
	_, err = d.Orm.Raw(sql, current, limit).QueryRows(&list)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//通过uid 更新余额
func (d *WealthDao) UpdateForBalanceSync(uid uint64, balance int64) (err error) {
	sql := fmt.Sprintf("UPDATE %s set balance = ? where uid = ?", models.TABLE_EosWealth)
	_, err = d.Orm.Raw(sql, balance, uid).Exec()

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//get all Eos account
func (d *WealthDao) GetEosAccount(page int, limit int) (accountList []*models.EosAccount, total int, err error) {

	sql := fmt.Sprintf("SELECT * from %s ORDER BY %s DESC LIMIT ?,?", models.TABLE_EosAccount,
		models.COLUMN_EosAccount_Ctime)
	_, err = d.Orm.Raw(sql, (page-1)*limit, limit).QueryRows(&accountList)

	sqlTotal := fmt.Sprintf("select count(*) from %s", models.TABLE_EosAccount)
	err = d.Orm.Raw(sqlTotal).QueryRow(&total)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
