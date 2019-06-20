package common

import (
	"github.com/mailru/easyjson"
)

const (
	// business name
	RabbitMQBusinessNameSvrPing    = "server.ping"
	RabbitMQBusinessNameTaskResult = "task.result"
	RabbitMQBusinessNameTask       = "task"

	// exchange name
	RabbitMQExchangePing       = "server.ping"
	RabbitMQExchangeTaskResult = "task.result"
	RabbitMQExchangeTask       = "task"

	// task routing key
	RabbitMQRoutingKeyTaskRun        = "run"
	RabbitMQRoutingKeyTaskSwitch     = "switch"
	RabbitMQRoutingKeyTaskDelete     = "delete"
	RabbitMQRoutingKeyTaskDistribute = "distibute"
)

const (
	PingMsgTypePing = iota
	PingMsgTypeRegister
)

//easyjson:json
type PingMsg struct {
	RegionId int64  `json:"region_id"` //region id
	ServerId int64  `json:"server_id"` //server id
	AppName  string `json:"app_name"`  //app name
	PingType int    `json:"ping_type"` //ping type
}

var pingMsg []byte

// SysInit system init
func SysInit(regionId, serverId int64, serverName string) (err error) {
	pingMsg, err = easyjson.Marshal(PingMsg{
		regionId,
		serverId,
		serverName,
		PingMsgTypePing,
	})
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	// send register message to admin
	err = Register(regionId, serverId, serverName)
	if err != nil {
		return
	}

	return
}

func Register(regionId, serverId int64, serverName string) (err error) {
	var registerMsg []byte
	registerMsg, err = easyjson.Marshal(PingMsg{
		regionId,
		serverId,
		serverName,
		PingMsgTypeRegister,
	})
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	return RabbitMQPublish(RabbitMQBusinessNameSvrPing, RabbitMQExchangePing, registerMsg)
}

// Ping send ping to admin
func Ping() (err error) {
	return RabbitMQPublish(RabbitMQBusinessNameSvrPing, RabbitMQExchangePing, pingMsg)
}
