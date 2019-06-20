package dao

import (
	"common"
	"fmt"
	"utils/admin/models"
)

type ActivityUserConfDao struct {
	common.BaseDao
}

func NewActivityUserConfDao(db string) *ActivityUserConfDao {
	return &ActivityUserConfDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var ActivityUserConfDaoEntity *ActivityUserConfDao

func (d *ActivityUserConfDao) Insert(playGameDay int32, betAmount, effectBetAmount int64) (err error) {
	sql := fmt.Sprintf("delete from %s", models.TABLE_ActivityUserConf)

	_, err = d.Orm.Raw(sql).Exec()
	if err != nil {
		common.LogFuncError("err %v ", err)
	}

	conf := &models.ActivityUserConf{
		PlayGameDay:        playGameDay,
		BetAmount:          betAmount,
		EffectiveBetAmount: effectBetAmount,
	}
	_, err = d.Orm.Insert(conf)
	if err != nil {
		common.LogFuncError("err %v ", err)
	}
	return
}

func (d *ActivityUserConfDao) GetConf() (conf models.ActivityUserConf, err error) {
	confs := make([]models.ActivityUserConf, 0)
	sql := fmt.Sprintf("select * from %s", models.TABLE_ActivityUserConf)
	_, err = d.Orm.Raw(sql).QueryRows(&confs)

	if len(confs) == 0 {
		common.LogFuncError("database has not activity user conf data")
	}
	if err != nil {
		common.LogFuncError("err %v")
	}
	conf = confs[0]
	return
}
