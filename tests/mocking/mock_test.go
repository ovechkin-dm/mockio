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

type MultiMethod interface {
	One(int) int
	Two(int) int
	Three(int) int
	Four(int) int
}

type ParentIface interface {
	Foo(int) int
}

type ChildIface interface {
	ParentIface
	Bar(int) int
}

type PrivateIface interface {
	privateMethod() bool
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

func TestMultiMethodOrder(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	myMock := Mock[MultiMethod]()
	WhenSingle(myMock.One(1)).ThenReturn(1)
	WhenSingle(myMock.Two(2)).ThenReturn(2)
	WhenSingle(myMock.Three(3)).ThenReturn(3)
	WhenSingle(myMock.Four(4)).ThenReturn(4)
	r.AssertEqual(myMock.One(1), 1)
	r.AssertEqual(myMock.Two(2), 2)
	r.AssertEqual(myMock.Three(3), 3)
	r.AssertEqual(myMock.Four(4), 4)
}

func TestMockSimpleCasting(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	myMock := Mock[OtherIface]()
	WhenSingle(myMock.SomeMethod()).ThenReturn(true)
	var casted any = myMock
	source := casted.(OtherIface)
	result := source.SomeMethod()
	r.AssertEqual(result, true)
}

func TestMockCasting(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	myMock := Mock[ChildIface]()
	WhenSingle(myMock.Foo(1)).ThenReturn(1)
	WhenSingle(myMock.Bar(1)).ThenReturn(2)
	var casted any = myMock
	source := casted.(ParentIface)
	result := source.Foo(1)
	r.AssertEqual(result, 1)
}

func TestMockPrivate(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	myMock := Mock[PrivateIface]()
	WhenSingle(myMock.privateMethod()).ThenReturn(true)
	var casted any = myMock.(PrivateIface)
	source := casted.(PrivateIface)
	source.privateMethod()
	r.AssertNoError()
}