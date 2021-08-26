package lib

import (
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
)

var CACHE *cache.Cache

func OpenCache() (err error) {
	bigcacheClient, err := bigcache.NewBigCache(bigcache.DefaultConfig(30 * time.Minute)) // TODO 过期时间是一样的，只有 Redis 才能设置单个 key 的
	if err != nil {
		return err
	}

	bigcacheStore := store.NewBigcache(bigcacheClient, nil) // No options provided (as second argument)
	CACHE = cache.New(bigcacheStore)
	return
}
