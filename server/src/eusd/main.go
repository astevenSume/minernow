package main

import (
	"common"
	"common/systoolbox"
	"eusd/cron"
	"eusd/eosplus"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
	utilsadmindao "utils/admin/dao"
	utilsadminmodels "utils/admin/models"
	utilsmodels "utils/common/models"
	utilseosdao "utils/eos/dao"
	utilseosmodels "utils/eos/models"
	utilseusddao "utils/eusd/dao"
	utilseusdmodels "utils/eusd/models"
	utilsotcdao "utils/otc/dao"
	utilsotcmodels "utils/otc/models"
	utilsusdtmodels "utils/usdt/models"
)

func main() {
	//init logger
	err := common.LogInit()
	if err != nil {
		panic(err)
	}
	if debug, err := beego.AppConfig.Bool("OrmDebug"); err != nil {
		orm.Debug = false
	} else {
		orm.Debug = debug
	}

	for k, f := range []func() error{
		eosplus.Cursvr.Init,
		common.RedisInit,
		modelsInit,
		common.DbInit,
		cronInit,
		daoInit,
		eosApiInit,
		mqInit,
		sysInit,
	} {
		if err := f(); err != nil {
			common.LogFuncError("init failed : %v; funcKey:%v", err, k)
			os.Exit(1)
		}
	}

	common.LogFuncDebug("EUSD CRON running..")

	//start task
	systoolbox.TaskMgr.Start(eosplus.Cursvr.ServerName, eosplus.Cursvr.RegionId, eosplus.Cursvr.ServerId)
	defer systoolbox.TaskMgr.Stop()

	select {}
}

func cronInit() (err error) {
	err = common.CronInit(&cron.FunctionContainer{})
	return
}

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

	if err = utilseusdmodels.ModelsInit(); err != nil {
		return
	}

	if err = utilsusdtmodels.ModelsInit(); err != nil {
		return
	}

	if err = utilseosmodels.ModelsInit(); err != nil {
		return
	}
	return
}

func daoInit() (err error) {

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

	err = utilseosdao.Init(nil)
	if err != nil {
		return
	}
	return
}

func eosApiInit() (err error) {
	err = eosplus.EosApiInit()
	if err != nil {
		return
	}
	err = eosplus.EosPlusAPI.Rpc.SetTransferAble()
	if r, err := beego.AppConfig.Bool("RpcDebug"); err == nil && r {
		eosplus.EosPlusAPI.Rpc.SetDebug()
	} else {
		common.LogFuncDebug("RpcDebug:%v,%v", r, err)
	}

	eosplus.FuncContainer = &cron.FunctionContainer{}

	return
}

func mqInit() (err error) {
	err = common.RabbitMQInit(&eosplus.AmqpFuncContainer{})
	if err != nil {
		return
	}
	return
}

// system init
func sysInit() (err error) {
	err = common.SysInit(eosplus.Cursvr.RegionId, eosplus.Cursvr.ServerId, eosplus.Cursvr.ServerName)
	if err != nil {
		return
	}

	return
}
