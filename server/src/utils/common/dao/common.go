package dao

import (
	"common"
	"github.com/astaxie/beego"
)

func Init(entityInitFunc common.EntityInitFunc) (err error) {
	////init token memory
	//adapter, err = NewCache("memory", `{"interval":20}`)
	//if err != nil {
	//	return
	//}

	//get access_token expired second configuration.
	AccessTokenExpiredSecs, err = beego.AppConfig.Int64("AccessTokenExpiredSecs")
	appName = beego.AppConfig.String("appname")

	//init entity (todo ： maybe this should be moved to business svr project?)
	TokenDaoEntity = NewTokenDao("otc")

	if entityInitFunc != nil {
		err = entityInitFunc()
		if err != nil {
			return
		}
	}

	////load token from db
	//err = TokenDaoEntity.Load()
	//if err != nil {
	//	return
	//}

	return
}
