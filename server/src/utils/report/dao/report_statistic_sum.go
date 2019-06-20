package dao

import (
	"common"
	"fmt"
	"strings"
	"utils/report/models"
)

type ReportStatisticSumDao struct {
	common.BaseDao
}

func NewReportStatisticSumDao(db string) (reportStatisticMonth *ReportStatisticSumDao) {
	return &ReportStatisticSumDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportStatisticSumDaoEntity *ReportStatisticSumDao

func (d *ReportStatisticSumDao) InsertData(reportStatisticSum models.ReportStatisticSum) (err error) {
	_, err = d.Orm.Insert(&reportStatisticSum)
	if err != nil {
		common.LogFuncError("insert err %v", err)
	}
	return
}

func (d *ReportStatisticSumDao) QuerySumReport(channelId uint32, start, over int64) (totalDataList models.ReportStatisticSum, err error) {
	var where = ""
	if start > 0 {
		where = fmt.Sprintf("and %s>=%d", models.COLUMN_ReportStatisticSum_Ctime, start)
	}
	if over > 0 {
		newWhere := fmt.Sprintf(" and %s<%d", models.COLUMN_ReportStatisticSum_Ctime, over)
		where = strings.Join([]string{where, newWhere}, "")
	}

	sql := fmt.Sprintf("select %s,%s,sum(%s) %s, sum(%s) %s,sum(%s) %s, %s from %s where %s=? %s order by %s desc",
		models.COLUMN_ReportStatisticSum_ChannelId, models.COLUMN_ReportStatisticSum_ChannelPositiveNums,
		models.COLUMN_ReportStatisticSum_ChannelSalaryDaily, models.COLUMN_ReportStatisticSum_ChannelSalaryDaily,
		models.COLUMN_ReportStatisticSum_ChannelWithdrawEusd, models.COLUMN_ReportStatisticSum_ChannelWithdrawEusd,
		models.COLUMN_ReportStatisticSum_ChannelRechargeEusd, models.COLUMN_ReportStatisticSum_ChannelRechargeEusd,
		models.COLUMN_ReportStatisticSum_Ctime, models.TABLE_ReportStatisticSum,
		models.COLUMN_ReportStatisticSum_ChannelId, where, models.COLUMN_ReportStatisticSum_Ctime)

	err = d.Orm.Raw(sql, channelId).QueryRow(&totalDataList)
	if err != nil {
		common.LogFuncError("sql err: %v", err)
		return
	}

	return
}
