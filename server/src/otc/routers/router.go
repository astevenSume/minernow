package routers

import (
	"fmt"
	"otc/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

const (
	ROUTER_PAYMENT_METHODS   = "/payment-methods"
	ROUTER_USDT              = "/usdt"
	ROUTER_EOS               = "/eos"
	ROUTER_EXCHANGE          = "/exchange"
	ROUTER_ORDERS            = "/orders"
	ROUTER_GAME              = "/game"
	ROUTER_SYSTEM            = "/common"
	ROUTER_AGENT             = "/agent"
	ROUTER_AGENT_COMMMISSION = "/agent-commission"
	ROUTER_API               = "/api"
	ROUTER_BUY               = "/buy"
	ROUTER_SELL              = "/sell"
	ROUTER_APPS              = "/apps"
	ROUTER_BANNERS           = "/banners"
	ROUTER_PAYMENT_PASSWORD  = "/payment-password"
	ROUTER_ANNOUNCEMENTS     = "/announcements"
	MONTH_DIVIDEND           = "/monthdividend"
	PROFIT_REPORT            = "/profitreport"
)

func init() {
	//服务状态检查
	beego.InsertFilter("*", beego.BeforeExec, ServerStop)

	// 添加防重放过滤器
	beego.InsertFilter("*", beego.BeforeExec, AntiReplayFilter)

	ns := beego.NewNamespace("/v1",
		beego.NSBefore(before),
		beego.NSNamespace("/sms",
			beego.NSRouter("/sendcode", &controllers.SmsCodeController{}, "post:SendCode"),
			beego.NSRouter("/usersendcode", &controllers.SmsCodeController{}, "post:UserSendCode"),
		),
		beego.NSNamespace("/user",
			beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
			beego.NSRouter("/relogin", &controllers.UserController{}, "post:ReLogin"),
			beego.NSRouter("/sign", &controllers.UserController{}, "post:Sign"),
			beego.NSRouter("/captcha", &controllers.UserController{}, "get:Captcha"),
			beego.NSRouter("/logout", &controllers.UserController{}, "get:Logout"),
			beego.NSRouter("/info", &controllers.UserController{}, "get:GetInfo;put:EditInfo"),
			beego.NSRouter("/setting", &controllers.UserController{}, "get:Setting;post:SetSetting"),
			beego.NSRouter("/check_pwd", &controllers.UserController{}, "post:CheckPassWord"),
			beego.NSRouter("/update_mobile", &controllers.UserController{}, "post:CheckVerificationCode"),
			beego.NSRouter("/sendcode_newmobile", &controllers.UserController{}, "post:GetNewMobile"),
		),
		beego.NSNamespace(ROUTER_EOS,
			beego.NSRouter("", &controllers.EosController{}, "get:Info"),
			beego.NSRouter("/transfer", &controllers.EosController{}, "post:Transfer"),
			beego.NSRouter("/records", &controllers.EosController{}, "get:Records"),
			beego.NSRouter("/records/:id", &controllers.EosController{}, "get:RecordInfo"),
		),
		beego.NSNamespace(ROUTER_EXCHANGE,
			beego.NSRouter("/apply", &controllers.ExchangeController{}, "post:Apply"),
			beego.NSRouter("/info", &controllers.ExchangeController{}, "get:Info"),
			beego.NSRouter("/transferInto", &controllers.ExchangeController{}, "post:TransferInto"),
			beego.NSRouter("/transferOut", &controllers.ExchangeController{}, "post:TransferOut"),
		),
		beego.NSNamespace(ROUTER_BUY,
			beego.NSRouter("/start", &controllers.BuyController{}, "post:Start"),
			beego.NSRouter("", &controllers.BuyController{}, "post:Post"),
			beego.NSRouter("/:order_id/pay", &controllers.BuyController{}, "post:Pay"),
			beego.NSRouter("/:order_id/cancel", &controllers.BuyController{}, "post:Cancel"),
			beego.NSRouter("/:order_id/confirm", &controllers.BuyController{}, "post:Confirm"),
			beego.NSRouter("/list", &controllers.BuyController{}, "get:List"),
		),
		beego.NSNamespace(ROUTER_SELL,
			beego.NSRouter("/start", &controllers.SellController{}, "post:Start"),
			beego.NSRouter("", &controllers.SellController{}, "post:Post"),
			beego.NSRouter("/:order_id/pay", &controllers.SellController{}, "post:Pay"),
			beego.NSRouter("/:order_id/cancel", &controllers.SellController{}, "post:Cancel"),
			beego.NSRouter("/:order_id/confirm", &controllers.SellController{}, "post:Confirm"),
			beego.NSRouter("/list", &controllers.SellController{}, "get:List"),
		),
		beego.NSNamespace(ROUTER_ORDERS,
			beego.NSRouter("", &controllers.OrdersController{}, "get:Get"),
			beego.NSRouter("/:order_id", &controllers.OrdersController{}, "get:Info"),
			beego.NSRouter("/exchanger", &controllers.OrdersController{}, "get:Exchanger"),

			beego.NSRouter("/:order_id/appeal", &controllers.AppealController{}, "get:GetAppeal;post:CreateAppeal;put:ResolveAppeal"),

			beego.NSRouter("/:order_id/messages", &controllers.OtcMessageMethodController{}, "get:Get"),
			beego.NSRouter("/add_msg", &controllers.OtcMessageMethodController{}, "post:AddMsg"),
			beego.NSRouter("/common_messages", &controllers.OtcMessageMethodController{}, "get:CommonMessages"),
			beego.NSRouter("/system_notification", &controllers.SystemNotificationController{}, "get:GetSystemNotification"),
			beego.NSRouter("/system_notification/is_read", &controllers.SystemNotificationController{}, "get:GetIsReadNum"),
		),
		beego.NSNamespace(ROUTER_PAYMENT_PASSWORD,
			beego.NSRouter("/setpassword", &controllers.UserPayPassController{}, "post:SetPassword"),
			beego.NSRouter("/status", &controllers.UserPayPassController{}, "get:GetPayPwdStatus"),
			beego.NSRouter("/verifybypassword", &controllers.UserPayPassController{}, "post:VerifyByPassword"),
			beego.NSRouter("/verifybysign", &controllers.UserPayPassController{}, "post:VerifyBySign"),
			beego.NSRouter("/setverifystep", &controllers.UserPayPassController{}, "post:SetVerifyStep"),
		),
		beego.NSNamespace(ROUTER_PAYMENT_METHODS,
			beego.NSRouter("", &controllers.PaymentMethodsController{}, "get:Get;post:Bind"),
			beego.NSRouter("/:id", &controllers.PaymentMethodsController{}, "get:Get;put:Edit;delete:Unbind"),
			beego.NSRouter("/:id/activate", &controllers.PaymentMethodsController{}, "post:Activate"),
			beego.NSRouter("/:id/deactivate", &controllers.PaymentMethodsController{}, "post:Deactivate"),
			beego.NSRouter("/reorder", &controllers.PaymentMethodsController{}, "post:ReOrder"),
		),
		beego.NSRouter(ROUTER_USDT, &controllers.UsdtController{}, "get:Get"),
		beego.NSNamespace(ROUTER_USDT,
			beego.NSRouter("/balance", &controllers.UsdtController{}, "get:Balance"),
			beego.NSRouter("/transfer", &controllers.UsdtController{}, "post:Transfer"),
			beego.NSRouter("/mortgage", &controllers.UsdtController{}, "post:Mortgage"),
			beego.NSRouter("/release", &controllers.UsdtController{}, "post:Release"),
			beego.NSRouter("/records", &controllers.UsdtController{}, "get:Records"),
			beego.NSRouter("/records/:id", &controllers.UsdtController{}, "get:Record"),
			beego.NSRouter("/transfer/:id/cancel", &controllers.UsdtController{}, "post:CancelTransfer"),
			beego.NSRouter("/calculatefee", &controllers.UsdtController{}, "get:CalculateFee"),
		),
		beego.NSNamespace(ROUTER_SYSTEM,
			beego.NSRouter("", &controllers.SystemController{}, "get:Get"),
			beego.NSRouter("/timestamp", &controllers.SystemController{}, "get:SystemMSecTime"),
			beego.NSRouter("/price", &controllers.SystemController{}, "get:OtcExUsdtPrice"),
			beego.NSRouter(fmt.Sprintf("/%s/version", controllers.KeySystemInput), &controllers.SystemController{}, "get:LastAppVersion"),
			beego.NSRouter("/endpoint", &controllers.SystemController{}, "get:GetEndpoint"),
		),
		beego.NSNamespace(ROUTER_GAME,
			beego.NSRouter("/login", &controllers.GameController{}, "post:Login"),
			beego.NSRouter("/transferin", &controllers.GameController{}, "post:TransferIn"),
			beego.NSRouter("/transferout", &controllers.GameController{}, "post:TransferOut"),
			beego.NSRouter("/balance", &controllers.GameController{}, "get:GetBalance"),
			beego.NSRouter("/logout", &controllers.GameController{}, "post:Logout"),
			beego.NSNamespace("/orders",
				beego.NSRouter("", &controllers.GameController{}, "get:TransferList"),
				beego.NSRouter("/:id", &controllers.GameController{}, "get:TransferDetail"),
			),
		),

		beego.NSRouter(ROUTER_AGENT, &controllers.AgentController{}, "get:Query"),
		beego.NSNamespace(ROUTER_AGENT,
			//beego.NSRouter("/:channel_id", &controllers.AgentController{}, "get:QueryChannel"),
			beego.NSRouter("/withdraw", &controllers.AgentController{}, "post:Withdraw"),
			beego.NSRouter("/withdraws", &controllers.AgentController{}, "get:Withdraws"),
			beego.NSRouter("/sumdetail", &controllers.AgentController{}, "get:SumDetail"),
			beego.NSNamespace("/game",
				beego.NSRouter("/salary", &controllers.AgentController{}, "get:SalaryList"),
				//beego.NSRouter("/daysalary", &controllers.AgentController{}, "post:GetSalary"),
				beego.NSRouter("/daysalary/:timestamp", &controllers.AgentController{}, "get:DaySalary"),
				beego.NSRouter("/transferInfo", &controllers.AgentController{}, "get:TransferInfo"),
			),
		),

		beego.NSNamespace(ROUTER_API,
			beego.NSRouter("/distributecommission", &controllers.ApiController{}, "post:DistributeCommission"),
			beego.NSRouter("/calccommission", &controllers.ApiController{}, "post:CalcCommission"),
			beego.NSRouter("/cancelorder", &controllers.ApiController{}, "post:CancelOrder"),
			beego.NSRouter("/confirmorderpay", &controllers.ApiController{}, "post:ConfirmOrderPay"),
			beego.NSRouter("/confirmorder", &controllers.ApiController{}, "post:ConfirmOrder"),
			beego.NSRouter("/syncusdttransaction", &controllers.ApiController{}, "post:SyncUsdtTransaction"),
			beego.NSRouter("/approvetransferout", &controllers.ApiController{}, "post:ApproveTransferOut"),
			beego.NSRouter("/rejecttransferout", &controllers.ApiController{}, "post:RejectTransferOut"),
			beego.NSRouter("/eustrecharge", &controllers.ApiController{}, "post:EusdRecharge"),
			beego.NSRouter("/changeusdtstatus", &controllers.ApiController{}, "post:ChangeUsdtStatus"),
			beego.NSNamespace("/admin",
				beego.NSRouter("/task/list", &controllers.ApiController{}, "post:AdminTaskList"),
				beego.NSRouter("/task/detail", &controllers.ApiController{}, "post:AdminTaskList"),
			),
			//beego.NSRouter("/ping", &controllers.ApiController{}, "post:Ping"),
		),

		beego.NSNamespace(ROUTER_APPS,
			beego.NSRouter("/channel", &controllers.AppController{}, "get:Channels"),
			beego.NSRouter(fmt.Sprintf("/%s/list", controllers.KeyChannelIdInput), &controllers.AppController{}, "get:Apps"),
		),
		beego.NSNamespace(MONTH_DIVIDEND,
			beego.NSRouter("/needrechargeamount", &controllers.MonthDividendController{}, "get:NeedRechargeAmount"),
			beego.NSRouter("/recharge", &controllers.MonthDividendController{}, "post:Recharge"),
			beego.NSRouter("/list", &controllers.MonthDividendController{}, "get:List"),
			beego.NSRouter("/details", &controllers.MonthDividendController{}, "get:Details"),
		),
		beego.NSRouter(ROUTER_BANNERS, &controllers.BannerController{}, "get:Banners"),
		beego.NSNamespace(ROUTER_ANNOUNCEMENTS,
			beego.NSRouter("/:type", &controllers.AnnouncementController{}, "get:Announcements"),
		),
		beego.NSNamespace(PROFIT_REPORT,
			beego.NSRouter("/getpersonreports", &controllers.ProfitReportController{}, "get:GetProfitReports"),
		),
	)

	if beego.BConfig.RunMode == "dev" {
		nsTest := beego.NewNamespace("/v1",
			beego.NSNamespace("/test",
				beego.NSRouter("/becomeExchanger", &controllers.TestController{}, "get:BecomeExchanger"),
				beego.NSRouter("/eos", &controllers.TestController{}, "get:Eos"),
				beego.NSRouter("/test", &controllers.TestController{}, "get:Test"),
				beego.NSRouter("/info", &controllers.TestController{}, "get:Info"),
				beego.NSRouter("/tx", &controllers.TestController{}, "get:Tx"),
				beego.NSRouter("/otcclear", &controllers.TestController{}, "get:OtcClear"),
				beego.NSRouter("/syncrechargetx", &controllers.TestController{}, "post:UsdtSyncRechargeTransactionByMobile"),
				beego.NSRouter("/usdt", &controllers.TestController{}, "get:UsdtDeposit"),
				beego.NSRouter("/eusd2Mq", &controllers.TestController{}, "get:EusdTransfer2Mq"),
			),
		)
		beego.AddNamespace(nsTest)
	}

	beego.AddNamespace(ns)

	return
}

func before(ctx *context.Context) {
	//set output Content-Type to be json
	ctx.Output.Header("Content-Type", "application/json;charset=utf-8")
}
