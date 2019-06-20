package common

import (
	"common"

	"github.com/astaxie/beego"
)

const (
	DefaultAccessTokenExpiredSecs    int64   = 7200 //default 2hours
	DefaultSignatureExpiredSecs      uint32  = 120  // default login signature expire seconds
	DefaultSignatureTokenExpiredSecs uint32  = 300  // default signature token expired seconds
	DefaultEusdPrecision             int     = 4
	DefaultUsdtPrecision             int     = 2
	DefaultUsdtPriceExpiredSecs      int     = 7200 // default 24 hours
	DefaultUsdtMinPrice              float64 = 4.0  // default usdt cny min price
	DefaultUsdtMaxPrice              float64 = 10.0 // default usdt cny max price
	DefaultBtcPrecision              int     = 2
	DefaultBtcPriceExpiredSecs       int     = 7200     // default 24 hours
	DefaultBtcMinPrice               float64 = 20700.0  // default btc cny min price
	DefaultBtcMaxPrice               float64 = 138000.0 // default btc cny max price
	EusdCommissionWithdrawLowLimit   float64 = 0.1
	DefaultAgentLockTimeout          int64   = 5   //default 5second
	DefaultAgentRetryDelay           int64   = 500 //default 0.5second
	DefaultAgentRetryCount           int     = 100 //default 100 count
)

type OtcSvr struct {
	common.SvrBase

	AccessTokenExpiredSecs int64 //access_token expire seconds

	SignatureExpiredSecs uint32 // login signature expire seconds

	SignatureTokenExpiredSecs uint32 // default signature token expired seconds

	EusdPrecision                                                       int   //eusd precision
	EusdCommissionWithDrawLowInteger, EusdCommissionWithDrawLowDecimals int32 //eusd commission withdraw low limit

	// game
	EusdExchangeRate      float64 // app eusd exchange rate
	EusdLimit             float64 // app eusd limit
	ChannelUrl            string  //game channel url
	PlatId                uint32  //game plat id
	ChannelAppKey         string  //game app key
	GameChannelPercentage float64 // game channel percentage

	// usdt
	UsdtPrecision        int     //usdt Precision
	UsdtPriceExpiredSecs int     //usdt price expired seconds
	UsdtMinPrice         float64 //usdt min price
	UsdtMaxPrice         float64 //usdt max price

	// btc
	BtcPrecision        int     //btc Precision
	BtcPriceExpiredSecs int     //btc price expired seconds
	BtcMinPrice         float64 //btc min price
	BtcMaxPrice         float64 //btc max price

	// agent
	AgentWithdrawFee int    //agent withdraw fee (percent)
	AgentInviteUrl   string //agent invite url

	// oss
	AccessKeyId     string //
	AccessKeySecret string //
	Endpoint        string //
	BucketName      string //

	//agent redis
	LockTimeout int64
	RetryDelay  int64
	RetryCount  int
}

func NewOtcSvr() *OtcSvr {
	return &OtcSvr{
		SvrBase: common.SvrBase{},
	}
}

