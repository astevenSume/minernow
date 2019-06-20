package controllers

import (
	"common"
	. "otc_error"
	"strconv"
	"strings"
	"usdt"
	admindao "utils/admin/dao"
	common2 "utils/common"
	"utils/otc/dao"
	usdtdao "utils/usdt/dao"
	"utils/usdt/models"
)

type UsdtController struct {
	BaseController
}

// get user usdt account
func (c *UsdtController) Get() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg, errCode := usdt.QueryByUid(uid)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	c.SuccessResponse(msg)
}

func (c *UsdtController) Balance() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	var errCodeUsdt ERROR_CODE
	balance, errCodeUsdt := usdt.BalanceByUid(uid)
	if errCodeUsdt != ERROR_CODE_SUCCESS {
		c.ErrorResponse(ERROR_CODE(errCodeUsdt))
		return
	}

	c.SuccessResponse(balance)
}

//easyjson:json
type UsdtTransferMsg struct {
	Method       string `json:"method"`
	Address      string `json:"address"`
	NationalCode string `json:"national_code"`
	Mobile       string `json:"mobile"`
	Amount       string `json:"amount"`
	Memo         string `json:"memo"`
	Fee          string `json:"fee"`
}

// start a transfer order
func (c *UsdtController) Transfer() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := UsdtTransferMsg{}
	err := c.GetPost(&msg)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}
	//check params
	switch msg.Method {
	case KeyUsdtTransferMethodAddress:
	case KeyUsdtTransferMethodMobile:
		user, err := dao.UserDaoEntity.InfoByMobile(msg.NationalCode, msg.Mobile)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_NO_USER)
			return
		}

		account, err := usdtdao.AccountDaoEntity.QueryByUid(user.Uid)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_NO_USER)
			return
		}

		msg.Address = account.Address
	default:
		c.ErrorResponse(ERROR_CODE_USDT_TRANSFER_METHOD_UNSUPPORTED)
		return
	}

	if len(msg.Address) <= 0 {
		common.LogFuncWarning("address is nil")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	if len(msg.Memo) > 256 {
		common.LogFuncWarning("memo len too long")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	transferInteger, err := common.CurrencyStrToInt64(msg.Amount)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_USDT_CURRENCY_PARAM_ERROR)
		return
	}

	errCode = c.checkTransferLowerLimit(transferInteger)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	// 支付密码验证
	// 二次验证 - 短信验证
	errCode = c.check2step(uid, false)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	var logId uint64
	// 生成转账申请单
	logId, errCode = usdt.Transfer(uid, msg.Address, msg.Amount, msg.Fee, msg.Memo)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	err = SmsWarning(admindao.ConfigWarningTypeUsdt, uid, msg.Amount)
	if err != nil {
		common.LogFuncError("error:%v", err)
	}

	auditLimitStr, err := common2.AppConfigMgr.String(common2.UsdtTransferAuditLimit)
	if err != nil {
		common.LogFuncError("error:%v", err)
	}

	if err == nil && len(auditLimitStr) > 0 {
		auditLimitInteger, err := common.CurrencyStrToInt64(auditLimitStr)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_USDT_TRANSFER_AUDIT_LIMIT_PARAM_ERROR)
			return
		}

		// 免审核
		if transferInteger <= auditLimitInteger {
			if errCode := usdt.ApproveTransferOutOrder(logId, true); errCode != ERROR_CODE_SUCCESS {
				c.ErrorResponse(errCode)
				return
			}
			// 丢入 mq 队列
			if err = common.RabbitMQPublish(RabbitMQExchangeUsdtTransfer, RabbitMQExchangeUsdtTransfer, []byte(strconv.FormatUint(logId, 10))); err != nil {
				common.LogFuncError("publish usdt.transfer failed:%v", err)
			}

		}
	}

	c.SuccessResponse(map[string]interface{}{
		KeyAmount: msg.Amount,
	})
}

