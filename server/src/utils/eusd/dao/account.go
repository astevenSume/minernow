package dao

import (
	"common"
	"fmt"
	"utils/eusd/models"
)

const (
	AccountStatusNoUse  = iota //未使用
	AccountStatusNormal        //正常
	AccountStatusLock          //锁定冻结
)

type AccountDao struct {
	common.BaseDao
}

func NewAccountDao(db string) *AccountDao {
	return &AccountDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *AccountDao) Create(data *models.EosAccount) (id int64, err error) {
	id, err = d.Orm.Insert(data)
	if err != nil {
		return
	}

	return
}

func (d *AccountDao) Bind(uid uint64) (account *models.EosAccount, err error) {
	if uid < 1 {
		return
	}
	account = &models.EosAccount{
		Uid: uid,
	}
	sql := fmt.Sprintf("Update %s Set %s=?,%s=?,%s=? Where %s=? Limit 1", models.TABLE_EosAccount,
		models.COLUMN_EosAccount_Uid, models.COLUMN_EosAccount_Utime, models.COLUMN_EosAccount_Status,
		models.COLUMN_EosAccount_Status)

	sqlRes, err := d.Orm.Raw(sql, uid, common.NowInt64MS(), AccountStatusNormal, AccountStatusNoUse).Exec()

	if err != nil {
		return
	}
	num, err := sqlRes.RowsAffected()
	if num == 0 {
		return
	}

	err = d.Orm.QueryTable(models.TABLE_EosAccount).Filter(models.COLUMN_EosAccount_Uid, uid).One(account)
	if err != nil {
		return
	}

	return
}

func (d *AccountDao) CountNotUse() (num int64, err error) {
	qs := d.Orm.QueryTable(models.TABLE_EosAccount).Filter(models.COLUMN_EosAccount_Status, AccountStatusNoUse)
	num, err = qs.Count()
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
	}
	return
}

// 获取type = 2 的account列表
func (d *AccountDao) GetAccount(current, limit int) (accList []*models.EosAccount, err error) {
	sql := fmt.Sprintf("SELECT account from %s WHERE %s = ? ORDER BY %s ASC LIMIT ?,?",
		models.TABLE_EosAccount, models.ATTRIBUTE_EosAccount_Status, models.ATTRIBUTE_EosAccount_Id)
	_, err = d.Orm.Raw(sql, AccountStatusNormal, current, limit).QueryRows(&accList)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//GetCount 获取总记录数
func (d *AccountDao) GetCount() (total int, err error) {

	sqlTotal := fmt.Sprintf("select count(*) from %s", models.TABLE_EosAccount)
	err = d.Orm.Raw(sqlTotal).QueryRow(&total)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}
