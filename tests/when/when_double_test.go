package when

import (
	"errors"
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type WhenDoubleInterface interface {
	Foo(a int) (int, error)
}

func TestWhenDoubleRet(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenDoubleInterface]()
	WhenDouble(m.Foo(Any[int]())).ThenReturn(42, nil)
	ret, _ := m.Foo(10)
	r.AssertEqual(42, ret)
}

func TestWhenDoubleAnswer(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenDoubleInterface]()
	WhenDouble(m.Foo(Any[int]())).ThenAnswer(func(args []any) (int, error) {
		return 42, nil
	})
	ret, _ := m.Foo(10)
	r.AssertEqual(42, ret)
}

func TestWhenDoubleAnswerWithArgs(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenDoubleInterface]()
	WhenDouble(m.Foo(Any[int]())).ThenAnswer(func(args []any) (int, error) {
		return args[0].(int) + 1, nil
	})
	ret, _ := m.Foo(10)
	r.AssertEqual(11, ret)
}

func TestWhenDoubleMultiAnswer(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenDoubleInterface]()
	WhenDouble(m.Foo(Any[int]())).
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

func TestWhenDoubleMultiReturn(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenDoubleInterface]()
	WhenDouble(m.Foo(Any[int]())).
		ThenReturn(10, nil).
		ThenReturn(11, nil)
	ret1, _ := m.Foo(12)
	ret2, _ := m.Foo(13)
	r.AssertEqual(10, ret1)
	r.AssertEqual(11, ret2)
}

func TestWhenDoubleAnswerAndReturn(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenDoubleInterface]()
	WhenDouble(m.Foo(Any[int]())).
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

func TestWhenDoubleReturnError(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenDoubleInterface]()
	WhenDouble(m.Foo(Any[int]())).
		ThenReturn(1, errors.New("err"))
	ret, err := m.Foo(12)
	r.AssertEqual(1, ret)
	r.AssertErrorContains(err, "err")
}

func TestWhenDoubleAnswerError(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenDoubleInterface]()
	WhenDouble(m.Foo(Any[int]())).
		ThenAnswer(func(args []any) (int, error) {
			return 0, errors.New("err")
		})
	ret, err := m.Foo(12)
	r.AssertEqual(0, ret)
	r.AssertErrorContains(err, "err")
}
