package check

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type St struct {
}

func TestNonInterfaceNotAllowed(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	_ = Mock[St]()
	r.AssertError()
}
