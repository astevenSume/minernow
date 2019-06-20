package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/otc/models"
)

const (
	PayModeWechat uint8 = 1 //按位递增
	PayModeAli    uint8 = 2
	PayModeBank   uint8 = 4
)

type PaymentMethodDao struct {
	common.BaseDao
}

func NewPaymentMethodDao(db string) *PaymentMethodDao {
	return &PaymentMethodDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var PaymentMethodDaoEntity *PaymentMethodDao

// get pay modes by uid
func (d *PaymentMethodDao) QueryByUid(uid, pmid uint64, mtype uint8) (list []*models.PaymentMethod, err error) {
	querySeter := d.Orm.QueryTable(models.TABLE_PaymentMethod)

	if pmid > 0 {
		querySeter = querySeter.Filter(models.COLUMN_PaymentMethod_Pmid, pmid)
	}

	if mtype > 0 {
		querySeter = querySeter.Filter(models.COLUMN_PaymentMethod_MType, mtype)
	}
	_, err = querySeter.
		Filter(models.COLUMN_PaymentMethod_Uid, uid).
		OrderBy(models.COLUMN_PaymentMethod_MType).
		All(&list, models.COLUMN_PaymentMethod_Pmid,
			models.COLUMN_PaymentMethod_MType,
			models.COLUMN_PaymentMethod_Name,
			models.COLUMN_PaymentMethod_Account,
			models.COLUMN_PaymentMethod_Bank,
			models.COLUMN_PaymentMethod_BankBranch,
			models.COLUMN_PaymentMethod_QRCode,
			models.COLUMN_PaymentMethod_Ord,
			models.COLUMN_PaymentMethod_Status,
			models.COLUMN_PaymentMethod_LowMoneyPerTxLimit,
			models.COLUMN_PaymentMethod_HighMoneyPerTxLimit,
			models.COLUMN_PaymentMethod_TimesPerDayLimit,
			models.COLUMN_PaymentMethod_MoneyPerDayLimit,
			models.COLUMN_PaymentMethod_MoneySumLimit,
			models.COLUMN_PaymentMethod_MoneyToday,
			models.COLUMN_PaymentMethod_TimesToday,
			models.COLUMN_PaymentMethod_MoneySum,
			models.COLUMN_PaymentMethod_UseTime)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
		}
		common.LogFuncError("%v", err)
		return
	}

	return
}

// add a payment method
func (d *PaymentMethodDao) Add(uid uint64, mType uint8, name, account, bank, bankBranch, qrCode, qrCodeContent string) (pm *models.PaymentMethod, err error) {
	// todo check if exist

	//all ord +1 (it is this safe )
	_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=%s+1 WHERE %s=?",
		models.TABLE_PaymentMethod, models.COLUMN_PaymentMethod_Ord, models.COLUMN_PaymentMethod_Ord, models.COLUMN_PaymentMethod_Uid), uid).Exec()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var id uint64
	id, err = common.IdManagerGen(IdTypePaymentMethod)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	pm = &models.PaymentMethod{
		Pmid:          id,
		Uid:           uid,
		MType:         mType,
		Ord:           1,
		Name:          name,
		Account:       account,
		Bank:          bank,
		BankBranch:    bankBranch,
		QRCode:        qrCode,
		QRCodeContent: qrCodeContent,
		Status:        PaymentMethodStatusActivated, //default is activated.
	}

	pm.Ctime = common.NowUint32()
	pm.Mtime = pm.Ctime
	pm.LowMoneyPerTxLimit = Config.PaymentMethod_LowMoneyPerTx
	pm.HighMoneyPerTxLimit = Config.PaymentMethod_HighMoneyPerTx
	pm.TimesPerDayLimit = Config.PaymentMethod_TimesPerDay
	pm.MoneyPerDayLimit = Config.PaymentMethod_MoneyPerDay
	pm.MoneySumLimit = Config.PaymentMethod_MoneySum

	_, err = d.Orm.Insert(pm)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *PaymentMethodDao) Remove(uid, pmId uint64) (err error) {
	// check if pm's owner is uid
	pm := &models.PaymentMethod{
		Pmid: pmId,
		Uid:  uid,
	}

	var n int64
	n, err = d.Orm.Delete(pm, models.COLUMN_PaymentMethod_Pmid, models.COLUMN_PaymentMethod_Uid)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	if n != 1 {
		err = orm.ErrNoRows
		common.LogFuncError("payment method of pmid(%d) uid(%d) no found.", pmId, uid)
		return
	}

	return
}

type Pmid2Ord struct {
	Pmid uint64 `json:"id"`
	Ord  uint32 `json:"ord"`
}

