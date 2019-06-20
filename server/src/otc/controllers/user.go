package controllers

import (
	"common"
	"eusd/eosplus"
	"fmt"
	common2 "otc/common"
	controllers "otc_error"
	"strconv"
	"time"
	"usdt"
	agentdao "utils/agent/dao"
	common3 "utils/common"
	"utils/otc/dao"
	"utils/otc/models"

	json "github.com/mailru/easyjson"
)

type UserController struct {
	BaseController
}

//easyjson:json
type UserLoginMsg struct {
	NationalCode string `json:"national_code"`
	Mobile       string `json:"mobile"`
	Sms          string `json:"verify_code,omitempty"`
	InviteCode   string `json:"invite_code"`
	Nick         string `json:"nick"`
}

//登录 & 注册
func (c *UserController) Login() {
	msg := UserLoginMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	//验证短信验证码
	check, err := common3.VerifySmsCode(msg.NationalCode, msg.Mobile, common3.SmsActionLogin, msg.Sms)
	if !check {
		if err == common3.ErrSmsOutTimes {
			c.ErrorResponse(controllers.ERROR_CODE_SMS_VERIFY_FAIL_TOO_MATH)
		} else {
			c.ErrorResponse(controllers.ERROR_CODE_SMS_ERR)
		}
		return
	}

	//用户信息获取
	user, err := dao.UserDaoEntity.InfoByMobile(msg.NationalCode, msg.Mobile)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	first := false
	if user.Uid == 0 {
		//邀请码验证
		errCode := VerifyInviteCode(msg.InviteCode, msg.NationalCode, msg.Mobile)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			c.ErrorResponse(errCode)
			return
		}

		//未注册，先注册
		user, err = dao.UserDaoEntity.Create(msg.NationalCode, msg.Mobile, msg.Nick, c.Ctx.Request.RemoteAddr)
		if err != nil {
			c.ErrorResponse(controllers.ERROR_CODE_REGISTER_ERR)
			return
		}
		first = true
		//初始化wealth
		eosplus.EosPlusAPI.Wealth.Create(user.Uid, false)
		_, _ = usdt.CreateAccount(user.Uid)
	}

	// sub login: 返回用户信息 & token & signSalt
	if resp, errCode := c.subLogin(user, first); errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	} else {
		c.SuccessResponse2(&resp)
	}

	if first {
		// create new agent data
		err = agentdao.AgentPathDaoEntity.Create(user.Uid, msg.InviteCode)
		if err != nil {
			common.LogFuncError("create agent_path error:%v", err)
			return
		}
		err = agentdao.AgentDaoEntity.Create(user.Uid)
		if err != nil {
			common.LogFuncError("create agent_path error:%v", err)
			return
		}
	}

	//登录日志
	c.loginLog(user)
}

//重新登录 token更新
func (c *UserController) ReLogin() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	//用户信息获取
	user, err := dao.UserDaoEntity.InfoByUId(uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	// sub login: 返回用户信息 & token & signSalt
	if resp, errCode := c.subLogin(user, false); errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	} else {
		c.SuccessResponse2(&resp)
	}

	//登录日志
	c.loginLog(user)
}

//easyjson:json
type UserSignMsg struct {
	Uid       string `json:"uid"`
	Timestamp uint32 `json:"timestamp"` // server time secs
	Signature string `json:"signature"`
}

