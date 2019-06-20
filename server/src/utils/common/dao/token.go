package dao

import (
	"common"
	"fmt"
	. "github.com/astaxie/beego/cache"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
	"utils/common/models"
)

//this entity can be used globally.
var (
	TokenDaoEntity *TokenDao
)

type TokenDao struct {
	common.BaseDao
}

func NewTokenDao(db string) *TokenDao {
	return &TokenDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var adapter Cache
var AccessTokenExpiredSecs int64
var appName string

//type TokenMemItem struct {
//	signature string
//}
//
//func NewTokenMemItem(signature string) *TokenMemItem {
//	return &TokenMemItem{
//		signature: signature,
//	}
//}

func TokenKey(uid uint64, clientType int, mac string) string {
	return fmt.Sprintf("%s.%d.%d.%s", appName, uid, clientType, mac)
}

func TokenKeyInterface(uid, clientType, mac interface{}) string {
	return fmt.Sprintf("%s.%v.%v.%v", appName, uid, clientType, mac)
}

//func (d *TokenMemItem) Signature() string {
//	return d.signature
//}

//// @Title GetFromMem
//// @Description get token info from memory.
//func (d *TokenDao) GetFromMem(uid uint64) (token *TokenMemItem) {
//	v := adapter.Get(fmt.Sprint(uid))
//	if v == nil {
//		return
//	}
//
//	token = v.(*TokenMemItem)
//	return
//}

// @Title Set
// @Description set token to db and memory.
func (d *TokenDao) Set(token *models.Token, isSaveToDb bool) (err error) {
	////save to db
	//if isSaveToDb {
	//	if _, err = d.Orm.InsertOrUpdate(token); err != nil {
	//		common.LogFuncError("save token %+v to db failed : %v", token, err)
	//		return
	//	}
	//}

	//save to memory.
	err = d.save2Mem(token)
	if err != nil {
		return
	}

	return
}

// @Title save2Mem
// @Description set token memory.
// @Params
func (d *TokenDao) save2Mem(token *models.Token) (err error) {
	//save to memory. only the signture part.
	accessTokenSegments := strings.Split(token.AccessToken, ".")
	if len(accessTokenSegments) != 3 || len(accessTokenSegments[2]) <= 0 {
		common.LogFuncError("invalid token %s", token.AccessToken)
		return
	}

	//err = adapter.Put(
	//	TokenKey(token.Uid, int(token.ClientType), token.Mac),
	//	NewTokenMemItem(accessTokenSegments[2]),
	//	time.Second*time.Duration(AccessTokenExpiredSecs))
	//if err != nil {
	//	common.LogFuncError("MEMORYCACHE_PUT CACHE_INDEX_ACCESS_TOKEN failed %v", err)
	//	return
	//}

	common.RedisManger.Set(TokenKey(token.Uid, int(token.ClientType), token.Mac), accessTokenSegments[2], time.Second*time.Duration(AccessTokenExpiredSecs))

	common.LogFuncDebug("token[%s] : %s", TokenKey(token.Uid, int(token.ClientType), token.Mac), accessTokenSegments[2])

	return
}

//// @Title Load
//// @Description load token from db to memory.
//func (d *TokenDao) Load() (err error) {
//	var maps []orm.Params
//	_, err = d.Orm.QueryTable(models.TABLE_Token).Values(&maps)
//	if err == nil {
//		now := time.Now().Unix()
//		for _, m := range maps {
//			d.parseTokenForMem(now, m)
//		}
//	}
//
//	return
//}

// @Description parse token string to jwt.Token.
func (d *TokenDao) ParseToken(tokenStr string) (t *jwt.Token, err error) {
	t, err = jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	if err != nil {
		//common.LogFuncDebug("parse token [%s] failed ： %v", tokenStr, err)
		return
	}

	return
}

//func (d *TokenDao) parseTokenForMem(now int64, m map[string]interface{}) {
//	var (
//		tokenStr string
//		ok       bool
//		v        interface{}
//	)
//
//	v = m[models.ATTRIBUTE_Token_AccessToken]
//	if v == nil {
//		common.LogFuncError("lack of TOKEN.")
//		return
//	}
//
//	if tokenStr, ok = v.(string); ok {
//		if len(tokenStr) > 0 {
//			t, err := d.ParseToken(tokenStr)
//			if err != nil {
//				//common.LogFuncInfo("wrong token str %s", tokenStr)
//				return
//			}
//
//			var claims *MyCustomClaims
//			claims, ok = t.Claims.(*MyCustomClaims)
//			if !ok {
//				return
//			}
//
//			if now > claims.NotBefore && now < claims.ExpiresAt {
//				//common.LogFuncDebug("cache token info %d, %s, %d(%d-%d) ", claims.Uid, t.Signature, claims.ExpiresAt-now, claims.ExpiresAt, now)
//				err := adapter.Put(TokenKeyInterface(m[models.ATTRIBUTE_Token_Uid],
//					m[models.ATTRIBUTE_Token_ClientType],
//					m[models.ATTRIBUTE_Token_Mac]),
//					NewTokenMemItem(t.Signature),
//					time.Second*time.Duration(claims.ExpiresAt-now))
//				if err != nil {
//					common.LogFuncError("add token to memory failed : %v", err)
//					return
//				}
//			}
//		}
//	}
//
//	return
//}

func (d *TokenDao) RemoveToken(uid uint64, clientType int) {
	//remove from db.
	//_, err := d.Orm.Delete(&models.Token{
	//	Uid:        uid,
	//	ClientType: uint32(clientType),
	//	Mac:        "",
	//}, models.COLUMN_Token_Uid, models.COLUMN_Token_ClientType, models.COLUMN_Token_Mac)
	//
	//if err != nil {
	//	common.LogFuncError("delete token data of %d.%d failed : %v", uid, clientType, err)
	//	return
	//}

	//return adapter.Delete(TokenKey(uid, clientType, mac))

	common.RedisManger.Del(TokenKey(uid, clientType, ""))
}

//use for jwt signature.
var SECRET_KEY = []byte("hi this is token secret for XXX server, welcome!")

type MyCustomClaims struct {
	Uid uint64 `json:"uid"`
	Mac string `json:"mac"`
	jwt.StandardClaims
}

//check token
func CheckTokenSignature(clientType int, uid uint64, mac, signature string) bool {
	//get cached token.
	tokenKey := TokenKey(uid, clientType, mac)
	v, err := common.RedisManger.Get(tokenKey).Result()
	if err != nil {
		common.LogFuncError("%v", err)
		return false
	}

	if signature != v {
		common.LogFuncWarning("%s signature (%s) doesn't  equal to cached one (%s), maybe user information is leaked, pls check!",
			tokenKey, signature, v)
		return false
	}

	return true
}

//clear token
func ClearTokenSignature(clientType int, uid uint64, mac string) bool {
	//get cached token.
	tokenKey := TokenKey(uid, clientType, mac)
	err := common.RedisManger.Del(tokenKey)
	if err != nil {
		common.LogFuncError("%v", err)
		return false
	}

	return true
}
