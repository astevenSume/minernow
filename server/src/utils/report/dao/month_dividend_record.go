package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	"utils/report/models"
)

const Received = 1   //上级已发放
const NoReceived = 2 //上级未发放状态

const Pay = 1   //已发放给下级
const NoPay = 2 //未发放给下级
type MonthDividendRecordDao struct {
	common.BaseDao
}

func NewMonthDividendRecordDao(db string) *MonthDividendRecordDao {
	return &MonthDividendRecordDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var MonthDividendRecordDaoEntity *MonthDividendRecordDao

func (d *MonthDividendRecordDao) FindByDataRange(start, end int64) (reports []*models.MonthDividendRecord, err error) {
	reports = make([]*models.MonthDividendRecord, 0)
	querySql := fmt.Sprintf("select * from %s where %s>=? and %s<=?", models.TABLE_MonthDividendRecord, models.COLUMN_MonthDividendRecord_Ctime, models.COLUMN_MonthDividendRecord_Ctime)
	_, err = d.Orm.Raw(querySql, start, end).QueryRows(&reports)
	if err != nil {
		common.LogFuncError("data err %v", err)
	}
	return
}

func (d *MonthDividendRecordDao) FindByRange(start, over int64) (report models.MonthDividendRecord, err error) {
	var where = ""
	if start > 0 {
		where = fmt.Sprintf("and %s>=%d", models.COLUMN_MonthDividendRecord_Ctime, start)
	}
	if over > 0 {
		newWhere := fmt.Sprintf(" and %s<%d", models.COLUMN_MonthDividendRecord_Ctime, over)
		where = strings.Join([]string{where, newWhere}, "")
	}

	querySql := fmt.Sprintf("select sum(%s) %s from %s where 1=1 %s",
		models.COLUMN_MonthDividendRecord_ResultDividend, models.COLUMN_MonthDividendRecord_ResultDividend,
		models.TABLE_MonthDividendRecord, where)
	err = d.Orm.Raw(querySql).QueryRow(&report)
	if err != nil {
		common.LogFuncError("data err %v", err)
	}
	return
}

func (d *MonthDividendRecordDao) Insert(records []*models.MonthDividendRecord) (err error) {
	_, err = d.Orm.InsertMulti(100, records)
	if err != nil {
		common.LogFuncError("err %v ", err)
	}
	return
}
func (d *MonthDividendRecordDao) UpdateReceiveStatus(uid uint64, nowStatus, beforeStatus int) (err error) {
	querySql := fmt.Sprintf("update %s set %s=? where %s=? and %s=? ",
		models.TABLE_MonthDividendRecord, models.COLUMN_MonthDividendRecord_ReceiveStatus, models.COLUMN_MonthDividendRecord_Uid, models.COLUMN_MonthDividendRecord_ReceiveStatus)
	_, err = d.Orm.Raw(querySql, nowStatus, uid, beforeStatus).Exec()
	if err != nil {
		common.LogFuncError("data err %v", err)
	}
	return
}

func (d *MonthDividendRecordDao) UpdatePayStatus(uid uint64, nowStatus, beforeStatus int) (err error) {
	querySql := fmt.Sprintf("update %s set %s=? where %s=? and %s=? ",
		models.TABLE_MonthDividendRecord, models.COLUMN_MonthDividendRecord_PayStatus, models.COLUMN_MonthDividendRecord_Uid, models.COLUMN_MonthDividendRecord_PayStatus)
	_, err = d.Orm.Raw(querySql, nowStatus, uid, beforeStatus).Exec()
	if err != nil {
		common.LogFuncError("data err %v", err)
	}
	return
}

func (d *MonthDividendRecordDao) GetById(id uint64) (reports *models.MonthDividendRecord, err error) {
	querySql := fmt.Sprintf("select * from %s where %s=?", models.TABLE_MonthDividendRecord, models.COLUMN_MonthDividendRecord_Id)
	err = d.Orm.Raw(querySql, id).QueryRow(&reports)

	if err != nil {
		common.LogFuncError("data err %v", err)
	}
	return
}

func (d *MonthDividendRecordDao) GetByUidAndCtime(uid uint64, ctime int64) (reports *models.MonthDividendRecord, err error) {
	querySql := fmt.Sprintf("select * from %s where %s=? and %s=?", models.TABLE_MonthDividendRecord, models.COLUMN_MonthDividendRecord_Id, models.COLUMN_MonthDividendRecord_Ctime)
	err = d.Orm.Raw(querySql, uid, ctime).QueryRow(&reports)

	if err != nil {
		common.LogFuncError("data err %v", err)
	}
	return
}
func (d *MonthDividendRecordDao) FindByUid(uid uint64) (reports []*models.MonthDividendRecord, err error) {
	reports = make([]*models.MonthDividendRecord, 0)
	querySql := fmt.Sprintf("select * from %s where %s=? order by %s desc", models.TABLE_MonthDividendRecord, models.COLUMN_MonthDividendRecord_Uid, models.COLUMN_MonthDividendRecord_Ctime)
	_, err = d.Orm.Raw(querySql, uid).QueryRows(&reports)
	if err != nil {
		common.LogFuncError("data err %v", err)
	}
	return
}

func (d *MonthDividendRecordDao) FindByUidsAndStatus(uids []uint64, status int32) (reports []*models.MonthDividendRecord, err error) {
	qs := d.Orm.QueryTable(models.TABLE_MonthDividendRecord)
	_, err = qs.Filter(models.COLUMN_MonthDividendRecord_Uid+"__in", uids).Filter(models.COLUMN_MonthDividendRecord_ReceiveStatus, status).All(&reports)
	if err != nil {
		common.LogFuncError("data err %v", err)
	}
	return
}

func (d *MonthDividendRecordDao) FindByUidsAndCtime(uids []uint64, cTime int64) (reports []*models.MonthDividendRecord, err error) {
	qs := d.Orm.QueryTable(models.TABLE_MonthDividendRecord)
	_, err = qs.Filter(models.COLUMN_MonthDividendRecord_Uid+"__in", uids).Filter(models.COLUMN_MonthDividendRecord_Ctime, cTime).All(&reports)
	if err != nil {
		common.LogFuncError("data err %v", err)
	}
	return
}
func (d *MonthDividendRecordDao) Count() (rowNum int64, err error) {
	rowNum, err = d.Orm.QueryTable(models.TABLE_MonthDividendRecord).Count()
	if err != nil {
		common.LogFuncError("err %v ", err)
	}
	return
}

func (d *MonthDividendRecordDao) DeleteAllData() (err error) {
	sql := fmt.Sprintf("delete from %s", models.TABLE_MonthDividendRecord)
	_, err = d.Orm.Raw(sql).Exec()
	if err != nil {
		common.LogFuncError("err %v", err)
	}
	return
}

type SalaryMonthSummary struct {
	SelfDividend int64 `orm:"column(self_dividend)" json:"self_dividend,omitempty"`
	Ctime        int64 `orm:"column(ctime)" json:"ctime,omitempty"`
	Status       uint8 `orm:"column(receive_status)" json:"receive_status,omitempty"`
}

func (d *MonthDividendRecordDao) GetListByTime(uid uint64, bTime, eTime int64) (salaryMonth []SalaryMonthSummary, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT %s, %s, %s FROM %s WHERE %s=? AND (%s>=? AND %s<=?) ORDER BY %s DESC",
		models.COLUMN_MonthDividendRecord_Ctime, models.COLUMN_MonthDividendRecord_SelfDividend,
		models.COLUMN_MonthDividendRecord_ReceiveStatus, models.TABLE_MonthDividendRecord, models.COLUMN_MonthDividendRecord_Uid,
		models.COLUMN_MonthDividendRecord_Ctime, models.COLUMN_MonthDividendRecord_Ctime,
		models.COLUMN_MonthDividendRecord_Ctime), uid, bTime, eTime).QueryRows(&salaryMonth)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (d *MonthDividendRecordDao) GetSumMonthResultDivided(uid uint64) (sumMonthDivided int64, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT SUM(%s) FROM %s WHERE %s=?",
		models.COLUMN_MonthDividendRecord_ResultDividend, models.TABLE_MonthDividendRecord,
		models.COLUMN_MonthDividendRecord_Uid), uid).QueryRow(&sumMonthDivided)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

