package mocking

import (
	"testing"

	"github.com/ovechkin-dm/mockio/tests/common"

	. "github.com/ovechkin-dm/mockio/mock"
)

type ByteArrInterface interface {
	DoSomething(b [16]byte) string
}

func TestByteArrayArgs(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	myMock := Mock[ByteArrInterface]()
	myBytes := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	WhenSingle(myMock.DoSomething(myBytes)).ThenReturn("test")
	result := myMock.DoSomething(myBytes)
	r.AssertEqual(result, "test")
}
