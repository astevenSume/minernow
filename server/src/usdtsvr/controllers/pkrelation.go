package controllers

import (
	"common"
	"usdt"
	usdtdao "utils/usdt/dao"
	usdtmodels "utils/usdt/models"

	"github.com/btcsuite/btcd/wire"
)

// process pk unrelation
func ProcessPkUnrelated() {
	//query unrelated usdt pk relation
	total, err := usdtdao.AccountDaoEntity.QueryTotalUnrelated()
	if err != nil {
		return
	}

	//common.LogFuncDebug("there are %d records pk unrelated.", total)

	if total <= 0 {
		return
	}

	var (
		perPage = 100
	)

	pages := (total / perPage)
	if total%perPage > 0 {
		pages += 1
	}

	now := common.NowInt64MS()
	for i := 1; i <= pages; i++ {
		list, err := usdtdao.AccountDaoEntity.QueryUnrelated(i, perPage)
		if err != nil {
			return
		}

		for _, l := range list {
			processSinglePkRelation(l, now)
		}
	}
}

func processSinglePkRelation(l usdtmodels.UsdtAccount, now int64) {
	// generate pk
	pri, addr, _, err := usdt.GenerateSegWitKey(wire.MainNet)
	if err != nil {
		return
	}

	id, err := common.IdManagerGen(usdtdao.IdTypePriKey)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	// encrypt to base64
	var priDb string
	priDb, err = common.EncryptToBase64(pri, usdt.PriAesKey)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	// save to db
	err = usdtdao.PriKeyDaoEntity.Add(id, priDb, addr)
	if err != nil {
		return
	}

	// relate pk to usdt account
	l.Pkid = id
	l.Mtime = now
	l.Address = addr
	err = usdtdao.AccountDaoEntity.UpdateRelation(l)
	if err != nil {
		return
	}
}
