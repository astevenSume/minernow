package controllers

import (
	"common"
	json "github.com/mailru/easyjson"
	common3 "otc/common"
	"otc_error"
	"regexp"
	"time"
	common2 "utils/common"
	"utils/otc/dao"
)

type UserPayPassController struct {
	BaseController
}

//easyjson:json
type UserPayPassSetPasswordMsg struct {
	Method      int    `json:"method"`
	Msm         string `json:"verify_code,omitempty"`
	OldPassword string `json:"oldpwd,omitempty"`
	NewPassword string `json:"newpwd"`
}

// set payment password by message code & old password
func (c *UserPayPassController) SetPassword() {

	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	// 判断用户是否存在
	user, err := dao.UserDaoEntity.InfoByUId(uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	if user.Uid == 0 || user.Mobile == "" || user.NationalCode == "" {
		c.ErrorResponse(controllers.ERROR_CODE_NO_USER)
		return
	}

	msg := UserPayPassSetPasswordMsg{}
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	if !(msg.Method > dao.UserPayPwdMethodUnknown && msg.Method < dao.UserPayPwdMethodMax) {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	// 判断原密码是否可用
	pwd, err := dao.UserPayPassDaoEntity.QueryPwdByUid(uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	if msg.Method == dao.UserPayPwdMethodPass && pwd.Status != dao.UserPayPwdActive {
		c.ErrorResponse(controllers.ERROR_CODE_PASSWORD_UNAVAILABLE_ERR)
		return
	}
	if dao.UserPayPassDaoEntity.IsLockPwd(uid) {
		c.ErrorResponse(controllers.ERROR_CODE_PASSWORD_LOCK)
		return
	}
	// 默认不开启二次认证
	verifyStep := int(pwd.VerifyStep)
	if verifyStep != dao.UserPayPwdStepTwice {
		verifyStep = dao.UserPayPwdStepOnce
	}

	type Ack struct {
		SignSalt   string `json:"sign_salt"`
		Status     int    `json:"status"`
		VerifyStep int    `json:"verify_step"`
	}

	// 旧密码修改方式验证
	ack := Ack{}
	if msg.Method == dao.UserPayPwdMethodPass && pwd.Status == dao.UserPayPwdActive {
		if errCode := checkPassword(msg.NewPassword); errCode != controllers.ERROR_CODE_SUCCESS {
			c.ErrorResponse(errCode)
			return
		}
		if msg.NewPassword == msg.OldPassword {
			c.ErrorResponse(controllers.ERROR_CODE_PASSWORD_REPEAT_ERR)
			return
		}

		if pass := dao.UserPayPassDaoEntity.ValidatePassword(uid, msg.OldPassword); !pass {
			dao.UserPayPassDaoEntity.VerifyPayPassResult(uid, false)
			c.ErrorResponse(controllers.ERROR_CODE_PASSWORD_ERR)
			return
		}
	}

	// 手机验证码修改方式验证
	if msg.Method == dao.UserPayPwdMethodSms {
		check, _ := common2.VerifySmsCode(user.NationalCode, user.Mobile, common2.SmsAcitionPayPassword, msg.Msm)
		if !check {
			if err == common2.ErrSmsOutTimes {
				c.ErrorResponse(controllers.ERROR_CODE_SMS_VERIFY_FAIL_TOO_MATH)
			} else {
				c.ErrorResponse(controllers.ERROR_CODE_SMS_ERR)
			}
			dao.UserPayPassDaoEntity.VerifyPayPassResult(uid, false)
			return
		}
	}

	if msg.Method != dao.UserPayPwdMethodSms && msg.Method != dao.UserPayPwdMethodPass {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	// 验证成功，设置密码到数据库
	pwd, err = dao.UserPayPassDaoEntity.InsertOrUpdatePwd(uid, msg.NewPassword, int8(msg.Method), int8(verifyStep))
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	ack.SignSalt = pwd.SignSalt
	ack.Status = int(pwd.Status)
	ack.VerifyStep = int(pwd.VerifyStep)
	c.SuccessResponse(ack)
	dao.UserPayPassDaoEntity.VerifyPayPassResult(uid, true)
	return
}

// check password format ^\d{6}$
func checkPassword(pwd string) controllers.ERROR_CODE {
	reg := regexp.MustCompile(`^\d{6}$`)
	strs := reg.FindAllString(pwd, -1)
	if len(strs) == 0 {
		return controllers.ERROR_CODE_PASSWORD_FORMAT_ERR
	}
	return controllers.ERROR_CODE_SUCCESS
}

//easyjson:json
type UserPayPassSetVerifyStepMsg struct {
	VerifyStep int `json:"verify_step"`
}

// 设置二次验证
func (c *UserPayPassController) SetVerifyStep() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	req := UserPayPassSetVerifyStepMsg{}
	if err := c.GetPost(&req); err != nil {
		return
	}

	if req.VerifyStep != dao.UserPayPwdStepOnce && req.VerifyStep != dao.UserPayPwdStepTwice {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	info, err := dao.UserPayPassDaoEntity.QueryPwdByUid(uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	if info.Status != dao.UserPayPwdActive {
		c.ErrorResponse(controllers.ERROR_CODE_PASSWORD_UNAVAILABLE_ERR)
		return
	}
	if dao.UserPayPassDaoEntity.IsLockPwd(uid) {
		c.ErrorResponse(controllers.ERROR_CODE_PASSWORD_LOCK)
		return
	}

	errCode = c.check2step(uid, true)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	err = dao.UserPayPassDaoEntity.UpdateVerifyStep(uid, int8(req.VerifyStep))
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	type Ack struct {
		Status     int `json:"status"`
		VerifyStep int `json:"verify_step"`
	}
	ack := Ack{
		Status:     int(info.Status),
		VerifyStep: int(req.VerifyStep),
	}
	c.SuccessResponse(ack)
}

type UserPayPassGetPayPwdStatusMsg struct {
	Status     int  `json:"status"`
	VerifyStep int  `json:"verify_step"`
	IsLock     bool `json:"is_lock"`
	LockSecond int  `json:"lock_second"`
	FailCnt    int  `json:"fail_cnt"`
}

// get payment password status
func (c *UserPayPassController) GetPayPwdStatus() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	status, step, err := dao.UserPayPassDaoEntity.QueryStatusByUid(uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	ack := UserPayPassGetPayPwdStatusMsg{
		Status:     status,
		VerifyStep: step,
	}
	ack.IsLock, ack.FailCnt, ack.LockSecond = dao.UserPayPassDaoEntity.GetLockInfo(uid)

	c.SuccessResponse(ack)
}

//easyjson:json
type UserPayPassVerifyBySignMsg struct {
	Timestamp uint32 `json:"timestamp"` // server time secs
	Signature string `json:"signature"`
}

// payment sign authorication
func (c *UserPayPassController) VerifyBySign() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := UserPayPassVerifyBySignMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	expireSecs := common3.Cursvr.SignatureTokenExpiredSecs
	nowSecs := common.NowUint32()
	// abs(timestamp - now) > expire
	if nowSecs+expireSecs < msg.Timestamp || msg.Timestamp+expireSecs < nowSecs {
		c.ErrorResponse(controllers.ERROR_CODE_PAY_SIGN_VERIFY_ERR)
		return
	}

	isValid, validator := dao.UserPayPassDaoEntity.ValidateBySign(uid, msg.Signature, msg.Timestamp, time.Duration(common3.Cursvr.SignatureTokenExpiredSecs)*time.Second)
	if !isValid {
		c.ErrorResponse(controllers.ERROR_CODE_PAY_SIGN_VERIFY_ERR)
		return
	}

	c.SuccessResponse(validator)
}

//easyjson:json
type UserPayPassVerifyByPasswordMsg struct {
	Password string `json:"password"`
}

// payment password authorication
func (c *UserPayPassController) VerifyByPassword() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := UserPayPassVerifyByPasswordMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	if dao.UserPayPassDaoEntity.IsLockPwd(uid) {
		c.ErrorResponse(controllers.ERROR_CODE_PASSWORD_LOCK)
		return
	}

	isValid, validator := dao.UserPayPassDaoEntity.ValidatePasswordAndRecord(uid, msg.Password, time.Duration(common3.Cursvr.SignatureTokenExpiredSecs)*time.Second)
	if !isValid {
		c.ErrorResponse(controllers.ERROR_CODE_PASSWORD_ERR)
		return
	}

	c.SuccessResponse(validator)
}
