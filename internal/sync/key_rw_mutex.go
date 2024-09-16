package sync

import "sync"

// KeyRWMutex is a read-write lock with keys
type KeyRWMutex[T any] struct {
	locks sync.Map // use sync.Map for concurrent access
}

// NewKeyRWMutex creates a new KeyRWMutex
func NewKeyRWMutex[T any]() *KeyRWMutex[T] {
	return &KeyRWMutex[T]{}
}

// GetLock retrieves or creates the lock for the given key
func (kl *KeyRWMutex[T]) GetLock(key T) *sync.RWMutex {
	if lock, ok := kl.locks.Load(key); ok {
		return lock.(*sync.RWMutex)
	}
	// Create a new RWMutex if not found
	newLock := &sync.RWMutex{}
	// Ensure only one RWMutex is added for the key
	actualLock, _ := kl.locks.LoadOrStore(key, newLock)
	return actualLock.(*sync.RWMutex)
}
