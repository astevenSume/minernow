package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"utils/admin/models"
)

type TaskDao struct {
	common.BaseDao
}

func NewTaskDao(db string) *TaskDao {
	return &TaskDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var TaskDaoEntity *TaskDao

// see generate id format
const (
	RegionIdImpossible int64 = 16   //region id won't larger than 15
	ServerIdImpossible int64 = 1024 //server id won't larger than 1023
)

func (d *TaskDao) Insert(t *models.Task) (err error) {
	now := common.NowUint32()
	t.Ctime = now
	t.Utime = now

	_, err = d.Orm.Insert(t)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	err = d.Orm.Read(t, models.COLUMN_Task_AppName, models.COLUMN_Task_Name)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *TaskDao) Update(t *models.Task) (err error) {
	now := common.NowUint32()
	t.Utime = now
	var n int64
	n, err = d.Orm.Update(t,
		models.COLUMN_Task_Name,
		models.COLUMN_Task_Alia,
		models.COLUMN_Task_AppName,
		models.COLUMN_Task_FuncName,
		models.COLUMN_Task_Status,
		models.COLUMN_Task_Spec,
		models.COLUMN_Task_Utime,
		models.COLUMN_Task_Desc)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	if n != 1 {
		err = common.ErrNoRowEffected
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *TaskDao) Query(appName string, page, perPage int) (meta PageInfo, list []models.Task, err error) {
	qs := d.Orm.QueryTable(models.TABLE_Task)
	if len(appName) > 0 {
		qs = qs.Filter(models.COLUMN_Task_AppName, appName)
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

func (d *TaskDao) QueryAll(appName string) (list []models.Task, err error) {
	qs := d.Orm.QueryTable(models.TABLE_Task)
	if appName != "" {
		qs = qs.Filter(models.COLUMN_Task_AppName, appName)
	}
	_, err = qs.All(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *TaskDao) QueryById(id uint32) (task models.Task, err error) {
	task.Id = id
	err = d.Orm.Read(&task)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

func (d *TaskDao) QueryByIdList(idList []uint32) (tasks []models.Task, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_Task).Filter(fmt.Sprintf("%s__in", models.COLUMN_Task_Id), idList).All(&tasks)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *TaskDao) QueryByAppNameAndName(appName, name string) (task models.Task, err error) {
	task.AppName = appName
	task.Name = name
	err = d.Orm.Read(&task, models.COLUMN_Task_AppName, models.COLUMN_Task_Name)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

func (d *TaskDao) QueryByAppNameAndNameList(appName string, nameList []string) (taskList []models.Task, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_Task).Filter(models.COLUMN_Task_AppName, appName).
		Filter(fmt.Sprintf("%s__in", models.COLUMN_Task_Name), nameList).All(&taskList)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

const (
	TaskStatusDisable uint8 = iota
	TaskStatusEnable
)

func (d *TaskDao) UpdateStatus(id uint32, status uint8) (err error) {
	t := models.Task{
		Id:     id,
		Status: status,
	}
	var n int64
	n, err = d.Orm.Update(&t, models.COLUMN_Task_Status)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	if n != 1 {
		err = common.ErrNoRowEffected
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *TaskDao) UpdateStatusMulti(idList []uint32, status uint8) (err error) {

	ids := ""
	for i, id := range idList {
		if i == 0 {
			ids += strconv.Itoa(int(id))
		} else {
			ids += "," + strconv.Itoa(int(id))
		}
	}

	_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=? WHERE %s IN (%s)",
		models.TABLE_Task, models.COLUMN_Task_Status, models.COLUMN_Task_Id, ids), status).
		Exec()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *TaskDao) Remove(id uint32) (err error) {
	t := models.Task{
		Id: id,
	}
	var n int64
	n, err = d.Orm.Delete(&t)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	if n != 1 {
		err = common.ErrNoRowEffected
		common.LogFuncError("%v", err)
		return
	}

	return
}
