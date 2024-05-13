package registry

import (
	"fmt"
	"reflect"

	"github.com/ovechkin-dm/go-dyno/proxy"

	"github.com/ovechkin-dm/mockio/matchers"
)

func AnyMatcher[T any]() matchers.Matcher[T] {
	return &matcherImpl[T]{
		f: func(values []any, a T) bool {
			return true
		},
		desc: fmt.Sprintf("Any[%s]", reflect.TypeOf(new(T)).Elem().String()),
	}
}

func FunMatcher[T any](description string, f func([]any, T) bool) matchers.Matcher[T] {
	return &matcherImpl[T]{
		f:    f,
		desc: description,
	}
}

type matcherImpl[T any] struct {
	f    func([]any, T) bool
	desc string
}

func (m *matcherImpl[T]) Description() string {
	return m.desc
}

func (m *matcherImpl[T]) Match(allArgs []any, actual T) bool {
	return m.f(allArgs, actual)
}

func untypedMatcher[T any](src matchers.Matcher[T]) matchers.Matcher[any] {
	return &matcherImpl[any]{
		f: func(args []any, a any) bool {
			var casted T
			if a == nil {
				return src.Match(args, casted)
			}
			switch v := a.(type) {
			case *proxy.DynamicStruct:
				casted = proxy.UnsafeCast[T](v)
			default:
				c, ok := a.(T)
				if !ok {
					return false
				}
				casted = c
			}
			return src.Match(args, casted)
		},
		desc: src.Description(),
	}
}
