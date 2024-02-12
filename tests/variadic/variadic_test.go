package variadic

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type myInterface interface {
	Foo(a ...int) int
}

func TestVariadicSimple(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[myInterface]()
	WhenSingle(m.Foo(1, 1)).ThenReturn(1)
	WhenSingle(m.Foo(1)).ThenReturn(2)
	ret := m.Foo(1)
	r.AssertEqual(2, ret)
	Verify(m, AtLeastOnce()).Foo(1)
	r.AssertNoError()
}

func TestCaptor(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[myInterface]()
	c1 := Captor[int]()
	c2 := Captor[int]()
	WhenSingle(m.Foo(c1.Capture(), c2.Capture())).ThenReturn(1)
	ret := m.Foo(1, 2)
	r.AssertEqual(1, ret)
	r.AssertEqual(c1.Last(), 1)
	r.AssertEqual(c2.Last(), 2)
	r.AssertNoError()
}