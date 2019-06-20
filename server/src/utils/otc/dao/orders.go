package dao

import (
	"common"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"utils/otc/models"
)

type OrdersDao struct {
	common.BaseDao
}

func NewOrdersDao(db string) *OrdersDao {
	return &OrdersDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var OrdersDaoEntity *OrdersDao

const (
	SideBuy  int8 = 1 //用户购买
	SideSell int8 = 2 //用户出售
)

const (
	OrderStatus             = iota
	OrderStatusCreated      // 刚创建，未付款
	OrderStatusPayed        // 已付款, 未确认
	OrderStatusConfirmed    // 已确认
	OrderStatusCanceled     // 取消订单
	OrderStatusAppeal       // 申诉
	OrderStatusExpired      // 已过期
	OrderStatusTransferring // 转账中
	OrderStatusFinish
)

const (
	OrderAppealStatusNil        = iota
	OrderAppealStatusPending    //等待处理
	OrderAppealStatusProcessing //正在处理
	OrderAppealStatusResolved   //已解决
	OrderAppealStatus
)

//RISK
const (
	RISK_SOLD_RATE  float64 = 0.8
	RISK_EUSD_VALUE float64 = 300
	RISK_TIME_RANGE int64   = 24 * 3600
)

type OrderAck struct {
	Id             string `orm:"column(id);pk" json:"id,omitempty"`
	Uid            string `orm:"column(uid)" json:"uid,omitempty"`
	Uip            string `orm:"column(uip);size(40)" json:"uip,omitempty"`
	EUid           string `orm:"column(euid)" json:"euid,omitempty"`
	Eip            string `orm:"column(eip);size(40)" json:"euip,omitempty"`
	Side           int8   `orm:"column(side)" json:"side,omitempty"`
	Amount         string `orm:"column(amount)" json:"amount,omitempty"`
	Price          string `orm:"column(price);size(100)" json:"price,omitempty"`
	Funds          string `orm:"column(funds)" json:"funds,omitempty"`
	Fee            string `orm:"column(fee)" json:"fee,omitempty"`
	PayId          string `orm:"column(pay_id)" json:"pay_id,omitempty"`
	PayType        int8   `orm:"column(pay_type)" json:"pay_type,omitempty"`
	PayName        string `orm:"column(pay_name);size(128)" json:"pay_name,omitempty"`
	PayAccount     string `orm:"column(pay_account);size(300)" json:"pay_account,omitempty"`
	PayBank        string `orm:"column(pay_bank);size(128)" json:"pay_bank,omitempty"`
	PayBankBranch  string `orm:"column(pay_bank_branch);size(300)" json:"pay_bank_branch,omitempty"`
	TransferId     string `orm:"column(transfer_id)" json:"rid,omitempty"`
	Ctime          int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	PayTime        int64  `orm:"column(pay_time)" json:"pay_time,omitempty"`
	FinishTime     int64  `orm:"column(finish_time)" json:"ftime,omitempty"`
	Utime          int64  `orm:"column(utime)" json:"utime,omitempty"`
	Status         int8   `orm:"column(status)" json:"status,omitempty"`
	EPayId         string `orm:"column(epay_id)" json:"epay_id,omitempty"`
	EPayType       int8   `orm:"column(epay_type)" json:"epay_type,omitempty"`
	EPayName       string `orm:"column(epay_name);size(128)" json:"epay_name,omitempty"`
	EPayAccount    string `orm:"column(epay_account);size(300)" json:"epay_account,omitempty"`
	EPayBank       string `orm:"column(epay_bank);size(128)" json:"epay_bank,omitempty"`
	EPayBankBranch string `orm:"column(epay_bank_branch);size(300)" json:"epay_bank_branch,omitempty"`
	AppealStatus   int8   `orm:"column(appeal_status)" json:"appeal_status,omitempty"`
	AdminId        uint32 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	QrCode         string `orm:"column(qr_code);size(300)" json:"qr_code,omitempty"`
	UMobile        string `json:"u_mobile"`
	UQrCode        string `json:"u_qr_code"`
	EuMobile       string `json:"eu_mobile"`
	EuQrCode       string `json:"eu_qr_code"`
}

type OrderExchangeTotal struct {
	Date   int32 `orm:"column(date)" json:"date"`
	Amount int64 `orm:"column(amount)" json:"amount"`
	Funds  int64 `orm:"column(funds)" json:"funds"`
}

func (d *OrdersDao) ClientOrder(data *models.OtcOrder) (ack OrderAck) {
	if data == nil {
		return
	}
	ack.Id = fmt.Sprintf("%v", data.Id)
	ack.Uid = fmt.Sprintf("%v", data.Uid)
	ack.Uip = data.Uip
	ack.EUid = fmt.Sprintf("%v", data.EUid)
	ack.Eip = data.Eip
	ack.Side = data.Side
	ack.Amount = fmt.Sprintf("%v", data.Amount)
	ack.Price = data.Price
	ack.Funds = fmt.Sprintf("%v", data.Funds)
	ack.Fee = fmt.Sprintf("%v", data.Fee)
	ack.PayId = fmt.Sprintf("%v", data.PayId)
	ack.PayType = data.PayType
	ack.PayName = data.PayName
	ack.PayAccount = data.PayAccount
	ack.PayBank = data.PayBank
	ack.PayBankBranch = data.PayBankBranch
	ack.TransferId = fmt.Sprintf("%v", data.TransferId)
	ack.Ctime = data.Ctime
	ack.Utime = data.Utime
	ack.PayTime = data.PayTime
	ack.FinishTime = data.FinishTime
	ack.Status = data.Status
	ack.EPayId = fmt.Sprintf("%v", data.EPayId)
	ack.EPayType = data.EPayType
	ack.EPayName = data.EPayName
	ack.EPayAccount = data.EPayAccount
	ack.EPayBank = data.EPayBank
	ack.EPayBankBranch = data.EPayBankBranch
	ack.AppealStatus = data.AppealStatus
	ack.AdminId = data.AdminId

	return
}

func (d *OrdersDao) Create(order *models.OtcOrder) (orderRes *models.OtcOrder, err error) {

	id, err := common.IdManagerGen(IdTypeOrders)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	t := time.Unix(order.Ctime/1000, 0)
	date, err := strconv.ParseInt(t.Format("20060102"), 10, 64)
	order.Date = int32(date)
	order.Id = id
	_, err = d.Orm.Insert(order)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	orderRes = order
	return
}

func (d *OrdersDao) Info(id uint64) (order *models.OtcOrder, err error) {
	order = &models.OtcOrder{
		Id: id,
	}

	err = d.Orm.Read(order)
	if err != nil {
		if err == orm.ErrNoRows {
			order.Id = 0
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}

func (d *OrdersDao) EditEPay(order *models.OtcOrder) (ok bool) {
	n, err := d.Orm.Update(order, models.COLUMN_OtcOrder_EPayAccount, models.COLUMN_OtcOrder_EPayBank,
		models.COLUMN_OtcOrder_EPayBankBranch, models.COLUMN_OtcOrder_EPayId,
		models.COLUMN_OtcOrder_EPayName, models.COLUMN_OtcOrder_EPayType)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	if n == 0 {
		return
	}
	ok = true
	return
}

//取消订单
func (d *OrdersDao) Cancel(id uint64) (ok bool) {
	sql := fmt.Sprintf("update %s set %s=?,%s=? where %s=? And %s=?",
		models.TABLE_OtcOrder,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_FinishTime,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Id)

	res, err := d.Orm.Raw(sql, OrderStatusCanceled, common.NowInt64MS(), OrderStatusCreated, id).Exec()

	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return
	}
	ok = true
	return
}

//申述取消订单
func (d *OrdersDao) AppealCancel(id uint64) (ok bool) {
	sql := fmt.Sprintf("update %s set %s=?,%s=? where %s=? And %s=?",
		models.TABLE_OtcOrder,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_FinishTime,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Id)

	res, err := d.Orm.Raw(sql, OrderStatusCanceled, common.NowInt64MS(), OrderStatusPayed, id).Exec()

	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return
	}
	ok = true
	return
}

//超时取消订单
func (d *OrdersDao) Timeout(id uint64) (ok bool) {
	sql := fmt.Sprintf("update %s set %s=?,%s=? where %s=? And %s=?",
		models.TABLE_OtcOrder,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_FinishTime,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Id)

	res, err := d.Orm.Raw(sql, OrderStatusExpired, common.NowInt64MS(), OrderStatusCreated, id).Exec()

	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		common.LogFuncError("Timeout ORDER NO Affected:%v", id)
		return
	}
	ok = true
	return
}

