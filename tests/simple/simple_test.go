package simple

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type myInterface interface {
	Foo(a int) int
}

func TestSimple(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[myInterface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(42)
	ret := m.Foo(10)
	r.AssertEqual(42, ret)
	Verify(m, AtLeastOnce()).Foo(10)
}
