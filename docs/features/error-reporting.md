# Error reporting

`Mockio` library supports providing custom error reporting in `SetUp()` function.
This can be helpful if you want to introduce custom error reporting or logging.
Reporter should implement `ErrorReporter` interface.
```go
type ErrorReporter interface {
	Fatalf(format string, args ...any)
	Errorf(format string, args ...any)
	Cleanup(func())
}

```

* `Fatalf` - should be used to report fatal errors. It should stop the test execution.
* `Errorf` - should be used to report non-fatal errors. It should continue the test execution.  
* `Cleanup` - should be used to register a cleanup function. It should be called after the test execution.

## Error output

### Incorrect `When` usage

Example:

```go
When(1)
```

Output:
```
At:
	/demo/error_reporting_test.go:22 +0xad
Cause:
	When() requires an argument which has to be 'a method call on a mock'.
	For example: When(mock.GetArticles()).ThenReturn(articles)
```

### Verify from different goroutine

Example:

```go
SetUp(r)
mock := Mock[Foo]()
wg := sync.WaitGroup{}
wg.Add(1)
go func() {
    SetUp(r)
    Verify(mock, Once())
    wg.Done()
}()
wg.Wait()
```

Output:
```
At:
	/demo/error_reporting_test.go:35 +0xc5
Cause:
	Argument passed to Verify() is {<nil> DynamicProxy[reporting.Foo] <nil>} and is not a mock, or a mock created in a different goroutine.
	Make sure you place the parenthesis correctly.
	Example of correct verification:
		Verify(mock, Times(10)).SomeMethod()
```

### Non-mock verification

Example:

```go
Verify(100, Once())
```

Output:
```
At:
	/demo/error_reporting_test.go:46 +0x105
Cause:
	Argument passed to Verify() is 100 and is not a mock, or a mock created in a different goroutine.
	Make sure you place the parenthesis correctly.
	Example of correct verification:
		Verify(mock, Times(10)).SomeMethod()
```

### Invalid use of matchers

Example:

```go
When(mock.Baz(AnyInt(), AnyInt(), 10)).ThenReturn(10)
```

Output:
```
At:
	/demo/error_reporting_test.go:55 +0x110
Cause:
	Invalid use of matchers
	3 expected, 2 recorded:
		/demo/error_reporting_test.go:55 +0xab
		/demo/error_reporting_test.go:55 +0xbc
	method:
		Foo.Baz(int, int, int) int
	expected:
		(int,int,int)
	got:
		(Any[int],Any[int])
	This can happen for 2 reasons:
		1. Declaration of matcher outside When() call
		2. Mixing matchers and exact values in When() call. Is this case, consider using "Exact" matcher.
```

### Expected method call

Example:

```go
When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10)
_ = mock.Baz(10, 10, 11)
Verify(mock, Once()).Baz(AnyInt(), AnyInt(), Exact(10))
```

Output:
```
At:
	/demo/error_reporting_test.go:88 +0x262
Cause:
	expected num method calls: 1, got : 0
		Foo.Baz(Any[int], Any[int], Exact(10))
	However, there were other interactions with this method:
		Foo.Baz(10, 10, 11) at /demo/error_reporting_test.go:87 +0x193
```

### Number of method calls

Example:

```go
When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10)
_ = mock.Baz(10, 10, 10)
Verify(mock, Times(20)).Baz(AnyInt(), AnyInt(), AnyInt())
```

Output:
```
At:
	/demo/error_reporting_test.go:121 +0x25a
Cause:
	expected num method calls: 20, got : 1
		Foo.Baz(Any[int], Any[int], Any[int])
	Invocations:
		/demo/error_reporting_test.go:120 +0x191
```

### Empty captor

Example:

```go
c := Captor[int]()
_ = c.Last()
```

Output:
```
At:
	/demo/error_reporting_test.go:130 +0x92
Cause:
	no values were captured for captor
```

### Invalid return values

Example:

```go
When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10, 20)
```

Output:
```
At:
	/demo/error_reporting_test.go:140 +0x1a7
Cause:
	invalid return values
expected:
	Foo.Baz(int, int, int) int
got:
	Foo.Baz(int, int, int) (string, int)
```

### No more interactions

Example:

```go
When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn("test", 10)
_ = mock.Baz(10, 10, 10)
_ = mock.Baz(10, 20, 10)
VerifyNoMoreInteractions(mock)
```

Output:
```
At:
	/demo/mockio/registry/registry.go:130 +0x45
Cause:
	No more interactions expected, but unverified interactions found:
		Foo.Baz(10, 10, 10) at /demo/error_reporting_test.go:150 +0x1a8
		Foo.Baz(10, 20, 10) at /demo/error_reporting_test.go:151 +0x1c6
```

### Unexpected matcher declaration

Example:

```go
When(mock.Baz(AnyInt(), AnyInt(), AnyInt())).ThenReturn(10)
mock.Baz(AnyInt(), AnyInt(), AnyInt())
Verify(mock, Once()).Baz(10, 10, 10)
```

```go
At:
	/demo/error_reporting_test.go:175 +0x23f
Cause:
	Unexpected matchers declaration.
		at /demo/error_reporting_test.go:174 +0x185
		at /demo/error_reporting_test.go:174 +0x196
		at /demo/error_reporting_test.go:174 +0x1a7
	Matchers can only be used inside When() method call.
```