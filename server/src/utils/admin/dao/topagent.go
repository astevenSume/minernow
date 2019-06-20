package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/admin/models"
)

//代理状态
const (
	TopAgentStatusNil = iota
	TopAgentStatusUnRegister
	TopAgentStatusRegistered
	TopAgentStatusMax
)

type TopAgentDao struct {
	common.BaseDao
}

func NewTopAgentDao(db string) *TopAgentDao {
	return &TopAgentDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var TopAgentDaoEntity *TopAgentDao

//获取一级代理
func (d *TopAgentDao) QueryById(id uint32) (topAgent *models.TopAgent, err error) {
	topAgent = &models.TopAgent{Id: id}
	err = d.Orm.Read(topAgent, models.COLUMN_TopAgent_Id)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

//分页条件查询
func (d *TopAgentDao) QueryPageCondition(mobile string, status int8, page int, perPage int) (total int64, topAgents []models.TopAgent, err error) {
	qs := d.Orm.QueryTable(models.TABLE_TopAgent)
	if qs == nil {
		common.LogFuncError("QueryTable fail")
		return
	}

	if status > TopAgentStatusNil && status < TopAgentStatusMax {
		qs = qs.Filter(models.COLUMN_TopAgent_Status, status)
	}
	if mobile != "" {
		qs = qs.Filter(models.COLUMN_TopAgent_Mobile+"__contains", mobile)
	}

	total, err = qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > total {
		err = nil
		return
	}
	_, err = qs.Limit(perPage, start).All(&topAgents)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *TopAgentDao) Create(nationalCode, mobile string) (isNew bool, topAgent models.TopAgent, err error) {
	topAgent = models.TopAgent{
		NationalCode: nationalCode,
		Mobile:       mobile,
		Status:       TopAgentStatusUnRegister,
		Ctime:        common.NowInt64MS(),
		Utime:        common.NowInt64MS(),
	}

	var id int64
	isNew, id, err = d.Orm.ReadOrCreate(&topAgent, models.COLUMN_TopAgent_Mobile)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}
	topAgent.Id = uint32(id)

	return
}

func (d *TopAgentDao) Update(id uint32, nationalCode, mobile string) (ok bool, topAgent *models.TopAgent, err error) {
	topAgent = &models.TopAgent{Mobile: mobile}
	err = d.Orm.Read(topAgent, models.COLUMN_TopAgent_Mobile)
	if err != nil {
		if err != orm.ErrNoRows {
			return
		}
		err = nil
	}
	if topAgent.Id > 0 {
		//号码已在列表
		return
	}

	topAgent.Id = id
	err = d.Orm.Read(topAgent, models.COLUMN_TopAgent_Id)
	if err != nil {
		return
	}
	if topAgent.Status == TopAgentStatusRegistered {
		//已注册 不能修改
		return
	}

	topAgent.Id = id
	topAgent.NationalCode = nationalCode
	topAgent.Mobile = mobile
	topAgent.Utime = common.NowInt64MS()
	_, err = d.Orm.Update(topAgent, models.COLUMN_TopAgent_NationalCode, models.COLUMN_TopAgent_Mobile, models.COLUMN_TopAgent_Utime)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}
	ok = true

	return
}

func (d *TopAgentDao) DelById(id uint32) (bool, error) {
	topAgent := &models.TopAgent{Id: id}

	err := d.Orm.Read(topAgent, models.COLUMN_TopAgent_Id)
	if err != nil {
		return false, err
	}
	if topAgent.Status == TopAgentStatusRegistered {
		return false, nil
	}

	_, err = d.Orm.Delete(topAgent, models.COLUMN_TopAgent_Id)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return false, err
	}

	return true, nil
}

func (d *TopAgentDao) IsTopAgent(nationalCode, mobile string) bool {
	topAgent := &models.TopAgent{Mobile: mobile}
	err := d.Orm.Read(topAgent, models.COLUMN_TopAgent_Mobile)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return false
	}

	err = d.Register(nationalCode, mobile)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
	}

	return true
}

func (d *TopAgentDao) Register(nationalCode, mobile string) error {
	sql := fmt.Sprintf("Update %s Set %s=%v where %s=?", models.TABLE_TopAgent, models.COLUMN_TopAgent_Status, TopAgentStatusRegistered, models.COLUMN_TopAgent_Mobile)
	_, err := d.Orm.Raw(sql, mobile).Exec()
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return nil
}