//easyjson:json
type UsdtSyncRechargeTransactionByMobileMsg struct {
	NationalCode string `json:"national_code"`
	Mobile       string `json:"mobile"`
}

//easyjson:json
type UsdtMortgageMsg struct {
	Amount string `json:"amount"`
}

// 抵押usdt
func (c *UsdtController) Mortgage() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(ERROR_CODE_NO_LOGIN)
		return
	}

	msg := UsdtMortgageMsg{}
	err := c.GetPost(&msg)
	if err != nil {
		return
	}

	// 抵押usdt
	errCode = usdt.Mortgage(uid, msg.Amount)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		"amount": msg.Amount,
	})
}

//easyjson:json
type UsdtReleaseMsg struct {
	Amount string `json:"amount"`
}

// 赎回usdt
func (c *UsdtController) Release() {

	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(ERROR_CODE_NO_LOGIN)
		return
	}

	msg := UsdtReleaseMsg{}
	err := c.GetPost(&msg)
	if err != nil {
		common.LogFuncError("%v", err)
		c.ErrorResponse(ERROR_CODE_DECODE_FAILED)
		return
	}

	// release usdt
	errCode = usdt.Release(uid, msg.Amount)
	if errCode != ERROR_CODE_SUCCESS {
		//common.LogFuncError("relese errCode,%d", errCode)
		c.ErrorResponse(errCode)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		"amount": msg.Amount,
	})
}

//获取usdt订单记录
func (c *UsdtController) Records() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(ERROR_CODE_NO_LOGIN)
		return
	}

	types := c.GetString("type")
	page, _ := c.GetInt64("page")
	limit, _ := c.GetInt64("per_page")

	typesList := []interface{}{}
	if types != "" {
		tmp := strings.Split(types, ",")
		for _, v := range tmp {
			typesList = append(typesList, v)
		}
	}

	list, meta, err := usdt.GetRecords(uid, typesList, page, limit)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		"list": list,
		"meta": meta,
	})
}

// 获取usdt订单记录
func (c *UsdtController) Record() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(ERROR_CODE_NO_LOGIN)
		return
	}

	id, err := c.GetUint64(KeyIdInput)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	record, err := usdt.GetRecord(id)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(record)
}

func (c *UsdtController) CancelTransfer() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	id, err := c.GetUint64(KeyIdInput)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	errCode = usdt.CancelTransferOutOrder(uid, id)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	record, _ := usdt.GetRecord(id)

	c.SuccessResponse(record)
}

// CalculateFee 计算手续费
func (c *UsdtController) CalculateFee() {

	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	var (
		err     error
		account *models.UsdtAccount
	)

	account, err = usdtdao.AccountDaoEntity.QueryByUid(uid)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_USDT_ACCOUNT_NO_FOUND)
		return
	}

	if account.Pkid <= 0 {
		c.ErrorResponse(ERROR_CODE_USDT_PRI_KEY_NO_FOUND)
		return
	}

	fee, errCode := usdt.GetRecommandedFee(0)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	c.SuccessResponse(fee)
}

func (c *UsdtController) checkTransferLowerLimit(transferInteger int64) (errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	var (
		err               error
		lowerLimitInteger int64
		lowerLimitStr     string
	)
	lowerLimitStr, err = common2.AppConfigMgr.String(common2.UsdtTransferLowerLimit)
	if err != nil {
		errCode = ERROR_CODE_USDT_TRANSFER_LOWER_LIMIT_PARAM_ERROR
		return
	}

	lowerLimitInteger, err = common.CurrencyStrToInt64(lowerLimitStr)
	if err != nil {
		errCode = ERROR_CODE_USDT_TRANSFER_LOWER_LIMIT_PARAM_ERROR
		return
	}

	if transferInteger < lowerLimitInteger {
		errCode = ERROR_CODE_USDT_TRANSFER_AMOUNT_LESS_THEN_LOWER_LIMIT
		return
	}
	return
}
