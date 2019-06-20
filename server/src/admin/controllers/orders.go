package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	otcerror "otc_error"
	"utils/otc/dao"
	otcdao "utils/otc/dao"
)

type OrdersController struct {
	BaseController
}

func (c *OrdersController) List() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadOrder, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	side, err := c.GetInt8("side", -1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOrder, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	status, err := c.GetUint8("status", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOrder, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	statusList := []interface{}{}
	if status > 0 {
		statusList = append(statusList, status)
	}

	appealStatus, err := c.GetUint8("appeal_status", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOrder, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	appealStatusList := []interface{}{}
	if appealStatus > 0 {
		appealStatusList = append(appealStatusList, appealStatus)
	}
	mobile := c.GetString(KEY_MOBILE)

	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOrder, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	limit, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOrder, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	if limit > MaxPerPage {
		limit = MaxPerPage
	}
	input := fmt.Sprintf("{\"mobile\":\"%v\",\"side\":%v,\"appeal_status\":%v,\"appeal_status\":%v,\"page\":%v,"+
		"\"limit\":%v}", mobile, side, status, appealStatus, page, limit)

	otcUid, errCode := GetOtcUidByMobile(mobile)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadOrder, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	total, list, err := dao.OrdersDaoEntity.FetchDetailByUid(otcUid, side, status, appealStatus, (page-1)*limit, limit)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOrder, controllers.ERROR_CODE_DB, input)
		return
	}

	meta := PageInfo{
		Limit: limit,
		Total: int(total),
		Page:  page,
	}
	c.SuccessResponseAndLog(OPActionReadOrder, input, map[string]interface{}{
		"list": list,
		"meta": meta,
	})
}

func (c *OrdersController) Info() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadOrder, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	oid, _ := c.GetParamUint64(":order_id")
	if oid < 1 {
		c.ErrorResponseAndLog(OPActionReadOrder, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"order_id\":%v}", oid)
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOrder, controllers.ERROR_CODE_DB, input)
		return
	}
	var data dao.OrderAck
	if order.Id > 0 {
		data = dao.OrdersDaoEntity.ClientOrder(order)
		data.UMobile, data.UQrCode = dao.UserDaoEntity.GetUserContact(order.Uid, order.PayId)
		data.EuMobile, data.EuQrCode = dao.UserDaoEntity.GetUserContact(order.EUid, order.EPayId)
	}

	c.SuccessResponseAndLog(OPActionReadOrder, input, data)
}

//easyjson:json
type adminOrderReq struct {
	Id  string `json:"id"`
	Uid string `json:"uid"`
}

//easyjson:json
type otcOrderReq struct {
	Id        string `json:"id"`
	Uid       string `json:"uid"`
	AdminId   uint32 `json:"admin_id"`
	Timestamp uint32 `json:"timestamp"`
}

type otcOrderAck struct {
	Code otcerror.ERROR_CODE `json:"code"`
	Msg  string              `json:"msg"`
}

