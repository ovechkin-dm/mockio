package threadlocal

import (
	"sync"

	"github.com/petermattis/goid"
)

type ThreadLocal[T any] interface {
	Get() T
	Set(t T)
	Clear()
}

type impl[T any] struct {
	data sync.Map
	init func() T
}

func (i *impl[T]) Get() T {
	id := goid.Get()
	v, ok := i.data.Load(id)
	if !ok {
		nv := i.init()
		i.data.Store(id, nv)
		return nv
	}
	return v.(T)
}

func (i *impl[T]) Set(t T) {
	id := goid.Get()
	i.data.Store(id, t)
}

func (i *impl[T]) Clear() {
	id := goid.Get()
	i.data.Delete(id)
}

func NewThreadLocal[T any](initFunc func() T) ThreadLocal[T] {
	return &impl[T]{
		data: sync.Map{},
		init: initFunc,
	}
}

func GoId() int64 {
	return goid.Get()
}
