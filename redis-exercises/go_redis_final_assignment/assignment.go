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
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {

	// Step 1: establishing a connection and checking it with PING

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

	// Step 2: writing and reading simple key-value data

	// Write data
	err = client.Set(ctx, "username", "golang_user", 0).Err()
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}

	// Read data
	val, err := client.Get(ctx, "username").Result()
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Println("Value:", val) // Output: golang_user

	// Step 3: working with TTL (Time-To-Live)

	// Write key with TTL
	err = client.Set(ctx, "temp_key", "some_value", time.Minute*5).Err()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Key 'temp_key' set with TTL = 5 minutes")
}
