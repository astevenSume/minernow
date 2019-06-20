package usdt

import (
	"bytes"
	"common"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	. "otc_error"
	"sort"
	"time"
	"usdt/prices"
	"utils/usdt/dao"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	lock "github.com/bsm/redis-lock"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

const REDIS_KEY_RECOMMENDED_FEES = "usdt_recommended_fees"

type Utxo struct {
	TxHash          string  `json:"tx_hash"`
	TxHashBigEndian string  `json:"tx_hash_big_endian"`
	TxIndex         int     `json:"tx_index"`
	TxOutputN       int     `json:"tx_output_n"`
	Script          string  `json:"script"`
	Value           float64 `json:"value"`
	ValueHex        string  `json:"value_hex"`
	Confirmations   int     `json:"confirmations"`
}

type Utxos []Utxo

func (u Utxos) Len() int {
	return len(u)
}

// 由高到低排序
func (u Utxos) Less(i, j int) bool {
	return u[i].Value > u[j].Value
}

func (u Utxos) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

var (
	ErrUsdtSenderUtxoInvalid   = errors.New("sender utxo invalid.")
	ErrUsdtPlatformUtxoInvalid = errors.New("platform utxo invalid.")
	ErrUsdtGetSignedTxFailed   = errors.New("get signed tx from node-offlinetx failed.")
)

type Msg struct {
	PayerPriKey         string  `json:"payer_pri_key"`
	PayerAddress        string  `json:"payer_address"`
	PayerUnspentOutputs Utxos   `json:"payer_unspent_outputs"`
	PriKey              string  `json:"pri_key"`
	From                string  `json:"from"`
	UnspentOutputs      Utxos   `json:"unspent_outputs"`
	UsdtValue           string  `json:"usdt_value"`
	FeeValue            float64 `json:"fee_value"`
	To                  string  `json:"to"`
}

func getOnceSignedTx(priKey string, usdtValue string, from, to string, fee int64) (signedTx string, err error) {
	msg := Msg{
		PriKey:    priKey,
		UsdtValue: usdtValue,
		FeeValue:  float64(0),
		From:      from,
		To:        to,
	}

	// add from utxos
	msg.UnspentOutputs, err = getUnspent(wire.MainNet, from)
	if err != nil {
		return
	}
	if len(msg.UnspentOutputs) <= 0 {
		err = ErrUsdtSenderUtxoInvalid
		common.LogFuncError("%v", err)
		return
	}
	//get platform utxos
	var uaid uint64 = UsdtConfig.PlatformUaid
	var platformPri, platformAddr string
	platformPri, platformAddr, err = priKeyMgr.Get(uaid)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	// add platform utxos
	msg.PayerUnspentOutputs, err = getUnspent(wire.MainNet, platformAddr)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	if len(msg.PayerUnspentOutputs) <= 0 {
		err = ErrUsdtPlatformUtxoInvalid
		common.LogFuncError("%v", err)
		return
	}

	msg.PayerPriKey = platformPri
	msg.PayerAddress = platformAddr
	msg.FeeValue = float64(fee) / 1e8

	signedTx, err = subGetSigned(msg)
	if err != nil {
		return
	}

	return
}

func getSignedTx2(priKey string, usdtValue string, from, to string, fee, maxFee int) (signedTx string, err error) {
	msg := Msg{
		PriKey:    priKey,
		UsdtValue: usdtValue,
		FeeValue:  float64(0),
		From:      from,
		To:        to,
	}

	// add from utxos
	msg.UnspentOutputs, err = getUnspent(wire.MainNet, from)
	if err != nil {
		return
	}
	if len(msg.UnspentOutputs) <= 0 {
		err = ErrUsdtSenderUtxoInvalid
		common.LogFuncError("%v", err)
		return
	}
	//get platform utxos
	var uaid uint64 = UsdtConfig.PlatformUaid
	var platformPri, platformAddr string
	platformPri, platformAddr, err = priKeyMgr.Get(uaid)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	// add platform utxos
	msg.PayerUnspentOutputs, err = getUnspent(wire.MainNet, platformAddr)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	if len(msg.PayerUnspentOutputs) <= 0 {
		err = ErrUsdtPlatformUtxoInvalid
		common.LogFuncError("%v", err)
		return
	}

	msg.PayerPriKey = platformPri
	msg.PayerAddress = platformAddr

	signedTx, err = subGetSigned(msg)
	if err != nil {
		return
	}

	// sizeof signedTx * 0.5 * fee(recommanded) is the best fee
	feeValue := float64(len(signedTx)/2) * float64(fee) / 100000000
	tmpValue := int(feeValue * 100000000)

	// adjust by MaxFee
	if tmpValue > maxFee {
		feeValue = float64(maxFee) / 100000000.0
	}
	msg.FeeValue = feeValue

	signedTx, err = subGetSigned(msg)
	if err != nil {
		return
	}

	return
}

// 参考链接 : 	https://www.thepolyglotdeveloper.com/2018/03/create-sign-bitcoin-transactions-golang/
func subGetSigned(msg Msg) (string, error) {
	payerPriKey, payerUtxos, payerFromAddr, priKey, utxos, fromAddr, feeValue, usdtValue, toAddress :=
		msg.PayerPriKey, msg.PayerUnspentOutputs, msg.PayerAddress, msg.PriKey, msg.UnspentOutputs, msg.From, msg.FeeValue, msg.UsdtValue, msg.To

	const TxVersion = 2

	txb := wire.NewMsgTx(TxVersion)

	wifPayer, err := btcutil.DecodeWIF(payerPriKey)
	if err != nil {
		return "", err
	}
	wif, err := btcutil.DecodeWIF(priKey)
	if err != nil {
		return "", err
	}
	const (
		DUST = 546
	)
	var (
		payerFundValue, senderFundValue, payerTotalUnspent, totalUnspent, payerChangeValue, changeValue float64 = DUST, 0, 0, 0, 0, 0

		payerUseUtxos, useUtxos Utxos
	)

	payerUseUtxos = make(Utxos, 0, len(payerUtxos))
	useUtxos = make(Utxos, 0, len(utxos))

	// 按金额从大到小排序
	sort.Sort(payerUtxos)
	sort.Sort(utxos)

	for _, v := range utxos {
		// 用户总的可用金额
		totalUnspent += v.Value
		// 金额不够支付 546 时
		if totalUnspent < payerFundValue {
			// 添加到待使用的 utxo
			useUtxos = append(useUtxos, v)
		} else {
			// 金额足够支付 546 时由用户支付
			senderFundValue, payerFundValue = payerFundValue, 0
			break
		}
	}

	// 用户账户足以支付 546 且总额等于 546 时,不使用用户支付
	if senderFundValue > 0 && totalUnspent == senderFundValue {
		// 546 由平台出
		senderFundValue, payerFundValue = payerFundValue, senderFundValue
		// 仅使用一个 utxo
		useUtxos = useUtxos[:1]
		totalUnspent = useUtxos[0].Value
	}

	for _, v := range payerUtxos {
		// 平台总的可用金额
		payerTotalUnspent += v.Value

		// 金额不够支付费用时
		if payerTotalUnspent < (payerFundValue - float64(feeValue*1e8)) {
			// 添加到待使用的 utxo
			payerUseUtxos = append(payerUseUtxos, v)
		} else {
			break
		}
	}

	// 平台金额不足以支付费用
	if payerTotalUnspent < float64(feeValue*1e8)+payerFundValue {
		return "", fmt.Errorf("Total %v less than feeValue*1e8 %v + payerFundValue %v", payerTotalUnspent, feeValue*1e8, payerFundValue)
	}

	// 平台余额变更
	payerChangeValue = payerTotalUnspent - payerFundValue - float64(feeValue*1e8)
	// 用户余额变更
	changeValue = changeValue - senderFundValue

	// 用户输入
	for _, v := range useUtxos {
		utxoHash, _ := chainhash.NewHashFromStr(v.TxHashBigEndian)
		utxo := wire.NewOutPoint(utxoHash, uint32(v.TxOutputN))
		txIn := wire.NewTxIn(utxo, nil, nil)
		txIn.Sequence = 0xfffffffe
		txb.AddTxIn(txIn)
	}

	// 平台输入
	for _, v := range payerUseUtxos {
		utxoHash, _ := chainhash.NewHashFromStr(v.TxHashBigEndian)
		utxo := wire.NewOutPoint(utxoHash, uint32(v.TxOutputN))
		txIn := wire.NewTxIn(utxo, nil, nil)
		txIn.Sequence = 0xfffffffe
		txb.AddTxIn(txIn)
	}

	// 平台输出
	{
		sourceAddress, err := btcutil.DecodeAddress(payerFromAddr, &chaincfg.MainNetParams)
		if err != nil {
			return "", err
		}
		pkScript, err := txscript.PayToAddrScript(sourceAddress)
		if err != nil {
			return "", err
		}
		txOut := wire.NewTxOut(int64(payerChangeValue), pkScript)
		txb.AddTxOut(txOut)
	}

	// usdt 输出
	var usdtInfo = fmt.Sprintf("%s%s%s%s", "6f6d6e69", "0000", "00000000001f", usdtValue)
	if data, err := hex.DecodeString(usdtInfo); err != nil {
		return "", err
	} else {
		sourcePkScript, err := txscript.NullDataScript(data)
		if err != nil {
			return "", err
		}
		txOut := wire.NewTxOut(0, sourcePkScript)
		txb.AddTxOut(txOut)
	}

	// 接收者输出
	{
		sourceAddress, err := btcutil.DecodeAddress(toAddress, &chaincfg.MainNetParams)
		if err != nil {
			return "", err
		}
		pkScript, err := txscript.PayToAddrScript(sourceAddress)
		if err != nil {
			return "", err
		}
		txOut := wire.NewTxOut(int64(DUST), pkScript)
		txb.AddTxOut(txOut)
	}

	{
		sourceAddress, err := btcutil.DecodeAddress(fromAddr, &chaincfg.MainNetParams)
		if err != nil {
			return "", err
		}
		pkScript, err := txscript.PayToAddrScript(sourceAddress)
		if err != nil {
			return "", err
		}
		txOut := wire.NewTxOut(int64(totalUnspent), pkScript)
		txb.AddTxOut(txOut)
	}

	for i, v := range useUtxos {
		switch {
		case fromAddr[0:1] == "3":
			p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.MainNetParams)
			if err != nil {
				return "", err
			}

			var witnessProgram []byte

			witnessProgram, err = txscript.PayToAddrScript(p2wkhAddr)
			if err != nil {
				return "", err
			}

			sigScript, err := txscript.NewScriptBuilder().AddData(witnessProgram).Script()
			if err != nil {
				return "", err
			}

			witnessScript, err := txscript.WitnessSignature(txb, txscript.NewTxSigHashes(txb), i, int64(v.Value), witnessProgram, txscript.SigHashAll, wif.PrivKey, wif.CompressPubKey)
			if err != nil {
				return "", err
			}

			txb.TxIn[i].SignatureScript = sigScript
			txb.TxIn[i].Witness = witnessScript
		default:
			script, err := hex.DecodeString(v.Script)
			if err != nil {
				return "", err
			}

			sigScript, err := txscript.SignatureScript(txb, i, script, txscript.SigHashAll, wif.PrivKey, wif.CompressPubKey)
			if err != nil {
				return "", err
			}

			txb.TxIn[i].SignatureScript = sigScript
		}
	}
	offset := len(utxos)
	for i, v := range payerUseUtxos {
		switch {
		case payerFromAddr[0:1] == "3":
			p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(
				btcutil.Hash160(wifPayer.SerializePubKey()), &chaincfg.MainNetParams,
			)
			if err != nil {
				return "", err
			}

			var witnessProgram []byte

			witnessProgram, err = txscript.PayToAddrScript(p2wkhAddr)
			if err != nil {
				return "", err
			}

			sigScript, err := txscript.NewScriptBuilder().AddData(witnessProgram).Script()
			if err != nil {
				return "", err
			}

			witnessScript, err := txscript.WitnessSignature(txb, txscript.NewTxSigHashes(txb),
				i+offset, int64(v.Value), witnessProgram,
				txscript.SigHashAll, wifPayer.PrivKey, wifPayer.CompressPubKey,
			)
			if err != nil {
				return "", err
			}

			txb.TxIn[i+offset].SignatureScript = sigScript
			txb.TxIn[i+offset].Witness = witnessScript
		default:
			script, err := hex.DecodeString(v.Script)
			if err != nil {
				return "", err
			}
			sigScript, err := txscript.SignatureScript(txb, i+offset, script, txscript.SigHashAll, wifPayer.PrivKey, wifPayer.CompressPubKey)
			if err != nil {
				return "", err
			}
			txb.TxIn[i+offset].SignatureScript = sigScript
		}
	}

	var signedTx bytes.Buffer
	txb.Serialize(&signedTx)

	return hex.EncodeToString(signedTx.Bytes()), nil
}

