package async

import (
	"context"
	"testing"
)

func BenchmarkAsyncExecution(b *testing.B) {
	runner := NewAsyncRunner()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result1, result2, result3 int
		
		err := runner.RunInAsync().
			Task(&result1, func() (any, error) {
				return 1, nil
			}).
			Task(&result2, func() (any, error) {
				return 2, nil
			}).
			Task(&result3, func() (any, error) {
				return 3, nil
			}).
			Go(context.Background())
		
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAsyncExecutionManyTasks(b *testing.B) {
	runner := NewAsyncRunner()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		async := runner.RunInAsync()
		
		// Add 10 tasks
		results := make([]int, 10)
		for j := 0; j < 10; j++ {
			j := j // capture loop variable
			async.Task(&results[j], func() (any, error) {
				return j, nil
			})
		}
		
		err := async.Go(context.Background())
		if err != nil {
			b.Fatal(err)
		}
	}
}