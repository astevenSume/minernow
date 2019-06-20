package controllers

import (
	"common"
	controllers "otc_error"
	admindao "utils/admin/dao"
	agentdao "utils/agent/dao"
	utils "utils/common"
	"utils/otc/dao"

	json "github.com/mailru/easyjson"
)

type SmsCodeController struct {
	BaseController
}

//easyjson:json
type SmsSendCodeMsg struct {
	NationalCode string `json:"national_code"`
	Mobile       string `json:"mobile"`
	Action       string `json:"action"`
	InviteCode   string `json:"invite_code"`
}

// 发送短信验证码
func (c *SmsCodeController) SendCode() {
	msg := SmsSendCodeMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json.Unmarshal this.Ctx.Input.RequestBody %s failed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	first := false
	if msg.Action == utils.SmsActionLogin {
		//邀请码验证
		if !admindao.TopAgentDaoEntity.IsTopAgent(msg.NationalCode, msg.Mobile) {
			if msg.InviteCode != "" {
				_, err := agentdao.AgentPathDaoEntity.GetUidByInviteCode(msg.InviteCode)
				if err != nil {
					c.ErrorResponse(controllers.ERROR_CODE_USER_INVITE_CODE_ERROR)
					return
				}
			} else {
				user, err := dao.UserDaoEntity.InfoByMobile(msg.NationalCode, msg.Mobile)
				if err != nil {
					c.ErrorResponse(controllers.ERROR_CODE_DB)
					return
				}
				if user.Uid == 0 {
					//新用户需要有邀请码
					first = true
				}
			}
		}
	}

	if !first {
		//发送短信验证码
		_, errCode := utils.AliSendSms(msg.NationalCode, msg.Mobile, msg.Action)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			c.ErrorResponse(errCode)
			return
		}
	}

	//返回用户信息 & token
	res := map[string]interface{}{
		"first": first,
	}
	c.SuccessResponse(res)
}

//easyjson:json
type SmsCodeUserSendCodeMsg struct {
	Action string `json:"Action"`
}

// 登录用户发送短信验证码
func (c *SmsCodeController) UserSendCode() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := SmsCodeUserSendCodeMsg{}
	err := c.GetPost(&msg)
	if err != nil {
		return
	}

	user, err := dao.UserDaoEntity.InfoByUId(uid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	if user.Status != dao.UserStatusActive {
		c.ErrorResponse(controllers.ERROR_CODE_USER_NO_ACTIVE)
		return
	}

	_, errCode = utils.AliSendSms(user.NationalCode, user.Mobile, msg.Action)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	c.SuccessResponseWithoutData()
}
