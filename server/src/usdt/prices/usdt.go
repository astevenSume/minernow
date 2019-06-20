package prices

import (
	"common"
	"encoding/json"
	"errors"
	"eusd/eosplus"
	"fmt"
	"math"
	common2 "otc/common"
	"time"
	"utils/usdt/dao"
	"utils/usdt/models"

	"github.com/astaxie/beego/httplib"
)

type PRICE_TRADE_TYPE int

const (
	PRICE_TRADE_BUY  PRICE_TRADE_TYPE = 0
	PRICE_TRADE_SELL PRICE_TRADE_TYPE = 1
)

type PRICE_CURRENCY_TYPE int

const (
	PRICE_CURRENCY_TYPE_BTC  PRICE_CURRENCY_TYPE = 1
	PRICE_CURRENCY_TYPE_USDT PRICE_CURRENCY_TYPE = 2
)

var priceCurrencyTypeName = map[PRICE_CURRENCY_TYPE]string{
	PRICE_CURRENCY_TYPE_BTC:  "BTC",
	PRICE_CURRENCY_TYPE_USDT: "USDT",
}

func (c PRICE_CURRENCY_TYPE) CurrencyName() string {
	return priceCurrencyTypeName[c]
}

const (
	PRICE_COMMON_POW = 8
	PRICE_USDT_KEY   = "price_otc_ex_usdt_key"
	PRICE_BTC_KEY    = "price_otc_ex_btc_key"

	PRICE_OKEX_TRADE_API    = "https://www.okex.com/v3/c2c/tradingOrders/book?side=%s&baseCurrency=%s&quoteCurrency=cny&userType=blockTrade&paymentMethod=all"
	PRICE_OKEX_CODE_SUCCESS = 0

	PRICE_HBG_TRADE_API = "https://otc-api.huobipro.com/v1/otc/trade/list/public?currPage=1&payWay=&country=0&payMethod=0&currency=1&merchant=1&online=1&range=0&coinId=%d&tradeType=%d"
	//PRICE_HBG_TRADE_API    = "https://otc-api.hbg.com/v1/data/trade-market?country=37&currency=1&payMethod=0&currPage=1&coinId=2&tradeType=%s&blockType=general&online=1"
	PRICE_HBG_CODE_SUCCESS = 200
)

// 同步otc兑换usdt cny 价格
func SyncUsdtPrice() (models.Prices, error) {
	return syncPrice(PRICE_CURRENCY_TYPE_USDT, PRICE_USDT_KEY, common2.Cursvr.UsdtPriceExpiredSecs, common2.Cursvr.UsdtMinPrice, common2.Cursvr.UsdtMaxPrice)
}

// 同步otc兑换btc cny 价格
func SyncBtcPrice() (models.Prices, error) {
	return syncPrice(PRICE_CURRENCY_TYPE_BTC, PRICE_BTC_KEY, common2.Cursvr.BtcPriceExpiredSecs, common2.Cursvr.BtcMinPrice, common2.Cursvr.BtcMaxPrice)
}

func syncPrice(curType PRICE_CURRENCY_TYPE, key string, expiredSecs int, minPrice, maxPrice float64) (models.Prices, error) {

	var ch chan models.MarketPrices
	ch = make(chan models.MarketPrices, 4)

	go common.SafeRun(func() {
		syncOkexPrice(curType, PRICE_TRADE_BUY, minPrice, maxPrice, ch)
	})()
	go common.SafeRun(func() {
		syncOkexPrice(curType, PRICE_TRADE_SELL, minPrice, maxPrice, ch)
	})()
	go common.SafeRun(func() {
		syncHbgPrice(curType, PRICE_TRADE_BUY, minPrice, maxPrice, ch)
	})()
	go common.SafeRun(func() {
		syncHbgPrice(curType, PRICE_TRADE_SELL, minPrice, maxPrice, ch)
	})()

	priceSum := 0.0
	count := 0
	for i := 0; i < 4; i++ {
		price := <-ch
		if price.Pow != 0 && price.PowPrice > 0 {
			priceSum += float64(price.PowPrice) / float64(math.Pow10(int(price.Pow)))
			count++
		}
	}
	if count > 0 {
		avgPrice := priceSum / float64(count)
		avgPrice = Max(avgPrice, minPrice)
		avgPrice = Min(avgPrice, maxPrice)
		price := models.Prices{
			Currency: curType.CurrencyName(),
			PowPrice: uint64(math.Pow10(PRICE_COMMON_POW) * avgPrice),
			Pow:      PRICE_COMMON_POW,
		}
		if err := dao.PricesDaoEntity.InsertPrice(&price); err != nil {
			common.LogFuncError("insert price err:%v", err)
			return models.Prices{}, err
		} else {
			if b, err := json.Marshal(price); err != nil {
				common.LogFuncError("Marshal failed :%v", err)
			} else {
				if err := common.RedisManger.
					Set(key, b, time.Duration(int64(expiredSecs))*time.Second).
					Err(); err != nil {
					common.LogFuncError("set %s price to redis failed:%v", curType.CurrencyName(), err)
				}
			}

			return price, nil
		}
	} else {
		common.LogFuncError("sycn price err")
	}

	return models.Prices{}, errors.New("Sync price err")
}

