package eosplus

import (
	"common"
	"common/systoolbox"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/streadway/amqp"
	"otc_error"
	"strconv"
	"time"
	utilscommon "utils/common"
)

const (
	RabbitMQEusdTransfer      = "eusd.transfer"
	RabbitMQEusdTransferCheck = "eusd.check"
	RabbitMQEusdCheckResource = "eusd.resource"
)

type AmqpFuncContainer struct {
}

func (c AmqpFuncContainer) Default(delivery amqp.Delivery) (err error) {
	// do nothing
	common.LogFuncDebug("an unrouted message from %s:%s", delivery.Exchange, delivery.RoutingKey)
	return
}

//发起交易
func (c AmqpFuncContainer) RunTxLog(delivery amqp.Delivery) (err error) {
	idStr := string(delivery.Body)

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.LogFuncError("MQ RunTxLog:%v,%v", delivery.Body, err)
	}
	MQRunTxLog(id)
	return
}

// TaskBroadcast process task distribute message
func (c AmqpFuncContainer) TaskDistribute(delivery amqp.Delivery) (err error) {
	return systoolbox.TaskDistribute(delivery, FuncContainer)
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
				respMsg.Code = int(controllers.ERROR_CODE_ENCODE_FAILED)
				return
			}

			respMsg.Code = int(controllers.ERROR_CODE_SUCCESS)
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

//交易检查
func (c AmqpFuncContainer) RunTxCheck(delivery amqp.Delivery) (err error) {
	idStr := string(delivery.Body)

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.LogFuncError("MQ CheckTxLog:%v,%v", idStr, err)
	}
	MQCheckTransfer(id)
	return
}

//检查账号资源
func (c AmqpFuncContainer) CheckResource(delivery amqp.Delivery) (err error) {
	acc := string(delivery.Body)
	rpc := EosPlusAPI.Rpc
	key := "eos_res_" + acc
	// redis lock,每天评估一次
	intcmd := common.RedisManger.IncrBy(key, -1)
	num := intcmd.Val()
	if num > 1 {
		return
	}

	accInfo, err := rpc.GetAccount(acc)
	common.LogFuncError("%v", ToJsonIndent(accInfo))
	if err != nil {
		common.LogFuncError("CheckResource:%v,%v", acc, err)
	}
	// 如果 内存小于250byte 转账可能失败
	if accInfo.RAMQuota-accInfo.RAMUsage < 250 {
		errCode := rpc.transferCheckResourceBuyRam(acc)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			return
		}
	}

	need := false
	//转账需要<300µs的CPU时间
	if accInfo.CPULimit.Available < 1000 {
		rpc.RexRentCpu(acc, "0.0001 EOS")
		need = true
	}
	// 转账需要16B的NET
	if accInfo.NetLimit.Available < 300 {
		rpc.RexRentNet(acc, "0.0001 EOS")
		need = true
	}
	if need == true {
		//补充资源后，删除key
		common.RedisManger.Del(key)
		return
	}

	t1 := accInfo.CPULimit.Available / 300
	t2 := accInfo.NetLimit.Available / 50
	common.LogFuncError("%v,%v", t1, t2)
	if t1 > t2 {
		common.RedisManger.Set(key, int(t2), time.Second*86400)
	} else {
		common.RedisManger.Set(key, int(t1), time.Second*86400)
	}
	return
}
