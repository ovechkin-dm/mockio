package matchers

// ReturnerSingle is interface that defines methods for returning a value or an answer
// for a mock function with one argument.
type ReturnerSingle[T any] interface {
	// ThenReturn sets the return value for the mock function with one argument.
	// The return value must be of type T.
	ThenReturn(value T) ReturnerSingle[T]
	// ThenAnswer sets a function that will be called when the mock function is
	// called with one argument. The function must take a variable number of
	// arguments of type interface{} and return a value of type T.
	ThenAnswer(func(args []any) T) ReturnerSingle[T]
}

// ReturnerDouble is an interface that provides methods to define the returned value and error of a mock function with a single argument.
// ThenReturn method sets the return value and error of the mocked function to the provided value and error respectively.
// ThenAnswer method sets the return value and error of the mocked function to the value and error returned by the provided function respectively.
type ReturnerDouble[A any, B any] interface {
	// ThenReturn sets the return value and error of the mocked function to the provided value and error respectively.
	ThenReturn(a A, b B) ReturnerDouble[A, B]
	// ThenAnswer sets the return value and error of the mocked function to the value and error returned by the provided function respectively.
	ThenAnswer(func(args []any) (A, B)) ReturnerDouble[A, B]
}

// ReturnerAll is a type that defines the methods for returning and answering values for
// a method call with multiple return values. It is returned by the When method.
type ReturnerAll interface {
	// ThenReturn sets the return values for the method call.
	// The number and types of the values should match the signature of the method being mocked.
	// This method can be called multiple times to set up different return values
	// for different calls to the same method with the same arguments.
	ThenReturn(values ...any) ReturnerAll

	// ThenAnswer sets a function that will be called to calculate the return values for the method call.
	// The function should have the same signature as the method being mocked.
	// This method can be called multiple times to set up different answer functions
	// for different calls to the same method with the same arguments.
	ThenAnswer(answer Answer) ReturnerAll
}
