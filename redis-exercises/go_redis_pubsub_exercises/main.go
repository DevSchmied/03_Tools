package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Task
// Implement a simple chat-like application in Go using Redis Pub/Sub, where:
// - A user can subscribe to multiple channels.
// - Messages published to these channels are automatically delivered to all subscribers.
// - Publishing runs in a separate goroutine.
// - The subscriber receives messages and exits after a given timeout.

func main() {
	ctx := context.Background()

	// Connect
	rdb := connectRedis()
	defer rdb.Close()

	// Subscribe
	pubsub := chanSubscribe(ctx, rdb, "news", "sport")
	defer pubsub.Close()

	// Channel to receive messages
	ch := pubsub.Channel()

	// Start publishing
	startPublishing(ctx, rdb)

	// Receive messages
	receiveMsgs(ch)
}

// Connect to Redis
func connectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// Subscribe to given channels
func chanSubscribe(ctx context.Context, rdb *redis.Client, channels ...string) *redis.PubSub {
	pubsub := rdb.Subscribe(ctx, channels...)

	// Wait for confirmation
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Subscribed to channels:", channels)
	return pubsub
}

// Start publisher in a goroutine
func startPublishing(ctx context.Context, rdb *redis.Client) {
	go func() {
		// Publish 3 messages to "news"
		for i := 1; i <= 3; i++ {
			msg := fmt.Sprintf("This is message %d from channel %s", i, "news")
			err := rdb.Publish(ctx, "news", msg).Err()
			if err != nil {
				fmt.Println("Publish error:", err)
			} else {
				fmt.Println("Published to channel news:", msg)
			}
			time.Sleep(1 * time.Second)
		}

		// Publish one message to "sport"
		msg := fmt.Sprintf("This is a message from channel %s", "sport")
		err := rdb.Publish(ctx, "sport", msg).Err()
		if err != nil {
			fmt.Println("Publish error:", err)
		} else {
			fmt.Println("Published to channel sport:", msg)
		}
	}()
}

// Receive messages with timeout
func receiveMsgs(ch <-chan *redis.Message) {
	timeout := time.After(5 * time.Second)
	for {
		select {
		case msg := <-ch:
			if msg == nil {
				fmt.Println("Channel closed")
				continue
			}
			fmt.Printf("Received from %s: %s\n", msg.Channel, msg.Payload)
		case <-timeout:
			fmt.Println("Timeout reached, exiting")
			return
		}
	}
}