func subGetSigned_bak(msg Msg) (signedTx string, err error) {
	req := httplib.Post("http://localhost:" + fmt.Sprint(UsdtConfig.OfflineTxPort))
	req, err = req.JSONBody(&msg)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	type MsgResp struct {
		Code     int    `json:"code,omitempty"`
		Msg      string `json:"msg,omitempty"`
		SignedTx string `json:"signed_tx,omitempty"`
	}

	resp := MsgResp{}

	err = req.ToJSON(&resp)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	if resp.Code != 0 {
		common.LogFuncError("get signedTx from node-offlinetx failed : %s", resp.Msg)
		err = ErrUsdtGetSignedTxFailed
		return
	}

	signedTx = resp.SignedTx

	return
}
func GetUnspent(netType wire.BitcoinNet, addr string) (utxos []Utxo, err error) {
	return getUnspent(netType, addr)
}

// get unspent utxos from blockchain.info/testnet.blockchain.info
// only the utxo confirmation >= 6 will be used.
// https://bitcoin.stackexchange.com/questions/1170/why-is-6-the-number-of-confirms-that-is-considered-secure
func getUnspent(netType wire.BitcoinNet, addr string) (utxos []Utxo, err error) {
	var host string
	switch netType {
	case wire.MainNet:
		{
			host = "blockchain.info"
		}
	case wire.TestNet, wire.TestNet3:
		{
			host = "testnet.blockchain.info"
		}
	default:
		common.LogFuncError("unkown net type %v", netType)
		err = fmt.Errorf("unkown net type %v", netType)
		return
	}

	url := "https://" + host + "/zh-cn/unspent?active=" + addr

	var respBody string
	respBody, err = httplib.Get(url).String()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	} else if respBody == "No free outputs to spend" {
		return
	}

	type Msg struct {
		Notice         string `json:"notice"`
		UnspentOutputs []Utxo `json:"unspent_outputs"`
	}

	msg := Msg{}

	err = json.Unmarshal([]byte(respBody), &msg)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	// only the utxo whose confirmation over limit will be used.
	for _, utxo := range msg.UnspentOutputs {
		if utxo.Confirmations >= UsdtConfig.ConfirmationLimit {
			utxos = append(utxos, utxo)
		}
	}

	return
}

