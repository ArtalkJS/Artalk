package cache

import (
	"context"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/allegro/bigcache/v3"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	bigcache_store "github.com/eko/gocache/store/bigcache/v4"
	memcache_store "github.com/eko/gocache/store/memcache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
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
		cacheStore = bigcache_store.NewBigcache(bigcacheClient) // No options provided (as second argument)

	case config.CacheTypeRedis:
		// Redis
		network := "tcp"
		if config.Instance.Cache.Redis.Network != "" {
			network = config.Instance.Cache.Redis.Network
		}

		cacheStore = redis_store.NewRedis(redis.NewClient(&redis.Options{
			Network:  network,
			Addr:     config.Instance.Cache.Server,
			Username: config.Instance.Cache.Redis.Username,
			Password: config.Instance.Cache.Redis.Password,
			DB:       config.Instance.Cache.Redis.DB,
		}))

	case config.CacheTypeMemcache:
		// Memcache
		servers := strings.Split(config.Instance.Cache.Server, ",")
		cacheStore = memcache_store.NewMemcache(
			memcache.New(servers...),
			store.WithExpiration(time.Duration(config.Instance.Cache.GetExpiresTime())),
		)

	default:
		logrus.Fatal("Invalid cache type `" + cacheType + "`, please check config option `cache.type`")

	}

	cacheInstance := cache.New[any](cacheStore)

	// marshaler wrapper
	// marshaler using VmihailencoMsgpack
	// @link https://github.com/vmihailenco/msgpack
	// Benchmarks
	// @link https://github.com/alecthomas/go_serialization_benchmarks
	CACHE = marshaler.New(cacheInstance)

	return
}
