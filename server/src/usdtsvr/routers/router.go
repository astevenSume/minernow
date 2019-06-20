package routers

import (
	"github.com/astaxie/beego/context"
)

func init() {
	return
}

func before(ctx *context.Context) {
	//set output Content-Type to be json
	ctx.Output.Header("Content-Type", "application/json;charset=utf-8")
}
