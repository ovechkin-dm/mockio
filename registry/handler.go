package registry

import (
	"fmt"
	"github.com/ovechkin-dm/go-dyno/pkg/dyno"
	"github.com/ovechkin-dm/mockio/matchers"
	"reflect"
	"sync"
	"sync/atomic"
)

type invocationHandler struct {
	ctx               *mockContext
	calls             []*methodRecorder
	instanceVerifiers []matchers.InstanceVerifier
	lock              sync.Mutex
	instanceType      reflect.Type
}

func (h *invocationHandler) Handle(method *dyno.Method, values []reflect.Value) []reflect.Value {
	h.lock.Lock()
	defer h.lock.Unlock()

	call := &MethodCall{
		Method: method,
		Values: values,
	}
	if h.ctx.getState().verifyState {
		return h.DoVerifyMethod(call)
	}
	h.calls[method.Num].calls = append(h.calls[method.Num].calls, call)
	return h.DoAnswer(call)
}

func (h *invocationHandler) DoAnswer(c *MethodCall) []reflect.Value {
	rec := h.calls[c.Method.Num]
	ok := h.VerifyInstance(c)
	if !ok {
		return createDefaultReturnValues(c.Method.Type)
	}
	h.ctx.getState().whenHandler = h
	h.ctx.getState().whenCall = c
	var matched bool
	for _, mm := range rec.methodMatches {
		matched = true
		for argIdx, matcher := range mm.matchers {
			if !matcher.matcher.Match(valueSliceToInterfaceSlice(c.Values), c.Values[argIdx].Interface()) {
				matched = false
				break
			}
		}
		if matched {
			ifaces := valueSliceToInterfaceSlice(c.Values)

			for i, m := range mm.matchers {
				if m.rec != nil {
					m.rec.Record(c, ifaces[i])
				}
			}

			ansWrapper := mm.popAnswer()
			if ansWrapper == nil {
				return createDefaultReturnValues(c.Method.Type)
			}

			retValues := ansWrapper.ans(ifaces)

			h.ctx.getState().whenAnswer = ansWrapper
			h.ctx.getState().whenMethodMatch = mm

			if !h.validateReturnValues(retValues, c.Method.Type) {
				h.ctx.reporter.ReportInvalidReturnValues(h.instanceType, c.Method.Type, retValues)
				return createDefaultReturnValues(c.Method.Type)
			}

			result := interfaceSliceToValueSlice(retValues, c.Method.Type)
			return result
		}
	}
	return createDefaultReturnValues(c.Method.Type)
}

func (h *invocationHandler) When() matchers.ReturnerAll {
	h.lock.Lock()
	defer h.lock.Unlock()

	whenCall := h.ctx.getState().whenCall
	whenAnswer := h.ctx.getState().whenAnswer
	whenMethodMatch := h.ctx.getState().whenMethodMatch

	if whenCall == nil {
		h.ctx.reporter.ReportIncorrectWhenUsage()
		return nil
	}
	whenCall.WhenCall = true

	if whenMethodMatch != nil {
		for _, m := range whenMethodMatch.matchers {
			if m.rec != nil {
				m.rec.RemoveRecord(whenCall)
			}
		}
	}

	h.ctx.getState().whenHandler = nil
	h.ctx.getState().whenCall = nil
	h.ctx.getState().whenMethodMatch = nil
	h.ctx.getState().whenAnswer = nil

	if whenAnswer != nil && whenMethodMatch != nil {
		whenMethodMatch.putBackAnswer(whenAnswer)
	}

	if !h.validateMatchers(whenCall) {
		return nil
	}

	rec := h.calls[whenCall.Method.Num]

	argMatchers := h.ctx.getState().matchers

	h.ctx.getState().matchers = make([]*matcherWrapper, 0)
	m := &methodMatch{
		matchers:   argMatchers,
		unanswered: make([]*answerWrapper, 0),
		answered:   make([]*answerWrapper, 0),
	}
	rec.methodMatches = append(rec.methodMatches, m)
	return NewReturnerAll(h.ctx, m)
}

func (h *invocationHandler) VerifyInstance(m *MethodCall) bool {
	data := &matchers.InvocationData{
		MethodType: m.Method.Type,
		MethodName: m.Method.Name,
		Args:       m.Values,
	}
	for _, v := range h.instanceVerifiers {
		err := v.RecordInteraction(data)
		if err != nil {
			h.ctx.reporter.FailNow(err)
			return false
		}
	}
	return true
}

