package check

import (
	"testing"

	"github.com/ovechkin-dm/mockio/v2/tests/common"

	. "github.com/ovechkin-dm/mockio/v2/mock"
)

type St struct{}

func TestNonInterfaceNotAllowed(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but code did not panic")
		}
	}()
	_ = Mock[St](ctrl)
	r.AssertError()
}
