package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/eusd/models"
)

type TransactionDao struct {
	common.BaseDao
}

const (
	TransactionFinish = 1 //"executed"
	TransactionFail   = 2 //"failed"
	TransactionCancel = 3 //"Cancel"
)

func NewTransactionDao(db string) *TransactionDao {
	return &TransactionDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *TransactionDao) Create(transaction *models.EosTransaction) (id int64, err error) {
	id, err = d.Orm.Insert(transaction)
	if err != nil {
		return
	}

	return
}

func (d *TransactionDao) CountByUid(transaction *models.EosTransaction) (id int64, err error) {
	if err != nil {
		return
	}

	return
}

type Record struct {
	Id        uint64 `json:"id"`
	Type      uint8  `json:"type"`
	Address   string `json:"address"`
	Address2  string `json:"address2"`
	Amount    string `json:"amount"`
	Memo      string `json:"memo"`
	Txid      string `json:"txid"`
	Block_num string `json:"block_num"`
	Status    int8   `json:"status"`
	Ctime     int64  `json:"ctime"`
	Utime     int64  `json:"utime"`
}

//get all eos transaction by page and limit dao
func (d *TransactionDao) GetEosTransactions(page int, limit int, rType uint8) (total int, RecordList []*Record, err error) {
	var cond string
	var param []interface{}
	if rType > 0 {
		cond = " WHERE " + models.COLUMN_EosTransaction_Type + " = ?"
		param = append(param, rType)
	}
	sqlTotal := fmt.Sprintf("select count(*) from %s %s", models.TABLE_EosTransaction, cond)
	err = d.Orm.Raw(sqlTotal, param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	param = append(param, (page-1)*limit)
	param = append(param, limit)
	sql := fmt.Sprintf("SELECT id, type, payer AS address, receiver AS address2, quantity AS amount, memo, "+
		"transaction_id AS txid, block_num, `status`, ctime, utime FROM eos_transaction %s ORDER BY ctime DESC LIMIT ?,?", cond)
	_, err = d.Orm.Raw(sql, param...).QueryRows(&RecordList)
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

//get eos transactions by id
func (d *TransactionDao) GetEosTransactionById(Id uint64) (record *Record, err error) {
	if Id <= 0 {
		common.LogFuncCritical("query EosTransaction id <= 0")
		return
	}
	sql := fmt.Sprintf("SELECT id, type, payer AS address, receiver AS address2, quantity AS amount, memo, transaction_id AS txid, block_num, `status`, ctime, utime FROM eos_transaction WHERE id = ?")
	err = d.Orm.Raw(sql, Id).QueryRow(&record)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//get eos transactions by payer or receiver
func (d *TransactionDao) GetEosTransactionsByAccount(account string, page int, limit int) (RecordList []*Record, total int, err error) {
	if account == "" {
		common.LogFuncCritical("query EosTransaction account is null")
		return
	}
	sql := fmt.Sprintf("SELECT id, type, payer AS address, receiver AS address2, quantity AS amount, memo, transaction_id AS txid, block_num, `status`, ctime, utime FROM eos_transaction WHERE payer = ? OR receiver = ? ORDER BY ctime DESC LIMIT ?,?")
	_, err = d.Orm.Raw(sql, account, account, (page-1)*limit, limit).QueryRows(&RecordList)

	sqlTotal := fmt.Sprintf("select count(*) from eos_transaction where payer = ?")
	err = d.Orm.Raw(sqlTotal, account).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

func (d *TransactionDao) Info(id uint64) (data *models.EosTransaction) {
	data = &models.EosTransaction{}
	if id <= 0 {
		return
	}
	data = &models.EosTransaction{
		Id: id,
	}
	err := d.Orm.Read(data)
	if err != nil {
		if err == orm.ErrNoRows {
			data.Id = 0
			common.LogFuncError("DBERR：%v", data.Id)
			return
		}
		common.LogFuncError("DBERR：%v", err)
		return
	}
	return
}

func (d *TransactionDao) Finish(id uint64, blockNum uint32) (data *models.EosTransaction) {
	data = &models.EosTransaction{
		Id:       id,
		Status:   TransactionFinish,
		BlockNum: blockNum,
		Utime:    common.NowInt64MS(),
	}
	_, err := d.Orm.Update(data, models.COLUMN_EosTransaction_Status, models.COLUMN_EosTransaction_Utime)
	if err != nil {
		common.LogFuncError("DBERR：%v", err)
		data.Id = 0
		return
	}
	return
}
