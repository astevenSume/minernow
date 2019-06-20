package dao

import (
	"common"
	"fmt"
	"strings"
	"utils/report/models"
)

type ReportTeamGameTransferDailyDao struct {
	common.BaseDao
}

func NewReportTeamGameTransferDailyDao(db string) *ReportTeamGameTransferDailyDao {
	return &ReportTeamGameTransferDailyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportTeamGameTransferDailyDaoEntity *ReportTeamGameTransferDailyDao

func (d *ReportTeamGameTransferDailyDao) InsertMul(timestamp int64, reportTeamGameTransferDaily []models.ReportTeamGameTransferDaily) (err error) {
	if len(reportTeamGameTransferDaily) == 0 {
		return
	}

	_, err = d.Orm.InsertMulti(InsertMulCount, reportTeamGameTransferDaily)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		_ = d.DelByTimestamp(timestamp)
		return err
	}

	return
}

func (d *ReportTeamGameTransferDailyDao) DelByTimestamp(timestamp int64) (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("DELETE FROM %s WHERE %s = ? ", models.TABLE_ReportTeamGameTransferDaily,
		models.COLUMN_ReportTeamGameTransferDaily_Ctime), timestamp).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	return
}

type ReportTeamGameTransfer struct {
	Uid          uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	TeamRecharge int64  `orm:"column(recharge)" json:"recharge,omitempty"`
	TeamWithdraw int64  `orm:"column(withdraw)" json:"withdraw,omitempty"`
}

func (d *ReportTeamGameTransferDailyDao) InfoByUIds(channelId uint32, uIds []string) (teamTransfer []ReportTeamGameTransfer, err error) {
	if len(uIds) == 0 {
		return
	}

	var param []interface{}
	sqlQuery := fmt.Sprintf("SELECT %s,SUM(%s) AS recharge, SUM(%s) AS withdraw FROM %s WHERE %s IN(%s)",
		models.COLUMN_ReportTeamGameTransferDaily_Uid, models.COLUMN_ReportTeamGameTransferDaily_TeamRecharge,
		models.COLUMN_ReportTeamGameTransferDaily_TeamWithdraw, models.TABLE_ReportTeamGameTransferDaily,
		models.COLUMN_ReportTeamGameTransferDaily_Uid, strings.Join(uIds, ","))
	if channelId > 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportTeamGameTransferDaily_ChannelId)
		param = append(param, channelId)
	}
	sqlQuery = fmt.Sprintf("%s GROUP BY %s ", sqlQuery, models.COLUMN_ReportTeamGameTransferDaily_Uid)

	_, err = d.Orm.Raw(sqlQuery, param...).QueryRows(&teamTransfer)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *ReportTeamGameTransferDailyDao) Sum(channelId uint32, startTime, endTime int64) (sumTransfer ReportTeamGameTransfer, err error) {
	var param []interface{}
	param = append(param, startTime)
	param = append(param, endTime)
	sqlQuery := fmt.Sprintf("SELECT SUM(%s) AS recharge, SUM(%s) AS withdraw FROM %s WHERE %s=1 AND %s>=? AND %s<=? ",
		models.COLUMN_ReportTeamGameTransferDaily_TeamRecharge, models.COLUMN_ReportTeamGameTransferDaily_TeamWithdraw,
		models.TABLE_ReportTeamGameTransferDaily, models.COLUMN_ReportTeamGameTransferDaily_Level,
		models.COLUMN_ReportTeamGameTransferDaily_Ctime, models.COLUMN_ReportTeamGameTransferDaily_Ctime)
	if channelId > 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportTeamGameTransferDaily_ChannelId)
		param = append(param, channelId)
	}

	err = d.Orm.Raw(sqlQuery, param...).QueryRow(&sumTransfer)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *ReportTeamGameTransferDailyDao) Personal(channelId uint32, startTime, endTime int64, uid uint64) (sumTransfer ReportTeamGameTransfer, err error) {
	var param []interface{}
	param = append(param, uid)
	param = append(param, startTime)
	param = append(param, endTime)
	sqlQuery := fmt.Sprintf("SELECT SUM(%s) AS recharge, SUM(%s) AS withdraw FROM %s WHERE %s=? AND %s>=? AND %s<=? ",
		models.COLUMN_ReportTeamGameTransferDaily_TeamRecharge, models.COLUMN_ReportTeamGameTransferDaily_TeamWithdraw,
		models.TABLE_ReportTeamGameTransferDaily, models.COLUMN_ReportTeamGameTransferDaily_Uid,
		models.COLUMN_ReportTeamGameTransferDaily_Ctime, models.COLUMN_ReportTeamGameTransferDaily_Ctime)
	if channelId > 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportTeamGameTransferDaily_ChannelId)
		param = append(param, channelId)
	}

	err = d.Orm.Raw(sqlQuery, param...).QueryRow(&sumTransfer)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}
