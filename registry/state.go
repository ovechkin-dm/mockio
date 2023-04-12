package registry

import (
	"github.com/ovechkin-dm/mockio/matchers"
	"github.com/timandy/routine"
	"sync"
)

type fiberState struct {
	matchers        []matchers.Matcher
	whenHandler     *invocationHandler
	verifyState     bool
	methodVerifier  matchers.MethodVerifier
	whenCall        *matchers.MethodCall
	whenAnswer      *answerWrapper
	whenMethodMatch *methodMatch
}

type mockContext struct {
	state              routine.ThreadLocal
	serviceMethodCalls map[string]struct{}
	reporter           *EnrichedReporter
	lock               sync.Mutex
}

func (ctx *mockContext) addServiceMethodCallId(id string) {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()
	ctx.serviceMethodCalls[id] = struct{}{}
}

func (ctx *mockContext) IsServiceCall(id string) bool {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()
	_, ok := ctx.serviceMethodCalls[id]
	return ok
}

type methodRecorder struct {
	methodMatches []*methodMatch
	calls         []*matchers.MethodCall
}

type methodMatch struct {
	matchers   []matchers.Matcher
	unanswered []*answerWrapper
	answered   []*answerWrapper
	lock       sync.Mutex
}

func (m *methodMatch) popAnswer() *answerWrapper {
	m.lock.Lock()
	defer m.lock.Unlock()
	if len(m.unanswered) == 0 {
		return nil
	}
	last := m.unanswered[0]
	m.unanswered = m.unanswered[1:]
	m.answered = append(m.answered, last)
	return last
}

func (m *methodMatch) addAnswer(wrapper *answerWrapper) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.unanswered = append(m.unanswered, wrapper)
}

func (m *methodMatch) putBackAnswer(wrapper *answerWrapper) {
	m.lock.Lock()
	defer m.lock.Unlock()
	foundIdx := -1
	for i := len(m.answered) - 1; i >= 0; i-- {
		if wrapper == m.answered[i] {
			foundIdx = i
			break
		}
	}
	if foundIdx == -1 {
		return
	}
	for i := foundIdx; i < len(m.unanswered)-1; i++ {
		m.answered[i] = m.answered[i+1]
	}
	m.answered = m.answered[0 : len(m.answered)-1]
	m.unanswered = append(m.unanswered, wrapper)
}

type answerWrapper struct {
	ans matchers.Answer
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
		serviceMethodCalls: make(map[string]struct{}, 0),
		lock:               sync.Mutex{},
	}
}
