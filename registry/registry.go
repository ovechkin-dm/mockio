package registry

import (
	"fmt"
	"github.com/ovechkin-dm/go-dyno/pkg/dyno"
	"github.com/ovechkin-dm/mockio/matchers"
	"github.com/timandy/routine"
	"log"
	"sync"
)

var instance = routine.NewInheritableThreadLocalWithInitial(newRegistry)
var lock sync.Mutex

type Registry struct {
	mockContext *mockContext
	mapping     map[any]*invocationHandler
}

func getInstance() *Registry {
	return instance.Get().(*Registry)
}

func SetUp(reporter matchers.ErrorReporter) {
	if reporter == nil {
		panic("call to SetUp with nil reporter")
	}
	if getInstance().mockContext.reporter != nil {
		getInstance().mockContext.reporter.Errorf("Mock registry is already set up. SetUp method should be called once")
		return
	}
	getInstance().mockContext = newMockContext(newEnrichedReporter(reporter))
}

func TearDown() {
	if getInstance().mockContext.reporter == nil {
		getInstance().mockContext.reporter.Errorf("Cannot TearDown since SetUp function wasn't called")
	}
	instance.Set(newRegistry())
}

func Mock[T any]() T {
	return withCheck[T](func() T {
		handler := newHandler[T](getInstance().mockContext)
		t, err := dyno.Dynamic[T](handler)
		if err != nil {
			getInstance().mockContext.reporter.FailNow(fmt.Errorf("error creating mock: %w", err))
			var zero T
			return zero
		}
		getInstance().mapping[t] = handler
		return t
	})
}

func AddMatcher(m matchers.Matcher) {
	withCheck[any](func() any {
		getInstance().mockContext.getState().matchers = append(getInstance().mockContext.getState().matchers, m)
		return nil
	})
}

func When() matchers.ReturnerAll {
	return withCheck(func() matchers.ReturnerAll {
		wh := getInstance().mockContext.getState().whenHandler
		if wh == nil {
			getInstance().mockContext.reporter.ReportIncorrectWhenUsage()
			return nil
		}
		return wh.When()
	})
}

func VerifyMethod(t any, v matchers.MethodVerifier) {
	withCheck(func() any {
		handler, ok := getInstance().mapping[t]
		if !ok {
			getInstance().mockContext.reporter.ReportUnregisteredMockVerify(t)
			return nil
		}
		handler.VerifyMethod(v)
		return nil
	})
}

func VerifyInstance(t any, v matchers.InstanceVerifier) {
	withCheck(func() any {
		handler, ok := getInstance().mapping[t]
		if !ok {
			getInstance().mockContext.reporter.ReportUnregisteredMockVerify(t)
			return nil
		}
		handler.AddInstanceVerifier(v)
		return nil
	})
}

func newRegistry() any {
	return &Registry{
		mockContext: newMockContext(nil),
		mapping:     make(map[any]*invocationHandler, 0),
	}
}

func withCheck[T any](f func() T) T {
	lock.Lock()
	defer lock.Unlock()

	if getInstance() == nil || getInstance().mockContext.reporter == nil {
		log.Fatalf("reporter is not initialized. You can initialize it with `mock.SetUp(*testing.T)`")
	}
	return f()
}

func NewArgumentCaptor[T any]() matchers.ArgumentCaptor[T] {
	return withCheck(func() matchers.ArgumentCaptor[T] {
		return &captorImpl[T]{
			values: make([]*capturedValue[T], 0),
			ctx:    getInstance().mockContext,
		}
	})
}
