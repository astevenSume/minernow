package controllers

// todo：delete

import (
	"common"
	"eusd/eosplus"
	"fmt"
	"otc/trade"
	. "otc_error"
	"time"
	"usdt"
	"utils/otc/dao"
	usdtdao "utils/usdt/dao"

	"github.com/eoscanada/eos-go"
)

type TestController struct {
	BaseController
}

// 测试接口
func (c *TestController) BecomeExchanger() {
	uid, _ := c.getUidFromToken()
	if uid == 0 {
		c.ErrorResponse(ERROR_CODE_NO_LOGIN)
		return
	}
	eosplus.EosPlusAPI.Wealth.BecomeExchanger(uid)
	dao.UserDaoEntity.BecomeExchange(uid)
	c.SuccessResponseWithoutData()
}

func (c *TestController) Eos() {
	uid, _ := c.getUidFromToken()
	mobile := c.GetString("mobile")
	if mobile != "" {
		res := common.RedisManger.SetNX("Eos_test_"+mobile, 1, 10*time.Second).Val()
		if !res {
			c.ErrorResponse(ERROR_CODE_OP_TOO_FAST)
			return
		}
		user, _ := dao.UserDaoEntity.InfoByMobile("86", mobile)
		_ = eosplus.EosPlusAPI.Wealth.DelegateUsdt(user.Uid, 1000)
		uid = user.Uid
	} else {
		res := common.RedisManger.SetNX(fmt.Sprintf("Eos_test_%v", uid), 1, 10*time.Second).Val()
		if !res {
			c.ErrorResponse(ERROR_CODE_OP_TOO_FAST)
			return
		}
		_ = eosplus.EosPlusAPI.Wealth.DelegateUsdt(uid, 1000)
		//dao2.WealthDaoEntity.TransferToAvailable(uid, 10000000, 0)
	}

	d, errCode := eosplus.EosPlusAPI.Wealth.InfoMap(uid)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}
	c.SuccessResponse(d)
}

func (c *TestController) Test() {
	err := common.RabbitMQPublishDelay(RabbitMQOrderCheck, RabbitMQOrderCheck,
		[]byte(fmt.Sprintf("%d", 193261004061147136)), fmt.Sprintf("%d", uint64(trade.PayExpire().Seconds())*10))
	if err != nil {
		common.LogFuncError("Order2MQ:%v", err)
	}
	//common.RabbitMQPublishDelay(eosplus.RabbitMQEusdTransferCheck, eosplus.RabbitMQEusdTransferCheck, []byte(fmt.Sprintf("%d", 1234568888)), "10000")
	//common.RabbitMQPublish("eusd.test", "eusd.test", []byte(fmt.Sprintf("%d", 999999)))
	c.SuccessResponseWithoutData()
}

func (c *TestController) EusdTransfer2Mq() {
	id := c.GetString("id")

	err := common.RabbitMQPublish(eosplus.RabbitMQEusdTransfer, eosplus.RabbitMQEusdTransfer, []byte(id))
	if err != nil {
		common.LogFuncError("Order2MQ:%v", err)
	}
	c.SuccessResponseWithoutData()
}

func (c *TestController) Info() {
	uid, _ := c.getUidFromToken()

	account := c.GetString("account")

	info := []eos.Asset{}
	if account == "" {
		w, _ := eosplus.EosPlusAPI.Wealth.Info(uid)

		info = eosplus.EosPlusAPI.Rpc.GetBalance(w.Account)
	} else {
		info = eosplus.EosPlusAPI.Rpc.GetBalance(account)
	}

	c.SuccessResponse(info)
}

func (c *TestController) Tx() {
	id := c.GetString("id")
	block, _ := c.GetInt("block")
	eosplus.EosPlusAPI.Rpc.SetDebug()
	info, err := eosplus.EosPlusAPI.Rpc.GetTransfer(id, block)

	common.LogFuncError("%v", err)
	c.SuccessResponse(info)

}

func (c *TestController) UsdtApproveOrder() {
	_, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	id, err := c.GetUint64(KeyIdInput)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_PARAMS_ERROR)
		return
	}

	errCode = usdt.ApproveTransferOutOrder(id, true)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	record, _ := usdt.GetRecord(id)

	c.SuccessResponse(record)
}

func (c *TestController) UsdtSyncBalance() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	errCodeUsdt := usdt.SyncBalance(uid)
	if errCodeUsdt != ERROR_CODE_SUCCESS {
		c.ErrorResponse(ERROR_CODE(errCodeUsdt))
		return
	}

	c.SuccessResponseWithoutData()

	return
}

func (c *TestController) OtcClear() {
	dao.OtcBuyDaoEntity.Orm.Raw("DELETE FROM otc_buy").Exec()
	dao.OtcBuyDaoEntity.Orm.Raw("DELETE FROM otc_sell").Exec()
	c.SuccessResponseWithoutData()
	return
}

// 同步充值记录（测试接口）
func (c *TestController) UsdtSyncRechargeTransactionByMobile() {
	_, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(ERROR_CODE_NO_LOGIN)
		return
	}

	msg := UsdtSyncRechargeTransactionByMobileMsg{}
	err := c.GetPost(&msg)
	if err != nil {
		return
	}
	err = usdt.SyncRechargeTransactionByMobile(msg.NationalCode, msg.Mobile)
	if err != nil {
		c.ErrorResponse(ERROR_CODE_DB)
		return
	}
	c.SuccessResponseWithoutData()
}

// usdt充值接口
func (c *TestController) UsdtDeposit() {
	uid, errCode := c.getUidFromToken()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponse(errCode)
		return
	}

	amount := int64(1000000000000) //10000 usdt
	mobile := c.GetString("mobile")

	// redisKey := mobile
	// if redisKey == "" {
	// 	redisKey = fmt.Sprintf("%v", uid)
	// }
	// res := common.RedisManger.SetNX("test_usdt_deposit"+redisKey, 1, 60*time.Second).Val()
	// if !res {
	// 	c.ErrorResponse(ERROR_CODE_OP_TOO_FAST)
	// 	return
	// }

	if mobile != "" {
		user, err := dao.UserDaoEntity.InfoByMobile("86", mobile)
		if err != nil {
			common.LogFuncError("%v", err)
			c.ErrorResponse(ERROR_CODE_USDT_ACCOUNT_NO_FOUND)
			return
		}

		err = usdtdao.AccountDaoEntity.Deposit(user.Uid, amount)
		if err != nil {
			common.LogFuncError("%v", err)
			c.ErrorResponse(ERROR_CODE_USDT_DEPOSIT_FAILED)
			return
		}

	} else {
		err := usdtdao.AccountDaoEntity.Deposit(uid, amount)
		if err != nil {
			common.LogFuncError("%v", err)
			c.ErrorResponse(ERROR_CODE_USDT_DEPOSIT_FAILED)
			return
		}
	}

	c.SuccessResponseWithoutData()
}
