package common

// 消息队列RPC消息命令码
const (
	RabbitMQRpcMsgCmdUnknown = iota

	// 定时任务 [1, 199]
	RabbitMQRpcMsgCmdTaskQuerySingle
)
