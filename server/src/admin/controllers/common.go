package controllers

import (
	"admin/controllers/errcode"
	"common"
	"common/systoolbox"
	"errors"
	"eusd/eosplus"
	"fmt"
	otcerror "otc_error"
	"regexp"
	"strconv"
	"strings"
	"usdt"
	"utils/admin/dao"
	"utils/admin/models"
	agentdao "utils/agent/dao"
	common2 "utils/common"
	otcDao "utils/otc/dao"
	otcModels "utils/otc/models"
	usdtDao "utils/usdt/dao"
	usdtmodels "utils/usdt/models"

	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
)

const (
	DEFAULT_PER_PAGE = 10
	MaxPerPage       = 20
)
const (
	KEY_ID_INPUT           = ":id"
	KEY_SYSTEM_INPUT       = ":system"
	KEY_KEY_ACTION         = ":action"
	KEY_START_TIME         = "stime"
	KEY_END_TIME           = "etime"
	KEY_ACTION             = "action"
	KEY_UID                = "uid"
	KEY_ID                 = "id"
	KEY_STATUS             = "status"
	KEY_EXCHANGER          = "exchanger"
	KEY_PAGE               = "page"
	KEY_LIMIT              = "limit"
	KEY_TOTAL              = "total"
	KEY_LIST               = "list"
	KEY_META               = "meta"
	KEY_TYPE               = "type"
	KEY_MOBILE             = "mobile"
	KeySign                = "sign"
	Key_EUID               = "euid"
	KeyAppName             = "appname"
	KeyRegionId            = "region_id"
	KeyServerId            = "server_id"
	KEY_GAME_TRANSFER_TYPE = "transfertype"

	RouterDistributeCommission = "distributecommission"
	RouterCalcCommission       = "calccommission"
	RouterCancelOrder          = "cancelorder"
	RouterConfirmOrderPay      = "confirmorderpay"
	RouterConfirmOrder         = "confirmorder"
	RouterSyncUsdtTransaction  = "syncusdttransaction"
	RouterApproveTransferOut   = "approvetransferout"
	RouterRejectTransferOut    = "rejecttransferout"
	RouterEustRecharge         = "eustrecharge"
	RouterChangeUsdtStatus     = "changeusdtstatus"

	NatinalChina = "86"
)

