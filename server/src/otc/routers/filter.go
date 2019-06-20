package routers

import (
	"common"
	"encoding/json"
	"fmt"
	"otc/controllers"
	"regexp"
	"strconv"
	"strings"
	"time"

	common2 "otc/common"
	otcerror "otc_error"
	utils "utils/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/go-redis/redis"
)

const (
	RedisKeyExpire = 600 * time.Second //上次请求时间戳保留时间
)

type ignoreRule struct {
	All             bool
	Static, Dynamic []string
	DynamicReg      []*regexp.Regexp
}

var (
	IGNORE         = map[string]ignoreRule{}
	SIGNATURE_SALT string
)

// InitFilter 初始化过滤器
func InitFilter() (err error) {

	SIGNATURE_SALT = beego.AppConfig.String("anti_replay::salt")

	if ignore := beego.AppConfig.String("anti_replay::ignore"); ignore != "" {
		if err = json.Unmarshal([]byte(ignore), &IGNORE); err != nil {
			return
		}

		for _, rule := range IGNORE {
			if len(rule.Dynamic) > 0 {
				rule.DynamicReg = make([]*regexp.Regexp, 0, len(rule.Dynamic))

				for _, pattern := range rule.Dynamic {
					reg, err := regexp.Compile(pattern)
					if err != nil {
						return err
					}
					rule.DynamicReg = append(rule.DynamicReg, reg)
				}
			}
		}

	}
	return
}

// AntiReplayFilter 防重放过滤器
func AntiReplayFilter(ctx *context.Context) {

	if rule, has := IGNORE[strings.ToLower(ctx.Request.Method)]; has {
		if rule.All {
			return
		}
		var path = ctx.Request.URL.Path
		for _, url := range rule.Static {
			if url == path {
				return
			}
		}
		for _, reg := range rule.DynamicReg {
			if reg.MatchString(path) {
				return
			}
		}
	}

	var (
		hash, sign, raw string
		err             error
		ctl             controllers.BaseController
		timestamp       int64
		uid             uint64
		ok              bool
	)
	ctl = controllers.BaseController{Controller: beego.Controller{Ctx: ctx}}
	sign = ctx.Request.Header.Get("sign")

	// TODO: 上线前删除
	if sign == "" {
		return
	}
	// TODO: 上线前删除

	raw = fmt.Sprintf("%s%s%s", ctx.Request.URL.String(), ctx.Input.RequestBody, ctx.Request.Header.Get("timestamp"))
	// 校验参数 sign = md5(md5(url+post_string+timstamp)+salt)
	if hash, err = common.GenerateDoubleMD5(raw, SIGNATURE_SALT); err != nil || hash != sign {
		common.LogFuncDebug("check sign failed: client is %s , server is %s", sign, hash)
		ctl.ErrorResponse(otcerror.ERROR_CODE_PARAM_FAILED)
		return
	}

	// 请求的时间戳与服务器当前时间间隔超过 30s 则返回错误
	if timestamp, err = strconv.ParseInt(ctx.Request.Header.Get("timestamp"), 10, 64); err != nil || common.NowInt64MS()-timestamp > 30000 {
		common.LogFuncDebug("check timestamp failed: client is %v , server is %v", timestamp, common.NowInt64MS())
		ctl.ErrorResponse(otcerror.ERROR_CODE_PARAM_FAILED)
		return
	}

	// 获取用户 uid
	clientType := ctl.ClientType()
	if clientType == common2.ClientTypeUnkown {
		ctl.ErrorResponse(otcerror.ERROR_CODE_CLIENT_TYPE_UNKOWN)
		return
	}

	if ok, uid, _ = utils.CheckToken(clientType, ctx.GetCookie(controllers.TokenKey)); !ok {
		ctl.ErrorResponse(otcerror.ERROR_CODE_TOKEN_VERIFY_ERR)
		return
	}

	var (
		redisKey      string
		lastTimestamp int64
	)
	redisKey = fmt.Sprintf("anti_replay_%d", uid)
	// 上一次的时间戳大于本次请求的时间戳 则返回错误
	if lastTimestamp, err = common.RedisManger.Get(redisKey).Int64(); err != nil && err != redis.Nil {
		ctl.ErrorResponse(otcerror.ERROR_CODE_REDIS)
		return
	}

	if lastTimestamp != 0 && lastTimestamp >= timestamp {
		ctl.ErrorResponse(otcerror.ERROR_CODE_REQUEST_OUTDATED)
		return
	}
	if err = common.RedisManger.Set(redisKey, timestamp, RedisKeyExpire).Err(); err != nil {
		ctl.ErrorResponse(otcerror.ERROR_CODE_REDIS)
		return
	}
	return

}
