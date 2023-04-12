package registry

import (
	"fmt"
	"github.com/ovechkin-dm/mockio/matchers"
	"reflect"
)

type capturedValue[T any] struct {
	value T
	call  *matchers.MethodCall
}

type captorImpl[T any] struct {
	values []*capturedValue[T]
	ctx    *mockContext
}

func (c *captorImpl[T]) Capture() T {
	tp := reflect.TypeOf(new(T)).Elem()
	AddMatcher(FunMatcher(fmt.Sprintf("Captor[%s]", tp), func(m *matchers.MethodCall, a any) bool {
		if c.ctx.IsServiceCall(m.ID) {
			return true
		}
		at, ok := a.(T)
		if !ok {
			c.ctx.reporter.Errorf("incorrect usage of argument captor. expected to capture type %s, got %v", tp.String(), a)
			return false
		}
		cv := &capturedValue[T]{
			value: at,
			call:  m,
		}
		c.values = append(c.values, cv)
		return true
	}))
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
	result := make([]T, 0)
	for i := range c.values {
		if !c.ctx.IsServiceCall(c.values[i].call.ID) {
			result = append(result, c.values[i].value)
		}
	}
	return result
}
