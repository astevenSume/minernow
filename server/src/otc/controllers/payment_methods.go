package controllers

import (
	"common"
	json "github.com/mailru/easyjson"
	. "otc_error"
	"strconv"
	"strings"
	common3 "utils/common"
	"utils/otc/dao"
	"utils/otc/models"
)

const (
	Payment = "payment"
)

type PaymentMethodsController struct {
	BaseController
}

// query user payment methods
func (c *PaymentMethodsController) Get() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	var id uint64
	idStr := c.Ctx.Input.Param(KeyIdInput)
	if len(idStr) > 0 {
		var err error
		id, err = strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
			return
		}
	}
	mtype, _ := c.GetUint8("type", 0)

	list, err := dao.PaymentMethodDaoEntity.QueryByUid(uid, id, mtype)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	today, _ := common.TodayTimeRange()
	todayUint32 := uint32(today)
	for _, v := range list {
		if v.UseTime < todayUint32 {
			v.TimesToday = 0
			v.MoneyToday = 0
		}
	}

	if id > 0 {
		if len(list) > 0 {
			c.SuccessResponse(list[0])
		} else {
			c.SuccessResponseWithoutData()
		}
		return
	}

	c.SuccessResponse(map[string]interface{}{
		KeyList: list,
	})
}

//easyjson:json
type PaymentMethodsBindMsg struct {
	Type       uint8  `json:"type"`
	Name       string `json:"name"`
	Account    string `json:"account"`
	Bank       string `json:"bank,omitempty"`
	BankBranch string `json:"bank_branch,omitempty"`
	QrCode     string `json:"qr_code,omitempty"`
	Sms        string `json:"verify_code"`
}

func (c *PaymentMethodsController) Bind() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	if uid < 1 {
		c.ErrorResponse(ERROR_CODE_NO_LOGIN)
		return
	}

	msg := PaymentMethodsBindMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DECODE_FAILED)
		return
	}
	if msg.Type != dao.PayModeWechat && msg.Type != dao.PayModeAli && msg.Type != dao.PayModeBank {
		c.ErrorResponse(ERROR_CODE_PAYMENT_METHOD_UNSUPPORTED)
		return
	}

	//验证短信验证码
	user, err := dao.UserDaoEntity.InfoByUId(uid)
	if err != nil {
		return
	}
	if user.Status != dao.UserStatusActive {
		c.ErrorResponse(ERROR_CODE_USER_NO_ACTIVE)
		return
	}
	check, err := common3.VerifySmsCode(user.NationalCode, user.Mobile, common3.SmsActionPayment, msg.Sms)
	if !check {
		common.LogFuncDebug("check check:%d, sms:%s", check, msg.Sms)
		// 短信验证
		if err == common3.ErrSmsOutTimes {
			c.ErrorResponse(ERROR_CODE_SMS_VERIFY_FAIL_TOO_MATH)
		} else {
			c.ErrorResponse(ERROR_CODE_SMS_ERR)
		}
		return
	}

	//common.LogFuncDebug("qrcode:%v", msg.QrCode)
	//兼容android上传json问题
	msg.QrCode = strings.Replace(msg.QrCode, "_", "/", -1)
	url, urlContent, errCode := UpFileToOss(uid, Payment, msg.QrCode)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	//save
	var pm *models.PaymentMethod
	pm, err = dao.PaymentMethodDaoEntity.Add(uid, msg.Type, msg.Name, msg.Account, msg.Bank, msg.BankBranch, url, urlContent)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(pm)
}

func (c *PaymentMethodsController) Unbind() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	var id uint64

	idStr := c.Ctx.Input.Param(KeyIdInput)
	if len(idStr) > 0 {
		var err error
		id, err = strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
			return
		}
	} else {
		common.LogFuncError("id parameter no found")
		c.ErrorResponse(ERROR_CODE_PARAM_FAILED)
		return
	}

	//
	err := dao.PaymentMethodDaoEntity.Remove(uid, id)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponse(map[string]interface{}{"id": id})
}

func (c *PaymentMethodsController) getId() (id uint64, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	idStr := c.Ctx.Input.Param(KeyIdInput)
	if len(idStr) > 0 {
		var err error
		id, err = strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			errCode = ERROR_CODE_PARAM_FAILED
			return
		}
	} else {
		common.LogFuncError("id parameter no found")
		errCode = ERROR_CODE_PARAM_FAILED
		return
	}

	return
}

//easyjson:json
type PaymentMethodsEditMsg struct {
	Type                uint8  `json:"mtype,omitempty"`
	Name                string `json:"name,omitempty"`
	Account             string `json:"account,omitempty"`
	Bank                string `json:"bank,omitempty"`
	BankBranch          string `json:"bank_branch,omitempty"`
	QrCode              string `json:"qr_code,omitempty"`
	LowMoneyPerTxLimit  int64  `json:"low_money_per_tx_limit,omitempty"`
	HighMoneyPerTxLimit int64  `json:"high_money_per_tx_limit,omitempty"`
	TimesPerDayLimit    int64  `json:"times_per_day_limit,omitempty"`
	MoneyPerDayLimit    int64  `json:"money_per_day_limit,omitempty"`
	MoneySumLimit       int64  `json:"money_sum_limit,omitempty"`
}

