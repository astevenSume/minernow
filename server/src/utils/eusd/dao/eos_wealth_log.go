package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/eusd/models"
)

const (
	WealthLogType               = iota
	WealthLogTypeOtcInto        //购买
	WealthLogTypeOtcOut         //出售
	WealthLogTypeTransferOut    //转出
	WealthLogTypeTransferInto   //转入
	WealthLogTypeToOtc          //转入承兑
	WealthLogTypeFromOtc        //承兑转出
	WealthLogTypeGameWin        //应用盈利（游戏提现）
	WealthLogTypeGameLost       //应用亏损（游戏充值）
	WealthLogTypeUsdtDelegate   //USDT抵押
	WealthLogTypeUsdtUnDelegate //USDT赎回
	WealthLogTypeCommission     //分润
)

const (
	WealthLogStatus       = iota //转账中
	WealthLogStatusFinish        //完成
	WealthLogStatusFail          //失败
)

type WealthLogDao struct {
	common.BaseDao
}

func NewWealthLogDao(db string) *WealthLogDao {
	return &WealthLogDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *WealthLogDao) Add(uid uint64, ttype uint8, quant int64, txid uint64) (id int64, err error) {
	data := &models.EosWealthLog{
		Uid:      uid,
		TType:    ttype,
		Txid:     txid,
		Quantity: quant,
		Ctime:    common.NowInt64MS(),
	}
	if data.TType == WealthLogTypeToOtc || data.TType == WealthLogTypeFromOtc {
		data.Status = WealthLogStatusFinish
	}
	id, err = d.Orm.Insert(data)
	if err != nil {
		common.LogFuncError("WealthLog DBERR:%v , ", err, ToJson(data))
		return
	}

	return
}

func (d *WealthLogDao) AddBoth(uid, uid2 uint64, ttype, ttype2 uint8, quant int64, txid uint64) (ids []int64, err error) {
	data := &models.EosWealthLog{
		Uid:      uid,
		Uid2:     uid2,
		TType:    ttype,
		Txid:     txid,
		Quantity: quant,
		Ctime:    common.NowInt64MS(),
	}
	data2 := &models.EosWealthLog{
		Uid:      uid2,
		Uid2:     uid,
		TType:    ttype2,
		Txid:     txid,
		Quantity: quant,
		Ctime:    common.NowInt64MS(),
	}
	id, err := d.Orm.Insert(data)
	if err != nil {
		common.LogFuncError("WealthLog DBERR:%v , ", err, ToJson(data))
		return
	}
	id2, err := d.Orm.Insert(data2)
	if err != nil {
		common.LogFuncError("WealthLog DBERR:%v , ", err, ToJson(data2))
		return
	}
	ids = []int64{id, id2}
	return
}

func (d *WealthLogDao) Count(uid uint64, types []interface{}) (num int64) {
	qs := d.Orm.QueryTable(models.TABLE_EosWealthLog).Filter(models.COLUMN_EosWealthLog_Uid, uid)
	if len(types) > 0 {
		qs = qs.Filter(models.COLUMN_EosWealthLog_TType+"__in", types...)
	}
	num, err := qs.Count()

	if err != nil {
		common.LogFuncError("DBERR：%v", err)
		return
	}
	return
}

func (d *WealthLogDao) Fetch(uid uint64, types []interface{}, limit int64, offset int64) (list []*models.EosWealthLog, err error) {
	qs := d.Orm.QueryTable(models.TABLE_EosWealthLog).Filter(models.COLUMN_EosWealthLog_Uid, uid)
	if len(types) > 0 {
		qs = qs.Filter(models.COLUMN_EosWealthLog_TType+"__in", types...)
	}

	qs = qs.Limit(limit).Offset(offset).OrderBy("-" + models.COLUMN_EosWealthLog_Id)

	list = []*models.EosWealthLog{}

	_, err = qs.All(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR：%v", err)

		return
	}

	return
}

func (d *WealthLogDao) UpdateTxid(txid uint64, ids ...int64) (ok bool) {
	qs := d.Orm.QueryTable(models.TABLE_EosWealthLog).Filter(models.COLUMN_EosWealthLog_Id+"__in", ids)
	n, err := qs.Update(orm.Params{models.COLUMN_EosWealthLog_Txid: txid})
	if err != nil {
		common.LogFuncError("DBERR：%v", err)
		return
	}
	if n == 0 {
		return
	}
	ok = true
	return
}

func (d *WealthLogDao) Info(id uint64) (data *models.EosWealthLog) {
	data = &models.EosWealthLog{
		Id: id,
	}
	err := d.Orm.Read(data)
	if err != nil {
		common.LogFuncError("DBERR：%v", err)
		data.Id = 0
		return
	}
	return
}

func (d *WealthLogDao) UpdateStatus(ids ...int64) (ok bool) {
	status := WealthLogStatusFinish
	qs := d.Orm.QueryTable(models.TABLE_EosWealthLog).Filter(models.COLUMN_EosWealthLog_Id+"__in", ids)
	n, err := qs.Update(orm.Params{models.COLUMN_EosWealthLog_Status: status})
	if err != nil {
		common.LogFuncError("DBERR：%v", err)
		return
	}
	if n == 0 {
		return
	}
	ok = true
	return
}

func (d *WealthLogDao) UpdateStatusFail(ids ...int64) (ok bool) {
	status := WealthLogStatusFail
	qs := d.Orm.QueryTable(models.TABLE_EosWealthLog).Filter(models.COLUMN_EosWealthLog_Id+"__in", ids)
	n, err := qs.Update(orm.Params{models.COLUMN_EosWealthLog_Status: status})
	if err != nil {
		common.LogFuncError("DBERR：%v", err)
		return
	}
	if n == 0 {
		return
	}
	ok = true
	return
}

func (d *WealthLogDao) FetchCheck(limit, offset int64) (list []*models.EosWealthLog, err error) {
	qs := d.Orm.QueryTable(models.TABLE_EosWealthLog).Filter(models.COLUMN_EosWealthLog_Status, WealthLogStatus)

	qs = qs.Limit(limit, offset)
	list = []*models.EosWealthLog{}

	_, err = qs.All(&list)
	if err != nil {
		common.LogFuncError("DBERR：%v", err)
		return
	}
	return
}
