package main

import (
	"fmt"
	"github.com/ovechkin-dm/mockio/matchers"
	. "github.com/ovechkin-dm/mockio/mock"
)

type User struct {
	Id   int
	Name string
	Age  int
}

type Storage interface {
	GetUser(id int) (*User, error)
}

type Service struct {
	storage Storage
}

func (s *Service) IsAdult(userId int) (bool, error) {
	user, err := s.storage.GetUser(userId)
	if err != nil {
		return false, err
	}
	return user.Age >= 18, nil
}

func main() {
	SetUp(&matchers.ConsoleReporter{})
	storage := NewMock[Storage]()
	service := &Service{
		storage: storage,
	}
	u := &User{
		Id:   10,
		Name: "Tony",
		Age:  19,
	}
	WhenE(storage.GetUser(Any[int]())).ThenReturn(u, nil)
	ad, err := service.IsAdult(10)
	Verify(storage, Never()).GetUser(Exact(11))
	VerifyNoMoreInteractions(storage)
	if err != nil {
		panic(err)
	}
	fmt.Println(ad)
}
