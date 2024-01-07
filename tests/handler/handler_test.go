package handler

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/tests/common"
	"testing"
)

type APIServiceA interface {
	GetAll() ([]string, error)
	Delete(ID string) error
	Add(ID string) error
}

type APIServiceX interface {
	SendMessage(message string) error
	IncrementCounterBy(incr int) error
}

type OtherService struct {
	serviceA APIServiceA
	serviceX APIServiceX
}

func (s OtherService) AddAllFromServiceA(strings []string) error {
	for _, str := range strings {
		err := s.serviceA.Add(str)
		if err != nil {
			return err
		}
	}
	s.serviceX.IncrementCounterBy(1)
	return nil
}

func TestGetAllFromServiceAUsingMockio(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	mockServiceA := Mock[APIServiceA]()

	When(mockServiceA.Add(AnyString()))

	mockServiceX := Mock[APIServiceX]()
	mockServiceX.SendMessage("TEST2")

	s := OtherService{
		serviceA: mockServiceA,
		serviceX: mockServiceX,
	}

	_ = s.AddAllFromServiceA([]string{"TEST", "OTHER"})

	_ = Verify(mockServiceA, Once()).Add("TEST")
	_ = Verify(mockServiceA, Once()).Add("OTHER")
	_ = Verify(mockServiceX, Once()).SendMessage("TEST2")
	Verify(mockServiceX, Once()).IncrementCounterBy(1)
	VerifyNoMoreInteractions(mockServiceA)
	VerifyNoMoreInteractions(mockServiceX)
	r.AssertNoError()
}
