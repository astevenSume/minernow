package common

import (
	"common"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"utils/common/dao"
	"utils/common/models"
)

func Init(entityInitFunc common.EntityInitFunc) (err error) {
	err = dao.Init(entityInitFunc)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// @Title GenerateToken
// @Description
// @Params
func GenerateToken(uid uint64, mac string) (token string, err error) {
	//generate jwt token
	claims := dao.MyCustomClaims{
		Uid: uid,
		Mac: mac,
		StandardClaims: jwt.StandardClaims{
			Id:        fmt.Sprint(uid),
			NotBefore: int64(time.Now().Unix()),
			ExpiresAt: int64(time.Now().Unix() + dao.AccessTokenExpiredSecs),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString(dao.SECRET_KEY)
	if err != nil {
		common.LogFuncError("generate access_token failed : %v", err)
		return
	}

	return

}

// @Description check token
func CheckToken(clientType int, token string) (ok bool, uid uint64, mac string) {
	t, err := ParseToken(token)
	if err != nil {
		return
	}

	var claims *dao.MyCustomClaims
	if claims, ok = t.Claims.(*dao.MyCustomClaims); ok && t.Valid {
		//check if token equals to cached token
		ok, uid, mac = dao.CheckTokenSignature(clientType, claims.Uid, claims.Mac, t.Signature), claims.Uid, claims.Mac
		if uid == 0 {
			common.LogFuncWarning("token.uid == 0, there must be something wrong with token.")
			ok = false
			return
		}

		return
	} else {
		ok = false
		common.LogFuncDebug("parse token failed : %v", err)
		return
	}

	return
}

// @Description clear token
func ClearToken(clientType int, token string) (ok bool) {
	t, err := ParseToken(token)
	if err != nil {
		return
	}

	var claims *dao.MyCustomClaims
	if claims, ok = t.Claims.(*dao.MyCustomClaims); ok && t.Valid {
		//check if token equals to cached token
		ok = dao.ClearTokenSignature(clientType, claims.Uid, claims.Mac)
		return
	} else {
		ok = false
		common.LogFuncDebug("parse token failed : %v", err)
		return
	}
}

//// @Description get uid from jwt token. no matter token is valid or not.
//func GetUidFromToken(token string) (uid uint64, err error) {
//	var t *jwt.Token
//	t, err = jwt.ParseWithClaims(token, &dao.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return dao.SECRET_KEY, nil
//	})
//
//	if t != nil {
//		if claims, ok := t.Claims.(*dao.MyCustomClaims); ok {
//			uid = claims.Uid
//			return
//		}
//	}
//
//	return
//}

// @Description parse token string to jwt.Token.
func ParseToken(tokenStr string) (t *jwt.Token, err error) {
	return dao.TokenDaoEntity.ParseToken(tokenStr)
}

func ResetToken(uid uint64, clientType int) (token string, err error) {
	token, err = GenerateToken(uid, "")
	if err != nil {
		return
	}

	err = dao.TokenDaoEntity.Set(&models.Token{
		Uid:         uid,
		ClientType:  uint32(clientType),
		AccessToken: token,
		MTime:       common.NowInt64MS(),
	}, true)
	if err != nil {
		return
	}

	return
}

// @Description generate token string.
func TokenString(t *jwt.Token) (s string) {
	if t != nil {
		if claims, ok := t.Claims.(*dao.MyCustomClaims); ok {
			s = fmt.Sprintf("[%d] isvalid %v,  ExpiresAt %s， NotBefore %s\n",
				claims.Uid,
				//TOKEN_TYPE_STR[claims.Type],
				t.Valid,
				time.Unix(claims.ExpiresAt, 0).Format("2006-01-02 15:04:05"),
				time.Unix(claims.NotBefore, 0).Format("2006-01-02 15:04:05"))
		}
	}

	return
}
