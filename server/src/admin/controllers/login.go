package controllers

import (
	common2 "admin/common"
	"admin/controllers/errcode"
	"admin/utils"
	"common"
	"encoding/json"
	"utils/admin/dao"
	"utils/admin/models"
)

type LoginController struct {
	BaseController
}

//easyjson:json
type AdminLoginMsg struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	VerifyCode string `json:"verify_code"`
}

//登录
func (c *LoginController) HandleLogin() {
	req := &AdminLoginMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil || req.Email == "" {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}
	if !VerifyEmailFormat(req.Email) {
		c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_USER_NAME_ERROR, "")
		return
	}

	adminUser := &models.AdminUser{
		Name:  req.Email,
		Email: req.Email,
		Pwd:   req.Password,
	}
	err = dao.AdminUserDaoEntity.QueryAdminUser(adminUser, models.COLUMN_AdminUser_Email)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_NO_USER, "")
		return
	}

	if adminUser.Status != dao.AdminUserStatusActive {
		c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_USER_NO_ACTIVE, "")
		return
	}

	if adminUser.Pwd != "" {
		//密码验证
		pwd, err := common.GenerateDoubleMD5(req.Password, dao.AdminPwdSalt)
		if err != nil {
			c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_PWD_ERROR, "")
			return
		}
		if pwd != adminUser.Pwd {
			c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_PWD_ERROR, "")
			return
		}
	}

	if adminUser.SecretId == "" {
		//生成秘钥及对应二维码
		secretId, base64Data, err := utils.CreateGoogleAuthQrCode(adminUser.Email)
		if err != nil {
			c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_GEN_GOOGLE_AUTH_FAIL, "")
			return
		}
		if err := dao.AdminUserDaoEntity.UpdateGoogleAuth(adminUser.Id, secretId, base64Data); err != nil {
			c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_DB, "")
			return
		}
		adminUser.SecretId = secretId
		adminUser.QrCode = base64Data
	}

	if adminUser.IsBind && req.VerifyCode == "" {
		c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_VERITY_GOOGLE_CODE_FAIL, "")
		return
	}

	if req.VerifyCode != "" {
		//google验证码验证
		if ok, err := utils.VerityGoogleCode(adminUser.Email, adminUser.SecretId, req.VerifyCode); err != nil || !ok {
			c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_VERITY_GOOGLE_CODE_FAIL, "")
			return
		}

		adminUser.LoginTime = common.NowInt64MS()
		err = dao.AdminUserDaoEntity.UpdateTableAdminUser(adminUser, models.COLUMN_AdminUser_LoginTime)
		if err != nil {
			common.LogFuncError("update logintime err:%v", err)
		}

		if errCode := c.subLogin(adminUser.Id); errCode != controllers.ERROR_CODE_SUCCESS {
			c.ErrorResponseAndLog(OPActionAddLogin, errCode, "")
			return
		}
		if !adminUser.IsBind {
			if err := dao.AdminUserDaoEntity.UpdateGoogleAuthBind(adminUser.Id, true); err != nil {
				common.LogFuncError("error:%v", err)
				c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_DB, "")
				return
			}
			adminUser.IsBind = true
		}
	}

	member, err := dao.GetLoginInfo(*adminUser)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddLogin, controllers.ERROR_CODE_DB, "")
		return
	}

	//返回c查询结果
	c.SuccessResponseAndLog(OPActionAddLogin, "", member)
}

func (c *LoginController) subLogin(id uint64) controllers.ERROR_CODE {
	// reset access token
	accessToken, errCode := c.resetToken(id)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		return errCode
	}
	// set to cookie
	c.Ctx.SetCookie(TokenKey, accessToken, common2.Cursvr.AccessTokenExpiredSecs)
	return controllers.ERROR_CODE_SUCCESS
}

//登出
func (c *LoginController) HandleLogout() {
	_, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddLogout, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	c.ClearCookieToken()
	c.SuccessResponseWithoutDataAndLog(OPActionAddLogout, string(c.Ctx.Input.RequestBody))
}

//重新登录 更新token
/*func (c *LoginController) HandleReLogin() {
	uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS || uid == 0 {
		c.ErrorResponseAndLog(OPActionAddReLogin, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if errCode := c.subLogin(uid); errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddReLogin, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddReLogin, string(c.Ctx.Input.RequestBody))
}*/

//easyjson:json
type AdminChangePwdMsg struct {
	OldPwd string `json:"oldpwd"`
	NewPwd string `json:"newpwd"`
}

//修改密码
func (c *LoginController) HandlePassword() {
	adminUser, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddPassword, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &AdminChangePwdMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddPassword, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	if req.NewPwd == "" {
		c.ErrorResponseAndLog(OPActionAddPassword, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	err = dao.AdminUserDaoEntity.ChangePassword(&adminUser.AdminUser, req.OldPwd, req.NewPwd)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddPassword, controllers.ERROR_CODE_DB, "")
		return
	}

	c.ClearCookieToken()
	c.SuccessResponseWithoutDataAndLog(OPActionAddPassword, "")
}
