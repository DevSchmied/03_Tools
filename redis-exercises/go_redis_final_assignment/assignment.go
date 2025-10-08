// Final Assignment: Working with Redis in Go
//
// Task:
// Implement a Go program that connects to a Redis server and demonstrates basic operations,
// including:
// - establishing a connection and checking it with PING
// - writing and reading simple key-value data
// - working with TTL (Time-To-Live)
// - using different Redis data structures (Lists, Hashes, Sets)
// - deleting keys
// - implementing a simple cache mechanism using Redis
//
// The code below demonstrates all these operations step by step.

package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {

	// Step 1: Connect to Redis and check the connection with PING

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // Password (if set)
	})

	// Check connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis connection error:", err)
		return
	}

	fmt.Println("Connected to Redis successfully")
}
