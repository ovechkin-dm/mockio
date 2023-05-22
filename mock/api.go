package mock

import (
	"fmt"
	"github.com/ovechkin-dm/mockio/matchers"
	"github.com/ovechkin-dm/mockio/registry"
	"reflect"
	"strings"
)

// SetUp initializes the mock library with the reporter.
// Example usage:
//
//	package simple
//
//	import (
//		. "github.com/ovechkin-dm/mockio/mock"
//		"testing"
//	)
//
//	type myInterface interface {
//		Foo(a int) int
//	}
//
//	func TestSimple(t *testing.T) {
//		SetUp(t)
//		m := Mock[myInterface]()
//		WhenSingle(m.Foo(Any[int]())).ThenReturn(42)
//		ret := m.Foo(10)
//		r.AssertEqual(42, ret)
//	}
func SetUp(t matchers.ErrorReporter) {
	registry.SetUp(t)
}

// Mock returns a mock object that implements the specified interface or type.
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
//	   SetUp(t)
//
//	   // Create a mock object that implements MyInterface
//	   myMock := Mock[MyInterface]()
//
//	   // Set up a mock behavior for the MyMethod method
//	   WhenSingle(myMock.MyMethod("foo", 42)).ThenReturn("bar")
//
//	   // Call the method on the mock object
//	   result, err := myMock.MyMethod("foo", 42)
//
//	   // Verify that the mock was called with the correct arguments
//	   Verify(myMock, Times(1)).MyMethod(Any[string](), Any[int]())
//	}
func Mock[T any]() T {
	return registry.Mock[T]()
}

// Any returns a mock value of type T that matches any value of type T.
// This can be useful when setting up mock behaviors for methods that take arguments of type T,
// but the specific argument value is not important for the test case.
//
// Example usage:
//
//	// Set up a mock behavior for a method that takes a string argument
//	WhenSingle(myMock.MyMethod(mock.Any[string]())).ThenReturn("bar")
//
//	// Set up a mock behavior for a method that takes an integer argument
//	WhenSingle(myMock.MyOtherMethod(mock.Any[int]())).ThenReturn("baz")
func Any[T any]() T {
	registry.AddMatcher(registry.AnyMatcher[T]())
	var t T
	return t
}

// AnyInt is an alias for Any[int]
// See Any for more description
func AnyInt() int {
	return Any[int]()
}

// AnyString is an alias for Any[string]
// See Any for more description
func AnyString() string {
	return Any[string]()
}

// AnyInterface is an alias for Any[any]
// See Any for more description
func AnyInterface() any {
	return Any[any]()
}

// AnyOfType is an alias for Any[T] for specific type
// Used for automatic type inference
func AnyOfType[T any](t T) T {
	return Any[T]()
}

// Exact returns a matcher that matches values of type T that are equal to the provided value.
// The value passed to Exact must be comparable with values of type T.
//
// Example usage:
//
//	// Set up a mock behavior for a method that takes a string argument equal to "foo"
//	WhenSingle(myMock.MyMethod(Exact("foo"))).ThenReturn("bar")
//
//	// Set up a mock behavior for a method that takes an integer argument equal to 42
//	WhenSingle(myMock.MyOtherMethod(Exact(42))).ThenReturn("baz")
func Exact[T comparable](value T) T {
	desc := fmt.Sprintf("Exact(%v)", value)
	m := registry.FunMatcher(desc, func(m []any, actual any) bool {
		return value == actual
	})
	registry.AddMatcher(m)
	var t T
	return t
}

// Equal returns a matcher that matches values of type T that are equal via reflect.DeepEqual to the provided value.
// The value passed to Equal must be of the exact same type as values of type T.
//
// Example usage:
//
//	// Set up a mock behavior for a method that takes a string argument exactly equal to "foo"
//	WhenSingle(myMock.MyMethod(Equal("foo"))).ThenReturn("bar")
//
//	// Set up a mock behavior for a method that takes an integer argument exactly equal to 42
//	WhenSingle(myMock.MyOtherMethod(Equal(42))).ThenReturn("baz")
func Equal[T any](value T) T {
	desc := fmt.Sprintf("Equal(%v)", value)
	m := registry.FunMatcher(desc, func(m []any, actual any) bool {
		return reflect.DeepEqual(value, actual)
	})
	registry.AddMatcher(m)
	var t T
	return t
}

// NotEqual returns a matcher that matches values of type T that are not equal via reflect.DeepEqual to the provided value.
// The value passed to NotEqual must be of the exact same type as values of type T.
//
// Example usage:
//
//	// Set up a mock behavior for a method that takes a string argument not equal to "foo"
//	WhenSingle(myMock.MyMethod(NotEqual("foo"))).ThenReturn("bar")
//
//	// Set up a mock behavior for a method that takes an integer argument not equal to 42
//	WhenSingle(myMock.MyOtherMethod(NotEqual(42))).ThenReturn("baz")
func NotEqual[T any](value T) T {
	desc := fmt.Sprintf("NotEqual(%v)", value)
	m := registry.FunMatcher(desc, func(m []any, actual any) bool {
		return !reflect.DeepEqual(value, actual)
	})
	registry.AddMatcher(m)
	var t T
	return t
}

