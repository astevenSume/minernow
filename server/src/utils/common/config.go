package common

import (
	"common"
	"encoding/json"
	"errors"
	"strconv"
	"utils/admin/dao"
	modoles2 "utils/admin/models"
)

const (
	// sys begin
	ContactWeChat   = "contact_wechat"
	ContactTelegram = "contact_telegram"
	ServiceWechat   = "service_wechat"
	InviteWechat    = "invite_wechat"
	BuyFeeRate      = "buy_fee_rate"  //购买交易手续费
	SellFeeRate     = "sell_fee_rate" // 出OtcTradePayExpire售手续费
	PlatformName    = "platform_name"
	OtcHost         = "otc_host"
	OtcPort         = "otc_port"
	AgentInviteUrl  = "agent_invite_url"      //agent invite url
	SmsAccessAppId  = "sys_sms_access_appid"  //短信配置
	SmsAccessSecret = "sys_sms_access_secret" //短信配置
	SmsSignName     = "sys_sms_sign_name"     //短信配置
	SmsChina        = "sys_sms_china"         //短信配置
	SmsOversea      = "sys_sms_oversea"       //短信配置
	// sys end

	OtcTradePayExpire     = "otc_trade_pay_expire"      //交易付款有效期 默认15M
	OtcTradeConfirmExpire = "otc_trade_confirm_expire"  //交易确认有效期 默认15M
	OtcTradeUpperLimitRmb = "otc_trade_upper_limit_rmb" //最大交易限制  默认5000000分
	OtcTradeLowerLimitRmb = "otc_trade_lower_limit_rmb" //最小交易限制  默认10000分

	EosConfigKeyTokenAccount     = "eos_token_account"        //代币发行账户
	EosConfigKeyResourcesAccount = "eos_resources_account"    //eos管理账户，有eos用于抵押的
	EosConfigKeyCpuEos           = "eos_create_cpu"           //抵押转账需要CPU的eos
	EosConfigKeyNetEos           = "eos_create_net"           //抵押转账需要NET的eos
	EosConfigKeyRamEos           = "eos_create_ram"           //购买新账号需要的eos设置
	EosNoUseAccountLimit         = "eos_no_use_account_limit" //EOS账号保持未用数量  默认10

	UsdtTransferWarning    = "usdt_transfer_warning"     //usdt转账超额预警值
	UsdtTransferLowerLimit = "usdt_transfer_lower_limit" //usdt转账最低限额
	UsdtTransferAuditLimit = "usdt_transfer_audit_limit" //usdt审核阈值
	ProfitThreshold        = "profit_threshold"          //分润审核阈值

	GameRgRate          = "game_rg_rate"
	GameRgRateToUser    = "game_rg_rate_to_user"
	GameRgRateToUpParam = "game_rg_rate_to_up_param"
	GameRgRateUpParam   = "game_rg_rate_up_param"
)

// 系统配置获取入口
// 请到后台配置相关的key对应的value, table otc_admin.config
// eg： key： “sell_fee_rate”，value：“0.002”
// eg： key： “sell_fee_rate”，value：“0.002”
// f64, err := AppConfigMgr.Float("sell_fee_rate")

type ConfigManagerInterface interface {
	String(key string) (string, error)
	Int64(key string) (int64, error)
	Int(key string) (int, error)
	Float(key string) (float64, error)
	Bool(key string) (bool, error)
	Json(key string, v interface{}) error
	FlushActionCache(int8, uint32) error
	FlushCache(key string) error
	CleanActionCache(int8, uint32) error
	CleanCache(key string) error
	FetchString(args ...string) (list map[string]string, err error)
}

type ConfigManager struct {
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}

var AppConfigMgr ConfigManagerInterface = NewConfigManager()

func (m *ConfigManager) String(key string) (string, error) {
	var result string
	var err error
	result, err = common.RedisManger.Get(dao.ConfigDaoEntity.GenerateConfigKey(key)).Result()
	// redis data is nil， try to get from db
	if err != nil || result == "" {
		var value modoles2.Config
		value, err = dao.ConfigDaoEntity.QueryConfigByKey(key)
		if err != nil {
			return "", err
		}
		if value.Value == "" {
			return "", errors.New("string is empty")
		}
		// update redis data
		if err := common.RedisManger.Set(dao.ConfigDaoEntity.GenerateConfigKey(key), value.Value, 0).Err(); err != nil {
			common.LogFuncError("redis update err: %s:%s,%v", key, value.Value, err)
		}
		return value.Value, nil

	}
	return result, err
}

