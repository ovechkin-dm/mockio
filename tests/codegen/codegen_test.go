package codegen

import (
	"testing"

	. "github.com/ovechkin-dm/mockio/v2/mock"
	"github.com/ovechkin-dm/mockio/v2/tests/common"
)

func TestGeneratedMock(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	m := NewMockGreeter(ctrl)
	WhenSingle(m.Greet(AnyString())).ThenAnswer(func (args []any) string  {
		return "Hello, " + args[0].(string)
	})
	result := m.Greet("John")
	r.AssertEqual("Hello, John", result)
	Verify(m, Once()).Greet("John")	
}

func TestGeneratedMockNoMoreInteractions(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	m := NewMockGreeter(ctrl)
	WhenSingle(m.Greet(AnyString())).ThenAnswer(func (args []any) string  {
		return "Hello, " + args[0].(string)
	})
	result := m.Greet("John")
	r.AssertEqual("Hello, John", result)
	VerifyNoMoreInteractions(m)
	
}
