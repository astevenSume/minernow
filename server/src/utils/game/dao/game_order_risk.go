package dao

import (
	"common"
	"github.com/astaxie/beego/logs"
	"utils/game/models"
	otcmodel "utils/otc/models"
)

type GameOrderRiskDao struct {
	common.BaseDao
}

func NewGameOrderRiskDao(db string) *GameOrderRiskDao {
	return &GameOrderRiskDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var GameOrderRiskDaoEntity *GameOrderRiskDao

func (d *GameOrderRiskDao) QueryInfoByRange(start, over int64) (info []models.GameOrderRisk) {
	d.Orm.QueryTable(models.TABLE_GameOrderRisk).All(&info)
	logs.Debug("info : ", info)
	return
}

func (d *GameOrderRiskDao) InsertRiskItem(items []otcmodel.OtcOrder, riskAleatId int64, ctime int64) (err error) {
	for _, o := range items {
		gameOrderDisk := models.GameOrderRisk{
			AlertId:    riskAleatId,
			Uid:        o.Uid,
			Amount:     o.Amount,
			PayAccount: o.PayAccount,
			PayType:    o.PayType,
			OrderTime:  o.Ctime,
			Ctime:      ctime,
		}
		InsertSet, err := d.Orm.QueryTable(models.TABLE_GameOrderRisk).PrepareInsert()
		if err != nil {
			common.LogFuncError("sql err is : %v", err)
		}

		_, err = InsertSet.Insert(&gameOrderDisk)
		if err != nil {
			common.LogFuncError("sql err is : %v", err)
		}
	}

	return
}
