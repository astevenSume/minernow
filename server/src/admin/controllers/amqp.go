package controllers

import (
	"common"
	"common/systoolbox"
	"encoding/json"
	"github.com/mailru/easyjson"
	"github.com/streadway/amqp"
	"utils/admin/dao"
	"utils/admin/models"
)

// this container used for managing alternate exchange function
type AmqpFuncContainer struct {
}

//
func (c AmqpFuncContainer) Default(delivery amqp.Delivery) (err error) {
	//

	common.LogFuncDebug("unrouted message %s.%s", delivery.Exchange, delivery.RoutingKey)

	return
}

// Ping process ping message from servers
func (c AmqpFuncContainer) Ping(delivery amqp.Delivery) (err error) {
	pingMsg := common.PingMsg{}
	err = json.Unmarshal(delivery.Body, &pingMsg)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	//common.LogFuncDebug("Received ping message from %v", pingMsg)
	switch pingMsg.PingType {
	case common.PingMsgTypeRegister:
		{
			err = c.pingProcess(pingMsg.AppName, pingMsg.RegionId, pingMsg.ServerId)
			if err != nil {
				return
			}

			//send task list to server node
			err = sendTask(pingMsg.AppName, pingMsg.RegionId, pingMsg.ServerId)
			if err != nil {
				return
			}
		}
	case common.PingMsgTypePing:
		{
			err = c.pingProcess(pingMsg.AppName, pingMsg.RegionId, pingMsg.ServerId)
			if err != nil {
				return
			}
		}
	default:
		common.LogFuncError("unkown ping message type %d", pingMsg.PingType)
	}

	return
}

func (c AmqpFuncContainer) pingProcess(appName string, regionId, serverId int64) (err error) {
	//update server node last ping timestamp
	err = dao.ServerNodeDaoEntity.UpdateLastPing(appName, regionId, serverId, common.NowUint32())
	if err != nil {
		return
	}

	// add task producer for server node
	businessName := svrTaskBusinessName(appName, regionId, serverId)
	exchange := businessName
	err = common.RabbitMQAddProducerIfNoExist(businessName, exchange)
	if err != nil {
		return
	}

	// add rpc client for server node
	err = common.RabbitMQAddRpcClientIfNoExist(appName, regionId, serverId)
	if err != nil {
		return
	}

	return
}

// TaskResult process task result from servers
func (c AmqpFuncContainer) TaskResult(delivery amqp.Delivery) (err error) {
	result := systoolbox.TaskResult{}

	err = easyjson.Unmarshal(delivery.Body, &result)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	//common.LogFuncDebug("received task result %+v", result)

	// store result
	err = dao.TaskResultDaoEntity.Add(models.TaskResult{
		AppName:   result.AppName,
		RegionId:  result.RegionId,
		ServerId:  result.ServerId,
		Name:      result.Name,
		Code:      result.Code,
		Detail:    result.Detail,
		BeginTime: result.BeginTime,
		EndTime:   result.EndTime,
		Ctime:     result.Ctime,
	})
	if err != nil {
		return
	}

	return
}

func sendTask(appName string, regionId, serverId int64) (err error) {
	var total int
	var page = 1
	total, err = sendSinglePageTask(appName, 1, regionId, serverId)
	if err != nil {
		return
	}

	for page = 2; page <= total/dao.DefaultPerPage; page++ {
		_, err = sendSinglePageTask(appName, page, regionId, serverId)
		if err != nil {
			return
		}
	}

	if total > dao.DefaultPerPage &&
		total%dao.DefaultPerPage > 0 {
		_, err = sendSinglePageTask(appName, page, regionId, serverId)
		if err != nil {
			return
		}
	}

	return
}

func sendSinglePageTask(appName string, page int, regionId, serverId int64) (total int, err error) {
	var (
		taskList    []models.Task
		meta        dao.PageInfo
		taskMsgList = systoolbox.TaskMsgList{}
		buf         []byte
	)
	meta, taskList, err = dao.TaskDaoEntity.Query(appName, page, dao.DefaultPerPage)
	if err != nil {
		return
	}

	total = meta.Total

	if total <= 0 {
		return
	}

	for _, dbTask := range taskList {
		taskMsgList.List = append(taskMsgList.List, systoolbox.TaskMsg{
			Name:     dbTask.Name,
			FuncName: dbTask.FuncName,
			Spec:     dbTask.Spec,
		})
	}

	buf, err = easyjson.Marshal(&taskMsgList)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	//send task list to server
	err = common.RabbitMQPublish(svrTaskBusinessName(appName, regionId, serverId),
		common.RabbitMQRoutingKeyTaskDistribute,
		buf,
	)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// Rpc rabbitmq rpc process function
func (c AmqpFuncContainer) Rpc(delivery amqp.Delivery) (buf []byte, err error) {
	switch delivery.RoutingKey {

	}

	//4debug
	return []byte("OK"), nil

	return
}

func (c AmqpFuncContainer) RpcClient(delivery amqp.Delivery) (err error) {
	return
}
