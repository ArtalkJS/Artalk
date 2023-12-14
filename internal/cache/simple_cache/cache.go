package simple_cache

import (
	"sync"
	"sync/atomic"
	"time"
)

type Cache struct {
	data sync.Map
	// when the number of entries reaches gcThold, GC will be triggered
	gcThold int
	len	 atomic.Int32
}

func New() *Cache {
	return &Cache{
		gcThold: 1000,
		len: atomic.Int32{},
	}
}

func NewWithGCThold(gcThold int) *Cache {
	return &Cache{
		gcThold: gcThold,
		len: atomic.Int32{},
	}
}

type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

func (c *Cache) Set(key string, value interface{}, expirationParam ...time.Duration) {
	defer c.GC(false)
	var expirationTime time.Time
	if len(expirationParam) > 0 && expirationParam[0] > 0 {
		expirationTime = time.Now().Add(expirationParam[0])
	}
	entry := cacheEntry{
		value:      value,
		expiration: expirationTime,
	}
	c.data.Store(key, entry)
	c.len.Add(1)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	defer c.GC(false)
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
	c.len.Add(-1)
}

func (c* Cache) GC(force bool) {
	if int(c.len.Load()) < c.gcThold && !force {
		return
	}
	c.data.Range(func(key, value interface{}) bool {
		entry := value.(cacheEntry)
		if !isExpired(entry) {
			return true
		}
		c.Delete(key.(string))
		return true
	})
}

func isExpired(entry cacheEntry) bool {
	return !entry.expiration.IsZero() && entry.expiration.Before(time.Now())
}
