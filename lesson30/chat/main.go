package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your nickname: ")
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)

	fmt.Print("Enter chat room name: ")
	chatRoom, _ := reader.ReadString('\n')
	chatRoom = strings.TrimSpace(chatRoom)

	channel := fmt.Sprintf("chatroom:%s", chatRoom)
	privateChannel := fmt.Sprintf("user:%s", strings.ToLower(nickname)) // for receiving private messages
	activeUsersKey := fmt.Sprintf("active_users:%s", chatRoom)
	historyKey := fmt.Sprintf("chat_history:%s", chatRoom)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	pubsub := rdb.Subscribe(ctx, channel, privateChannel)
	defer pubsub.Close()

	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Joined chat room: %s\n", chatRoom)
	rdb.Publish(ctx, channel, fmt.Sprintf("%s has joined the chat room", nickname))
	rdb.SAdd(ctx, activeUsersKey, nickname)
	defer func() {
		rdb.Publish(ctx, channel, fmt.Sprintf("%s has left the chat room", nickname))
		rdb.SRem(ctx, activeUsersKey, nickname)
	}()

	fmt.Println("Chat history:")
	history, _ := rdb.LRange(ctx, historyKey, 0, -1).Result()
	for _, msg := range history {
		fmt.Println(msg)
	}
	fmt.Println("*******You can use `/help` command to use special commands!*******")
	fmt.Println("Type your messages below. Press Ctrl+C to exit.")

	ch := pubsub.Channel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		for msg := range ch {
			if msg.Payload == "" {
				continue
			}
			if msg.Channel == privateChannel {
				fmt.Printf("\r[Private] %s\n> ", msg.Payload)
			} else if !strings.HasPrefix(msg.Payload, nickname+":") {
				fmt.Printf("\r%s\n> ", msg.Payload)
			}
		}
	}()

	for {
		select {
		case <-sigCh:
			fmt.Println("\nExiting chat...")
			return
		default:
			fmt.Print("> ")
			text, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			text = strings.TrimSpace(text)
			if text == "" {
				continue
			}

			if text == "/help" {
				fmt.Println("Commands:")
				fmt.Println("/help       - Show this help message")
				fmt.Println("/users      - List active users")
				fmt.Println("/msg <user> <message> - Send a private message")
				continue
			}

			if text == "/users" {
				users, _ := rdb.SMembers(ctx, activeUsersKey).Result()
				fmt.Println("Active users:", strings.Join(users, ", "))
				continue
			}

			if strings.HasPrefix(text, "/msg ") {
				parts := strings.SplitN(text, " ", 3)
				if len(parts) < 3 {
					fmt.Println("Invalid private message format. Use /msg <nickname> <message>")
					continue
				}
				targetNickname := parts[1]
				privateMessage := fmt.Sprintf("%s: %s", nickname, parts[2])
				err = rdb.Publish(ctx, "user:"+strings.ToLower(targetNickname), privateMessage).Err()
				if err != nil {
					fmt.Println("Error publishing private message:", err)
				}
				fmt.Println("Your private message sent to", targetNickname)
				continue
			}

			message := fmt.Sprintf("%s: %s", nickname, text)
			err = rdb.Publish(ctx, channel, message).Err()
			if err != nil {
				fmt.Println("Error publishing message:", err)
			}

			rdb.RPush(ctx, historyKey, message)
			rdb.LTrim(ctx, historyKey, -100, -1)
		}
	}
}
