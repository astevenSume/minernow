package dao

import (
	"common"
	utils "eusd/eosplus"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
	"utils/report/models"
)

type ReportGameUserDailyDao struct {
	common.BaseDao
}

func NewReportGameUserDailyDao(db string) *ReportGameUserDailyDao {
	return &ReportGameUserDailyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportGameUserDailyDaoEntity *ReportGameUserDailyDao

//领取状态
const (
	ReportGameUserDailyGet     = iota + 1 //已领取
	ReportGameUserDailyUnGet              //未领取
	ReportGameUserDailyUnGrant            //未发放
)

type ReportGameUserSalary struct {
	Uid    uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Salary int64  `orm:"column(salary)" json:"salary,omitempty"`
}

func (d *ReportGameUserDailyDao) DelByTimestamp(channelId uint32, timestamp int64) (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("DELETE FROM %s WHERE %s=? AND %s = ? ", models.TABLE_ReportGameUserDaily,
		models.COLUMN_ReportGameUserDaily_ChannelId, models.COLUMN_ReportGameUserDaily_Ctime), channelId, timestamp).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	return
}

func (d *ReportGameUserDailyDao) QueryTotalByChannelId(channelId uint32, timestamp int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(*) FROM %s WHERE %s=? AND %s=?",
		models.TABLE_ReportGameUserDaily, models.COLUMN_ReportGameUserDaily_ChannelId,
		models.COLUMN_ReportGameUserDaily_Ctime), channelId, timestamp).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ReportGameUserDailyDao) QueryByChannelId(channelId uint32, timestamp int64, page, limit int) (list []models.ReportGameUserDaily, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s=? AND %s=? LIMIT ? OFFSET ?",
		models.TABLE_ReportGameUserDaily, models.COLUMN_ReportGameUserDaily_ChannelId,
		models.COLUMN_ReportGameUserDaily_Ctime), channelId, timestamp, limit, (page-1)*limit).QueryRows(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ReportGameUserDailyDao) InsertMul(channelId uint32, timestamp int64, reportGameUserDailys []models.ReportGameUserDaily) (err error) {
	if len(reportGameUserDailys) == 0 {
		return
	}
	common.LogFuncDebug("ReportGameUserDailyDao InsertMul")
	_, err = d.Orm.InsertMulti(InsertMulCount, reportGameUserDailys)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		d.DelByTimestamp(channelId, timestamp)
		return err
	}

	return
}

type SalaryDaySummary struct {
	Salary int64 `orm:"column(salary)" json:"salary,omitempty"`
	Status uint8 `orm:"column(status)" json:"is_get,omitempty"`
	Ctime  int64 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (d *ReportGameUserDailyDao) GetListByTime(uid uint64, bTime, eTime int64) (salaryDay []SalaryDaySummary, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT %s, %s,SUM(%s)AS salary FROM %s WHERE %s=? AND (%s>=? AND %s<=?) GROUP BY "+
		"%s,%s ORDER BY %s DESC", models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Status,
		models.COLUMN_ReportGameUserDaily_Salary, models.TABLE_ReportGameUserDaily, models.COLUMN_ReportGameUserDaily_Uid,
		models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Uid,
		models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Ctime), uid, bTime, eTime).QueryRows(&salaryDay)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (d *ReportGameUserDailyDao) GetYesterdayTotalBet(uid uint64) (totalBet int64, err error) {
	timestamp := common.GetZeroTime(time.Now().Unix()) - common.DaySeconds
	err = d.Orm.Raw(fmt.Sprintf("SELECT SUM(%s) FROM %s WHERE %s=? AND %s=?",
		models.COLUMN_ReportGameUserDaily_TotalValidBet, models.TABLE_ReportGameUserDaily,
		models.COLUMN_ReportGameUserDaily_Uid, models.COLUMN_ReportGameUserDaily_Ctime), uid, timestamp).QueryRow(&totalBet)
	if err != nil {
		common.LogFuncError("error:%v")
		return
	}
	return
}

func (d *ReportGameUserDailyDao) GetDayTotalSalary(uid uint64, timestamp int64) (totalSalary int64, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT SUM(%s) FROM %s WHERE %s=? AND %s=?",
		models.COLUMN_ReportGameUserDaily_Salary, models.TABLE_ReportGameUserDaily,
		models.COLUMN_ReportGameUserDaily_Uid, models.COLUMN_ReportGameUserDaily_Ctime), uid, timestamp).QueryRow(&totalSalary)
	if err != nil {
		common.LogFuncError("error:%v")
		return
	}
	return
}

