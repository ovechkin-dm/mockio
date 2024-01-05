package matchers

import (
	"fmt"
	"reflect"
)

type MethodVerificationData struct {
	NumMethodCalls int
}

type InvocationData struct {
	MethodType reflect.Method
	MethodName string
	Args       []reflect.Value
}

type MethodVerifier interface {
	Verify(data *MethodVerificationData) error
}

func AtLeastOnce() MethodVerifier {
	return MethodVerifierFromFunc(func(data *MethodVerificationData) error {
		if data.NumMethodCalls <= 0 {
			return fmt.Errorf("expected num method calls: atLeastOnce, got: %d", data.NumMethodCalls)
		}
		return nil
	})
}

func Times(n int) MethodVerifier {
	return MethodVerifierFromFunc(func(data *MethodVerificationData) error {
		if data.NumMethodCalls != n {
			return fmt.Errorf("expected num method calls: %d, got : %d", n, data.NumMethodCalls)
		}
		return nil
	})
}

func MethodVerifierFromFunc(f func(data *MethodVerificationData) error) MethodVerifier {
	return &methodVerifierImpl{
		f: f,
	}
}

type methodVerifierImpl struct {
	f func(data *MethodVerificationData) error
}

func (m *methodVerifierImpl) Verify(data *MethodVerificationData) error {
	return m.f(data)
}
