package controllers

import (
	"common"
	json "github.com/mailru/easyjson"
	"otc/trade"
	. "otc_error"
	admindao "utils/admin/dao"
	"utils/otc/dao"
)

type AppealController struct {
	BaseController
}

//easyjson:json
type AppealCreateAppealMsg struct {
	Type    int8   `json:"type"`
	Context string `json:"context"`
	WeChat  string `json:"wechat"`
}

//创建申诉
func (c *AppealController) CreateAppeal() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := AppealCreateAppealMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil || msg.Type <= dao.AppealTypeNil || msg.Type >= dao.AppealTypeMax {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	orderId, _ := c.GetUint64(":order_id", 0)
	if orderId <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	order, err := dao.OrdersDaoEntity.Info(orderId)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	if order.Uid != uid && order.EUid != uid {
		c.ErrorResponse(ERROR_CODE_NO_AUTH)
		return
	}
	if order.Status != dao.OrderStatusCanceled && order.Status != dao.OrderStatusPayed {
		c.ErrorResponse(ERROR_CODE_CAN_NOT_APPEAL)
		return
	}

	oldAppeal, _ := dao.AppealDaoEntity.QueryByOrderId(orderId)
	if oldAppeal.Id > 0 {
		c.ErrorResponse(ERROR_CODE_APPEAL_EXIST)
		return
	}

	adminID, weChat, qrCode := admindao.AppealServiceDaoEntity.GetAppealAdminInfo()
	if adminID == 0 {
		c.ErrorResponse(ERROR_CODE_APPEAL_SERVICE_LESS)
		return
	}

	appeal, err := dao.AppealDaoEntity.Create(msg.Type, adminID, orderId, uid, msg.WeChat, msg.Context)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	err = dao.OrdersDaoEntity.SetOrderAppealStatus(orderId, uid, dao.OrderAppealStatusPending, qrCode)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_OTC_ORDER_NOT_FOUND)
		return
	}
	data := dao.AppealDaoEntity.ClientAppeal(&appeal)
	data.ContactWechat = weChat
	data.QrCode = qrCode

	c.SuccessResponse(data)
}

//获取申诉
func (c *AppealController) GetAppeal() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	orderId, _ := c.GetUint64(":order_id", 0)
	if orderId <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	order, err := dao.OrdersDaoEntity.Info(orderId)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	if order.Uid != uid && order.EUid != uid {
		c.ErrorResponse(ERROR_CODE_NO_AUTH)
		return
	}

	appeal, err := dao.AppealDaoEntity.QueryByOrderId(orderId)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_OTC_ORDER_NOT_FOUND)
		return
	}

	data := dao.AppealDaoEntity.ClientAppeal(&appeal)
	data.ContactWechat, data.QrCode = admindao.AppealServiceDaoEntity.GetContactByAdminID(appeal.AdminId)

	c.SuccessResponse(data)
}

//解决申诉
func (c *AppealController) ResolveAppeal() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	orderId, _ := c.GetUint64(":order_id", 0)
	if orderId <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	errCode = trade.CheckAppealOrder(uid, orderId)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	appeal, err := dao.AppealDaoEntity.Resolve(orderId, uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(dao.AppealDaoEntity.ClientAppeal(&appeal))
}
