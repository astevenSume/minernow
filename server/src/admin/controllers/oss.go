package controllers

import (
	common2 "admin/common"
	"admin/controllers/errcode"
	"common"
)

type OssController struct {
	BaseController
}

func (c *OssController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadOss, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	/*uid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadOss, errCode, uid,  "")
		return
	}*/

	token, err := common.GetPolicyToken(common2.Cursvr.OssAccessKeyId, common2.Cursvr.OssAccessKeySecret, common2.Cursvr.OssHost, "", common2.Cursvr.OssCallbackUrl, common2.Cursvr.OssExpireTime)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOss, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	c.SuccessResponseAndLog(OPActionReadOss, "", token)
}

//回调处理
func (c *OssController) Post() {
	// Get PublicKey bytes
	bytePublicKey, err := common.GetPublicKey(c.Ctx.Request)
	if err != nil {
		return
	}

	// Get Authorization bytes : decode from Base64String
	byteAuthorization, err := common.GetAuthorization(c.Ctx.Request)
	if err != nil {
		return
	}

	// Get MD5 bytes from Newly Constructed Authrization String.
	byteMD5, err := common.GetMD5FromNewAuthString(c.Ctx.Request)
	if err != nil {
		common.ResponseFailed(c.Ctx.ResponseWriter)
		return
	}

	// verifySignature and response to client
	if common.VerifySignature(bytePublicKey, byteMD5, byteAuthorization) {
		common.ResponseSuccess(c.Ctx.ResponseWriter) // response OK : 200
	} else {
		common.ResponseFailed(c.Ctx.ResponseWriter) // response FAILED : 400
	}
}
