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
		DB:   0,
	})

	sessionHistoryKey := "session:12345:history"

	// Add user actions to session history
	rdb.LPush(ctx, sessionHistoryKey, "viewed_item1")
	rdb.LPush(ctx, sessionHistoryKey, "added_to_cart:item2")

	// Get all user actions (latest first)
	history, err := rdb.LRange(ctx, sessionHistoryKey, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Session history:", history)
}
