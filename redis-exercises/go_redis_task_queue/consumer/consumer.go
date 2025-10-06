package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Counter for loop iterations
	counter := 0
	maxAttempts := 10

	for {
		// Stop after 10 attempts
		if counter >= maxAttempts {
			fmt.Println("Reached maximum number of attempts, stopping consumer.")
			break
		}

		// Pop one task from the beginning of the queue
		task, err := rdb.LPop(ctx, "task_queue").Result()
		if err == redis.Nil {
			// Queue is empty
			fmt.Println("Queue is empty, waiting for tasks...")
			time.Sleep(1 * time.Second)
			counter++
			continue
		} else if err != nil {
			panic(err)
		}

		// Simulate task processing
		fmt.Println("Processing:", task)
		time.Sleep(2 * time.Second)

		counter++
	}

	fmt.Println("Consumer finished.")
}
