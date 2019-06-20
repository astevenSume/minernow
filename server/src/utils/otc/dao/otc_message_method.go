package dao

import (
	"common"
	"fmt"
	"utils/otc/models"
)

type MessageMethodDao struct {
	common.BaseDao
}

func NewMessageMethodDao(db string) *MessageMethodDao {
	return &MessageMethodDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var MessageMethodDaoEntity *MessageMethodDao

// add new message
func (d *MessageMethodDao) Add(order_id int64, uid uint64, msg_type string, content string) (mm *models.OtcMsg, err error) {

	var id uint64
	id, err = common.IdManagerGen(IdTypeMessageMethod)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	mm = &models.OtcMsg{
		Id:          id,
		OrderId:     order_id,
		Uid:         uid,
		Content:     content,
		IsRead:      0,
		MessageType: msg_type,
		Ctime:       common.NowInt64MS(),
	}
	_, err = d.Orm.Insert(mm)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

type OtcMsg struct {
	Id          string `orm:"column(id);pk" json:"id,omitempty"`
	OrderId     string `orm:"column(order_id)" json:"order_id,omitempty"`
	Uid         string `orm:"column(uid)" json:"uid,omitempty"`
	Content     string `orm:"column(content)" json:"content,omitempty"`
	IsRead      uint8  `orm:"column(is_read)" json:"is_read,omitempty"`
	MessageType string `orm:"column(msg_type);size(200)" json:"msg_type,omitempty"`
	Ctime       int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

// get msg by order_id
func (d *MessageMethodDao) QueryByOrderId(order_id int64, page int, limit int) (list []*OtcMsg, total int64, err error) {
	querySeter := d.Orm.QueryTable(models.TABLE_OtcMsg)
	if order_id > 0 {
		querySeter = querySeter.Filter(models.ATTRIBUTE_OtcMsg_OrderId, order_id).OrderBy("ctime")
	}
	//querySeter = querySeter.Filter(models.ATTRIBUTE_OtcMsg_OrderId, order_id)
	total, err = querySeter.Count()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	var otcMsg []*models.OtcMsg
	_, err = querySeter.Limit(limit, (page-1)*limit).
		All(&otcMsg, models.COLUMN_OtcMsg_Id,
			models.COLUMN_OtcMsg_OrderId,
			models.COLUMN_OtcMsg_Uid,
			models.COLUMN_OtcMsg_Content,
			models.COLUMN_OtcMsg_IsRead,
			models.COLUMN_OtcMsg_MessageType,
			models.COLUMN_OtcMsg_Ctime)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	for _, item := range otcMsg {
		msg := &OtcMsg{
			Id:          fmt.Sprintf("%v", item.Id),
			OrderId:     fmt.Sprintf("%v", item.OrderId),
			Uid:         fmt.Sprintf("%v", item.Uid),
			Content:     item.Content,
			IsRead:      item.IsRead,
			MessageType: item.MessageType,
			Ctime:       item.Ctime,
		}
		list = append(list, msg)
	}
	return
}

// get msg by id = 0
func (d *MessageMethodDao) QueryByOrderIdZero(order_id int64) (list []*models.OtcMsg, err error) {
	sql := fmt.Sprintf("select * from otc_msg WHERE order_id = ? ORDER BY id DESC LIMIT 0,10")
	_, err = d.Orm.Raw(sql, order_id).QueryRows(&list)
	if list == nil {
		list = []*models.OtcMsg{}
	}
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//id > 0 direction == "up"
func (d *MessageMethodDao) QueryUpMsg(id uint64, order_id int64) (list []*models.OtcMsg, err error) {
	sql := fmt.Sprintf("select * from otc_msg WHERE order_id = ? and id < ? ORDER BY id DESC LIMIT 0,10")
	_, err = d.Orm.Raw(sql, order_id, id).QueryRows(&list)
	if list == nil {
		list = []*models.OtcMsg{}
	}
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//id > 0 direction == "up" all = 1
func (d *MessageMethodDao) QueryUpMsgAll(id uint64, order_id int64) (list []*models.OtcMsg, err error) {
	sql := fmt.Sprintf("select * from otc_msg WHERE order_id = ? and id < ? ORDER BY id DESC LIMIT 0,100")
	_, err = d.Orm.Raw(sql, order_id, id).QueryRows(&list)

	if list == nil {
		list = []*models.OtcMsg{}
	}

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//id > 0 direction == "down"
func (d *MessageMethodDao) QueryDownMsg(id uint64, order_id int64) (list []*models.OtcMsg, err error) {
	sql := fmt.Sprintf("select * from otc_msg WHERE order_id = ? and id > ? ORDER BY id ASC LIMIT 0,10")
	_, err = d.Orm.Raw(sql, order_id, id).QueryRows(&list)

	if list == nil {
		list = []*models.OtcMsg{}
	}

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}
