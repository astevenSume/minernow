package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/report/models"
)

type ProfitReportDailyDao struct {
	common.BaseDao
}

func NewProfitReportDailyDao(db string) *ProfitReportDailyDao {
	return &ProfitReportDailyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ProfitReportDailyDaoEntity *ProfitReportDailyDao

func (d *ProfitReportDailyDao) InsertMulti(reports []*models.ProfitReportDaily) (err error) {
	_, err = d.Orm.InsertMulti(100, &reports)
	if err != nil {
		common.LogFuncError("db err %v", err)
	}
	return
}
func (d *ProfitReportDailyDao) QueryTotalByDateUid(uid uint64, startDay int64, endDay int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(*) FROM %s WHERE %s=?  AND %s>=? AND %s<=?",
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_Uid,
		models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime), uid, startDay, endDay).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ProfitReportDailyDao) QueryTotalByDateUidChannel(uid uint64, channelId uint32, startDay int64, endDay int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(*) FROM %s WHERE %s=? AND %s=? AND %s>=? AND %s<=?",
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_ChannelId,
		models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime), uid, channelId, startDay, endDay).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ProfitReportDailyDao) FindByDataRangeAndUid(uid uint64, startDay int64, endDay int64, page, limit int) (reports []*models.ProfitReportDaily, err error) {
	querySql := fmt.Sprintf("select * from %s WHERE %s=?  AND %s>=? AND %s<=?  ORDER BY %s DESC LIMIT ?,?",
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_Ctime,
		models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime)
	_, err = d.Orm.Raw(querySql, uid, startDay, endDay, (page-1)*limit, limit).QueryRows(&reports)
	if err != nil {
		common.LogFuncError("db err %v", err)
	}
	return
}

func (d *ProfitReportDailyDao) FindByDataRangeAndUidChannel(uid uint64, channelId uint32, startDay int64, endDay int64, page, limit int) (reports []*models.ProfitReportDaily, err error) {
	querySql := fmt.Sprintf("select * from %s WHERE %s=? AND %s=? AND %s>=? AND %s<=?  ORDER BY %s DESC LIMIT ?,?",
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_ChannelId, models.COLUMN_ProfitReportDaily_Ctime,
		models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime)
	_, err = d.Orm.Raw(querySql, uid, channelId, startDay, endDay, (page-1)*limit, limit).QueryRows(&reports)
	if err != nil {
		common.LogFuncError("db err %v", err)
	}
	return
}

type ProfitReportStat struct {
	Uid                uint64 `json:"uid"`
	Bet                int64  `json:"bet"`
	TotalValidBet      int64  `json:"total_valid_bet"`
	Profit             int64  `json:"profit"`
	Salary             int64  `json:"salary"`
	SelfDividend       int64  `json:"self_dividend"`
	AgentDividend      int64  `json:"agent_dividend"`
	ResultDividend     int64  `json:"result_dividend"`
	WithdrawAmount     uint32 `json:"withdraw_amount"`
	GameRechargeAmount uint32 `json:"game_recharge_amount"`
}

//统计这个用户这个时间段的某个渠道总和
func (d *ProfitReportDailyDao) StatProfitReportByUidChannelId(uid uint64, channelId uint32, startTime, endTime int64) (stat *models.ProfitReportDaily, err error) {
	stat = new(models.ProfitReportDaily)
	querySql := fmt.Sprintf("SELECT %s,SUM(%s) as bet,SUM(%s) as total_valid_bet,SUM(%s) as profit,"+
		"SUM(%s) as salary,SUM(%s) as self_dividend,SUM(%s) as agent_dividend,SUM(%s) as result_dividend,"+
		"SUM(%s) as withdraw_amount,SUM(%s) as game_recharge_amount FROM %s WHERE %s=? AND %s=? AND %s>=? AND %s<=? group by %s",
		models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_Bet, models.COLUMN_ProfitReportDaily_TotalValidBet,
		models.COLUMN_ProfitReportDaily_Profit, models.COLUMN_ProfitReportDaily_Salary, models.COLUMN_ProfitReportDaily_SelfDividend,
		models.COLUMN_ProfitReportDaily_AgentDividend, models.COLUMN_ProfitReportDaily_ResultDividend, models.COLUMN_ProfitReportDaily_GameWithdrawAmount, models.COLUMN_ProfitReportDaily_GameRechargeAmount,
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_ChannelId,
		models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Uid)
	err = d.Orm.Raw(querySql, uid, channelId, startTime, endTime).QueryRow(&stat)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("db error %v", err)
		return
	}
	return
}

