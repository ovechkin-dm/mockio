package matchers

// ArgumentCaptor is interface that allows capturing arguments
// passed to a mock method call.
//
// The Capture method is used to capture and store a single argument passed to a method call.
// The Last method is used to retrieve the last captured argument.
// The Values method is used to retrieve all captured arguments.
//
// Example usage:
//
// // Create a mock object
//
//	m := Mock[Iface]()
//
//	// Create captor for int value
//	c := Captor[int]()
//
//	// Use captor.Capture() inside When expression
//	WhenSingle(m.Foo(AnyInt(), c.Capture())).ThenReturn(10)
//
//	m.Foo(10, 20)
//	capturedValue := c.Last()
//
//	fmt.Printf("Captured value: %v\n", capturedValue)
//
// Output:
//
//	Captured value: 20
type ArgumentCaptor[T any] interface {
	// Capture captures and stores a single argument passed to a method call.
	Capture() T
	// Last retrieves the last captured argument.
	Last() T
	// Values retrieves all captured arguments.
	Values() []T
}
