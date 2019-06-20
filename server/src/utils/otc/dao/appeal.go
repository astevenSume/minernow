package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/otc/models"
)

//申诉类型
const (
	AppealTypeNil             = iota
	AppealTypeNoPay           //买家未付款
	AppealTypeNoConfirm       //卖家未确认
	AppealTypeNoResponse      //长时间无回应
	AppealTypeFraud           //涉嫌诈骗
	AppealTypeMoneyLaundering //涉嫌洗钱
	AppealTypeOther           //其它
	AppealTypeMax
)

//处理状态
const (
	AppealStatusNil = iota
	AppealStatusPending
	AppealStatusProcessing
	AppealStatusResolved
	AppealStatusMax
)

type AppealDao struct {
	common.BaseDao
}

func NewAppealDao(db string) *AppealDao {
	return &AppealDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AppealDaoEntity *AppealDao

type AppealAck struct {
	Id            string `json:"id"`
	Type          int8   `json:"type"`
	UserId        string `json:"user_id"`
	AdminId       uint32 `json:"admin_id"`
	OrderId       string `json:"order_id"`
	Context       string `json:"context"`
	ContactWechat string `json:"contact_wechat"`
	QrCode        string `json:"qr_code"`
	Status        int8   `json:"status"`
	Ctime         int64  `json:"ctime"`
	Utime         int64  `json:"utime"`
	Mobile        string `json:"mobile"`
	Wechat        string `json:"wechat"`
}

func (d *AppealDao) ClientAppeal(data *models.Appeal) (ack AppealAck) {
	if data == nil {
		return
	}

	ack.Id = fmt.Sprintf("%v", data.Id)
	ack.UserId = fmt.Sprintf("%v", data.UserId)
	ack.OrderId = fmt.Sprintf("%v", data.OrderId)
	ack.AdminId = data.AdminId
	ack.Type = data.Type
	ack.Context = data.Context
	ack.Status = data.Status
	ack.Ctime = data.Ctime
	ack.Utime = data.Utime
	ack.Wechat = data.Wechat
	var err error
	ack.Mobile, err = UserDaoEntity.GetMobileByUid(data.UserId)
	if err != nil {
		common.LogFuncError("error:%v", err)
	}
	return
}

/*func (d *AppealDao) ClientAppeals(datas []AppealMobile) (acks []AppealAck) {
	for _, data := range datas {
		var ack AppealAck
		ack.Id = fmt.Sprintf("%v", data.Id)
		ack.UserId = fmt.Sprintf("%v", data.UserId)
		ack.OrderId = fmt.Sprintf("%v", data.OrderId)
		ack.AdminId = data.AdminId
		ack.Type = data.Type
		ack.Context = data.Context
		ack.Status = data.Status
		ack.Ctime = data.Ctime
		ack.Utime = data.Utime
		ack.Mobile = data.Mobile
		acks = append(acks, ack)
	}

	return
}*/

//获取申诉信息
func (d *AppealDao) QueryById(id uint64) (appeal models.Appeal, err error) {
	appeal = models.Appeal{
		Id: id,
	}

	err = d.Orm.Read(&appeal, models.COLUMN_Appeal_Id)
	if err != nil {
		if err == orm.ErrNoRows {
			appeal.Id = 0
			err = nil
			return
		}
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

//获取申诉信息
func (d *AppealDao) QueryByOrderId(oid uint64) (appeal models.Appeal, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_Appeal).Filter(models.COLUMN_Appeal_OrderId, oid).
		Filter(models.COLUMN_Appeal_Status+"__in", AppealStatusPending, AppealStatusProcessing).All(&appeal)
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

type AppealMobile struct {
	Id      string `orm:"column(id);pk" json:"id,omitempty"`
	Type    int8   `orm:"column(type)" json:"type,omitempty"`
	UserId  string `orm:"column(user_id)" json:"user_id,omitempty"`
	AdminId uint32 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	OrderId string `orm:"column(order_id)" json:"order_id,omitempty"`
	Context string `orm:"column(context);size(256)" json:"context,omitempty"`
	Status  int8   `orm:"column(status)" json:"status,omitempty"`
	Ctime   int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Utime   int64  `orm:"column(utime)" json:"utime,omitempty"`
	Mobile  string `orm:"column(mobile)" json:"mobile,omitempty"`
	Wechat  string `orm:"column(wechat)" json:"wechat,omitempty"`
}

//分页条件查询
func (d *AppealDao) QueryPageCondition(uid uint64, typeId int8, status int8, page int, perPage int) (total int, appeals []AppealMobile, err error) {
	var qbQuery orm.QueryBuilder
	var qbTotal orm.QueryBuilder
	qbTotal, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	// 构建查询对象
	qbTotal.Select("Count(*) AS total").
		From(models.TABLE_Appeal).
		Where("1 = 1")
	qbQuery.Select("T1.*",
		"T2."+models.COLUMN_User_Mobile).
		From("((").Select("*").From(models.TABLE_Appeal).Where("1=1")
	var param []interface{}

	if uid > 0 {
		qbQuery.And(models.COLUMN_Appeal_UserId + "= ?")
		qbTotal.And(models.COLUMN_Appeal_UserId + "= ?")
		param = append(param, uid)
	}
	if typeId > AppealTypeNil && typeId < AppealTypeMax {
		qbQuery.And(models.COLUMN_Appeal_Type + "= ?")
		qbTotal.And(models.COLUMN_Appeal_Type + "= ?")
		param = append(param, typeId)
	}
	if status > AppealStatusNil && status < AppealStatusMax {
		qbQuery.And(models.COLUMN_Appeal_Status + "= ?")
		qbTotal.And(models.COLUMN_Appeal_Status + "= ?")
		param = append(param, status)
	}

	sqlTotal := qbTotal.String()
	err = d.Orm.Raw(sqlTotal, param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error: %v", err)
		return
	}

	qbQuery.OrderBy("-" + models.COLUMN_Appeal_Id)
	qbQuery.Limit(perPage).Offset((page - 1) * perPage)
	sqlQuery := qbQuery.String()
	sqlQuery = fmt.Sprintf("%s) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s)", sqlQuery, models.TABLE_User, models.COLUMN_Appeal_UserId, models.COLUMN_User_Uid)
	_, err = d.Orm.Raw(sqlQuery, param...).QueryRows(&appeals)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *AppealDao) Create(typeId int8, adminID uint32, orderID uint64, uId uint64, weChat, context string) (appeal models.Appeal, err error) {
	appeal = models.Appeal{
		UserId:  uId,
		Type:    typeId,
		Status:  AppealStatusPending,
		OrderId: orderID,
		AdminId: adminID,
		Context: context,
		Wechat:  weChat,
		Ctime:   common.NowInt64MS(),
		Utime:   common.NowInt64MS(),
	}

	var id int64
	id, err = d.Orm.Insert(&appeal)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}
	appeal.Id = uint64(id)

	return
}

func (d *AppealDao) Update(id uint64, adminID uint32, sType int8, status int8, context string) (appeal models.Appeal, err error) {
	appeal = models.Appeal{Id: id}
	err = d.Orm.Read(&appeal, models.COLUMN_Appeal_Id)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}

	appeal.Type = sType
	appeal.Status = status
	appeal.AdminId = adminID
	appeal.Context = context
	appeal.Utime = common.NowInt64MS()
	_, err = d.Orm.Update(&appeal, models.COLUMN_Appeal_Status, models.COLUMN_Appeal_Utime, models.COLUMN_Appeal_AdminId, models.COLUMN_Appeal_Context, models.COLUMN_Appeal_Type)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

//解决申诉
func (d *AppealDao) Resolve(orderId, uid uint64) (appeal models.Appeal, err error) {
	appeal, err = d.QueryByOrderId(orderId)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}
	if appeal.Id == 0 {
		err = orm.ErrNoRows
		return
	}

	appeal.Status = AppealStatusResolved
	appeal.Utime = common.NowInt64MS()
	_, err = d.Orm.Update(&appeal, models.COLUMN_Appeal_Status, models.COLUMN_Appeal_Utime)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	err = OrdersDaoEntity.SetOrderAppealStatus(orderId, uid, OrderAppealStatusResolved, "")
	if err != nil {
		return
	}

	return
}

//解决申诉
func (d *AppealDao) ResolveById(id, uid uint64) (appeal models.Appeal, err error) {
	appeal = models.Appeal{Id: id}
	appeal.Status = AppealStatusResolved
	appeal.Utime = common.NowInt64MS()
	_, err = d.Orm.Update(&appeal, models.COLUMN_Appeal_Status, models.COLUMN_Appeal_Utime)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	err = OrdersDaoEntity.SetOrderAppealStatus(appeal.OrderId, uid, OrderAppealStatusResolved, "")
	if err != nil {
		return
	}

	return
}
