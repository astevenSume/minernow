package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	agentdao "utils/agent/dao"
	gamemodels "utils/game/models"
)

type GameUserDailyDao struct {
	common.BaseDao
}

func NewGameUserDailyDao(db string) *GameUserDailyDao {
	return &GameUserDailyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var GameUserDailyDaoEntity *GameUserDailyDao

func (d *GameUserDailyDao) RemoveByTimestamp(timestamp int64) (err error) {
	dbDaily := gamemodels.GameUserDaily{
		Ctime: timestamp,
	}
	_, err = d.Orm.Delete(&dbDaily, gamemodels.COLUMN_GameUserDaily_Ctime)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *GameUserDailyDao) QueryTotal(timestamp int64) (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(*) FROM %s WHERE %s=?",
		gamemodels.TABLE_GameUserDaily, gamemodels.COLUMN_GameUserDaily_Ctime), timestamp).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// get users' chips
func (d *GameUserDailyDao) Chips(channelId uint32, timestamp int64, uids []uint64) (chipsInteger, chipsDecimals int32, err error) {
	var uidStrs []string
	for _, uid := range uids {
		uidStrs = append(uidStrs, fmt.Sprint(uid))
	}
	var qb orm.QueryBuilder
	qb, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var integer, decimals uint32
	qb.Select(fmt.Sprintf("sum(%s), sum(%s)",
		gamemodels.COLUMN_GameUserDaily_ChipsInteger,
		gamemodels.COLUMN_GameUserDaily_ChipsDecimals)).
		From(gamemodels.TABLE_GameUserDaily).
		Where(gamemodels.COLUMN_GameUserDaily_Uid).In(uidStrs...).
		And(fmt.Sprintf("%s=?", gamemodels.COLUMN_GameUserDaily_ChannelId)).
		And(fmt.Sprintf("%s=?", gamemodels.COLUMN_GameUserDaily_Ctime))
	err = d.Orm.Raw(qb.String(), channelId, timestamp).QueryRow(&integer, &decimals)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	sum := int64(integer)*int64(common.MaxDecimals) + int64(decimals)
	chipsInteger = int32(sum / common.MaxDecimals)
	chipsDecimals = int32(sum % common.MaxDecimals)

	return
}

// get taxes
func (d *GameUserDailyDao) Taxes(timestamp int64) (taxInteger, taxDecimals int32, err error) {
	var integer, decimals int32
	err = d.Orm.Raw(fmt.Sprintf("SELECT sum(%s), sum(%s) FROM %s WHERE %s=?",
		gamemodels.COLUMN_GameUserDaily_TaxInteger,
		gamemodels.COLUMN_GameUserDaily_TaxDecimals,
		gamemodels.TABLE_GameUserDaily,
		gamemodels.COLUMN_GameUserDaily_Ctime), timestamp).QueryRow(&integer, &decimals)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	sum := common.AmountToInt64(integer, decimals)
	taxInteger, taxDecimals = common.Int64ToAmount(sum)

	return
}

// sync user daily
func (d *GameUserDailyDao) Sync(host string, platId uint32, timestamp, now int64, start, end string) (err error) {
	// try remove old data
	err = d.RemoveByTimestamp(timestamp)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	// loop all channels
	for _, channelId := range []uint32{
		agentdao.CHANNEL_ID_FIRST,
	} {
		// sync new data
		var curPage, pages int

		pages, err = d.syncAndSave(host, platId, channelId, timestamp, now, start, end, curPage)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}

		for curPage = 1; curPage < pages; curPage++ {
			_, err = d.syncAndSave(host, platId, channelId, timestamp, now, start, end, curPage)
			if err != nil {
				common.LogFuncError("%v", err)
				return
			}
		}
	}

	return
}

