package sync

import "sync"

// KeyMutex provides a per-key mutex mechanism
type KeyMutex[T any] struct {
	locks sync.Map // sync.Map for concurrent access
}

// NewKeyMutex creates a new KeyMutex
func NewKeyMutex[T any]() *KeyMutex[T] {
	return &KeyMutex[T]{}
}

// GetLock returns the mutex associated with the given key, creating it if necessary
func (l *KeyMutex[T]) GetLock(key T) *sync.Mutex {
	if lock, ok := l.locks.Load(key); ok {
		return lock.(*sync.Mutex)
	}
	// Create a new mutex if not found
	newLock := &sync.Mutex{}
	// Ensure only one mutex is added for the key
	actualLock, _ := l.locks.LoadOrStore(key, newLock)
	return actualLock.(*sync.Mutex)
}
