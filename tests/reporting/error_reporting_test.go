package reporting

import (
	"sync"
	"testing"

	"github.com/ovechkin-dm/mockio/v2/mockopts"
	"github.com/ovechkin-dm/mockio/v2/tests/common"

	. "github.com/ovechkin-dm/mockio/v2/mock"
)

type Foo interface {
	Bar()
	Baz(a int, b int, c int) int
	VarArgs(a string, b ...int) int
}

func TestReportIncorrectWhenUsage(t *testing.T) {	
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Expected panic but got none")
		}
	}()
	When(1)	
}

func TestVerifyFromDifferentGoroutine(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {		
		Verify(mock, Once())
		wg.Done()
	}()
	wg.Wait()
	r.AssertNoError()
}

func TestReportVerifyNotAMock(t *testing.T) {	
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Expected panic but got none")
		}
	}()
	Verify(100, Once())
}

func TestInvalidUseOfMatchers(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	When(mock.Baz(AnyInt(), AnyInt(), 10)).ThenReturn(10)
	mock.Baz(1, 2, 3)
	r.AssertError()
	r.PrintError()
}

func TestInvalidUseOfMatchersVarArgs(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	When(mock.VarArgs(AnyString(), AnyInt(), 10)).ThenReturn(10)
	mock.VarArgs("a", 2)
	r.AssertError()
	r.PrintError()
}

func TestCaptorInsideVerify(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10)
	c := Captor[int]()
	Verify(mock, Once()).Baz(AnyInt(), AnyInt(), c.Capture())
	r.AssertError()
	r.PrintError()
}

func TestVerify(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10)
	_ = mock.Baz(10, 10, 11)
	Verify(mock, Once()).Baz(AnyInt(), AnyInt(), Exact(10))
	r.AssertError()
	r.PrintError()
}

func TestVerifyVarArgs(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	When(mock.VarArgs(AnyString(), AnyInt(), AnyInt())).ThenReturn(10)
	_ = mock.VarArgs("a", 10, 11)
	Verify(mock, Once()).VarArgs(AnyString(), AnyInt(), Exact(10))
	r.AssertError()
	r.PrintError()
}

func TestVerifyDifferentVarArgs(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	When(mock.VarArgs(AnyString(), AnyInt(), AnyInt())).ThenReturn(10)
	_ = mock.VarArgs("a", 10, 11)
	Verify(mock, Once()).VarArgs(AnyString(), AnyInt(), AnyInt(), AnyInt())
	r.AssertError()
	r.PrintError()
}

func TestVerifyTimes(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10)
	_ = mock.Baz(10, 10, 10)
	Verify(mock, Times(20)).Baz(AnyInt(), AnyInt(), AnyInt())
	r.AssertError()
	r.PrintError()
}

func TestEmptyCaptor(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Expected panic but got none")
		}
	}()	
	c := Captor[int]()
	_ = c.Last()
}

func TestInvalidReturnValues(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn("test", 10)
	_ = mock.Baz(10, 10, 10)
	r.AssertError()
	r.PrintError()
}

func TestNoMoreInteractions(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn("test", 10)
	_ = mock.Baz(10, 10, 10)
	_ = mock.Baz(10, 20, 10)
	VerifyNoMoreInteractions(mock)
	r.AssertError()
	r.PrintError()
}

func TestNoMoreInteractionsVarArgs(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	When(mock.VarArgs(AnyString(), AnyInt(), AnyInt())).ThenReturn("test", 10)
	_ = mock.Baz(10, 10, 10)
	_ = mock.Baz(10, 20, 10)
	VerifyNoMoreInteractions(mock)
	r.AssertError()
	r.PrintError()
}

func TestUnexpectedMatchers(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10)
	mock.Baz(AnyInt(), AnyInt(), AnyInt())
	Verify(mock, Once()).Baz(10, 10, 10)
	r.AssertError()
	r.PrintError()
}

func TestStackTraceDisabled(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r, mockopts.WithoutStackTrace())
	mock := Mock[Foo](ctrl)
	WhenSingle(mock.Baz(1, 2, AnyInt())).ThenReturn(10)
	_ = mock.Baz(1, 2, 3)
	r.AssertError()
	r.PrintError()
}

func TestStackTraceEnabled(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	WhenSingle(mock.Baz(1, 2, AnyInt())).ThenReturn(10)
	_ = mock.Baz(1, 2, 3)
	r.AssertError()
	r.PrintError()
}

func TestImplicitMatchersLogValue(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	mock := Mock[Foo](ctrl)
	WhenSingle(mock.Baz(1, 2, 3)).ThenReturn(10).Verify(Once())
	r.TriggerCleanup()
	r.AssertError()
	r.PrintError()
}
