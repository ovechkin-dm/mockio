package matchers

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
type Matcher[T any] interface {
	// Description returns a string describing the criteria defined by the Matcher.
	Description() string

	// Match returns true if the given method call satisfies the criteria defined by the Matcher.
	// The actual parameter represents the actual value passed to method.
	// The allArgs parameter represents all the arguments that were passed to a method.
	Match(allArgs []any, actual T) bool
}
