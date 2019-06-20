package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/usdt/models"
)

const (
	TransactionTypeUnknown uint8 = iota
	TransactionTypeTransferin
	TransactionTypeTransferout
	TransactionTypeCollection
)

const (
	TransactionStatusUnknown uint32 = iota
	TransactionStatusTransferred
	TransactionStatusConfirmed
	TransactionStatusFailure
	TransactionStatusRevoke
)

type TransactionDao struct {
	common.BaseDao
}

func NewTransactionDao(db string) *TransactionDao {
	return &TransactionDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var TransactionDaoEntity *TransactionDao

func (d *TransactionDao) Create(transaction *models.UsdtTransaction) (id int64, err error) {
	id, err = d.Orm.Insert(transaction)
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
	Status    string `json:"status"`
	Ctime     int64  `json:"ctime"`
	Utime     int64  `json:"utime"`
}

// 获取对应地址充值记录
func (d *TransactionDao) QueryAllRecharge(uaid uint64) (list []models.UsdtTransaction, rechargeInterger int64, err error) {
	list = []models.UsdtTransaction{}
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM `%s` WHERE `%s`=?",
		models.TABLE_UsdtTransaction,
		models.ATTRIBUTE_UsdtTransaction_Uaid,
	), uaid).QueryRows(&list)

	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
		}
		return
	}

	// 计算充值总额
	if len(list) > 0 {
		for _, v := range list {
			if v.Type == TransactionTypeTransferin {
				rechargeInterger += v.AmountInteger
			} else {
				rechargeInterger -= v.AmountInteger
			}
		}
	}

	return
}

// 批量插入充值记录
func (d *TransactionDao) InsertMulti(txs []models.UsdtTransaction) (err error) {
	_, err = d.Orm.InsertMulti(common.BulkCount, txs)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}
