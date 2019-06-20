package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"eusd/eosplus"
	eusd_dao "utils/eusd/dao"
	eusd_models "utils/eusd/models"
)

type EusdRetireController struct {
	BaseController
}

func (c *EusdRetireController) Retire() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddEusdRetire, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type Msg struct {
		FromUid  string `json:"from_uid"`
		Quantity int64  `json:"quantity"`
		Memo     string `json:"memo"`
	}

	req := Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncInfo("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddEusdRetire, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	var from_uid uint64
	common.StrToUint(req.FromUid, &from_uid)
	from, err := eusd_dao.WealthDaoEntity.Info(from_uid)
	to := &eusd_models.EosWealth{
		Uid:     0,
		Account: "retire",
	}
	//增加一条销毁记录， cron执行销毁
	eosplus.TransferLogic([]int64{}, from, to, req.Quantity*10000)

	id, err := eusd_dao.EusdRetireDaoEntity.Add(from_uid, req.Quantity, from.Account)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}
	if id <= 0 {
		c.ErrorResponseAndLog(OPActionAddEusdRetire, controllers.ERROR_CODE_UPDATE_FAIL, string(c.Ctx.Input.RequestBody))
		return
	}
	res := eusd_models.EusdRetire{
		Id:       uint64(id),
		From:     from.Account,
		FromUid:  from_uid,
		Quantity: req.Quantity,
	}

	c.SuccessResponseAndLog(OPActionAddEusdRetire, string(c.Ctx.Input.RequestBody), res)
}

func (c *EusdRetireController) TransferToTokenAccount() {

	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddEusdTransfer, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type Msg struct {
		FromUid  string `json:"from_uid"`
		Quantity int64  `json:"quantity"`
		Memo     string `json:"memo"`
	}

	req := Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncInfo("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddEusdTransfer, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	if req.Quantity <= 0 {
		common.LogFuncInfo("req.Quantity <= 0")
		c.ErrorResponseAndLog(OPActionAddEusdTransfer, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var from_uid uint64
	common.StrToUint(req.FromUid, &from_uid)
	from, err := eusd_dao.WealthDaoEntity.Info(from_uid)
	to, err := eusd_dao.WealthDaoEntity.Info(1)
	// 平台账号执行转账
	eosplus.TransferLogic([]int64{}, from, to, req.Quantity*10000)

	c.SuccessResponseAndLog(OPActionAddEusdTransfer, string(c.Ctx.Input.RequestBody), "transfer ok")
}

func (c *EusdRetireController) GetBalanceByUid() {

	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddEusdBalance, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type Msg struct {
		Uid string `json:"uid"`
	}

	req := Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncInfo("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddEusdBalance, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var uid_c uint64
	common.StrToUint(req.Uid, &uid_c)
	wealth, err := eusd_dao.WealthDaoEntity.Info(uid_c)
	// 查询余额
	api := eosplus.EosRpc
	balance := api.GetBalance(wealth.Account)

	res := map[string]interface{}{}
	res["balance"] = balance
	c.SuccessResponseAndLog(OPActionAddEusdBalance, string(c.Ctx.Input.RequestBody), res)
}
