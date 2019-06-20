package cron

import (
	"common"
	"eusd/eosplus"
	. "otc_error"
	"time"
	"utils/eusd/dao"
)

const (
	limit    = 10000       // 每次处理条数
	interval = 60 * 60 * 6 // 6小时
)

// todo 待整体重新修改
func CheckResourceAction() (err error) {
	common.GoSafeRun(StartWork)
	return
}

//开始
func StartWork() {
	for true {
		total, err := dao.AccountDaoEntity.GetCount()
		if err != nil {
			common.LogFuncError("BalanceSync ERR:%v", err)
			return
		}
		if total > 0 {
			times := total/limit + 1
			go common.SafeRun(func() {
				CheckResource(limit, times)
			})()
		}
		time.Sleep(time.Second * interval)
	}
}

//检查资源
func CheckResource(limit, times int) {
	for i := 1; i <= times; i++ {
		acclist, err := dao.AccountDaoEntity.GetAccount((i-1)*limit, limit)
		if err != nil {
			common.LogFuncError("CheckResourceAction CheckResource ERR:%v", err)
			return
		}
		if len(acclist) < 0 {
			return
		}
		for _, acc := range acclist {
			account := acc.Account
			// todo 多了赎回
			//检查资源，不够会买入
			errCode := eosplus.EosPlusAPI.Rpc.CheckResource(account)

			if errCode != ERROR_CODE_SUCCESS {
				common.LogFuncError("FetchAndUpdate GetBalanceByName ERR:%v", errCode)
				continue
			}
		}
	}
}
