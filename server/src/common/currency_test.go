package common

/*import "testing"

func TestDecodeCurrency(t *testing.T) {

	var err error
	var testStr string

	var precision = 10
	_, err = EncodeCurrency(0, 0, precision)
	if err != ErrInvalidPrecision {
		t.Fatalf("precision %d(> 9), should return ErrInvalidPrecision", precision)
	}

	testStr = ""
	_, _, err = DecodeCurrency(testStr, 4)
	if err != ErrEncodeCurrencyStrNil {
		t.Fatalf("check \"%s\" failed.", testStr)
	}

	testStr = "1.2.34"
	_, _, err = DecodeCurrency(testStr, 4)
	if err != ErrEncodeCurrencyInvalidSepNum {
		t.Fatalf("check \"%s\" failed.", testStr)
	}

	testStr = "1234"
	_, _, err = DecodeCurrency(testStr, 4)
	if err != ErrEncodeCurrencyInvalidSepNum {
		t.Fatalf("check \"%s\" failed.", testStr)
	}

	testStr = ".1234"
	_, _, err = DecodeCurrency(testStr, 4)
	if err != ErrEncodeCurrencyInvalidSep {
		t.Fatalf("check \"%s\" failed.", testStr)
	}

	testStr = "1234."
	_, _, err = DecodeCurrency(testStr, 4)
	if err != ErrEncodeCurrencyInvalidSep {
		t.Fatalf("check \"%s\" failed.", testStr)
	}

	testStr = "1234.567890"
	_, _, err = DecodeCurrency(testStr, 4)
	if err != ErrNoFixPrecision {
		t.Fatalf("%v", err)
	}

	testStr = "1234.5"
	_, _, err = DecodeCurrency(testStr, 4)
	if err != ErrNoFixPrecision {
		t.Fatalf("%v", err)
	}

	testStr = "1234.5678"
	var (
		v1, v2 int32
	)
	v1, v2, err = DecodeCurrency(testStr, 4)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if v1 != 1234 {
		t.Fatalf("v1 should be 1234")
	}

	if v2 != 56780000 {
		t.Fatalf("v1 should be 56780000")
	}
}

func TestEncodeCurrency(t *testing.T) {
	var (
		s         string
		err       error
		v1, v2    int32
		precision int
	)

	precision = 10
	s, err = EncodeCurrency(v1, v2, precision)
	if err != ErrInvalidPrecision {
		t.Fatalf("precision %d(> 9), should return ErrInvalidPrecision", precision)
	}

	v1, v2, precision = 1234, 5678988, 4
	s, err = EncodeCurrency(v1, v2, precision)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if s != "1234.0567" {
		t.Fatalf("s is %s, expect %s", s, "1234.5678")
	}

	v1, v2, precision = 1234, 567898888, 4
	s, err = EncodeCurrency(v1, v2, precision)
	if err != ErrNoFixPrecision {
		t.Fatalf("%v", err)
	}

	v1, v2, precision = 1234, 0, 4
	s, err = EncodeCurrency(v1, v2, precision)
	if err != nil {
		t.Fatalf("should be no error")
	}

	if s != "1234.0000" {
		t.Fatalf("should be 1234.0000 but %s", s)
	}

	v1, v2, precision = 10, 20, 4
	s, err = EncodeCurrency(v1, v2, precision)
	if err != nil {
		t.Fatalf("should be no error")
	}

	if s != "10.0000" {
		t.Fatalf("should be 10.0020 but %s", s)
	}
}

func TestCompareCurrency(t *testing.T) {
	var precision = 10
	var err error
	_, err = CompareCurrency(0, 0, 0, 0, precision)
	if err != ErrInvalidPrecision {
		t.Fatalf("precision %d(> 9), should return ErrInvalidPrecision", precision)
	}

	var (
		leftV1, leftv2, rightV1, rightV2 int32
		result                           int
	)

	leftV1, leftv2, rightV1, rightV2, precision = 1, 99999999, 1, 99999998, 4
	result, err = CompareCurrency(leftV1, leftv2, rightV1, rightV2, precision)
	if err != nil {
		t.Fatalf("should be no error.")
	}
	if result != CURRENCY_BIGGER {
		t.Fatalf("%d.%d is bigger than %d.%d", leftV1, leftv2, rightV1, rightV2)
	}

	leftV1, leftv2, rightV1, rightV2, precision = 0, 0, 1, 1, 4
	result, err = CompareCurrency(leftV1, leftv2, rightV1, rightV2, precision)
	if err != nil {
		t.Fatalf("should be no error.")
	}
	if result != CURRENCY_SMALLER {
		t.Fatalf("%d.%d is smaller than %d.%d", leftV1, leftv2, rightV1, rightV2)
	}

	leftV1, leftv2, rightV1, rightV2, precision = 1, 0, 1, 0, 4
	result, err = CompareCurrency(leftV1, leftv2, rightV1, rightV2, precision)
	if err != nil {
		t.Fatalf("should be no error.")
	}
	if result != CURRENCY_EQUAL {
		t.Fatalf("%d.%d is equal to %d.%d", leftV1, leftv2, rightV1, rightV2)
	}

	leftV1, leftv2, rightV1, rightV2, precision = 0, 1234, 0, 1232, 4
	result, err = CompareCurrency(leftV1, leftv2, rightV1, rightV2, precision)
	if err != nil {
		t.Fatalf("should be no error.")
	}
	if result != CURRENCY_BIGGER {
		t.Fatalf("%d.%d is bigger than %d.%d", leftV1, leftv2, rightV1, rightV2)
	}

	leftV1, leftv2, rightV1, rightV2, precision = 0, 1234, 0, 0, 4
	result, err = CompareCurrency(leftV1, leftv2, rightV1, rightV2, precision)
	if err != nil {
		t.Fatalf("should be no error.")
	}
	if result != CURRENCY_BIGGER {
		t.Fatalf("%d.%d is bigger than %d.%d", leftV1, leftv2, rightV1, rightV2)
	}

	leftV1, leftv2, rightV1, rightV2, precision = 0, 0, 0, 0, 4
	result, err = CompareCurrency(leftV1, leftv2, rightV1, rightV2, precision)
	if err != nil {
		t.Fatalf("should be no error.")
	}
	if result != CURRENCY_EQUAL {
		t.Fatalf("%d.%d should equal to %d.%d", leftV1, leftv2, rightV1, rightV2)
	}

}

func TestAddCurrency(t *testing.T) {
	//
}

func TestSubCurrency(t *testing.T) {

}

func TestDecodeCurrencyNoCarePrecision(t *testing.T) {
	var (
		err               error
		integer, decimals int32
	)
	//  "1.2.34"
	_, _, err = DecodeCurrencyNoCarePrecision("1.2.34")
	if err != ErrEncodeCurrencyInvalidSepNum {
		t.Fatalf("expect return ErrEncodeCurrencyInvalidSepNum but %v", err)
	}

	//  ""
	integer, decimals, err = DecodeCurrencyNoCarePrecision("")
	if err != nil {
		t.Fatalf("expect no error but %v", err)
	}
	if integer != 0 || decimals != 0 {
		t.Fatalf("expect integer(0), decimal(0) but integer(%d), decimal(%d)", integer, decimals)
	}
	//  "1234"
	integer, decimals, err = DecodeCurrencyNoCarePrecision("1234")
	if err != nil {
		t.Fatalf("expect no error but %v", err)
	}
	if integer != 1234 || decimals != 0 {
		t.Fatalf("expect integer(1234), decimal(0) but integer(%d), decimal(%d)", integer, decimals)
	}

	//	".1234"
	integer, decimals, err = DecodeCurrencyNoCarePrecision(".1234")
	if err != nil {
		t.Fatalf("expect no error but %v", err)
	}
	if integer != 0 || decimals != 12340000 {
		t.Fatalf("expect integer(0), decimal(12340000) but integer(%d), decimal(%d)", integer, decimals)
	}

	//  "1.234"
	integer, decimals, err = DecodeCurrencyNoCarePrecision("1.234")
	if err != nil {
		t.Fatalf("expect no error but %v", err)
	}
	if integer != 1 || decimals != 23400000 {
		t.Fatalf("expect integer(1), decimal(23400000) but integer(%d), decimal(%d)", integer, decimals)
	}

	//  "12.34"
	integer, decimals, err = DecodeCurrencyNoCarePrecision("12.34")
	if err != nil {
		t.Fatalf("expect no error but %v", err)
	}
	if integer != 12 || decimals != 34000000 {
		t.Fatalf("expect integer(12), decimal(34000000) but integer(%d), decimal(%d)", integer, decimals)
	}

	//  "1.234567891"
	integer, decimals, err = DecodeCurrencyNoCarePrecision("1.234567891")
	if err != nil {
		t.Fatalf("expect no error but %v", err)
	}
	if integer != 1 || decimals != 23456789 {
		t.Fatalf("expect integer(1), decimal(23456789) but integer(%d), decimal(%d)", integer, decimals)
	}

	//  "1.23456789"
	integer, decimals, err = DecodeCurrencyNoCarePrecision("1.23456789")
	if err != nil {
		t.Fatalf("expect no error but %v", err)
	}
	if integer != 1 || decimals != 23456789 {
		t.Fatalf("expect integer(1), decimal(23456789) but integer(%d), decimal(%d)", integer, decimals)
	}

	//  "1234."
	integer, decimals, err = DecodeCurrencyNoCarePrecision("1234.")
	if err != nil {
		t.Fatalf("expect no error but %v", err)
	}
	if integer != 1234 || decimals != 0 {
		t.Fatalf("expect integer(1234), decimal(0) but integer(%d), decimal(%d)", integer, decimals)
	}

	// "1 .2 3 4"
	integer, decimals, err = DecodeCurrencyNoCarePrecision("1 .2 3 4")
	if err != nil {
		t.Fatalf("expect no error but %v", err)
	}
	if integer != 1 || decimals != 23400000 {
		t.Fatalf("expect integer(1), decimal(23400000) but integer(%d), decimal(%d)", integer, decimals)
	}

	// " .123"
	integer, decimals, err = DecodeCurrencyNoCarePrecision(" .234")
	if err != nil {
		t.Fatalf("expect no error but %v", err)
	}
	if integer != 0 || decimals != 23400000 {
		t.Fatalf("expect integer(1), decimal(23400000) but integer(%d), decimal(%d)", integer, decimals)
	}
}

func TestDivideCurrency(t *testing.T) {
	integer, decimals := DivideCurrency(10, 20, 2)
	if integer != 5 || decimals != 10 {
		t.Fatalf("should be integer 5, decimals 10")
	}

	integer, decimals = DivideCurrency(10, 20, 0.5)
	if integer != 20 || decimals != 40 {
		t.Fatalf("should be integer 20, decimals 40")
	}
}

func TestMultiplyCurrency(t *testing.T) {
	integer, decimals := MultiplyCurrency(10, 20, 2)
	if integer != 20 || decimals != 40 {
		t.Fatalf("should be integer 20, decimals 40")
	}

	integer, decimals = MultiplyCurrency(10, 20, 0.5)
	if integer != 5 || decimals != 10 {
		t.Fatalf("should be integer 5, decimals 10")
	}
}

func TestEncodeCurrency3(t *testing.T) {
	//
	integer := int64(9223372036750000000)
	precision := 8
	s := DecodeCurrency3Str(integer, precision)
	if s != "92233720367.50000000" {
		t.Fatalf("s should be 92233720367.50000000 but %s", s)
	}

	integer = int64(9223372036744444444)
	precision = 8
	s = DecodeCurrency3Str(integer, precision)
	if s != "92233720367.44444444" {
		t.Fatalf("s should be 92233720367.44444444 but %s", s)
	}
}*/
