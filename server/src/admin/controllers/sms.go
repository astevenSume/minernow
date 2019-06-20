package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"strconv"
	"utils/admin/dao"
	"utils/admin/models"
)

type SmsController struct {
	BaseController
}

func (c *SmsController) GetPageInfo() (uint64, int, int, error) {
	//请求参数
	var id uint64
	strId := c.GetString(KEY_ID)
	if len(strId) > 0 {
		var err error
		id, err = strconv.ParseUint(strId, 10, 64)
		if err != nil {
			return 0, 0, 0, err
		}
	}
	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		common.LogFuncError("err:%v", err)
		return 0, 0, 0, err
	}
	perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		common.LogFuncError("err:%v", err)
		return 0, 0, 0, err
	}

	return id, page, perPage, nil
}

//获取短信模板配置
func (c *SmsController) GetSmsTemplates() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadSmsTemplate, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, page, perPage, err := c.GetPageInfo()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadSmsTemplate, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	//数据获取
	if id == 0 {
		res := map[string]interface{}{}
		input = fmt.Sprintf("{\"id\":%v,\"page\":%v,\"limit\":%v}", id, page, perPage)
		//分页查询
		data, pageInf, err := dao.SmsTemplateDaoEntity.QueryPageSmsTemplates(page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadSmsTemplate, controllers.ERROR_CODE_DB, input)
			return
		}
		res["meta"] = pageInf
		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadSmsTemplate, input, res)
	} else {
		input = fmt.Sprintf("{\"id\":%v}", id)
		smstemplates := &models.Smstemplates{Id: int64(id)}
		err := dao.SmsTemplateDaoEntity.QuerySmsTemplates(smstemplates, models.COLUMN_Smstemplates_Id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadSmsTemplate, errCode, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadSmsTemplate, input, *smstemplates)
	}
}

//新增短信模板配置
func (c *SmsController) AddSmsTemplates() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddSmsTemplate, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Name     string `json:"name"`
		Type     int8   `json:"type"`
		Template string `json:"template"`
	}
	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddSmsTemplate, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//插入数据
	data := new(models.Smstemplates)
	data.Type = req.Type
	data.Name = req.Name
	data.Template = req.Template
	data.Ctime = common.NowInt64MS()
	data.Utime = common.NowInt64MS()
	err = dao.SmsTemplateDaoEntity.AddSmsTemplates(data)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddSmsTemplate, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回c查询结果
	c.SuccessResponseAndLog(OPActionAddSmsTemplate, string(c.Ctx.Input.RequestBody), *data)
}

//更新短信模板配置
func (c *SmsController) UpdateSmsTemplates() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditSmsTemplate, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Type     int8   `json:"type"`
		Template string `json:"template"`
	}
	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditSmsTemplate, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	data := models.Smstemplates{
		Id:       req.Id,
		Name:     req.Name,
		Type:     req.Type,
		Template: req.Template,
		Utime:    common.NowInt64MS(),
	}
	err = dao.SmsTemplateDaoEntity.UpdateSmsTemplates(&data, models.COLUMN_Smstemplates_Type, models.COLUMN_Smstemplates_Name, models.COLUMN_Smstemplates_Template, models.COLUMN_Smstemplates_Utime)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditSmsTemplate, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回c查询结果
	c.SuccessResponseAndLog(OPActionEditSmsTemplate, string(c.Ctx.Input.RequestBody), data)
}

//删除短信模板配置
func (c *SmsController) DelSmsTemplates() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelSmsTemplate, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt("id", 0)
	if err != nil || id == 0 {
		c.ErrorResponseAndLog(OPActionDelSmsTemplate, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)

	err = dao.SmsTemplateDaoEntity.DelSmsTemplates(int64(id))
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelSmsTemplate, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回c查询结果
	data := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelSmsTemplate, input, data)
}

//获取短信验证码
func (c *SmsController) GetSmsCodes() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadSmsCode, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, page, perPage, err := c.GetPageInfo()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadSmsCode, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	if id == 0 {
		status, err := c.GetInt8(KEY_STATUS, 0)
		if err != nil {
			common.LogFuncError("error:%v", err)
			return
		}
		mobile := c.GetString(KEY_MOBILE)
		action := c.GetString(KEY_ACTION)
		input = fmt.Sprintf("{\"status\":%v,\"mobile\":%v,\"action\":%v,\"page\":%v,\"per_page\":%v}",
			status, mobile, action, page, perPage)
		//分页查询
		total, data, err := dao.SmsCodeDaoEntity.QueryPageSmsCodes(page, perPage, status, mobile, action)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadSmsCode, controllers.ERROR_CODE_DB, input)
			return
		}
		res := map[string]interface{}{}
		meta := PageInfo{
			Limit: perPage,
			Total: int(total),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadSmsCode, input, res)
	} else {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data := &models.Smscodes{Id: int64(id)}
		err = dao.SmsCodeDaoEntity.QuerySmsCode(data, models.COLUMN_Smscodes_Id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadSmsCode, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadSmsCode, input, *data)
	}
}
