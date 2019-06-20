package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
	otcDao "utils/otc/dao"
)

type AppealServiceController struct {
	BaseController
}

//获取申述客服
func (c *AppealServiceController) Get() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadAppealService, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAppealService, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	var input string
	if id > 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := dao.AppealServiceDaoEntity.QueryById(id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppealService, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadAppealService, input, data)
	} else {
		status, err := c.GetInt8(KEY_STATUS, 0)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppealService, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		page, err := c.GetInt(KEY_PAGE, 1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppealService, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		limit, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppealService, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		input = fmt.Sprintf("{\"status\":%v,\"page\":%v,\"limit\":%v}", status, page, limit)

		total, data, err := dao.AppealServiceDaoEntity.QueryByPage(status, page, limit)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAppealService, controllers.ERROR_CODE_DB, input)
			return
		}

		res := map[string]interface{}{}
		meta := dao.PageInfo{
			Limit: limit,
			Total: int(total),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = data
		c.SuccessResponseAndLog(OPActionReadAppealService, input, res)
	}
}

func (c *AppealServiceController) Create() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAppealService, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		AdminId uint32 `json:"admin_id"`
		Wechat  string `json:"wechat"`
		QrCode  string `json:"qr_code"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddAppealService, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	err = dao.AppealServiceDaoEntity.Create(req.AdminId, req.Wechat, req.QrCode)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionAddAppealService, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionAddAppealService, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionAddAppealService, string(c.Ctx.Input.RequestBody))
}

//更新
func (c *AppealServiceController) Update() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditAppealService, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id     uint32 `json:"id"`
		Wechat string `json:"wechat"`
		QrCode string `json:"qr_code"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditAppealService, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	err = dao.AppealServiceDaoEntity.Update(req.Id, req.Wechat, req.QrCode)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditAppealService, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionEditAppealService, string(c.Ctx.Input.RequestBody))
}

//删除
func (c *AppealServiceController) Del() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelAppealService, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionDelAppealService, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)

	if dao.AppealServiceDaoEntity.IsLast(id) {
		c.ErrorResponseAndLog(OPActionDelAppealService, controllers.ERROR_CODE_APPEAL_SERVICE_LESS, input)
		return
	}

	err = dao.AppealServiceDaoEntity.DelById(id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAppealService, controllers.ERROR_CODE_DB, input)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionDelAppealService, input)
}

func (c *AppealServiceController) Suspend() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelAppealServiceSuspend, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionDelAppealServiceSuspend, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)

	if dao.AppealServiceDaoEntity.IsLast(id) {
		c.ErrorResponseAndLog(OPActionDelAppealServiceSuspend, controllers.ERROR_CODE_APPEAL_SERVICE_LESS, input)
		return
	}

	err = dao.AppealServiceDaoEntity.SetStatus(id, dao.AppealServiceStatusSuspend)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAppealServiceSuspend, controllers.ERROR_CODE_DB, input)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionDelAppealServiceSuspend, input)
}

func (c *AppealServiceController) Restore() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAppealServiceRestore, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint32(KEY_ID, 0)
	if err != nil || id <= 0 {
		c.ErrorResponseAndLog(OPActionAddAppealServiceRestore, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"id\":%v}", id)
	err = dao.AppealServiceDaoEntity.SetStatus(id, dao.AppealServiceStatusRestore)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddAppealServiceRestore, controllers.ERROR_CODE_DB, input)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddAppealServiceRestore, input)
}

//easyjson:json
type AppealDealLog struct {
	Id       string `orm:"column(id);pk" json:"id,omitempty"`
	AppealId string `orm:"column(appeal_id)" json:"appeal_id,omitempty"`
	OrderId  string `orm:"column(order_id)" json:"order_id,omitempty"`
	AdminId  uint32 `orm:"column(admin_id)" json:"admin_id,omitempty"`
	Action   int8   `orm:"column(action)" json:"action,omitempty"`
	Ctime    int64  `orm:"column(ctime)" json:"ctime,omitempty"`
}

//获取申述处理记录
func (c *AppealServiceController) Record() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadAppealServiceRecord, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	appealId, err := c.GetUint64("appeal_id", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAppealServiceRecord, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}
	adminId, err := c.GetUint32("admin_id", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAppealServiceRecord, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}
	orderId, err := c.GetUint64("order_id", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAppealServiceRecord, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}
	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAppealServiceRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	limit, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAppealServiceRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input = fmt.Sprintf("{\"appeal_id\":%v,\"admin_id\":%v,\"page\":%v,\"limit\":%v}", appealId, adminId, page, limit)

	total, data, err := otcDao.AppealDealLogDaoEntity.QueryByPage(appealId, orderId, adminId, page, limit)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAppealServiceRecord, controllers.ERROR_CODE_DB, input)
		return
	}
	var logs []AppealDealLog
	for _, item := range data {
		logs = append(logs, AppealDealLog{
			Id:       fmt.Sprintf("%v", item.Id),
			AppealId: fmt.Sprintf("%v", item.AppealId),
			OrderId:  fmt.Sprintf("%v", item.OrderId),
			AdminId:  item.AdminId,
			Action:   item.Action,
			Ctime:    item.Ctime,
		})
	}

	res := map[string]interface{}{}
	meta := PageInfo{
		Limit: limit,
		Total: int(total),
		Page:  page,
	}
	res["meta"] = meta
	res["list"] = logs
	c.SuccessResponseAndLog(OPActionReadAppealServiceRecord, input, res)
}
