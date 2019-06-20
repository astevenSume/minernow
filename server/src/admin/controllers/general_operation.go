package controllers

import (
	"admin/controllers/errcode"
	GCommon "common"
	"eusd/eosplus"
	"strconv"
	"umeng_push/uemng_plus"
	"utils/common/dao"
	otcDao "utils/otc/dao"
	usdtDao "utils/usdt/dao"
)

type GeneralOperationController struct {
	BaseController
}

type userOp struct {
	Uid string `json:"uid"`
}

func (c *GeneralOperationController) getPost() (uid uint64, err error) {
	post := &userOp{}
	err = c.GetPost(post)
	if err != nil {
		return
	}
	uid, err = strconv.ParseUint(post.Uid, 10, 64)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAM_FAILED)
		return
	}
	return
}

//冻结账户
func (c *GeneralOperationController) FrozenUser() {
	c.setOPAction(OPActionGeneralFrozenUser)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	uid, err := c.getPost()
	if err != nil {
		return
	}
	// 冻结EUSD账号
	errCode2 := eosplus.EosPlusAPI.Wealth.Lock(uid)
	errCode = controllers.ERROR_CODE(errCode2)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	// 冻结承兑OTC账号
	errCode2 = eosplus.EosPlusAPI.Otc.Lock(uid)
	errCode = controllers.ERROR_CODE(errCode2)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	// 冻结USDT账号
	err = usdtDao.AccountDaoEntity.UpdateStatus(uid, usdtDao.STATUS_LOCKED)
	if err != nil {
		c.ErrorResponse(errCode)
		return
	}

	//推送通知
	go func() {
		title := "您的账号被冻结了!"
		content := "您的账号被冻结了，如有疑问，请联系客服。"
		p := new(uemng_plus.UPushPlus)
		p.PushSysNotice(uid, content, title)
		_, _ = otcDao.SystemNotificationdDaoEntity.InsertSystemNotification("system", content, uid)
	}()
	c.SuccessResponseWithoutData()
}

//解冻用户
func (c *GeneralOperationController) UnFrozenUser() {
	c.setOPAction(OPActionGeneralUnFrozenUser)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	uid, err := c.getPost()
	if err != nil {
		return
	}
	// 解冻EUSD账号
	errCode2 := eosplus.EosPlusAPI.Wealth.Unlock(uid)
	errCode = controllers.ERROR_CODE(errCode2)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	// 解冻承兑OTC账号
	errCode2 = eosplus.EosPlusAPI.Otc.Unlock(uid)
	errCode = controllers.ERROR_CODE(errCode2)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	// 解冻USDT账号
	err = usdtDao.AccountDaoEntity.UpdateStatus(uid, usdtDao.STATUS_WORKING)
	if err != nil {
		c.ErrorResponse(errCode)
		return
	}

	//推送通知
	go func() {
		title := "您的账号被解冻了!"
		content := "您的账号被解冻了"
		p := new(uemng_plus.UPushPlus)
		p.PushSysNotice(uid, content, title)
		_, _ = otcDao.SystemNotificationdDaoEntity.InsertSystemNotification("system", content, uid)
	}()
	c.SuccessResponseWithoutData()
}

//锁定用户(禁止登录)
func (c *GeneralOperationController) LockUser() {
	c.setOPAction(OPActionGeneralLockUser)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	uid, err := c.getPost()
	if err != nil {
		return
	}

	err = otcDao.UserDaoEntity.UpdateStatus(uid, otcDao.UserStatusSuspended)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	dao.TokenDaoEntity.RemoveToken(uid, GCommon.ClientTypeApp) //OC的APP用户TOKEN
	//推送通知
	go func() {
		title := "您的账号被禁用了!"
		content := "您的账号被禁用了，如有疑问，请联系客服。"
		p := new(uemng_plus.UPushPlus)
		p.PushSysNotice(uid, content, title)
		_, _ = otcDao.SystemNotificationdDaoEntity.InsertSystemNotification("system", content, uid)
	}()
	c.SuccessResponseWithoutData()
}

//解锁用户
func (c *GeneralOperationController) UnlockUser() {
	c.setOPAction(OPActionGeneralUnlockUser)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	uid, err := c.getPost()
	if err != nil {
		return
	}
	err = otcDao.UserDaoEntity.UpdateStatus(uid, otcDao.UserStatusActive)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	//推送通知
	go func() {
		title := "您的账号被解禁了!"
		content := "您的账号被解禁了。"
		p := new(uemng_plus.UPushPlus)
		p.PushSysNotice(uid, content, title)
		_, _ = otcDao.SystemNotificationdDaoEntity.InsertSystemNotification("system", content, uid)
	}()
	c.SuccessResponseWithoutData()
}
