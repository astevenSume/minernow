package main

import (
	"common"
	"otc/usdt"
)

func main() {
	//resp, err := usdt.NewExplorer().Balances([]string{"1GPQ4yPjL3xUEae4wWgwsqwJpjrZHkZRhv"})
	//resp, err := usdt.NewExplorer().AddrDetails("1GPQ4yPjL3xUEae4wWgwsqwJpjrZHkZRhv")
	//if err == nil {
	//	common.LogFuncDebug("%v", resp)
	//}
	priviteKey := "5JPzYaNNyrUKTo7pZ5rMH1grWHkhSYSD6ieyXPuXGtSKbMZSD2S"
	toAddr := "1DXDso8dNotaWxKu9MDt48ZNNgqgxJ1ggC"
	txHash := "7d2f02bee6792903529d821fecc73e9b76bb871de566f1515469ba7b504ad4b2"
	transaction, err := usdt.CreateTransaction(priviteKey, toAddr, 546, txHash)
	if err != nil {
		return
	}

	common.LogFuncDebug("%+v", transaction)
}
