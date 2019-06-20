package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"utils/game/models"
)

type GameTransferDao struct {
	common.BaseDao
}

func NewGameTransferDao(db string) *GameTransferDao {
	return &GameTransferDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var GameTransferDaoEntity *GameTransferDao

type TRANSFER_TYPE int

const (
	TRANSFER_TYPE_UNKOWN TRANSFER_TYPE = iota
	TRANSFER_TYPE_IN                   //充值
	TRANSFER_TYPE_OUT                  //提现
)

var TRANSFER_TYPE_STR = map[TRANSFER_TYPE]string{
	TRANSFER_TYPE_UNKOWN: "UNKOWN",
	TRANSFER_TYPE_IN:     "IN",
	TRANSFER_TYPE_OUT:    "OUT",
}

type TRANSFER_STATUS int

const (
	TRANSFER_STATUS_INIT TRANSFER_STATUS = iota
	TRANSFER_STATUS_DONE
	TRANSFER_STATUS_FAILED
)

func (d *GameTransferDao) Add(uid uint64, channelId uint32, account string, transferType TRANSFER_TYPE, coin, eusd int64, step string) (transfer models.GameTransfer, err error) {
	var id uint64
	id, err = common.IdManagerGen(IdTypeIndexGameTransfer)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	now := common.NowInt64MS()
	transfer.Id = id
	transfer.Uid = uid
	transfer.ChannelId = channelId
	transfer.Account = account
	transfer.TransferType = uint32(transferType)
	transfer.CoinInteger = coin
	transfer.EusdInteger = eusd
	transfer.Ctime = now
	transfer.Order = d.transferOrder(transferType)
	transfer.Status = uint32(TRANSFER_STATUS_INIT)
	transfer.Step = step
	_, err = d.Orm.Insert(&transfer)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *GameTransferDao) transferOrder(transferType TRANSFER_TYPE) string {
	if k, ok := TRANSFER_TYPE_STR[transferType]; ok {
		return k + common.NowStringMillS()
	}
	return TRANSFER_TYPE_STR[TRANSFER_TYPE_UNKOWN] + common.NowString()
}

func (d *GameTransferDao) Update(transfer models.GameTransfer) (err error) {
	_, err = d.Orm.Update(&transfer)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

type GameTransferGrandTotal struct {
	ChannelId    uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	TransferType uint32 `orm:"column(transfer_type)" json:"transfer_type,omitempty"`
	CoinInteger  int64  `orm:"column(coin_integer)" json:"coin_integer,omitempty"`
	EusdInteger  int64  `orm:"column(eusd_integer)" json:"eusd_integer,omitempty"`
}

// GrandTotal 统计累计充值提现金额
func (d *GameTransferDao) GrandTotal(uid uint64, groupByChannel bool, begin, end int64) (result []GameTransferGrandTotal, err error) {
	var (
		sqlStr, whereStr        string
		selectCols, groupbyCols []string
		params                  []interface{}
		qbQuery                 orm.QueryBuilder
	)

	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}

	selectCols = []string{
		models.COLUMN_GameTransfer_TransferType,
		fmt.Sprintf("SUM(%s)AS %s", models.COLUMN_GameTransfer_CoinInteger, models.COLUMN_GameTransfer_CoinInteger),
		fmt.Sprintf("SUM(%s)AS %s", models.COLUMN_GameTransfer_EusdInteger, models.COLUMN_GameTransfer_EusdInteger),
	}

	groupbyCols = []string{models.COLUMN_GameTransfer_TransferType}

	if groupByChannel {
		selectCols = append(selectCols, models.COLUMN_GameTransfer_ChannelId)
		groupbyCols = append(groupbyCols, models.COLUMN_GameTransfer_ChannelId)
	}

	whereStr = fmt.Sprintf(" %s=? and %s=? ", models.COLUMN_GameTransfer_Uid, models.COLUMN_GameTransfer_Status)
	params = []interface{}{uid, TRANSFER_STATUS_DONE}

	if begin > 0 {
		whereStr = fmt.Sprintf(" %s and %s>=? ", whereStr, models.COLUMN_GameTransfer_Ctime)
		params = append(params, begin)
	}
	if end > 0 {
		whereStr = fmt.Sprintf(" %s and %s<=? ", whereStr, models.COLUMN_GameTransfer_Ctime)
		params = append(params, end)
	}

	sqlStr = qbQuery.Select(selectCols...).From(models.TABLE_GameTransfer).Where(whereStr).GroupBy(groupbyCols...).String()
	_, err = d.Orm.Raw(sqlStr, params...).QueryRows(&result)
	return
}

//充值 & 提现 总额
func (d *GameTransferDao) CountByTime(start, end int64) (recharge, withdrawal uint32) {
	sql := fmt.Sprintf("Select sum(%s) as eusd,%s as t From %s Where %s>=? And %s<? Group By %s",
		models.COLUMN_GameTransfer_EusdInteger, models.COLUMN_GameTransfer_TransferType, models.TABLE_GameTransfer,
		models.COLUMN_GameTransfer_Ctime, models.COLUMN_GameTransfer_Ctime, models.COLUMN_GameTransfer_TransferType)

	type C struct {
		Eusd uint32 `orm:"eusd"`
		Type int    `orm:"t"`
	}
	res := []*C{}
	_, err := d.Orm.Raw(sql, start, end).QueryRows(&res)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}
	for _, v := range res {
		if v.Type == int(TRANSFER_TYPE_IN) {
			recharge = v.Eusd
		} else if v.Type == int(TRANSFER_TYPE_OUT) {
			withdrawal = v.Eusd
		}
	}

	return
}

