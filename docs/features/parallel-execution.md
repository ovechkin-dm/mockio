# Parallel execution

## Parallelism

It is possible to run multiple tests with mockio in parallel using the `--parallel` option. This option is available in the `test` and `run` commands.

## Concurrency

Library supports invoking stubbed methods from different goroutine.

```go
func TestParallelSuccess(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	wg := sync.WaitGroup{}
	wg.Add(2)
	When(greeter.Greet("John")).ThenReturn("hello world")
	go func() {
		greeter.Greet("John")
		wg.Done()
	}()
	go func() {
		greeter.Greet("John")
		wg.Done()
	}()
	wg.Wait()
	Verify(greeter, Times(2)).Greet("John")
}
```

However, library does not support stubbing methods from different goroutine. 
This test will result in error:

```go
func TestParallelFail(t *testing.T) {
	SetUp(t)
	greeter := Mock[Greeter]()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		When(greeter.Greet("John")).ThenReturn("hello world")
		wg.Done()
	}()
	go func() {
		When(greeter.Greet("John")).ThenReturn("hello world")
		wg.Done()
	}()
	wg.Wait()
	if greeter.Greet("John") != "hello world" {
		t.Error("Expected 'hello world'")
	}
}
```

The main rule is that call to `When` should be in the same goroutine in which the mock is created.

Also, each time you create a mock in a newly created goroutine, you need to call `SetUp(t)` again to initialize the mockio library in that goroutine.

```go
func TestParallelSuccess(t *testing.T) {
	SetUp(t)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		SetUp(t)
		greeter := Mock[Greeter]()
		When(greeter.Greet("John")).ThenReturn("hello world")
		wg.Done()
	}()
	wg.Wait()
}
```