//分页统计所有用户这个时间段某个渠道的总和
func (d *ProfitReportDailyDao) StatProfitReportByChannelId(channelId uint32, startTime, endTime int64, page, limit int) (stat []*models.ProfitReportDaily, err error) {
	querySql := fmt.Sprintf("SELECT %s,SUM(%s) as bet,SUM(%s) as total_valid_bet,SUM(%s) as profit,"+
		"SUM(%s) as salary,SUM(%s) as self_dividend,SUM(%s) as agent_dividend,SUM(%s) as result_dividend,"+
		"SUM(%s) as withdraw_amount,SUM(%s) as game_recharge_amount FROM %s WHERE %s=? AND %s>=? AND %s<=? group by %s ORDER BY %s DESC LIMIT ?,?",
		models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_Bet, models.COLUMN_ProfitReportDaily_TotalValidBet,
		models.COLUMN_ProfitReportDaily_Profit, models.COLUMN_ProfitReportDaily_Salary, models.COLUMN_ProfitReportDaily_SelfDividend,
		models.COLUMN_ProfitReportDaily_AgentDividend, models.COLUMN_ProfitReportDaily_ResultDividend, models.COLUMN_ProfitReportDaily_GameWithdrawAmount, models.COLUMN_ProfitReportDaily_GameRechargeAmount,
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_ChannelId, models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_Uid)
	_, err = d.Orm.Raw(querySql, channelId, startTime, endTime, (page-1)*limit, limit).QueryRows(&stat)
	if err != nil {
		common.LogFuncError("db error %v", err)
		return
	}
	fmt.Println("stat ", stat)
	return
}

//统计这个用户这个时间段的所有渠道总和
func (d *ProfitReportDailyDao) StatProfitReport(uid uint64, startTime, endTime int64) (stat *ProfitReportStat, err error) {
	querySql := fmt.Sprintf("SELECT %s,SUM(%s) as bet,SUM(%s) as total_valid_bet,SUM(%s) as profit,"+
		"SUM(%s) as salary,SUM(%s) as self_dividend,SUM(%s) as agent_dividend,SUM(%s) as result_dividend,"+
		"SUM(%s) as withdraw_amount,SUM(%s) as game_recharge_amount FROM %s WHERE %s=? AND %s>=? AND %s<=? group by %s",
		models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_Bet, models.COLUMN_ProfitReportDaily_TotalValidBet,
		models.COLUMN_ProfitReportDaily_Profit, models.COLUMN_ProfitReportDaily_Salary, models.COLUMN_ProfitReportDaily_SelfDividend,
		models.COLUMN_ProfitReportDaily_AgentDividend, models.COLUMN_ProfitReportDaily_ResultDividend, models.COLUMN_ProfitReportDaily_GameWithdrawAmount, models.COLUMN_ProfitReportDaily_GameRechargeAmount,
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Uid)
	err = d.Orm.Raw(querySql, uid, startTime, endTime).QueryRow(&stat)
	if err != nil {
		common.LogFuncError("db error %v", err)
		return
	}
	fmt.Println("stat ", stat)
	return
}

