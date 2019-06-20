package dao

import (
	"common"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"utils/admin/dao"
	common2 "utils/common"
	otcmodels "utils/otc/models"
	"utils/usdt/models"

	"github.com/astaxie/beego/orm"
)

const (
	STATUS_LOCKED uint8 = iota
	STATUS_WORKING
)

const (
	SweepStatusInitial = iota
	SweepStatusPreSweep
)

type AccountDao struct {
	common.BaseDao
}

func NewAccountDao(db string) *AccountDao {
	return &AccountDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AccountDaoEntity *AccountDao

// generate usdt account
func (d *AccountDao) Create(uid uint64) (account *models.UsdtAccount, err error) {

	var id uint64
	id, err = common.IdManagerGen(IdTypeUsdtAccount)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	now := common.NowInt64MS()

	// create usdt account data
	account = &models.UsdtAccount{
		Uaid:   id,
		Uid:    uid,
		Status: STATUS_WORKING,
		Ctime:  now,
	}

	if account.Sign, err = d.calcSign(account); err != nil {
		common.LogFuncError("%v", err)
		return
	}

	_, err = d.Orm.Insert(account)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// generate hot wallet account
func (d *AccountDao) CreateHotWalletAccount(pkid uint64, address string) (account *models.UsdtAccount, err error) {

	now := common.NowInt64MS()

	// create usdt account data
	account = &models.UsdtAccount{
		Uaid:    pkid,
		Uid:     pkid,
		Pkid:    pkid,
		Status:  STATUS_WORKING,
		Address: address,
		Ctime:   now,
	}

	if account.Sign, err = d.calcSign(account); err != nil {
		common.LogFuncError("%v", err)
		return
	}

	_, err = d.Orm.Insert(account)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

type AccountMsg struct {
	Address   string `json:"address"`
	Status    uint8  `json:"status"`
	Symbol    string `json:"symbol"`
	Precision int    `json:"precision"`
}

// query count
func (d *AccountDao) QueryTotal() (total int, err error) {
	var tmp int64
	tmp, err = d.Orm.QueryTable(models.TABLE_UsdtAccount).Count()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	total = int(tmp)

	return
}

// query one page
func (d *AccountDao) QueryPage(page, perPage int) (list []models.UsdtAccount, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_UsdtAccount).Limit(perPage).Offset((page - 1) * perPage).All(&list)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
		return
	}

	return
}

func (d *AccountDao) Query(uaid uint64) (account *models.UsdtAccount, err error) {
	account = &models.UsdtAccount{
		Uaid: uaid,
	}
	err = d.Orm.Read(account)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
		return
	}
	return
}

func (d *AccountDao) UpdateCols(account *models.UsdtAccount, cols ...string) (err error) {
	if account.Uid == 0 {
		return dao.ErrSql
	}

	defer func() {
		if err != nil {
			common.LogFuncError("%v", err)
		}
	}()

	if err = d.checkSign(account.Uid); err != nil {
		return
	}

	o := orm.NewOrm()
	if err = o.Begin(); err != nil {
		return
	}
	defer o.Rollback()

	_, err = o.Update(account, cols...)
	if err != nil {
		return
	}

	if err = d.updateSign(o, account.Uid); err != nil {
		return
	}

	err = o.Commit()

	return
}

func (d *AccountDao) UpdateStatus(uid uint64, status uint8) (err error) {
	if uid <= 0 {
		return dao.ErrSql
	}

	_, err = d.Orm.QueryTable(models.TABLE_UsdtAccount).Filter(models.COLUMN_UsdtAccount_Uid, uid).Update(orm.Params{
		models.COLUMN_UsdtAccount_Status: status,
	})
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *AccountDao) QueryByUid(uid uint64) (account *models.UsdtAccount, err error) {
	account = &models.UsdtAccount{
		Uid: uid,
	}
	err = d.Orm.Read(account, models.COLUMN_UsdtAccount_Uid)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("%v", err)
		return
	}

	return
}

const (
	ModifyTypeUnknown = iota
	ModifyTypeFrozen
	ModifyTypeUnfrozen
	ModifyTypeClearFrozen
	ModifyTypeMortgage
	ModifyTypeRelease
	ModifyTypeDeposit
	ModifyTypeDepositClearTransferFrozen
	ModifyTypeTransferFrozen
	ModifyTypeTransferUnfrozen
	ModifyTypeTransferClearFrozen
	ModifyTypeTransferFromPlatform
	ModifyTypeTransferServiceCharge
)

type modifyParam struct {
	AmountInteger         int64
	TransferFrozenInteger int64
}

