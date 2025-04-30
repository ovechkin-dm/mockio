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

type HandlerHolder interface {
	Handler() matchers.Handler
}

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
	tp := reflect.TypeOf(new(T)).Elem()
	handler := ctrl.MockFactory.BuildHandler(ctrl.Env, tp)
	t, err := dyno.DynamicByType(handler.Handle, tp, dynoopts.WithPayload(handler))
	if err != nil {
		getInstance().reporter.FailNow(fmt.Errorf("error creating mock: %w", err))
		var zero T
		return zero
	}
	result, ok := t.(T)
	if !ok {
		getInstance().reporter.FailNow(fmt.Errorf("error casting mock to type %s", tp.String()))
		var zero T
		return zero
	}
	return result
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
	handler := UnwrapHandler(t)
	if handler == nil {
		return
	}
	handler.VerifyMethod(v)
}

func VerifyNoMoreInteractions(t any) {
	handler := UnwrapHandler(t)
	if handler == nil {
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

func NewMockController(reporter matchers.ErrorReporter, opts ...config.Option) *matchers.MockController {
	cfg := config.NewConfig()
	if reporter == nil {
		panic("MockController: provided a nil error reporting (*testing.T) instance")
	}
	for _, opt := range opts {
		opt(cfg)
	}
	env := &matchers.MockEnv{
		Reporter: reporter,
		Config:   cfg,
	}
	factory := &mockFactoryImpl{}
	ctrl := &matchers.MockController{
		Env:         env,
		MockFactory: factory,
	}
	return ctrl
}

func UnwrapHandler(mock any) *invocationHandler {
	if mock == nil {
		getInstance().reporter.ReportUnregisteredMockVerify(mock)
	}
	handlerHolder, ok := mock.(HandlerHolder)
	var handler *invocationHandler
	if ok {
		handler, handlerOk := handlerHolder.Handler().(*invocationHandler)
		if handlerOk {
			return handler
		}
	}
	payload, err := dyno.UnwrapPayload(mock)
	if err != nil {
		getInstance().reporter.ReportUnregisteredMockVerify(mock)
		return nil
	}
	handler, ok = payload.(*invocationHandler)
	if !ok {
		getInstance().reporter.ReportUnregisteredMockVerify(mock)
		return nil
	}
	return handler
}

type mockFactoryImpl struct{}

func (m *mockFactoryImpl) BuildHandler(env *matchers.MockEnv, ifaceType reflect.Type) matchers.Handler {
	handler := newHandler(ifaceType, getInstance().mockContext, env)
	env.Reporter.Cleanup(handler.TearDown)
	return handler
}
