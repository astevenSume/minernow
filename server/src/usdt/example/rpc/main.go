package main

import (
	"common"
	"otc/usdt"
)

func main() {

	common.IdManagerInit(common.IdTypeMax, 0, 0)

	usdt.RpcConfig.Host = "localhost:8334"
	usdt.RpcConfig.User = "minerhub"
	usdt.RpcConfig.Password = "minerhub"

	//user1:pri:936ykxjirXUdAsaiiF6FAtciwLrfzSEqm15LFthuaCHP7jBqjVs, addr:mzGztTPXrsS1gCSx6H5E6zHa8GDFbd3mKs
	//user2:pri:92gD8HNBi5XNbDCe7PihGPKPGLALVUFJnV8Qk7B2ndvHxBffdMU, addr:n1Tmq7C6PWHnsxGHPLUj5xZtkU1JgwEfrm

	//result, err := usdt.GetBalance("n1Tmq7C6PWHnsxGHPLUj5xZtkU1JgwEfrm", usdt.BTC_PROPERTY_ID_TEST_OMNI)
	//result, err := usdt.ListProperties()
	result, err := usdt.Generate(10)
	if err != nil {
		return
	}

	common.LogFuncDebug("%v", result)
}
