package dao

import (
	"common"
	"fmt"
	"time"
	"utils/admin/models"
)

type ProfitThresholdDao struct {
	common.BaseDao
}

func NewProfitThresholdDao(db string) *ProfitThresholdDao {
	return &ProfitThresholdDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ProfitThresholdDaoEntity *ProfitThresholdDao

func (d *ProfitThresholdDao) Add(threshold int64, adminId uint64) (thresholdData *models.ProfitThreshold, err error) {
	now := time.Now().Unix()
	thresholdData = new(models.ProfitThreshold)
	thresholdData.Threshold = threshold
	thresholdData.AdminId = adminId
	thresholdData.Ctime = now
	thresholdData.Utime = now
	_, err = d.Orm.Insert(thresholdData)
	if err != nil {
		common.LogFuncError("add ProfitThreshold error %v", err)
		return
	}
	return
}

func (d *ProfitThresholdDao) DeleteAll() (err error) {
	querySql := fmt.Sprintf("delete from %s", models.TABLE_ProfitThreshold)
	_, err = d.Orm.Raw(querySql).Exec()
	if err != nil {
		common.LogFuncError("DeleteAll ProfitThreshold error %v", err)
		return
	}
	return
}

func (d *ProfitThresholdDao) Delete(id uint32) (err error) {
	querySql := fmt.Sprintf("delete from %s where %s=?", models.TABLE_ProfitThreshold, models.COLUMN_ProfitThreshold_Id)
	_, err = d.Orm.Raw(querySql, id).Exec()
	if err != nil {
		common.LogFuncError("Delete ProfitThreshold error %v", err)
		return
	}
	return
}

func (d *ProfitThresholdDao) Update(id uint32, threshold int64, adminId uint64) (thresholdData *models.ProfitThreshold, err error) {
	now := time.Now().Unix()
	thresholdData = new(models.ProfitThreshold)
	thresholdData.Id = id
	thresholdData.Threshold = threshold
	thresholdData.AdminId = adminId
	thresholdData.Ctime = now
	thresholdData.Utime = now
	_, err = d.Orm.Update(thresholdData)
	if err != nil {
		common.LogFuncError("Update ProfitThreshold error %v", err)
		return
	}
	return
}

func (d *ProfitThresholdDao) Find() (thresholdData []*models.ProfitThreshold, err error) {
	thresholdData = make([]*models.ProfitThreshold, 0)
	querySql := fmt.Sprintf("select * from %s", models.TABLE_ProfitThreshold)
	_, err = d.Orm.Raw(querySql).QueryRows(&thresholdData)
	if err != nil {
		common.LogFuncError("Find ProfitThreshold error %v", err)
		return
	}
	return
}
