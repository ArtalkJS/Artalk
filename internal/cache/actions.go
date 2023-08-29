package cache

import (
	"reflect"

	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/eko/gocache/lib/v4/store"
	"golang.org/x/sync/singleflight"
)

var (
	CacheFindGroup = new(singleflight.Group)
)

func FindAndStoreCache(name string, dest interface{}, queryDBResult func() interface{}) error {
	// SingleFlight 防止缓存击穿 (Cache breakdown)
	v, err, _ := CacheFindGroup.Do(name, func() (interface{}, error) {
		err := FindCache(name, dest)

		// cache hit 直接返回结果
		if err == nil {
			return dest, nil
		}

		// cache miss 查数据库
		result := queryDBResult()
		if err := StoreCache(name, result); err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		return err
	}

	if v != nil {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(v).Elem()) // similar to `*dest = &v`
	}

	return nil
}

func FindCache(name string, dest interface{}) error {
	if std == nil {
		// return fmt.Errorf("cache not initialize")
		return nil
	}

	// `Get()` is Thread Safe, so no need to add Mutex
	// @see https://github.com/go-redis/redis/issues/23
	_, err := std.Get(ctx, name, dest)
	if err != nil {
		return err
	}

	log.Debug("[Cache Hit] " + name)

	return nil
}

func StoreCache(name string, source interface{}) error {
	if std == nil {
		// return fmt.Errorf("cache not initialize")
		return nil
	}

	// `Set()` is Thread Safe too, no need to add Mutex either
	err := std.Set(ctx, name, source,
		store.WithExpiration(ttl),
	)
	if err != nil {
		return err
	}

	log.Debug("[写入缓存] " + name)

	return nil
}

func DelCache(name string) error {
	if std == nil {
		// return fmt.Errorf("cache not initialize")
		return nil
	}

	return std.Delete(ctx, name)
}
