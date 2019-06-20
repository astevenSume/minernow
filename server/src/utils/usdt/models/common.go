package models

import "github.com/astaxie/beego/orm"

func ModelsInit() (err error) {
	orm.RegisterModel(
		new(MarketPrices),
		new(Prices),
		new(UsdtAccount),
		new(UsdtOnchainBalance),
		new(UsdtOnChainData),
		new(UsdtOnchainLog),
		new(UsdtOnChainSyncPos),
		new(UsdtOnchainTransaction),
		new(PriKey),
		new(UsdtSweepLog),
		new(UsdtTransaction),
		new(UsdtWealthLog),
	)

	return
}
