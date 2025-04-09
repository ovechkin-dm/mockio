package captor

import (
	"testing"

	"github.com/ovechkin-dm/mockio/v2/tests/common"

	. "github.com/ovechkin-dm/mockio/v2/mock"
)

type iface interface {
	Foo(i int) int
	VoidFoo(i int, j int)
}

func TestCaptorBasic(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	m := Mock[iface](ctrl)
	c := Captor[int]()
	WhenSingle(m.Foo(c.Capture())).ThenReturn(10)
	m.Foo(11)
	r.AssertEqual(c.Last(), 11)
}

func TestCaptorMatches(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)	
	m := Mock[iface](ctrl)	
	c := Captor[int]()
	WhenSingle(m.Foo(c.Capture())).ThenReturn(10)
	ans := m.Foo(11)
	r.AssertEqual(ans, 10)
}

func TestCaptorMultiCalls(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)	
	m := Mock[iface](ctrl)	
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
	ctrl := NewMockController(r)	
	m1 := Mock[iface](ctrl)	
	m2 := Mock[iface](ctrl)
	c := Captor[int]()
	WhenSingle(m1.Foo(c.Capture())).ThenReturn(10)
	WhenSingle(m2.Foo(c.Capture())).ThenReturn(10)
	m1.Foo(10)
	m2.Foo(11)
	r.AssertEqual(c.Values()[0], 10)
	r.AssertEqual(c.Values()[1], 11)
}

func TestCaptorVerify(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)	
	m := Mock[iface](ctrl)	
	c := Captor[int]()
	m.VoidFoo(10, 20)
	Verify(m, Once()).VoidFoo(c.Capture(), Exact(20))
	r.AssertNoError()
	r.AssertEqual(10, c.Last())
}
