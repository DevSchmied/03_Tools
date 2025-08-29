package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// Task:
// 1. Add message to stream
// 2. Read messages from stream
// 3. Create consumer group
// 4. Read with consumer group
// 5. Acknowledge message

func connectRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func main() {

	ctx := context.Background()

	// Connect to Redis
	rdb := connectRedis("localhost:6379")
	defer rdb.Close()

	// 1. Add message to stream
	id, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "mystream",
		Values: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}).Result()
	if err != nil {
		fmt.Printf("XADD failed: %v\n", err)
		return
	}
	fmt.Println("Added message ID:", id)

	// 2. Read messages from stream
	msgs, err := rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"mystream", "0"},
		Count:   2,
		Block:   0,
	}).Result()
	if err != nil {
		fmt.Printf("XREAD failed: %v\n", err)
		return
	}

	for _, stream := range msgs {
		for _, message := range stream.Messages {
			fmt.Printf("Message ID=%s Values=%v\n", message.ID, message.Values)
		}
	}

	// 3. Create consumer group
	err = rdb.XGroupCreateMkStream(ctx, "mystream", "mygroup", "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		fmt.Printf("XGROUP CREATE failed: %v", err)
		return
	}

	// 4. Read with consumer group
	groupMsgs, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    "mygroup",
		Consumer: "consumer1",
		Streams:  []string{"mystream", ">"},
		Count:    2,
		Block:    0,
	}).Result()
	if err != nil {
		fmt.Printf("XREADGROUP failed: %v\n", err)
		return
	}

	for _, stream := range groupMsgs {
		for _, message := range stream.Messages {
			fmt.Printf("[Group] Message ID=%s Values=%v\n", message.ID, message.Values)

			// 5. Acknowledge message
			if err := rdb.XAck(ctx, "mystream", "mygroup", message.ID).Err(); err != nil {
				fmt.Printf("XACK failed: %v\n", err)
				return
			}
		}
	}
}
