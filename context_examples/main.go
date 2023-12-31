package main

import (
	"context"
	"fmt"
	"time"
)

// TimedTask simulates a task that takes time to complete. It respects context timeout.
func TimedTask(ctx context.Context) {
	// Simulating a long-running task
	operation := func() {
		for i := 0; i < 10; i++ {
			if ctx.Err() != nil {
				fmt.Printf("Task interrupted, err: %v.\n", ctx.Err())
				return
			}
			fmt.Printf("Working (%d/10).\n", i+1)
			time.Sleep(1 * time.Second)
		}
		fmt.Println("Task completed successfully")
	}

	select {
	case <-ctx.Done():
		fmt.Printf("Operation cancelled, err: %v.\n", ctx.Err())
	case <-func() <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			defer close(ch)
			operation()
		}()
		return ch
	}():
	}
}

func main() {
	// Create a context that cancels automatically after a 5-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run the timed task with the context.
	TimedTask(ctx)
}
