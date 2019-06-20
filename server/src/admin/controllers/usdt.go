package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	otcerror "otc_error"
	"strconv"
	"time"
	"usdt"
	"utils/admin/dao"
	otcDao "utils/otc/dao"
	usdtDao "utils/usdt/dao"
)

type UsdtController struct {
	BaseController
}

//获取usdt钱包
func (c *UsdtController) GetWallet() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadUsdtWallet, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUId, err := c.GetUint64(KEY_UID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtWallet, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	if otcUId > 0 {
		input = fmt.Sprintf("{\"uid\":%v,\"mobile\":%v}", otcUId)
		data, err := usdtDao.AccountDaoEntity.QueryByUid(otcUId)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadUsdtWallet, controllers.ERROR_CODE_DB, input)
			return
		}
		clientData := ClientUsdtAccount(data)
		clientData.Mobile, _ = otcDao.UserDaoEntity.GetMobileByUid(otcUId)
		c.SuccessResponseAndLog(OPActionReadUsdtWallet, input, clientData)
	} else {
		mobile := c.GetString(KEY_MOBILE)
		if len(mobile) > 0 {
			otcUId, errCode = GetOtcUidByMobile(mobile)
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

		page, err := c.GetInt(KEY_PAGE, 1)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadUsdtWallet, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadUsdtWallet, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}

		//分页查询
		res := map[string]interface{}{}
		input = fmt.Sprintf("{\"mobile\":\"%v\",\"status\":%v,\"page\":%v,\"per_page\":%v}", mobile, status, page, perPage)
		count, data, err := usdtDao.AccountDaoEntity.QueryByPage(page, perPage, otcUId, status)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadUsdtWallet, controllers.ERROR_CODE_DB, input)
			return
		}

		meta := dao.PageInfo{
			Limit: perPage,
			Total: int(count),
			Page:  page,
		}

		res["meta"] = meta
		res["list"] = ClientUsdtAccounts(data)
		c.SuccessResponseAndLog(OPActionReadUsdtWallet, input, res)
	}
}