/*
操作日志action定义,往后加通知前端保持一致
urlMethod  flag
get	       Read
put 	   Edit
post	   Add
delete     Del
*/
const (
	OPActionReadEosRecord                   = 1   //查询EOS记录
	OPActionReadEosWealth                   = 2   //查询EOS钱包
	OPActionReadEosAddress                  = 3   //查询EOS地址
	OPActionReadOrderMsg                    = 4   //查询订单聊天内容
	OPActionReadPermission                  = 5   //查询权限
	OPActionEditPermission                  = 6   //编辑权限
	OPActionAddPermission                   = 7   //添加权限
	OPActionDelPermission                   = 8   //删除权限
	OPActionReadRole                        = 9   //查询角色
	OPActionAddRole                         = 10  //创建角色
	OPActionEditRole                        = 11  //编辑角色
	OPActionDelRole                         = 12  //删除角色
	OPActionReadRolePermissions             = 13  //获取角色权限
	OPActionEditRolePermissions             = 14  //添加角色权限
	OPActionDelRolePermissions              = 15  //删除角色权限
	OPActionReadRoleMember                  = 16  //查询角色成员
	OPActionReadAdmin                       = 17  //查询管理员
	OPActionAddAdmin                        = 18  //创建管理员
	OPActionEditAdmin                       = 19  //编辑管理员
	OPActionDelAdmin                        = 20  //删除管理员
	OPActionAddSuspendAdmin                 = 21  //禁用管理员
	OPActionAddRestoreAdmin                 = 22  //恢复管理员
	OPActionReadAdminsRole                  = 23  //查询管理员角色
	OPActionAddAdminsRole                   = 24  //添加管理员角色
	OPActionDelAdminsRole                   = 25  //删除管理员角色
	OPActionReadAgentWhiteList              = 26  //查询代理白名单档位
	OPActionAddAgentWhiteList               = 27  //创建代理白名单档位
	OPActionEditAgentWhiteList              = 28  //编辑代理白名单档位
	OPActionDelAgentWhiteList               = 29  //删除代理白名单档位
	OPActionReadAgentWhiteListAgent         = 30  //查询代理白名单档位代理
	OPActionEditAgentWhiteListAgent         = 31  //编辑代理白名单档位代理
	OPActionDelAgentWhiteListAgent          = 32  //删除代理白名单档位代理
	OPActionReadAppChannel                  = 33  //查询应用渠道
	OPActionAddAppChannel                   = 34  //创建应用渠道
	OPActionEditAppChannel                  = 35  //编辑应用渠道
	OPActionDelAppChannel                   = 36  //删除应用渠道
	OPActionReadAppType                     = 37  //查询应用分类
	OPActionAddAppType                      = 38  //创建应用分类
	OPActionEditAppType                     = 39  //编辑应用分类
	OPActionDelAppType                      = 40  //删除应用分类
	OPActionReadVersion                     = 41  //查询应用版本
	OPActionAddVersion                      = 42  //创建应用版本
	OPActionEditVersion                     = 43  //编辑应用版本
	OPActionDelVersion                      = 44  //删除应用版本
	OPActionAddVersionPublish               = 45  //上架应用版本
	OPActionAddVersionUnPublish             = 46  //下架应用版本
	OPActionReadAppeal                      = 47  //查询申述
	OPActionEditAppeal                      = 48  //编辑申述
	OPActionReadAppealService               = 49  //查询申述客服
	OPActionAddAppealService                = 50  //创建申述客服
	OPActionEditAppealService               = 51  //编辑申述客服
	OPActionDelAppealService                = 52  //删除申述客服
	OPActionDelAppealServiceSuspend         = 53  //禁用申述客服
	OPActionAddAppealServiceRestore         = 54  //恢复申述客服
	OPActionReadApp                         = 55  //查询第三方应用
	OPActionAddApp                          = 56  //创建第三方应用
	OPActionEditApp                         = 57  //编辑第三方应用
	OPActionDelApp                          = 58  //删除第三方应用
	OPActionAddAppFeature                   = 59  //推荐第三方应用
	OPActionAddAppUnFeature                 = 60  //推荐第三方应用
	OPActionAddAppPublish                   = 61  //上线应用
	OPActionAddAppUnPublish                 = 62  //取消上线应用
	OPActionReadBanner                      = 63  //查询广告
	OPActionAddBanner                       = 64  //创建广告
	OPActionEditBanner                      = 65  //编辑广告
	OPActionDelBanner                       = 66  //删除广告
	OPActionAddBannerPublish                = 67  //发布广告
	OPActionAddBannerUnPublish              = 68  //取消发布广告
	OPActionReadCommissionStat              = 69  //查询佣金统计
	OPActionAddCommissionDistribute         = 70  //发放佣金
	OPActionAddCommissionCalc               = 71  //计算佣金
	OPActionReadCommissionRate              = 72  //查询返佣等级
	OPActionEditCommissionRate              = 73  //编辑返佣等级
	OPActionReadConfig                      = 74  //查询缓存配置
	OPActionAddConfig                       = 75  //创建缓存配置
	OPActionEditConfig                      = 76  //编辑缓存配置
	OPActionDelConfig                       = 77  //删除缓存配置
	OPActionEditConfigRefresh               = 78  //刷新缓存
	OPActionDelConfigClean                  = 79  //清除缓存
	OPActionAddEusdTransfer                 = 80  //eusd转账
	OPActionAddEusdRetire                   = 81  //eusd退款
	OPActionAddEusdBalance                  = 82  //eusd余额
	OPActionAddLogin                        = 83  //登录
	OPActionAddReLogin                      = 84  //重新登录
	OPActionAddLogout                       = 85  //登出
	OPActionAddPassword                     = 86  //修改密码
	OPActionReadExchanger                   = 87  //查询承兑商
	OPActionAddExchanger                    = 88  //创建承兑商
	OPActionEditExchanger                   = 89  //更新承兑商
	OPActionAddExchangerSuspend             = 90  //禁用承兑商
	OPActionAddExchangerRestore             = 91  //恢复承兑商
	OPActionReadOperation                   = 92  //查询操作日志
	OPActionReadOrder                       = 93  //查询订单
	OPActionAddOrderCancel                  = 94  //取消订单
	OPActionAddOrderConfirm                 = 95  //确认订单
	OPActionReadOss                         = 96  //oss直传认证
	OPActionReadExchangersVerify            = 97  //查询承兑商审核
	OPActionReadExchangersVerifyApprove     = 98  //通过承兑商审核
	OPActionReadExchangersVerifyReject      = 99  //拒绝承兑商审核
	OPActionDelUserBulk                     = 100 //匹配删除otc用户
	OPActionReadUser                        = 101 //查询otc用户
	OPActionEditUser                        = 102 //编辑otc用户
	OPActionDelUser                         = 103 //删除otc用户
	OPActionAddUserSyncUsdtbyuid            = 104 //同步otc用户usdt链上数据
	OPActionAddUserEusdRecharge             = 105 //otc用户eusd充值
	OPActionAddUserSuspend                  = 106 //禁用otc用户
	OPActionAddUserRestore                  = 107 //恢复otc用户
	OPActionReadUserPayment                 = 108 //查询otc用户支付方式
	OPActionReadUserLoginLog                = 109 //查询otc登录日志
	OPActionReadPrice                       = 110 //查询指导价
	OPActionReadSmsTemplate                 = 111 //查询短信模板
	OPActionEditSmsTemplate                 = 112 //编辑短信模板
	OPActionAddSmsTemplate                  = 113 //添加短信模板
	OPActionDelSmsTemplate                  = 114 //删除短信模板
	OPActionReadSmsCode                     = 115 //查询短信验证码
	OPActionReadNotification                = 116 //查询系统通知
	OPActionEditNotification                = 117 //编辑系统通知
	OPActionAddNotification                 = 118 //添加系统通知
	OPActionDelNotification                 = 119 //删除系统通知
	OPActionAddNotificationPublish          = 120 //发布系统通知
	OPActionAddNotificationUnPublish        = 121 //取消发布系统通知
	OPActionReadSystemMessage               = 122 //查询系统消息
	OPActionAddSystemMessages               = 123 //添加系统消息
	OPActionEditSystemMessages              = 124 //编辑系统消息
	OPActionDelSystemMessages               = 125 //删除系统消息
	OPActionReadTopAgent                    = 126 //查询一级代理
	OPActionAddTopAgent                     = 127 //添加一级代理
	OPActionEditTopAgent                    = 128 //编辑一级代理
	OPActionDelTopAgent                     = 129 //删除一级代理
	OPActionReadUsdtAddresses               = 130 //查询usdt地址
	OPActionReadUsdtWallet                  = 131 //查询usdt钱包
	OPActionReadUsdtTransRecord             = 132 //查询usdt转入转出
	OPActionReadUsdtTurnRecord              = 133 //查询usdt抵押赎回
	OPActionReadUsdtCashRecord              = 134 //查询usdt提现申请
	OPActionAddUsdtCashRecordApprove        = 135 //通过usdt提现申请
	OPActionAddUsdtCashRecordReject         = 136 //拒绝usdt提现申请
	OPActionAddOrderResolve                 = 137 //解决申诉订单
	OPActionReadAppealServiceRecord         = 138 //查询申述客服处理记录
	OPActionAddOrderConfirmPay              = 139 //申述客服确认买家已付款
	OPActionAddOrderFrozenSell              = 140 //申述客服冻结卖家
	OPActionAddOrderUnFrozenSell            = 141 //申述客服解冻卖家
	OPActionAddOrderForbidSell              = 142 //申述客服禁止卖家
	OPActionAddOrderActiveSell              = 143 //申述客服激活卖家
	OPActionAddOrderForbidBuy               = 144 //申述客服禁止买家
	OPActionAddOrderActiveBuy               = 145 //申述客服激活买家
	OPActionAddEosFrozen                    = 146 //冻结用户EOS
	OPActionAddEosUnFrozen                  = 147 //解冻用户EOS
	OPActionReadCurPrice                    = 148 //查询当前指导价
	OPActionReadEndPoint                    = 149 //查询域名配置
	OPActionAddEndPoint                     = 150 //创建域名配置
	OPActionEditEndPoint                    = 151 //编辑域名配置
	OPActionDelEndPoint                     = 152 //删除域名配置
	OPActionAddAdminsUnBind                 = 153 //解除管理员google绑定
	OPActionAddPlatformUserCateAdd          = 154 //增加平台账号分类
	OPActionReadPlatformUserCate            = 155 //读取平台账号分类
	OPActionReadAnnouncement                = 156 //查询公告
	OPActionAddAnnouncement                 = 157 //创建公告
	OPActionEditAnnouncement                = 158 //编辑公告
	OPActionDelAnnouncement                 = 159 //删除公告
	OPActionGetOtcOrder                     = 160 //获取承兑商otc订单
	OPActionGetOtcWealth                    = 161 //获取承兑商otc资产
	OPActionReadPlatformUser                = 162 //增加平台账号
	OPActionAddPlatformUserAdd              = 163 //读取平台账号
	OPActionAddPlatformUserStatus           = 164 //编辑平台账号状态
	OPActionReadConfigWarning               = 165 //查询预警配置
	OPActionAddConfigWarning                = 166 //创建预警配置
	OPActionEditConfigWarning               = 167 //编辑预警配置
	OPActionDelConfigWarning                = 168 //删除预警配置
	OPActionAddUsdtLock                     = 169 //usdt账户锁定
	OPActionAddUsdtUnlock                   = 170 //usdt账户解锁
	OPActionReadAppWhite                    = 171 //查询第三方应用白名单
	OPActionAddAppWhite                     = 172 //创建第三方应用白名单
	OPActionDelAppWhite                     = 173 //删除第三方应用白名单
	OPActionServerStopGet                   = 174 //服务状态获取
	OPActionServerStopSet                   = 175 //服务状态设置
	OPActionServerStopLog                   = 176 //服务状态日志获取
	OPActionOtcStat                         = 177 //查看OTC统计数据
	OPActionGeneralFrozenUser               = 178 //冻结用户账户
	OPActionGeneralUnFrozenUser             = 179 //解冻用户账户
	OPActionGeneralLockUser                 = 180 //锁定用户
	OPActionGeneralUnlockUser               = 181 //解锁用户
	OPActionReadMonthDividendCfg            = 182 //查看月分红配置
	OPActionEditMonthDividendCfg            = 183 ////编辑月分红档位配置
	OPActionEditAgentDividendPosition       = 184 //编辑代理月分红白名单档位代理
	OPActionDelAgentDividendPosition        = 185 //删除代理月分红白名单档位代理
	OPActionServerNodeGet                   = 186 //服务节点状态获取
	OPActionTaskRead                        = 187 //定时任务读操作
	OPActionTaskWrite                       = 188 //定时任务写操作
	OPActionTaskDelete                      = 189 //定时任务删除操作
	OPActionReadMonthDividendPosition       = 190 //查看月分红档位配置
	OPActionReadAgentDividendPosition       = 191 //分页查看月分红指定档位下的代理
	OPActionDeleteMonthDividendCfg          = 192 ////删除月分红档位配置
	OPActionQueryGameCharge                 = 193 //查询游戏充值或提现费用
	OPActionGetUsersByNick                  = 194 //根据用户名模糊查询用户
	OPActionGenReportDaily                  = 195 //每天生成ag平台数据
	OPActionGetGameReport                   = 196 //读取游戏数据
	OPActionReadMonthDividendWhiteList      = 197 //获取月分红白名单
	OPActionAddMonthDividendWhiteList       = 198 //增加月分红白名单
	OPActionEditMonthDividendWhiteList      = 199 //编辑月分红白名单
	OPActionDelMonthDividendWhiteList       = 200 //删除月分红白名单
	OPActionReadMonthDividendWhiteListAgent = 201 //获取月分红白名单代理列表
	OPActionEditMonthDividendWhiteListAgent = 202 //添加代理到指定白名单
	OPActionDelMonthDividendWhiteListAgent  = 203 //从白名单里删除代理
	OPActionAddProfitThreshold              = 204 //添加分润阈值
	OPActionDelProfitThreshold              = 205 //删除分润阈值
	OPActionUpdateProfitThreshold           = 206 //更新分润阈值
	OPActionFindProfitThreshold             = 207 //查找分润阈值
	OPActionAddMenu                         = 208 //添加菜单
	OPActionDelMenu                         = 209 //删除菜单
	OPActionUpdateMenu                      = 210 //更新菜单
	OPActionGetAllMenu                      = 211 //获取所有菜单
	OPActionGetAccessMenus                  = 212 //获取管理员可见菜单
	OPActionUpdateMenuAccess                = 213 //更新管理员的菜单权限
	OPActionReadReportTeamList              = 214 //读取团队报表列表
	OPActionReadReportTeamSum               = 215 //读取团队报表总计
	OPActionReadReportTeamPersonal          = 216 //读取团队报表个人
)