func (d *AccountDao) modify(o orm.Ormer, uid uint64, mt int, param modifyParam) (err error) {

	if uid == 0 {
		return dao.ErrSql
	}

	defer func() {
		if err != nil {
			common.LogFuncError("%v", err)
		}
	}()

	if err = d.checkSign(uid); err != nil {
		return
	}

	var deferFunc = func() {}

	if o == nil {
		o = orm.NewOrm()
		o.Begin()
		deferFunc = func() {
			if err != nil {
				o.Rollback()
			} else {
				err = o.Commit()
			}
		}
	}

	defer deferFunc()

	switch mt {
	case ModifyTypeFrozen:
		{
			err = d.subModify(o, fmt.Sprintf("UPDATE %s SET %s=%s+?, %s=%s-?, %s=? WHERE %s=? AND %s>=?",
				d.BusinessTable(models.TABLE_UsdtAccount),
				models.COLUMN_UsdtAccount_FrozenInteger, models.COLUMN_UsdtAccount_FrozenInteger,
				models.COLUMN_UsdtAccount_AvailableInteger, models.COLUMN_UsdtAccount_AvailableInteger,
				models.COLUMN_UsdtAccount_Mtime,
				models.COLUMN_UsdtAccount_Uid,
				models.COLUMN_UsdtAccount_AvailableInteger),
				param.AmountInteger, param.AmountInteger, common.NowInt64MS(), uid, param.AmountInteger)
		}
	case ModifyTypeUnfrozen:
		{
			err = d.subModify(o, fmt.Sprintf("UPDATE %s SET %s=%s+?, %s=%s-?, %s=? WHERE %s=? AND %s>=?",
				d.BusinessTable(models.TABLE_UsdtAccount),
				models.COLUMN_UsdtAccount_AvailableInteger, models.COLUMN_UsdtAccount_AvailableInteger,
				models.COLUMN_UsdtAccount_FrozenInteger, models.COLUMN_UsdtAccount_FrozenInteger,
				models.COLUMN_UsdtAccount_Mtime,
				models.COLUMN_UsdtAccount_Uid,
				models.COLUMN_UsdtAccount_FrozenInteger),
				param.AmountInteger, param.AmountInteger, common.NowInt64MS(), uid, param.AmountInteger)
		}
	case ModifyTypeClearFrozen:
		{
			err = d.subModify(o, fmt.Sprintf("UPDATE %s SET %s=%s-?, %s=? WHERE %s=? AND %s>=?",
				d.BusinessTable(models.TABLE_UsdtAccount),
				models.COLUMN_UsdtAccount_FrozenInteger, models.COLUMN_UsdtAccount_FrozenInteger,
				models.COLUMN_UsdtAccount_Mtime,
				models.COLUMN_UsdtAccount_Uid, models.COLUMN_UsdtAccount_FrozenInteger),
				param.AmountInteger, common.NowInt64MS(), uid, param.AmountInteger)
		}
	case ModifyTypeTransferFrozen:
		{
			err = d.subModify(o, fmt.Sprintf("UPDATE %s SET %s=%s+?, %s=%s-?, %s=? WHERE %s=? AND %s>=?",
				d.BusinessTable(models.TABLE_UsdtAccount),
				models.COLUMN_UsdtAccount_TransferFrozenInteger, models.COLUMN_UsdtAccount_TransferFrozenInteger,
				models.COLUMN_UsdtAccount_AvailableInteger, models.COLUMN_UsdtAccount_AvailableInteger,
				models.COLUMN_UsdtAccount_Mtime,
				models.COLUMN_UsdtAccount_Uid,
				models.COLUMN_UsdtAccount_AvailableInteger),
				param.AmountInteger, param.AmountInteger, common.NowInt64MS(), uid, param.AmountInteger)
		}
	case ModifyTypeTransferUnfrozen:
		{
			err = d.subModify(o, fmt.Sprintf("UPDATE %s SET %s=%s+?, %s=%s-?, %s=? WHERE %s=? AND %s>=?",
				d.BusinessTable(models.TABLE_UsdtAccount),
				models.COLUMN_UsdtAccount_AvailableInteger, models.COLUMN_UsdtAccount_AvailableInteger,
				models.COLUMN_UsdtAccount_TransferFrozenInteger, models.COLUMN_UsdtAccount_TransferFrozenInteger,
				models.COLUMN_UsdtAccount_Mtime,
				models.COLUMN_UsdtAccount_Uid,
				models.COLUMN_UsdtAccount_TransferFrozenInteger),
				param.AmountInteger, param.AmountInteger, common.NowInt64MS(), uid, param.AmountInteger)
		}
	case ModifyTypeTransferClearFrozen:
		{
			err = d.subModify(o, fmt.Sprintf("UPDATE %s SET %s=%s-?, %s=? WHERE %s=? AND %s>=?",
				d.BusinessTable(models.TABLE_UsdtAccount),
				models.COLUMN_UsdtAccount_TransferFrozenInteger, models.COLUMN_UsdtAccount_TransferFrozenInteger,
				models.COLUMN_UsdtAccount_Mtime,
				models.COLUMN_UsdtAccount_Uid, models.COLUMN_UsdtAccount_TransferFrozenInteger),
				param.AmountInteger, common.NowInt64MS(), uid, param.AmountInteger)
		}
	case ModifyTypeMortgage:
		{
			err = d.subModify(o, fmt.Sprintf("UPDATE %s SET %s=%s+?, %s=%s-?, %s=? WHERE %s=? AND %s>=?",
				d.BusinessTable(models.TABLE_UsdtAccount),
				models.COLUMN_UsdtAccount_MortgagedInteger, models.COLUMN_UsdtAccount_MortgagedInteger,
				models.COLUMN_UsdtAccount_AvailableInteger, models.COLUMN_UsdtAccount_AvailableInteger,
				models.COLUMN_UsdtAccount_Mtime,
				models.COLUMN_UsdtAccount_Uid,
				models.COLUMN_UsdtAccount_AvailableInteger),
				param.AmountInteger, param.AmountInteger, common.NowInt64MS(), uid, param.AmountInteger)
		}
	case ModifyTypeRelease:
		{
			err = d.subModify(o, fmt.Sprintf("UPDATE %s SET %s=%s+?, %s=%s-?,%s=? WHERE %s=? AND %s>=?",
				d.BusinessTable(models.TABLE_UsdtAccount),
				models.COLUMN_UsdtAccount_AvailableInteger, models.COLUMN_UsdtAccount_AvailableInteger,
				models.COLUMN_UsdtAccount_MortgagedInteger, models.COLUMN_UsdtAccount_MortgagedInteger,
				models.COLUMN_UsdtAccount_Mtime,
				models.COLUMN_UsdtAccount_Uid,
				models.COLUMN_UsdtAccount_MortgagedInteger),
				param.AmountInteger, param.AmountInteger, common.NowInt64MS(), uid, param.AmountInteger)
		}
	case ModifyTypeDeposit:
		{
			err = d.subModify(o, fmt.Sprintf("UPDATE %s SET %s=%s+?,%s=? WHERE %s=?",
				d.BusinessTable(models.TABLE_UsdtAccount),
				models.COLUMN_UsdtAccount_AvailableInteger, models.COLUMN_UsdtAccount_AvailableInteger,
				models.COLUMN_UsdtAccount_Mtime,
				models.COLUMN_UsdtAccount_Uid),
				param.AmountInteger, common.NowInt64MS(), uid)
		}
	case ModifyTypeDepositClearTransferFrozen:
		{
			err = d.subModify(o, fmt.Sprintf("UPDATE %s SET %s=%s+?, %s=%s+?, %s=? WHERE %s=? AND %s>=?",
				d.BusinessTable(models.TABLE_UsdtAccount),
				models.COLUMN_UsdtAccount_AvailableInteger, models.COLUMN_UsdtAccount_AvailableInteger,
				models.COLUMN_UsdtAccount_TransferFrozenInteger, models.COLUMN_UsdtAccount_TransferFrozenInteger,
				models.COLUMN_UsdtAccount_Mtime,
				models.COLUMN_UsdtAccount_Uid, models.COLUMN_UsdtAccount_TransferFrozenInteger),
				param.AmountInteger, param.TransferFrozenInteger, common.NowInt64MS(), uid, param.TransferFrozenInteger)
		}
	case ModifyTypeTransferFromPlatform:
		{
			err = d.subModify(o,
				fmt.Sprintf("UPDATE %s SET  %s=%s-?, %s=%s+?, %s=? WHERE %s=? AND %s >=?",
					d.BusinessTable(models.TABLE_UsdtAccount),
					models.COLUMN_UsdtAccount_TransferFrozenInteger, models.COLUMN_UsdtAccount_TransferFrozenInteger,
					models.COLUMN_UsdtAccount_OwnedByPlatformInteger, models.COLUMN_UsdtAccount_OwnedByPlatformInteger,
					models.COLUMN_UsdtAccount_Mtime,
					// where
					models.COLUMN_UsdtAccount_Uid,
					models.COLUMN_UsdtAccount_TransferFrozenInteger,
				),
				param.AmountInteger, param.AmountInteger, common.NowInt64MS(), uid, param.AmountInteger, param.AmountInteger)

		}
	case ModifyTypeTransferServiceCharge:
		{

			err = d.subModify(o,
				fmt.Sprintf("UPDATE %s SET  %s=%s+?, %s=%s+?, %s=? WHERE %s=? ",
					d.BusinessTable(models.TABLE_UsdtAccount),
					models.COLUMN_UsdtAccount_WaitingCashSweepInteger, models.COLUMN_UsdtAccount_WaitingCashSweepInteger,
					models.COLUMN_UsdtAccount_OwnedByPlatformInteger, models.COLUMN_UsdtAccount_OwnedByPlatformInteger,
					models.COLUMN_UsdtAccount_Mtime,
					// where
					models.COLUMN_UsdtAccount_Uid,
				),
				param.AmountInteger, param.AmountInteger, common.NowInt64MS(), uid)

		}
	default:
		err = errors.New("invalid modify type")
		return
	}

	if err != nil {
		return
	}

	if err = d.updateSign(o, uid); err != nil {
		return
	}
	return
}

