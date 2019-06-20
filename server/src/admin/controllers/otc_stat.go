package controllers

import (
	"admin/controllers/cron"
	"admin/controllers/errcode"
	"common"
	"utils/admin/dao"
	"utils/admin/models"
)

type OtcStatController struct {
	BaseController
}

func (c *OtcStatController) Get() {
	c.setOPAction(OPActionOtcStat)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	page, _ := c.GetInt64("page", 1)
	limit, _ := c.GetInt64("limit", 20)

	count := dao.OtcStatDaoEntity.Count()

	m := common.MakeMeta(count, page, limit)

	list := dao.OtcStatDaoEntity.Fetch(m.Offset, m.Limit)

	today := &models.OtcStat{}
	if m.Page == 1 {
		today = cron.OtcStat(0)
	}

	c.SuccessResponse(map[string]interface{}{
		"list":  list,
		"meta":  m,
		"today": today,
	})
}

func (c *OtcStatController) User() {
	c.setOPAction(OPActionOtcStat)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	uid, err := c.GetUint64("uid")
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}
	res := cron.OtcStatPeople(uid)
	c.SuccessResponse(res)
}