//上一个月分红
func (d *MonthDividendRecordDao) GetPreMonthResultDivided(uid uint64, begin, end int64) (monthDivided int64, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT %s FROM %s WHERE %s=? AND (%s>=? AND %s<?) ",
		models.COLUMN_MonthDividendRecord_ResultDividend, models.TABLE_MonthDividendRecord,
		models.COLUMN_MonthDividendRecord_Uid, models.COLUMN_MonthDividendRecord_Ctime,
		models.COLUMN_MonthDividendRecord_Ctime), uid, begin, end).QueryRow(&monthDivided)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("error:%v", err)
		return
	}

	return
}
func (d *MonthDividendRecordDao) Update(record *models.MonthDividendRecord) (err error) {
	_, err = d.Orm.Update(record, models.COLUMN_MonthDividendRecord_ReceiveStatus)
	if err != nil {
		common.LogFuncError("Update MonthDividendRecord db err %v", err)
	}
	return
}

type ReportTeamMonthBonus struct {
	Uid            uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	ResultDividend int64  `orm:"column(result_dividend)" json:"month_bonus,omitempty"`
}

func (d *MonthDividendRecordDao) InfoByUIds(uIds []string) (teamMonthBonus []ReportTeamMonthBonus, err error) {
	if len(uIds) == 0 {
		return
	}
	sqlQuery := fmt.Sprintf("SELECT %s,%s FROM %s WHERE %s IN(%s)", models.COLUMN_MonthDividendRecord_Uid,
		models.COLUMN_MonthDividendRecord_ResultDividend, models.TABLE_MonthDividendRecord,
		models.COLUMN_MonthDividendRecord_Uid, strings.Join(uIds, ","))
	_, err = d.Orm.Raw(sqlQuery).QueryRows(&teamMonthBonus)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *MonthDividendRecordDao) Sum(startTime, endTime int64) (sumMonthBonus ReportTeamMonthBonus, err error) {
	sqlQuery := fmt.Sprintf("SELECT SUM(%s) AS result_dividend FROM %s WHERE %s=1 AND %s>=? AND %s<=? ",
		models.COLUMN_MonthDividendRecord_ResultDividend, models.TABLE_MonthDividendRecord,
		models.COLUMN_MonthDividendRecord_Level, models.COLUMN_MonthDividendRecord_Ctime, models.COLUMN_MonthDividendRecord_Ctime)
	err = d.Orm.Raw(sqlQuery, startTime, endTime).QueryRow(&sumMonthBonus)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *MonthDividendRecordDao) Personal(startTime, endTime int64, uid uint64) (sumMonthBonus ReportTeamMonthBonus, err error) {
	sqlQuery := fmt.Sprintf("SELECT SUM(%s) AS result_dividend FROM %s WHERE %s=? AND %s>=? AND %s<=? ",
		models.COLUMN_MonthDividendRecord_ResultDividend, models.TABLE_MonthDividendRecord,
		models.COLUMN_MonthDividendRecord_Uid, models.COLUMN_MonthDividendRecord_Ctime, models.COLUMN_MonthDividendRecord_Ctime)
	err = d.Orm.Raw(sqlQuery, uid, startTime, endTime).QueryRow(&sumMonthBonus)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}
