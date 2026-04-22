package async

import (
	"context"
	"fmt"
)

// Example code for different types demonstrating Bind

func ExampleBind_int() {
	runner := NewAsyncRunner()
	var result int

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (int, error) {
			return 42, nil
		})).
		Go(context.Background())

	fmt.Println(result)
	// Output: 42
}

func ExampleBind_string() {
	runner := NewAsyncRunner()
	var result string

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (string, error) {
			return "Hello Bind", nil
		})).
		Go(context.Background())

	fmt.Println(result)
	// Output: Hello Bind
}

func ExampleBind_sliceOfInt() {
	runner := NewAsyncRunner()
	var result []int

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) ([]int, error) {
			return []int{1, 2, 3}, nil
		})).
		Go(context.Background())

	fmt.Println(result)
	// Output: [1 2 3]
}

func ExampleBind_sliceOfString() {
	runner := NewAsyncRunner()
	var result []string

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) ([]string, error) {
			return []string{"apple", "banana"}, nil
		})).
		Go(context.Background())

	fmt.Println(result)
	// Output: [apple banana]
}

func ExampleBind_struct() {
	runner := NewAsyncRunner()

	type User struct {
		ID   int
		Name string
	}
	var result User

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (User, error) {
			return User{ID: 1, Name: "Alice"}, nil
		})).
		Go(context.Background())

	fmt.Printf("%+v\n", result)
	// Output: {ID:1 Name:Alice}
}

func ExampleBind_arrayOfStruct() {
	runner := NewAsyncRunner()

	type User struct {
		ID   int
		Name string
	}
	var result []User

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) ([]User, error) {
			return []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}, nil
		})).
		Go(context.Background())

	fmt.Printf("%+v\n", result)
	// Output: [{ID:1 Name:Alice} {ID:2 Name:Bob}]
}

// Data represents a generic wrapper
type Data[T any] struct {
	Value T
}

func ExampleBind_genericType() {
	runner := NewAsyncRunner()
	var result Data[string]

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (Data[string], error) {
			return Data[string]{Value: "Generic Binding"}, nil
		})).
		Go(context.Background())

	fmt.Printf("%+v\n", result)
	// Output: {Value:Generic Binding}
}

func ExampleBind_pointerStruct() {
	runner := NewAsyncRunner()

	type User struct {
		ID   int
		Name string
	}
	var result *User

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (*User, error) {
			return &User{ID: 1, Name: "Alice Pointer"}, nil
		})).
		Go(context.Background())

	if result != nil {
		fmt.Printf("&{ID:%d Name:%s}\n", result.ID, result.Name)
	}
	// Output: &{ID:1 Name:Alice Pointer}
}

func ExampleBind_pointerInt() {
	runner := NewAsyncRunner()
	var result *int

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (*int, error) {
			val := 42
			return &val, nil
		})).
		Go(context.Background())

	if result != nil {
		fmt.Println(*result)
	}
	// Output: 42
}

func ExampleBind_pointerString() {
	runner := NewAsyncRunner()
	var result *string

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (*string, error) {
			val := "Hello Pointer String"
			return &val, nil
		})).
		Go(context.Background())

	if result != nil {
		fmt.Println(*result)
	}
	// Output: Hello Pointer String
}

func ExampleBind_pointerSliceOfInt() {
	runner := NewAsyncRunner()
	var result *[]int

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (*[]int, error) {
			val := []int{10, 20, 30}
			return &val, nil
		})).
		Go(context.Background())

	if result != nil {
		fmt.Println(*result)
	}
	// Output: [10 20 30]
}

func ExampleBind_pointerSliceOfString() {
	runner := NewAsyncRunner()
	var result *[]string

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (*[]string, error) {
			val := []string{"carrot", "potato"}
			return &val, nil
		})).
		Go(context.Background())

	if result != nil {
		fmt.Println(*result)
	}
	// Output: [carrot potato]
}

func ExampleBind_pointerArrayOfStruct() {
	runner := NewAsyncRunner()

	type User struct {
		ID   int
		Name string
	}
	var result *[]User

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (*[]User, error) {
			val := []User{{ID: 3, Name: "Charlie"}, {ID: 4, Name: "Dave"}}
			return &val, nil
		})).
		Go(context.Background())

	if result != nil {
		fmt.Printf("%+v\n", *result)
	}
	// Output: [{ID:3 Name:Charlie} {ID:4 Name:Dave}]
}

func ExampleBind_pointerGenericType() {
	runner := NewAsyncRunner()
	var result *Data[string]

	runner.RunInAsync().
		Task(Bind(&result, func(ctx context.Context) (*Data[string], error) {
			val := Data[string]{Value: "Generic Pointer"}
			return &val, nil
		})).
		Go(context.Background())

	if result != nil {
		fmt.Printf("%+v\n", *result)
	}
	// Output: {Value:Generic Pointer}
}