// 签名认证
func (c *UserController) Sign() {
	msg := UserSignMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	uid, err1 := strconv.ParseUint(msg.Uid, 10, 64)
	if err != nil || err1 != nil || uid == 0 || msg.Timestamp == 0 || msg.Signature == "" {
		common.LogFuncDebug("json decode: %s \nfailed : %v, params err: %d, %d, %s", string(c.Ctx.Input.RequestBody), err, msg.Uid, msg.Timestamp, msg.Signature)
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	expireSecs := common2.Cursvr.SignatureExpiredSecs
	nowSecs := common.NowUint32()
	// abs(timestamp - now) > expire
	if nowSecs+expireSecs < msg.Timestamp || msg.Timestamp+expireSecs < nowSecs {
		c.ErrorResponse(controllers.ERROR_CODE_LOGIN_SIGN_VERIFY_ERR)
		return
	}

	//用户信息获取
	user, err := dao.UserDaoEntity.InfoByUId(uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	if user.Uid == 0 || user.Mobile == "" {
		c.ErrorResponse(controllers.ERROR_CODE_NO_USER)
		return
	}

	// 签名认证
	if user.SignSalt == "" {
		c.ErrorResponse(controllers.ERROR_CODE_LOGIN_SIGN_VERIFY_ERR)
		return
	}
	src := map[string]string{
		"uid": msg.Uid,
	}
	if signature, err := common.AppSignMgr.GenerateMSign(src, msg.Timestamp, user.SignSalt); err != nil || signature != msg.Signature {
		c.ErrorResponse(controllers.ERROR_CODE_LOGIN_SIGN_VERIFY_ERR)
		return
	}

	// sub login: 返回用户信息 & token & signSalt
	if resp, errCode := c.subLogin(user, false); errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	} else {
		c.SuccessResponse(&resp)
	}

	//登录日志
	c.loginLog(user)
}

//easyjson:json
type UserLoginRespMsg struct {
	Uid           string `json:"uid"`
	NationalCode  string `json:"national_code"`
	Mobile        string `json:"mobile"`
	First         bool   `json:"first"`
	Status        int8   `json:"status"`
	Name          string `json:"name"`
	Ctime         int64  `json:"ctime"`
	Utime         int64  `json:"utime"`
	LastLoginTime int64  `json:"ltime"`
	IsExchanger   int8   `json:"exchanger"`
	SignSalt      string `json:"sign_salt"`
}

// sub login: 返回用户信息 & token & signSalt
//func (c *UserController) subLogin(user *models.User, first bool) (data map[string]interface{}, errCode controllers.ERROR_CODE) {
func (c *UserController) subLogin(user *models.User, first bool) (resp UserLoginRespMsg, errCode controllers.ERROR_CODE) {
	if user.Status != dao.UserStatusActive {
		errCode = controllers.ERROR_CODE_USER_NO_ACTIVE
		return
	}

	// reset access token
	var accessToken string
	accessToken, errCode = c.resetToken(user)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		errCode = controllers.ERROR_CODE_TOKEN_ERR
		return
	}

	// reset sign salt
	if u, eCode := c.resetSignSalt(user.Uid); eCode != controllers.ERROR_CODE_SUCCESS {
		errCode = eCode
		return
	} else {
		user = u
	}

	// set to cookie
	c.Ctx.SetCookie(TokenKey, accessToken, common2.Cursvr.AccessTokenExpiredSecs)

	user.LastLoginTime = common.NowInt64MS()
	err := dao.UserDaoEntity.UpdateLoginTime(user.Uid, user.LastLoginTime)
	if err != nil {
		common.LogFuncError("create login log fail err:%v", err)
	}
	//返回用户信息 & token
	resp.First = first
	resp.Uid = fmt.Sprint(user.Uid)
	resp.NationalCode = user.NationalCode
	resp.Mobile = user.Mobile
	resp.Name = user.Nick
	resp.Status = user.Status
	resp.IsExchanger = user.IsExchanger
	resp.Ctime = user.Ctime
	resp.LastLoginTime = user.LastLoginTime
	resp.Utime = user.Utime
	resp.SignSalt = user.SignSalt

	errCode = controllers.ERROR_CODE_SUCCESS

	return
}

func (c *UserController) loginLog(user *models.User) {
	//登录日志
	err := dao.UserLoginLogDaoEntity.Create(user.Uid, c.Ctx.Input.UserAgent(), common.ClientIP(c.Ctx))
	if err != nil {
		common.LogFuncError("create login log fail err:%v", err)
	}
}

//获取登录图像验证码
func (c *UserController) Captcha() {
	nationalCode := c.GetString("national_code")
	mobile := c.GetString("mobile")
	out := c.GetString("out")

	if out == "png" {
		image := common.GetCaptchaPng(nationalCode, mobile, common.CaptchaActionLogin)

		c.Ctx.Output.ContentType("png")
		_ = c.Ctx.Output.Body([]byte(image))
		return
	}

	image := common.GetCaptcha(nationalCode, mobile, common.CaptchaActionLogin)
	res := map[string]interface{}{}
	res["image"] = image
	c.SuccessResponse(res)
}

func (c *UserController) Logout() {
	c.ClearCookieToken()
	c.SuccessResponseWithoutData()
}

//获取用户信息
func (c *UserController) GetInfo() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	user, err := dao.UserDaoEntity.InfoByUId(uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	//返回用户信息
	data := map[string]interface{}{
		"uid":           fmt.Sprintf("%v", user.Uid),
		"national_code": user.NationalCode,
		"name":          user.Nick,
		"mobile":        user.Mobile,
		"status":        user.Status,
		"exchanger":     user.IsExchanger,
		"ctime":         user.Ctime,
		"ltime":         user.LastLoginTime,
		"utime":         user.Utime,
	}
	c.SuccessResponse(data)
}

//easyjson:json
type UserEditInfoMsg struct {
	Name string `json:"name"`
}