//确认支付
func (d *OrdersDao) Pay(id uint64) (ok bool) {
	sql := fmt.Sprintf("update %s set %s=?,%s=? where %s=? And %s=?",
		models.TABLE_OtcOrder,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_PayTime,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Id)

	res, err := d.Orm.Raw(sql, OrderStatusPayed, common.NowInt64MS(), OrderStatusCreated, id).Exec()

	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return
	}
	ok = true
	return
}

func (d *OrdersDao) EditEPay2(order *models.OtcOrder) (ok bool) {
	n, err := d.Orm.Update(order, models.COLUMN_OtcOrder_EPayAccount, models.COLUMN_OtcOrder_EPayBank,
		models.COLUMN_OtcOrder_EPayBankBranch, models.COLUMN_OtcOrder_EPayId,
		models.COLUMN_OtcOrder_EPayName, models.COLUMN_OtcOrder_EPayType)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	if n == 0 {
		return
	}
	ok = true
	return
}

//确认支付
func (d *OrdersDao) ExchangerPay(id, pmid uint64, mType uint8, account, name, bank, bankBranch string) (ok bool) {
	sql := fmt.Sprintf("update %s set %s=?,%s=?,%s=?,%s=?,%s=?,%s=?,%s=?,%s=? where %s=? And %s=?",
		models.TABLE_OtcOrder,
		models.COLUMN_OtcOrder_EPayId, models.COLUMN_OtcOrder_EPayType,
		models.COLUMN_OtcOrder_EPayAccount, models.COLUMN_OtcOrder_EPayName,
		models.COLUMN_OtcOrder_EPayBank, models.COLUMN_OtcOrder_EPayBankBranch,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_PayTime,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Id)

	res, err := d.Orm.Raw(sql, pmid, mType, account, name, bank, bankBranch, OrderStatusPayed, common.NowInt64MS(),
		OrderStatusCreated, id).Exec()

	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return
	}
	ok = true
	return
}

