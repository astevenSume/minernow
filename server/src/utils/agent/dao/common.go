package dao

import (
	"common"
	"fmt"
	"github.com/bsm/redis-lock"
)

const (
	IdTypeAgentWithdraw = iota //agent withdraw record
	IdTypeMax
)

const (
	QueryPage = 50

	//redis key
	AgentCanDraw = "agent_can_draw"
	AgentPath    = "agent_path"
)

func Init(entityInitFunc common.EntityInitFunc) (err error) {
	const dbOtc = "otc"
	AgentDaoEntity = NewAgentDao(dbOtc)
	AgentPathDaoEntity = NewAgentPathDao(dbOtc)
	AgentWithdrawDaoEntity = NewAgentWithdrawDao(dbOtc)
	AgentChannelCommissionDaoEntity = NewAgentChannelCommissionDao(dbOtc)
	InviteCodeDaoEntity = NewInviteCodeDao(dbOtc)

	if entityInitFunc != nil {
		err = entityInitFunc()
		if err != nil {
			return
		}
	}
	return
}

//堵塞锁
func agentRedisJamLock(uid uint64, pre string, opt lock.Options) (l *lock.Locker, err error) {
	l, err = common.RedisLock2(fmt.Sprintf("%s%d", pre, uid), opt)
	return
}

func agentRedisUnJamLock(l *lock.Locker) (err error) {
	if l == nil {
		return
	}
	err = common.RedisUnlock(l)
	l = nil
	return
}
