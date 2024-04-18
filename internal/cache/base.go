package cache

import (
	"context"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/log"
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

type Cache struct {
	ttl      time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	instance *lib_cache.Cache[any]
	marshal  *marshaler.Marshaler
}

func (cache *Cache) Close() {
	cache.cancel()
	cache.marshal = nil
	cache.instance = nil
}

func New(conf config.CacheConf) (*Cache, error) {
	// create new context
	ctx, cancel := context.WithCancel(context.Background())

	cache := &Cache{
		ttl:    time.Duration(conf.GetExpiresTime()),
		ctx:    ctx,
		cancel: cancel,
	}

	var cacheStore store.StoreInterface

	switch conf.Type {

	case config.CacheTypeBuiltin:
		// 内建缓存
		bigcacheClient, err := bigcache.New(context.Background(), bigcache.DefaultConfig(
			// Tip: 内建缓存过期时间是一样的，只有 Redis/Memcache 才能设置单个 item 的
			cache.ttl,
		))
		if err != nil {
			return nil, err
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
			store.WithExpiration(cache.ttl),
		)

	default:
		log.Fatal(`Invalid cache type "` + conf.Type + `", please check config option "cache.type"`)

	}

	cache.instance = lib_cache.New[any](cacheStore)

	// marshaler wrapper
	// marshaler using VmihailencoMsgpack
	// @link https://github.com/vmihailenco/msgpack
	// Benchmarks
	// @link https://github.com/alecthomas/go_serialization_benchmarks
	cache.marshal = marshaler.New(cache.instance)

	return cache, nil
}
