package cron

import (
	"common"
	"eusd/eosplus"
	"github.com/eoscanada/eos-go/ecc"
	"math/rand"
	"otc_error"
	"time"
	dao2 "utils/eos/dao"
	models2 "utils/eos/models"
	"utils/eusd/dao"
	"utils/eusd/models"
)

func AccountNumCheck(lowLimit int64) {
	num, err := dao.AccountDaoEntity.CountNotUse()
	if err != nil {
		return
	}

	// 账号小于限定值创建新账号
	if num >= lowLimit {
		return
	}

	for ; num < lowLimit; num++ {
		for i := 0; i < 10; i++ {
			ok := CreateRandName()
			if ok {
				break
			}
		}
	}
}

func randName(n int) string {
	rand.Seed(time.Now().UnixNano())
	letterBytes := "abcdefghijklmnopqrstuvwxyz12345"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func CreateRandName() (ok bool) {
	name := randName(12)

	api := eosplus.EosPlusAPI.Rpc

	ownerPrivateKey, err := ecc.NewRandomPrivateKey()
	ownerPublicKey := ownerPrivateKey.PublicKey()

	activePrivateKey, err := ecc.NewRandomPrivateKey()
	activePublicKey := activePrivateKey.PublicKey()

	_, errCode, _ := api.AccountCreate(eosplus.EosConfig.ResourcesAccount, name, ownerPublicKey.String(), activePublicKey.String())
	if errCode != controllers.ERROR_CODE_SUCCESS {
		return
	}
	id, err := dao2.EosAccountKeysEntity.Create(&models2.EosAccountKeys{
		Account:          name,
		PublicKeyOwner:   ownerPublicKey.String(),
		PrivateKeyOwner:  ownerPrivateKey.String(),
		PublicKeyActive:  activePublicKey.String(),
		PrivateKeyActive: activePrivateKey.String(),
		Ctime:            common.NowInt64MS(),
	})

	acc := &models.EosAccount{
		Id:      uint64(id),
		Uid:     0,
		Account: name,
		Status:  dao.AccountStatusNoUse,
		Ctime:   common.NowInt64MS(),
	}
	_, err = dao.AccountDaoEntity.Create(acc)
	if err != nil {
		return
	}
	ok = true
	return
}
