# Argument Captors

Argument captors are a powerful feature that allow you to capture the arguments passed to a method when it is
called. This is useful when you want to verify that a method was called with specific arguments, but you don't know what
those arguments will be ahead of time.

## Creating a Captor

To create a captor, you simply call the `Captor` function with the type of the argument you want to capture:

```go
c := Captor[string]()
```

## Using a Captor

To use a captor, you pass it as an argument to the `When` function. When the method is called, the captor will capture the
argument and store it in the captor's value:

```go
When(greeter.Greet(c.Capture())).ThenReturn("Hello, world!")
```

## Retrieving the Captured Values

Argument captor records an argument on each stub call. You can retrieve the captured values by calling the `Values` method

```go
capturedValues := c.Values()
```

If you want to retrieve just the last captured value, you can call the `Last` method

```go
capturedValue := c.Last()
```

## Example usage

In this example we will create a mock, and use an argument captor to capture the arguments passed to the `Greet` method:

```go
package main

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"testing"
)

type Greeter interface {
	Greet(name any) string
}

func TestSimple(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	c := Captor[string]()
	When(greeter.Greet(c.Capture())).ThenReturn("Hello, world!")
	_ = greeter.Greet("John")
	_ = greeter.Greet("Jane")
	if c.Values()[0] != "John" {
		t.Error("Expected John, got", c.Values()[0])
	}
	if c.Values()[1] != "Jane" {
		t.Error("Expected Jane, got", c.Values()[1])
	}
}
```