type FeeMsg struct {
	HalfHourFee int `json:"halfHourFee"`
	HourFee     int `json:"hourFee"`
	TwoHourFee  int `json:"twoHourFee"`
	FourHourFee int `json:"fourHourFee"`
}

// get fee rate as follow :
// - https://bitcoinfees.earn.com/api/v1/fees/recommended
// + https://bitcoinfees.earn.com/api/v1/fees/list
func getFeeRecommanded() (fee *FeeMsg, err error) {
	defer func() {
		if err != nil {
			fee.HalfHourFee = UsdtConfig.UnitFee
			fee.HourFee = UsdtConfig.UnitFee
			fee.TwoHourFee = UsdtConfig.UnitFee
			fee.FourHourFee = UsdtConfig.UnitFee
			common.LogFuncError("get fee failed, will set default 20 : %v", err)
			return
		}
	}()

	fee = &FeeMsg{}
	var recommendedFees string
	if recommendedFees, err = common.RedisManger.Get(REDIS_KEY_RECOMMENDED_FEES).Result(); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(recommendedFees), fee); err != nil {
		return
	}

	return
}

const (
	FeeModeUnkown = iota
	FeeModeHalfHour
	FeeModeHour
	FeeModeTwoHour
	FeeModeFourHour
	FeeModeNormal
)

