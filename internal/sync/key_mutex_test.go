package sync_test

import (
	go_sync "sync"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/sync"
	"github.com/stretchr/testify/assert"
)

func TestKeyMutex_LockAndUnlock(t *testing.T) {
	keyMutex := sync.NewKeyMutex[int]()

	// Test that the same key is locked/unlocked correctly
	key := 1
	mutex := keyMutex.GetLock(key)
	mutex.Lock()

	unlocked := false

	go func() {
		mutex.Lock()
		unlocked = true
		mutex.Unlock()
	}()

	// Ensure the key is still locked
	time.Sleep(100 * time.Millisecond)
	assert.False(t, unlocked, "Key should be locked")

	// Unlock the key
	mutex.Unlock()

	// Ensure the key is unlocked now
	time.Sleep(100 * time.Millisecond)
	assert.True(t, unlocked, "Key should be unlocked now")
}

func TestKeyMutex_ConcurrentAccess(t *testing.T) {
	keyMutex := sync.NewKeyMutex[string]()

	var wg go_sync.WaitGroup
	key := "test-key"

	counter := 0
	// Use the same key in multiple goroutines to test concurrent access
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mutex := keyMutex.GetLock(key)
			mutex.Lock()
			defer mutex.Unlock()

			// Critical section
			counter++
		}()
	}

	wg.Wait()
	assert.Equal(t, 100, counter, "Counter should be incremented 100 times")
}

func TestKeyMutex_DifferentKeys(t *testing.T) {
	keyMutex := sync.NewKeyMutex[int]()

	var wg go_sync.WaitGroup
	counter1, counter2 := 0, 0
	key1, key2 := 1, 2

	// Lock different keys in different goroutines
	for i := 0; i < 50; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			mutex := keyMutex.GetLock(key1)
			mutex.Lock()
			defer mutex.Unlock()
			counter1++
		}()

		go func() {
			defer wg.Done()
			mutex := keyMutex.GetLock(key2)
			mutex.Lock()
			defer mutex.Unlock()
			counter2++
		}()
	}

	wg.Wait()

	assert.Equal(t, 50, counter1, "Counter1 should be incremented 50 times")
	assert.Equal(t, 50, counter2, "Counter2 should be incremented 50 times")
}
