package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
	"utils/report/models"
)

type ReportGameRecordKyDao struct {
	common.BaseDao
}

func NewReportGameRecordKyDao(db string) *ReportGameRecordKyDao {
	return &ReportGameRecordKyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportGameRecordKyDaoEntity *ReportGameRecordKyDao

func (d *ReportGameRecordKyDao) DelByTimestamp(timestamp int64) (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("DELETE FROM %s WHERE %s = ? ", models.TABLE_ReportGameRecordKy,
		models.COLUMN_ReportGameRecordKy_Ctime), timestamp).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	return
}

func (d *ReportGameRecordKyDao) InsertMul(timestamp int64, reportGameRecordKys []models.ReportGameRecordKy) (err error) {
	if len(reportGameRecordKys) == 0 {
		return
	}
	_, err = d.Orm.InsertMulti(InsertMulCount, reportGameRecordKys)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return
}

func (d *ReportGameRecordKyDao) QueryTotalByTimestamp(timestamp int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT COUNT(%s) FROM %s WHERE %s=?", models.COLUMN_ReportGameRecordKy_Ctime,
		models.TABLE_ReportGameRecordKy, models.COLUMN_ReportGameRecordKy_Ctime), timestamp).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ReportGameRecordKyDao) QueryByTimestamp(timestamp int64, page, limit int) (list []BetInfo, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT %s,%s,%s,%s FROM %s WHERE %s=? LIMIT ? OFFSET ?",
		models.COLUMN_ReportGameRecordKy_Uid, models.COLUMN_ReportGameRecordKy_Bet, models.COLUMN_ReportGameRecordKy_ValidBet,
		models.COLUMN_ReportGameRecordKy_Profit, models.TABLE_ReportGameRecordKy, models.COLUMN_ReportGameRecordKy_Ctime),
		timestamp, limit, (page-1)*limit).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ReportGameRecordKyDao) IsBetPeriod(uid uint64, channelId uint32, start, over int64) (boolRes bool) {
	boolRes = false

	qs := d.Orm.QueryTable(models.TABLE_ReportGameRecordKy).Filter(models.COLUMN_ReportGameRecordKy_Uid, uid).Filter(models.ATTRIBUTE_ReportGameRecordKy_KindId, channelId)
	qs = qs.Filter(fmt.Sprintf("%s__gte", models.COLUMN_ReportGameRecordKy_StartTime), start)
	qs = qs.Filter(fmt.Sprintf("%s__lt", models.COLUMN_ReportGameRecordKy_StartTime), over)
	cnt, err := qs.Filter(fmt.Sprintf("%s__gt", models.COLUMN_ReportGameRecordKy_Bet), 0).Count()
	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
		return
	}
	if cnt > 0 {
		boolRes = true
	}
	return
}

//查询所有用户写入tmpuser表
func (d *ReportGameRecordKyDao) GetBetUsers(start, over int64, uid uint32, hasBet bool, gameId uint32) (num int64, reportGameRecordKy []interface{}, err error) {
	var reportGameRecord []models.ReportGameRecordKy
	where := ""
	if uid > 0 {
		myUid := strconv.Itoa(int(uid))
		where = fmt.Sprintf(" and %s=%s", models.COLUMN_ReportGameRecordKy_Uid, myUid)
	}
	if hasBet {
		newWhere := fmt.Sprintf(" and %s>0", models.COLUMN_ReportGameRecordKy_Bet)
		where = strings.Join([]string{where, newWhere}, "")

	}
	if gameId > 0 {
		newWhere := fmt.Sprintf(" and %s=%d", models.COLUMN_ReportGameRecordKy_KindId, gameId)
		where = strings.Join([]string{where, newWhere}, "")
	}
	sql := fmt.Sprintf(strings.Join([]string{"SELECT * FROM %s where %s>=? AND %s<?", where}, ""),
		models.TABLE_ReportGameRecordKy, models.COLUMN_ReportGameRecordKy_Ctime, models.COLUMN_ReportGameRecordKy_Ctime)

	num, err = d.Orm.Raw(sql, start, over).QueryRows(&reportGameRecord)

	reportGameRecordKy = make([]interface{}, len(reportGameRecord))
	for k, v := range reportGameRecord {
		reportGameRecordKy[k] = v
	}
	//logs.Trace("result : ", reportGameRecordAg)
	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
	}
	return
}

