package when

import (
	"errors"
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type WhenEInterface interface {
	FooE(a int) (int, error)
}

func TestWhenERet(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).ThenReturn(42)
	ret := m.Foo(10)
	r.AssertEqual(42, ret)
}

func TestWhenEAnswer(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).ThenAnswer(func(args []any) int {
		return 42
	})
	ret := m.Foo(10)
	r.AssertEqual(42, ret)
}

func TestWhenEAnswerWithArgs(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).ThenAnswer(func(args []any) int {
		return args[0].(int) + 1
	})
	ret := m.Foo(10)
	r.AssertEqual(11, ret)
}

func TestWhenEMultiAnswer(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenEInterface]()
	WhenE(m.FooE(Any[int]())).
		ThenAnswer(func(args []any) (int, error) {
			return 0, errors.New("err")
		}).
		ThenAnswer(func(args []any) (int, error) {
			return 0, errors.New("err")
		})
	ret1 := m.FooE(10)
	ret2 := m.FooE(11)
	r.AssertEqual(11, ret1)
	r.AssertEqual(13, ret2)
}

func TestWhenEMultiReturn(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).
		ThenReturn(10).
		ThenReturn(11)
	ret1 := m.Foo(12)
	ret2 := m.Foo(13)
	r.AssertEqual(10, ret1)
	r.AssertEqual(11, ret2)
}

func TestWhenEAnswerAndReturn(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := NewMock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).
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
