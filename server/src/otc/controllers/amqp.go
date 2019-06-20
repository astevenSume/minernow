package controllers

import (
	"common"
	"common/systoolbox"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/streadway/amqp"
	"otc/trade"
	. "otc_error"
	"strconv"
	utilscommon "utils/common"
	"utils/otc/dao"
)

const (
	RabbitMQExchangeUsdtTransfer = "usdt.transfer"
	RabbitMQOrderCheck           = "order.check"
)

// this container used for managing alternate exchange function
type AmqpFuncContainer struct {
}

//
func (c AmqpFuncContainer) Default(delivery amqp.Delivery) (err error) {
	// do nothing
	common.LogFuncDebug("an unrouted message from %s:%s", delivery.Exchange, delivery.RoutingKey)
	return
}

// TaskBroadcast process task distribute message
func (c AmqpFuncContainer) TaskDistribute(delivery amqp.Delivery) (err error) {
	return systoolbox.TaskDistribute(delivery, &FunctionContainer{})
}

// Rpc process rabbitmq rpc request
func (c AmqpFuncContainer) Rpc(delivery amqp.Delivery) (msg []byte, err error) {
	common.LogFuncDebug("Received rpc request from %s : %s", delivery.ReplyTo, delivery.CorrelationId)

	// process specific request
	// decode
	reqMsg := common.RabbitMQRpcMsg{}
	err = easyjson.Unmarshal(delivery.Body, &reqMsg)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	respMsg := common.RabbitMQRpcMsg{
		Cmd: reqMsg.Cmd,
	}

	switch reqMsg.Cmd {
	case utilscommon.RabbitMQRpcMsgCmdTaskQuerySingle: //query task detail from single server node
		{
			detail := systoolbox.TaskMgr.Detail()
			var buf []byte
			buf, err = easyjson.Marshal(&detail)
			if err != nil {
				common.LogFuncError("%v", err)
				respMsg.Code = int(ERROR_CODE_ENCODE_FAILED)
				return
			}

			respMsg.Code = int(ERROR_CODE_SUCCESS)
			respMsg.Body = buf
		}
	default:
		common.LogFuncWarning("unkown cmd %d", reqMsg.Cmd)
		err = fmt.Errorf("unkown cmd %d", reqMsg.Cmd)
		return
	}

	msg, err = easyjson.Marshal(&respMsg)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (c AmqpFuncContainer) str2Uint64(delivery amqp.Delivery) (id uint64, err error) {
	idStr := string(delivery.Body)
	id, err = strconv.ParseUint(idStr, 10, 64)
	return
}

//Order check paid or cancel order
func (c AmqpFuncContainer) OrderCheck(delivery amqp.Delivery) (err error) {
	id, err := c.str2Uint64(delivery)
	if err != nil {
		common.LogFuncError("OrderCheck str2id err:%v,%v", string(delivery.Body), err)
		return
	}
	order, err := dao.OrdersDaoEntity.Info(id)
	if err != nil {
		common.LogFuncError("OrderCheck not found:%v,%v", string(delivery.Body), err)
		return
	}
	if order.Status != dao.OrderStatusCreated { //状态已经改变
		return
	}
	if order.Side == dao.SideBuy {
		trade.BuyLogicEntity.TimeOutOrder(order.Id)
	} else if order.Side == dao.SideSell {
		trade.SellLogicEntity.TimeOutOrder(order.Id)
	}
	return
}