const (
	ERROR_CODERISK_ALERT = 201
)

//easyjson:json
type PageInfo struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func isString(str string) bool {
	_, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return true
	}
	return false
}

//获取url所需权限
func GetUrlPermission(url string, urlMothod string) (permissions []string) {
	if url == "" || urlMothod == "" {
		return
	}

	var permission string
	m := strings.ToUpper(urlMothod)
	switch m {
	case "GET":
		permission = "READ"
	case "PUT":
		permission = "EDIT"
	case "POST":
		permission = "ADD"
	case "DELETE":
		permission = "DELETE"
	default:
		return
	}

	var parent string
	start := -1
	s := strings.Split(url, "/")
	for i, v := range s {
		if start > 0 {
			if isString(v) {
				permission = permission + "_" + strings.ToUpper(v)
				parent = parent + "_" + strings.ToUpper(v)
				parentPermission := "MANAGE" + parent
				permissions = append(permissions, parentPermission)
			}
		}
		if v == "admin" {
			start = i
		}
	}
	permissions = append(permissions, permission)
	permissions = append(permissions, "ALL")

	return
}

func VerifyEmailFormat(email string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9@]{4,16}$", email); !ok {
		return false
	}
	return true
	//pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	//reg := regexp.MustCompile(pattern)
	//return reg.MatchString(email)
}

