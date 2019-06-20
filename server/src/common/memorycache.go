package common

import (
	"fmt"
	. "github.com/astaxie/beego/cache"
	"time"
)

//
var adapters []Cache

func checkCacheIndex(index int) (err error) {
	if len(adapters) < index+1 {
		err = fmt.Errorf("adaper of index %d no found.", index)
		return
	}
	return
}

func MemoryCacheInit(num int) (err error) {
	adapters = make([]Cache, num, num)
	for i := 0; i < num; i++ {
		adapters[i], err = NewCache("memory", `{"interval":20}`)
		if err != nil || adapters[i] == nil {
			LogFuncError("init cache %d failed", i)
			return
		}
	}

	return
}

func MemorycacheGet(index int, k string) (err error, v interface{}) {
	err = checkCacheIndex(index)
	if err == nil {
		v = adapters[index].Get(k)
	}

	return
}

func MemorycachePut(index int, k string, v interface{}, duration time.Duration) (err error) {
	err = checkCacheIndex(index)
	if err == nil {
		err = adapters[index].Put(k, v, duration)
	}

	return
}

func MemorycacheRemove(index int, k string) (err error) {
	err = checkCacheIndex(index)
	if err == nil {
		err = adapters[index].Delete(k)
	}

	return
}
