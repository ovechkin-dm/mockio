package registry

import (
	"fmt"
	"github.com/ovechkin-dm/mockio/matchers"
	"reflect"
)

func AnyMatcher[T any]() matchers.Matcher {
	return &matcherImpl{
		f: func(values []any, a any) bool {
			_, ok := a.(T)
			return ok
		},
		desc: fmt.Sprintf("Any[%s]", reflect.TypeOf(new(T)).Elem().String()),
	}
}

func FunMatcher(description string, f func([]any, any) bool) matchers.Matcher {
	return &matcherImpl{
		f:    f,
		desc: description,
	}
}

type matcherImpl struct {
	f    func([]any, any) bool
	desc string
}

func (m *matcherImpl) Description() string {
	return m.desc
}

func (m *matcherImpl) Match(allArgs []any, actual any) bool {
	return m.f(allArgs, actual)
}
