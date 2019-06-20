package controllers

import (
	controllers "admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
)

type AppController struct {
	BaseController
}

//获取应用
func (c *AppController) GetApp() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadApp, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	id, err := c.GetInt(KEY_ID, -1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadApp, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := dao.AppDaoEntity.QueryAppById(uint32(id))
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadApp, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadApp, input, *data)
	} else {
		appType, err := c.GetInt8("type", -1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadApp, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		featured, err := c.GetInt8("featured", -1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadApp, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		status, err := c.GetInt8(KEY_STATUS, -1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadApp, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		channelId, err := c.GetInt8("channel_id", -1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadApp, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		page, err := c.GetInt(KEY_PAGE, 1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadApp, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadApp, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		//分页查询
		res := map[string]interface{}{}
		input = fmt.Sprintf("{\"channelId\":%v, \"status\":%v,\"type\":%v,\"featured\":%v,\"page\":%v,\"limit\":%v}",
			channelId, status, appType, featured, page, perPage)
		count, data, err := dao.AppDaoEntity.QueryPageCondition(channelId, appType, featured, status, page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadApp, controllers.ERROR_CODE_DB, input)
			return
		}
		meta := dao.PageInfo{
			Limit: perPage,
			Total: int(count),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadApp, input, res)
	}
}

//创建
func (c *AppController) CreateApp() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddApp, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Name        string `json:"name"`
		Desc        string `json:"desc"`
		Url         string `json:"url"`
		IconUrl     string `json:"icon_url"`
		TypeId      int8   `json:"type_id"`
		Featured    int8   `json:"Featured"`
		ChannelId   uint32 `json:"channel_id"`
		AppId       string `json:"app_id"`
		Orientation int8   `json:"orientation"`
		Position    uint32 `json:"position"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddApp, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	app, err := dao.AppDaoEntity.CreateApp(req.TypeId, req.Orientation, req.Featured, req.ChannelId, req.AppId, req.Position,
		req.Name, req.Desc, req.Url, req.IconUrl)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionAddApp, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionAddApp, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionAddApp, string(c.Ctx.Input.RequestBody), app)
}

//更新
func (c *AppController) UpdateApp() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditApp, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id          uint32 `json:"id"`
		Name        string `json:"name"`
		Desc        string `json:"desc"`
		Url         string `json:"url"`
		IconUrl     string `json:"icon_url"`
		TypeId      int8   `json:"type_id"`
		Featured    int8   `json:"featured"`
		Status      int8   `json:"status"`
		ChannelId   uint32 `json:"channel_id"`
		AppId       string `json:"app_id"`
		Orientation int8   `json:"orientation"`
		Position    uint32 `json:"position"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditApp, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	app, err := dao.AppDaoEntity.UpdateApp(req.Id, req.Orientation, req.Status, req.TypeId, req.Featured, req.ChannelId,
		req.AppId, req.Position, req.Name, req.Desc, req.Url, req.IconUrl)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditApp, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionEditApp, string(c.Ctx.Input.RequestBody), app)
}

//删除
func (c *AppController) DelApp() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelApp, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, -1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelApp, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.AppDaoEntity.DelAppById(uint32(id))
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelApp, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelApp, input, ack)
}

//推荐
func (c *AppController) FeatureApp() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAppFeature, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, -1)
	if err != nil || id < 0 {
		c.ErrorResponseAndLog(OPActionAddAppFeature, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.AppDaoEntity.UpdateFeatured(uint32(id), dao.AppFeaturedYes)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddAppFeature, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionAddAppFeature, input, ack)
}

//不推荐
func (c *AppController) UnFeatureApp() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAppUnFeature, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, -1)
	if err != nil || id < 0 {
		c.ErrorResponseAndLog(OPActionAddAppUnFeature, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.AppDaoEntity.UpdateFeatured(uint32(id), dao.AppFeaturedNo)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddAppUnFeature, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionAddAppUnFeature, input, ack)
}

//上线
func (c *AppController) PublishApp() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAppPublish, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, -1)
	if err != nil || id < 0 {
		c.ErrorResponseAndLog(OPActionAddAppPublish, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.AppDaoEntity.UpdateStatus(uint32(id), dao.AppStatusPublish)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddAppPublish, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionAddAppPublish, input, ack)
}

//取消上线
func (c *AppController) UnPublishApp() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAppUnPublish, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, -1)
	if err != nil || id < 0 {
		c.ErrorResponseAndLog(OPActionAddAppUnPublish, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.AppDaoEntity.UpdateStatus(uint32(id), dao.AppStatusUnPublish)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddAppUnPublish, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionAddAppUnPublish, input, ack)
}
