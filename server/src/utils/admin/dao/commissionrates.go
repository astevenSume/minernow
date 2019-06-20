package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"sort"
	"utils/admin/models"
)

type CommissionRateDao struct {
	common.BaseDao
}

func NewCommissionRateDao(db string) *CommissionRateDao {
	return &CommissionRateDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var CommissionRateDaoEntity *CommissionRateDao

const CommissPrecision = 10000 //与eosplus.EosPrecision保持一致

//查询返佣等级表
func (d *CommissionRateDao) QueryCommRate(commissionrates *models.Commissionrates, cols ...string) error {
	if commissionrates == nil {
		return ErrParam
	}

	err := d.Orm.Read(commissionrates, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	commissionrates.Min = commissionrates.Min % uint64(commissionrates.Precision)
	commissionrates.Max = commissionrates.Max % uint64(commissionrates.Precision)

	return nil
}

//分页查询返佣等级表
func (d *CommissionRateDao) QueryPageCommRate() (data []models.Commissionrates, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_Commissionrates).All(&data)
	if err != nil {
		if err != orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	for i := 0; i < len(data); i++ {
		data[i].Max = data[i].Max / uint64(data[i].Precision)
		data[i].Min = data[i].Min / uint64(data[i].Precision)
		//common.LogFuncDebug("%v,%v,%v",d.Min , d.Max,d.Precision)
	}

	return
}

//添加返佣等级配置
func (d *CommissionRateDao) EditCommRate(commissioncfg Commissioncfgs) error {
	lenCfg := len(commissioncfg)
	if lenCfg == 0 {
		common.LogFuncError("mysql error:%v", commissioncfg)
		return ErrParam
	}

	mysql := fmt.Sprintf("delete from %s ", models.TABLE_Commissionrates)
	_, err := d.Orm.Raw(mysql).Exec()
	//_, err := d.Orm.QueryTable(models.TABLE_Commissionrates).Delete()
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	var commissionrates []models.Commissionrates
	sort.Sort(commissioncfg)
	for i := 0; i < lenCfg-1; i++ {
		if (commissioncfg[i].Min < commissioncfg[i].Max) && (commissioncfg[i].Max <= commissioncfg[i+1].Min) {
			cfg := models.Commissionrates{
				Id:         int64(i + 1),
				Min:        commissioncfg[i].Min * CommissPrecision,
				Max:        commissioncfg[i].Max * CommissPrecision,
				Commission: commissioncfg[i].Commission,
				Precision:  CommissPrecision,
				Ctime:      common.NowInt64MS(),
				Utime:      common.NowInt64MS(),
			}
			commissionrates = append(commissionrates, cfg)
		} else {
			common.LogFuncError("mysql error:%v", commissioncfg[i])
			return ErrParam
		}
	}
	cfg := models.Commissionrates{
		Id:         int64(lenCfg),
		Min:        commissioncfg[lenCfg-1].Min * CommissPrecision,
		Max:        commissioncfg[lenCfg-1].Max * CommissPrecision,
		Commission: commissioncfg[lenCfg-1].Commission,
		Precision:  CommissPrecision,
		Ctime:      common.NowInt64MS(),
		Utime:      common.NowInt64MS(),
	}
	commissionrates = append(commissionrates, cfg)

	_, err = d.Orm.InsertMulti(InsertMulCount, commissionrates)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return nil
}

// get all commission rates.
func (d *CommissionRateDao) All() (list []models.Commissionrates, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_Commissionrates).OrderBy(models.COLUMN_Commissionrates_Min).All(&list,
		models.COLUMN_Commissionrates_Min, models.COLUMN_Commissionrates_Max, models.COLUMN_Commissionrates_Commission, models.ATTRIBUTE_Commissionrates_Precision)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

/*//查询满足条件配置
func (d *CommissionRateDao) QueryConditionCommRate(min, max uint64) (commissionrates models.Commissionrates, err error) {
	mysql := fmt.Sprintf("SELECT * from %s WHERE (%s<=? and %s >?) or (%s<? and %s >=?)", models.TABLE_Commissionrates, models.COLUMN_Commissionrates_Min, models.COLUMN_Commissionrates_Max, models.COLUMN_Commissionrates_Min, models.COLUMN_Commissionrates_Max)
	err = d.Orm.Raw(mysql, min, min, max, max).QueryRow(&commissionrates)
	if err != nil {
		common.LogFuncDebug("mysql error:%v", err)
	}
	return
}

//添加返佣等级配置
func (d *CommissionRateDao) AddCommRate(commissionrates *models.Commissionrates) error {
	if commissionrates == nil || commissionrates.Min >= commissionrates.Max {
		return ErrParam
	}
	if commissionrates.Precision == 0 {
		commissionrates.Precision = CommissPrecision
	}

	commissionrates.Min = commissionrates.Min * commissionrates.Precision
	commissionrates.Max = commissionrates.Max * commissionrates.Precision
	_, err := d.QueryConditionCommRate(commissionrates.Min, commissionrates.Max)
	if err == nil {
		return ErrParam
	}

	id, err := d.Orm.Insert(commissionrates)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}
	commissionrates.Id = id

	return nil
}

//更新返佣等级配置
func (d *CommissionRateDao) UpdateCommRate(newCom *models.Commissionrates, cols ...string) error {
	if newCom == nil || newCom.Min >= newCom.Max {
		return ErrParam
	}
	if newCom.Precision == 0 {
		newCom.Precision = CommissPrecision
	}

	commissionrates := &models.Commissionrates{Id: newCom.Id}
	err := d.QueryCommRate(commissionrates, models.COLUMN_Commissionrates_Id)
	if err != nil {
		return err
	}

	commissionrates = newCom
	commissionrates.Min = commissionrates.Min * commissionrates.Precision
	commissionrates.Max = commissionrates.Max * commissionrates.Precision
	_, err = d.QueryConditionCommRate(commissionrates.Min, commissionrates.Max)
	if err == nil {
		return ErrParam
	}

	_, err = d.Orm.Update(commissionrates, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return nil
}

//删除返佣等级配置
func (d *CommissionRateDao) DelCommRate(id int64) error {
	commissionrates := &models.Commissionrates{Id: id}
	_, err := d.Orm.Delete(commissionrates)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return nil
}*/

//根据业绩金额获取业绩佣金
func (d *CommissionRateDao) GetPrecision(cur uint64) (precision, commissionrate uint64, err error) {
	mysql := fmt.Sprintf("SELECT %v,%v from %v WHERE %v<? and %v>=?", models.COLUMN_Commissionrates_Precision, models.COLUMN_Commissionrates_Commission, models.TABLE_Commissionrates, models.COLUMN_Commissionrates_Min, models.COLUMN_Commissionrates_Max)
	err = d.Orm.Raw(mysql, cur, cur).QueryRow(&precision, &commissionrate)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
	}
	return
}
