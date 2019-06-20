package controllers

import (
	"common"
	"common/systoolbox"
	"fmt"
	"github.com/mailru/easyjson"
	. "otc_error"
	"strconv"
	"usdt"
	utilscommon "utils/common"
	"utils/usdt/dao"
	"utils/usdt/models"

	"github.com/streadway/amqp"
)

// this container used for managing alternate exchange function

type AmqpFuncContainer struct {
}

//
func (c AmqpFuncContainer) Default(delivery amqp.Delivery) {
	//    do nothing
	common.LogFuncDebug("unprocessed : routing key %s , body %s", delivery.RoutingKey, string(delivery.Body))

}

func (c AmqpFuncContainer) Transfer(delivery amqp.Delivery) (err error) {

	defer func() {
		if err != nil {
			common.LogFuncDebug("transfer failed: routing key %s , body %s , error %v", delivery.RoutingKey, string(delivery.Body), err)
		}
	}()

	if delivery.RoutingKey != "usdt.transfer" {
		err = fmt.Errorf("wrong routing key")
		return
	}

	var (
		id  uint64
		log *models.UsdtWealthLog
	)
	if id, err = strconv.ParseUint(string(delivery.Body), 10, 64); err != nil {
		return
	}

	if log, err = dao.WealthLogDaoEntity.QueryById(id); err != nil {
		return
	}

	if _, errCode := usdt.TransferByUser(log.Uid, log.Id, log.To, log.AmountInteger, log.Memo); errCode != ERROR_CODE_SUCCESS {
		err = fmt.Errorf("transfer failed , ERROR_CODE : %v", errCode)
		return
	}
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
