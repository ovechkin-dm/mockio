package registry

import (
	"fmt"
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

func (e *EnrichedReporter) FailNow(err error) {
	e.Errorf(err.Error())
}

func (e *EnrichedReporter) Fatal(format string) {
	e.reporter.Fatalf(format)
}

func (e *EnrichedReporter) ReportIncorrectWhenUsage() {
	e.Fatal("incorrect usage of `When`. You can only use it with method call: When(foo.Bar()).ThenReturn(...)")
}

func (e *EnrichedReporter) ReportUnregisteredMockVerify(t any) {
	e.Errorf("unregistered mock instance during Verify call: %v", t)
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
	e.Errorf(`invalid use of matchers
method:
%v
expected:
(%s)
got:
(%s)
you can only use matchers within When() call: When(foo.Bar(Any[Int])).
Possible cause is the mixing of matchers with exact values. In this case use "Exact" method instead. 
`, methodSig, inArgsStr, matchersString)
}

func (e *EnrichedReporter) ReportCaptorInsideVerify(call *MethodCall, m []*matcherWrapper) {
	e.Fatal("Unexpected use of captor. `captor.Capture()` should not be used inside `Verify` method")
}

func (e *EnrichedReporter) ReportVerifyMethodError(call *MethodCall, err error) {
	e.FailNow(err)
}

func (e *EnrichedReporter) ReportEmptyCaptor() {
	e.Fatal("no values were captured")
}

func (e *EnrichedReporter) ReportInvalidCaptorValue(expectedType reflect.Type, actualType reflect.Type) {
	e.Fatal("no values were captured")
}

func (e *EnrichedReporter) ReportInvalidReturnValues(instanceType reflect.Type, method reflect.Method, ret []any) {
	retStrValues := make([]string, len(ret))
	for i := range retStrValues {
		if ret[i] == nil {
			retStrValues[i] = "nil"
		} else {
			retStrValues[i] = reflect.ValueOf(ret[i]).Type().Name()
		}
	}
	retStr := strings.Join(retStrValues, ",")
	tp := method.Type
	outTypes := make([]string, 0)
	for i := 0; i < tp.NumOut(); i++ {
		outTypes = append(outTypes, tp.Out(i).Name())
	}
	outTypesStr := strings.Join(outTypes, ", ")
	methodSig := prettyPrintMethodSignature(instanceType, method)
	e.Errorf(`invalid return values
method:
%v
expected:
(%s)
got:
(%s)
`, methodSig, outTypesStr, retStr)
}

func (e *EnrichedReporter) ReportWantedButNotInvoked(
	instanceType reflect.Type,
	methodType reflect.Method,
	match *methodMatch,
	calls []*MethodCall,
) {
	m := match.matchers
	matcherArgs := make([]string, len(m))
	for i := range m {
		matcherArgs[i] = m[i].matcher.Description()
	}
	matchersString := strings.Join(matcherArgs, ", ")
	interfaceName := instanceType.Name()
	methodName := methodType.Name
	methodSig := interfaceName + "." + methodName
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Wanted, but not invoked:\n\t%s(%s)\n", methodSig, matchersString))
	if len(calls) == 0 {
		sb.WriteString("There were zero invocations on this method")
	} else {
		sb.WriteString(fmt.Sprintf("There were %v invocations on this method:\n", len(calls)))
		for _, c := range calls {
			invocation := PrettyPrintMethodInvocation(instanceType, methodType, c.Values)
			sb.WriteString("\t" + invocation)
			sb.WriteString("\n")
		}
	}
	e.Errorf(sb.String())
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

func PrettyPrintMethodInvocation(interfaceType reflect.Type, method reflect.Method, values []reflect.Value) string {
	sb := strings.Builder{}
	interfaceName := interfaceType.Name()
	methodName := method.Name
	sb.WriteString(interfaceName + "." + methodName)
	sb.WriteRune('(')
	anyvals := valueSliceToInterfaceSlice(values)
	for i, v := range anyvals {
		sb.WriteString(fmt.Sprintf("%v", v))
		if i != len(anyvals)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteRune(')')
	return sb.String()
}
