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
	UpdateWeatherKeyWithTTL()
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

func UpdateWeatherKeyWithTTL() {
	fmt.Println("\n\n-------------------Exercise 2-------------------")
	/*
	  Requirements:
	   - Create a key weather:moscow with the value +25°C and a TTL of 1 hour
	   - Check the remaining time to live of the key
	   - Update the value to +28°C without changing the TTL
	   - Delete the key before it expires
	*/

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//  - Create a key weather:moscow with the value +25°C and a TTL of 1 hour
	const key = "weather:moscow"

	err := rdb.SetEX(ctx, key, "+25°C", 3600*time.Second).Err()
	if err != nil {
		fmt.Printf("Error in SetEX: %v\n", err)
		return
	}
	fmt.Println("Value stored in Redis successfully")

	// - Check the remaining time to live of the key
	time.Sleep(2 * time.Second)
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		fmt.Printf("Error in TTL: %v\n", err)
		return
	}
	fmt.Println("The remaining time to live of the key is", ttl)

	// - Update the value to +28°C without changing the TTL
	time.Sleep(1 * time.Second)
	err = rdb.Set(ctx, key, "+28°C", redis.KeepTTL).Err()
	if err != nil {
		fmt.Printf("Error in Set: %v\n", err)
		return
	}
	ttl, err = rdb.TTL(ctx, key).Result()
	if err != nil {
		fmt.Printf("Error in TTL: %v\n", err)
		return
	}
	fmt.Println("The remaining time to live of the key is", ttl)

	// - Delete the key before it expires
	err = rdb.Del(ctx, key).Err()
	if err != nil {
		fmt.Printf("Error in Del: %v\n", err)
		return
	}
	exists, _ := rdb.Exists(ctx, key).Result()
	fmt.Printf("Key exists after deletion? %v\n", exists > 0)
	fmt.Println("The key deleted successfully.")
}
