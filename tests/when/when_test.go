package when

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type WhenInterface interface {
	Foo(a int) (int, string)
	Bar(a int, b string, c string) (int, string)
	Empty() int
}

type WhenInterface2 interface {
	Bar(a int, b string, c string) (int, string)
}

type WhenStruct struct {
}

func (w *WhenStruct) foo() int {
	return 10
}

func TestWhenRet(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenInterface]()
	When(m.Foo(Any[int]())).ThenReturn(42, "test")
	i, s := m.Foo(10)
	r.AssertEqual(42, i)
	r.AssertEqual("test", s)
}

func TestEmptyWhenErr(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	ws := &WhenStruct{}
	When(ws.foo())
	r.AssertError()
}

func TestIncorrectNumMatchers(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenInterface]()
	When(m.Bar(10, Any[string](), Any[string]()))
	r.AssertError()
}

func TestIncorrectMatchersReuse(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenInterface]()
	anyS := Any[string]()
	When(m.Bar(10, anyS, anyS))
	r.AssertError()
}

func TestNoMatchersAreExactOnReturn(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenInterface]()
	When(m.Bar(10, "test1", "test2")).ThenReturn(10, "2")
	r.AssertNoError()
	i, s := m.Bar(10, "test1", "test2")
	r.AssertEqual(10, i)
	r.AssertEqual("2", s)
}

func TestNoMatchersAreExactOnAnswer(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenInterface]()
	When(m.Bar(10, "test1", "test2")).ThenAnswer(func(args []any) []any {
		return []any{args[0].(int) + 1, "2"}
	})
	r.AssertNoError()
	i, s := m.Bar(10, "test1", "test2")
	r.AssertEqual(11, i)
	r.AssertEqual("2", s)
}

func TestEmptyArgs(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[WhenInterface]()
	When(m.Empty()).ThenReturn(10)
	ret := m.Empty()
	r.AssertEqual(10, ret)
}

func TestWhenMultipleIfaces(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m1 := Mock[WhenInterface]()
	m2 := Mock[WhenInterface2]()
	When(m1.Bar(10, "test", "test")).ThenReturn(10, "test")
	When(m2.Bar(10, "test", "test")).ThenReturn(11, "test1")
	i1, s1 := m1.Bar(10, "test", "test")
	i2, s2 := m2.Bar(10, "test", "test")
	r.AssertEqual(10, i1)
	r.AssertEqual("test", s1)
	r.AssertEqual(11, i2)
	r.AssertEqual("test1", s2)
	r.AssertNoError()
}
