package gameapi

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego/orm"
	"time"

	"github.com/astaxie/beego"
	reportdao "utils/report/dao"
)

const (
	pageLimit   = 100
	ReqInterval = 12 * time.Second //12s
	ReqTimeSlot = 3600 * 1000      //时间段
)

// GameAPI 游戏对接 api
type GameAPI interface {
	// Login 登录接口 account:账号 password:密码 KindID:游戏ID
	Login(account, password, ip string, KindID string) (*LoginReply, error)
	Register(account, password, ip string) error
	Logout(account, password string)
	GetBalance(account, password string) (float64, error)
	TransferIn(account, password, orderNum string, money float64) (*TransferInReply, error)
	TransferOut(account, password, orderNum string, money float64) (*TransferOutReply, error)
	//GetBetRecords(pageIndex, pageSize int, startTime, endTime string) (*LotteryListReply, error)
	DayBetRecords(timestamp int64, mapWhite orm.Params) error

	//明细个数
	GetTotalByTimestamp(timestamp int64) (int, error)
	//明细
	GetRecordByTimestamp(timestamp int64, page, limit int) ([]reportdao.BetInfo, error)
	//删除明细
	DelBetRecord(timestamp int64) error
}
type LoginReply struct {
	URL string
}
type TransferInReply struct {
	Success bool
	Order   string
}
type TransferOutReply struct {
	Success bool
	Order   string
}

func md5hex(data ...[]byte) (string, error) {
	hash := md5.New()
	for _, d := range data {
		_, err := hash.Write(d)
		if err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

const (
	GAME_CHANNEL_KY  = uint32(iota + 1) //开元棋牌
	GAME_CHANNEL_AG                     //AG
	GAME_CHANNEL_RG                     //彩票
	GAME_CHANNEL_MAX                    //
)

var (
	supportedGameAPI map[uint32]GameAPI
)

func Init() error {
	supportedGameAPI = make(map[uint32]GameAPI)
	supportedGameAPI[GAME_CHANNEL_KY] = NewKaiYuanAPI(
		beego.AppConfig.String("ky_game::apiurl"),
		beego.AppConfig.String("ky_game::recordapiurl"),
		beego.AppConfig.String("ky_game::agent"),
		beego.AppConfig.String("ky_game::deskey"),
		beego.AppConfig.String("ky_game::md5key"),
		beego.AppConfig.String("ky_game::linecode"),
	)

	supportedGameAPI[GAME_CHANNEL_AG] = NewAsiaGamingAPI(
		beego.AppConfig.String("ag_game::apiurl"),
		beego.AppConfig.String("ag_game::prefix"),
		beego.AppConfig.String("ag_game::deskey"),
	)

	api, err := NewRoyalGameAPI(
		beego.AppConfig.String("rg_game::apiurl"),
		beego.AppConfig.String("rg_game::prefix"),
		beego.AppConfig.String("rg_game::deskey"),
	)
	if err != nil {
		return err
	}

	supportedGameAPI[GAME_CHANNEL_RG] = api
	return nil

}

// GetApi get api by game channel
func GetApi(gameChannel uint32) GameAPI {
	v, ok := supportedGameAPI[gameChannel]
	if ok {
		return v
	}

	return nil
}
