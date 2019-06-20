package uemng_plus

import (
	"common"
	"fmt"
	"github.com/astaxie/beego"
	"umeng"
)

type UPushPlus struct {
	upush *umeng.UPush
}

const (
	GO_APP                = "go_app"                           //打开应用
	GO_URL                = "go_url"                           //跳转URL
	GO_ACTIVITY           = "go_activity"                      //打开特定的activity
	GO_CUSTOM             = "go_custom"                        //用户自定义内容
	GO_Order              = "go_order"                         //跳转订单页
	GO_Message            = "go_message"                       //跳转聊天页
	GO_Notice             = "go_notice"                        //系统通知
	Android_App_key       = "5cb0877d0cafb2d0860018ab"         //安卓的APP KEY
	Android_Master_Secret = "czlkqzd5ab2xathgiopig2k8mgiqw5vk" //安卓的 master_secret
	IOS_App_key           = "5caddd2a570df3438900046c"         //iOS的APP KEY
	IOS_Master_Secret     = "keijjrbzydg92aatqhlxiq7lvuu8bu0y" //iOS的 master_secret
	SysNotify             = "系统通知"
	Production            = true
)

//广播 ps:暂未使用
func (p *UPushPlus) Broadcast(PushType int, MsgType int, params map[string]interface{}, Platform int) {
	AppKey := beego.AppConfig.String("umeng_push::appkey")
	MasterSecret := beego.AppConfig.String("umeng_push::master_secret")
	fmt.Println(AppKey, MasterSecret)

	upush := &umeng.UPush{}
	if Platform == umeng.Platform_Android {
		upush = umeng.NewPush(umeng.PushType_Broadcast, umeng.MsgType_Notification, Android_App_key, Android_Master_Secret)
		upush.Body("ticker", params["ticker"])
		upush.Body("title", params["title"])
		upush.Body("text", params["text"])
		upush.Body("icon", "4")
		upush.Body("after_open", params["after_open"])
		upush.Body("play_vibrate", params["play_vibrate"])
		upush.Body("play_lights", params["play_lights"])
		upush.Body("play_sound", params["play_sound"])
		if params["production"] == false {
			upush.Mode(false)
		} else {
			upush.Mode(true)
		}
	} else if Platform == umeng.Platform_iOS {
		upush = umeng.NewPush(umeng.PushType_Broadcast, umeng.MsgType_Notification, IOS_App_key, IOS_Master_Secret)
		upush.APNs("alert", params["alert"])
		//upush.APNs("badge", body["badge"])
		upush.APNs("sound", params["sound"])
		//upush.APNs("content-available", body["content-available"])
		//upush.APNs("category", body["category"])
		//upush.Policy("expire_time", "2019-05-16 17:10:15")
		upush.Description("desc")
		if params["production"] == false {
			upush.Mode(false)
		} else {
			upush.Mode(true)
		}
	} else {
		return
	}

	res, err := upush.Push(Platform)
	fmt.Println(res)
	fmt.Println(err)
}

//系统通知广播
func BroadcastAll(text string, isProductionMode bool) (err error) {
	var result string
	//安卓
	upush := &umeng.UPush{}
	upush = umeng.NewPush(umeng.PushType_Broadcast, umeng.MsgType_Notification, Android_App_key, Android_Master_Secret)
	upush.Body("ticker", "ticker")
	upush.Body("title", SysNotify)
	upush.Body("text", text)
	upush.Body("after_open", GO_APP)
	upush.Body("play_vibrate", "true")
	upush.Body("play_lights", "true")
	upush.Body("play_sound", "true")
	upush.Mode(isProductionMode)
	result, err = upush.Push(umeng.Platform_Android)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	println("result:", result)

	//ios
	upush = umeng.NewPush(umeng.PushType_Broadcast, umeng.MsgType_Notification, IOS_App_key, IOS_Master_Secret)
	ans := make(map[string]interface{})
	alert := make(map[string]interface{})
	ans["alert"] = alert // content-available=1时可选 否则必填
	alert["title"] = SysNotify
	//alert["subtitle"] = SysNotify
	alert["body"] = text
	ans["sound"] = "default"
	ans["description"] = "Description"
	upush.Mode(isProductionMode)
	upush.APNs("alert", alert)
	result, err = upush.Push(umeng.Platform_iOS)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	println("result:", result)

	return
}

