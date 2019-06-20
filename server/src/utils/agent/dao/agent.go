package dao

import (
	"common"
	"errors"
	"fmt"
	"github.com/bsm/redis-lock"
	"time"

	//"fmt"
	"github.com/astaxie/beego/orm"
	"utils/agent/models"
)

// todo channel data should be maintained by admin system.
const (
	CHANNEL_ID_UNKOWN uint32 = iota
	CHANNEL_ID_FIRST
)

const WithdrawSalt = "ZyGYFWIWO1BWYl9lpBKaNtKmXxFRrHwu5PgJD9V332AEWweZY1QdrRyTWithDraw"

type AgentDao struct {
	common.BaseDao
}

func NewAgentDao(db string) *AgentDao {
	return &AgentDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AgentDaoEntity *AgentDao

type AgentSummary struct {
	SumSalary      int64  `orm:"column(sum_salary)" json:"sum_salary,omitempty"`
	SumCanWithdraw int64  `orm:"column(sum_can_withdraw)" json:"sum_can_withdraw,omitempty"`
	InviteCode     string `orm:"column(invite_code);size(100)" json:"invite_code,omitempty"`
	InviteNum      uint32 `orm:"column(invite_num)" json:"invite_num,omitempty"`
}

// query all of agent's commission
func (d *AgentDao) QuerySummary(uid uint64) (agentSummary AgentSummary, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT `%s`, `%s`, `%s`, `%s` "+
		"FROM (SELECT `%s`, `%s`, `%s` "+
		"FROM `%s` WHERE `%s`=?) AS T1 "+
		"LEFT JOIN `%s` AS T2 ON T1.`%s` = T2.`%s`",
		models.COLUMN_Agent_SumSalary,
		models.COLUMN_Agent_SumCanWithdraw,
		models.COLUMN_AgentPath_InviteNum, models.COLUMN_AgentPath_InviteCode,
		models.COLUMN_AgentPath_Uid,
		models.COLUMN_AgentPath_InviteNum, models.COLUMN_AgentPath_InviteCode,
		models.TABLE_AgentPath, models.COLUMN_AgentPath_Uid,
		models.TABLE_Agent, models.COLUMN_AgentPath_Uid, models.COLUMN_Agent_Uid),
		uid).QueryRow(&agentSummary)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
		return
	}

	return
}

