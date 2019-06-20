package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	otcmodels "utils/otc/models"
)

//承兑商审核状态
const (
	CommissionStatStatusNil = iota
	CommissionStatStatusUnDistribute
	CommissionStatStatusDistributed
	CommissionStatStatusStatusMax
)

const (
	CommissionStatusUnsent uint8 = iota
	CommissionStatusSent
)

type CommissionStatDao struct {
	common.BaseDao
}

func NewCommissionStatDao(db string) *CommissionStatDao {
	return &CommissionStatDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var CommissionStatDaoEntity *CommissionStatDao

func (d *CommissionStatDao) Query(timestamp int64) (stat otcmodels.CommissionStat, err error) {
	err = d.Orm.Read(&stat)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
	}
	return
}

func (d *CommissionStatDao) Update(timestamp int64,
	channelInteger, channelDecimals,
	commissionInteger, commissionDecimals,
	profitInteger, profitDecimals int32) (err error) {
	stat := otcmodels.CommissionStat{
		Ctime:              timestamp,
		ChannelInteger:     channelInteger,
		ChannelDecimals:    channelDecimals,
		CommissionInteger:  commissionInteger,
		CommissionDecimals: commissionDecimals,
		ProfitInteger:      profitInteger,
		ProfitDecimals:     profitDecimals,
		Mtime:              common.NowInt64MS(),
	}

	_, err = d.Orm.InsertOrUpdate(&stat, "tax_integer=tax_integer", "tax_decimals=tax_decimals", "`status`=`status`")
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *CommissionStatDao) UpdateTax(timestamp int64, taxInteger, taxDecimals int32) (err error) {
	stat := otcmodels.CommissionStat{
		Ctime:       timestamp,
		TaxInteger:  taxInteger,
		TaxDecimals: taxDecimals,
		Mtime:       common.NowInt64MS(),
	}

	_, err = d.Orm.InsertOrUpdate(&stat,
		"channel_integer=channel_integer",
		"channel_decimals=channel_decimals",
		"commission_integer=commission_integer",
		"commission_decimals=commission_decimals",
		"profit_integer=profit_integer",
		"profit_decimals=profit_decimals",
		"`status`=`status`")
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// set commission distributed.
func (d *CommissionStatDao) SetDistributed(timestamp int64) (err error) {
	s := otcmodels.CommissionStat{
		Ctime:  timestamp,
		Status: CommissionStatusSent,
	}

	_, err = d.Orm.Update(&s, otcmodels.COLUMN_CommissionStat_Status)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *CommissionStatDao) IsDistributed(ctime int64) (bool, error) {
	commissionStat := otcmodels.CommissionStat{
		Ctime: ctime,
	}
	err := d.Orm.Read(&commissionStat, otcmodels.COLUMN_CommissionStat_Ctime)
	if err != nil {
		return false, err
	}

	if commissionStat.Status == CommissionStatStatusDistributed {
		return true, nil
	}

	return false, nil
}

//分页条件查询
func (d *CommissionStatDao) QueryByPage(status int8, page int, perPage int) (total int64, commissionStats []otcmodels.CommissionStat, err error) {
	qs := d.Orm.QueryTable(otcmodels.TABLE_CommissionStat)
	if status > CommissionStatStatusNil && status < CommissionStatStatusStatusMax {
		qs = qs.Filter(otcmodels.COLUMN_CommissionStat_Status, status)
	}
	qs = qs.OrderBy("-" + otcmodels.COLUMN_CommissionStat_Ctime)

	total, err = qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > total {
		err = nil
		return
	}
	_, err = qs.Limit(perPage, start).All(&commissionStats)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}
