package dao

import (
	"common"
	"fmt"
	"time"
	otcmodels "utils/otc/models"
)

type CommissionDistributeDao struct {
	common.BaseDao
}

func NewCommissionDistributeDao(db string) *CommissionDistributeDao {
	return &CommissionDistributeDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var CommissionDistributeDaoEntity *CommissionDistributeDao

func (d *CommissionDistributeDao) Add(start, end string) (id uint64, err error) {

	id, err = common.IdManagerGen(IdTypeCommissionDistribute)
	if err != nil {
		common.LogFuncError("IdManagerGen failed : %v", err)
		return
	}

	calc := otcmodels.CommissionDistribute{
		Id:              id,
		Start:           start,
		End:             end,
		Status:          CommissionStatusDoing,
		DistributeStart: time.Now().String(),
	}
	_, err = d.Orm.Insert(&calc)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *CommissionDistributeDao) Update(id uint64, status int, isEnd bool, desc string) (err error) {
	sql := fmt.Sprintf("UPDATE %s SET `%s`=?, `%s`=?", otcmodels.TABLE_CommissionDistribute,
		otcmodels.COLUMN_CommissionDistribute_Status,
		otcmodels.COLUMN_CommissionDistribute_Desc)
	params := []interface{}{status, desc}
	if isEnd {
		sql += fmt.Sprintf(", `%s`=? ", otcmodels.COLUMN_CommissionDistribute_DistributeEnd)
		params = append(params, time.Now().String())
	}

	sql += fmt.Sprintf(" WHERE `%s`=?", otcmodels.COLUMN_CommissionDistribute_Id)
	params = append(params, id)

	_, err = d.Orm.Raw(sql, params...).Exec()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
