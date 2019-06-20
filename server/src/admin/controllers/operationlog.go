package controllers

import (
	"admin/controllers/errcode"
	"common"
	"fmt"
	"strconv"
	"utils/admin/dao"
)

type OperationLogController struct {
	BaseController
}

func (c *OperationLogController) GetPageInfo() (uint64, int, int, error) {
	//请求参数
	var id uint64
	strId := c.GetString(KEY_ID)
	if len(strId) > 0 {
		var err error
		id, err = strconv.ParseUint(strId, 10, 64)
		if err != nil {
			return 0, 0, 0, err
		}
	}
	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		common.LogFuncError("err:%v", err)
		return 0, 0, 0, err
	}
	perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		common.LogFuncError("err:%v", err)
		return 0, 0, 0, err
	}

	return id, page, perPage, nil
}

//获取管理员操作日志
func (c *OperationLogController) HandleGetAdminOperationLog() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadOperation, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	_, page, perPage, err := c.GetPageInfo()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOperation, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	adminId, err := c.GetUint64("admin_id", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOperation, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	action, err := c.GetInt32(KEY_ACTION, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOperation, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"action\":%v,\"page\":%v,\"limit\":%v}", action, page, perPage)
	//获取数据
	total, logs, err := dao.OperationLogDaoEntity.QueryPageOperationLog(adminId, page, perPage, action)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadOperation, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回c查询结果
	res := map[string]interface{}{}
	meta := PageInfo{
		Limit: perPage,
		Total: total,
		Page:  page,
	}
	res["meta"] = meta
	res["list"] = logs

	c.SuccessResponseAndLog(OPActionReadOperation, input, res)
}
