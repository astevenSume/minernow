package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
	"utils/admin/models"
	common2 "utils/common"
)

type ConfigController struct {
	BaseController
}

// 获取系统配置列表
func (c *ConfigController) List() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadConfig, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadConfig, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	limit, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadConfig, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	action := c.GetString(KEY_KEY_ACTION)
	actionType, ok := dao.ConfigDaoEntity.GetActionType(action)
	if !ok {
		c.ErrorResponseAndLog(OPActionReadConfig, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"page\":%v,\"page\":%v}", page, limit)

	total, config, err := dao.ConfigDaoEntity.QueryPageConfig(actionType, page, limit)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadConfig, controllers.ERROR_CODE_DB, input)
		return
	}

	res := map[string]interface{}{}
	meta := PageInfo{
		Limit: limit,
		Total: int(total),
		Page:  page,
	}
	res["data"] = config
	res["meta"] = meta

	c.SuccessResponseAndLog(OPActionReadConfig, input, res)
}

//easyjson:json
type ConfigCreateReq struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Desc  string `json:"desc"`
}

// 新增系统配置
func (c *ConfigController) Create() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddConfig, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := ConfigCreateReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncInfo("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddConfig, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	action := c.GetString(KEY_KEY_ACTION)
	actionType, ok := dao.ConfigDaoEntity.GetActionType(action)
	if !ok {
		c.ErrorResponseAndLog(OPActionAddConfig, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	res := models.Config{
		Key:    req.Key,
		Action: actionType,
		Desc:   req.Desc,
		Value:  req.Value,
	}

	if err := dao.ConfigDaoEntity.InsertConfig(actionType, req.Key, req.Value, req.Desc); err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionAddConfig, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		c.ErrorResponseAndLog(OPActionAddConfig, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseAndLog(OPActionAddConfig, string(c.Ctx.Input.RequestBody), res)
}

//easyjson:json
type ConfigUpdateReq struct {
	Id     uint32 `orm:"column(id);pk" json:"id,omitempty"`
	Action string `orm:"column(action);size(32)" json:"action,omitempty"`
	Key    string `orm:"column(key);size(256)" json:"key,omitempty"`
	Value  string `orm:"column(value)" json:"value,omitempty"`
	Desc   string `orm:"column(desc);size(256)" json:"desc,omitempty"`
	Ctime  int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

// 修改系统配置
func (c *ConfigController) Update() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditConfig, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := ConfigUpdateReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncInfo("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditConfig, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	action := c.GetString(KEY_KEY_ACTION)
	actionType, ok := dao.ConfigDaoEntity.GetActionType(action)
	if !ok || req.Id == 0 {
		c.ErrorResponseAndLog(OPActionEditConfig, controllers.ERROR_CODE_PARAM_FAILED, string(c.Ctx.Input.RequestBody))
		return
	}

	if err := dao.ConfigDaoEntity.UpdateConfig(req.Id, actionType, req.Key, req.Value, req.Desc); err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionEditConfig, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		c.ErrorResponseAndLog(OPActionEditConfig, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseAndLog(OPActionEditConfig, string(c.Ctx.Input.RequestBody), req)
}

// 删除系统配置
func (c *ConfigController) Delete() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelConfig, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	action := c.GetString(KEY_KEY_ACTION)
	actionType, ok := dao.ConfigDaoEntity.GetActionType(action)
	if !ok {
		c.ErrorResponseAndLog(OPActionDelConfig, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelConfig, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)

	if err := dao.ConfigDaoEntity.DelConfigById(id, actionType); err != nil {
		c.ErrorResponseAndLog(OPActionDelConfig, controllers.ERROR_CODE_DB, input)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionDelConfig, input)
}

// 刷新缓存
func (c *ConfigController) Flush() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditConfigRefresh, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	action := c.GetString(KEY_KEY_ACTION)
	actionType, ok := dao.ConfigDaoEntity.GetActionType(action)
	if !ok {
		c.ErrorResponseAndLog(OPActionEditConfigRefresh, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditConfigRefresh, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)

	if err := common2.AppConfigMgr.FlushActionCache(actionType, id); err != nil {
		c.ErrorResponseAndLog(OPActionEditConfigRefresh, controllers.ERROR_CODE_UPDATE_FAIL, input)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionEditConfigRefresh, input)
}

// 清除缓存
func (c *ConfigController) Clean() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelConfigClean, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	action := c.GetString(KEY_KEY_ACTION)
	actionType, ok := dao.ConfigDaoEntity.GetActionType(action)
	if !ok {
		c.ErrorResponseAndLog(OPActionDelConfig, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelConfigClean, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)

	if err := common2.AppConfigMgr.CleanActionCache(actionType, id); err != nil {
		c.ErrorResponseAndLog(OPActionDelConfigClean, controllers.ERROR_CODE_UPDATE_FAIL, input)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionDelConfigClean, string(c.Ctx.Input.RequestBody))
}

// 获取预警配置列表
func (c *ConfigController) ListWarning() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadConfigWarning, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadConfigWarning, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	limit, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadConfigWarning, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	cType, err := c.GetInt8(KEY_TYPE, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadConfigWarning, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	mobile := c.GetString(KEY_MOBILE)
	input := fmt.Sprintf("{\"type\":%v,\"mobile\":%v,\"page\":%v,\"page\":%v}", cType, mobile, page, limit)

	total, config, err := dao.ConfigWarningDaoEntity.QueryPageConfigWarning(page, limit, cType, mobile)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadConfigWarning, controllers.ERROR_CODE_DB, input)
		return
	}

	res := map[string]interface{}{}
	meta := PageInfo{
		Limit: limit,
		Total: int(total),
		Page:  page,
	}
	res["data"] = config
	res["meta"] = meta
	c.SuccessResponseAndLog(OPActionReadConfigWarning, input, res)
}

//easyjson:json
type ConfigWarningCreateReq struct {
	CType        int8   `json:"type"`
	NationalCode string `json:"national_code"`
	Mobile       string `json:"mobile"`
	SmsType      int8   `json:"sms_type"`
}

// 新增预警配置
func (c *ConfigController) CreateWarning() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddConfigWarning, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := ConfigWarningCreateReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncInfo("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddConfigWarning, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	if req.NationalCode == NatinalChina {
		c.ErrorResponseAndLog(OPActionAddConfigWarning, controllers.ERROR_CODE_SMS_NO_SUPPORT_CHINA_MOBILE, string(c.Ctx.Input.RequestBody))
		return
	}

	if err := dao.ConfigWarningDaoEntity.Create(req.CType, req.SmsType, req.NationalCode, req.Mobile); err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionAddConfigWarning, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		c.ErrorResponseAndLog(OPActionAddConfigWarning, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddConfigWarning, string(c.Ctx.Input.RequestBody))
}

//easyjson:json
type ConfigWarningUpdateReq struct {
	Id           uint32 `orm:"column(id);pk" json:"id,omitempty"`
	CType        int8   `json:"type"`
	NationalCode string `json:"national_code"`
	Mobile       string `json:"mobile"`
	SmsType      int8   `json:"sms_type"`
}

// 修改预警配置
func (c *ConfigController) UpdateWarning() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditConfigWarning, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := ConfigWarningUpdateReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncInfo("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditConfigWarning, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	if err := dao.ConfigWarningDaoEntity.Update(req.CType, req.SmsType, req.Id, req.NationalCode, req.Mobile); err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionEditConfigWarning, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		c.ErrorResponseAndLog(OPActionEditConfigWarning, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionEditConfigWarning, string(c.Ctx.Input.RequestBody))
}

// 删除预警配置
func (c *ConfigController) DelWarning() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelConfigWarning, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelConfigWarning, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)
	if err := dao.ConfigWarningDaoEntity.DelById(id); err != nil {
		c.ErrorResponseAndLog(OPActionDelConfigWarning, controllers.ERROR_CODE_DB, input)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionDelConfigWarning, input)
}
