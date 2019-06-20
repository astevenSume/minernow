package usdt

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
)

func Test_GetUnspent(t *testing.T) {
	addr1 := "1Ck1XZW82ZHqv69WZtzDZVHzmh6KyfuUDp"
	utxos, err := getUnspent(wire.MainNet, addr1)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(utxos) > 0 {
		t.Fatalf("you think you are Satoshi, right?")
	}

	t.Logf("unspent outputs : %d", len(utxos))
}

func Test_GetUnspent_Satoshi(t *testing.T) {
	addr1 := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	utxos, err := getUnspent(wire.MainNet, addr1)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(utxos) <= 0 {
		t.Fatalf("Satoshi got no btc ?! R U kiding?!")
	}

	t.Logf("unspent outputs : %d", len(utxos))
}

func Test_GetFeeRate(t *testing.T) {
	fastestFee, err := getFee(FeeModeHalfHour)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if fastestFee <= 0 {
		t.Fatalf("fastestFee <= 0")
	}

	t.Logf("fastestFee : %d", fastestFee)
}
