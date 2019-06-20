package common

import (
	//"os"
	"testing"
	//"common"
	//utilsadmindao "utils/admin/dao"
	//utilsadminmodels "utils/admin/models"
	//utilsdao "utils/common/dao"
	//utilsmodels "utils/common/models"
	//utilseusddao "utils/eusd/dao"
	//utilseusdmodels "utils/eusd/models"
	//utilsotcdao "utils/otc/dao"
	//utilsotcmodels "utils/otc/models"
)

func TestConfigManager_GetBool(t *testing.T) {

	//for _, f := range []func() error{
	//	modelsInit,
	//	common.DbInit,
	//	common.RedisInit,
	//	daoInit,
	//} {
	//	if err := f(); err != nil {
	//		common.LogFuncError("init failed : %v", err)
	//		os.Exit(1)
	//	}
	//}

	//var AppConfigMgr ConfigManagerInterface = NewConfigManager()
	//if v, e := AppConfigMgr.Float("buy_fee_rate"); e !=nil {
	//	t.Fatalf(fmt.Sprintf("%f.2,%v", v, e))
	//}
	//if v, e := AppConfigMgr.String("aws_access_key_id"); e !=nil {
	//	t.Fatalf(fmt.Sprintf("%s,%v", v, e))
	//}
	//if v, e := AppConfigMgr.Float("xxx"); e !=nil {
	//	t.Fatalf(fmt.Sprintf("%f.2,%v", v, e))
	//}
	//var v4 interface{}
	//if e := AppConfigMgr.Json("xxx", &v4); e !=nil {
	//	t.Fatalf(fmt.Sprintf("%f.2,%v", v4, e))
	//}

}

//// models init
//func modelsInit() (err error) {
//	if err = utilsotcmodels.ModelsInit(); err != nil {
//		return
//	}
//
//	if err = utilsadminmodels.ModelsInit(); err != nil {
//		return
//	}
//
//	if err = utilsmodels.ModelsInit(); err != nil {
//		return
//	}
//
//	if err = utilseusdmodels.ModelsInit(); err != nil {
//		return
//	}
//	return
//}
//
//// dao init
//func daoInit() (err error) {
//	// utils common dao init
//	err = Init(func() error {
//		utilsdao.TokenDaoEntity = utilsdao.NewTokenDao("otc_admin")
//		return nil
//	})
//	if err != nil {
//		return
//	}
//
//	// utils admin dao init
//	err = utilsadmindao.Init(nil)
//	if err != nil {
//		return
//	}
//
//	// utils otc dao init
//	err = utilsotcdao.Init(nil)
//	if err != nil {
//		return
//	}
//
//	// utils eusd dao init
//	err = utilseusddao.Init(nil)
//	if err != nil {
//		return
//	}
//	return
//}
//
