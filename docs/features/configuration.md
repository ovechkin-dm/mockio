# Configuration

## Using options

MockIO library can be configured by providing options from `mockopts` package inside `SetUp` function like this:
```go
package main

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/mockopts"
	"testing"
)

func TestSimple(t *testing.T) {
	SetUp(t, mockopts.WithoutStackTrace())
}

```

## StrictVerify()
**StrictVerify** adds extra checks on each test teardown. 
It will fail the test if there are any unverified calls.
It will also fail the test if there are any calls that were not expected.

### Unverified calls check

Consider the following example:
```go
package main

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"github.com/ovechkin-dm/mockio/mockopts"
	"testing"
)

type Greeter interface {
	Greet(name string) string
}

func TestSimple(t *testing.T) {
	SetUp(t, mockopts.StrictVerify())
	greeter := Mock[Greeter]()
	When(greeter.Greet("John")).ThenReturn("Hello, John!")
}

```
In this case, the test will fail because the `Greet` method was not called with the expected argument.
If we want this test to pass, we need to call greeter with the expected argument:
```go
func TestSimple(t *testing.T) {
	SetUp(t, mockopts.StrictVerify())
	greeter := Mock[Greeter]()
	When(greeter.Greet("John")).ThenReturn("Hello, John!")
	greeter.Greet("John")
}
```

### Unexpected calls check

Consider the following example:

```go
func TestSimple(t *testing.T) {
    SetUp(t, mockopts.StrictVerify())
    greeter := Mock[Greeter]()
    When(greeter.Greet("John")).ThenReturn("Hello, John!")
    greeter.Greet("John")
    greeter.Greet("Jane")
}
```

In this case, the test will fail because the `Greet` method was called with an unexpected argument.
If we want this test to pass, we need to remove the unexpected call, or add an expectation for it:
```go
func TestSimple(t *testing.T) {
    SetUp(t, mockopts.StrictVerify())
    greeter := Mock[Greeter]()
    When(greeter.Greet("John")).ThenReturn("Hello, John!")
    When(greeter.Greet("Jane")).ThenReturn("Hello, Jane!")
    greeter.Greet("John")
    greeter.Greet("Jane")
}
```

## WithoutStackTrace()
**WithoutStackTrace** option disables stack trace printing in case of test failure.

Consider the following example:
```go
package main

import (
	. "github.com/ovechkin-dm/mockio/mock"
	"testing"
)

type Greeter interface {
	Greet(name string) string
}

func TestSimple(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	When(greeter.Greet("Jane")).ThenReturn("hello world")
	greeter.Greet("John")
	VerifyNoMoreInteractions(greeter)
}

```

If we run this test, we will see the following error:
```
=== RUN   TestSimple
    reporter.go:75: At:
        	go/pkg/mod/github.com/ovechkin-dm/mockio@v0.7.2/registry/registry.go:130 +0x45
        Cause:
        	No more interactions expected, but unverified interactions found:
        		Greeter.Greet(John) at demo/hello_test.go:16 +0xf2
        Trace:
        demo.TestSimple.VerifyNoMoreInteractions.VerifyNoMoreInteractions.func1()
        	go/pkg/mod/github.com/ovechkin-dm/mockio@v0.7.2/registry/registry.go:130 +0x45
        demo.TestSimple(0xc00018c4e0?)
        	demo/hello_test.go:17 +0x15a
        testing.tRunner(0xc00018c4e0, 0x647ca0)
        	/usr/local/go/src/testing/testing.go:1689 +0xfb
        created by testing.(*T).Run in goroutine 1
        	/usr/local/go/src/testing/testing.go:1742 +0x390
        
--- FAIL: TestSimple (0.00s)

FAIL
```

By adding `mockopts.WithoutStackTrace()` to the `SetUp` function, we can disable stack trace printing:
```go
func TestSimple(t *testing.T) {
	SetUp(t, mockopts.WithoutStackTrace())
	greeter := Mock[Greeter]()
	When(greeter.Greet("Jane")).ThenReturn("hello world")
	greeter.Greet("John")
	VerifyNoMoreInteractions(greeter)
}
```

Now the error will look like this:
```
=== RUN   TestSimple
    reporter.go:75: At:
        	go/pkg/mod/github.com/ovechkin-dm/mockio@v0.7.2/registry/registry.go:130 +0x45
        Cause:
        	No more interactions expected, but unverified interactions found:
        		Greeter.Greet(John) at demo/hello_test.go:17 +0x10b
--- FAIL: TestSimple (0.00s)

FAIL
```