func (d *AccountDao) subModify(o orm.Ormer, s string, args ...interface{}) (err error) {
	if o == nil {
		o = d.Orm
	}

	var result sql.Result
	result, err = o.Raw(s, args).Exec()

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var n int64
	if n, err = result.RowsAffected(); err != nil {
		common.LogFuncError("%v", err)
		return
	} else if n != 1 {
		err = common.ErrNoRowEffected
		common.LogFuncError("effected row %d no equal to 1", n)
		return
	}

	return
}

// frozen
func (d *AccountDao) Frozen(uid uint64, amount int64) (err error) {
	return d.modify(nil, uid, ModifyTypeFrozen, modifyParam{AmountInteger: amount})
}

// frozen with transaction
func (d *AccountDao) FrozenTx(uid uint64, amount int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeFrozen, modifyParam{AmountInteger: amount})
}

// unfrozen
func (d *AccountDao) Unfrozen(uid uint64, amount int64) (err error) {
	return d.modify(nil, uid, ModifyTypeUnfrozen, modifyParam{AmountInteger: amount})
}

// unfrozen with transaction
func (d *AccountDao) UnfrozenTx(uid uint64, amount int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeUnfrozen, modifyParam{AmountInteger: amount})
}

// clear frozen
func (d *AccountDao) ClearFrozen(uid uint64, amount int64) (err error) {
	return d.modify(nil, uid, ModifyTypeClearFrozen, modifyParam{AmountInteger: amount})
}

