package controllers

import (
	"common"
	"otc_error"
	"strings"
	dao2 "utils/admin/dao"
	"utils/otc/dao"
	"utils/otc/models"
)

type OrdersController struct {
	BaseController
}

//所有订单
func (c *OrdersController) Get() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}
	side, _ := c.GetInt8("side")
	status := c.GetString("status")
	appealStatus := c.GetString("appeal_status")
	page, _ := c.GetInt64("page")
	limit, _ := c.GetInt64("limit")

	statusList := []interface{}{}
	if status != "" {
		tmp := strings.Split(status, ",")
		for _, v := range tmp {
			statusList = append(statusList, v)
		}
	}

	appealStatusList := []interface{}{}
	if status != "" {
		tmp := strings.Split(appealStatus, ",")
		for _, v := range tmp {
			appealStatusList = append(appealStatusList, v)
		}
	}

	num, err := dao.OrdersDaoEntity.Count(uid, side, statusList, appealStatusList)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	meta := common.MakeMeta(num, page, limit)

	list := []*models.OtcOrder{}
	if meta.Offset < meta.Total {
		list = dao.OrdersDaoEntity.FetchByUid(uid, side, statusList, appealStatusList, meta.Offset, meta.Limit)
	}
	c.SuccessResponse(map[string]interface{}{
		"list": list,
		"meta": meta,
	})
}

//获取订单详情
func (c *OrdersController) Info() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}
	oid, _ := c.GetParamUint64(":order_id")
	if oid < 1 {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}
	order, err := dao.OrdersDaoEntity.Info(oid)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	if order.Id == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_OTC_ORDER_NOT_FOUND)
		return
	}
	if order.Uid != uid && order.EUid != uid {
		c.ErrorResponse(controllers.ERROR_CODE_NO_AUTH)
		return
	}

	type OrderAck struct {
		*models.OtcOrder
		UMobile     string `json:"u_mobile"`
		UQrCode     string `json:"u_qr_code"`
		EuMobile    string `json:"eu_mobile"`
		EuQrCode    string `json:"eu_qr_code"`
		AppealId    uint64 `json:"appeal_id"`
		AdminName   string `json:"admin_name"`
		AdminWechat string `json:"admin_wechat"`
		AdminQr     string `json:"admin_qr"`
		AppealTime  int64  `json:"appeal_time"`
		AppealUId   uint64 `json:"appeal_uid"`
	}
	data := OrderAck{OtcOrder: order}
	data.UMobile, data.UQrCode = dao.UserDaoEntity.GetUserContact(order.Uid, order.PayId)
	data.EuMobile, data.EuQrCode = dao.UserDaoEntity.GetUserContact(order.EUid, order.EPayId)
	if data.AppealStatus != dao.OrderAppealStatusNil {
		appealData, err := dao.AppealDaoEntity.QueryByOrderId(data.Id)
		data.AppealId = appealData.Id
		data.AppealUId = appealData.UserId
		data.AppealTime = appealData.Ctime
		if err == nil {
			admin, err := dao2.AppealServiceDaoEntity.QueryById(data.AdminId)
			if err == nil {
				data.AdminName = admin.AdminNick
				data.AdminQr = admin.QrCode
				data.AdminWechat = admin.Wechat
			}
		}
	}

	c.SuccessResponse(data)
}

func (c *OrdersController) Exchanger() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}
	side, _ := c.GetInt8("side")
	status, _ := c.GetUint8("status")
	payType, _ := c.GetUint8("pay_type")
	page, _ := c.GetInt64("page")
	limit, _ := c.GetInt64("limit")
	date, _ := c.GetInt32("date")

	num, err := dao.OrdersDaoEntity.CountExchange(uid, side, status, payType, date)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}

	meta := common.MakeMeta(num, page, limit)

	list := []*models.OtcOrder{}
	total := []*dao.OrderExchangeTotal{}
	if meta.Offset < meta.Total {
		list = dao.OrdersDaoEntity.FetchExchangeByUid(uid, side, status, payType, date, meta.Offset, meta.Limit)
		dates := []int32{}
		temp := map[int32]struct{}{}

		for _, item := range list {
			if _, ok := temp[item.Date]; !ok {
				temp[item.Date] = struct{}{}
				dates = append(dates, item.Date)
			}
		}
		tmp, err := dao.OrdersDaoEntity.GroupFinishByUidDates(uid, side, payType, dates)
		if err != nil {
			c.ErrorResponse(controllers.ERROR_CODE_DB)
			return
		}
		for _, d := range dates {
			t := &dao.OrderExchangeTotal{
				Date: d,
			}

			for _, v := range tmp {
				if v.Date == d {
					t.Amount = v.Amount
					t.Funds = v.Funds
					break
				}
			}

			total = append(total, t)
		}

	}
	c.SuccessResponse(map[string]interface{}{
		"list":  list,
		"total": total,
		"meta":  meta,
	})
}
