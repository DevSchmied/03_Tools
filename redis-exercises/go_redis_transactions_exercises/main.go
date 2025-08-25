package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {

	// Go exercise: Redis transactions
	// Requirements:
	// 1. Connect to Redis
	// 2. Use MULTI ... EXEC to set key1="value1" and key2="value2"
	// 3. Simulate DISCARD by canceling a transaction under a condition
	// 4. Use WATCH on key1: if value is unchanged, update it to "new_value"
	// 5. Simulate "buying a book":
	//    - Decrease user:1000:balance by 50
	//    - Add "Spent 50 on a new book" to user:1000:transactions list

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// MULTI ... EXEC
	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SetEX(ctx, "key1", "value1", time.Hour)
		pipe.SetEX(ctx, "key2", "value2", time.Hour)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// DISCARD example
	pipe := rdb.TxPipeline()
	pipe.SetEX(ctx, "discard:key1", "temp1", time.Minute)
	pipe.SetEX(ctx, "discard:key2", "temp2", time.Minute)

	// Simulate a condition that cancels the transaction
	discardCondition := true
	if discardCondition {
		err = pipe.Discard()
		if err != nil {
			fmt.Println("Discard error:", err)
		} else {
			fmt.Println("Transaction was discarded, keys were not set")
		}
	} else {
		_, err = pipe.Exec(ctx)
		if err != nil {
			fmt.Println("Exec error:", err)
		}
	}

	// WATCH key1
	err = rdb.Watch(ctx, func(tx *redis.Tx) error {
		v, err := tx.Get(ctx, "key1").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if v == "value1" {
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, "key1", "new_value", time.Hour)
				return nil
			})
		}
		return err
	}, "key1")
	if err != nil {
		fmt.Println("WATCH error:", err)
	}

	// Buying a book
	_, err = rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.DecrBy(ctx, "user:1000:balance", 50)
		pipe.LPush(ctx, "user:1000:transactions", "Spent 50 on a new book")
		return nil
	})
	if err != nil {
		fmt.Println("Buy error:", err)
	}
}