//统计这个用户同一天的所有渠道总和
func (d *ProfitReportDailyDao) StatProfitReportInOneDay(uid uint64, cTime int64) (stat *models.ProfitReportDaily, err error) {
	querySql := fmt.Sprintf("SELECT %s,SUM(%s) as bet,SUM(%s) as total_valid_bet,SUM(%s) as profit,"+
		"SUM(%s) as salary,SUM(%s) as self_dividend,SUM(%s) as agent_dividend,SUM(%s) as result_dividend,"+
		"SUM(%s) as game_withdraw_amount,SUM(%s) as game_recharge_amount FROM %s WHERE %s=? AND %s=? group by %s",
		models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_Bet, models.COLUMN_ProfitReportDaily_TotalValidBet,
		models.COLUMN_ProfitReportDaily_Profit, models.COLUMN_ProfitReportDaily_Salary, models.COLUMN_ProfitReportDaily_SelfDividend,
		models.COLUMN_ProfitReportDaily_AgentDividend, models.COLUMN_ProfitReportDaily_ResultDividend, models.COLUMN_ProfitReportDaily_GameWithdrawAmount, models.COLUMN_ProfitReportDaily_GameRechargeAmount,
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Uid)
	err = d.Orm.Raw(querySql, uid, cTime).QueryRow(&stat)
	if err != nil {
		common.LogFuncError("db error %v", err)
		return
	}
	return
}

func (d *ProfitReportDailyDao) QueryTotalByDateChannel(channelId uint32, startDay int64, endDay int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(*) FROM %s WHERE %s=? AND %s>=? AND %s<=?",
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_ChannelId,
		models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime), channelId, startDay, endDay).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ProfitReportDailyDao) FindByDataRange(startDay int64, endDay int64, page, limit int) (reports []*models.ProfitReportDaily, err error) {
	querySql := fmt.Sprintf("select * from %s WHERE  %s>=? AND %s<=? ORDER BY %s DESC LIMIT ?,?",
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_Ctime,
		models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime)
	_, err = d.Orm.Raw(querySql, startDay, endDay, (page-1)*limit, limit).QueryRows(&reports)
	if err != nil {
		common.LogFuncError("db err %v", err)
	}
	return
}
func (d *ProfitReportDailyDao) FindByDataRangeAndChannel(channelId uint32, startDay int64, endDay int64, page, limit int) (reports []*models.ProfitReportDaily, err error) {
	querySql := fmt.Sprintf("select * from %s WHERE  %s=? AND %s>=? AND %s<=? ORDER BY %s DESC LIMIT ?,?",
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_ChannelId, models.COLUMN_ProfitReportDaily_Ctime,
		models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime)
	_, err = d.Orm.Raw(querySql, channelId, startDay, endDay, (page-1)*limit, limit).QueryRows(&reports)
	if err != nil {
		common.LogFuncError("db err %v", err)
	}
	return
}

func (d *ProfitReportDailyDao) QueryTotalByDateAndChannelId(channelId uint32, startDay int64, endDay int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(DISTINCT %s) FROM %s WHERE %s=? AND %s>=? AND %s<=?", models.COLUMN_ProfitReportDaily_Uid,
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_ChannelId, models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime), channelId, startDay, endDay).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
func (d *ProfitReportDailyDao) QueryTotalByDateAndUid(uid uint64, startDay int64, endDay int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(*) FROM %s WHERE %s=? AND %s>=? AND %s<=?",
		models.TABLE_ProfitReportDaily, models.COLUMN_ProfitReportDaily_Uid, models.COLUMN_ProfitReportDaily_Ctime,
		models.COLUMN_ProfitReportDaily_Ctime), uid, startDay, endDay).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
func (d *ProfitReportDailyDao) FindByDayRange(startDay int64, endDay int64, page, limit int) (reports []*models.ProfitReportDaily, err error) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s>=? AND %s<=? ORDER BY %s DESC LIMIT ?,?", models.TABLE_ProfitReportDaily,
		models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime, models.COLUMN_ProfitReportDaily_Ctime)
	_, err = d.Orm.Raw(sql, startDay, endDay, (page-1)*limit, limit).QueryRows(&reports)
	if err != nil {
		common.LogFuncError("err %v", err)
	}
	return
}
