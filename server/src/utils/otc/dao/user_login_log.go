package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/otc/models"
)

type UserLoginLogDao struct {
	common.BaseDao
}

func NewUserLoginLogDao(db string) *UserLoginLogDao {
	return &UserLoginLogDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var UserLoginLogDaoEntity *UserLoginLogDao

//auto_models_start
type UserLoginLog struct {
	Id        string `orm:"column(id);pk" json:"id,omitempty"`
	Userid    string `orm:"column(user_id)" json:"user_id,omitempty"`
	Mobile    string `orm:"column(mobile)" json:"mobile,omitempty"`
	UserAgent string `orm:"column(user_agent);null;size(256)" json:"user_agent,omitempty"`
	Ips       string `orm:"column(ips);null;size(256)" json:"ips,omitempty"`
	Ctime     int64  `orm:"column(ctime)" json:"ltime,omitempty"`
}

//分页查询
func (d *UserLoginLogDao) QueryByPage(uid uint64, page int, perPage int) (total int, userLoginLog []UserLoginLog, err error) {
	var qbQuery orm.QueryBuilder
	var qbTotal orm.QueryBuilder
	qbTotal, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	// 构建查询对象
	qbTotal.Select("Count(*) AS total").
		From(models.TABLE_UserLoginLog).
		Where("1 = 1")
	qbQuery.Select("T1.*",
		"T2."+models.COLUMN_User_Mobile).
		From("((").Select("*").From(models.TABLE_UserLoginLog).Where("1=1")
	var param []interface{}

	if uid > 0 {
		qbQuery.And(models.COLUMN_UserLoginLog_Userid + "= ?")
		qbTotal.And(models.COLUMN_UserLoginLog_Userid + "= ?")
		param = append(param, uid)
	}

	sqlTotal := qbTotal.String()
	err = d.Orm.Raw(sqlTotal, param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error: %v", err)
		return
	}

	qbQuery.OrderBy("-" + models.COLUMN_UserLoginLog_Id)
	qbQuery.Limit(perPage).Offset((page - 1) * perPage)
	sqlQuery := qbQuery.String()
	sqlQuery = fmt.Sprintf("%s) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s)", sqlQuery, models.TABLE_User,
		models.COLUMN_UserLoginLog_Userid, models.COLUMN_User_Uid)
	_, err = d.Orm.Raw(sqlQuery, param...).QueryRows(&userLoginLog)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//创建
func (d *UserLoginLogDao) Create(userId uint64, agent string, ip string) (err error) {
	userLoginLog := &models.UserLoginLog{
		Userid:    userId,
		UserAgent: agent,
		Ips:       ip,
		Ctime:     common.NowInt64MS(),
	}

	_, err = d.Orm.Insert(userLoginLog)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

//统计登录人数
func (d *UserLoginLogDao) CountByTime(start, end int64) (c uint32) {
	sql := fmt.Sprintf("select count(distinct(%s)) as c from %s where %s>=? and %s<? Group by %s",
		models.COLUMN_UserLoginLog_Userid, models.TABLE_UserLoginLog,
		models.COLUMN_UserLoginLog_Ctime, models.COLUMN_UserLoginLog_Ctime, models.COLUMN_UserLoginLog_Userid)

	type count struct {
		C uint32 `json:"c"`
	}
	res := &count{}
	if err := d.Orm.Raw(sql, start, end).QueryRow(res); err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}
	c = res.C
	return
}