func (h *invocationHandler) AddInstanceVerifier(v matchers.InstanceVerifier) {
	h.instanceVerifiers = append(h.instanceVerifiers, v)
}

func (h *invocationHandler) VerifyMethod(verifier matchers.MethodVerifier) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.ctx.getState().verifyState = true
	h.ctx.getState().methodVerifier = verifier
}

func (h *invocationHandler) DoVerifyMethod(call *MethodCall) []reflect.Value {

	matchersOk := h.validateMatchers(call)

	verifyMatchersOk := true

	if matchersOk {
		verifyMatchersOk = h.validateVerifyMatchers(call)
	}

	argMatchers := h.ctx.getState().matchers

	h.ctx.getState().matchers = make([]*matcherWrapper, 0)
	h.ctx.getState().verifyState = false

	if !matchersOk {
		return createDefaultReturnValues(call.Method.Type)
	}
	if !verifyMatchersOk {
		return createDefaultReturnValues(call.Method.Type)
	}

	numMethodCalls := 0
	rec := h.calls[call.Method.Num]
	for _, c := range rec.calls {
		if c.WhenCall {
			continue
		}
		matches := true
		if c.Method.Type != call.Method.Type {
			continue
		}

		for i := range argMatchers {
			if !argMatchers[i].matcher.Match(valueSliceToInterfaceSlice(c.Values), c.Values[i].Interface()) {
				matches = false
				break
			}
		}

		if matches {
			numMethodCalls += 1
		}
	}
	verifyData := &matchers.MethodVerificationData{
		NumMethodCalls: numMethodCalls,
	}
	err := h.ctx.getState().methodVerifier.Verify(verifyData)
	h.ctx.getState().methodVerifier = nil
	if err != nil {
		h.ctx.reporter.ReportVerifyMethodError(call, err)
	}
	return createDefaultReturnValues(call.Method.Type)
}

func newHandler[T any](holder *mockContext) *invocationHandler {
	tp := reflect.TypeOf(new(T)).Elem()
	recorders := make([]*methodRecorder, tp.NumMethod())
	for i := range recorders {
		recorders[i] = &methodRecorder{
			methodMatches: make([]*methodMatch, 0),
			calls:         make([]*MethodCall, 0),
			methodType:    tp.Method(i),
		}
	}
	return &invocationHandler{
		ctx:               holder,
		calls:             recorders,
		instanceVerifiers: make([]matchers.InstanceVerifier, 0),
		instanceType:      tp,
	}
}

func (h *invocationHandler) validateMatchers(call *MethodCall) bool {
	argMatchers := h.ctx.getState().matchers
	if len(argMatchers) == 0 {
		ifaces := valueSliceToInterfaceSlice(call.Values)
		for _, v := range ifaces {
			cur := v
			desc := fmt.Sprintf("Equal[%s]", reflect.TypeOf(v).String())
			fm := FunMatcher(desc, func(call []any, a any) bool {
				return reflect.DeepEqual(cur, a)
			})
			mw := &matcherWrapper{
				matcher: fm,
				rec:     nil,
			}
			argMatchers = append(argMatchers, mw)
		}
		h.ctx.getState().matchers = argMatchers
	}
	mt := call.Method.Type
	if len(argMatchers) != mt.Type.NumIn() {
		h.ctx.reporter.ReportInvalidUseOfMatchers(h.instanceType, call, argMatchers)
		return false
	}
	return true
}

func (h *invocationHandler) validateReturnValues(result []any, method reflect.Method) bool {
	if method.Type.NumOut() != len(result) {
		return false
	}
	for i := range result {
		if reflect.Zero(method.Type.Out(i)).Interface() == result[i] {
			continue
		}
		retExpected := method.Type.Out(i)
		retActual := reflect.TypeOf(result[i])
		if retActual == nil {
			return false
		}
		if !retActual.AssignableTo(retExpected) {
			return false
		}
	}
	return true
}

func (h *invocationHandler) validateVerifyMatchers(call *MethodCall) bool {
	argMatchers := h.ctx.getState().matchers
	for _, a := range argMatchers {
		if a.rec != nil {
			h.ctx.reporter.ReportCaptorInsideVerify(call, argMatchers)
			return false
		}
	}
	return true
}

func (h *invocationHandler) CheckUnusedStubs() {
	for _, rec := range h.calls {
		for _, m := range rec.methodMatches {
			if atomic.LoadInt64(&m.invocations) == 0 {
				h.ctx.reporter.ReportWantedButNotInvoked(h.instanceType, rec.methodType, m)
			}
		}
	}
}
