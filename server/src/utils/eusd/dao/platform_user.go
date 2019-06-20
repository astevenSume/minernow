package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/eusd/models"
)

const (
	PlatformCateNone     = iota // 未分配
	PlatformCateToken           // TOKEN管理账户 发币
	PlatformCateResource        // 资源账户  有EOS，可以创建账号 & 购买RAM & 抵押资源
	PlatformCatePlatform        // 平台总池
	PlatformCateGame            // 游戏API账号
	PlatformCateSalesman        // 推销总账号
	PlatformCateDividend        // 分红
)

const (
	PlatformStatusLock   = 0
	PlatformStatusActive = 1
)

type PlatformUserDao struct {
	common.BaseDao
}

func NewPlatformUserDao(db string) *PlatformUserDao {
	return &PlatformUserDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *PlatformUserDao) All() (list []*models.PlatformUser, err error) {
	list = []*models.PlatformUser{}

	_, err = d.Orm.Raw("SELECT * from platform_user").QueryRows(&list)
	if err != nil {
		if err != nil {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

func (d *PlatformUserDao) Add(uid uint64, pid int32) (data *models.PlatformUser, err error) {
	data = &models.PlatformUser{
		Pid:    pid,
		Uid:    uid,
		Ctime:  common.NowUint32(),
		Status: PlatformStatusActive,
	}

	_, err = d.Orm.Insert(data)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

func (d *PlatformUserDao) FetchActive(pid int32) (list []*models.PlatformUser, err error) {
	list = []*models.PlatformUser{}

	_, err = d.Orm.Raw("SELECT * from platform_user Where pid=? and status=?", pid, PlatformStatusActive).QueryRows(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

func (d *PlatformUserDao) Lock(uid uint64) (err error) {
	data := &models.PlatformUser{
		Uid:    uid,
		Status: PlatformStatusLock,
	}

	_, err = d.Orm.Update(data, models.COLUMN_PlatformUser_Status)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

func (d *PlatformUserDao) Active(uid uint64) (err error) {
	data := &models.PlatformUser{
		Uid:    uid,
		Status: PlatformStatusActive,
	}

	_, err = d.Orm.Update(data, models.COLUMN_PlatformUser_Status)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}
