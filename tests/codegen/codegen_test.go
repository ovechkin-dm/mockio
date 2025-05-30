package codegen

import (
	"testing"

	"github.com/ovechkin-dm/mockio/v2/tests/common"

	. "github.com/ovechkin-dm/mockio/v2/mock"
)

func TestGeneratedMock(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	m := NewMockGreeter(ctrl)
	WhenSingle(m.Greet(AnyString())).ThenAnswer(func(args []any) string {
		return "Hello, " + args[0].(string)
	})
	result := m.Greet("John")
	r.AssertEqual("Hello, John", result)
	Verify(m, Once()).Greet("John")
	r.AssertNoError()
}

func TestGeneratedMockNoMoreInteractionsSuccess(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	m := NewMockGreeter(ctrl)
	WhenSingle(m.Greet(AnyString())).ThenAnswer(func(args []any) string {
		return "Hello, " + args[0].(string)
	})
	result := m.Greet("John")
	r.AssertEqual("Hello, John", result)
	Verify(m, Once()).Greet("John")
	VerifyNoMoreInteractions(m)
	r.AssertNoError()
}

func TestGeneratedMockNoMoreInteractionsFail(t *testing.T) {
	r := common.NewMockReporter(t)
	ctrl := NewMockController(r)
	m := NewMockGreeter(ctrl)
	WhenSingle(m.Greet(AnyString())).ThenAnswer(func(args []any) string {
		return "Hello, " + args[0].(string)
	})
	result := m.Greet("John")
	r.AssertEqual("Hello, John", result)
	VerifyNoMoreInteractions(m)
	r.AssertError()
}