//获取转入转出记录
func (c *UsdtController) GetTransRecords() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadUsdtTransRecord, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUId, err := c.GetUint64(KEY_UID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtTransRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	mobile := c.GetString(KEY_MOBILE)
	if otcUId == 0 && len(mobile) > 0 {
		otcUId, errCode = GetOtcUidByMobile(mobile)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			c.ErrorResponseAndLog(OPActionReadUsdtTransRecord, errCode, string(c.Ctx.Input.RequestBody))
			return
		}
	}
	status, err := c.GetUint32(KEY_STATUS, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtTurnRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	page, err := c.GetInt64(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtTransRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	limit, err := c.GetInt64(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtTransRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"status\":%v,\"mobile\":\"%v\",\"uid\":%v,\"page\":%v,\"limit\":%v}", status, mobile, otcUId, page, limit)
	types := []interface{}{
		usdtDao.WealthLogTypeTransferIn,
		usdtDao.WealthLogTypeTransferOut,
	}
	strTypes := []string{
		fmt.Sprintf("%v", usdtDao.WealthLogTypeTransferIn),
		fmt.Sprintf("%v", usdtDao.WealthLogTypeTransferOut),
	}
	list, meta, err := usdt.GetDetailRecords(otcUId, status, types, strTypes, page, limit)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtTransRecord, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	res := map[string]interface{}{}
	res["meta"] = meta
	res["list"] = list
	c.SuccessResponseAndLog(OPActionReadUsdtTransRecord, input, res)
}

//获取抵押赎回记录
func (c *UsdtController) GetTurnRecords() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadUsdtTurnRecord, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUId, err := c.GetUint64(KEY_UID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtTurnRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	status, err := c.GetUint32(KEY_STATUS, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtTurnRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	mobile := c.GetString(KEY_MOBILE)
	if otcUId == 0 && len(mobile) > 0 {
		otcUId, errCode = GetOtcUidByMobile(mobile)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			c.ErrorResponseAndLog(OPActionReadUsdtTurnRecord, errCode, string(c.Ctx.Input.RequestBody))
			return
		}
	}

	page, err := c.GetInt64(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtTurnRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	limit, err := c.GetInt64(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtTurnRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"status\":%v,\"mobile\":\"%v\",\"uid\":\"%v\",\"page\":%v,\"limit\":%v}", status, mobile, otcUId, page, limit)
	types := []interface{}{
		usdtDao.WealthLogTypeMortgage,
		usdtDao.WealthLogTypeRelease,
	}
	strTypes := []string{
		fmt.Sprintf("%v", usdtDao.WealthLogTypeMortgage),
		fmt.Sprintf("%v", usdtDao.WealthLogTypeRelease),
	}
	list, meta, err := usdt.GetDetailRecords(otcUId, status, types, strTypes, page, limit)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtTurnRecord, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	res := map[string]interface{}{}
	res["meta"] = meta
	res["list"] = list
	c.SuccessResponseAndLog(OPActionReadUsdtTurnRecord, input, res)
}

//获取提现申请记录
func (c *UsdtController) GetCashRecords() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadUsdtCashRecord, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	otcUId, err := c.GetUint64(KEY_UID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtCashRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	mobile := c.GetString(KEY_MOBILE)
	if otcUId == 0 && len(mobile) > 0 {
		otcUId, errCode = GetOtcUidByMobile(mobile)
		if errCode != controllers.ERROR_CODE_SUCCESS {
			c.ErrorResponseAndLog(OPActionReadUsdtCashRecord, errCode, string(c.Ctx.Input.RequestBody))
			return
		}
	}

	page, err := c.GetInt64(KEY_PAGE, 1)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtCashRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	limit, err := c.GetInt64(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtCashRecord, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	input := fmt.Sprintf("{\"mobile\":\"%v\",\"uid\":\"%v\",\"page\":%v,\"limit\":%v}", mobile, otcUId, page, limit)
	list, meta, err := usdt.GetRecordsByStatus(otcUId, usdtDao.WealthLogStatusOutSubmitted, page, limit)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUsdtCashRecord, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	res := map[string]interface{}{}
	res["meta"] = meta
	res["list"] = list
	c.SuccessResponseAndLog(OPActionReadUsdtCashRecord, input, res)

}

type adminUsdtReq struct {
	Id  string `json:"id"`
	Uid string `json:"uid"`
}

type otcUsdtReq struct {
	Id        string `json:"id"`
	Uid       string `json:"uid"`
	Timestamp uint32 `json:"timestamp"`
}

type otcUsdtAck struct {
	Code otcerror.ERROR_CODE `json:"code"`
}

//通过提现申请
func (c *UsdtController) ApproveTransferOut() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUsdtCashRecordApprove, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &adminUsdtReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddUsdtCashRecordApprove, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//otc请求包
	reqOtc := &otcUsdtReq{
		Id:        req.Id,
		Uid:       req.Uid,
		Timestamp: uint32(time.Now().Unix()),
	}
	//otc响应数据
	ack := otcUsdtAck{}

	//otc post请求
	errCode = PostOtc(RouterApproveTransferOut, map[string]string{
		"id":  req.Id,
		"uid": req.Uid,
	}, reqOtc, &ack, reqOtc.Timestamp)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUsdtCashRecordApprove, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if ack.Code != otcerror.ERROR_CODE_SUCCESS {
		common.LogFuncError("errCode:%v", ack.Code)
		c.ErrorResponseAndLog(OPActionAddUsdtCashRecordApprove, controllers.ERROR_CODE_OTC_FAIL, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddUsdtCashRecordApprove, string(c.Ctx.Input.RequestBody))
}

//拒绝提现申请
func (c *UsdtController) RejectTransferOut() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUsdtCashRecordReject, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &adminUsdtReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddUsdtCashRecordReject, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//otc请求包
	reqOtc := &otcUsdtReq{
		Id:        req.Id,
		Uid:       req.Uid,
		Timestamp: uint32(time.Now().Unix()),
	}
	//otc响应数据
	ack := otcUsdtAck{}

	//otc post请求
	errCode = PostOtc(RouterRejectTransferOut, map[string]string{
		"id":  req.Id,
		"uid": req.Uid,
	}, reqOtc, &ack, reqOtc.Timestamp)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUsdtCashRecordReject, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if ack.Code != otcerror.ERROR_CODE_SUCCESS {
		common.LogFuncError("errCode:%v", ack.Code)
		c.ErrorResponseAndLog(OPActionAddUsdtCashRecordReject, controllers.ERROR_CODE_OTC_FAIL, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddUsdtCashRecordReject, string(c.Ctx.Input.RequestBody))
}

type adminUsdtLockReq struct {
	Uid string `json:"uid"`
}
type otcUsdtStatusReq struct {
	Status    uint8  `json:"status"`
	Uid       uint64 `json:"uid"`
	Timestamp uint32 `json:"timestamp"`
}

//usdt账户锁定
func (c *UsdtController) Lock() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUsdtLock, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &adminUsdtLockReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddUsdtLock, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	otcUid, err := strconv.ParseUint(req.Uid, 10, 64)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddUsdtLock, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//otc请求包
	reqOtc := &otcUsdtStatusReq{
		Status:    usdtDao.STATUS_LOCKED,
		Uid:       otcUid,
		Timestamp: uint32(time.Now().Unix()),
	}
	//otc响应数据
	ack := otcUsdtAck{}

	//otc post请求
	errCode = PostOtc(RouterChangeUsdtStatus, map[string]string{
		"status": fmt.Sprintf("%v", reqOtc.Status),
		"uid":    req.Uid,
	}, reqOtc, &ack, reqOtc.Timestamp)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUsdtLock, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if ack.Code != otcerror.ERROR_CODE_SUCCESS {
		common.LogFuncError("errCode:%v", ack.Code)
		c.ErrorResponseAndLog(OPActionAddUsdtLock, controllers.ERROR_CODE_OTC_FAIL, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddUsdtLock, string(c.Ctx.Input.RequestBody))
}

type adminUsdtUnLockReq struct {
	Uid string `json:"uid"`
}

//usdt账户锁定
func (c *UsdtController) Unlock() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUsdtUnlock, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	req := &adminUsdtUnLockReq{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddUsdtLock, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	otcUid, err := strconv.ParseUint(req.Uid, 10, 64)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddUsdtUnlock, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//otc请求包
	reqOtc := &otcUsdtStatusReq{
		Status:    usdtDao.STATUS_WORKING,
		Uid:       otcUid,
		Timestamp: uint32(time.Now().Unix()),
	}
	//otc响应数据
	ack := otcUsdtAck{}

	//otc post请求
	errCode = PostOtc(RouterChangeUsdtStatus, map[string]string{
		"status": fmt.Sprintf("%v", reqOtc.Status),
		"uid":    req.Uid,
	}, reqOtc, &ack, reqOtc.Timestamp)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUsdtUnlock, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if ack.Code != otcerror.ERROR_CODE_SUCCESS {
		common.LogFuncError("errCode:%v", ack.Code)
		c.ErrorResponseAndLog(OPActionAddUsdtUnlock, controllers.ERROR_CODE_OTC_FAIL, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddUsdtUnlock, string(c.Ctx.Input.RequestBody))
}
