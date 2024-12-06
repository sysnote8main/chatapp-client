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

type Msg struct {
	Username string
	Message  string
}

var (
	dialer = websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func main() {
	exitSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Username: ")
	scanner.Scan()
	username := scanner.Text()

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

	go func() {
		for {
			msgType, msgByte, err := conn.ReadMessage()
			if err != nil {
				select {
				case <-exitSignal.Done():
					return
				default:
					slog.Error("Failed to read message", slog.Any("error", err))
					return
				}
			}
			slog.Info("Message received!", slog.String("message", string(msgByte)), slog.Int("msgtype", msgType))
		}
	}()

	for {
		scanner.Scan()

		err = conn.WriteJSON(Msg{Username: username, Message: scanner.Text()})
		if err != nil {
			slog.Error("Failed to send message", slog.Any("error", err))
			return
		}
	}
}
