package dao

import (
	"common"
	"fmt"
	"utils/admin/dao"
	common2 "utils/common"
	"utils/usdt/models"

	"github.com/astaxie/beego/orm"
)

type OnchainTransactionDao struct {
	common.BaseDao
}

func NewOnchainTransactionDao(db string) *OnchainTransactionDao {
	return &OnchainTransactionDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var OnchainTransactionDaoEntity *OnchainTransactionDao

func (d OnchainTransactionDao) UpdateMulti(uid uint64, txs []models.UsdtOnchainTransaction, addr string) (err error) {
	tableName := models.TABLE_UsdtOnchainTransaction
	for _, tx := range txs {
		_, err = d.Orm.Raw(fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) "+
			"ON DUPLICATE KEY UPDATE %s=?,  %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?",
			tableName,
			models.COLUMN_UsdtOnchainTransaction_TxId,
			models.COLUMN_UsdtOnchainTransaction_Uaid,
			models.COLUMN_UsdtOnchainTransaction_Type,
			models.COLUMN_UsdtOnchainTransaction_PropertyId,
			models.COLUMN_UsdtOnchainTransaction_PropertyName,
			models.COLUMN_UsdtOnchainTransaction_TxType,
			models.COLUMN_UsdtOnchainTransaction_TxTypeInt,
			models.COLUMN_UsdtOnchainTransaction_AmountInteger,
			models.COLUMN_UsdtOnchainTransaction_Block,
			models.COLUMN_UsdtOnchainTransaction_BlockHash,
			models.COLUMN_UsdtOnchainTransaction_BlockTime,
			models.COLUMN_UsdtOnchainTransaction_Confirmations,
			models.COLUMN_UsdtOnchainTransaction_Divisible,
			models.COLUMN_UsdtOnchainTransaction_FeeAmountInteger,
			models.COLUMN_UsdtOnchainTransaction_IsMine,
			models.COLUMN_UsdtOnchainTransaction_PositionInBlock,
			models.COLUMN_UsdtOnchainTransaction_ReferenceAddress,
			models.COLUMN_UsdtOnchainTransaction_SendingAddress,
			models.COLUMN_UsdtOnchainTransaction_Version,
			models.COLUMN_UsdtOnchainTransaction_Mtime,
			models.COLUMN_UsdtOnchainTransaction_TxId,
			models.COLUMN_UsdtOnchainTransaction_Uaid,
			models.COLUMN_UsdtOnchainTransaction_Type,
			models.COLUMN_UsdtOnchainTransaction_PropertyId,
			models.COLUMN_UsdtOnchainTransaction_PropertyName,
			models.COLUMN_UsdtOnchainTransaction_TxType,
			models.COLUMN_UsdtOnchainTransaction_TxTypeInt,
			models.COLUMN_UsdtOnchainTransaction_AmountInteger,
			models.COLUMN_UsdtOnchainTransaction_Block,
			models.COLUMN_UsdtOnchainTransaction_BlockHash,
			models.COLUMN_UsdtOnchainTransaction_BlockTime,
			models.COLUMN_UsdtOnchainTransaction_Confirmations,
			models.COLUMN_UsdtOnchainTransaction_Divisible,
			models.COLUMN_UsdtOnchainTransaction_FeeAmountInteger,
			models.COLUMN_UsdtOnchainTransaction_IsMine,
			models.COLUMN_UsdtOnchainTransaction_PositionInBlock,
			models.COLUMN_UsdtOnchainTransaction_ReferenceAddress,
			models.COLUMN_UsdtOnchainTransaction_SendingAddress,
			models.COLUMN_UsdtOnchainTransaction_Version,
			models.COLUMN_UsdtOnchainTransaction_Mtime),
			tx.TxId, tx.Uaid, tx.Type, tx.PropertyId, tx.PropertyName, tx.TxType, tx.TxTypeInt, tx.AmountInteger, tx.Block, tx.BlockHash, tx.BlockTime, tx.Confirmations,
			tx.Divisible, tx.FeeAmountInteger, tx.IsMine, tx.PositionInBlock, tx.ReferenceAddress, tx.SendingAddress, tx.Version, tx.Mtime,
			tx.TxId, tx.Uaid, tx.Type, tx.PropertyId, tx.PropertyName, tx.TxType, tx.TxTypeInt, tx.AmountInteger, tx.Block, tx.BlockHash, tx.BlockTime, tx.Confirmations,
			tx.Divisible, tx.FeeAmountInteger, tx.IsMine, tx.PositionInBlock, tx.ReferenceAddress, tx.SendingAddress, tx.Version, tx.Mtime).
			Exec()
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

func (d *OnchainTransactionDao) QueryByTxid(txid string) (transaction *models.UsdtOnchainTransaction, err error) {
	transaction = &models.UsdtOnchainTransaction{
		TxId: txid,
	}

	err = d.Orm.Read(transaction)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
		return
	}
	return
}