// 获取ustd价格
func GetPrice(curType PRICE_CURRENCY_TYPE) (models.Prices, error) {
	price := models.Prices{}
	var redisKey string

	switch curType {
	case PRICE_CURRENCY_TYPE_BTC:
		redisKey = PRICE_BTC_KEY
	case PRICE_CURRENCY_TYPE_USDT:
		redisKey = PRICE_USDT_KEY
	default:
		return price, fmt.Errorf("unknow price currency type : %v", curType)
	}

	if b, err := common.RedisManger.Get(redisKey).Bytes(); err != nil {
		return dao.PricesDaoEntity.QueryLastPrice(curType.CurrencyName())
	} else {
		if err := json.Unmarshal(b, &price); err != nil {
			return dao.PricesDaoEntity.QueryLastPrice(curType.CurrencyName())
		}
	}

	return price, nil
}

func GetPriceFloat64(curType PRICE_CURRENCY_TYPE) float64 {
	price, err := GetPrice(curType)
	if err != nil {
		return 0
	}

	res := float64(price.PowPrice) / float64(math.Pow10(int(price.Pow)))

	return res
}

func GetBtc2Usdt() float64 {

	btc2cny := GetPriceFloat64(PRICE_CURRENCY_TYPE_BTC)
	usdt2cny := GetPriceFloat64(PRICE_CURRENCY_TYPE_USDT)

	if btc2cny == 0 || usdt2cny == 0 {
		return 0
	}

	return btc2cny / usdt2cny
}

// 同步火币网usdt价格
func syncHbgPrice(curType PRICE_CURRENCY_TYPE, tradeType PRICE_TRADE_TYPE, minPrice, maxPrice float64, ch chan models.MarketPrices) {
	type MarketPrice struct {
		Price float64 `json:"price"`
	}
	type MarketHbgRes struct {
		Code int           `json:"code"`
		Data []MarketPrice `json:"data"`
	}
	buyMsg := MarketHbgRes{}
	buyMsg.Code = -1

	var side int // side 0:bid, 1:ask
	switch tradeType {
	case PRICE_TRADE_BUY:
		side = 0
	case PRICE_TRADE_SELL:
		side = 1
	default:
		side = 1
	}

	// 需要代理获取，或将服务器部署到墙外
	if err := httplib.Get(fmt.Sprintf(PRICE_HBG_TRADE_API, int(curType), side)).
		SetUserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.104 Safari/537.36 Core/1.53.4882.400 QQBrowser/9.7.13059.400").
		//SetProxy(func(request *http.Request) (*url.URL, error) {
		//	u := new(url.URL)
		//	u.Scheme = "http"
		//	u.Host = "127.0.0.1:1080"
		//	return u, nil}).
		SetTimeout(time.Second*15, time.Second*10).
		ToJSON(&buyMsg);
	//Response();
	err != nil {
		common.LogFuncError(fmt.Sprintf("err %v, %s", err, eosplus.ToJson(buyMsg)))
	} else {
		common.LogFuncInfo(fmt.Sprintf(" %d, %s", tradeType, eosplus.ToJson(buyMsg)))
		if buyMsg.Code == PRICE_HBG_CODE_SUCCESS {
			var priceTmp = 0.0
			switch tradeType {
			case PRICE_TRADE_BUY:
				for _, v := range buyMsg.Data {
					if v.Price > priceTmp {
						priceTmp = v.Price
					}
				}
			case PRICE_TRADE_SELL:
				if len(buyMsg.Data) > 0 {
					priceTmp = buyMsg.Data[0].Price
				}
				for _, v := range buyMsg.Data {
					if v.Price < priceTmp {
						priceTmp = v.Price
					}
				}
			}

			if priceTmp > 0.0 {
				priceTmp = Max(priceTmp, minPrice)
				priceTmp = Min(priceTmp, maxPrice)
				marketPrice := models.MarketPrices{
					Currency:    curType.CurrencyName(),
					TradeMethod: int8(tradeType),
					Market:      dao.PRICE_MARKET_HBG,
					PowPrice:    uint64(math.Pow10(PRICE_COMMON_POW) * priceTmp),
					Pow:         PRICE_COMMON_POW,
				}
				if err := dao.PricesDaoEntity.InsertMarketPrice(&marketPrice); err != nil {
					common.LogFuncError(fmt.Sprintf("sql err %v", err))
				} else {
					// sync price event
					ch <- marketPrice
					return
				}
			}
		}
	}
	ch <- models.MarketPrices{}

	return
}

