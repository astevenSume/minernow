package dao

import (
	"common"
	"fmt"
	"strings"
	"utils/report/models"
)

type ReportCommissionDao struct {
	common.BaseDao
}

func NewReportCommissionDao(db string) *ReportCommissionDao {
	return &ReportCommissionDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportCommissionDaoEntity *ReportCommissionDao

func (d *ReportCommissionDao) InsertMul(reportCommission []models.ReportCommission) (err error) {
	if len(reportCommission) == 0 {
		return
	}

	_, err = d.Orm.InsertMulti(InsertMulCount, reportCommission)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		_ = d.Del()
		return err
	}

	return
}

func (d *ReportCommissionDao) Del() (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("TRUNCATE TABLE %s", models.TABLE_ReportCommission)).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	return
}

type ReportCommission struct {
	Uid             uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	TeamWithdraw    int64  `orm:"column(team_withdraw)" json:"team_withdraw,omitempty"`
	TeamCanWithdraw int64  `orm:"column(team_can_withdraw)" json:"team_can_withdraw,omitempty"`
}

func (d *ReportCommissionDao) InfoByUIds(uIds []string) (reportCommission []ReportCommission, err error) {
	if len(uIds) == 0 {
		return
	}

	sqlQuery := fmt.Sprintf("SELECT %s,%s,%s FROM %s WHERE %s IN(%s)",
		models.COLUMN_ReportCommission_Uid, models.COLUMN_ReportCommission_TeamCanWithdraw,
		models.COLUMN_ReportCommission_TeamWithdraw, models.TABLE_ReportCommission,
		models.COLUMN_ReportCommission_Uid, strings.Join(uIds, ","))
	_, err = d.Orm.Raw(sqlQuery).QueryRows(&reportCommission)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *ReportCommissionDao) Sum() (sumCommission ReportCommission, err error) {
	sqlQuery := fmt.Sprintf("SELECT SUM(%s) AS team_withdraw, SUM(%s) AS team_can_withdraw FROM %s WHERE %s=1",
		models.COLUMN_ReportCommission_TeamCanWithdraw, models.COLUMN_ReportCommission_TeamWithdraw,
		models.TABLE_ReportCommission, models.COLUMN_ReportCommission_Level)
	err = d.Orm.Raw(sqlQuery).QueryRow(&sumCommission)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *ReportCommissionDao) Personal(uid uint64) (sumCommission ReportCommission, err error) {
	sqlQuery := fmt.Sprintf("SELECT SUM(%s) AS team_withdraw, SUM(%s) AS team_can_withdraw FROM %s WHERE %s=?",
		models.COLUMN_ReportCommission_TeamCanWithdraw, models.COLUMN_ReportCommission_TeamWithdraw,
		models.TABLE_ReportCommission, models.COLUMN_ReportCommission_Uid)
	err = d.Orm.Raw(sqlQuery, uid).QueryRow(&sumCommission)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}
