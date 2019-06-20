package dao

import (
	"common"
	"fmt"
	"utils/report/models"
)

type GameUserMonthReportDao struct {
	common.BaseDao
}

func NewGameUserMonthReportDao(db string) *GameUserMonthReportDao {
	return &GameUserMonthReportDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var GameUserMonthReportDaoEntity *GameUserMonthReportDao

func (d *GameUserMonthReportDao) FindByDataRange(start, end int64) (reports []*models.GameUserMonthReport, err error) {
	reports = make([]*models.GameUserMonthReport, 0)
	querySql := fmt.Sprintf("select * from %s where %s>=? and %s<=?", models.TABLE_GameUserMonthReport, models.COLUMN_GameUserMonthReport_Ctime, models.COLUMN_GameUserMonthReport_Ctime)
	_, err = d.Orm.Raw(querySql, start, end).QueryRows(&reports)
	if err != nil {
		common.LogFuncError("db err %v", err)
	}
	return
}
func (d *GameUserMonthReportDao) Insert(reports []*models.GameUserMonthReport) (err error) {
	_, err = d.Orm.InsertMulti(100, reports)
	if err != nil {
		common.LogFuncError("db err %v", err)
	}
	return
}
func (d *GameUserMonthReportDao) Count() (rowNum int64, err error) {
	rowNum, err = d.Orm.QueryTable(models.TABLE_GameUserMonthReport).Count()
	if err != nil {
		common.LogFuncError("db err %v ", err)
	}
	return
}

func (d *GameUserMonthReportDao) DeleteAllData() (err error) {
	sql := fmt.Sprintf("delete from %s", models.TABLE_GameUserMonthReport)
	_, err = d.Orm.Raw(sql).Exec()
	if err != nil {
		common.LogFuncError("db err %v ", err)
	}
	return
}
