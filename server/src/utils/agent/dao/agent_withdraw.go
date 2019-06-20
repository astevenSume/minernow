package dao

import (
	"common"
	"fmt"
	"utils/agent/models"
)

type AgentWithdrawDao struct {
	common.BaseDao
}

func NewAgentWithdrawDao(db string) *AgentWithdrawDao {
	return &AgentWithdrawDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AgentWithdrawDaoEntity *AgentWithdrawDao

const (
	AgentWithdrawStatusUnkown uint8 = iota
	AgentWithdrawStatusInit
	AgentWithdrawStatusDecreased
	AgentWithdrawStatusDone //提现已完成 需要记录日报,接口ReportAgentDailyDaoEntity.AddWithDraw
	AgentWithdrawStatusFailed
	AgentWithdrawStatusChecking //提现额超过阈值进入审核状态
)

func (d *AgentWithdrawDao) Add(uid uint64, amount, now int64) (r models.AgentWithdraw, err error) {
	var id uint64
	id, err = common.IdManagerGen(IdTypeAgentWithdraw)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	r = models.AgentWithdraw{
		Id:     id,
		Uid:    uid,
		Amount: amount,
		Ctime:  now,
		Mtime:  now,
		Status: AgentWithdrawStatusInit,
	}

	_, err = d.Orm.Insert(&r)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *AgentWithdrawDao) UpdateStatus(id uint64, status uint8, desc string) (err error) {
	r := models.AgentWithdraw{
		Id:     id,
		Status: status,
		Desc:   desc,
	}

	_, err = d.Orm.Update(&r, models.COLUMN_AgentWithdraw_Status, models.COLUMN_AgentWithdraw_Desc)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *AgentWithdrawDao) ExitByUidAndStatus(uid uint64, status uint8) (exit bool) {
	qs := d.Orm.QueryTable(models.TABLE_AgentWithdraw)
	exit = qs.Filter(models.COLUMN_AgentWithdraw_Uid, uid).Filter(models.COLUMN_AgentWithdraw_Status, status).Exist()
	return
}
func (d *AgentWithdrawDao) Get(id uint64) (agentWithdraw *models.AgentWithdraw, err error) {
	agentWithdraw = new(models.AgentWithdraw)
	querySql := fmt.Sprintf("select * from %s where %s=?", models.TABLE_AgentWithdraw, models.COLUMN_AgentWithdraw_Id)
	err = d.Orm.Raw(querySql, id).QueryRow(&agentWithdraw)
	if err != nil {
		common.LogFuncError("AgentWithdrawDao Get err %v", err)
		return
	}
	return
}
func (d *AgentWithdrawDao) FindByStatus(status uint8, page, limit int) (r []*models.AgentWithdraw, err error) {
	r = make([]*models.AgentWithdraw, 0)
	querySql := fmt.Sprintf("select * from %s where %s=?  LIMIT ? OFFSET ?", models.TABLE_AgentWithdraw, models.COLUMN_AgentWithdraw_Status)
	_, err = d.Orm.Raw(querySql, status, limit, limit*(page-1)).QueryRows(&r)
	if err != nil {
		common.LogFuncError("AgentWithdrawDao FindByStatus err %v", err)
		return
	}
	return
}

type AgentWithdrawMsg struct {
	Id     uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Amount int64  `orm:"column(amount)" json:"amount,omitempty"`
	Status uint8  `orm:"column(status)" json:"status,omitempty"`
	Ctime  int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Mtime  int64  `orm:"column(mtime)" json:"utime,omitempty"`
}

func (d *AgentWithdrawDao) Query(uid uint64, page, limit int) (total int, records []AgentWithdrawMsg, err error) {
	// query total
	total, err = d.queryTotal(uid)
	if err != nil {
		return
	}

	_, err = d.Orm.Raw(fmt.Sprintf("SELECT `%s`, `%s`, `%s`, `%s`,`%s` FROM `%s` WHERE `%s`=? LIMIT ? OFFSET ?",
		models.COLUMN_AgentWithdraw_Id,
		models.COLUMN_AgentWithdraw_Amount,
		models.COLUMN_AgentWithdraw_Status,
		models.COLUMN_AgentWithdraw_Ctime,
		models.COLUMN_AgentWithdraw_Mtime,
		models.TABLE_AgentWithdraw,
		models.COLUMN_AgentWithdraw_Uid), uid, limit, limit*(page-1)).QueryRows(&records)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

//
func (d *AgentWithdrawDao) queryTotal(uid uint64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(%s) FROM %s WHERE %s=?",
		models.COLUMN_AgentWithdraw_Id, models.TABLE_AgentWithdraw, models.COLUMN_AgentWithdraw_Uid), uid).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *AgentWithdrawDao) QueryTotalWithStatus(status uint8) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(%s) FROM %s WHERE %s=?",
		models.COLUMN_AgentWithdraw_Id, models.TABLE_AgentWithdraw, models.COLUMN_AgentWithdraw_Status), status).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
