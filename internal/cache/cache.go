package cache

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/allegro/bigcache/v3"
	"github.com/bradfitz/gomemcache/memcache"
	lib_cache "github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	bigcache_store "github.com/eko/gocache/store/bigcache/v4"
	memcache_store "github.com/eko/gocache/store/memcache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
)

var (
	std *marshaler.Marshaler
	ttl time.Duration
)

var ctx = context.Background()

func OpenCache(conf config.CacheConf) (err error) {
	ttl = time.Duration(conf.GetExpiresTime())

	var cacheStore store.StoreInterface

	switch conf.Type {

	case config.CacheTypeBuiltin:
		// 内建缓存
		bigcacheClient, err := bigcache.New(context.Background(), bigcache.DefaultConfig(
			// Tip: 内建缓存过期时间是一样的，只有 Redis/Memcache 才能设置单个 item 的
			ttl,
		))
		if err != nil {
			return err
		}
		cacheStore = bigcache_store.NewBigcache(bigcacheClient) // No options provided (as second argument)

	case config.CacheTypeRedis:
		// Redis
		network := "tcp"
		if conf.Redis.Network != "" {
			network = conf.Redis.Network
		}

		cacheStore = redis_store.NewRedis(redis.NewClient(&redis.Options{
			Network:  network,
			Addr:     conf.Server,
			Username: conf.Redis.Username,
			Password: conf.Redis.Password,
			DB:       conf.Redis.DB,
		}))

	case config.CacheTypeMemcache:
		// Memcache
		servers := strings.Split(conf.Server, ",")
		cacheStore = memcache_store.NewMemcache(
			memcache.New(servers...),
			store.WithExpiration(ttl),
		)

	default:
		log.Fatal(`Invalid cache type "` + conf.Type + `", please check config option "cache.type"`)

	}

	cacheInstance := lib_cache.New[any](cacheStore)

	// marshaler wrapper
	// marshaler using VmihailencoMsgpack
	// @link https://github.com/vmihailenco/msgpack
	// Benchmarks
	// @link https://github.com/alecthomas/go_serialization_benchmarks
	std = marshaler.New(cacheInstance)

	return
}

func CloseCache() {
	// std.Clear(Ctx)
	std = nil
}
