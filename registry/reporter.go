package registry

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ovechkin-dm/mockio/v2/config"
	"github.com/ovechkin-dm/mockio/v2/matchers"
)

type panicReporter struct{}

func (p *panicReporter) Cleanup(f func()) {
}

func (p *panicReporter) Fatalf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

func (p *panicReporter) Errorf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

type EnrichedReporter struct {
	reporter matchers.ErrorReporter
	cfg      *config.MockConfig
}

func (e *EnrichedReporter) Errorf(format string, args ...any) {
	e.reporter.Errorf(format, args...)
}

func (e *EnrichedReporter) StackTraceFatalf(format string, args ...any) {
	e.StackTraceErrorf(nil, true, format, args...)
}

func (e *EnrichedReporter) StackTraceErrorf(s *StackTrace, fatal bool, format string, args ...any) {
	if s == nil {
		s = NewStackTrace()
	}
	result := fmt.Sprintf(format, args...)
	var st string
	if e.cfg.PrintStackTrace {
		st = fmt.Sprintf(`At:
	%s
Cause:
	%s
Trace:
%s
`, s.CallerLine(), result, s.WithoutLibraryCalls().String())
	} else {
		st = fmt.Sprintf(`At:
	%s
Cause:
	%s
`, s.CallerLine(), result)
	}
	if fatal {
		e.Fatalf(st)
	} else {
		e.reporter.Errorf(st)
	}
}

func (e *EnrichedReporter) FailNow(err error) {
	e.Fatalf(err.Error())
}

func (e *EnrichedReporter) Fatal(format string) {
	e.Fatalf(format)
}

func (e *EnrichedReporter) Fatalf(format string, args ...any) {
	e.reporter.Fatalf(format, args...)
}

func (e *EnrichedReporter) ReportIncorrectWhenUsage() {
	e.StackTraceFatalf(`When() requires an argument which has to be 'a method call on a mock'.
	For example: When(mock.GetArticles()).ThenReturn(articles)`)
}

func (e *EnrichedReporter) ReportUnregisteredMockVerify(t any) {
	e.StackTraceFatalf(`Argument passed to Verify() is %v and is not a mock, or a mock created in a different goroutine.
	Make sure you place the parenthesis correctly.
	Example of correct verification:
		Verify(mock, Times(10)).SomeMethod()`, t)
}

func (e *EnrichedReporter) ReportInvalidUseOfMatchers(instanceType reflect.Type, call *MethodCall, m []*matcherWrapper) {
	matcherArgs := make([]string, len(m))
	for i := range m {
		matcherArgs[i] = m[i].matcher.Description()
	}
	matchersString := strings.Join(matcherArgs, ",")
	tp := call.Method.Type
	inArgs := make([]string, 0)
	methodSig := prettyPrintMethodSignature(instanceType, call.Method)
	for i := 0; i < tp.NumIn(); i++ {
		inArgs = append(inArgs, tp.In(i).String())
	}
	inArgsStr := strings.Join(inArgs, ",")
	numExpected := call.Method.Type.NumIn()
	numActual := len(m)
	declarationLines := make([]string, 0)
	for i := range m {
		declarationLines = append(declarationLines, "\t\t"+m[i].stackTrace.CallerLine())
	}
	decl := strings.Join(declarationLines, "\n")
	expectedStr := fmt.Sprintf("%v expected, %v recorded:\n", numExpected, numActual)
	if call.Method.Type.IsVariadic() {
		expectedStr = ""
	}
	e.StackTraceFatalf(`Invalid use of matchers
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

func (e *EnrichedReporter) ReportVerifyMethodError(
	fatal bool,
	tp reflect.Type,
	method reflect.Method,
	invocations []*MethodCall,
	argMatchers []*matcherWrapper,
	recorder *methodRecorder,
	err error,
	stackTrace *StackTrace,
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
	callStr := PrettyPrintMethodInvocation(tp, method, args)

	other := strings.Builder{}
	calls := recorder.calls.GetCopy()
	for j, c := range calls {
		if c.WhenCall {
			continue
		}
		callArgs := make([]string, len(c.Values))
		for i := range c.Values {
			callArgs[i] = fmt.Sprintf("%v", c.Values[i])
		}
		pretty := PrettyPrintMethodInvocation(tp, c.Method, callArgs)
		other.WriteString(fmt.Sprintf("\t\t%s at %s", pretty, c.StackTrace.CallerLine()))
		if j != len(calls)-1 {
			other.WriteString("\n")
		}
	}
	if len(other.String()) == 0 && len(sb.String()) == 0 {
		e.StackTraceErrorf(stackTrace, fatal, `%v
		%v
`, err, callStr)
	} else if len(invocations) == 0 {
		e.StackTraceErrorf(stackTrace, fatal, `%v
		%v
	However, there were other interactions with this method:
%v`, err, callStr, other.String())
	} else {
		e.StackTraceErrorf(stackTrace, fatal, `%v
		%v
	Invocations:
%v`, err, callStr, sb.String())
	}
}

func (e *EnrichedReporter) ReportEmptyCaptor() {
	e.StackTraceFatalf("no values were captured for captor")
}

func (e *EnrichedReporter) ReportInvalidCaptorValue(expectedType reflect.Type, actualType reflect.Type) {
	e.StackTraceFatalf("captor contains unexpected type")
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

	e.StackTraceFatalf(`invalid return values
expected:
	%v
got:
	%s
`, methodSig, outTypesSB.String())
}

func newEnrichedReporter(reporter matchers.ErrorReporter, cfg *config.MockConfig) *EnrichedReporter {
	return &EnrichedReporter{
		reporter: reporter,
		cfg:      cfg,
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

func (e *EnrichedReporter) ReportNoMoreInteractionsExpected(fatal bool, instanceType reflect.Type, calls []*MethodCall) {
	sb := strings.Builder{}
	for i, c := range calls {
		args := make([]string, 0)
		for _, v := range c.Values {
			args = append(args, fmt.Sprintf("%v", v))
		}
		s := PrettyPrintMethodInvocation(instanceType, c.Method, args)
		line := fmt.Sprintf("\t\t%s at %s", s, c.StackTrace.CallerLine())
		sb.WriteString(line)
		if i != len(calls)-1 {
			sb.WriteString("\n")
		}

	}
	e.StackTraceErrorf(nil, fatal, `No more interactions expected, but unverified interactions found:
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
	e.StackTraceFatalf(`Unexpected matchers declaration.
%s
	Matchers can only be used inside When() method call.`, sb.String())
}
