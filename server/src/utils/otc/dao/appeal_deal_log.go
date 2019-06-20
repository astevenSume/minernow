package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/otc/models"
)

//承兑商审核状态
const (
	AppealDealLogActionNil          = iota
	AppealDealLogActionResolve      //解决订单
	AppealDealLogActionCancel       //取消订单
	AppealDealLogActionConfirmPay   //确认已付款
	AppealDealLogActionConfirm      //确认放币
	AppealDealLogActionFrozenSell   //冻结卖家
	AppealDealLogActionUnFrozenSell //解冻卖家
	AppealDealLogActionForbidSell   //禁止卖家
	AppealDealLogActionActionSell   //激活卖家
	AppealDealLogActionForbidBuy    //禁止买家
	AppealDealLogActionActionBuy    //激活买家
	AppealDealLogActionMax
)

type AppealDealLogDao struct {
	common.BaseDao
}

func NewAppealDealLogDao(db string) *AppealDealLogDao {
	return &AppealDealLogDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AppealDealLogDaoEntity *AppealDealLogDao

//申述处理记录
func (d *AppealDealLogDao) Create(action int8, adminId uint32, appealId uint64, orderId uint64) (data models.AppealDealLog, err error) {
	data = models.AppealDealLog{
		Action:   action,
		AdminId:  adminId,
		AppealId: appealId,
		OrderId:  orderId,
		Ctime:    common.NowInt64MS(),
	}

	id, err := d.Orm.Insert(&data)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
		return
	}
	data.Id = uint64(id)

	return
}

//申述处理记录
func (d *AppealDealLogDao) QueryByPage(appealId uint64, orderId uint64, adminId uint32, page int, limit int) (total int64, data []models.AppealDealLog, err error) {
	qs := d.Orm.QueryTable(models.TABLE_AppealDealLog)
	if qs == nil {
		common.LogFuncError("QueryTable fail")
		return
	}

	if appealId > 0 {
		qs = qs.Filter(models.COLUMN_AppealDealLog_AppealId, appealId)
	}
	if orderId > 0 {
		qs = qs.Filter(models.COLUMN_AppealDealLog_OrderId, orderId)
	}
	if adminId > 0 {
		qs = qs.Filter(models.COLUMN_AppealDealLog_AdminId, adminId)
	}

	total, err = qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	//获取当前页数据的起始位置
	start := (page - 1) * limit
	if int64(start) > total {
		err = nil
		return
	}
	_, err = qs.OrderBy("-"+models.COLUMN_AppealDealLog_Id).Limit(limit, start).All(&data)
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
