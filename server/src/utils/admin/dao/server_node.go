package dao

import (
	"common"
	"github.com/astaxie/beego"
	"utils/admin/models"
)

type ServerNodeDao struct {
	DeadSecs uint32
	common.BaseDao
}

func NewServerNodeDao(db string) (d *ServerNodeDao) {

	d = &ServerNodeDao{
		BaseDao: common.NewBaseDao(db),
	}

	if secs, err := beego.AppConfig.Int("ServerDeadSecs"); err != nil && secs <= 0 {
		d.DeadSecs = 30 //default 30 seconds
	} else {
		d.DeadSecs = uint32(secs)
	}

	return
}

var ServerNodeDaoEntity *ServerNodeDao

func (d *ServerNodeDao) UpdateLastPing(appName string, regionId, serverId int64, lastPing uint32) (err error) {
	node := models.ServerNode{
		AppName:  appName,
		RegionId: regionId,
		ServerId: serverId,
		LastPing: lastPing,
	}
	_, err = d.Orm.InsertOrUpdate(&node)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// easyjson:json
type ServerNodeMsg struct {
	Id       uint32 `json:"id"`
	RegionId int64  `json:"region_id"`
	ServerId int64  `json:"server_id"`
	AppName  string `json:"app_name"`
	LastPing uint32 `json:"last_ping"`
	Alive    bool   `json:"alive"`
}

func (d *ServerNodeDao) Query(appName string, regionId, serverId int64, page, perPage int) (total int64, list []ServerNodeMsg, err error) {
	qs := d.Orm.QueryTable(models.TABLE_ServerNode)
	if appName != "" {
		qs = qs.Filter(models.COLUMN_ServerNode_AppName, appName)
	}

	if regionId >= 0 && regionId < RegionIdImpossible {
		qs = qs.Filter(models.COLUMN_ServerNode_RegionId, regionId)
	}

	if serverId >= 0 && serverId < ServerIdImpossible {
		qs = qs.Filter(models.COLUMN_ServerNode_ServerId, serverId)
	}

	total, err = qs.Count()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	data := make([]models.ServerNode, 0)
	_, err = qs.Limit(perPage).Offset((page - 1) * perPage).All(&data)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	for _, n := range data {
		list = append(list, ServerNodeMsg{
			AppName:  n.AppName,
			RegionId: n.RegionId,
			ServerId: n.ServerId,
			LastPing: n.LastPing,
			Alive:    n.LastPing > common.NowUint32()-d.DeadSecs,
		})
	}

	return
}
