package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"eusd/eosplus"
	"fmt"
	otcerror "otc_error"
	"strconv"
	"utils/eusd/dao"
	otcDao "utils/otc/dao"
)

type AdminEosMethods struct {
	BaseController
}

// get all eos transaction
func (c *AdminEosMethods) GetEosTransactions() {
	c.setOPAction(OPActionReadEosRecord)
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadEosRecord, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt(KEY_ID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadEosRecord, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	page, err := c.GetInt("page", 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadEosRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	perPage, err := c.GetInt("per_page", 10)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadEosRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	uid_c, err := c.GetUint64("uid", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadEosRecord, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	type Meta struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
		Total int `json:"total"`
	}

	meta := &Meta{}

	if id > 0 {
		// 通过id查询
		record, err := dao.TransactionDaoEntity.GetEosTransactionById(uint64(id))
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadEosRecord, controllers.ERROR_CODE_DB, "EOS get transaction by id error")
			return
		}

		c.SuccessResponse(record)
	} else if uid_c > 0 {
		//通过uid查询
		wealth, err := dao.WealthDaoEntity.Info(uid_c)
		listTr, total, err := dao.TransactionDaoEntity.GetEosTransactionsByAccount(wealth.Account, page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadEosRecord, controllers.ERROR_CODE_DB, "EOS get transaction by account error")
			return
		}

		meta.Total = total
		meta.Limit = perPage
		meta.Page = page

		res := map[string]interface{}{}
		res["list"] = listTr
		res["meta"] = meta
		c.SuccessResponse(res)
	} else {
		//查询所有
		rType, err := c.GetUint8(KEY_TYPE, 0)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadEosRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		total, trList, err := dao.TransactionDaoEntity.GetEosTransactions(page, perPage, rType)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadEosRecord, controllers.ERROR_CODE_DB, "EOS get all transaction error")
			return
		}

		meta.Total = total
		meta.Limit = perPage
		meta.Page = page

		res := map[string]interface{}{}
		res["list"] = trList
		res["meta"] = meta
		c.SuccessResponse(res)
	}
}

func (c *AdminEosMethods) GetEosWealth() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadEosWealth, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadEosWealth, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadEosWealth, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	uidCheck, err := c.GetUint64("uid", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadEosWealth, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	type Meta struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
		Total int `json:"total"`
	}

	type WealthRecord struct {
		Uid       string `json:"uid"`
		Account   string `json:"account"`
		Status    int8   `json:"status"`
		Balance   string `json:"balance"`
		Available string `json:"available"`
		Game      string `json:"game"`
		Trade     string `json:"trade"`
		Ctime     int64  `json:"ctime"`
		Utime     int64  `json:"utime"`
		Mobile    string `json:"mobile"`
	}

	meta := &Meta{}

	var input string
	if uidCheck > 0 {
		//根据 uid查询
		input = fmt.Sprintf("{\"uid\":%v}", uidCheck)
		wealth, err := dao.WealthDaoEntity.Info(uidCheck)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadEosWealth, controllers.ERROR_CODE_DB, input)
			return
		}
		wealthRecord := &WealthRecord{}
		wealthRecord.Uid = fmt.Sprintf("%v", wealth.Uid)
		wealthRecord.Account = wealth.Account
		wealthRecord.Status = wealth.Status
		wealthRecord.Balance = eosplus.QuantityInt64ToString(wealth.Balance)
		wealthRecord.Available = eosplus.QuantityInt64ToString(wealth.Available)
		wealthRecord.Game = eosplus.QuantityInt64ToString(wealth.Game)
		wealthRecord.Trade = eosplus.QuantityInt64ToString(wealth.Trade)
		wealthRecord.Ctime = wealth.Ctime
		wealthRecord.Utime = wealth.Utime
		wealthRecord.Mobile, _ = otcDao.UserDaoEntity.GetMobileByUid(wealth.Uid)

		c.SuccessResponseAndLog(OPActionReadEosWealth, input, wealthRecord)
	} else {
		//查询所有
		mobile := c.GetString(KEY_MOBILE)
		if len(mobile) > 0 {
			uidCheck, errCode = GetOtcUidByMobile(mobile)
			if errCode != controllers.ERROR_CODE_SUCCESS {
				c.ErrorResponseAndLog(OPActionReadUsdtWallet, errCode, string(c.Ctx.Input.RequestBody))
				return
			}
		}
		status, err := c.GetInt8(KEY_STATUS, -1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadUsdtWallet, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		input = fmt.Sprintf("{\"mobile\":\"%v\",\"status\":%v,\"page\":%v,\"limit\":%v}", mobile, status, page, perPage)
		wealthList, total, err := dao.WealthDaoEntity.GetEosWealth(page, perPage, uidCheck, status)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadEosWealth, controllers.ERROR_CODE_DB, input)
			return
		}

		var uIds []uint64
		for _, r := range wealthList {
			uIds = append(uIds, r.Uid)
		}
		users, err := otcDao.UserDaoEntity.GetMobileByUIds(uIds)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadEosWealth, controllers.ERROR_CODE_DB, input)
			return
		}

		meta.Total = int(total)
		meta.Limit = perPage
		meta.Page = page
		data := make([]*WealthRecord, 0, len(wealthList))
		for _, r := range wealthList {
			for _, u := range users {
				if r.Uid == u.Uid {
					wealthRecord := &WealthRecord{}
					wealthRecord.Uid = fmt.Sprintf("%v", r.Uid)
					wealthRecord.Status = r.Status
					wealthRecord.Balance = eosplus.QuantityInt64ToString(r.Balance)
					wealthRecord.Available = eosplus.QuantityInt64ToString(r.Available)
					wealthRecord.Game = eosplus.QuantityInt64ToString(r.Game)
					wealthRecord.Trade = eosplus.QuantityInt64ToString(r.Trade)
					wealthRecord.Ctime = r.Ctime
					wealthRecord.Utime = r.Utime
					wealthRecord.Mobile = u.Mobile
					data = append(data, wealthRecord)
					break
				}
			}
		}

		res := map[string]interface{}{}
		res["list"] = data
		res["meta"] = meta
		c.SuccessResponseAndLog(OPActionReadEosWealth, input, res)
	}
}

