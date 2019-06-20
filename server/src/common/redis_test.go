package common

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestRedisLock(t *testing.T) {
	RedisManger = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	TokenPrefix = "otc.0.0."
	RetryCount = 10
	RetryDelay = time.Millisecond * 10
	ExpiredTime = time.Second * 30

	key := "123456"

	l, err := RedisLock("123456")
	if err != nil {
		t.Fatalf("RedisLock %s failed : %v", key, err)
	}
	if l == nil {
		t.Fatalf("RedisLock %s failed", key)
	}
	v := RedisManger.Get(key).String()
	t.Logf("%s : %s", key, RedisManger.Get(key).String())
	time.Sleep(time.Second * 10)
	v1 := RedisManger.Get(key).String()
	if v != v1 {
		t.Fatalf("value changed : %v -> %v", v, v1)
	}
	t.Logf("after 10s, %s : %s", key, RedisManger.Get(key).String())
	l, err = RedisLock(key)
	if err == nil {
		t.Fatalf("should lock failed.")
	}
	if l != nil {
		t.Fatalf("should lock failed.")
	}
	time.Sleep(time.Second * 50)
	v2 := RedisManger.Get(key).String()
	if v2 != "" {
		t.Fatalf("should be nil")
	}
	l, err = RedisLock(key)
	if err != nil {
		t.Fatalf("RedisLock %s failed : %v", key, err)
	}
	if l == nil {
		t.Fatalf("RedisLock %s failed", key)
	}
}

func TestRedisUnlock(t *testing.T) {
	RedisManger = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	TokenPrefix = "otc.0.0."
	RetryCount = 10
	RetryDelay = time.Millisecond * 10
	ExpiredTime = time.Second * 60

	key := "1234567"

	l, err := RedisLock("1234567")
	if err != nil {
		t.Fatalf("RedisLock %s failed : %v", key, err)
	}
	if l == nil {
		t.Fatalf("RedisLock %s failed", key)
	}
	t.Logf("lock value -> %s", RedisManger.Get(key).Val())
	err = RedisUnlock(l)
	if err != nil {
		t.Fatalf("RedisUnlock failed : %v", err)
	}
	v1 := RedisManger.Get(key).Val()
	if v1 != "" {
		t.Fatalf("should be nil")
	}
	t.Logf("unlock value -> %s", v1)
	l, err = RedisLock(key)
	if err != nil {
		t.Fatalf("RedisLock %s failed : %v", key, err)
	}
	if l == nil {
		t.Fatalf("RedisLock %s failed", key)
	}
	v2 := RedisManger.Get(key).Val()
	t.Logf("lock again value -> %s", v2)
}

func TestJamLock(t *testing.T) {
	RedisManger = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	TokenPrefix = "otc.0.0."
	RetryCount = 300
	RetryDelay = time.Millisecond * 10
	ExpiredTime = time.Second * 60

	go func() {
		fmt.Println("into 1")
		l, err := RedisLock("test_lock")
		if err != nil {
			t.Fatalf("RedisLock failed:%s", err)
			return
		}
		fmt.Println("start sleep")
		time.Sleep(time.Second * 3)
		fmt.Println("up")
		RedisUnlock(l)
	}()
	time.Sleep(time.Second)
	go func() {
		fmt.Println("into 2")
		l, err := RedisLock("test_lock")
		if err != nil {
			t.Fatalf("RedisLock failed:%s", err)
			return
		}
		fmt.Println("start 2")

		fmt.Println("end 2")
		RedisUnlock(l)
	}()

	time.Sleep(time.Second * 5)
}
