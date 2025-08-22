package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func printHeader(n int) {
	fmt.Printf("\n------------------------------%d. exercise------------------------------\n", n)
}

func main() {
	HashExercise()
	UserVisitsHash()
}
func HashExercise() {

	printHeader(1)
	// Requirements:
	// Create a hash "user:1000" with fields:
	// "name" → "John Doe"
	// "email" → "john@example.com"
	// "login_count" → 0
	// Increment login_count by 1
	// Get all fields of the hash

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err := rdb.HSet(ctx, "user:1000",
		"name", "John Doe",
		"email", "john@example.com",
		"login_count", 0,
	).Err()
	if err != nil {
		fmt.Println("Error in HSet:", err)
		return
	}
	fmt.Println("Hash created successfully")

	// Increment login_count by 1
	err = rdb.HIncrBy(ctx, "user:1000", "login_count", 1).Err()
	if err != nil {
		fmt.Println("Error in HIncrBy:", err)
		return
	}
	fmt.Println("login_count incremented")

	// Get all fields of the hash
	values, err := rdb.HGetAll(ctx, "user:1000").Result()
	if err != nil {
		fmt.Println("Error in HGetAll:", err)
		return
	}
	fmt.Println("Hash values:", values)
}

func UserVisitsHash() {
	printHeader(2)
	// Requirements:
	// - Create a hash user:100 with fields:
	//   - name (string)
	//   - visits (integer, initialized to 0)
	// - Increment visits by 1 atomically
	// - Get all fields in a single request
	// - Delete the field name

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// Create a hash user:100 with fields name and visits
	err := rdb.HSet(ctx, "user:100",
		"name", "Mustermann",
		"visits", 0).Err()
	if err != nil {
		fmt.Println("Error creating hash:", err)
		return
	}
	fmt.Println("Hash user:100 created successfully with fields name and visits=0")

	// Increment visits by 1 atomically
	err = rdb.HIncrBy(ctx, "user:100", "visits", 1).Err()
	if err != nil {
		fmt.Println("Error incrementing visits:", err)
		return
	}
	fmt.Println("Field visits incremented successfully")

	// Get all fields in a single request
	values, err := rdb.HGetAll(ctx, "user:100").Result()
	if err != nil {
		fmt.Println("Error fetching fields:", err)
		return
	}
	fmt.Printf("All fields after increment: %+v\n", values)

	// Delete the field name
	err = rdb.HDel(ctx, "user:100", "name").Err()
	if err != nil {
		fmt.Println("Error deleting field name:", err)
		return
	}
	fmt.Println("Field name deleted successfully")

	// Get all fields again after deletion
	values, err = rdb.HGetAll(ctx, "user:100").Result()
	if err != nil {
		fmt.Println("Error fetching fields after deletion:", err)
		return
	}
	fmt.Printf("All fields after deleting name: %+v\n", values)
}
