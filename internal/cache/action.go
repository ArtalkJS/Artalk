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
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		panic("The 'dest' param in 'QueryDBWithCache' func is expected to pointer type to update its data.")
	}

	// use SingleFlight to prevent Cache Breakdown
	v, err, _ := cacheFindStoreGroup.Do(name, func() (any, error) {
		// query cache
		err := c.FindCache(name, dest)

		if err != nil { // cache miss

			// call queryDB() the dest value will be updated
			queryDB()

			if err := c.StoreCache(dest, name); err != nil {
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
	// if cache miss, `dest` has been updated in `queryDB()` so no need to set again
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
		log.Debug("[CacheMis] " + name)
		return err
	}

	log.Debug("[CacheHit] " + name)
	return nil
}

func (c *Cache) StoreCache(source any, name ...string) (err error) {
	for _, n := range name {
		log.Debug("[StoreCache] " + n)

		// `Set()` is Thread Safe, no need to add Mutex
		if setErr := c.marshal.Set(c.ctx, n, source,
			store.WithExpiration(c.ttl),
		); setErr != nil {
			err = setErr
		}
	}
	return
}

func (c *Cache) DelCache(name ...string) (err error) {
	for _, n := range name {
		log.Debug("[DelCache] " + n)
		if delErr := c.marshal.Delete(c.ctx, n); delErr != nil {
			err = delErr
		}
	}
	return
}