func (d *ReportGameUserDailyDao) GetDetailDaySalaryByUid(uid uint64, timestamp int64) (salary []models.ReportGameUserDaily, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_ReportGameUserDaily).
		Filter(models.COLUMN_ReportGameUserDaily_Uid, uid).
		Filter(models.COLUMN_ReportGameUserDaily_Ctime, timestamp).
		All(&salary)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (d *ReportGameUserDailyDao) UpdateSalaryGet(uid uint64, timestamp int64) (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=? AND %s=? ", models.TABLE_ReportGameUserDaily,
		models.COLUMN_ReportGameUserDaily_Status, models.COLUMN_ReportGameUserDaily_Uid, models.COLUMN_ReportGameUserDaily_Ctime),
		ReportGameUserDailyGet, uid, timestamp).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (d *ReportGameUserDailyDao) SetSalary(uid uint64, timestamp, salary int64) (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=? AND %s=? ", models.TABLE_ReportGameUserDaily,
		models.COLUMN_ReportGameUserDaily_Salary, models.COLUMN_ReportGameUserDaily_Uid,
		models.COLUMN_ReportGameUserDaily_Ctime), salary, uid, timestamp).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (d *ReportGameUserDailyDao) QueryTotalByData(startDay int64, endDay int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(*) FROM %s WHERE %s>=? AND %s<=?",
		models.TABLE_ReportGameUserDaily, models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Ctime), startDay, endDay).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

//todo ctime存的是时间戳
func (d *ReportGameUserDailyDao) FindByDayRangeAndChannelId(startDay int64, endDay int64, channelId uint32, page, limit int) (reports []models.ReportGameUserDaily, err error) {
	querySql := fmt.Sprintf("SELECT * FROM %s WHERE %s>=? AND %s<=? AND %s=? ORDER BY %s DESC LIMIT ?,?", models.TABLE_ReportGameUserDaily,
		models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Ctime, models.ATTRIBUTE_ReportGameUserDaily_ChannelId,
		models.ATTRIBUTE_ReportGameUserDaily_Ctime)

	_, err = d.Orm.Raw(querySql, startDay, endDay, channelId, (page-1)*limit, limit).QueryRows(&reports)
	if err != nil {
		common.LogFuncError("err %v", err)
	}
	return
}

//todo 改成分页查询，ctime存的是时间戳
func (d *ReportGameUserDailyDao) FindByDayRange(startDay int64, endDay int64, page, limit int) (reports []models.ReportGameUserDaily, err error) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s>=? AND %s<=? ORDER BY %s DESC LIMIT ?,?", models.TABLE_ReportGameUserDaily,
		models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Ctime)
	_, err = d.Orm.Raw(sql, startDay, endDay, (page-1)*limit, limit).QueryRows(&reports)
	if err != nil {
		common.LogFuncError("err %v", err)
	}
	return
}
func (d *ReportGameUserDailyDao) InsertForTest(reports []*models.ReportGameUserDaily) (err error) {
	num, err := d.Orm.InsertMulti(100, &reports)
	if err != nil {
		common.LogFuncError("err %s", err)
	}
	fmt.Println("num ", num)
	return
}

func (d *ReportGameUserDailyDao) DeleteForTest(uid []uint64) (err error) {
	sql := fmt.Sprintf("delete from %s", models.TABLE_ReportGameUserDaily)
	if len(uid) != 0 {
		sql = fmt.Sprintf("delete from %s where %s in ?", models.TABLE_ReportGameUserDaily, models.COLUMN_GameUserMonthReport_Uid)
	}
	_, err = d.Orm.Raw(sql).Exec()
	if err != nil {
		common.LogFuncError("err %v", err)
	}
	return
}

const (
	GAME_CHANNEL_KY = uint32(iota + 1) //开元棋牌
	GAME_CHANNEL_AG                    //AG
	GAME_CHANNEL_RG                    //彩票
)

type PositiveUser struct {
	Uid      uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Bet      int64  `orm:"column(bet)" json:"bet,omitempty"`
	ValidBet int64  `orm:"column(valid_bet)" json:"valid_bet,omitempty"`
	Datenums uint32 `orm:"column(date_nums)" json:"date_nums,omitempty"`
}

