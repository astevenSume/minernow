package dao

import (
	"common"
	"fmt"
	"utils/game/models"
	gamemodels "utils/game/models"

	"github.com/astaxie/beego/orm"
)

const MaxQueryNum = 100

type GameUserDao struct {
	common.BaseDao
}

func NewGameUserDao(db string) *GameUserDao {
	return &GameUserDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var GameUserDaoEntity *GameUserDao

func (d *GameUserDao) QueryByUid(uid uint64, channelId uint32) (gameUser gamemodels.GameUser, err error) {
	gameUser = gamemodels.GameUser{
		Uid:       uid,
		ChannelId: channelId,
	}

	err = d.Orm.Read(&gameUser, gamemodels.COLUMN_GameUser_Uid, gamemodels.COLUMN_GameUser_ChannelId)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

func (d *GameUserDao) QueryByAccount(account string, channelId uint32) (gameUser gamemodels.GameUser, err error) {
	gameUser = gamemodels.GameUser{
		Account:   account,
		ChannelId: channelId,
	}

	err = d.Orm.Read(&gameUser, gamemodels.COLUMN_GameUser_ChannelId, gamemodels.COLUMN_GameUser_Account)
	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

func (d *GameUserDao) Add(gameUser gamemodels.GameUser) (err error) {
	if gameUser.Ctime == 0 {
		gameUser.Ctime = common.NowInt64MS()
	}
	_, err = d.Orm.Insert(&gameUser)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *GameUserDao) Update(gameUser gamemodels.GameUser) (err error) {
	if gameUser.Mtime == 0 {
		gameUser.Mtime = common.NowInt64MS()
	}
	_, err = d.Orm.Update(&gameUser)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

type GameUser struct {
}

func (d *GameUserDao) QueryUidByAccounts(channelId uint32, accounts []string) (mapUser orm.Params, err error) {
	if len(accounts) == 0 {
		return
	}
	var qb orm.QueryBuilder
	qb, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	qb.Select(gamemodels.COLUMN_GameUser_Uid,
		gamemodels.COLUMN_GameUser_Account).
		From(gamemodels.TABLE_GameUser).
		Where(gamemodels.COLUMN_GameUser_ChannelId + "=?").
		And(gamemodels.COLUMN_GameUser_Account).In(accounts...)

	mapUser = make(orm.Params)
	_, err = d.Orm.Raw(qb.String(), channelId).RowsToMap(&mapUser, gamemodels.COLUMN_GameUser_Account, gamemodels.COLUMN_GameUser_Uid)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (d *GameUserDao) QueryTotalByChannelId(channelId uint32) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(*) FROM %s WHERE %s=?", gamemodels.TABLE_GameUser,
		gamemodels.COLUMN_GameUser_ChannelId), channelId).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *GameUserDao) QueryByChannelId(channelId uint32, page, limit int) (list []gamemodels.GameUser, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s=? LIMIT ? OFFSET ?", gamemodels.TABLE_GameUser,
		gamemodels.COLUMN_GameUser_ChannelId), channelId, limit, (page-1)*limit).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *GameUserDao) QueryAllByUid(uid uint64) (list []gamemodels.GameUser, err error) {

	_, err = d.Orm.QueryTable(models.TABLE_GameUser).Filter(gamemodels.COLUMN_GameUser_Uid, uid).All(&list)

	if err != nil {
		if err != orm.ErrNoRows {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

//返回今天新注册用户数
func (d *GameUserDao) QueryUserByTimestamp(platformId uint32, start int64, over int64) (todayUsers []gamemodels.GameUser, err error) {
	_, err = d.Orm.QueryTable(gamemodels.TABLE_GameUser).
		Filter(gamemodels.COLUMN_GameUser_ChannelId, platformId).
		Filter(fmt.Sprintf("%s__gte", gamemodels.COLUMN_GameUser_Ctime), start).
		Filter(fmt.Sprintf("%s__lt", gamemodels.COLUMN_GameUser_Ctime), over).All(&todayUsers)
	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err : %v", err))
		return
	}
	return
}
