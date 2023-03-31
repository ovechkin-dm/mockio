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
//	// Define a mock interface
//	type MyInterface interface {
//	      MyMethod(param1 int, param2 string) (string, error)
//	}
//
//	// Create a mock object
//	mockObj := mock.Mock[MyInterface]()
//
//	// Call the method with arguments to be captured
//	param1 := 123
//	param2 := "test string"
//	mockObj.MyMethod(param1, param2)
//
//	// Create an argument captor
//	captor := mock.Captor[string]()
//
//	// Verify that the method was called with the expected arguments
//	mock.WhenE(mockObj.MyMethod(mock.Any[Int](), captor.Capture())).ThenReturn("", nil)
//	mockObj.MyMethod(1, "")
//	capturedValue := captor.Last()
//	fmt.Printf("Captured value: %v\n", capturedValue)
//
// Output:
//
//	Captured value: test string
type ArgumentCaptor[T any] interface {
	// Capture captures and stores a single argument passed to a method call.
	Capture() T
	// Last retrieves the last captured argument.
	Last() T
	// Values retrieves all captured arguments.
	Values() []T
}
