package dao

import (
	"common"
	"fmt"
	"utils/otc/models"
)

type SystemNotification struct {
	common.BaseDao
}

func NewSystemNotificationdDao(db string) *SystemNotification {
	return &SystemNotification{
		BaseDao: common.NewBaseDao(db),
	}
}

var SystemNotificationdDaoEntity *SystemNotification

//新增系统通知
func (d *SystemNotification) InsertSystemNotification(notification_type string, content string, uid uint64) (sn *models.SystemNotification, err error) {
	var nid uint64
	nid, err = common.IdManagerGen(IdTypeSystemNotification)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	sn = &models.SystemNotification{
		Nid:              nid,
		NotificationType: notification_type,
		Content:          content,
		Uid:              uid,
		IsRead:           0,
		Ctime:            common.NowInt64MS(),
	}
	_, err = d.Orm.Insert(sn)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//获取系统通知
func (d *SystemNotification) GetSystemNotificationByUid(uid uint64, page int, limit int) (snList []*models.SystemNotification, err error) {

	if uid <= 0 {
		common.LogFuncCritical("GetSystemNotificationByUid uid <= 0")
		return
	}
	sql := fmt.Sprintf("SELECT * FROM system_notification WHERE uid = ? ORDER BY ctime DESC LIMIT ?,?")
	_, err = d.Orm.Raw(sql, uid, (page-1)*limit, page*limit).QueryRows(&snList)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//修改为已读
func (d *SystemNotification) UpdateSystemNotificationByNid(NidListStr string) (err error) {

	sql := fmt.Sprintf("update system_notification set is_read =1 where nid in (%s)", NidListStr)
	_, err = d.Orm.Raw(sql).Exec()

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//获取未读条数通知
func (d *SystemNotification) GetSystemNotificationIsRead(uid uint64) (num int, err error) {

	sql := fmt.Sprintf("select count(*) from system_notification where is_read = 0 and uid =?")
	_ = d.Orm.Raw(sql, uid).QueryRow(&num)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}
