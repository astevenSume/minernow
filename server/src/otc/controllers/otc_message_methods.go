package controllers

import (
	"common"
	"fmt"
	json "github.com/mailru/easyjson"
	otc_common "otc/common"
	. "otc_error"
	"strconv"
	"umeng_push/uemng_plus"
	"utils/otc/dao"
)

type OtcMessageMethodController struct {
	BaseController
}

func (c *OtcMessageMethodController) Get() {
	uid, err1 := c.getUidFromToken()
	if err1 != ERROR_CODE_SUCCESS {
		c.ErrorResponse(ERROR_CODE_NO_LOGIN)
		return
	}

	id, err := c.GetUint64("id")
	if err != nil {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}
	if id < 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}
	direction := c.GetString("direction")
	if direction != "up" && direction != "down" {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}
	all, err := c.GetInt("all", 0)

	if err != nil {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}
	if all != 1 {
		all = 0
	}

	var orderId int64

	orderIdStr := c.Ctx.Input.Param(KeyOrderIdInput)
	if len(orderIdStr) > 0 {
		var err error
		orderId, err = strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
			return
		}
	}

	order, err := dao.OrdersDaoEntity.Info(uint64(orderId))
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	if order.Uid != uid && order.EUid != uid {
		c.ErrorResponse(ERROR_CODE_NO_AUTH)
		return
	}
	res := map[string]interface{}{}
	if id == 0 {
		//取最新30条记录
		msgList, err := dao.MessageMethodDaoEntity.QueryByOrderIdZero(orderId)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_DB)
			return
		}
		res["list"] = msgList
	} else if id > 0 && direction == "up" {
		//取id往上10条，返回倒序
		if all == 0 {
			msgList, err := dao.MessageMethodDaoEntity.QueryUpMsg(id, orderId)
			if err != nil {
				c.ErrorResponse(ERROR_CODE_DB)
				return
			}
			res["list"] = msgList
		} else if all == 1 {
			msgList, err := dao.MessageMethodDaoEntity.QueryUpMsgAll(id, orderId)

			if err != nil {
				c.ErrorResponse(ERROR_CODE_DB)
				return
			}
			res["list"] = msgList
		}
	} else if id > 0 && direction == "down" {
		//取id往下30条，返回顺序
		msgList, err := dao.MessageMethodDaoEntity.QueryDownMsg(id, orderId)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_DB)
			return
		}
		res["list"] = msgList
	}
	c.SuccessResponse(res)

}

//easyjson:json
type OtcMessageMethodAddMsg struct {
	OrderId uint64 `json:"order_id"`
	Uid     uint64 `json:"uid"`
	MsgType string `json:"msg_type"`
	Content string `json:"content"`
}

func (c *OtcMessageMethodController) AddMsg() {
	uid, err1 := c.getUidFromToken()
	if err1 != ERROR_CODE_SUCCESS {
		c.ErrorResponse(ERROR_CODE_NO_LOGIN)
		return
	}

	data := OtcMessageMethodAddMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &data)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}
	order, err := dao.OrdersDaoEntity.Info(data.OrderId)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	if order.Uid != uid && order.EUid != uid {
		c.ErrorResponse(ERROR_CODE_NO_AUTH)
		return
	}

	msg, err := dao.MessageMethodDaoEntity.Add(int64(data.OrderId), uid, data.MsgType, data.Content)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	//向另一方推送消息
	p := new(uemng_plus.UPushPlus)
	title := otc_common.NewMessage
	body := fmt.Sprintf(msg.Content)
	go common.SafeRun(func() {
		if order.Uid == uid {
			//通知承兑商
			p.PushMessage(order.EUid, order.Id, body, title, "1", order.Status)
		} else if uid == order.EUid {
			//通知用户
			p.PushMessage(order.Uid, order.Id, body, title, "0", order.Status)
		}
	})()

	res := map[string]interface{}{}
	res["msg"] = msg
	c.SuccessResponse(res)
}

func (c *OtcMessageMethodController) CommonMessages() {
	_, err := c.getUidFromToken()
	if err != ERROR_CODE_SUCCESS {
		c.ErrorResponse(ERROR_CODE_NO_LOGIN)
		return
	}

	res := map[string]interface{}{}
	var created = [...]string{"您已成功下单，请及时付款。", "已成功下单，请等待买家付款。"}
	var unconfirmed = [...]string{"你已付款，请等待卖家确认收款。", "买家已付款，请及时查收并确认收款。"}
	var confirmed = [...]string{"卖家已确认收款，系统正在转账代币到您的钱包。", "您已确认收款，系统正在转账代币到买家的钱包。"}
	var finished = [...]string{"转账完成，您已收到您购买的代币。", "转账完成，买家已收到您出售的代币。"}
	var buyer_canceled = [...]string{"您已取消订单，请不要再付款给卖家。", "买家取消订单，如收到买家付款，请及时退还。"}
	var seller_canceled = [...]string{"卖家取消订单，请不要再付款给卖家。", "您已取消订单，如收到买家付款，请及时退还。"}
	var admin_canceled = [...]string{"客服取消订单，请不要再付款给卖家。", "客服取消订单，如收到买家付款，请及时退还。"}
	var pay_expired = [...]string{"您的订单因为超时已被系统取消，如有疑问，请联系客服。", "买家付款超过时间限制，订单已被系统取消。"}
	var confirm_expired = [...]string{"卖家确认收款时间已过期，可以通过申诉联系客服。", "您的确认收款时间已过期，请尽快确认。"}
	var buyer_appealed = [...]string{"您已成功提交申诉，请联系客服处理。", "买家申诉，请配合客服沟通协调。"}
	var seller_appealed = [...]string{"卖家申诉，请配合客服沟通协调。", "您已成功提交申诉，请联系客服处理。"}
	var buyer_appeal_resolved = [...]string{"您已取消申诉，请继续完成订单。", "买家已取消申诉，您可以继续完成订单。"}
	var seller_appeal_resolved = [...]string{"卖家已取消申诉，您可以继续完成订单。", "您已取消申诉，请继续完成订单。"}

	res["created"] = created
	res["unconfirmed"] = unconfirmed
	res["confirmed"] = confirmed
	res["finished"] = finished
	res["buyer_canceled"] = buyer_canceled
	res["seller_canceled"] = seller_canceled
	res["admin_canceled"] = admin_canceled
	res["pay_expired"] = pay_expired
	res["confirm_expired"] = confirm_expired
	res["buyer_appealed"] = buyer_appealed
	res["seller_appealed"] = seller_appealed
	res["buyer_appeal_resolved"] = buyer_appeal_resolved
	res["seller_appeal_resolved"] = seller_appeal_resolved
	c.SuccessResponse(res)
}
