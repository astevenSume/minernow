package main

import (
	"encoding/json"
	"eusd/eosplus"
	"fmt"
	"github.com/astaxie/beego/orm"
)

func main() {
	//url := "https://proxy.eosnode.tools"

	//主网
	url := "https://api.eosnewyork.io"

	//EOS测试网的另一个url
	//url := "http://88.99.193.44:8888"

	//url := "https://api.kylin.alohaeos.com"
	walletUrl := ""

	//err := beego.LoadAppConfig("ini", "/code/go/wallet/server/src/otc/conf/app.conf")
	//err := beego.LoadAppConfig("ini", "D:/code/server/src/otc/conf/app.conf")

	//if err != nil {
	//	common.LogFuncError("%v", err)
	//}

	orm.Debug = true

	//models.ModelsInit()
	//A2, _ := eosplus.ApiInit(url, walletUrl)
	//fmt.Println(A2)

	//modelsInit()
	//common.DbInit()
	//common.RedisInit()
	//
	//eosplus.DaoInit()
	//daoInit()

	api := new(eosplus.RpcKeyBag)
	api.InitApi(url, walletUrl)
	api.SetDebug()
	// winpaytoken2 信息
	//api.TokenAccount = "winpaytoken2"
	//api.KeyBag.ImportPrivateKey("5K9RHHfGGqSGVWGuH3oYZgrvF3MY83fhhK5tN4ESKeQPWqAtJcn")

	//hyf 信息
	api.TokenAccount = "xxx"
	api.KeyBag.ImportPrivateKey("xxx")

	api.Api.SetSigner(api.KeyBag)
	api.SetTransferAble()
	//res, err := api.GetInfo()
	//fmt.Println(res)
	//res1 := api.GetAccount("winpaytoken1")
	//fmt.Println(res1)
	//api.CreateKey()

	//res2 := api.GetBalance("fows42pj3w5l")
	//fmt.Println(res2)
	//res5 := api.GetBalance("winpaytoken2")
	//fmt.Println(res5)
	//res3 := api.GetBalanceEos("winpaytoken2")
	//fmt.Println(res3)

	//api.IssueToken("winpaytoken2","winpaytoken2", "10000.0000 EUSD", "winpaytoken2", "memo") // 发行token
	//api.IssueToken("winpaytoken2","winpaytoken2", "10000.0000 EUSD") // 发行token
	//api.RetireToken("winpaytoken2", "10000.0000 EUSD") //销毁token
	//res4 := api.GetBalance("fows42pj3w5l")
	//fmt.Println(res4)

	//api.IssueToken("winpaytoken2", "1.0000 EUSD")
	//api.RetireToken("1.0000 EUSD", "")
	//res5 := api.GetBalanceEos("hyfeos111111")
	//fmt.Println(res5)

	//充值储备金 体现储备金
	//res, t, errCode := api.Deposit("0.0001 EOS")
	//res, t, errCode := api.Withdraw("0.0001 EOS")

	//购买cpu，net
	//res, t, errCode := api.Rentcpu("hyfeos111111", "0.0001 EOS", "0.0000 EOS")
	//res, t, errCode := api.Rentnet("hyfeos111111", "0.0001 EOS", "0.0000 EOS")
	//fmt.Println(res)
	//fmt.Println(t)
	//fmt.Println(errCode)

	//res_tr, err := api.GetTransferTest("bee3bff85bb7fc953372efbbe27a053540121946789257e0b5a7078769a1b1ad")
	//fmt.Println(res_tr)

	//test_transferRaw(api) //有点小问题
	//test_newAccount(api)
	//test_delegate(api)
	//test_undelegate(api)
	//test_buyRam(api)

	//批量新增链上账户 ok
	//Account := eosplus.Account{}
	//Account.SetApi(api)
	//Account.CronNewAccount(1)

	//acc, errCode, err := Account.Bind(1)
	//fmt.Printf("\n\n acc:%v\n\n \n\n errCode:%v\n\n err:%v\n\n", acc, errCode)

	//errCode := A2.Transaction.TransferByUids(1, 2, 10000)
	//_, errCode := A2.Wealth.Create(1, true)
	//errCode := A2.Wealth.BecomeExchanger(1)
	//fmt.Printf("\n\n  errCode:%v\n\n", errCode)
	//for i := 0; i < 1000; i++ {
	//	test_icu(api)
	//}
	//api.TransferCheckResource("nineninesix1", "nineninesix1")

	//account := "123"
	//api.ResourcesAccount = eos.AccountName(account)
	//eosplus.EosPlusAPI.Rpc.ImportKeys(account)
	//eosplus.EosPlusAPI.Rpc.RexRentCpu(account, "0.0001 EOS")
	//return

}

