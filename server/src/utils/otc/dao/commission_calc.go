package dao

import (
	"common"
	"fmt"
	"time"
	otcmodels "utils/otc/models"
)

type CommissionCalcDao struct {
	common.BaseDao
}

func NewCommissionCalcDao(db string) *CommissionCalcDao {
	return &CommissionCalcDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var CommissionCalcDaoEntity *CommissionCalcDao

const (
	CommissionStatusUnkown = iota
	CommissionStatusDoing
	CommissionStatusDone
	CommissionStatusFailed
)

func (d *CommissionCalcDao) Add(start, end string) (id uint64, err error) {

	id, err = common.IdManagerGen(IdTypeCommissionCalc)
	if err != nil {
		common.LogFuncError("IdManagerGen failed : %v", err)
		return
	}

	calc := otcmodels.CommissionCalc{
		Id:        id,
		Start:     start,
		End:       end,
		Status:    CommissionStatusDoing,
		CalcStart: time.Now().String(),
	}
	_, err = d.Orm.Insert(&calc)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *CommissionCalcDao) Update(id uint64, status int, isEnd bool, desc string) (err error) {
	sql := fmt.Sprintf("UPDATE %s SET `%s`=?, `%s`=?", otcmodels.TABLE_CommissionCalc,
		otcmodels.COLUMN_CommissionCalc_Status,
		otcmodels.COLUMN_CommissionCalc_Desc)
	params := []interface{}{status, desc}
	if isEnd {
		sql += fmt.Sprintf(", `%s`=? ", otcmodels.COLUMN_CommissionCalc_CalcEnd)
		params = append(params, time.Now().String())
	}

	sql += fmt.Sprintf(" WHERE `%s`=?", otcmodels.COLUMN_CommissionCalc_Id)
	params = append(params, id)

	_, err = d.Orm.Raw(sql, params...).Exec()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
