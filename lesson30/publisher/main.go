package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Adjust the address if necessary
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	ctx := context.Background()
	channel := "mychannel"

	// Publish messages to the channel every second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		message := fmt.Sprintf("Current time: %s", t.Format(time.RFC3339))
		err := rdb.Publish(ctx, channel, message).Err()
		if err != nil {
			fmt.Println("Error publishing message:", err)
		} else {
			fmt.Println("Published message:", message)
		}
	}
}
