package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/otc/models"
)

type EosOtcReportDao struct {
	common.BaseDao
}

func NewEosOtcReportDao(db string) *EosOtcReportDao {
	return &EosOtcReportDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var EosOtcReportDaoEntity *EosOtcReportDao

func (d *EosOtcReportDao) Create(report []models.EosOtcReport) (err error) {
	if len(report) == 0 {
		return
	}
	_, err = d.Orm.InsertMulti(100, report)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	return
}

func (d *EosOtcReportDao) QueryByUIdAndDate(uid uint64, date int32) (has bool, reportRes *models.EosOtcReport, err error) {
	reportRes = &models.EosOtcReport{
		Uid:  uid,
		Date: date,
	}

	err = d.Orm.Read(reportRes)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DB_ERR:%v", err)
		return
	}
	has = true

	return
}

func (d *EosOtcReportDao) DeleteByDate(date int32) (delRow int64, err error) {
	//var report = new(models.EosOtcReport)
	//report.Date = date
	//delRow, err = d.Orm.Delete(report)
	sql := fmt.Sprintf("DELETE from %s where date=?", models.TABLE_EosOtcReport)
	_, err = d.Orm.Raw(sql, date).Exec()
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	return
}
