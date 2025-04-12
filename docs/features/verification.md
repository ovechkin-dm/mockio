# Verification

We will use the following interface for the examples:
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
	ctrl := NewMockController(t)
	greeter := Mock[Greeter](ctrl)
	When(greeter.Greet("Jane")).ThenReturn("hello world")
	greeter.Greet("John")

}
```

## Verify

To verify that a method was called, use the `Verify` function. 
If the method was called, the test will pass. If the method was not called, the test will fail.

This test will succeed:
```go
func TestSimple(t *testing.T) {
	ctrl := NewMockController(t)
	greeter := Mock[Greeter](ctrl)
	When(greeter.Greet("Jane")).ThenReturn("hello world")
	greeter.Greet("John")
	Verify(greeter, Once()).Greet("John")
}
```

This test will fail:
```go
func TestSimple(t *testing.T) {
	ctrl := NewMockController(t)
	greeter := Mock[Greeter](ctrl)
	When(greeter.Greet("Jane")).ThenReturn("hello world")
	greeter.Greet("John")
	Verify(greeter, Once()).Greet("Jane")
}
```

### AtLeastOnce

Verify that a method was called at least once:
```go
func TestSimple(t *testing.T) {
    ctrl := NewMockController(t)
    greeter := Mock[Greeter](t)
    When(greeter.Greet("Jane")).ThenReturn("hello world")
    greeter.Greet("John")
    Verify(greeter, AtLeastOnce()).Greet("John")
}
```

### Once

Verify that a method was called exactly once:
```go
func TestSimple(t *testing.T) {
    ctrl := NewMockController(t)
    greeter := Mock[Greeter](t)
    When(greeter.Greet("Jane")).ThenReturn("hello world")
    greeter.Greet("John")
    Verify(greeter, Once()).Greet("John")
}
```

### Times

Verify that a method was called a specific number of times:
```go
func TestSimple(t *testing.T) {
    ctrl := NewMockController(t)
    greeter := Mock[Greeter](t)
    When(greeter.Greet("Jane")).ThenReturn("hello world")
    greeter.Greet("John")
    greeter.Greet("John")
    Verify(greeter, Times(2)).Greet("John")
}
```


## VerifyNoMoreInteractions

To verify that no other methods were called on the mock object, use the `VerifyNoMoreInteractions` function.
It will fail the test if there are any unverified calls.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    ctrl := NewMockController(t)
    greeter := Mock[Greeter](ctrl)
    When(greeter.Greet("Jane")).ThenReturn("hello world")
    greeter.Greet("John")
    Verify(greeter, Once()).Greet("John")
    VerifyNoMoreInteractions(greeter)
}
```

This test will fail:
```go
func TestSimple(t *testing.T) {
    ctrl := NewMockController(t)
    greeter := Mock[Greeter](ctrl)
    When(greeter.Greet("John")).ThenReturn("hello world")
    greeter.Greet("John")
    VerifyNoMoreInteractions(greeter)
}
```

## Verify after `ThenReturn`

Since it is common to actually verify that a stub was used correctly, you can use the `Verify` function after the `ThenReturn` function:
```go
func TestSimple(t *testing.T) {
	ctrl := NewMockController(t)
	greeter := Mock[Greeter](ctrl)
	When(greeter.Greet("John")).ThenReturn("hello world").Verify(Once())
	greeter.Greet("John")
	VerifyNoMoreInteractions(greeter)
}
```