func (d *ReportGameUserDailyDao) GetPositiveUserNums(channelId uint32) (num uint32) {
	t := time.Now()
	monthFirstDayUnix := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Unix()
	monthFirstDate := time.Date(t.Year(), t.Month(), 1, 23, 59, 59, 0, t.Location())
	monthLastDayUnix := monthFirstDate.AddDate(0, 1, -1).Unix()
	var positiveUser = []PositiveUser{}
	sql := fmt.Sprintf("select %s,sum(1) %s, sum(%s) bet, sum(%s) valid_bet from %s where %s=? and %s>=? and %s<? group by %s,%s",
		models.COLUMN_ReportGameUserDaily_Uid, "date_nums", models.COLUMN_ReportGameUserDaily_Bet, models.COLUMN_ReportGameUserDaily_ValidBet,
		models.TABLE_ReportGameUserDaily, models.COLUMN_ReportGameUserDaily_ChannelId,
		models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Ctime,
		models.COLUMN_ReportGameUserDaily_Uid, models.COLUMN_ReportGameUserDaily_ChannelId)
	sql_again := fmt.Sprintf("select t.* from (%s) as t where t.date_nums>? and t.bet>=? or t.valid_bet>=?", sql)
	d.Orm.Raw(sql_again, channelId, monthFirstDayUnix, monthLastDayUnix, POSITIVE_DAY_NUMS, utils.QuantityFloat64ToInt64(POSITIVE_BET_LOWEST), utils.QuantityFloat64ToInt64(POSITIVE_VALIDBET_LOWEST)).QueryRows(&positiveUser)

	num = uint32(len(positiveUser))
	return
}

//get salary by channelId
func (d *ReportGameUserDailyDao) GetMonthSalaryByChannelId(channelId uint32, ctime int64) (salary int64, err error) {
	var salarys = models.ReportGameUserDaily{}

	sql := fmt.Sprintf("SELECT SUM(%s) %s FROM %s WHERE %s=? and %s=?",
		models.COLUMN_ReportGameUserDaily_Salary, models.COLUMN_ReportGameUserDaily_Salary, models.TABLE_ReportGameUserDaily,
		models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_ChannelId)

	if err := d.Orm.Raw(sql, ctime, channelId).QueryRow(&salarys); err != nil {
		common.LogFuncError("sql err %v", err)
	}
	logs.Warn("salary record :", salarys)
	salary = salarys.Salary

	return
}

//get salary by channelId
func (d *ReportGameUserDailyDao) GetSalaryByChannelId(channelId uint32, ctime int64) (salary int64, err error) {
	var salarys = models.ReportGameUserDaily{}

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s=? and %s=?",
		models.COLUMN_ReportGameUserDaily_Salary, models.TABLE_ReportGameUserDaily,
		models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_ChannelId)

	if err := d.Orm.Raw(sql, ctime, channelId).QueryRow(&salarys); err != nil {
		common.LogFuncError("sql err %v", err)
	}
	logs.Warn("salary record :", salarys)
	salary = salarys.Salary
	return
}

func (d *ReportGameUserDailyDao) QueryTotalTeamByTime(channelId uint32, startTime, endTime int64, uid, pUid uint64) (total int, err error) {
	var param []interface{}
	sqlQuery := fmt.Sprintf("SELECT COUNT(%s) FROM %s WHERE %s>=? AND %s<=?", models.COLUMN_ReportGameUserDaily_Uid,
		models.TABLE_ReportGameUserDaily, models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Ctime)
	param = append(param, startTime)
	param = append(param, endTime)
	if channelId > 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportGameUserDaily_ChannelId)
		param = append(param, channelId)
	}
	if uid > 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportGameUserDaily_Uid)
		param = append(param, uid)
	}
	if pUid > 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportGameUserDaily_PUid)
		param = append(param, pUid)
	}
	if uid == 0 && pUid == 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportGameUserDaily_Level)
		param = append(param, 1)
	}

	err = d.Orm.Raw(sqlQuery, param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
	/*var qbTotal orm.QueryBuilder
	qbTotal, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}

	// 构建查询对象
	qbTotal.Select("Count(*) AS total").
		From(models.TABLE_ReportGameUserDaily).
		Where("1 = 1")

	var param []interface{}
	if channelId > 0 {
		qbTotal.And(models.COLUMN_ReportGameUserDaily_ChannelId + "=?")
		param = append(param, channelId)
	}
	if uid > 0 {
		qbTotal.And(models.COLUMN_ReportGameUserDaily_Uid + "=?")
		param = append(param, uid)
	}
	qbTotal.And(models.COLUMN_ReportGameUserDaily_Ctime + "=?")
	qbTotal.And(models.COLUMN_ReportGameUserDaily_Ctime + "=?")
	qbTotal.GroupBy(models.COLUMN_ReportGameUserDaily_Uid)
	param = append(param, startTime)
	param = append(param, endTime)

	err = d.Orm.Raw(qbTotal.String(), param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error: %v", err)
		return
	}

	return*/
}

