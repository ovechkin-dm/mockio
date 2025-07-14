package utils

import "sync"

type SyncList[T any] struct {
	lock  sync.Mutex
	items []T
}

func (l *SyncList[T]) Add(item T) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.items = append(l.items, item)
}

func (l *SyncList[T]) GetCopy() []T {
	l.lock.Lock()
	defer l.lock.Unlock()
	itemsCopy := make([]T, len(l.items))
	copy(itemsCopy, l.items)
	return itemsCopy
}

func NewSyncList[T any]() *SyncList[T] {
	return &SyncList[T]{
		lock:  sync.Mutex{},
		items: make([]T, 0),
	}
}
