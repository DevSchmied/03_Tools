package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

/*
Order processing implementation using Redis Streams and consumer group.

Steps:
1) Add several orders to the "orders" stream.
2) Create a consumer group "order_group".
3) Start two workers (goroutines): consumer1 and consumer2.
4) Each worker:
   - performs XREADGROUP
   - starts WATCH on the key "processed:<message-id>"
   - if the key does not exist — inside a transaction (MULTI/EXEC) marks the message as processed (SET)
     and performs a side effect (LPUSH into "processed_orders")
   - after the transaction performs XACK to acknowledge the message
   - after successful processing sends a signal into the "done" channel
5) After receiving signals for all orders, print the contents of the "processed_orders" list.
*/

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	totalOrders := 5
	done := make(chan struct{}, totalOrders)

	// 1) Add several orders to the "orders" stream
	for i := 1; i <= totalOrders; i++ {
		err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: "orders",
			Values: map[string]interface{}{"item": fmt.Sprintf("Order-%d", i)},
		}).Err()
		if err != nil {
			log.Fatalf("XADD error: %v", err)
		}
	}

	// 2) Create consumer group (ignore error if it already exists)
	err := rdb.XGroupCreateMkStream(ctx, "orders", "order_group", "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Error creating consumer group: %v", err)
	}

	var wg sync.WaitGroup
	// 3) Start two workers
	wg.Add(1)
	go worker(rdb, "consumer1", done, &wg)
	wg.Add(1)
	go worker(rdb, "consumer2", done, &wg)

	// Wait until all orders are processed (N signals)
	for i := 0; i < totalOrders; i++ {
		<-done
	}

	// 5) Show the contents of the "processed_orders" list
	processed, err := rdb.LRange(ctx, "processed_orders", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error reading processed_orders: %v", err)
	}
	fmt.Println("\nProcessed orders:", processed)

	wg.Wait()
}

func worker(rdb *redis.Client, consumerName string, done chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// Read messages for a specific consumer
		messages, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "order_group",
			Consumer: consumerName,
			Streams:  []string{"orders", ">"},
			Block:    3 * time.Second,
			Count:    1,
		}).Result()

		if err == redis.Nil || len(messages) == 0 {
			// No new messages: keep waiting
			continue
		}
		if err != nil {
			log.Printf("[%s] XREADGROUP error: %v\n", consumerName, err)
			continue
		}

		for _, stream := range messages {
			for _, msg := range stream.Messages {
				orderID := msg.ID
				processedKey := fmt.Sprintf("processed:%s", orderID)

				// WATCH ensures key check before transaction
				err := rdb.Watch(ctx, func(tx *redis.Tx) error {
					exists, err := tx.Exists(ctx, processedKey).Result()
					if err != nil {
						return err
					}
					if exists == 1 {
						// Already processed: just acknowledge
						_, err = tx.XAck(ctx, "orders", "order_group", orderID).Result()
						if err != nil {
							return err
						}
						log.Printf("[%s] Order already processed: ID=%s\n", consumerName, orderID)
						return nil
					}

					// MULTI/EXEC block
					_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
						pipe.Set(ctx, processedKey, 1, 24*time.Hour)
						pipe.LPush(ctx, "processed_orders", orderID)
						pipe.XAck(ctx, "orders", "order_group", orderID)
						return nil
					})
					if err != nil {
						return err
					}

					log.Printf("[%s] Order processed: ID=%s Values=%v\n", consumerName, orderID, msg.Values)

					// Send signal that one order is done
					done <- struct{}{}
					return nil
				}, processedKey)

				if err != nil {
					log.Printf("[%s] Transaction error: %v\n", consumerName, err)
				}

				// Trim old messages from the stream
				_, err = rdb.XTrimMaxLen(ctx, "orders", 1000).Result()
				if err != nil {
					log.Printf("[%s] XTRIM error: %v\n", consumerName, err)
				}
			}
		}
	}
}
