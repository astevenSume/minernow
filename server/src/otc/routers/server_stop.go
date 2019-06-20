package routers

import (
	common2 "common"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"otc/common"
	"otc/controllers"
	otcerror "otc_error"
	"strings"
)

var filterUrl = []string{
	"/v1/buy",
	"/v1/sell",
	"/v1/eos/transfer",
	"/v1/agent/withdraw",
}

var filterUrlStopTrade = []string{
	"/v1/buy",
	"/v1/sell",
}

var filterUrlGame = []string{
	"/v1/game/login",
}

func ServerStop(ctx *context.Context) {
	ctl := controllers.BaseController{Controller: beego.Controller{Ctx: ctx}}

	path := strings.ToLower(ctx.Request.URL.Path)
	// otc交易服务关闭
	if !common.ServerOtcTradeRunning {
		for _, v := range filterUrlStopTrade {
			if strings.LastIndex(path, v) > -1 {
				ctl.ErrorResponse(otcerror.ERROR_CODE_SERVER_STOPPED)
				return
			}
		}
	}

	//游戏服务关闭
	if !common.ServerGameRunning {
		for _, v := range filterUrlGame {
			if v == path {
				ctl.ErrorResponse(otcerror.ERROR_CODE_SERVER_STOPPING)
				return
			}
		}
	}

	// otc服务状态判断
	if common.ServerOtcState == common2.ServerStateRunning {
		//正常运行
		return
	}

	if common.ServerOtcState == common2.ServerStateStopSecond {
		//停止阶段2
		ctl.ErrorResponse(otcerror.ERROR_CODE_SERVER_STOPPED)
		return
	}
	//停止阶段1 - 禁止创建订单
	for _, v := range filterUrl {
		if v == path {
			ctl.ErrorResponse(otcerror.ERROR_CODE_SERVER_STOPPING)
			return
		}
	}

	return
}
