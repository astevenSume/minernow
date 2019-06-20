package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"utils/admin/dao"
	otcDao "utils/otc/dao"
)

type OtcUserController struct {
	BaseController
}

func (c *OtcUserController) GetPageInfo() (uint64, int, int, error) {
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

//获取otc用户
func (c *OtcUserController) GetOtcUser() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadUser, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	_, page, perPage, err := c.GetPageInfo()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUser, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	exchanger, err := c.GetInt(KEY_EXCHANGER, -1)
	if err != nil {
		common.LogFuncError("param err:%v", err)
		return
	}
	status, err := c.GetInt(KEY_STATUS, -1)
	if err != nil {
		common.LogFuncError("param err:%v", err)
		return
	}

	otcUid, err := c.GetUint64(KEY_UID, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUser, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	name := c.GetString("name")
	mobile := c.GetString(KEY_MOBILE)

	//数据获取
	var input string
	if otcUid == 0 {
		//分页查询
		input = fmt.Sprintf("{\"name\":\"%v\",\"mobile\":\"%v\",\"uid\":\"%v\",\"status\":%v,\"page\":%v,"+
			"\"per_page\":%v}", name, mobile, otcUid, status, page, perPage)
		count, data, err := otcDao.UserDaoEntity.QueryPageUser(name, mobile, int8(exchanger), int8(status), page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadUser, controllers.ERROR_CODE_DB, input)
			return
		}

		var list []otcDao.OtcUserInfoClient
		for _, item := range data {
			list = append(list, otcDao.UserDaoEntity.GetOtcUserInfo(&item))
		}
		res := map[string]interface{}{}
		meta := dao.PageInfo{
			Limit: perPage,
			Total: int(count),
			Page:  page,
		}
		res["meta"] = meta
		res["list"] = list
		c.SuccessResponseAndLog(OPActionReadUser, input, res)
	} else {
		input = fmt.Sprintf("{\"uid\":%v}", otcUid)
		user, err := otcDao.UserDaoEntity.UserInfoByUId(otcUid)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadUser, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadUser, input, otcDao.UserDaoEntity.GetOtcUserInfo(user))
	}
}

