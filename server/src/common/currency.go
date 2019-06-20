package common

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"strconv"
	"strings"
)

const (
	CURRENCY_SEP = "."
)

type AmountStruct struct {
	Integer  int32
	Decimals int32
}

var (
	ErrEncodeCurrencyStrNil           = errors.New("input string is \"\"")
	ErrEncodeCurrencyInvalidSepNum    = errors.New("number of \".\" no equals to 1")
	ErrEncodeCurrencyInvalidSep       = errors.New("invalid currency string, \".\" is prefix or subfix.")
	ErrEncodeCurrencyInvalidSubStrNum = errors.New("invalid currency sub string number.")
	ErrNoEnoughCurrency               = errors.New("no enough currency.")
)

// @Description get the sum of amount list
func SumCurrency(list []AmountStruct) (int32, int32) {
	var sum int64
	for _, v := range list {
		sum += int64(v.Integer)*int64(MaxDecimals) + int64(v.Decimals)
	}

	return Int64ToAmount(sum)
}

// @Description encode currency
func DecodeCurrency(s string, precision int) (valueInteger, valueDecimals int32, err error) {
	if !checkPrecision(precision) {
		err = ErrInvalidPrecision
		return
	}

	if len(s) <= 0 {
		err = ErrEncodeCurrencyStrNil
		return
	}

	//only one "." .
	if strings.Count(s, CURRENCY_SEP) != 1 {
		err = ErrEncodeCurrencyInvalidSepNum
		return
	}

	// "." should not be prefix or suffix.
	if strings.HasPrefix(s, CURRENCY_SEP) || strings.HasSuffix(s, CURRENCY_SEP) {
		err = ErrEncodeCurrencyInvalidSep
		return
	}

	//
	subStrs := strings.Split(s, CURRENCY_SEP)
	if len(subStrs) != 2 {
		err = ErrEncodeCurrencyInvalidSubStrNum
		return
	}
	var tmp int
	tmp, err = strconv.Atoi(subStrs[0])
	if err != nil {
		return
	}

	valueInteger = int32(tmp)

	//adjust decimals.
	if len(subStrs[1]) != precision {
		err = ErrNoFixPrecision
		return
	}

	subStrs[1] = adjustDecimalsStr(subStrs[1])
	tmp, err = strconv.Atoi(subStrs[1])
	if err != nil {
		return
	}
	valueDecimals = int32(tmp)

	return
}

func adjustDecimalsStr(s string) (adjustedS string) {
	length := len(s)
	switch {
	case length > 8: // cut to 8 bytes
		{
			adjustedS = string([]byte(s)[:8])
		}
	case length < 8: //fill 0
		{
			tmp := []byte(s)
			for i := length; i < 8; i++ {
				tmp = append(tmp, '0')
			}
			adjustedS = string(tmp)
		}
	default:
		{
			adjustedS = s
			//do nothing
		}
	}

	return
}

func adjustDecimalsInteger(decimals int32) (adjustedS string) {
	s := fmt.Sprint(decimals)
	length := len(s)
	switch {
	case length > 8: // cut to tail 8 bytes
		{
			adjustedS = string([]byte(s)[length-8 : length])
		}
	case length < 8: //fill 0
		{
			adjustedS = s
			for i := 0; i < 8-length; i++ {
				adjustedS = "0" + adjustedS
			}
		}
	default:
		{
			adjustedS = s
			//do nothing
		}
	}

	LogFuncDebug("before adjust %s, after adjust %s", s, adjustedS)

	return
}

// @Description decode currency no care of precision
func DecodeCurrencyNoCarePrecision(s string) (valueInteger, valueDecimals int32, err error) {
	s = strings.Replace(s, " ", "", -1)

	// just an empty string
	if len(strings.TrimSpace(s)) <= 0 {
		return
	}

	//only one "." .
	if strings.Count(s, CURRENCY_SEP) > 1 {
		err = ErrEncodeCurrencyInvalidSepNum
		return
	}

	pos := strings.Index(s, CURRENCY_SEP)

	switch {
	case pos == 0: // eg : ".1234"
		{
			sDecimals := string([]byte(s)[1:])
			sDecimals = adjustDecimalsStr(sDecimals)
			var tmp int
			tmp, err = strconv.Atoi(sDecimals)
			if err != nil {
				return
			}
			valueDecimals = int32(tmp)
		}
	case pos == -1: // eg : "1234"
		{
			var tmp int
			tmp, err = strconv.Atoi(s)
			if err != nil {
				return
			}
			valueInteger = int32(tmp)

		}
	case pos > 0: // eg : "1.234", "123.4", "1234."
		{
			var tmp int
			subStrs := strings.Split(s, CURRENCY_SEP)
			if len(subStrs[0]) > 0 {
				tmp, err = strconv.Atoi(subStrs[0])
				if err != nil {
					return
				}
				valueInteger = int32(tmp)
			}

			if len(subStrs[1]) > 0 {
				subStrs[1] = adjustDecimalsStr(subStrs[1])
				tmp, err = strconv.Atoi(subStrs[1])
				if err != nil {
					return
				}
				valueDecimals = int32(tmp)
			}
		}
	default:
		//
		err = fmt.Errorf("no ganna happen")
		return
	}

	if valueInteger < 0 {
		valueDecimals = -1 * valueDecimals
	}

	return
}

