package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var (
	dialer = websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func main() {
	exitSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	conn, _, err := dialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		slog.Error("Failed to dial", slog.Any("error", err))
		return
	}

	// Graceful shutdown
	go func() {
		<-exitSignal.Done()
		err = conn.Close()
		if err != nil {
			slog.Error("Failed to close connection", slog.Any("error", err))
		}
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Message: ")
		scanner.Scan()

		err = conn.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
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
	}
}
