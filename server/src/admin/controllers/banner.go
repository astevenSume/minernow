package controllers

import (
	controllers "admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
)

type BannerController struct {
	BaseController
}

//获取广告
func (c *BannerController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadBanner, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadBanner, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	var input string
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := dao.BannerDaoEntity.QueryById(uint32(id))
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadBanner, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadBanner, input, data)
	} else {
		status, err := c.GetInt("status", 0)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadBanner, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		page, err := c.GetInt(KEY_PAGE, 1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadBanner, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadBanner, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		//分页查询
		res := map[string]interface{}{}
		input = fmt.Sprintf("{\"status\":%v,\"page\":%v,\"limit\":%v}", status, page, perPage)
		count, data, err := dao.BannerDaoEntity.QueryByPage(int8(status), page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadBanner, controllers.ERROR_CODE_DB, input)
			return
		}

		meta := dao.PageInfo{
			Limit: perPage,
			Total: int(count),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadBanner, input, res)
	}
}

func (c *BannerController) Create() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddBanner, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Subject string `json:"subject"`
		Image   string `json:"image"`
		Url     string `json:"url"`
		Stime   int64  `json:"stime"`
		Etime   int64  `json:"etime"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddBanner, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	data, err := dao.BannerDaoEntity.Create(req.Subject, req.Image, req.Url, req.Stime, req.Etime)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionAddBanner, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionAddBanner, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionAddBanner, string(c.Ctx.Input.RequestBody), data)
}

//更新
func (c *BannerController) Update() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditBanner, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id      uint32 `json:"id"`
		Subject string `json:"subject"`
		Image   string `json:"image"`
		Url     string `json:"url"`
		Stime   int64  `json:"stime"`
		Etime   int64  `json:"etime"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditBanner, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	data, err := dao.BannerDaoEntity.Update(req.Id, req.Subject, req.Image, req.Url, req.Stime, req.Etime)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionEditBanner, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionEditBanner, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionEditBanner, string(c.Ctx.Input.RequestBody), data)
}

//删除
func (c *BannerController) Del() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelBanner, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, 0)
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionDelBanner, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.BannerDaoEntity.DelById(uint32(id))
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelBanner, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelBanner, input, ack)
}

//上线
func (c *BannerController) Publish() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddBannerPublish, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, 0)
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionAddBannerPublish, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.BannerDaoEntity.UpdateStatus(uint32(id), dao.BannerStatusPublished)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddBannerPublish, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionAddBannerPublish, input, ack)
}

//取消上线
func (c *BannerController) UnPublish() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddBannerUnPublish, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, 0)
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionAddBannerUnPublish, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.BannerDaoEntity.UpdateStatus(uint32(id), dao.BannerStatusPending)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddBannerUnPublish, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionAddBannerUnPublish, input, ack)
}
