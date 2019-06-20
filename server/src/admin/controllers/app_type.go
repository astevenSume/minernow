package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
)

type AppTypeController struct {
	BaseController
}

//获取app类型
func (c *AppTypeController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadAppType, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAppType, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	var input string
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := dao.AppTypeDaoEntity.QueryById(uint32(id))
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppType, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadAppType, input, data)
	} else {
		res := map[string]interface{}{}
		data, err := dao.AppTypeDaoEntity.QueryAll()
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppType, controllers.ERROR_CODE_DB, input)
			return
		}

		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadAppType, input, res)
	}
}

func (c *AppTypeController) Create() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAppType, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id   uint32 `json:"id"`
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddAppType, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	isNew, data, err := dao.AppTypeDaoEntity.Create(req.Id, req.Name, req.Desc)
	if !isNew {
		c.ErrorResponseAndLog(OPActionAddAppType, controllers.ERROR_CODE_ID_EXSIT, string(c.Ctx.Input.RequestBody))
		return
	}
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionAddAppType, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionAddAppType, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionAddAppType, string(c.Ctx.Input.RequestBody), data)
}

//更新
func (c *AppTypeController) Update() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditAppType, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id   uint32 `json:"id"`
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditAppType, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	data, err := dao.AppTypeDaoEntity.Update(req.Id, req.Name, req.Desc)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionEditAppType, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionEditAppType, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionEditAppType, string(c.Ctx.Input.RequestBody), data)
}

//删除
func (c *AppTypeController) Del() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelAppType, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, 0)
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionDelAppType, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.AppTypeDaoEntity.DelById(uint32(id))
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAppType, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelAppType, input, ack)
}