var (
	ErrInvalidPrecision = errors.New("invalid precision")
	ErrNoFixPrecision   = errors.New("input parameter no fix precision")
)

// @Description decode currency
func EncodeCurrency(valueInteger, valueDecimals int32, precision int) (s string, err error) {
	if !checkPrecision(precision) {
		err = ErrInvalidPrecision
		return
	}

	if !checkDecimals(valueDecimals) {
		err = ErrNoFixPrecision
		return
	}

	// remove the '-' if smaller than 0.
	if valueDecimals < 0 {
		valueDecimals *= -1
	}

	tmp := []byte(adjustDecimalsInteger(valueDecimals))
	s = fmt.Sprint(valueInteger) + "." + string(tmp[:precision])

	return
}

// @Description encode currency
func EncodeCurrencyFLoat64(valueInteger, valueDecimals int32, precision int) (amount float64, err error) {
	if !checkPrecision(precision) {
		err = ErrInvalidPrecision
		return
	}

	if !checkDecimals(valueDecimals) {
		err = ErrNoFixPrecision
		return
	}

	amount = float64(valueInteger) + float64(valueDecimals)/MaxDecimals

	return
}

// @Description encode currency
func EncodeCurrencyFLoat642(valueInteger, valueDecimals int32) (amount float64) {
	amount = float64(valueInteger) + float64(valueDecimals)/MaxDecimals
	return
}

// @Description encode currency
func DecodeCurrencyFLoat64(amount float64) (valueInteger, valueDecimals int32) {
	valueInteger = int32(amount)
	valueDecimals = int32((amount - float64(valueInteger)) * MaxDecimals)
	return
}

const (
	CURRENCY_BIGGER  = 1
	CURRENCY_SMALLER = -1
	CURRENCY_EQUAL   = 0
)

const (
	MaxDecimals = 100000000 // max decimals is 100 million
)

// @Description compare currency.
// left > right , return 1
// left < right , return -1
// left == right , return 0
func CompareCurrency(leftInteger, leftDecimals, rightInteger, rightDecimals int32, precision int) (result int, err error) {
	var left, right int64

	left, right, err = checkAndGetLeftRight(leftInteger, leftDecimals, rightInteger, rightDecimals, precision)
	if err != nil {
		return
	}

	switch {
	case left > right:
		result = CURRENCY_BIGGER
	case left < right:
		result = CURRENCY_SMALLER
	default:
		result = CURRENCY_EQUAL
	}

	return
}

// @Description compare currency no check precision
// left > right , return 1
// left < right , return -1
// left == right , return 0
func CompareCurrency2(leftInteger, leftDecimals, rightInteger, rightDecimals int32) (result int) {
	left, right := getLeftRight(leftInteger, leftDecimals, rightInteger, rightDecimals)
	switch {
	case left > right:
		result = CURRENCY_BIGGER
	case left < right:
		result = CURRENCY_SMALLER
	default:
		result = CURRENCY_EQUAL
	}
	return
}

//@Descirption if currency is zero
func ZeroCurrency(integer, decimals int32) bool {
	if AmountToInt64(integer, decimals) <= 0 {
		return true
	}

	return false
}

func AddCurrency(leftInteger, leftDecimals, rightInteger, rightDecimals int32, precision int) (integer, decimals int32, err error) {
	var left, right int64
	left, right, err = checkAndGetLeftRight(leftInteger, leftDecimals, rightInteger, rightDecimals, precision)
	if err != nil {
		return
	}

	sum := left + right

	integer, decimals = Int64ToAmount(sum)

	return
}

