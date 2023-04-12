package mock

import (
	"fmt"
	"github.com/ovechkin-dm/mockio/matchers"
	"github.com/ovechkin-dm/mockio/registry"
	"reflect"
)

// SetUp initializes the mock library and registers it with the global registry.
// It should be called once before using any of the mocking functions in this library.
//
// Example usage:
//
//	type MyInterface interface {
//	   MyMethod(arg1 string, arg2 int) (string, error)
//	}
//
//	func TestMyFunction(t *testing.T) {
//	   // Set up the mock library
//	   mock.SetUp(t)
//
//	   // Create a mock object that implements MyInterface
//	   myMock := mock.NewMock[MyInterface]()
//
//	   // Set up a mock behavior for the MyMethod method
//	   mock.WhenA(myMock.MyMethod("foo", 42)).ThenReturn("bar")
//
//	   // Call the method on the mock object
//	   result, err := myMock.MyMethod("foo", 42)
//
//	   // Verify that the mock was called with the correct arguments
//	   mock.Verify(myMock, mock.Times(1)).MyMethod(mock.Any[string](), mock.Any[int]())
//	}
func SetUp(t matchers.ErrorReporter) {
	registry.SetUp(t)
}

// NewMock returns a mock object that implements the specified interface or type.
// The returned object can be used to set up mock behaviors for its methods.
//
// Example usage:
//
//	type MyInterface interface {
//	   MyMethod(arg1 string, arg2 int) (string, error)
//	}
//
//	func TestMyFunction(t *testing.T) {
//	   // Set up the mock library
//	   mock.SetUp(t)
//
//	   // Create a mock object that implements MyInterface
//	   myMock := mock.NewMock[MyInterface]()
//
//	   // Set up a mock behavior for the MyMethod method
//	   mock.WhenA(myMock.MyMethod("foo", 42)).ThenReturn("bar")
//
//	   // Call the method on the mock object
//	   result, err := myMock.MyMethod("foo", 42)
//
//	   // Verify that the mock was called with the correct arguments
//	   mock.Verify(myMock, mock.Times(1)).MyMethod(mock.Any[string](), mock.Any[int]())
//	}
func NewMock[T any]() T {
	return registry.Mock[T]()
}

// Any returns a mock value of type T that matches any value of type T.
// This can be useful when setting up mock behaviors for methods that take arguments of type T,
// but the specific argument value is not important for the test case.
//
// Example usage:
//
//	// Set up a mock behavior for a method that takes a string argument
//	mock.WhenA(myMock.MyMethod(mock.Any[string]())).ThenReturn("bar")
//
//	// Set up a mock behavior for a method that takes an integer argument
//	mock.WhenA(myMock.MyOtherMethod(mock.Any[int]())).ThenReturn("baz")
func Any[T any]() T {
	registry.AddMatcher(registry.AnyMatcher[T]())
	var t T
	return t
}

// Equal returns a matcher that matches values of type T that are equal to the provided value.
// The value passed to Equal must be comparable with values of type T.
//
// Example usage:
//
//	// Set up a mock behavior for a method that takes a string argument equal to "foo"
//	mock.WhenA(myMock.MyMethod(mock.Equal("foo"))).ThenReturn("bar")
//
//	// Set up a mock behavior for a method that takes an integer argument equal to 42
//	mock.WhenA(myMock.MyOtherMethod(mock.Equal(42))).ThenReturn("baz")
func Equal[T comparable](value T) T {
	m := registry.FunMatcher("mock.Equal", func(m *matchers.MethodCall, actual any) bool {
		return value == actual
	})
	registry.AddMatcher(m)
	var t T
	return t
}

// Exact returns a matcher that matches values of type T that are equal via reflect.DeepEqual to the provided value.
// The value passed to Exact must be of the exact same type as values of type T.
//
// Example usage:
//
//	// Set up a mock behavior for a method that takes a string argument exactly equal to "foo"
//	mock.WhenA(myMock.MyMethod(mock.Exact("foo"))).ThenReturn("bar")
//
//	// Set up a mock behavior for a method that takes an integer argument exactly equal to 42
//	mock.WhenA(myMock.MyOtherMethod(mock.Exact(42))).ThenReturn("baz")
func Exact[T any](value T) T {
	m := registry.FunMatcher("mock.Exact", func(m *matchers.MethodCall, actual any) bool {
		return reflect.DeepEqual(value, actual)
	})
	registry.AddMatcher(m)
	var t T
	return t
}