// 获取地址对应的链上充值记录
func (d *OnchainTransactionDao) QueryOnchainAllRecharge(uaid uint64, confirmsLimit uint32) (listOnChain, listSweep, listWealth []models.UsdtOnchainTransaction, rechargeInteger, rechargeIntegerTransferFrozen int64, err error) {
	//查询充值记录
	listIn := []models.UsdtOnchainTransaction{}
	listIn, err = d.queryOnchainTxByType(uaid, confirmsLimit, TransactionTypeTransferin)
	if err != nil {
		return
	}

	//查询转账记录
	listOut := []models.UsdtOnchainTransaction{}
	listOut, err = d.queryOnchainTxByType(uaid, confirmsLimit, TransactionTypeTransferout)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
		} else {
			return
		}
	}

	//计算充值总额
	//充值加和
	if len(listIn) > 0 {
		for _, v := range listIn {
			rechargeInteger += v.AmountInteger
		}
	}

	listSweep = make([]models.UsdtOnchainTransaction, 0, len(listOut))
	listWealth = make([]models.UsdtOnchainTransaction, 0, len(listOut))
	//扣除转账
	if len(listOut) > 0 {
		for _, v := range listOut {

			var (
				isExist bool
			)
			// db 错误或者存在归集记录则跳过
			if isExist, err = SweepLogDaoEntity.IsExistTransferred(v.TxId); err != nil {
				return
			}
			if isExist {
				listSweep = append(listSweep, v)
				continue
			}

			if isExist, err = WealthLogDaoEntity.IsExistTransfered(v.TxId); err != nil { // 校验交易是否存在
				return
			}

			if isExist { // 存在，则扣除transfer_integer
				rechargeIntegerTransferFrozen -= v.AmountInteger
				listWealth = append(listSweep, v)
			} else { // 不存在，则扣除available
				rechargeInteger -= v.AmountInteger
				common.LogFuncError("同步用户链上数据出现系统内不存在的转出记录 \t uaid:%v \t txID:%v", v.Uaid, v.TxId)
				common2.SmsWarning(dao.ConfigWarningTypeUsdtPriKeyLeakage, map[string]string{
					"UAID": fmt.Sprint(v.Uaid),
				})
			}

		}
	}

	// 拼接交易记录
	listOnChain = append(listIn, listOut...)

	return
}

func (d *OnchainTransactionDao) queryOnchainTxByType(uaid uint64, confirmsLimit uint32, t uint8) (list []models.UsdtOnchainTransaction, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM `%s` WHERE `%s`=? AND `%s`=? AND `%s`>=?",
		models.TABLE_UsdtOnchainTransaction,
		models.COLUMN_UsdtOnchainTransaction_Uaid,
		models.COLUMN_UsdtOnchainTransaction_Type,
		models.COLUMN_UsdtOnchainTransaction_Confirmations,
	), uaid, t, confirmsLimit).QueryRows(&list)

	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
		}
		return
	}

	return
}
