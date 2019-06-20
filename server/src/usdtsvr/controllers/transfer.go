package controllers

import (
	"usdt"
	usdtdao "utils/usdt/dao"
	usdtmodels "utils/usdt/models"
)

func ProcessApprovedTransfer() {
	// query approved transfer count
	total, err := usdtdao.WealthLogDaoEntity.QueryApprovedTotal()
	if err != nil || total <= 0 {
		return
	}

	//
	perPage := 100
	pages := total / perPage

	if total%perPage > 0 {
		pages += 1
	}

	var list []usdtmodels.UsdtWealthLog
	for i := 1; i <= pages; i++ {
		// query approved logs
		list, err = usdtdao.WealthLogDaoEntity.QueryApproved(i, perPage)
		if err != nil {
			return
		}

		for _, l := range list {
			// todo maybe sometimes transfer by system
			usdt.TransferByUser(l.Uid, l.Id, l.To, l.AmountInteger, "")
		}
	}
}
