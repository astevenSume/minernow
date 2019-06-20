package dao

import (
	"common"
	"fmt"
	"time"
	"utils/game/models"
)

type DO_GET_NOTICE int

const (
	NOT_GET_NOTICE DO_GET_NOTICE = iota
	GET_NOTICE
)

const (
	RISK_TYPE_WITHDRAW = iota + 1 //提现风险
	RISK_TYPE_ORDER               //订单审核，洗钱风险
)

const (
	WARN_GRADE uint8 = iota
	WARN_GRADE_ONE
	WARN_GRADE_TWO
	WARN_GRADE_THREE
)
const ONE_GRADE_ALERT_NUM int64 = 30000 * 100000000

type GameRiskAlertDao struct {
	common.BaseDao
}

func NewGameRiskAlertDao(db string) *GameRiskAlertDao {
	return &GameRiskAlertDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var GameRiskAlertDaoEntiry *GameRiskAlertDao

func (d *GameRiskAlertDao) QureyAll() (allRisk []models.GameRiskAlert, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_GameRiskAlert).All(&allRisk)
	if err != nil {
		common.LogFuncDebug("%v", err)
		return
	}
	return
}

//写入转帐风险
func (d *GameRiskAlertDao) InsertOrderRiskInfo(uid uint64, eusdNum, orderTime, alertTime int64, warnGrade, riskType uint8, orderRiskId uint64) (lastId int64, err error) {
	//避免重复插入
	cnt, err := d.Orm.QueryTable(models.TABLE_GameRiskAlert).Filter(models.COLUMN_GameRiskAlert_OrderRiskId, orderRiskId).Count()
	if err != nil {
		common.LogFuncError("sql err : %v", err)
		return
	}
	if cnt > 0 {
		return
	}

	sql := fmt.Sprintf("insert into %s(%s,%s,%s,%s,%s,%s,%s) values(%d, %d, %d, %d, %d, %d,%d)", models.TABLE_GameRiskAlert,
		models.COLUMN_GameRiskAlert_Uid, models.COLUMN_GameRiskAlert_EusdNum,
		models.COLUMN_GameRiskAlert_OrderTime, models.COLUMN_GameRiskAlert_AlertTime,
		models.COLUMN_GameRiskAlert_WarnGrade, models.COLUMN_GameRiskAlert_RiskType, models.COLUMN_GameRiskAlert_OrderRiskId,
		uid, eusdNum, orderTime, alertTime, warnGrade, riskType, orderRiskId)

	res, err := d.Orm.Raw(sql).Exec()
	nums, _ := res.RowsAffected()
	if nums > 0 {
		lastId, _ = res.LastInsertId()
	}
	if err != nil {
		common.LogFuncError("sql err : %v", err)
	}
	return
}

/**
 * 插入今天风控预警信息， 大于3w eusd预警
 */
func (d *GameRiskAlertDao) InsertRiskInfo(riskList []models.GameTransfer) (riskResult []models.GameRiskAlert, err error) {
	nowTime := time.Now().Unix()
	for _, v := range riskList {
		if v.EusdInteger > ONE_GRADE_ALERT_NUM {
			//无需重复写入
			num, err := d.Orm.QueryTable(models.TABLE_GameRiskAlert).
				Filter(models.COLUMN_GameTransfer_Uid, v.Uid).
				Filter(models.COLUMN_GameRiskAlert_OrderTime, v.Ctime).Count()
			if err == nil && num < 1 {
				GameAlertAlert := models.GameRiskAlert{
					Uid:       v.Uid,
					EusdNum:   v.EusdInteger,
					OrderTime: uint64(v.Ctime),
					AlertTime: uint64(nowTime),
					WarnGrade: WARN_GRADE_ONE,
					RiskType:  RISK_TYPE_WITHDRAW,
				}

				_, err = d.Orm.Insert(&GameAlertAlert)
				if err != nil {
					common.LogFuncError(fmt.Sprintf("sql err: %v", err))
				}
				riskResult = append(riskResult, GameAlertAlert)
			}

		}
	}
	return
}