type ChannelTeamReport struct {
	Uid      uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Level    uint32 `orm:"column(level)" json:"level,omitempty"`
	ValidBet int64  `orm:"column(valid_bet)" json:"valid_bet,omitempty"`
	BetNum   int32  `orm:"column(bet_num)" json:"bet_num,omitempty"`
	Profit   int64  `orm:"column(profit)" json:"profit,omitempty"`
	Salary   int64  `orm:"column(salary)" json:"salary,omitempty"`
}

func (d *ReportGameUserDailyDao) QueryTeamDataByTime(page, limit int, channelId uint32, startTime, endTime int64, uid, pUid uint64) (channelTeamReport []ChannelTeamReport, err error) {
	var param []interface{}
	param = append(param, startTime)
	param = append(param, endTime)
	sqlQuery := fmt.Sprintf("SELECT %s,%s,SUM(%s) AS valid_bet,SUM(%s) AS bet_num,SUM(%s) AS profit,"+
		"SUM(%s) AS salary FROM %s WHERE %s>=? AND %s<=?", models.COLUMN_ReportGameUserDaily_Uid,
		models.COLUMN_ReportGameUserDaily_Level, models.COLUMN_ReportGameUserDaily_TotalValidBet,
		models.COLUMN_ReportGameUserDaily_TotalBetNum, models.COLUMN_ReportGameUserDaily_TotalProfit,
		models.COLUMN_ReportGameUserDaily_TeamSalary, models.TABLE_ReportGameUserDaily,
		models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Ctime)
	if channelId > 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportGameUserDaily_ChannelId)
		param = append(param, channelId)
	}
	if uid > 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportGameUserDaily_Uid)
		param = append(param, uid)
	}
	if pUid > 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportGameUserDaily_PUid)
		param = append(param, pUid)
	}
	if uid == 0 && pUid == 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportGameUserDaily_Level)
		param = append(param, 1)
	}
	sqlQuery = fmt.Sprintf("%s GROUP BY %s LIMIT ? OFFSET ?", sqlQuery, models.COLUMN_ReportGameUserDaily_Uid)
	param = append(param, limit)
	param = append(param, (page-1)*limit)

	_, err = d.Orm.Raw(sqlQuery, param...).QueryRows(&channelTeamReport)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *ReportGameUserDailyDao) Sum(channelId uint32, startTime, endTime int64) (channelTeamReport ChannelTeamReport, err error) {
	var param []interface{}
	param = append(param, startTime)
	param = append(param, endTime)
	sqlQuery := fmt.Sprintf("SELECT SUM(%s) AS valid_bet,SUM(%s) AS bet_num,SUM(%s) AS profit,"+
		"SUM(%s) AS salary FROM %s WHERE %s=1 AND %s>=? AND %s<=?", models.COLUMN_ReportGameUserDaily_TotalValidBet,
		models.COLUMN_ReportGameUserDaily_TotalBetNum, models.COLUMN_ReportGameUserDaily_TotalProfit,
		models.COLUMN_ReportGameUserDaily_TeamSalary, models.TABLE_ReportGameUserDaily, models.COLUMN_ReportGameUserDaily_Level,
		models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Ctime)
	if channelId > 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportGameUserDaily_ChannelId)
		param = append(param, channelId)
	}

	err = d.Orm.Raw(sqlQuery, param...).QueryRow(&channelTeamReport)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *ReportGameUserDailyDao) Personal(channelId uint32, startTime, endTime int64, uid uint64) (channelTeamReport ChannelTeamReport, err error) {
	var param []interface{}
	param = append(param, uid)
	param = append(param, startTime)
	param = append(param, endTime)
	sqlQuery := fmt.Sprintf("SELECT %s,SUM(%s) AS valid_bet,SUM(%s) AS bet_num,SUM(%s) AS profit,"+
		"SUM(%s) AS salary FROM %s WHERE %s=? AND %s>=? AND %s<=?", models.COLUMN_ReportGameUserDaily_Level,
		models.COLUMN_ReportGameUserDaily_ValidBet, models.COLUMN_ReportGameUserDaily_TotalBetNum,
		models.COLUMN_ReportGameUserDaily_Profit, models.COLUMN_ReportGameUserDaily_Salary,
		models.TABLE_ReportGameUserDaily, models.COLUMN_ReportGameUserDaily_Uid,
		models.COLUMN_ReportGameUserDaily_Ctime, models.COLUMN_ReportGameUserDaily_Ctime)
	if channelId > 0 {
		sqlQuery = fmt.Sprintf("%s AND %s=? ", sqlQuery, models.COLUMN_ReportGameUserDaily_ChannelId)
		param = append(param, channelId)
	}

	err = d.Orm.Raw(sqlQuery, param...).QueryRow(&channelTeamReport)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}
