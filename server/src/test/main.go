package main

import (
	"bytes"
	"common"
	"encoding/hex"
	"errors"
	"fmt"
	"usdt"
	"utils/usdt/models"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	omni "github.com/ibclabs/omnilayer-go"
	"github.com/ibclabs/omnilayer-go/omnijson"
)

var key = [32]byte{'T', 'r', 'W', '2', '4', 'f', 'x', 'M',
	'P', '6', 'Q', 'L', '7', 's', 'B', '1',
	'G', 'H', 'c', 'r', 'l', 'K', 'R', 'O',
	'A', 'U', 'C', '3', 'D', 'S', 'c', '='}

func testSplitUtxo() (err error) {

	var (
		FromPriKey  string
		FromAddress string
		FromUtxo    []usdt.Utxo
		FeeValue    float64
	)

	FromPriKey, err = common.DecryptFromBase64("eVAOEFqgQLln7Zl5VviZ7I7i/F/TnOHsW/MVlx1szzV8E7EF0SNvYwOvJ+ihhY3Ble3ghIe1mJH2d3uLglS4r5SbVyTA46ZPQE+B++KIkA==", usdt.PriAesKey)
	FromAddress = "1HgAb9RcnkxWv7p7CgJ6af9QkRoUVy3qNg"

	FromUtxo, err = usdt.GetUnspent(wire.MainNet, FromAddress)

	txb := wire.NewMsgTx(TxVersion)
	var setPayer *btcutil.WIF
	// 平台支付账户私钥
	setPayer, err = btcutil.DecodeWIF(FromPriKey)
	if err != nil {
		return err
	}

	var (
		fundValue, totalUnspent, changeValue float64 = 546, 0, 0
	)

	for _, v := range FromUtxo {
		totalUnspent += v.Value
	}

	FeeValue = float64(2000) / 1e8

	if totalUnspent < float64(FeeValue*1e8)+fundValue {
		return errors.New(fmt.Sprintf("Total %v less than feeValue*1e8 %v + fundValue %v  %v", totalUnspent, FeeValue*1e8, fundValue, (float64(FeeValue*1e8) + fundValue)))
	}

	// payer's change
	changeValue = totalUnspent - fundValue - float64(FeeValue*1e8)
	fmt.Println(float64(FeeValue*1e8) + fundValue)

	// add payer utxos
	for _, v := range FromUtxo {
		utxoHash, _ := chainhash.NewHashFromStr(v.TxHashBigEndian)
		utxo := wire.NewOutPoint(utxoHash, uint32(v.TxOutputN))
		txIn := wire.NewTxIn(utxo, nil, nil)
		txIn.Sequence = 0xfffffffe
		txb.AddTxIn(txIn)
	}

	halfValue := changeValue / 4

	var address btcutil.Address
	address, err = btcutil.DecodeAddress(FromAddress, &chaincfg.MainNetParams)
	if err != nil {
		return err
	}

	fmt.Println(address.EncodeAddress(), address.String(), address.ScriptAddress())

	pkScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return err
	}
	txOut := wire.NewTxOut(int64(halfValue), pkScript)
	txb.AddTxOut(txOut)

	txOut2 := wire.NewTxOut(int64(halfValue), pkScript)
	txb.AddTxOut(txOut2)

	txOut3 := wire.NewTxOut(int64(halfValue), pkScript)
	txb.AddTxOut(txOut3)

	txOut4 := wire.NewTxOut(int64(halfValue), pkScript)
	txb.AddTxOut(txOut4)

	for i, v := range FromUtxo {
		script, err := hex.DecodeString(v.Script)
		if err != nil {
			return err
		}
		sigScript, err := txscript.SignatureScript(txb, i, script, txscript.SigHashAll, setPayer.PrivKey, setPayer.CompressPubKey)
		if err != nil {
			return err
		}
		txb.TxIn[i].SignatureScript = sigScript
	}

	var signedTx bytes.Buffer
	txb.Serialize(&signedTx)
	fmt.Println(hex.EncodeToString(signedTx.Bytes()))
	fmt.Println("txid", txb.TxHash().String())
	return nil

}
func testOmni() {

	c := omni.New(&omni.ConnConfig{
		Host:     "",
		Endpoint: "",
		User:     "",
		Pass:     "",
	})

	result, err := c.OmniGetBalance(omnijson.OmniGetBalanceCommand{
		Address:    "",
		PropertyID: 31,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}

func pubKeyHashToScript(pubKey []byte) []byte {
	pubKeyHash := btcutil.Hash160(pubKey)
	script, err := txscript.NewScriptBuilder().
		AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
	if err != nil {
		panic(err)
	}
	return script
}

func ttttttt() {

	var priKey *btcec.PrivateKey
	priKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var param *chaincfg.Params = &chaincfg.MainNetParams

	// wif, err := btcutil.NewWIF(priKey, param, true)
	// if err != nil {
	// 	return
	// }

	pubKey := priKey.PubKey().SerializeCompressed()
	hash160 := btcutil.Hash160(pubKey)
	addr, err := btcutil.NewAddressWitnessPubKeyHash(hash160, param)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("bc1:", addr.String())
	addr3(hash160)

	w, err := btcutil.NewAddressScriptHash(pubKeyHashToScript(pubKey), &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("3:", w.String())

	wif, err := btcutil.NewWIF(priKey, &chaincfg.MainNetParams, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("pri key:", wif.String())
}

func addr3(pubKeyHash []byte) {

	witAddr, _ := btcutil.NewAddressWitnessPubKeyHash(
		pubKeyHash, &chaincfg.MainNetParams,
	)

	witnessProgram, _ := txscript.PayToAddrScript(witAddr)

	address, _ := btcutil.NewAddressScriptHash(
		witnessProgram, &chaincfg.MainNetParams,
	)

	fmt.Println("3:", address.EncodeAddress())
}
func testTran() {
	var PayerAddress = "1GPa2CEu887nocR2MVzc993q925xHMRmpJ"
	var From = "1DXDso8dNotaWxKu9MDt48ZNNgqgxJ1ggC"
	payerUtxo, err := usdt.GetUnspent(wire.MainNet, PayerAddress)
	if err != nil {
		return
	}
	utxo, err := usdt.GetUnspent(wire.MainNet, From)
	if err != nil {
		return
	}
	PayerPriKey, err := common.DecryptFromBase64("JFSex+yFAqwLisFa/7wATm8bErHfsuiEPw0lx/njqZUNeh56RoGpHFMmM75WkGkzHIUqmqn2ICEDFFGxN7eWT+7+7AXb2QZbcP3i4hZilQ==", usdt.PriAesKey)
	if err != nil {
		return
	}
	fmt.Println(PayerPriKey)
	CreateTransaction(Msg{
		PayerPriKey:         PayerPriKey,
		PayerAddress:        PayerAddress,
		PayerUnspentOutputs: payerUtxo,
		PriKey:              "5KGwqLKeTM44iSkHdmka5VraxP2iG86JERKZYaeQjU8UfNpRjNK",
		From:                From,
		UnspentOutputs:      utxo,
		UsdtValue:           fmt.Sprintf("%.16x", int64(1*1e8)),
		FeeValue:            0.00002,
		To:                  "353gNLXJgGxtnR1SX51zqDm3HTa4BQQ8vY",
	})

}
func testTran2() {
	var PayerAddress = "1HgAb9RcnkxWv7p7CgJ6af9QkRoUVy3qNg"
	var From = "1M2RzpoQ2XT5csVFd3q4m2A9maHtdnDNsR"
	payerUtxo, err := usdt.GetUnspent(wire.MainNet, PayerAddress)
	if err != nil {
		return
	}
	utxo, err := usdt.GetUnspent(wire.MainNet, From)
	if err != nil {
		return
	}
	PayerPriKey, err := common.DecryptFromBase64("eVAOEFqgQLln7Zl5VviZ7I7i/F/TnOHsW/MVlx1szzV8E7EF0SNvYwOvJ+ihhY3Ble3ghIe1mJH2d3uLglS4r5SbVyTA46ZPQE+B++KIkA==", usdt.PriAesKey)
	if err != nil {
		return
	}

	PriKey, err := common.DecryptFromBase64("L89gX491hHZFt8Cp7UuuNeu02khvHPNO2dw9euY/OefH0mKmDZJ1Q5vPZRvXON3bgFedxB4PIEPvb9RXViinawt2aq6noGjcOxTO2dXTIw==", usdt.PriAesKey)
	if err != nil {
		return
	}
	fmt.Println(PayerPriKey)
	CreateTransaction(Msg{
		PayerPriKey:         PayerPriKey,
		PayerAddress:        PayerAddress,
		PayerUnspentOutputs: payerUtxo,
		PriKey:              PriKey,
		From:                From,
		UnspentOutputs:      utxo,
		UsdtValue:           fmt.Sprintf("%.16x", int64(0.02*1e8)),
		FeeValue:            0.0002,
		To:                  "353gNLXJgGxtnR1SX51zqDm3HTa4BQQ8vY",
	})

}
func testTran3() {
	var (
		payerUtxo, utxo []Utxo = make([]Utxo, 0), make([]Utxo, 0)
	)
	usdt.UsdtConfig.ConfirmationLimit = 6
	var PayerAddress = "1ARxT5FMttHorqkXro1WEtgegKd87o8bi9"
	var From = "353gNLXJgGxtnR1SX51zqDm3HTa4BQQ8vY"
	tmpPayerUtxo, err := usdt.GetUnspent(wire.MainNet, PayerAddress)
	if err != nil {
		return
	}

	for _, u := range tmpPayerUtxo {
		if u.Value > 50000 {
			continue
		}
		if u.Confirmations >= 6 {
			payerUtxo = append(payerUtxo, u)
			break
		}
	}

	tmpUtxo, err := usdt.GetUnspent(wire.MainNet, From)
	if err != nil {
		return
	}

	for _, u := range tmpUtxo {
		if u.Confirmations >= 6 {
			utxo = append(utxo, u)
			break
		}
	}

	PayerPriKey, err := common.DecryptFromBase64("vvBb6wYrK5Bb60/u9xdmSBUD2ZzZeh9Z727mX1d1eG3w7hdPaYKBRWicRGR07je68zDsRr2h07ei0zeyCaqqj/L1OVDJCA+Dueo7kBl4rR4=", usdt.PriAesKey)
	if err != nil {
		return
	}

	PriKey := "L1Dcr55jsbzVJonKZHWXoaKEYaGcHzggUCRtScs8pJGJfHTsJf6r"

	PayFrom3Addr(Msg{
		PayerPriKey:         PayerPriKey,
		PayerAddress:        PayerAddress,
		PayerUnspentOutputs: payerUtxo,
		PriKey:              PriKey,
		From:                From,
		UnspentOutputs:      utxo,
		UsdtValue:           fmt.Sprintf("%.16x", int64(0.01*1e8)),
		FeeValue:            0.0003,
		To:                  "32tWWKNRBwTLeymNYaARbydJNNuf9wv62s",
	})

}
func main() {
	testTran3()
	return

	ttttttt()
	return
	// ttttttt()
	// // return

	// priKey, err := common.DecryptFromBase64("vvBb6wYrK5Bb60/u9xdmSBUD2ZzZeh9Z727mX1d1eG3w7hdPaYKBRWicRGR07je68zDsRr2h07ei0zeyCaqqj/L1OVDJCA+Dueo7kBl4rR4=", usdt.PriAesKey)
	// if err != nil {
	// 	return
	// }
	// wif, err := btcutil.DecodeWIF(priKey)
	// if err != nil {
	// 	return
	// }
	// hash160 := btcutil.Hash160(wif.SerializePubKey())
	// addr, err := btcutil.NewAddressWitnessPubKeyHash(hash160, &chaincfg.MainNetParams)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("bc1:", addr.String())

	// addr3(wif.SerializePubKey())
	// w, err := btcutil.NewAddressScriptHash(pubKeyHashToScript(wif.SerializePubKey()), &chaincfg.MainNetParams)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println("3:", w.String())

	// pubKeyAddress, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &chaincfg.MainNetParams)
	// fmt.Println("EncodeAddress:", pubKeyAddress.EncodeAddress())
	// return

	//
	//
	//
	//
	//
	// btconf := &rpcclient.ConnConfig{
	// 	Host:         "172.104.41.79:8888",
	// 	User:         "qdd123usdt",
	// 	Pass:         "Qdd.usdt567.rpc",
	// 	HTTPPostMode: true,
	// 	DisableTLS:   true,
	// }

	// c, err := rpcclient.New(btconf, nil)

	// if err != nil {
	// 	return
	// }
	// tx, err := hex.DecodeString("020000000569c67af4a4e352c4b9be1dab01e2f7b2d45e2bd5ee225986fd639c2d03dfac96030000008b483045022100907d4a8710d01d18d872af6655a82a3ff47b0959a35f37a143992ba5084f3ae2022064bbc1ceb5fa8c163d15eb81df191dbc0cd89f867e46524949470a5a96d64221014104a64c9def0d50e49092427c02ed54811475245a40a110344f2f3b423fd339af6bbafab2b71d7a08ba1a87d44a90dbc041f78492ca2da38efcf3819c1fae1c329dfefffffffc8913f7fa6f9969913797d5b6074b287b4cd000af9752594abe5597ccedb4b7020000008a4730440220143b47e648b17ebe212a7c76360d674921bf50b6a6a520247358dc85083152b802205f318a85c23c3eef45a5c808bb0f737af8850283b95d48f80ab834affa1ea386014104a64c9def0d50e49092427c02ed54811475245a40a110344f2f3b423fd339af6bbafab2b71d7a08ba1a87d44a90dbc041f78492ca2da38efcf3819c1fae1c329dfeffffff2409a65a43449763cd8418f76ce1db42f436ed52def36c4b12d68e66405fb2ac020000008a47304402206a73ed491284eaaa9eca3596ca1e2203230fd041d978e3e48e628e5c58b7d3f6022032b19fccc689498d6ca2ad0a6adec49443a0c29e601585219b299861dfb86514014104a64c9def0d50e49092427c02ed54811475245a40a110344f2f3b423fd339af6bbafab2b71d7a08ba1a87d44a90dbc041f78492ca2da38efcf3819c1fae1c329dfeffffff26eae0c43b3eb0663f2b1db8dbbb4b4ee0f15945f74c167b30fe20bf3b457bda020000008a47304402202339d8d07b23806fe13a3c4d24404deb239fd305664143407772230593b5b93a022030cc7ea1f6b3984bf5be8cf2bd1d247c1373b90a71b7473652251471955a5ecd014104a64c9def0d50e49092427c02ed54811475245a40a110344f2f3b423fd339af6bbafab2b71d7a08ba1a87d44a90dbc041f78492ca2da38efcf3819c1fae1c329dfeffffffcf88c1629a360ca1dc22473210cc3d1a6bb184d8c5f05eea8c52043a3352c018000000006a473044022051a6e9dcb55baafcd046a71dfce34d0ea7d3435679e532b4733600df4dbb8a0d02200375909c8844364b4a4f2f43ca01f7c31f2907aa1b68a5543765f517be91edbd012102e4914683606a153e6067f5565082cf0ac61a5b80cdc504a83236231794d0a853feffffff04fa318a01000000001976a914677166eec2d94b670a9084e727028e7b0762869988ac0000000000000000166a146f6d6e69000000000000001f00000000447888b022020000000000001976a914677166eec2d94b670a9084e727028e7b0762869988ac30330000000000001976a914a8ce9a23e4e8036fb6dd07f02aa32d0bdf21abec88ac00000000")
	// if err != nil {
	// 	return
	// }

	return
	testOmni()
	return

	fmt.Println(fmt.Sprintf("UPDATE %s SET  %s=%s-?, %s=%s+?, %s=? WHERE %s=? AND %s >=?",
		models.TABLE_UsdtAccount,
		models.COLUMN_UsdtAccount_TransferFrozenInteger, models.COLUMN_UsdtAccount_TransferFrozenInteger,
		models.COLUMN_UsdtAccount_OwnedByPlatformInteger, models.COLUMN_UsdtAccount_OwnedByPlatformInteger,
		models.COLUMN_UsdtAccount_Mtime,
		// where
		models.COLUMN_UsdtAccount_Uid,
		models.COLUMN_UsdtAccount_TransferFrozenInteger,
	))

	return
	// 分割钱包中的utxo
	fmt.Println(testSplitUtxo())
	return
	// fmt.Println(fmt.Sprintf("%.16x", int64(2*1e8)))
	// fmt.Println(usdt.GetUnspent(wire.MainNet,"1ARxT5FMttHorqkXro1WEtgegKd87o8bi9"))
	// return
	//api, _ := gameapi.NewRoyalGameAPI("http://www.devbj.com/api", "qdd_", "2cac7adc65935858c29a8488d6b7aecb")
	//
	//api.Login("test_account", "2cac7adc65935858c29a8488d6b7aecb", "", "506")
	//api.TransferIn("test_account", "2cac7adc65935858c29a8488d6b7aecb", 12.34)
	//api.TransferOut("test_account", "2cac7adc65935858c29a8488d6b7aecb", 12.34)
	//api.GetBalance("test_account", "2cac7adc65935858c29a8488d6b7aecb")
	//// api := gameapi.NewAsiaGamingAPI("http://203.78.143.203:8004")
	//return
	//fmt.Println("")
	//fmt.Println("")
	//fmt.Println("")
	//fmt.Println("")
	// var (
	// 	account  = "test_account"
	// 	password = "2cac7adc65935858c29a8488d6b7aecb"
	// )
	// account = "test_account2"
	// password = "2cac7adc65935858c29a8488d6b7aec2"

	// fmt.Println(api.Register(account, password, ""))
	// fmt.Println(api.Login(account, password, "", "0"))
	// fmt.Println(api.GetBalance(account, password))

	// fmt.Println(api.TransferIn(account, password, 1.02))
	// fmt.Println(api.GetBalance(account, password))

	// fmt.Println(api.TransferOut(account, password, 1.02))
	// fmt.Println(api.GetBalance(account, password))

	// api := gameapi.NewKaiYuanAPI()
	// api.Login("test_account", "127.0.0.1")
	// api.GetAccount("test_account")
	// api.TransferIn("test_account", 10)
	// api.TransferOut("test_account", 10)
	// return
	//testTransaction()
	//return
	//var src string
	////var key string
	//flag.StringVar(&src, "src", "", "")
	////flag.StringVar(&key, "key", "", "")
	//flag.Parse()
	//fmt.Println("Use  --src=  --key=")
	//if len(key) != 32 {
	//	fmt.Println("KEY must be 32 bytes long.")
	//	return
	//}
	//b := [32]byte{}
	//for i, v := range key {
	//	b[i] = byte(v)
	//}
	//fmt.Printf("Src %s\n", string(src))
	//
	//dst, err := common.EncryptToBase64(src, key)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Printf("Encrypt %s\n", string(dst))
	//
	//tmp, err := common.DecryptFromBase64(dst, key)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Printf("Check %s\n", string(tmp))
	//
	//dec2, err := common.DecryptFromBase64(src, key)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Printf("Decrypt %s\n", string(dec2))

	TestSafeRun()
	forever := make(chan bool)
	<-forever
}

func testTransaction() {
	CreateTransaction(Msg{
		PayerPriKey:  "KzmVvhF3J83d5nUbakCu83TL3Fk16qQGxEL48sBy2g67guNWjQq8",
		PayerAddress: "1ARxT5FMttHorqkXro1WEtgegKd87o8bi9",
		PayerUnspentOutputs: []Utxo{
			Utxo{
				TxHashBigEndian: "acb25f40668ed6124b6cf3de52ed36f442dbe16cf71884cd639744435aa60924",
				TxOutputN:       0,
				Script:          "76a914677166eec2d94b670a9084e727028e7b0762869988ac",
				Value:           26026257,
			},
		},
		PriKey: "5KGwqLKeTM44iSkHdmka5VraxP2iG86JERKZYaeQjU8UfNpRjNK",
		From:   "1DXDso8dNotaWxKu9MDt48ZNNgqgxJ1ggC",
		UnspentOutputs: []Utxo{Utxo{
			TxHashBigEndian: "acb25f40668ed6124b6cf3de52ed36f442dbe16cf71884cd639744435aa60924",
			TxOutputN:       3,
			Script:          "76a9148958bc86122f5cc454c8ce1a4bc7a6a540a3661f88ac",
			Value:           1092,
		}},
		UsdtValue: fmt.Sprintf("%.16x", int64(2*1e8)),
		FeeValue:  0.0002,
		To:        "1GPa2CEu887nocR2MVzc993q925xHMRmpJ",
	})
	// CreateTransaction(Msg{
	// 	PayerPriKey:  "5JPzYaNNyrUKTo7pZ5rMH1grWHkhSYSD6ieyXPuXGtSKbMZSD2S",
	// 	PayerAddress: "1HgAb9RcnkxWv7p7CgJ6af9QkRoUVy3qNg",
	// 	PayerUnspentOutputs: []Utxo{Utxo{
	// 		TxHash:          "4658675900c5341c9957bd5e6bf2fe5f5010a034a2d711923538eee21a560b68",
	// 		TxHashBigEndian: "680b561ae2ee38359211d7a234a010505ffef26b5ebd57991c34c50059675846",
	// 		TxOutputN:       0,
	// 		Script:          "76a914b6ea3d42f9fd8ef1cc44c60cb6abc9d6c082355a88ac",
	// 		Value:           189101,
	// 		ValueHex:        "02e2ad",
	// 		Confirmations:   6,
	// 		TxIndex:         445383079,
	// 	}},
	// 	PriKey: "5KGwqLKeTM44iSkHdmka5VraxP2iG86JERKZYaeQjU8UfNpRjNK",
	// 	From:   "1DXDso8dNotaWxKu9MDt48ZNNgqgxJ1ggC",
	// 	UnspentOutputs: []Utxo{Utxo{
	// 		TxHash:          "4658675900c5341c9957bd5e6bf2fe5f5010a034a2d711923538eee21a560b68",
	// 		TxHashBigEndian: "680b561ae2ee38359211d7a234a010505ffef26b5ebd57991c34c50059675846",
	// 		TxOutputN:       2,
	// 		Script:          "76a9148958bc86122f5cc454c8ce1a4bc7a6a540a3661f88ac",
	// 		Value:           1092,
	// 		ValueHex:        "0444",
	// 		Confirmations:   6,
	// 		TxIndex:         445383079,
	// 	}},
	// 	UsdtValue: fmt.Sprintf("%.16x", int64(2*1e8)),
	// 	FeeValue:  0.0002,
	// 	To:        "1GPa2CEu887nocR2MVzc993q925xHMRmpJ",
	// })
}

type Transaction struct {
	TxId               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	UnsignedTx         string `json:"unsignedtx"`
	SignedTx           string `json:"signedtx"`
}

type Msg struct {
	PayerPriKey         string  `json:"payer_pri_key"`
	PayerAddress        string  `json:"payer_address"`
	PayerUnspentOutputs []Utxo  `json:"payer_unspent_outputs"`
	PriKey              string  `json:"pri_key"`
	From                string  `json:"from"`
	UnspentOutputs      []Utxo  `json:"unspent_outputs"`
	UsdtValue           string  `json:"usdt_value"`
	FeeValue            float32 `json:"fee_value"`
	To                  string  `json:"to"`
}
type Utxo = usdt.Utxo

// type Utxo struct {
// 	TxHash          string  `json:"tx_hash"`
// 	TxHashBigEndian string  `json:"tx_hash_big_endian"`
// 	TxIndex         int     `json:"tx_index"`
// 	TxOutputN       int     `json:"tx_output_n"`
// 	Script          string  `json:"script"`
// 	Value           float64 `json:"value"`
// 	ValueHex        string  `json:"value_hex"`
// 	Confirmations   int     `json:"confirmations"`
// }

const (
	TxVersion = 2
)

func CreateTransaction(msg Msg) (err error) {
	defer func() {
		if err != nil {
			panic(err)
		}
	}()
	payerPriKey, payerUtxos, payerFromAddr, priKey, utxos, fromAddr, feeValue, usdtValue, toAddress :=
		msg.PayerPriKey, msg.PayerUnspentOutputs, msg.PayerAddress, msg.PriKey, msg.UnspentOutputs, msg.From, msg.FeeValue, msg.UsdtValue, msg.To

	txb := wire.NewMsgTx(TxVersion)
	// 平台支付账户私钥
	setPayer, err := btcutil.DecodeWIF(payerPriKey)
	if err != nil {
		return err
	}
	// sender 账户私钥
	set, err := btcutil.DecodeWIF(priKey)
	if err != nil {
		return err
	}
	var (
		fundValue, payerTotalUnspent, totalUnspent, payerChangeValue float64 = 546, 0, 0, 0
	)

	for _, v := range payerUtxos {
		payerTotalUnspent += v.Value
	}

	for _, v := range utxos {
		totalUnspent += v.Value
	}

	if payerTotalUnspent < float64(feeValue*1e8)+fundValue {
		return errors.New(fmt.Sprintf("Total %v less than feeValue*1e8 %v + fundValue %v", payerTotalUnspent, feeValue*1e8, fundValue))
	}

	// payer's change
	payerChangeValue = payerTotalUnspent - fundValue - float64(feeValue*1e8)
	fmt.Println(float64(feeValue*1e8) + fundValue)

	// add sender utxos
	for _, v := range utxos {
		utxoHash, _ := chainhash.NewHashFromStr(v.TxHashBigEndian)
		utxo := wire.NewOutPoint(utxoHash, uint32(v.TxOutputN))
		txIn := wire.NewTxIn(utxo, nil, nil)
		txIn.Sequence = 0xfffffffe

		// pkScripts, err := hex.DecodeString(v.Script)
		// if err != nil {
		// 	return err
		// }
		// script, err := txscript.SignatureScript(txb, len(txb.TxIn), pkScripts, txscript.SigHashAll, set.PrivKey, false)
		// if err != nil {
		// 	return err
		// }
		// txIn.SignatureScript = script
		// fmt.Println(script)
		txb.AddTxIn(txIn)
	}

	// add payer utxos
	for _, v := range payerUtxos {
		utxoHash, _ := chainhash.NewHashFromStr(v.TxHashBigEndian)
		utxo := wire.NewOutPoint(utxoHash, uint32(v.TxOutputN))
		txIn := wire.NewTxIn(utxo, nil, nil)
		txIn.Sequence = 0xfffffffe

		// pkScripts, err := hex.DecodeString(v.Script)
		// if err != nil {
		// 	return err
		// }
		// script, err := txscript.SignatureScript(txb, len(txb.TxIn), pkScripts, txscript.SigHashAll, setPayer.PrivKey, false)
		// if err != nil {
		// 	return err
		// }
		// txIn.SignatureScript = script

		txb.AddTxIn(txIn)
	}

	{
		sourceAddress, err := btcutil.DecodeAddress(payerFromAddr, &chaincfg.MainNetParams)
		if err != nil {
			return err
		}

		fmt.Println(sourceAddress.EncodeAddress(), sourceAddress.String(), sourceAddress.ScriptAddress())

		pkScript, err := txscript.PayToAddrScript(sourceAddress)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(int64(payerChangeValue), pkScript)
		txb.AddTxOut(txOut)
	}

	var usdtInfo = fmt.Sprintf("%s%s%s%s", "6f6d6e69", "0000", "00000000001f", usdtValue)
	if data, err := hex.DecodeString(usdtInfo); err != nil {
		return err
	} else {
		sourcePkScript, err := txscript.NullDataScript(data)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(0, sourcePkScript)
		txb.AddTxOut(txOut)
	}

	{
		sourceAddress, err := btcutil.DecodeAddress(toAddress, &chaincfg.MainNetParams)
		if err != nil {
			return err
		}
		fmt.Println(sourceAddress.EncodeAddress(), sourceAddress.String(), sourceAddress.ScriptAddress())
		pkScript, err := txscript.PayToAddrScript(sourceAddress)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(int64(fundValue), pkScript)
		txb.AddTxOut(txOut)
	}

	{
		sourceAddress, err := btcutil.DecodeAddress(fromAddr, &chaincfg.MainNetParams)
		if err != nil {
			return err
		}
		pkScript, err := txscript.PayToAddrScript(sourceAddress)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(int64(totalUnspent), pkScript)
		txb.AddTxOut(txOut)
	}

	for i, v := range utxos {
		script, err := hex.DecodeString(v.Script)
		if err != nil {
			return err
		}
		// script = txb.TxOut[i].PkScript
		fmt.Println(script)

		fmt.Println(set.PrivKey.Serialize())
		fmt.Println(set.PrivKey.PubKey().SerializeCompressed())
		fmt.Println(set.PrivKey.PubKey().SerializeUncompressed())
		// sigScript, err := txscript.RawTxInSignature(txb, i, script, txscript.SigHashAll, set.PrivKey)
		// sigScript, err := txscript.SignatureScript(txb, i, script, txscript.SigHashAll, set.PrivKey, true)
		sigScript, err := txscript.SignatureScript(txb, i, script, txscript.SigHashAll, set.PrivKey, set.CompressPubKey)
		fmt.Println(sigScript)
		fmt.Println(txb.TxIn[i].Sequence)
		if err != nil {
			return err
		}
		txb.TxIn[i].SignatureScript = sigScript
	}
	offset := len(utxos)
	for i, v := range payerUtxos {
		script, err := hex.DecodeString(v.Script)
		if err != nil {
			return err
		}
		// script = txb.TxOut[i+offset].PkScript
		// sigScript, err := txscript.RawTxInSignature(txb, i+offset, script, txscript.SigHashAll, setPayer.PrivKey)
		sigScript, err := txscript.SignatureScript(txb, i+offset, script, txscript.SigHashAll, setPayer.PrivKey, setPayer.CompressPubKey)
		// sigScript, err := txscript.SignatureScript(txb, i+offset, script, txscript.SigHashAll, setPayer.PrivKey, false)
		if err != nil {
			return err
		}
		txb.TxIn[i+offset].SignatureScript = sigScript
	}

	var signedTx bytes.Buffer
	txb.Serialize(&signedTx)
	fmt.Println(hex.EncodeToString(signedTx.Bytes()))

	fmt.Println("txid", txb.TxHash().String())
	return nil
}

func TestSafeRun() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestSafeRun", r)
		}
	}()

	go common.SafeRun(func() {
		var a int
		fmt.Println(1 / a)
	})()
}

func PayFrom3Addr(msg Msg) (err error) {
	defer func() {
		if err != nil {
			panic(err)
		}
	}()
	payerPriKey, payerUtxos, payerFromAddr, priKey, utxos, fromAddr, feeValue, usdtValue, toAddress :=
		msg.PayerPriKey, msg.PayerUnspentOutputs, msg.PayerAddress, msg.PriKey, msg.UnspentOutputs, msg.From, msg.FeeValue, msg.UsdtValue, msg.To

	txb := wire.NewMsgTx(TxVersion)
	// 平台支付账户私钥
	setPayer, err := btcutil.DecodeWIF(payerPriKey)
	if err != nil {
		return err
	}
	// sender 账户私钥
	set, err := btcutil.DecodeWIF(priKey)
	if err != nil {
		return err
	}
	var (
		fundValue, payerTotalUnspent, totalUnspent, payerChangeValue float64 = 546, 0, 0, 0
		fromAddress, payerAddress                                    btcutil.Address
	)

	for _, v := range payerUtxos {
		payerTotalUnspent += v.Value
	}

	for _, v := range utxos {
		totalUnspent += v.Value
	}

	if payerTotalUnspent < float64(feeValue*1e8)+fundValue {
		return errors.New(fmt.Sprintf("Total %v less than feeValue*1e8 %v + fundValue %v", payerTotalUnspent, feeValue*1e8, fundValue))
	}

	// payer's change
	payerChangeValue = payerTotalUnspent - fundValue - float64(feeValue*1e8)
	fmt.Println(float64(feeValue*1e8) + fundValue)

	// add sender utxos
	for _, v := range utxos {
		utxoHash, _ := chainhash.NewHashFromStr(v.TxHashBigEndian)
		utxo := wire.NewOutPoint(utxoHash, uint32(v.TxOutputN))
		txIn := wire.NewTxIn(utxo, nil, nil)
		txIn.Sequence = 0xfffffffe
		txb.AddTxIn(txIn)
	}

	// add payer utxos
	for _, v := range payerUtxos {
		utxoHash, _ := chainhash.NewHashFromStr(v.TxHashBigEndian)
		utxo := wire.NewOutPoint(utxoHash, uint32(v.TxOutputN))
		txIn := wire.NewTxIn(utxo, nil, nil)
		txIn.Sequence = 0xfffffffe
		txb.AddTxIn(txIn)
	}

	payerAddress, err = btcutil.DecodeAddress(payerFromAddr, &chaincfg.MainNetParams)
	if err != nil {
		return err
	} else {

		pkScript, err := txscript.PayToAddrScript(payerAddress)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(int64(payerChangeValue), pkScript)
		txb.AddTxOut(txOut)
	}

	var usdtInfo = fmt.Sprintf("%s%s%s%s", "6f6d6e69", "0000", "00000000001f", usdtValue)
	if data, err := hex.DecodeString(usdtInfo); err != nil {
		return err
	} else {
		sourcePkScript, err := txscript.NullDataScript(data)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(0, sourcePkScript)
		txb.AddTxOut(txOut)
	}

	{
		sourceAddress, err := btcutil.DecodeAddress(toAddress, &chaincfg.MainNetParams)
		if err != nil {
			return err
		}
		fmt.Println(sourceAddress.EncodeAddress(), sourceAddress.String(), sourceAddress.ScriptAddress())
		pkScript, err := txscript.PayToAddrScript(sourceAddress)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(int64(fundValue), pkScript)
		txb.AddTxOut(txOut)
	}

	fromAddress, err = btcutil.DecodeAddress(fromAddr, &chaincfg.MainNetParams)
	if err != nil {
		return err
	} else {
		pkScript, err := txscript.PayToAddrScript(fromAddress)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(int64(totalUnspent), pkScript)
		txb.AddTxOut(txOut)
	}

	for i, v := range utxos {
		switch {
		case fromAddr[0:1] == "3":
			p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(
				btcutil.Hash160(set.SerializePubKey()), &chaincfg.MainNetParams,
			)
			if err != nil {
				return err
			}

			var witnessProgram []byte

			witnessProgram, err = txscript.PayToAddrScript(p2wkhAddr)
			if err != nil {
				return err
			}

			sigScript, err := txscript.NewScriptBuilder().AddData(witnessProgram).Script()
			if err != nil {
				return err
			}

			witnessScript, err := txscript.WitnessSignature(txb, txscript.NewTxSigHashes(txb),
				i, int64(v.Value), witnessProgram,
				txscript.SigHashAll, set.PrivKey, set.CompressPubKey,
			)
			if err != nil {
				return err
			}

			txb.TxIn[i].SignatureScript = sigScript
			txb.TxIn[i].Witness = witnessScript
		default:
			script, err := hex.DecodeString(v.Script)
			if err != nil {
				return err
			}

			sigScript, err := txscript.SignatureScript(txb, i, script, txscript.SigHashAll, set.PrivKey, set.CompressPubKey)
			if err != nil {
				return err
			}

			txb.TxIn[i].SignatureScript = sigScript
		}

	}
	offset := len(utxos)
	for i, v := range payerUtxos {
		switch {
		case payerFromAddr[0:1] == "3":
			p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(
				btcutil.Hash160(setPayer.SerializePubKey()), &chaincfg.MainNetParams,
			)
			if err != nil {
				return err
			}

			var witnessProgram []byte

			witnessProgram, err = txscript.PayToAddrScript(p2wkhAddr)
			if err != nil {
				return err
			}

			sigScript, err := txscript.NewScriptBuilder().AddData(witnessProgram).Script()
			if err != nil {
				return err
			}

			witnessScript, err := txscript.WitnessSignature(txb, txscript.NewTxSigHashes(txb),
				i+offset, int64(v.Value), witnessProgram,
				txscript.SigHashAll, setPayer.PrivKey, setPayer.CompressPubKey,
			)
			if err != nil {
				return err
			}

			txb.TxIn[i+offset].SignatureScript = sigScript
			txb.TxIn[i+offset].Witness = witnessScript
		default:
			script, err := hex.DecodeString(v.Script)
			if err != nil {
				return err
			}
			sigScript, err := txscript.SignatureScript(txb, i+offset, script, txscript.SigHashAll, setPayer.PrivKey, setPayer.CompressPubKey)
			if err != nil {
				return err
			}
			txb.TxIn[i+offset].SignatureScript = sigScript
		}

	}

	var signedTx bytes.Buffer
	txb.Serialize(&signedTx)
	fmt.Println(hex.EncodeToString(signedTx.Bytes()))

	fmt.Println("txid", txb.TxHash().String())
	return nil
}