//单播 ps:暂未使用
func (p *UPushPlus) Unicast(DeviceToken string, params map[string]interface{}, Platform int) {
	upush := &umeng.UPush{}
	if Platform == umeng.Platform_Android {
		upush = umeng.NewPush(umeng.PushType_Unicast, umeng.MsgType_Notification, Android_App_key, Android_Master_Secret)
		upush.Token(DeviceToken)
		upush.Body("ticker", params["ticker"])
		upush.Body("title", params["title"])
		upush.Body("text", params["text"])
		upush.Body("icon", "4")
		upush.Body("after_open", params["after_open"])
		upush.Body("play_vibrate", params["play_vibrate"])
		upush.Body("play_lights", params["play_lights"])
		upush.Body("play_sound", params["play_sound"])
		upush.Body("after_open", params["after_open"])
		if params["production"] == false {
			upush.Mode(false)
		} else {
			upush.Mode(true)
		}
	} else if Platform == umeng.Platform_iOS {
		upush = umeng.NewPush(umeng.PushType_Unicast, umeng.MsgType_Notification, IOS_App_key, IOS_Master_Secret)
		upush.Token(DeviceToken)
		upush.APNs("alert", params["alert"])
		//upush.APNs("badge", body["badge"])
		upush.APNs("sound", params["sound"])
		//upush.APNs("content-available", body["content-available"])
		//upush.APNs("category", body["category"])
		//upush.Policy("expire_time", "2019-04-15 17:10:15")
		upush.Description("desc")
		if params["production"] == false {
			upush.Mode(false)
		} else {
			upush.Mode(true)
		}
	} else {
		return
	}
	res, err := upush.Push(Platform)
	fmt.Println(res)
	fmt.Println(err)
}

//自定义播
func (p *UPushPlus) Alias(typ, alias string, params map[string]interface{}, Platform int) {
	upush := &umeng.UPush{}
	if Platform == umeng.Platform_Android {
		upush = umeng.NewPush(umeng.PushType_Customizedcast, umeng.MsgType_Notification, Android_App_key, Android_Master_Secret)
		upush.Body("ticker", params["ticker"])
		upush.Body("title", params["title"])
		upush.Body("text", params["text"])
		upush.Body("icon", "4")
		upush.Body("after_open", params["after_open"])
		upush.Body("play_vibrate", params["play_vibrate"])
		upush.Body("play_lights", params["play_lights"])
		upush.Body("play_sound", params["play_sound"])
		upush.Body("uid_type", params["uid_type"])
		upush.Body("order_status", params["order_status"])
		upush.Body("order_id", params["order_id"])
		upush.Mode(true)
		upush.Description("desc")
		upush.Alias(typ, alias)
	} else if Platform == umeng.Platform_iOS {
		upush = umeng.NewPush(umeng.PushType_Customizedcast, umeng.MsgType_Notification, IOS_App_key, IOS_Master_Secret)
		upush.APNs("alert", params["alert"])
		upush.APNs("badge", params["badge"])
		upush.APNs("sound", params["sound"])
		upush.APNs("content-available", params["content-available"])
		upush.APNs("category", params["category"])
		//upush.Policy("expire_time", "2019-04-11 17:10:15")
		upush.Description("desc")
		upush.Mode(true)
		upush.Alias(typ, alias)
	} else {
		return
	}
	res, err := upush.Push(Platform)
	fmt.Println(res)
	fmt.Println(err)
}

