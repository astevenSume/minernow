package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/bsm/redis-lock"
	"utils/agent/models"
	utilscommon "utils/common"
)

type InviteCodeDao struct {
	common.BaseDao
	region uint32
}

func NewInviteCodeDao(db string) (d *InviteCodeDao) {
	d = &InviteCodeDao{
		BaseDao: common.NewBaseDao(db),
	}

	tmp, err := beego.AppConfig.Int("InviteCodeRegion")
	if err != nil {
		panic("configuration InviteCodeRegion no found.")
	}

	d.region = uint32(tmp)

	return
}

var InviteCodeDaoEntity *InviteCodeDao

const (
	INVITE_CODE_STATUS_UNUSED uint32 = iota
	INVITE_CODE_STATUS_FROZEN
	INVITE_CODE_STATUS_USED
)

// try to get one unused invite code
func (d *InviteCodeDao) TryGetUnusedOne() (id uint32, inviteCode string, err error) {
	var l *lock.Locker
	l, err = common.RedisLock(utilscommon.KeyRedisInviteCode)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	defer common.RedisUnlock(l)

	// query a free one
	sql := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s=?",
		models.COLUMN_InviteCode_Id, models.COLUMN_InviteCode_Code, models.TABLE_InviteCode,
		models.COLUMN_InviteCode_Status)
	err = d.Orm.Raw(sql, INVITE_CODE_STATUS_UNUSED).QueryRow(&id, &inviteCode)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	//	frozen it
	err = d.updateStatusUnsafe(id, INVITE_CODE_STATUS_FROZEN)
	if err != nil {
		return
	}

	return
}

func (d *InviteCodeDao) updateStatusUnsafe(id, status uint32) (err error) {
	sql := fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=?", models.TABLE_InviteCode,
		models.COLUMN_InviteCode_Status, models.COLUMN_InviteCode_Id)
	_, err = d.Orm.Raw(sql, status, id).Exec()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// update invite code status
func (d *InviteCodeDao) UpdateStatus(id, status uint32) (err error) {
	var l *lock.Locker
	l, err = common.RedisLock(utilscommon.KeyRedisInviteCode)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	defer common.RedisUnlock(l)

	return d.updateStatusUnsafe(id, status)
}

// release invite code
func (d *InviteCodeDao) Release(id uint32) (err error) {
	var l *lock.Locker
	l, err = common.RedisLock(utilscommon.KeyRedisInviteCode)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	defer common.RedisUnlock(l)

	return d.updateStatusUnsafe(id, INVITE_CODE_STATUS_UNUSED)
}
