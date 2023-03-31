package registry

import (
	"fmt"
	"github.com/ovechkin-dm/mockio/matchers"
	"reflect"
)

func AnyMatcher[T any]() matchers.Matcher {
	var t T
	return &matcherImpl{
		f: func(m *matchers.MethodCall, a any) bool {
			_, ok := a.(T)
			return ok
		},
		desc: fmt.Sprintf("Any[%s]", reflect.TypeOf(t).String()),
	}
}

func FunMatcher(description string, f func(*matchers.MethodCall, any) bool) matchers.Matcher {
	return &matcherImpl{
		f:    f,
		desc: description,
	}
}

type matcherImpl struct {
	f    func(*matchers.MethodCall, any) bool
	desc string
}

func (m *matcherImpl) Description() string {
	return m.desc
}

func (m *matcherImpl) Match(methodCall *matchers.MethodCall, actual interface{}) bool {
	return m.f(methodCall, actual)
}