// OneOf returns a matcher that matches at least one of values of type T that are equal via reflect.DeepEqual to the provided value.
// The value passed to OneOf must be of the exact same type as values of type T.
//
// Example usage:
//
//	// Set up a mock behavior for a method that takes a string argument equal to either "foo" or "bar"
//	WhenSingle(myMock.MyMethod(OneOf("foo", "bar"))).ThenReturn("bar")
//
//	// Set up a mock behavior for a method that takes an integer argument equal to either 41 or 42
//	WhenSingle(myMock.MyOtherMethod(OneOf(41, 42))).ThenReturn("baz")
func OneOf[T any](values ...T) T {
	vs := make([]string, len(values))
	for i := range values {
		vs[i] = fmt.Sprintf("%v", values[i])
	}

	desc := fmt.Sprintf("OneOf(%s)", strings.Join(vs, ","))
	m := registry.FunMatcher[T](desc, func(args []any, t T) bool {
		for i := range values {
			if reflect.DeepEqual(values[i], t) {
				return true
			}
		}
		return false
	})
	registry.AddMatcher(m)
	var t T
	return t
}

// CreateMatcher returns a Matcher that matches values of type T using the provided Matcher implementation.
// The provided Matcher implementation must implement the Matcher interface.
func CreateMatcher[T any](description string, f func(allArgs []any, actual T) bool) matchers.Matcher[T] {
	m := registry.FunMatcher(description, f)
	return m
}

// Match provides matching for method argument with a matcher that was created via CreateMatcher
// The provided Matcher implementation must implement the Matcher interface.
func Match[T any](m matchers.Matcher[T]) T {
	registry.AddMatcher(m)
	var t T
	return t
}

// WhenSingle takes an argument of type T and returns a ReturnerSingle interface
// that allows for specifying a return value for a method call that has that argument.
// This function should be used for method that returns exactly one return value
// It acts like When, but also provides additional type check on return value
// For more than on value consider using WhenDouble or When
func WhenSingle[T any](t T) matchers.ReturnerSingle[T] {
	return registry.ToReturnerSingle[T](registry.When())
}

// WhenDouble takes an arguments of type A and B and  returns a ReturnerDouble interface
// that allows for specifying two return values for a method call that has that argument.
// This function should be used for method that returns exactly two return values
// It acts like When, but also provides additional type check on return values
// For more multiple return values consider using When
func WhenDouble[A any, B any](a A, b B) matchers.ReturnerDouble[A, B] {
	return registry.ToReturnerDouble[A, B](registry.When())
}

// When sets up a method call expectation on a mocked object with a specified set of arguments
// and returns a ReturnerAll object that allows specifying the return values or answer function
// for the method call. Arguments can be any values, and the method call expectation is matched
// based on the types and values of the arguments passed. If multiple expectations match the same
// method call, the first matching expectation will be used.
func When(args ...any) matchers.ReturnerAll {
	return registry.When()
}

// Captor returns an ArgumentCaptor, which can be used to capture arguments
// passed to a mocked method. ArgumentCaptor is a generic type, which means
// that the type of the arguments to be captured should be specified when
// calling Captor.
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
//	package simple
//
//	import (
//		. "github.com/ovechkin-dm/mockio/mock"
//		"testing"
//	)
//
//	type myInterface interface {
//		Foo(a int) int
//	}
//
//	func TestSimple(t *testing.T) {
//		r := common.NewMockReporter(t)
//		SetUp(r)
//		m := Mock[myInterface]()
//		WhenSingle(m.Foo(Any[int]())).ThenReturn(42)
//		ret := m.Foo(10)
//		r.AssertEqual(42, ret)
//		Verify(m, AtLeastOnce()).Foo(10)
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
//	mockObj := Mock[MyInterface]()
//	mockObj.MyMethod("arg1")
//	mockObj.MyMethod("arg2")
//	Verify(mockObj, AtLeastOnce()).MyMethod(Any[string]())
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
//	mockObj := Mock[MyInterface]()
//
//	// Call a method on the mock object
//	mockObj.MyMethod()
//
//	// Verify that MyMethod was called exactly once
//	Verify(mockObj, Times(1)).MyMethod()
//
//	// Call the method again
//	mockObj.MyMethod()
//
//	// Verify that MyMethod was called exactly twice
//	Verify(mockObj, Times(2)).MyMethod()
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
//	mockObj := Mock[MyInterface]()
//
//	// Verify that MyMethod was never called
//	Verify(mockObj, Never()).MyMethod()
//
//	// Call the method
//	mockObj.MyMethod()
//
//	// Verify that MyMethod was called at least once
//	Verify(mockObj, AtLeastOnce()).MyMethod()
func Never() matchers.MethodVerifier {
	return matchers.Times(0)
}

// VerifyNoMoreInteractions verifies that there are no more unverified interactions with the mock object.
//
// Example usage:
//
//	// Create a mock object for testing
//	mockObj := Mock[MyInterface]()
//
//	// Call the method
//	mockObj.MyMethod()
//
//	// Verify that MyMethod was called exactly once
//	Verify(mockObj, Once()).MyMethod()
//
//	// Verify that there are no more unverified interactions
//	VerifyNoMoreInteractions(mockObj)
func VerifyNoMoreInteractions(value any) {
	registry.VerifyInstance(value, matchers.InstanceVerifierFromFunc(func(data *matchers.InvocationData) error {
		return fmt.Errorf("no more interactions should be recorded for mock")
	}))
}
