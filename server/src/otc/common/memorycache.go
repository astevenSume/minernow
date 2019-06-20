package common

import "common"

type CACHE_INDEX int

const (
	CacheIndexCaptcha CACHE_INDEX = iota
	CacheIndexAccessToken
	CacheIndexMax
)

var CACHE_NAME = map[CACHE_INDEX]string{
	CacheIndexCaptcha:     "CACHE_CAPTCHA",
	CacheIndexAccessToken: "CACHE_ACCESS_TOKEN",
}

func MemoryCacheInit() error {
	common.MemoryCacheInit(int(CacheIndexMax))
	return nil
}