type TransferRecord struct {
	Uid       uint64 `orm:"uid" `
	ChannelId uint32 `orm:"channel_id"`
	Eusd      uint32 `orm:"eusd"`
	Type      int    `orm:"t"`
}

//玩家的充值 & 提现 记录
func (d *GameTransferDao) FindByTime(start, end int64) (res []*TransferRecord, err error) {
	sql := fmt.Sprintf("Select sum(%s) as eusd,%s as t From %s Where %s>=? And %s<? Group By %s",
		models.COLUMN_GameTransfer_EusdInteger, models.COLUMN_GameTransfer_TransferType, models.TABLE_GameTransfer,
		models.COLUMN_GameTransfer_Ctime, models.COLUMN_GameTransfer_Ctime, models.COLUMN_GameTransfer_TransferType)

	_, err = d.Orm.Raw(sql, start, end).QueryRows(&res)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

func (d *GameTransferDao) QueryPageTransferList(uid uint64, channelID, status string, page, limit int) (total int64, list []models.GameTransfer, err error) {

	qs := d.Orm.QueryTable(models.TABLE_GameTransfer)

	qs = qs.Filter(models.COLUMN_GameTransfer_Uid, uid)

	if channelID != "" {
		qs = qs.Filter(models.COLUMN_GameTransfer_ChannelId, channelID)
	}

	if status != "" {
		qs = qs.Filter(models.COLUMN_GameTransfer_Status, status)
	}

	total, err = qs.Count()
	if err != nil && err != orm.ErrNoRows {
		common.LogFuncError("%v", err)
		return
	}

	_, err = qs.Offset((page - 1) * limit).Limit(limit).All(&list)
	if err != nil && err != orm.ErrNoRows {
		common.LogFuncError("%v", err)
		return
	}

	err = nil
	return
}

func (d *GameTransferDao) QueryPageTransferCharge(uid string, transferType, page, limit int) (total int64, list []models.GameTransfer, err error) {

	qs := d.Orm.QueryTable(models.TABLE_GameTransfer)

	qs = qs.Filter(models.COLUMN_GameTransfer_Uid, uid)
	if transferType != 0 {
		qs = qs.Filter(models.COLUMN_GameTransfer_TransferType, transferType)
	}
	qs = qs.Offset((page - 1) * limit).Limit(limit)
	total, err = qs.Count()
	if err != nil && err != orm.ErrNoRows {
		common.LogFuncError("%v", err)
		return
	}

	_, err = qs.Offset((page - 1) * limit).Limit(limit).All(&list)
	if err != nil && err != orm.ErrNoRows {
		common.LogFuncError("%v", err)
		return
	}

	err = nil
	return
}

func getTodayStartTime() (ts int64) {
	t := time.Now()
	tm := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return tm.Unix()
}

func (d *GameTransferDao) QueryTodayRiskInfo() (userRiskList []models.GameTransfer) {
	nowTime := strconv.FormatInt(time.Now().Unix(), 10)
	t := strconv.FormatInt(getTodayStartTime(), 10)

	//models.COLUMN_GameTransfer_CoinInteger is the sum
	sql := fmt.Sprintf("select %s,sum(%s) %s, %s from %s where ctime between %s and %s and %s=%d group by uid",
		models.COLUMN_GameTransfer_Uid, models.COLUMN_GameTransfer_EusdInteger,
		models.COLUMN_GameTransfer_EusdInteger, models.COLUMN_GameTransfer_Ctime,
		models.TABLE_GameTransfer, t, nowTime, models.COLUMN_GameTransfer_TransferType, TRANSFER_TYPE_OUT)
	_, err := d.Orm.Raw(sql).QueryRows(&userRiskList)
	if err != nil {
		common.LogFuncError("mysql_err: %v", err)
	}
	return
}

//get recharge or withdraw by channelID and transferType
func (d *GameTransferDao) QueryFundRwByChannelID(channelId, transferType uint32, start, over int64) (fund int64, err error) {
	var fundObj = models.GameTransfer{}
	sql := fmt.Sprintf("SELECT SUM(%s) %s FROM %s WHERE %s=? AND %s=? AND %s>=? AND %s<?",
		models.COLUMN_GameTransfer_EusdInteger, models.COLUMN_GameTransfer_EusdInteger, models.TABLE_GameTransfer,
		models.COLUMN_GameTransfer_ChannelId, models.COLUMN_GameTransfer_TransferType,
		models.COLUMN_GameTransfer_Ctime, models.COLUMN_GameTransfer_Ctime)

	if err = d.Orm.Raw(sql, channelId, transferType, start, over).QueryRow(&fundObj); err != nil {
		common.LogFuncError("mysql_err: %v", err)
	}
	fund = fundObj.EusdInteger
	return
}
