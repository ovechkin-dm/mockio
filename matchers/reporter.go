package matchers

import "fmt"

// ErrorReporter is an interface for reporting errors during test execution.
// Implementations of this interface should provide a way to fail the test with a message.
type ErrorReporter interface {
	// Fatalf reports an error and fails the test execution.
	// It formats the message according to a format specifier and arguments,
	// then passes it to the panic function.
	// The panic function should be intercepted by the testing framework to fail the test.
	// It can be used to report an error and provide additional context about the error.
	Fatalf(format string, args ...any)
}

type ConsoleReporter struct {
}

func (c *ConsoleReporter) Fatalf(format string, args ...any) {
	panic(fmt.Sprintf(format+"\n", args...))
}