type UsdtAccountClient struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	Status    uint8  `json:"status"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Mortgaged string `json:"mortgaged"`
	Address   string `json:"address"`
	Mobile    string `json:"mobile"`
	Ctime     int64  `json:"ctime"`
	Utime     int64  `json:"utime"`
}

func ClientUsdtAccount(usdtAccount *usdtmodels.UsdtAccount) (clientData UsdtAccountClient) {
	if usdtAccount == nil {
		return
	}
	clientData.Id = fmt.Sprintf("%v", usdtAccount.Uaid)
	clientData.UserId = fmt.Sprintf("%v", usdtAccount.Uid)
	clientData.Status = usdtAccount.Status
	clientData.Ctime = usdtAccount.Ctime
	clientData.Utime = usdtAccount.Mtime
	clientData.Address = usdtAccount.Address
	clientData.Available = common.CurrencyInt64ToStr(usdtAccount.AvailableInteger, usdt.UsdtConfig.Precision)
	clientData.Frozen = common.CurrencyInt64ToStr(usdtAccount.FrozenInteger, usdt.UsdtConfig.Precision)
	clientData.Mortgaged = common.CurrencyInt64ToStr(usdtAccount.MortgagedInteger, usdt.UsdtConfig.Precision)
	clientData.Balance = common.CurrencyInt64ToStr(usdtAccount.AvailableInteger+usdtAccount.FrozenInteger+usdtAccount.MortgagedInteger, usdt.UsdtConfig.Precision)
	return
}
func ClientUsdtAccounts(usdtAccounts []usdtDao.UsdtAccount) (clentDatas []UsdtAccountClient) {
	for _, usdtAccount := range usdtAccounts {
		var clentData UsdtAccountClient
		clentData.Id = usdtAccount.Uaid
		clentData.UserId = usdtAccount.Uid
		clentData.Status = usdtAccount.Status
		clentData.Ctime = usdtAccount.Ctime
		clentData.Utime = usdtAccount.Mtime
		clentData.Address = usdtAccount.Address
		clentData.Available = common.CurrencyInt64ToStr(usdtAccount.AvailableInteger, usdt.UsdtConfig.Precision)
		clentData.Frozen = common.CurrencyInt64ToStr(usdtAccount.FrozenInteger, usdt.UsdtConfig.Precision)
		clentData.Mortgaged = common.CurrencyInt64ToStr(usdtAccount.MortgagedInteger, usdt.UsdtConfig.Precision)
		clentData.Balance = common.CurrencyInt64ToStr(usdtAccount.AvailableInteger+usdtAccount.FrozenInteger+usdtAccount.FrozenInteger, usdt.UsdtConfig.Precision)
		clentData.Mobile = usdtAccount.Mobile
		clentDatas = append(clentDatas, clentData)
	}

	return
}

type AgentInfoClient struct {
	Uid              string `json:"uid"`
	UMobile          string `json:"u_mobile"`
	ParentUid        string `json:"puid"`
	PMobile          string `json:"p_mobile"`
	WhiteListId      uint32 `json:"whitelist_id"`
	DividendPosition uint32 `json:"dividend_position"`
	InviteCode       string `json:"invite_code"`
	InviteNum        uint32 `json:"invite_num"`
	Commission       string `json:"commission"`
	Balance          string `json:"balance"`
	Ctime            int64  `json:"ctime"`
	Mtime            int64  `json:"utime"`
}

func ClientAgentInfos(agentInfos []agentdao.AgentInfo) (clentDatas []AgentInfoClient) {
	for _, agentInfo := range agentInfos {
		var clentData AgentInfoClient
		clentData.Uid = agentInfo.Uid
		clentData.UMobile = agentInfo.UMobile
		clentData.ParentUid = agentInfo.ParentUid
		clentData.PMobile = agentInfo.PMobile
		clentData.WhiteListId = agentInfo.WhitelistId
		clentData.DividendPosition = agentInfo.DividendPosition
		clentData.InviteCode = agentInfo.InviteCode
		clentData.InviteNum = agentInfo.InviteNum
		clentData.Mtime = agentInfo.Mtime
		clentData.Commission = eosplus.QuantityToString(uint64(agentInfo.SumSalary))
		clentData.Balance = eosplus.QuantityToString(uint64(agentInfo.SumCanWithdraw))
		clentDatas = append(clentDatas, clentData)
	}
	return
}

type CommissionStatClient struct {
	Ctime      int64  `orm:"column(ctime);pk" json:"ctime,omitempty"`
	Tax        string `orm:"column(tax_integer)" json:"tax,omitempty"`
	Channel    string `orm:"column(channel_integer)" json:"channel,omitempty"`
	Commission string `orm:"column(commission_integer)" json:"commission,omitempty"`
	Profit     string `orm:"column(profit_integer)" json:"profit,omitempty"`
	Mtime      int64  `orm:"column(mtime)" json:"utime,omitempty"`
	Status     uint8  `orm:"column(status)" json:"status,omitempty"`
}

func ClientCommissionStat(items []otcModels.CommissionStat) (clients []CommissionStatClient) {
	for _, data := range items {
		var clientData CommissionStatClient
		clientData.Status = data.Status
		clientData.Ctime = data.Ctime
		clientData.Mtime = data.Mtime

		var err error
		clientData.Tax, err = common.EncodeCurrency(data.TaxInteger, data.TaxDecimals, usdt.UsdtConfig.Precision)
		if err != nil {
			common.LogFuncError("%v", err)
		}
		clientData.Channel, err = common.EncodeCurrency(data.ChannelInteger, data.ChannelDecimals, usdt.UsdtConfig.Precision)
		if err != nil {
			common.LogFuncError("%v", err)
		}
		clientData.Commission, err = common.EncodeCurrency(data.CommissionInteger, data.CommissionDecimals, usdt.UsdtConfig.Precision)
		if err != nil {
			common.LogFuncError("%v", err)
		}
		clientData.Profit, err = common.EncodeCurrency(data.ProfitInteger, data.ProfitDecimals, usdt.UsdtConfig.Precision)
		if err != nil {
			common.LogFuncError("%v", err)
		}

		clients = append(clients, clientData)
	}

	return
}

//成为承兑商
func BecomeExchanger(uid uint64, from int8, mobile, wechat, telegram string) (data otcModels.OtcExchanger, err error) {
	data, err = otcDao.OtcExchangerDaoEntity.Create(uid, from, mobile, wechat, telegram)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	errCode := eosplus.EosPlusAPI.Wealth.BecomeExchanger(uid)
	if errCode != otcerror.ERROR_CODE_SUCCESS {
		err = errors.New("db error")
		common.LogFuncError("errCode:%v", errCode)
		return
	}
	_, err = otcDao.UserDaoEntity.BecomeExchange(uid)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func PostOtc(subRoute string, mapParams map[string]string, req interface{}, ack interface{}, timestamp uint32) controllers.ERROR_CODE {
	otcHost, err := common2.AppConfigMgr.String(common2.OtcHost)
	if err != nil {
		common.LogFuncError("get otc_host error:%v", err)
		return controllers.ERROR_CODE_OTC_REQ_FAIL
	}
	otcPort, err := common2.AppConfigMgr.String(common2.OtcPort)
	if err != nil {
		common.LogFuncError("get otcPort error:%v", err)
		return controllers.ERROR_CODE_OTC_REQ_FAIL
	}

	sign, err := common.GenerateDoubleMD5ByParams(mapParams, otcDao.SIGNATURE_SALT, timestamp)
	if err != nil {
		return controllers.ERROR_CODE_SIGN_FAIL
	}

	url := fmt.Sprintf("%s:%s/v1/api/%s?%s=%s", otcHost, otcPort, subRoute, KeySign, sign)
	b := httplib.Post(url)
	_, err = b.JSONBody(req)
	if err != nil {
		common.LogFuncError("JSONBody error:%v", err)
		return controllers.ERROR_CODE_OTC_REQ_FAIL
	}

	err = b.ToJSON(ack)
	if err != nil {
		common.LogFuncError("ToJSON error:%v", err)
		return controllers.ERROR_CODE_OTC_REQ_FAIL
	}

	return controllers.ERROR_CODE_SUCCESS
}

func GetOtcUidByMobile(mobile string) (uint64, controllers.ERROR_CODE) {
	if len(mobile) > 0 {
		otcUid, err := otcDao.UserDaoEntity.GetUidByMobile(mobile)
		if err != nil {
			if err == orm.ErrNoRows {
				return 0, controllers.ERROR_CODE_NO_USER
			}
			return 0, controllers.ERROR_CODE_DB
		}
		return otcUid, controllers.ERROR_CODE_SUCCESS
	}

	return 0, controllers.ERROR_CODE_SUCCESS
}

// cron function container
var FuncContainer interface{}

func Init() (err error) {
	//load tasks
	var tasks []models.Task
	tasks, err = dao.TaskDaoEntity.QueryAll(adminSvr)
	if err != nil {
		return
	}

	if len(tasks) > 0 {
		var taskMsgList systoolbox.TaskMsgList
		for _, t := range tasks {
			taskMsgList.List = append(taskMsgList.List, systoolbox.TaskMsg{
				Id:       t.Id,
				Name:     t.Name,
				Spec:     t.Spec,
				FuncName: t.FuncName,
				Switch:   dao.TaskStatusEnable,
			})
		}

		if taskList, ok := systoolbox.CheckTaskFunc(taskMsgList, true, FuncContainer, taskResultStoreLocal); ok {
			if len(taskList) > 0 {
				systoolbox.TaskMgr.AddTaskList(taskList)
			}
		}
	}

	return
}
