package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"strconv"
	"utils/admin/dao"
)

type AnnouncementController struct {
	BaseController
}

//获取公告
func (c *AnnouncementController) Get() {
	c.setOPAction(OPActionReadAnnouncement)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	if id > 0 {
		c.setRequestData(fmt.Sprintf("{\"id\":%v}", id))
		data, err := dao.AnnouncementDaoEntity.QueryById(id)
		if err != nil {
			c.ErrorResponse(controllers.ERROR_CODE_DB)
			return
		}
		c.SuccessResponse(data)
		return
	}
	c.setOPAction(OPActionReadAppeal)

	aType, err := c.GetInt8(KEY_TYPE, 0)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	sTime, err := c.GetInt64("stime", 0)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	eTime, err := c.GetInt64("etime", 0)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}
	limit, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	c.setRequestData(fmt.Sprintf("{\"type\":%v,\"stime\":%v,\"etime\":%v,\"page\":%v,\"limit\":%v}", aType, sTime, eTime, page, limit))
	//分页查询
	res := map[string]interface{}{}
	count, data, err := dao.AnnouncementDaoEntity.QueryByPage(aType, sTime, eTime, page, limit)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	meta := PageInfo{
		Limit: limit,
		Total: int(count),
		Page:  page,
	}
	res["meta"] = meta
	res["list"] = data
	c.SuccessResponse(res)
}

//easyjson:json
type AnnouncementCreateReq struct {
	Type    int8   `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Stime   int64  `json:"stime"`
	Etime   int64  `json:"etime"`
}

func (c *AnnouncementController) Create() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAnnouncement, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	params := &AnnouncementCreateReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, params)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddAnnouncement, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	data, err := dao.AnnouncementDaoEntity.Create(params.Type, params.Title, params.Content, params.Stime, params.Etime)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionAddAnnouncement, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionAddAnnouncement, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionAddAnnouncement, string(c.Ctx.Input.RequestBody), data)
}

//easyjson:json
type AnnouncementUpdateReq struct {
	Id      uint32 `json:"id"`
	Type    int8   `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Stime   int64  `json:"stime"`
	Etime   int64  `json:"etime"`
}

//更新
func (c *AnnouncementController) Update() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditAnnouncement, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	params := &AnnouncementUpdateReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, params)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditAnnouncement, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	data, err := dao.AnnouncementDaoEntity.Update(params.Type, params.Id, params.Title, params.Content, params.Stime, params.Etime)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionEditAnnouncement, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionEditAnnouncement, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionEditAnnouncement, string(c.Ctx.Input.RequestBody), data)
}

//删除
func (c *AnnouncementController) Del() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelAnnouncement, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, 0)
	if _id := c.Ctx.Input.Param(":id"); id == 0 && _id != "" {
		id, _ = strconv.Atoi(_id)
	}
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionDelAnnouncement, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.AnnouncementDaoEntity.DelById(uint32(id))
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAnnouncement, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelAnnouncement, input, ack)
}
