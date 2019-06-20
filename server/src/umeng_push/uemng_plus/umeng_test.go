package uemng_plus

import (
	"fmt"
	"strconv"
	"testing"
	"umeng"
)

const UID = 181631811439296512

func TestAndroidBroadcast(t *testing.T) {
	p := new(UPushPlus)
	var body map[string]interface{}
	body = make(map[string]interface{})
	body["ticker"] = "ticker"  //通知栏文字
	body["title"] = "这是一个广播标题" //标题
	body["text"] = "这是广播消息内容"  //内容
	//body["icon"] = "default"
	//body["sound"] = "default"
	//body["builder_id"] = "default"
	//body["after_open"] = uemng_plus.GO_APP
	body["play_vibrate"] = "false"
	body["play_lights"] = "false"
	body["play_sound"] = "true"
	body["production"] = false
	p.Broadcast(umeng.PushType_Broadcast, umeng.MsgType_Notification, body, umeng.Platform_Android)
}

func TestAndroidBroadcastStatus(t *testing.T) {
	upush := umeng.NewPush(umeng.PushType_Broadcast, umeng.MsgType_Notification, Android_App_key, Android_Master_Secret)
	result, err := upush.Status("umqkpjl155797163569300")
	println("result:", result, "err:", err)
}

func TestIosBroadcast(t *testing.T) {
	p := new(UPushPlus)
	var ans map[string]interface{}
	ans = make(map[string]interface{})
	var alert map[string]interface{}
	alert = make(map[string]interface{})
	ans["alert"] = alert // content-available=1时可选 否则必填
	alert["title"] = "这是一个iOS广播标题"
	alert["subtitle"] = "subtitle"
	alert["body"] = "body"
	//ans["badge"] = "badge"
	ans["sound"] = "default"
	//ans["content-available"] = 0 //1为静默推送
	//ans["category"] = "category"
	ans["production"] = false //false为测试
	ans["description"] = "Description"
	p.Broadcast(umeng.PushType_Broadcast, umeng.MsgType_Notification, ans, umeng.Platform_iOS)
}

func TestIosBroadcastStatus(t *testing.T) {
	upush := umeng.NewPush(umeng.PushType_Broadcast, umeng.MsgType_Notification, IOS_App_key, IOS_Master_Secret)
	result, err := upush.Status("umozhjj155797163579700")
	println("result:", result, "err:", err)
}

func TestAliasInterface(t *testing.T) {
	p := new(UPushPlus)
	uid := 181631811439296512 //169891261749133312
	body := "测试：您订单号为『1276423986』，数量为『200』个，总价为『%s』CNY的订单对方已完成支付，请及时查账，确认到账后请及时放币，若未收到款项，可申请客诉处理。"
	title := "买家已付款!!"
	p.AliasInterface(string(strconv.Itoa(int(uid))), body, title, "go order", "1234", "0", 4)
}

func TestBroadcastAll(t *testing.T) {
	err := BroadcastAll("广播测试", true)
	if err != nil {
		println("err:", err)
	}
}

func TestPushSysNotice(t *testing.T) {
	p := new(UPushPlus)
	p.PushSysNotice(UID, "系统通知测试", "系统通知title")
}

func TestIosPushByUid(t *testing.T) {
	p := new(UPushPlus)
	ios := make(map[string]interface{})
	alert := make(map[string]interface{})
	ios["alert"] = alert // content-available=1时可选 否则必填
	alert["title"] = "IOS Test"
	alert["subtitle"] = "" //字幕
	alert["body"] = "IOS Push Test"
	alert["after_open"] = GO_Notice
	p.iosPushByUid(fmt.Sprintf("%d", UID), ios)
}
