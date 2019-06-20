package eosplus

import "testing"

func TestEosQuantityToUint64(t *testing.T) {
	quant := "1234.12345 EOS"
	res, _ := QuantityToUint64(quant)
	if res != 12341234 {
		t.Fatalf("result is ERR,%v => %v", quant, res)
	}
}
