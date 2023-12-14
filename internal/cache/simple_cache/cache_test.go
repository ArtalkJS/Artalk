package simple_cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		cache := New()
		cache.Set("key", "value")
		value, ok := cache.Get("key")
		if !ok {
			t.Error("cache should be hit")
		}
		if value != "value" {
			t.Error("cache value should be correct")
		}
	})

	t.Run("Expired", func(t *testing.T) {
		cache := New()
		cache.Set("key", "value", 1*time.Millisecond)
		time.Sleep(2 * time.Millisecond)
		_, ok := cache.Get("key")
		if ok {
			t.Error("cache should be expired")
		}
	})

	t.Run("Concurrent", func(t *testing.T) {
		cache := New()

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				cache.Set(fmt.Sprint(i), i)
				value, ok := cache.Get(fmt.Sprint(i))
				if !ok {
					t.Errorf("cache should be hit for key %d", i)
				}
				if value != i {
					t.Errorf("cache value should be %d for key %d", i, i)
				}
			}(i)
		}
		wg.Wait()
	})

	t.Run("Delete", func(t *testing.T) {
		cache := New()
		cache.Set("key", "value")
		cache.Delete("key")
		_, ok := cache.Get("key")
		if ok {
			t.Error("cache should be deleted")
		}
	})

	t.Run("GC thold reached", func(t *testing.T) {
		cache := NewWithGCThold(1)
		cache.Set("key", "value", 1*time.Millisecond)
		time.Sleep(2 * time.Millisecond)
		cache.GC(false)
		_, ok := cache.Get("key")
		if ok {
			t.Error("cache should be expired")
		}
	})

	t.Run("GC not expired", func(t *testing.T) {
		cache := NewWithGCThold(1)
		cache.Set("key", "value")
		cache.GC(false)
		_, ok := cache.Get("key")
		if !ok {
			t.Error("cache should not be expired")
		}
	})

	t.Run("GC thold not reached", func(t *testing.T) {
		cache := NewWithGCThold(2)
		cache.Set("key", "value", 1*time.Millisecond)
		time.Sleep(2 * time.Millisecond)
		cache.GC(false)
		len := cache.len.Load()
		if len != 1 {
			t.Errorf("cache length should be 1, but got %d", len)
		}
	})

	t.Run("GC thold not reached but force GC", func(t *testing.T) {
		cache := NewWithGCThold(2)
		cache.Set("key", "value", 1*time.Millisecond)
		time.Sleep(2 * time.Millisecond)
		cache.GC(true)
		_, ok := cache.Get("key")
		if ok {
			t.Error("cache should be expired")
		}
	})
}
