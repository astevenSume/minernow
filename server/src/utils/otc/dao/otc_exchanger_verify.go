package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/otc/models"
)

//承兑商审核状态
const (
	ExchangerVerifyStatusNil = iota
	ExchangerVerifyStatusPending
	ExchangerVerifyStatusActive
	ExchangerVerifyStatusReject
	ExchangerVerifyStatusMax
)

//承兑商申请渠道
const (
	ExchangerApply = iota
	ExchangerAssign
)

type OtcExchangerVerifyDao struct {
	common.BaseDao
}

func NewOtcExchangerVerifyDao(db string) *OtcExchangerVerifyDao {
	return &OtcExchangerVerifyDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var OtcExchangerVerifyDaoEntity *OtcExchangerVerifyDao

type OtcExchangerVerifyAck struct {
	Id       int32  `json:"id"`
	Uid      string `json:"user_id"`
	Mobile   string `json:"mobile"`
	Wechat   string `json:"wechat"`
	Telegram string `json:"telegram"`
	Status   int8   `json:"status"`
	From     int8   `json:"from"`
	Ctime    int64  `json:"ctime"`
	Utime    int64  `json:"utime"`
}

func (d *OtcExchangerVerifyDao) ClientExchangerVerify(data *models.OtcExchangerVerify) (ack OtcExchangerVerifyAck) {
	if data == nil {
		return
	}
	ack.Id = data.Id
	ack.Uid = fmt.Sprintf("%v", data.Uid)
	ack.Mobile = data.Mobile
	ack.Wechat = data.Wechat
	ack.Telegram = data.Telegram
	ack.Status = data.Status
	ack.From = data.From
	ack.Ctime = data.Ctime
	ack.Utime = data.Utime
	return
}

func (d *OtcExchangerVerifyDao) ClientExchangerVerifys(datas []models.OtcExchangerVerify) (acks []OtcExchangerVerifyAck) {
	for _, data := range datas {
		var ack OtcExchangerVerifyAck
		ack.Id = data.Id
		ack.Uid = fmt.Sprintf("%v", data.Uid)
		ack.Mobile = data.Mobile
		ack.Wechat = data.Wechat
		ack.Telegram = data.Telegram
		ack.Status = data.Status
		ack.From = data.From
		ack.Ctime = data.Ctime
		ack.Utime = data.Utime
		acks = append(acks, ack)
	}

	return
}

//承兑商审核
func (d *OtcExchangerVerifyDao) Create(uid uint64, from int8, mobile, wechat, telegram string) (data models.OtcExchangerVerify, err error) {
	data = models.OtcExchangerVerify{
		Uid:      uid,
		Mobile:   mobile,
		Wechat:   wechat,
		Telegram: telegram,
		From:     from,
		Ctime:    common.NowInt64MS(),
		Utime:    common.NowInt64MS(),
		Status:   ExchangerVerifyStatusPending,
	}
	if from == ExchangerAssign {
		data.Status = ExchangerVerifyStatusActive
	}

	id, err := d.Orm.Insert(&data)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
		return
	}
	data.Id = int32(id)

	return
}

func (d *OtcExchangerVerifyDao) LastRow(uid uint64) (data *models.OtcExchangerVerify, err error) {
	data = &models.OtcExchangerVerify{}

	err = d.Orm.QueryTable(data).Filter(models.COLUMN_OtcExchangerVerify_Uid, uid).
		OrderBy("-" + models.COLUMN_OtcExchangerVerify_Id).Limit(1).One(data)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DB_ERR:%v", err)
		return
	}

	return
}

func (d *OtcExchangerVerifyDao) QueryById(id int32) (otcExchangerVerify *models.OtcExchangerVerify, err error) {
	otcExchangerVerify = &models.OtcExchangerVerify{
		Id: id,
	}

	err = d.Orm.Read(otcExchangerVerify)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DB_ERR:%v", err)
		return
	}

	return
}

//条件查询
func (d *OtcExchangerVerifyDao) QueryCondition(mobile string, wechat string, status int, page int, perPage int) (total int, otcExchangerVerify []models.OtcExchangerVerify, err error) {
	qs := d.Orm.QueryTable(models.TABLE_OtcExchangerVerify)
	if qs == nil {
		common.LogFuncError("QueryTable fail")
		return
	}

	if mobile != "" {
		qs = qs.Filter(models.COLUMN_OtcExchangerVerify_Mobile+"__contains", mobile)
	}
	if wechat != "" {
		qs = qs.Filter(models.COLUMN_OtcExchangerVerify_Wechat, wechat)
	}
	if status > ExchangerVerifyStatusNil && status < ExchangerVerifyStatusMax {
		qs = qs.Filter(models.COLUMN_OtcExchangerVerify_Status, status)
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
	_, err = qs.Limit(perPage, start).All(&otcExchangerVerify)
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

//更新承兑商审核状态
func (d *OtcExchangerVerifyDao) UpdateStatus(id int32, status int8) (models.OtcExchangerVerify, error) {
	otcExchangerVerify := models.OtcExchangerVerify{Id: id}
	err := d.Orm.Read(&otcExchangerVerify, models.COLUMN_OtcExchangerVerify_Id)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return otcExchangerVerify, err
	}

	otcExchangerVerify.Status = status
	otcExchangerVerify.Utime = common.NowInt64MS()
	_, err = d.Orm.Update(&otcExchangerVerify, models.COLUMN_OtcExchangerVerify_Status, models.COLUMN_OtcExchangerVerify_Utime)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return otcExchangerVerify, err
	}

	return otcExchangerVerify, nil
}
