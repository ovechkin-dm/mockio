package registry

import (
	"fmt"
	"github.com/ovechkin-dm/go-dyno/proxy"
	"github.com/ovechkin-dm/mockio/matchers"
	"reflect"
	"strings"
)

type panicReporter struct {
}

func (p *panicReporter) Cleanup(f func()) {

}

func (p *panicReporter) Fatalf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

type EnrichedReporter struct {
	reporter matchers.ErrorReporter
}

func (e *EnrichedReporter) Errorf(format string, args ...any) {
	e.reporter.Fatalf(format, args...)
}

func (e *EnrichedReporter) StackTraceErrorf(format string, args ...any) {
	s := NewStackTrace()
	result := fmt.Sprintf(format, args...)
	st := fmt.Sprintf(`At:
	%s
Cause:
	%s
Trace:
%s
`, s.CallerLine(), result, s.WithoutLibraryCalls().String())
	e.reporter.Fatalf(st)
}

func (e *EnrichedReporter) FailNow(err error) {
	e.Errorf(err.Error())
}

func (e *EnrichedReporter) Fatal(format string) {
	e.reporter.Fatalf(format)
}

func (e *EnrichedReporter) ReportIncorrectWhenUsage() {
	e.StackTraceErrorf(`When() requires an argument which has to be 'a method call on a mock'.
	For example: When(mock.GetArticles()).ThenReturn(articles)`)
}

func (e *EnrichedReporter) ReportUnregisteredMockVerify(t any) {
	switch t.(type) {
	case *proxy.DynamicStruct:
		e.StackTraceErrorf(`Argument passed to Verify() is a mock from different goroutine.
	Make sure you made call to Mock() and Verify() from the same goroutine.`)
	default:
		e.StackTraceErrorf(`Argument passed to Verify() is %v and is not a mock.
	Make sure you place the parenthesis correctly.
	Example of correct verification:
		Verify(mock, Times(10)).SomeMethod()`, t)
	}

}

func (e *EnrichedReporter) ReportInvalidUseOfMatchers(instanceType reflect.Type, call *MethodCall, m []*matcherWrapper) {
	matcherArgs := make([]string, len(m))
	for i := range m {
		matcherArgs[i] = m[i].matcher.Description()
	}
	matchersString := strings.Join(matcherArgs, ",")
	tp := call.Method.Type.Type
	inArgs := make([]string, 0)
	methodSig := prettyPrintMethodSignature(instanceType, call.Method.Type)
	for i := 0; i < tp.NumIn(); i++ {
		inArgs = append(inArgs, tp.In(i).String())
	}
	inArgsStr := strings.Join(inArgs, ",")
	numExpected := call.Method.Type.Type.NumIn()
	numActual := len(m)
	declarationLines := make([]string, 0)
	for i := range m {
		declarationLines = append(declarationLines, "\t\t" + m[i].stackTrace.CallerLine())
	}
	decl := strings.Join(declarationLines, "\n")
	expectedStr := fmt.Sprintf("%v expected, %v recorded:\n", numExpected, numActual)
	if call.Method.Type.Type.IsVariadic() {
		expectedStr = ""
	}
	e.StackTraceErrorf(`Invalid use of matchers
	%s%v
	method:
		%v
	expected:
		(%s)
	got:
		(%s)
	This can happen for 2 reasons:
		1. Declaration of matcher outside When() call
		2. Mixing matchers and exact values in When() call. Is this case, consider using "Exact" matcher.`,
		expectedStr, decl, methodSig, inArgsStr, matchersString)
}

func (e *EnrichedReporter) ReportCaptorInsideVerify(call *MethodCall, m []*matcherWrapper) {
	e.StackTraceErrorf("Unexpected use of captor. `captor.Capture()` should not be used inside `Verify` method")
}

func (e *EnrichedReporter) ReportVerifyMethodError(
	tp reflect.Type,
	call *MethodCall,
	invocations []*MethodCall,
	argMatchers []*matcherWrapper,
	recorder *methodRecorder,
	err error,
) {
	sb := strings.Builder{}
	for i, c := range invocations {
		if c.WhenCall {
			continue
		}
		sb.WriteString("\t\t" + c.StackTrace.CallerLine())
		if i != len(invocations)-1 {
			sb.WriteString("\n")
		}
	}
	args := make([]string, len(argMatchers))
	for i := range argMatchers {
		args[i] = argMatchers[i].matcher.Description()
	}
	callStr := PrettyPrintMethodInvocation(tp, call.Method.Type, args)

	other := strings.Builder{}
	for j, c := range recorder.calls {
		if c.WhenCall {
			continue
		}
		callArgs := make([]string, len(c.Values))
		for i := range c.Values {
			callArgs[i] = fmt.Sprintf("%v", c.Values[i])
		}
		pretty := PrettyPrintMethodInvocation(tp, c.Method.Type, callArgs)
		other.WriteString(fmt.Sprintf("\t\t%s at %s", pretty, c.StackTrace.CallerLine()))
		if j != len(recorder.calls)-1 {
			other.WriteString("\n")
		}
	}

	if len(invocations) == 0 {
		e.StackTraceErrorf(`%v
		%v
	However, there were other interactions with this method:
%v`, err, callStr, other.String())
	} else {
		e.StackTraceErrorf(`%v
		%v
	Invocations:
%v`, err, callStr, sb.String())
	}

}

