package dao

import (
	"common"
	"fmt"
	"utils/tmp/models"
)

const (
	NOTBET = iota
	HASBET
)

type TmpGamebetersDao struct {
	common.BaseDao
}

func NewTmpGameBetersDao(db string) *TmpGamebetersDao {
	return &TmpGamebetersDao{
		common.NewBaseDao(db),
	}
}

var TmpGameBetersDaoEntity *TmpGamebetersDao

func (d *TmpGamebetersDao) InsertUser(uid uint64, channel_id, game_id uint32, isBet uint32, ctime int64) (err error) {

	if _, err = d.Orm.Insert(&models.TmpGamebeters{Uid: uid, ChannelId: channel_id, GameId: game_id, Ctime: ctime, BetType: isBet}); err != nil {
		common.LogFuncError(fmt.Sprintf("sql err is : %v", err))
	}
	return
}
func (d *TmpGamebetersDao) UpdateUser(uid uint64, gameId uint32, channel_id uint32, betType uint32, ctime int64) (err error) {
	if _, err := d.Orm.Update(&models.TmpGamebeters{Uid: uid, ChannelId: channel_id, Ctime: ctime, BetType: betType}); err != nil {
		common.LogFuncError(fmt.Sprintf("sql err is : %v", err))
	}
	return
}

func (d *TmpGamebetersDao) QueryUnbetUser(channelId uint32) (allUsers []models.TmpGamebeters, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_TmpGamebeters).Filter(models.COLUMN_TmpGamebeters_ChannelId, channelId).Filter(models.COLUMN_TmpGamebeters_BetType, NOTBET).All(&allUsers)

	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err is : %v", err))
	}
	return
}

type NewersList struct {
	Nums   uint32 `orm:"column(nums)" json:"nums,omitempty"` //增加数量
	GameId uint32 `orm:"column(game_id)" json:"game_id,omitempty"`
}

func (d *TmpGamebetersDao) GetTodayNewerList(channelId uint32, start, over int64) (newersList []NewersList, err error) {
	sql := fmt.Sprintf("SELECT sum(1) nums,game_id FROM %s where %s>=? and %s<? and %s=%d and %s=%d group by %s",
		models.TABLE_TmpGamebeters, models.COLUMN_TmpGamebeters_Ctime,
		models.COLUMN_TmpGamebeters_Ctime, models.COLUMN_TmpGamebeters_BetType, HASBET,
		models.COLUMN_TmpGamebeters_ChannelId, channelId,
		models.COLUMN_TmpGamebeters_GameId)

	_, err = d.Orm.Raw(sql, start, over).QueryRows(&newersList)
	if err != nil {
		common.LogFuncError(fmt.Sprintf("sql err is : %v", err))
	}
	return
}

func (d *TmpGamebetersDao) DeleteByObj(tmpUser models.TmpGamebeters) (err error) {
	if _, err := d.Orm.Delete(&tmpUser); err != nil {
		common.LogFuncError(fmt.Sprintf("sql err is : %v", err))
	}
	return
}
