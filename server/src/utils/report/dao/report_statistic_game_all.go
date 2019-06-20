package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strings"
	"utils/report/models"
)

type ReportStatisticGameAllDao struct {
	common.BaseDao
}

func NewReportStatisticGameAllDao(db string) (reportStatisticGameKy *ReportStatisticGameAllDao) {
	return &ReportStatisticGameAllDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ReportStatisticGameAllDaoEntity *ReportStatisticGameAllDao

func (d *ReportStatisticGameAllDao) FindByChannalId(channelId uint32, start int64) (num int64, err error) {
	qs := d.Orm.QueryTable(models.TABLE_ReportStatisticGameAll)
	qs = qs.Filter(models.COLUMN_ReportStatisticGameAll_ChannelId, channelId).Filter(models.COLUMN_ReportStatisticGameAll_Ctime, start)
	num, err = qs.Count()
	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
	}
	return
}

func (d *ReportStatisticGameAllDao) InertReport(reportList []models.ReportStatisticGameAll) (err error) {
	if len(reportList) != 0 {
		_, err = d.Orm.InsertMulti(100, reportList)
		if err != nil {
			common.LogFuncError(fmt.Sprintf("sql err : %v", err))
		}
	}

	return
}

func queryGameReportEveryMonth(d *ReportStatisticGameAllDao, channelId uint32, start, over int64) (totalList models.ReportStatisticGameAll, err error) {

	totalSql := fmt.Sprintf("select %s, from %s where %s>=? and %s<? order by %s desc limit 1",
		models.COLUMN_ReportStatisticGameAll_ChannelId, models.COLUMN_ReportStatisticGameAll_Ctime,
		models.TABLE_ReportStatisticGameAll, models.COLUMN_ReportStatisticGameAll_Ctime,
		models.COLUMN_ReportStatisticGameAll_Ctime, models.COLUMN_ReportStatisticGameAll_Ctime)
	d.Orm.Raw(totalSql, start, over).QueryRow(&totalList)
	logs.Debug("total : ", totalList)

	if channelId == GAME_CHANNEL_AG {

		md := models.MonthDividendRecord{}
		querySql := fmt.Sprintf("select %s from %s where %s>=? and %s<=?", models.COLUMN_MonthDividendRecord_ResultDividend, models.TABLE_MonthDividendRecord, models.COLUMN_MonthDividendRecord_Ctime, models.COLUMN_MonthDividendRecord_Ctime)
		err = common.DbOrms["otc_admin"].Raw(querySql, start, over).QueryRow(&md)
		if err != nil {
			common.LogFuncError("data err %v", err)
		}

		//(&totalList).ChannelRgDividend = md.ResultDividend
	}
	return
}

func (d *ReportStatisticGameAllDao) QueryGameReport(channelId uint32, start, over int64, page, limit int) (total int64, list []models.ReportStatisticGameAll, err error) {
	logs.Debug("request data : ", channelId, start, over, page, limit)
	var where = ""
	if start > 0 {
		where = fmt.Sprintf("and %s>=%d", models.COLUMN_ReportStatisticGameAll_Ctime, start)
	}
	if over > 0 {
		newWhere := fmt.Sprintf(" and %s<%d", models.COLUMN_ReportStatisticGameAll_Ctime, over)
		where = strings.Join([]string{where, newWhere}, "")
	}

	//先得到表中所有数据
	sql := fmt.Sprintf("select %s,%s,%s,%s,%s,%s,%s,%s from %s where %s=%d %s order by %s desc limit %d offset %d",
		models.COLUMN_ReportStatisticGameAll_GameId, models.COLUMN_ReportStatisticGameAll_NewerNums, models.COLUMN_ReportStatisticGameAll_Bet,
		models.COLUMN_ReportStatisticGameAll_ValidBet, models.COLUMN_ReportStatisticGameAll_Profit, models.COLUMN_ReportStatisticGameAll_Revenue,
		models.COLUMN_ReportStatisticGameAll_Ctime, models.COLUMN_ReportStatisticGameAll_Note, models.TABLE_ReportStatisticGameAll,
		models.COLUMN_ReportStatisticGameAll_ChannelId, channelId,
		where, models.COLUMN_ReportStatisticGameAll_Ctime, limit, (page-1)*limit)
	//logs.Debug("list sql : ", sql)
	_, err = d.Orm.Raw(sql).QueryRows(&list)

	if err != nil && err != orm.ErrNoRows {
		common.LogFuncError("%v", err)
		return
	}

	return
}
