package when

import (
	"errors"
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type WhenEInterface interface {
	Foo(a int) (int, error)
}

func TestWhenERet(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).ThenReturn(42, nil)
	ret, _ := m.Foo(10)
	r.AssertEqual(42, ret)
}

func TestWhenEAnswer(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).ThenAnswer(func(args []any) (int, error) {
		return 42, nil
	})
	ret, _ := m.Foo(10)
	r.AssertEqual(42, ret)
}

func TestWhenEAnswerWithArgs(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).ThenAnswer(func(args []any) (int, error) {
		return args[0].(int) + 1, nil
	})
	ret, _ := m.Foo(10)
	r.AssertEqual(11, ret)
}

func TestWhenEMultiAnswer(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).
		ThenAnswer(func(args []any) (int, error) {
			return args[0].(int) + 1, nil
		}).
		ThenAnswer(func(args []any) (int, error) {
			return args[0].(int) + 2, nil
		})
	ret1, _ := m.Foo(10)
	ret2, _ := m.Foo(11)
	r.AssertEqual(11, ret1)
	r.AssertEqual(13, ret2)
}

func TestWhenEMultiReturn(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).
		ThenReturn(10, nil).
		ThenReturn(11, nil)
	ret1, _ := m.Foo(12)
	ret2, _ := m.Foo(13)
	r.AssertEqual(10, ret1)
	r.AssertEqual(11, ret2)
}

func TestWhenEAnswerAndReturn(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).
		ThenReturn(10, nil).
		ThenAnswer(func(args []any) (int, error) {
			return args[0].(int) + 1, nil
		}).
		ThenReturn(11, nil)
	ret1, _ := m.Foo(12)
	ret2, _ := m.Foo(14)
	ret3, _ := m.Foo(15)
	r.AssertEqual(10, ret1)
	r.AssertEqual(15, ret2)
	r.AssertEqual(11, ret3)
}

func TestWhenEReturnError(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).
		ThenReturn(0, errors.New("err"))
	ret, err := m.Foo(12)
	r.AssertEqual(0, ret)
	r.AssertErrorContains(err, "err")
}

func TestWhenEAnswerError(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenEInterface]()
	WhenE(m.Foo(Any[int]())).
		ThenAnswer(func(args []any) (int, error) {
			return 0, errors.New("err")
		})
	ret, err := m.Foo(12)
	r.AssertEqual(0, ret)
	r.AssertErrorContains(err, "err")
}
