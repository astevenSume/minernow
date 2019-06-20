package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
	"utils/admin/models"
)

type AppVersionController struct {
	BaseController
}

// 获取app最新版本
func (c *AppVersionController) GetLastAppVersion() {
	c.setOPAction(OPActionReadVersion)
	c.setRequestData(string(c.Ctx.Input.RequestBody))
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadVersion, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	system, _ := c.GetInt8(KEY_SYSTEM_INPUT)
	appVersion, err := dao.AppVersionDaoEntity.QueryLastAppVersion(system)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	c.SuccessResponse(appVersion)
}

// 获取app版本列表 arg: msg.GetAppVersionsReq{}
func (c *AppVersionController) GetAppVersions() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadVersion, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type GetAppVersionsReq struct {
		System  int8 `json:"system"`
		Status  int8 `json:"status"`
		Page    int  `json:"page"`
		PerPage int  `json:"per_page"`
	}
	req := &GetAppVersionsReq{}

	req.Page, _ = c.GetInt(KEY_PAGE)
	req.PerPage, _ = c.GetInt(KEY_LIMIT)
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PerPage == 0 {
		req.PerPage = DEFAULT_PER_PAGE
	}
	req.Status, _ = c.GetInt8("status")
	req.System, _ = c.GetInt8("system")
	input := fmt.Sprintf("{\"status\":%d,\"system\":%d,\"page\":%d,\"limit\":%d}", req.Status, req.System, req.Page, req.PerPage)

	appVersions, meta, err := dao.AppVersionDaoEntity.QueryPageAppVersions(req.System, req.Status, req.Page, req.PerPage)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadVersion, controllers.ERROR_CODE_DB, input)
		return
	}

	res := map[string]interface{}{}
	res["list"] = appVersions
	res["meta"] = meta

	c.SuccessResponseAndLog(OPActionReadVersion, input, res)
}

// 新增app版本信息
func (c *AppVersionController) CreateAppVersion() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddVersion, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type Msg struct {
		Version    string `json:"version"`
		ChangeLog  string `json:"changelog"`
		Download   string `json:"download"`
		System     int8   `json:"system"`
		VersionNum int32  `json:"version_num"`
	}
	req := &Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddVersion, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	appVersion := models.AppVersion{
		Version:    req.Version,
		System:     req.System,
		Download:   req.Download,
		ChangeLog:  req.ChangeLog,
		VersionNum: req.VersionNum,
	}

	if err := dao.AppVersionDaoEntity.CreateAppVersion(&appVersion); err != nil {
		c.ErrorResponseAndLog(OPActionAddVersion, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseAndLog(OPActionAddVersion, string(c.Ctx.Input.RequestBody), appVersion)
}

// 修改app版本信息
func (c *AppVersionController) UpdateAppVersion() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditVersion, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type Msg struct {
		Version    string `json:"version"`
		ChangeLog  string `json:"changelog"`
		Download   string `json:"download"`
		System     int8   `json:"system"`
		VersionNum int32  `json:"version_num"`
	}
	req := &Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditVersion, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	id, _ := c.GetInt64(KEY_ID_INPUT)

	appVersion := models.AppVersion{
		Id:         id,
		Version:    req.Version,
		System:     req.System,
		Download:   req.Download,
		ChangeLog:  req.ChangeLog,
		VersionNum: req.VersionNum,
	}

	if err := dao.AppVersionDaoEntity.UpdateAppVersion(&appVersion); err != nil {
		c.ErrorResponseAndLog(OPActionEditVersion, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseAndLog(OPActionAddVersion, string(c.Ctx.Input.RequestBody), appVersion)
}

// 发布app版本信息
func (c *AppVersionController) PublishAppVersion() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddVersionPublish, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, _ := c.GetInt64(KEY_ID_INPUT)

	if err := dao.AppVersionDaoEntity.PublishAppVersion(id); err != nil {
		c.ErrorResponseAndLog(OPActionAddVersionPublish, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	type ResMsg struct {
		Id int64 `json:"id"`
	}
	ack := ResMsg{
		Id: id,
	}
	c.SuccessResponseAndLog(OPActionAddVersionPublish, string(c.Ctx.Input.RequestBody), ack)
}

// 下架app版本信息
func (c *AppVersionController) UnPublishAppVersion() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddVersionUnPublish, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, _ := c.GetInt64(KEY_ID_INPUT)

	if err := dao.AppVersionDaoEntity.PendAppVersion(id); err != nil {
		c.ErrorResponseAndLog(OPActionAddVersionUnPublish, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	type ResMsg struct {
		Id int64 `json:"id"`
	}
	ack := ResMsg{
		Id: id,
	}
	c.SuccessResponseAndLog(OPActionAddVersionUnPublish, string(c.Ctx.Input.RequestBody), ack)
}

// 删除app版本信息
func (c *AppVersionController) DelAppVersion() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelVersion, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, _ := c.GetInt64(KEY_ID_INPUT)

	if err := dao.AppVersionDaoEntity.DelAppVersion(id); err != nil {
		c.ErrorResponseAndLog(OPActionDelVersion, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	type ResMsg struct {
		Id int64 `json:"id"`
	}
	ack := ResMsg{
		Id: id,
	}
	c.SuccessResponseAndLog(OPActionDelVersion, string(c.Ctx.Input.RequestBody), ack)
}
