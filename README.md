# Mockio - A Golang Mocking Library

Mockio is a Golang library that provides functionality for mocking and stubbing functions and methods in tests inspired by mockito. The library is designed to simplify the testing process by allowing developers to easily create test doubles for their code, which can then be used to simulate different scenarios.

## Installation

To use Mockio in your Golang project, you can install it using the go get command:

```bash
go get github.com/ovechkin-dm/mockio
```

## Getting Started
To use Mockio, you first need to import the `mockio` package into your test file:
```go
import (
    "github.com/ovechkin-dm/mockio/mock"
)
```

You can then start creating mock objects using the mock.Mock[T]() function. This function returns a mock object of type T, which you can then use to stub methods or functions.
```
// Create a mock object for interface
mockObject := mock.Mock[MyInterface]()
```


## Stubs and Matchers
Once you have a mock object, you can start stubbing methods or functions using the `WhenA`, `WhenE`, and `When` functions. These functions allow you to define a set of conditions under which the stubbed method or function will return a certain value.

```go
// Stub a method to always return a specific value
mockObject := mock.Mock[MyInterface]()
mock.WhenA(mockObject.MethodCall(mock.Exact[Int](1))).ThenReturn("value")
```

Mockio also provides a set of matchers that you can use to define more complex conditions for your stubbed methods.
```go
// Use a matcher to stub a method with a specific input
mockObject := mock.Mock[MyInterface]()
mock.WhenA(mockObject.MethodCall(matchers.Equal("input"))).ThenReturn("value")
```

## Verification
Mockio also provides functionality for verifying that certain methods were called with specific inputs. You can use the `Verify` function to verify that a specific method or function was called a certain number of times, or with a specific set of inputs.
```go
// Verify that a method was called exactly once
mockObject := mock.Mock[MyInterface]()
mockObject.MethodCall("input")
mock.Verify(mockObject, mock.Once()).MethodCall(matchers.Equal("input"))
```

## Argument Captors
Mockio also provides functionality for capturing the arguments passed to a mocked method or function. You can use the Captor function to create an argument captor, which you can then use to retrieve the captured arguments.
```go
// Use an argument captor to capture the argument passed to a function
mockObject := mock.Mock[MyInterface]()
argumentCaptor := mockio.Captor[int]()
mockio.WhenA(mockObject.MethodCall(argumentCaptor.Capture())).ThenReturn("value")
mockFunction(42)
capturedArgument := argumentCaptor.Last() // 42
```

## Reporting
The `ErrorReporter` interface defines how errors should be reported in the library. It has a single method `Fatalf` which takes a format string and its arguments and panics with a formatted error message. The `ConsoleReporter` is a concrete implementation of this interface that reports the error to the console.

Returner
The `Returner` interfaces define how the mocked function should return values. There are three different `Returner` interfaces:

- `Returner1[T any]` for functions with one return value
- `ReturnerE[T any]` for functions with two return values, where the second one is an error
- `ReturnerAll` for functions that return multiple values
Each of these interfaces provides methods ThenReturn and ThenAnswer for setting the mocked function's return value(s). ThenReturn takes one or more values and sets them as the return value(s) of the mocked function. ThenAnswer takes a function that returns the value(s) to be used as the return value(s) of the mocked function.

## Call verifiers

The `MethodVerifier` interface defines how method calls should be verified. There are four concrete implementations of this interface:

`AtLeastOnce()` verifies that the mocked method was called at least once
`Once()` verifies that the mocked method was called exactly once
`Times(n int)` verifies that the mocked method was called n times
`Never()` verifies that the mocked method was never called
The `Verify` method of the MethodVerifier interface takes a MethodVerificationData object which contains information about the method call, such as the number of times it was called. If the verification fails, an error is returned.

The `InstanceVerifier` interface defines how instances should be verified. It has a single method RecordInteraction which takes an InvocationData object containing information about the method call. If the verification fails, error is being reported.

## Conclusion
The mock package provides a powerful library for creating and managing mock objects in Go. With its support for capturing arguments, matching arguments, and verifying method calls, it makes it easy to test complex systems with many dependencies. Its well-defined interfaces and clear documentation make it easy to use and extend, and its support for multiple return values and errors makes it suitable for a wide range of use cases.