func (c *OrdersController) ResolveOrder() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddOrderResolve, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &adminOrderReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddOrderResolve, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUid, err := strconv.ParseUint(req.Uid, 10, 64)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddOrderResolve, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	appealID, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddOrderResolve, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	appeal, err := otcdao.AppealDaoEntity.QueryById(appealID)
	if err != nil || appeal.OrderId == 0 {
		common.LogFuncError("err:%v", err)
		c.ErrorResponseAndLog(OPActionAddOrderResolve, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	if appeal.Status == otcdao.AppealStatusResolved {
		c.ErrorResponseAndLog(OPActionAddOrderResolve, controllers.ERROR_CODE_OTC_ORDER_APPEAL_RESOLVED, string(c.Ctx.Input.RequestBody))
		return
	}

	order, err := dao.OrdersDaoEntity.Info(appeal.OrderId)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddOrderResolve, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	if order.Status != otcdao.OrderStatusCanceled && order.Status != otcdao.OrderStatusExpired {
		c.ErrorResponseAndLog(OPActionAddOrderResolve, controllers.ERROR_CODE_OTC_ORDER_STATUS_ERR, string(c.Ctx.Input.RequestBody))
		return
	}

	_, err = dao.AppealDaoEntity.Resolve(appeal.OrderId, otcUid)
	if err != nil {
		common.LogFuncError("err:%v", err)
		c.ErrorResponseAndLog(OPActionAddOrderResolve, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	_, err = otcdao.AppealDealLogDaoEntity.Create(otcdao.AppealDealLogActionResolve, uint32(c.adminId), appealID, appeal.OrderId)
	if err != nil {
		common.LogFuncError("err:%v", err)
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddOrderResolve, string(c.Ctx.Input.RequestBody))
}

//客服申述取消订单(买家点已付款但实际未付款)
func (c *OrdersController) CancelOrder() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddOrderCancel, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &adminOrderReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddOrderCancel, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	//otc请求包
	reqOtc := &otcOrderReq{
		Id:        req.Id,
		Uid:       req.Uid,
		AdminId:   uint32(c.adminId),
		Timestamp: uint32(time.Now().Unix()),
	}

	//otc响应数据
	ack := otcOrderAck{}

	//otc post请求
	errCode = PostOtc(RouterCancelOrder, map[string]string{
		"id":       reqOtc.Id,
		"uid":      reqOtc.Uid,
		"admin_id": fmt.Sprintf("%d", reqOtc.AdminId),
	}, reqOtc, &ack, reqOtc.Timestamp)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddOrderCancel, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if ack.Code != otcerror.ERROR_CODE_SUCCESS {
		common.LogFuncError("errCode:%v", ack.Code)
		c.ResponseAndLog(OPActionAddOrderCancel, string(c.Ctx.Input.RequestBody), ack)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddOrderCancel, string(c.Ctx.Input.RequestBody))
}

//客服申述确认已付款(买家点取消或系统超时自动取消)
func (c *OrdersController) ConfirmOrderPay() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddOrderConfirmPay, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &adminOrderReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddOrderConfirmPay, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//otc请求包
	reqOtc := &otcOrderReq{
		Id:        req.Id,
		Uid:       req.Uid,
		AdminId:   uint32(c.adminId),
		Timestamp: uint32(time.Now().Unix()),
	}
	//otc响应数据
	ack := otcOrderAck{}

	//otc post请求
	errCode = PostOtc(RouterConfirmOrderPay, map[string]string{
		"id":       reqOtc.Id,
		"uid":      reqOtc.Uid,
		"admin_id": fmt.Sprintf("%d", reqOtc.AdminId),
	}, reqOtc, &ack, reqOtc.Timestamp)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddOrderConfirmPay, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if ack.Code != otcerror.ERROR_CODE_SUCCESS {
		common.LogFuncError("errCode:%v", ack.Code)
		c.ResponseAndLog(OPActionAddOrderConfirmPay, string(c.Ctx.Input.RequestBody), ack)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddOrderConfirmPay, string(c.Ctx.Input.RequestBody))
}

//客服申述确认放币(买家已付款卖家未确认)
func (c *OrdersController) ConfirmOrder() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddOrderConfirm, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &adminOrderReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddOrderConfirm, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//otc请求包
	reqOtc := &otcOrderReq{
		Id:        req.Id,
		Uid:       req.Uid,
		AdminId:   uint32(c.adminId),
		Timestamp: uint32(time.Now().Unix()),
	}
	//otc响应数据
	ack := otcOrderAck{}

	//otc post请求
	errCode = PostOtc(RouterConfirmOrder, map[string]string{
		"id":       reqOtc.Id,
		"uid":      reqOtc.Uid,
		"admin_id": fmt.Sprintf("%d", reqOtc.AdminId),
	}, reqOtc, &ack, reqOtc.Timestamp)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddOrderConfirm, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if ack.Code != otcerror.ERROR_CODE_SUCCESS {
		common.LogFuncError("errCode:%v", ack.Code)
		c.ResponseAndLog(OPActionAddOrderConfirm, string(c.Ctx.Input.RequestBody), ack)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddOrderConfirm, string(c.Ctx.Input.RequestBody))
}
