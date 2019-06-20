package dao

import (
	"common"
	"fmt"
	"time"
	"utils/report/models"
)

type ReportEusdDailyDao struct {
	common.BaseDao
}

func NewReportEusdDailyDao(db string) *ReportEusdDailyDao {
	return &ReportEusdDailyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportEusdDailyDaoEntity *ReportEusdDailyDao

func (d *ReportEusdDailyDao) Buy(uid uint64, amount int64) (err error) {
	curTime := time.Now().Unix()
	zeroTime := common.GetZeroTime(curTime)
	expireSecond := zeroTime + common.DaySeconds - curTime
	ok := common.RedisSetNX(fmt.Sprintf("%v_%v_%v", ReportEusd, uid, zeroTime), time.Duration(expireSecond))
	if ok {
		//插入
		reportGameTransferDaily := models.ReportEusdDaily{
			Uid:   uid,
			Buy:   amount,
			Sell:  0,
			Ctime: zeroTime,
		}
		_, err = d.Orm.Insert(&reportGameTransferDaily)
		if err != nil {
			common.LogFuncError("uid:%v,amount:%v,ReportGameTransferDaily insert err:%v", uid, amount)
			return
		}
	} else {
		//更新
		_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=%s+? WHERE %s=? AND %s=?", models.TABLE_ReportEusdDaily,
			models.COLUMN_ReportEusdDaily_Buy, models.COLUMN_ReportEusdDaily_Buy, models.COLUMN_ReportEusdDaily_Uid,
			models.COLUMN_ReportEusdDaily_Ctime), amount, uid, zeroTime).Exec()
		if err != nil {
			common.LogFuncError("uid:%v,amount:%v,ReportGameTransferDaily update err:%v", uid, amount)
			return
		}
	}
	return
}

func (d *ReportEusdDailyDao) Sell(uid uint64, amount int64) (err error) {
	curTime := time.Now().Unix()
	zeroTime := common.GetZeroTime(curTime)
	expireSecond := zeroTime + common.DaySeconds - curTime
	ok := common.RedisSetNX(fmt.Sprintf("%v_%v_%v", ReportEusd, uid, zeroTime), time.Duration(expireSecond))
	if ok {
		//插入
		reportGameTransferDaily := models.ReportEusdDaily{
			Uid:   uid,
			Buy:   0,
			Sell:  amount,
			Ctime: zeroTime,
		}
		_, err = d.Orm.Insert(&reportGameTransferDaily)
		if err != nil {
			common.LogFuncError("uid:%v,amount:%v,ReportGameTransferDaily insert err:%v", uid, amount)
			return
		}
	} else {
		//更新
		_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=%s+? WHERE %s=? AND %s=?", models.TABLE_ReportEusdDaily,
			models.COLUMN_ReportEusdDaily_Sell, models.COLUMN_ReportEusdDaily_Sell, models.COLUMN_ReportEusdDaily_Uid,
			models.COLUMN_ReportEusdDaily_Ctime), amount, uid, zeroTime).Exec()
		if err != nil {
			common.LogFuncError("uid:%v,amount:%v,ReportGameTransferDaily update err:%v", uid, amount)
			return
		}
	}
	return
}

func (d *ReportEusdDailyDao) QueryTotalByTimestamp(timestamp int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT COUNT(%s) FROM %s WHERE %s=?", models.COLUMN_ReportEusdDaily_Ctime,
		models.TABLE_ReportEusdDaily, models.COLUMN_ReportEusdDaily_Ctime), timestamp).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ReportEusdDailyDao) QueryByTimestamp(timestamp int64, page, limit int) (list []models.ReportEusdDaily, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s=? LIMIT ? OFFSET ?", models.TABLE_ReportEusdDaily,
		models.COLUMN_ReportEusdDaily_Ctime), timestamp, limit, (page-1)*limit).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
