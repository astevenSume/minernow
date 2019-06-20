package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"sort"
	"utils/admin/models"
	//models2 "utils/report/models"
)

const MonthPDividendPrecision = 10000
const PositionMaxNum = 40 // 每个等级档位的最大个数
type MonthDividendPositionConfDao struct {
	common.BaseDao
}

func NewMonthDividendPositionConfDao(db string) *MonthDividendPositionConfDao {
	return &MonthDividendPositionConfDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var MonthDividendPositionConfDaoEntity *MonthDividendPositionConfDao

//添加月分红配置
func (d *MonthDividendPositionConfDao) EditMonthDividendConf(monthDividendCfg MonthDividendCfgs) error {
	lenCfg := len(monthDividendCfg)
	if lenCfg == 0 {
		common.LogFuncError("mysql error:%v", monthDividendCfg)
		return ErrParam
	}
	agentLv := monthDividendCfg[0].AgentLv
	mysql := fmt.Sprintf("delete from %s where %s=?", models.TABLE_MonthDividendPositionConf, models.COLUMN_MonthDividendPositionConf_AgentLv)
	_, err := d.Orm.Raw(mysql, agentLv).Exec()
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	var positionConf []models.MonthDividendPositionConf
	sort.Sort(monthDividendCfg)
	perPositionNum := PositionMaxNum * int(agentLv-1)
	for i := 0; i < lenCfg-1; i++ {
		id := int64(perPositionNum + i + 1)
		if (monthDividendCfg[i].Min < monthDividendCfg[i].Max) && (monthDividendCfg[i].Max <= monthDividendCfg[i+1].Min) {
			cfg := models.MonthDividendPositionConf{
				Id:            id,
				AgentLv:       monthDividendCfg[i].AgentLv,
				Position:      monthDividendCfg[i].Position,
				Min:           monthDividendCfg[i].Min,
				Max:           monthDividendCfg[i].Max,
				ActivityNum:   monthDividendCfg[i].ActivityNum,
				DividendRatio: monthDividendCfg[i].DividendRatio,
				Ctime:         common.NowInt64MS(),
				Utime:         common.NowInt64MS(),
			}
			positionConf = append(positionConf, cfg)
		} else {
			common.LogFuncError("data error:%v data %v", monthDividendCfg)
			return ErrParam
		}
	}
	lastIndex := lenCfg - 1

	cfg := models.MonthDividendPositionConf{
		Id:            int64(perPositionNum + lenCfg),
		AgentLv:       monthDividendCfg[lastIndex].AgentLv,
		Position:      monthDividendCfg[lastIndex].Position,
		Min:           monthDividendCfg[lastIndex].Min,
		Max:           monthDividendCfg[lastIndex].Max,
		ActivityNum:   monthDividendCfg[lastIndex].ActivityNum,
		DividendRatio: monthDividendCfg[lastIndex].DividendRatio,
		Ctime:         common.NowInt64MS(),
		Utime:         common.NowInt64MS(),
	}
	positionConf = append(positionConf, cfg)

	_, err = d.Orm.InsertMulti(InsertMulCount, positionConf)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return nil
}

//查询月分红表
func (d *MonthDividendPositionConfDao) QueryPageDividendConf(conf *models.MonthDividendPositionConf, cols ...string) error {
	if conf == nil {
		return ErrParam
	}

	err := d.Orm.Read(conf, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	conf.Min = conf.Min % MonthPDividendPrecision
	conf.Max = conf.Max % MonthPDividendPrecision

	return nil
}

//分页查询月分红表
func (d *MonthDividendPositionConfDao) QueryPageDividendConfs() (respData map[int][]models.MonthDividendPositionConf, err error) {
	data := make([]models.MonthDividendPositionConf, 0)
	_, err = d.Orm.QueryTable(models.TABLE_MonthDividendPositionConf).All(&data)
	if err != nil {
		if err != orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql_err:%v", err)
		return
	}
	respData = make(map[int][]models.MonthDividendPositionConf, 0)
	for i := 0; i < len(data); i++ {
		lv := int(data[i].AgentLv)
		if _, ok := respData[lv]; !ok {
			respData[lv] = make([]models.MonthDividendPositionConf, 0)
		}
		//data[i].Max = data[i].Max / MonthPDividendPrecision
		//data[i].Min = data[i].Min / MonthPDividendPrecision

		respData[lv] = append(respData[lv], data[i])
	}

	return
}

type MonthDividendPosition struct {
	Position      int32 `json:"position"`
	DividendRatio int32 `json:"dividend_ratio"`
}

func (d *MonthDividendPositionConfDao) QueryDividendPosition() (respData map[int][]MonthDividendPosition, err error) {
	data := make([]models.MonthDividendPositionConf, 0)
	_, err = d.Orm.QueryTable(models.TABLE_MonthDividendPositionConf).All(&data)
	if err != nil {
		if err != orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql_err:%v", err)
		return
	}
	respData = make(map[int][]MonthDividendPosition, 0)
	for i := 0; i < len(data); i++ {
		lv := int(data[i].AgentLv)
		if _, ok := respData[lv]; !ok {
			respData[lv] = make([]MonthDividendPosition, 0)
		}
		position := MonthDividendPosition{
			Position:      data[i].Position,
			DividendRatio: data[i].DividendRatio,
		}
		respData[lv] = append(respData[lv], position)
	}

	return
}

func (d *MonthDividendPositionConfDao) FindByAgentLv(agentLv int32) (respData []models.MonthDividendPositionConf, err error) {
	sql := fmt.Sprintf("select * from %s where %s=?", models.TABLE_MonthDividendPositionConf, models.COLUMN_MonthDividendPositionConf_AgentLv)
	_, err = d.Orm.Raw(sql, agentLv).QueryRows(&respData)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}
	return
}

func (d *MonthDividendPositionConfDao) InsertForTest(cfgs []*models.MonthDividendPositionConf) (err error) {
	_, err = d.Orm.InsertMulti(100, cfgs)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
	}
	return
}

func (d *MonthDividendPositionConfDao) DeleteByLevel(lv int32) (err error) {
	querySql := fmt.Sprintf("delete from %s where %s=?", models.TABLE_MonthDividendPositionConf, models.COLUMN_MonthDividendPositionConf_AgentLv)
	_, err = d.Orm.Raw(querySql, lv).Exec()
	if err != nil {
		common.LogFuncError("err %v", err)
	}
	return
}

func (d *MonthDividendPositionConfDao) CountLowLevelConfNum(lv int32) (rows int, err error) {
	querySql := fmt.Sprintf("select count(*) from %s where %s>?", models.TABLE_MonthDividendPositionConf, models.COLUMN_MonthDividendPositionConf_AgentLv)
	err = d.Orm.Raw(querySql, lv).QueryRow(&rows)

	if err != nil {
		common.LogFuncError("err %v", err)
	}
	return
}

func (d *MonthDividendPositionConfDao) DeleteAllData() (err error) {
	sql := fmt.Sprintf("delete from %s", models.TABLE_MonthDividendPositionConf)
	_, err = d.Orm.Raw(sql).Exec()
	if err != nil {
		common.LogFuncError("err %v", err)
	}
	return
}
