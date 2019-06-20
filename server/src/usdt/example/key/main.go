package main

import (
	"fmt"
	"usdt"
)

func main() {
	priKey, addr, pubKey, err := usdt.GenerateKey(usdt.BTC_MainNet)
	//priKey, addr, pubKey, err := usdt.GenerateKey(usdt.BTC_TestNet)
	//priKey, addr, pubKey, err := usdt.GenerateKey(usdt.BTC_RegNet)
	fmt.Printf("%+v (size %d)\n%+v size(%d)\n%+v (size %d)\n%+v\n", priKey, len(priKey), pubKey, len(pubKey), addr, len(addr), err)
	//time.Sleep(time.Second * 10)
}
