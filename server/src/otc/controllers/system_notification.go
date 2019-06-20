package controllers

import (
	"common"
	. "otc_error"
	"strconv"
	"strings"
	otc_dao "utils/otc/dao"
)

type SystemNotificationController struct {
	BaseController
}

func (c *SystemNotificationController) GetSystemNotification() {
	_, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	uid_c, err := c.GetUint64("uid")
	if err != nil {
		c.ErrorResponse(ERROR_CODE_SUCCESS)
		return
	}

	page, err := c.GetInt("page", 1)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_SUCCESS)
		return
	}
	limit, err := c.GetInt("limit", 10)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_SUCCESS)
		return
	}
	snList, err := otc_dao.SystemNotificationdDaoEntity.GetSystemNotificationByUid(uid_c, page, limit)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_SYSTEM_NOTIFICATION_FALIED)
		return
	}
	type Meta struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
		Total int `json:"total"`
	}

	go common.SafeRun(func() {
		var nidList []string
		for _, sn := range snList {
			nidList = append(nidList, strconv.Itoa(int(sn.Nid)))
		}
		nidListStr := strings.Join(nidList, ", ")
		otc_dao.SystemNotificationdDaoEntity.UpdateSystemNotificationByNid(nidListStr)
	})()
	meta := &Meta{}
	meta.Total = len(snList)
	meta.Limit = limit
	meta.Page = page

	res := map[string]interface{}{}
	res["list"] = snList
	res["meta"] = meta
	c.SuccessResponse(res)
}

func (c *SystemNotificationController) GetIsReadNum() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	num, err := otc_dao.SystemNotificationdDaoEntity.GetSystemNotificationIsRead(uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_SYSTEM_NOTIFICATION_ISREAD_FALIED)
		return
	}
	res := map[string]interface{}{}
	res["num"] = num
	c.SuccessResponse(res)
}
