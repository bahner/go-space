package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// handleClient handles the WebSocket client connection
func handleClient(conn *websocket.Conn, topic *Topic) {
	msg := fmt.Sprintf("Welcome to the chat room %q!", topic.PubSubTopic)
	sendText(conn, msg)

	sub, err := topic.PubSubTopic.Subscribe()
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine for reading from the WebSocket
	go func() {
		defer wg.Done()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				// Log the error and return from the goroutine
				log.Printf("read error: %v", err)
				return
			}

			// Publish the message to the pubsub topic
			err = topic.PubSubTopic.Publish(context.Background(), message)
			if err != nil {
				log.Printf("publish error: %v", err)
				return
			}
		}
	}()

	// Goroutine for writing to the WebSocket
	go func() {
		defer wg.Done()
		for {
			msg, err := sub.Next(context.Background())
			if err != nil {
				// Log the error and return from the goroutine
				log.Printf("subscription error: %v", err)
				return
			}

			// Write the message back to the WebSocket
			err = conn.WriteMessage(websocket.TextMessage, msg.GetData())
			if err != nil {
				log.Printf("write error: %v", err)
				return
			}
		}
	}()

	wg.Wait()
}

// sendText sends a text message through the WebSocket connection
func sendText(c *websocket.Conn, text string) error {
	return c.WriteMessage(websocket.TextMessage, []byte(text))
}