//更新otc用户
func (c *OtcUserController) UpdateOtcUser() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditUser, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Uid       string `json:"uid"`
		Name      string `json:"name"`
		Mobile    string `json:"mobile"`
		Exchanger int8   `json:"exchanger"`
		Status    int8   `json:"status"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditUser, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var otcUid uint64
	if len(req.Uid) > 0 {
		var err error
		otcUid, err = strconv.ParseUint(req.Uid, 10, 64)
		if err != nil {
			c.ErrorResponseAndLog(OPActionEditUser, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
	}

	_, err = otcDao.UserDaoEntity.UpdateUser(otcUid, req.Exchanger, req.Status, req.Name, req.Mobile)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditUser, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseWithoutDataAndLog(OPActionEditUser, string(c.Ctx.Input.RequestBody))
}

//删除otc用户
func (c *OtcUserController) DelOctUser() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelUser, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	var otcUid uint64
	strOtcUid := c.GetString(KEY_UID)
	if len(strOtcUid) > 0 {
		var err error
		otcUid, err = strconv.ParseUint(strOtcUid, 10, 64)
		if err != nil {
			c.ErrorResponseAndLog(OPActionDelUser, controllers.ERROR_CODE_PARAMS_ERROR, strOtcUid)
			return
		}
	}
	input := fmt.Sprintf("{\"id\":%v}", otcUid)

	err := otcDao.UserDaoEntity.UpdateStatus(otcUid, otcDao.UserStatusDeleted)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelUser, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	data := map[string]interface{}{
		"id": fmt.Sprintf("%v", otcUid),
	}
	c.SuccessResponseAndLog(OPActionDelUser, input, data)
}

//批量删除otc用户
func (c *OtcUserController) DelOctUsers() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelUserBulk, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Uids []string `json:"uids"`
	}
	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionDelUserBulk, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var ids []uint64
	for _, v := range req.Uids {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			c.ErrorResponseAndLog(OPActionDelUserBulk, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
			return
		}
		ids = append(ids, id)
	}
	err = otcDao.UserDaoEntity.DelUsers(ids)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelUserBulk, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	data := map[string]interface{}{
		"uids": req.Uids,
	}
	c.SuccessResponseAndLog(OPActionDelUserBulk, string(c.Ctx.Input.RequestBody), data)
}

//获取收付款方式
func (c *OtcUserController) GetPayment() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadUserPayment, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	var id uint64
	var otcUid uint64
	strOtcUid := c.GetString(KEY_UID)
	if len(strOtcUid) > 0 {
		var err error
		otcUid, err = strconv.ParseUint(strOtcUid, 10, 64)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadUserPayment, controllers.ERROR_CODE_PARAMS_ERROR, strOtcUid)
			return
		}
	}
	strId := c.GetString(KEY_ID)
	if len(strId) > 0 {
		var err error
		id, err = strconv.ParseUint(strId, 10, 64)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadUserPayment, controllers.ERROR_CODE_PARAMS_ERROR, strId)
			return
		}
	}
	input := fmt.Sprintf("{\"uid\":\"%v\",\"id\":%v}", otcUid, id)

	list, err := otcDao.PaymentMethodDaoEntity.QueryByUid(otcUid, id, 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadUserPayment, controllers.ERROR_CODE_DB, input)
		return
	}

	c.SuccessResponseAndLog(OPActionReadUserPayment, input, list)
}

//// 根据用户手机同步usdt链上数据
//func (c *OtcUserController) SyncTransactionByMobile() {
//	_, errCode := c.getUidFromToken()
//	if errCode != controllers.ERROR_CODE_SUCCESS {
//		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
//		return
//	}
//
//	type Msg struct {
//		NationalCode string `json:"national_code"`
//		Mobile       string `json:"mobile"`
//	}
//
//	msg := Msg{}
//	err := c.GetPost(&msg)
//	if err != nil {
//		return
//	}
//	err = usdt.SyncRechargeTransactionByMobile(msg.NationalCode, msg.Mobile)
//	if err != nil {
//		c.ErrorResponse(controllers.ERROR_CODE_PARAM_FAILED)
//		return
//	}
//	c.SuccessResponseWithoutData()
//}

// 根据用户uid同步usdt链上数据
func (c *OtcUserController) SyncTransactionByUid() {
	/*adminUid, errCode := c.getUidFromToken()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		c.ErrorResponseAndLog(OPActionAddUserSyncUsdtbyuid, controllers.ERROR_CODE_NO_LOGIN, adminUid, string(c.Ctx.Input.RequestBody))
		return
	}*/
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUserSyncUsdtbyuid, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type Msg struct {
		Uid string `json:"uid"`
	}

	msg := Msg{}
	err := c.GetPost(&msg)
	if err != nil {
		return
	}

	var uid uint64 = 0
	uid, err = strconv.ParseUint(msg.Uid, 10, 64)
	if err != nil || uid == 0 {
		c.ErrorResponseAndLog(OPActionAddUserSyncUsdtbyuid, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//otc请求包
	type otcReq struct {
		Uid       uint64 `json:"uid"`
		Timestamp uint32 `json:"timestamp"`
	}
	reqOtc := &otcReq{
		Uid:       uid,
		Timestamp: uint32(time.Now().Unix()),
	}

	//otc响应数据
	type otcAck struct {
		Code controllers.ERROR_CODE `json:"code"`
	}
	ack := otcAck{}

	//otc post请求
	errCode = PostOtc(RouterSyncUsdtTransaction, map[string]string{
		"uid": fmt.Sprint(reqOtc.Uid),
	}, reqOtc, &ack, reqOtc.Timestamp)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUserSyncUsdtbyuid, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if ack.Code != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUserSyncUsdtbyuid, controllers.ERROR_CODE_PARAM_FAILED, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddUserSyncUsdtbyuid, string(c.Ctx.Input.RequestBody))
}

// eusd 充值
func (c *OtcUserController) EusdRecharge() {
	c.setOPAction(OPActionAddExchangerRestore)
	c.setRequestData(string(c.Ctx.Input.RequestBody))
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUserEusdRecharge, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type Msg struct {
		Uid      string  `json:"uid"`
		Quantity float64 `json:"quantity"`
	}

	msg := Msg{}
	err := c.GetPost(&msg)
	if err != nil {
		return
	}

	var uid uint64 = 0
	uid, err = strconv.ParseUint(msg.Uid, 10, 64)
	if err != nil || uid == 0 || msg.Quantity <= 0 {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		common.LogFuncDebug("err:%v, uid:%v", err, uid)
		return
	}

	//otc请求包
	type otcReq struct {
		Uid       uint64  `json:"uid"`
		Timestamp uint32  `json:"timestamp"`
		Quantity  float64 `json:"quantity"`
	}
	reqOtc := otcReq{
		Uid:       uid,
		Timestamp: uint32(time.Now().Unix()),
		Quantity:  msg.Quantity,
	}

	//otc响应数据
	type otcAck struct {
		Code controllers.ERROR_CODE `json:"code"`
	}
	ack := otcAck{}

	//otc post请求
	errCode = PostOtc(RouterEustRecharge, map[string]string{
		"uid":      fmt.Sprint(reqOtc.Uid),
		"quantity": fmt.Sprint(reqOtc.Quantity),
	}, &reqOtc, &ack, reqOtc.Timestamp)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUserEusdRecharge, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if ack.Code != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddUserEusdRecharge, ack.Code, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddUserEusdRecharge, string(c.Ctx.Input.RequestBody))
}