func (c *AdminEosMethods) GetEosAccount() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadEosAddress, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadEosAddress, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	perPage, err := c.GetInt(KEY_LIMIT, 10)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadEosAddress, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"page\":%v,\"limit\":%v}", page, perPage)

	type Meta struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
		Total int `json:"total"`
	}

	meta := &Meta{}

	accountList, total, err := dao.WealthDaoEntity.GetEosAccount(page, perPage)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadEosAddress, controllers.ERROR_CODE_DB, input)
		return
	}
	meta.Total = total
	meta.Limit = perPage
	meta.Page = page

	res := map[string]interface{}{}
	res["list"] = accountList
	res["meta"] = meta
	c.SuccessResponseAndLog(OPActionReadEosAddress, input, res)
}

//easyjson:json
type EosFrozenMsg struct {
	Uid string `json:"uid"`
}

//冻结EOS
func (c *AdminEosMethods) EosFrozen() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddEosFrozen, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &EosFrozenMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddEosFrozen, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUid, err := strconv.ParseUint(req.Uid, 10, 64)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddEosFrozen, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUser, err := otcDao.UserDaoEntity.InfoByUId(otcUid)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddEosFrozen, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	if otcUser.IsExchanger == otcDao.Exchanger {
		//承兑商专有
		if errCode := eosplus.EosPlusAPI.Otc.Lock(otcUid); errCode != otcerror.ERROR_CODE_SUCCESS {
			common.LogFuncError("errcode:%v", errCode)
			c.ErrorResponseAndLog(OPActionAddEosFrozen, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
			return
		}
	}

	if errCode := eosplus.EosPlusAPI.Wealth.Lock(otcUid); errCode != otcerror.ERROR_CODE_SUCCESS {
		common.LogFuncError("errcode:%v", errCode)
		c.ErrorResponseAndLog(OPActionAddEosFrozen, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionAddEosFrozen, string(c.Ctx.Input.RequestBody))
}

//easyjson:json
type EosUnFrozenMsg struct {
	Uid string `json:"uid"`
}

//解冻EOS
func (c *AdminEosMethods) EosUnFrozen() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddEosUnFrozen, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &EosUnFrozenMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddEosUnFrozen, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUid, err := strconv.ParseUint(req.Uid, 10, 64)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddEosUnFrozen, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUser, err := otcDao.UserDaoEntity.InfoByUId(otcUid)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddEosUnFrozen, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	if otcUser.IsExchanger == otcDao.Exchanger {
		//承兑商专有
		if errCode := eosplus.EosPlusAPI.Otc.Unlock(otcUid); errCode != otcerror.ERROR_CODE_SUCCESS {
			common.LogFuncError("errcode:%v", errCode)
			c.ErrorResponseAndLog(OPActionAddEosUnFrozen, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
			return
		}
	}

	if errCode := eosplus.EosPlusAPI.Wealth.Unlock(otcUid); errCode != otcerror.ERROR_CODE_SUCCESS {
		common.LogFuncError("errcode:%v", errCode)
		c.ErrorResponseAndLog(OPActionAddEosUnFrozen, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionAddEosUnFrozen, string(c.Ctx.Input.RequestBody))
}
