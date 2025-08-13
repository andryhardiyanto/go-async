# Go Async - Concurrent Function Execution Library

[![Go Version](https://img.shields.io/badge/Go-1.24.3-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen.svg)](#testing)

A simple, efficient, and type-safe Go library for executing multiple functions concurrently with comprehensive error handling and timeout support.

## Features

- 🚀 **Concurrent Execution**: Run multiple functions simultaneously using goroutines
- 🔒 **Type Safety**: Compile-time and runtime type checking with reflection
- ⏱️ **Timeout Support**: Set timeouts for async operations
- 🛡️ **Error Handling**: Comprehensive error propagation and context cancellation
- 🔗 **Method Chaining**: Fluent API for easy usage
- 📦 **Zero Dependencies**: Only uses Go standard library and `golang.org/x/sync/errgroup`

## Installation

```bash
go get github.com/andryhardiyanto/go-async
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/andryhardiyanto/go-async"
)

func main() {
    runner := async.NewAsyncRunner()
    
    var result1 int
    var result2 string
    
    err := runner.RunInAsync().
        Task(&result1, func() (any, error) {
            time.Sleep(100 * time.Millisecond)
            return 42, nil
        }).
        Task(&result2, func() (any, error) {
            time.Sleep(50 * time.Millisecond)
            return "Hello, World!", nil
        }).
        Go(context.Background())
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Result 1: %d\n", result1) // Output: Result 1: 42
    fmt.Printf("Result 2: %s\n", result2) // Output: Result 2: Hello, World!
}
```

## API Reference

### Core Interfaces

#### `AsyncRunner`

```go
type AsyncRunner interface {
    RunInAsync() Async
}
```

Factory interface for creating async operation batches.

#### `Async`

```go
type Async interface {
    Task(dest any, asyncFunc AsyncFunc) Async
    WithTimeout(timeout time.Duration) Async
    Go(ctx context.Context) error
}
```

Main interface for building and executing async operations.

#### `AsyncFunc`

```go
type AsyncFunc func() (any, error)
```

Function signature for async operations.

### Methods

#### `NewAsyncRunner() AsyncRunner`

Creates a new AsyncRunner instance.

#### `Task(dest any, asyncFunc AsyncFunc) Async`

Adds an async function to the execution queue.

- `dest`: Pointer to store the result (must be a pointer type)
- `asyncFunc`: Function to execute asynchronously
- Returns: Same Async instance for method chaining

#### `WithTimeout(timeout time.Duration) Async`

Sets a timeout for all async operations in the batch.

- `timeout`: Maximum duration to wait for all operations
- Returns: Same Async instance for method chaining

#### `Go(ctx context.Context) error`

Executes all queued async operations concurrently.

- `ctx`: Context for cancellation and timeout control
- Returns: Error if any operation fails, times out, or context is cancelled

## Usage Examples

### Basic Usage

```go
runner := async.NewAsyncRunner()

var num int
var text string

err := runner.RunInAsync().
    Task(&num, func() (any, error) {
        return 100, nil
    }).
    Task(&text, func() (any, error) {
        return "async result", nil
    }).
    Go(context.Background())

if err != nil {
    log.Fatal(err)
}
```

### With Timeout

```go
runner := async.NewAsyncRunner()

var result int

err := runner.RunInAsync().
    WithTimeout(5 * time.Second).
    Task(&result, func() (any, error) {
        time.Sleep(2 * time.Second)
        return 42, nil
    }).
    Go(context.Background())

if err != nil {
    log.Printf("Operation failed: %v", err)
}
```

### Error Handling

```go
runner := async.NewAsyncRunner()

var result int

err := runner.RunInAsync().
    Task(&result, func() (any, error) {
        return nil, errors.New("something went wrong")
    }).
    Go(context.Background())

if err != nil {
    log.Printf("Async operation failed: %v", err)
}
```

### Context Cancellation

```go
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(1 * time.Second)
    cancel() // Cancel after 1 second
}()

runner := async.NewAsyncRunner()
var result int

err := runner.RunInAsync().
    Task(&result, func() (any, error) {
        time.Sleep(5 * time.Second) // This will be cancelled
        return 42, nil
    }).
    Go(ctx)

if err != nil {
    log.Printf("Operation cancelled: %v", err)
}
```

### Multiple Operations

```go
runner := async.NewAsyncRunner()

var (
    userID   int
    userName string
    userAge  int
    isActive bool
)

err := runner.RunInAsync().
    Task(&userID, func() (any, error) {
        return fetchUserID(), nil
    }).
    Task(&userName, func() (any, error) {
        return fetchUserName(), nil
    }).
    Task(&userAge, func() (any, error) {
        return fetchUserAge(), nil
    }).
    Task(&isActive, func() (any, error) {
        return fetchUserStatus(), nil
    }).
    Go(context.Background())

if err != nil {
    log.Fatal(err)
}
```

## Error Types

The library provides detailed error messages for different failure scenarios:

- **Type Mismatch**: `type mismatch: cannot assign int to string`
- **Invalid Destination**: `destination must be a pointer, received int (kind: int)`
- **Nil Pointer**: `destination pointer cannot be nil`
- **Timeout**: `async operation was cancelled: context deadline exceeded`
- **Cancellation**: `async operation was cancelled: context canceled`
- **Assignment Error**: `failed to assign result: [specific error]`

## Testing

Run the test suite:

```bash
go test -v
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

### Test Coverage

The library includes comprehensive tests covering:

- ✅ Basic functionality
- ✅ Error handling
- ✅ Context cancellation
- ✅ Timeout operations
- ✅ Type compatibility
- ✅ Nil pointer handling
- ✅ Performance benchmarks

## Performance

Benchmark results on Apple M1 Pro:

```
BenchmarkAsyncExecution-8                 647395    1855 ns/op    640 B/op    19 allocs/op
BenchmarkAsyncExecutionManyTasks-8        261549    4996 ns/op   1720 B/op    50 allocs/op
```

## Best Practices

1. **Always use pointers** for destination variables
2. **Handle errors** returned by the `Go()` method
3. **Use context** for cancellation and timeout control
4. **Set reasonable timeouts** for long-running operations
5. **Avoid shared state** between async functions without proper synchronization

## Requirements

- Go 1.24.3 or later
- `golang.org/x/sync/errgroup` package

## Contributing

We welcome contributions! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## Support

If you encounter any issues or have questions, please open an issue on GitHub.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.