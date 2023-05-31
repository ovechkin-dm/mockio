package registry

import (
	"fmt"
	"github.com/ovechkin-dm/go-dyno/pkg/dyno"
	"github.com/ovechkin-dm/mockio/matchers"
	"github.com/timandy/routine"
	"log"
	"reflect"
	"sync"
)

var instance = routine.NewThreadLocalWithInitial(newRegistry)
var lock sync.Mutex

type Registry struct {
	mockContext *mockContext
	mapping     map[any]*invocationHandler
}

func getInstance() *Registry {
	v := instance.Get()
	if v == nil {
		v = newRegistry()
		instance.Set(v)
	}
	return v.(*Registry)
}

func SetUp(reporter matchers.ErrorReporter) {
	if reporter == nil {
		log.Println("Warn: call to SetUp with nil reporter")
	}
	getInstance().mockContext = newMockContext(newEnrichedReporter(reporter))
	reporter.Cleanup(TearDown)
}

func TearDown() {
	reg := getInstance()
	instance.Remove()

	if reg.mockContext.reporter == nil {
		reg.mockContext.reporter.Errorf("Cannot TearDown since SetUp function wasn't called")
	}

	for _, v := range reg.mapping {
		v.CheckUnusedStubs()
	}
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

func AddMatcher[T any](m matchers.Matcher[T]) {
	withCheck[any](func() any {
		w := &matcherWrapper{
			matcher: untypedMatcher(m),
			rec:     nil,
		}
		getInstance().mockContext.getState().matchers = append(getInstance().mockContext.getState().matchers, w)
		return nil
	})
}

func AddCaptor[T any](c *captorImpl[T]) {
	withCheck[any](func() any {
		tp := reflect.TypeOf(new(T)).Elem()
		w := &matcherWrapper{
			matcher: FunMatcher(fmt.Sprintf("Captor[%s]", tp), func(call []any, a any) bool {
				return true
			}),
			rec: c,
		}
		getInstance().mockContext.getState().matchers = append(getInstance().mockContext.getState().matchers, w)
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
	reporter := &EnrichedReporter{&panicReporter{}}
	return &Registry{
		mockContext: newMockContext(reporter),
		mapping:     make(map[any]*invocationHandler, 0),
	}
}

func withCheck[T any](f func() T) T {
	lock.Lock()
	defer lock.Unlock()
	rep, ok := getInstance().mockContext.reporter.reporter.(*panicReporter)
	if ok {
		log.Println("Warning: reporter is not initialized. You can initialize it with `SetUp(*testing.T)`. Defaulting to the panic reporter. This could also happen when using mocks concurrently")
	}
	initRoutineID := getInstance().mockContext.routineID
	if initRoutineID != routine.Goid() {
		rep.Fatalf("Call to mock api from a different goroutine. `When` or `Verify` can only be used from the initial goroutine.")
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
