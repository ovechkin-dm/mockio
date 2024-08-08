# Matchers
MockIO library provides a lot of ways to match arguments of the method calls.
Matchers are used to define the expected arguments of the method calls.

## Test template

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
	SetUp(t)
	greeter := Mock[Greeter]()
	When(greeter.Greet("Jane")).ThenReturn("hello world")
	greeter.Greet("John")
	VerifyNoMoreInteractions(greeter)
}
```

## Any\[T]()
The `Any[T]()` matcher matches any value of the type `T`.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(Any[string]())).ThenReturn("hello world")
    if greeter.Greet("John") != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## AnyInt()
The `AnyInt()` matcher matches any integer value. 

This test will succeed:
```go
func TestSimple(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	When(greeter.Greet(Any[int]())).ThenReturn("hello world")
	if greeter.Greet(10) != "hello world" {
		t.Error("Expected 'hello world'")
	}
}
```
This test will fail, because the argument is not an integer:
```go
func TestSimple(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	When(greeter.Greet(Any[int]())).ThenReturn("hello world")
	if greeter.Greet("John") != "hello world" {
		t.Error("Expected 'hello world'")
	}
}
```

## AnyString()
The `AnyString()` matcher matches any string value.

This test will succeed:
```go
func TestSimple(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	When(greeter.Greet(Any[string]())).ThenReturn("hello world")
	if greeter.Greet("John") != "hello world" {
		t.Error("Expected 'hello world'")
	}
}
```
This test will fail, because the argument is not a string:
```go
func TestSimple(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	When(greeter.Greet(Any[int]())).ThenReturn("hello world")
	if greeter.Greet(10) != "hello world" {
		t.Error("Expected 'hello world'")
	}
}
```

## AnyInterface()
The `AnyInterface()` matcher matches any value of any type.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(Any[interface{}]())).ThenReturn("hello world")
    if greeter.Greet("John") != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

This test will also succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(Any[interface{}]())).ThenReturn("hello world")
    if greeter.Greet(10) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## AnyContext()
The `AnyContext()` matcher matches any context.Context value.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(Any[context.Context]())).ThenReturn("hello world")
    if greeter.Greet(context.Background()) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## AnyOfType\[T](t T)
The `AnyOfType[T](t T)` matcher matches any value of the type `T` or its subtype. It is useful for type inference.

This test will succeed:
```go
func TestSimple(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	When(greeter.Greet(AnyOfType(10))).ThenReturn("hello world")
	if greeter.Greet(10) != "hello world" {
		t.Error("Expected 'hello world'")
	}
}
```
Note that when we are using AnyOfType, we don't need to specify the type explicitly.

##  Nil\[T\]()
The `Nil[T]()` matcher matches any nil value of the type `T`.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(Nil[any]())).ThenReturn("hello world")
    if greeter.Greet(nil) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## NotNil\[T\]()
The `NotNil[T]()` matcher matches any non-nil value of the type `T`.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(NotNil[any]())).ThenReturn("hello world")
    if greeter.Greet("John") != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

This test will fail:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(NotNil[any]())).ThenReturn("hello world")
    if greeter.Greet(nil) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## Regex(pattern string)
The `Regex(pattern string)` matcher matches any string that matches the regular expression `pattern`.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(Regex("J.*"))).ThenReturn("hello world")
    if greeter.Greet("John") != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## Substring(sub string)
The `Substring(sub string)` matcher matches any string that contains the substring `sub`.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(Substring("oh"))).ThenReturn("hello world")
    if greeter.Greet("John") != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## SliceLen(length int)
The `SliceLen(length int)` matcher matches any slice with the length `length`.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(SliceLen[int](2))).ThenReturn("hello world")
    if greeter.Greet([]int{1, 2}) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

This test will fail:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(SliceLen[int](2))).ThenReturn("hello world")
    if greeter.Greet([]int{1, 2, 3}) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## MapLen(length int)
The `MapLen(length int)` matcher matches any map with the length `length`.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(MapLen[int, string](2))).ThenReturn("hello world")
    if greeter.Greet(map[int]string{1: "one", 2: "two"}) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

This test will fail:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(MapLen[int, string](2))).ThenReturn("hello world")
    if greeter.Greet(map[int]string{1: "one", 2: "two", 3: "three"}) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## SliceContains\[T](values ...T)
The `SliceContains[T any](values ...T)` matcher matches any slice that contains all the values `values`.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(SliceContains[int](1, 2))).ThenReturn("hello world")
    if greeter.Greet([]int{1, 2, 3}) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

This test will fail:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(SliceContains[int](1, 2))).ThenReturn("hello world")
    if greeter.Greet([]int{1, 3}) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## MapContains\[K comparable, V](keys ...K)
The `MapContains[K any, V any](keys ...K)` matcher matches any map that contains all the keys `keys`.
 
This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(MapContains[int, string](1, 2))).ThenReturn("hello world")
    if greeter.Greet(map[int]string{1: "one", 2: "two", 3: "three"}) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

This test will fail:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(MapContains[int, string](1, 2))).ThenReturn("hello world")
    if greeter.Greet(map[int]string{1: "one", 3: "three"}) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## SliceEqualUnordered\[T](values \[]T)

The `SliceEqualUnordered[T any](values []T)` matcher matches any slice that contains the same elements as `values`, but in any order.

This test will succeed:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(SliceEqualUnordered[int](1, 2))).ThenReturn("hello world")
    if greeter.Greet([]int{2, 1}) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

This test will fail:
```go
func TestSimple(t *testing.T) {
    SetUp(t)
    greeter := Mock[Greeter]()
    When(greeter.Greet(SliceEqualUnordered[int](1, 2))).ThenReturn("hello world")
    if greeter.Greet([]int{1, 3}) != "hello world" {
        t.Error("Expected 'hello world'")
    }
}
```

## Exact

The `Exact` matcher matches any value that is equal to the expected value.
`Exact` uses `==` operator to compare values.

This test will succeed:
```go
func TestSimple(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	world1 := "world"
	When(greeter.Greet(Exact(&world1))).ThenReturn("hello world")
	if greeter.Greet(&world1) != "hello world" {
		t.Error("Expected hello world")
	}
}
```

However, this test will fail, because although the values are equal, they are different pointers:
```go
func TestSimple(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	world1 := "world"
	world2 := "world"
	When(greeter.Greet(Exact(&world1))).ThenReturn("hello world")
	if greeter.Greet(world2) != "hello world" {
		t.Error("Expected hello world")
	}
}
```