// 确认收款
func (d *OrdersDao) Confirmed(id uint64, ip string) (ok bool) {
	sql := fmt.Sprintf("update %s set %s=?,%s=?,%s=? where %s=? And %s=?",
		models.TABLE_OtcOrder,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Eip, models.COLUMN_OtcOrder_Utime,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Id)

	res, err := d.Orm.Raw(sql, OrderStatusConfirmed, ip, common.NowInt64MS(), OrderStatusPayed, id).Exec()

	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return
	}
	ok = true
	return
}

// 申诉状态修改
func (d *OrdersDao) AppealChangeStatus(id uint64, ip string, from int8, to int8) (ok bool) {
	sql := fmt.Sprintf("update %s set %s=?,%s=?,%s=? where %s=? And %s=?",
		models.TABLE_OtcOrder,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Eip, models.COLUMN_OtcOrder_Utime,
		models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Id)

	res, err := d.Orm.Raw(sql, to, ip, common.NowInt64MS(), from, id).Exec()
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return
	}
	ok = true
	return
}

func (d *OrdersDao) Count(uid uint64, side int8, status []interface{}, appealStatus []interface{}) (num int64, err error) {
	qs := d.Orm.QueryTable(models.TABLE_OtcOrder)
	cond := orm.NewCondition()
	if uid > 0 {
		cond = cond.And(models.COLUMN_OtcOrder_Uid, uid)
	}
	if side == SideBuy || side == SideSell {
		cond = cond.And(models.COLUMN_OtcOrder_Side, side)
	}
	if len(status) > 0 {
		cond = cond.And(models.COLUMN_OtcOrder_Status+"__in", status...)
	}
	if len(appealStatus) > 0 {
		cond = cond.And(models.COLUMN_OtcOrder_AppealStatus+"__in", appealStatus...)
	}
	qs = qs.SetCond(cond)
	num, err = qs.Count()
	return
}

