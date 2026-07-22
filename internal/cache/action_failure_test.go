package cache

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	lib_cache "github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryDBWithCacheReturnsDatabaseValueWhenStoreFails(t *testing.T) {
	storeErr := errors.New("cache store unavailable")
	cacheStore := &failureStore{
		getErr: errors.New("cache miss"),
		setErr: storeErr,
	}
	cacheInstance := newCacheWithStore(cacheStore)
	defer cacheInstance.Close()

	type user struct {
		Name  string
		Email string
	}
	want := user{Name: "admin", Email: "admin@example.com"}

	got, err := QueryDBWithCache(cacheInstance, t.Name(), func() (user, error) {
		return want, nil
	})

	require.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, int32(1), cacheStore.setCalls.Load())
}

func TestQueryDBWithCacheReturnsDatabaseError(t *testing.T) {
	dbErr := errors.New("database unavailable")
	cacheStore := &failureStore{getErr: errors.New("cache miss")}
	cacheInstance := newCacheWithStore(cacheStore)
	defer cacheInstance.Close()

	got, err := QueryDBWithCache(cacheInstance, t.Name(), func() (string, error) {
		return "", dbErr
	})

	require.ErrorIs(t, err, dbErr)
	assert.Empty(t, got)
	assert.Equal(t, int32(0), cacheStore.setCalls.Load())
}

func newCacheWithStore(cacheStore store.StoreInterface) *Cache {
	ctx, cancel := context.WithCancel(context.Background())
	instance := lib_cache.New[any](cacheStore)

	return &Cache{
		ttl:      time.Minute,
		ctx:      ctx,
		cancel:   cancel,
		instance: instance,
		marshal:  marshaler.New(instance),
	}
}

type failureStore struct {
	getErr   error
	setErr   error
	setCalls atomic.Int32
}

func (s *failureStore) Get(context.Context, any) (any, error) {
	return nil, s.getErr
}

func (s *failureStore) GetWithTTL(context.Context, any) (any, time.Duration, error) {
	return nil, 0, s.getErr
}

func (s *failureStore) Set(context.Context, any, any, ...store.Option) error {
	s.setCalls.Add(1)
	return s.setErr
}

func (*failureStore) Delete(context.Context, any) error {
	return nil
}

func (*failureStore) Invalidate(context.Context, ...store.InvalidateOption) error {
	return nil
}

func (*failureStore) Clear(context.Context) error {
	return nil
}

func (*failureStore) GetType() string {
	return "failure"
}
