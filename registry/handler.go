package registry

import (
	"github.com/ovechkin-dm/go-dyno/pkg/dyno"
	"github.com/ovechkin-dm/mockio/matchers"
	"reflect"
	"sync"
)

type invocationHandler struct {
	ctx               *mockContext
	calls             []*matchers.MethodCall
	methodMatches     []*methodMatch
	instanceVerifiers []matchers.InstanceVerifier
	lock              sync.Mutex
}

func (h *invocationHandler) Handle(method *dyno.Method, values []reflect.Value) []reflect.Value {
	h.lock.Lock()
	defer h.lock.Unlock()

	call := &matchers.MethodCall{
		ID:     GenerateId(),
		Method: method,
		Values: values,
	}
	if h.ctx.getState().verifyState {
		h.ctx.addServiceMethodCallId(call.ID)
		return h.DoVerifyMethod(call)
	}
	h.calls = append(h.calls, call)
	return h.DoAnswer(call)
}

func (h *invocationHandler) DoAnswer(c *matchers.MethodCall) []reflect.Value {
	ok := h.VerifyInstance(c)
	if !ok {
		return createDefaultReturnValues(c.Method.Type)
	}
	h.ctx.getState().whenHandler = h
	h.ctx.getState().whenCall = c
	h.ctx.getState().lastMatch = nil
	matched := true
	for _, mm := range h.methodMatches {
		matched = true
		for argIdx, matcher := range mm.matchers {
			if !matcher.Match(c, c.Values[argIdx].Interface()) {
				matched = false
				break
			}
		}
		if matched {
			for argIdx, matcher := range mm.matchers {
				matcher.Match(c, c.Values[argIdx].Interface())
			}
			h.ctx.getState().lastMatch = mm
			ifaces := valueSliceToInterfaceSlice(c.Values)

			if len(mm.answers) == 0 {
				return createDefaultReturnValues(c.Method.Type)
			}

			actualCalls := make([]*matchers.MethodCall, 0)
			serviceCallIds := h.ctx.getServiceMethodCallIds()
			for _, c := range mm.calls {
				isServiceCall := false
				for _, cid := range serviceCallIds {
					if cid == c.ID {
						isServiceCall = true
						break
					}
				}
				if !isServiceCall {
					actualCalls = append(actualCalls, c)
				}
			}
			curAns := len(actualCalls)
			if curAns >= len(mm.answers) {
				curAns = len(mm.answers) - 1
			}
			ans := mm.answers[curAns](ifaces)
			result := interfaceSliceToValueSlice(ans, c.Method.Type)
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
	if h.ctx.getState().whenCall == nil {
		h.ctx.reporter.ReportIncorrectWhenUsage()
		return nil
	}
	if !h.validateMatchers(h.ctx.getState().whenCall) {
		return nil
	}

	h.ctx.addServiceMethodCallId(h.ctx.getState().whenCall.ID)
	h.ctx.getState().whenHandler = nil
	h.ctx.getState().whenCall = nil

	argMatchers := h.ctx.getState().matchers

	h.ctx.getState().matchers = make([]matchers.Matcher, 0)
	h.ctx.getState().lastMatch = nil
	m := &methodMatch{
		matchers: argMatchers,
		answers:  make([]matchers.Answer, 0),
	}
	h.methodMatches = append(h.methodMatches, m)
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
	ok := h.validateMatchers(call)
	h.ctx.getState().matchers = make([]matchers.Matcher, 0)
	h.ctx.getState().verifyState = false
	if !ok {
		return createDefaultReturnValues(call.Method.Type)
	}
	numMethodCalls := 0
	for _, c := range h.calls {
		matches := true
		if c.Method.Type != call.Method.Type {
			continue
		}

		for i := range argMatchers {
			if !argMatchers[i].Match(c, c.Values[i].Interface()) {
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

func NewHandler(holder *mockContext) *invocationHandler {
	return &invocationHandler{
		ctx:               holder,
		calls:             make([]*matchers.MethodCall, 0),
		methodMatches:     make([]*methodMatch, 0),
		instanceVerifiers: make([]matchers.InstanceVerifier, 0),
	}
}

func (h *invocationHandler) validateMatchers(call *matchers.MethodCall) bool {
	argMatchers := h.ctx.getState().matchers
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
