package verify

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type iface interface {
	Foo(a int) int
}

func TestVerifySimple(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	Verify(m, Once()).Foo(10)
	r.AssertNoError()
}

func TestVerifyAny(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	Verify(m, Once()).Foo(Any[int]())
	r.AssertNoError()
}

func TestVerifyMultipleAny(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	m.Foo(11)
	Verify(m, Times(2)).Foo(Any[int]())
	r.AssertNoError()
}

func TestVerifyNever(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	m.Foo(11)
	Verify(m, Never()).Foo(13)
	r.AssertNoError()
}

func TestVerifyNeverFails(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	m.Foo(11)
	Verify(m, Never()).Foo(10)
	r.AssertError()
}

func TestNoMoreInteractions(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	VerifyNoMoreInteractions(m)
	m.Foo(10)
	r.AssertError()
}
