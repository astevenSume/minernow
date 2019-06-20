package controllers

import (
	. "otc_error"
	admindao "utils/admin/dao"
)

type AnnouncementController struct {
	BaseController
}

func (c *AnnouncementController) Announcements() {
	_, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	typeInput, err := c.GetInt8(KeyTypeInput, 0)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}
	// get apps from admin
	banner, err := admindao.AnnouncementDaoEntity.EffectiveAnnouncement(typeInput)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		KeyList: banner,
	})
}
