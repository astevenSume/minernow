package controllers

import (
	. "admin/controllers/errcode"
	"common"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"strconv"
	"utils/admin/dao"
)

type SystemMessageMethodControllers struct {
	BaseController
}

//get all sysmsg
func (c *SystemMessageMethodControllers) GetAllSystemMessage() {

	_, errCode := c.CheckPermission()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadSystemMessage, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	SysMsgList, err := dao.SystemMessageMethodDaoEntity.QuerySysMsgAll()

	if err != nil {
		c.ErrorResponseAndLog(OPActionReadSystemMessage, ERROR_CODE_DB, "")
		return
	}

	res := map[string]interface{}{}
	res["list"] = SysMsgList
	c.SuccessResponseAndLog(OPActionReadSystemMessage, "", res)
}

//add new sysmsg
func (c *SystemMessageMethodControllers) AddNewSystemMessage() {

	_, errCode := c.CheckPermission()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddSystemMessages, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type Data struct {
		Key    string `json:"key"`
		Buyer  string `json:"buyer"`
		Seller string `json:"seller"`
		Admin  string `json:"admin"`
	}

	data := &Data{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, data)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddSystemMessages, ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	sys_msg, err := dao.SystemMessageMethodDaoEntity.AddSysMsg(data.Key, data.Buyer, data.Seller, data.Admin)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddSystemMessages, ERROR_CODE_DB, "")
		return
	}

	c.SuccessResponseAndLog(OPActionAddSystemMessages, string(c.Ctx.Input.RequestBody), sys_msg)
}

// 通过id获取系统消息
func (c *SystemMessageMethodControllers) GetSystemMessageById() {

	_, errCode := c.CheckPermission()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadSystemMessage, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	var Id uint64
	IdStr := c.Ctx.Input.Param(KEY_ID_INPUT)
	if len(IdStr) > 0 {
		var err error
		Id, err = strconv.ParseUint(IdStr, 10, 64)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadSystemMessage, ERROR_CODE_PARAM_FAILED, IdStr)
			return
		}
	}
	SysMsg, err := dao.SystemMessageMethodDaoEntity.QuerySysMsgById(Id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.ErrorResponseAndLog(OPActionReadSystemMessage, ERROR_CODE_NO_SYSMSG, IdStr)
			return
		}
		c.ErrorResponseAndLog(OPActionReadSystemMessage, ERROR_CODE_DB, IdStr)
		return
	}
	c.SuccessResponseAndLog(OPActionReadSystemMessage, IdStr, SysMsg)
}

// 通过id更新系统消息
func (c *SystemMessageMethodControllers) UpdateSystemMessageById() {

	_, errCode := c.CheckPermission()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditSystemMessages, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	var Id uint64
	IdStr := c.Ctx.Input.Param(KEY_ID_INPUT)
	if len(IdStr) > 0 {
		var err error
		Id, err = strconv.ParseUint(IdStr, 10, 64)
		if err != nil {
			c.ErrorResponseAndLog(OPActionEditSystemMessages, ERROR_CODE_PARAM_FAILED, string(c.Ctx.Input.RequestBody))
			return
		}
	}

	type Data struct {
		Key    string `json:"key"`
		Buyer  string `json:"buyer"`
		Seller string `json:"seller"`
		Admin  string `json:"admin"`
	}
	data := &Data{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, data)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditSystemMessages, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	_, err = dao.SystemMessageMethodDaoEntity.QuerySysMsgById(Id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.ErrorResponseAndLog(OPActionEditSystemMessages, ERROR_CODE_NO_SYSMSG, string(c.Ctx.Input.RequestBody))
			return
		}
		c.ErrorResponseAndLog(OPActionEditSystemMessages, ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	SysMsg, err := dao.SystemMessageMethodDaoEntity.UpdateSysMsg(Id, data.Key, data.Buyer, data.Seller, data.Admin)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditSystemMessages, ERROR_CODE_UPDATE_FAIL, string(c.Ctx.Input.RequestBody))
		return
	}
	c.SuccessResponseAndLog(OPActionEditSystemMessages, string(c.Ctx.Input.RequestBody), SysMsg)
}

// 通过id删除系统消息
func (c *SystemMessageMethodControllers) DelSystemMessageById() {

	_, errCode := c.CheckPermission()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelSystemMessages, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	var Id uint64
	IdStr := c.Ctx.Input.Param(KEY_ID_INPUT)
	if len(IdStr) > 0 {
		var err error
		Id, err = strconv.ParseUint(IdStr, 10, 64)
		if err != nil {
			c.ErrorResponseAndLog(OPActionDelSystemMessages, ERROR_CODE_PARAM_FAILED, IdStr)
			return
		}
	}

	err := dao.SystemMessageMethodDaoEntity.ReadSysMsgById(Id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.ErrorResponseAndLog(OPActionDelSystemMessages, ERROR_CODE_NO_SYSMSG, IdStr)
			return
		}
		c.ErrorResponseAndLog(OPActionDelSystemMessages, ERROR_CODE_DB, IdStr)
		return
	}

	err = dao.SystemMessageMethodDaoEntity.DelSysMsg(Id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelSystemMessages, ERROR_CODE_DEL_FAIL, IdStr)
		return
	}

	res := map[string]interface{}{}
	res["id"] = Id
	c.SuccessResponseAndLog(OPActionDelSystemMessages, IdStr, res)
}
