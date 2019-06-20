package controllers

import (
	"common"
	"eusd/eosplus"
	"fmt"
	"otc/trade"
	. "otc_error"
	"strconv"
	"strings"
	"sync"
	"usdt"
	"utils/admin/dao"
	adminmodels "utils/admin/models"
	gamedao "utils/game/dao"
	otcdao "utils/otc/dao"
	usdtdao "utils/usdt/dao"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	json "github.com/mailru/easyjson"
)

type AdminIpWhiteListManager struct {
	lock  sync.RWMutex
	items []adminmodels.IpWhiteList
}

func NewAdminIpWhiteListManager() *AdminIpWhiteListManager {
	return &AdminIpWhiteListManager{}
}

// load ip white list
func (m *AdminIpWhiteListManager) load() (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	//
	m.items, err = dao.IpWhiteListDaoEntity.All()
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		return
	}

	return
}

// check if ip is in white list
func (m *AdminIpWhiteListManager) check(ip string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	//
	if ip == "127.0.0.1" {
		return true
	}

	for _, v := range m.items {
		if v.Ip == ip {
			return true
		}
	}

	return false
}

var adminIpWhiteListMgr = NewAdminIpWhiteListManager()

type ApiController struct {
	BaseController
}

func (c *ApiController) check(params map[string]string, timestamp uint32, sign string) (ok bool) {
	// check md5 sign
	ok = common.CompareMd5(sign, common.GenerateSource(params, timestamp), otcdao.SIGNATURE_SALT)
	if !ok {
		return
	}

	// check ip
	ip := common.ClientIP(c.Ctx)
	ok = adminIpWhiteListMgr.check(ip)
	if !ok {
		common.LogFuncWarning("ip %s is not in white list", ip)
		return
	}

	return
}

//easyjson:json
type ApiDistributeCommissionMsg struct {
	Time      int64  `json:"time"`
	Timestamp uint32 `json:"timestamp"`
}

func (c *ApiController) DistributeCommission() {
	return
}

//easyjson:json
type ApiCalcCommissionMsg struct {
	Time      int64  `json:"time"`
	Timestamp uint32 `json:"timestamp"`
}

func (c *ApiController) CalcCommission() {
	return
}

//easyjson:json
type ApiCancelOrderMsg struct {
	Id        string `json:"id"`
	Uid       string `json:"uid"`
	AdminId   uint32 `json:"admin_id"`
	Timestamp uint32 `json:"timestamp"`
}

