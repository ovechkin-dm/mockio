package registry

import (
	"reflect"
	"sync"
)

type recordable interface {
	Record(call *MethodCall, value any)
	RemoveRecord(call *MethodCall)
}

type capturedValue[T any] struct {
	value T
	call  *MethodCall
}

type captorImpl[T any] struct {
	values []*capturedValue[T]
	ctx    *mockContext
	lock   sync.Mutex
}

func (c *captorImpl[T]) Capture() T {
	AddCaptor[T](c)
	var t T
	return t
}

func (c *captorImpl[T]) Last() T {
	values := c.Values()
	if len(values) == 0 {
		c.ctx.reporter.ReportEmptyCaptor()
		var t T
		return t
	}
	return values[len(values)-1]
}

func (c *captorImpl[T]) Values() []T {
	c.lock.Lock()
	defer c.lock.Unlock()
	result := make([]T, len(c.values))
	for i := range c.values {
		result[i] = c.values[i].value
	}
	return result
}

func (c *captorImpl[T]) Record(call *MethodCall, value any) {
	c.lock.Lock()
	defer c.lock.Unlock()
	t, ok := value.(T)
	if !ok {
		tp := reflect.TypeOf(new(T)).Elem()
		c.ctx.reporter.ReportInvalidCaptorValue(tp, reflect.TypeOf(value))
		return
	}
	cv := &capturedValue[T]{
		value: t,
		call:  call,
	}
	c.values = append(c.values, cv)
}

func (c *captorImpl[T]) RemoveRecord(call *MethodCall) {
	c.lock.Lock()
	defer c.lock.Unlock()
	wo := make([]*capturedValue[T], 0)
	for _, v := range c.values {
		if v.call != call {
			wo = append(wo, v)
		}
	}
	c.values = wo
}
