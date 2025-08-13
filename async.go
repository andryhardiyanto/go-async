package async

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"golang.org/x/sync/errgroup"
)

// AsyncFunc represents a function that can be executed asynchronously.
// It returns a result of any type and an error if the operation fails.
type AsyncFunc func() (any, error)

// Async interface defines the contract for building and executing async operations.
type Async interface {
	Task(dest any, asyncFunc AsyncFunc) Async
	WithTimeout(timeout time.Duration) Async
	Go(ctx context.Context) error
}

// AsyncRunner interface provides a factory method to create new async operation batches.
type AsyncRunner interface {
	RunInAsync() Async
}

// asyncRunner implements the AsyncRunner interface.
type asyncRunner struct{}

// NewAsyncRunner creates a new instance of AsyncRunner.
// This is the main entry point for using the async package.
func NewAsyncRunner() AsyncRunner {
	return &asyncRunner{}
}

// RunInAsync creates a new async operation batch with an empty task queue.
func (a *asyncRunner) RunInAsync() Async {
	return &async{
		funcs: make([]*asyncHolder, 0),
	}
}

// Task adds an async function to the execution queue.
// dest: pointer to store the result (must be a pointer type)
// asyncFunc: function to execute asynchronously
// Returns the same Async instance for method chaining.
func (a *async) Task(dest any, asyncFunc AsyncFunc) Async {
	a.funcs = append(a.funcs, &asyncHolder{
		dest: dest,
		fun:  asyncFunc,
	})
	return a
}

// WithTimeout sets a timeout for all async operations in this batch.
// If any operation takes longer than the specified duration, all operations will be cancelled.
// timeout: maximum duration to wait for all operations to complete
// Returns the same Async instance for method chaining.
func (a *async) WithTimeout(timeout time.Duration) Async {
	a.timeout = &timeout
	return a
}

// Go executes all queued async operations concurrently.
// ctx: context for cancellation and timeout control
// Returns an error if any operation fails, times out, or if the context is cancelled.
// All operations are executed concurrently using errgroup for proper error handling.
func (a *async) Go(ctx context.Context) error {
	// Apply timeout if specified
	if a.timeout != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *a.timeout)
		defer cancel()
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, fn := range a.funcs {
		asyncTask(ctx, g, fn.dest, fn.fun)
	}

	// Wait for all tasks to complete and return any error
	return g.Wait()
}

// async implements the Async interface and holds the state for a batch of async operations.
type async struct {
	funcs   []*asyncHolder // Queue of async functions to execute
	timeout *time.Duration // Optional timeout for all operations
}

// asyncHolder holds a single async operation with its destination and function.
type asyncHolder struct {
	dest any       // Pointer to store the result
	fun  AsyncFunc // Function to execute asynchronously
}

// asyncTask schedules a single async function for execution within an errgroup.
// ctx: context for cancellation control
// g: errgroup to manage concurrent execution
// dest: destination pointer for storing the result
// fn: async function to execute
func asyncTask(ctx context.Context, g *errgroup.Group, dest any, fn AsyncFunc) {
	g.Go(func() error {
		res, err := fn()

		select {
		case <-ctx.Done():
			return fmt.Errorf("async operation was cancelled: %w", ctx.Err())
		default:
		}

		// Handle result assignment with improved type safety
		if dest != nil && res != nil {
			if assignErr := assignResult(dest, res); assignErr != nil {
				return fmt.Errorf("failed to assign result: %w", assignErr)
			}
		}

		return err
	})
}

// assignResult safely assigns the result to destination with comprehensive type checking.
// This function uses reflection to ensure type safety and provide detailed error messages.
// dest: destination pointer where the result will be stored
// res: result value to assign
// Returns an error if the assignment is not possible due to type mismatch or invalid destination.
func assignResult(dest any, res any) error {
	val := reflect.ValueOf(dest)

	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer, received %T (kind: %v)", dest, val.Kind())
	}

	if val.IsNil() {
		return fmt.Errorf("destination pointer cannot be nil")
	}

	destElem := val.Elem()
	resVal := reflect.ValueOf(res)

	if !resVal.Type().AssignableTo(destElem.Type()) {
		return fmt.Errorf("type mismatch: cannot assign %T to %T", res, destElem.Interface())
	}

	destElem.Set(resVal)
	return nil
}
