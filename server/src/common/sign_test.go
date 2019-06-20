package common

import (
	"testing"
)

func TestSignGenerator_GenerateSalt(t *testing.T) {
	v := AppSignMgr.GenerateSalt()
	if (len(v) <= 0) {
		t.Fatalf("generator salt err")
	}
	for i:=0; i<10 ; i++ {
		v1 := AppSignMgr.GenerateSalt()
		if (v1 == v) {
			t.Fatalf("generator salt equal err")
		}
		v = v1
	}
}

func TestSignGenerator_GenerateSource(t *testing.T) {
	//v := AppSignMgr.GenerateSalt()
	mv := map[string]string{
		"phone": "phone",
		"uid": "101",
	}
	v := AppSignMgr.GenerateSource(mv, NowUint32())

	if (v== "") {
		t.Fatalf("generator Source err")
	}
	LogFuncInfo("source: %s", v)

}

func TestSignGenerator_GenerateMSign(t *testing.T) {
	//v := AppSignMgr.GenerateSalt()
	v := "926009653746449"
	mv := map[string]string{
		"uid": "166717822901157888",
	}
	nowt := uint32(1554868412)
	s, e := AppSignMgr.GenerateMSign(mv, nowt, v)

	if (e != nil) {
		t.Fatalf("generator Source err")
	}
	LogFuncInfo("GenerateMSign: salt: %s, error: %v, sign: %s, time: %d", v, e, s, nowt)
}