func (o *OtcSvr) Init() (err error) {

	if err = o.SvrBase.Init(); err != nil {
		return
	}

	{
		if expiredSecs, err := beego.AppConfig.Int64("CaptchaExpiredSecs"); err != nil {
			o.AccessTokenExpiredSecs = DefaultAccessTokenExpiredSecs
		} else {
			o.AccessTokenExpiredSecs = expiredSecs
		}
	}

	{
		if expiredSecs, err := beego.AppConfig.Int64("SignatureExpiredSecs"); err != nil {
			o.SignatureExpiredSecs = DefaultSignatureExpiredSecs
		} else {
			o.SignatureExpiredSecs = uint32(expiredSecs)
		}
	}

	{
		if expiredSecs, err := beego.AppConfig.Int64("SignatureTokenExpiredSecs"); err != nil {
			o.SignatureTokenExpiredSecs = DefaultSignatureTokenExpiredSecs
		} else {
			o.SignatureTokenExpiredSecs = uint32(expiredSecs)
		}
	}

	{
		if precision, err := beego.AppConfig.Int("eusd::Precision"); err != nil {
			o.EusdPrecision = DefaultEusdPrecision
		} else {
			o.EusdPrecision = precision
		}
	}

	{
		if limit, err := beego.AppConfig.Float("eusd::CommissionWithdrawLowLimit"); err != nil {
			o.EusdCommissionWithDrawLowDecimals = 10000000
		} else {
			o.EusdCommissionWithDrawLowInteger, o.EusdCommissionWithDrawLowDecimals = common.DecodeCurrencyFLoat64(limit)
		}
	}

	{
		if precision, err := beego.AppConfig.Int("usdt::PricePrecision"); err != nil {
			o.UsdtPrecision = DefaultUsdtPrecision
		} else {
			o.UsdtPrecision = precision
		}
	}

	{
		if secs, err := beego.AppConfig.Int("usdt::PriceExpiredSecs"); err != nil {
			o.UsdtPriceExpiredSecs = DefaultUsdtPriceExpiredSecs
		} else {
			o.UsdtPriceExpiredSecs = secs
		}
	}

	{
		if price, err := beego.AppConfig.Float("usdt::MinPrice"); err != nil {
			o.UsdtMinPrice = DefaultUsdtMinPrice
		} else {
			o.UsdtMinPrice = price
		}
	}

	{
		if price, err := beego.AppConfig.Float("usdt::MaxPrice"); err != nil {
			o.UsdtMaxPrice = DefaultUsdtMaxPrice
		} else {
			o.UsdtMaxPrice = price
		}
	}

	{
		if precision, err := beego.AppConfig.Int("btc::PricePrecision"); err != nil {
			o.BtcPrecision = DefaultBtcPrecision
		} else {
			o.BtcPrecision = precision
		}
	}

	{
		if secs, err := beego.AppConfig.Int("btc::PriceExpiredSecs"); err != nil {
			o.BtcPriceExpiredSecs = DefaultBtcPriceExpiredSecs
		} else {
			o.BtcPriceExpiredSecs = secs
		}
	}

	{
		if price, err := beego.AppConfig.Float("btc::MinPrice"); err != nil {
			o.BtcMinPrice = DefaultBtcMinPrice
		} else {
			o.BtcMinPrice = price
		}
	}

	{
		if price, err := beego.AppConfig.Float("btc::MaxPrice"); err != nil {
			o.BtcMaxPrice = DefaultBtcMaxPrice
		} else {
			o.BtcMaxPrice = price
		}
	}

	{
		if rate, err := beego.AppConfig.Float("game::EusdExchangeRate"); err != nil {
			o.EusdExchangeRate = 6.0000
		} else {
			o.EusdExchangeRate = rate
		}

		if limit, err := beego.AppConfig.Float("game::EusdLimit"); err != nil {
			o.EusdLimit = 8000.00
		} else {
			o.EusdLimit = limit
		}

		if percentage, err := beego.AppConfig.Float("game::GameChannelPercentage"); err != nil {
			o.GameChannelPercentage = 0.1
		} else {
			o.GameChannelPercentage = percentage
		}
	}

	{
		o.ChannelUrl = beego.AppConfig.String("game::ChannelUrl")
		o.ChannelAppKey = beego.AppConfig.String("game::AppKey")
	}

	{
		if platId, err := beego.AppConfig.Int("game::PlatId"); err != nil {
			o.PlatId = 57 //default test plat id
		} else {
			o.PlatId = uint32(platId)
		}
	}

	{
		o.AccessKeyId = beego.AppConfig.String("oss::accessKeyId")
		o.AccessKeySecret = beego.AppConfig.String("oss::accessKeySecret")
		o.Endpoint = beego.AppConfig.String("oss::endpoint")
		o.BucketName = beego.AppConfig.String("oss::bucketName")
	}
	{
		if lockTimeout, err := beego.AppConfig.Int64("agent::lockTimeout"); err != nil {
			o.LockTimeout = DefaultAgentLockTimeout
		} else {
			o.LockTimeout = lockTimeout
		}
		if retryDelay, err := beego.AppConfig.Int64("agent::retryDelay"); err != nil {
			o.RetryDelay = DefaultAgentRetryDelay
		} else {
			o.RetryDelay = retryDelay
		}
		if retryCount, err := beego.AppConfig.Int("agent::retryCount"); err != nil {
			o.RetryCount = DefaultAgentRetryCount
		} else {
			o.RetryCount = retryCount
		}
	}

	{
		o.AgentInviteUrl = beego.AppConfig.String("agent::InviteUrl")
	}

	return
}

//current server
var Cursvr *OtcSvr = NewOtcSvr()
var ServerOtcState = common.ServerStateRunning
var ServerGameRunning = true
var ServerOtcTradeRunning = true
