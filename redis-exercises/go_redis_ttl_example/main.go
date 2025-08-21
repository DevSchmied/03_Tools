package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	fmt.Println("Test")
	ManageTestKey()
}

func ManageTestKey() {

	fmt.Println("-------------------Exercise 1-------------------")
	/*
	   Goal: Create a Golang script that:
	   - Connects to Redis
	   - Stores the value "testValue" with the key "testKey" and a TTL of 60 seconds
	   - Reads the value
	   - Deletes the key
	   - Checks whether the key exists after deletion
	*/

	ctx := context.Background()

	// - Connects to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	fmt.Println("Redis client created:", rdb)

	// - Stores the value "testValue" with the key "testKey" and a TTL of 60 seconds
	err := rdb.SetEX(ctx, "testKey:1000", "testValue", 60*time.Second).Err()
	if err != nil {
		fmt.Printf("Error in SetEX: %v\n", err)
		return
	}
	fmt.Println("Value stored in Redis successfully")

	// - Reads the value
	val, err := rdb.Get(ctx, "testKey:1000").Result()
	if err != nil {
		fmt.Printf("Error in Get: %v\n", err)
		return
	}
	fmt.Println("Value for testKey:1000:", val)

	// - Deletes the key
	err = rdb.Del(ctx, "testKey:1000").Err()
	if err != nil {
		fmt.Printf("Error in Del: %v\n", err)
		return
	}
	fmt.Println("Key deleted successfully")

	// - Checks whether the key exists after deletion
	result, err := rdb.Exists(ctx, "testKey:1000").Result()
	if err != nil {
		fmt.Printf("Error in Exists: %v\n", err)
		return
	}
	fmt.Println("Exists after deletion:", result)

}
