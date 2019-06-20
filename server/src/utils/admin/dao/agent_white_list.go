package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/admin/models"
)

type AgentWhiteListDao struct {
	common.BaseDao
}

func NewAgentWhiteListDao(db string) *AgentWhiteListDao {
	return &AgentWhiteListDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AgentWhiteListDaoEntity *AgentWhiteListDao

//获取代理白名单档位
func (d *AgentWhiteListDao) QueryById(id uint32) (agentWhiteList *models.AgentWhiteList, err error) {
	agentWhiteList = &models.AgentWhiteList{
		Id: id,
	}

	err = d.Orm.Read(agentWhiteList, models.COLUMN_AgentWhiteList_Id)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

//分页查询
func (d *AgentWhiteListDao) QueryByPage(page int, perPage int) (total int, agentWhiteList []models.AgentWhiteList, err error) {
	qs := d.Orm.QueryTable(models.TABLE_AgentWhiteList)
	if qs == nil {
		common.LogFuncError("QueryTable fail")
		return
	}

	var count int64
	count, err = qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}
	total = int(count)

	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > count {
		err = nil
		return
	}
	_, err = qs.Limit(perPage, start).All(&agentWhiteList)
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

func (d *AgentWhiteListDao) Create(commission, precision int32, name string) (agentWhiteList models.AgentWhiteList, err error) {
	agentWhiteList = models.AgentWhiteList{
		Name:       name,
		Precision:  precision,
		Commission: commission,
		Ctime:      common.NowInt64MS(),
		Utime:      common.NowInt64MS(),
	}
	agentWhiteList.Precision = CommissPrecision

	var id int64
	id, err = d.Orm.Insert(&agentWhiteList)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}
	agentWhiteList.Id = uint32(id)

	return
}

func (d *AgentWhiteListDao) UpdateById(id uint32, commission, precision int32, name string) (models.AgentWhiteList, error) {
	agentWhiteList := models.AgentWhiteList{Id: id}
	err := d.Orm.Read(&agentWhiteList, models.COLUMN_AgentWhiteList_Id)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return agentWhiteList, err
	}

	agentWhiteList.Commission = commission
	agentWhiteList.Name = name
	agentWhiteList.Utime = common.NowInt64MS()
	if precision > 0 {
		agentWhiteList.Precision = precision
	}
	_, err = d.Orm.Update(&agentWhiteList)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return agentWhiteList, err
	}

	return agentWhiteList, nil
}

func (d *AgentWhiteListDao) DelById(id uint32) error {
	agentWhiteList := &models.AgentWhiteList{
		Id: id,
	}

	_, err := d.Orm.Delete(agentWhiteList, models.COLUMN_AgentWhiteList_Id)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

func (d *AgentWhiteListDao) All() (list []models.AgentWhiteList, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_AgentWhiteList).All(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}
