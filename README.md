# Go Async - Concurrent Function Execution Library

[![Go Version](https://img.shields.io/badge/Go-1.26.2-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen.svg)](#testing)

A simple, efficient, and type-safe Go library for executing multiple functions concurrently with built-in panic recovery, timeout support, and compile-time generics for result binding.

## Features

- 🚀 **Concurrent Execution**: Run multiple functions simultaneously using goroutines
- 🔒 **Compile-Time Type Safety**: Generic `Bind[T]` helper ensures type safety without reflection
- ⏱️ **Timeout Support**: Set timeouts for async operations
- 🛡️ **Panic Recovery**: Goroutine panics are caught and returned as errors
- 🔗 **Method Chaining**: Fluent API for easy usage
- 🧩 **Context Propagation**: Each task receives the parent context for cancellation awareness
- 📦 **Minimal Dependencies**: Only uses Go standard library and `golang.org/x/sync/errgroup`

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
        Task(async.Bind(&result1, func(ctx context.Context) (int, error) {
            time.Sleep(100 * time.Millisecond)
            return 42, nil
        })).
        Task(async.Bind(&result2, func(ctx context.Context) (string, error) {
            time.Sleep(50 * time.Millisecond)
            return "Hello, World!", nil
        })).
        Go(context.Background())

    if err != nil {
        panic(err)
    }

    fmt.Printf("Result 1: %d\n", result1) // Output: Result 1: 42
    fmt.Printf("Result 2: %s\n", result2) // Output: Result 2: Hello, World!
}
```

## API Reference

### Core Types

#### `AsyncFunc`

```go
type AsyncFunc func(ctx context.Context) error
```

A function that can be executed concurrently. It receives a context to handle graceful cancellations.

#### `Async`

```go
type Async interface {
    Task(fn AsyncFunc) Async
    WithTimeout(timeout time.Duration) Async
    Go(ctx context.Context) error
}
```

Interface for building and executing a batch of async operations.

#### `AsyncRunner`

```go
type AsyncRunner interface {
    RunInAsync() Async
}
```

Factory interface for creating async operation batches.

### Functions

#### `NewAsyncRunner() AsyncRunner`

Creates a new AsyncRunner instance.

#### `Bind[T any](dest *T, fn func(ctx context.Context) (T, error)) AsyncFunc`

A generic helper (also referred to as `Await` in some patterns) that bridges a function's result to a destination pointer. It ensures type safety at compile-time without the overhead of reflection.

- `dest`: Pointer to store the result (pass `nil` to discard the result)
- `fn`: Function that returns a typed result and an error
- Returns: An `AsyncFunc` that can be passed to `Task()`

### Methods

#### `Task(fn AsyncFunc) Async`

Adds a function to the execution queue.

- `fn`: An `AsyncFunc` to execute concurrently (use `Bind()` to capture results)
- Returns: Same Async instance for method chaining

#### `WithTimeout(timeout time.Duration) Async`

Sets a maximum duration for the entire batch to complete.

- `timeout`: Maximum duration to wait for all operations
- Returns: Same Async instance for method chaining

#### `Go(ctx context.Context) error`

Executes all queued tasks concurrently and waits for completion or the first error.

- `ctx`: Context for cancellation and timeout control
- Returns: Error if any operation fails, panics, times out, or context is cancelled

## Usage Examples

### Basic Usage with Bind

```go
runner := async.NewAsyncRunner()

var num int
var text string

err := runner.RunInAsync().
    Task(async.Bind(&num, func(ctx context.Context) (int, error) {
        return 100, nil
    })).
    Task(async.Bind(&text, func(ctx context.Context) (string, error) {
        return "async result", nil
    })).
    Go(context.Background())

if err != nil {
    log.Fatal(err)
}
```

### Raw Task (No Result Binding)

For tasks that don't need to return a value, pass an `AsyncFunc` directly:

```go
runner := async.NewAsyncRunner()

err := runner.RunInAsync().
    Task(func(ctx context.Context) error {
        // perform side-effect work, e.g. send an email
        return sendNotification(ctx)
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
    Task(async.Bind(&result, func(ctx context.Context) (int, error) {
        time.Sleep(2 * time.Second)
        return 42, nil
    })).
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
    Task(async.Bind(&result, func(ctx context.Context) (int, error) {
        return 0, errors.New("something went wrong")
    })).
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
    Task(async.Bind(&result, func(ctx context.Context) (int, error) {
        time.Sleep(5 * time.Second) // This will be cancelled
        return 42, nil
    })).
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
    Task(async.Bind(&userID, func(ctx context.Context) (int, error) {
        return fetchUserID()
    })).
    Task(async.Bind(&userName, func(ctx context.Context) (string, error) {
        return fetchUserName()
    })).
    Task(async.Bind(&userAge, func(ctx context.Context) (int, error) {
        return fetchUserAge()
    })).
    Task(async.Bind(&isActive, func(ctx context.Context) (bool, error) {
        return fetchUserStatus()
    })).
    Go(context.Background())

if err != nil {
    log.Fatal(err)
}
```

## Error Handling

The library handles the following error scenarios:

- **Task Error**: Any error returned by an `AsyncFunc` is propagated
- **Panic Recovery**: `async task panicked: <panic value>`
- **Timeout**: `context deadline exceeded`
- **Cancellation**: `context canceled`

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

- ✅ Basic functionality with `Bind`
- ✅ Error handling
- ✅ Context cancellation
- ✅ Timeout operations
- ✅ Raw task execution (without `Bind`)
- ✅ Nil destination with `Bind`
- ✅ Panic recovery
- ✅ Context propagation to tasks

## Best Practices

1. **Use `Bind[T]` for result capture** — it provides compile-time type safety without reflection
2. **Use raw `AsyncFunc` for side-effects** — when no result is needed, pass a `func(ctx context.Context) error` directly
3. **Handle errors** returned by the `Go()` method
4. **Use context** for cancellation and timeout control
5. **Set reasonable timeouts** for long-running operations
6. **Avoid shared state** between async functions without proper synchronization
7. **Leverage context in tasks** — `AsyncFunc` receives context, use it for downstream calls (HTTP, DB, etc.)

## Requirements

- Go 1.26 or later (generics support required)
- `golang.org/x/sync/errgroup` package

## Contributing

We welcome contributions! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## Support

If you encounter any issues or have questions, please open an issue on GitHub.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
