package registry

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/ovechkin-dm/mockio/v2/matchers"
	"github.com/ovechkin-dm/mockio/v2/utils"
)

type invocationHandler struct {
	ctx          *mockContext
	methods      map[string]*methodRecorder
	instanceType reflect.Type
	lock         sync.Mutex
	env          *matchers.MockEnv
	reporter     *EnrichedReporter
}

func (h *invocationHandler) Handle(method reflect.Method, values []reflect.Value) []reflect.Value {
	values = h.refineValues(method, values)
	call := &MethodCall{
		Method:     method,
		Values:     values,
		StackTrace: NewStackTrace(),
	}
	if h.ctx.getState().verifyState {
		return h.DoVerifyMethod(call)
	}
	h.methods[method.Name].calls.Add(call)
	return h.DoAnswer(call)
}

func (h *invocationHandler) DoAnswer(c *MethodCall) []reflect.Value {
	rec := h.methods[c.Method.Name]
	h.ctx.getState().whenHandler = h
	h.ctx.getState().whenCall = c
	var matched bool
	methodMatches := rec.methodMatches.GetCopy()
	for _, mm := range methodMatches {
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
				return createDefaultReturnValues(c.Method)
			}

			retValues := ansWrapper.ans(ifaces)

			h.ctx.getState().whenAnswer = ansWrapper
			h.ctx.getState().whenMethodMatch = mm

			if !h.validateReturnValues(retValues, c.Method) {
				h.reporter.ReportInvalidReturnValues(h.instanceType, c.Method, retValues)
				return createDefaultReturnValues(c.Method)
			}

			result := interfaceSliceToValueSlice(retValues, c.Method)
			return result
		}
	}
	return createDefaultReturnValues(c.Method)
}

func (h *invocationHandler) When() matchers.ReturnerAll {
	h.lock.Lock()
	defer h.lock.Unlock()
	whenCall := h.ctx.getState().whenCall
	whenAnswer := h.ctx.getState().whenAnswer
	whenMethodMatch := h.ctx.getState().whenMethodMatch

	if whenCall == nil {
		h.reporter.ReportIncorrectWhenUsage()
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

	rec := h.methods[whenCall.Method.Name]

	argMatchers := h.ctx.getState().matchers

	h.ctx.getState().matchers = make([]*matcherWrapper, 0)
	m := &methodMatch{
		matchers:   argMatchers,
		unanswered: make([]*answerWrapper, 0),
		answered:   make([]*answerWrapper, 0),
		stackTrace: NewStackTrace(),
	}
	rec.methodMatches.Add(m)
	return NewReturnerAll(h.ctx, m)
}

func (h *invocationHandler) VerifyMethod(verifier matchers.MethodVerifier) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.ctx.getState().verifyState = true
	h.ctx.getState().methodVerifier = verifier
	if len(h.ctx.getState().matchers) != 0 {
		h.reporter.ReportUnexpectedMatcherDeclaration(h.ctx.getState().matchers)
	}
}

func (h *invocationHandler) DoVerifyMethod(call *MethodCall) []reflect.Value {
	h.lock.Lock()
	defer h.lock.Unlock()
	matchersOk := h.validateMatchers(call)
	argMatchers := h.ctx.getState().matchers

	h.ctx.getState().matchers = make([]*matcherWrapper, 0)
	h.ctx.getState().verifyState = false

	if !matchersOk {
		return createDefaultReturnValues(call.Method)
	}

	rec := h.methods[call.Method.Name]
	matchedInvocations := make([]*MethodCall, 0)
	calls := rec.calls.GetCopy()
	for _, c := range calls {
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
		h.reporter.ReportVerifyMethodError(
			true,
			h.instanceType,
			call.Method,
			matchedInvocations,
			argMatchers,
			h.methods[call.Method.Name],
			err,
			nil,
		)
	}
	for i, m := range argMatchers {
		if m.rec != nil {
			for _, inv := range matchedInvocations {
				argMatchers[i].rec.Record(inv, valueToInterface(inv.Values[i]))
			}
		}
	}
	return createDefaultReturnValues(call.Method)
}

// newHandler creates a new invocationHandler.
// The `tp` parameter represents the reflector type for the target interface.
func newHandler(tp reflect.Type, holder *mockContext, env *matchers.MockEnv) *invocationHandler {
	recorders := make(map[string]*methodRecorder)
	for i := 0; i < tp.NumMethod(); i++ {
		recorders[tp.Method(i).Name] = &methodRecorder{
			methodMatches: utils.NewSyncList[*methodMatch](),
			calls:         utils.NewSyncList[*MethodCall](),
			methodType:    tp.Method(i),
		}
	}
	return newInvocationHandler(holder, recorders, tp, env)
}

func (h *invocationHandler) validateMatchers(call *MethodCall) bool {
	argMatchers := h.ctx.getState().matchers
	if len(argMatchers) == 0 {
		ifaces := valueSliceToInterfaceSlice(call.Values)
		for _, v := range ifaces {
			cur := v
			desc := fmt.Sprintf("Equal(%v)", v)
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
		h.reporter.ReportInvalidUseOfMatchers(h.instanceType, call, argMatchers)
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
		if !retActual.AssignableTo(retExpected) {
			return false
		}
	}
	return true
}

func (h *invocationHandler) VerifyNoMoreInteractions(tearDown bool) {
	h.PostponedVerify(tearDown)
	unexpected := make([]*MethodCall, 0)
	for _, rec := range h.methods {
		calls := rec.calls.GetCopy()
		for _, call := range calls {
			if call.WhenCall {
				continue
			}
			if !call.Verified {
				unexpected = append(unexpected, call)
			}
		}
	}
	reportFatal := !tearDown
	if len(unexpected) > 0 {
		h.reporter.ReportNoMoreInteractionsExpected(reportFatal, h.instanceType, unexpected)
	}
}

func (h *invocationHandler) refineValues(method reflect.Method, values []reflect.Value) []reflect.Value {
	tp := method.Type
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

func (h *invocationHandler) PostponedVerify(tearDown bool) {
	for _, rec := range h.methods {
		methodMatches := rec.methodMatches.GetCopy()
		for _, match := range methodMatches {
			if len(match.verifiers) == 0 {
				continue
			}
			matchedInvocations := make([]*MethodCall, 0)
			calls := rec.calls.GetCopy()
			for _, call := range calls {
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
					var stackTrace *StackTrace
					if tearDown {
						stackTrace = match.stackTrace
					}
					h.reporter.ReportVerifyMethodError(
						!tearDown,
						h.instanceType,
						rec.methodType,
						matchedInvocations,
						match.matchers,
						rec,
						err,
						stackTrace,
					)
				}
			}
		}
	}
}

func (h *invocationHandler) TearDown() {
	if h.env.Config.StrictVerify {
		for _, m := range h.methods {
			methodMatches := m.methodMatches.GetCopy()
			for _, mm := range methodMatches {
				if len(mm.verifiers) == 0 {
					mm.verifiers = append(mm.verifiers, matchers.AtLeastOnce())
				}
			}
		}
		h.VerifyNoMoreInteractions(true)
	} else {
		h.PostponedVerify(true)
	}
}

func newInvocationHandler(
	ctx *mockContext,
	methods map[string]*methodRecorder,
	instanceType reflect.Type,
	env *matchers.MockEnv,
) *invocationHandler {
	handler := &invocationHandler{
		ctx:          ctx,
		methods:      methods,
		instanceType: instanceType,
		lock:         sync.Mutex{},
		env:          env,
		reporter:     newEnrichedReporter(env.Reporter, env.Config),
	}
	return handler
}
