package main

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

var (
	dialer = websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func main() {
	conn, _, err := dialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		slog.Error("Failed to dial", slog.Any("error", err))
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello World!!!!!!!!"))
	if err != nil {
		slog.Error("Failed to send message", slog.Any("error", err))
		return
	}

	msgType, msgByte, err := conn.ReadMessage()
	if err != nil {
		slog.Error("Failed to read message", slog.Any("error", err))
		return
	}

	slog.Info("Message received!", slog.String("message", string(msgByte)), slog.Int("msgtype", msgType))
	err = conn.Close()
	if err != nil {
		slog.Error("Failed to close connection", slog.Any("error", err))
	}
}