//编辑用户信息
func (c *UserController) EditInfo() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := UserEditInfoMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil || msg.Name == "" {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	user, err := dao.UserDaoEntity.EditNickByUId(uid, msg.Name)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	//返回用户信息
	data := map[string]interface{}{
		"uid":           fmt.Sprintf("%v", user.Uid),
		"national_code": user.NationalCode,
		"name":          user.Nick,
		"mobile":        user.Mobile,
		"status":        "",
		"exchanger":     user.IsExchanger,
		"ctime":         user.Ctime,
		"ltime":         user.LastLoginTime,
		"utime":         user.Utime,
	}
	c.SuccessResponse(data)
}

//用户设置
func (c *UserController) Setting() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	info := dao.UserConfigDaoEntity.Info(uid)

	c.SuccessResponse(info)
}

//easyjson:json
type UserSetSettingMsg struct {
	WealthNotice bool `json:"wealth_notice"`
	OrderNotice  bool `json:"order_notice"`
}

// 设置用户设置
func (c *UserController) SetSetting() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := UserSetSettingMsg{}
	err := c.GetPost(&msg)
	if err != nil {
		return
	}

	dao.UserConfigDaoEntity.Update(uid, msg.WealthNotice, msg.OrderNotice)

	c.SuccessResponseWithoutData()
}

//easyjson:json
type UpdateMobileInfo struct {
	Password string `json:"password"`
}

// 更换手机,验证密码
func (c *UserController) CheckPassWord() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	fmt.Println(uid)
	updateInfo := UpdateMobileInfo{}
	err := c.GetPost(&updateInfo)
	if err != nil {
		return
	}
	isValid, passValidator := dao.UserPayPassDaoEntity.ValidateByPassword(uid, updateInfo.Password, time.Duration(common2.Cursvr.SignatureTokenExpiredSecs)*time.Second)

	res := map[string]interface{}{}
	if isValid == true {
		res["is_valid"] = true
		res["pass_validator"] = passValidator
	} else {
		res["is_valid"] = false
	}
	c.SuccessResponse(res)
}

//easyjson:json
type NewMobileInfo struct {
	Token     string `json:"token"`
	NewMobile string `json:"new_mobile"`
}

// 输入新手机号,发送短信验证码
func (c *UserController) GetNewMobile() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	NewMobileInfo := NewMobileInfo{}
	err := c.GetPost(&NewMobileInfo)
	if err != nil {
		return
	}

	// 检查认证token
	if isValid := dao.UserPayPassDaoEntity.ValidateCacheToken(uid, NewMobileInfo.Token, false); !isValid {
		return
	}
	// 判断新手机号是否已存在
	num, err := dao.UserDaoEntity.NewMobileIsPresence(NewMobileInfo.NewMobile)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	if num != 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NEW_MOBILE_ALREADY_PRESENCE)
		return
	}

	user, err := dao.UserDaoEntity.InfoByUId(uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	// 发送验证码
	_, errCode = common3.AliSendSms(user.NationalCode, NewMobileInfo.NewMobile, common3.SmsActionPaySecond)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	res := map[string]interface{}{}
	res["send_is_valid"] = "success"
	c.SuccessResponse(res)
}

//easyjson:json
type VerificationCodeInfo struct {
	NewMobile        string `json:"new_mobile"`
	VerificationCode string `json:"verification_code,omitempty"`
}

// 验证短信验证码
func (c *UserController) CheckVerificationCode() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	user, err := dao.UserDaoEntity.InfoByUId(uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	info := VerificationCodeInfo{}
	err = c.GetPost(&info)
	if err != nil {
		return
	}
	isVerify, err := common3.VerifySmsCode(user.NationalCode, info.NewMobile, common3.SmsActionPaySecond, info.VerificationCode)
	common.LogFuncDebug("CheckVerificationCode isVerify err %v", err)
	res := map[string]interface{}{}
	if !isVerify {
		// 短信验证
		if err == common3.ErrSmsOutTimes {
			c.ErrorResponse(controllers.ERROR_CODE_SMS_VERIFY_FAIL_TOO_MATH)
			res["is_verify"] = controllers.ERROR_CODE_SMS_VERIFY_FAIL_TOO_MATH
		} else {
			c.ErrorResponse(controllers.ERROR_CODE_SMS_ERR)
			res["is_verify"] = controllers.ERROR_CODE_SMS_ERR
		}
	} else {
		err = dao.UserDaoEntity.UpdateMobile(uid, info.NewMobile)
		if err != nil {
			c.ErrorResponse(controllers.ERROR_CODE_USER_UPDATE_MOBILE_FAIL)
			return
		}
		res["is_verify"] = controllers.ERROR_CODE_SUCCESS
		res["update_mobile"] = "success"
	}
	c.SuccessResponse(res)
}
