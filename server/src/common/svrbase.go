package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// @Description server base info.
type SvrBase struct {
	RegionId   int64  //region id
	ServerId   int64  //server id
	ServerName string //server name
	ServerUrl  string //server url

	LogPath     string
	ProfilePath string
}

func (this *SvrBase) Init() (err error) {

	{
		if regionId, err := beego.AppConfig.Int64("RegionId"); err != nil {
			panic("no specific RegionId")
		} else {
			this.RegionId = regionId
		}
	}

	{
		if serverId, err := beego.AppConfig.Int64("ServerId"); err != nil {
			panic("no specific ServerId")
		} else {
			this.ServerId = serverId
		}
	}

	{
		this.ServerName = beego.AppConfig.String("appname")
		this.ServerUrl = beego.AppConfig.String("appurl")
	}

	{
		this.LogPath = beego.AppConfig.String("log::path")
		this.ProfilePath = this.LogPath + "/profile"
	}

	if debug, err := beego.AppConfig.Bool("OrmDebug"); err != nil {
		orm.Debug = false
	} else {
		orm.Debug = debug
	}

	return
}
