package controllers

import (
	"common"
	"fmt"
	common2 "otc/common"
	. "otc_error"
	"strconv"
	"strings"
	utils "utils/common"
	"utils/otc/dao"
	"utils/otc/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	json "github.com/mailru/easyjson"

	json2 "encoding/json"
)

const (
	TokenKey = "token"
)

var (
	AesKey = [32]byte{'T', 'r', 'W', '2', '4', 'f', 'x', 'M',
		'P', '6', 'Q', 'L', '7', 's', 'B', '1',
		'G', 'H', 'c', 'r', 'l', 'K', 'R', 'O',
		'A', 'U', 'C', '3', 'D', 'S', 'c', '='}
)

type BaseController struct {
	clientType int
	beego.Controller
}

// encryptResult 加密返回数据
func (c *BaseController) encryptResult(data interface{}) (interface{}, ERROR_CODE) {

	var (
		buf       []byte
		err       error
		res       interface{}
		marshaler bool
	)

	// TODO: 上线前删除
	if flag := c.Ctx.Request.Header.Get("enable-encrypt"); flag == "" {

		if _, ok := data.(json.Marshaler); ok {

			if buf, err = json.Marshal(data.(json.Marshaler)); err != nil {
				return data, ERROR_CODE_ENCODE_FAILED
			}
			return string(buf), ERROR_CODE_SUCCESS
		}

		return data, ERROR_CODE_SUCCESS
	}
	// TODO: 上线前删除

	// 实现了 json.Marshaler 的数据
	if _, ok := data.(json.Marshaler); ok {
		marshaler = true
		if buf, err = json.Marshal(data.(json.Marshaler)); err != nil {
			return data, ERROR_CODE_ENCODE_FAILED
		}
	} else if buf, err = json2.Marshal(data); err != nil {
		return data, ERROR_CODE_ENCODE_FAILED
	}
	// 对数据进行 aes-gcm 加密
	if res, err = common.EncryptToBase64(string(buf), AesKey); err != nil {
		return data, ERROR_CODE_ENCRYPT_FAILED
	}

	if marshaler {
		if buf, err = json2.Marshal(res); err != nil {
			return data, ERROR_CODE_ENCODE_FAILED
		}
		return string(buf), ERROR_CODE_SUCCESS
	}

	return res, ERROR_CODE_SUCCESS
}

// 成功返回
// this is old func for interface{} input, no recommanded.
func (c *BaseController) SuccessResponse(result interface{}) {

	data, errCode := c.encryptResult(result)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	params := map[string]interface{}{
		"code": ERROR_CODE_SUCCESS,
		"data": data,
	}

	_ = c.Ctx.Output.JSON(params, false, false)
}

// 成功返回
// this is old func for interface{} input, no recommanded.
func (c *BaseController) SuccessResponse2(result json.Marshaler) {

	data, errCode := c.encryptResult(result)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	_ = c.Ctx.Output.Body([]byte(fmt.Sprintf("{\"code\":%d, \"data\":%s}", ERROR_CODE_SUCCESS, data)))
}

// 成功返回（无结果）
func (c *BaseController) SuccessResponseWithoutData() {
	_ = c.Ctx.Output.Body([]byte(fmt.Sprintf("{\"code\":%d}", ERROR_CODE_SUCCESS)))
}

//错误返回
func (c *BaseController) ErrorResponse(errCode ERROR_CODE) {
	lang := "zh"
	msg := ErrorMsg(errCode, lang)
	_ = c.Ctx.Output.Body([]byte(fmt.Sprintf("{\"code\":%d, \"msg\":\"%s\"}", errCode, msg)))
}

//easyjson:json
type BaseCheck2StepMsg struct {
	Token      string `json:"token"`
	VerifyCode string `json:"verify_code,omitempty"`
}

