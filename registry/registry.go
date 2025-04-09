package registry

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/ovechkin-dm/go-dyno/pkg/dyno"
	"github.com/ovechkin-dm/go-dyno/pkg/dynoopts"
	"github.com/ovechkin-dm/mockio/v2/config"
	"github.com/ovechkin-dm/mockio/v2/matchers"
	"github.com/ovechkin-dm/mockio/v2/threadlocal"
)

var instance = threadlocal.NewThreadLocal(newRegistry)	

type Registry struct {
	mockContext *mockContext
	reporter    *EnrichedReporter
}

func getInstance() *Registry {
	v := instance.Get()
	if v == nil {
		v = newRegistry()
		v.mockContext = newMockContext()
		instance.Set(v)
	}	
	return v
}
func Mock[T any](ctrl *matchers.MockController) T {
	handler := newHandler[T](getInstance().mockContext, ctrl)
	ctrl.Reporter.Cleanup(handler.TearDown)
	t, err := dyno.Dynamic[T](handler.Handle, dynoopts.WithPayload(handler))
	if err != nil {
		getInstance().reporter.FailNow(fmt.Errorf("error creating mock: %w", err))
		var zero T
		return zero
	}
	return t
}

func AddMatcher[T any](m matchers.Matcher[T]) {
	w := &matcherWrapper{
		matcher:    untypedMatcher(m),
		rec:        nil,
		stackTrace: NewStackTrace(),
	}
	getInstance().mockContext.getState().matchers = append(getInstance().mockContext.getState().matchers, w)
}

func AddCaptor[T any](c *captorImpl[T]) {
	tp := reflect.TypeOf(new(T)).Elem()
	w := &matcherWrapper{
		matcher: FunMatcher(fmt.Sprintf("Captor[%s]", tp), func(call []any, a any) bool {
			return true
		}),
		rec:        c,
		stackTrace: NewStackTrace(),
	}
	getInstance().mockContext.getState().matchers = append(getInstance().mockContext.getState().matchers, w)
}

func When() matchers.ReturnerAll {
	wh := getInstance().mockContext.getState().whenHandler
	if wh == nil {
		getInstance().reporter.ReportIncorrectWhenUsage()
		return nil
	}
	return wh.When()
}

func VerifyMethod(t any, v matchers.MethodVerifier) {
	payload, err := dyno.UnwrapPayload(t)
	if err != nil {
		getInstance().reporter.ReportUnregisteredMockVerify(t)
		return
	}
	handler, ok := payload.(*invocationHandler)
	if !ok {
		getInstance().reporter.ReportUnregisteredMockVerify(t)
		return
	}
	handler.VerifyMethod(v)
}

func VerifyNoMoreInteractions(t any) {
	payload, err := dyno.UnwrapPayload(t)
	if err != nil {
		getInstance().reporter.ReportUnregisteredMockVerify(t)
		return
	}
	handler, ok := payload.(*invocationHandler)
	if !ok {
		getInstance().reporter.ReportUnregisteredMockVerify(t)
		return
	}
	handler.VerifyNoMoreInteractions(false)
}

func newRegistry() *Registry {
	cfg := &config.MockConfig{
		PrintStackTrace: false,
	}
	reporter := newEnrichedReporter(&panicReporter{}, cfg)
	return &Registry{
		mockContext: newMockContext(),
		reporter:    reporter,
	}
}

func NewArgumentCaptor[T any]() matchers.ArgumentCaptor[T] {
	return &captorImpl[T]{
		values:   make([]*capturedValue[T], 0),
		ctx:      getInstance().mockContext,
		lock:     sync.Mutex{},
		reporter: getInstance().reporter,
	}
}
