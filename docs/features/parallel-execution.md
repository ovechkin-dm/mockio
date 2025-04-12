# Parallel execution

## Parallelism

It is possible to run multiple tests with mockio in parallel using the `--parallel` option. This option is available in the `test` and `run` commands.

## Concurrency

Library supports invoking stubbed methods from different goroutine.

```go
func TestParallelSuccess(t *testing.T) {
	ctrl := NewMockController(t)
	greeter := Mock[Greeter](ctrl)
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