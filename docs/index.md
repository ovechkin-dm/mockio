# Mockio

Golang library for mocking without code generation, inspired by Mockito.

## Installing library

Install latest version of the library using `go get` command:

```bash
go get -u github.com/ovechkin-dm/mockio
```

## Creating test

Let's create an interface that we want to mock:

```go
type Greeter interface {
    Greet(name string) string
}
```

Now we will use `dot import` to simplify the usage of the library:

```go
import (
    ."github.com/ovechkin-dm/mockio/mock"
    "testing"
)
```

Now we can create a mock for the `Greeter` interface, and test it's method `Greet`:

```go
func TestGreet(t *testing.T) {
    ctrl := NewMockController(t)
    m := Mock[Greeter](ctrl)
    WhenSingle(m.Greet("John")).ThenReturn("Hello, John!")
    if m.Greet("John") != "Hello, John!" {
        t.Fail()
    }
}
```

## Full example
Here is the full listing for our simple test:

```go
package main

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"testing"
)

type Greeter interface {
	Greet(name string) string
}

func TestGreet(t *testing.T) {
	ctrl := NewMockController(t)
	m := Mock[Greeter](ctrl)
	WhenSingle(m.Greet("John")).ThenReturn("Hello, John!")
	if m.Greet("John") != "Hello, John!" {
		t.Fail()
	}
}

```

That's it! You have created a mock for the `Greeter` interface without any code generation.
As you can see, the library is very simple and easy to use.
And no need to generate mocks for your interfaces.