// Match returns a Matcher that matches values of type T using the provided Matcher implementation.
// The provided Matcher implementation must implement the Matcher interface.
func Match[T any](m matchers.Matcher) T {
	registry.AddMatcher(m)
	var t T
	return t
}

// WhenA takes an argument of type T and returns a Returner1 interface
// that allows for specifying a return value for a method call that has that argument.
//
// Example usage:
//
//	// Set up a mock behavior for a method call that takes a string argument
//	mock.WhenA(myMock.MyMethod("some string")).ThenReturn("bar")
//
//	// Set up a mock behavior for a method call that takes a custom struct as an argument
//	myStruct := MyStruct{Field1: "value1", Field2: "value2"}
//	mock.WhenA(myMock.MyOtherMethod(mock.Exact(myStruct))).ThenReturn("baz")
func WhenA[T any](t T) matchers.Returner1[T] {
	return registry.ToReturner1[T](registry.When())
}

// WhenE takes an argument of type T and an error value and returns a ReturnerE interface
// that allows for specifying a return value and an error for a method call that has that argument.
//
// Example usage:
//
//	// Set up a mock behavior for a method call that takes a string argument and returns an error
//	mock.WhenE(myMock.MyMethod(mock.Exact("some string"))).ThenReturn("", errors.New("some other error"))
func WhenE[T any](t T, err error) matchers.ReturnerE[T] {
	return registry.ToReturnerE[T](registry.When())
}

// When sets up a method call expectation on a mocked object with a specified set of arguments
// and returns a ReturnerAll object that allows specifying the return values or answer function
// for the method call. Arguments can be any values, and the method call expectation is matched
// based on the types and values of the arguments passed. If multiple expectations match the same
// method call, the first matching expectation will be used.
//
// Args:
//
//	args: List of arguments that are expected to be passed to the method call.
//
// Returns:
//
//	A ReturnerAll object that allows specifying the return values or answer function for
//	the method call.
//
// Example Usage:
//
//	// Given an interface
//	type MyInterface interface {
//	    MyMethod(a int, b string) bool
//	}
//
//	// And a mocked implementation
//	mockMyInterface := mock.NewMock[MyInterface]()
//
//	// Set up a method call expectation
//	mockWhen := mock.When(mockMyInterface.MyMethod(mock.Any[int](), mock.Exact[String]("test"))).ThenReturn(true)
//
//	// Call the method on the mocked object
//	result := mockMyInterface.MyMethod(123, "test")
//
//	// Verify that the method was called with the expected arguments
//	mock.Verify(mockMyInterface, mockMyInterface.MyMethod(mock.Equal(123), mock.Exact("test"))).Once()
//
//	// Verify that the method was called with any int and the string "test"
//	mock.Verify(mockMyInterface, mockMyInterface.MyMethod(mock.Any[int](), mock.Exact("test"))).Once()
func When(args ...any) matchers.ReturnerAll {
	return registry.When()
}

// Captor returns an ArgumentCaptor, which can be used to capture arguments
// passed to a mocked method. ArgumentCaptor is a generic type, which means
// that the type of the arguments to be captured should be specified when
// calling Captor.
//
// Example Usage:
//
//	// Create a mock object for a Foo interface
//	mockFoo := mock.NewMock[Foo]()
//
//	// Call a method on the mock object, passing in some arguments
//	mockFoo.Bar("hello", 42)
//
//	// Create an argument captor for string arguments
//	stringCaptor := mock.Captor[string]()
//	mock.When(mockFoo.Bar(stringCaptor.Capture())
//
//	// Get the captured string value
//	capturedString := stringCaptor.Last()
//
//	// Do something with the captured string value
//	fmt.Println(capturedString)
func Captor[T any]() matchers.ArgumentCaptor[T] {
	return registry.NewArgumentCaptor[T]()
}

