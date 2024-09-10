package match

import (
	"testing"

	"github.com/ovechkin-dm/mockio/tests/common"

	. "github.com/ovechkin-dm/mockio/mock"
)

type Iface interface {
	Test(i interface{}) bool
}

type Greeter interface {
	Greet(name any) string
}

type St struct {
	value int
}

type MyStruct struct {
	items any
}

type MyInterface interface {
	Test(m *MyStruct) int
}

type SliceInterface interface {
	Test(m []int) int
}

type MapInterface interface {
	Test(m map[int]int) int
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

func TestDeepEqual(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[MyInterface]()
	s1 := MyStruct{
		items: &[]int{1, 2, 3},
	}
	s2 := MyStruct{
		items: &[]int{1, 2, 3},
	}
	WhenSingle(m.Test(&s1)).ThenReturn(9)
	result := m.Test(&s2)
	r.AssertEqual(result, 9)
}

func TestNilMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(Nil[any]())).ThenReturn(true)
	ret := m.Test(nil)
	r.AssertEqual(true, ret)
}

func TestNilNoMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(Nil[any]())).ThenReturn(true)
	ret := m.Test(10)
	r.AssertEqual(false, ret)
}

func TestSubstringMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(Substring("test"))).ThenReturn(true)
	ret := m.Test("123test123")
	r.AssertEqual(true, ret)
}

func TestSubstringNoMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(Substring("test321"))).ThenReturn(true)
	ret := m.Test("123test123")
	r.AssertEqual(false, ret)
}

func TestNotNilMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(NotNil[any]())).ThenReturn(true)
	ret := m.Test(10)
	r.AssertEqual(true, ret)
}

func TestNotNilNoMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(NotNil[any]())).ThenReturn(true)
	ret := m.Test(nil)
	r.AssertEqual(false, ret)
}

func TestRegexMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(Regex("test"))).ThenReturn(true)
	ret := m.Test("123test123")
	r.AssertEqual(true, ret)
}

func TestRegexNoMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	WhenSingle(m.Test(Regex("test321"))).ThenReturn(true)
	ret := m.Test("123test123")
	r.AssertEqual(false, ret)
}

func TestSliceLenMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[SliceInterface]()
	WhenSingle(m.Test(SliceLen[int](3))).ThenReturn(3)
	ret := m.Test([]int{1, 2, 3})
	r.AssertEqual(3, ret)
}

func TestSliceLenNoMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[SliceInterface]()
	WhenSingle(m.Test(SliceLen[int](4))).ThenReturn(3)
	ret := m.Test([]int{1, 2, 3})
	r.AssertEqual(0, ret)
}

func TestMapLenMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[MapInterface]()
	WhenSingle(m.Test(MapLen[int, int](3))).ThenReturn(3)
	ret := m.Test(map[int]int{1: 1, 2: 2, 3: 3})
	r.AssertEqual(3, ret)
}

func TestMapLenNoMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[MapInterface]()
	WhenSingle(m.Test(MapLen[int, int](3))).ThenReturn(3)
	ret := m.Test(map[int]int{1: 1, 2: 2, 3: 3, 4: 4})
	r.AssertEqual(0, ret)
}

func TestSliceContainsMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[SliceInterface]()
	WhenSingle(m.Test(SliceContains[int](3))).ThenReturn(3)
	ret := m.Test([]int{1, 2, 3})
	r.AssertEqual(3, ret)
}

func TestSliceContainsNoMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[SliceInterface]()
	WhenSingle(m.Test(SliceContains[int](4))).ThenReturn(3)
	ret := m.Test([]int{1, 2, 3})
	r.AssertEqual(0, ret)
}

func TestMapContainsMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[MapInterface]()
	WhenSingle(m.Test(MapContains[int, int](3))).ThenReturn(3)
	ret := m.Test(map[int]int{1: 1, 2: 2, 3: 3})
	r.AssertEqual(3, ret)
}

func TestMapContainsNoMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[MapInterface]()
	WhenSingle(m.Test(MapContains[int, int](4))).ThenReturn(3)
	ret := m.Test(map[int]int{1: 1, 2: 2, 3: 3})
	r.AssertEqual(0, ret)
}

func TestSliceEqualUnorderedMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[SliceInterface]()
	WhenSingle(m.Test(SliceEqualUnordered[int]([]int{1, 2, 3}))).ThenReturn(3)
	ret := m.Test([]int{3, 2, 1})
	r.AssertEqual(3, ret)
}

func TestSliceEqualUnorderedNoMatch(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[SliceInterface]()
	WhenSingle(m.Test(SliceEqualUnordered[int]([]int{1, 2, 3}))).ThenReturn(3)
	ret := m.Test([]int{3, 2, 1, 4})
	r.AssertEqual(0, ret)
}

func TestUnexpectedUseOfMatchers(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[Iface]()
	m.Test(AnyString())
	Verify(m, Once()).Test("test")
	r.AssertErrorContains(r.GetError(), "Unexpected matchers declaration")
}

func TestExactNotComparable(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	var data any = []int{1, 2}
	When(greeter.Greet(Exact(data))).ThenReturn("hello world")
	greeter.Greet(data)
}
