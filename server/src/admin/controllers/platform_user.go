package controllers

import (
	"admin/controllers/errcode"
	"common"
	"eusd/eosplus"
	"fmt"
	"strconv"
	"utils/eusd/dao"
	"utils/eusd/models"
	dao2 "utils/otc/dao"
	models2 "utils/otc/models"
)

type PlatformUserController struct {
	BaseController
}

//获取用户list
func (c *PlatformUserController) List() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadOss, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	list, err := dao.PlatformUserDaoEntity.All()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadPlatformUser, controllers.ERROR_CODE_DB, "")
	}

	type platformUser struct {
		*models.EosWealth
		*models.PlatformUserCate
	}
	//读取wealth的信息
	uids := []uint64{}

	for _, v := range list {
		uids = append(uids, v.Uid)
	}
	users, err := dao.WealthDaoEntity.FetchByIds(uids...)
	usersTmp := map[uint64]*models.EosWealth{}
	for _, v := range users {
		usersTmp[v.Uid] = v
	}

	cateList, err := dao.PlatformUserCateDaoEntity.All()
	cateListTmp := map[int32]*models.PlatformUserCate{}
	for _, v := range cateList {
		cateListTmp[v.Id] = v
	}

	resList := []*platformUser{}
	for _, v := range list {
		tmp := &platformUser{}
		tmp.EosWealth = usersTmp[v.Uid]
		tmp.PlatformUserCate = cateListTmp[v.Pid]
		resList = append(resList, tmp)
	}

	c.SuccessResponseAndLog(OPActionReadPlatformUser, "", map[string]interface{}{"list": resList})
}

// 新增用户
func (c *PlatformUserController) Add() {
	c.setOPAction(OPActionAddPlatformUserAdd)
	c.setRequestData(string(c.Ctx.Input.RequestBody))
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddPlatformUserAdd, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type post struct {
		Pid int32  `json:"pid"`
		Uid string `json:"uid"`
	}

	msg := post{}

	err := c.GetPost(&msg)
	if err != nil {
		return
	}
	uid, err := strconv.ParseUint(msg.Uid, 10, 64)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAM_FAILED)
		return
	}
	user := &models2.User{}
	if uid > 0 {
		user, err = dao2.UserDaoEntity.InfoByUId(uid)
		if err != nil {
			c.ErrorResponse(controllers.ERROR_CODE_DB)
			return
		}
	} else {
		//创建登录账号
		phone := fmt.Sprintf("9%d", common.NowUint32())
		user, err = dao2.UserDaoEntity.Create("1", phone, phone, "")
		if err != nil {
			c.ErrorResponse(controllers.ERROR_CODE_DB)
			return
		}
	}

	//创建账号资产
	eosplus.EosPlusAPI.Wealth.Create(user.Uid, true)

	// 写入平台表里
	_, err = dao.PlatformUserDaoEntity.Add(user.Uid, msg.Pid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddPlatformUserAdd, string(c.Ctx.Input.RequestBody))
}

// 设置用户状态
func (c *PlatformUserController) Status() {
	c.setOPAction(OPActionAddPlatformUserStatus)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddPlatformUserStatus, errCode, string(c.Ctx.Input.RequestBody))
		return
	}
	type post struct {
		Uid    string `json:"uid"`
		Status int8   `json:"status"`
	}

	msg := post{}

	err := c.GetPost(&msg)
	if err != nil {
		return
	}
	uid, err := strconv.ParseUint(msg.Uid, 10, 64)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAM_FAILED)
		return
	}
	// 写入平台表里
	if msg.Status == 1 {
		err = dao.PlatformUserDaoEntity.Active(uid)
	} else {
		err = dao.PlatformUserDaoEntity.Lock(uid)
	}
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddPlatformUserStatus, string(c.Ctx.Input.RequestBody))
}
