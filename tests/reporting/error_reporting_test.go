package reporting

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"sync"
	"testing"
)

type Foo interface {
	Bar()
	Baz(a int, b int, c int) int
}

func TestReportIncorrectWhenUsage(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	When(1)
	r.AssertError()
	r.PrintError()
}

func TestReportVerifyFromDifferentGoroutine(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	mock := Mock[Foo]()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		SetUp(r)
		Verify(mock, Once())
		wg.Done()
	}()
	wg.Wait()
	r.AssertError()
	r.PrintError()
}

func TestReportVerifyNotAMock(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	Verify(100, Once())
	r.AssertError()
	r.PrintError()
}

func TestInvalidUseOfMatchers(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	mock := Mock[Foo]()
	When(mock.Baz(AnyInt(), AnyInt(), 10)).ThenReturn(10)
	mock.Baz(1, 2, 3)
	r.AssertError()
	r.PrintError()
}

func TestCaptorInsideVerify(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	mock := Mock[Foo]()
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10)
	c := Captor[int]()
	Verify(mock, Once()).Baz(AnyInt(), AnyInt(), c.Capture())
	r.AssertError()
	r.PrintError()
}

func TestVerify(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	mock := Mock[Foo]()
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10)
	_ = mock.Baz(10, 10, 11)
	Verify(mock, Once()).Baz(AnyInt(), AnyInt(), Exact(10))
	r.AssertError()
	r.PrintError()
}

func TestVerifyTimes(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	mock := Mock[Foo]()
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10)
	_ = mock.Baz(10, 10, 10)
	Verify(mock, Times(20)).Baz(AnyInt(), AnyInt(), AnyInt())
	r.AssertError()
	r.PrintError()
}

func TestEmptyCaptor(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	c := Captor[int]()
	_ = c.Last()
	r.AssertError()
	r.PrintError()
}

func TestInvalidReturnValues(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	mock := Mock[Foo]()
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn("test", 10)
	_ = mock.Baz(10, 10, 10)
	r.AssertError()
	r.PrintError()
}

func TestNoMoreInteractions(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	mock := Mock[Foo]()
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn("test", 10)
	_ = mock.Baz(10, 10, 10)
	_ = mock.Baz(10, 20, 10)
	VerifyNoMoreInteractions(mock)
	r.AssertError()
	r.PrintError()
}

func TestUnexpectedMatchers(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	mock := Mock[Foo]()
	When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10)
	mock.Baz(AnyInt(), AnyInt(), AnyInt())
	Verify(mock, Once()).Baz(10, 10, 10)
	r.AssertError()
	r.PrintError()
}