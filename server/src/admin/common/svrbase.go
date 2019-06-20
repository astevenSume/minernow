package common

import (
	"common"
	"github.com/astaxie/beego"
)

const (
	DefaultAccessTokenExpiredSecs int64  = 7200 //default 2hours
	DefaultOssExpiredSecs         int64  = 30   //default 30s
	DefaultServerDeadSecs         uint32 = 30   //default 30s
	DefaultAgentLockTimeout       int64  = 5    //default 5second
	DefaultAgentRetryDelay        int64  = 500  //default 0.5second
	DefaultAgentRetryCount        int    = 100  //default 100 count
)

type AdminSvr struct {
	common.SvrBase

	AccessTokenExpiredSecs int64 //access_token expire seconds

	//OSS
	OssAccessKeyId     string
	OssAccessKeySecret string
	OssHost            string
	OssCallbackUrl     string
	OssExpireTime      int64
}

func NewAdminSvr() *AdminSvr {
	return &AdminSvr{
		SvrBase: common.SvrBase{},
	}
}

func (o *AdminSvr) Init() (err error) {
	if err = o.SvrBase.Init(); err != nil {
		return
	}

	{
		if expiredSecs, err := beego.AppConfig.Int64("AccessTokenExpiredSecs"); err != nil {
			o.AccessTokenExpiredSecs = DefaultAccessTokenExpiredSecs
		} else {
			o.AccessTokenExpiredSecs = expiredSecs
		}
	}

	{
		o.OssAccessKeyId = beego.AppConfig.String("oss::accessKeyId")
		o.OssAccessKeySecret = beego.AppConfig.String("oss::accessKeySecret")
		o.OssHost = beego.AppConfig.String("oss::host")
		o.OssCallbackUrl = beego.AppConfig.String("oss::callbackUrl")
		if expiredSecs, err := beego.AppConfig.Int64("expireTime"); err != nil {
			o.OssExpireTime = DefaultOssExpiredSecs
		} else {
			o.OssExpireTime = expiredSecs
		}
	}

	return
}

//current server
var Cursvr *AdminSvr = NewAdminSvr()