func (m *ConfigManager) Int(key string) (int, error) {
	i, e := m.Int64(key)
	return int(i), e
}

func (m *ConfigManager) Int64(key string) (int64, error) {
	if value, err := m.String(key); err != nil {
		return 0, err
	} else {
		return strconv.ParseInt(value, 10, 64)
	}
}

func (m *ConfigManager) Float(key string) (float64, error) {
	if value, err := m.String(key); err != nil {
		return 0, err
	} else {
		return strconv.ParseFloat(value, 64)
	}
}

func (m *ConfigManager) Bool(key string) (bool, error) {
	if value, err := m.String(key); err != nil {
		return false, err
	} else {
		return strconv.ParseBool(value)
	}
}

func (m *ConfigManager) Json(key string, v interface{}) error {
	if value, err := m.String(key); err != nil {
		return err
	} else {
		return json.Unmarshal([]byte(value), &v)
	}
}

func (m *ConfigManager) FlushActionCache(action int8, id uint32) error {
	if list, err := dao.ConfigDaoEntity.DumpConfig(action, id); err != nil {
		return err
	} else {
		values := []string{}
		for _, v := range list {
			values = append(values, dao.ConfigDaoEntity.GenerateConfigKey(v.Key))
			values = append(values, v.Value)
		}
		if _, err := common.RedisManger.MSet(values).Result(); err != nil {
			common.LogFuncError("flush all config data to redis err:%v", err)
			return err
		}
	}
	return nil
}

func (m *ConfigManager) FlushCache(key string) error {
	if config, err := dao.ConfigDaoEntity.QueryConfigByKey(key); err != nil {
		return err
	} else {
		// update redis data
		if err := common.RedisManger.Set(dao.ConfigDaoEntity.GenerateConfigKey(key), config.Value, 0).Err(); err != nil {
			common.LogFuncError("redis update err: %s:%s,%v", key, config.Value, err)
		}
	}
	return nil
}

func (m *ConfigManager) CleanActionCache(action int8, id uint32) error {
	if list, err := dao.ConfigDaoEntity.DumpConfig(action, id); err != nil {
		return err
	} else {
		keys := []string{}
		for _, v := range list {
			keys = append(keys, dao.ConfigDaoEntity.GenerateConfigKey(v.Key))
		}
		if err := common.RedisManger.Del(keys...).Err(); err != nil {
			common.LogFuncError("clear all config data to redis err:%v", err)
			return err
		}
	}
	return nil
}

func (m *ConfigManager) CleanCache(key string) error {
	if err := common.RedisManger.Del(dao.ConfigDaoEntity.GenerateConfigKey(key)).Err(); err != nil {
		common.LogFuncError("clear all config data to redis err:%v", err)
		return nil
	}
	return nil
}

func (m *ConfigManager) FetchString(args ...string) (list map[string]string, err error) {
	list = map[string]string{}
	keysList := []string{}
	for _, v := range args {
		keysList = append(keysList, dao.ConfigDaoEntity.GenerateConfigKey(v))
	}

	res := common.RedisManger.MGet(keysList...)

	resList := res.Val()
	isNull := false

	if len(list) <= 0 {
		isNull = true
	}
	for k, v := range resList {
		if v == nil {
			isNull = true
			break
		}
		list[args[k]] = v.(string)
	}

	//有空数据
	if !isNull {
		return
	}
	configList := dao.ConfigDaoEntity.Fetch(args...)
	values := []string{}
	for _, v := range configList {
		values = append(values, dao.ConfigDaoEntity.GenerateConfigKey(v.Key))
		values = append(values, v.Value)
		list[v.Key] = v.Value
	}

	if len(values) > 0 {
		if _, err := common.RedisManger.MSet(values).Result(); err != nil {
			common.LogFuncError("flush all config data to redis err:%v", err)
		}
	}

	return
}
