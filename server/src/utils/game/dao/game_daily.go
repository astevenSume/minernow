package dao

import (
	"common"
	"fmt"
	agentdao "utils/agent/dao"
	gamemodels "utils/game/models"
)

type GameDailyDao struct {
	common.BaseDao
}

func NewGameDailyDao(db string) *GameDailyDao {
	return &GameDailyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var GameDailyDaoEntity *GameDailyDao

//easyjson:json
type GameDailyApiReq struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	Curpage int    `json:"curpage"`
	Perpage int    `json:"perpage"`
}

//easyjson:json
type GameDailyItem struct {
	Date         string  `json:"date"` //eg:20170818
	WinLoseMoney float64 `json:"winlosemoney"`
	Chips        float64 `json:"chips"`
}

//easyjson:json
type GameDailyApiResp struct {
	Status  int             `json:"status"`
	Desc    string          `json:"desc"`
	Start   int             `json:"start"`
	End     int             `json:"end"`
	Curpage int             `json:"curpage"`
	Perpage int             `json:"perpage"`
	Maxpage int             `json:"maxpage"`
	Data    []GameDailyItem `json:"data"`
}

func (d *GameDailyDao) Sync(host string, platId uint32, timestamp, now int64, start, end string) (integer, decimals int32, err error) {
	//
	// try remove old data
	err = d.RemoveByTimestamp(timestamp)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var integerTmp, decimalsTmp int32
	// loop all channels
	for _, channelId := range []uint32{
		agentdao.CHANNEL_ID_FIRST,
	} {
		// sync new data
		var curPage, pages int
		// page 1
		pages, integerTmp, decimalsTmp, err = d.syncAndSave(host, platId, channelId, timestamp, now, start, end, curPage)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
		integer += integerTmp
		decimals += decimalsTmp

		// page n
		for curPage = 1; curPage < pages; curPage++ {
			_, integerTmp, decimalsTmp, err = d.syncAndSave(host, platId, channelId, timestamp, now, start, end, curPage)
			if err != nil {
				common.LogFuncError("%v", err)
				return
			}
			integer += integerTmp
			decimals += decimalsTmp
		}
	}

	return
}

func (d *GameDailyDao) RemoveByTimestamp(timestamp int64) (err error) {
	dbDaily := gamemodels.ChannelDaily{
		ChannelId: agentdao.CHANNEL_ID_FIRST,
		Ctime:     timestamp,
	}
	_, err = d.Orm.Delete(&dbDaily, gamemodels.COLUMN_ChannelDaily_Ctime)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// sync and save one page
func (d *GameDailyDao) syncAndSave(host string, platId uint32, channelId uint32, timestamp, now int64, start, end string, curPage int) (pages int, winLoseMoneyInteger, winLoseMoneyDecimals int32, err error) {

	var resp GameDailyApiResp
	resp, err = d.subSync(host, platId, start, end, curPage)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	if resp.Status != 200 {
		err = fmt.Errorf("get game daily failed : %+v", resp)
		return
	}

	pages = resp.Maxpage

	winLoseMoneyInteger, winLoseMoneyDecimals, err = d.save(channelId, timestamp, now, resp.Data)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *GameDailyDao) subSync(host string, platId uint32, start, end string, curPage int) (resp GameDailyApiResp, err error) {

	//4debug to delete
	//return GameDailyResp{
	//	Status: 200, Desc: "", Start: 20190414000000, End: 20190414235959, Curpage: 0, Perpage: 20, Maxpage: 1,
	//	Data: []GameDailyItem{
	//		{
	//			Date: "2019/04/14", WinLoseMoney: -2279.8, Chips: 12638.2,
	//		},
	//	},
	//}, nil

	msg := GameDailyApiReq{
		Start:   start,
		End:     end,
		Perpage: 20,
		Curpage: curPage,
	}

	err = GameDaoEntity.SendToChannel(host, PROTO_CHANNEL_DAILY, platId, &msg, &resp)
	if err != nil {
		return
	}

	return
}

// save one page data
func (d *GameDailyDao) save(channelId uint32, timestamp, now int64, list []GameDailyItem) (winLoseMoneyInteger, winLoseMoneyDecimals int32, err error) {
	dbDailies := make([]gamemodels.ChannelDaily, 0, len(list))
	for _, item := range list {
		var tmpTimestamp uint32
		tmpTimestamp, err = common.TimeStrToUint32Plus(item.Date)
		if int64(tmpTimestamp) != timestamp { //skip time no fit one
			continue
		}

		dbDaily := gamemodels.ChannelDaily{
			ChannelId: channelId,
			Ctime:     timestamp,
			Mtime:     now,
		}

		dbDaily.WinLoseMoneyInteger, dbDaily.WinLoseMoneyDecimals, err = common.DecodeCurrencyNoCarePrecision(fmt.Sprint(item.WinLoseMoney))
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
		dbDaily.ChipsInteger, dbDaily.ChipsDecimals, err = common.DecodeCurrencyNoCarePrecision(fmt.Sprint(item.Chips))
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}

		winLoseMoneyInteger, winLoseMoneyDecimals = common.AddCurrency2(winLoseMoneyInteger, winLoseMoneyDecimals,
			dbDaily.WinLoseMoneyInteger, dbDaily.WinLoseMoneyDecimals)

		if winLoseMoneyInteger < 0 {
			winLoseMoneyInteger *= -1
		}

		if winLoseMoneyDecimals < 0 {
			winLoseMoneyDecimals *= -1
		}

		dbDailies = append(dbDailies, dbDaily)
	}

	_, err = d.Orm.InsertMulti(common.BulkCount, &dbDailies)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
