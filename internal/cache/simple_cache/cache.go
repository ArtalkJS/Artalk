package simple_cache

import (
	"sync"
	"time"
)

type Cache struct {
	data sync.Map
}

func New() *Cache {
	return &Cache{}
}

type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

func (c *Cache) Set(key string, value interface{}, expirationParam ...time.Duration) {
	var expirationTime time.Time
	if len(expirationParam) > 0 && expirationParam[0] > 0 {
		expirationTime = time.Now().Add(expirationParam[0])
	}
	entry := cacheEntry{
		value:      value,
		expiration: expirationTime,
	}
	c.data.Store(key, entry)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	if val, found := c.data.Load(key); found {
		entry := val.(cacheEntry)
		if !isExpired(entry) {
			return entry.value, true
		}
		c.data.Delete(key)
	}
	return nil, false
}

func (c *Cache) Delete(key string) {
	c.data.Delete(key)
}

func (c* Cache) GC() {
	c.data.Range(func(key, value interface{}) bool {
		entry := value.(cacheEntry)
		if !isExpired(entry) {
			return true
		}
		c.data.Delete(key)
		return true
	})
}

func isExpired(entry cacheEntry) bool {
	return !entry.expiration.IsZero() && entry.expiration.Before(time.Now())
}
