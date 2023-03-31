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
	var t T
	AddMatcher(FunMatcher(fmt.Sprintf("Captor[%s]", reflect.TypeOf(t)), func(m *matchers.MethodCall, a any) bool {
		at, ok := a.(T)
		if !ok {
			c.ctx.reporter.Fatalf("incorrect usage of argument captor. expected to capture type %s, got %v", reflect.TypeOf(t).String(), a)
			return false
		}
		cv := &capturedValue[T]{
			value: at,
			call:  m,
		}
		c.values = append(c.values, cv)
		return true
	}))
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
		serviceCall := false
		for _, cid := range c.ctx.serviceMethodCalls {
			if cid == c.values[i].call.ID {
				serviceCall = true
				break
			}
		}
		if !serviceCall {
			result = append(result, c.values[i].value)
		}
	}
	return result
}
