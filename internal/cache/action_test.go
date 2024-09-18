package cache_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/cache"
	"github.com/stretchr/testify/assert"
)

func TestAction(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	doCrudTest := func(t *testing.T, testKey string, testData any) {
		// Find Miss
		{
			var data any
			err := cache.FindCache(testKey, &data)
			assert.Error(t, err, "a find cache miss err should be occurred before store")
			assert.Zero(t, data, "cache should not be hit before store it")
		}

		// Store
		storeErr := cache.StoreCache(testData, testKey)

		if assert.NoError(t, storeErr, "error occurred while cache store") {
			// Find Hit
			var data any
			findErr := cache.FindCache(testKey, &data)

			if assert.NoError(t, findErr, "cache data should be hit") {
				assert.Equal(t, testData, data, "cache should equal to original data while cache hit")
			}
		}

		// Del
		delErr := cache.DelCache(testKey)
		if assert.NoError(t, delErr, "error occurred while cache delete") {
			// Find again
			var dataAgain any
			findErr := cache.FindCache(testKey, &dataAgain)
			assert.Error(t, findErr, "a find cache miss err should be occurred after the cache had been deleted")
			assert.Zero(t, dataAgain, "cache should not be found while cache deleted")
		}
	}

	type testCase struct {
		key   string
		value any
	}

	doTest := func(testCases []testCase) {
		for _, tt := range testCases {
			doCrudTest(t, tt.key, tt.value)
		}
	}

	doTest([]testCase{
		{key: "String", value: "12345"},
		{key: "Integer", value: int64(114514)},
		{key: "Float", value: float64(2.3333)},
	})

	t.Run("Concurrency", func(t *testing.T) {
		const numRoutines = 100

		var wg sync.WaitGroup
		wg.Add(numRoutines)

		for i := 0; i < numRoutines; i++ {
			go func(num int) {
				defer wg.Done()

				doTest([]testCase{
					{key: fmt.Sprintf("StringConcurrency_%d", num), value: "12345"},
					{key: fmt.Sprintf("IntegerConcurrency_%d", num), value: int64(114514)},
					{key: fmt.Sprintf("FloatConcurrency_%d", num), value: float64(2.3333)},
				})
			}(i)
		}

		wg.Wait()
	})
}

func TestQueryDBWithCache(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		cacheInstance := newTestCache(t)
		defer cacheInstance.Close()

		type user struct {
			Name  string
			Email string
		}

		const key = "data_key_233"
		var value = user{
			Name:  "qwqcode",
			Email: "artalkjs@example.com",
		}

		doCachedFind := func() bool {
			dbQueried := false

			data, err := cache.QueryDBWithCache(cacheInstance, key, func() (data user, err error) {
				// simulate db query result
				data = user{
					Name:  value.Name,
					Email: value.Email,
				}
				dbQueried = true
				return data, nil
			})

			if assert.NoError(t, err) {
				assert.Equal(t, value, data)
			}

			return dbQueried
		}

		if dbQueried := doCachedFind(); dbQueried {
			assert.True(t, dbQueried, "first call `QueryDBWithCache`, db should be queried")
		}

		if dbQueried := doCachedFind(); dbQueried {
			assert.False(t, dbQueried, "second call `QueryDBWithCache`, db should not be queried")
		}
	})

	t.Run("Concurrency", func(t *testing.T) {
		cacheInstance := newTestCache(t)
		defer cacheInstance.Close()

		type mockStruct struct {
			Value string
		}

		wg := sync.WaitGroup{}
		ready := make(chan struct{})
		findCallTimes := int32(0)

		const n = 1000
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-ready // make sure all goroutines start at the same time

				data, err := cache.QueryDBWithCache(cacheInstance, "key", func() (data mockStruct, err error) {
					atomic.AddInt32(&findCallTimes, 1)

					return mockStruct{
						Value: "concurrency_value",
					}, nil
				})

				if assert.NoError(t, err) {
					assert.Equal(t, "concurrency_value", data.Value, "data should be equal to the value returned by the query function")
				}
			}()
		}

		close(ready) // start all goroutines at the same time

		wg.Wait()

		if got := atomic.LoadInt32(&findCallTimes); got != 1 {
			t.Errorf("expected findCallTimes to be 1, got %d", got)
		}
	})
}
