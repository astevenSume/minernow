package offline_tx_bk

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type Transaction struct {
	TxId               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	UnsignedTx         string `json:"unsignedtx"`
	SignedTx           string `json:"signedtx"`
}

// Create Transaction offline
func CreateTransaction(priKey string, destAddr string, amount int64, txHash string) (transaction Transaction, err error) {
	//get wif private key string
	wif, err := btcutil.DecodeWIF(priKey)
	if err != nil {
		return
	}

	//obtain public key and
	var addresspubkey *btcutil.AddressPubKey
	addresspubkey, err = btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.MainNetParams)
	if err != nil {
		return
	}

	sourceTx := wire.NewMsgTx(wire.TxVersion)
	sourceUtxoHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return
	}

	sourceUtxo := wire.NewOutPoint(sourceUtxoHash, 0)
	sourceTxIn := wire.NewTxIn(sourceUtxo, nil, nil)

	//
	destinationAddress, err := btcutil.DecodeAddress(destAddr, &chaincfg.MainNetParams)
	if err != nil {
		return
	}
	sourceAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), &chaincfg.MainNetParams)
	if err != nil {
		return
	}
	destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		return
	}
	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		return
	}

	//generate new source tx out
	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)
	sourceTx.AddTxIn(sourceTxIn)
	sourceTx.AddTxOut(sourceTxOut)
	sourceTxHash := sourceTx.TxHash()

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&sourceTxHash, 0)
	redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(redeemTxIn)
	redeemTxOut := wire.NewTxOut(amount, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)
	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return
	}
	redeemTx.TxIn[0].SignatureScript = sigScript
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)
	if err != nil {
		return
	}
	if err = vm.Execute(); err != nil {
		return
	}
	var unsignedTx bytes.Buffer
	var signedTx bytes.Buffer
	sourceTx.Serialize(&unsignedTx)
	redeemTx.Serialize(&signedTx)
	transaction.TxId = sourceTxHash.String()
	transaction.UnsignedTx = hex.EncodeToString(unsignedTx.Bytes())
	transaction.Amount = amount
	transaction.SignedTx = hex.EncodeToString(signedTx.Bytes())
	transaction.SourceAddress = sourceAddress.EncodeAddress()
	transaction.DestinationAddress = destinationAddress.EncodeAddress()
	return
}

// Create offline signed transaction
// from : sender address
// to : recipient address
// amount : amount to transfer
// fee : fee for
//func CreateTransaction2(to, priKey string, amount int64, fee uint64, propertyId int, utxos []Utxo) (transaction Transaction, err error) {
//	//
//	wif, err := btcutil.DecodeWIF(priKey)
//	if err != nil {
//		return
//	}
//
//	pubKey, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.MainNetParams)
//
//	sourceTx := wire.NewMsgTx(wire.TxVersion)
//	for _, utxo := range utxos {
//		sourceUtxoHash, _ := chainhash.NewHashFromStr(utxo.TxHash)
//		sourceUtxo := wire.NewOutPoint(sourceUtxoHash, 0)
//		sourceTxIn := wire.NewTxIn(sourceUtxo, nil, nil)
//		sourceTx.AddTxIn(sourceTxIn)
//	}
//
//	toAddress, err := btcutil.DecodeAddress(to, &chaincfg.MainNetParams)
//	fromAddress, err := btcutil.DecodeAddress(pubKey.EncodeAddress(), &chaincfg.MainNetParams)
//	if err != nil {
//		return
//	}
//	toPkScript, _ := txscript.PayToAddrScript(toAddress)
//	sourcePkScript, _ := txscript.PayToAddrScript(fromAddress)
//	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)
//	sourceTx.AddTxOut(sourceTxOut)
//	sourceTxHash := sourceTx.TxHash()
//	redeemTx := wire.NewMsgTx(wire.TxVersion)
//	prevOut := wire.NewOutPoint(&sourceTxHash, 0)
//	redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
//	redeemTx.AddTxIn(redeemTxIn)
//	redeemTxOut := wire.NewTxOut(amount, toPkScript)
//	redeemTx.AddTxOut(redeemTxOut)
//	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)
//	if err != nil {
//		return
//	}
//	redeemTx.TxIn[0].SignatureScript = sigScript
//	flags := txscript.StandardVerifyFlags
//	vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)
//	if err != nil {
//		return
//	}
//	if err = vm.Execute(); err != nil {
//		return
//	}
//	var unsignedTx bytes.Buffer
//	var signedTx bytes.Buffer
//	sourceTx.Serialize(&unsignedTx)
//	redeemTx.Serialize(&signedTx)
//	transaction.TxId = sourceTxHash.String()
//	transaction.UnsignedTx = hex.EncodeToString(unsignedTx.Bytes())
//	transaction.Amount = amount
//	transaction.SignedTx = hex.EncodeToString(signedTx.Bytes())
//	transaction.SourceAddress = fromAddress.EncodeAddress()
//	transaction.DestinationAddress = toAddress.EncodeAddress()
//	return
//}
//
//func CreateTransaction3(from string) {
//
//}
