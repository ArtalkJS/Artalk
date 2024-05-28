package sync

import (
	"sync"
)

// KeyRWMutex is a read-write lock with keys
type KeyRWMutex[T any] struct {
	mutex sync.Mutex
	locks map[any]*sync.RWMutex
}

// NewKeyRWMutex creates a new KeyRWMutex
func NewKeyRWMutex[T any]() *KeyRWMutex[T] {
	return &KeyRWMutex[T]{locks: make(map[any]*sync.RWMutex)}
}

// getLockByKey retrieves the lock for the given key
func (kl *KeyRWMutex[T]) getLockByKey(key string) *sync.RWMutex {
	kl.mutex.Lock()
	defer kl.mutex.Unlock()

	val, ok := kl.locks[key]
	if !ok {
		lock := &sync.RWMutex{}
		kl.locks[key] = lock
		return lock
	}
	return val
}

// Lock acquires the write lock for the given key
func (kl *KeyRWMutex[T]) Lock(key string) {
	kl.getLockByKey(key).Lock()
}

// Unlock releases the write lock for the given key
func (kl *KeyRWMutex[T]) Unlock(key string) {
	kl.getLockByKey(key).Unlock()
}

// RLock acquires the read lock for the given key
func (kl *KeyRWMutex[T]) RLock(key string) {
	kl.getLockByKey(key).RLock()
}

// RUnlock releases the read lock for the given key
func (kl *KeyRWMutex[T]) RUnlock(key string) {
	kl.getLockByKey(key).RUnlock()
}
