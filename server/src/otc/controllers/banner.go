package controllers

import (
	. "otc_error"
	admindao "utils/admin/dao"
)

type BannerController struct {
	BaseController
}

func (c *BannerController) Banners() {
	_, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	// get apps from admin
	banner, err := admindao.BannerDaoEntity.EffectiveBanners()
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		KeyList: banner,
	})
}
