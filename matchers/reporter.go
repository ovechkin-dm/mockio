package matchers

// ErrorReporter is an interface for reporting errors during test execution.
// Implementations of this interface should provide a way to fail the test with a message.
type ErrorReporter interface {
	// Fatalf reports an error and fails the test execution.
	// It formats the message according to a format specifier and arguments
	// It can be used to report an error and provide additional context about the error.
	Fatalf(format string, args ...any)
	// Cleanup adds hooks that are used to clean up data after test was executed.
	Cleanup(func())
}
