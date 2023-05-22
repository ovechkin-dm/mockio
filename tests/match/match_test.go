package match

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type Iface interface {
	Test(i interface{}) bool
}

type St struct {
	value int
}

func TestAny(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(Any[string]())).ThenReturn(true)
	ret := m.Test("test")
	r.AssertEqual(true, ret)
}

func TestAnyStruct(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(Any[*St]())).ThenReturn(true)
	st := &St{}
	ret := m.Test(st)
	r.AssertEqual(true, ret)
}

func TestAnyWrongType(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(Any[int]())).ThenReturn(true)
	ret := m.Test("test")
	r.AssertEqual(false, ret)
}

func TestExactStruct(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	a := St{}
	m := Mock[Iface]()
	WhenSingle(m.Test(Exact(&a))).ThenReturn(true)
	ret := m.Test(&a)
	r.AssertEqual(true, ret)
}

func TestExactWrongStruct(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	a := &St{10}
	b := &St{10}
	m := Mock[Iface]()
	WhenSingle(m.Test(Exact(a))).ThenReturn(true)
	ret := m.Test(b)
	r.AssertEqual(false, ret)
}

func TestEqualStruct(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	a := &St{10}
	b := &St{10}
	m := Mock[Iface]()
	WhenSingle(m.Test(Equal(a))).ThenReturn(true)
	ret := m.Test(b)
	r.AssertEqual(true, ret)
}

func TestNonEqualStruct(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	a := &St{11}
	b := &St{10}
	m := Mock[Iface]()
	WhenSingle(m.Test(Equal(a))).ThenReturn(true)
	ret := m.Test(b)
	r.AssertEqual(false, ret)
}

func TestCustomMatcher(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	evenm := CreateMatcher[int]("even", func(allArgs []any, actual int) bool {
		return actual%2 == 0
	})
	m := Mock[Iface]()
	WhenSingle(m.Test(Match(evenm))).ThenReturn(true)
	ret1 := m.Test(10)
	ret2 := m.Test(11)
	r.AssertEqual(ret1, true)
	r.AssertEqual(ret2, false)
}

func TestNotEqual(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(NotEqual("test"))).ThenReturn(true)
	ret := m.Test("test1")
	r.AssertEqual(true, ret)
}

func TestOneOf(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(OneOf("test1", "test2"))).ThenReturn(true)
	ret := m.Test("test2")
	r.AssertEqual(true, ret)
}
