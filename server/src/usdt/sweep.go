package usdt

import (
	"common"
	"fmt"
	"usdt/explorer"
	"usdt/prices"
	"utils/usdt/dao"
	"utils/usdt/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/mitchellh/mapstructure"
)

// Sweep 归集
func Sweep() (err error) {
	var (
		// 归集门槛
		limitAmount int64
		records     []dao.SweepRecord
		account     models.UsdtAccount
	)

	if limitAmount, err = beego.AppConfig.Int64("usdt::ProcessSweepLimitAmount"); err != nil {
		return err
	}

	limitAmount = limitAmount * 1e8

	for {
		// 查询待归集记录
		if account, err = dao.AccountDaoEntity.QueryPreSweep(limitAmount); err != nil {
			if err == orm.ErrNoRows {
				break
			}
			return
		}
		// 锁定待归集金额
		if err = dao.AccountDaoEntity.PreSweep(account, limitAmount); err != nil {
			return
		}
	}

	total, err := dao.AccountDaoEntity.QueryTotalWaitingSweep(limitAmount)
	if err != nil {
		return
	}

	if total <= 0 {
		return
	}

	var (
		perPage = 100
	)

	pages := (total / perPage)
	if total%perPage > 0 {
		pages++
	}

	// TODO 消息预警,根据时间统计异常数量发送预警
	for i := 1; i <= pages; i++ {
		records, err = dao.AccountDaoEntity.QueryWaitingSweep(limitAmount, i, perPage)
		if err != nil {
			return
		}

		for _, l := range records {
			// 失败目前不重试,待下次触发归集时重新调用
			if err = sweepHandle(l); err != nil {
				common.LogFuncError("%v", err)
				continue
			}
		}
	}

	return
}

func sweepHandle(record dao.SweepRecord) (err error) {
	var (
		account      *models.UsdtAccount
		platformAddr string
		logID        uint64
		sweepAmount  int64
	)
	account, err = dao.AccountDaoEntity.QueryByUid(record.Uid)
	if err != nil {
		return
	}
	if account.Pkid <= 0 {
		err = fmt.Errorf("can't found account by uid %d", record.Uid)
		return
	}
	sweepAmount = record.WaitingCashSweepInteger
	// 根据归集金额获取落入的冷热钱包地址
	platformAddr, err = WalletMgr.GetSweepWallet()
	if err != nil {
		common.LogFuncWarning("get wallet failed : %v", err)
		return
	}

	var (
		pkStr, fromAddr string
	)
	// 根据用户 uaid 获取用户私钥和地址
	pkStr, fromAddr, err = priKeyMgr.Get(account.Uaid)
	if err != nil {
		if err == orm.ErrNoRows {
			common.LogFuncWarning("priKey of uaid %d no found.", account.Uaid)
			err = fmt.Errorf("priKey of uaid %d no found", account.Uaid)
		}
		// update sweep log step
		dao.SweepLogDaoEntity.UpdateStatus(logID, dao.SweepLogStatusFailure, dao.SweepStepGetPriKey, fmt.Sprint(err))
		return
	}

	// 创建一条 sweep log
	logID, err = dao.SweepLogDaoEntity.Add(record.Uid, dao.SweepLogTypeToPlatform, dao.SweepLogStatus, sweepAmount, "", fromAddr, platformAddr)
	if err != nil {
		common.LogFuncWarning("add sweep log failed : %v", err)
		return
	}

	// 待归集金额
	amountTx := fmt.Sprintf("%.16x", sweepAmount)

	//calc fee by recommand fee and tx size.
	var (
		fee, maxFee int
	)

	fee, err = getFee(FeeModeFourHour)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	maxFee, err = getSweepFeeLimit()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var signedTx string
	// 离线签名
	signedTx, err = getSignedTx2(pkStr, amountTx, fromAddr, platformAddr, fee, maxFee)
	if err != nil {
		common.LogFuncError("getSignedTx2 failed : %v", err)
		dao.SweepLogDaoEntity.UpdateStatus(logID, dao.SweepLogStatusFailure, dao.SweepStepSignedTx, fmt.Sprint(err))
		return
	}

	common.LogFuncDebug("signedTx : %s", signedTx)

	var resp map[string]interface{}
	resp, err = explorer.NewExplorer().TransactionPush(signedTx)
	if err != nil {
		common.LogFuncError("transaction push failed : %v", err)
		dao.SweepLogDaoEntity.UpdateStatus(logID, dao.SweepLogStatusFailure, dao.SweepStepTransferTxPushed, fmt.Sprint(err))
		return
	}

	type Msg struct {
		Status string `json:"status"`
		Pushed string `json:"pushed"`
		Tx     string `json:"tx"`
	}

	msg := Msg{}

	err = mapstructure.Decode(resp, &msg)
	if err != nil {
		//信息验证不通过后，发出交易前：恢复提交状态。 (redo again ?)
		dao.SweepLogDaoEntity.UpdateStatus(logID, dao.SweepLogStatusFailure, dao.SweepStepTransferTxPushed, fmt.Sprint(err))
		common.LogFuncError("decode %v to struct failed : %v", resp, err)
		return
	}

	if msg.Status != "OK" {
		//信息验证不通过后，发出交易前：恢复提交状态。 (redo again ?)
		dao.SweepLogDaoEntity.UpdateStatus(logID, dao.SweepLogStatusFailure, dao.SweepStepTransferTxPushed, fmt.Sprintf("push signedTx to %s status : %v", explorer.OMNI_EXPLORER_URL, msg.Status))
		common.LogFuncError("push signedTx to %s status : %v", explorer.OMNI_EXPLORER_URL, msg.Status)
		return
	}

	if msg.Pushed != "success" {
		//信息验证不通过后，发出交易前：恢复提交状态。 (redo again ?)
		dao.SweepLogDaoEntity.UpdateStatus(logID, dao.SweepLogStatusFailure, dao.SweepStepTransferTxPushed, fmt.Sprintf("push signedTx to %s result : %v", explorer.OMNI_EXPLORER_URL, msg.Pushed))
		common.LogFuncError("push signedTx to %s result : %v", explorer.OMNI_EXPLORER_URL, msg.Pushed)
		return
	}

	// 更新 txid
	err = dao.SweepLogDaoEntity.UpdateStatusAndTxid(logID, dao.SweepLogStatusTransferred, msg.Tx)
	if err != nil {
		common.LogFuncError("update txid failed : %v", err)
		return
	}

	// 更新待归集金额到归集金额中
	err = dao.AccountDaoEntity.FinishSweep(record.Uid, sweepAmount)
	if err != nil {
		common.LogFuncError("update txid failed : %v", err)
		return
	}

	return
}

// 获取归集最高手续费
func getSweepFeeLimit() (fee int, err error) {
	var (
		maxFeeCNY, btc2cny float64
	)
	if maxFeeCNY, err = beego.AppConfig.Float("usdt::ProcessSweepMaxFeeCNY"); err != nil {
		return
	}

	if btc2cny = prices.GetPriceFloat64(prices.PRICE_CURRENCY_TYPE_BTC); btc2cny == 0 {
		err = fmt.Errorf("can't get btc to cny rate")
		return
	}

	return int((maxFeeCNY / btc2cny) * 1e8), nil
}