// clear frozen with transaction
func (d *AccountDao) ClearFrozenTx(uid uint64, amount int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeClearFrozen, modifyParam{AmountInteger: amount})
}

// mortgage
func (d *AccountDao) Mortgage(uid uint64, amount int64) (err error) {
	return d.modify(nil, uid, ModifyTypeMortgage, modifyParam{AmountInteger: amount})
}

// mortgage with transaction
func (d *AccountDao) MortgageTx(uid uint64, amount int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeMortgage, modifyParam{AmountInteger: amount})
}

// release
func (d *AccountDao) Release(uid uint64, amount int64) (err error) {
	return d.modify(nil, uid, ModifyTypeRelease, modifyParam{AmountInteger: amount})
}

// release with transaction
func (d *AccountDao) ReleaseTx(uid uint64, amount int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeRelease, modifyParam{AmountInteger: amount})
}

// deposit
func (d *AccountDao) Deposit(uid uint64, amount int64) (err error) {
	return d.modify(nil, uid, ModifyTypeDeposit, modifyParam{AmountInteger: amount})
}

// deposit with transaction
func (d *AccountDao) DepositTx(uid uint64, amount int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeDeposit, modifyParam{AmountInteger: amount})
}

// deposit and transfer_frozen
func (d *AccountDao) DepositAndClearTransferFrozen(uid uint64, amount, transferFrozen int64) (err error) {
	return d.modify(nil, uid, ModifyTypeDepositClearTransferFrozen, modifyParam{AmountInteger: amount, TransferFrozenInteger: transferFrozen})
}

// deposit and transfer_frozen with transaction
func (d *AccountDao) DepositAndClearTransferFrozenTx(uid uint64, amount, transferFrozen int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeDepositClearTransferFrozen, modifyParam{AmountInteger: amount, TransferFrozenInteger: transferFrozen})
}

// transfer frozen
func (d *AccountDao) TransferFrozen(uid uint64, amount int64) (err error) {
	return d.modify(nil, uid, ModifyTypeTransferFrozen, modifyParam{AmountInteger: amount})
}

// transfer frozen with transaction
func (d *AccountDao) TransferFrozenTx(uid uint64, amount int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeTransferFrozen, modifyParam{AmountInteger: amount})
}

// transfer unfrozen
func (d *AccountDao) TransferUnfrozen(uid uint64, amount int64) (err error) {
	return d.modify(nil, uid, ModifyTypeTransferUnfrozen, modifyParam{AmountInteger: amount})
}

// transfer unfrozen with transaction
func (d *AccountDao) TransferUnfrozenTx(uid uint64, amount int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeTransferUnfrozen, modifyParam{AmountInteger: amount})
}

// clear transfer frozen
func (d *AccountDao) TransferClearFrozen(uid uint64, amount int64) (err error) {
	return d.modify(nil, uid, ModifyTypeTransferClearFrozen, modifyParam{AmountInteger: amount})
}

// clear transfer frozen with transaction
func (d *AccountDao) TransferClearFrozenTx(uid uint64, amount int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeTransferClearFrozen, modifyParam{AmountInteger: amount})
}