// get fee by feeMode
func getFee(feeMode int) (fee int, err error) {
	switch feeMode {
	case FeeModeHalfHour, FeeModeHour, FeeModeTwoHour, FeeModeFourHour:
		{
			f := &FeeMsg{}
			f, err = getFeeRecommanded()
			// 取不到fee用默认的fee
			if err != nil { //record log, but not break.
				fee = UsdtConfig.UnitFee
				common.LogFuncError("%v", err)
				err = nil
				return
			}

			switch feeMode {
			case FeeModeHalfHour:
				fee = f.HalfHourFee
			case FeeModeHour:
				fee = f.HourFee
			case FeeModeTwoHour:
				fee = f.TwoHourFee
			case FeeModeFourHour:
				fee = f.FourHourFee
			}
		}
	case FeeModeNormal:
		{
			fee = UsdtConfig.UnitFee
		}
	default:
		err = fmt.Errorf("unkown FeeMode %d", feeMode)
		return
	}

	return
}

func Init() (err error) {
	var pUaid int64
	pUaid, err = beego.AppConfig.Int64("usdt::platformuaid")
	UsdtConfig.PlatformUaid = uint64(pUaid)
	if err != nil {
		common.LogFuncWarning("may be lack of usdt::platformuaid.")
		return
	}

	UsdtConfig.Precision, err = beego.AppConfig.Int("usdt::precision")
	if err != nil {
		common.LogFuncWarning("may be lack of usdt::precision.")
		return
	}

	UsdtConfig.OfflineTxPort, err = beego.AppConfig.Int("usdt::otxport")
	if err != nil {
		common.LogFuncWarning("may be lack of usdt::otxport.")
		return
	}

	UsdtConfig.MinFee, err = beego.AppConfig.Int("usdt::MinFee")
	if err != nil {
		common.LogFuncWarning("may be lack of usdt::MinFee.")
		return
	}

	UsdtConfig.MaxFee, err = beego.AppConfig.Int("usdt::MaxFee")
	if err != nil {
		common.LogFuncWarning("may be lack of usdt::MaxFee.")
		return
	}

	UsdtConfig.UnitFee, err = beego.AppConfig.Int("usdt::unitFee")
	if err != nil {
		common.LogFuncWarning("may be lack of usdt::unitFee.")
		return
	}

	UsdtConfig.Symbol = beego.AppConfig.String("usdt::Symbol")
	if UsdtConfig.Symbol == "" {
		UsdtConfig.Symbol = "USDT"
	}

	UsdtConfig.ConfirmationLimit, err = beego.AppConfig.Int("usdt::ConfirmationLimit")
	if err != nil {
		common.LogFuncWarning("may be lack of usdt::ConfirmationLimit.")
		return
	}

	UsdtConfig.FeeMode, err = beego.AppConfig.Int("usdt::feeMode")
	if err != nil {
		common.LogFuncWarning("may be lack of usdt::feeMode.")
		return
	}

	UsdtConfig.SyncFrequency, err = beego.AppConfig.Int64("usdt::syncFrequency")
	if err != nil {
		common.LogFuncWarning("may be lack of usdt::syncFrequency.")
		return
	}

	err = dao.Init(nil)
	if err != nil {
		return
	}

	return
}

