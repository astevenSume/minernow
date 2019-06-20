package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	eusddao "utils/eusd/dao"
	eusdmodels "utils/eusd/models"
	"utils/otc/models"
)

type OtcExchangerDao struct {
	common.BaseDao
}

func NewOtcExchangerDao(db string) *OtcExchangerDao {
	return &OtcExchangerDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var OtcExchangerDaoEntity *OtcExchangerDao

type OtcExchanger struct {
	Uid      uint64 `json:"user_id"`
	Mobile   string `json:"mobile"`
	Wechat   string `json:"wechat"`
	Telegram string `json:"telegram"`
	Status   int8   `json:"status"`
	From     int8   `json:"from"`
	Ctime    int64  `json:"ctime"`
	Utime    int64  `json:"utime"`
}

type OtcExchangerAck struct {
	Uid      string `json:"user_id"`
	Mobile   string `json:"mobile"`
	Wechat   string `json:"wechat"`
	Telegram string `json:"telegram"`
	Status   int8   `json:"status"`
	From     int8   `json:"from"`
	Ctime    int64  `json:"ctime"`
	Utime    int64  `json:"utime"`
}

func (d *OtcExchangerDao) ClientExchanger(data *OtcExchanger) (ack OtcExchangerAck) {
	if data == nil {
		return
	}

	ack.Uid = fmt.Sprintf("%v", data.Uid)
	ack.Mobile = data.Mobile
	ack.Wechat = data.Wechat
	ack.Telegram = data.Telegram
	ack.Status = data.Status
	ack.From = data.From
	ack.Ctime = data.Ctime
	ack.Utime = data.Utime

	return
}

func (d *OtcExchangerDao) Create(uid uint64, from int8, mobile, wechat, telegram string) (data models.OtcExchanger, err error) {
	data = models.OtcExchanger{
		Uid:      uid,
		Mobile:   mobile,
		Wechat:   wechat,
		Telegram: telegram,
		From:     from,
		Ctime:    common.NowInt64MS(),
		Utime:    common.NowInt64MS(),
	}

	_, _, err = d.Orm.ReadOrCreate(&data, models.COLUMN_OtcExchanger_Uid)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
		return
	}

	return
}

