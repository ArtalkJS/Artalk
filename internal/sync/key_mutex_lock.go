package sync

import "sync"

type KeyMutex[T any] struct {
	locks map[any]*sync.Mutex

	mapLock sync.Mutex // to make the map safe concurrently
}

func NewKeyMutex[T any]() *KeyMutex[T] {
	return &KeyMutex[T]{locks: make(map[any]*sync.Mutex)}
}

func (l *KeyMutex[T]) getLockBy(key T) *sync.Mutex {
	l.mapLock.Lock()
	defer l.mapLock.Unlock()

	ret, found := l.locks[key]
	if found {
		return ret
	}

	ret = &sync.Mutex{}
	l.locks[key] = ret
	return ret
}

func (l *KeyMutex[T]) Lock(key T) {
	l.getLockBy(key).Lock()
}

func (l *KeyMutex[T]) Unlock(key T) {
	l.getLockBy(key).Unlock()
}
