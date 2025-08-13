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
		Task(&result1, func() (any, error) {
			return 42, nil
		}).
		Task(&result2, func() (any, error) {
			return "hello", nil
		}).
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
		Task(&result, func() (any, error) {
			return nil, errors.New("test error")
		}).
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
		Task(&result, func() (any, error) {
			time.Sleep(100 * time.Millisecond)
			return 42, nil
		}).
		Go(ctx)
	
	if err == nil {
		t.Fatal("Expected cancellation error, got nil")
	}
}

func TestAsyncTypeCompatibility(t *testing.T) {
	runner := NewAsyncRunner()
	
	var result string
	
	// This should fail because we're trying to assign int to string
	err := runner.RunInAsync().
		Task(&result, func() (any, error) {
			return 42, nil // int instead of string
		}).
		Go(context.Background())
	
	if err == nil {
		t.Fatal("Expected type compatibility error, got nil")
	}
}

func TestAsyncNilPointer(t *testing.T) {
	runner := NewAsyncRunner()
	
	var result *int // nil pointer
	
	err := runner.RunInAsync().
		Task(result, func() (any, error) {
			return 42, nil
		}).
		Go(context.Background())
	
	if err == nil {
		t.Fatal("Expected nil pointer error, got nil")
	}
}

func TestAsyncNonPointerDest(t *testing.T) {
	runner := NewAsyncRunner()
	
	var result int
	
	// Pass non-pointer (should fail)
	err := runner.RunInAsync().
		Task(result, func() (any, error) {
			return 42, nil
		}).
		Go(context.Background())
	
	if err == nil {
		t.Fatal("Expected non-pointer error, got nil")
	}
}

func TestAsyncNilResult(t *testing.T) {
	runner := NewAsyncRunner()
	
	var result int
	
	// Return nil result (should not cause error)
	err := runner.RunInAsync().
		Task(&result, func() (any, error) {
			return nil, nil
		}).
		Go(context.Background())
	
	if err != nil {
		t.Fatalf("Expected no error for nil result, got %v", err)
	}
	
	// result should remain unchanged (zero value)
	if result != 0 {
		t.Errorf("Expected result to remain 0, got %d", result)
	}
}

func TestAsyncWithTimeout(t *testing.T) {
	runner := NewAsyncRunner()
	
	var result int
	
	// Task that takes longer than timeout
	err := runner.RunInAsync().
		WithTimeout(50 * time.Millisecond).
		Task(&result, func() (any, error) {
			time.Sleep(100 * time.Millisecond) // Sleep longer than timeout
			return 42, nil
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
		WithTimeout(100 * time.Millisecond).
		Task(&result, func() (any, error) {
			time.Sleep(10 * time.Millisecond) // Sleep less than timeout
			return 42, nil
		}).
		Go(context.Background())
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if result != 42 {
		t.Errorf("Expected result to be 42, got %d", result)
	}
}