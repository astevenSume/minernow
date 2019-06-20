package usdt

import (
	"common"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

const (
	BTC_MainNet = iota
	BTC_TestNet
	BTC_RegNet
)

// generate private/public/address
// https://studygolang.com/articles/13625
// https://my.oschina.net/lifephp/blog/1614964?hmsr=studygolang.com&utm_medium=studygolang.com&utm_source=studygolang.com
func GenerateKey(netType wire.BitcoinNet) (priKeyStr, address string, pubKey []byte, err error) {
	var priKey *btcec.PrivateKey
	priKey, err = btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var param *chaincfg.Params
	switch netType {
	case wire.MainNet:
		{
			param = &chaincfg.MainNetParams
		}
	case wire.TestNet, wire.TestNet3:
		{
			param = &chaincfg.TestNet3Params
		}
	case wire.SimNet:
		{
			param = &chaincfg.RegressionNetParams
		}
	default:
		common.LogFuncError("unkown net type %d", netType)
		fmt.Errorf("unkown net type %d", netType)
		return
	}

	var priKeyWif *btcutil.WIF
	priKeyWif, err = btcutil.NewWIF(priKey, param, false)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	pubKey = priKey.PubKey().SerializeUncompressed()

	var pubKeyAddr *btcutil.AddressPubKey
	pubKeyAddr, err = btcutil.NewAddressPubKey(pubKey, param)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	address = pubKeyAddr.EncodeAddress()

	priKeyStr = priKeyWif.String()

	return
}

func GenerateSegWitKey(netType wire.BitcoinNet) (priKeyStr, address string, pubKey []byte, err error) {
	var priKey *btcec.PrivateKey
	priKey, err = btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var param *chaincfg.Params
	switch netType {
	case wire.MainNet:
		{
			param = &chaincfg.MainNetParams
		}
	case wire.TestNet, wire.TestNet3:
		{
			param = &chaincfg.TestNet3Params
		}
	case wire.SimNet:
		{
			param = &chaincfg.RegressionNetParams
		}
	default:
		common.LogFuncError("unkown net type %d", netType)
		fmt.Errorf("unkown net type %d", netType)
		return
	}

	var priKeyWif *btcutil.WIF
	priKeyWif, err = btcutil.NewWIF(priKey, param, true)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	pubKey = priKeyWif.SerializePubKey()

	var (
		pubKeyAddr *btcutil.AddressScriptHash
		script     []byte
	)

	script, err = txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(btcutil.Hash160(pubKey)).Script()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	pubKeyAddr, err = btcutil.NewAddressScriptHash(script, param)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	address = pubKeyAddr.EncodeAddress()

	priKeyStr = priKeyWif.String()

	return
}

func GenerateSegWitKeyV2(netType wire.BitcoinNet) (priKeyStr, address string, pubKey []byte, err error) {
	var priKey *btcec.PrivateKey
	priKey, err = btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var param *chaincfg.Params
	switch netType {
	case wire.MainNet:
		{
			param = &chaincfg.MainNetParams
		}
	case wire.TestNet, wire.TestNet3:
		{
			param = &chaincfg.TestNet3Params
		}
	case wire.SimNet:
		{
			param = &chaincfg.RegressionNetParams
		}
	default:
		common.LogFuncError("unkown net type %d", netType)
		fmt.Errorf("unkown net type %d", netType)
		return
	}

	var priKeyWif *btcutil.WIF
	priKeyWif, err = btcutil.NewWIF(priKey, param, true)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	pubKey = priKeyWif.SerializePubKey()

	var pubKeyAddr *btcutil.AddressWitnessPubKeyHash
	pubKeyAddr, err = btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey), param)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	address = pubKeyAddr.EncodeAddress()

	priKeyStr = priKeyWif.String()

	return
}
