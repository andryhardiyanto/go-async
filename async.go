package async

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

// AsyncFunc represents a function that can be executed concurrently.
// It receives a context to handle graceful cancellations.
type AsyncFunc func(ctx context.Context) error

// Async defines the contract for building and executing a batch of async operations.
type Async interface {
	// Task adds a function to the execution queue.
	Task(fn AsyncFunc) Async
	// WithTimeout sets a maximum duration for the entire batch to complete.
	WithTimeout(timeout time.Duration) Async
	// Go executes all queued tasks and waits for completion or the first error.
	Go(ctx context.Context) error
}

// AsyncRunner provides a factory method to create new async operation batches.
type AsyncRunner interface {
	RunInAsync() Async
}

type asyncRunner struct{}

// NewAsyncRunner creates a new instance of AsyncRunner.
func NewAsyncRunner() AsyncRunner {
	return &asyncRunner{}
}

// RunInAsync initializes a new batch of async operations.
func (a *asyncRunner) RunInAsync() Async {
	return &async{
		funcs: make([]AsyncFunc, 0),
	}
}

// Bind is a generic helper that bridges a function's result to a destination pointer.
// It ensures type safety at compile-time without the overhead of reflection.
func Bind[T any](dest *T, fn func(ctx context.Context) (T, error)) AsyncFunc {
	return func(ctx context.Context) error {
		res, err := fn(ctx)
		if err != nil {
			return err
		}
		if dest != nil {
			*dest = res
		}
		return nil
	}
}

// async implements the Async interface and manages the state of the task batch.
type async struct {
	funcs   []AsyncFunc
	timeout *time.Duration
}

// Task appends a function to the execution list.
func (a *async) Task(fn AsyncFunc) Async {
	a.funcs = append(a.funcs, fn)
	return a
}

// WithTimeout applies an optional timeout to the operation context.
func (a *async) WithTimeout(timeout time.Duration) Async {
	a.timeout = &timeout
	return a
}

// Go executes all tasks concurrently using an errgroup.
func (a *async) Go(ctx context.Context) error {
	// Apply timeout if specified to prevent goroutine leaks
	if a.timeout != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *a.timeout)
		defer cancel()
	}

	// Use errgroup for concurrency management and error propagation
	g, ctx := errgroup.WithContext(ctx)

	for _, fn := range a.funcs {
		// Re-bind the function variable to avoid closure capture issues in loops
		f := fn
		g.Go(func() (err error) {
			// Panic Recovery: Prevents the entire application from crashing on unexpected errors
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("async task panicked: %v", r)
				}
			}()

			// Pre-check if context is already cancelled before execution
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return f(ctx)
			}
		})
	}

	// Wait for all tasks to finish or return the first error encountered
	return g.Wait()
}