//easyjson:json
type GameUserDailyApiReq struct {
	Account   string `json:"account"`
	Subplatid string `json:"subplatid"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Curpage   int    `json:"curpage"`
	Perpage   int    `json:"perpage"`
}

//easyjson:json
type GameUserDailyApiItem struct {
	Date         string  `json:"date"`
	Userid       uint64  `json:"userid"`
	Account      string  `json:"account"`
	Winlosemoney float64 `json:"winlosemoney"`
	Chips        float64 `json:"chips"`
	Tax          float64 `json:"tax"`
}

//easyjson:json
type GameUserDailyApiResp struct {
	Status  int                    `json:"status"`
	Desc    string                 `json:"desc"`
	Start   string                 `json:"start"`
	End     string                 `json:"end"`
	Curpage int                    `json:"curpage"`
	Perpage int                    `json:"perpage"`
	Maxpage int                    `json:"maxpage"`
	Data    []GameUserDailyApiItem `json:"data"`
}

// sync and save one page
func (d *GameUserDailyDao) syncAndSave(host string, platId uint32, channelId uint32, timestamp, now int64, start, end string, curPage int) (pages int, err error) {

	var resp GameUserDailyApiResp
	resp, err = d.subSync(host, platId, start, end, curPage)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	if resp.Status != 200 {
		err = fmt.Errorf("get game user daily failed : %+v", resp)
		return
	}

	pages = resp.Maxpage

	err = d.save(channelId, timestamp, now, resp.Data)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// sync from game channel
func (d *GameUserDailyDao) subSync(host string, platId uint32, start, end string, curPage int) (resp GameUserDailyApiResp, err error) {
	//4debug, to delete
	//return &GameUserDailyResp{
	//	Status:  200,
	//	Start:   "20190414000000",
	//	End:     "20190414235959",
	//	Curpage: 0,
	//	Perpage: 20,
	//	Maxpage: 1,
	//	Data: []GameUserDailyItem{
	//		{Date: "2019/04/14", Userid: 57180, Account: "U169444543387664384", Winlosemoney: 133.2, Chips: 1530, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57180, Account: "U169444543387664384", Winlosemoney: 285, Chips: 3315, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57185, Account: "U169454230048866304", Winlosemoney: -7.8, Chips: 12.6, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57188, Account: "U169445993824124928", Winlosemoney: 196, Chips: 200, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57188, Account: "U169445993824124928", Winlosemoney: -200, Chips: 500, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57188, Account: "U169445993824124928", Winlosemoney: 138, Chips: 1, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57188, Account: "U169445993824124928", Winlosemoney: 9, Chips: 2, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57188, Account: "U169445993824124928", Winlosemoney: -1280, Chips: 1560, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57192, Account: "U169442933404073984", Winlosemoney: -2.2, Chips: 2.6, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57193, Account: "U169446009552764928", Winlosemoney: -75, Chips: 205, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57197, Account: "U169447921719181312", Winlosemoney: 49, Chips: 50, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57199, Account: "U169643538705809408", Winlosemoney: -2500, Chips: 3500, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57203, Account: "U169723495444381696", Winlosemoney: 0, Chips: 55, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57203, Account: "U169723495444381696", Winlosemoney: -65, Chips: 105, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57209, Account: "U169751779917955072", Winlosemoney: 970, Chips: 1010, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57209, Account: "U169751779917955072", Winlosemoney: -30, Chips: 30, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57209, Account: "U169751779917955072", Winlosemoney: -50, Chips: 70, Tax: 0},
	//		{Date: "2019/04/14", Userid: 57212, Account: "U169750166759276544", Winlosemoney: 150, Chips: 490, Tax: 0},
	//	},
	//}, nil

	msg := GameUserDailyApiReq{
		Start:   start,
		End:     end,
		Perpage: 20,
		Curpage: curPage,
	}

	//resp = GameUserDailyApiResp{}

	err = GameDaoEntity.SendToChannel(host, PROTO_USER_DAILY, platId, &msg, &resp)
	if err != nil {
		return
	}

	return
}

// save one page data
func (d *GameUserDailyDao) save(channelId uint32, timestamp, now int64, list []GameUserDailyApiItem) (err error) {
	tmpMap := make(map[string]*gamemodels.GameUserDaily)
	for _, item := range list {
		var gameUser gamemodels.GameUser
		gameUser, err = GameUserDaoEntity.QueryByAccount(item.Account, channelId)
		if err != nil {
			common.LogFuncWarning("query game user of %s %v", item.Account, err)
			err = nil
			continue
		}

		dbDaily := gamemodels.GameUserDaily{
			ChannelId: channelId,
			Uid:       gameUser.Uid,
			Ctime:     timestamp,
			Mtime:     now,
		}

		dbDaily.TaxInteger, dbDaily.TaxDecimals, err = common.DecodeCurrencyNoCarePrecision(fmt.Sprint(item.Tax))
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
		dbDaily.ChipsInteger, dbDaily.ChipsDecimals, err = common.DecodeCurrencyNoCarePrecision(fmt.Sprint(item.Chips))
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
		dbDaily.WinloseInteger, dbDaily.WinloseDecimals, err = common.DecodeCurrencyNoCarePrecision(fmt.Sprint(item.Winlosemoney))
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}

		// todo to delete, add the same user's data up.
		key := fmt.Sprintf("%d%d%s", channelId, timestamp, dbDaily.Uid)
		if _, ok := tmpMap[key]; !ok {
			tmpMap[key] = &dbDaily
		} else {
			tmpMap[key].TaxInteger, tmpMap[key].TaxDecimals = common.AddCurrency2(tmpMap[key].TaxInteger, tmpMap[key].TaxDecimals, dbDaily.TaxInteger, dbDaily.TaxDecimals)
			tmpMap[key].ChipsInteger, tmpMap[key].ChipsDecimals = common.AddCurrency2(tmpMap[key].ChipsInteger, tmpMap[key].ChipsDecimals, dbDaily.ChipsInteger, dbDaily.ChipsDecimals)
			tmpMap[key].WinloseInteger, tmpMap[key].WinloseDecimals = common.AddCurrency2(tmpMap[key].WinloseInteger, tmpMap[key].WinloseDecimals, dbDaily.WinloseInteger, dbDaily.WinloseDecimals)
		}

		//dbDailies = append(dbDailies, dbDaily)
	}

	gameUserInMapStr := "game user daily in map : \n"
	dbDailies := make([]gamemodels.GameUserDaily, 0, len(tmpMap))
	for k, v := range tmpMap {
		gameUserInMapStr += fmt.Sprintf("%s : %+v\n", k, v)
		dbDailies = append(dbDailies, *v)
	}

	dbDailiesStr := "game user dailys : \n"
	for _, daily := range dbDailies {
		dbDailiesStr += fmt.Sprintf("%+v\n", daily)
		err = d.checkAndUpdateSingle(daily)
		if err != nil {
			return
		}
	}
	common.LogFuncDebug(dbDailiesStr)

	//_, err = d.Orm.InsertMulti(common.BulkCount, &dbDailies)
	//if err != nil {
	//	common.LogFuncError("%v", err)
	//	return
	//}

	return
}

func (d *GameUserDailyDao) checkAndUpdateSingle(gameUserDaily gamemodels.GameUserDaily) (err error) {
	// try get
	tmp := gamemodels.GameUserDaily{
		ChannelId: gameUserDaily.ChannelId,
		Uid:       gameUserDaily.Uid,
		Ctime:     gameUserDaily.Ctime,
	}
	err = d.Orm.Read(&tmp, gamemodels.COLUMN_GameUserDaily_ChannelId, gamemodels.COLUMN_GameUserDaily_Uid, gamemodels.COLUMN_GameUserDaily_Ctime)
	if err != nil {
		if err == orm.ErrNoRows { //no exist, just insert
			_, err = d.Orm.Insert(&gameUserDaily)
			if err != nil {
				common.LogFuncError("%v", err)
				return
			}
			return
		}
		common.LogFuncError("%v", err)
		return
	}

	gameUserDaily.TaxInteger, gameUserDaily.TaxDecimals = common.AddCurrency2(gameUserDaily.TaxInteger, gameUserDaily.TaxDecimals, tmp.TaxInteger, tmp.TaxDecimals)
	gameUserDaily.ChipsInteger, gameUserDaily.ChipsDecimals = common.AddCurrency2(gameUserDaily.ChipsInteger, gameUserDaily.ChipsDecimals, tmp.ChipsInteger, tmp.ChipsDecimals)
	gameUserDaily.WinloseInteger, gameUserDaily.WinloseDecimals = common.AddCurrency2(gameUserDaily.WinloseInteger, gameUserDaily.WinloseDecimals, tmp.WinloseInteger, tmp.WinloseDecimals)
	// exist, add currency
	_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=?,%s=?,%s=?,%s=?,%s=?,%s=? WHERE %s=? AND %s=? AND %s=?",
		gamemodels.TABLE_GameUserDaily,
		gamemodels.COLUMN_GameUserDaily_TaxInteger, gamemodels.COLUMN_GameUserDaily_TaxDecimals,
		gamemodels.COLUMN_GameUserDaily_ChipsInteger, gamemodels.COLUMN_GameUserDaily_ChipsDecimals,
		gamemodels.COLUMN_GameUserDaily_WinloseInteger, gamemodels.COLUMN_GameUserDaily_WinloseDecimals,
		gamemodels.COLUMN_GameUserDaily_ChannelId, gamemodels.COLUMN_GameUserDaily_Uid, gamemodels.COLUMN_GameUserDaily_Ctime),
		gameUserDaily.TaxInteger, gameUserDaily.TaxDecimals, gameUserDaily.ChipsInteger, gameUserDaily.ChipsDecimals,
		gameUserDaily.WinloseInteger, gameUserDaily.WinloseDecimals, gameUserDaily.ChannelId, gameUserDaily.Uid, gameUserDaily.Ctime).Exec()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
