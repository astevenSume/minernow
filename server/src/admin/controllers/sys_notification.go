package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"umeng_push/uemng_plus"
	"utils/admin/dao"
)

type SysNotificationController struct {
	BaseController
}

//获取系统通知
func (c *SysNotificationController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadNotification, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadNotification, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	var input string
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := dao.SysNotificationDaoEntity.QueryById(id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadNotification, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadNotification, input, data)
	} else {
		status, err := c.GetInt8(KEY_STATUS, -1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadNotification, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		page, err := c.GetInt(KEY_PAGE, 1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadNotification, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		limit, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadNotification, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		total, data, err := dao.SysNotificationDaoEntity.QueryByPage(status, page, limit)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadNotification, controllers.ERROR_CODE_DB, input)
			return
		}

		res := map[string]interface{}{}
		meta := dao.PageInfo{
			Limit: limit,
			Total: int(total),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadNotification, input, res)
	}
}

func (c *SysNotificationController) Create() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddNotification, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Content string `json:"content"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddNotification, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	data, err := dao.SysNotificationDaoEntity.Create(uint32(c.adminId), req.Content)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddNotification, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionAddNotification, string(c.Ctx.Input.RequestBody), data)
}

//更新
func (c *SysNotificationController) Update() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditNotification, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id      uint32 `json:"id"`
		Content string `json:"content"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditNotification, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	data, err := dao.SysNotificationDaoEntity.Update(req.Id, uint32(c.adminId), req.Content)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditNotification, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionEditNotification, string(c.Ctx.Input.RequestBody), data)
}

//删除
func (c *SysNotificationController) Del() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelNotification, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionDelNotification, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.SysNotificationDaoEntity.DelById(id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelNotification, controllers.ERROR_CODE_DB, input)
		return
	}
	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelNotification, input, ack)
}

//发布
func (c *SysNotificationController) Publish() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddNotificationPublish, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil || id == 0 {
		c.ErrorResponseAndLog(OPActionAddNotificationPublish, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)

	notify, err := dao.SysNotificationDaoEntity.QueryById(id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadNotification, controllers.ERROR_CODE_DB, input)
		return
	}

	err = dao.SysNotificationDaoEntity.SetStatus(id, dao.SysNotificatioStatusPublish)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddNotificationPublish, controllers.ERROR_CODE_DB, input)
		return
	}

	err = uemng_plus.BroadcastAll(notify.Content, true)
	if err != nil {
		common.LogFuncError("error:%v", err)
	}
	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionAddNotificationPublish, input, ack)
}

//取消发布
func (c *SysNotificationController) UnPublish() {
	/*_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddNotificationUnPublish, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil || id == 0 {
		c.ErrorResponseAndLog(OPActionAddNotificationUnPublish, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.SysNotificationDaoEntity.SetStatus(id, dao.SysNotificatioStatusUnPublish)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddNotificationUnPublish, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionAddNotificationUnPublish, input, ack)*/
}
