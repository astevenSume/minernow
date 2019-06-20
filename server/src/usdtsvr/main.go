package main

import (
	"common"
	"common/systoolbox"
	"os"
	"usdt"
	usdtsvrcommon "usdtsvr/common"
	"usdtsvr/controllers"
	admindao "utils/admin/dao"
	adminmodels "utils/admin/models"
	usdtdao "utils/usdt/dao"
	usdtmodels "utils/usdt/models"
)

func main() {
	//init logger
	err := common.LogInit()
	if err != nil {
		panic(err)
	}

	for _, f := range []func() error{
		usdtsvrcommon.Cursvr.Init,
		modelsInit,
		common.DbInit,
		common.RedisInit,
		usdt.Init,
		daoInit,
		controllers.Init,
		sysInit,
	} {
		if err := f(); err != nil {
			common.LogFuncError("init failed : %v", err)
			os.Exit(1)
		}
	}

	//start task
	systoolbox.TaskMgr.Start(usdtsvrcommon.Cursvr.ServerName, usdtsvrcommon.Cursvr.RegionId, usdtsvrcommon.Cursvr.ServerId)
	defer systoolbox.TaskMgr.Stop()

	//beego.Run()

	// block the main runtine
	forever := make(chan bool)
	<-forever

}

// models init
func modelsInit() (err error) {
	if err = usdtmodels.ModelsInit(); err != nil {
		return
	}

	if err = adminmodels.ModelsInit(); err != nil {
		return
	}
	return
}

// dao init
func daoInit() (err error) {
	if err = usdtdao.Init2(nil); err != nil {
		return
	}

	if err = admindao.Init(nil); err != nil {
		return
	}
	return
}

// system init
func sysInit() (err error) {
	err = common.SysInit(usdtsvrcommon.Cursvr.RegionId, usdtsvrcommon.Cursvr.ServerId, usdtsvrcommon.Cursvr.ServerName)
	if err != nil {
		return
	}

	return
}
