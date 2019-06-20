package main

import (
	"common"
	"otc/usdt"
)

func main() {

	common.IdMangerInit()

	omni := usdt.Client{usdt.Config{
		Host:     "localhost:8334",
		User:     "minerhub",
		Password: "minerhub",
	}}

	result, err := omni.Call("omni_getproperty", 2)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	common.LogFuncDebug("%v", result)
}
