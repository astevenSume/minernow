package main

import (
	"common"
	"common/systoolbox"
	"eusd/eosplus"
	"github.com/astaxie/beego"
	"math/rand"
	"os"
	common2 "otc/common"
	"otc/controllers"
	"otc/routers"
	"time"
	"usdt"
	admindao "utils/admin/dao"
	adminmodels "utils/admin/models"
	agentdao "utils/agent/dao"
	agentmodels "utils/agent/models"
	utils "utils/common"
	utilsmodels "utils/common/models"
	eusddao "utils/eusd/dao"
	eusdmodels "utils/eusd/models"
	gamedao "utils/game/dao"
	"utils/game/dao/gameapi"
	gamemodels "utils/game/models"
	otcmodels "utils/otc/models"
	reportdao "utils/report/dao"
	reportmodels "utils/report/models"
	usdtmodels "utils/usdt/models"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	//init logger
	err := common.LogInit()
	if err != nil {
		panic(err)
	}

	for _, f := range []func() error{
		common2.Cursvr.Init,
		routers.InitFilter,
		modelsInit,
		common.DbInit,
		common.RedisInit,
		common.CaptchaInit,
		common2.MemoryCacheInit,
		usdt.Init,
		usdt.InitSyncRechargeTransaction,
		daoInit,
		controllers.Init,
		eosplus.EosApiInit,
		gameapi.Init,
		sysInit,
	} {
		if err := f(); err != nil {
			common.LogFuncError("init failed : %v", err)
			os.Exit(1)
		}
	}

	//toolbox.StartTask()
	//defer toolbox.StopTask()

	//new task
	systoolbox.TaskMgr.Start(common2.Cursvr.ServerName, common2.Cursvr.RegionId, common2.Cursvr.ServerId)
	defer systoolbox.TaskMgr.Stop()

	beego.Run()
}

// models init
func modelsInit() (err error) {
	if err = otcmodels.ModelsInit(); err != nil {
		return
	}

	if err = usdtmodels.ModelsInit(); err != nil {
		return
	}

	if err = utilsmodels.ModelsInit(); err != nil {
		return
	}

	if err = adminmodels.ModelsInit(); err != nil {
		return
	}

	if err = gamemodels.ModelsInit(); err != nil {
		return
	}

	if err = eusdmodels.ModelsInit(); err != nil {
		return
	}

	if err = agentmodels.ModelsInit(); err != nil {
		return
	}

	if err = reportmodels.ModelsInit(); err != nil {
		return
	}

	return
}

// dao init
func daoInit() (err error) {
	if err = utils.Init(nil); err != nil {
		return
	}

	if err = admindao.Init(nil); err != nil {
		return
	}

	if err = gamedao.Init(nil); err != nil {
		return
	}

	if err = eusddao.Init(nil); err != nil {
		return
	}

	if err = agentdao.Init(nil); err != nil {
		return
	}

	if err = reportdao.Init(nil); err != nil {
		return
	}

	return
}

// system init
func sysInit() (err error) {
	err = common.SysInit(common2.Cursvr.RegionId, common2.Cursvr.ServerId, common2.Cursvr.ServerName)
	if err != nil {
		return
	}

	return
}
