package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Using Redis Set for unique tags
	// Requirements:
	// - Add 5 tags to the set article:123:tags
	// - Check if the tag 'golang' exists
	// - Remove 2 tags and get a random tag

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	articleTagsKey := "article:123:tags"

	// - Add 5 tags to the set
	err := rdb.SAdd(ctx, articleTagsKey, "tag1", "tag2", "tag3", "tag4", "tag5").Err()
	if err != nil {
		fmt.Println("Error in SAdd:", err)
		return
	}
	fmt.Println("5 tags successfully added to the set.")

	// - Check if 'golang' tag exists
	isGolangPresent, err := rdb.SIsMember(ctx, articleTagsKey, "golang").Result()
	if err != nil {
		fmt.Println("Error in SIsMember:", err)
		return
	}
	fmt.Printf("Is 'golang' tag present? %t\n", isGolangPresent)

	// - Remove 2 tags from the set
	err = rdb.SRem(ctx, articleTagsKey, "tag3", "tag1").Err()
	if err != nil {
		fmt.Println("Error in SRem:", err)
		return
	}
	fmt.Println("Tags 'tag3' and 'tag1' removed from the set.")

	// - Get a random tag from the set
	randomTag, err := rdb.SRandMember(ctx, articleTagsKey).Result()
	if err != nil {
		fmt.Println("Error in SRandMember:", err)
		return
	}
	fmt.Println("Random tag from the set:", randomTag)

	currentTags, err := rdb.SMembers(ctx, articleTagsKey).Result()
	if err != nil {
		fmt.Println("Error fetching current tags:", err)
		return
	}
	fmt.Println("Current tags in the set:", currentTags)
}
