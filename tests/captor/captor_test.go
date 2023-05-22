package captor

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type iface interface {
	Foo(i int) int
}

func TestCaptorBasic(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	c := Captor[int]()
	WhenSingle(m.Foo(c.Capture())).ThenReturn(10)
	m.Foo(11)
	r.AssertEqual(c.Last(), 11)
}

func TestCaptorMatches(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	c := Captor[int]()
	WhenSingle(m.Foo(c.Capture())).ThenReturn(10)
	ans := m.Foo(11)
	r.AssertEqual(ans, 10)
}

func TestCaptorMultiCalls(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	c := Captor[int]()
	WhenSingle(m.Foo(c.Capture())).ThenReturn(10)
	m.Foo(11)
	m.Foo(12)
	r.AssertEqual(c.Last(), 12)
	r.AssertEqual(c.Values()[0], 11)
	r.AssertEqual(c.Values()[1], 12)
}

func TestCaptorMultiUsage(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m1 := Mock[iface]()
	m2 := Mock[iface]()
	c := Captor[int]()
	WhenSingle(m1.Foo(c.Capture())).ThenReturn(10)
	WhenSingle(m2.Foo(c.Capture())).ThenReturn(10)
	m1.Foo(10)
	m2.Foo(11)
	r.AssertEqual(c.Values()[0], 10)
	r.AssertEqual(c.Values()[1], 11)
}
