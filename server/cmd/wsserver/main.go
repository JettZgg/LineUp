package main

import (
	"log"

	"github.com/JettZgg/LineUp/internal/config"
	"github.com/JettZgg/LineUp/internal/websocket"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	hub := websocket.NewHub()
	go hub.Run()

	wsServer := websocket.NewServer(cfg, hub)
	log.Fatal(wsServer.Start())
}