// ios 指定用户推送
func (p *UPushPlus) iosPushByUid(uid string, params map[string]interface{}) (res string, err error) {
	upush := umeng.NewPush(umeng.PushType_Customizedcast, umeng.MsgType_Notification, IOS_App_key, IOS_Master_Secret)
	upush.Alias("Alias_UserId", uid)
	upush.Description("desc")
	upush.Mode(true)

	upush.APNs("badge", "badge")
	upush.APNs("sound", "default")
	upush.APNs("content-available", 0) //1为静默推送
	upush.APNs("category", "类别")
	for k, v := range params {
		upush.APNs(k, v)
	}

	for k, v := range params {
		upush.APNs(k, v)
	}
	res, err = upush.Push(umeng.Platform_iOS)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

// android 指定用户推送
func (p *UPushPlus) androidPushByUid(uid string, params map[string]interface{}) (res string, err error) {
	upush := umeng.NewPush(umeng.PushType_Customizedcast, umeng.MsgType_Notification, Android_App_key, Android_Master_Secret)
	upush.Alias("Alias_UserId", uid)
	upush.Description("desc")
	upush.Mode(true)

	upush.Body("ticker", "ticker")
	upush.Body("icon", "4")
	upush.Body("sound", "default")
	upush.Body("play_vibrate", "true") //true为震动
	upush.Body("play_lights", "false") //true为闪灯
	upush.Body("play_sound", "true")   //true发出声音

	for k, v := range params {
		upush.Body(k, v)
	}
	res, err = upush.Push(umeng.Platform_Android)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

func (p *UPushPlus) AliasInterface(uid, text, title, goToUrl, orderId, uidType string, orderStatus int8) {
	// 安卓推送
	var body map[string]interface{}
	body = make(map[string]interface{})
	body["ticker"] = "ticker"
	body["title"] = title
	body["text"] = text
	body["icon"] = "default"
	body["sound"] = "default"
	body["builder_id"] = 2
	body["after_open"] = GO_CUSTOM
	body["uid_type"] = uidType
	body["order_status"] = orderStatus
	body["order_id"] = orderId
	body["play_vibrate"] = "true" //true为震动
	body["play_lights"] = "false" //true为闪灯
	body["play_sound"] = "true"   //true发出声音
	body["production"] = "true"   //true 为正式环境 false为测试环境
	p.Alias("Alias_UserId", uid, body, umeng.Platform_Android)
	// iOS推送
	var ans map[string]interface{}
	ans = make(map[string]interface{})
	var alert map[string]interface{}
	alert = make(map[string]interface{})
	ans["alert"] = alert // content-available=1时可选 否则必填
	alert["title"] = title
	alert["subtitle"] = "" //字幕
	alert["body"] = text
	alert["after_open"] = goToUrl
	alert["order_id"] = orderId
	alert["uid_type"] = uidType
	alert["order_status"] = orderStatus
	ans["badge"] = "badge"
	ans["sound"] = "default"
	ans["content-available"] = 0 //1为静默推送
	ans["category"] = "类别"
	ans["production"] = true //false为测试
	ans["description"] = "详情"
	p.Alias("Alias_UserId", uid, ans, umeng.Platform_iOS)
}

// 订单推送
func (p *UPushPlus) PushOtcOrder(uid, orderId uint64, body, title string, uidType string, orderStatus int8) {
	p.AliasInterface(fmt.Sprintf("%d", uid), body, title, GO_Order, fmt.Sprintf("%d", orderId), uidType, orderStatus)
}

// 聊天消息
func (p *UPushPlus) PushMessage(uid, orderId uint64, body, title string, uidType string, orderStatus int8) {
	p.AliasInterface(fmt.Sprintf("%d", uid), body, title, GO_Message, fmt.Sprintf("%d", orderId), uidType, orderStatus)
}

// 系统消息通知
func (p *UPushPlus) PushSysNotice(uid uint64, content, title string) {
	uidstr := fmt.Sprintf("%d", uid)
	and := make(map[string]interface{})

	and["title"] = title
	and["text"] = content
	and["after_open"] = GO_Notice
	_, _ = p.androidPushByUid(uidstr, and)

	ios := make(map[string]interface{})
	alert := make(map[string]interface{})
	ios["alert"] = alert // content-available=1时可选 否则必填
	alert["title"] = title
	alert["subtitle"] = "" //字幕
	alert["body"] = content
	alert["after_open"] = GO_Notice

	_, _ = p.iosPushByUid(uidstr, ios)
}