func (c *PaymentMethodsController) Edit() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	//get id
	var id uint64
	id, errCode = c.getId()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	//parse body

	msg := PaymentMethodsEditMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DECODE_FAILED)
		return
	}

	msg.QrCode = strings.Replace(msg.QrCode, "_", "/", -1)
	url, urlContent, errCode := UpFileToOss(uid, Payment, msg.QrCode)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	msg.QrCode = url

	// check whether owned by user
	var list []*models.PaymentMethod
	list, err = dao.PaymentMethodDaoEntity.QueryByUid(uid, id, 0)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	if len(list) <= 0 {
		common.LogFuncDebug("payment method %d no owned by %d", id, uid)
		c.ErrorResponse(ERROR_CODE_PAYMENT_METHOD_OWNER_ERROR)
		return
	}

	// get changed parameters
	pm := models.PaymentMethod{
		Pmid: id,
	}

	isChanged := false
	columnNames := []string{}
	if msg.Type != 0 {
		if msg.Type != dao.PayModeWechat && msg.Type != dao.PayModeAli && msg.Type != dao.PayModeBank {
			c.ErrorResponse(ERROR_CODE_PAYMENT_METHOD_UNSUPPORTED)
			return
		}
		pm.MType = msg.Type
		columnNames = append(columnNames, models.COLUMN_PaymentMethod_MType)
		isChanged = true
	}

	type kvpairStr2Str struct {
		key  string
		v    *string
		name string
	}

	for _, v := range []kvpairStr2Str{
		{msg.Name, &pm.Name, models.COLUMN_PaymentMethod_Name},
		{msg.Account, &pm.Account, models.COLUMN_PaymentMethod_Account},
		{msg.Bank, &pm.Bank, models.COLUMN_PaymentMethod_Bank},
		{msg.BankBranch, &pm.BankBranch, models.COLUMN_PaymentMethod_BankBranch},
		{msg.QrCode, &pm.QRCode, models.COLUMN_PaymentMethod_QRCode},
		{urlContent, &pm.QRCodeContent, models.COLUMN_PaymentMethod_QRCodeContent},
	} {
		if v.key != "" {
			*v.v = v.key
			columnNames = append(columnNames, v.name)
			isChanged = true
		}
	}

	type kvpairint642Str struct {
		value int64
		v     *int64
		name  string
	}
	for _, v := range []kvpairint642Str{
		{msg.LowMoneyPerTxLimit, &pm.LowMoneyPerTxLimit, models.COLUMN_PaymentMethod_LowMoneyPerTxLimit},
		{msg.HighMoneyPerTxLimit, &pm.HighMoneyPerTxLimit, models.COLUMN_PaymentMethod_HighMoneyPerTxLimit},
		{msg.TimesPerDayLimit, &pm.TimesPerDayLimit, models.COLUMN_PaymentMethod_TimesPerDayLimit},
		{msg.MoneyPerDayLimit, &pm.MoneyPerDayLimit, models.COLUMN_PaymentMethod_MoneyPerDayLimit},
		{msg.MoneySumLimit, &pm.MoneySumLimit, models.COLUMN_PaymentMethod_MoneySumLimit},
	} {
		if v.value > 0 {
			//todo : every limit has a range, it must be warned while out of range.
			*v.v = v.value
			columnNames = append(columnNames, v.name)
			isChanged = true
		}
	}

	if isChanged { //changed
		err = dao.PaymentMethodDaoEntity.Edit(pm, columnNames...)
		if err != nil {
			c.ErrorResponse(ERROR_CODE_PAYMENT_METHOD_EDIT_FAILED)
			return
		}
	} else { //warn client if nothing changed.
		c.ErrorResponse(ERROR_CODE_PAYMENT_METHOD_EDIT_NOTHING)
		return
	}

	c.SuccessResponseWithoutData()
}

func (c *PaymentMethodsController) Activate() {
	c.changeStatus(dao.PaymentMethodStatusActivated)
}

func (c *PaymentMethodsController) Deactivate() {
	c.changeStatus(dao.PaymentMethodStatusDeactivated)
}

func (c *PaymentMethodsController) changeStatus(status uint8) {
	_, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	//get id
	id, errCode := c.getId()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	//	try change status
	err := dao.PaymentMethodDaoEntity.ChangeStatus(id, status)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_PAYMENT_METHOD_CHANGE_STATUS_FAILED)
		return
	}

	c.SuccessResponse(map[string]interface{}{
		KeyId: id,
	})
}

//easyjson:json
type PaymentMethodsReOrderMsg struct {
	Data []dao.Pmid2Ord `json:"data"`
}

func (c *PaymentMethodsController) ReOrder() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	msg := PaymentMethodsReOrderMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DECODE_FAILED)
		return
	}

	err = dao.PaymentMethodDaoEntity.ReOrder(uid, msg.Data)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}

	c.SuccessResponseWithoutData()
}
