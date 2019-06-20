package eosplus

import (
	"common"
	. "otc_error"
	"utils/eusd/dao"
	"utils/eusd/models"
)

type Account struct {
}

func AccountRpc() *Account {
	return &Account{}
}

func (a *Account) Bind(uid uint64) (acc *models.EosAccount, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	acc, err := dao.AccountDaoEntity.Bind(uid)
	if err != nil {
		errCode = ERROR_CODE_BIND_WALLET_ERROR
		common.LogFuncError("Bind Wallet Account ERR:%v", err)
		return
	}
	if acc.Account == "" {
		errCode = ERROR_CODE_BIND_WALLET_ERROR
		common.LogFuncError("Bind Wallet Account ERR:NOT ACCOUNT")
		return
	}

	return
}