// TransferFromPlatform 从平台账号扣款,无事务
func (d *AccountDao) TransferFromPlatform(uid uint64, amount int64) (err error) {
	return d.modify(nil, uid, ModifyTypeTransferFromPlatform, modifyParam{AmountInteger: amount})
}

// TransferFromPlatform 从平台账号扣款,有事务
func (d *AccountDao) TransferFromPlatformTx(uid uint64, amount int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeTransferFromPlatform, modifyParam{AmountInteger: amount})
}

// TransferServiceCharge 转账服务费,无事务
func (d *AccountDao) TransferServiceCharge(uid uint64, amount int64) (err error) {
	return d.modify(nil, uid, ModifyTypeTransferServiceCharge, modifyParam{AmountInteger: amount})
}

// TransferServiceChargeTx 转账服务费,有事务
func (d *AccountDao) TransferServiceChargeTx(uid uint64, amount int64, o orm.Ormer) (err error) {
	return d.modify(o, uid, ModifyTypeTransferServiceCharge, modifyParam{AmountInteger: amount})
}

type UsdtAddress struct {
	Id      string `json:"id"`
	Address string `json:"address"`
	Uid     string `json:"uid"`
	Mobile  string `json:"mobile"`
	Ctime   int64  `json:"ctime"`
	Utime   int64  `json:"utime"`
}

func (d *AccountDao) GetAddress(uid uint64) (usdtAddress UsdtAddress, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT T1.id, T1.address, T1.uid, T2.mobile, T1.ctime, T1.utime FROM (SELECT %s "+
		"AS id, ? AS %s, %s, %s , %s AS utime FROM %s WHERE %s=?) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s",
		models.COLUMN_UsdtAccount_Uaid, models.COLUMN_UsdtAccount_Uid,
		models.COLUMN_UsdtAccount_Address, models.COLUMN_UsdtAccount_Ctime, models.COLUMN_UsdtAccount_Mtime,
		models.TABLE_UsdtAccount, models.COLUMN_UsdtAccount_Uid,
		otcmodels.TABLE_User, models.COLUMN_UsdtAccount_Uid, otcmodels.COLUMN_User_Uid), uid, uid).QueryRow(&usdtAddress)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *AccountDao) GetPageAddress(page int, limit int, uid uint64) (total int, usdtAddress []UsdtAddress, err error) {
	var param []interface{}
	var cond string
	if uid > 0 {
		cond = " AND " + models.COLUMN_UsdtAccount_Uid + "=? "
		param = append(param, uid)
	}
	sqlTotal := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE 1=1 %s", models.TABLE_UsdtAccount, cond)
	err = d.Orm.Raw(sqlTotal, param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error: %v", err)
		return
	}
	param = append(param, limit)
	param = append(param, (page-1)*limit)

	_, err = d.Orm.Raw(fmt.Sprintf("SELECT T1.id, T1.address, T1.uid, T2.mobile, T1.ctime, T1.utime FROM (SELECT"+
		" %s AS id, %s, %s, %s , %s AS utime FROM %s WHERE 1=1 %s LIMIT ? OFFSET ?) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s",
		models.COLUMN_UsdtAccount_Uaid, models.COLUMN_UsdtAccount_Uid,
		models.COLUMN_UsdtAccount_Address, models.COLUMN_UsdtAccount_Ctime, models.COLUMN_UsdtAccount_Mtime,
		models.TABLE_UsdtAccount, cond, otcmodels.TABLE_User,
		models.COLUMN_UsdtAccount_Uid, otcmodels.COLUMN_User_Uid), param...).QueryRows(&usdtAddress)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *AccountDao) QueryByUaid(uaid uint64) (usdtAccount models.UsdtAccount, err error) {
	usdtAccount = models.UsdtAccount{
		Uaid: uaid,
	}
	err = d.Orm.Read(&usdtAccount, models.COLUMN_UsdtAccount_Uaid)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
		return
	}
	return
}

type UsdtAccount struct {
	Uaid                string `orm:"column(uaid);pk" json:"uaid,omitempty"`
	Uid                 string `orm:"column(uid)" json:"uid,omitempty"`
	Status              uint8  `orm:"column(status)" json:"status,omitempty"`
	AvailableInteger    int64  `orm:"column(available_integer)" json:"available_integer,omitempty"`
	FrozenInteger       int64  `orm:"column(frozen_integer)" json:"frozen_integer,omitempty"`
	MortgagedInteger    int64  `orm:"column(mortgaged_integer)" json:"mortgaged_integer,omitempty"`
	BtcAvailableInteger int64  `orm:"column(btc_available_integer)" json:"btc_available_integer,omitempty"`
	BtcFrozenInteger    int64  `orm:"column(btc_frozen_integer)" json:"btc_frozen_integer,omitempty"`
	BtcMortgagedInteger int64  `orm:"column(btc_mortgaged_integer)" json:"btc_mortgaged_integer,omitempty"`
	Pkid                uint64 `orm:"column(pkid)" json:"pkid,omitempty"`
	Address             string `orm:"column(address);size(100)" json:"address,omitempty"`
	Ctime               int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Mtime               int64  `orm:"column(mtime)" json:"mtime,omitempty"`
	Mobile              string `orm:"column(mobile)" json:"mobile,omitempty"`
}