func (d *PaymentMethodDao) ReOrder(uid uint64, list []Pmid2Ord) (err error) {
	s := fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=? AND %s=?",
		models.TABLE_PaymentMethod, models.COLUMN_PaymentMethod_Ord, models.COLUMN_PaymentMethod_Pmid, models.COLUMN_PaymentMethod_Uid)

	for _, v := range list {
		_, err = d.Orm.Raw(s, v.Ord, v.Pmid, uid).Exec()
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

const (
	PaymentMethodStatusDeactivated uint8 = iota
	PaymentMethodStatusActivated
)

func (d *PaymentMethodDao) ChangeStatus(pmid uint64, status uint8) (err error) {
	pm := models.PaymentMethod{
		Pmid:   pmid,
		Status: status,
	}
	_, err = d.Orm.Update(&pm, models.COLUMN_PaymentMethod_Status)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *PaymentMethodDao) Edit(pm models.PaymentMethod, columns ...string) (err error) {
	_, err = d.Orm.Update(&pm, columns...)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// 获取用户所有开启的收款方式信息
func (d *PaymentMethodDao) FetchByUid(uid uint64) (list []*models.PaymentMethod) {
	s := fmt.Sprintf("Select * from %s Where %s=? AND %s=?", models.TABLE_PaymentMethod,
		models.COLUMN_PaymentMethod_Uid, models.COLUMN_PaymentMethod_Status)

	list = []*models.PaymentMethod{}
	_, err := d.Orm.Raw(s, uid, PaymentMethodStatusActivated).QueryRows(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		return
	}

	return
}

// 获取指定的用户有开启的收款方式信息
func (d *PaymentMethodDao) FetchByUidPayType(uid uint64, pay uint8) (list []*models.PaymentMethod) {
	s := fmt.Sprintf("Select * from %s Where %s=? And %s=? And %s=? Order By %s DESC", models.TABLE_PaymentMethod,
		models.COLUMN_PaymentMethod_Uid, models.COLUMN_PaymentMethod_MType, models.COLUMN_PaymentMethod_Status,
		models.COLUMN_PaymentMethod_Ord)

	list = []*models.PaymentMethod{}
	_, err := d.Orm.Raw(s, uid, pay, PaymentMethodStatusActivated).QueryRows(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

// 更新收款方式使用
func (d *PaymentMethodDao) Use(pmid uint64, amount int64, TimesToday int64) bool {
	s := fmt.Sprintf("Update %s set %s=%s+?,%s=%s+?,%s=%s+1 Where %s=? And %s=?", models.TABLE_PaymentMethod,
		models.COLUMN_PaymentMethod_MoneyToday, models.COLUMN_PaymentMethod_MoneyToday,
		models.COLUMN_PaymentMethod_MoneySum, models.COLUMN_PaymentMethod_MoneySum,
		models.COLUMN_PaymentMethod_TimesToday, models.COLUMN_PaymentMethod_TimesToday,
		models.COLUMN_PaymentMethod_Pmid,
		models.COLUMN_PaymentMethod_TimesToday,
	)

	res, err := d.Orm.Raw(s, amount, amount, pmid, TimesToday).Exec()
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return false
	}

	if n, _ := res.RowsAffected(); n != 1 {
		return false
	}

	return true
}

// 回滚收款方式使用
func (d *PaymentMethodDao) RollUse(pmid uint64, amount int64) bool {
	s := fmt.Sprintf("Update %s set %s=%s-?,%s=%s-?,%s=%s-1 Where %s=?", models.TABLE_PaymentMethod,
		models.COLUMN_PaymentMethod_MoneyToday, models.COLUMN_PaymentMethod_MoneyToday,
		models.COLUMN_PaymentMethod_MoneySum, models.COLUMN_PaymentMethod_MoneySum,
		models.COLUMN_PaymentMethod_TimesToday, models.COLUMN_PaymentMethod_TimesToday,
		models.COLUMN_PaymentMethod_Pmid,
	)

	res, err := d.Orm.Raw(s, amount, amount, pmid).Exec()
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return false
	}

	if n, _ := res.RowsAffected(); n != 1 {
		return false
	}

	return true
}

// 刷新使用信息
func (d *PaymentMethodDao) FlushUse(pmid uint64, today uint32) bool {
	s := fmt.Sprintf("Update %s set %s=0,%s=0 Where %s=? And %s<?", models.TABLE_PaymentMethod,
		models.COLUMN_PaymentMethod_MoneyToday,
		models.COLUMN_PaymentMethod_TimesToday,
		models.COLUMN_PaymentMethod_Pmid,
		models.COLUMN_PaymentMethod_Mtime,
	)

	res, err := d.Orm.Raw(s, pmid, today).Exec()
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return false
	}

	if n, _ := res.RowsAffected(); n != 1 {
		return false
	}

	return true
}

func (d *PaymentMethodDao) Info(id uint64) (data *models.PaymentMethod) {
	data = &models.PaymentMethod{
		Pmid: id,
	}

	err := d.Orm.Read(data)
	if err != nil {
		if err == orm.ErrNoRows {
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

func (d *PaymentMethodDao) ResetToday(pmid uint64) (err error) {
	data := &models.PaymentMethod{
		Pmid:       pmid,
		UseTime:    common.NowUint32(),
		TimesToday: 0,
		MoneyToday: 0,
	}

	_, err = d.Orm.Update(data, models.COLUMN_PaymentMethod_UseTime, models.COLUMN_PaymentMethod_TimesToday,
		models.COLUMN_PaymentMethod_MoneyToday)
	if err != nil {
		if err == orm.ErrNoRows {
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}
	return
}