//转账记录
func test_transfer(api *eosplus.RpcKeyBag) {
	from := "winpaytoken2"
	to := "fows42pj3w5l"
	//from := "ihyfwinpay13"
	//to := "ihyfwinpay12"
	//fmt.Printf("\n\n eosio:%v", eosplus.ToJsonIndent(api.GetAccount(from)))
	//response, errCode, err := api.Transfer("eosio.token", from, to, "10.1000 EOS", "memo")
	response, errCode, err := api.Transfer("winpaytoken2", from, to, "2.0000 EUSD", "memo")

	fmt.Printf("\n\n errCode:%v\n\n", errCode)
	fmt.Printf("\n\n err:%v\n\n", err)
	fmt.Printf("\n\n response:%v\n\n", eosplus.ToJsonIndent(response))
	//fmt.Printf("\n\n eosio:%v", eosplus.ToJsonIndent(api.GetAccount(from)))

}

//转账记录
func test_transferRaw(api *eosplus.RpcKeyBag) {
	transpId := "83524ebb5b15f113dd6a46fe9ab34c66ff3159713da4b76af1a8a12504667991"
	blockNum := 23407041
	res, err := api.GetTransfer(transpId, blockNum)
	fmt.Printf("\n\n transp json:\n%v   \nerr:%v\n", eosplus.ToJsonIndent(res), err)
}

//创建新账号
func test_newAccount(api *eosplus.RpcKeyBag) {
	//owner := "EOS5XPWBgu4w5nNyEENRnjWbek9ZsThURnHLSDBjGcVn6kB3TUk3a"
	//active := "EOS53VMeWJJpMT61WpGHRpFtHhka9diRK9UjMZQmgEwayLeboAoha"
	owner := "EOS71PnvbaRF5T6wQ614YuUGETaZKYYUr1bj3RsxLFAkt9nn8tviq"
	active := "EOS71PnvbaRF5T6wQ614YuUGETaZKYYUr1bj3RsxLFAkt9nn8tviq"
	res, errCode, err := api.AccountCreate("hyfeos111111", "nineninesix1", owner, active)
	//res, errCode, err := api.AccountCreate("ihyfwinpay12", "ihyfiiizyhm1", owner, active)

	content, _ := json.MarshalIndent(res, "", "  ")
	fmt.Printf("errCode:%v\n  err:%v\n res:%v\n", errCode, err, string(content))
}

func test_delegate(api *eosplus.RpcKeyBag) {
	from := "winpaytoken2"
	to := "fows42pj3w5l"

	//fmt.Printf("\n\n eosio:%v", eosplus.ToJsonIndent(api.GetAccount(to)))
	res, errCode, err := api.DelegateBw(from, to, "2.1000 EOS", "2.1000 EOS")
	fmt.Printf("\n\n errCode:%v \n\nerr:%v\n \n\n res:\n%v", errCode, err, eosplus.ToJsonIndent(res))
	//fmt.Printf("\n\n eosio:%v", eosplus.ToJsonIndent(api.GetAccount(to)))
}

func test_undelegate(api *eosplus.RpcKeyBag) {
	from := "winpaytoken2"
	//to := "qqqwwweee222"
	to := "fows42pj3w5l"

	//fmt.Printf("\n\n eosio:%v", eosplus.ToJsonIndent(api.GetAccount(to)))
	res, errCode, err := api.UnDelegateBw(from, to, "0.5000 EOS", "1.0000 EOS")
	fmt.Printf("\n\n errCode:%v \n\nerr:%v\n \n\n res:\n%v", errCode, err, eosplus.ToJsonIndent(res))
	//fmt.Printf("\n\n eosio:%v", eosplus.ToJsonIndent(api.GetAccount(to)))
}

func test_buyRam(api *eosplus.RpcKeyBag) {
	from := "winpaytoken2"
	res, errCode, err := api.BuyRam(from, "fows42pj3w5l", "5.0000 EOS")
	fmt.Printf("\n\n errCode:%v \n\nerr:%v\n \n\n res:\n%v", errCode, err, eosplus.ToJsonIndent(res))
}

func test_icu(api *eosplus.RpcKeyBag) {
	api.CheckResource("fows42pj3w5l")
	//response, errCode, err := api.Transfer("hyfeos111111", "nineninesix1", "hyfeos111111", "10.0000 ICU", "memo")
	//fmt.Println(response, errCode, err)
}