//
func (d *OrdersDao) FetchByUid(uid uint64, side int8, status []interface{}, appealStatusList []interface{}, offset int64, limit int64) (list []*models.OtcOrder) {
	qs := d.Orm.QueryTable(models.TABLE_OtcOrder)
	cond := orm.NewCondition()
	if uid > 0 {
		cond = cond.And(models.COLUMN_OtcOrder_Uid, uid)
	}
	if side == SideBuy || side == SideSell {
		cond = cond.And(models.COLUMN_OtcOrder_Side, side)
	}
	if len(status) > 0 {
		cond = cond.And(models.COLUMN_OtcOrder_Status+"__in", status...)
	}
	if len(appealStatusList) > 0 {
		cond = cond.And(models.COLUMN_OtcOrder_AppealStatus+"__in", appealStatusList...)
	}
	qs = qs.SetCond(cond)
	qs = qs.OrderBy("-" + models.COLUMN_OtcOrder_Id)
	list = []*models.OtcOrder{}
	qs = qs.Limit(limit, offset)
	_, err := qs.All(&list)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	return
}

func (d *OrdersDao) FetchDetailByUid(uid uint64, side int8, status uint8, appealStatusList uint8, offset int, limit int) (total uint64, list []*OrderAck, err error) {
	var qbQuery orm.QueryBuilder
	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	var qbTotal orm.QueryBuilder
	qbTotal, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}

	qbQuery.Select("T1.*",
		"T2."+models.COLUMN_User_Mobile+" AS u_mobile ",
		"T3."+models.COLUMN_User_Mobile+" AS eu_mobile ").
		From("((").Select("*").From(models.TABLE_OtcOrder).Where("1=1")
	qbTotal.Select("COUNT(*)").From(models.TABLE_OtcOrder).Where("1=1")
	var param []interface{}
	if uid > 0 {
		qbQuery.And("(" + models.COLUMN_OtcOrder_Uid + "=?").Or(models.COLUMN_OtcOrder_EUid + "=?)")
		qbTotal.And("(" + models.COLUMN_OtcOrder_Uid + "=?").Or(models.COLUMN_OtcOrder_EUid + "=?)")
		param = append(param, uid)
		param = append(param, uid)
	}
	if side == SideBuy || side == SideSell {
		qbQuery.And(models.COLUMN_OtcOrder_Side + "=?")
		qbTotal.And(models.COLUMN_OtcOrder_Side + "=?")
		param = append(param, side)
	}
	if status > 0 {
		qbQuery.And(models.COLUMN_OtcOrder_Status + "=?")
		qbTotal.And(models.COLUMN_OtcOrder_Status + "=?")
		param = append(param, status)
	}
	if appealStatusList > 0 {
		qbQuery.And(models.COLUMN_OtcOrder_AppealStatus + "=?")
		qbTotal.And(models.COLUMN_OtcOrder_AppealStatus + "=?")
		param = append(param, appealStatusList)
	}

	err = d.Orm.Raw(qbTotal.String(), param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	qbQuery.OrderBy("-" + models.COLUMN_OtcOrder_Id).Limit(limit).Offset(offset)
	sqlQuery := qbQuery.String()
	sqlQuery = fmt.Sprintf("%s) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s LEFT JOIN %s AS T3 ON T1.%s=T3.%s)",
		sqlQuery, models.TABLE_User, models.COLUMN_OtcOrder_Uid, models.COLUMN_User_Uid, models.TABLE_User, models.COLUMN_OtcOrder_EUid, models.COLUMN_User_Uid)
	_, err = d.Orm.Raw(sqlQuery, param...).QueryRows(&list)
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

func (d *OrdersDao) CountExchange(uid uint64, side int8, status, payType uint8, date int32) (num int64, err error) {
	qs := d.Orm.QueryTable(models.TABLE_OtcOrder)
	qs = qs.Filter(models.COLUMN_OtcOrder_EUid, uid)
	if side == SideBuy || side == SideSell {
		qs = qs.Filter(models.COLUMN_OtcOrder_Side, side)
	}
	if status > OrderStatus {
		qs = qs.Filter(models.COLUMN_OtcOrder_Status, status)
	}
	if payType > 0 {
		qs = qs.Filter(models.COLUMN_OtcOrder_PayType, payType)
	}
	if date > 0 {
		qs = qs.Filter(models.COLUMN_OtcOrder_Date, date)
	}
	num, err = qs.Count()
	return
}

func (d *OrdersDao) FetchExchangeByUid(uid uint64, side int8, status, payType uint8, date int32, offset int64, limit int64) (list []*models.OtcOrder) {
	qs := d.Orm.QueryTable(models.TABLE_OtcOrder)
	qs = qs.Filter(models.COLUMN_OtcOrder_EUid, uid)
	if side == SideBuy || side == SideSell {
		qs = qs.Filter(models.COLUMN_OtcOrder_Side, side)
	}
	if status > OrderStatus {
		qs = qs.Filter(models.COLUMN_OtcOrder_Status, status)
	}
	if payType > 0 {
		qs = qs.Filter(models.COLUMN_OtcOrder_PayType, payType)
	}
	if date > 0 {
		qs = qs.Filter(models.COLUMN_OtcOrder_Date, date)
	}
	qs = qs.OrderBy("-" + models.COLUMN_OtcOrder_Id)
	list = []*models.OtcOrder{}
	qs = qs.Offset(offset)
	_, err := qs.All(&list)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	return
}

func (d *OrdersDao) GroupFinishByUidDates(uid uint64, side int8, payType uint8, dates []int32) (list []*OrderExchangeTotal, err error) {
	if len(dates) == 0 {
		err = errors.New("date is empty")
		common.LogFuncError("DBERR:%v", err)
		return
	}

	sql := fmt.Sprintf("Select `date`,sum(amount) as amount,sum(funds) as funds from %s Where %s=? and status=?",
		models.TABLE_OtcOrder, models.COLUMN_OtcOrder_EUid)
	params := []interface{}{uid, OrderStatusConfirmed}

	wen := "?"
	params = append(params, dates[0])
	for i := 1; i < len(dates); i++ {
		wen += ",?"
		params = append(params, dates[i])
	}

	sql += " and `date` in (" + wen + ")"

	if side == SideBuy || side == SideSell {
		sql += " and " + models.COLUMN_OtcOrder_PayType + "=?"
		params = append(params, side)
	}

	if payType > 0 {
		sql += " and " + models.COLUMN_OtcOrder_PayType + "=?"
		params = append(params, payType)
	}

	sql += " group by date"
	list = []*OrderExchangeTotal{}
	_, err = d.Orm.Raw(sql, params...).QueryRows(&list)

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

//订单申述状态
func (d *OrdersDao) SetOrderAppealStatus(oid, uid uint64, appealStatus int8, qrCode string) (err error) {
	order := &models.OtcOrder{
		Id: oid,
	}

	err = d.Orm.Read(order)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	if order.Uid != uid && order.EUid != uid {
		return errors.New("order id error")
	}

	cols := []string{
		models.COLUMN_OtcOrder_Utime,
		models.COLUMN_OtcOrder_AppealStatus,
	}
	if qrCode != "" {
		order.QrCode = qrCode
		cols = append(cols, models.COLUMN_OtcOrder_QrCode)
	}

	order.AppealStatus = appealStatus
	order.Utime = common.NowInt64MS()
	_, err = d.Orm.Update(order, cols...)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (d *OrdersDao) FetchCheck(id uint64) (list []*models.OtcOrder, err error) {
	qs := d.Orm.QueryTable(models.TABLE_OtcOrder)
	if id > 0 {
		qs = qs.Filter(models.COLUMN_OtcOrder_Id+"__gt", id).OrderBy(models.COLUMN_OtcOrder_Id).Limit(1000)
	} else {
		qs = qs.OrderBy("-" + models.COLUMN_OtcOrder_Id).Limit(1000)
	}

	list = []*models.OtcOrder{}
	_, err = qs.All(&list, models.COLUMN_OtcOrder_Id, models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Side,
		models.COLUMN_OtcOrder_AppealStatus, models.COLUMN_OtcOrder_Ctime)

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

func (d *OrdersDao) LastCreatedOrder() (list []*models.OtcOrder, err error) {
	last := &models.OtcOrder{}
	err = d.Orm.QueryTable(models.TABLE_OtcOrder).Limit(1, 1000).
		OrderBy("-"+models.COLUMN_OtcOrder_Id).One(last, models.COLUMN_OtcOrder_Id)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("error:%v", err)
		return
	}

	qs := d.Orm.QueryTable(models.TABLE_OtcOrder).Filter(models.COLUMN_OtcOrder_Status, OrderStatusCreated).
		Filter(models.COLUMN_OtcOrder_Id+"__gt", last.Id)

	list = []*models.OtcOrder{}

	_, err = qs.All(&list, models.COLUMN_OtcOrder_Id, models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Ctime, models.COLUMN_OtcOrder_Side)

	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	return
}

type ConfirmedOrderInfo struct {
}

// 获取承兑商未处理订单数量
func (d *OrdersDao) GetNumOrderStatusConfirmed(euid uint64) (num int, err error) {
	sql := fmt.Sprintf("select count(*) from %s where euid=? and status=?", models.TABLE_OtcOrder)
	_ = d.Orm.Raw(sql, euid, OrderStatusConfirmed).QueryRow(&num)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

type OrderTimeOutInfo struct {
	Id           uint64 `json:"id"`
	Side         int8   `json:"side"`
	AppealStatus int8   `json:"appeal_status"`
	Ctime        int64  `json:"ctime"`
}

// 获取超时订单
func (d *OrdersDao) GetOrderTimeOut(start, limit int) (orderInfoList []*OrderTimeOutInfo, err error) {
	sql := fmt.Sprintf("select %s,%s,%s,%s from otc_order where status=?", models.COLUMN_OtcOrder_Id, models.COLUMN_OtcOrder_Side,
		models.COLUMN_OtcOrder_AppealStatus, models.COLUMN_OtcOrder_Ctime)
	_, err = d.Orm.Raw(sql, OrderStatusCreated).QueryRows(&orderInfoList)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

// 获取承兑商Otc Order
func (d *OrdersDao) GetOtcOrder(euid uint64, page, limit int) (orderInfoList []*models.OtcOrder, err error) {
	sql := fmt.Sprintf("select * from %s where euid=? ORDER BY ctime DESC LIMIT ?,?", models.TABLE_OtcOrder)
	_, err = d.Orm.Raw(sql, euid, (page-1)*limit, limit).QueryRows(&orderInfoList)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

// 统计
func (d *OrdersDao) OtcStat(date string) (res map[string]uint32, err error) {
	// 找出时间内第一个订单 ctime没有索引，查询时间越久越慢
	sql := fmt.Sprintf("Select id From %s Where `date`=? ORDER BY `id` DESC LIMIT 1", models.TABLE_OtcOrder)

	orderFirst := &models.OtcOrder{}
	err = d.Orm.Raw(sql, date).QueryRow(orderFirst)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}

	// 找出时间内最后一个订单
	sql = fmt.Sprintf("Select id From %s Where `date`=? ORDER BY `id` ASC LIMIT 1", models.TABLE_OtcOrder)
	orderLast := &models.OtcOrder{}
	err = d.Orm.Raw(sql, date).QueryRow(orderLast)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}
	res = map[string]uint32{}
	//统计订单数
	qs := d.Orm.QueryTable(models.TABLE_OtcOrder).Filter("id__gte", orderFirst.Id).Filter("id__lte", orderLast.Id)

	//全部订单数
	tmp, err := qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}
	res["num_order"] = uint32(tmp)

	// 成交订单、金额、手续费
	sql = "Select side,count(id) as orders, sum(fee) as fee, sum(amount) as amount, sum(funds) as funds " +
		" From otc_order Where id>=? and id<=? and status=? group by side"

	type stat struct {
		Side   int32  `orm:"side"`
		Order  uint32 `orm:"orders"`
		Fee    uint32 `orm:"fee"`
		Amount uint32 `orm:"amount"`
		Funds  uint32 `orm:"funds"`
	}

	s := []*stat{}
	_, err = d.Orm.Raw(sql, orderFirst.Id, orderLast.Id, OrderStatusConfirmed).QueryRows(&s)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}
	for _, v := range s {
		if v.Side == 1 { //用户购买
			res["num_order_buy"] = v.Order
			res["num_fee_buy"] = v.Fee
			res["num_amount_buy"] = v.Amount
			res["num_funds_buy"] = v.Funds
		} else {
			res["num_order_sell"] = v.Order
			res["num_fee_sell"] = v.Fee
			res["num_amount_sell"] = v.Amount
			res["num_funds_sell"] = v.Funds
		}
	}
	return
}

// 统计用户订单信息
func (d *OrdersDao) OtcStatPeople(uid uint64) (res map[string]uint32) {
	// 成交订单、金额、手续费
	sql := "Select side,count(id) as orders, sum(fee) as fee, sum(amount) as amount, sum(funds) as funds " +
		" From otc_order Where  `uid`=? and status=? group by side"

	type stat struct {
		Side   int32  `orm:"side"`
		Order  uint32 `orm:"orders"`
		Fee    uint32 `orm:"fee"`
		Amount uint32 `orm:"amount"`
		Funds  uint32 `orm:"funds"`
	}

	s := []*stat{}
	_, err := d.Orm.Raw(sql, uid, OrderStatusConfirmed).QueryRows(&s)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}
	res = map[string]uint32{}
	for _, v := range s {
		if v.Side == 1 { //用户购买
			res["num_order_buy"] = v.Order
			res["num_fee_buy"] = v.Fee
			res["num_amount_buy"] = v.Amount
			res["num_funds_buy"] = v.Funds
		} else {
			res["num_order_sell"] = v.Order
			res["num_fee_sell"] = v.Fee
			res["num_amount_sell"] = v.Amount
			res["num_funds_sell"] = v.Funds
		}
	}
	return
}

//统计用户Eusd出售和购买
func (d *OrdersDao) StatPeopleByDateRange(uid uint64, startTime, endTime int64) (res map[string]uint32) {
	querySql := fmt.Sprintf("select %s as side,sum(%s) as amount from %s where %s=? and %s=? and %s>=? and %s<=? group by %s",
		models.COLUMN_OtcOrder_Side, models.COLUMN_OtcOrder_Amount, models.TABLE_OtcOrder,
		models.COLUMN_OtcOrder_Uid, models.COLUMN_OtcOrder_Status, models.COLUMN_OtcOrder_Ctime, models.COLUMN_OtcOrder_Ctime, models.COLUMN_OtcOrder_Side)
	type stat struct {
		Side   uint8  `json:"side"`
		Amount uint32 `json:"amount"`
	}

	list := []*stat{}
	_, err := d.Orm.Raw(querySql, uid, OrderStatusConfirmed, startTime, endTime).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	res = map[string]uint32{
		"eusd_buy":  0,
		"eusd_sell": 0,
	}

	for _, v := range list {
		if v.Side == 1 {
			res["eusd_buy"] = v.Amount
		} else {
			res["eusd_sell"] = v.Amount
		}
	}

	return
}

func (d *OrdersDao) GetOrdersByDate(date int32, page, limit int) (orderInfoList []*models.OtcOrder, err error) {
	querySql := fmt.Sprintf("select * from %s where date=? ORDER BY date DESC LIMIT ?,?", models.TABLE_OtcOrder)

	_, err = d.Orm.Raw(querySql, date, (page-1)*limit, limit).QueryRows(&orderInfoList)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

func (d *OrdersDao) GetOrdersNumByDate(date int32) (total int, err error) {
	var qbTotal orm.QueryBuilder
	qbTotal, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}

	qbTotal.Select("COUNT(*)").From(models.TABLE_OtcOrder).Where("" + models.COLUMN_OtcOrder_Date + "=?")
	var param []interface{}
	param = append(param, date)

	err = d.Orm.Raw(qbTotal.String(), param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	return
}

//获取承兑商未处理订单数
func (d *OrdersDao) GetExchangerUnDealOrders(uid uint64) (sell, buy int, err error) {
	var qbTotal orm.QueryBuilder
	qbTotal, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}

	qbTotal.Select("side,id").From(models.TABLE_OtcOrder).
		Where(fmt.Sprintf("%s=? And %s<%v", models.COLUMN_OtcOrder_EUid, models.COLUMN_OtcOrder_Status, OrderStatusConfirmed))
	list := []*models.OtcOrder{}
	_, err = d.Orm.Raw(qbTotal.String(), uid).QueryRows(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}
	for _, v := range list {
		if v.Side == SideBuy {
			buy += 1
		} else {
			sell += 1
		}
	}

	return
}

//
func (d *OrdersDao) QueryDataByTime(t int64) (orders map[uint64][]models.OtcOrder, err error) {
	var orderUids []models.OtcOrder
	_, err = d.Orm.QueryTable(models.TABLE_OtcOrder).Filter(fmt.Sprintf("%s__gte", models.COLUMN_OtcOrder_Ctime), t).GroupBy(models.COLUMN_OtcOrder_Uid).All(&orderUids, "uid")
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	if len(orderUids) == 0 {
		return
	}

	var tmpOrders []models.OtcOrder
	orders = make(map[uint64][]models.OtcOrder, len(orderUids))
	for _, v := range orderUids {
		_, err = d.Orm.QueryTable(models.TABLE_OtcOrder).Filter(fmt.Sprintf("%s__gte", models.COLUMN_OtcOrder_Ctime), t).
			Filter(fmt.Sprintf("%s__exact", models.COLUMN_OtcOrder_Uid), v.Uid).
			OrderBy("-"+models.COLUMN_OtcOrder_PayType, "-"+models.COLUMN_OtcOrder_Ctime).
			All(&tmpOrders, models.COLUMN_OtcOrder_Id, models.COLUMN_OtcOrder_Uid,
				models.COLUMN_OtcOrder_Amount, models.COLUMN_OtcOrder_Funds,
				models.COLUMN_OtcOrder_PayType, models.COLUMN_OtcOrder_PayAccount, models.COLUMN_OtcOrder_Ctime)
		if err != nil {
			common.LogFuncError("mysql error:%v", err)
			return
		}
		orders[v.Uid] = tmpOrders
	}

	return
}