//后台分页条件查询
func (d *AccountDao) QueryByPage(page int, limit int, uid uint64, status int8) (total int64, usdtAccounts []UsdtAccount, err error) {
	var cond string
	var param []interface{}
	if uid > 0 {
		cond = cond + " AND " + models.COLUMN_UsdtAccount_Uid + " =? "
		param = append(param, uid)
	}
	if status >= 0 {
		cond = cond + " AND " + models.COLUMN_UsdtAccount_Status + " =? "
		param = append(param, status)
	}

	sqlTotal := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE 1=1 %s ", models.TABLE_UsdtAccount, cond)
	err = d.Orm.Raw(sqlTotal, param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	param = append(param, limit)
	param = append(param, (page-1)*limit)

	sqlQuery := fmt.Sprintf("SELECT T1.*,T2.%s FROM ((SELECT * FROM %s WHERE 1=1 %s LIMIT ? OFFSET ?) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s)",
		otcmodels.COLUMN_User_Mobile, models.TABLE_UsdtAccount, cond, otcmodels.TABLE_User, models.COLUMN_UsdtAccount_Uid, otcmodels.COLUMN_User_Uid)
	_, err = d.Orm.Raw(sqlQuery, param...).QueryRows(&usdtAccounts)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

// @Description query unrelated records.
func (d *AccountDao) QueryUnrelated(page, limit int) (list []models.UsdtAccount, err error) {
	return d.queryByIsRelated(page, limit, Unrelated)
}

// @Description query total number of unrelated records.
func (d *AccountDao) QueryTotalUnrelated() (total int, err error) {
	return d.queryTotalByIsRelated(Unrelated)
}

// @Description query related records.
func (d *AccountDao) QueryRelated(page, limit int) (list []models.UsdtAccount, err error) {
	return d.queryByIsRelated(page, limit, Related)
}

// @Description query total number of related records.
func (d *AccountDao) QueryTotalRelated() (total int, err error) {
	return d.queryTotalByIsRelated(Related)
}

// @Description update related columns.
func (d *AccountDao) UpdateRelation(l models.UsdtAccount) (err error) {

	if l.Uid == 0 {
		return dao.ErrSql
	}
	if err = d.checkSign(l.Uid); err != nil {
		return
	}

	o := orm.NewOrm()
	if err = o.Begin(); err != nil {
		return
	}
	defer o.Rollback()

	var n int64
	n, err = o.Update(&l, models.COLUMN_UsdtAccount_Pkid, models.COLUMN_UsdtAccount_Address, models.COLUMN_UsdtAccount_Mtime)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	if n != 1 {
		err = common.ErrNoRowEffected
		common.LogFuncError("%v", err)
		return
	}

	if err = d.updateSign(o, l.Uid); err != nil {
		return
	}

	err = o.Commit()

	return
}

const (
	Related   = true
	Unrelated = !Related
)

func (d *AccountDao) queryTotalByIsRelated(isRelated bool) (count int, err error) {
	var tmp int64
	tmp, err = d.relatedFilter(isRelated).Count()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	count = int(tmp)

	return
}

func (d *AccountDao) queryByIsRelated(page, limit int, isRelated bool) (list []models.UsdtAccount, err error) {
	_, err = d.relatedFilter(isRelated).Limit(limit).Offset((page - 1) * limit).All(&list)
	if err != nil {
		common.LogFuncError("%v", err)
	}

	return
}

const (
	SystemAccountBoundary = 100000
)

func (d *AccountDao) relatedFilter(isRelated bool) (qs orm.QuerySeter) {
	qs = d.Orm.QueryTable(models.TABLE_UsdtAccount)
	if isRelated {
		qs = qs.Filter(fmt.Sprintf("%s__gt", models.COLUMN_UsdtAccount_Uaid), SystemAccountBoundary).Filter(fmt.Sprintf("%s__gt", models.COLUMN_UsdtAccount_Pkid), 0)
	} else {
		qs = qs.Filter(fmt.Sprintf("%s__gt", models.COLUMN_UsdtAccount_Uaid), SystemAccountBoundary).Filter(models.COLUMN_UsdtAccount_Pkid, 0)
	}

	return
}

func (d *AccountDao) QueryWalletAccount(beginPkid, endPkid uint64) (list []models.UsdtAccount, err error) {

	_, err = d.Orm.QueryTable(models.TABLE_UsdtAccount).
		Filter(fmt.Sprintf("%s__gte", models.COLUMN_UsdtAccount_Pkid), beginPkid).
		Filter(fmt.Sprintf("%s__lte", models.COLUMN_UsdtAccount_Pkid), endPkid).
		All(&list)

	return
}

const (
	ACCOUNT_SIGN_SALT = "ACCOUNT_SIGN_SALT"
)

func (d *AccountDao) calcSign(account *models.UsdtAccount) (sign string, err error) {

	hash := md5.New()

	if _, err = hash.Write([]byte(fmt.Sprint(
		ACCOUNT_SIGN_SALT,
		account.Uaid,
		account.Uid,
		account.AvailableInteger,
		account.FrozenInteger,
		account.TransferFrozenInteger,
		account.MortgagedInteger,
		account.BtcAvailableInteger,
		account.BtcFrozenInteger,
		account.BtcMortgagedInteger,
		account.WaitingCashSweepInteger,
		account.CashSweepInteger,
		account.OwnedByPlatformInteger,
		account.SweepStatus,
		account.Pkid, account.Address,
		account.Ctime,
		account.Mtime,
		ACCOUNT_SIGN_SALT,
	))); err != nil {
		common.LogFuncError("%v", err)
		return
	}
	sign = hex.EncodeToString(hash.Sum(nil))

	hash.Reset()

	if _, err = hash.Write([]byte(sign + ACCOUNT_SIGN_SALT)); err != nil {
		common.LogFuncError("%v", err)
		return
	}

	sign = hex.EncodeToString(hash.Sum(nil))

	return
}

func (d *AccountDao) checkSign(uid uint64) (err error) {

	var (
		account *models.UsdtAccount
		sign    string
	)

	defer func() {
		if err != nil {
			common.LogFuncError("%v", err)
		}
	}()

	if account, err = d.QueryByUid(uid); err != nil {
		return
	}

	if sign, err = d.calcSign(account); err != nil {
		return
	}

	if sign != account.Sign {
		// TODO: 校验失败 需要告警
		err = fmt.Errorf("account data has been tampered")

		smsErr := common2.SmsWarning(dao.ConfigWarningTypeUsdtAccountTampered, map[string]string{
			"UID": fmt.Sprint(uid),
		})
		if smsErr != nil {
			common.LogFuncError("send warning sms failed :%v", err)
		}

		return
	}
	// USDT账户(uid)数据签名校验失败,账户有被篡改风险
	return
}

func (d *AccountDao) updateSign(o orm.Ormer, uid uint64) (err error) {
	if o == nil {
		return fmt.Errorf("ormer is nil")
	}

	var account *models.UsdtAccount

	account.Uid = uid
	if err = o.Read(account, models.COLUMN_UsdtAccount_Uid); err != nil {
		return
	}

	if account.Sign, err = d.calcSign(account); err != nil {
		return
	}

	_, err = o.Update(account, models.COLUMN_UsdtAccount_Sign)
	return
}

func (d *AccountDao) QueryPreSweep(limitAmount int64) (account models.UsdtAccount, err error) {

	err = d.Orm.Raw(
		fmt.Sprintf("select * from %s where (%s - %s)>=? and %s = ?",
			models.TABLE_UsdtAccount,
			models.COLUMN_UsdtAccount_AvailableInteger,
			models.COLUMN_UsdtAccount_CashSweepInteger,
			models.COLUMN_UsdtAccount_SweepStatus),
		limitAmount,
		SweepStatusInitial,
	).QueryRow(&account)

	return
}

func (d *AccountDao) PreSweep(account models.UsdtAccount, limitAmount int64) (err error) {

	if err = d.checkSign(account.Uid); err != nil {
		return
	}

	var (
		rawSQL string
		o      orm.Ormer
	)

	o = orm.NewOrm()
	if err = o.Begin(); err != nil {
		return
	}
	defer o.Rollback()

	rawSQL = fmt.Sprintf("update %s set %s=(%s + %s - %s) , %s = ? where %s = ? and (%s - %s)>=? and %s = ?",
		models.TABLE_UsdtAccount,
		models.COLUMN_UsdtAccount_WaitingCashSweepInteger, models.COLUMN_UsdtAccount_WaitingCashSweepInteger, models.COLUMN_UsdtAccount_AvailableInteger, models.COLUMN_UsdtAccount_CashSweepInteger,
		models.COLUMN_UsdtAccount_SweepStatus,

		models.COLUMN_UsdtAccount_Uid,
		models.COLUMN_UsdtAccount_AvailableInteger, models.COLUMN_UsdtAccount_CashSweepInteger,
		models.COLUMN_UsdtAccount_SweepStatus,
	)

	result, err := o.Raw(rawSQL, SweepStatusPreSweep, account.Uid, limitAmount, SweepStatusInitial).Exec()
	if err != nil {
		return
	}
	row, err := result.RowsAffected()
	if err != nil {
		return
	}
	if row != 1 {
		return common.ErrNoRowEffected
	}

	if err = d.updateSign(o, account.Uid); err != nil {
		return
	}

	return o.Commit()

}

func (d *AccountDao) QueryTotalWaitingSweep(limitAmount int64) (count int, err error) {

	err = d.Orm.Raw(
		fmt.Sprintf("select count(0) as totalCount from %s where %s >= ? and %s = ? ",
			models.TABLE_UsdtAccount,
			models.COLUMN_UsdtAccount_WaitingCashSweepInteger,
			models.COLUMN_UsdtAccount_SweepStatus),
		limitAmount,
		SweepStatusPreSweep,
	).QueryRow(&count)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// SweepRecord 归集记录模型
type SweepRecord struct {
	Uaid                    uint64 `orm:"column(uaid);pk" json:"uaid,omitempty"`
	Uid                     uint64 `orm:"column(uid)" json:"uid,omitempty"`
	AvailableInteger        int64  `orm:"column(available_integer)" json:"available_integer,omitempty"`
	WaitingCashSweepInteger int64  `orm:"column(waiting_cash_sweep_integer)" json:"waiting_cash_sweep_integer,omitempty"`
	CashSweepInteger        int64  `orm:"column(cash_sweep_integer)" json:"cash_sweep_integer,omitempty"`
	OwnedByPlatformInteger  int64  `orm:"column(owned_by_platform_integer)" json:"owned_by_platform_integer,omitempty"`
	Pkid                    uint64 `orm:"column(pkid)" json:"pkid,omitempty"`
	Address                 string `orm:"column(address);size(100)" json:"address,omitempty"`
}

func (d *AccountDao) QueryWaitingSweep(limitAmount int64, page, limit int) (records []SweepRecord, err error) {

	var qbQuery orm.QueryBuilder
	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	qbQuery.
		Select(models.COLUMN_UsdtAccount_Uaid, models.COLUMN_UsdtAccount_Uid, models.COLUMN_UsdtAccount_Status,
			models.COLUMN_UsdtAccount_AvailableInteger, models.COLUMN_UsdtAccount_WaitingCashSweepInteger, models.COLUMN_UsdtAccount_CashSweepInteger,
			models.COLUMN_UsdtAccount_OwnedByPlatformInteger, models.COLUMN_UsdtAccount_Pkid, models.COLUMN_UsdtAccount_Address).
		From(models.TABLE_UsdtAccount).
		Where(fmt.Sprintf("%s >= ? and %s = ? ", models.COLUMN_UsdtAccount_WaitingCashSweepInteger, models.COLUMN_UsdtAccount_SweepStatus)).
		OrderBy(models.COLUMN_UsdtAccount_WaitingCashSweepInteger).Desc().Limit(limit).Offset((page - 1) * limit)
	_, err = d.Orm.Raw(qbQuery.String(), limitAmount, SweepStatusPreSweep).QueryRows(&records)
	if err != nil {
		return
	}
	return

}

func (d *AccountDao) FinishSweep(uid uint64, sweepAmount int64) (err error) {

	if err = d.checkSign(uid); err != nil {
		return
	}

	var (
		rawSQL string
		rows   int64
	)

	o := orm.NewOrm()
	if err = o.Begin(); err != nil {
		return
	}
	defer o.Rollback()

	rawSQL = fmt.Sprintf("update %s set %s=%s - ? , %s=%s + ? - %s , %s=0 , %s = ? where %s >= ? and (%s + ?) >= %s and %s = ?",
		models.TABLE_UsdtAccount,
		models.COLUMN_UsdtAccount_WaitingCashSweepInteger, models.COLUMN_UsdtAccount_WaitingCashSweepInteger,
		models.COLUMN_UsdtAccount_CashSweepInteger, models.COLUMN_UsdtAccount_CashSweepInteger, models.COLUMN_UsdtAccount_OwnedByPlatformInteger,
		models.COLUMN_UsdtAccount_OwnedByPlatformInteger,
		models.COLUMN_UsdtAccount_SweepStatus,
		// where
		models.COLUMN_UsdtAccount_WaitingCashSweepInteger,
		models.COLUMN_UsdtAccount_CashSweepInteger,
		models.COLUMN_UsdtAccount_OwnedByPlatformInteger,
		models.COLUMN_UsdtAccount_Uid,
	)
	res, err := o.Raw(rawSQL, sweepAmount, sweepAmount, SweepStatusInitial, sweepAmount, sweepAmount, uid).Exec()
	if err != nil {
		return
	}
	rows, err = res.RowsAffected()
	if err != nil {
		return
	}
	if rows <= 0 {
		err = fmt.Errorf("0 rows affected")
	}

	if err = d.updateSign(o, uid); err != nil {
		return
	}

	return o.Commit()
}
