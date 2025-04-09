package simple

import (
	"testing"

	"github.com/ovechkin-dm/mockio/v2/tests/common"

	. "github.com/ovechkin-dm/mockio/v2/mock"
)

type myInterface interface {
	Foo(a int) int
}

func TestSimple(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	m := Mock[myInterface](ctrl)
	WhenSingle(m.Foo(Any[int]())).ThenReturn(42)
	ret := m.Foo(10)
	r.AssertEqual(42, ret)
	Verify(m, AtLeastOnce()).Foo(10)
}
