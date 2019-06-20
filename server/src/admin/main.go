package main

import (
	common2 "admin/common"
	"admin/controllers"
	"admin/controllers/cron"
	_ "admin/routers"
	"common"
	"common/systoolbox"
	"eusd/eosplus"
	"github.com/astaxie/beego"
	"os"
	"usdt"
	utilsadmindao "utils/admin/dao"
	utilsadminmodels "utils/admin/models"
	agentdao "utils/agent/dao"
	agentmodels "utils/agent/models"
	utils "utils/common"
	utilsdao "utils/common/dao"
	utilsmodels "utils/common/models"
	utilseusddao "utils/eusd/dao"
	utilseusdmodels "utils/eusd/models"
	utilsGameDao "utils/game/dao"
	utilsGameModels "utils/game/models"
	utilsotcdao "utils/otc/dao"
	utilsotcmodels "utils/otc/models"
	utilsusdtdao "utils/usdt/dao"
	utilsusdtmodels "utils/usdt/models"

	utilsreportdao "utils/report/dao"
	utilsreportmodels "utils/report/models"
	utilsTmpdao "utils/tmp/dao"
	utilstmpmodels "utils/tmp/models"
)

func main() {
	//init logger
	err := common.LogInit()
	if err != nil {
		panic(err)
	}

	for _, f := range []func() error{
		common2.Cursvr.Init,
		modelsInit,
		common.DbInit,
		common.RedisInit,
		daoInit,
		cronInit,
		usdt.Init,
		eosplus.DaoInit,
		eosplus.EosApiInit,
		initRabbitMQ,
		controllers.Init,
	} {
		if err := f(); err != nil {
			common.LogFuncError("init failed : %v", err)
			os.Exit(1)
		}
	}

	//start task
	systoolbox.TaskMgr.Start(common2.Cursvr.ServerName, common2.Cursvr.RegionId, common2.Cursvr.ServerId)
	defer systoolbox.TaskMgr.Stop()

	beego.Run()
}

// cron init
func cronInit() (err error) {
	controllers.FuncContainer = &cron.FunctionContainer{}
	//// init cron function
	//err = common.CronInit(&cron.FunctionContainer{})
	//if err != nil {
	//	return
	//}

	return
}

// models init
func modelsInit() (err error) {
	if err = utilsotcmodels.ModelsInit(); err != nil {
		return
	}

	if err = utilsadminmodels.ModelsInit(); err != nil {
		return
	}

	if err = utilsmodels.ModelsInit(); err != nil {
		return
	}

	if err = utilsusdtmodels.ModelsInit(); err != nil {
		return
	}

	if err = utilseusdmodels.ModelsInit(); err != nil {
		return
	}

	if err = agentmodels.ModelsInit(); err != nil {
		return
	}

	if err = utilsGameModels.ModelsInit(); err != nil {
		return
	}
	if err = utilsreportmodels.ModelsInit(); err != nil {
		return
	}
	if err = utilstmpmodels.ModelsInit(); err != nil {
		return
	}
	return
}

// dao init
func daoInit() (err error) {
	// utils common dao init
	err = utils.Init(func() error {
		utilsdao.TokenDaoEntity = utilsdao.NewTokenDao("otc_admin")
		return nil
	})
	if err != nil {
		return
	}

	// utils admin dao init
	err = utilsadmindao.Init(nil)
	if err != nil {
		return
	}

	// utils otc dao init
	err = utilsotcdao.Init(nil)
	if err != nil {
		return
	}

	// utils eusd dao init
	err = utilseusddao.Init(nil)
	if err != nil {
		return
	}

	// utils usdt dao init
	err = utilsusdtdao.Init(nil)
	if err != nil {
		return
	}

	if err = agentdao.Init(nil); err != nil {
		return
	}
	if err = utilsGameDao.Init(nil); err != nil {
		return
	}
	if err = utilsreportdao.Init(nil); err != nil {
		return
	}
	if err = utilsTmpdao.Init(nil); err != nil {
		return
	}
	return
}

// init rabbit mq
func initRabbitMQ() (err error) {
	err = common.RabbitMQInit(&controllers.AmqpFuncContainer{})
	if err != nil {
		return
	}

	return
}
