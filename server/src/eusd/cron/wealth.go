package cron

import (
	"common"
	"eusd/eosplus"
	"time"
	"utils/eusd/dao"
)

// 先获取account 然后获取链上余额 最后更新db
func EosWealthFetchAndUpdate() (err error) {

	limit := 10000
	// 获取count 分成n份
	total, err := dao.WealthDaoEntity.GetCount()
	if err != nil {
		common.LogFuncError("BalanceSync ERR:%v", err)
		return
	}
	if total <= 0 {
		return
	}
	total += limit

	for i := 0; i <= total; i += limit {
		list, err := dao.WealthDaoEntity.GetAccountForLimit(i, limit)
		if err != nil {
			common.LogFuncError("EosWealthFetchAndUpdate GetAccountForLimit ERR:%v", err)
			return err
		}
		if len(list) <= 0 {
			common.LogFuncError("EosWealthFetchAndUpdate len(list) <= 0 ERR:%v", err)
			return err
		}
		for _, s := range list {
			account := s.Account
			res := eosplus.EosPlusAPI.Rpc.GetBalance(account)

			if len(res) > 0 {
				balance := int64(res[0].Amount)
				err := dao.WealthDaoEntity.UpdateForBalanceSync(s.Uid, balance)

				if err != nil {
					common.LogFuncError("FetchAndUpdate UpdateForBalanceSync ERR:%v", err)
					continue
				}
			} else {
				common.LogFuncError("FetchAndUpdate len(res) < 0, this account have not EUSD")
			}
		}
	}
	time.Sleep(time.Second * 2)
	return
}
