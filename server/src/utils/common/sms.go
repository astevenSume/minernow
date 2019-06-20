package common

import (
	"bytes"
	"common"
	"errors"
	"fmt"
	"math/rand"
	. "otc_error"
	"strconv"
	"text/template"
	"time"
	admindao "utils/admin/dao"

	"github.com/astaxie/beego"
)

const (
	SmsExpire          = 600 * time.Second //短信验证码有效时间
	SmsCodeLen         = 6
	SmsSendInterval    = 60
	SmsVerifyFailTimes = 5                   // 短信验证可以失败的次数，超过需要换验证码
	SmsVerifyExpire    = 86400 * time.Second //短信验证码失败次数统计有效时间
	NatinalChina       = "86"

	SmsActionOverSea      = "oversea"
	SmsActionLogin        = "login"
	SmsActionPayment      = "payment"
	SmsActionPaySecond    = "2-step"
	SmsAcitionPayPassword = "change-payment-password" // 设置或修改支付密码
)

var (
	ErrSmsError    = errors.New("sms error")
	ErrSmsOutTimes = errors.New("out fail times")
)

//ali短信服务
func SendSms(accessKeyID, accessSecret, signName, smsChina, smsOverSea, nationalCode, mobile, action string) (string, ERROR_CODE) {
	//common.LogFuncDebug("keyID:%v, accessSecret:%v, signName:%v", accessKeyID, accessSecret, signName)
	keyMobile := smsMobileKey(nationalCode, mobile)
	ok := common.RedisSetNX(keyMobile, SmsSendInterval)
	if !ok {
		// 60s 默认发送成功
		if action == SmsActionLogin {
			return "", ERROR_CODE_SUCCESS
		} else {
			return "", ERROR_CODE_SMS_SEND_REQ_TOO_MUCH
		}
	}
	key := smsRedisKey(nationalCode, mobile, action)
	rand.Seed(time.Now().UnixNano())
	//生成短信验证码,有效时间600s
	sms := common.RandomNum(SmsCodeLen)
	err := common.RedisManger.Set(key, sms, SmsExpire).Err()
	if err != nil {
		return sms, ERROR_CODE_REDIS
	}

	common.LogFuncDebug("mobile:%s, sms:%s", mobile, sms)

	phoneNum := mobile
	jsonParam := fmt.Sprintf("{\"code\":\"%s\"}", sms)
	templateCode := smsChina
	if nationalCode != NatinalChina {
		templateCode = smsOverSea
		phoneNum = fmt.Sprintf("%v%v", nationalCode, mobile)
	}
	err = common.SendSms(accessKeyID, accessSecret, phoneNum, signName, jsonParam, templateCode)
	if err != nil {
		common.LogFuncError("sendalisms fail err:%v", err)
		return sms, ERROR_CODE_SMS_SEND_FAIL
	}

	id, err := admindao.SmsCodeDaoEntity.InsertSmsCode(nationalCode, mobile, action, sms, int64(SmsExpire))
	if err != nil {
		return sms, ERROR_CODE_DB
	}
	keyId := smsCodeIdRedisKey(nationalCode, mobile, action)
	err = common.RedisManger.Set(keyId, id, SmsExpire).Err()
	if err != nil {
		return sms, ERROR_CODE_REDIS
	}

	return sms, ERROR_CODE_SUCCESS
}

//国外ali账号短信服务
func AliSendSms(nationalCode, mobile, action string) (string, ERROR_CODE) {
	var err error
	defer func() {
		if err != nil {
			common.LogFuncError("%v", err)
		}
	}()
	accessKeyID, err := AppConfigMgr.String(SmsAccessAppId)
	if err != nil {
		return "", ERROR_CODE_CONFIG_LACK
	}

	accessSecret, err := AppConfigMgr.String(SmsAccessSecret)
	if err != nil {
		return "", ERROR_CODE_CONFIG_LACK
	}

	signName, err := AppConfigMgr.String(SmsSignName)
	if err != nil {
		return "", ERROR_CODE_CONFIG_LACK
	}

	smsChina, err := AppConfigMgr.String(SmsChina)
	if err != nil {
		return "", ERROR_CODE_CONFIG_LACK
	}

	// todo: 上线删除
	if beego.BConfig.RunMode == "dev" {
		return "", ERROR_CODE_SUCCESS
	}

	keyMobile := smsMobileKey(nationalCode, mobile)
	ok := common.RedisSetNX(keyMobile, SmsSendInterval)
	if !ok {
		// 60s 默认发送成功
		if action == SmsActionLogin {
			return "", ERROR_CODE_SUCCESS
		} else {
			return "", ERROR_CODE_SMS_SEND_REQ_TOO_MUCH
		}
	}
	key := smsRedisKey(nationalCode, mobile, action)
	rand.Seed(time.Now().UnixNano())
	//生成短信验证码,有效时间600s
	sms := common.RandomNum(SmsCodeLen)
	err = common.RedisManger.Set(key, sms, SmsExpire).Err()
	if err != nil {
		return sms, ERROR_CODE_REDIS
	}
	common.LogFuncDebug("mobile:%s, sms:%s", mobile, sms)

	phoneNum := fmt.Sprintf("%v%v", nationalCode, mobile)
	if nationalCode != NatinalChina {
		msg, err := getSmsParam(admindao.SmsTemplateTypeVerifyCode, map[string]string{"Code": sms})
		if err != nil {
			common.LogFuncError("sendsms:%v", err)
			return sms, ERROR_CODE_SMS_SEND_FAIL
		}
		err = common.NewSendSmsOverSea(accessKeyID, accessSecret, phoneNum, msg)
	} else {
		jsonParam := fmt.Sprintf("{\"code\":\"%s\"}", sms)
		err = common.NewSendSmsChina(accessKeyID, accessSecret, phoneNum, signName, jsonParam, smsChina)
	}
	if err != nil {
		common.LogFuncError("sendalisms fail err:%v", err)
		return sms, ERROR_CODE_SMS_SEND_FAIL
	}

	id, err := admindao.SmsCodeDaoEntity.InsertSmsCode(nationalCode, mobile, action, sms, int64(SmsExpire))
	if err != nil {
		return sms, ERROR_CODE_DB
	}
	keyId := smsCodeIdRedisKey(nationalCode, mobile, action)
	err = common.RedisManger.Set(keyId, id, SmsExpire).Err()
	if err != nil {
		return sms, ERROR_CODE_REDIS
	}

	return sms, ERROR_CODE_SUCCESS
}

