package gameapi

import (
	"fmt"
	"strconv"
	"testing"
)

func TestRgLogin(t *testing.T) {
	api, _ := NewRoyalGameAPI("http://www.devbj.com/api", "qdd_", "2cac7adc65935858c29a8488d6b7aecb")
	res, err := api.Login("hello1", "2cac7adc65935858c29a8488d6b7aecb", "", "506")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("url", res.URL)
}

func TestRgTransferIn(t *testing.T) {
	api, _ := NewRoyalGameAPI("http://www.devbj.com/api", "qdd_", "2cac7adc65935858c29a8488d6b7aecb")
	res, err := api.TransferIn("hello1", "2cac7adc65935858c29a8488d6b7aecb", "asdasd", 12.34)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("Order:", res.Order)
}

func TestRgLotteryList(t *testing.T) {
	/*api, _ := NewRoyalGameAPI("http://www.devbj.com/api", "qdd_", "2cac7adc65935858c29a8488d6b7aecb")

	endTime := common.GetZeroTime(1558540800)
	startTime := endTime - common.DaySeconds
	res, err := api.dayLotteryRecords(1, 10, common.GetTimeFormat(startTime), common.GetTimeFormat(endTime))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("res:", res)*/
}

func TestAgRegister(t *testing.T) {
	api := NewAsiaGamingAPI("http://203.78.143.173:35842", "QianDuoDuo", "f8bf5527577b41c08f4711cfae323016")
	err := api.Register("hello1", "2cac7adc65935858c29a8488d6b7aecb", "")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}

func TestAgTransferIn(t *testing.T) {
	api := NewAsiaGamingAPI("http://203.78.143.173:35842", "QianDuoDuo", "f8bf5527577b41c08f4711cfae323016")
	res, err := api.TransferIn("hello1", "2cac7adc65935858c29a8488d6b7aecb", "", 50)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("Order:", res.Order)
}

func TestAgLogin(t *testing.T) {
	api := NewAsiaGamingAPI("http://203.78.143.203:8004", "QianDuoDuo", "f8bf5527577b41c08f4711cfae323016")
	res, err := api.Login("hello1", "2cac7adc65935858c29a8488d6b7aecb", "", "21")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("url", res.URL)
}

func TestKyLogin(t *testing.T) {
	api := NewKaiYuanAPI(
		"https://kyapi.ky206.com:189",
		"https://kyrecord.ky206.com:190",
		"71038",
		"BEFEFD9D87B44B8A",
		"A5C4E767559AE9F8",
		"test001")
	res, err := api.Login("kaka1", "2cac7adc65935858c29a8488d6b7aecb", "117.30.55.100", "600")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("url", res.URL)
}

func TestKyTransferIn(t *testing.T) {
	f, err := strconv.ParseFloat("-1.90", 64)
	if err != nil {
		return
	}
	fmt.Println(f)
	api := NewKaiYuanAPI(
		"https://kyapi.ky206.com:189",
		"https://kyrecord.ky206.com:190",
		"71038",
		"BEFEFD9D87B44B8A",
		"A5C4E767559AE9F8",
		"test001")
	res, err := api.TransferIn("kaka1", "2cac7adc65935858c29a8488d6b7aecb", "", 20)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("Order:", res.Order)
}

func TestKyGetBetRecord(t *testing.T) {
	/*api := NewKaiYuanAPI(
		"https://kyapi.ky206.com:189",
		"https://kyrecord.ky206.com:190",
		"71038",
		"BEFEFD9D87B44B8A",
		"A5C4E767559AE9F8",
		"test001")
	_, err := api.GetBetRecord(1559295000000)
	if err != nil {
		fmt.Println("error:", err)
		return
	}*/

}
