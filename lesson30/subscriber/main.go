package main

import (
	"context"
	"fmt"

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
	channel := "mychan"

	// Subscribe to the channel
	pubsub := rdb.Subscribe(ctx, channel)

	// Wait for confirmation that subscription is created
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Subscribed to channel:", channel)

	// Create a channel to receive messages
	ch := pubsub.Channel()

	// Consume messages in a goroutine
	for msg := range ch {
		fmt.Printf("Received message from %s: %s\n", msg.Channel, msg.Payload)
	}
}
