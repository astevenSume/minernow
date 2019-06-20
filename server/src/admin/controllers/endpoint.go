package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
)

type EndPointController struct {
	BaseController
}

//获取域名
func (c *EndPointController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadEndPoint, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	data, err := dao.EndPointDaoEntity.GetAll()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadEndPoint, controllers.ERROR_CODE_DB, "")
		return
	}

	res := map[string]interface{}{}
	res["list"] = data
	c.SuccessResponseAndLog(OPActionReadEndPoint, "", res)

}

//easyjson:json
type EndPointCreateReq struct {
	Endpoint string `json:"endpoint"`
}

func (c *EndPointController) Create() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddEndPoint, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &EndPointCreateReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddEndPoint, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	err = dao.EndPointDaoEntity.Create(req.Endpoint)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddEndPoint, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionAddEndPoint, string(c.Ctx.Input.RequestBody))
}

//easyjson:json
type EndPointUpdateReq struct {
	Id       uint32 `json:"id"`
	Endpoint string `json:"endpoint"`
}

//更新
func (c *EndPointController) Update() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditEndPoint, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &EndPointUpdateReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditEndPoint, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	err = dao.EndPointDaoEntity.Update(req.Id, req.Endpoint)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditEndPoint, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionEditEndPoint, string(c.Ctx.Input.RequestBody))
}

//删除
func (c *EndPointController) Del() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelEndPoint, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionDelEndPoint, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)

	err = dao.EndPointDaoEntity.DelById(id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelEndPoint, controllers.ERROR_CODE_DB, input)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionDelEndPoint, input)
}
