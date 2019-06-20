package common

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

const SmsResponseOk = "OK"

//-----------------------  阿里云短信服务api代码 begin --------------------------
//SendSmsReply 发送短信返回
type SendSmsReply struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
}

func replace(in string) string {
	rep := strings.NewReplacer("+", "%20", "*", "%2A", "~", "%7E")
	return rep.Replace(url.QueryEscape(in))
}

// SendSms 发送短信
// eg: go SendSms(Cursvr.SmsAccessAppId, Cursvr.SmsAccessSecret, nationalCode+mobile, Cursvr.SmsSignName, fmt.Sprintf("{\"code\":\"%s\",\"expired\":%d}", code, Cursvr.CaptchaExpiredSecs/60), templateCode)
func SendSms(accessKeyID, accessSecret, phoneNumbers, signName, templateParam, templateCode string) error {
	paras := map[string]string{
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   fmt.Sprintf("%d", rand.Int63()),
		"AccessKeyId":      accessKeyID,
		"SignatureVersion": "1.0",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Format":           "JSON",

		"Action":        "SendSms",
		"Version":       "2017-05-25",
		"RegionId":      "cn-hangzhou",
		"PhoneNumbers":  phoneNumbers,
		"SignName":      signName,
		"TemplateParam": templateParam,
		"TemplateCode":  templateCode,
	}

	var keys []string

	for k := range paras {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var sortQueryString string

	for _, v := range keys {
		sortQueryString = fmt.Sprintf("%s&%s=%s", sortQueryString, replace(v), replace(paras[v]))
	}

	stringToSign := fmt.Sprintf("GET&%s&%s", replace("/"), replace(sortQueryString[1:]))

	mac := hmac.New(sha1.New, []byte(fmt.Sprintf("%s&", accessSecret)))
	mac.Write([]byte(stringToSign))
	sign := replace(base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	str := fmt.Sprintf("http://dysmsapi.aliyuncs.com/?Signature=%s%s", sign, sortQueryString)

	resp, err := http.Get(str)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ssr := &SendSmsReply{}

	if err := json.Unmarshal(body, ssr); err != nil {
		return err
	}

	if ssr.Code == "SignatureNonceUsed" {
		return SendSms(accessKeyID, accessSecret, phoneNumbers, signName, templateParam, templateCode)
	} else if ssr.Code != "OK" {
		return errors.New(ssr.Code)
	}

	return nil
}

const aliParamLimit int = 20

// @Description format sms params
// ali require single param length no more than 20 bytes
func SmsParamsFormat(in []string) (out []interface{}) {
	for _, param := range in {
		//use rune in case of cut one Chinese character to two
		runeIn := []rune(param)
		length := 0
		runeOut := []rune{}
		for _, r := range runeIn {
			length += utf8.RuneLen(r)
			runeOut = append(runeOut, r)
			if length >= aliParamLimit {
				break
			}
		}
		outParam := string(runeOut)
		out = append(out, outParam)
	}

	return
}

//SendSmsReply 发送短信返回
type NewSendSmsReply struct {
	ResponseCode string `ResponseCode:"Code,omitempty"`
}

// NewSendSmsChina 国外阿里云账号发送国内短信,国内短信必须使用模板
// eg: go SendSms(Cursvr.SmsAccessAppId, Cursvr.SmsAccessSecret, nationalCode+mobile, Cursvr.SmsSignName, fmt.Sprintf("{\"code\":\"%s\",\"expired\":%d}", code, Cursvr.CaptchaExpiredSecs/60), templateCode)
func NewSendSmsChina(accessKeyID, accessSecret, phoneNumbers, signName, templateParam, templateCode string) error {
	client, err := sdk.NewClientWithAccessKey("ap-southeast-1", accessKeyID, accessSecret)
	if err != nil {
		LogFuncError("error:%v", err)
		return err
	}

	request := requests.NewCommonRequest() // 构造一个公共请求
	request.AcceptFormat = "json"
	request.Method = "POST"                                 // 设置请求方式
	request.Domain = "dysmsapi.ap-southeast-1.aliyuncs.com" // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2018-05-01"                          // 指定产品版本
	request.ApiName = "SendMessageWithTemplate"             // 指定接口名

	request.QueryParams["To"] = phoneNumbers             // 发送号码
	request.QueryParams["From"] = signName               // 签名
	request.QueryParams["TemplateCode"] = templateCode   // 指定请求的区域，不指定则使用客户端区域、默认区域
	request.QueryParams["TemplateParam"] = templateParam // 指定请求的区域，不指定则使用客户端区域、默认区域

	request.TransToAcsRequest()
	response := responses.NewCommonResponse()

	err = client.DoAction(request, response)
	if err != nil {
		LogFuncError("error:%v", err)
		return err
	}

	ssr := &NewSendSmsReply{}
	err = json.Unmarshal([]byte(response.GetHttpContentString()), ssr)
	if err != nil {
		LogFuncError("Umarshal failed:%v", err)
		return err
	}

	if ssr.ResponseCode != SmsResponseOk {
		LogFuncError("SendFail:%v", ssr.ResponseCode)
		return err
	}

	return nil
}

// NewSendSmsOverSea 国外阿里云账号发送国外短信
// eg: go SendSms(Cursvr.SmsAccessAppId, Cursvr.SmsAccessSecret, nationalCode+mobile, Cursvr.SmsSignName, fmt.Sprintf("{\"code\":\"%s\",\"expired\":%d}", code, Cursvr.CaptchaExpiredSecs/60), templateCode)
func NewSendSmsOverSea(accessKeyID, accessSecret, phoneNumbers, msg string) error {
	client, err := sdk.NewClientWithAccessKey("ap-southeast-1", accessKeyID, accessSecret)
	if err != nil {
		LogFuncError("error:%v", err)
		return err
	}

	request := requests.NewCommonRequest() // 构造一个公共请求
	request.AcceptFormat = "json"
	request.Method = "POST"                                 // 设置请求方式
	request.Domain = "dysmsapi.ap-southeast-1.aliyuncs.com" // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2018-05-01"                          // 指定产品版本
	request.ApiName = "SendMessageToGlobe"                  // 指定接口名

	request.QueryParams["To"] = phoneNumbers // 电话号码
	//request.QueryParams["From"] = ""  // 指定请求的区域，不指定则使用客户端区域、默认区域
	request.QueryParams["Message"] = msg // 消息

	request.TransToAcsRequest()
	response := responses.NewCommonResponse()

	err = client.DoAction(request, response)
	if err != nil {
		LogFuncError("error:%v", err)
		return err
	}

	ssr := &NewSendSmsReply{}
	err = json.Unmarshal([]byte(response.GetHttpContentString()), ssr)
	if err != nil {
		LogFuncError("Umarshal failed:%v", err)
		return err
	}

	if ssr.ResponseCode != SmsResponseOk {
		LogFuncError("SendFail:%v", ssr.ResponseCode)
		return err
	}

	return nil
}
