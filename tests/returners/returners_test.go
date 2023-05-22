package returners

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type iface interface {
	Test(i interface{}) bool
	Foo(i int) int
}

func TestReturnSimple(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Test(10)).ThenReturn(true)
	ret := m.Test(10)
	r.AssertEqual(true, ret)
}

func TestAnswerSimple(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Test(10)).ThenAnswer(func(args []any) bool {
		return args[0].(int) > 0
	})
	ret1 := m.Test(10)
	ret2 := m.Test(-10)
	r.AssertEqual(true, ret1)
	r.AssertEqual(false, ret2)
}

func TestMultiReturn(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(10)).
		ThenReturn(1).
		ThenReturn(2)
	ret1 := m.Foo(10)
	ret2 := m.Foo(10)
	ret3 := m.Foo(10)
	r.AssertEqual(1, ret1)
	r.AssertEqual(2, ret2)
	r.AssertEqual(2, ret3)
}

func TestMultiAnswer(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(10)).
		ThenAnswer(func(args []any) int {
			return 1
		}).
		ThenAnswer(func(args []any) int {
			return 2
		})
	ret1 := m.Foo(10)
	ret2 := m.Foo(10)
	ret3 := m.Foo(10)
	r.AssertEqual(1, ret1)
	r.AssertEqual(2, ret2)
	r.AssertEqual(2, ret3)
}

func TestReturnBetweenCalls(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	ret := WhenSingle(m.Foo(10))
	ret.ThenReturn(1)
	r1 := m.Foo(10)
	ret.ThenReturn(2)
	r2 := m.Foo(10)
	r3 := m.Foo(10)
	r.AssertEqual(1, r1)
	r.AssertEqual(2, r2)
	r.AssertEqual(2, r3)
}

func TestReturnWrongType(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	When(m.Test(Any[any]())).ThenReturn(10)
	m.Test(10)
	r.AssertError()
}
