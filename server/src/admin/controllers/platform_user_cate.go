package controllers

import (
	"admin/controllers/errcode"
	"fmt"
	"utils/eusd/dao"
)

type PlatformUserCateController struct {
	BaseController
}

func (c *PlatformUserCateController) List() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadPlatformUserCate, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	list, err := dao.PlatformUserCateDaoEntity.All()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadPlatformUserCate, controllers.ERROR_CODE_DB, "")
	}

	c.SuccessResponseAndLog(OPActionReadPlatformUserCate, "", map[string]interface{}{"list": list})
}

// 新增
func (c *PlatformUserCateController) Add() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddPlatformUserCateAdd, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	type post struct {
		Name string `json:"name"`
	}

	msg := post{}

	err := c.GetPost(&msg)
	if err != nil {
		return
	}

	data, err := dao.PlatformUserCateDaoEntity.Add(msg.Name, 0)

	c.SuccessResponseWithoutDataAndLog(OPActionAddPlatformUserCateAdd, fmt.Sprintf("{\"id\":%v}", data.Id))
}