// query agent's commission of specific channel
func (d *AgentDao) QuerySpecificChannel(uid uint64, channelId uint32) (agent models.AgentChannelCommission, err error) {
	agent.Uid = uid
	agent.ChannelId = channelId
	err = d.Orm.Read(&agent, models.COLUMN_AgentChannelCommission_Uid, models.COLUMN_AgentChannelCommission_ChannelId)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		} else {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

var (
	ErrCurrencyNoEnough = errors.New("decrease balance no enough")
)

func (d *AgentDao) Create(uid uint64) (err error) {
	pwd, _ := d.GeneratePwdMD5(uid, 0)
	agent := models.Agent{
		Uid:            uid,
		SumSalary:      0,
		SumWithdraw:    0,
		SumCanWithdraw: 0,
		Ctime:          time.Now().Unix(),
		Mtime:          time.Now().Unix(),
		Pwd:            pwd,
	}
	_, err = d.Orm.Insert(&agent)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (d *AgentDao) Info(uid uint64) (agent models.Agent, err error) {
	agent = models.Agent{
		Uid: uid,
	}
	err = d.Orm.Read(&agent, models.COLUMN_Agent_Uid)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (d *AgentDao) CheckPwdMD5(uid uint64, sumCanWithdraw int64, curPwd string) (ok bool, err error) {
	if sumCanWithdraw > 0 {
		var pwd string
		pwd, err = common.GenerateDoubleMD5(fmt.Sprintf("%v%v", uid, sumCanWithdraw), WithdrawSalt)
		if err != nil {
			common.LogFuncError("md5 error:[%v]", uid)
			return
		}

		if pwd != curPwd {
			ok = false
			return
		}
	}

	ok = true
	return
}

func (d *AgentDao) GeneratePwdMD5(uid uint64, sumCanWithdraw int64) (pwd string, err error) {
	pwd, err = common.GenerateDoubleMD5(fmt.Sprintf("%v%v", uid, sumCanWithdraw), WithdrawSalt)
	if err != nil {
		common.LogFuncError("md5 error:[%v]", uid)
		return
	}
	return
}

// add commission
func (d *AgentDao) AddBalance(uid uint64, amount int64, opt lock.Options) (before, after int64, err error) {
	var l *lock.Locker
	l, err = agentRedisJamLock(uid, AgentCanDraw, opt)
	if err != nil {
		common.LogFuncError("err:%v\n", err)
		return
	}
	defer func() {
		_ = agentRedisUnJamLock(l)
	}()

	agent := models.Agent{
		Uid: uid,
	}

	err = d.Orm.Read(&agent)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var ok bool
	ok, err = d.CheckPwdMD5(uid, agent.SumCanWithdraw, agent.Pwd)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	if !ok {
		//警告：金额被修改
		err = errors.New(fmt.Sprintf("warning uid[%v] sumcanwithdraw modify", uid))
		return
	}

	before = agent.SumCanWithdraw
	agent.SumCanWithdraw = agent.SumCanWithdraw + amount
	after = agent.SumCanWithdraw
	agent.SumSalary = agent.SumSalary + amount
	pwd, err := d.GeneratePwdMD5(agent.Uid, agent.SumCanWithdraw)
	if err != nil {
		common.LogFuncError("md5 error:[%v]", agent.Uid)
		return
	}
	agent.Pwd = pwd
	agent.Mtime = time.Now().Unix()

	_, err = d.Orm.Update(&agent, models.COLUMN_Agent_SumCanWithdraw, models.COLUMN_Agent_SumSalary,
		models.COLUMN_Agent_Pwd, models.COLUMN_Agent_Mtime)
	if err != nil {
		common.LogFuncError("md5 error:[%v]", agent.Uid)
		return
	}

	/*_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=%s+?,%s=%s+? WHERE %s=?", models.TABLE_Agent,
		models.COLUMN_Agent_SumSalary, models.COLUMN_Agent_SumSalary, models.COLUMN_Agent_SumCanWithdraw,
		models.COLUMN_Agent_SumCanWithdraw, models.COLUMN_Agent_Uid), amount, amount, uid).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}*/

	return
}

func (d *AgentDao) DecreaseBalance(uid uint64, amount int64, opt lock.Options) (err error) {
	var l *lock.Locker
	l, err = agentRedisJamLock(uid, AgentCanDraw, opt)
	if err != nil {
		common.LogFuncError("err:%v\n", err)
		return
	}
	defer func() {
		_ = agentRedisUnJamLock(l)
	}()

	agent := models.Agent{
		Uid: uid,
	}

	err = d.Orm.Read(&agent)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var ok bool
	ok, err = d.CheckPwdMD5(uid, agent.SumCanWithdraw, agent.Pwd)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	if !ok {
		//警告：金额被修改
		err = errors.New(fmt.Sprintf("warning uid[%v] sumcanwithdraw modify", uid))
		return
	}

	agent.SumCanWithdraw = agent.SumCanWithdraw - amount
	agent.SumSalary = agent.SumSalary - amount
	pwd, err := d.GeneratePwdMD5(agent.Uid, agent.SumCanWithdraw)
	if err != nil {
		common.LogFuncError("md5 error:[%v]", agent.Uid)
		return
	}
	agent.Pwd = pwd
	agent.Mtime = time.Now().Unix()

	_, err = d.Orm.Update(&agent, models.COLUMN_Agent_SumCanWithdraw, models.COLUMN_Agent_SumSalary,
		models.COLUMN_Agent_Pwd, models.COLUMN_Agent_Mtime)
	if err != nil {
		common.LogFuncError("md5 error:[%v]", agent.Uid)
		return
	}

	/*_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=%s-?,%s=%s-? WHERE %s=?", models.TABLE_Agent,
		models.COLUMN_Agent_SumSalary, models.COLUMN_Agent_SumSalary, models.COLUMN_Agent_SumCanWithdraw,
		models.COLUMN_Agent_SumCanWithdraw, models.COLUMN_Agent_Uid), amount, amount, uid).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}*/

	return
}

func (d *AgentDao) DecreaseCommission(uid uint64, amount int64, opt lock.Options) (err error) {
	var l *lock.Locker
	l, err = agentRedisJamLock(uid, AgentCanDraw, opt)
	if err != nil {
		common.LogFuncError("err:%v\n", err)
		return
	}
	defer func() {
		_ = agentRedisUnJamLock(l)
	}()

	agent := models.Agent{
		Uid: uid,
	}

	err = d.Orm.Read(&agent)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var ok bool
	ok, err = d.CheckPwdMD5(uid, agent.SumCanWithdraw, agent.Pwd)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	if !ok {
		//警告：金额被修改
		err = errors.New(fmt.Sprintf("warning uid[%v] sumcanwithdraw modify", uid))
		return
	}

	if agent.SumCanWithdraw < amount {
		common.LogFuncWarning("amount[%v] bigger than balance [%v]", amount, agent.SumCanWithdraw)
		err = ErrCurrencyNoEnough
		return
	}

	agent.SumCanWithdraw = agent.SumCanWithdraw - amount
	agent.SumWithdraw = agent.SumWithdraw + amount
	pwd, err := d.GeneratePwdMD5(agent.Uid, agent.SumCanWithdraw)
	if err != nil {
		common.LogFuncError("md5 error:[%v]", agent.Uid)
		return
	}
	agent.Pwd = pwd
	agent.Mtime = time.Now().Unix()

	_, err = d.Orm.Update(&agent, models.COLUMN_Agent_SumCanWithdraw, models.COLUMN_Agent_SumWithdraw,
		models.COLUMN_Agent_Pwd, models.COLUMN_Agent_Mtime)
	if err != nil {
		common.LogFuncError("md5 error:[%v]", agent.Uid)
		return
	}

	/*_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=%s+?,%s=%s-? WHERE %s=?", models.TABLE_Agent,
		models.COLUMN_Agent_SumWithdraw, models.COLUMN_Agent_SumWithdraw, models.COLUMN_Agent_SumCanWithdraw,
		models.COLUMN_Agent_SumCanWithdraw, models.COLUMN_Agent_Uid), amount, amount, uid).Exec()
	if err != nil {
		common.LogFuncError("uid:%v,error:%v", uid, err)
		return
	}*/

	return
}

func (d *AgentDao) UpdateCanWithdraw(uid uint64, amount int64, opt lock.Options) (canWithdraw int64, err error) {
	var l *lock.Locker
	l, err = agentRedisJamLock(uid, AgentCanDraw, opt)
	if err != nil {
		common.LogFuncError("err:%v\n", err)
		return
	}
	defer func() {
		_ = agentRedisUnJamLock(l)
	}()

	agent := models.Agent{
		Uid: uid,
	}

	err = d.Orm.Read(&agent)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var ok bool
	ok, err = d.CheckPwdMD5(uid, agent.SumCanWithdraw, agent.Pwd)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	if !ok {
		//警告：金额被修改
		err = errors.New(fmt.Sprintf("warning uid[%v] sumcanwithdraw modify", uid))
		return
	}

	agent.SumCanWithdraw = agent.SumCanWithdraw - amount
	pwd, err := d.GeneratePwdMD5(agent.Uid, agent.SumCanWithdraw)
	if err != nil {
		common.LogFuncError("md5 error:[%v]", agent.Uid)
		return
	}
	agent.Pwd = pwd
	agent.Mtime = time.Now().Unix()

	_, err = d.Orm.Update(&agent, models.COLUMN_Agent_SumCanWithdraw, models.COLUMN_Agent_Pwd, models.COLUMN_Agent_Mtime)
	if err != nil {
		common.LogFuncError("md5 error:[%v]", agent.Uid)
		return
	}
	canWithdraw = agent.SumCanWithdraw

	/*_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=%s+? WHERE %s = ?", models.TABLE_Agent,
		models.COLUMN_Agent_SumCanWithdraw, models.COLUMN_Agent_SumCanWithdraw, models.COLUMN_Agent_Uid), salary, uid).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	canWithdraw = agent.SumCanWithdraw */

	return
}

func (d *AgentDao) Total() (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT COUNT(%s) FROM %s", models.COLUMN_Agent_Uid, models.TABLE_Agent)).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *AgentDao) QueryByPage(page, limit int) (list []models.Agent, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", models.TABLE_Agent),
		limit, (page-1)*limit).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
