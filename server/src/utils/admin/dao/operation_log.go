package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/admin/models"
)

type OperationLogDao struct {
	common.BaseDao
}

func NewOperationLogDao(db string) *OperationLogDao {
	return &OperationLogDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var OperationLogDaoEntity *OperationLogDao

//auto_models_start
type OperationLog struct {
	Id           string `orm:"column(id);pk" json:"id,omitempty"`
	AdminId      uint64 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	AdminNick    string `orm:"column(name)" json:"name,omitempty"`
	Method       string `orm:"column(method);size(100)" json:"method,omitempty"`
	Route        string `orm:"column(route);size(100)" json:"route,omitempty"`
	Action       int32  `orm:"column(action)" json:"action,omitempty"`
	Input        string `orm:"column(input);size(512)" json:"input,omitempty"`
	UserAgent    string `orm:"column(user_agent);size(512)" json:"user_agent,omitempty"`
	Ips          string `orm:"column(ips);size(100)" json:"ips,omitempty"`
	ResponseCode int32  `orm:"column(response_code)" json:"response_code,omitempty"`
	Ctime        int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

//分页查询日志
func (d *OperationLogDao) QueryPageOperationLog(id uint64, page int, perPage int, action int32) (total int, logs []OperationLog, err error) {
	var qbQuery orm.QueryBuilder
	var qbTotal orm.QueryBuilder
	qbTotal, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	// 构建查询对象
	qbTotal.Select("Count(*) AS total").
		From(models.TABLE_OperationLog).
		Where("1 = 1")
	qbQuery.Select("T1.*",
		"T2."+models.COLUMN_AdminUser_Name).
		From("((").Select("*").From(models.TABLE_OperationLog).Where("1=1")
	var param []interface{}

	if id > 0 {
		qbQuery.And(models.COLUMN_OperationLog_AdminId + "= ?")
		qbTotal.And(models.COLUMN_OperationLog_AdminId + "= ?")
		param = append(param, id)
	}
	if action > 0 {
		qbQuery.And(models.COLUMN_OperationLog_Action + "= ?")
		qbTotal.And(models.COLUMN_OperationLog_Action + "= ?")
		param = append(param, action)
	}

	sqlTotal := qbTotal.String()
	err = d.Orm.Raw(sqlTotal, param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error: %v", err)
		return
	}

	qbQuery.OrderBy("-" + models.COLUMN_OperationLog_Id)
	qbQuery.Limit(perPage).Offset((page - 1) * perPage)
	sqlQuery := qbQuery.String()
	sqlQuery = fmt.Sprintf("%s) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s)", sqlQuery, models.TABLE_AdminUser,
		models.COLUMN_OperationLog_AdminId, models.COLUMN_AdminUser_Id)
	_, err = d.Orm.Raw(sqlQuery, param...).QueryRows(&logs)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//创建日志
func (d *OperationLogDao) CreateOperationLog(operationLog *models.OperationLog) error {
	if operationLog == nil {
		return ErrParam
	}

	_, err := d.Orm.Insert(operationLog)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	return nil
}
