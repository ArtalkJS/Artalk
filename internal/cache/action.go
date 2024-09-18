package cache

import (
	"fmt"
	"reflect"

	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/eko/gocache/lib/v4/store"
	"golang.org/x/sync/singleflight"
)

var cacheSingleflightGroup = new(singleflight.Group)

func QueryDBWithCache[T any](c *Cache, name string, queryDB func() (T, error)) (T, error) {
	// Use SingleFlight to prevent Cache Breakdown
	v, err, _ := cacheSingleflightGroup.Do(name, func() (any, error) {
		var val T

		// Query from cache
		err := c.FindCache(name, &val)

		if err != nil {
			// Miss cache

			// Query from db
			val, err := queryDB()
			if err != nil {
				return nil, err
			}

			// Store cache
			if err := c.StoreCache(val, name); err != nil {
				return nil, err
			}

			return val, nil
		} else {
			// Hit cache
			return val, nil
		}
	})
	if err != nil {
		return *new(T), err
	}

	return v.(T), err
}

func (c *Cache) FindCache(name string, dest any) error {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("[FindCache] dest must be a pointer")
	}

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
