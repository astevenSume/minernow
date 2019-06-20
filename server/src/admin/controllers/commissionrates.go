package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"utils/admin/dao"
	"utils/admin/models"
)

type CommissionRateController struct {
	BaseController
}

//获取返佣等级配置
func (c *CommissionRateController) GetCommRates() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadCommissionRate, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt64(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadCommissionRate, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	res := map[string]interface{}{}
	var input string
	//数据获取
	if id == 0 {
		input = fmt.Sprintf("{\"id\":%v}", id)
		data, err := dao.CommissionRateDaoEntity.QueryPageCommRate()
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadCommissionRate, controllers.ERROR_CODE_DB, input)
			return
		}
		res["list"] = data
	} else {
		input = fmt.Sprintf("{\"id\":%v}", id)
		commissionrates := &models.Commissionrates{Id: int64(id)}
		err := dao.CommissionRateDaoEntity.QueryCommRate(commissionrates, models.COLUMN_Permission_Id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadCommissionRate, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadCommissionRate, input, *commissionrates)
		return
	}

	c.SuccessResponseAndLog(OPActionReadCommissionRate, input, res)
}

//新增返佣等级配置
func (c *CommissionRateController) EditCommRates() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditCommissionRate, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Data dao.Commissioncfgs `json:"data"`
	}

	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditCommissionRate, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//插入数据
	err = dao.CommissionRateDaoEntity.EditCommRate(req.Data)
	if err != nil {
		if err == dao.ErrParam {
			c.ErrorResponseAndLog(OPActionEditCommissionRate, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		} else {
			c.ErrorResponseAndLog(OPActionEditCommissionRate, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		}
		return
	}

	//返回c查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionEditCommissionRate, string(c.Ctx.Input.RequestBody))
}

/*//新增返佣等级配置
func (c *CommissionRateController) AddCommRates() {
	_, uid, errCode := c.PrePare()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(errCode, uid, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Min        uint64 `json:"min"`
		Max        uint64 `json:"max"`
		Commission uint64 `json:"commission"`
	}

	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(controllers.ERROR_CODE_PARAMS_ERROR, uid, string(c.Ctx.Input.RequestBody))
		return
	}

	//插入数据
	commissionrates := new(models.Commissionrates)
	commissionrates.Min = req.Min
	commissionrates.Max = req.Max
	commissionrates.Commission = req.Commission
	commissionrates.Ctime = common.NowInt64MS()
	commissionrates.Utime = common.NowInt64MS()
	err = dao.CommissionRateDaoEntity.AddCommRate(commissionrates)
	if err != nil {
		c.ErrorResponseAndLog(controllers.ERROR_CODE_DB, uid, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回c查询结果
	c.SuccessResponseAndLog(uid, string(c.Ctx.Input.RequestBody), *commissionrates)
}

//更新返佣等级配置
func (c *CommissionRateController) UpdateCommRates() {
	_, uid, errCode := c.PrePare()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(errCode, uid, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id         int64  `json:"id"`
		Min        uint64 `json:"min"`
		Max        uint64 `json:"max"`
		Commission uint64 `json:"commission"`
	}

	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(controllers.ERROR_CODE_PARAMS_ERROR, uid, string(c.Ctx.Input.RequestBody))
		return
	}

	commissionrates := &models.Commissionrates{
		Id:         req.Id,
		Min:        req.Min,
		Max:        req.Max,
		Commission: req.Commission,
		Utime:      common.NowInt64MS(),
	}
	err = dao.CommissionRateDaoEntity.UpdateCommRate(commissionrates, models.COLUMN_Commissionrates_Min, models.COLUMN_Commissionrates_Max, models.COLUMN_Commissionrates_Commission, models.COLUMN_Commissionrates_Utime)
	if err != nil {
		c.ErrorResponseAndLog(controllers.ERROR_CODE_DB, uid, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回c查询结果
	c.SuccessResponseAndLog(uid, string(c.Ctx.Input.RequestBody), *commissionrates)
}

//删除返佣等级配置
func (c *CommissionRateController) DelCommRates() {
	_, uid, errCode := c.PrePare()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(errCode, uid, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt("id", 0)
	if err != nil || id == 0 {
		c.ErrorResponseAndLog(controllers.ERROR_CODE_PARAMS_ERROR, uid, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)

	err = dao.CommissionRateDaoEntity.DelCommRate(int64(id))
	if err != nil {
		c.ErrorResponseAndLog(controllers.ERROR_CODE_DB, uid, input)
		return
	}

	//返回c查询结果
	data := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(uid, input, data)
}
*/
