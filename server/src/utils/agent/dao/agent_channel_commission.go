package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	otccommon "otc/common"
	"utils/agent/models"
)

type AgentChannelCommissionDao struct {
	common.BaseDao
}

func NewAgentChannelCommissionDao(db string) *AgentChannelCommissionDao {
	return &AgentChannelCommissionDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AgentChannelCommissionDaoEntity *AgentChannelCommissionDao

// query user agent channel commission
func (d *AgentChannelCommissionDao) QueryByUid(uid uint64) (integer, decimals int32, err error) {

	var list []models.AgentChannelCommission
	_, err = d.Orm.QueryTable(models.TABLE_AgentChannelCommission).Filter(models.COLUMN_AgentChannelCommission_Uid, uid).All(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			return
		}
		common.LogFuncError("%v", err)
		return
	}

	for _, v := range list {
		integer, decimals, err = common.AddCurrency(integer, decimals, v.Integer, v.Decimals, otccommon.Cursvr.EusdPrecision)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

const (
	CommissionStatusUnsent uint8 = iota
	CommissionStatusSent
)

// add agent channel commission
func (d *AgentChannelCommissionDao) InsertOrUpdate(uid uint64, channelId uint32, integer, decimals int32, timestamp, now int64) (err error) {
	commission := models.AgentChannelCommission{
		Uid:       uid,
		ChannelId: channelId,
		Ctime:     timestamp,
		Integer:   integer,
		Decimals:  decimals,
		Mtime:     now,
	}

	_, err = d.Orm.InsertOrUpdate(&commission, "status=status") //don't change status
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *AgentChannelCommissionDao) QueryTotal(timestamp int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(*) FROM "+
		"(SELECT %s, sum(`%s`),sum(`%s`) FROM %s WHERE %s=? AND %s=? GROUP BY %s) as T",
		models.COLUMN_AgentChannelCommission_Uid,
		models.COLUMN_AgentChannelCommission_Integer,
		models.COLUMN_AgentChannelCommission_Decimals,
		d.BusinessTable(models.TABLE_AgentChannelCommission),
		models.COLUMN_AgentChannelCommission_Ctime,
		models.COLUMN_AgentChannelCommission_Status,
		models.COLUMN_AgentChannelCommission_Uid), timestamp, CommissionStatusUnsent).QueryRow(&total)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
		}
		return
	}

	return
}

type AgentCommission struct {
	Uid      uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Integer  int32  `orm:"column(integer)" json:"integer,omitempty"`
	Decimals int32  `orm:"column(decimals)" json:"decimals,omitempty"`
}

func (d *AgentChannelCommissionDao) Query(timestamp int64, page, limit int) (list []AgentCommission, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT %s, sum(`%s`) as `%s`,sum(`%s`) as `%s` FROM %s WHERE %s=? AND %s=? GROUP BY %s LIMIT ? OFFSET ?",
		models.COLUMN_AgentChannelCommission_Uid,
		models.COLUMN_AgentChannelCommission_Integer,
		models.COLUMN_AgentChannelCommission_Integer,
		models.COLUMN_AgentChannelCommission_Decimals,
		models.COLUMN_AgentChannelCommission_Decimals,
		d.BusinessTable(models.TABLE_AgentChannelCommission),
		models.COLUMN_AgentChannelCommission_Ctime,
		models.COLUMN_AgentChannelCommission_Status,
		models.COLUMN_AgentChannelCommission_Uid), timestamp, CommissionStatusUnsent, limit, (page-1)*limit).QueryRows(&list)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

func (d *AgentChannelCommissionDao) UpdateStatusToSent(uid uint64, timestamp int64) (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=? AND %s=? AND %s=?",
		d.BusinessTable(models.TABLE_AgentChannelCommission),
		models.COLUMN_AgentChannelCommission_Status,
		models.COLUMN_AgentChannelCommission_Uid,
		models.COLUMN_AgentChannelCommission_Ctime,
		models.COLUMN_AgentChannelCommission_Status),
		CommissionStatusSent, uid, timestamp, CommissionStatusUnsent).Exec()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