//取消订单
func (c *ApiController) CancelOrder() {
	sign := c.GetString(gamedao.KeySign)
	if len(sign) <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	msg := ApiCancelOrderMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_DECODE_FAILED)
		return
	}

	if !c.check(map[string]string{
		"id":       msg.Id,
		"uid":      msg.Uid,
		"admin_id": fmt.Sprintf("%d", msg.AdminId),
	}, msg.Timestamp, sign) {
		common.LogFuncError("check fail")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	id, err := strconv.ParseUint(msg.Id, 10, 64)
	if err != nil {
		common.LogFuncError("err:%v", err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	appeal, err := otcdao.AppealDaoEntity.QueryById(id)
	if err != nil || appeal.OrderId == 0 {
		common.LogFuncError("err:%v", err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}
	if appeal.Status == otcdao.AppealStatusResolved {
		c.ErrorResponse(ERROR_CODE_OTC_ORDER_APPEAL_RESOLVED)
		return
	}

	uid, err := strconv.ParseUint(msg.Uid, 10, 64)
	if err != nil {
		common.LogFuncError("err:%v", err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	ipIdx := strings.Index(c.Ctx.Request.RemoteAddr, ":")
	ip := c.Ctx.Request.RemoteAddr[:ipIdx]
	errCode := trade.AppealLogicEntity.CancelOrder(appeal.OrderId, uid, ip)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	_, err = otcdao.AppealDealLogDaoEntity.Create(otcdao.AppealDealLogActionCancel, msg.AdminId, id, appeal.OrderId)
	if err != nil {
		common.LogFuncError("err:%v", err)
	}

	c.SuccessResponseWithoutData()
}

//easyjson:json
type ApiConfirmOrderPayMsg struct {
	Id        string `json:"id"`
	Uid       string `json:"uid"`
	AdminId   uint32 `json:"admin_id"`
	Timestamp uint32 `json:"timestamp"`
}

//客服申述确认已付款(买家点取消或系统超时自动取消)
func (c *ApiController) ConfirmOrderPay() {
	sign := c.GetString(gamedao.KeySign)
	if len(sign) <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	msg := ApiConfirmOrderPayMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_DECODE_FAILED)
		return
	}

	if !c.check(map[string]string{
		"id":       msg.Id,
		"uid":      msg.Uid,
		"admin_id": fmt.Sprintf("%d", msg.AdminId),
	}, msg.Timestamp, sign) {
		common.LogFuncError("check fail")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	id, err := strconv.ParseUint(msg.Id, 10, 64)
	if err != nil {
		common.LogFuncError("err:%v", err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	appeal, err := otcdao.AppealDaoEntity.QueryById(id)
	if err != nil || appeal.OrderId == 0 {
		common.LogFuncError("err:%v", err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}
	if appeal.Status == otcdao.AppealStatusResolved {
		c.ErrorResponse(ERROR_CODE_OTC_ORDER_APPEAL_RESOLVED)
		return
	}

	uid, err := strconv.ParseUint(msg.Uid, 10, 64)
	if err != nil {
		common.LogFuncError("err:%v", err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	ipIdx := strings.Index(c.Ctx.Request.RemoteAddr, ":")
	ip := c.Ctx.Request.RemoteAddr[:ipIdx]
	errCode := trade.AppealLogicEntity.ConfirmOrderPay(appeal.OrderId, uid, ip)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	_, err = otcdao.AppealDealLogDaoEntity.Create(otcdao.AppealDealLogActionConfirmPay, msg.AdminId, id, appeal.OrderId)
	if err != nil {
		common.LogFuncError("err:%v", err)
	}

	c.SuccessResponseWithoutData()
}

//easyjson:json
type ApiConfirmOrderMsg struct {
	Id        string `json:"id"`
	Uid       string `json:"uid"`
	AdminId   uint32 `json:"admin_id"`
	Timestamp uint32 `json:"timestamp"`
}

//确认订单
func (c *ApiController) ConfirmOrder() {
	sign := c.GetString(gamedao.KeySign)
	if len(sign) <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	msg := ApiConfirmOrderMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_DECODE_FAILED)
		return
	}

	if !c.check(map[string]string{
		"id":       msg.Id,
		"uid":      msg.Uid,
		"admin_id": fmt.Sprintf("%d", msg.AdminId),
	}, msg.Timestamp, sign) {
		common.LogFuncError("check fail")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	id, err := strconv.ParseUint(msg.Id, 10, 64)
	if err != nil {
		common.LogFuncError("err:%v", err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	appeal, err := otcdao.AppealDaoEntity.QueryById(id)
	if err != nil || appeal.OrderId == 0 {
		common.LogFuncError("err:%v", err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}
	if appeal.Status == otcdao.AppealStatusResolved {
		c.ErrorResponse(ERROR_CODE_OTC_ORDER_APPEAL_RESOLVED)
		return
	}

	uid, err := strconv.ParseUint(msg.Uid, 10, 64)
	if err != nil {
		common.LogFuncError("err:%v", err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	ipIdx := strings.Index(c.Ctx.Request.RemoteAddr, ":")
	ip := c.Ctx.Request.RemoteAddr[:ipIdx]
	errCode := trade.AppealLogicEntity.ConfirmOrder(appeal.OrderId, uid, ip)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	_, err = otcdao.AppealDealLogDaoEntity.Create(otcdao.AppealDealLogActionConfirm, msg.AdminId, id, appeal.OrderId)
	if err != nil {
		common.LogFuncError("err:%v", err)
	}

	c.SuccessResponseWithoutData()
}

//easyjson:json
type ApiSyncUsdtTransactionMsg struct {
	Uid       uint64 `json:"uid"`
	Timestamp uint32 `json:"timestamp"`
}

// 同步usdt链上交易记录
func (c *ApiController) SyncUsdtTransaction() {
	sign := c.GetString(gamedao.KeySign)
	if len(sign) <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	msg := ApiSyncUsdtTransactionMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	// 签名认证
	if !c.check(map[string]string{
		"uid": fmt.Sprintf("%d", msg.Uid),
	}, msg.Timestamp, sign) {
		common.LogFuncError("check fail")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	//同步交易数据
	err = usdt.SyncRechargeTransactionByUid(msg.Uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}
	c.SuccessResponseWithoutData()
}

//easyjson:json
type ApiApproveTransferOutMsg struct {
	Id        string `json:"id"`
	Uid       string `json:"uid"`
	Timestamp uint32 `json:"timestamp"`
}

//批准usdt提现
func (c *ApiController) ApproveTransferOut() {
	sign := c.GetString(KeySign)
	if len(sign) <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	msg := ApiApproveTransferOutMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_DECODE_FAILED)
		return
	}

	if !c.check(map[string]string{
		"id":  msg.Id,
		"uid": msg.Uid,
	}, msg.Timestamp, sign) {
		common.LogFuncError("check fail")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	id, err := strconv.ParseUint(msg.Id, 10, 64)
	if err != nil {
		common.LogFuncError("err:%v", err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	if errCode := usdt.ApproveTransferOutOrder(id, true); errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	if err = common.RabbitMQPublish(RabbitMQExchangeUsdtTransfer, RabbitMQExchangeUsdtTransfer, []byte(msg.Id)); err != nil {
		common.LogFuncError("publish usdt_transfer failed:%v", err)
	}

	c.SuccessResponseWithoutData()
}

//easyjson:json
type ApiRejectTransferOutMsg struct {
	Id        string `json:"id"`
	Uid       string `json:"uid"`
	Timestamp uint32 `json:"timestamp"`
}

//拒绝usdt提现
func (c *ApiController) RejectTransferOut() {
	sign := c.GetString(gamedao.KeySign)
	if len(sign) <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	msg := ApiRejectTransferOutMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_DECODE_FAILED)
		return
	}

	if !c.check(map[string]string{
		"id":  msg.Id,
		"uid": msg.Uid,
	}, msg.Timestamp, sign) {
		common.LogFuncError("check fail")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	id, err := strconv.ParseUint(msg.Id, 10, 64)
	if err != nil {
		common.LogFuncError("err:%v", err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	if errCode := usdt.ApproveTransferOutOrder(id, false); errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	c.SuccessResponseWithoutData()
}

//easyjson:json
type ApiEusdRechargeMsg struct {
	Uid       uint64  `json:"uid"`
	Timestamp uint32  `json:"timestamp"`
	Quantity  float64 `json:"quantity"`
}

// eusd 充值
func (c *ApiController) EusdRecharge() {
	sign := c.GetString(gamedao.KeySign)
	if len(sign) <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	msg := ApiEusdRechargeMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	// 签名认证
	if !c.check(map[string]string{
		"uid":      fmt.Sprint(msg.Uid),
		"quantity": fmt.Sprint(msg.Quantity),
	}, msg.Timestamp, sign) {
		common.LogFuncError("check fail")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	errCode := eosplus.EosPlusAPI.Wealth.DelegateUsdt(msg.Uid, msg.Quantity)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	c.SuccessResponseWithoutData()
}

//easyjson:json
type apiChangeUsdtStatusMsg struct {
	UID       uint64 `json:"uid"`
	Status    uint8  `json:"status"`
	Timestamp uint32 `json:"timestamp"`
}

// ChangeUsdtStatus usdt账户锁定/解锁
func (c *ApiController) ChangeUsdtStatus() {

	sign := c.GetString(gamedao.KeySign)
	if len(sign) <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	msg := apiChangeUsdtStatusMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	// 签名认证
	if !c.check(map[string]string{
		"uid":    fmt.Sprint(msg.UID),
		"status": fmt.Sprint(msg.Status),
	}, msg.Timestamp, sign) {
		common.LogFuncError("check fail")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	checkStatus := false
	for _, v := range []uint8{usdtdao.STATUS_LOCKED, usdtdao.STATUS_WORKING} {
		if msg.Status == v {
			checkStatus = true
			break
		}
	}

	if !checkStatus {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	if err = usdtdao.AccountDaoEntity.UpdateStatus(msg.UID, msg.Status); err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponseWithoutData()
}

//easyjson:json
type apiAdminTask struct {
	Name   string `json:"name"`
	Spec   string `json:"spec"`
	Status string `json:"status"`
	Prev   string `json:"prev"`
	Next   string `json:"next"`
}

//easyjson:json
type apiAdminTaskList struct {
	Items []apiAdminTask `json:"items"`
}

//easyjson:json
type apiAdminTaskMsg struct {
	Name      string `json:"name"`
	Timestamp uint32 `json:"timestamp"`
}

func (c *ApiController) AdminTaskList() {
	sign := c.GetString(KeySign)
	if len(sign) <= 0 {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	msg := apiAdminTaskMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	// 签名认证
	if !c.check(map[string]string{
		"name": fmt.Sprint(msg.Name),
	}, msg.Timestamp, sign) {
		common.LogFuncError("check fail")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	if len(msg.Name) <= 0 {
		list := apiAdminTaskList{}
		for name, t := range toolbox.AdminTaskList {
			list.Items = append(list.Items, apiAdminTask{
				Name:   name,
				Spec:   t.GetSpec(),
				Status: t.GetStatus(),
				Prev:   t.GetPrev().Format("2006-01-02 15:04:05"),
				Next:   t.GetNext().Format("2006-01-02 15:04:05"),
			})
		}
		c.SuccessResponse(map[string]interface{}{
			"list": list,
		})
	} else {
		for name, t := range toolbox.AdminTaskList {
			if msg.Name == name {
				c.SuccessResponse(
					apiAdminTask{
						Name:   name,
						Spec:   t.GetSpec(),
						Status: t.GetStatus(),
						Prev:   t.GetPrev().Format("2006-01-02 15:04:05"),
						Next:   t.GetNext().Format("2006-01-02 15:04:05"),
					},
				)
				return
			}
		}

		c.ErrorResponse(ERROR_CODE_SUCCESS)
	}
}