func (e *EnrichedReporter) ReportEmptyCaptor() {
	e.StackTraceErrorf("no values were captured for captor")
}

func (e *EnrichedReporter) ReportInvalidCaptorValue(expectedType reflect.Type, actualType reflect.Type) {
	e.StackTraceErrorf("captor contains unexpected type")
}

func (e *EnrichedReporter) ReportInvalidReturnValues(instanceType reflect.Type, method reflect.Method, ret []any) {
	tp := method.Type
	outTypesSB := strings.Builder{}

	interfaceName := instanceType.Name()
	methodName := method.Name
	outTypesSB.WriteString(interfaceName + "." + methodName)
	outTypesSB.WriteString("(")
	for i := 0; i < tp.NumIn(); i++ {
		outTypesSB.WriteString(tp.In(i).Name())
		if i != tp.NumIn()-1 {
			outTypesSB.WriteString(", ")
		}
	}
	outTypesSB.WriteString(")")
	if len(ret) > 0 {
		outTypesSB.WriteString(" ")
	}
	if len(ret) > 1 {
		outTypesSB.WriteString("(")
	}
	for i := 0; i < len(ret); i++ {
		of := reflect.ValueOf(ret[i])
		if of.Kind() == reflect.Invalid {
			outTypesSB.WriteString("nil")
		} else {
			outTypesSB.WriteString(of.Type().String())
		}

		if i != len(ret)-1 {
			outTypesSB.WriteString(", ")
		}
	}
	if len(ret) > 1 {
		outTypesSB.WriteString(")")
	}

	methodSig := prettyPrintMethodSignature(instanceType, method)

	e.StackTraceErrorf(`invalid return values
expected:
	%v
got:
	%s
`, methodSig, outTypesSB.String())
}

func newEnrichedReporter(reporter matchers.ErrorReporter) *EnrichedReporter {
	return &EnrichedReporter{
		reporter: reporter,
	}
}

func prettyPrintMethodSignature(interfaceType reflect.Type, method reflect.Method) string {
	var signature string

	interfaceName := interfaceType.Name()
	methodName := method.Name
	methodType := method.Type
	signature += interfaceName + "." + methodName

	numParams := methodType.NumIn()
	signature += "("
	for i := 0; i < numParams; i++ {
		paramType := methodType.In(i)
		signature += paramType.String()
		if i != numParams-1 {
			signature += ", "
		}
	}
	signature += ")"

	numReturns := methodType.NumOut()
	if numReturns > 0 {
		signature += " "
	}
	if numReturns > 1 {
		signature += "("
	}
	for i := 0; i < numReturns; i++ {
		returnType := methodType.Out(i)
		signature += returnType.String()
		if i != numReturns-1 {
			signature += ", "
		}
	}
	if numReturns > 1 {
		signature += ")"
	}

	return signature
}

func PrettyPrintMethodInvocation(interfaceType reflect.Type, method reflect.Method, args []string) string {
	sb := strings.Builder{}
	interfaceName := interfaceType.Name()
	methodName := method.Name
	sb.WriteString(interfaceName + "." + methodName)
	sb.WriteRune('(')
	for i, v := range args {
		sb.WriteString(fmt.Sprintf("%v", v))
		if i != len(args)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteRune(')')
	return sb.String()
}

func (e *EnrichedReporter) ReportNoMoreInteractionsExpected(instanceType reflect.Type, calls []*MethodCall) {
	sb := strings.Builder{}
	for i, c := range calls {
		args := make([]string, 0)
		for _, v := range c.Values {
			args = append(args, fmt.Sprintf("%v", v))
		}
		s := PrettyPrintMethodInvocation(instanceType, c.Method.Type, args)
		line := fmt.Sprintf("\t\t%s at %s", s, c.StackTrace.CallerLine())
		sb.WriteString(line)
		if i != len(calls)-1 {
			sb.WriteString("\n")
		}

	}
	e.StackTraceErrorf(`No more interactions expected, but unverified interactions found:
%v`, sb.String())
}

func (e *EnrichedReporter) ReportUnexpectedMatcherDeclaration(m []*matcherWrapper) {
	sb := strings.Builder{}
	for i, v := range m {
		sb.WriteString("\t\tat " + v.stackTrace.CallerLine())
		if i != len(m)-1 {
			sb.WriteString("\n")
		}
	}
	e.StackTraceErrorf(`Unexpected matchers declaration.
%s
	Matchers can only be used inside When() method call.`, sb.String())
}
