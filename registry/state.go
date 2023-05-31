package registry

import (
	"github.com/ovechkin-dm/go-dyno/pkg/dyno"
	"github.com/ovechkin-dm/mockio/matchers"
	"github.com/timandy/routine"
	"reflect"
	"sync"
	"sync/atomic"
)

type fiberState struct {
	matchers        []*matcherWrapper
	whenHandler     *invocationHandler
	verifyState     bool
	methodVerifier  matchers.MethodVerifier
	whenCall        *MethodCall
	whenAnswer      *answerWrapper
	whenMethodMatch *methodMatch
}

type mockContext struct {
	state     routine.ThreadLocal
	reporter  *EnrichedReporter
	lock      sync.Mutex
	routineID int64
}

type methodRecorder struct {
	methodMatches []*methodMatch
	calls         []*MethodCall
	methodType    reflect.Method
}

type methodMatch struct {
	matchers    []*matcherWrapper
	unanswered  []*answerWrapper
	answered    []*answerWrapper
	lock        sync.Mutex
	lastAnswer  *answerWrapper
	invocations int64
}

func (m *methodMatch) popAnswer() *answerWrapper {
	m.lock.Lock()
	defer m.lock.Unlock()
	atomic.AddInt64(&m.invocations, 1)
	if len(m.unanswered) == 0 {
		return m.lastAnswer
	}
	last := m.unanswered[0]
	m.unanswered = m.unanswered[1:]
	m.answered = append(m.answered, last)
	m.lastAnswer = last
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
	atomic.AddInt64(&m.invocations, -1)
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

type matcherWrapper struct {
	matcher matchers.Matcher[any]
	rec     recordable
}

func (ctx *mockContext) getState() *fiberState {
	return ctx.state.Get().(*fiberState)
}

func newMockContext(reporter *EnrichedReporter) *mockContext {
	return &mockContext{
		state: routine.NewThreadLocalWithInitial(func() any {
			return &fiberState{
				matchers:       make([]*matcherWrapper, 0),
				whenHandler:    nil,
				whenCall:       nil,
				methodVerifier: nil,
				verifyState:    false,
			}
		}),
		reporter:  reporter,
		lock:      sync.Mutex{},
		routineID: routine.Goid(),
	}
}

type MethodCall struct {
	Method   *dyno.Method
	Values   []reflect.Value
	WhenCall bool
}
