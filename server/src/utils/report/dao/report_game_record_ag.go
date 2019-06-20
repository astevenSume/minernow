package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
	"utils/report/models"
)

type ReportGameRecordAgDao struct {
	common.BaseDao
}

func NewReportGameRecordAgDao(db string) *ReportGameRecordAgDao {
	return &ReportGameRecordAgDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportGameRecordAgDaoEntity *ReportGameRecordAgDao

//const COLUMN_RESULTVALUE = "result_value"
//type SumAgInfo struct {
//	ResultValue int64 `orm:"column(result_value)" json:"result_value,omitempty"`
//}

func (d *ReportGameRecordAgDao) DelByTimestamp(timestamp int64) (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("DELETE FROM %s WHERE %s = ? ", models.TABLE_ReportGameRecordAg,
		models.COLUMN_ReportGameRecordAg_Ctime), timestamp).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	return
}

func (d *ReportGameRecordAgDao) InsertMul(timestamp int64, reportGameRecordAgs []models.ReportGameRecordAg) (err error) {
	if len(reportGameRecordAgs) == 0 {
		return
	}
	_, err = d.Orm.InsertMulti(InsertMulCount, reportGameRecordAgs)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return
}

func (d *ReportGameRecordAgDao) QueryTotalByTimestamp(timestamp int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT COUNT(%s) FROM %s WHERE %s=?", models.COLUMN_ReportGameRecordAg_Ctime,
		models.TABLE_ReportGameRecordAg, models.COLUMN_ReportGameRecordAg_Ctime), timestamp).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ReportGameRecordAgDao) QueryByTimestamp(timestamp int64, page, limit int) (list []BetInfo, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT %s,%s,%s,%s FROM %s WHERE %s=? LIMIT ? OFFSET ?",
		models.COLUMN_ReportGameRecordAg_Uid, models.COLUMN_ReportGameRecordAg_Bet, models.COLUMN_ReportGameRecordAg_ValidBet,
		models.COLUMN_ReportGameRecordAg_Profit, models.TABLE_ReportGameRecordAg, models.COLUMN_ReportGameRecordAg_Ctime),
		timestamp, limit, (page-1)*limit).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

//查询所有用户写入tmpuser表
func (d *ReportGameRecordAgDao) GetBetUsers(start, over int64, uid uint32, hasBet bool, gameId uint32) (num int64, reportGameRecordAg []interface{}, err error) {
	var reportGameRecord []models.ReportGameRecordAg
	where := ""
	if uid > 0 {
		myUid := strconv.Itoa(int(uid))
		where = fmt.Sprintf(" and %s=%s", models.COLUMN_ReportGameRecordAg_Uid, myUid)
	}
	if hasBet {
		newWhere := fmt.Sprintf(" and %s>0", models.COLUMN_ReportGameRecordAg_Bet)
		where = strings.Join([]string{where, newWhere}, "")
	}
	if gameId > 0 {
		newWhere := fmt.Sprintf(" and %s=%d", models.COLUMN_ReportGameRecordAg_GameType, gameId)
		where = strings.Join([]string{where, newWhere}, "")
	}
	sql := fmt.Sprintf(strings.Join([]string{"SELECT * FROM %s where %s>=? AND %s<?", where}, ""),
		models.TABLE_ReportGameRecordAg, models.COLUMN_ReportGameRecordAg_Ctime, models.COLUMN_ReportGameRecordAg_Ctime)

	num, err = d.Orm.Raw(sql, start, over).QueryRows(&reportGameRecord)

	//if num <= 0 {
	//	//tmpdao.TmpGameBetersDaoEntity.InsertUser(u.Uid, channelId, 0, tmpdao.NOTBET, start)
	//}

	reportGameRecordAg = make([]interface{}, len(reportGameRecord))
	for k, v := range reportGameRecord {
		reportGameRecordAg[k] = v
	}
	//logs.Trace("result : ", reportGameRecordAg)
	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
	}
	return
}

func (d *ReportGameRecordAgDao) IsBetPeriod(uid uint64, start, over int64) (boolRes bool) {
	boolRes = false

	qs := d.Orm.QueryTable(models.TABLE_ReportGameRecordAg).Filter(models.COLUMN_ReportGameRecordAg_Uid, uid)
	qs = qs.Filter(fmt.Sprintf("%s__gte", models.COLUMN_ReportGameRecordAg_Ctime), start)
	qs = qs.Filter(fmt.Sprintf("%s__lt", models.COLUMN_ReportGameRecordAg_Ctime), over)
	cnt, err := qs.Filter(fmt.Sprintf("%s__gt", models.COLUMN_ReportGameRecordAg_Bet), 0).Count()
	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
		return
	}
	if cnt > 0 {
		boolRes = true
	}
	return
}

//
func (d *ReportGameRecordAgDao) GetInfo(start, over int64) (reportAgList []models.ReportGameRecordAg, err error) {
	sql := fmt.Sprintf("SELECT %s,sum(%s) %s,sum(%s) %s,sum(%s) %s FROM %s where %s>=? and %s<? and %s>0 group by %s",
		models.COLUMN_ReportGameRecordAg_GameType, models.COLUMN_ReportGameRecordAg_Bet, models.COLUMN_ReportGameRecordAg_Bet,
		models.COLUMN_ReportGameRecordAg_ValidBet, models.COLUMN_ReportGameRecordAg_ValidBet,
		models.COLUMN_ReportGameRecordAg_Profit, models.COLUMN_ReportGameRecordAg_Profit,
		models.TABLE_ReportGameRecordAg, models.COLUMN_ReportGameRecordAg_Ctime,
		models.COLUMN_ReportGameRecordAg_Ctime, models.COLUMN_ReportGameRecordAg_Bet, models.COLUMN_ReportGameRecordAg_GameType)
	//logs.Debug("sql :", sql)
	_, err = d.Orm.Raw(sql, start, over).QueryRows(&reportAgList)

	logs.Debug("reportAgList is  : ", reportAgList)
	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
		return
	}

	return
}

//
func (d *ReportGameRecordAgDao) QueryAllByTime(start, over int64, page, limit int) (total int64, list []models.ReportGameRecordAg, err error) {
	qs := d.Orm.QueryTable(models.TABLE_ReportGameRecordAg)
	if start > 0 {
		qs = qs.Filter(fmt.Sprintf("%s__gte", models.COLUMN_ReportGameRecordAg_Ctime), start)
	}
	if over > 0 {
		qs = qs.Filter(fmt.Sprintf("%s__lt", models.COLUMN_ReportGameRecordAg_Ctime), over)
	}
	qs = qs.OrderBy(fmt.Sprintf("-%s", models.COLUMN_ReportGameRecordAg_Ctime)).Limit(limit)
	total, _ = qs.Count()
	_, err = qs.Offset((page - 1) * limit).All(&list)

	if err != nil {
		common.LogFuncError("sql err: %v", err)
		return
	}
	return
}
