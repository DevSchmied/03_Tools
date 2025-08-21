package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Implementing a simple FIFO queue in Redis
	// Requirements:
	// - Add 3 tasks to the list "tasks:queue"
	// - Get the last 2 tasks without removing them
	// - Pop one task for processing
	// - Trim the list to 2 elements

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// - Add 3 tasks to the queue
	err := rdb.RPush(ctx, "tasks:queue", "task1", "task2", "task3").Err()
	if err != nil {
		fmt.Println("Error adding tasks:", err)
		return
	}
	fmt.Println("Tasks successfully added to the queue.")

	// - Get the last 2 tasks without removing
	tasks, err := rdb.LRange(ctx, "tasks:queue", -2, -1).Result()
	if err != nil {
		fmt.Println("Error fetching tasks:", err)
		return
	}
	fmt.Println("Last 2 tasks in the queue:")
	for i, task := range tasks {
		fmt.Printf("%d. %s\n", i+1, task)
	}

	// - Pop one task for processing
	currentTask, err := rdb.LPop(ctx, "tasks:queue").Result()
	if err != nil {
		fmt.Println("Error popping task:", err)
		return
	}
	fmt.Println("Processing task:", currentTask)

	// - Trim the list to keep only the first 2 elements
	err = rdb.LTrim(ctx, "tasks:queue", 0, 1).Err()
	if err != nil {
		fmt.Println("Error trimming the queue:", err)
		return
	}
	fmt.Println("Queue trimmed to 2 tasks.")
}
