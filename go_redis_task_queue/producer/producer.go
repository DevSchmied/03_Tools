package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// List of tasks to add
	tasks := []string{"task1", "task2", "task3"}

	// Add tasks to the queue (FIFO: add to the end)
	for _, task := range tasks {
		err := rdb.RPush(ctx, "task_queue", task).Err()
		if err != nil {
			panic(err)
		}
		fmt.Println("Added:", task)
	}
}
