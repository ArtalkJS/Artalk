package cache

import (
	"reflect"

	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/eko/gocache/lib/v4/store"
	"golang.org/x/sync/singleflight"
)

var (
	cacheFindStoreGroup = new(singleflight.Group)
)

func (c *Cache) QueryDBWithCache(name string, dest any, queryDB func()) error {
	// use SingleFlight to prevent Cache Breakdown
	v, err, _ := cacheFindStoreGroup.Do(name, func() (any, error) {
		// query cache
		err := c.FindCache(name, dest)

		if err != nil { // cache miss

			// call queryDB() the dest value will be updated
			queryDB()

			if err := c.StoreCache(name, dest); err != nil {
				return nil, err
			}

			// because queryDB() had update dest value,
			// no need to update it again, so return nil
			return nil, nil
		}

		return dest, nil
	})

	if err != nil {
		return err
	}

	// update dest value only if cache hit,
	// if cache miss, dest has been updated in queryDB()
	if v != nil {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(v).Elem()) // similar to `*dest = &v`
	}

	return nil
}

func (c *Cache) FindCache(name string, dest any) error {
	// `Get()` is Thread Safe, so no need to add Mutex
	// @see https://github.com/go-redis/redis/issues/23
	_, err := c.marshal.Get(c.ctx, name, dest)
	if err != nil {
		return err
	}

	log.Debug("[Cache Hit] " + name)

	return nil
}

func (c *Cache) StoreCache(name string, source any) error {
	// `Set()` is Thread Safe too, no need to add Mutex either
	err := c.marshal.Set(c.ctx, name, source,
		store.WithExpiration(c.ttl),
	)
	if err != nil {
		return err
	}

	log.Debug("[写入缓存] " + name)

	return nil
}

func (c *Cache) DelCache(name string) error {
	return c.marshal.Delete(c.ctx, name)
}
