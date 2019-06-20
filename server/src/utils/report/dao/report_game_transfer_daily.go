package dao

import (
	"common"
	"fmt"
	"time"
	"utils/report/models"
)

type ReportGameTransferDailyDao struct {
	common.BaseDao
}

func NewReportGameTransferDailyDao(db string) *ReportGameTransferDailyDao {
	return &ReportGameTransferDailyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportGameTransferDailyDaoEntity *ReportGameTransferDailyDao

/*func (d *ReportGameTransferDailyDao) InsertMul(timestamp int64, reportGameTransferDaily []models.ReportGameTransferDaily) (err error) {
	if len(reportGameTransferDaily) == 0 {
		return
	}
	common.LogFuncDebug("ReportGameTransferDailyDao InsertMul")
	_, err = d.Orm.InsertMulti(InsertMulCount, reportGameTransferDaily)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		d.DelByTimestamp(timestamp)
		return err
	}

	return
}

func (d *ReportGameTransferDailyDao) DelByTimestamp(timestamp int64) (err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("DELETE FROM %s WHERE %s = ? ", models.TABLE_ReportGameTransferDaily,
		models.COLUMN_ReportGameTransferDaily_Ctime), timestamp).Exec()
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	return
}*/

func (d *ReportGameTransferDailyDao) Recharge(channelID uint32, uid uint64, amount int64) (err error) {
	curTime := time.Now().Unix()
	zeroTime := common.GetZeroTime(curTime)
	expireSecond := zeroTime + common.DaySeconds - curTime
	ok := common.RedisSetNX(fmt.Sprintf("%v_%v_%v", ReportGameTransfer, uid, zeroTime), time.Duration(expireSecond))
	if ok {
		//插入
		reportGameTransferDaily := models.ReportGameTransferDaily{
			Uid:       uid,
			ChannelId: channelID,
			Recharge:  amount,
			Withdraw:  0,
			Ctime:     zeroTime,
		}
		_, err = d.Orm.Insert(&reportGameTransferDaily)
		if err != nil {
			common.LogFuncError("channelID:%v,uid:%v,amount:%v,ReportGameTransferDaily insert err:%v", channelID, uid, amount)
			return
		}
	} else {
		//更新
		_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=%s+? WHERE %s=? AND %s=? AND %s=?", models.TABLE_ReportGameTransferDaily,
			models.COLUMN_ReportGameTransferDaily_Recharge, models.COLUMN_ReportGameTransferDaily_Recharge,
			models.COLUMN_ReportGameTransferDaily_Uid, models.COLUMN_ReportGameTransferDaily_ChannelId,
			models.COLUMN_ReportGameTransferDaily_Ctime), amount, uid, channelID, zeroTime).Exec()
		if err != nil {
			common.LogFuncError("channelID:%v,uid:%v,amount:%v,ReportGameTransferDaily update err:%v", channelID, uid, amount)
			return
		}
	}
	return
}

func (d *ReportGameTransferDailyDao) Withdraw(channelID uint32, uid uint64, amount int64) (err error) {
	curTime := time.Now().Unix()
	zeroTime := common.GetZeroTime(curTime)
	expireSecond := zeroTime + common.DaySeconds - curTime
	ok := common.RedisSetNX(fmt.Sprintf("%v_%v_%v", ReportGameTransfer, uid, zeroTime), time.Duration(expireSecond))
	if ok {
		//插入
		reportGameTransferDaily := models.ReportGameTransferDaily{
			Uid:       uid,
			ChannelId: channelID,
			Recharge:  0,
			Withdraw:  amount,
			Ctime:     zeroTime,
		}
		_, err = d.Orm.Insert(&reportGameTransferDaily)
		if err != nil {
			common.LogFuncError("channelID:%v,uid:%v,amount:%v,ReportGameTransferDaily insert err:%v", channelID, uid, amount)
			return
		}
	} else {
		//更新
		_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=%s+? WHERE %s=? AND %s=? AND %s=?", models.TABLE_ReportGameTransferDaily,
			models.COLUMN_ReportGameTransferDaily_Withdraw, models.COLUMN_ReportGameTransferDaily_Withdraw,
			models.COLUMN_ReportGameTransferDaily_Uid, models.COLUMN_ReportGameTransferDaily_ChannelId,
			models.COLUMN_ReportGameTransferDaily_Ctime), amount, uid, channelID, zeroTime).Exec()
		if err != nil {
			common.LogFuncError("channelID:%v,uid:%v,amount:%v,ReportGameTransferDaily update err:%v", channelID, uid, amount)
			return
		}
	}

	return
}

func (d *ReportGameTransferDailyDao) QueryTotalByTimestamp(timestamp int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT COUNT(%s) FROM %s WHERE %s=?",
		models.COLUMN_ReportGameTransferDaily_Uid, models.TABLE_ReportGameTransferDaily,
		models.COLUMN_ReportGameTransferDaily_Ctime), timestamp).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *ReportGameTransferDailyDao) QueryByTimestamp(timestamp int64, page, limit int) (list []models.ReportGameTransferDaily, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s=? LIMIT ? OFFSET ?",
		models.TABLE_ReportGameTransferDaily, models.COLUMN_ReportGameTransferDaily_Ctime), timestamp, limit,
		(page-1)*limit).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
