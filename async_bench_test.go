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
			Task(Bind(&result1, func() (int, error) {
				return 1, nil
			})).
			Task(Bind(&result2, func() (int, error) {
				return 2, nil
			})).
			Task(Bind(&result3, func() (int, error) {
				return 3, nil
			})).
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
		a := runner.RunInAsync()

		// Add 10 tasks
		results := make([]int, 10)
		for j := 0; j < 10; j++ {
			j := j // capture loop variable
			a.Task(Bind(&results[j], func() (int, error) {
				return j, nil
			}))
		}

		err := a.Go(context.Background())
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAsyncRawTask(b *testing.B) {
	runner := NewAsyncRunner()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := runner.RunInAsync().
			Task(func(ctx context.Context) error {
				return nil
			}).
			Task(func(ctx context.Context) error {
				return nil
			}).
			Task(func(ctx context.Context) error {
				return nil
			}).
			Go(context.Background())

		if err != nil {
			b.Fatal(err)
		}
	}
}