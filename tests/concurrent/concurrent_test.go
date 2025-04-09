package concurrent

import (
	"sync"
	"testing"

	"github.com/ovechkin-dm/mockio/v2/tests/common"

	. "github.com/ovechkin-dm/mockio/v2/mock"
)

type myInterface interface {
	Foo(a int) int
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