// no check precision
func AddCurrency2(leftInteger, leftDecimals, rightInteger, rightDecimals int32) (integer, decimals int32) {
	left, right := getLeftRight(leftInteger, leftDecimals, rightInteger, rightDecimals)
	return Int64ToAmount(left + right)
}

func SubCurrency(leftInteger, leftDecimals, rightInteger, rightDecimals int32, precision int) (integer, decimals int32, err error) {
	var result int
	result, err = CompareCurrency(leftInteger, leftDecimals, rightInteger, rightDecimals, precision)
	if err != nil {
		return
	}

	switch result {
	case CURRENCY_BIGGER:
		{

		}
	case CURRENCY_EQUAL:
		{
			//no error, and integer, decimals are 0.
			return
		}
	default:
		err = ErrNoEnoughCurrency
		return
	}

	left, right := getLeftRight(leftInteger, leftDecimals, rightInteger, rightDecimals)

	sub := left - right

	integer, decimals = Int64ToAmount(sub)

	return
}

func SubCurrency2(leftInteger, leftDecimals, rightInteger, rightDecimals int32) (integer, decimals int32) {
	left, right := getLeftRight(leftInteger, leftDecimals, rightInteger, rightDecimals)
	integer, decimals = Int64ToAmount(left - right)
	return
}

func DivideCurrency(integerIn, decimalsIn int32, divisor float64) (int32, int32) {
	return Int64ToAmount(int64(float64(AmountToInt64(integerIn, decimalsIn)) / divisor))
}

func MultiplyCurrency(integerIn, decimalsIn int32, multiplier float64) (int32, int32) {
	return Int64ToAmount(int64(float64(AmountToInt64(integerIn, decimalsIn)) * multiplier))
}

func checkAndGetLeftRight(leftInteger, leftDecimals, rightInteger, rightDecimals int32, precision int) (left, right int64, err error) {
	if !checkPrecision(precision) {
		err = ErrInvalidPrecision
		return
	}

	if !checkDecimals(leftDecimals) || !checkDecimals(rightDecimals) {
		err = ErrNoFixPrecision
		return
	}

	left, right = getLeftRight(leftInteger, leftDecimals, rightInteger, rightDecimals)

	return
}

func getLeftRight(leftInteger, leftDecimals, rightInteger, rightDecimals int32) (int64, int64) {
	return AmountToInt64(leftInteger, leftDecimals), AmountToInt64(rightInteger, rightDecimals)
}

func AmountToInt64(integer, decimals int32) int64 {
	return int64(integer)*int64(MaxDecimals) + int64(decimals)
}

func Int64ToAmount(amount int64) (integer, decimals int32) {
	integer = int32(amount / MaxDecimals)
	decimals = int32(amount - int64(integer*MaxDecimals))
	return
}

// make sure no overflow
func checkPrecision(precision int) bool {
	return math.Pow10(precision) <= MaxDecimals
}

// decimals must in [0, 100000000)
func checkDecimals(decimals int32) bool {
	return decimals >= 0 && decimals < MaxDecimals
}

//
func CurrencyInt64ToStr(integer int64, precision int) (s string) {
	return CurrencyInt64ToDecimal(integer).StringFixed(int32(precision))
}

func CurrencyInt64ToDecimal(integer int64) (d decimal.Decimal) {
	return decimal.New(integer, 0).Div(decimal.New(MaxDecimals, 0))
}

func CurrencyInt64ToFloat64(integer int64) (f float64) {
	f, _ = CurrencyInt64ToDecimal(integer).Float64()
	return
}

func CurrencyStrToInt64(s string) (integer int64, err error) {
	var d decimal.Decimal
	d, err = decimal.NewFromString(s)
	if err != nil {
		return
	}
	integer = d.Mul(decimal.NewFromFloat(float64(MaxDecimals))).IntPart()
	return
}

func CurrencyFloat64ToInt64(f float64) (integer int64) {
	return decimal.NewFromFloat(f).Mul(decimal.NewFromFloat(float64(MaxDecimals))).IntPart()
}

func GameEusd2GameCoin(eusd float64, ExchangeRate, Precision int32) (coin float64) {
	if ExchangeRate == 0 || Precision == 0 {
		return 0
	}
	return eusd * float64(ExchangeRate) / float64(math.Pow10(int(Precision)))
}

func GameGameCoin2Eusd(coin float64, ExchangeRate, Precision int32) (eusd float64) {
	if ExchangeRate == 0 || Precision == 0 {
		return 0
	}
	return coin * float64(math.Pow10(int(Precision))) / float64(ExchangeRate)
}
