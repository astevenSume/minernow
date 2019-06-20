package dao

import (
	"common"
	"fmt"
	"utils/usdt/models"
)

type OnchainBalanceDao struct {
	common.BaseDao
}

func NewOnchainBalanceDao(db string) *OnchainBalanceDao {
	return &OnchainBalanceDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var OnchainBalanceDaoEntity *OnchainBalanceDao

func (d OnchainBalanceDao) UpdateMulti(uid uint64, balances []models.UsdtOnchainBalance) (err error) {
	tableName := OnchainTableMgr.GetOnChainDataTableName(uid, models.TABLE_UsdtOnchainBalance)
	for _, b := range balances {
		_, err = d.Orm.Raw(fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?, %s=?",
			tableName,
			models.COLUMN_UsdtOnchainBalance_Address,
			models.COLUMN_UsdtOnchainBalance_PropertyId,
			models.COLUMN_UsdtOnchainBalance_PendingPos,
			models.COLUMN_UsdtOnchainBalance_Reserved,
			models.COLUMN_UsdtOnchainBalance_Divisible,
			models.COLUMN_UsdtOnchainBalance_AmountInteger,
			models.COLUMN_UsdtOnchainBalance_Frozen,
			models.COLUMN_UsdtOnchainBalance_PendingNeg,
			models.COLUMN_UsdtOnchainBalance_Mtime,
			models.COLUMN_UsdtOnchainBalance_Address,
			models.COLUMN_UsdtOnchainBalance_PropertyId,
			models.COLUMN_UsdtOnchainBalance_PendingPos,
			models.COLUMN_UsdtOnchainBalance_Reserved,
			models.COLUMN_UsdtOnchainBalance_Divisible,
			models.COLUMN_UsdtOnchainBalance_AmountInteger,
			models.COLUMN_UsdtOnchainBalance_Frozen,
			models.COLUMN_UsdtOnchainBalance_PendingNeg,
			models.COLUMN_UsdtOnchainBalance_Mtime),
			b.Address, b.PropertyId, b.PendingPos, b.Reserved, b.Divisible, b.AmountInteger, b.Frozen, b.PendingNeg, b.Mtime,
			b.Address, b.PropertyId, b.PendingPos, b.Reserved, b.Divisible, b.AmountInteger, b.Frozen, b.PendingNeg, b.Mtime).
			Exec()
		//_, err = DbOrm.InsertOrUpdate(&b)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
	}
	return
}
