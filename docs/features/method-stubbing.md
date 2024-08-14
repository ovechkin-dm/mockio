# Method stubbing

Method stubbing is a technique used in unit testing to replace a method with a stub. A stub is a small piece of code
that simulates the behavior of the method it replaces. This allows you to test the behavior of the code that calls the
method without actually executing the method itself.

Basic usage of method stubbing in Mockio looks like this:

```go
When(mock.SomeMethod(AnyInt())).ThenReturn("some value")
```

* `When` is a function that takes a method call as an argument and returns a `Returner` object.
* Inside the method call argument you can use any matcher from the library's API. In this example we used `AnyInt()` matcher.
* `ThenReturn` is a method of the `Returner`

This is basic usage of method stubbing. But there are also some useful extensions to this API.

## When

`When` is a function that allows you to stub a method.
Keep in mind, that `When` is a generic function, so it does not provide any type check on return value.


## WhenSingle 

`WhenSingle` is a function that allows you to stub a method to return a single value. 
It is almost the same as `When`, but it provides additional type check on return value.

Consider Following interface:
```go
type Foo interface {
    Bar(int) string
}
```

You can stub `Bar` method like this:
```go
WhenSingle(mock.Bar(AnyInt())).ThenReturn("some value")
```

However, this will not compile:
```go
WhenSingle(mock.Bar(AnyInt())).ThenReturn(42)
```

But this will:
```go
When(mock.Bar(AnyInt())).ThenReturn(42)
```

## WhenDouble

`WhenDouble` is a function that allows you to stub a method to return two values.
It is almost the same as `When`, but it provides additional type check on return values.

Consider Following interface:
```go
type Foo interface {
    Bar(int) (string, error)
}
```

You can stub `Bar` method like this:
```go
WhenDouble(mock.Bar(AnyInt())).ThenReturn("some value", nil)
```

However, this will not compile:
```go
WhenDouble(mock.Bar(AnyInt())).ThenReturn("some value", 42)
```
 
But this will:
```go
When(mock.Bar(AnyInt())).ThenReturn("some value", 42)
```

## ThenAnswer

`Answer` is a function that allows you to stub a method to return a value based on the arguments passed to the method.

Consider following interface:
```go
type Foo interface {
    Bar(int) string
}
```

You can stub `Bar` method like this:
```go
mock := Mock[Foo]()
WhenSingle(mock.Bar(AnyInt())).ThenAnswer(func(args []any) string {
    return fmt.Sprintf("Hello, %d", args[0].(int))
})
```

When `Bar` method is called with argument `42`, it will return `"Hello, 42"`.

## ThenReturn

You can chain multiple `ThenReturn` calls to return different values on subsequent calls:

```go
When(mock.SomeMethod(AnyInt())).
    ThenReturn("first value").
    ThenReturn("second value")
```

Calling `SomeMethod` first time will return `"first value"`, second time `"second value"`, and so on.

## Implicit `Exact` matchers

Consider following interface:

```go
type Foo interface {
    Bar(int, int) string
}

```

To stub `Bar` method, we can use something like this:
```go
When(mock.Bar(Exact(1), Exact(2))).ThenReturn("some value")
```

However, this can be simplified to:
```go
When(mock.Bar(1, 2)).ThenReturn("some value")
```

In short, you can omit using matchers when you want to match exact values, but they all should be exact.
For example, this will not work:
```go
When(mock.Bar(1, Exact(2))).ThenReturn("some value")
```
