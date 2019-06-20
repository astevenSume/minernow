package main

import (
	"bytes"
	"common"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"usdt"
	"usdt/explorer"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

var (
	CONFIG_PATH                                      string
	platformPriKey, platformAddress, priKey, address string
	confirmationLimit, utxoQuantity                  int
	satoshiLimit                                     float64
	cfg                                              config
)

type config struct {
	Confirmation int    `json:"confirmation"`
	PriKey       string `json:"prikey"`
	Address      string `json:"address"`
	Data         []struct {
		PriKey   string  `json:"prikey"`
		Address  string  `json:"address"`
		Quantity int     `json:"quantity"`
		Satoshi  float64 `json:"satoshi"`
	} `json:"data"`
}

func main() {

	flag.StringVar(&CONFIG_PATH, "config", "./data.json", "directory of config file")
	flag.Parse()
	log.Flags()

	var (
		data []byte
		err  error
	)
	if data, err = ioutil.ReadFile(CONFIG_PATH); err != nil {
		log.Fatalf("read file %s failed : %v", CONFIG_PATH, err)
		return
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("json.Unmarshal failed : %v", err)
		return
	}

	confirmationLimit = cfg.Confirmation
	platformPriKey = cfg.PriKey
	platformAddress = cfg.Address

	for _, d := range cfg.Data {
		err = splitUtxo(d.PriKey, d.Address, d.Quantity, d.Satoshi)
		if err != nil {
			fmt.Println(err)
		}
	}

}

const TxVersion = 2

func splitUtxo(priKey, address string, utxoQuantity int, satoshiLimit float64) (err error) {
	priKey, err = common.DecryptFromBase64(priKey, usdt.PriAesKey)
	platformPriKey, err = common.DecryptFromBase64(platformPriKey, usdt.PriAesKey)
	if err != nil {
		return
	}

	var (
		utxos, useUtxos, payerUtxos, payerUseUtxos []usdt.Utxo
		totalAmount, payerAmount                   float64
		availableQuantity                          int
		addr, payerAddr                            btcutil.Address
		wif, payerWif                              *btcutil.WIF
	)

	// 平台支付账户私钥
	payerWif, err = btcutil.DecodeWIF(platformPriKey)
	if err != nil {
		return err
	}
	// 待分割账号私钥
	wif, err = btcutil.DecodeWIF(priKey)
	if err != nil {
		return err
	}

	// 支付手续费账号地址
	payerAddr, err = btcutil.DecodeAddress(platformAddress, &chaincfg.MainNetParams)
	if err != nil {
		return err
	}

	// 待分割账号地址
	addr, err = btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err != nil {
		return err
	}

	// 手续费 暂时先定死20000
	const FeeValue = float64(20000)

	// 主网拉取待分割地址的 utxo
	utxos, err = usdt.GetUnspent(wire.MainNet, address)
	if err != nil {
		return
	}

	// 主网拉取支付手续费地址的 utxo
	payerUtxos, err = usdt.GetUnspent(wire.MainNet, platformAddress)
	if err != nil {
		return
	}

	payerUseUtxos = make([]usdt.Utxo, 0, len(payerUtxos))
	for _, utxo := range payerUtxos {
		// 未达到确认数限制
		if utxo.Confirmations < confirmationLimit {
			continue
		}
		payerAmount += utxo.Value
		payerUseUtxos = append(payerUseUtxos, utxo)
		if payerAmount >= FeeValue {
			break
		}
	}

	useUtxos = make([]usdt.Utxo, 0, len(utxos))
	for _, utxo := range utxos {
		// 未达到确认数限制
		if utxo.Confirmations < confirmationLimit {
			continue
		}

		// utxo 金额不等于 satoshi 的不统计进可用数量
		if utxo.Value != satoshiLimit {
			useUtxos = append(useUtxos, utxo)
			totalAmount += utxo.Value
			continue
		}

		// 可用 utxo 数量
		availableQuantity++
	}

	utxoQuantity -= availableQuantity

	if utxoQuantity <= 0 {
		fmt.Println("当前可用 utxo 数量 %v , 无须分割", availableQuantity)
		return
	}

	if totalAmount < satoshiLimit {
		fmt.Println("剩余 utxo 金额 %v , 不足以继续分割 ", totalAmount)
		return
	}

	txb := wire.NewMsgTx(TxVersion)

	// ========================== begin 输入 ==========================
	// 输入支付手续费的 utxo
	for _, v := range payerUseUtxos {
		utxoHash, _ := chainhash.NewHashFromStr(v.TxHashBigEndian)
		utxo := wire.NewOutPoint(utxoHash, uint32(v.TxOutputN))
		txIn := wire.NewTxIn(utxo, nil, nil)
		txIn.Sequence = 0xfffffffe
		txb.AddTxIn(txIn)
	}

	// 输入待分割的 utxo
	for _, v := range useUtxos {
		utxoHash, _ := chainhash.NewHashFromStr(v.TxHashBigEndian)
		utxo := wire.NewOutPoint(utxoHash, uint32(v.TxOutputN))
		txIn := wire.NewTxIn(utxo, nil, nil)
		txIn.Sequence = 0xfffffffe
		txb.AddTxIn(txIn)
	}
	// ========================== end 输入 ==========================
	//
	//
	//
	// ========================== begin 输出 ==========================
	const QuantityLimit = 400
	// 实际可分数量
	avaQuantity := int(math.Floor(totalAmount / satoshiLimit))
	// 分割数量上限
	if avaQuantity > utxoQuantity {
		avaQuantity = utxoQuantity
	}
	// 一次最大分割上限
	if avaQuantity > QuantityLimit {
		avaQuantity = QuantityLimit
	}
	// 输出支付手续费后剩余金额
	if payerLastAmount := payerAmount - FeeValue; payerLastAmount > 0 {
		pkScript, err := txscript.PayToAddrScript(payerAddr)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(int64(payerLastAmount), pkScript)
		txb.AddTxOut(txOut)
	}

	// 输出分割
	for i := 0; i < avaQuantity; i++ {
		pkScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(int64(satoshiLimit), pkScript)
		txb.AddTxOut(txOut)
	}

	// 输出分割剩余金额
	if lastAmount := totalAmount - float64(avaQuantity)*satoshiLimit; lastAmount > 0 {
		pkScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(int64(lastAmount), pkScript)
		txb.AddTxOut(txOut)
	}
	// ========================== end 输出 ==========================
	//
	//
	// ========================== begin 签名 ==========================
	offset := 0
	for i, v := range payerUseUtxos {
		script, err := hex.DecodeString(v.Script)
		if err != nil {
			return err
		}
		sigScript, err := txscript.SignatureScript(txb, i, script, txscript.SigHashAll, payerWif.PrivKey, payerWif.CompressPubKey)
		if err != nil {
			return err
		}
		txb.TxIn[i].SignatureScript = sigScript
		offset++
	}
	for i, v := range useUtxos {
		ii := i + offset
		script, err := hex.DecodeString(v.Script)
		if err != nil {
			return err
		}
		sigScript, err := txscript.SignatureScript(txb, ii, script, txscript.SigHashAll, wif.PrivKey, wif.CompressPubKey)
		if err != nil {
			return err
		}
		txb.TxIn[ii].SignatureScript = sigScript
	}
	// ========================== end 签名 ==========================

	var signedTx bytes.Buffer
	txb.Serialize(&signedTx)
	fmt.Println(hex.EncodeToString(signedTx.Bytes()))
	fmt.Println("txid", txb.TxHash().String())

	var resp map[string]interface{}
	resp, err = explorer.NewExplorer().TransactionPush(hex.EncodeToString(signedTx.Bytes()))
	if err != nil {
		return
	}
	result, err := json.Marshal(resp)
	if err != nil {
		return
	}
	fmt.Println(string(result))

	return
}
