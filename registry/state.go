package registry

import (
	"github.com/ovechkin-dm/mockio/matchers"
	"github.com/timandy/routine"
	"sync"
)

type fiberState struct {
	matchers       []matchers.Matcher
	whenHandler    *invocationHandler
	verifyState    bool
	methodVerifier matchers.MethodVerifier
	lastMatch      *methodMatch
	whenCall       *matchers.MethodCall
}

type mockContext struct {
	state              routine.ThreadLocal
	serviceMethodCalls []string
	reporter           *EnrichedReporter
	lock               sync.Mutex
}

func (ctx *mockContext) addServiceMethodCallId(id string) {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()
	ctx.serviceMethodCalls = append(ctx.serviceMethodCalls, id)
}

func (ctx *mockContext) getServiceMethodCallIds() []string {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()
	result := make([]string, len(ctx.serviceMethodCalls))
	copy(result, ctx.serviceMethodCalls)
	return result
}

type methodMatch struct {
	matchers []matchers.Matcher
	answers  []matchers.Answer
	calls    []*matchers.MethodCall
}

func (ctx *mockContext) getState() *fiberState {
	return ctx.state.Get().(*fiberState)
}

func newMockContext(reporter *EnrichedReporter) *mockContext {
	return &mockContext{
		state: routine.NewThreadLocalWithInitial(func() any {
			return &fiberState{
				matchers:       make([]matchers.Matcher, 0),
				whenHandler:    nil,
				whenCall:       nil,
				methodVerifier: nil,
				verifyState:    false,
			}
		}),
		reporter:           reporter,
		serviceMethodCalls: make([]string, 0),
		lock:               sync.Mutex{},
	}
}
