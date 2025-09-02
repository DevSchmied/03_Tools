// Manual Redis sharding example: Distribute keys across two Redis instances using a hash function

package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/go-redis/redis/v8"
)

/*
Task: Manual Redis Sharding

1. Connect to two Redis instances (shards).
2. Implement a simple hash function to assign keys to shards.
3. Use the hash function to select the shard for each key.
4. Store key-value pairs with expiration on the chosen shard.
5. Verify that keys are set correctly by printing confirmation.
*/

// Step 1: Connect to a Redis instance
func connectRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

// Step 2: Implement a simple hash function to determine shard
func hashKey(key string) uint32 {
	h := fnv.New32()
	h.Write([]byte(key))
	return h.Sum32()
}

// Step 3: Choose the correct Redis instance based on hash
func getRedisInstance(key string, rdb1 *redis.Client, rdb2 *redis.Client) *redis.Client {
	hashVal := hashKey(key)
	if hashVal%2 == 0 {
		return rdb1
	}
	return rdb2
}

func main() {
	ctx := context.Background()

	// Step 1: Connect to two Redis instances (shards)
	rdb1 := connectRedis("192.168.1.101:6379")
	rdb2 := connectRedis("192.168.1.102:6379")

	// Keys to store
	key1 := "key1"
	key2 := "key2"

	// Step 3: Select shard for each key
	rdbInstance1 := getRedisInstance(key1, rdb1, rdb2)
	rdbInstance2 := getRedisInstance(key2, rdb1, rdb2)

	// Step 4: Set keys with expiration on the selected shard
	err := rdbInstance1.SetEX(ctx, key1, "value1", 60*time.Second).Err()
	if err != nil {
		fmt.Println("Error setting key1:", err)
	} else {
		fmt.Println("Key1 set successfully on shard")
	}

	err = rdbInstance2.SetEX(ctx, key2, "value2", 60*time.Second).Err()
	if err != nil {
		fmt.Println("Error setting key2:", err)
	} else {
		fmt.Println("Key2 set successfully on shard")
	}

	// Step 5: Confirmation printed for each key
}
