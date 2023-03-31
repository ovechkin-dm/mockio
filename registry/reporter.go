package registry

import (
	"github.com/ovechkin-dm/mockio/matchers"
	"reflect"
	"strings"
)

type EnrichedReporter struct {
	reporter matchers.ErrorReporter
}

func (e *EnrichedReporter) Fatalf(format string, args ...any) {
	e.reporter.Fatalf(format, args...)
}

func (e *EnrichedReporter) FailNow(err error) {
	e.Fatalf(err.Error())
}

func (e *EnrichedReporter) Fatal(format string) {
	e.reporter.Fatalf(format)
}

func (e *EnrichedReporter) ReportIncorrectWhenUsage() {
	e.Fatal("incorrect usage of mock.When. You can only use it with method call: mock.When(foo.Bar()).ThenReturn(...)")
}

func (e *EnrichedReporter) ReportUnregisteredMockVerify(t any) {
	e.Fatalf("unregistered mock instance during Verify call: %v", t)
}

func (e *EnrichedReporter) ReportInvalidUseOfMatchers(call *matchers.MethodCall, m []matchers.Matcher) {
	matcherArgs := make([]string, len(m))
	for i := range m {
		matcherArgs[i] = m[i].Description()
	}
	matchersString := strings.Join(matcherArgs, ",")
	tp := call.Method.Type.Type
	inArgs := make([]string, 0)
	for i := 0; i < tp.NumIn(); i++ {
		inArgs = append(inArgs, tp.In(i).String())
	}
	inArgsStr := strings.Join(inArgs, ",")
	e.Fatalf(`invalid use of matchers
method:
%v
expected:
(%s)
got:
(%s)
you can only use matchers within When() call: mock.When(foo.Bar(mock.Any[Int]))
`, call.Method.Type, inArgsStr, matchersString)
}

func (e *EnrichedReporter) ReportVerifyMethodError(call *matchers.MethodCall, err error) {
	e.FailNow(err)
}

func (e *EnrichedReporter) ReportEmptyCaptor() {
	e.Fatal("no values were captured")
}

func (e *EnrichedReporter) ReportInvalidReturnValues(ret []reflect.Value, method reflect.Method) {
	retStrValues := make([]string, len(ret))
	for i := range retStrValues {
		retStrValues[i] = ret[i].String()
	}
	retStr := strings.Join(retStrValues, ",")
	tp := method.Type
	outTypes := make([]string, 0)
	for i := 0; i < tp.NumOut(); i++ {
		outTypes = append(outTypes, tp.Out(i).String())
	}
	outTypesStr := strings.Join(outTypes, ",")
	e.Fatalf(`invalid return values
method:
%v
expected number of return values:
%d
actual number of return values:
%d
expected types:
(%s)
got:
(%s)
`, tp, method.Type.NumIn(), len(ret), outTypesStr, retStr)
}

func newEnrichedReporter(reporter matchers.ErrorReporter) *EnrichedReporter {
	return &EnrichedReporter{
		reporter: reporter,
	}
}
