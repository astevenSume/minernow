package eosplus

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ToJson(arg interface{}) string {
	data, _ := json.Marshal(arg)

	return string(data)
}

func ToJsonIndent(arg interface{}) string {
	data, _ := json.MarshalIndent(arg, "", "  ")

	return string(data)
}

//  exp:quantity = "1234.1234 EOS"
func QuantityToUint64(quantity string) (res uint64, err error) {
	tmp := strings.Split(quantity, " ")
	n, err := strconv.ParseFloat(tmp[0], 32)
	if err != nil {
		return
	}
	res = uint64(n * math.Pow10(EosPrecision))
	return
}

func QuantityFloat64ToUint64(quantity float64) (res uint64) {
	res = uint64(quantity * math.Pow10(EosPrecision))
	return
}

func QuantityFloat64ToInt64(quantity float64) (res int64) {
	res = int64(quantity * math.Pow10(EosPrecision))
	return
}

//数字转字符串(除去精度) 用于显示
func QuantityToString(q uint64) string {
	return fmt.Sprintf("%.4f", float64(q)/math.Pow10(EosPrecision))
}

func QuantityInt64ToString(q int64) string {
	return fmt.Sprintf("%.4f", float64(q)/math.Pow10(EosPrecision))
}

//字符串转数字（乘上精度）用户存储运算
func QuantityStringToUint64(q string) uint64 {
	f, err := strconv.ParseFloat(q, 64)
	if err != nil {
		return 0
	}
	n := f * math.Pow10(EosPrecision)
	return uint64(n)
}

// 返回float64的eusd
func QuantityUint64ToFloat64(q uint64) float64 {
	n := float64(q) / math.Pow10(EosPrecision)
	return n
}

// 返回float64的eusd
func QuantityInt64ToFloat64(q int64) float64 {
	n := float64(q) / math.Pow10(EosPrecision)
	return n
}

//字符串转数字（乘上精度）用户存储运算 int64
func QuantityStringToint64(q string) int64 {
	f, err := strconv.ParseFloat(q, 64)
	if err != nil {
		return 0
	}
	n := f * math.Pow10(EosPrecision)
	return int64(n)
}