package common

import (
	"fmt"
	"time"
)

const (
	ServerAdmin    = 0
	ServerOtc      = 1
	ServerEUSD     = 2
	ServerUSDT     = 3
	ServerGame     = 4
	ServerOtcTrade = 5 // otc交易
)

const (
	ServerStateRunning    = iota //默认 正常
	ServerStateStopFirst         //准备停服
	ServerStateStopSecond        //已经停服
)

const RedisStopKey = "server_stop"
const RedisStopLogKey = "server_stop_log"

var ServerRunning = true
var ServerState = ServerStateRunning

func IsServerStop(server int) int {
	state, err := RedisManger.Get(fmt.Sprintf("%s_%d", RedisStopKey, server)).Int()
	if err != nil {
		return 0
	}
	return state
}

func StopServer(server int) (err error) {
	key := fmt.Sprintf("%s_%d", RedisStopKey, server)
	if server == ServerOtc {
		_, err = RedisManger.Set(key, ServerStateStopFirst, 0).Result()
		if err != nil {
			return
		}
		//30分钟后完全停止服务 针对 otc
		time.AfterFunc(30*time.Minute, func() {
			_, err = RedisManger.Set(key, ServerStateStopSecond, 0).Result()
			if err != nil {
				return
			}
		})
	} else {
		_, err = RedisManger.Set(key, ServerStateStopSecond, 0).Result()
		if err != nil {
			return
		}
	}
	return
}

func StartServer(server int) (err error) {
	key := fmt.Sprintf("%s_%d", RedisStopKey, server)

	_, err = RedisManger.Del(key).Result()
	return
}

//cron 定时同步 服务状态
func SyncServerState(server int) {
	res := IsServerStop(server)

	if res != ServerStateRunning {
		ServerRunning = false
	}
	ServerState = res
	return
}

//写日志
func ServerStopWriteLog(str string) {
	RedisManger.HSet(RedisStopLogKey, str, NowString())
	RedisManger.Expire(RedisStopLogKey, 3*time.Hour)
}
