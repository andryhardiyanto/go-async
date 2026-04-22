package async

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAsyncBasicFunctionality(t *testing.T) {
	runner := NewAsyncRunner()

	var result1 int
	var result2 string

	err := runner.RunInAsync().
		Task(Bind(&result1, func(ctx context.Context) (int, error) {
			return 42, nil
		})).
		Task(Bind(&result2, func(ctx context.Context) (string, error) {
			return "hello", nil
		})).
		Go(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result1 != 42 {
		t.Errorf("Expected result1 to be 42, got %d", result1)
	}

	if result2 != "hello" {
		t.Errorf("Expected result2 to be 'hello', got %s", result2)
	}
}

func TestAsyncWithError(t *testing.T) {
	runner := NewAsyncRunner()

	var result int

	err := runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (int, error) {
			return 0, errors.New("test error")
		})).
		Go(context.Background())

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got %v", err)
	}
}

func TestAsyncWithCancellation(t *testing.T) {
	runner := NewAsyncRunner()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var result int

	err := runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (int, error) {
			time.Sleep(100 * time.Millisecond)
			return 42, nil
		})).
		Go(ctx)

	if err == nil {
		t.Fatal("Expected cancellation error, got nil")
	}
}

func TestAsyncWithRawTask(t *testing.T) {
	runner := NewAsyncRunner()

	called := false

	// Test using a raw AsyncFunc without Bind
	err := runner.RunInAsync().
		Task(func(ctx context.Context) error {
			called = true
			return nil
		}).
		Go(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !called {
		t.Error("Expected task to be called")
	}
}

func TestAsyncWithRawTaskError(t *testing.T) {
	runner := NewAsyncRunner()

	err := runner.RunInAsync().
		Task(func(ctx context.Context) error {
			return errors.New("raw task error")
		}).
		Go(context.Background())

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "raw task error" {
		t.Errorf("Expected 'raw task error', got %v", err)
	}
}

func TestAsyncBindNilDest(t *testing.T) {
	runner := NewAsyncRunner()

	// Bind with nil dest should not panic, result is discarded
	err := runner.RunInAsync().
		Task(Bind[int](nil, func(ctx context.Context) (int, error) {
			return 42, nil
		})).
		Go(context.Background())

	if err != nil {
		t.Fatalf("Expected no error for nil dest, got %v", err)
	}
}

func TestAsyncBindNilResult(t *testing.T) {
	runner := NewAsyncRunner()

	var result int

	// Return zero value (should not cause error)
	err := runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (int, error) {
			return 0, nil
		})).
		Go(context.Background())

	if err != nil {
		t.Fatalf("Expected no error for zero-value result, got %v", err)
	}

	// result should remain zero value
	if result != 0 {
		t.Errorf("Expected result to remain 0, got %d", result)
	}
}

func TestAsyncWithTimeout(t *testing.T) {
	runner := NewAsyncRunner()

	// Task that blocks until context is done, ensuring timeout is detected
	err := runner.RunInAsync().
		WithTimeout(50 * time.Millisecond).
		Task(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(200 * time.Millisecond):
				return nil
			}
		}).
		Go(context.Background())

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	// Should be a context deadline exceeded error
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

func TestAsyncWithTimeoutSuccess(t *testing.T) {
	runner := NewAsyncRunner()

	var result int

	// Task that completes within timeout
	err := runner.RunInAsync().
		WithTimeout(200 * time.Millisecond).
		Task(Bind(&result, func(ctx context.Context) (int, error) {
			return 42, nil
		})).
		Go(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != 42 {
		t.Errorf("Expected result to be 42, got %d", result)
	}
}

func TestAsyncPanicRecovery(t *testing.T) {
	runner := NewAsyncRunner()

	err := runner.RunInAsync().
		Task(func(ctx context.Context) error {
			panic("something went wrong")
		}).
		Go(context.Background())

	if err == nil {
		t.Fatal("Expected error from panic, got nil")
	}

	expected := "async task panicked: something went wrong"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got %v", expected, err)
	}
}

func TestAsyncContextPassedToTask(t *testing.T) {
	runner := NewAsyncRunner()

	type ctxKey string
	key := ctxKey("test-key")
	ctx := context.WithValue(context.Background(), key, "test-value")

	var received string

	err := runner.RunInAsync().
		Task(func(ctx context.Context) error {
			val, ok := ctx.Value(key).(string)
			if !ok {
				return errors.New("context value not found")
			}
			received = val
			return nil
		}).
		Go(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if received != "test-value" {
		t.Errorf("Expected 'test-value', got %s", received)
	}
}