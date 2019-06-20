package usdt

import (
	"common"
	"fmt"
	"utils/usdt/dao"
	"utils/usdt/models"
)

var PriAesKey = [32]byte{'r', 'f', 'B', 'd', '5', '6', 't', 'i',
	'2', 'S', 'M', 't', 'Y', 'v', 'S', 'g',
	'D', '5', 'x', 'A', 'V', '0', 'Y', 'U',
	'9', '=', 'z', 'a', 'm', 'p', 't', 'a'}

type PriKey struct {
	//lock sync.RWMutex
}

func NewPriKey() *PriKey {
	return &PriKey{}
}

var priKeyMgr = NewPriKey()

//// try to bind a private key.
//func (p *PriKey) TryBind(uaid uint64) (pk *models.PriKey, err error) {
//	// must be locked, or some priKey may be bound more than once.
//	p.lock.Lock()
//	defer p.lock.Unlock()
//
//	// get single unused
//	var pkid uint64
//	pkid, err = dao.PriKeyDaoEntity.GetSingleUnused()
//	if err != nil {
//		return
//	}
//
//	//bind
//	err = dao.PriKeyDaoEntity.Bind(uaid, pkid)
//	if err != nil {
//		return
//	}
//
//	//get result
//	pk, err = dao.PriKeyDaoEntity.Query(pkid)
//	if err != nil {
//		return
//	}
//
//	return
//}

// Get get usdt account private key
func (p *PriKey) Get(uaid uint64) (priKey, addr string, err error) {
	var (
		pk      models.PriKey
		account *models.UsdtAccount
	)
	account, err = dao.AccountDaoEntity.Query(uaid)
	if err != nil {
		return
	}
	if account.Pkid <= 0 {
		err = fmt.Errorf("pri key not found")
		return
	}
	pk, err = dao.PriKeyDaoEntity.QueryByPkid(account.Pkid)
	if err != nil {
		return
	}

	priKey, err = common.DecryptFromBase64(pk.Pri, PriAesKey)
	if err != nil {
		return
	}

	addr = pk.Address

	return
}

// GetByAddr get usdt account private key by addr
func (p *PriKey) GetByAddr(Addr string) (priKey, addr string, err error) {
	addr = Addr
	var pk models.PriKey
	pk, err = dao.PriKeyDaoEntity.QueryByAddr(addr)
	if err != nil {
		return
	}

	priKey, err = common.DecryptFromBase64(pk.Pri, PriAesKey)
	if err != nil {
		return
	}

	return
}

// GetByAddr get usdt account private key by addr
func (p *PriKey) GetByPKID(pkid uint64) (priKey, addr string, err error) {
	var pk models.PriKey
	pk, err = dao.PriKeyDaoEntity.QueryByPkid(pkid)
	if err != nil {
		return
	}

	priKey, err = common.DecryptFromBase64(pk.Pri, PriAesKey)
	if err != nil {
		return
	}

	addr = pk.Address

	return
}
