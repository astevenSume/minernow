package dao

import (
	"common"
	"fmt"
	"strings"
	"utils/report/models"
)

type ReportTeamDailyDao struct {
	common.BaseDao
}

func NewReportTeamDailyDao(db string) *ReportTeamDailyDao {
	return &ReportTeamDailyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportTeamDailyDaoEntity *ReportTeamDailyDao

func (d *ReportTeamDailyDao) InsertMul(timestamp int64, reportTeamDailys []models.ReportTeamDaily) (err error) {
	if len(reportTeamDailys) == 0 {
		return
	}

	_, err = d.Orm.InsertMulti(InsertMulCount, reportTeamDailys)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		_ = d.DelByTimestamp(timestamp)
		return err
	}

	return
}

func (d *ReportTeamDailyDao) DelByTimestamp(timestamp int64) (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("DELETE FROM %s WHERE %s = ? ", models.TABLE_ReportTeamDaily,
		models.COLUMN_ReportTeamDaily_Ctime), timestamp).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	return
}

type ReportTeam struct {
	Uid      uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	EusdBuy  int64  `orm:"column(eusd_buy)" json:"eusd_buy,omitempty"`
	EusdSell int64  `orm:"column(eusd_sell)" json:"eusd_sell,omitempty"`
}

func (d *ReportTeamDailyDao) InfoByUIds(uIds []string) (ReportTeam []ReportTeam, err error) {
	if len(uIds) == 0 {
		return
	}
	sqlQuery := fmt.Sprintf("SELECT %s,SUM(%s) AS eusd_buy, SUM(%s) AS eusd_sell FROM %s WHERE %s IN(%s) GROUP BY %s",
		models.COLUMN_ReportTeamDaily_Uid, models.COLUMN_ReportTeamDaily_EusdBuy, models.COLUMN_ReportTeamDaily_EusdSell,
		models.TABLE_ReportTeamDaily, models.COLUMN_ReportTeamDaily_Uid, strings.Join(uIds, ","), models.COLUMN_ReportTeamDaily_Uid)
	_, err = d.Orm.Raw(sqlQuery).QueryRows(&ReportTeam)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *ReportTeamDailyDao) Sum(startTime, endTime int64) (sumEusd ReportTeam, err error) {
	sqlQuery := fmt.Sprintf("SELECT SUM(%s) AS eusd_buy, SUM(%s) AS eusd_sell FROM %s WHERE %s=1 AND %s>=? AND %s<=? ",
		models.COLUMN_ReportTeamDaily_EusdBuy, models.COLUMN_ReportTeamDaily_EusdSell, models.TABLE_ReportTeamDaily,
		models.COLUMN_ReportTeamGameTransferDaily_Level, models.COLUMN_ReportTeamDaily_Ctime, models.COLUMN_ReportTeamDaily_Ctime)
	err = d.Orm.Raw(sqlQuery, startTime, endTime).QueryRow(&sumEusd)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *ReportTeamDailyDao) Personal(startTime, endTime int64, uid uint64) (sumEusd ReportTeam, err error) {
	sqlQuery := fmt.Sprintf("SELECT SUM(%s) AS eusd_buy, SUM(%s) AS eusd_sell FROM %s WHERE %s=? AND %s>=? AND %s<=? ",
		models.COLUMN_ReportTeamDaily_EusdBuy, models.COLUMN_ReportTeamDaily_EusdSell, models.TABLE_ReportTeamDaily,
		models.COLUMN_ReportTeamGameTransferDaily_Uid, models.COLUMN_ReportTeamDaily_Ctime, models.COLUMN_ReportTeamDaily_Ctime)
	err = d.Orm.Raw(sqlQuery, uid, startTime, endTime).QueryRow(&sumEusd)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}
