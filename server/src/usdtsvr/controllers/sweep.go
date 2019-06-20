package controllers

import (
	"common"
	"usdt"
)

// ProcessSweep sweep usdt account
func ProcessSweep() {
	err := usdt.Sweep()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
}
