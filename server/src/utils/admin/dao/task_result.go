package dao

import (
	"common"
	adminmodels "utils/admin/models"
)

type TaskResultDao struct {
	common.BaseDao
}

func NewTaskResultDao(db string) *TaskResultDao {
	return &TaskResultDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var TaskResultDaoEntity *TaskResultDao

func (d *TaskResultDao) Add(r adminmodels.TaskResult) (err error) {
	_, err = d.Orm.Insert(&r)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *TaskResultDao) Query(appName string, regionId, serverId int64, page, perPage int) (meta PageInfo,
	list []adminmodels.TaskResult, err error) {
	qs := d.Orm.QueryTable(adminmodels.TABLE_TaskResult)
	if len(appName) > 0 {
		qs = qs.Filter(adminmodels.COLUMN_TaskResult_AppName, appName)
	}

	if regionId < RegionIdImpossible && regionId >= 0 {
		qs = qs.Filter(adminmodels.COLUMN_TaskResult_RegionId, regionId)
	}

	if serverId < ServerIdImpossible && serverId >= 0 {
		qs = qs.Filter(adminmodels.COLUMN_TaskResult_ServerId, serverId)
	}

	meta.Page = page
	meta.Limit = perPage

	var total int64
	total, err = qs.Count()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	meta.Total = int(total)

	if total <= 0 {
		return
	}

	_, err = qs.Limit(perPage).Offset((page - 1) * perPage).All(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
