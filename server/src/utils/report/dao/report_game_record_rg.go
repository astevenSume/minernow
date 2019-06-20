package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
	"utils/report/models"
)

type ReportGameRecordRgDao struct {
	common.BaseDao
}

func NewReportGameRecordRgDao(db string) *ReportGameRecordRgDao {
	return &ReportGameRecordRgDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportGameRecordRgDaoEntity *ReportGameRecordRgDao

func (d *ReportGameRecordRgDao) DelByTimestamp(timestamp int64) (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("DELETE FROM %s WHERE %s = ? ", models.TABLE_ReportGameRecordRg,
		models.COLUMN_ReportGameRecordRg_Ctime), timestamp).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	return
}

func (d *ReportGameRecordRgDao) InsertMul(timestamp int64, reportGameRecordRgs []models.ReportGameRecordRg) (err error) {
	if len(reportGameRecordRgs) == 0 {
		return
	}
	_, err = d.Orm.InsertMulti(InsertMulCount, reportGameRecordRgs)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return
}

type BetInfo struct {
	Uid      uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Bet      int64  `orm:"column(bet)" json:"bet,omitempty"`
	ValidBet int64  `orm:"column(valid_bet)" json:"valid_bet,omitempty"`
	Profit   int64  `orm:"column(profit)" json:"profit,omitempty"`
}

func (d *ReportGameRecordRgDao) QueryTotalByTimestamp(timestamp int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT COUNT(%s) FROM %s WHERE %s=?", models.COLUMN_ReportGameRecordRg_Ctime,
		models.TABLE_ReportGameRecordRg, models.COLUMN_ReportGameRecordRg_Ctime), timestamp).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ReportGameRecordRgDao) QueryByTimestamp(timestamp int64, page, limit int) (list []BetInfo, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT %s,%s,%s,%s FROM %s WHERE %s=? LIMIT ? OFFSET ?",
		models.COLUMN_ReportGameRecordRg_Uid, models.COLUMN_ReportGameRecordRg_Bet, models.COLUMN_ReportGameRecordRg_ValidBet,
		models.COLUMN_ReportGameRecordRg_Profit, models.TABLE_ReportGameRecordRg, models.COLUMN_ReportGameRecordRg_Ctime),
		timestamp, limit, (page-1)*limit).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ReportGameRecordRgDao) IsBetPeriod(uid uint64, channelId uint32, start, over int64) (boolRes bool) {
	boolRes = false

	qs := d.Orm.QueryTable(models.TABLE_ReportGameRecordRg).Filter(models.COLUMN_ReportGameRecordRg_Uid, uid).Filter(models.COLUMN_ReportGameRecordRg_GameNameID, channelId)
	qs = qs.Filter(fmt.Sprintf("%s__gte", models.COLUMN_ReportGameRecordRg_BetTime), start)
	qs = qs.Filter(fmt.Sprintf("%s__lt", models.COLUMN_ReportGameRecordRg_BetTime), over)
	cnt, err := qs.Filter(fmt.Sprintf("%s__gt", models.COLUMN_ReportGameRecordRg_Bet), 0).Count()
	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
		return
	}
	if cnt > 0 {
		boolRes = true
	}
	return
}

//查询投注用户
func (d *ReportGameRecordRgDao) GetBetUsers(start, over int64, uid uint32, hasBet bool, gameId uint32) (num int64, reportGameRecordRg []interface{}, err error) {
	var reportGameRecord []models.ReportGameRecordRg
	where := ""
	if uid > 0 {
		myUid := strconv.Itoa(int(uid))
		where = fmt.Sprintf(" and %s=%s", models.COLUMN_ReportGameRecordRg_Uid, myUid)
	}
	if hasBet {
		newWhere := fmt.Sprintf(" and %s>0", models.COLUMN_ReportGameRecordRg_Bet)
		where = strings.Join([]string{where, newWhere}, "")
	}
	if gameId > 0 {
		newWhere := fmt.Sprintf(" and %s=%d", models.COLUMN_ReportGameRecordRg_GameNameID, gameId)
		where = strings.Join([]string{where, newWhere}, "")
	}
	sql := fmt.Sprintf(strings.Join([]string{"SELECT * FROM %s where %s>=? AND %s<?", where}, ""),
		models.TABLE_ReportGameRecordRg, models.COLUMN_ReportGameRecordRg_Ctime, models.COLUMN_ReportGameRecordRg_Ctime)

	num, err = d.Orm.Raw(sql, start, over).QueryRows(&reportGameRecord)

	reportGameRecordRg = make([]interface{}, len(reportGameRecord))
	for k, v := range reportGameRecord {
		reportGameRecordRg[k] = v
	}

	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
	}
	return
}

func (d *ReportGameRecordRgDao) GetInfo(start, over int64) (reportRgList []models.ReportGameRecordRg, err error) {
	sql := fmt.Sprintf("SELECT %s,sum(%s) %s,sum(%s) %s,sum(%s) %s FROM %s where %s>=? and %s<? and %s>0 group by %s",
		models.COLUMN_ReportGameRecordRg_GameNameID, models.COLUMN_ReportGameRecordRg_Bet, models.COLUMN_ReportGameRecordRg_Bet,
		models.COLUMN_ReportGameRecordRg_ValidBet, models.COLUMN_ReportGameRecordRg_ValidBet,
		models.COLUMN_ReportGameRecordRg_Profit, models.COLUMN_ReportGameRecordRg_Profit,
		models.TABLE_ReportGameRecordRg, models.COLUMN_ReportGameRecordRg_Ctime,
		models.COLUMN_ReportGameRecordRg_Ctime, models.COLUMN_ReportGameRecordRg_Bet, models.COLUMN_ReportGameRecordRg_GameNameID)
	//logs.Debug("sql :", sql)
	_, err = d.Orm.Raw(sql, start, over).QueryRows(&reportRgList)

	logs.Debug("reportRgList is  : ", reportRgList)
	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
		return
	}

	return
}

//
func (d *ReportGameRecordRgDao) QueryAllByTime(start, over int64, page, limit int) (total int64, list []models.ReportGameRecordRg, err error) {
	qs := d.Orm.QueryTable(models.TABLE_ReportGameRecordRg)
	if start > 0 {
		qs = qs.Filter(fmt.Sprintf("%s__gte", models.COLUMN_ReportGameRecordRg_Ctime), start)
	}
	if over > 0 {
		qs = qs.Filter(fmt.Sprintf("%s__lt", models.COLUMN_ReportGameRecordRg_Ctime), over)
	}
	qs = qs.OrderBy(fmt.Sprintf("-%s", models.COLUMN_ReportGameRecordRg_Ctime)).Limit(limit)
	total, _ = qs.Count()
	_, err = qs.Offset((page - 1) * limit).All(&list)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}
