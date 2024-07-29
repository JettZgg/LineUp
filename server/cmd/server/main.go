package main

import (
	"log"

	"github.com/JettZgg/LineUp/internal/config"
	"github.com/JettZgg/LineUp/internal/server"
	"github.com/JettZgg/LineUp/internal/utils/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Create and start server
	srv, err := server.New(cfg, hub)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Fatal(srv.Start())
}