//查询投注用户
//func (d *ReportGameRecordKyDao) GetBetUsers(start, over int64, uid uint32, hasBet bool, gameId uint32) (reportGameRecordKy []interface{}, err error) {
//	where := ""
//	if uid > 0 {
//		myUid := strconv.Itoa(int(uid))
//		where = fmt.Sprintf(" and %s=%s", models.COLUMN_ReportGameRecordKy_Uid, myUid)
//	}
//	if hasBet {
//		newWhere := fmt.Sprintf(" and %s>0", models.COLUMN_ReportGameRecordKy_Bet)
//		where = strings.Join([]string{where, newWhere}, "")
//	}
//	if gameId > 0 {
//		newWhere := fmt.Sprintf(" and %s=%d", models.COLUMN_ReportGameRecordKy_KindId, gameId)
//		where = strings.Join([]string{where, newWhere}, "")
//	}
//	sql := fmt.Sprintf(strings.Join([]string{"SELECT * FROM %s where %s>=? AND %s<?", where}, ""),
//		models.TABLE_ReportGameRecordKy, models.COLUMN_ReportGameRecordKy_Ctime, models.COLUMN_ReportGameRecordKy_Ctime)
//
//	_, err = d.Orm.Raw(sql, start, over).QueryRows(&reportGameRecordKy)
//	if err != nil {
//		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
//	}
//	return
//}

//
func (d *ReportGameRecordKyDao) GetInfo(start, over int64) (reportKyList []models.ReportGameRecordKy, err error) {
	sql := fmt.Sprintf("SELECT %s,sum(%s) %s,sum(%s) %s,sum(%s) %s,sum(%s) %s  FROM %s where %s>=? and %s<? and %s>0 group by %s",
		models.COLUMN_ReportGameRecordKy_KindId, models.COLUMN_ReportGameRecordKy_Bet, models.COLUMN_ReportGameRecordKy_Bet,
		models.COLUMN_ReportGameRecordKy_ValidBet, models.COLUMN_ReportGameRecordKy_ValidBet,
		models.COLUMN_ReportGameRecordKy_Profit, models.COLUMN_ReportGameRecordKy_Profit,
		models.COLUMN_ReportGameRecordKy_Revenue, models.COLUMN_ReportGameRecordKy_Revenue,
		models.TABLE_ReportGameRecordKy, models.COLUMN_ReportGameRecordKy_Ctime,
		models.COLUMN_ReportGameRecordKy_Ctime, models.COLUMN_ReportGameRecordKy_Bet, models.COLUMN_ReportGameRecordKy_KindId)
	//logs.Debug("sql :", sql)
	_, err = d.Orm.Raw(sql, start, over).QueryRows(&reportKyList)

	logs.Debug("reportKyList is  : ", reportKyList)
	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
		return
	}

	return
}

//
func (d *ReportGameRecordKyDao) QueryAllByTime(start, over int64, page, limit int) (total int64, list []models.ReportGameRecordKy, err error) {
	qs := d.Orm.QueryTable(models.TABLE_ReportGameRecordKy)
	if start > 0 {
		qs = qs.Filter(fmt.Sprintf("%s__gte", models.COLUMN_ReportGameRecordKy_Ctime), start)
	}
	if over > 0 {
		qs = qs.Filter(fmt.Sprintf("%s__lt", models.COLUMN_ReportGameRecordKy_Ctime), over)
	}
	qs = qs.OrderBy(fmt.Sprintf("-%s", models.COLUMN_ReportGameRecordKy_Ctime)).Limit(limit)
	total, _ = qs.Count()
	_, err = qs.Offset((page - 1) * limit).All(&list)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}
