package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Message represents a structured message format
type Message struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
	Source    string      `json:"source"`
}

// WeatherUpdate represents weather-specific data
type WeatherUpdate struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Location    string  `json:"location"`
}

type StockUpdate struct {
	Company string  `json:"company"`
	Price   float64 `json:"price"`
	Volume  int     `json:"volume"`
}

// Publisher represents a component that publishes messages
type Publisher struct {
	client  *redis.Client
	source  string
	channel string
}

// NewPublisher creates a new publisher instance
func NewPublisher(client *redis.Client, source, channel string) *Publisher {
	return &Publisher{
		client:  client,
		source:  source,
		channel: channel,
	}
}

// Publish sends a message to the specified channel
func (p *Publisher) Publish(ctx context.Context, msgType string, data interface{}) error {
	message := Message{
		Type:      msgType,
		Timestamp: time.Now(),
		Data:      data,
		Source:    p.source,
	}

	// Convert message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	// Publish to Redis channel
	err = p.client.Publish(ctx, p.channel, jsonData).Err()
	if err != nil {
		return fmt.Errorf("error publishing message: %v", err)
	}

	return nil
}

// Subscriber represents a component that receives messages
type Subscriber struct {
	client   *redis.Client
	channels []string
	handlers map[string]func(Message)
}

// NewSubscriber creates a new subscriber instance
func NewSubscriber(client *redis.Client, channels []string) *Subscriber {
	return &Subscriber{
		client:   client,
		channels: channels,
		handlers: make(map[string]func(Message)),
	}
}

// AddHandler adds a message handler for a specific message type
func (s *Subscriber) AddHandler(msgType string, handler func(Message)) {
	s.handlers[msgType] = handler
}

// Start begins listening for messages
func (s *Subscriber) Start(ctx context.Context) error {
	// Subscribe to channels
	pubsub := s.client.PSubscribe(ctx, s.channels...)
	defer pubsub.Close()

	// Listen for messages
	channel := pubsub.Channel()
	for msg := range channel {
		// Parse message
		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Handle message based on type
		if handler, ok := s.handlers[message.Type]; ok {
			handler(message)
		}
	}

	return nil
}

func main() {
	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()

	// Create publisher
	weatherPublisher := NewPublisher(client, "weather-station-1", "weather.updates")

	stockPublisher := NewPublisher(client, "stock-exchange-1", "stock.updates")

	// Create subscriber
	subscriber := NewSubscriber(client, []string{"weather.updates", "stock.updates"})

	// Add message handlers
	subscriber.AddHandler("temperature", func(msg Message) {
		fmt.Println(msg.Data)
		if update, ok := msg.Data.(WeatherUpdate); ok {
			fmt.Printf("Temperature Update: %.1f°C in %s\n",
				update.Temperature,
				update.Location)
		}
	})

	subscriber.AddHandler("stock", func(msg Message) {
		fmt.Println(msg.Data)
		if update, ok := msg.Data.(StockUpdate); ok {
			fmt.Printf("Stock Update: %s, Price: %.2f, Volume: %d\n",
				update.Company, update.Price, update.Volume)
		}
	})

	// Start subscriber in a goroutine
	go func() {
		if err := subscriber.Start(ctx); err != nil {
			log.Printf("Subscriber error: %v", err)
		}
	}()

	// Simulate publishing weather updates
	go func() {
		for i := 0; i < 5; i++ {
			update := WeatherUpdate{
				Temperature: 20.5 + float64(i),
				Humidity:    65.0,
				Location:    "New York",
			}

			if err := weatherPublisher.Publish(ctx, "temperature", update); err != nil {
				log.Printf("Error publishing update: %v", err)
			}

			time.Sleep(time.Second)
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			stockUpdate := StockUpdate{
				Company: "Apple",
				Price:   150.5 + float64(i),
				Volume:  1000 + (i * 100),
			}

			if err := stockPublisher.Publish(ctx, "stock", stockUpdate); err != nil {
				log.Printf("Error publishing stock update: %v", err)
			}

			time.Sleep(time.Second)
		}
	}()

}
