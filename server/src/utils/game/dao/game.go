package dao

import (
	"common"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	//"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	json "github.com/mailru/easyjson"
	"io/ioutil"
	otccommon "otc/common"
	"time"
)

const (
	PROTO_CHANNEL_DAILY  = "sdk/callback_channeldaily"
	PROTO_USER_DAILY     = "sdk/callback_userdaily"
	PROTO_REGISTER       = "sdk/register"
	PROTO_LOGIN          = "sdk/login"
	PROTO_LOGOUT         = "sdk/logout"
	PROTO_REDEEM_IN      = "sdk/callback_redeemin"
	PROTO_REDEEM_OUT     = "sdk/callback_redeemout"
	PROTO_REDEEM_CONFIRM = "sdk/callback_redeemconfirm"
	PROTO_BALANCE        = "sdk/callback_balance"
	PROTO_ONLINE_PLAYERS = "sdk/callback_onlineplayers"

	GAMEID = 315 //game id will never be changed
)

const (
	KeyGameId    = "gameid"
	KeyPlatId    = "platid"
	KeyChannelId = "channelid"
	KeySign      = "sign"
	KeyAccount   = "account"
	KeySubplatid = "subplatid"
	KeyStart     = "start"
	KeyEnd       = "end"
	KeyCurpage   = "curpage"
	KeyPerpage   = "perpage"
	KeyStatus    = "status"
	KeyDesc      = "desc"
	KeyData      = "data"
	KeyMaxpage   = "maxpage"

	KeyUserid       = "userid"
	KeyWinlosemoney = "winlosemoney"
	KeyChips        = "chips"
	KeyTax          = "tax"
)

type GameDao struct {
	common.BaseDao
}

func NewGameDao(db string) *GameDao {
	return &GameDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var GameDaoEntity *GameDao

func (d *GameDao) SendToChannel(host, proto string, platId uint32, msg json.Marshaler, resp json.Unmarshaler) (err error) {
	//body
	var buf []byte
	buf, err = json.Marshal(msg)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	body := base64.URLEncoding.EncodeToString(buf)

	//sign
	h := md5.New()
	h.Write(append(buf, []byte(otccommon.Cursvr.ChannelAppKey)...))
	sign := hex.EncodeToString(h.Sum(nil))

	// request
	// discussion about game id, as follow：
	// "你直接写死315 即可，这里不允许修改。"
	// "生产环境也是写死吗？"
	// "是的。"
	url := fmt.Sprintf("%s/%s?%s=%v&%s=%v&%s=%s", host, proto, KeyGameId, GAMEID, KeyPlatId, platId, KeySign, sign)
	req := httplib.Post(url)
	req.Body(body)
	common.LogFuncDebug("url [%s] \n body [%s]", url, body)
	//err = req.ToJSON(resp)
	var data []byte
	data, err = req.Bytes()
	if err != nil {
		common.LogFuncError("%v", err)
		return err
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	common.LogFuncDebug("\nreq [%+v]\nresp [%+v]", msg, resp)

	ioutil.WriteFile("game.log", []byte(fmt.Sprintf("%s : url [%s] \n body [%s] \n req [%s] \n resp [%s] \n",
		time.Now().String(), url, body, req, resp)), 0666)

	return
}
