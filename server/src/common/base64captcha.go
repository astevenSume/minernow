// example of HTTP server that uses the captcha package.
package common

import (
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"github.com/mojocn/base64Captcha"
	"time"
)

const (
	CaptchaActionLogin = "login"
)

var CaptchaConfig = base64Captcha.ConfigCharacter{
	Height:             60,
	Width:              130,
	Mode:               base64Captcha.CaptchaModeNumberAlphabet,
	IsUseSimpleFont:    true,
	ComplexOfNoiseDot:  10,
	ComplexOfNoiseText: 2,
	IsShowHollowLine:   true,
	IsShowNoiseDot:     true,
	IsShowNoiseText:    true,
	IsShowSineLine:     true,
	IsShowSlimeLine:    true,
	CaptchaLen:         4,
}

//customizeRdsStore An object implementing Store interface
type customizeRdsStore struct {
	Cache *redis.Client
}

// customizeRdsStore implementing Set method of  Store interface
func (s *customizeRdsStore) Set(id string, value string) {
	err := s.Cache.Set(id, value, time.Minute*10)
	if err != nil {

	}
}

// customizeRdsStore implementing Get method of  Store interface
func (s *customizeRdsStore) Get(id string, clear bool) (value string) {
	value, _ = s.Cache.Get(id).Result()
	LogFuncError("%v", value)
	if clear {
		_ = s.Cache.Del(id)
	}
	return
}

func CaptchaInit() error {
	//init redis store
	customeStore := customizeRdsStore{
		Cache: RedisManger,
	}
	w, _ := beego.AppConfig.Int("width")
	h, _ := beego.AppConfig.Int("height")
	m, _ := beego.AppConfig.Int("mode")
	l, _ := beego.AppConfig.Int("len")
	if w > 0 {
		CaptchaConfig.Width = w
	}
	if h > 0 {
		CaptchaConfig.Height = h
	}
	if m > 0 {
		CaptchaConfig.Mode = m
	}
	if l > 0 {
		CaptchaConfig.CaptchaLen = l
	}

	base64Captcha.SetCustomStore(&customeStore)
	return nil
}

func createCaptchaId(nationalCode, mobile string, action string) string {
	return "captcha_" + nationalCode + mobile + action
}

//获取验证码
func GetCaptcha(nationalCode, mobile string, action string) (base64Png string) {
	captchaId := createCaptchaId(nationalCode, mobile, action)
	_, image := base64Captcha.GenerateCaptcha(captchaId, CaptchaConfig)
	base64Png = base64Captcha.CaptchaWriteToBase64Encoding(image)

	return
}

//图片二进制流
func GetCaptchaPng(nationalCode, mobile string, action string) (imageByte []byte) {
	captchaId := createCaptchaId(nationalCode, mobile, action)
	_, image := base64Captcha.GenerateCaptcha(captchaId, CaptchaConfig)
	imageByte = image.BinaryEncodeing()

	return
}

//验证验证码
func VerifyCaptcha(nationalCode, mobile string, action string, value string) bool {
	captchaId := createCaptchaId(nationalCode, mobile, action)

	return base64Captcha.VerifyCaptcha(captchaId, value)
}
