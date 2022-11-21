package cache

import (
	"github.com/gogf/gcache-adapter/adapter"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	"time"
)

//缓存适配器
//如果配置了redis缓存切为开启状态，就使用redis缓存，否则使用本进程的内存缓存

func CurrUseRedisCache() bool {
	return g.Cfg().GetBool("redis.open") && g.Cfg().GetString("redis.default") != ""
}

var cache = gcache.New()

func Get(key string) (string, error) {
	if CurrUseRedisCache() {
		cache.SetAdapter(adapter.NewRedis(g.Redis()))
		value, err := cache.Get(key)
		if value != nil {
			return value.(string), err
		}
		return "", err
	}

	value, err := cache.Get(key)
	if value != nil {
		return value.(string), err
	}
	return "", err
}

func Set(key, value string, expire time.Duration) error {
	if CurrUseRedisCache() {
		cache.SetAdapter(adapter.NewRedis(g.Redis()))
	}

	return cache.Set(key, value, expire)
}

func Del(key string) error {
	if CurrUseRedisCache() {
		cache.SetAdapter(adapter.NewRedis(g.Redis()))
	}
	_, err := cache.Remove(key)
	return err
}
