package usdt

import (
	"common"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
	"utils/usdt/dao"
	"utils/usdt/models"

	"github.com/astaxie/beego"
	lock "github.com/bsm/redis-lock"
	"github.com/btcsuite/btcd/wire"
)

const (
	REDIS_KEY_COLD_WALLET_GROUP = "usdt_cold_wallet_group"
	REDIS_KEY_HOT_WALLET_GROUP  = "usdt_hot_wallet_group"
)

type (
	walletMgr struct {
		lock sync.Mutex
	}
	Wallet struct {
		Addr string `json:"address"`
		BTC  int64  `json:"btc"`
		USDT int64  `json:"usdt"`
	}
	Wallets      []Wallet
	UsdtAccounts []models.UsdtAccount
)

func (w Wallets) Len() int           { return len(w) }
func (w Wallets) Less(i, j int) bool { return w[i].USDT < w[j].USDT }
func (w Wallets) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }

func (ua UsdtAccounts) Len() int           { return len(ua) }
func (ua UsdtAccounts) Less(i, j int) bool { return ua[i].AvailableInteger < ua[j].AvailableInteger }
func (ua UsdtAccounts) Swap(i, j int)      { ua[i], ua[j] = ua[j], ua[i] }

// WalletMgr 钱包组
var WalletMgr = &walletMgr{}

// GetTransferWallet 平台转账热钱包地址
func (w *walletMgr) GetTransferWallet() (priKey, addr string, l *lock.Locker, err error) {

	// var (
	// 	data       []byte
	// 	hotWallets Wallets
	// )

	// w.lock.Lock()
	// defer func() {
	// 	w.lock.Unlock()
	// }()

	// if data, err = common.RedisManger.Get(REDIS_KEY_HOT_WALLET_GROUP).Bytes(); err != nil {
	// 	return
	// }

	// if err = json.Unmarshal(data, &hotWallets); err != nil {
	// 	return
	// }

	// if len(hotWallets) <= 0 {
	// 	err = fmt.Errorf("hot wallet not found")
	// 	return
	// }

	// sort.Sort(hotWallets)

	// for i := len(hotWallets) - 1; i >= 0; i-- {

	// 	l, err = common.RedisLock2("usdt_transfer_using_wallet_"+hotWallets[i].Addr, lock.Options{
	// 		LockTimeout: time.Second * time.Duration(60),
	// 		RetryCount:  common.RetryCount,
	// 		RetryDelay:  common.RetryDelay,
	// 	})
	// 	if err != nil {
	// 		continue
	// 	}

	// 	priKey, addr, err = priKeyMgr.GetByAddr(hotWallets[i].Addr)
	// 	return
	// }

	// err = fmt.Errorf("hot wallet all using")
	// return
	var (
		hotWallets UsdtAccounts
		utxos      []Utxo
	)

	w.lock.Lock()
	defer func() {
		w.lock.Unlock()
	}()

	hotWallets, err = dao.AccountDaoEntity.QueryWalletAccount(dao.HOT_WALLET_MIN_PKID, dao.HOT_WALLET_MAX_PKID)

	if len(hotWallets) <= 0 {
		err = fmt.Errorf("hot wallet not found")
		return
	}

	sort.Sort(hotWallets)

	for i := len(hotWallets) - 1; i >= 0; i-- {

		l, err = common.RedisLock2("usdt_transfer_using_wallet_"+hotWallets[i].Address, lock.Options{
			LockTimeout: time.Second * time.Duration(60),
			RetryCount:  common.RetryCount,
			RetryDelay:  common.RetryDelay,
		})
		if err != nil {
			continue
		}

		priKey, addr, err = priKeyMgr.GetByPKID(hotWallets[i].Pkid)

		utxos, err = getUnspent(wire.MainNet, addr)
		if err != nil {
			continue
		}

		if len(utxos) == 0 {
			continue
		}

		return
	}

	err = fmt.Errorf("hot wallet all using")
	return
}

// GetSweepWallet 根据金额从钱包组中获取归集的冷钱包或者热钱包的地址
// params:
// 		amount 本次转账的金额
func (w *walletMgr) GetSweepWallet() (address string, err error) {

	var (
		data                            []byte
		hotWallets, coldWallets         Wallets
		hotTotalAmount, coldTotalAmount int64
		hotWalletProportion             float64
	)

	hotWalletProportion, err = beego.AppConfig.Float("wallet::hot_proportion")
	if err != nil {
		return
	}

	if hotWalletProportion <= 0 || hotWalletProportion >= 1 {
		err = fmt.Errorf("hot wallet proportion error : %v", hotWalletProportion)
	}

	if data, err = common.RedisManger.Get(REDIS_KEY_HOT_WALLET_GROUP).Bytes(); err != nil {
		return
	}

	if err = json.Unmarshal(data, &hotWallets); err != nil {
		return
	}

	if len(hotWallets) <= 0 {
		err = fmt.Errorf("hot wallet not found")
		return
	}

	if data, err = common.RedisManger.Get(REDIS_KEY_COLD_WALLET_GROUP).Bytes(); err != nil {
		return
	}

	if err = json.Unmarshal(data, &coldWallets); err != nil {
		return
	}

	if len(coldWallets) <= 0 {
		err = fmt.Errorf("cold wallet not found")
		return
	}

	for _, v := range hotWallets {
		hotTotalAmount += v.USDT
	}

	for _, v := range coldWallets {
		coldTotalAmount += v.USDT
	}

	// 热钱包占比大于配置,使用冷钱包
	if float64(hotTotalAmount)/float64(hotTotalAmount+coldTotalAmount) > hotWalletProportion {
		sort.Sort(coldWallets)
		address = coldWallets[0].Addr
		return
	}

	// 热钱包地址
	sort.Sort(hotWallets)
	address = hotWallets[0].Addr
	return
}
