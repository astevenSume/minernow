package controllers

import (
	"common"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"usdt"
	"usdt/explorer"
	"utils/usdt/dao"
	"utils/usdt/models"

	"github.com/astaxie/beego"
	"github.com/btcsuite/btcd/wire"
)

// HotWalletGroupGenerate 生成热钱包私钥
func HotWalletGroupGenerate() {

	hotWalletCount, err := beego.AppConfig.Int("wallet::hot_wallet_count")
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	keys, err := dao.PriKeyDaoEntity.FetchPriKey(dao.HOT_WALLET_MIN_PKID, dao.HOT_WALLET_MAX_PKID)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	// 数量大于等于配置
	if len(keys) >= hotWalletCount {
		return
	}

	var (
		newKey *models.PriKey
	)
	hotWalletCount = hotWalletCount - len(keys)

	for hotWalletCount > 0 {
		pri, addr, _, err := usdt.GenerateSegWitKey(wire.MainNet)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}

		// encrypt to base64
		var priDb string
		priDb, err = common.EncryptToBase64(pri, usdt.PriAesKey)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}

		newKey, err = dao.PriKeyDaoEntity.AddHotWalletPriKey(priDb, addr)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
		dao.AccountDaoEntity.CreateHotWalletAccount(newKey.Pkid, newKey.Address)
		hotWalletCount--
	}
	return
}

// WalletGroupSyncBalance 同步钱包组余额
func WalletGroupSyncBalance() {

	var err error

	if err = syncBalance("hot"); err != nil {
		common.LogFuncError("%v", err)
	}

	if err = syncBalance("cold"); err != nil {
		common.LogFuncError("%v", err)
	}
	return
}

func syncBalance(walletType string) (err error) {

	var (
		addrs            []string
		prikeys          []models.PriKey
		resp             map[string]interface{}
		wallets          []usdt.Wallet
		data             []byte
		minPKID, maxPKID uint64
		redisKey         string
	)

	switch walletType {
	case "hot":
		minPKID, maxPKID = dao.HOT_WALLET_MIN_PKID, dao.HOT_WALLET_MAX_PKID
		redisKey = usdt.REDIS_KEY_HOT_WALLET_GROUP
	case "cold":
		minPKID, maxPKID = dao.COLD_WALLET_MIN_PKID, dao.COLD_WALLET_MAX_PKID
		redisKey = usdt.REDIS_KEY_COLD_WALLET_GROUP
	default:
		return fmt.Errorf("unknow wallet type")
	}

	prikeys, err = dao.PriKeyDaoEntity.FetchPriKey(minPKID, maxPKID)
	if err != nil {
		return
	}

	wallets = make([]usdt.Wallet, 0, len(prikeys))
	loopCount := int(math.Floor(float64(len(prikeys) / 20)))
	for b := 0; b <= loopCount; b++ {
		if b*20 >= len(prikeys) {
			break
		}
		addrs = make([]string, 0, 20)
		begin := b * 20
		limitLen := len(prikeys) - begin
		if limitLen > 20 {
			limitLen = 20
		}
		for _, v := range prikeys[b*20:] {
			addrs = append(addrs, v.Address)
			if len(addrs) == limitLen {

				resp, err = explorer.NewExplorer().Balances(addrs)
				if err != nil {
					return
				}

				if data, err = json.Marshal(resp); err != nil {
					common.LogFuncError("%v", string(data))
					return
				}

				if _, ok := resp["error"]; ok {
					return fmt.Errorf(resp["error"].(string))
				}

				for k, v := range resp {
					vv, ok := v.(map[string]interface{})["balance"].([]interface{})
					if !ok {
						return fmt.Errorf("response data error: %v %v", resp, vv)
					}
					wallet := usdt.Wallet{
						Addr: k,
					}

					for _, vvv := range vv {
						vvvv := vvv.(map[string]interface{})
						switch vvvv["id"].(string) {
						case "0":
							wallet.BTC, err = strconv.ParseInt(vvvv["value"].(string), 10, 64)
							if err != nil {
								return err
							}
						case "31":
							wallet.USDT, err = strconv.ParseInt(vvvv["value"].(string), 10, 64)
							if err != nil {
								return err
							}
						}
					}

					wallets = append(wallets, wallet)
				}
				break
			}
		}

	}

	data, err = json.Marshal(wallets)
	if err != nil {
		return
	}

	if err = common.RedisManger.Set(redisKey, string(data), 0).Err(); err != nil {
		return
	}

	return
}
