package dao

import (
	"common"
	"fmt"
	"time"
	"utils/report/models"
)

type ReportAgentDailyDao struct {
	common.BaseDao
}

func NewReportAgentDailyDao(db string) *ReportAgentDailyDao {
	return &ReportAgentDailyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportAgentDailyDaoEntity *ReportAgentDailyDao

/*
func (d *ReportAgentDailyDao) InsertMul(timestamp int64, reportAgentDaily []models.ReportAgentDaily) (err error) {
	if len(reportAgentDaily) == 0 {
		return
	}
	common.LogFuncDebug("ReportAgentDailyDao InsertMul")
	_, err = d.Orm.InsertMulti(InsertMulCount, reportAgentDaily)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		d.DelByTimestamp(timestamp)
		return err
	}

	return
}

func (d *ReportAgentDailyDao) DelByTimestamp(timestamp int64) (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("DELETE FROM %s WHERE %s = ? ", models.TABLE_ReportAgentDaily,
		models.COLUMN_ReportAgentDaily_Ctime), timestamp).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	return
}*/

func (d *ReportAgentDailyDao) AddWithDraw(uid uint64, amount int64) (err error) {
	curTime := time.Now().Unix()
	zeroTime := common.GetZeroTime(curTime)
	expireSecond := zeroTime + common.DaySeconds - curTime
	ok := common.RedisSetNX(fmt.Sprintf("%v_%v_%v", ReportAgent, uid, zeroTime), time.Duration(expireSecond))
	if ok {
		//插入
		reportGameTransferDaily := models.ReportAgentDaily{
			Uid:         uid,
			SumWithdraw: amount,
			Ctime:       zeroTime,
		}
		_, err = d.Orm.Insert(&reportGameTransferDaily)
		if err != nil {
			common.LogFuncError("uid:%v,amount:%v,ReportAgentDaily insert err:%v", uid, amount)
			return
		}
	} else {
		//更新
		_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=%s+? WHERE %s=? AND %s=?", models.TABLE_ReportAgentDaily,
			models.COLUMN_ReportAgentDaily_SumWithdraw, models.COLUMN_ReportAgentDaily_SumWithdraw,
			models.COLUMN_ReportAgentDaily_Uid, models.COLUMN_ReportAgentDaily_Ctime), amount, uid, zeroTime).Exec()
		if err != nil {
			common.LogFuncError("uid:%v,amount:%v,ReportAgentDaily update err:%v", uid, amount)
			return
		}
	}
	return
}
