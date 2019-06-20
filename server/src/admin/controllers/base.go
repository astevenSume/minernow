package controllers

import (
	. "admin/controllers/errcode"
	"common"
	"encoding/json"
	"github.com/astaxie/beego"
	"strconv"
	"utils/admin/dao"
	"utils/admin/models"
	utils "utils/common"
)

const (
	TokenKey = "token"
)

type BaseController struct {
	clientType  int
	adminId     uint64
	OpAction    int32
	RequestData string
	beego.Controller
}

// 设置操作
func (c *BaseController) setOPAction(op int32) {
	c.OpAction = op
}

// 设置请求参数
func (c *BaseController) setRequestData(data string) {
	c.RequestData = data
}

// 成功返回
func (c *BaseController) SuccessResponse(result interface{}) {
	params := map[string]interface{}{
		"code": ERROR_CODE_SUCCESS,
		"data": result,
	}

	_ = c.Ctx.Output.JSON(params, false, false)
	if c.RequestData == "" {
		c.setRequestData(string(c.Ctx.Input.RequestBody))
	}

	c.CreateLog(int32(ERROR_CODE_SUCCESS), c.OpAction, c.RequestData)
}

// 错误返回
func (c *BaseController) ErrorResponse(errCode ERROR_CODE) {
	params := map[string]interface{}{
		"code": errCode,
	}
	_ = c.Ctx.Output.JSON(params, false, false)
	// 没有参数信息，尝试加载post参数
	if c.RequestData == "" {
		c.setRequestData(string(c.Ctx.Input.RequestBody))
	}
	c.CreateLog(int32(ERROR_CODE_SUCCESS), c.OpAction, c.RequestData)
}

// 成功返回（无结果）
func (c *BaseController) SuccessResponseWithoutData() {
	params := map[string]interface{}{
		"code": ERROR_CODE_SUCCESS,
	}
	_ = c.Ctx.Output.JSON(params, false, false)
}

func (c *BaseController) SuccessResponseAndLog(action int32, reqPara string, result interface{}) {
	params := map[string]interface{}{
		"code": ERROR_CODE_SUCCESS,
		"data": result,
	}

	_ = c.Ctx.Output.JSON(params, false, false)
	c.CreateLog(int32(ERROR_CODE_SUCCESS), action, reqPara)
}

func (c *BaseController) ErrorResponseAndLog(action int32, errCode ERROR_CODE, reqPara string) {
	lang := "zh"
	msg := ErrorMsg(errCode, lang)
	params := map[string]interface{}{
		"code": errCode,
		"msg":  msg,
	}

	_ = c.Ctx.Output.JSON(params, false, false)
	c.CreateLog(int32(errCode), action, reqPara)
}

// 成功返回（无结果）
func (c *BaseController) SuccessResponseWithoutDataAndLog(action int32, reqPara string) {
	c.SuccessResponseWithoutData()
	c.CreateLog(int32(ERROR_CODE_SUCCESS), action, reqPara)
}

func (c *BaseController) ResponseAndLog(action int32, reqPara string, result interface{}) {
	_ = c.Ctx.Output.JSON(result, false, false)
	c.CreateLog(int32(ERROR_CODE_SUCCESS), action, reqPara)
}

func (c *BaseController) CreateLog(errCode int32, action int32, reqPara string) {
	log := &models.OperationLog{
		AdminId:      c.adminId,
		Method:       c.Ctx.Request.Method,
		Route:        c.Ctx.Input.URL(),
		Action:       action,
		Input:        reqPara,
		UserAgent:    c.Ctx.Input.UserAgent(),
		Ips:          common.ClientIP(c.Ctx),
		ResponseCode: errCode,
		Ctime:        common.NowInt64MS(),
	}
	err := dao.OperationLogDaoEntity.CreateOperationLog(log)
	if err != nil {
		common.LogFuncError("error:%v", err)
	}
}

// 生成cookie
func (c *BaseController) resetToken(id uint64) (t string, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	var err error
	t, err = utils.ResetToken(id, 0)
	if err != nil {
		errCode = ERROR_CODE_TOKEN_ERR
		return
	}

	return
}

// 清理cookie
func (c *BaseController) ClearCookieToken() {
	ok := utils.ClearToken(0, c.Ctx.GetCookie(TokenKey))
	if !ok {
		common.LogFuncError("ClearToken fail")
	}
	c.Ctx.SetCookie(TokenKey, "", -1)
	return
}

// 解析TOKEN
func (c *BaseController) getUidFromToken() (uid uint64, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	var (
		ok bool
	)
	ok, uid, _ = utils.CheckToken(0, c.Ctx.GetCookie(TokenKey))
	if !ok {
		errCode = ERROR_CODE_TOKEN_VERIFY_ERR
		return
	}
	c.adminId = uid

	return
}

//检查权限
func (c *BaseController) CheckPermission() (adminUser *dao.AdminDao, errCode ERROR_CODE) {
	_, errCode = c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS || c.adminId == 0 {
		return
	}

	var err error
	adminUser, err = dao.GetAdminInfo(c.adminId)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}

	//权限判断
	needPermission := GetUrlPermission(c.Ctx.Input.URL(), c.Ctx.Request.Method)
	if !adminUser.HasOneOf(needPermission...) {
		errCode = ERROR_CODE_NO_RIGHT_ERROR
		return
	}
	errCode = ERROR_CODE_SUCCESS

	return
}

func (c *BaseController) GetParamUint64(key string) (id uint64, err error) {
	idStr := c.Ctx.Input.Param(key)
	if len(idStr) > 0 {
		id, err = strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
			return
		}
	}

	return
}

func (c *BaseController) GetPost(post interface{}) (err error) {
	err = json.Unmarshal(c.Ctx.Input.RequestBody, post)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
	}
	return
}