func InitSyncRechargeTransaction() error {
	// run sync recharge transaction go runtine
	go common.SafeRun(syncRechargeTransactionFunc)()
	// run sync wallet transaction go runtine
	go common.SafeRun(syncWalletTransactionFunc)()

	return nil
}

func syncRechargeTransactionFunc() {
	interval, err := beego.AppConfig.Int64("usdt::SyncOnChainTxInterval")
	if err != nil {
		interval = 1800 //default 30 minutes
	}
	var lockSecs int64
	lockSecs, err = beego.AppConfig.Int64("usdt::SyncOnChainTxLockSecs")
	if err != nil {
		lockSecs = 10800 //default 3 hours
	}

	for {
		l, err := common.RedisLock2("usdt_sync_onchain_tx", lock.Options{
			LockTimeout: time.Second * time.Duration(lockSecs),
			RetryCount:  common.RetryCount,
			RetryDelay:  common.RetryDelay,
		})
		if err != nil { //doesn't get redis lock, skip
			if err != lock.ErrLockNotObtained {
				common.LogFuncError("%v", err)
			}
		} else { //got redis lock, do the job
			SyncRechargeTransaction()
			common.RedisUnlock(l)
		}

		//take a rest
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func syncWalletTransactionFunc() {
	interval, err := beego.AppConfig.Int64("usdt::SyncWalletOnChainTxInterval")
	if err != nil {
		interval = 1800 //default 30 minutes
	}
	var lockSecs int64
	lockSecs, err = beego.AppConfig.Int64("usdt::SyncWalletOnChainTxLockSecs")
	if err != nil {
		lockSecs = 10800 //default 3 hours
	}

	for {
		l, err := common.RedisLock2("usdt_sync_wallet_onchain_tx", lock.Options{
			LockTimeout: time.Second * time.Duration(lockSecs),
			RetryCount:  common.RetryCount,
			RetryDelay:  common.RetryDelay,
		})
		if err != nil { //doesn't get redis lock, skip
			if err != lock.ErrLockNotObtained {
				common.LogFuncError("%v", err)
			}
		} else { //got redis lock, do the job
			SyncWalletTransaction()
			common.RedisUnlock(l)
		}

		//take a rest
		time.Sleep(time.Second * time.Duration(interval))
	}
}

type UsdtRecommandedFee struct {
	HalfHourCNY string `json:"half_hour_cny"`
	HourCNY     string `json:"hour_cny"`
	TwoHourCNY  string `json:"two_hour_cny"`
	FourHourCNY string `json:"four_hour_cny"`

	// HalfHourBTC string `json:"half_hour_btc"`
	// HourBTC     string `json:"hour_btc"`
	// TwoHourBTC  string `json:"two_hour_btc"`
	// FourHourBTC string `json:"four_hour_btc"`

	HalfHourUSDT string `json:"half_hour_usdt"`
	HourUSDT     string `json:"hour_usdt"`
	TwoHourUSDT  string `json:"two_hour_usdt"`
	FourHourUSDT string `json:"four_hour_usdt"`
}

func GetRecommandedFee(amount int64) (fee UsdtRecommandedFee, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	Premium, err := beego.AppConfig.Float("usdt::Premium")
	if err != nil {
		errCode = ERROR_CODE_CONFIG_LACK
		return
	}
	if Premium > 1 || Premium < 0 {
		errCode = ERROR_CODE_CONFIG_LACK
		return
	}

	Premium = Premium + 1

	//generate transaction and broadcast it to mainnet
	amountTx := fmt.Sprintf("%.16x", amount)

	msg := Msg{
		UsdtValue: amountTx,
		FeeValue:  float64(0),
	}
	var (
		signedTx          string
		feeMsg            *FeeMsg
		btc2cny, usdt2cny float64
	)
	if signedTx, err = mockSign(msg); err != nil {
		errCode = ERROR_CODE_USDT_GET_SIGNED_TX_FAILED
		return
	}

	if btc2cny = prices.GetPriceFloat64(prices.PRICE_CURRENCY_TYPE_BTC); btc2cny <= 0 {
		errCode = ERROR_CODE_GET_BTC_TO_USDT_RATE_FAILED
		return
	}

	if usdt2cny = prices.GetPriceFloat64(prices.PRICE_CURRENCY_TYPE_USDT); usdt2cny <= 0 {
		errCode = ERROR_CODE_GET_BTC_TO_USDT_RATE_FAILED
		return
	}

	// return int((maxFeeCNY / btc2cny) * 1e8), nil

	if feeMsg, err = getFeeRecommanded(); err != nil {
		errCode = ERROR_CODE_USDT_GET_FEE_FAILED
		return
	}

	var (
		half, hour, two, four = float64(len(signedTx)/2) * float64(feeMsg.HalfHourFee) / 1e8,
			float64(len(signedTx)/2) * float64(feeMsg.HourFee) / 1e8,
			float64(len(signedTx)/2) * float64(feeMsg.TwoHourFee) / 1e8,
			float64(len(signedTx)/2) * float64(feeMsg.FourHourFee) / 1e8
	)

	// fee.HalfHourBTC = fmt.Sprintf("%.4f", half)
	// fee.HourBTC = fmt.Sprintf("%.4f", hour)
	// fee.TwoHourBTC = fmt.Sprintf("%.4f", two)
	// fee.FourHourBTC = fmt.Sprintf("%.4f", four)

	fee.HalfHourCNY = fmt.Sprintf("%.4f", half*btc2cny*Premium)
	fee.HourCNY = fmt.Sprintf("%.4f", hour*btc2cny*Premium)
	fee.TwoHourCNY = fmt.Sprintf("%.4f", two*btc2cny*Premium)
	fee.FourHourCNY = fmt.Sprintf("%.4f", four*btc2cny*Premium)

	fee.HalfHourUSDT = fmt.Sprintf("%.4f", half*btc2cny*Premium/usdt2cny)
	fee.HourUSDT = fmt.Sprintf("%.4f", hour*btc2cny*Premium/usdt2cny)
	fee.TwoHourUSDT = fmt.Sprintf("%.4f", two*btc2cny*Premium/usdt2cny)
	fee.FourHourUSDT = fmt.Sprintf("%.4f", four*btc2cny*Premium/usdt2cny)

	return

}

func mockSign(msg Msg) (signedTx string, err error) {
	msg.From = `1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa`
	msg.UnspentOutputs, err = getUnspent(wire.MainNet, msg.From)
	msg.PayerUnspentOutputs = msg.UnspentOutputs
	if err != nil {
		return
	}
	msg.PayerAddress = `1taEvJUixDUn4ZjiFC9sCv1vTpRQjX3yS`
	msg.To = msg.PayerAddress
	msg.PayerPriKey, err = common.DecryptFromBase64(`jZfSlpfObAmny32a9Cg1Oq+5dRi4T2b5JuJ3Kw/wG3ClUghoW3ZhI57L1A8DYxQB+GMziETzfilS/qM6reLCumWLbERJNPc9Uz2GMFzMBQ==`, PriAesKey)
	msg.PriKey = msg.PayerPriKey
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	signedTx, err = subGetSigned(msg)
	if err != nil {
		return
	}
	return
}
