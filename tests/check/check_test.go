package check

import (
	"testing"

	"github.com/ovechkin-dm/mockio/tests/common"

	. "github.com/ovechkin-dm/mockio/mock"
)

type St struct{}

func TestNonInterfaceNotAllowed(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	_ = Mock[St]()
	r.AssertError()
}
