package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Task:
// Write a Go program that connects to a Redis server, executes a simple Lua script to set a key, retrieves the key’s value, loads the script to get its SHA1, executes it again via EVALSHA for a different key, and prints the results. Use a context with timeout for all Redis operations.

// connectRedis creates and returns a Redis client
func connectRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func main() {

	rdb := connectRedis("localhost:6379")

	// Create a context with timeout for all Redis operations
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// time.Sleep(3 * time.Second) // test context timeout

	script := `return redis.call('set', KEYS[1], ARGV[1])`

	// Eval Lua script
	cmdCode, err := rdb.Eval(ctx, script, []string{"mykey"}, "myvalue").Result()
	if err != nil {
		fmt.Println("Error executing Eval:", err)
		return
	}
	fmt.Println("Command code:", cmdCode)

	// Get the key value
	res, err := rdb.Get(ctx, "mykey").Result()
	if err != nil {
		fmt.Println("Error getting 'mykey':", err)
		return
	}
	fmt.Println("Value of 'mykey':", res)

	// Load script to get SHA1
	sha, err := rdb.ScriptLoad(ctx, script).Result()
	if err != nil {
		fmt.Println("Error loading script:", err)
		return
	}
	fmt.Println("SHA1:", sha)

	// EvalSha with timeout context
	err = rdb.EvalSha(ctx, sha, []string{"mykey2"}, "myvalue2").Err()
	if err != nil {
		fmt.Println("Error executing EvalSha:", err)
		return
	}
	fmt.Println("Lua script successfully evaluated on 'mykey2'")

	// Check value of second key
	val2, err := rdb.Get(ctx, "mykey2").Result()
	if err != nil {
		fmt.Println("Error getting 'mykey2':", err)
		return
	}
	fmt.Println("Value of 'mykey2':", val2)
}
