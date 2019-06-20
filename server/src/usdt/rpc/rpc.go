package rpc

// Package usdt provide a USDT interface

const (
	BTC_PROPERTY_ID_OMNI      = 1
	BTC_PROPERTY_ID_TEST_OMNI = 2
	BTC_PROPERTY_ID_USDT      = 31
)

// generate blocks, for regtest
func Generate(num int) (result []string, err error) {
	err = NewClient().Call(&result, "generate", num)
	return
}

type Property struct {
	PropertyId int    `json:"propertyid"`
	Name       string `json:"name"`
	Divisible  bool   `json:"divisible"`
}

// list properties
func ListProperties() (result []Property, err error) {
	err = NewClient().Call(&result, "omni_listproperties")
	return
}

type Balance struct {
	Balance  string `json:"balance"`
	Reserved string `json:"reserved"`
	Frozen   string `json:"frozen"`
}

func GetBalance(addr string, propertyId int) (result Balance, err error) {
	err = NewClient().Call(&result, "omni_getbalance", addr, propertyId)
	return
}

func GetUsdtBalance(addr string) (result Balance, err error) {
	return GetBalance(addr, BTC_PROPERTY_ID_USDT)
}
