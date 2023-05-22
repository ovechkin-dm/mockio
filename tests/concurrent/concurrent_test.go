package concurrent

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"sync"
	"testing"
)

type myInterface interface {
	Foo(a int) int
}

func TestNewMockInOtherFiber(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[myInterface]()
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
