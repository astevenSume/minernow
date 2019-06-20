package controllers

import (
	. "otc_error"
	admindao "utils/admin/dao"
)

type AppController struct {
	BaseController
}

/*func (c *AppController) Apps() {
	_, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	// get apps from admin
	var apps []adminmodels.Apps
	var err error
	apps, err = admindao.AppDaoEntity.AllPublished()
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	var channels []adminmodels.AppChannel
	channels, err = admindao.AppDaoEntity.AllChannels()
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	var types []adminmodels.AppType
	types, err = admindao.AppDaoEntity.AllTypes()
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		KeyList:           apps,
		KeyListAppChannel: channels,
		KeyListAppType:    types,
	})
}*/

func (c *AppController) Apps() {
	_, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	channelIdInput, err := c.GetUint32(KeyChannelIdInput, 0)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	// get apps from admin
	apps, err := admindao.AppDaoEntity.ChannelApp(channelIdInput)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	types, err := admindao.AppDaoEntity.AllTypes()
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		KeyList:        apps,
		KeyListAppType: types,
	})
}

func (c *AppController) Channels() {
	_, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	channels, err := admindao.AppDaoEntity.AllChannels()
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	appList, err := admindao.AppDaoEntity.AllFeatured()
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		KeyList:    channels,
		KeyListApp: appList,
	})
}
