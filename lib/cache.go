package lib

import (
	"context"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/allegro/bigcache/v3"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/marshaler"
	"github.com/eko/gocache/v2/store"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var (
	CACHE *marshaler.Marshaler
)

var Ctx = context.Background()

func OpenCache() (err error) {
	cacheType := config.Instance.Cache.Type

	var cacheStore store.StoreInterface

	switch cacheType {

	case config.CacheTypeBuiltin:
		// 内建缓存
		bigcacheClient, err := bigcache.NewBigCache(bigcache.DefaultConfig(
			// Tip: 内建缓存过期时间是一样的，只有 Redis/Memcache 才能设置单个 item 的
			time.Duration(config.Instance.Cache.GetExpiresTime()),
		))
		if err != nil {
			return err
		}
		cacheStore = store.NewBigcache(bigcacheClient, nil) // No options provided (as second argument)

	case config.CacheTypeRedis:
		// Redis
		network := "tcp"
		if config.Instance.Cache.Redis.Network != "" {
			network = config.Instance.Cache.Redis.Network
		}

		cacheStore = store.NewRedis(redis.NewClient(&redis.Options{
			Network:  network,
			Addr:     config.Instance.Cache.Server,
			Username: config.Instance.Cache.Redis.Username,
			Password: config.Instance.Cache.Redis.Password,
			DB:       config.Instance.Cache.Redis.DB,
		}), nil)

	case config.CacheTypeMemcache:
		// Memcache
		servers := strings.Split(config.Instance.Cache.Server, ",")
		cacheStore = store.NewMemcache(
			memcache.New(servers...),
			&store.Options{
				Expiration: time.Duration(config.Instance.Cache.GetExpiresTime()),
			},
		)

	default:
		logrus.Fatal("请检查配置文件 `cache.type` 无效缓存类型：", cacheType)

	}

	cacheInstance := cache.New(cacheStore)

	// marshaler wrapper
	// marshaler using VmihailencoMsgpack
	// @link https://github.com/vmihailenco/msgpack
	// Benchmarks
	// @link https://github.com/alecthomas/go_serialization_benchmarks
	CACHE = marshaler.New(cacheInstance)

	return
}