//aws短信服务
func SendAwsSms(nationalCode, mobile, action string, smsType int8, param map[string]string) (sms string, err error) {
	key := smsRedisKey(nationalCode, mobile, action)
	rand.Seed(time.Now().UnixNano())
	//生成短信验证码,有效时间600s
	sms = common.RandomNum(SmsCodeLen)
	err = common.RedisManger.Set(key, sms, SmsExpire).Err()
	if err != nil {
		return
	}

	//短信发送
	var msg string
	msg, err = getSmsParam(smsType, map[string]string{"Code": sms})
	if err != nil {
		common.LogFuncError("sendsms:%v", err)
		return
	}

	phoneNum := fmt.Sprintf("%v%v", nationalCode, mobile)
	err = Send(msg, phoneNum)
	if err != nil {
		return
	}

	id, err := admindao.SmsCodeDaoEntity.InsertSmsCode(nationalCode, mobile, action, sms, int64(SmsExpire))
	if err != nil {
		return
	}
	keyId := smsCodeIdRedisKey(nationalCode, mobile, action)
	err = common.RedisManger.Set(keyId, id, SmsExpire).Err()
	if err != nil {
		return
	}

	return
}

func getSmsParam(smsType int8, param map[string]string) (msg string, err error) {
	tpl, err := admindao.SmsTemplateDaoEntity.GetSmsTemplates(smsType)
	if err != nil {
		common.LogFuncError("err=%v", err)
		return
	}
	t := template.New("sms template")
	t, err = t.Parse(tpl)
	if err != nil {
		common.LogFuncError("err=%v", err)
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, param)
	if err != nil {
		common.LogFuncError("sendsms error:%v", err)
		return
	}
	msg = buf.String()
	return
}

func smsMobileKey(nationalCode, mobile string) string {
	return "sms_" + nationalCode + mobile
}

func smsRedisKey(nationalCode, mobile, action string) string {
	return "sms_" + action + nationalCode + mobile
}

func smsVerifyTimesRedisKey(nationalCode, mobile, action string) string {
	return "sms_verify_" + action + nationalCode + mobile
}

func smsCodeIdRedisKey(nationalCode, mobile, action string) string {
	return "sms_id_" + action + nationalCode + mobile
}

func VerifySmsCode(nationalCode, mobile, action, sms string) (res bool, err error) {
	if beego.AppConfig.String("runmode") == "dev" {
		//todo 开发模式下万能码 测试用
		if sms == "000000" {
			res = true
			return
		}
	}

	res = false
	if sms == "" {
		err = ErrSmsError
		return
	}

	key := smsRedisKey(nationalCode, mobile, action)
	timesKey := smsVerifyTimesRedisKey(nationalCode, mobile, action)
	smsRedis, _ := common.RedisManger.Get(key).Result()

	if smsRedis != sms {
		err = ErrSmsError
		times, _ := common.RedisManger.Incr(timesKey).Result()
		if times >= SmsVerifyFailTimes {
			err = ErrSmsOutTimes
			//失败太多，清空验证码 & 次数统计
			_ = common.RedisManger.Del(key, timesKey)
		} else if times == 1 {
			common.RedisManger.Expire(timesKey, SmsVerifyExpire)
		}
		return
	}
	//验证成功直接删除key
	keyId := smsCodeIdRedisKey(nationalCode, mobile, action)
	keyIdRedis, err := common.RedisManger.Get(keyId).Result()
	if err == nil {
		//取到ID,设置验证码状态
		kId, err := strconv.Atoi(keyIdRedis)
		if err == nil {
			err = admindao.SmsCodeDaoEntity.SetSmsCodeUsed(int64(kId))
			if err != nil {
				common.LogFuncError("SetSmsCodeUsed error:%v", err)
				err = nil
			}
		}
	}

	_ = common.RedisManger.Del(key, timesKey, keyId)
	res = true

	return
}

func GenerateNationalCode(sourceCode string) string {
	return "+" + sourceCode
}
