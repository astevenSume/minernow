var express = require('express')
var bitcoin = require('bitcoinjs-lib')

var server = express()

server.listen(4040, 'localhost')

server.post('/', function (req, resp) {
    console.log(req)

    req.rawBody = ''
    req.setEncoding('utf8')
    req.on('data', function (chunk) {
        req.rawBody += chunk
    })

    req.on('end', function () {
        var obj = JSON.parse(req.rawBody)
        // var signedTx = usdtSign(obj.pri_key, obj.unspent_outputs, obj.fee_value, obj.usdt_value, obj.from, obj.to)
        var msg = usdtSignMulti(obj.payer_pri_key, obj.payer_unspent_outputs, obj.payer_address,
            obj.pri_key, obj.unspent_outputs, obj.from,
            obj.fee_value, obj.usdt_value, obj.to)
        resp.send(msg)
    })
})

function addPreZero(num){
    var t = (num+'').length,
        s = ''
    for(var i=0; i<16-t; i++){
        s += '0'
    }
    console.log("s+num : "+(s+num))
    return s+num
}

function usdtSign(privateKey, utxo, feeValue, usdtValue, fromAddress, toAddress) {
    var txb = new bitcoin.TransactionBuilder()
    var set = bitcoin.ECPair.fromWIF(privateKey)
    const fundValue = 546
    var usdtAmount = parseInt(usdtValue*1e8).toString(16)
    var totalUnspent = 0
    for(var i = 0; i < utxo.length; i++){
        totalUnspent = totalUnspent + utxo[i].value
    }
    const changeValue = totalUnspent - fundValue - (feeValue*1e8)
    if (totalUnspent < feeValue + fundValue) {
        console.log("Total less than fee")
        return constant.LessValue
    }
    for(var i = 0; i< utxo.length; i++){
        txb.addInput(utxo[i].tx_hash_big_endian, utxo[i].tx_output_n, 0xfffffffe)
    }
    const usdtInfo = [
        "6f6d6e69",
        "0000",
        "00000000001f",
        addPreZero(usdtAmount)
    ].join('')
    const data = Buffer.from(usdtInfo, "hex")
    const omniOutput = bitcoin.script.compile([
        bitcoin.opcodes.OP_RETURN,
        data
    ])
    txb.addOutput(toAddress, fundValue)
    txb.addOutput(omniOutput, 0)
    txb.addOutput(fromAddress, changeValue)
    for(var i = 0;i < utxo.length; i++){
        txb.sign(0, set)
    }
    return txb.buildIncomplete().toHex()
}

// "payer" means the guy who pays for the fee.
function usdtSignMulti(payerPriKey, payerUtxos, payerFromAddr, priKey, utxos, fromAddr, feeValue, usdtValue, toAddress) {
    console.log(usdtValue)
    // check param
    if (payerPriKey == "") {
        return "{\"code\" : 1, \"msg\" : \"payerPriKey is ''\"}"
    }

    if (payerUtxos.length <= 0) {
        return "{\"code\" : 1, \"msg\" : \"len(payerUtxos) <= 0\"}"
    }

    if (payerFromAddr == "") {
        return "{\"code\" : 1, \"msg\" : \"payerFromAddr is ''\"}"
    }

    if (priKey == "") {
        return "{\"code\" : 1, \"msg\" : \"priKey is ''\"}"
    }

    if (utxos.length <= 0) {
        return "{\"code\" : 1, \"msg\" : \"len(utxos) <= 0\"}"
    }

    if (fromAddr == "") {
        return "{\"code\" : 1, \"msg\" : \"fromAddr is ''\"}"
    }

    if (toAddress == "") {
        return "{\"code\" : 1, \"msg\" : \"toAddress is ''\"}"
    }

    var txb = new bitcoin.TransactionBuilder()
    var setPayer = bitcoin.ECPair.fromWIF(payerPriKey)
    var set = bitcoin.ECPair.fromWIF(priKey)
    const fundValue = 546
    // var usdtAmount = parseInt(usdtValue*1e8).toString(16)
    var payerTotalUnspent = 0
    var totalUnspent = 0
    for(var i = 0; i < payerUtxos.length; i++){
        payerTotalUnspent = payerTotalUnspent + payerUtxos[i].value
    }

    for(var i = 0; i < utxos.length; i++){
        totalUnspent = totalUnspent + utxos[i].value
    }

    if (payerTotalUnspent < feeValue*1e8 + fundValue) {
        console.log("Total less than fee")
        return "{\"code\" : 1, \"msg\" : \"Total "+ payerTotalUnspent + "less than feeValue*1e8 "+ feeValue*1e8 +" + fundValue "+ fundValue +" \"}"
    }

    // payer's change
    const payerChangeValue = payerTotalUnspent - fundValue - (feeValue*1e8)


    // add sender utxos
    for(var i = 0; i< utxos.length; i++){
        txb.addInput(utxos[i].tx_hash_big_endian, utxos[i].tx_output_n, 0xfffffffe)
    }

    // add payer utxos
    for(var i = 0; i< payerUtxos.length; i++){
        txb.addInput(payerUtxos[i].tx_hash_big_endian, payerUtxos[i].tx_output_n, 0xfffffffe)
    }

    const usdtInfo = [
        "6f6d6e69",
        "0000",
        "00000000001f",
        usdtValue
        // addPreZero(usdtAmount)
    ].join('')
    const data = Buffer.from(usdtInfo, "hex")
    const omniOutput = bitcoin.script.compile([
        bitcoin.opcodes.OP_RETURN,
        data
    ])

    // =_=!!! multi OP_RETURN is unexpected. when you send multi OP_RETURN transaction, the net will response with this error :
    // 64: multi-op-return

    // const usdtPayerInfo = [
    //     "6f6d6e69",
    //     "0000",
    //     "00000000001f",
    //     addPreZero(1)
    // ].join('')
    // const dataPayer = Buffer.from(usdtPayerInfo, "hex")
    // const omniOutputPayer = bitcoin.script.compile([
    //     bitcoin.opcodes.OP_RETURN,
    //     dataPayer
    // ])
    //
    // txb.addOutput(omniOutputPayer, 0)//add usdt output

    txb.addOutput(payerFromAddr, payerChangeValue)//add payer btc output
    txb.addOutput(omniOutput, 0)//add usdt output
    txb.addOutput(toAddress, fundValue)
    txb.addOutput(fromAddr, totalUnspent)//add sender btc output

    //sign sender utxos
    for (var i = 0; i < utxos.length; i++){
        txb.sign(i, set)
    }

    //sign payer utxos
    for(var i = utxos.length;i < utxos.length + payerUtxos.length; i++){
        txb.sign(i, setPayer)
    }


    return "{\"code\" : 0, \"msg\" : \"\", \"signed_tx\" : \""+ txb.buildIncomplete().toHex() +"\"}"
}