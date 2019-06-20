package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/eos/models"
)

type EosAccountKeysDao struct {
	common.BaseDao
}

func NewEosAccountKeysDao(db string) *EosAccountKeysDao {
	return &EosAccountKeysDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *EosAccountKeysDao) Create(data *models.EosAccountKeys) (id int64, err error) {
	id, err = d.Orm.Insert(data)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

func (d *EosAccountKeysDao) Info(account string) (data *models.EosAccountKeys, err error) {
	sql := fmt.Sprintf("Select * from %s where %s=?", models.TABLE_EosAccountKeys, models.COLUMN_EosAccountKeys_Account)
	data = &models.EosAccountKeys{}
	err = d.Orm.Raw(sql, account).QueryRow(data)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			data.Account = ""
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}
