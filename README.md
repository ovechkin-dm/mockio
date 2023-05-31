# Mockio 

[![Build Status](https://github.com/ovechkin-dm/mockio/actions/workflows/build.yml/badge.svg)](https://github.com/ovechkin-dm/mockio/actions)
[![Codecov](https://codecov.io/gh/ovechkin-dm/mockio/branch/main/graph/badge.svg)](https://app.codecov.io/gh/ovechkin-dm/mockio)
[![Go Report Card](https://goreportcard.com/badge/github.com/ovechkin-dm/mockio)](https://goreportcard.com/report/github.com/ovechkin-dm/mockio)
[![Documentation](https://pkg.go.dev/badge/github.com/ovechkin-dm/mockio.svg)](https://pkg.go.dev/github.com/ovechkin-dm/mockio)
[![Release](https://img.shields.io/github/release/ovechkin-dm/mockio.svg)](https://github.com/ovechkin-dm/mockio/releases)
[![License](https://img.shields.io/github/license/ovechkin-dm/mockio.svg)](https://github.com/ovechkin-dm/mockio/blob/main/LICENSE)

# Mock library for golang without code generation
Mockio is a Golang library that provides functionality for mocking and stubbing functions and methods in tests inspired by mockito. The library is designed to simplify the testing process by allowing developers to easily create test doubles for their code, which can then be used to simulate different scenarios.

# Features
* No code generation required, mocks are created at runtime
* Simple and easy to use API
* Support for parallel test running
* Extensive use of generics, which provides additional type check at compile time

## Installation

To use Mockio in your Golang project, you can install it using the go get command:

```bash
go get github.com/ovechkin-dm/mockio
```

## Quick start
```go
package simple

import (
  . "github.com/ovechkin-dm/mockio/mock"
  "testing"
)

type myInterface interface {
  Foo(a int) int
}

func TestSimple(t *testing.T) {
  SetUp(t)
  m := Mock[myInterface]()
  WhenSingle(m.Foo(Any[int]())).ThenReturn(42)
  _ = m.Foo(10)
  Verify(m, AtLeastOnce()).Foo(10)
}


```

## Stubs and Matchers
Once you have a mock object, you can start stubbing methods or functions using the `WhenSingle`, `WhenDouble`, and `When` functions. These functions allow you to define a set of conditions under which the stubbed method or function will return a certain value.

Because golang does not support method overloading, and we still want additional type check on returning values three separate methods were introduced for stubbing:
* `WhenSingle` is used when there is only one return value for the method
* `WhenDouble` is used when there is a `(A, B)` tuple as return value for the method
* `When` for multiple return values

Example usage: 
```go
// Stub a method to always return a specific value
mockObject := Mock[MyInterface]()
WhenSingle(mockObject.MethodCall(Exact[Int](1))).ThenReturn("value")
```

Mockio also provides a set of matchers that you can use to define more complex conditions for your stubbed methods.
```go
// Use a matcher to stub a method with a specific input
mockObject := Mock[MyInterface]()
WhenSingle(mockObject.MethodCall(Equal("input"))).ThenReturn("value")
```

## Verification
Mockio also provides functionality for verifying that certain methods were called with specific inputs. You can use the `Verify` function to verify that a specific method or function was called a certain number of times, or with a specific set of inputs.
```go
// Verify that a method was called exactly once
mockObject := Mock[MyInterface]()
mockObject.MethodCall("input")
Verify(mockObject, Once()).MethodCall(Equal("input"))
```

## Argument Captors
Mockio also provides functionality for capturing the arguments passed to a mocked method or function. You can use the Captor function to create an argument captor, which you can then use to retrieve the captured arguments.
```go
// Use an argument captor to capture the argument passed to a function
mockObject := Mock[MyInterface]()
argumentCaptor := Captor[int]()
WhenSingle(mockObject.MethodCall(argumentCaptor.Capture())).ThenReturn("value")
mockObject.MethodCall(42)
capturedArgument := argumentCaptor.Last() // 42
```

## Reporting
The `ErrorReporter` interface defines how errors should be reported in the library. It has a single method `Fatalf` which takes a format string and its arguments and panics with a formatted error message.
```go

```

## Returner
The `Returner` interfaces define how the mocked function should return values. There are three different `Returner` interfaces:

- `Returner1[T any]` for functions with one return value
- `ReturnerE[T any]` for functions with two return values, where the second one is an error
- `ReturnerAll` for functions that return multiple values
Each of these interfaces provides methods ThenReturn and ThenAnswer for setting the mocked function's return value(s). ThenReturn takes one or more values and sets them as the return value(s) of the mocked function. ThenAnswer takes a function that returns the value(s) to be used as the return value(s) of the mocked function.

## Call verifiers

The `MethodVerifier` interface defines how method calls should be verified. There are four concrete implementations of this interface:

* `AtLeastOnce()` verifies that the mocked method was called at least once
* `Once()` verifies that the mocked method was called exactly once
* `Times(n int)` verifies that the mocked method was called n times
* `Never()` verifies that the mocked method was never called
The `Verify` method of the MethodVerifier interface takes a MethodVerificationData object which contains information about the method call, such as the number of times it was called. If the verification fails, an error is returned.

The `InstanceVerifier` interface defines how instances should be verified. It has a single method RecordInteraction which takes an InvocationData object containing information about the method call. If the verification fails, error is being reported.

## Conclusion
The mock package provides a powerful library for creating and managing mock objects in Go. With its support for capturing arguments, matching arguments, and verifying method calls, it makes it easy to test complex systems with many dependencies. Its well-defined interfaces and clear documentation make it easy to use and extend, and its support for multiple return values and errors makes it suitable for a wide range of use cases.


## Limitations
* **Restricted support for processor architectures**. For now library only supports amd64 architecture, but can be extended to others if there is demand for it. 
* **Go >= 1.18**
* **Concurrency limitations**
  * For now, you have to use every call to library in the same goroutine, on which  `SetUp()` was called.