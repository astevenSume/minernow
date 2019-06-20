package controllers

import (
	"common"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	common2 "otc/common"
	controllers "otc_error"
	"strings"
	"time"
	admindao "utils/admin/dao"
	agentdao "utils/agent/dao"
	common3 "utils/common"
	otcdao "utils/otc/dao"
)

const (
	KeyCode        = "code"
	KeyAccessToken = "access_token"

	KeyData           = "data"
	KeyIdInput        = ":id"
	KeyTypeInput      = ":type"
	KeyId             = "id"
	KeyTimestampInput = ":timestamp"
	KeyTimestamp      = "timestamp"

	KeySystemInput = ":system"

	KeyClientTypeCookie = "client-type"
	KeyLang             = "lang"
	KeyInviteCode       = "invite_code"
	KeyInviteNum        = "invite_num"
	KeyInviteUrl        = "invite_url"
	KeyOrderIdInput     = ":order_id"

	KeyChannelId                 = "channel_id"
	KeyChannelIdInput            = ":channel_id"
	KeyCommission                = "commission"
	KeyBalance                   = "balance"
	KeyAmount                    = "amount"
	KeyYesterdayChips            = "yesterday_chips"
	KeyFee                       = "fee"
	KeyCtime                     = "ctime"
	KeyUtime                     = "utime"
	KeyStatus                    = "status"
	KeyUsdtTransferMethodAddress = "address"
	KeyUsdtTransferMethodMobile  = "mobile"

	KeyPage           = "page"
	KeyPerPage        = "per_page"
	KeyTotal          = "total"
	KeyMeta           = "meta"
	KeyList           = "list"
	KeyListAppChannel = "channel_list"
	KeyListAppType    = "type_list"
	KeyListApp        = "app_list"

	KeySign    = "sign"
	KeyToken   = "token"
	KeyGameUrl = "url"

	FileExtJpg = "jpg"
	FileExtPng = "png"
)

// redis lock key prefix
const (
	REDIS_LOCK_PREFIX_GAME = "game"
)

type MetaMsg struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func Init() (err error) {
	//init id manager
	err = common.IdMgrInit(otcdao.IdTypeMax)
	if err != nil {
		return
	}

	err = otcdao.Init(nil)
	if err != nil {
		return
	}

	// init cron function
	err = common.CronInit(&FunctionContainer{})
	if err != nil {
		return
	}

	// load chip to commission config
	/*err = Chip2CommissionConfigMgr.Load()
	if err != nil {
		return
	}

	err = WhiteListCommissionConfigMgr.Load()
	if err != nil {
		return
	}*/

	err = adminIpWhiteListMgr.load()
	if err != nil {
		return
	}

	err = common.RabbitMQInit(&AmqpFuncContainer{})
	if err != nil {
		return
	}

	return nil
}

func UpFileToOss(uid uint64, prefix string, fileData string) (url string, urlContent string, errCode controllers.ERROR_CODE) {
	if fileData == "" {
		errCode = controllers.ERROR_CODE_SUCCESS
		return
	}

	extIdx := strings.Index(fileData, ";")
	if extIdx <= 0 {
		errCode = controllers.ERROR_CODE_UP_FILE_EXT
		return
	}
	fileExt := fileData[:extIdx]
	if fileExt != FileExtJpg && fileExt != FileExtPng {
		errCode = controllers.ERROR_CODE_UP_FILE_EXT
		return
	}
	fileBody := fileData[extIdx+1:]

	//common.LogFuncDebug("fileBody:%v", fileBody)
	encodeBody, err := base64.StdEncoding.DecodeString(fileBody)
	if err != nil {
		errCode = controllers.ERROR_CODE_UP_FILE_FAIL
		common.LogFuncError("DecodeString fail:%v", err)
		return
	}

	imageFileName := fmt.Sprintf("%v_%v_%v.%v", prefix, uid, time.Now().Unix(), fileExt)
	err = ioutil.WriteFile(imageFileName, encodeBody, 0444)
	if err != nil {
		errCode = controllers.ERROR_CODE_UP_FILE_FAIL
		common.LogFuncError("error:%v", err)
		return
	}

	//二维码解析
	urlContent, err = common.ZBarImgDecode(imageFileName)
	if err != nil {
		errCode = controllers.ERROR_CODE_QR_CODE_DECODE_FAIL
		common.LogFuncError("error:%v", err)
		return
	}

	//上传二维码图片到OSS
	url, err = common.UpFile(common2.Cursvr.AccessKeyId, common2.Cursvr.AccessKeySecret, common2.Cursvr.Endpoint, common2.Cursvr.BucketName, imageFileName, imageFileName)
	if err != nil {
		errCode = controllers.ERROR_CODE_UP_FILE_FAIL
		common.LogFuncError("error:%v", err)
		return
	}

	result := os.Remove(imageFileName)
	if result != nil {
		common.LogFuncError("error:%v", result)
	}
	errCode = controllers.ERROR_CODE_SUCCESS

	return
}

//邀请码验证
func VerifyInviteCode(inviteCode, nationalCode, mobile string) controllers.ERROR_CODE {
	if admindao.TopAgentDaoEntity.IsTopAgent(nationalCode, mobile) {
		return controllers.ERROR_CODE_SUCCESS
	}

	_, err := agentdao.AgentPathDaoEntity.GetUidByInviteCode(inviteCode)
	if err != nil {
		return controllers.ERROR_CODE_USER_INVITE_CODE_ERROR
	}

	return controllers.ERROR_CODE_SUCCESS
}

//短信预警
func SmsWarning(wType int8, uid uint64, amount string) error {
	var key string
	switch wType {
	case admindao.ConfigWarningTypeUsdt:
		key = common3.UsdtTransferWarning
	default:
		return errors.New("config_warning type error")
	}

	strWarningValue, _ := common3.AppConfigMgr.String(key)
	if len(strWarningValue) > 0 {
		transferInteger, err := common.CurrencyStrToInt64(amount)
		if err != nil {
			common.LogFuncError("error:%v", err)
			return err
		}
		WarningValueInteger, err := common.CurrencyStrToInt64(strWarningValue)
		if err != nil {
			common.LogFuncError("error:%v", err)
			return err
		}
		//common.LogFuncInfo("transferInteger=%v, WarningValueInteger:%v", transferInteger, WarningValueInteger)
		if transferInteger > WarningValueInteger {
			mobile, err := otcdao.UserDaoEntity.GetMobileByUid(uid)
			if err != nil {
				common.LogFuncError("error:%v", err)
				return err
			}

			err = common3.SmsWarning(wType, map[string]string{
				"Mobile":   mobile,
				"TranNum":  amount,
				"LimitNum": strWarningValue,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