// Verify checks if the method call on the provided mock object matches the expected verification conditions.
//
// It takes two arguments: the mock object to be verified and a method verifier. The method verifier defines the conditions
// that should be matched during the verification process. If the verification passes, Verify returns the mock object.
// If it fails, it reports an error.
//
// The method verifier can be created using one of the following functions:
//
//   - AtLeastOnce() MethodVerifier: Matches if the method is called at least once.
//
//   - Once() MethodVerifier: Matches if the method is called exactly once.
//
//   - Times(n int) MethodVerifier: Matches if the method is called n times.
//
//   - Never() MethodVerifier: Matches if the method is never called.
//
// The Verify function is typically used to assert that a method is called with the correct arguments and/or that it is
// called the correct number of times during a unit test.
//
// Example usage:
//
//	package main
//
//	import (
//	        "testing"
//	        "github.com/ovechkin-dm/mockio/mock"
//	)
//
//	func TestMyFunction(t *testing.T) {
//	        // Create a mock object
//	        mockObj := mock.NewMock[MyObject]()
//
//	        // Call a method on the mock object
//	        mockObj.MyMethod("arg1", "arg2")
//
//	        // Verify that the MyMethod was called exactly once
//	        mock.Verify(mockObj, mock.Once()).MyMethod(mock.Match(matchers.Any[String]()), mock.Match(matchers.Exact[String]("arg2")))
//	}
func Verify[T any](t T, v matchers.MethodVerifier) T {
	registry.VerifyMethod(t, v)
	return t
}

// AtLeastOnce returns a MethodVerifier that verifies if the number of method calls
// is greater than zero. It can be used to verify that a method has been called at least once.
//
// Example usage:
//
//	mockObj := mock.NewMock[MyInterface]()
//	mockObj.MyMethod("arg1")
//	mockObj.MyMethod("arg2")
//	mock.Verify(mockObj, mock.AtLeastOnce()).MyMethod(matchers.Any[string])
//
// This verifies that the MyMethod function of mockObj was called at least once.
func AtLeastOnce() matchers.MethodVerifier {
	return matchers.AtLeastOnce()
}

// Once returns a MethodVerifier that expects a method to be called exactly once.
// If the method is not called, or called more than once, an error will be returned during verification.
func Once() matchers.MethodVerifier {
	return matchers.Times(1)
}

// Times returns a MethodVerifier that verifies the number of times a method has been called.
// It takes an integer 'n' as an argument, which specifies the expected number of method calls.
//
// Example usage:
//
//	// Create a mock object for testing
//	mockObj := mock.NewMock[MyStruct]()
//
//	// Call a method on the mock object
//	mockObj.MyMethod()
//
//	// Verify that MyMethod was called exactly once
//	mock.Verify(mockObj, mock.Times(1)).MyMethod()
//
//	// Call the method again
//	mockObj.MyMethod()
//
//	// Verify that MyMethod was called exactly twice
//	mock.Verify(mockObj, mock.Times(2)).MyMethod()
//
// If the number of method calls does not match the expected number of method calls, an error is returned.
// The error message will indicate the expected and actual number of method calls.
func Times(n int) matchers.MethodVerifier {
	return matchers.Times(n)
}

// Never returns a MethodVerifier that verifies that a method has never been called.
//
// Example usage:
//
//	// Create a mock object for testing
//	mockObj := mock.NewMock[MyInterface]()
//
//	// Verify that MyMethod was never called
//	mock.Verify(mockObj, mock.Never()).MyMethod()
//
//	// Call the method
//	mockObj.MyMethod()
//
//	// Verify that MyMethod was called at least once
//	mock.Verify(mockObj, mock.AtLeastOnce()).MyMethod()
func Never() matchers.MethodVerifier {
	return matchers.Times(0)
}

// VerifyNoMoreInteractions verifies that there are no more unverified interactions with the mock object.
//
// Example usage:
//
//	// Create a mock object for testing
//	mockObj := mock.NewMock[MyInterface]()
//
//	// Call the method
//	mockObj.MyMethod()
//
//	// Verify that MyMethod was called exactly once
//	mock.Verify(mockObj, mock.Once()).MyMethod()
//
//	// Verify that there are no more unverified interactions
//	mock.VerifyNoMoreInteractions(mockObj)
func VerifyNoMoreInteractions(value any) {
	registry.VerifyInstance(value, matchers.InstanceVerifierFromFunc(func(data *matchers.InvocationData) error {
		return fmt.Errorf("no more interactions should be recorded for mock")
	}))
}
