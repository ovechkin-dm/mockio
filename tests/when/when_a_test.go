package when

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type WhenAInterface interface {
	Foo(a int) int
}

func TestWhenARet(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenAInterface]()
	WhenA(m.Foo(Any[int]())).ThenReturn(42)
	ret := m.Foo(10)
	r.AssertEqual(42, ret)
}

func TestWhenAAnswer(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenAInterface]()
	WhenA(m.Foo(Any[int]())).ThenAnswer(func(args []any) int {
		return 42
	})
	ret := m.Foo(10)
	r.AssertEqual(42, ret)
}

func TestWhenAAnswerWithArgs(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenAInterface]()
	WhenA(m.Foo(Any[int]())).ThenAnswer(func(args []any) int {
		return args[0].(int) + 1
	})
	ret := m.Foo(10)
	r.AssertEqual(11, ret)
}

func TestWhenAMultiAnswer(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenAInterface]()
	WhenA(m.Foo(Any[int]())).
		ThenAnswer(func(args []any) int {
			return args[0].(int) + 1
		}).
		ThenAnswer(func(args []any) int {
			return args[0].(int) + 2
		})
	ret1 := m.Foo(10)
	ret2 := m.Foo(11)
	r.AssertEqual(11, ret1)
	r.AssertEqual(13, ret2)
}

func TestWhenAMultiReturn(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenAInterface]()
	WhenA(m.Foo(Any[int]())).
		ThenReturn(10).
		ThenReturn(11)
	ret1 := m.Foo(12)
	ret2 := m.Foo(13)
	r.AssertEqual(10, ret1)
	r.AssertEqual(11, ret2)
}

func TestWhenAAnswerAndReturn(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenAInterface]()
	WhenA(m.Foo(Any[int]())).
		ThenReturn(10).
		ThenAnswer(func(args []any) int {
			return args[0].(int) + 1
		}).
		ThenReturn(11)
	ret1 := m.Foo(12)
	ret2 := m.Foo(14)
	ret3 := m.Foo(15)
	r.AssertEqual(10, ret1)
	r.AssertEqual(15, ret2)
	r.AssertEqual(11, ret3)
}