package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// SessionData
type SessionData struct {
	UserID   int      `json:"user_id"`
	Username string   `json:"username"`
	Cart     []string `json:"cart"`
}

func main() {
	ctx := context.Background()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	sessionID := "session:12345"

	// Example session data
	data := SessionData{
		UserID:   1,
		Username: "john_doe",
		Cart:     []string{"item1", "item2"},
	}

	// Struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// with TTL = 30 minutes
	err = rdb.Set(ctx, sessionID, jsonData, 30*time.Minute).Err()
	if err != nil {
		panic(err)
	}

	// Retrieve session data
	val, err := rdb.Get(ctx, sessionID).Result()
	if err != nil {
		panic(err)
	}

	// JSON back to struct
	var session SessionData
	json.Unmarshal([]byte(val), &session)

	fmt.Println("Session data:", session)
}
