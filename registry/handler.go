package registry

import (
	"fmt"
	"github.com/ovechkin-dm/go-dyno/pkg/dyno"
	"github.com/ovechkin-dm/mockio/matchers"
	"reflect"
	"sync"
)

type invocationHandler struct {
	ctx               *mockContext
	calls             []*methodRecorder
	instanceVerifiers []matchers.InstanceVerifier
	lock              sync.Mutex
}

func (h *invocationHandler) Handle(method *dyno.Method, values []reflect.Value) []reflect.Value {
	h.lock.Lock()
	defer h.lock.Unlock()

	call := &matchers.MethodCall{
		Method: method,
		Values: values,
	}
	if h.ctx.getState().verifyState {
		return h.DoVerifyMethod(call)
	}
	h.calls[method.Num].calls = append(h.calls[method.Num].calls, call)
	return h.DoAnswer(call)
}

func (h *invocationHandler) DoAnswer(c *matchers.MethodCall) []reflect.Value {
	rec := h.calls[c.Method.Num]
	ok := h.VerifyInstance(c)
	if !ok {
		return createDefaultReturnValues(c.Method.Type)
	}
	h.ctx.getState().whenHandler = h
	h.ctx.getState().whenCall = c
	matched := true
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

			result := interfaceSliceToValueSlice(retValues, c.Method.Type)
			if !h.validateReturnValues(result, c.Method.Type) {
				h.ctx.reporter.ReportInvalidReturnValues(result, c.Method.Type)
			}
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

func (h *invocationHandler) VerifyInstance(m *matchers.MethodCall) bool {
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

func (h *invocationHandler) DoVerifyMethod(call *matchers.MethodCall) []reflect.Value {
	argMatchers := h.ctx.getState().matchers

	matchersOk := h.validateMatchers(call)
	verifyMatchersOk := true

	if matchersOk {
		verifyMatchersOk = h.validateVerifyMatchers(call)
	}

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
			calls:         make([]*matchers.MethodCall, 0),
		}
	}
	return &invocationHandler{
		ctx:               holder,
		calls:             recorders,
		instanceVerifiers: make([]matchers.InstanceVerifier, 0),
	}
}

func (h *invocationHandler) validateMatchers(call *matchers.MethodCall) bool {
	argMatchers := h.ctx.getState().matchers
	if len(argMatchers) == 0 {
		ifaces := valueSliceToInterfaceSlice(call.Values)
		for _, v := range ifaces {
			cur := v
			desc := fmt.Sprintf("Exact[%s]", reflect.TypeOf(v).String())
			fm := FunMatcher(desc, func(call []any, a any) bool {
				return cur == a
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
		h.ctx.reporter.ReportInvalidUseOfMatchers(call, argMatchers)
		return false
	}
	return true
}

func (h *invocationHandler) validateReturnValues(result []reflect.Value, method reflect.Method) bool {
	if method.Type.NumOut() != len(result) {
		return false
	}
	for i := range result {
		retExpected := method.Type.Out(i)
		retActual := result[i].Type()
		if retExpected != retActual {
			return false
		}
	}
	return true
}

func (h *invocationHandler) validateVerifyMatchers(call *matchers.MethodCall) bool {
	argMatchers := h.ctx.getState().matchers
	for _, a := range argMatchers {
		if a.rec != nil {
			h.ctx.reporter.ReportInvalidUseOfCaptors(call, argMatchers)
			return false
		}
	}
	return true
}
