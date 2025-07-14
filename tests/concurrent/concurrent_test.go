package concurrent

import (
	"sync"
	"testing"

	"github.com/ovechkin-dm/mockio/v2/tests/common"

	. "github.com/ovechkin-dm/mockio/v2/mock"
)

type myInterface interface {
	Foo(a int) int
	Bar(a int) int
}

func TestNewMockInOtherFiber(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	m := Mock[myInterface](ctrl)
	wg := sync.WaitGroup{}
	wg.Add(1)
	WhenSingle(m.Foo(Any[int]())).ThenReturn(42)
	ans := 0
	go func() {
		ans = m.Foo(10)
		wg.Done()
	}()
	wg.Wait()

	r.AssertEqual(42, ans)
	r.AssertNoError()
}

func TestRecursiveCall(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	m := Mock[myInterface](ctrl)
	WhenSingle(m.Foo(Any[int]())).ThenAnswer(func(args []any) int {
		return m.Bar(args[0].(int))
	})
	WhenSingle(m.Bar(Any[int]())).ThenAnswer(func(args []any) int {
		return args[0].(int) + 1
	})
	ans := m.Foo(10)
	r.AssertEqual(11, ans)
	r.AssertNoError()
}