func (d *OtcExchangerDao) QueryById(uid uint64) (otcExchanger OtcExchanger, err error) {
	var qbQuery orm.QueryBuilder
	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	fields := fmt.Sprintf("T1.%s,T1.%s,T1.%s,T1.%s,T1.%s,T1.%s,T1.%s,T2.%s",
		models.COLUMN_OtcExchanger_Uid,
		models.COLUMN_OtcExchanger_Wechat,
		models.COLUMN_OtcExchanger_Telegram,
		models.COLUMN_OtcExchanger_From,
		models.COLUMN_OtcExchanger_Utime,
		models.COLUMN_OtcExchanger_Ctime,
		models.COLUMN_OtcExchanger_Mobile,
		eusdmodels.COLUMN_EosOtc_Status)
	qbQuery.Select(fields).From("((").Select("*").From(models.TABLE_OtcExchanger).Where(models.COLUMN_OtcExchanger_Uid + "=?")
	sqlQuery := qbQuery.String()
	sqlQuery = fmt.Sprintf("%s) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s)", sqlQuery, eusdmodels.TABLE_EosOtc, models.COLUMN_OtcExchanger_Uid, eusdmodels.COLUMN_EosOtc_Uid)
	err = d.Orm.Raw(sqlQuery, uid).QueryRow(&otcExchanger)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//条件查询
func (d *OtcExchangerDao) QueryCondition(mobile string, wechat string, status int8, page int, limit int) (total int64, otcExchangers []OtcExchanger, err error) {
	var qbQuery orm.QueryBuilder
	var qbTotal orm.QueryBuilder
	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	qbTotal, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}

	var sqlTotal string
	var sqlQuery string
	var sqlSub string
	qbTotal.Select("Count(*) AS total")
	qbQuery.Select("T1."+models.COLUMN_OtcExchanger_Uid,
		"T1."+models.COLUMN_OtcExchanger_Wechat,
		"T1."+models.COLUMN_OtcExchanger_Telegram,
		"T1."+models.COLUMN_OtcExchanger_From,
		"T1."+models.COLUMN_OtcExchanger_Mobile,
		"T1."+models.COLUMN_OtcExchanger_Utime,
		"T2."+models.COLUMN_OtcExchanger_Ctime,
		"T2."+eusdmodels.COLUMN_EosOtc_Status)
	var param []interface{}

	if status == eusddao.WealthStatusWorking || status == eusddao.WealthStatusLock {
		qbTotal.From("((").Select("*").From(eusdmodels.TABLE_EosOtc).Where(eusdmodels.COLUMN_EosOtc_Status + "=?")
		qbQuery.From("((").Select("*").From(eusdmodels.TABLE_EosOtc).Where(eusdmodels.COLUMN_EosOtc_Status + "=?")
		param = append(param, status)

		sqlTotal = qbTotal.String()
		sqlQuery = qbQuery.String()
		qbQuery.Limit(limit).Offset((page - 1) * limit)

		var qbSub orm.QueryBuilder
		qbSub, err = orm.NewQueryBuilder("mysql")
		if err != nil {
			common.LogFuncError("err:%v", err)
			return
		}

		qbSub.Select("*").From(models.TABLE_OtcExchanger).Where("1=1")
		if mobile != "" {
			qbSub = qbSub.And(models.COLUMN_OtcExchanger_Mobile + " REGEXP  ?")
			param = append(param, mobile)
		}
		if wechat != "" {
			qbSub = qbSub.And(models.COLUMN_OtcExchanger_Wechat + " REGEXP ?")
			param = append(param, wechat)
		}

		sqlSub = qbSub.String()
		sqlTotal = fmt.Sprintf("%s) AS T2 LEFT JOIN (%s) AS T1 ON T1.%s=T2.%s)", sqlTotal, sqlSub, models.COLUMN_OtcExchanger_Uid, eusdmodels.COLUMN_EosOtc_Uid)
		sqlQuery = fmt.Sprintf("%s) AS T2 LEFT JOIN (%s) AS T1 ON T1.%s=T2.%s)", sqlQuery, sqlSub, models.COLUMN_OtcExchanger_Uid, eusdmodels.COLUMN_EosOtc_Uid)
	} else {
		qbTotal.From(models.TABLE_OtcExchanger).Where("1=1")
		qbQuery.From("((").Select("*").From(models.TABLE_OtcExchanger).Where("1=1")

		if mobile != "" {
			cond := models.COLUMN_OtcExchanger_Mobile + " = ?"
			qbTotal.And(cond)
			qbQuery.And(cond)
			param = append(param, mobile)
		}
		if wechat != "" {
			cond := models.COLUMN_OtcExchanger_Wechat + " = ?"
			qbTotal.And(cond)
			qbQuery.And(models.COLUMN_OtcExchanger_Wechat + " = ?")
			param = append(param, wechat)
		}

		qbQuery.Limit(limit).Offset((page - 1) * limit)
		sqlTotal = qbTotal.String()
		sqlQuery = qbQuery.String()
		sqlQuery = fmt.Sprintf("%s) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s)", sqlQuery, eusdmodels.TABLE_EosOtc, models.COLUMN_OtcExchanger_Uid, eusdmodels.COLUMN_EosOtc_Uid)
	}

	err = d.Orm.Raw(sqlTotal, param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error: %v", err)
		return
	}

	_, err = d.Orm.Raw(sqlQuery, param...).QueryRows(&otcExchangers)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//更新承兑商
func (d *OtcExchangerDao) UpdateExchanger(uid uint64, mobile string, wechat string, telegram string) (models.OtcExchanger, error) {
	otcExchanger := models.OtcExchanger{Uid: uid}
	err := d.Orm.Read(&otcExchanger, models.COLUMN_OtcExchanger_Uid)
	if err != nil {
		return otcExchanger, nil
	}
	var cols []string
	if mobile != "" {
		otcExchanger.Mobile = mobile
		cols = append(cols, models.COLUMN_OtcExchanger_Mobile)
	}
	if wechat != "" {
		otcExchanger.Wechat = wechat
		cols = append(cols, models.COLUMN_OtcExchanger_Wechat)
	}
	if telegram != "" {
		otcExchanger.Telegram = telegram
		cols = append(cols, models.COLUMN_OtcExchanger_Telegram)
	}
	if len(cols) == 0 {
		return otcExchanger, nil
	}
	otcExchanger.Utime = common.NowInt64MS()
	cols = append(cols, models.COLUMN_OtcExchanger_Utime)

	_, err = d.Orm.Update(&otcExchanger, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return otcExchanger, err
	}

	return otcExchanger, nil
}

//更新时间
func (d *OtcExchangerDao) UpdateUtime(uid uint64) {
	otcExchanger := models.OtcExchanger{Uid: uid}
	otcExchanger.Utime = common.NowInt64MS()
	_, err := d.Orm.Update(&otcExchanger, models.COLUMN_OtcExchanger_Utime)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}
	return
}
