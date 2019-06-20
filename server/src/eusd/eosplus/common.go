package eosplus

import (
	"common"
	"errors"
	"github.com/astaxie/beego"
	. "otc_error"
	common2 "utils/common"
	"utils/eusd/dao"
	"utils/eusd/models"
)

type EosPlusApi struct {
	Account     *Account
	Transaction *Transaction
	Wealth      *Wealth
	Otc         *Otc
	Game        *Game
	Rpc         *RpcKeyBag
}

// eosApi初始化
var EosPlusAPI *EosPlusApi

// cron function container
var FuncContainer interface{}

func EosApiInit() (err error) {
	//初始化
	url := beego.AppConfig.String("eos::main_rpc")
	if url == "" {
		err = errors.New("eos config lack:main_rpc")
		common.LogFuncError("Eos config lack:main_rpc")
		return
	}
	walletUrl := beego.AppConfig.String("eos::wallet_rpc")
	if walletUrl == "" {
		err = errors.New("eos config lack: wallet_rpc")
		common.LogFuncError("Eos config lack: wallet_rpc")
		return
	}
	err = initEosConfig()
	if err != nil {
		return
	}
	EosPlusAPI, err = apiInit(url, walletUrl)

	return
}

var EosRpc RpcKeyBag

func apiInit(mainRpcUrl, walletRpcUrl string) (obj *EosPlusApi, err error) {

	EosRpc.InitApi(mainRpcUrl, walletRpcUrl)
	//EosRpc.SetDebug()
	obj = &EosPlusApi{
		Account:     &Account{},
		Transaction: &Transaction{},
		Wealth:      &Wealth{},
		Otc:         &Otc{},
		Game:        &Game{},
		Rpc:         &EosRpc,
	}

	return
}

func DaoInit() (err error) {
	err = dao.Init(nil)
	return
}

func (api *EosPlusApi) GetErrorMsg(code ERROR_CODE, lang string) string {
	return ErrorMsg(code, lang)
}

// 转账类型
const (
	TransactionType            = iota
	TransactionTypeTransfer    //交易
	TransactionTypeNewAccount  //新建账号
	TransactionTypeBuyRam      //购买存储
	TransactionTypeDelegate    //抵押
	TransactionTypeUnDelegate  //赎回
	TransactionTypeIssueToken  //发币
	TransactionTypeRetireToken //销毁币
	TransactionTypeRexDeposit  //充值到REX
	TransactionTypeRexWithdraw //从REX提现
	TransactionTypeRexRentNet  //从REX租用NET
	TransactionTypeRexRentCpu  //从REX租用CPU
)

// 抵押状态
const (
	EosUseStatus        = iota
	EosUseStatusIng     // 抵押中
	EosUseStatusRecover // 已经取出
)
const EosPrecision = 4 //eos 精度值

type EosConfigS struct {
	TokenAccount        string
	ResourcesAccount    string
	ResourcesAccountKey string
	AdminPrivateKey     string
	Symbol              string
	CpuEos              string
	NetEos              string
	RamEos              string
	PlatformUid         uint64
	CommissionUid       uint64
	SystemUid           uint64
}

var EosConfig = EosConfigS{}

var EosKey = [32]byte{'T', 'r', 'W', '2', '4', 'f', 'x', 'M',
	'P', '6', 'Q', 'L', '7', 's', 'B', '1',
	'G', 'H', 'c', 'r', 'l', 'K', 'R', 'O',
	'A', 'U', 'C', '3', 'D', 'S', 'c', '='}

func initEosConfig() (err error) {
	data, err := common2.AppConfigMgr.FetchString(
		common2.EosConfigKeyTokenAccount,
		common2.EosConfigKeyResourcesAccount,
		common2.EosConfigKeyCpuEos, common2.EosConfigKeyNetEos, common2.EosConfigKeyRamEos,
		common2.EosNoUseAccountLimit,
	)
	if err != nil {
		common.LogFuncError("EosRpc Setting ERR:%v", err)
		return
	}

	for k, v := range data {
		if v == "" {
			err = errors.New(k + " is empty")
			return
		}
	}
	tokenUser, err := TokenAccount()
	if err != nil {
		common.LogFuncError("initEosConfig TokenAccount err: %v", err)
		return
	}
	EosConfig.TokenAccount = tokenUser.Account
	EosConfig.ResourcesAccount = data[common2.EosConfigKeyResourcesAccount]
	if EosConfig.ResourcesAccount == "" {
		EosConfig.ResourcesAccount = EosConfig.TokenAccount
	}

	EosConfig.Symbol = "EUSD"
	EosConfig.CpuEos = data[common2.EosConfigKeyCpuEos]
	EosConfig.NetEos = data[common2.EosConfigKeyNetEos]
	EosConfig.RamEos = data[common2.EosConfigKeyRamEos]

	return
}

//获取持币账号TokenAccount
func TokenAccount() (user *models.EosWealth, err error) {
	if EosConfig.SystemUid == 0 {
		list := []*models.PlatformUser{}
		list, err = dao.PlatformUserDaoEntity.FetchActive(dao.PlatformCateToken)
		if err != nil {
			return
		}
		if len(list) == 0 {
			return
		}

		EosConfig.SystemUid = list[0].Uid
	}
	user, err = dao.WealthDaoEntity.Info(EosConfig.SystemUid)
	return
}

//获取分红账号
func CommissionAccount() (user *models.EosWealth, err error) {
	// todo  按规则分配账号
	if EosConfig.CommissionUid == 0 {
		list := []*models.PlatformUser{}

		list, err = dao.PlatformUserDaoEntity.FetchActive(dao.PlatformCateSalesman)
		if err != nil {
			return
		}
		if len(list) == 0 {
			return
		}

		EosConfig.CommissionUid = list[0].Uid
	}
	user, err = dao.WealthDaoEntity.Info(EosConfig.CommissionUid)

	return
}

//获取游戏平台账号
func GamePlatformAccount() (user *models.EosWealth, err error) {
	// todo  按规则分配账号
	if EosConfig.PlatformUid == 0 {
		list := []*models.PlatformUser{}

		list, err = dao.PlatformUserDaoEntity.FetchActive(dao.PlatformCatePlatform)
		if err != nil {
			return
		}
		if len(list) == 0 {
			return
		}

		EosConfig.PlatformUid = list[0].Uid
	}
	user, err = dao.WealthDaoEntity.Info(EosConfig.PlatformUid)

	return
}
