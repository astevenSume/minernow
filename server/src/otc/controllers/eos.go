package controllers

import (
	"eusd/eosplus"
	"otc_error"
	"strings"
	"utils/otc/dao"
)

type EosController struct {
	BaseController
}

func (c *EosController) Info() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}
	d, errCode := eosplus.EosPlusAPI.Wealth.InfoMap(uid)
	if controllers.ERROR_CODE(errCode) != controllers.ERROR_CODE_SUCCESS {

	}
	c.SuccessResponse(d)
}

//easyjson:json
type EosTransferMsg struct {
	NationalCode string `json:"national_code"`
	Mobile       string `json:"mobile"`
	Amount       int64  `json:"amount"`
	Memo         string `json:"memo"`
}

func (c *EosController) Transfer() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	p := &EosTransferMsg{}
	err := c.GetPost(p)
	if err != nil {
		return
	}

	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}
	if p.Amount == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}

	userTo, err := dao.UserDaoEntity.InfoByMobile(p.NationalCode, p.Mobile)
	if err != nil {
		c.ErrorResponse(controllers.ERROR_CODE_DB)
		return
	}
	if userTo.Uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_USER)
		return
	}
	if userTo.Uid == uid {
		c.ErrorResponse(controllers.ERROR_CODE_EUSD_TRANSFER_SELF)
		return
	}

	// 支付密码验证
	// 二次验证 - 短信验证
	//if errCode := c.check2step(uid, false); errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponse(errCode)
	//	return
	//}

	errCode2 := eosplus.EosPlusAPI.Transaction.TransferByUids(uid, userTo.Uid, p.Amount, p.Memo)
	errCode := controllers.ERROR_CODE(errCode2)
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	c.SuccessResponseWithoutData()
}

func (c *EosController) Records() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	types := c.GetString("type")
	page, _ := c.GetInt64("page")
	limit, _ := c.GetInt64("limit")

	typesList := []interface{}{}
	if types != "" {
		tmp := strings.Split(types, ",")
		for _, v := range tmp {
			typesList = append(typesList, v)
		}
	}

	list, meta := eosplus.EosPlusAPI.Wealth.Records(uid, typesList, page, limit)

	c.SuccessResponse(map[string]interface{}{
		"list": list,
		"meta": meta,
	})
}

func (c *EosController) RecordInfo() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_NO_LOGIN)
		return
	}

	id, _ := c.GetParamUint64(":id")
	if id < 1 {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}
	transfer := eosplus.EosPlusAPI.Wealth.RecordInfo(uid, id)
	c.SuccessResponse(transfer)
}