// 二次认证,forceTwiceVerify:是否必须二次认证
func (c *BaseController) check2step(uid uint64, forceTwiceVerify bool) ERROR_CODE {
	req := BaseCheck2StepMsg{}
	if err := c.GetPost(&req); err != nil {
		return ERROR_CODE_PARAMS_ERROR
	}

	info, err := dao.UserPayPassDaoEntity.QueryPwdByUid(uid)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		return ERROR_CODE_DB
	}

	//密码是否可用
	if info.Status != dao.UserPayPwdActive {
		return ERROR_CODE_PASSWORD_UNAVAILABLE_ERR
	}

	if dao.UserPayPassDaoEntity.IsLockPwd(uid) {
		return ERROR_CODE_PASSWORD_LOCK
	}

	if req.Token == "" {
		return ERROR_CODE_PASSWORD_TOKEN_ERR
	}

	// 检查认证token
	if isValid := dao.UserPayPassDaoEntity.ValidateCacheToken(uid, req.Token, false); !isValid {
		return ERROR_CODE_PASSWORD_TOKEN_ERR
	}

	// 需要二次认证
	if forceTwiceVerify || info.VerifyStep == dao.UserPayPwdStepTwice {
		user, err := dao.UserDaoEntity.InfoByUId(uid)
		if err != nil {
			return ERROR_CODE_DB
		}
		//短信验证
		check, _ := utils.VerifySmsCode(user.NationalCode, user.Mobile, utils.SmsActionPaySecond, req.VerifyCode)
		common.LogFuncDebug("check check:%d, sms:%s", check, req.VerifyCode)
		if !check {
			// 短信验证
			if err == utils.ErrSmsOutTimes {
				c.ErrorResponse(ERROR_CODE_SMS_VERIFY_FAIL_TOO_MATH)
				return ERROR_CODE_SMS_VERIFY_FAIL_TOO_MATH
			} else {
				c.ErrorResponse(ERROR_CODE_SMS_ERR)
				return ERROR_CODE_SMS_VERIFY_FAIL_TOO_MATH
			}
		}
	}

	// 验证成功，清除token
	dao.UserPayPassDaoEntity.DelCacheToken(uid)
	return ERROR_CODE_SUCCESS
}

// generate signature salt by uid
func (c *BaseController) resetSignSalt(uid uint64) (user *models.User, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	clientType := c.ClientType()
	if clientType == common2.ClientTypeUnkown {
		common.LogFuncWarning("get %s failed", KeyClientTypeCookie)
		errCode = ERROR_CODE_CLIENT_TYPE_UNKOWN
		return
	}

	var err error
	user, err = dao.UserDaoEntity.ResetSignSaltByUId(uid)
	if err != nil {
		if err == orm.ErrNoRows {
			errCode = ERROR_CODE_NO_USER
			return
		}
		errCode = ERROR_CODE_DB
		return
	}

	return
}

// 生成cookie
func (c *BaseController) resetToken(u *models.User) (t string, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	clientType := c.ClientType()
	if clientType == common2.ClientTypeUnkown {
		common.LogFuncWarning("get %s failed", KeyClientTypeCookie)
		errCode = ERROR_CODE_CLIENT_TYPE_UNKOWN
		return
	}

	var err error
	t, err = utils.ResetToken(u.Uid, clientType)
	if err != nil {
		errCode = ERROR_CODE_TOKEN_ERR
		return
	}

	return
}

// 清理cookie
func (c *BaseController) ClearCookieToken() (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	clientType := c.ClientType()
	if clientType == common2.ClientTypeUnkown {
		errCode = ERROR_CODE_CLIENT_TYPE_UNKOWN
		return
	}
	ok := utils.ClearToken(clientType, c.Ctx.GetCookie(TokenKey))
	if !ok {
		common.LogFuncError("ClearToken fail")
	}
	c.Ctx.SetCookie(TokenKey, "", -1)
	return
}

// 解析TOKEN
func (c *BaseController) getUidFromToken() (uid uint64, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	clientType := c.ClientType()
	if clientType == common2.ClientTypeUnkown {
		errCode = ERROR_CODE_CLIENT_TYPE_UNKOWN
		return
	}

	var (
		ok bool
	)
	ok, uid, _ = utils.CheckToken(clientType, c.Ctx.GetCookie(TokenKey))
	if !ok {
		errCode = ERROR_CODE_TOKEN_VERIFY_ERR
		return
	}

	return
}

// @Description get client type
func (c BaseController) ClientType() int {
	if c.clientType == common2.ClientTypeUnkown {

		s := c.Ctx.Input.Header(KeyClientTypeCookie)
		if len(s) <= 0 {
			return common2.ClientTypeUnkown
		}

		c.clientType = common2.GetClientType(s)
		if c.clientType == common2.ClientTypeUnkown {
			common.LogFuncError("clientType %s unkown.", s)
			return c.clientType
		}
	}

	return c.clientType
}

func (c *BaseController) GetPost(post json.Unmarshaler) (err error) {
	err = json.Unmarshal(c.Ctx.Input.RequestBody, post)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
	}
	return
}

// @Description get remote ip address
func (c *BaseController) GetIP() string {
	ip := c.Ctx.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = c.Ctx.Request.Header.Get("X-real-ip")
	}

	if ip == "" {
		return "127.0.0.1"
	}

	return ip
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
