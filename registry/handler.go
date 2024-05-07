package registry

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/ovechkin-dm/go-dyno/pkg/dyno"
	"github.com/ovechkin-dm/go-dyno/proxy"

	"github.com/ovechkin-dm/mockio/matchers"
)

type invocationHandler struct {
	ctx          *mockContext
	methods      []*methodRecorder
	lock         sync.Mutex
	instanceType reflect.Type
}

func (h *invocationHandler) Handle(method *dyno.Method, values []reflect.Value) []reflect.Value {
	h.lock.Lock()
	defer h.lock.Unlock()
	values = h.refineValues(method, values)
	call := &MethodCall{
		Method:     method,
		Values:     values,
		StackTrace: NewStackTrace(),
	}
	if h.ctx.getState().verifyState {
		return h.DoVerifyMethod(call)
	}
	h.methods[method.Num].calls = append(h.methods[method.Num].calls, call)
	return h.DoAnswer(call)
}

func (h *invocationHandler) DoAnswer(c *MethodCall) []reflect.Value {
	rec := h.methods[c.Method.Num]
	h.ctx.getState().whenHandler = h
	h.ctx.getState().whenCall = c
	var matched bool
	for _, mm := range rec.methodMatches {
		matched = true
		if len(mm.matchers) != len(c.Values) {
			continue
		}
		for argIdx, matcher := range mm.matchers {
			if !matcher.matcher.Match(valueSliceToInterfaceSlice(c.Values), valueToInterface(c.Values[argIdx])) {
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
		return NewEmptyReturner()
	}

	rec := h.methods[whenCall.Method.Num]

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

func (h *invocationHandler) VerifyMethod(verifier matchers.MethodVerifier) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.ctx.getState().verifyState = true
	h.ctx.getState().methodVerifier = verifier
	if len(h.ctx.getState().matchers) != 0 {
		h.ctx.reporter.ReportUnexpectedMatcherDeclaration(h.ctx.getState().matchers)
	}
}

func (h *invocationHandler) DoVerifyMethod(call *MethodCall) []reflect.Value {
	matchersOk := h.validateMatchers(call)
	argMatchers := h.ctx.getState().matchers

	h.ctx.getState().matchers = make([]*matcherWrapper, 0)
	h.ctx.getState().verifyState = false

	if !matchersOk {
		return createDefaultReturnValues(call.Method.Type)
	}

	rec := h.methods[call.Method.Num]
	matchedInvocations := make([]*MethodCall, 0)
	for _, c := range rec.calls {
		if c.WhenCall {
			continue
		}
		matches := true
		if c.Method.Type != call.Method.Type {
			continue
		}
		if len(argMatchers) != len(c.Values) {
			continue
		}

		for i := range argMatchers {
			if !argMatchers[i].matcher.Match(valueSliceToInterfaceSlice(c.Values), valueToInterface(c.Values[i])) {
				matches = false
				break
			}
		}

		if matches {
			c.Verified = true
			matchedInvocations = append(matchedInvocations, c)
		}
	}
	verifyData := &matchers.MethodVerificationData{
		NumMethodCalls: len(matchedInvocations),
	}
	err := h.ctx.getState().methodVerifier.Verify(verifyData)
	h.ctx.getState().methodVerifier = nil
	if err != nil {
		h.ctx.reporter.ReportVerifyMethodError(h.instanceType, call.Method.Type, matchedInvocations, argMatchers, h.methods[call.Method.Num], err)
	}
	for i, m := range argMatchers {
		if m.rec != nil {
			for _, inv := range matchedInvocations {
				argMatchers[i].rec.Record(inv, valueToInterface(inv.Values[i]))
			}
		}
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
		ctx:          holder,
		methods:      recorders,
		instanceType: tp,
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
	if len(argMatchers) != len(call.Values) {
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
		if result[i] == nil {
			continue
		}
		retExpected := method.Type.Out(i)
		retActual := reflect.TypeOf(result[i])
		if retActual == nil {
			return false
		}
		switch v := result[i].(type) {
		case *proxy.DynamicStruct:
			retActual = v.IFaceValue.Type()
		default:
		}

		if !retActual.AssignableTo(retExpected) {
			return false
		}
	}
	return true
}

func (h *invocationHandler) VerifyNoMoreInteractions() {
	h.PostponedVerify()
	unexpected := make([]*MethodCall, 0)
	for _, rec := range h.methods {
		for _, call := range rec.calls {
			if call.WhenCall {
				continue
			}
			if !call.Verified {
				unexpected = append(unexpected, call)
			}
		}
	}
	if len(unexpected) > 0 {
		h.ctx.reporter.ReportNoMoreInteractionsExpected(h.instanceType, unexpected)
	}
}

func (h *invocationHandler) refineValues(method *dyno.Method, values []reflect.Value) []reflect.Value {
	tp := method.Type.Type
	if tp.IsVariadic() {
		result := make([]reflect.Value, 0)
		for i := 0; i < tp.NumIn()-1; i++ {
			result = append(result, values[i])
		}
		last := values[len(values)-1]
		for i := 0; i < last.Len(); i++ {
			result = append(result, last.Index(i))
		}
		return result
	}
	return values
}

func (h *invocationHandler) PostponedVerify() {
	for _, rec := range h.methods {
		for _, match := range rec.methodMatches {
			if len(match.verifiers) == 0 {
				continue
			}
			matchedInvocations := make([]*MethodCall, 0)
			for _, call := range rec.calls {
				if call.WhenCall {
					continue
				}
				matches := true
				for i := range match.matchers {
					if !match.matchers[i].matcher.Match(valueSliceToInterfaceSlice(call.Values), valueToInterface(call.Values[i])) {
						matches = false
						break
					}
				}
				if matches {
					call.Verified = true
					matchedInvocations = append(matchedInvocations, call)
				}
			}
			verifyData := &matchers.MethodVerificationData{
				NumMethodCalls: len(matchedInvocations),
			}
			for _, v := range match.verifiers {
				err := v.Verify(verifyData)
				if err != nil {
					h.ctx.reporter.ReportVerifyMethodError(h.instanceType, rec.methodType, matchedInvocations, match.matchers, rec, err)
				}
			}
		}
	}
}

func (h *invocationHandler) TearDown() {
	h.PostponedVerify()
}
