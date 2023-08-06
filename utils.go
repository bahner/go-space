package main

import (
	"github.com/gorilla/websocket"
)

// SendText sends a text message through the WebSocket connection
func SendText(conn *websocket.Conn, text string) error {
	return conn.WriteMessage(websocket.TextMessage, []byte(text))
}
