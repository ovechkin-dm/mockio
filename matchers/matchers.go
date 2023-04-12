package matchers

import (
	"github.com/ovechkin-dm/go-dyno/pkg/dyno"
	"reflect"
)

// MethodCall represents a recorded method call with its unique identifier, method object, and arguments.
type MethodCall struct {
	// ID is the unique identifier for the recorded method call.
	ID string
	// Method is a pointer to the dyno.Method object representing the method being called.
	Method *dyno.Method
	// Values is a slice containing the argument values passed to the method.
	Values []reflect.Value
}

// Answer is a type alias for a function that can be used as a return value for mock function calls.
// This function takes a variable number of interface{} arguments and returns a slice of interface{} values.
// Each value in the returned slice corresponds to a return value for the mock function call.
// This type can be used to provide dynamic return values based on the input arguments passed to the mock function call.
type Answer = func(args []any) []any

// Matcher interface represents an object capable of matching method calls to specific criteria.
//
// A Matcher should implement the Match method, which takes a MethodCall and an actual parameter, and returns true
// if the parameter satisfies the criteria defined by the Matcher.
//
// A Matcher should also implement the Description method, which returns a string describing the criteria defined by
// the Matcher.
//
// Matchers can be used in conjunction with the Match function to create flexible and powerful method call matching
// capabilities.
type Matcher interface {
	// Description returns a string describing the criteria defined by the Matcher.
	Description() string

	// Match returns true if the given method call satisfies the criteria defined by the Matcher.
	// The actual parameter represents the expected value or type, depending on the Matcher implementation.
	Match(methodCall *MethodCall, actual any) bool
}
