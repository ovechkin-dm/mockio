package mocking

import (
	"errors"
	"testing"

	"github.com/ovechkin-dm/mockio/tests/common"

	. "github.com/ovechkin-dm/mockio/mock"
)

type ByteArrInterface interface {
	DoSomething(b [16]byte) string
}
type OtherIface interface {
	SomeMethod() bool
}

type CallingIface interface {
	GetMocked(appClient OtherIface) OtherIface
}

type SingleArgIface interface {
	SingleArgMethod(other OtherIface) error
}

func TestMockWithMockedArg(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)

	callingMock := Mock[CallingIface]()
	otherMock := Mock[OtherIface]()

	WhenSingle(callingMock.GetMocked(Exact[OtherIface](otherMock))).ThenReturn(otherMock)

	res := callingMock.GetMocked(otherMock)

	Verify(callingMock, Times(1)).GetMocked(Exact[OtherIface](otherMock))

	VerifyNoMoreInteractions(callingMock)

	r.AssertEqual(otherMock, res)

	r.AssertNoError()
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

func TestNilArgs(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	myMock := Mock[SingleArgIface]()
	WhenSingle(myMock.SingleArgMethod(Any[OtherIface]())).ThenReturn(errors.New("test"))
	result := myMock.SingleArgMethod(nil)
	r.AssertEqual(result.Error(), "test")
}