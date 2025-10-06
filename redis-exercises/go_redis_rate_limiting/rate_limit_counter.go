package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Rate limit parameters
	userID := "user:1234"
	limit := 5       // Max number of requests
	windowTime := 60 // Time window in seconds (1 minute)
	key := userID + ":requests"

	// Get current request count
	current, err := rdb.Get(ctx, key).Result()

	if err != nil && err != redis.Nil {
		panic(err)
	}

	// Convert current count from string to int
	currentCount := 0
	if current != "" {
		currentCount, _ = strconv.Atoi(current)
	}

	// Check if limit exceeded
	if currentCount >= limit {
		fmt.Println("Rate limit exceeded. Try again later.")
		return
	}

	// Increase request counter
	newCount, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		panic(err)
	}

	// Set TTL only for the first request in the time window
	if currentCount == 0 {
		rdb.Expire(ctx, key, time.Duration(windowTime)*time.Second)
	}

	fmt.Printf("Request processed. Total requests in current window: %d\n", newCount)
}
