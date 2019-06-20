package controllers

import (
	controllers "admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
)

type AppChannelController struct {
	BaseController
}

//获取app渠道
func (c *AppChannelController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadAppChannel, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAppChannel, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	var input string
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := dao.AppChannelDaoEntity.QueryById(uint32(id))
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppChannel, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadAppChannel, input, data)
	} else {
		res := map[string]interface{}{}
		data, err := dao.AppChannelDaoEntity.QueryAll()
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppChannel, controllers.ERROR_CODE_DB, input)
			return
		}

		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadAppChannel, input, res)
	}
}

func (c *AppChannelController) Create() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAppChannel, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id           uint32 `json:"id"`
		IsThirdHall  int8   `json:"is_third_hall"`
		Name         string `json:"name"`
		Desc         string `json:"desc"`
		ExchangeRate int32  `json:"exchange_rate"`
		Precision    int32  `json:"precision"`
		ProfitRate   int32  `json:"profit_rate"`
		IconUrl      string `json:"icon_url"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddAppChannel, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	isNew, data, err := dao.AppChannelDaoEntity.Create(req.IsThirdHall, req.Id, req.Name, req.Desc, req.IconUrl, req.ExchangeRate, req.Precision, req.ProfitRate)
	if !isNew {
		c.ErrorResponseAndLog(OPActionAddAppChannel, controllers.ERROR_CODE_ID_EXSIT, string(c.Ctx.Input.RequestBody))
		return
	}
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionAddAppChannel, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionAddAppChannel, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionAddAppChannel, string(c.Ctx.Input.RequestBody), data)
}

//更新
func (c *AppChannelController) Update() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditAppChannel, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id           uint32 `json:"id"`
		IsThirdHall  int8   `json:"is_third_hall"`
		Name         string `json:"name"`
		Desc         string `json:"desc"`
		ExchangeRate int32  `json:"exchange_rate"`
		Precision    int32  `json:"precision"`
		ProfitRate   int32  `json:"profit_rate"`
		IconUrl      string `json:"icon_url"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditAppChannel, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	data, err := dao.AppChannelDaoEntity.Update(req.IsThirdHall, req.Id, req.Name, req.Desc, req.IconUrl, req.ExchangeRate, req.Precision, req.ProfitRate)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionEditAppChannel, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionEditAppChannel, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionEditAppChannel, string(c.Ctx.Input.RequestBody), data)
}

//删除
func (c *AppChannelController) Del() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelAppChannel, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, 0)
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionDelAppChannel, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.AppChannelDaoEntity.DelById(uint32(id))
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAppChannel, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	ack := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelAppChannel, input, ack)
}