// 同步okex网usdt价格
func syncOkexPrice(curType PRICE_CURRENCY_TYPE, tradeType PRICE_TRADE_TYPE, minPrice, maxPrice float64, ch chan models.MarketPrices) {

	type MarketPrice struct {
		Price float64 `json:"price"`
	}
	type Data struct {
		Buy  []MarketPrice `json:"buy"`
		Sell []MarketPrice `json:"sell"`
	}
	type MarketOkexRes struct {
		Code int  `json:"code"`
		Data Data `json:"data"`
	}
	buyMsg := MarketOkexRes{}
	buyMsg.Code = -1

	var side string
	switch tradeType {
	case PRICE_TRADE_BUY:
		side = "buy"
	case PRICE_TRADE_SELL:
		side = "sell"
	default:
		side = "sell"
	}

	var cur string
	switch curType {
	case PRICE_CURRENCY_TYPE_BTC:
		cur = "btc"
	case PRICE_CURRENCY_TYPE_USDT:
		cur = "usdt"
	}
	// 需要代理获取，或将服务器部署到墙外
	if err := httplib.Get(fmt.Sprintf(PRICE_OKEX_TRADE_API, side, cur)).
		//SetProxy(func(request *http.Request) (*url.URL, error) {
		//	u := new(url.URL)
		//	u.Scheme = "socks5"
		//	u.Host = "127.0.0.1:1080"
		//	return u, nil}).
		SetTimeout(time.Second*15, time.Second*10).
		ToJSON(&buyMsg); err != nil {
		common.LogFuncDebug(fmt.Sprintf("err %v, %s", err, eosplus.ToJson(buyMsg)))
	} else {
		common.LogFuncInfo(fmt.Sprintf(" %d, %s", tradeType, eosplus.ToJson(buyMsg)))
		if buyMsg.Code == PRICE_OKEX_CODE_SUCCESS {
			var priceTmp = 0.0
			var list = []MarketPrice{}
			switch tradeType {
			case PRICE_TRADE_BUY:
				list = buyMsg.Data.Buy
				for _, v := range list {
					if v.Price > priceTmp {
						priceTmp = v.Price
					}
				}
			case PRICE_TRADE_SELL:
				list = buyMsg.Data.Sell
				if len(list) > 0 {
					priceTmp = list[0].Price
				}
				for _, v := range list {
					if v.Price < priceTmp {
						priceTmp = v.Price
					}
				}
			}

			if priceTmp > 0.0 {
				priceTmp = Max(priceTmp, minPrice)
				priceTmp = Min(priceTmp, maxPrice)
				marketPrice := models.MarketPrices{
					Currency:    curType.CurrencyName(),
					TradeMethod: int8(tradeType),
					Market:      dao.PRICE_MARKET_OKEX,
					PowPrice:    uint64(math.Pow10(PRICE_COMMON_POW) * priceTmp),
					Pow:         PRICE_COMMON_POW,
				}
				if err := dao.PricesDaoEntity.InsertMarketPrice(&marketPrice); err != nil {
					common.LogFuncError(fmt.Sprintf("sql err %v", err))
				} else {
					ch <- marketPrice
					return
				}
			}
		}
	}

	ch <- models.MarketPrices{}

	return
}

func Min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}
